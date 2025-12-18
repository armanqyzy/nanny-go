package admin

import (
	"database/sql"
	"fmt"

	"nanny-backend/internal/common/models"
)

type Repository interface {
	GetPendingSitters() ([]models.Sitter, error)
	ApproveSitter(sitterID int) error
	RejectSitter(sitterID int) error
	GetAllUsers() ([]models.User, error)
	GetUserByID(userID int) (*models.User, error)
	DeleteUser(userID int) error
	GetSitterDetails(sitterID int) (*SitterDetails, error)
	UpdateSitterStatus(sitterID int, status string) error
}

type SitterDetails struct {
	models.Sitter
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Rating   float64 `json:"rating"`
	Reviews  int     `json:"reviews"`
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetPendingSitters() ([]models.Sitter, error) {
	rows, err := r.db.Query(`
		SELECT sitter_id, experience_years, certificates, preferences, location, status
		FROM sitters
		WHERE status = 'pending'
		ORDER BY sitter_id DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения заявок: %w", err)
	}
	defer rows.Close()

	var sitters []models.Sitter
	for rows.Next() {
		var sitter models.Sitter
		err := rows.Scan(
			&sitter.SitterID,
			&sitter.ExperienceYears,
			&sitter.Certificates,
			&sitter.Preferences,
			&sitter.Location,
			&sitter.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning nanny: %w", err)
		}
		sitters = append(sitters, sitter)
	}

	return sitters, nil
}

func (r *repository) ApproveSitter(sitterID int) error {
	_, err := r.db.Exec(`
		UPDATE sitters
		SET status = 'approved'
		WHERE sitter_id = $1
	`, sitterID)

	if err != nil {
		return fmt.Errorf("could not approve a nanny: %w", err)
	}

	return nil
}

func (r *repository) RejectSitter(sitterID int) error {
	_, err := r.db.Exec(`
		UPDATE sitters
		SET status = 'rejected'
		WHERE sitter_id = $1
	`, sitterID)

	if err != nil {
		return fmt.Errorf("could not decline a nanny: %w", err)
	}

	return nil
}

func (r *repository) GetAllUsers() ([]models.User, error) {
	rows, err := r.db.Query(`
		SELECT user_id, full_name, email, phone, role, created_at
		FROM users
		ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.UserID,
			&user.FullName,
			&user.Email,
			&user.Phone,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning users: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *repository) GetUserByID(userID int) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT user_id, full_name, email, phone, role, created_at
		FROM users
		WHERE user_id = $1
	`, userID).Scan(
		&user.UserID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}

	return user, nil
}

func (r *repository) DeleteUser(userID int) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("coould not try to delete user: %w", err)
	}
	return nil
}

func (r *repository) GetSitterDetails(sitterID int) (*SitterDetails, error) {
	details := &SitterDetails{}

	err := r.db.QueryRow(`
		SELECT 
			s.sitter_id, s.experience_years, s.certificates, s.preferences, s.location, s.status,
			u.full_name, u.email, u.phone,
			COALESCE(AVG(r.rating), 0) as rating,
			COUNT(r.review_id) as reviews
		FROM sitters s
		JOIN users u ON s.sitter_id = u.user_id
		LEFT JOIN reviews r ON s.sitter_id = r.sitter_id
		WHERE s.sitter_id = $1
		GROUP BY s.sitter_id, u.full_name, u.email, u.phone
	`, sitterID).Scan(
		&details.SitterID,
		&details.ExperienceYears,
		&details.Certificates,
		&details.Preferences,
		&details.Location,
		&details.Status,
		&details.FullName,
		&details.Email,
		&details.Phone,
		&details.Rating,
		&details.Reviews,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("nanny not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting details of a nanny: %w", err)
	}

	return details, nil
}

func (r *repository) UpdateSitterStatus(sitterID int, status string) error {
	_, err := r.db.Exec(`
		UPDATE sitters
		SET status = $1
		WHERE sitter_id = $2
	`, status, sitterID)

	if err != nil {
		return fmt.Errorf("could not refresh the status of a nanny: %w", err)
	}

	return nil
}
