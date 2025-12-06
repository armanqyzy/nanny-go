package pets

import (
	"encoding/json"
	"net/http"
	"strconv"

	"nanny-backend/pkg/validator"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreatePetRequest struct {
	OwnerID int    `json:"owner_id" validate:"required,gt=0"`
	Name    string `json:"name" validate:"required,min=1,max=100"`
	Type    string `json:"type" validate:"required,pet_type"`
	Age     int    `json:"age" validate:"required,gte=0,lte=30"`
	Notes   string `json:"notes,omitempty" validate:"max=500"`
}

type UpdatePetRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=100"`
	Type  string `json:"type" validate:"required,pet_type"`
	Age   int    `json:"age" validate:"required,gte=0,lte=30"`
	Notes string `json:"notes,omitempty" validate:"max=500"`
}

func (h *Handler) CreatePet(w http.ResponseWriter, r *http.Request) {
	var req CreatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	petID, err := h.service.CreatePet(req.OwnerID, req.Name, req.Type, req.Age, req.Notes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "питомец создан успешно",
		"pet_id":  petID,
	})
}

func (h *Handler) GetPet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID питомца")
		return
	}

	if petID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID питомца должен быть положительным числом")
		return
	}

	pet, err := h.service.GetPetByID(petID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, pet)
}

func (h *Handler) GetOwnerPets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerID, err := strconv.Atoi(vars["owner_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID владельца")
		return
	}

	if ownerID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID владельца должен быть положительным числом")
		return
	}

	pets, err := h.service.GetPetsByOwner(ownerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, pets)
}

func (h *Handler) UpdatePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID питомца")
		return
	}

	if petID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID питомца должен быть положительным числом")
		return
	}

	var req UpdatePetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdatePet(petID, req.Name, req.Type, req.Age, req.Notes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "питомец обновлён успешно",
	})
}

func (h *Handler) DeletePet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	petID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID питомца")
		return
	}

	if petID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID питомца должен быть положительным числом")
		return
	}

	err = h.service.DeletePet(petID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "питомец удалён успешно",
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
