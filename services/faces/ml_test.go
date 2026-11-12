package faces

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmbeddingRoundTripAndDistance(t *testing.T) {
	a, err := normalize(Embedding{3, 4})
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(float64(a[0])-0.6) > 1e-6 || math.Abs(float64(a[1])-0.8) > 1e-6 {
		t.Fatalf("normalize: %v", a)
	}
	back, err := decode(encode(a))
	if err != nil || back[0] != a[0] || back[1] != a[1] {
		t.Fatalf("round trip: %v %v", back, err)
	}
	if d := distance(a, a); math.Abs(float64(d)) > 1e-6 {
		t.Fatalf("same face distance %v", d)
	}
	b, _ := normalize(Embedding{-3, -4})
	if d := distance(a, b); math.Abs(float64(d)-2) > 1e-6 {
		t.Fatalf("opposite distance %v", d)
	}
	if d := distance(a, Embedding{1, 0, 0}); d != 2 {
		t.Fatalf("other model must never match, got %v", d)
	}
	if _, err := normalize(Embedding{0, 0}); err == nil {
		t.Fatal("zero vector accepted")
	}
	if _, err := decode([]byte{1, 2, 3}); err == nil {
		t.Fatal("truncated embedding accepted")
	}
}

func TestMLClientSpeaksImmichPredict(t *testing.T) {
	var gotEntries map[string]any
	var gotImage []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/predict" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		if err := json.Unmarshal([]byte(r.FormValue("entries")), &gotEntries); err != nil {
			t.Errorf("entries: %v", err)
		}
		f, _, err := r.FormFile("image")
		if err != nil {
			t.Errorf("image: %v", err)
			return
		}
		gotImage, _ = io.ReadAll(f)
		// Immich sends the embedding as a JSON string.
		_, _ = w.Write([]byte(`{"facial-recognition":[
			{"boundingBox":{"x1":100,"y1":50,"x2":300,"y2":250},"embedding":"[3,4]","score":0.98},
			{"boundingBox":{"x1":0,"y1":0,"x2":10,"y2":10},"embedding":[0,2],"score":0.71}],
			"imageWidth":1000,"imageHeight":500}`))
	}))
	defer srv.Close()

	c := NewMLClient(srv.URL+"/", "buffalo_l", 0.7)
	d, err := c.Detect(context.Background(), []byte("jpeg-bytes"))
	if err != nil {
		t.Fatal(err)
	}
	if string(gotImage) != "jpeg-bytes" {
		t.Fatalf("image sent: %q", gotImage)
	}
	fr := gotEntries["facial-recognition"].(map[string]any)
	det := fr["detection"].(map[string]any)
	if det["modelName"] != "buffalo_l" || fr["recognition"].(map[string]any)["modelName"] != "buffalo_l" {
		t.Fatalf("entries: %v", gotEntries)
	}
	if len(d.Faces) != 2 || d.ImageWidth != 1000 || d.ImageHeight != 500 {
		t.Fatalf("detection: %+v", d)
	}
	f := d.Faces[0]
	if f.X1 != 0.1 || f.Y1 != 0.1 || f.X2 != 0.3 || f.Y2 != 0.5 {
		t.Fatalf("relative box: %+v", f)
	}
	if math.Abs(float64(f.Embedding[0])-0.6) > 1e-6 || d.Faces[1].Embedding[1] != 1 {
		t.Fatalf("embeddings not normalised: %v %v", f.Embedding, d.Faces[1].Embedding)
	}
}

func TestMLClientErrors(t *testing.T) {
	status := http.StatusUnprocessableEntity
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"detail":"cannot identify image file"}`))
	}))
	defer srv.Close()
	c := NewMLClient(srv.URL, "buffalo_l", 0.7)
	if _, err := c.Detect(context.Background(), []byte("x")); !errors.Is(err, ErrUnreadableImage) {
		t.Fatalf("422: %v", err)
	}
	status = http.StatusInternalServerError
	if _, err := c.Detect(context.Background(), []byte("x")); err == nil || errors.Is(err, ErrUnreadableImage) {
		t.Fatalf("500 must be a retryable error: %v", err)
	}
}

func face(x1, y1, x2, y2, score float32) Face {
	return Face{X1: x1, Y1: y1, X2: x2, Y2: y2, Score: score, Embedding: Embedding{1, 0}}
}

func TestPickReference(t *testing.T) {
	det := func(faces ...Face) Detection { return Detection{ImageWidth: 1000, ImageHeight: 1000, Faces: faces} }
	cases := []struct {
		name string
		d    Detection
		want error
	}{
		{"one clear face", det(face(0.3, 0.2, 0.7, 0.7, 0.99)), nil},
		{"no face", det(), ErrNoFace},
		{"only a blurry face", det(face(0.3, 0.2, 0.7, 0.7, 0.5)), ErrFaceTooSmall},
		{"tiny face", det(face(0.5, 0.5, 0.55, 0.55, 0.99)), ErrFaceTooSmall},
		{"two similar faces", det(face(0.1, 0.1, 0.4, 0.4, 0.99), face(0.6, 0.1, 0.85, 0.4, 0.99)), ErrManyFaces},
		{"someone small in the background", det(face(0.2, 0.1, 0.8, 0.8, 0.99), face(0.85, 0.1, 0.95, 0.2, 0.9)), nil},
	}
	for _, c := range cases {
		if _, err := pickReference(c.d); !errors.Is(err, c.want) {
			t.Errorf("%s: got %v, want %v", c.name, err, c.want)
		}
	}
}
