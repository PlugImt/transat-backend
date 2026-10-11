package photos

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/photos/repository"
	"github.com/plugimt/transat-backend/services/photos"
	"github.com/plugimt/transat-backend/utils"
)

const (
	defaultLimit = 60
	maxLimit     = 120
	albumLimit   = 50
	maxAlbums    = 100

	gridSize    = 320
	coverSize   = 400
	previewSize = 1200

	listMaxAge   = 60
	detailMaxAge = 300
)

// Google serves thumbnails of public files straight from lh3.googleusercontent.com: cacheable for a day, no hop.
// drive.google.com/thumbnail answers with an uncacheable 302 to the same image, so it only serves as a fallback.
const (
	thumbTemplate    = "https://lh3.googleusercontent.com/d/{id}=w{size}"
	fallbackTemplate = "https://drive.google.com/thumbnail?id={id}&sz=w{size}"
	downloadTemplate = "https://drive.google.com/uc?export=download&id={id}"
)

func expand(template, id string, size int) string {
	r := strings.NewReplacer("{id}", id, "{size}", strconv.Itoa(size))
	return r.Replace(template)
}

type PhotosHandler struct {
	repo  *repository.PhotosRepository
	sync  *photos.SyncService // nil when Drive credentials are not configured
	cache *responseCache
}

func NewPhotosHandler(db *sql.DB, sync *photos.SyncService) *PhotosHandler {
	repo := repository.NewPhotosRepository(db)
	return &PhotosHandler{repo: repo, sync: sync, cache: newResponseCache(repo.DataVersion)}
}

type photoDTO struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	SortAt    string `json:"sort_at"`
	AlbumID   *int64 `json:"album_id"`
	// The event's name, from the same album as AlbumID (null when the photo has no published event).
	AlbumTitle *string `json:"album_title"`
	// Null when Drive has no thumbnail for the file (RAW images, …): show a placeholder.
	ThumbnailURL         *string `json:"thumbnail_url"`
	ThumbnailFallbackURL *string `json:"thumbnail_fallback_url"`
	PreviewURL           *string `json:"preview_url"`
	// Opens the original in Google Drive. Videos are never streamed by this API.
	DriveURL string `json:"drive_url"`
}

type photoDetailDTO struct {
	photoDTO
	Name                      string     `json:"name"`
	MimeType                  string     `json:"mime_type"`
	SizeBytes                 int64      `json:"size_bytes"`
	TakenAt                   *time.Time `json:"taken_at"`
	TakenAtSource             string     `json:"taken_at_source"`
	Path                      string     `json:"path"`
	Section                   string     `json:"section"`
	FolderURL                 string     `json:"folder_url"`
	DownloadURL               string     `json:"download_url"`
	ThumbnailTemplate         string     `json:"thumbnail_template"`
	ThumbnailFallbackTemplate string     `json:"thumbnail_fallback_template"`
	Album                     *albumRef  `json:"album"`
}

type albumRef struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	DriveURL string `json:"drive_url"`
}

type albumDTO struct {
	ID           int64      `json:"id"`
	ParentID     *int64     `json:"parent_id"`
	Title        string     `json:"title"`
	EventDate    *string    `json:"event_date"`
	AcademicYear string     `json:"academic_year"`
	Category     string     `json:"category"`
	PhotoCount   int        `json:"photo_count"`
	CoverURL     *string    `json:"cover_url"`
	FirstTakenAt *time.Time `json:"first_taken_at"`
	LastTakenAt  *time.Time `json:"last_taken_at"`
	// Always the event's own Drive folder, never a generic parent such as "Sports" or "Autres".
	DriveURL string `json:"drive_url"`
}

type photosPage struct {
	Photos     []photoDTO `json:"photos"`
	NextCursor *string    `json:"next_cursor"`
}

type albumsPage struct {
	Albums     []albumDTO `json:"albums"`
	NextCursor *string    `json:"next_cursor"`
}

func strPtr(s string) *string { return &s }

func toPhotoDTO(p repository.PhotoRow) photoDTO {
	d := photoDTO{
		ID: p.ID, MediaType: p.MediaType, Width: p.Width, Height: p.Height,
		SortAt: p.SortAt.UTC().Format(time.RFC3339), DriveURL: p.DriveURL,
	}
	if p.AlbumID.Valid {
		d.AlbumID = &p.AlbumID.Int64
	}
	if p.AlbumTitle.Valid {
		d.AlbumTitle = &p.AlbumTitle.String
	}
	if p.HasThumbnail {
		d.ThumbnailURL = strPtr(expand(thumbTemplate, p.ID, gridSize))
		d.ThumbnailFallbackURL = strPtr(expand(fallbackTemplate, p.ID, gridSize))
		d.PreviewURL = strPtr(expand(thumbTemplate, p.ID, previewSize))
	}
	return d
}

func toAlbumDTO(a repository.AlbumRow) albumDTO {
	d := albumDTO{
		ID: a.ID, Title: a.Title, AcademicYear: a.AcademicYear, Category: a.Category, PhotoCount: a.PhotoCount,
		DriveURL: a.DriveURL,
	}
	if a.ParentID.Valid {
		d.ParentID = &a.ParentID.Int64
	}
	if a.EventDate.Valid {
		d.EventDate = strPtr(a.EventDate.Time.Format("2006-01-02"))
	}
	if a.CoverID.Valid {
		d.CoverURL = strPtr(expand(thumbTemplate, a.CoverID.String, coverSize))
	}
	if a.FirstTakenAt.Valid {
		d.FirstTakenAt = &a.FirstTakenAt.Time
	}
	if a.LastTakenAt.Valid {
		d.LastTakenAt = &a.LastTakenAt.Time
	}
	return d
}

// ---- cursors: opaque, stable under inserts because they are positions, not offsets.

func encodeCursor(parts ...string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strings.Join(parts, "|")))
}

func decodeCursor(s string, n int) ([]string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(string(raw), "|", n)
	if len(parts) != n || parts[n-1] == "" {
		return nil, errors.New("malformed cursor")
	}
	return parts, nil
}

func photoCursor(p repository.PhotoRow) string {
	return encodeCursor(p.SortAt.UTC().Format(time.RFC3339Nano), strconv.FormatInt(p.AlbumKey, 10), p.ID)
}

func parsePhotoCursor(s string) (*repository.PhotoCursor, error) {
	if s == "" {
		return nil, nil
	}
	parts, err := decodeCursor(s, 3)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, err
	}
	key, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &repository.PhotoCursor{SortAt: t, AlbumKey: key, ID: parts[2]}, nil
}

func albumCursor(a repository.AlbumRow) string {
	return encodeCursor(a.SortDate.Format("2006-01-02"), strconv.FormatInt(a.ID, 10))
}

func parseAlbumCursor(s string) (*repository.AlbumCursor, error) {
	if s == "" {
		return nil, nil
	}
	parts, err := decodeCursor(s, 2)
	if err != nil {
		return nil, err
	}
	d, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &repository.AlbumCursor{SortDate: d, ID: id}, nil
}

func limitParam(c *fiber.Ctx, def, max int) int {
	n, err := strconv.Atoi(c.Query("limit", strconv.Itoa(def)))
	if err != nil || n < 1 {
		return def
	}
	return min(n, max)
}

func photosResponse(rows []repository.PhotoRow, limit int) photosPage {
	out := photosPage{Photos: []photoDTO{}}
	// One extra row was requested to know whether another page exists.
	if len(rows) > limit {
		rows = rows[:limit]
		out.NextCursor = strPtr(photoCursor(rows[limit-1]))
	}
	for _, r := range rows {
		out.Photos = append(out.Photos, toPhotoDTO(r))
	}
	return out
}

func fail(c *fiber.Ctx, status int, msg string, err error) error {
	if err != nil {
		utils.LogMessage(utils.LevelError, msg)
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
	}
	return c.Status(status).JSON(fiber.Map{"error": msg})
}

// serve answers from the shared in-process cache and lets clients and the ETag middleware reuse the body.
func (h *PhotosHandler) serve(c *fiber.Ctx, maxAge int, build func(ctx context.Context) (any, error)) error {
	key := c.Path() + "?" + string(c.Request().URI().QueryString())
	body, err := h.cache.get(c.Context(), key, func() ([]byte, error) {
		v, err := build(c.Context())
		if err != nil {
			return nil, err
		}
		return json.Marshal(v)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fail(c, fiber.StatusNotFound, "Not found", nil)
		}
		return fail(c, fiber.StatusInternalServerError, "Failed to fetch photos", err)
	}
	return send(c, maxAge, body)
}

// serveDirect answers without the shared cache. Photo details are a single primary-key lookup, and thousands of
// distinct ids would otherwise evict the hot gallery pages; clients still cache them through Cache-Control.
func (h *PhotosHandler) serveDirect(c *fiber.Ctx, maxAge int, build func(ctx context.Context) (any, error)) error {
	v, err := build(c.Context())
	if err == nil {
		var body []byte
		if body, err = json.Marshal(v); err == nil {
			return send(c, maxAge, body)
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	return fail(c, fiber.StatusInternalServerError, "Failed to fetch photo", err)
}

func send(c *fiber.Ctx, maxAge int, body []byte) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	c.Set(fiber.HeaderCacheControl, fmt.Sprintf("private, max-age=%d", maxAge))
	return c.Send(body)
}

// GetGallery handles GET /photos?cursor=&limit= (general gallery, newest first).
func (h *PhotosHandler) GetGallery(c *fiber.Ctx) error {
	cursor, err := parsePhotoCursor(c.Query("cursor"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid cursor", nil)
	}
	limit := limitParam(c, defaultLimit, maxLimit)

	return h.serve(c, listMaxAge, func(ctx context.Context) (any, error) {
		rows, err := h.repo.Gallery(ctx, cursor, limit+1)
		return photosResponse(rows, limit), err
	})
}

// GetAlbums handles GET /photos/albums?cursor=&limit= (newest event first).
func (h *PhotosHandler) GetAlbums(c *fiber.Ctx) error {
	cursor, err := parseAlbumCursor(c.Query("cursor"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid cursor", nil)
	}
	limit := limitParam(c, albumLimit, maxAlbums)
	return h.serve(c, listMaxAge, func(ctx context.Context) (any, error) {
		rows, err := h.repo.Albums(ctx, cursor, limit+1)
		if err != nil {
			return nil, err
		}
		out := albumsPage{Albums: []albumDTO{}}
		if len(rows) > limit {
			rows = rows[:limit]
			out.NextCursor = strPtr(albumCursor(rows[limit-1]))
		}
		for _, a := range rows {
			out.Albums = append(out.Albums, toAlbumDTO(a))
		}
		return out, nil
	})
}

// GetAlbum handles GET /photos/albums/:id.
func (h *PhotosHandler) GetAlbum(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid album id", nil)
	}
	return h.serve(c, listMaxAge, func(ctx context.Context) (any, error) {
		a, err := h.repo.Album(ctx, int64(id))
		return toAlbumDTO(a), err
	})
}

// GetAlbumPhotos handles GET /photos/albums/:id/photos?cursor=&limit= (newest first).
func (h *PhotosHandler) GetAlbumPhotos(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid album id", nil)
	}
	cursor, err := parsePhotoCursor(c.Query("cursor"))
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid cursor", nil)
	}
	limit := limitParam(c, defaultLimit, maxLimit)
	return h.serve(c, listMaxAge, func(ctx context.Context) (any, error) {
		rows, err := h.repo.AlbumPhotos(ctx, int64(id), cursor, limit+1)
		return photosResponse(rows, limit), err
	})
}

// GetPhoto handles GET /photos/:id (full metadata, links to the original and its folders).
func (h *PhotosHandler) GetPhoto(c *fiber.Ctx) error {
	id := c.Params("id")
	if len(id) < 10 || len(id) > 100 {
		return fail(c, fiber.StatusBadRequest, "Invalid photo id", nil)
	}
	return h.serveDirect(c, detailMaxAge, func(ctx context.Context) (any, error) {
		d, err := h.repo.Photo(ctx, id)
		if err != nil {
			return nil, err
		}
		out := photoDetailDTO{
			photoDTO: toPhotoDTO(d.PhotoRow), Name: d.Name, MimeType: d.MimeType, SizeBytes: d.SizeBytes,
			TakenAtSource: d.TakenSource, Path: d.FolderPath, Section: d.Section, FolderURL: d.FolderURL,
			DownloadURL:               expand(downloadTemplate, d.ID, 0),
			ThumbnailTemplate:         strings.ReplaceAll(thumbTemplate, "{id}", d.ID),
			ThumbnailFallbackTemplate: strings.ReplaceAll(fallbackTemplate, "{id}", d.ID),
		}
		if d.TakenAt.Valid {
			out.TakenAt = &d.TakenAt.Time
		}
		if d.AlbumID.Valid {
			out.Album = &albumRef{ID: d.AlbumID.Int64, Title: d.AlbumTitle.String, DriveURL: d.AlbumURL.String}
		}
		return out, nil
	})
}

// ---- admin

// GetSyncStatus handles GET /admin/photos/sync.
func (h *PhotosHandler) GetSyncStatus(c *fiber.Ctx) error {
	s, err := h.repo.SyncStatus(c.Context())
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Failed to fetch sync status", err)
	}
	type run struct {
		ID         int64           `json:"id"`
		Mode       string          `json:"mode"`
		StartedAt  time.Time       `json:"started_at"`
		FinishedAt *time.Time      `json:"finished_at"`
		Complete   bool            `json:"complete"`
		Stats      json.RawMessage `json:"stats"`
		Error      *string         `json:"error"`
	}
	runs := make([]run, 0, len(s.Runs))
	for _, r := range s.Runs {
		item := run{ID: r.ID, Mode: r.Mode, StartedAt: r.StartedAt, Complete: r.Complete, Stats: json.RawMessage(r.Stats)}
		if r.FinishedAt.Valid {
			item.FinishedAt = &r.FinishedAt.Time
		}
		if r.Error.Valid {
			item.Error = &r.Error.String
		}
		runs = append(runs, item)
	}
	out := fiber.Map{"data_version": s.DataVersion, "full_sync_in_progress": s.FullInProgress, "runs": runs, "enabled": h.sync != nil, "cache": h.cache.Stats()}
	if s.LastFullAt.Valid {
		out["last_full_at"] = s.LastFullAt.Time
	}
	if s.LastIncrementalAt.Valid {
		out["last_incremental_at"] = s.LastIncrementalAt.Time
	}
	return c.JSON(out)
}

// TriggerSync handles POST /admin/photos/sync?mode=full|incremental; it runs in the background.
func (h *PhotosHandler) TriggerSync(c *fiber.Ctx) error {
	if h.sync == nil {
		return fail(c, fiber.StatusServiceUnavailable, "Drive credentials are not configured", nil)
	}
	mode := photos.Mode(c.Query("mode"))
	if mode != "" && mode != photos.ModeFull && mode != photos.ModeIncremental {
		return fail(c, fiber.StatusBadRequest, "mode must be full or incremental", nil)
	}
	go func() {
		var err error
		if mode == "" {
			_, err = h.sync.Tick(context.Background())
		} else {
			_, err = h.sync.Run(context.Background(), mode)
		}
		if err != nil && !errors.Is(err, photos.ErrSyncRunning) {
			utils.LogMessage(utils.LevelError, "Photos sync failed")
			utils.LogLineKeyValue(utils.LevelError, "Error", err)
		}
	}()
	return c.SendStatus(fiber.StatusAccepted)
}
