package auth

import (
	"database/sql"
	"fmt"

	"nanny-backend/internal/common/models"
)

type Repository interface {
	CreateUser(user *models.User) (int, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateSitter(sitter *models.Sitter) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *models.User) (int, error) {
	var userID int
	err := r.db.QueryRow(`
		INSERT INTO users (full_name, email, phone, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id
	`, user.FullName, user.Email, user.Phone, user.PasswordHash, user.Role).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать пользователя: %w", err)
	}

	return userID, nil
}

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(`
		SELECT user_id, full_name, email, phone, password_hash, role, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.UserID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("пользователь не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}

	return user, nil
}

func (r *repository) CreateSitter(sitter *models.Sitter) error {
	_, err := r.db.Exec(`
		INSERT INTO sitters (sitter_id, experience_years, certificates, preferences, location, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, sitter.SitterID, sitter.ExperienceYears, sitter.Certificates, sitter.Preferences, sitter.Location, sitter.Status)

	if err != nil {
		return fmt.Errorf("не удалось создать профиль няни: %w", err)
	}

	return nil
}
