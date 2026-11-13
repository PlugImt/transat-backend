package photos

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/photos/repository"
	"github.com/plugimt/transat-backend/services/photos"
)

const adminAlbumsPageSize = 50

type adminAlbumDTO struct {
	ID           int64      `json:"id"`
	ParentID     *int64     `json:"parent_id"`
	Title        string     `json:"title"`
	EventDate    *string    `json:"event_date"`
	AcademicYear string     `json:"academic_year"`
	Category     string     `json:"category"`
	PhotoCount   int        `json:"photo_count"`
	Visibility   string     `json:"visibility"`
	Hidden       bool       `json:"hidden"`
	HiddenAt     *time.Time `json:"hidden_at"`
	HiddenBy     *string    `json:"hidden_by"`
	FirstTakenAt *time.Time `json:"first_taken_at"`
	// Shown in the app and on the web: published by the classifier and not hidden.
	Shown bool `json:"shown"`
	// The album's transat.dev page; it answers only while the album is shown.
	ShareURL string `json:"share_url"`
	// The cover image for the panel (needs the panel's token), relative to the API; null when there is none.
	CoverPath *string `json:"cover_path"`
	DriveURL  string  `json:"drive_url"`
}

func toAdminAlbum(a repository.AdminAlbumRow) adminAlbumDTO {
	d := adminAlbumDTO{
		ID: a.ID, Title: a.Title, AcademicYear: a.AcademicYear, Category: a.Category, PhotoCount: a.PhotoCount,
		Visibility: a.Visibility, Hidden: a.Hidden, DriveURL: a.DriveURL,
		Shown:    a.Visibility == "published" && !a.Hidden,
		ShareURL: photos.AlbumShareURL(a.ID, a.Title),
	}
	if a.ParentID.Valid {
		d.ParentID = &a.ParentID.Int64
	}
	if a.EventDate.Valid {
		d.EventDate = strPtr(a.EventDate.Time.Format("2006-01-02"))
	}
	if a.HiddenAt.Valid {
		d.HiddenAt = &a.HiddenAt.Time
	}
	if a.HiddenBy.Valid {
		d.HiddenBy = &a.HiddenBy.String
	}
	if a.FirstTakenAt.Valid {
		d.FirstTakenAt = &a.FirstTakenAt.Time
	}
	if a.CoverID.Valid {
		d.CoverPath = strPtr("/admin/photos/albums/" + strconv.FormatInt(a.ID, 10) + "/cover")
	}
	return d
}

// ListAdminAlbums handles GET /admin/photos/albums?q=&status=&page= (status: shown, hidden, review).
func (h *PhotosHandler) ListAdminAlbums(c *fiber.Ctx) error {
	status := c.Query("status")
	switch status {
	case "", "shown", "hidden", "review":
	default:
		return fail(c, fiber.StatusBadRequest, "status must be shown, hidden or review", nil)
	}
	page := max(c.QueryInt("page", 1), 1)
	q := c.Query("q")
	if len(q) > 200 {
		q = q[:200]
	}
	rows, counts, err := h.repo.AdminAlbums(c.Context(), repository.AdminAlbumFilter{
		Query: q, Status: status, Limit: adminAlbumsPageSize, Offset: (page - 1) * adminAlbumsPageSize,
	})
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Failed to fetch albums", err)
	}
	albums := make([]adminAlbumDTO, 0, len(rows))
	for _, a := range rows {
		albums = append(albums, toAdminAlbum(a))
	}
	total := counts.All
	switch status {
	case "shown":
		total = counts.Shown
	case "hidden":
		total = counts.Hidden
	case "review":
		total = counts.Review
	}
	return c.JSON(fiber.Map{
		"albums": albums, "counts": counts, "page": page, "page_size": adminAlbumsPageSize, "total": total,
	})
}

// SetAlbumHidden handles PUT /admin/photos/albums/:id/hidden with {"hidden": true|false}.
func (h *PhotosHandler) SetAlbumHidden(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, fiber.StatusBadRequest, "Invalid album id", nil)
	}
	var body struct {
		Hidden *bool `json:"hidden"`
	}
	if err := c.BodyParser(&body); err != nil || body.Hidden == nil {
		return fail(c, fiber.StatusBadRequest, "hidden (boolean) is required", nil)
	}
	email, _ := c.Locals("email").(string)
	a, err := h.repo.SetAlbumHidden(c.Context(), id, *body.Hidden, email, time.Now())
	if errors.Is(err, sql.ErrNoRows) {
		return fail(c, fiber.StatusNotFound, "Album not found", nil)
	}
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Failed to update album", err)
	}
	h.cache.refresh()
	return c.JSON(toAdminAlbum(a))
}

// GetAdminAlbumCover handles GET /admin/photos/albums/:id/cover: the cover of any album, hidden ones included.
func (h *PhotosHandler) GetAdminAlbumCover(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id <= 0 {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	if h.media == nil {
		return fail(c, fiber.StatusServiceUnavailable, "Photos are not available", nil)
	}
	a, err := h.repo.AdminAlbum(c.Context(), id)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && !a.CoverID.Valid) {
		return fail(c, fiber.StatusNotFound, "Not found", nil)
	}
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "Failed to fetch album", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), mediaTimeout)
	m, err := h.media.Thumbnail(ctx, a.CoverID.String, 480)
	if err := h.sendMedia(c, cancel, m, err, "inline", ""); err != nil {
		return err
	}
	c.Set(fiber.HeaderCacheControl, "private, max-age=3600")
	return nil
}
