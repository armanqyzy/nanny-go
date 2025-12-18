package auth

import (
	"encoding/json"
	"net/http"

	"nanny-backend/pkg/validator"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type RegisterOwnerRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required,phone_kz"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RegisterSitterRequest struct {
	FullName        string `json:"full_name" validate:"required,min=2,max=100"`
	Email           string `json:"email" validate:"required,email"`
	Phone           string `json:"phone" validate:"required,phone_kz"`
	Password        string `json:"password" validate:"required,min=8,max=72"`
	ExperienceYears int    `json:"experience_years" validate:"required,gte=0,lte=50"`
	Certificates    string `json:"certificates" validate:"max=500"`
	Preferences     string `json:"preferences" validate:"max=500"`
	Location        string `json:"location" validate:"required,min=2,max=200"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

func (h *Handler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	var req RegisterOwnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect data")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.RegisterOwner(req.FullName, req.Email, req.Phone, req.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "owner registered succesfully",
	})
}

func (h *Handler) RegisterSitter(w http.ResponseWriter, r *http.Request) {
	var req RegisterSitterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect data")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.service.RegisterSitter(
		req.FullName,
		req.Email,
		req.Phone,
		req.Password,
		req.ExperienceYears,
		req.Certificates,
		req.Preferences,
		req.Location,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "nanny registered, expecting acceptance",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect data")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "login happened",
		"user_id":   user.UserID,
		"role":      user.Role,
		"email":     user.Email,
		"full_name": user.FullName,
		"token":     token,
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
