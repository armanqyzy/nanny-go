package pets

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nanny-backend/internal/common/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreatePet(ownerID int, name, petType string, age int, notes string) (int, error) {
	args := m.Called(ownerID, name, petType, age, notes)
	return args.Int(0), args.Error(1)
}

func (m *MockService) GetPetByID(petID int) (*models.Pet, error) {
	args := m.Called(petID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Pet), args.Error(1)
}

func (m *MockService) GetPetsByOwner(ownerID int) ([]models.Pet, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]models.Pet), args.Error(1)
}

func (m *MockService) UpdatePet(petID int, name, petType string, age int, notes string) error {
	args := m.Called(petID, name, petType, age, notes)
	return args.Error(0)
}

func (m *MockService) DeletePet(petID int) error {
	args := m.Called(petID)
	return args.Error(0)
}

func TestHandler_CreatePet_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := CreatePetRequest{
		OwnerID: 1,
		Name:    "Мурка",
		Type:    "кошка", // ВАЖНО: pet_type на русском
		Age:     3,
		Notes:   "спокойная",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On("CreatePet", 1, "Мурка", "кошка", 3, "спокойная").
		Return(10, nil)

	handler.CreatePet(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(10), resp["pet_id"])

	mockService.AssertExpectations(t)
}

func TestHandler_GetPet_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	pet := &models.Pet{
		PetID:   1,
		OwnerID: 2,
		Name:    "Мурка",
		Type:    "кошка",
		Age:     3,
	}

	mockService.
		On("GetPetByID", 1).
		Return(pet, nil)

	req := httptest.NewRequest(http.MethodGet, "/pets/1", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/pets/{id}", handler.GetPet)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetOwnerPets_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	pets := []models.Pet{
		{PetID: 1, Name: "Мурка"},
		{PetID: 2, Name: "Барсик"},
	}

	mockService.
		On("GetPetsByOwner", 5).
		Return(pets, nil)

	req := httptest.NewRequest(http.MethodGet, "/owners/5/pets", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/owners/{owner_id}/pets", handler.GetOwnerPets)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_UpdatePet_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := UpdatePetRequest{
		Name:  "Мурка обновлённая",
		Type:  "кошка", // ВАЖНО
		Age:   4,
		Notes: "обновлено",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/pets/10", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On("UpdatePet", 10, "Мурка обновлённая", "кошка", 4, "обновлено").
		Return(nil)

	router := mux.NewRouter()
	router.HandleFunc("/pets/{id}", handler.UpdatePet)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_DeletePet_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.
		On("DeletePet", 10).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/pets/10", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/pets/{id}", handler.DeletePet)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}
