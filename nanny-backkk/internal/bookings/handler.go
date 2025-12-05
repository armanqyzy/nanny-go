package bookings

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateBookingRequest struct {
	OwnerID   int    `json:"owner_id"`
	SitterID  int    `json:"sitter_id"`
	PetID     int    `json:"pet_id"`
	ServiceID int    `json:"service_id"`
	StartTime string `json:"start_time"` // ISO 8601 format
	EndTime   string `json:"end_time"`   // ISO 8601 format
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	// Парсим время
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный формат времени начала")
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный формат времени окончания")
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

	booking, err := h.service.GetBookingByID(bookingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, booking)
}
func (h *Handler) GetOwnerBookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ownerIDStr := vars["owner_id"]

	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID владельца")
		return
	}

	// Получаем бронирования у сервиса
	bookings, err := h.service.GetOwnerBookings(ownerID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 🔴 ВАЖНО:
	// Если сервис вернул nil-срез, json.Encoder закодирует его как null.
	// Фронт ожидает массив, поэтому подменяем на пустой массив.
	var resp interface{} = bookings
	if bookings == nil {
		// []any{} сериализуется в "[]"
		resp = []any{}
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetSitterBookings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sitterID, err := strconv.Atoi(vars["sitter_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID няни")
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
