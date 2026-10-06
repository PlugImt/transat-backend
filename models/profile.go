package models

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	MaxBioLength       = 160
	MaxEmojiRunes      = 10 // enough for ZWJ sequences such as family or profession emojis
	MaxInterests       = 6
	PioneerMemberLimit = 100 // accounts numbered up to this value are pioneers
	CriticReviews      = 10  // RU reviews needed for the "critic" badge
)

// Interests a user can pick. The app maps each id to an icon and a label.
var InterestCatalog = []string{
	"music", "sport", "gaming", "movies", "series", "reading", "cooking", "travel",
	"tech", "art", "photography", "nature", "board_games", "fashion", "volunteering", "parties",
}

// Badges are computed server-side from account data, never stored or sent by clients.
const (
	BadgePioneer    = "pioneer"
	BadgeClubLeader = "club_leader"
	BadgeCritic     = "critic"
)

var publicIDPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// IsPublicID reports whether s looks like a user public id.
func IsPublicID(s string) bool { return publicIDPattern.MatchString(s) }

type ProfileClub struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
	IsRespo bool   `json:"is_respo"`
}

type ProfileBassine struct {
	Score int `json:"score"`
	Rank  int `json:"rank"`
}

type ProfileStats struct {
	Reviews int             `json:"reviews"`
	Events  int             `json:"events"`
	Bassine *ProfileBassine `json:"bassine"`
}

// PublicProfile is everything another student may see about a user. It never contains an email,
// phone number, language, roles or any authentication data.
type PublicProfile struct {
	ID             string        `json:"id"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	ProfilePicture *string       `json:"profile_picture"`
	FormationName  *string       `json:"formation_name"`
	GraduationYear *int          `json:"graduation_year"`
	Campus         *string       `json:"campus"`
	JoinedAt       string        `json:"joined_at"`
	MemberNumber   int           `json:"member_number"`
	Badges         []string      `json:"badges"`
	Bio            string        `json:"bio"`
	Emoji          string        `json:"emoji"`
	Interests      []string      `json:"interests"`
	Clubs          []ProfileClub `json:"clubs"`
	Stats          ProfileStats  `json:"stats"`
	IsMe           bool          `json:"is_me"`
}

// ProfileUpdateRequest holds the profile fields a user may edit; nil fields are left unchanged.
type ProfileUpdateRequest struct {
	Bio       *string   `json:"bio"`
	Emoji     *string   `json:"emoji"`
	Interests *[]string `json:"interests"`
}

// Normalize trims and validates the request in place.
func (r *ProfileUpdateRequest) Normalize() error {
	if r.Bio != nil {
		bio := strings.Map(func(c rune) rune {
			if unicode.IsControl(c) && c != '\n' {
				return -1
			}
			return c
		}, strings.TrimSpace(*r.Bio))
		if utf8.RuneCountInString(bio) > MaxBioLength {
			return fmt.Errorf("bio must be at most %d characters", MaxBioLength)
		}
		r.Bio = &bio
	}

	if r.Emoji != nil {
		emoji := strings.TrimSpace(*r.Emoji)
		if !isEmoji(emoji) {
			return fmt.Errorf("emoji must be a single emoji")
		}
		r.Emoji = &emoji
	}

	if r.Interests != nil {
		if len(*r.Interests) > MaxInterests {
			return fmt.Errorf("at most %d interests", MaxInterests)
		}
		known := make(map[string]bool, len(InterestCatalog))
		for _, id := range InterestCatalog {
			known[id] = true
		}
		seen := make(map[string]bool, len(*r.Interests))
		unique := make([]string, 0, len(*r.Interests))
		for _, id := range *r.Interests {
			if !known[id] {
				return fmt.Errorf("unknown interest %q", id)
			}
			if !seen[id] {
				seen[id] = true
				unique = append(unique, id)
			}
		}
		r.Interests = &unique
	}
	return nil
}

// isEmoji accepts the empty string (clears the emoji) or a short run of symbol characters.
func isEmoji(s string) bool {
	if s == "" {
		return true
	}
	if utf8.RuneCountInString(s) > MaxEmojiRunes {
		return false
	}
	for _, c := range s {
		if c < 0x80 || unicode.IsLetter(c) || unicode.IsDigit(c) || unicode.IsSpace(c) || unicode.IsControl(c) || unicode.IsPunct(c) {
			return false
		}
	}
	return true
}
