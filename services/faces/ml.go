package faces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Face is one face found on an image, its box relative to the image (0..1).
type Face struct {
	X1, Y1, X2, Y2 float32
	Score          float32
	Embedding      Embedding // unit length
}

// Width and Height of the box, relative to the image.
func (f Face) Width() float32  { return f.X2 - f.X1 }
func (f Face) Height() float32 { return f.Y2 - f.Y1 }

// Detection is what the model returns for one image.
type Detection struct {
	ImageWidth, ImageHeight int
	Faces                   []Face
}

// Recognizer finds faces and computes their embeddings. MLClient talks to Immich's machine-learning container.
type Recognizer interface {
	Detect(ctx context.Context, image []byte) (Detection, error)
	Model() string
}

// ErrUnreadableImage means the model could not decode the image: retrying will not help.
var ErrUnreadableImage = errors.New("image could not be decoded")

// MLClient calls POST /predict of ghcr.io/immich-app/immich-machine-learning, the facial-recognition task only.
// Nothing is kept by the container: it decodes the image, runs the models and answers.
type MLClient struct {
	base     string
	model    string
	minScore float32
	http     *http.Client
}

var _ Recognizer = (*MLClient)(nil)

func NewMLClient(baseURL, model string, minScore float32) *MLClient {
	return &MLClient{
		base:     strings.TrimRight(baseURL, "/"),
		model:    model,
		minScore: minScore,
		// The first call after an idle period loads the models (tens of seconds on CPU).
		http: &http.Client{Timeout: 3 * time.Minute},
	}
}

func (c *MLClient) Model() string { return c.model }

type mlEntries map[string]map[string]any

func (c *MLClient) entries() ([]byte, error) {
	return json.Marshal(mlEntries{
		"facial-recognition": {
			"detection":   map[string]any{"modelName": c.model, "options": map[string]any{"minScore": c.minScore}},
			"recognition": map[string]any{"modelName": c.model},
		},
	})
}

func (c *MLClient) Detect(ctx context.Context, image []byte) (Detection, error) {
	entries, err := c.entries()
	if err != nil {
		return Detection{}, err
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("entries", string(entries)); err != nil {
		return Detection{}, err
	}
	part, err := w.CreateFormFile("image", "image.jpg")
	if err != nil {
		return Detection{}, err
	}
	if _, err := part.Write(image); err != nil {
		return Detection{}, err
	}
	if err := w.Close(); err != nil {
		return Detection{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+"/predict", &body)
	if err != nil {
		return Detection{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := c.http.Do(req)
	if err != nil {
		return Detection{}, fmt.Errorf("face model unreachable: %w", err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 16<<20))
	if err != nil {
		return Detection{}, err
	}
	switch {
	case res.StatusCode == http.StatusUnprocessableEntity || res.StatusCode == http.StatusBadRequest:
		return Detection{}, fmt.Errorf("%w: %s", ErrUnreadableImage, snippet(raw))
	case res.StatusCode != http.StatusOK:
		return Detection{}, fmt.Errorf("face model answered %d: %s", res.StatusCode, snippet(raw))
	}
	return parsePrediction(raw)
}

func snippet(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}

type mlResponse struct {
	Faces []struct {
		BoundingBox struct {
			X1, Y1, X2, Y2 float64
		} `json:"boundingBox"`
		// Immich serialises the vector as a JSON string ("[0.1,...]"); a plain array is accepted too.
		Embedding json.RawMessage `json:"embedding"`
		Score     float64         `json:"score"`
	} `json:"facial-recognition"`
	ImageWidth  int `json:"imageWidth"`
	ImageHeight int `json:"imageHeight"`
}

func parsePrediction(raw []byte) (Detection, error) {
	var r mlResponse
	if err := json.Unmarshal(raw, &r); err != nil {
		return Detection{}, fmt.Errorf("face model response: %w", err)
	}
	if r.ImageWidth <= 0 || r.ImageHeight <= 0 {
		return Detection{}, errors.New("face model response without image size")
	}
	d := Detection{ImageWidth: r.ImageWidth, ImageHeight: r.ImageHeight}
	w, h := float64(r.ImageWidth), float64(r.ImageHeight)
	for _, f := range r.Faces {
		vec, err := parseVector(f.Embedding)
		if err != nil {
			return Detection{}, err
		}
		unit, err := normalize(vec)
		if err != nil {
			return Detection{}, err
		}
		d.Faces = append(d.Faces, Face{
			X1: clamp01(f.BoundingBox.X1 / w), Y1: clamp01(f.BoundingBox.Y1 / h),
			X2: clamp01(f.BoundingBox.X2 / w), Y2: clamp01(f.BoundingBox.Y2 / h),
			Score: float32(f.Score), Embedding: unit,
		})
	}
	return d, nil
}

func parseVector(raw json.RawMessage) (Embedding, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) > 0 && raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return nil, errBadEmbedding
		}
		raw = []byte(s)
	}
	var v Embedding
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, errBadEmbedding
	}
	return v, nil
}

func clamp01(v float64) float32 {
	return float32(min(1, max(0, v)))
}
