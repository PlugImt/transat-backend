package photos

import "testing"

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Soirée de Gala 2026":            "soiree-de-gala-2026",
		"  DEST'IMT — Apocalypse Tonus ": "dest-imt-apocalypse-tonus",
		"Œuvre ß 日本":                     "uvre",
		"!!!":                            "",
		"Ça c'est l'été":                 "ca-c-est-l-ete",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
	if got := Slug("a very long title that keeps going and going and going well past the sixty characters"); len(got) > maxSlugLen || got[len(got)-1] == '-' {
		t.Errorf("long slug = %q", got)
	}
}

func TestAlbumShareURL(t *testing.T) {
	t.Setenv("PUBLIC_SITE_URL", "")
	if got := AlbumShareURL(42, "Soirée de Gala"); got != "https://transat.dev/albums/42-soiree-de-gala" {
		t.Errorf("default = %q", got)
	}
	t.Setenv("PUBLIC_SITE_URL", "https://develop.transat.dev/")
	if got := AlbumShareURL(7, "!!!"); got != "https://develop.transat.dev/albums/7" {
		t.Errorf("override = %q", got)
	}
}
