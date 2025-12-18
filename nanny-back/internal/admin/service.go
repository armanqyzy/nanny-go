package admin

import (
	"fmt"

	"nanny-backend/internal/common/models"
)

type Service interface {
	GetPendingSitters() ([]models.Sitter, error)
	ApproveSitter(sitterID int) error
	RejectSitter(sitterID int) error
	GetAllUsers() ([]models.User, error)
	GetUser(userID int) (*models.User, error)
	DeleteUser(userID int) error
	GetSitterDetails(sitterID int) (*SitterDetails, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetPendingSitters() ([]models.Sitter, error) {
	return s.repo.GetPendingSitters()
}

func (s *service) ApproveSitter(sitterID int) error {
	details, err := s.repo.GetSitterDetails(sitterID)
	if err != nil {
		return err
	}

	if details.Status != "pending" {
		return fmt.Errorf("you can approve only request in status 'pending'")
	}

	return s.repo.ApproveSitter(sitterID)
}

func (s *service) RejectSitter(sitterID int) error {
	details, err := s.repo.GetSitterDetails(sitterID)
	if err != nil {
		return err
	}

	if details.Status != "pending" {
		return fmt.Errorf("you can reject only request in status 'pending'")
	}

	return s.repo.RejectSitter(sitterID)
}

func (s *service) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *service) GetUser(userID int) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *service) DeleteUser(userID int) error {
	_, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteUser(userID)
}

func (s *service) GetSitterDetails(sitterID int) (*SitterDetails, error) {
	return s.repo.GetSitterDetails(sitterID)
}
