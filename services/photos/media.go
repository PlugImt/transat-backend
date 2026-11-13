package photos

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// ErrMediaNotFound means Drive has no such file, or no image of it at that size.
var ErrMediaNotFound = errors.New("media not found")

// Media is one file, or one rendition of it, read from Drive. The caller closes Body.
type Media struct {
	Body        io.ReadCloser
	ContentType string
	Size        int64 // -1 when unknown
}

// MediaSource reads the photos of the shared drive for the public album pages. The Drive files are private: the
// API reads them with its own service account and streams them, so no Drive link ever reaches a visitor.
type MediaSource interface {
	// Thumbnail is Drive's own rendition of the file, at most size pixels on its long side.
	Thumbnail(ctx context.Context, id string, size int) (Media, error)
	// Original is the file exactly as it was uploaded.
	Original(ctx context.Context, id string) (Media, error)
}

// thumbLinkTTL keeps Drive's thumbnailLink of a file for a while: it is valid for hours, and fetching it costs
// one Drive call per image otherwise.
const (
	thumbLinkTTL     = 30 * time.Minute
	thumbLinkMaxKeep = 20000
)

// DriveMedia is the MediaSource of the photo sync's read-only service account.
type DriveMedia struct {
	svc  *drive.Service
	http *http.Client

	mu    sync.Mutex
	links map[string]thumbLink
}

type thumbLink struct {
	url     string
	expires time.Time
}

var _ MediaSource = (*DriveMedia)(nil)

// MediaFromEnv builds the media source from GOOGLE_DRIVE_CREDENTIALS_JSON, or returns nil when it is not set.
func MediaFromEnv(ctx context.Context) (*DriveMedia, error) {
	creds := os.Getenv("GOOGLE_DRIVE_CREDENTIALS_JSON")
	if creds == "" {
		return nil, nil
	}
	return NewDriveMedia(ctx, creds)
}

// NewDriveMedia expects the service-account key JSON in base64, like the photo sync. Read-only scope.
func NewDriveMedia(ctx context.Context, credsBase64 string) (*DriveMedia, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(credsBase64))
	if err != nil {
		return nil, fmt.Errorf("decode drive credentials: %w", err)
	}
	creds, err := google.CredentialsFromJSON(ctx, key, drive.DriveReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("drive credentials: %w", err)
	}
	// No client timeout: originals and album archives stream for as long as the visitor reads them. Every call
	// carries its own request context instead.
	client := oauth2.NewClient(context.Background(), creds.TokenSource)
	svc, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("create drive service: %w", err)
	}
	return &DriveMedia{svc: svc, http: client, links: map[string]thumbLink{}}, nil
}

func (d *DriveMedia) Original(ctx context.Context, id string) (Media, error) {
	res, err := d.svc.Files.Get(id).SupportsAllDrives(true).Context(ctx).Download()
	if driveNotFound(err) {
		return Media{}, ErrMediaNotFound
	}
	if err != nil {
		return Media{}, err
	}
	ct := res.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	return Media{Body: res.Body, ContentType: ct, Size: res.ContentLength}, nil
}

func (d *DriveMedia) Thumbnail(ctx context.Context, id string, size int) (Media, error) {
	link, err := d.thumbnailLink(ctx, id)
	if err != nil {
		return Media{}, err
	}
	u, err := url.Parse(fmt.Sprintf("%s=s%d", link, size))
	if err != nil || u.Scheme != "https" || !strings.HasSuffix(u.Hostname(), ".googleusercontent.com") {
		return Media{}, ErrMediaNotFound
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Media{}, err
	}
	res, err := d.http.Do(req)
	if err != nil {
		return Media{}, err
	}
	if res.StatusCode != http.StatusOK {
		_ = res.Body.Close()
		// An expired link answers 403/404: forget it so the next request asks Drive for a fresh one.
		d.forget(id)
		return Media{}, ErrMediaNotFound
	}
	ct := res.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "image/") {
		ct = "image/jpeg"
	}
	return Media{Body: res.Body, ContentType: ct, Size: res.ContentLength}, nil
}

func (d *DriveMedia) thumbnailLink(ctx context.Context, id string) (string, error) {
	now := time.Now()
	d.mu.Lock()
	if l, ok := d.links[id]; ok && now.Before(l.expires) {
		d.mu.Unlock()
		return l.url, nil
	}
	d.mu.Unlock()

	f, err := d.svc.Files.Get(id).SupportsAllDrives(true).Fields("thumbnailLink,trashed").Context(ctx).Do()
	if driveNotFound(err) {
		return "", ErrMediaNotFound
	}
	if err != nil {
		return "", err
	}
	if f.Trashed || f.ThumbnailLink == "" {
		return "", ErrMediaNotFound
	}
	link := f.ThumbnailLink
	if i := strings.LastIndex(link, "="); i > 0 && strings.HasPrefix(link[i:], "=s") {
		link = link[:i]
	}

	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.links) >= thumbLinkMaxKeep {
		for k, l := range d.links {
			if now.After(l.expires) || len(d.links) >= thumbLinkMaxKeep {
				delete(d.links, k)
			}
		}
	}
	d.links[id] = thumbLink{url: link, expires: now.Add(thumbLinkTTL)}
	return link, nil
}

func (d *DriveMedia) forget(id string) {
	d.mu.Lock()
	delete(d.links, id)
	d.mu.Unlock()
}

func driveNotFound(err error) bool {
	var gerr *googleapi.Error
	return errors.As(err, &gerr) && gerr.Code == http.StatusNotFound
}
