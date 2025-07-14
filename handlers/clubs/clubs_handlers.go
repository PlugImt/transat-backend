package clubs

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/plugimt/transat-backend/handlers/clubs/repository"
	"github.com/plugimt/transat-backend/models"
	"github.com/plugimt/transat-backend/utils"
)

type ClubsHandler struct {
	ClubsRepository *repository.ClubsRepository
}

func NewClubsHandler(db *sql.DB) *ClubsHandler {
	clubsRepo := repository.NewClubsRepository(db)

	return &ClubsHandler{
		ClubsRepository: clubsRepo,
	}
}

// CreateClub creates a new club (admin only)
func (h *ClubsHandler) CreateClub(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Create Club")

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user is admin
	isAdmin, err := h.ClubsRepository.IsUserAdmin(userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check admin status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	if !isAdmin {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Non-admin user %s attempted to create club", userEmail))
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin privileges required",
		})
	}

	var req models.CreateClubRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Name == "" || req.Picture == "" {
		utils.LogMessage(utils.LevelError, "Missing required fields: name and picture")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and picture are required",
		})
	}

	result, err := h.ClubsRepository.CreateClub(req)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to create club: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create club",
		})
	}

	if !result.Success {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Club creation failed: %s", result.Message))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Club '%s' created successfully by %s", req.Name, userEmail))
	utils.LogFooter()
	return c.Status(fiber.StatusCreated).JSON(result)
}

// UpdateClub updates club details (admin or club respo only)
func (h *ClubsHandler) UpdateClub(c *fiber.Ctx) error {
	utils.LogHeader("✏️ Update Club")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user is admin or club respo
	isAdmin, err := h.ClubsRepository.IsUserAdmin(userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check admin status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	isClubRespo, err := h.ClubsRepository.IsUserClubRespo(userEmail, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check respo status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	if !isAdmin && !isClubRespo {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("User %s attempted to update club %d without permission", userEmail, clubID))
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin or club respo privileges required",
		})
	}

	var req models.UpdateClubRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	result, err := h.ClubsRepository.UpdateClub(clubID, req)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to update club: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update club",
		})
	}

	if !result.Success {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Club update failed: %s", result.Message))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Club %d updated successfully by %s", clubID, userEmail))
	utils.LogFooter()
	return c.JSON(result)
}

// AddRespo adds a respo to a club (admin or club respo only)
func (h *ClubsHandler) AddRespo(c *fiber.Ctx) error {
	utils.LogHeader("👑 Add Club Respo")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user is admin or club respo
	isAdmin, err := h.ClubsRepository.IsUserAdmin(userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check admin status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	isClubRespo, err := h.ClubsRepository.IsUserClubRespo(userEmail, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check respo status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	if !isAdmin && !isClubRespo {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("User %s attempted to add respo to club %d without permission", userEmail, clubID))
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin or club respo privileges required",
		})
	}

	var req models.AddRespoRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Email == "" {
		utils.LogMessage(utils.LevelError, "Missing required field: email")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	result, err := h.ClubsRepository.AddRespo(clubID, req.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to add respo: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add respo",
		})
	}

	if !result.Success {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Add respo failed: %s", result.Message))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("User %s added as respo to club %d by %s", req.Email, clubID, userEmail))
	utils.LogFooter()
	return c.JSON(result)
}

// RemoveRespo removes a respo from a club (admin or club respo only)
func (h *ClubsHandler) RemoveRespo(c *fiber.Ctx) error {
	utils.LogHeader("🚫 Remove Club Respo")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	// Check if user is admin or club respo
	isAdmin, err := h.ClubsRepository.IsUserAdmin(userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check admin status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	isClubRespo, err := h.ClubsRepository.IsUserClubRespo(userEmail, clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to check respo status: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to verify permissions",
		})
	}

	if !isAdmin && !isClubRespo {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("User %s attempted to remove respo from club %d without permission", userEmail, clubID))
		utils.LogFooter()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin or club respo privileges required",
		})
	}

	var req models.RemoveRespoRequest
	if err := c.BodyParser(&req); err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to parse request body: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	if req.Email == "" {
		utils.LogMessage(utils.LevelError, "Missing required field: email")
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	result, err := h.ClubsRepository.RemoveRespo(clubID, req.Email)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to remove respo: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove respo",
		})
	}

	if !result.Success {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Remove respo failed: %s", result.Message))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("User %s removed as respo from club %d by %s", req.Email, clubID, userEmail))
	utils.LogFooter()
	return c.JSON(result)
}

// JoinClub allows a user to join a club
func (h *ClubsHandler) JoinClub(c *fiber.Ctx) error {
	utils.LogHeader("✅ Join Club")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	result, err := h.ClubsRepository.JoinClub(clubID, userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to join club: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to join club",
		})
	}

	if !result.Success {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Join club failed: %s", result.Message))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("User %s joined club %d", userEmail, clubID))
	utils.LogFooter()
	return c.JSON(result)
}

// LeaveClub allows a user to leave a club
func (h *ClubsHandler) LeaveClub(c *fiber.Ctx) error {
	utils.LogHeader("❌ Leave Club")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	result, err := h.ClubsRepository.LeaveClub(clubID, userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to leave club: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to leave club",
		})
	}

	if !result.Success {
		utils.LogMessage(utils.LevelWarn, fmt.Sprintf("Leave club failed: %s", result.Message))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("User %s left club %d", userEmail, clubID))
	utils.LogFooter()
	return c.JSON(result)
}

// GetUserClubs gets all clubs the signed-in user belongs to
func (h *ClubsHandler) GetUserClubs(c *fiber.Ctx) error {
	utils.LogHeader("📋 Get User Clubs")

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	response, err := h.ClubsRepository.GetUserClubs(userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get user clubs: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user clubs",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved clubs for user %s: %d member, %d respo", userEmail, len(response.MemberOf), len(response.RespoOf)))
	utils.LogFooter()
	return c.JSON(response)
}

// GetAllClubRespos gets the list of all club respos
func (h *ClubsHandler) GetAllClubRespos(c *fiber.Ctx) error {
	utils.LogHeader("👑 Get All Club Respos")

	response, err := h.ClubsRepository.GetAllClubRespos()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get club respos: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve club respos",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved %d club respos", response.Total))
	utils.LogFooter()
	return c.JSON(response)
}

// GetClubMembers gets all members of a specific club
func (h *ClubsHandler) GetClubMembers(c *fiber.Ctx) error {
	utils.LogHeader("👥 Get Club Members")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	response, err := h.ClubsRepository.GetClubMembers(clubID)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get club members: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve club members",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved %d members for club %d", response.Total, clubID))
	utils.LogFooter()
	return c.JSON(response)
}

// GetAllClubs gets the list of all clubs with their details
func (h *ClubsHandler) GetAllClubs(c *fiber.Ctx) error {
	utils.LogHeader("🏛️ Get All Clubs")

	response, err := h.ClubsRepository.GetAllClubs()
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get clubs: %v", err))
		utils.LogFooter()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve clubs",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved %d clubs", response.Total))
	utils.LogFooter()
	return c.JSON(response)
}

// GetClubDetails gets detailed information about a specific club
func (h *ClubsHandler) GetClubDetails(c *fiber.Ctx) error {
	utils.LogHeader("🔍 Get Club Details")

	clubIDParam := c.Params("id")
	clubID, err := strconv.Atoi(clubIDParam)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Invalid club ID: %s", clubIDParam))
		utils.LogFooter()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid club ID format",
		})
	}

	userEmail, ok := c.Locals("email").(string)
	if !ok || userEmail == "" {
		utils.LogMessage(utils.LevelError, "No email found in JWT token")
		utils.LogFooter()
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	response, err := h.ClubsRepository.GetClubDetails(clubID, userEmail)
	if err != nil {
		utils.LogMessage(utils.LevelError, fmt.Sprintf("Failed to get club details: %v", err))
		utils.LogFooter()
		if err.Error() == "club not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Club not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve club details",
		})
	}

	utils.LogMessage(utils.LevelInfo, fmt.Sprintf("Retrieved details for club %d", clubID))
	utils.LogFooter()
	return c.JSON(response)
}
