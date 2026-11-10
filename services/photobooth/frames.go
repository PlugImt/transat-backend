package photobooth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"image"
	_ "image/png" // registers the PNG decoder for image.DecodeConfig
	"math"
	"regexp"
	"time"
)

const (
	// FrameSessionTTL is how long the QR code shown by the booth accepts uploads.
	FrameSessionTTL = 2 * time.Hour
	// MaxFrameBytes stays under the API's 15 MB request limit, multipart envelope included.
	MaxFrameBytes = 14 << 20

	maxFrameLayouts = 30
	maxFrameSlots   = 16
	maxCanvasPx     = 6000
	// A frame may differ from the canvas size (the booth scales it) but not from its shape.
	maxAspectDrift = 0.02
)

var (
	layoutTypeRe = regexp.MustCompile(`^[A-Z0-9_]{1,40}$`)

	errFrameNotPNG   = invalid("the file must be a PNG image")
	errFrameTooLarge = invalid("the image is larger than %d MB", MaxFrameBytes>>20)
)

// FrameZone is a rectangle of a layout in canvas pixels; Shot numbers the photo shown in it (1 = first shot).
type FrameZone struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	W    int `json:"w"`
	H    int `json:"h"`
	Shot int `json:"shot,omitempty"`
}

// FrameLayout is one layout of the event as the booth prints it, enough for the page to draw an exact template.
type FrameLayout struct {
	Type      string      `json:"type"`
	Title     string      `json:"title"`
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	TitleZone *FrameZone  `json:"title_zone,omitempty"`
	Photos    []FrameZone `json:"photos"`
	// HasFrame tells the page the booth already has a frame for this layout.
	HasFrame bool `json:"has_frame,omitempty"`
}

// FrameSession is what the booth gets back: the link to show as a QR code.
type FrameSession struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// FrameUpload describes a received frame, without its bytes.
type FrameUpload struct {
	LayoutType string    `json:"layout_type"`
	SHA256     string    `json:"sha256"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Size       int       `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// PublicFrameSession is what the upload page shows.
type PublicFrameSession struct {
	EventName string                 `json:"event_name"`
	ExpiresAt time.Time              `json:"expires_at"`
	Layouts   []FrameLayout          `json:"layouts"`
	Uploaded  map[string]FrameUpload `json:"uploaded"`
}

func validateFrameLayouts(layouts []FrameLayout) error {
	if len(layouts) == 0 || len(layouts) > maxFrameLayouts {
		return invalid("between 1 and %d layouts", maxFrameLayouts)
	}
	seen := map[string]bool{}
	inside := func(z FrameZone, l FrameLayout) bool {
		return z.X >= 0 && z.Y >= 0 && z.W > 0 && z.H > 0 && z.X+z.W <= l.Width && z.Y+z.H <= l.Height
	}
	for _, l := range layouts {
		if !layoutTypeRe.MatchString(l.Type) || seen[l.Type] {
			return invalid("invalid or repeated layout type")
		}
		seen[l.Type] = true
		if CleanEventName(l.Title) == "" || len([]rune(l.Title)) > 60 {
			return invalid("invalid layout title")
		}
		if l.Width < 100 || l.Height < 100 || l.Width > maxCanvasPx || l.Height > maxCanvasPx {
			return invalid("invalid canvas size")
		}
		if len(l.Photos) == 0 || len(l.Photos) > maxFrameSlots {
			return invalid("between 1 and %d photo zones per layout", maxFrameSlots)
		}
		for _, z := range l.Photos {
			if !inside(z, l) || z.Shot < 1 || z.Shot > maxFrameSlots {
				return invalid("photo zone outside the canvas")
			}
		}
		if l.TitleZone != nil && !inside(*l.TitleZone, l) {
			return invalid("title zone outside the canvas")
		}
	}
	return nil
}

// CreateFrameSession opens an upload link for the booth's layouts. Older links of the same booth stay valid until
// they expire.
func (s *Service) CreateFrameSession(ctx context.Context, device, eventName string, layouts []FrameLayout) (FrameSession, error) {
	name := CleanEventName(eventName)
	if name == "" {
		return FrameSession{}, invalid("event name is required")
	}
	for i := range layouts {
		layouts[i].Title = CleanEventName(layouts[i].Title)
	}
	if err := validateFrameLayouts(layouts); err != nil {
		return FrameSession{}, err
	}
	raw, err := json.Marshal(layouts)
	if err != nil {
		return FrameSession{}, err
	}
	token, err := newToken()
	if err != nil {
		return FrameSession{}, err
	}
	expires := s.now().Add(FrameSessionTTL)
	var id int64
	if err := s.db.QueryRowContext(ctx, `INSERT INTO photobooth_frame_sessions (device_key_id, token, event_name, layouts, expires_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, device, token, name, raw, expires).Scan(&id); err != nil {
		return FrameSession{}, err
	}
	return FrameSession{ID: id, URL: s.shareBaseURL + "/frames/" + token, ExpiresAt: expires}, nil
}

// FrameUploads lists what was sent to one of the booth's sessions; the booth polls it.
func (s *Service) FrameUploads(ctx context.Context, device string, sessionID int64) ([]FrameUpload, time.Time, error) {
	var expires time.Time
	err := s.db.QueryRowContext(ctx, `SELECT expires_at FROM photobooth_frame_sessions WHERE id = $1 AND device_key_id = $2`,
		sessionID, device).Scan(&expires)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, time.Time{}, errNotFound
	}
	if err != nil {
		return nil, time.Time{}, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT layout_type, sha256, width, height, length(png), uploaded_at
		FROM photobooth_frame_uploads WHERE session_id = $1 ORDER BY uploaded_at`, sessionID)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer rows.Close()
	out := []FrameUpload{}
	for rows.Next() {
		var u FrameUpload
		if err := rows.Scan(&u.LayoutType, &u.SHA256, &u.Width, &u.Height, &u.Size, &u.UploadedAt); err != nil {
			return nil, time.Time{}, err
		}
		out = append(out, u)
	}
	return out, expires, rows.Err()
}

// FrameImage returns the PNG received for one layout of the booth's session.
func (s *Service) FrameImage(ctx context.Context, device string, sessionID int64, layoutType string) ([]byte, error) {
	var png []byte
	err := s.db.QueryRowContext(ctx, `SELECT u.png FROM photobooth_frame_uploads u
		JOIN photobooth_frame_sessions f ON f.id = u.session_id
		WHERE u.session_id = $1 AND f.device_key_id = $2 AND u.layout_type = $3`, sessionID, device, layoutType).Scan(&png)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}
	return png, err
}

// CloseFrameSession ends a link before it expires (the admin closed the QR code). Idempotent.
func (s *Service) CloseFrameSession(ctx context.Context, device string, sessionID int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM photobooth_frame_sessions WHERE id = $1 AND device_key_id = $2`, sessionID, device)
	return err
}

func (s *Service) openFrameSession(ctx context.Context, token string) (int64, PublicFrameSession, error) {
	var out PublicFrameSession
	if !shareTokenRe.MatchString(token) {
		return 0, out, errNotFound
	}
	var id int64
	var raw []byte
	err := s.db.QueryRowContext(ctx, `SELECT id, event_name, layouts, expires_at FROM photobooth_frame_sessions
		WHERE token = $1 AND expires_at > $2`, token, s.now()).Scan(&id, &out.EventName, &raw, &out.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, out, errNotFound
	}
	if err != nil {
		return 0, out, err
	}
	if err := json.Unmarshal(raw, &out.Layouts); err != nil {
		return 0, out, err
	}
	return id, out, nil
}

// PublicFrameSession is the upload page's view of a live link.
func (s *Service) PublicFrameSession(ctx context.Context, token string) (PublicFrameSession, error) {
	id, out, err := s.openFrameSession(ctx, token)
	if err != nil {
		return out, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT layout_type, sha256, width, height, length(png), uploaded_at
		FROM photobooth_frame_uploads WHERE session_id = $1`, id)
	if err != nil {
		return out, err
	}
	defer rows.Close()
	out.Uploaded = map[string]FrameUpload{}
	for rows.Next() {
		var u FrameUpload
		if err := rows.Scan(&u.LayoutType, &u.SHA256, &u.Width, &u.Height, &u.Size, &u.UploadedAt); err != nil {
			return out, err
		}
		out.Uploaded[u.LayoutType] = u
	}
	return out, rows.Err()
}

// UploadFrame stores the frame of one layout: a PNG with the layout's shape (any size; the booth scales it).
func (s *Service) UploadFrame(ctx context.Context, token, layoutType string, data []byte) (FrameUpload, error) {
	id, session, err := s.openFrameSession(ctx, token)
	if err != nil {
		return FrameUpload{}, err
	}
	var layout *FrameLayout
	for i := range session.Layouts {
		if session.Layouts[i].Type == layoutType {
			layout = &session.Layouts[i]
		}
	}
	if layout == nil {
		return FrameUpload{}, errNotFound
	}
	if len(data) > MaxFrameBytes {
		return FrameUpload{}, errFrameTooLarge
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || format != "png" {
		return FrameUpload{}, errFrameNotPNG
	}
	want := float64(layout.Width) / float64(layout.Height)
	got := float64(cfg.Width) / float64(cfg.Height)
	if math.Abs(got-want)/want > maxAspectDrift {
		return FrameUpload{}, invalid("this layout needs a %d×%d px image (or the same proportions); this one is %d×%d",
			layout.Width, layout.Height, cfg.Width, cfg.Height)
	}
	if cfg.Width < layout.Width/2 {
		return FrameUpload{}, invalid("the image is too small: export it at %d×%d px", layout.Width, layout.Height)
	}
	if cfg.Width > maxCanvasPx || cfg.Height > maxCanvasPx {
		return FrameUpload{}, invalid("the image is too large: export it at %d×%d px", layout.Width, layout.Height)
	}
	sum := sha256.Sum256(data)
	u := FrameUpload{LayoutType: layoutType, SHA256: hex.EncodeToString(sum[:]), Width: cfg.Width, Height: cfg.Height, Size: len(data)}
	err = s.db.QueryRowContext(ctx, `INSERT INTO photobooth_frame_uploads (session_id, layout_type, png, sha256, width, height)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (session_id, layout_type) DO UPDATE
		SET png = EXCLUDED.png, sha256 = EXCLUDED.sha256, width = EXCLUDED.width, height = EXCLUDED.height, uploaded_at = now()
		RETURNING uploaded_at`, id, layoutType, data, u.SHA256, cfg.Width, cfg.Height).Scan(&u.UploadedAt)
	return u, err
}

// PurgeFrameSessions removes links (and their images) an hour after they expired, the booth having had time to
// download the last uploads.
func (s *Service) PurgeFrameSessions(ctx context.Context) (int64, error) {
	res, err := s.db.ExecContext(ctx, `DELETE FROM photobooth_frame_sessions WHERE expires_at < $1`, s.now().Add(-time.Hour))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
