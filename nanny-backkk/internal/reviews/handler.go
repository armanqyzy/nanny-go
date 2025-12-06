package reviews

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

type CreateReviewRequest struct {
	BookingID int    `json:"booking_id" validate:"required,gt=0"`
	OwnerID   int    `json:"owner_id" validate:"required,gt=0"`
	SitterID  int    `json:"sitter_id" validate:"required,gt=0"`
	Rating    int    `json:"rating" validate:"required,gte=1,lte=5"`
	Comment   string `json:"comment,omitempty" validate:"max=1000"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
	Comment string `json:"comment,omitempty" validate:"max=1000"`
}

func (h *Handler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	reviewID, err := h.service.CreateReview(
		req.BookingID,
		req.OwnerID,
		req.SitterID,
		req.Rating,
		req.Comment,
	)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":   "отзыв создан успешно",
		"review_id": reviewID,
	})
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID отзыва")
		return
	}

	if reviewID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID отзыва должен быть положительным числом")
		return
	}

	review, err := h.service.GetReview(reviewID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, review)
}

func (h *Handler) GetSitterReviews(w http.ResponseWriter, r *http.Request) {
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

	reviews, err := h.service.GetSitterReviews(sitterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, reviews)
}

func (h *Handler) GetSitterRating(w http.ResponseWriter, r *http.Request) {
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

	avgRating, count, err := h.service.GetSitterRating(sitterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"sitter_id":      sitterID,
		"average_rating": avgRating,
		"review_count":   count,
	})
}

func (h *Handler) GetBookingReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["booking_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID бронирования")
		return
	}

	if bookingID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID бронирования должен быть положительным числом")
		return
	}

	review, err := h.service.GetBookingReview(bookingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, review)
}

func (h *Handler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID отзыва")
		return
	}

	if reviewID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID отзыва должен быть положительным числом")
		return
	}

	var req UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "неверные данные")
		return
	}

	if err := validator.Validate(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UpdateReview(reviewID, req.Rating, req.Comment)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "отзыв обновлён успешно",
	})
}

func (h *Handler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "неверный ID отзыва")
		return
	}

	if reviewID <= 0 {
		respondWithError(w, http.StatusBadRequest, "ID отзыва должен быть положительным числом")
		return
	}

	err = h.service.DeleteReview(reviewID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "отзыв удалён успешно",
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
