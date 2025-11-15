package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"nanny-backend/internal/common/models"
)

type Service interface {
	RegisterOwner(fullName, email, phone, password string) error
	RegisterSitter(fullName, email, phone, password string, experienceYears int, certificates, preferences, location string) error
	Login(email, password string) (*models.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) RegisterOwner(fullName, email, phone, password string) error {
	// Хешируем пароль
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
	// Хешируем пароль
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

	// Создаём пользователя
	userID, err := s.repo.CreateUser(user)
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	// Создаём профиль няни
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

func (s *service) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	return user, nil
}
