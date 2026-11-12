package faces

import (
	"encoding/binary"
	"errors"
	"math"
)

// Embedding is a face descriptor from the recognition model (512 values for buffalo_l).
type Embedding []float32

var errBadEmbedding = errors.New("invalid face embedding")

// normalize returns a unit-length copy, so that comparing two faces is a dot product.
func normalize(e Embedding) (Embedding, error) {
	var sum float64
	for _, v := range e {
		if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			return nil, errBadEmbedding
		}
		sum += float64(v) * float64(v)
	}
	if len(e) == 0 || sum == 0 {
		return nil, errBadEmbedding
	}
	n := math.Sqrt(sum)
	out := make(Embedding, len(e))
	for i, v := range e {
		out[i] = float32(float64(v) / n)
	}
	return out, nil
}

// distance is the cosine distance between two unit embeddings: 0 for the same face, up to 2.
// Embeddings of different lengths (another model) never match.
func distance(a, b Embedding) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 2
	}
	var dot float32
	for i := range a {
		dot += a[i] * b[i]
	}
	return 1 - dot
}

// encode stores an embedding as little-endian float32, 2 KiB for 512 values.
func encode(e Embedding) []byte {
	b := make([]byte, 4*len(e))
	for i, v := range e {
		binary.LittleEndian.PutUint32(b[4*i:], math.Float32bits(v))
	}
	return b
}

func decode(b []byte) (Embedding, error) {
	if len(b) == 0 || len(b)%4 != 0 {
		return nil, errBadEmbedding
	}
	e := make(Embedding, len(b)/4)
	for i := range e {
		e[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[4*i:]))
	}
	return e, nil
}
