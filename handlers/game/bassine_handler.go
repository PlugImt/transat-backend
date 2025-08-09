package game

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type BassineHandler struct {
	db *sql.DB
}

func NewBassineHandler(db *sql.DB) *BassineHandler {
	return &BassineHandler{db: db}
}

// checkUserExists verifies if a user exists in the database
func (h *BassineHandler) checkUserExists(email string) error {
	var exists bool
	if err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM newf WHERE email = $1)", email).Scan(&exists); err != nil {
		utils.LogMessage(utils.LevelError, "Failed to check user existence")
		return err
	}
	if !exists {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	return nil
}

// SetUserScore handles PUT /game/bassine/:email with body {"score": number}
func (h *BassineHandler) SetUserScore(c *fiber.Ctx) error {
    utils.LogHeader("🎮 Bassine Set Score")

    email := c.Params("email")
    if email == "" {
        utils.LogMessage(utils.LevelWarn, "Email is required")
        utils.LogFooter()
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
    }

    	// Parse score from JSON body only
	var req models.BassineScoreRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelWarn, "Invalid request body")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

    if req.Score < 0 {
        utils.LogMessage(utils.LevelWarn, "Score cannot be negative")
        utils.LogFooter()
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Score cannot be negative"})
    }

    	// Ensure user exists
	if err := h.checkUserExists(email); err != nil {
		utils.LogMessage(utils.LevelWarn, "User not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

    // Upsert score
    _, err := h.db.Exec(`
        INSERT INTO bassine_scores (email, score)
        VALUES ($1, $2)
        ON CONFLICT (email) DO UPDATE SET score = EXCLUDED.score
    `, email, req.Score)
    if err != nil {
        utils.LogMessage(utils.LevelError, "Failed to set score")
        utils.LogFooter()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set score"})
    }

    utils.LogMessage(utils.LevelInfo, "Score updated")
    utils.LogFooter()
    return c.Status(fiber.StatusOK).JSON(models.BassineScoreResponse{Email: email, Score: req.Score})
}

// GetUserScore handles GET /game/bassine/:email and returns the current score (default 0 if not set)
func (h *BassineHandler) GetUserScore(c *fiber.Ctx) error {
    utils.LogHeader("🎮 Bassine Get Score")

    email := c.Params("email")
    if email == "" {
        utils.LogMessage(utils.LevelWarn, "Email is required")
        utils.LogFooter()
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email is required"})
    }

    	// Ensure user exists
	if err := h.checkUserExists(email); err != nil {
		utils.LogMessage(utils.LevelWarn, "User not found")
		utils.LogFooter()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

    var score int
    err := h.db.QueryRow("SELECT score FROM bassine_scores WHERE email = $1", email).Scan(&score)
    if err == sql.ErrNoRows {
        score = 0
    } else if err != nil {
        utils.LogMessage(utils.LevelError, "Failed to get score")
        utils.LogFooter()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get score"})
    }

    // Fetch top 3 leaderboard profile pictures
    rows, qerr := h.db.Query(`
        SELECT n.profile_picture
        FROM bassine_scores bs
        JOIN newf n ON n.email = bs.email
        WHERE n.profile_picture IS NOT NULL AND n.profile_picture <> ''
        ORDER BY bs.score DESC
        LIMIT 3
    `)
    if qerr != nil {
        utils.LogMessage(utils.LevelError, "Failed to get leaderboard")
        utils.LogFooter()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get leaderboard"})
    }
    defer rows.Close()

    top := make([]string, 0, 3)
    for rows.Next() {
        var pic string
        if err := rows.Scan(&pic); err != nil {
            utils.LogMessage(utils.LevelError, "Failed to scan leaderboard row")
            utils.LogFooter()
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get leaderboard"})
        }
        top = append(top, pic)
    }
    if err := rows.Err(); err != nil {
        utils.LogMessage(utils.LevelError, "Leaderboard rows error")
        utils.LogFooter()
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get leaderboard"})
    }

    utils.LogMessage(utils.LevelInfo, "Score retrieved")
    utils.LogFooter()
    return c.Status(fiber.StatusOK).JSON(models.BassineScoreResponse{Email: email, Score: score, TopLeaderboard: top})
}


