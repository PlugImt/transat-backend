package profile

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type ProfileHandler struct {
	DB *sql.DB
}

func NewProfileHandler(db *sql.DB) *ProfileHandler {
	return &ProfileHandler{DB: db}
}

var errProfileNotFound = errors.New("profile not found")

// Shared by both lookups: only student accounts have public profiles, like the other student features.
const profileSelect = `
	SELECT
		n.public_id::text,
		n.id_newf,
		n.email,
		n.first_name,
		n.last_name,
		NULLIF(n.profile_picture, ''),
		NULLIF(n.formation_name, ''),
		n.graduation_year,
		NULLIF(n.campus, ''),
		TO_CHAR(n.creation_date AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
		COALESCE(p.bio, ''),
		COALESCE(p.emoji, ''),
		COALESCE(p.interests, '{}')
	FROM newf n
	LEFT JOIN user_profiles p ON p.email = n.email
	WHERE EXISTS (
		SELECT 1 FROM newf_roles nr
		JOIN roles r ON nr.id_roles = r.id_roles
		WHERE nr.email = n.email AND r.name IN ('NEWF', 'ADMIN')
	)
`

// load builds a profile from a row filter ("n.public_id = $1::uuid" or "n.email = $1").
func (h *ProfileHandler) load(filter string, arg string, viewerEmail string) (*models.PublicProfile, error) {
	var (
		profile                    models.PublicProfile
		email                      string
		memberNumber               int
		picture, formation, campus sql.NullString
		graduationYear             sql.NullInt64
		interests                  []string
	)

	err := h.DB.QueryRow(profileSelect+" AND "+filter, arg).Scan(
		&profile.ID, &memberNumber, &email, &profile.FirstName, &profile.LastName,
		&picture, &formation, &graduationYear, &campus, &profile.JoinedAt,
		&profile.Bio, &profile.Emoji, pq.Array(&interests),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errProfileNotFound
	}
	if err != nil {
		return nil, err
	}

	profile.MemberNumber = memberNumber
	profile.IsMe = strings.EqualFold(email, viewerEmail)
	profile.ProfilePicture = nullable(picture)
	profile.FormationName = nullable(formation)
	profile.Campus = nullable(campus)
	if graduationYear.Valid {
		year := int(graduationYear.Int64)
		profile.GraduationYear = &year
	}
	profile.Interests = interests
	if profile.Interests == nil {
		profile.Interests = []string{}
	}

	if err := h.loadActivity(&profile, email); err != nil {
		return nil, err
	}
	profile.Badges = badgesFor(memberNumber, profile)
	return &profile, nil
}

// loadActivity fills data that is already visible elsewhere in the app (club and event member lists,
// reviews, the bassine leaderboard), only aggregated.
func (h *ProfileHandler) loadActivity(profile *models.PublicProfile, email string) error {
	profile.Clubs = []models.ProfileClub{}

	rows, err := h.DB.Query(`
		SELECT c.id_clubs, c.name, COALESCE(c.picture, ''), EXISTS (
			SELECT 1 FROM newf_roles nr
			JOIN roles r ON nr.id_roles = r.id_roles
			WHERE nr.email = cm.email AND r.name = LOWER(REPLACE(c.name, ' ', '')) || '_respo'
		) AS is_respo
		FROM clubs_members cm
		JOIN clubs c ON c.id_clubs = cm.id_clubs
		WHERE cm.email = $1
		ORDER BY is_respo DESC, c.name
	`, email)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var club models.ProfileClub
		if err := rows.Scan(&club.ID, &club.Name, &club.Picture, &club.IsRespo); err != nil {
			return err
		}
		profile.Clubs = append(profile.Clubs, club)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if err := h.DB.QueryRow(`SELECT COUNT(*) FROM restaurant_articles_notes WHERE email = $1`, email).Scan(&profile.Stats.Reviews); err != nil {
		return err
	}
	if err := h.DB.QueryRow(`SELECT COUNT(*) FROM events_attendents WHERE email = $1`, email).Scan(&profile.Stats.Events); err != nil {
		return err
	}

	var bassine models.ProfileBassine
	err = h.DB.QueryRow(`
		SELECT score, rank FROM (
			SELECT email, score, RANK() OVER (ORDER BY score DESC, email ASC) AS rank FROM bassine_scores
		) ranked WHERE email = $1
	`, email).Scan(&bassine.Score, &bassine.Rank)
	if err == nil {
		profile.Stats.Bassine = &bassine
	} else if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}

func badgesFor(memberNumber int, profile models.PublicProfile) []string {
	badges := []string{}
	if memberNumber <= models.PioneerMemberLimit {
		badges = append(badges, models.BadgePioneer)
	}
	for _, club := range profile.Clubs {
		if club.IsRespo {
			badges = append(badges, models.BadgeClubLeader)
			break
		}
	}
	if profile.Stats.Reviews >= models.CriticReviews {
		badges = append(badges, models.BadgeCritic)
	}
	return badges
}

func nullable(s sql.NullString) *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// GetMyProfile returns the logged-in user's own profile, as others see it.
func (h *ProfileHandler) GetMyProfile(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	return h.respond(c, "n.email = $1", email, email)
}

// GetProfile returns the public profile of the user with the given public id.
func (h *ProfileHandler) GetProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	if !models.IsPublicID(id) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return h.respond(c, "n.public_id = $1::uuid", id, c.Locals("email").(string))
}

func (h *ProfileHandler) respond(c *fiber.Ctx, filter, arg, viewerEmail string) error {
	profile, err := h.load(filter, arg, viewerEmail)
	if errors.Is(err, errProfileNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to load profile")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load profile"})
	}
	return c.JSON(profile)
}

// UpdateMyProfile changes the customizable part of the logged-in user's profile.
func (h *ProfileHandler) UpdateMyProfile(c *fiber.Ctx) error {
	email := c.Locals("email").(string)

	var req models.ProfileUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	if err := req.Normalize(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Absent fields keep their stored value.
	var interests interface{}
	if req.Interests != nil {
		interests = pq.Array(*req.Interests)
	}
	_, err := h.DB.Exec(`
		INSERT INTO user_profiles (email, bio, emoji, interests)
		VALUES ($1, COALESCE($2::text, ''), COALESCE($3::text, ''), COALESCE($4::text[], '{}'))
		ON CONFLICT (email) DO UPDATE SET
			bio = COALESCE($2::text, user_profiles.bio),
			emoji = COALESCE($3::text, user_profiles.emoji),
			interests = COALESCE($4::text[], user_profiles.interests),
			updated_at = CURRENT_TIMESTAMP
	`, email, req.Bio, req.Emoji, interests)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to save profile")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save profile"})
	}

	return h.respond(c, "n.email = $1", email, email)
}
