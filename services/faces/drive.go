package faces

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// ErrImageUnavailable means the file is gone or Drive has no image for it: the photo is skipped, not retried.
var ErrImageUnavailable = errors.New("no image available for this file")

// ImageSource gives the worker a JPEG of a Drive photo. Photos are only read, never copied anywhere.
type ImageSource interface {
	Image(ctx context.Context, driveFileID string, size int) ([]byte, error)
}

const maxImageBytes = 25 << 20

// DriveImages downloads Drive's own resized rendition of a photo (its thumbnail at the requested size), which
// is far smaller than the original and is always a JPEG or PNG the model can decode.
type DriveImages struct {
	svc  *drive.Service
	http *http.Client
}

var _ ImageSource = (*DriveImages)(nil)

// NewDriveImages expects the service-account key JSON in base64, like the photo sync. Read-only scope.
func NewDriveImages(ctx context.Context, credsBase64 string) (*DriveImages, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(credsBase64))
	if err != nil {
		return nil, fmt.Errorf("decode drive credentials: %w", err)
	}
	creds, err := google.CredentialsFromJSON(ctx, key, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("drive credentials: %w", err)
	}
	client := oauth2.NewClient(ctx, creds.TokenSource)
	client.Timeout = 2 * time.Minute
	svc, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("create drive service: %w", err)
	}
	return &DriveImages{svc: svc, http: client}, nil
}

func (d *DriveImages) Image(ctx context.Context, id string, size int) ([]byte, error) {
	f, err := d.svc.Files.Get(id).SupportsAllDrives(true).Fields("thumbnailLink,trashed").Context(ctx).Do()
	var gerr *googleapi.Error
	if errors.As(err, &gerr) && gerr.Code == http.StatusNotFound {
		return nil, ErrImageUnavailable
	}
	if err != nil {
		return nil, err
	}
	if f.Trashed || f.ThumbnailLink == "" {
		return nil, ErrImageUnavailable
	}
	link := f.ThumbnailLink
	if i := strings.LastIndex(link, "=s"); i > 0 {
		link = link[:i]
	}
	u, err := url.Parse(fmt.Sprintf("%s=s%d", link, size))
	if err != nil || u.Scheme != "https" || !strings.HasSuffix(u.Hostname(), ".googleusercontent.com") {
		return nil, ErrImageUnavailable
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	res, err := d.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	switch {
	case res.StatusCode == http.StatusNotFound:
		return nil, ErrImageUnavailable
	case res.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("drive thumbnail answered %d", res.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(res.Body, maxImageBytes+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxImageBytes {
		return nil, ErrImageUnavailable
	}
	return b, nil
}
