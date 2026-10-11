package photos

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	MimeFolder   = "application/vnd.google-apps.folder"
	MimeShortcut = "application/vnd.google-apps.shortcut"

	// maxInFlight caps concurrent Drive calls; the sync is sequential today, this guards future callers.
	maxInFlight = 4
	pageSize    = 1000
)

// itemFields is the whole metadata the index needs. thumbnailLink is requested only to know a thumbnail
// exists: it expires within hours and is never stored.
const itemFields = "id,name,mimeType,parents,size,md5Checksum,createdTime,modifiedTime,webViewLink,trashed," +
	"hasThumbnail,shortcutDetails(targetId,targetMimeType),imageMediaMetadata(time,width,height),videoMediaMetadata(width,height)"

// DriveItem is the subset of Drive file metadata the index relies on.
type DriveItem struct {
	ID             string
	Name           string
	MimeType       string
	ParentID       string
	Size           int64
	MD5            string
	CreatedTime    time.Time
	ModifiedTime   time.Time
	WebViewLink    string
	Trashed        bool
	HasThumbnail   bool
	ShortcutTarget string
	ShortcutMime   string
	ExifTime       string // camera wall clock, no zone
	Width, Height  int
}

func (i DriveItem) IsFolder() bool   { return i.MimeType == MimeFolder }
func (i DriveItem) IsShortcut() bool { return i.MimeType == MimeShortcut }

// Change is one entry of the Drive changes feed.
type Change struct {
	FileID  string
	Removed bool
	Item    *DriveItem // nil when removed
}

type ChangePage struct {
	Changes           []Change
	NextPageToken     string
	NewStartPageToken string // set on the last page only
}

// DriveAPI is everything the sync needs from Google Drive (read-only).
type DriveAPI interface {
	StartPageToken(ctx context.Context) (string, error)
	ListFolders(ctx context.Context) ([]DriveItem, error)
	ListFilesPage(ctx context.Context, pageToken string) ([]DriveItem, string, error)
	Changes(ctx context.Context, pageToken string) (ChangePage, error)
	GetItem(ctx context.Context, id string) (DriveItem, error)
	ListChildren(ctx context.Context, folderID string) ([]DriveItem, error)
	Calls() map[string]int64
}

// DriveClient reads (never writes) the shared Photos drive with a read-only service account.
type DriveClient struct {
	svc     *drive.Service
	driveID string
	sem     chan struct{}

	mu    sync.Mutex
	calls map[string]*atomic.Int64
}

var _ DriveAPI = (*DriveClient)(nil)

// NewDriveClient expects the service-account key JSON encoded in base64.
func NewDriveClient(ctx context.Context, credsBase64, driveID string) (*DriveClient, error) {
	if credsBase64 == "" || driveID == "" {
		return nil, errors.New("GOOGLE_DRIVE_CREDENTIALS_JSON and PHOTOS_DRIVE_ID are required")
	}
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(credsBase64))
	if err != nil {
		return nil, fmt.Errorf("decode drive credentials: %w", err)
	}
	svc, err := drive.NewService(ctx, option.WithCredentialsJSON(key), option.WithScopes(drive.DriveReadonlyScope))
	if err != nil {
		return nil, fmt.Errorf("create drive service: %w", err)
	}
	return &DriveClient{svc: svc, driveID: driveID, sem: make(chan struct{}, maxInFlight), calls: map[string]*atomic.Int64{}}, nil
}

// Calls returns how many Drive requests were issued per operation since the client was created.
func (c *DriveClient) Calls() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.calls))
	for k, v := range c.calls {
		out[k] = v.Load()
	}
	return out
}

func (c *DriveClient) do(ctx context.Context, op string, fn func() error) error {
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return ctx.Err()
	}
	c.mu.Lock()
	ctr, ok := c.calls[op]
	if !ok {
		ctr = &atomic.Int64{}
		c.calls[op] = ctr
	}
	c.mu.Unlock()
	return retry(ctx, func() error {
		ctr.Add(1)
		return fn()
	})
}

func (c *DriveClient) StartPageToken(ctx context.Context) (string, error) {
	var tok *drive.StartPageToken
	err := c.do(ctx, "changes.getStartPageToken", func() (err error) {
		tok, err = c.svc.Changes.GetStartPageToken().DriveId(c.driveID).SupportsAllDrives(true).Context(ctx).Do()
		return
	})
	if err != nil {
		return "", err
	}
	return tok.StartPageToken, nil
}

const (
	qFolders = "trashed = false and (mimeType = '" + MimeFolder + "' or mimeType = '" + MimeShortcut + "')"
	qFiles   = "trashed = false and mimeType != '" + MimeFolder + "' and mimeType != '" + MimeShortcut + "'"
)

// ListFolders returns every folder and shortcut of the drive: the whole structure fits in a couple of calls.
func (c *DriveClient) ListFolders(ctx context.Context) ([]DriveItem, error) {
	var out []DriveItem
	token := ""
	for {
		items, next, err := c.list(ctx, "files.list(folders)", qFolders, token, true)
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
		if next == "" {
			return out, nil
		}
		token = next
	}
}

// ListFilesPage returns one page of non-folder files; pass the returned token to resume.
func (c *DriveClient) ListFilesPage(ctx context.Context, pageToken string) ([]DriveItem, string, error) {
	return c.list(ctx, "files.list(files)", qFiles, pageToken, true)
}

// ListChildren returns the direct children of a folder that lives outside the shared drive.
func (c *DriveClient) ListChildren(ctx context.Context, folderID string) ([]DriveItem, error) {
	var out []DriveItem
	token := ""
	q := fmt.Sprintf("'%s' in parents and trashed = false", strings.ReplaceAll(folderID, "'", ""))
	for {
		items, next, err := c.list(ctx, "files.list(children)", q, token, false)
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
		if next == "" {
			return out, nil
		}
		token = next
	}
}

func (c *DriveClient) list(ctx context.Context, op, q, token string, inDrive bool) ([]DriveItem, string, error) {
	call := c.svc.Files.List().IncludeItemsFromAllDrives(true).SupportsAllDrives(true).
		Q(q).PageSize(pageSize).Fields("nextPageToken, files(" + itemFields + ")").Context(ctx)
	if inDrive {
		call = call.Corpora("drive").DriveId(c.driveID)
	}
	if token != "" {
		call = call.PageToken(token)
	}
	var page *drive.FileList
	if err := c.do(ctx, op, func() (err error) { page, err = call.Do(); return }); err != nil {
		return nil, "", err
	}
	items := make([]DriveItem, 0, len(page.Files))
	for _, f := range page.Files {
		items = append(items, toItem(f))
	}
	return items, page.NextPageToken, nil
}

// Changes returns one page of the changes feed starting at pageToken.
func (c *DriveClient) Changes(ctx context.Context, pageToken string) (ChangePage, error) {
	call := c.svc.Changes.List(pageToken).DriveId(c.driveID).
		IncludeItemsFromAllDrives(true).SupportsAllDrives(true).IncludeRemoved(true).PageSize(pageSize).
		Fields("nextPageToken,newStartPageToken,changes(changeType,removed,fileId,file(" + itemFields + "))").Context(ctx)
	var res *drive.ChangeList
	if err := c.do(ctx, "changes.list", func() (err error) { res, err = call.Do(); return }); err != nil {
		return ChangePage{}, err
	}
	page := ChangePage{NextPageToken: res.NextPageToken, NewStartPageToken: res.NewStartPageToken}
	for _, ch := range res.Changes {
		if ch.ChangeType != "file" {
			continue
		}
		out := Change{FileID: ch.FileId, Removed: ch.Removed}
		if ch.File != nil && !ch.Removed {
			it := toItem(ch.File)
			out.Item = &it
		}
		page.Changes = append(page.Changes, out)
	}
	return page, nil
}

// GetItem returns a single file or folder, or an error when the caller cannot read it.
func (c *DriveClient) GetItem(ctx context.Context, id string) (DriveItem, error) {
	var f *drive.File
	err := c.do(ctx, "files.get", func() (err error) {
		f, err = c.svc.Files.Get(id).SupportsAllDrives(true).Fields(itemFields).Context(ctx).Do()
		return
	})
	if err != nil {
		return DriveItem{}, err
	}
	return toItem(f), nil
}

// retry backs off on Drive quota errors (403 rate limit, 429) and transient 5xx.
func retry(ctx context.Context, fn func() error) error {
	delay := 2 * time.Second
	for attempt := 0; ; attempt++ {
		err := fn()
		var gerr *googleapi.Error
		if err == nil || attempt >= 6 || !errors.As(err, &gerr) || !retryable(gerr) {
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

func retryable(e *googleapi.Error) bool {
	if e.Code == 429 || e.Code >= 500 {
		return true
	}
	for _, d := range e.Errors {
		if d.Reason == "rateLimitExceeded" || d.Reason == "userRateLimitExceeded" {
			return true
		}
	}
	return false
}

func toItem(f *drive.File) DriveItem {
	it := DriveItem{
		ID: f.Id, Name: f.Name, MimeType: f.MimeType, Size: f.Size, MD5: f.Md5Checksum,
		WebViewLink: f.WebViewLink, HasThumbnail: f.HasThumbnail, Trashed: f.Trashed,
	}
	if len(f.Parents) > 0 {
		it.ParentID = f.Parents[0]
	}
	it.CreatedTime, _ = time.Parse(time.RFC3339, f.CreatedTime)
	it.ModifiedTime, _ = time.Parse(time.RFC3339, f.ModifiedTime)
	if f.ShortcutDetails != nil {
		it.ShortcutTarget, it.ShortcutMime = f.ShortcutDetails.TargetId, f.ShortcutDetails.TargetMimeType
	}
	if m := f.ImageMediaMetadata; m != nil {
		it.ExifTime, it.Width, it.Height = m.Time, int(m.Width), int(m.Height)
	} else if m := f.VideoMediaMetadata; m != nil {
		it.Width, it.Height = int(m.Width), int(m.Height)
	}
	return it
}
