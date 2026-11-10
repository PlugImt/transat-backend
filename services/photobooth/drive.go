package photobooth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const mimeFolder = "application/vnd.google-apps.folder"

// ErrNotFound is returned when a Drive file or folder does not exist (or the service account cannot see it).
var ErrNotFound = errors.New("drive item not found")

// Folder is a Drive folder seen by the resolver.
type Folder struct {
	ID   string
	Name string
}

// FileInfo is what the service checks on an uploaded file and serves on the share page.
type FileInfo struct {
	ID            string
	Name          string
	MimeType      string
	Parents       []string
	Size          int64
	MD5           string
	Trashed       bool
	ThumbnailLink string
	Width, Height int
}

// Drive is every Drive operation the photobooth needs. The real implementation is DriveClient; tests use a fake.
type Drive interface {
	ListChildFolders(ctx context.Context, parentID string) ([]Folder, error)
	ListChildFilesNamed(ctx context.Context, parentID, name string) ([]FileInfo, error)
	CreateFolder(ctx context.Context, parentID, name string) (string, error)
	// StartUpload opens a resumable upload session; whoever holds the returned URL can send exactly one file
	// of that size into parentID, for about a week.
	StartUpload(ctx context.Context, parentID, name, mimeType string, size int64) (string, error)
	GetFile(ctx context.Context, id string) (FileInfo, error)
	Trash(ctx context.Context, id string) error
	Download(ctx context.Context, id string) (io.ReadCloser, error)
	Thumbnail(ctx context.Context, info FileInfo, size int) (io.ReadCloser, string, error)
}

// DriveClient writes into the shared Photos drive with a service account that is "Content manager" on it.
type DriveClient struct {
	svc  *drive.Service
	http *http.Client
}

var _ Drive = (*DriveClient)(nil)

var driveIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{10,200}$`)

// NewDriveClient expects the service-account key JSON encoded in base64, like the photo sync.
func NewDriveClient(ctx context.Context, credsBase64 string) (*DriveClient, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(credsBase64))
	if err != nil {
		return nil, fmt.Errorf("decode drive credentials: %w", err)
	}
	creds, err := google.CredentialsFromJSON(ctx, key, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("drive credentials: %w", err)
	}
	// No client timeout: downloads stream for as long as their request context allows.
	client := oauth2.NewClient(ctx, creds.TokenSource)
	svc, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("create drive service: %w", err)
	}
	return &DriveClient{svc: svc, http: client}, nil
}

func quoteQ(s string) string {
	return strings.NewReplacer(`\`, `\\`, `'`, `\'`).Replace(s)
}

func (c *DriveClient) ListChildFolders(ctx context.Context, parentID string) ([]Folder, error) {
	var out []Folder
	q := fmt.Sprintf("'%s' in parents and mimeType = '%s' and trashed = false", quoteQ(parentID), mimeFolder)
	err := c.list(ctx, q, func(f *drive.File) { out = append(out, Folder{ID: f.Id, Name: f.Name}) })
	return out, err
}

func (c *DriveClient) ListChildFilesNamed(ctx context.Context, parentID, name string) ([]FileInfo, error) {
	var out []FileInfo
	q := fmt.Sprintf("'%s' in parents and name = '%s' and trashed = false", quoteQ(parentID), quoteQ(name))
	err := c.list(ctx, q, func(f *drive.File) { out = append(out, toInfo(f)) })
	return out, err
}

func (c *DriveClient) list(ctx context.Context, q string, each func(*drive.File)) error {
	token := ""
	for {
		call := c.svc.Files.List().Q(q).Corpora("allDrives").IncludeItemsFromAllDrives(true).SupportsAllDrives(true).
			PageSize(1000).Fields("nextPageToken, files(" + fileFields + ")").Context(ctx)
		if token != "" {
			call = call.PageToken(token)
		}
		var page *drive.FileList
		if err := retry(ctx, func() (err error) { page, err = call.Do(); return }); err != nil {
			return err
		}
		for _, f := range page.Files {
			each(f)
		}
		if page.NextPageToken == "" {
			return nil
		}
		token = page.NextPageToken
	}
}

func (c *DriveClient) CreateFolder(ctx context.Context, parentID, name string) (string, error) {
	var f *drive.File
	err := retry(ctx, func() (err error) {
		f, err = c.svc.Files.Create(&drive.File{Name: name, MimeType: mimeFolder, Parents: []string{parentID}}).
			SupportsAllDrives(true).Fields("id").Context(ctx).Do()
		return
	})
	if err != nil {
		return "", err
	}
	return f.Id, nil
}

func (c *DriveClient) StartUpload(ctx context.Context, parentID, name, mimeType string, size int64) (string, error) {
	meta, err := json.Marshal(map[string]any{"name": name, "parents": []string{parentID}, "mimeType": mimeType})
	if err != nil {
		return "", err
	}
	endpoint := "https://www.googleapis.com/upload/drive/v3/files?uploadType=resumable&supportsAllDrives=true" +
		"&fields=" + url.QueryEscape("id,name,size,md5Checksum,parents")
	var location string
	err = retry(ctx, func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(meta))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("X-Upload-Content-Type", mimeType)
		req.Header.Set("X-Upload-Content-Length", fmt.Sprint(size))
		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		if err := googleapi.CheckResponse(res); err != nil {
			return err
		}
		location = res.Header.Get("Location")
		if location == "" {
			return errors.New("drive did not return an upload session")
		}
		return nil
	})
	return location, err
}

const fileFields = "id,name,mimeType,parents,size,md5Checksum,trashed,thumbnailLink,imageMediaMetadata(width,height)"

func (c *DriveClient) GetFile(ctx context.Context, id string) (FileInfo, error) {
	if !driveIDRe.MatchString(id) {
		return FileInfo{}, ErrNotFound
	}
	var f *drive.File
	err := retry(ctx, func() (err error) {
		f, err = c.svc.Files.Get(id).SupportsAllDrives(true).Fields(fileFields).Context(ctx).Do()
		return
	})
	if isNotFound(err) {
		return FileInfo{}, ErrNotFound
	}
	if err != nil {
		return FileInfo{}, err
	}
	return toInfo(f), nil
}

// Trash moves a file to the drive's trash, where Drive purges it after 30 days. A file already gone is fine.
func (c *DriveClient) Trash(ctx context.Context, id string) error {
	if !driveIDRe.MatchString(id) {
		return nil
	}
	err := retry(ctx, func() error {
		_, err := c.svc.Files.Update(id, &drive.File{Trashed: true}).SupportsAllDrives(true).Fields("id").Context(ctx).Do()
		return err
	})
	if isNotFound(err) {
		return nil
	}
	return err
}

func (c *DriveClient) Download(ctx context.Context, id string) (io.ReadCloser, error) {
	res, err := c.svc.Files.Get(id).SupportsAllDrives(true).Context(ctx).Download()
	if isNotFound(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

// Thumbnail returns a JPEG at most size pixels on its long side, rendered by Drive. Drive may need a few
// seconds after an upload before it has one; the caller then falls back to the original.
func (c *DriveClient) Thumbnail(ctx context.Context, info FileInfo, size int) (io.ReadCloser, string, error) {
	link := info.ThumbnailLink
	if link == "" {
		return nil, "", ErrNotFound
	}
	if i := strings.LastIndex(link, "=s"); i > 0 {
		link = link[:i]
	}
	u, err := url.Parse(fmt.Sprintf("%s=s%d", link, size))
	if err != nil || u.Scheme != "https" || !strings.HasSuffix(u.Hostname(), ".googleusercontent.com") {
		return nil, "", ErrNotFound
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, "", err
	}
	res, err := c.http.Do(req)
	if err != nil {
		return nil, "", err
	}
	if res.StatusCode != http.StatusOK {
		_ = res.Body.Close()
		return nil, "", ErrNotFound
	}
	return res.Body, res.Header.Get("Content-Type"), nil
}

func toInfo(f *drive.File) FileInfo {
	info := FileInfo{
		ID: f.Id, Name: f.Name, MimeType: f.MimeType, Parents: f.Parents, Size: f.Size, MD5: f.Md5Checksum,
		Trashed: f.Trashed, ThumbnailLink: f.ThumbnailLink,
	}
	if m := f.ImageMediaMetadata; m != nil {
		info.Width, info.Height = int(m.Width), int(m.Height)
	}
	return info
}

func isNotFound(err error) bool {
	var gerr *googleapi.Error
	return errors.As(err, &gerr) && gerr.Code == http.StatusNotFound
}

// retry backs off on Drive quota errors and transient failures.
func retry(ctx context.Context, fn func() error) error {
	delay := time.Second
	for attempt := 0; ; attempt++ {
		err := fn()
		if err == nil || attempt >= 4 || !retryable(err) {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay *= 2
	}
}

func retryable(err error) bool {
	var gerr *googleapi.Error
	if !errors.As(err, &gerr) {
		return !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded)
	}
	if gerr.Code == 429 || gerr.Code >= 500 {
		return true
	}
	for _, d := range gerr.Errors {
		if d.Reason == "rateLimitExceeded" || d.Reason == "userRateLimitExceeded" {
			return true
		}
	}
	return false
}
