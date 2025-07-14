package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/plugimt/transat-backend/models"
)

type ClubsRepository struct {
	DB *sql.DB
}

func NewClubsRepository(db *sql.DB) *ClubsRepository {
	return &ClubsRepository{
		DB: db,
	}
}

// CreateClub creates a new club and automatically creates the respo role
func (r *ClubsRepository) CreateClub(req models.CreateClubRequest) (*models.ClubOperationResponse, error) {
	// Check if club name already exists
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE LOWER(name) = LOWER($1))", req.Name).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check club existence: %w", err)
	}
	if exists {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "A club with this name already exists",
		}, nil
	}

	// Insert the club
	var clubID int
	err = r.DB.QueryRow(`
		INSERT INTO clubs (name, picture, description, location, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id_clubs
	`, req.Name, req.Picture, req.Description, req.Location, req.Link).Scan(&clubID)

	if err != nil {
		return nil, fmt.Errorf("failed to create club: %w", err)
	}

	return &models.ClubOperationResponse{
		Success: true,
		Message: "Club created successfully",
		ClubID:  clubID,
	}, nil
}

// UpdateClub updates club details
func (r *ClubsRepository) UpdateClub(clubID int, req models.UpdateClubRequest) (*models.ClubOperationResponse, error) {
	// Check if club exists
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE id_clubs = $1)", clubID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check club existence: %w", err)
	}
	if !exists {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "Club not found",
		}, nil
	}

	// Build dynamic update query
	var setParts []string
	var args []interface{}
	argIndex := 1

	if req.Name != "" {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}
	if req.Picture != "" {
		setParts = append(setParts, fmt.Sprintf("picture = $%d", argIndex))
		args = append(args, req.Picture)
		argIndex++
	}
	if req.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, req.Description)
		argIndex++
	}
	if req.Location != "" {
		setParts = append(setParts, fmt.Sprintf("location = $%d", argIndex))
		args = append(args, req.Location)
		argIndex++
	}
	if req.Link != "" {
		setParts = append(setParts, fmt.Sprintf("link = $%d", argIndex))
		args = append(args, req.Link)
		argIndex++
	}

	if len(setParts) == 0 {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "No fields to update",
		}, nil
	}

	query := fmt.Sprintf("UPDATE clubs SET %s WHERE id_clubs = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, clubID)

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update club: %w", err)
	}

	return &models.ClubOperationResponse{
		Success: true,
		Message: "Club updated successfully",
		ClubID:  clubID,
	}, nil
}

// AddRespo adds a respo to a club and grants the club respo role
func (r *ClubsRepository) AddRespo(clubID int, email string) (*models.ClubOperationResponse, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if club exists and get club name
	var clubName string
	err = tx.QueryRow("SELECT name FROM clubs WHERE id_clubs = $1", clubID).Scan(&clubName)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.ClubOperationResponse{
				Success: false,
				Message: "Club not found",
			}, nil
		}
		return nil, fmt.Errorf("failed to get club: %w", err)
	}

	// Check if user exists
	var userExists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM newf WHERE email = $1)", email).Scan(&userExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !userExists {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	// Create the club respo role (e.g., "club_1_respo" for club ID 1)
	roleName := fmt.Sprintf("club_%d_respo", clubID)

	// Check if role already exists, if not create it
	_, err = tx.Exec(`
		INSERT INTO roles (role_name) 
		VALUES ($1) 
		ON CONFLICT (role_name) DO NOTHING
	`, roleName)

	if err != nil {
		return nil, fmt.Errorf("failed to create club respo role: %w", err)
	}

	// Grant the club respo role to the user
	_, err = tx.Exec(`
		INSERT INTO user_roles (email, role_name)
		VALUES ($1, $2)
		ON CONFLICT (email, role_name) DO NOTHING
	`, email, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to grant respo role: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.ClubOperationResponse{
		Success: true,
		Message: "Respo added successfully",
		ClubID:  clubID,
	}, nil
}

// RemoveRespo removes a respo from a club and revokes the club respo role
func (r *ClubsRepository) RemoveRespo(clubID int, email string) (*models.ClubOperationResponse, error) {
	// Check if club exists
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE id_clubs = $1)", clubID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check club existence: %w", err)
	}
	if !exists {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "Club not found",
		}, nil
	}

	// Create the club respo role name
	roleName := fmt.Sprintf("club_%d_respo", clubID)

	// Remove the club respo role from the user
	result, err := r.DB.Exec("DELETE FROM user_roles WHERE email = $1 AND role_name = $2", email, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke respo role: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "User is not a respo of this club",
		}, nil
	}

	return &models.ClubOperationResponse{
		Success: true,
		Message: "Respo removed successfully",
		ClubID:  clubID,
	}, nil
}

// JoinClub adds a user to a club
func (r *ClubsRepository) JoinClub(clubID int, email string) (*models.ClubOperationResponse, error) {
	// Check if club exists
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs WHERE id_clubs = $1)", clubID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check club existence: %w", err)
	}
	if !exists {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "Club not found",
		}, nil
	}

	// Check if user exists
	var userExists bool
	err = r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM newf WHERE email = $1)", email).Scan(&userExists)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if !userExists {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	// Add as member
	_, err = r.DB.Exec(`
		INSERT INTO clubs_members (email, id_clubs)
		VALUES ($1, $2)
		ON CONFLICT (email, id_clubs) DO NOTHING
	`, email, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to join club: %w", err)
	}

	return &models.ClubOperationResponse{
		Success: true,
		Message: "Successfully joined club",
		ClubID:  clubID,
	}, nil
}

// LeaveClub removes a user from a club
func (r *ClubsRepository) LeaveClub(clubID int, email string) (*models.ClubOperationResponse, error) {
	result, err := r.DB.Exec("DELETE FROM clubs_members WHERE email = $1 AND id_clubs = $2", email, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to leave club: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &models.ClubOperationResponse{
			Success: false,
			Message: "You are not a member of this club",
		}, nil
	}

	return &models.ClubOperationResponse{
		Success: true,
		Message: "Successfully left club",
		ClubID:  clubID,
	}, nil
}

// GetUserClubs gets all clubs a user is member of or respo of
func (r *ClubsRepository) GetUserClubs(email string) (*models.UserClubsResponse, error) {
	// Get clubs user is member of
	memberQuery := `
		SELECT c.id_clubs, c.name, c.picture, c.description, c.location, c.link
		FROM clubs c
		JOIN clubs_members cm ON c.id_clubs = cm.id_clubs
		WHERE cm.email = $1
		ORDER BY c.name
	`

	memberRows, err := r.DB.Query(memberQuery, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get member clubs: %w", err)
	}
	defer memberRows.Close()

	var memberClubs []models.Club
	for memberRows.Next() {
		var club models.Club
		err := memberRows.Scan(&club.ID, &club.Name, &club.Picture, &club.Description, &club.Location, &club.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member club: %w", err)
		}
		memberClubs = append(memberClubs, club)
	}

	// Get clubs user is respo of by checking user_roles for club_X_respo roles
	respoQuery := `
		SELECT c.id_clubs, c.name, c.picture, c.description, c.location, c.link
		FROM clubs c
		JOIN user_roles ur ON ur.role_name = 'club_' || c.id_clubs || '_respo'
		WHERE ur.email = $1
		ORDER BY c.name
	`

	respoRows, err := r.DB.Query(respoQuery, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get respo clubs: %w", err)
	}
	defer respoRows.Close()

	var respoClubs []models.Club
	for respoRows.Next() {
		var club models.Club
		err := respoRows.Scan(&club.ID, &club.Name, &club.Picture, &club.Description, &club.Location, &club.Link)
		if err != nil {
			return nil, fmt.Errorf("failed to scan respo club: %w", err)
		}
		respoClubs = append(respoClubs, club)
	}

	return &models.UserClubsResponse{
		MemberOf: memberClubs,
		RespoOf:  respoClubs,
	}, nil
}

// GetAllClubRespos gets all respos across all clubs
func (r *ClubsRepository) GetAllClubRespos() (*models.ClubResposResponse, error) {
	query := `
		SELECT ur.email, n.first_name, n.last_name, n.profile_picture, n.graduation_year, n.campus, c.name, c.id_clubs
		FROM user_roles ur
		JOIN newf n ON ur.email = n.email
		JOIN clubs c ON ur.role_name = 'club_' || c.id_clubs || '_respo'
		ORDER BY c.name, n.first_name, n.last_name
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get club respos: %w", err)
	}
	defer rows.Close()

	var respos []struct {
		models.ClubMemberWithDetails
		ClubName string `json:"club_name"`
		ClubID   int    `json:"club_id"`
	}

	for rows.Next() {
		var respo struct {
			models.ClubMemberWithDetails
			ClubName string `json:"club_name"`
			ClubID   int    `json:"club_id"`
		}

		err := rows.Scan(&respo.Email, &respo.FirstName, &respo.LastName, &respo.ProfilePicture,
			&respo.GraduationYear, &respo.Campus, &respo.ClubName, &respo.ClubID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan respo: %w", err)
		}
		respos = append(respos, respo)
	}

	return &models.ClubResposResponse{
		Respos: respos,
		Total:  len(respos),
	}, nil
}

// GetClubMembers gets all members of a specific club
func (r *ClubsRepository) GetClubMembers(clubID int) (*models.ClubMembersResponse, error) {
	query := `
		SELECT cm.email, n.first_name, n.last_name, n.profile_picture, n.graduation_year, n.campus
		FROM clubs_members cm
		JOIN newf n ON cm.email = n.email
		WHERE cm.id_clubs = $1
		ORDER BY n.first_name, n.last_name
	`

	rows, err := r.DB.Query(query, clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get club members: %w", err)
	}
	defer rows.Close()

	var members []models.ClubMemberWithDetails
	for rows.Next() {
		var member models.ClubMemberWithDetails
		err := rows.Scan(&member.Email, &member.FirstName, &member.LastName,
			&member.ProfilePicture, &member.GraduationYear, &member.Campus)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, member)
	}

	return &models.ClubMembersResponse{
		ClubID:  clubID,
		Members: members,
		Total:   len(members),
	}, nil
}

// GetAllClubs gets all clubs with statistics
func (r *ClubsRepository) GetAllClubs() (*models.ClubListResponse, error) {
	query := `
		SELECT c.id_clubs, c.name, c.picture, c.description, c.location, c.link,
			   COALESCE(member_counts.count, 0) as member_count,
			   COALESCE(respo_counts.count, 0) as respo_count
		FROM clubs c
		LEFT JOIN (
			SELECT id_clubs, COUNT(*) as count
			FROM clubs_members
			GROUP BY id_clubs
		) member_counts ON c.id_clubs = member_counts.id_clubs
		LEFT JOIN (
			SELECT id_clubs, COUNT(*) as count
			FROM clubs_respos
			GROUP BY id_clubs
		) respo_counts ON c.id_clubs = respo_counts.id_clubs
		ORDER BY c.name
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get clubs: %w", err)
	}
	defer rows.Close()

	var clubs []models.ClubWithStats
	for rows.Next() {
		var club models.ClubWithStats
		err := rows.Scan(&club.ID, &club.Name, &club.Picture, &club.Description,
			&club.Location, &club.Link, &club.MemberCount, &club.RespoCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan club: %w", err)
		}
		clubs = append(clubs, club)
	}

	return &models.ClubListResponse{
		Clubs: clubs,
		Total: len(clubs),
	}, nil
}

// GetClubDetails gets detailed information about a specific club
func (r *ClubsRepository) GetClubDetails(clubID int, userEmail string) (*models.ClubDetailsResponse, error) {
	// Get club basic info
	var club models.Club
	err := r.DB.QueryRow(`
		SELECT id_clubs, name, picture, description, location, link
		FROM clubs
		WHERE id_clubs = $1
	`, clubID).Scan(&club.ID, &club.Name, &club.Picture, &club.Description, &club.Location, &club.Link)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("club not found")
		}
		return nil, fmt.Errorf("failed to get club: %w", err)
	}

	// Get members
	membersResp, err := r.GetClubMembers(clubID)
	if err != nil {
		return nil, fmt.Errorf("failed to get club members: %w", err)
	}

	// Get respos from user_roles table
	respoQuery := `
		SELECT ur.email, n.first_name, n.last_name, n.profile_picture, n.graduation_year, n.campus
		FROM user_roles ur
		JOIN newf n ON ur.email = n.email
		WHERE ur.role_name = $1
		ORDER BY n.first_name, n.last_name
	`

	respoRows, err := r.DB.Query(respoQuery, fmt.Sprintf("club_%d_respo", clubID))
	if err != nil {
		return nil, fmt.Errorf("failed to get club respos: %w", err)
	}
	defer respoRows.Close()

	var respos []models.ClubMemberWithDetails
	for respoRows.Next() {
		var respo models.ClubMemberWithDetails
		err := respoRows.Scan(&respo.Email, &respo.FirstName, &respo.LastName,
			&respo.ProfilePicture, &respo.GraduationYear, &respo.Campus)
		if err != nil {
			return nil, fmt.Errorf("failed to scan respo: %w", err)
		}
		respos = append(respos, respo)
	}

	// Check if user is member or respo
	var isMember, isRespo bool
	r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM clubs_members WHERE email = $1 AND id_clubs = $2)", userEmail, clubID).Scan(&isMember)
	r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_roles WHERE email = $1 AND role_name = $2)", userEmail, fmt.Sprintf("club_%d_respo", clubID)).Scan(&isRespo)

	return &models.ClubDetailsResponse{
		Club:        club,
		MemberCount: membersResp.Total,
		RespoCount:  len(respos),
		Members:     membersResp.Members,
		Respos:      respos,
		IsMember:    isMember,
		IsRespo:     isRespo,
	}, nil
}

// IsUserAdmin checks if a user has admin role
func (r *ClubsRepository) IsUserAdmin(email string) (bool, error) {
	var isAdmin bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_roles 
			WHERE email = $1 AND role_name = 'admin'
		)
	`, email).Scan(&isAdmin)

	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %w", err)
	}

	return isAdmin, nil
}

// IsUserClubRespo checks if a user is a respo of a specific club
func (r *ClubsRepository) IsUserClubRespo(email string, clubID int) (bool, error) {
	var isRespo bool
	err := r.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM user_roles 
			WHERE email = $1 AND role_name = $2
		)
	`, email, fmt.Sprintf("club_%d_respo", clubID)).Scan(&isRespo)

	if err != nil {
		return false, fmt.Errorf("failed to check respo status: %w", err)
	}

	return isRespo, nil
}
