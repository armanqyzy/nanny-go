package pets

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"nanny-backend/internal/common/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreatePet_InvalidJSON(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewBuffer([]byte("invalid json")))
	rec := httptest.NewRecorder()

	handler.CreatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверные данные")
}

func TestHandler_CreatePet_ValidationError(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"owner_id": 1,
		"name":     "Buddy",
		"type":     "invalid-type",
		"age":      5,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CreatePet_ServiceError(t *testing.T) {
	mockSvc := &mockPetService{
		createPetFunc: func(ownerID int, name, petType string, age int, notes string) (int, error) {
			return 0, errors.New("database error")
		},
	}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"owner_id": 1,
		"name":     "Max",
		"type":     "dog",
		"age":      3,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/pets", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.CreatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_GetPet_InvalidID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/pets/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.GetPet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID питомца")
}

func TestHandler_GetPet_ZeroID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/pets/0", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	rec := httptest.NewRecorder()

	handler.GetPet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "ID питомца должен быть положительным числом")
}

func TestHandler_GetPet_NotFound(t *testing.T) {
	mockSvc := &mockPetService{
		getPetByIDFunc: func(petID int) (*models.Pet, error) {
			return nil, errors.New("питомец не найден")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/pets/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	handler.GetPet(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandler_GetOwnerPets_InvalidID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/abc/pets", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "abc"})
	rec := httptest.NewRecorder()

	handler.GetOwnerPets(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID владельца")
}

func TestHandler_GetOwnerPets_ZeroID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/0/pets", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "0"})
	rec := httptest.NewRecorder()

	handler.GetOwnerPets(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_GetOwnerPets_ServiceError(t *testing.T) {
	mockSvc := &mockPetService{
		getPetsByOwnerFunc: func(ownerID int) ([]models.Pet, error) {
			return nil, errors.New("database error")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/1/pets", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetOwnerPets(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandler_UpdatePet_InvalidID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"name": "Updated",
		"type": "dog",
		"age":  5,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/pets/invalid", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.UpdatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_UpdatePet_ZeroID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"name": "Updated",
		"type": "dog",
		"age":  5,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/pets/0", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	rec := httptest.NewRecorder()

	handler.UpdatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_UpdatePet_InvalidJSON(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/pets/1", bytes.NewBuffer([]byte("bad json")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	handler.UpdatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверные данные")
}

func TestHandler_UpdatePet_ValidationError(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"name": "X",
		"type": "invalid-type",
		"age":  -5,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/pets/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	handler.UpdatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_UpdatePet_ServiceError(t *testing.T) {
	mockSvc := &mockPetService{
		updatePetFunc: func(petID int, name, petType string, age int, notes string) error {
			return errors.New("pet not found")
		},
	}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"name": "Updated",
		"type": "dog",
		"age":  5,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/pets/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	handler.UpdatePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_DeletePet_InvalidID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/pets/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.DeletePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_DeletePet_ZeroID(t *testing.T) {
	mockSvc := &mockPetService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/pets/0", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	rec := httptest.NewRecorder()

	handler.DeletePet(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_DeletePet_ServiceError(t *testing.T) {
	mockSvc := &mockPetService{
		deletePetFunc: func(petID int) error {
			return errors.New("pet not found")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/pets/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	handler.DeletePet(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

type mockPetService struct {
	createPetFunc      func(int, string, string, int, string) (int, error)
	getPetByIDFunc     func(int) (*models.Pet, error)
	getPetsByOwnerFunc func(int) ([]models.Pet, error)
	updatePetFunc      func(int, string, string, int, string) error
	deletePetFunc      func(int) error
}

func (m *mockPetService) CreatePet(ownerID int, name, petType string, age int, notes string) (int, error) {
	if m.createPetFunc != nil {
		return m.createPetFunc(ownerID, name, petType, age, notes)
	}
	return 1, nil
}

func (m *mockPetService) GetPetByID(petID int) (*models.Pet, error) {
	if m.getPetByIDFunc != nil {
		return m.getPetByIDFunc(petID)
	}
	return &models.Pet{PetID: petID}, nil
}

func (m *mockPetService) GetPetsByOwner(ownerID int) ([]models.Pet, error) {
	if m.getPetsByOwnerFunc != nil {
		return m.getPetsByOwnerFunc(ownerID)
	}
	return []models.Pet{}, nil
}

func (m *mockPetService) UpdatePet(petID int, name, petType string, age int, notes string) error {
	if m.updatePetFunc != nil {
		return m.updatePetFunc(petID, name, petType, age, notes)
	}
	return nil
}

func (m *mockPetService) DeletePet(petID int) error {
	if m.deletePetFunc != nil {
		return m.deletePetFunc(petID)
	}
	return nil
}
