package admin

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetPendingSitters возвращает список заявок нянь на рассмотрении
func (h *Handler) GetPendingSitters(w http.ResponseWriter, r *http.Request) {
	sitters, err := h.service.GetPendingSitters()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, sitters)
}

// ApproveSitter одобряет заявку няни
func (h *Handler) ApproveSitter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sitterID, err := strconv.Atoi(vars["sitter_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID няни")
		return
	}

	err = h.service.ApproveSitter(sitterID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "няня одобрена успешно",
	})
}

// RejectSitter отклоняет заявку няни
func (h *Handler) RejectSitter(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sitterID, err := strconv.Atoi(vars["sitter_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID няни")
		return
	}

	err = h.service.RejectSitter(sitterID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "няня отклонена",
	})
}

// GetAllUsers возвращает список всех пользователей
func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

// GetUser возвращает информацию о конкретном пользователе
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID пользователя")
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser удаляет пользователя
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID пользователя")
		return
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "пользователь удалён успешно",
	})
}

// GetSitterDetails возвращает детальную информацию о няне
func (h *Handler) GetSitterDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sitterID, err := strconv.Atoi(vars["sitter_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID няни")
		return
	}

	details, err := h.service.GetSitterDetails(sitterID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, details)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
