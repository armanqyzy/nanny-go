package auth

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type RegisterOwnerRequest struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterSitterRequest struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	ExperienceYears int    `json:"experience_years"`
	Certificates    string `json:"certificates"`
	Preferences     string `json:"preferences"`
	Location        string `json:"location"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) RegisterOwner(w http.ResponseWriter, r *http.Request) {
	var req RegisterOwnerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	err := h.service.RegisterOwner(req.FullName, req.Email, req.Phone, req.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "владелец зарегистрирован успешно",
	})
}

func (h *Handler) RegisterSitter(w http.ResponseWriter, r *http.Request) {
	var req RegisterSitterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
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
		"message": "няня зарегистрирована, ожидает подтверждения",
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	// сервис теперь возвращает: user, token, err
	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Отправляем токен и данные пользователя
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "вход выполнен",
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
