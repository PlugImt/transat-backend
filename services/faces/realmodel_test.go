package faces

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/plugimt/transat-backend/internal/database/dbtest"
)

// End-to-end check against the real model, skipped unless both are set:
//
//	FACES_ML_TEST_URL=http://localhost:3003   (docker run -p 3003:3003 ghcr.io/immich-app/immich-machine-learning:release)
//	FACES_TEST_IMAGES=/path/to/portraits      (one person per JPEG, file name = person; never committed)
//
// Each portrait is the person's selfie; other "shots" of them are derived with macOS sips (mirrored, smaller,
// slightly rotated). The album of each person must hold their own shots and nobody else's.

type portrait struct {
	person string
	path   string
}

func realModel(t *testing.T) (*MLClient, []portrait) {
	t.Helper()
	url, dir := os.Getenv("FACES_ML_TEST_URL"), os.Getenv("FACES_TEST_IMAGES")
	if url == "" || dir == "" {
		t.Skip("FACES_ML_TEST_URL and FACES_TEST_IMAGES not set")
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.jpg"))
	if len(files) < 2 {
		t.Fatalf("need at least two portraits in %s", dir)
	}
	var out []portrait
	for _, f := range files {
		out = append(out, portrait{strings.TrimSuffix(filepath.Base(f), ".jpg"), f})
	}
	return NewMLClient(url, defaultModel, 0.7), out
}

// variants derives other shots of a portrait with sips.
func variants(t *testing.T, p portrait) map[string]string {
	t.Helper()
	if _, err := exec.LookPath("sips"); err != nil {
		t.Skip("sips (macOS) is needed to derive shots")
	}
	dir := t.TempDir()
	out := map[string]string{}
	for name, args := range map[string][]string{
		"mirror": {"-f", "horizontal"},
		"small":  {"-Z", "480"},
		"tilted": {"-r", "8"},
	} {
		dst := filepath.Join(dir, p.person+"-"+name+".jpg")
		cmd := exec.Command("sips", append(args, p.path, "--out", dst)...)
		if b, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("sips %v: %v %s", args, err, b)
		}
		out[p.person+"-"+name] = dst
	}
	return out
}

func read(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestRealModelSeparatesPeople(t *testing.T) {
	ml, people := realModel(t)
	ctx := context.Background()
	refs := map[string]Embedding{}
	for _, p := range people {
		det, err := ml.Detect(ctx, read(t, p.path))
		if err != nil {
			t.Fatalf("%s: %v", p.person, err)
		}
		f, err := pickReference(det)
		if err != nil {
			t.Fatalf("%s is not a usable selfie (%d faces): %v", p.person, len(det.Faces), err)
		}
		if len(f.Embedding) != 512 {
			t.Fatalf("%s: embedding of %d values", p.person, len(f.Embedding))
		}
		refs[p.person] = f.Embedding
	}

	var same, other []float32
	for _, p := range people {
		for shot, path := range variants(t, p) {
			det, err := ml.Detect(ctx, read(t, path))
			if err != nil || len(det.Faces) == 0 {
				t.Fatalf("%s: %v (%d faces)", shot, err, len(det.Faces))
			}
			face := det.Faces[0]
			for person, ref := range refs {
				d := distance(ref, face.Embedding)
				if person == p.person {
					same = append(same, d)
					t.Logf("%-16s vs %-8s %.3f (same)", shot, person, d)
				} else {
					other = append(other, d)
				}
			}
		}
	}
	sort.Slice(same, func(i, j int) bool { return same[i] < same[j] })
	sort.Slice(other, func(i, j int) bool { return other[i] < other[j] })
	t.Logf("same person: max %.3f; different people: min %.3f; threshold %.2f",
		same[len(same)-1], other[0], DefaultConfig().MaxDistance)
	if same[len(same)-1] > DefaultConfig().MaxDistance {
		t.Errorf("a shot of the same person is above the threshold: %.3f", same[len(same)-1])
	}
	if other[0] <= DefaultConfig().MaxDistance {
		t.Errorf("two different people are under the threshold: %.3f", other[0])
	}
}

// files serves photos from local paths, as Drive would.
type files map[string]string

func (f files) Image(_ context.Context, id string, _ int) ([]byte, error) {
	p, ok := f[id]
	if !ok {
		return nil, ErrImageUnavailable
	}
	return os.ReadFile(p)
}

func TestRealModelBuildsAlbums(t *testing.T) {
	ml, people := realModel(t)
	db := dbtest.Open(t)
	ctx := context.Background()
	photos := files{}
	want := map[string][]string{}
	minute := 0
	for _, p := range people {
		for shot, path := range variants(t, p) {
			id := fmt.Sprintf("PHOTO-%s", shot)
			photos[id] = path
			want[p.person] = append(want[p.person], id)
			minute++
			if _, err := db.Exec(`INSERT INTO photos (drive_file_id, name, mime_type, media_type, md5, sort_at, sort_at_source, web_url, last_seen_run)
				VALUES ($1, $1, 'image/jpeg', 'image', $1, now() - make_interval(mins => $2), 'x', 'u', 1)`, id, minute); err != nil {
				t.Fatal(err)
			}
		}
	}
	cfg := DefaultConfig()
	cfg.TickBudget = 0
	svc := NewService(db, ml, photos, cfg)

	// Half of the people opt in; the others never do.
	for i, p := range people {
		email := p.person + "@imt.fr"
		if _, err := db.Exec(`INSERT INTO newf (email, password, first_name, last_name) VALUES ($1, 'x', 'Jo', 'Doe')`, email); err != nil {
			t.Fatal(err)
		}
		if i%2 == 1 {
			continue
		}
		if err := svc.Consent(ctx, email, ConsentVersion); err != nil {
			t.Fatal(err)
		}
		if err := svc.AddReference(ctx, email, read(t, p.path)); err != nil {
			t.Fatalf("%s selfie: %v", p.person, err)
		}
	}
	st, err := svc.Tick(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("tick: %+v", st)

	for i, p := range people {
		var got []string
		rows, err := db.Query(`SELECT DISTINCT drive_file_id FROM face_matches WHERE email = $1 ORDER BY 1`, p.person+"@imt.fr")
		if err != nil {
			t.Fatal(err)
		}
		for rows.Next() {
			var id string
			_ = rows.Scan(&id)
			got = append(got, id)
		}
		rows.Close()
		expected := []string{}
		if i%2 == 0 {
			expected = append(expected, want[p.person]...)
			sort.Strings(expected)
		}
		if strings.Join(got, ",") != strings.Join(expected, ",") {
			t.Errorf("%s album: got %v, want %v", p.person, got, expected)
		}
	}
}
