package photos

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/photos/repository"
	"github.com/plugimt/transat-backend/services/photos"
)

// Public album pages (transat.dev/albums/<id>-<slug>): no account needed, only albums that are published and not
// hidden, only their images. The site's server is the only intended caller; it proxies every image so the
// visitor never sees a Drive link, and the Drive files themselves stay private.

const (
	publicGridSize    = 640
	publicPreviewSize = 2048
	publicListMaxAge  = 60
	publicMediaMaxAge = 3600
	mediaTimeout      = 2 * time.Minute
)

var driveIDRe = regexp.MustCompile(`^[A-Za-z0-9_-]{10,100}$`)

type publicAlbumDTO struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Slug         string     `json:"slug"`
	EventDate    *string    `json:"event_date"`
	AcademicYear string     `json:"academic_year"`
	Category     string     `json:"category"`
	PhotoCount   int        `json:"photo_count"`
	FirstTakenAt *time.Time `json:"first_taken_at"`
	LastTakenAt  *time.Time `json:"last_taken_at"`
	// The photo id of the cover, served by the cover endpoint; null when the album has none.
	CoverID  *string `json:"cover_id"`
	ShareURL string  `json:"share_url"`
}

type publicPhotoDTO struct {
	ID           string `json:"id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	SortAt       string `json:"sort_at"`
	HasThumbnail bool   `json:"has_thumbnail"`
}

type publicPhotosPage struct {
	Photos     []publicPhotoDTO `json:"photos"`
	NextCursor *string          `json:"next_cursor"`
}

type publicFileDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	MimeType  string `json:"mime_type"`
	SizeBytes int64  `json:"size_bytes"`
}

func toPublicAlbum(a repository.AlbumRow) publicAlbumDTO {
	d := publicAlbumDTO{
		ID: a.ID, Title: a.Title, Slug: photos.Slug(a.Title), AcademicYear: a.AcademicYear, Category: a.Category,
		PhotoCount: a.PhotoCount, ShareURL: photos.AlbumShareURL(a.ID, a.Title),
	}
	if a.EventDate.Valid {
		d.EventDate = strPtr(a.EventDate.Time.Format("2006-01-02"))
	}
	if a.FirstTakenAt.Valid {
		d.FirstTakenAt = &a.FirstTakenAt.Time
	}
	if a.LastTakenAt.Valid {
		d.LastTakenAt = &a.LastTakenAt.Time
	}
	if a.CoverID.Valid {
		d.CoverID = &a.CoverID.String
	}
	return d
}

// albumID reads ":id", which may be the whole "<id>-<slug>" path segment of a share link.
func albumID(c *fiber.Ctx) (int64, bool) {
	raw := c.Params("id")
	if i := strings.IndexByte(raw, '-'); i > 0 {
		raw = raw[:i]
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	return id, err == nil && id > 0
}

func sendPublic(c *fiber.Ctx, maxAge int, body []byte) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", maxAge))
	c.Set("X-Robots-Tag", "noindex")
	return c.Send(body)
}

// GetPublicAlbum handles GET /public/albums/:id.
func (h *PhotosHandler) GetPublicAlbum(c *fiber.Ctx) error {
	id, ok := albumID(c)
	if !ok {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	return h.servePublic(c, "album", id, func(ctx context.Context) (any, error) {
		a, err := h.repo.Album(ctx, id)
		if err == nil && a.PhotoCount == 0 {
			err = sql.ErrNoRows
		}
		return toPublicAlbum(a), err
	})
}

// GetPublicAlbumPhotos handles GET /public/albums/:id/photos?cursor=&limit= (newest first, images only).
func (h *PhotosHandler) GetPublicAlbumPhotos(c *fiber.Ctx) error {
	id, ok := albumID(c)
	if !ok {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	cursor, err := parsePhotoCursor(c.Query("cursor"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid cursor", nil)
	}
	limit := limitParam(c, defaultLimit, maxLimit)
	return h.servePublic(c, "photos", id, func(ctx context.Context) (any, error) {
		// The album must exist and be shown: an empty page would otherwise still answer 200.
		if _, err := h.repo.Album(ctx, id); err != nil {
			return nil, err
		}
		rows, err := h.repo.AlbumImagesPage(ctx, id, cursor, limit+1)
		if err != nil {
			return nil, err
		}
		out := publicPhotosPage{Photos: []publicPhotoDTO{}}
		if len(rows) > limit {
			rows = rows[:limit]
			out.NextCursor = strPtr(photoCursor(rows[limit-1]))
		}
		for _, r := range rows {
			out.Photos = append(out.Photos, publicPhotoDTO{
				ID: r.ID, Width: r.Width, Height: r.Height, SortAt: r.SortAt.UTC().Format(time.RFC3339),
				HasThumbnail: r.HasThumbnail,
			})
		}
		return out, nil
	})
}

// GetPublicAlbumFiles handles GET /public/albums/:id/files: every image of the album, oldest first, for the site
// to build the whole-album download one original at a time.
func (h *PhotosHandler) GetPublicAlbumFiles(c *fiber.Ctx) error {
	id, ok := albumID(c)
	if !ok {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	return h.servePublic(c, "files", id, func(ctx context.Context) (any, error) {
		a, err := h.repo.Album(ctx, id)
		if err != nil {
			return nil, err
		}
		files, err := h.repo.AlbumImages(ctx, id)
		if err != nil {
			return nil, err
		}
		out := struct {
			Title string          `json:"title"`
			Slug  string          `json:"slug"`
			Files []publicFileDTO `json:"files"`
		}{Title: a.Title, Slug: photos.Slug(a.Title), Files: make([]publicFileDTO, 0, len(files))}
		for _, f := range files {
			out.Files = append(out.Files, publicFileDTO{ID: f.ID, Name: f.Name, MimeType: f.MimeType, SizeBytes: f.SizeBytes})
		}
		return out, nil
	})
}

func (h *PhotosHandler) servePublic(c *fiber.Ctx, kind string, id int64, build func(ctx context.Context) (any, error)) error {
	// Keyed by the album id, not by the path: "42" and "42-gala" are the same album.
	key := fmt.Sprintf("public:%s:%d?%s", kind, id, c.Request().URI().QueryString())
	body, err := h.cache.get(c.Context(), key, func() ([]byte, error) {
		v, err := build(c.Context())
		if err != nil {
			return nil, err
		}
		return json.Marshal(v)
	})
	if errors.Is(err, sql.ErrNoRows) {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Failed to fetch album", err)
	}
	return sendPublic(c, publicListMaxAge, body)
}

// GetPublicCover handles GET /public/albums/:id/cover: the album's cover at preview size (share cards, header).
func (h *PhotosHandler) GetPublicCover(c *fiber.Ctx) error {
	id, ok := albumID(c)
	if !ok {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	a, err := h.repo.Album(c.Context(), id)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && !a.CoverID.Valid) {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Failed to fetch album", err)
	}
	return h.streamThumbnail(c, id, a.CoverID.String, 1200)
}

// GetPublicPhoto handles GET /public/albums/:id/photos/:photo/:variant where variant is "grid", "preview" or
// "original" (?download=1 to save it under its own name).
func (h *PhotosHandler) GetPublicPhoto(c *fiber.Ctx) error {
	id, ok := albumID(c)
	photoID := c.Params("photo")
	if !ok || !driveIDRe.MatchString(photoID) {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	switch c.Params("variant") {
	case "grid":
		return h.streamThumbnail(c, id, photoID, publicGridSize)
	case "preview":
		return h.streamThumbnail(c, id, photoID, publicPreviewSize)
	case "original":
		return h.streamOriginal(c, id, photoID)
	}
	return fail(c, fiber.StatusNotFound, "Not found", nil)
}

// mediaFile finds an image of a shown album. When ok is false the response is already written.
func (h *PhotosHandler) mediaFile(c *fiber.Ctx, albumID int64, photoID string) (f repository.PhotoFile, ok bool, err error) {
	if h.media == nil {
		return f, false, fail(c, fiber.StatusServiceUnavailable, "Photos are not available", nil)
	}
	f, err = h.repo.AlbumPhotoFile(c.Context(), albumID, photoID)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && f.MediaType != "image") {
		return f, false, fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	if err != nil {
		return f, false, fail(c, fiber.StatusInternalServerError, "Failed to fetch photo", err)
	}
	return f, true, nil
}

func (h *PhotosHandler) streamThumbnail(c *fiber.Ctx, albumID int64, photoID string, size int) error {
	f, ok, err := h.mediaFile(c, albumID, photoID)
	if !ok {
		return err
	}
	if !f.HasThumbnail {
		return fail(c, fiber.StatusNotFound, "No preview for this photo", nil)
	}
	ctx, cancel := context.WithTimeout(context.Background(), mediaTimeout)
	m, err := h.media.Thumbnail(ctx, f.ID, size)
	return h.sendMedia(c, cancel, m, err, "inline", "")
}

func (h *PhotosHandler) streamOriginal(c *fiber.Ctx, albumID int64, photoID string) error {
	f, ok, err := h.mediaFile(c, albumID, photoID)
	if !ok {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), mediaTimeout)
	m, err := h.media.Original(ctx, f.ID)
	if err == nil && f.MimeType != "" {
		m.ContentType = f.MimeType
	}
	disposition := "inline"
	if c.Query("download") == "1" {
		disposition = "attachment"
	}
	return h.sendMedia(c, cancel, m, err, disposition, f.Name)
}

// sendMedia streams m; the body is read after the handler returns, so cancel runs when the stream is closed.
func (h *PhotosHandler) sendMedia(c *fiber.Ctx, cancel context.CancelFunc, m photos.Media, err error, disposition, name string) error {
	if err != nil {
		cancel()
		if errors.Is(err, photos.ErrMediaNotFound) {
			return fail(c, fiber.StatusNotFound, "Not found", nil)
		}
		return fail(c, fiber.StatusBadGateway, "Drive is unavailable", err)
	}
	c.Set(fiber.HeaderContentType, m.ContentType)
	c.Set(fiber.HeaderCacheControl, fmt.Sprintf("public, max-age=%d", publicMediaMaxAge))
	c.Set("X-Robots-Tag", "noindex")
	c.Set("X-Content-Type-Options", "nosniff")
	if name != "" {
		c.Set(fiber.HeaderContentDisposition, contentDisposition(disposition, name))
	}
	body := &cancelOnClose{ReadCloser: m.Body, cancel: cancel}
	if m.Size >= 0 {
		return c.SendStream(body, int(m.Size))
	}
	return c.SendStream(body)
}

func contentDisposition(disposition, name string) string {
	name = path.Base(strings.ReplaceAll(name, `\`, "/"))
	var ascii strings.Builder
	for _, r := range name {
		if r < 0x20 || r > 0x7e || r == '"' || r == '\\' {
			ascii.WriteByte('_')
		} else {
			ascii.WriteRune(r)
		}
	}
	return fmt.Sprintf(`%s; filename="%s"; filename*=UTF-8''%s`, disposition, ascii.String(), url.PathEscape(name))
}

type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (r *cancelOnClose) Close() error {
	err := r.ReadCloser.Close()
	r.cancel()
	return err
}
