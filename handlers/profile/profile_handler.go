package profile

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/services"
	"github.com/plugimt/transat-backend/utils"
)

type ProfileHandler struct {
	DB        *sql.DB
	R2Service *services.R2Service
}

func NewProfileHandler(db *sql.DB, r2Service *services.R2Service) *ProfileHandler {
	return &ProfileHandler{DB: db, R2Service: r2Service}
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
		COALESCE(p.interests, '{}'),
		COALESCE(p.decoration_image_ids, '{}'),
		COALESCE(p.decoration_emojis, '{}')
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
		decorationImageIDs         []int64
		decorationEmojis           []string
	)

	err := h.DB.QueryRow(profileSelect+" AND "+filter, arg).Scan(
		&profile.ID, &memberNumber, &email, &profile.FirstName, &profile.LastName,
		&picture, &formation, &graduationYear, &campus, &profile.JoinedAt,
		&profile.Bio, &profile.Emoji, pq.Array(&interests),
		pq.Array(&decorationImageIDs), pq.Array(&decorationEmojis),
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
	profile.DecorationEmojis = decorationEmojis
	if profile.DecorationEmojis == nil {
		profile.DecorationEmojis = []string{}
	}
	if err := h.loadDecorationImages(&profile, email, decorationImageIDs); err != nil {
		return nil, err
	}

	if err := h.loadActivity(&profile, email); err != nil {
		return nil, err
	}
	profile.Badges = badgesFor(memberNumber, profile)
	return &profile, nil
}

func (h *ProfileHandler) loadDecorationImages(profile *models.PublicProfile, email string, ids []int64) error {
	profile.DecorationImages = []models.ProfileDecorationImage{}
	if len(ids) == 0 {
		return nil
	}

	rows, err := h.DB.Query(`
		SELECT selected.file_id, f.path
		FROM unnest($1::integer[]) WITH ORDINALITY AS selected(file_id, position)
		JOIN files f ON f.id_files = selected.file_id AND f.email = $2
			AND LOWER(f.path) ~ '[.](jpg|jpeg|png|webp)$'
		ORDER BY selected.position
	`, pq.Array(ids), email)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var path string
		if err := rows.Scan(&id, &path); err != nil {
			return err
		}
		profile.DecorationImages = append(profile.DecorationImages, models.ProfileDecorationImage{
			ID:  id,
			URL: h.R2Service.GetPublicURL(path),
		})
	}
	return rows.Err()
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
	var decorationImageIDs interface{}
	if req.DecorationImageIDs != nil {
		if len(*req.DecorationImageIDs) > 0 {
			var owned int
			if err := h.DB.QueryRow(`
				SELECT COUNT(*) FROM files
				WHERE email = $1 AND id_files = ANY($2::integer[])
					AND LOWER(path) ~ '[.](jpg|jpeg|png|webp)$'
			`, email, pq.Array(*req.DecorationImageIDs)).Scan(&owned); err != nil {
				utils.LogMessage(utils.LevelError, "Failed to verify profile decoration images")
				utils.LogLineKeyValue(utils.LevelError, "Error", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save profile"})
			}
			if owned != len(*req.DecorationImageIDs) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Decoration images must belong to your account"})
			}
		}
		decorationImageIDs = pq.Array(*req.DecorationImageIDs)
	}
	var decorationEmojis interface{}
	if req.DecorationEmojis != nil {
		decorationEmojis = pq.Array(*req.DecorationEmojis)
	}
	_, err := h.DB.Exec(`
		INSERT INTO user_profiles (email, bio, emoji, interests, decoration_image_ids, decoration_emojis)
		VALUES ($1, COALESCE($2::text, ''), COALESCE($3::text, ''), COALESCE($4::text[], '{}'),
			COALESCE($5::integer[], '{}'), COALESCE($6::text[], '{}'))
		ON CONFLICT (email) DO UPDATE SET
			bio = COALESCE($2::text, user_profiles.bio),
			emoji = COALESCE($3::text, user_profiles.emoji),
			interests = COALESCE($4::text[], user_profiles.interests),
			decoration_image_ids = COALESCE($5::integer[], user_profiles.decoration_image_ids),
			decoration_emojis = COALESCE($6::text[], user_profiles.decoration_emojis),
			updated_at = CURRENT_TIMESTAMP
	`, email, req.Bio, req.Emoji, interests, decorationImageIDs, decorationEmojis)
	if err != nil {
		utils.LogMessage(utils.LevelError, "Failed to save profile")
		utils.LogLineKeyValue(utils.LevelError, "Error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save profile"})
	}

	return h.respond(c, "n.email = $1", email, email)
}
