package bookings

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"nanny-backend/pkg/validator"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateBookingRequest struct {
	OwnerID   int    `json:"owner_id" validate:"required,gt=0"`
	SitterID  int    `json:"sitter_id" validate:"required,gt=0"`
	PetID     int    `json:"pet_id" validate:"required,gt=0"`
	ServiceID int    `json:"service_id" validate:"required,gt=0"`
	StartTime string `json:"start_time" validate:"required"` // ISO 8601 format
	EndTime   string `json:"end_time" validate:"required"`   // ISO 8601 format
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный формат времени начала (используйте ISO 8601)")
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный формат времени окончания (используйте ISO 8601)")
		return
	}

	if endTime.Before(startTime) {
		respondWithError(w, http.StatusBadRequest, "время окончания должно быть позже времени начала")
		return
	}

	if startTime.Before(time.Now()) {
		respondWithError(w, http.StatusBadRequest, "время начала не может быть в прошлом")
		return
	}

	duration := endTime.Sub(startTime)
	if duration.Hours() > 24 {
		respondWithError(w, http.StatusBadRequest, "максимальная длительность бронирования - 24 часа")
		return
	}

	if duration.Minutes() < 30 {
		respondWithError(w, http.StatusBadRequest, "минимальная длительность бронирования - 30 минут")
		return
	}

	bookingID, err := h.service.CreateBooking(
		req.OwnerID,
		req.SitterID,
		req.PetID,
		req.ServiceID,
		startTime,
		endTime,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "бронирование создано успешно",
		"booking_id": bookingID,
	})
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID бронирования")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID бронирования должен быть положительным числом")
		return
	}

	booking, err := h.service.GetBookingByID(bookingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, booking)
}

func (h *Handler) GetOwnerBookings(w http.ResponseWriter, r *http.Request) {
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

	bookings, err := h.service.GetOwnerBookings(ownerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, bookings)
}

func (h *Handler) GetSitterBookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sitterID, err := strconv.Atoi(vars["sitter_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID няни")
		return
	}

	if sitterID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID няни должен быть положительным числом")
		return
	}

	bookings, err := h.service.GetSitterBookings(sitterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, bookings)
}

func (h *Handler) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID бронирования")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID бронирования должен быть положительным числом")
		return
	}

	err = h.service.ConfirmBooking(bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "бронирование подтверждено",
	})
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID бронирования")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID бронирования должен быть положительным числом")
		return
	}

	err = h.service.CancelBooking(bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "бронирование отменено",
	})
}

func (h *Handler) CompleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID бронирования")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID бронирования должен быть положительным числом")
		return
	}

	err = h.service.CompleteBooking(bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "бронирование завершено",
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
