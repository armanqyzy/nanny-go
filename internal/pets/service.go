package pets

import (
	"fmt"

	"nanny-backend/internal/common/models"
)

type Service interface {
	CreatePet(ownerID int, name, petType string, age int, notes string) (int, error)
	GetPetByID(petID int) (*models.Pet, error)
	GetPetsByOwner(ownerID int) ([]models.Pet, error)
	UpdatePet(petID int, name, petType string, age int, notes string) error
	DeletePet(petID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreatePet(ownerID int, name, petType string, age int, notes string) (int, error) {
	validTypes := map[string]bool{"cat": true, "dog": true, "rodent": true}
	if !validTypes[petType] {
		return 0, fmt.Errorf("неверный тип питомца. Допустимые значения: cat, dog, rodent")
	}

	pet := &models.Pet{
		OwnerID: ownerID,
		Name:    name,
		Type:    petType,
		Age:     age,
		Notes:   notes,
	}

	petID, err := s.repo.Create(pet)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания питомца: %w", err)
	}

	return petID, nil
}

func (s *service) GetPetByID(petID int) (*models.Pet, error) {
	return s.repo.GetByID(petID)
}

func (s *service) GetPetsByOwner(ownerID int) ([]models.Pet, error) {
	return s.repo.GetByOwnerID(ownerID)
}

func (s *service) UpdatePet(petID int, name, petType string, age int, notes string) error {
	validTypes := map[string]bool{"cat": true, "dog": true, "rodent": true}
	if !validTypes[petType] {
		return fmt.Errorf("неверный тип питомца. Допустимые значения: cat, dog, rodent")
	}

	pet := &models.Pet{
		PetID: petID,
		Name:  name,
		Type:  petType,
		Age:   age,
		Notes: notes,
	}

	return s.repo.Update(pet)
}

func (s *service) DeletePet(petID int) error {
	return s.repo.Delete(petID)
}
