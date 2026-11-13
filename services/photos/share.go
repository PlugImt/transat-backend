package photos

import (
	"os"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// DefaultSiteURL is the public website that hosts the album pages; PUBLIC_SITE_URL overrides it (the develop
// API points at its own site, for instance).
const DefaultSiteURL = "https://transat.dev"

const maxSlugLen = 60

// SiteURL is the public website's origin, without a trailing slash.
func SiteURL() string {
	if v := strings.TrimSpace(os.Getenv("PUBLIC_SITE_URL")); v != "" {
		return strings.TrimRight(v, "/")
	}
	return DefaultSiteURL
}

// AlbumShareURL is the link given out for an album: its page on transat.dev, which also opens the app on the
// album when it is installed. The id identifies the album; the slug is only there to read well and may go stale
// when the album is renamed, which the site tolerates.
func AlbumShareURL(id int64, title string) string {
	return SiteURL() + "/albums/" + AlbumPathSegment(id, title)
}

// AlbumPathSegment is "<id>-<slug>", or "<id>" when the title has nothing to slugify.
func AlbumPathSegment(id int64, title string) string {
	seg := strconv.FormatInt(id, 10)
	if s := Slug(title); s != "" {
		seg += "-" + s
	}
	return seg
}

// Slug turns a title into lowercase ASCII words joined by dashes: "Soirée de Gala 2026" → "soiree-de-gala-2026".
func Slug(title string) string {
	var b strings.Builder
	dash := false
	for _, r := range norm.NFD.String(title) {
		switch {
		case unicode.Is(unicode.Mn, r):
			continue // accents, once decomposed
		case r < unicode.MaxASCII && (unicode.IsLetter(r) || unicode.IsDigit(r)):
			if dash && b.Len() > 0 {
				b.WriteByte('-')
			}
			dash = false
			b.WriteRune(unicode.ToLower(r))
		default:
			dash = true
		}
		if b.Len() > maxSlugLen {
			break
		}
	}
	slug := b.String()
	if len(slug) > maxSlugLen {
		slug = slug[:maxSlugLen]
		// Cut at a word boundary when one is not too far back.
		if i := strings.LastIndexByte(slug, '-'); i > maxSlugLen/2 {
			slug = slug[:i]
		}
	}
	return strings.Trim(slug, "-")
}
