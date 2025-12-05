package auth

import (
	"fmt"
	"time"

	"nanny-backend/internal/common/models"
	"nanny-backend/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterOwner(fullName, email, phone, password string) error
	RegisterSitter(fullName, email, phone, password string, experienceYears int, certificates, preferences, location string) error
	Login(email, password string) (*models.User, string, error) // ← token added
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository) Service {
	cfg := config.Load() // ← берём JWT_SECRET из config.go

	return &service{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
}

func (s *service) RegisterOwner(fullName, email, phone, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	user := &models.User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Role:         "owner",
	}

	_, err = s.repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("ошибка регистрации владельца: %w", err)
	}

	return nil
}

func (s *service) RegisterSitter(fullName, email, phone, password string, experienceYears int, certificates, preferences, location string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	user := &models.User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: string(hashedPassword),
		Role:         "sitter",
	}

	userID, err := s.repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	sitter := &models.Sitter{
		SitterID:        userID,
		ExperienceYears: experienceYears,
		Certificates:    certificates,
		Preferences:     preferences,
		Location:        location,
		Status:          "pending",
	}

	err = s.repo.CreateSitter(sitter)
	if err != nil {
		return fmt.Errorf("ошибка создания профиля няни: %w", err)
	}

	return nil
}

func (s *service) Login(email, password string) (*models.User, string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, "", fmt.Errorf("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", fmt.Errorf("неверный email или пароль")
	}

	// --- Генерация JWT ---
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // токен живёт 72 часа
	})

	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", fmt.Errorf("ошибка генерации токена: %w", err)
	}

	return user, signedToken, nil
}
