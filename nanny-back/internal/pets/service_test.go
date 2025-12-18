package pets

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"nanny-backend/internal/common/models"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(pet *models.Pet) (int, error) {
	args := m.Called(pet)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetByID(petID int) (*models.Pet, error) {
	args := m.Called(petID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pet), args.Error(1)
}

func (m *MockRepository) GetByOwnerID(ownerID int) ([]models.Pet, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]models.Pet), args.Error(1)
}

func (m *MockRepository) Update(pet *models.Pet) error {
	args := m.Called(pet)
	return args.Error(0)
}

func (m *MockRepository) Delete(petID int) error {
	args := m.Called(petID)
	return args.Error(0)
}

func TestCreatePet_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("Create", mock.Anything).
		Return(1, nil)

	petID, err := service.CreatePet(
		1,
		"Buddy",
		"dog",
		3,
		"Friendly dog",
	)

	assert.NoError(t, err)
	assert.Equal(t, 1, petID)
	mockRepo.AssertExpectations(t)
}

func TestCreatePet_InvalidType(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	petID, err := service.CreatePet(
		1,
		"Buddy",
		"dragon",
		3,
		"",
	)

	assert.Error(t, err)
	assert.Equal(t, 0, petID)
}

func TestGetPetByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedPet := &models.Pet{
		PetID:   1,
		OwnerID: 10,
		Name:    "Buddy",
		Type:    "dog",
		Age:     3,
		Notes:   "Good boy",
	}

	mockRepo.
		On("GetByID", 1).
		Return(expectedPet, nil)

	pet, err := service.GetPetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedPet, pet)
	mockRepo.AssertExpectations(t)
}

func TestGetPetByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("GetByID", 99).
		Return(nil, errors.New("pet not found"))

	pet, err := service.GetPetByID(99)

	assert.Error(t, err)
	assert.Nil(t, pet)
}

func TestGetPetsByOwner_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedPets := []models.Pet{
		{PetID: 1, OwnerID: 5, Name: "Catty"},
		{PetID: 2, OwnerID: 5, Name: "Doggy"},
	}

	mockRepo.
		On("GetByOwnerID", 5).
		Return(expectedPets, nil)

	pets, err := service.GetPetsByOwner(5)

	assert.NoError(t, err)
	assert.Len(t, pets, 2)
	assert.Equal(t, expectedPets, pets)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePet_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("Update", mock.Anything).
		Return(nil)

	err := service.UpdatePet(
		1,
		"NewName",
		"cat",
		4,
		"Updated notes",
	)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePet_InvalidType(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	err := service.UpdatePet(
		1,
		"Name",
		"dragon",
		4,
		"",
	)

	assert.Error(t, err)
}

func TestDeletePet_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("Delete", 1).
		Return(nil)

	err := service.DeletePet(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
