package photobooth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/lib/pq"
)

// Events prepared in the admin panel. The booths download them (settings, layouts switched on, frames) instead of
// every operator setting each tablet up by hand. The layouts' geometry belongs to the booth's app: the booths
// publish it as a catalog, which the panel uses to draw templates and check frames.

// PhotoAspects are the photo ratios a booth knows (the names of its PhotoAspect enum).
var PhotoAspects = []string{"RATIO_4_3", "RATIO_3_2", "RATIO_16_9"}

const (
	// BoothEventsLookback is how far back the booths still see past events (a party ending after midnight, a
	// late review on the tablet).
	BoothEventsLookback = 30 * 24 * time.Hour
	// MaxCatalogJSON fits every layout of every ratio.
	MaxCatalogJSON = 512 << 10

	minCountdown, maxCountdown = 1, 15
	minPrints, maxPrints       = 1, 20
)

// CatalogLayout is one layout of the booth's catalog.
type CatalogLayout struct {
	FrameLayout
	EnabledByDefault bool `json:"enabled_by_default"`
}

// LayoutCatalog is every layout of every photo ratio, in the booth's order.
type LayoutCatalog struct {
	Aspects   map[string][]CatalogLayout `json:"aspects"`
	UpdatedAt *time.Time                 `json:"updated_at"`
	Device    string                     `json:"device,omitempty"`
}

// EventSettings is what the admin panel edits.
type EventSettings struct {
	Name                string   `json:"name"`
	Date                string   `json:"date"` // yyyy-mm-dd
	PhotoAspect         string   `json:"photo_aspect"`
	CountdownSeconds    int      `json:"countdown_seconds"`
	MaxPrintsPerSession int      `json:"max_prints_per_session"`
	AllowRetakes        bool     `json:"allow_retakes"`
	DigitalShare        bool     `json:"digital_share"`
	AttractSlideshow    bool     `json:"attract_slideshow"`
	MirrorLivePreview   bool     `json:"mirror_live_preview"`
	EnabledLayouts      []string `json:"enabled_layouts"`
}

// DefaultEventSettings are the booth's own defaults for a new event.
func DefaultEventSettings() EventSettings {
	return EventSettings{PhotoAspect: PhotoAspects[0], CountdownSeconds: 5, MaxPrintsPerSession: 5, AllowRetakes: true,
		DigitalShare: true, AttractSlideshow: true, MirrorLivePreview: true}
}

// EventFrame describes the frame of one layout, without its bytes.
type EventFrame struct {
	SHA256     string    `json:"sha256"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Size       int       `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// EventLayout is one layout of an event: switched on or not, and its frame if one was sent. Layouts never touched
// are absent (off, no frame).
type EventLayout struct {
	Type    string      `json:"type"`
	Enabled bool        `json:"enabled"`
	Frame   *EventFrame `json:"frame"`
}

// BoothEvent is an event as the panel and the booths see it. EnabledLayouts repeats the layouts switched on.
type BoothEvent struct {
	ID int64 `json:"id"`
	EventSettings
	Layouts   []EventLayout `json:"layouts"`
	CreatedBy string        `json:"created_by,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// ---- Catalog --------------------------------------------------------------------------------------------

// PublishLayoutCatalog replaces the catalog with the booth's: its app is the reference for the geometry.
func (s *Service) PublishLayoutCatalog(ctx context.Context, device string, aspects map[string][]CatalogLayout) error {
	if len(aspects) == 0 {
		return invalid("the catalog is empty")
	}
	for aspect, layouts := range aspects {
		if !slices.Contains(PhotoAspects, aspect) {
			return invalid("unknown photo ratio %q", truncate(aspect, 20))
		}
		plain := make([]FrameLayout, len(layouts))
		for i := range layouts {
			layouts[i].Title = CleanEventName(layouts[i].Title)
			layouts[i].HasFrame = false
			plain[i] = layouts[i].FrameLayout
		}
		if err := validateFrameLayouts(plain); err != nil {
			return err
		}
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `DELETE FROM photobooth_layout_catalog`); err != nil {
		return err
	}
	for aspect, layouts := range aspects {
		for i, l := range layouts {
			raw, err := json.Marshal(l.FrameLayout)
			if err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO photobooth_layout_catalog
				(photo_aspect, layout_type, position, enabled_by_default, layout, device_key_id) VALUES ($1, $2, $3, $4, $5, $6)`,
				aspect, l.Type, i, l.EnabledByDefault, raw, device); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// LayoutCatalog returns the last catalog a booth published (empty before the first one).
func (s *Service) LayoutCatalog(ctx context.Context) (LayoutCatalog, error) {
	out := LayoutCatalog{Aspects: map[string][]CatalogLayout{}}
	rows, err := s.db.QueryContext(ctx, `SELECT photo_aspect, enabled_by_default, layout, device_key_id, updated_at
		FROM photobooth_layout_catalog ORDER BY photo_aspect, position`)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	for rows.Next() {
		var aspect string
		var raw []byte
		var at time.Time
		var l CatalogLayout
		if err := rows.Scan(&aspect, &l.EnabledByDefault, &raw, &out.Device, &at); err != nil {
			return out, err
		}
		if err := json.Unmarshal(raw, &l.FrameLayout); err != nil {
			return out, err
		}
		out.Aspects[aspect] = append(out.Aspects[aspect], l)
		if out.UpdatedAt == nil || at.After(*out.UpdatedAt) {
			out.UpdatedAt = &at
		}
	}
	return out, rows.Err()
}

func (s *Service) catalogLayouts(ctx context.Context, aspect string) (map[string]FrameLayout, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT layout FROM photobooth_layout_catalog WHERE photo_aspect = $1`, aspect)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]FrameLayout{}
	for rows.Next() {
		var raw []byte
		var l FrameLayout
		if err := rows.Scan(&raw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(raw, &l); err != nil {
			return nil, err
		}
		out[l.Type] = l
	}
	return out, rows.Err()
}

// ---- Events ---------------------------------------------------------------------------------------------

func (s *Service) validateSettings(ctx context.Context, in *EventSettings) error {
	in.Name = CleanEventName(in.Name)
	if in.Name == "" {
		return invalid("event name is required")
	}
	date, err := time.Parse("2006-01-02", in.Date)
	if err != nil || date.Year() < 2000 || date.Year() > 2100 {
		return invalid("event date must be yyyy-mm-dd")
	}
	if !slices.Contains(PhotoAspects, in.PhotoAspect) {
		return invalid("unknown photo ratio")
	}
	if in.CountdownSeconds < minCountdown || in.CountdownSeconds > maxCountdown {
		return invalid("countdown must be between %d and %d seconds", minCountdown, maxCountdown)
	}
	if in.MaxPrintsPerSession < minPrints || in.MaxPrintsPerSession > maxPrints {
		return invalid("prints per session must be between %d and %d", minPrints, maxPrints)
	}
	if len(in.EnabledLayouts) == 0 {
		return invalid("switch on at least one layout")
	}
	if len(in.EnabledLayouts) > maxFrameLayouts {
		return invalid("too many layouts")
	}
	catalog, err := s.catalogLayouts(ctx, in.PhotoAspect)
	if err != nil {
		return err
	}
	// Before any booth published its catalog the panel has nothing to offer, but the API stays usable.
	open := false
	if len(catalog) == 0 {
		if err := s.db.QueryRowContext(ctx, `SELECT NOT EXISTS (SELECT 1 FROM photobooth_layout_catalog)`).Scan(&open); err != nil {
			return err
		}
		if !open {
			return invalid("the booths have no layouts for this photo ratio")
		}
	}
	seen := map[string]bool{}
	for _, t := range in.EnabledLayouts {
		if !layoutTypeRe.MatchString(t) || seen[t] {
			return invalid("invalid or repeated layout type")
		}
		if _, ok := catalog[t]; !open && !ok {
			return invalid("the booths have no layout %s", t)
		}
		seen[t] = true
	}
	return nil
}

// CreateEvent records an event prepared in the admin panel by [by].
func (s *Service) CreateEvent(ctx context.Context, by string, in EventSettings) (BoothEvent, error) {
	if err := s.validateSettings(ctx, &in); err != nil {
		return BoothEvent{}, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return BoothEvent{}, err
	}
	defer func() { _ = tx.Rollback() }()
	var id int64
	if err := tx.QueryRowContext(ctx, `INSERT INTO photobooth_events (name, event_date, photo_aspect, countdown_seconds,
		max_prints_per_session, allow_retakes, digital_share, attract_slideshow, mirror_live_preview, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		in.Name, in.Date, in.PhotoAspect, in.CountdownSeconds, in.MaxPrintsPerSession, in.AllowRetakes, in.DigitalShare,
		in.AttractSlideshow, in.MirrorLivePreview, by).Scan(&id); err != nil {
		return BoothEvent{}, err
	}
	if err := setEnabledLayouts(ctx, tx, id, in.EnabledLayouts); err != nil {
		return BoothEvent{}, err
	}
	if err := tx.Commit(); err != nil {
		return BoothEvent{}, err
	}
	return s.Event(ctx, id)
}

// UpdateEvent replaces an event's settings and the layouts switched on; frames are kept, even on layouts
// switched off.
func (s *Service) UpdateEvent(ctx context.Context, id int64, in EventSettings) (BoothEvent, error) {
	if err := s.validateSettings(ctx, &in); err != nil {
		return BoothEvent{}, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return BoothEvent{}, err
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.ExecContext(ctx, `UPDATE photobooth_events SET name = $2, event_date = $3, photo_aspect = $4,
		countdown_seconds = $5, max_prints_per_session = $6, allow_retakes = $7, digital_share = $8,
		attract_slideshow = $9, mirror_live_preview = $10, updated_at = now() WHERE id = $1`,
		id, in.Name, in.Date, in.PhotoAspect, in.CountdownSeconds, in.MaxPrintsPerSession, in.AllowRetakes,
		in.DigitalShare, in.AttractSlideshow, in.MirrorLivePreview)
	if err != nil {
		return BoothEvent{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return BoothEvent{}, errNotFound
	}
	if err := setEnabledLayouts(ctx, tx, id, in.EnabledLayouts); err != nil {
		return BoothEvent{}, err
	}
	if err := tx.Commit(); err != nil {
		return BoothEvent{}, err
	}
	return s.Event(ctx, id)
}

func setEnabledLayouts(ctx context.Context, tx *sql.Tx, id int64, enabled []string) error {
	if _, err := tx.ExecContext(ctx, `UPDATE photobooth_event_layouts SET enabled = (layout_type = ANY($2))
		WHERE event_id = $1`, id, pq.Array(enabled)); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO photobooth_event_layouts (event_id, layout_type, enabled)
		SELECT $1, t, TRUE FROM unnest($2::text[]) AS t ON CONFLICT (event_id, layout_type) DO NOTHING`, id, pq.Array(enabled))
	return err
}

// DeleteEvent removes an event and its frames. The booths keep their copy (and the photos taken with it).
func (s *Service) DeleteEvent(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM photobooth_events WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errNotFound
	}
	return nil
}

// Event returns one event.
func (s *Service) Event(ctx context.Context, id int64) (BoothEvent, error) {
	list, err := s.listEvents(ctx, `WHERE id = $1`, id)
	if err != nil {
		return BoothEvent{}, err
	}
	if len(list) == 0 {
		return BoothEvent{}, errNotFound
	}
	return list[0], nil
}

// Events lists every event, the most recent date first, for the admin panel.
func (s *Service) Events(ctx context.Context) ([]BoothEvent, error) {
	return s.listEvents(ctx, `ORDER BY event_date DESC, id DESC`)
}

// BoothEvents lists what a booth may run: events of the last BoothEventsLookback and upcoming ones, soonest
// first.
func (s *Service) BoothEvents(ctx context.Context) ([]BoothEvent, error) {
	from := s.now().Add(-BoothEventsLookback).Format("2006-01-02")
	return s.listEvents(ctx, `WHERE event_date >= $1 ORDER BY event_date, id`, from)
}

func (s *Service) listEvents(ctx context.Context, tail string, args ...any) ([]BoothEvent, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, to_char(event_date, 'YYYY-MM-DD'), photo_aspect,
		countdown_seconds, max_prints_per_session, allow_retakes, digital_share, attract_slideshow, mirror_live_preview,
		created_by, created_at, updated_at FROM photobooth_events `+tail, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BoothEvent{}
	index := map[int64]int{}
	var ids []int64
	for rows.Next() {
		var e BoothEvent
		if err := rows.Scan(&e.ID, &e.Name, &e.Date, &e.PhotoAspect, &e.CountdownSeconds, &e.MaxPrintsPerSession,
			&e.AllowRetakes, &e.DigitalShare, &e.AttractSlideshow, &e.MirrorLivePreview, &e.CreatedBy,
			&e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		e.Layouts = []EventLayout{}
		e.EnabledLayouts = []string{}
		index[e.ID] = len(out)
		ids = append(ids, e.ID)
		out = append(out, e)
	}
	if err := rows.Err(); err != nil || len(ids) == 0 {
		return out, err
	}
	rows.Close()

	lrows, err := s.db.QueryContext(ctx, `SELECT event_id, layout_type, enabled, frame_sha256, frame_width, frame_height,
		length(frame_png), frame_uploaded_at FROM photobooth_event_layouts WHERE event_id = ANY($1)
		ORDER BY event_id, layout_type`, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer lrows.Close()
	for lrows.Next() {
		var id int64
		var l EventLayout
		var sha sql.NullString
		var w, h, size sql.NullInt64
		var at sql.NullTime
		if err := lrows.Scan(&id, &l.Type, &l.Enabled, &sha, &w, &h, &size, &at); err != nil {
			return nil, err
		}
		if sha.Valid {
			l.Frame = &EventFrame{SHA256: sha.String, Width: int(w.Int64), Height: int(h.Int64), Size: int(size.Int64), UploadedAt: at.Time}
		}
		e := &out[index[id]]
		e.Layouts = append(e.Layouts, l)
		if l.Enabled {
			e.EnabledLayouts = append(e.EnabledLayouts, l.Type)
		}
	}
	return out, lrows.Err()
}

// ---- Frames ---------------------------------------------------------------------------------------------

// SetEventFrame stores the frame of one layout of an event, checked against the booths' catalog.
func (s *Service) SetEventFrame(ctx context.Context, id int64, layoutType string, data []byte) (EventFrame, error) {
	var aspect string
	err := s.db.QueryRowContext(ctx, `SELECT photo_aspect FROM photobooth_events WHERE id = $1`, id).Scan(&aspect)
	if errors.Is(err, sql.ErrNoRows) {
		return EventFrame{}, errNotFound
	}
	if err != nil {
		return EventFrame{}, err
	}
	catalog, err := s.catalogLayouts(ctx, aspect)
	if err != nil {
		return EventFrame{}, err
	}
	layout, ok := catalog[layoutType]
	if !ok {
		return EventFrame{}, invalid("no booth has published this layout yet: open the booth's app once, connected to the server")
	}
	u, err := checkFrame(layout, data)
	if err != nil {
		return EventFrame{}, err
	}
	return s.storeEventFrame(ctx, id, u, data, false)
}

// storeEventFrame files a checked frame in an event; a layout new to the event starts switched on if [enable].
func (s *Service) storeEventFrame(ctx context.Context, id int64, u FrameUpload, data []byte, enable bool) (EventFrame, error) {
	f := EventFrame{SHA256: u.SHA256, Width: u.Width, Height: u.Height, Size: u.Size}
	err := s.db.QueryRowContext(ctx, `INSERT INTO photobooth_event_layouts
		(event_id, layout_type, enabled, frame_png, frame_sha256, frame_width, frame_height, frame_uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		ON CONFLICT (event_id, layout_type) DO UPDATE SET frame_png = EXCLUDED.frame_png,
		frame_sha256 = EXCLUDED.frame_sha256, frame_width = EXCLUDED.frame_width, frame_height = EXCLUDED.frame_height,
		frame_uploaded_at = EXCLUDED.frame_uploaded_at
		RETURNING frame_uploaded_at`, id, u.LayoutType, enable, data, u.SHA256, u.Width, u.Height).Scan(&f.UploadedAt)
	if err != nil {
		return f, err
	}
	_, err = s.db.ExecContext(ctx, `UPDATE photobooth_events SET updated_at = now() WHERE id = $1`, id)
	return f, err
}

// DeleteEventFrame removes the frame of one layout (the booths then print it without one).
func (s *Service) DeleteEventFrame(ctx context.Context, id int64, layoutType string) error {
	res, err := s.db.ExecContext(ctx, `UPDATE photobooth_event_layouts SET frame_png = NULL, frame_sha256 = NULL,
		frame_width = NULL, frame_height = NULL, frame_uploaded_at = NULL WHERE event_id = $1 AND layout_type = $2
		AND frame_sha256 IS NOT NULL`, id, layoutType)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errNotFound
	}
	_, err = s.db.ExecContext(ctx, `UPDATE photobooth_events SET updated_at = now() WHERE id = $1`, id)
	return err
}

// EventFrameImage returns the PNG of one layout of an event.
func (s *Service) EventFrameImage(ctx context.Context, id int64, layoutType string) ([]byte, error) {
	var png []byte
	err := s.db.QueryRowContext(ctx, `SELECT frame_png FROM photobooth_event_layouts
		WHERE event_id = $1 AND layout_type = $2 AND frame_png IS NOT NULL`, id, layoutType).Scan(&png)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}
	return png, err
}
