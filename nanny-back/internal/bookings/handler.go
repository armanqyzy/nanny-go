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
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect data")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect format start date (use ISO 8601)")
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect format end date (use ISO 8601)")
		return
	}

	if endTime.Before(startTime) {
		respondWithError(w, http.StatusBadRequest, "end time must be later than start time")
		return
	}

	if startTime.Before(time.Now()) {
		respondWithError(w, http.StatusBadRequest, "start time cannot be past time")
		return
	}

	duration := endTime.Sub(startTime)
	if duration.Hours() > 24 {
		respondWithError(w, http.StatusBadRequest, "max duration booking  - 24 hours")
		return
	}

	if duration.Minutes() < 30 {
		respondWithError(w, http.StatusBadRequest, "min duration booking - 30 min")
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
		"message":    "booking created successfully",
		"booking_id": bookingID,
	})
}

func (h *Handler) GetBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID booking")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID booking must be positive")
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
		respondWithError(w, http.StatusBadRequest, "incorrect ID owner")
		return
	}

	if ownerID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID owner must be positive")
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
		respondWithError(w, http.StatusBadRequest, "incorrect ID nanny")
		return
	}

	if sitterID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID nanny must be positive")
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
		respondWithError(w, http.StatusBadRequest, "incorrect ID booking")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID booking must be positive")
		return
	}

	err = h.service.ConfirmBooking(bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "booking confirmed",
	})
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID booking")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID booking must be positive")
		return
	}

	err = h.service.CancelBooking(bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "booking declined",
	})
}

func (h *Handler) CompleteBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID booking")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID booking must be positive")
		return
	}

	err = h.service.CompleteBooking(bookingID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "booking completed",
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
