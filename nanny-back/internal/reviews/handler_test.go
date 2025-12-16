package reviews

import (
	"bytes"
	"encoding/json"
	"nanny-backend/internal/common/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateReview(bookingID, ownerID, sitterID, rating int, comment string) (int, error) {
	args := m.Called(bookingID, ownerID, sitterID, rating, comment)
	return args.Int(0), args.Error(1)
}

func (m *MockService) GetReview(reviewID int) (*models.Review, error) {
	args := m.Called(reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Review), args.Error(1)
}

func (m *MockService) GetSitterReviews(sitterID int) ([]models.Review, error) {
	args := m.Called(sitterID)
	return args.Get(0).([]models.Review), args.Error(1)
}

func (m *MockService) GetBookingReview(bookingID int) (*models.Review, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Review), args.Error(1)
}

func (m *MockService) UpdateReview(reviewID, rating int, comment string) error {
	args := m.Called(reviewID, rating, comment)
	return args.Error(0)
}

func (m *MockService) DeleteReview(reviewID int) error {
	args := m.Called(reviewID)
	return args.Error(0)
}

func (m *MockService) GetSitterRating(sitterID int) (float64, int, error) {
	args := m.Called(sitterID)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func TestHandler_CreateReview_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := CreateReviewRequest{
		BookingID: 1,
		OwnerID:   2,
		SitterID:  3,
		Rating:    5,
		Comment:   "Отличная няня",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/reviews", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On("CreateReview", 1, 2, 3, 5, "Отличная няня").
		Return(10, nil)

	handler.CreateReview(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(10), resp["review_id"])

	mockService.AssertExpectations(t)
}

func TestHandler_GetReview_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	review := &models.Review{
		ReviewID:  1,
		BookingID: 10,
		OwnerID:   2,
		SitterID:  3,
		Rating:    5,
		Comment:   "Отлично",
	}

	mockService.
		On("GetReview", 1).
		Return(review, nil)

	req := httptest.NewRequest(http.MethodGet, "/reviews/1", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/reviews/{id}", handler.GetReview)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetSitterReviews_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reviews := []models.Review{
		{ReviewID: 1, Rating: 5},
		{ReviewID: 2, Rating: 4},
	}

	mockService.
		On("GetSitterReviews", 3).
		Return(reviews, nil)

	req := httptest.NewRequest(http.MethodGet, "/sitters/3/reviews", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/sitters/{sitter_id}/reviews", handler.GetSitterReviews)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_UpdateReview_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := UpdateReviewRequest{
		Rating:  4,
		Comment: "Хорошо",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/reviews/5", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On("UpdateReview", 5, 4, "Хорошо").
		Return(nil)

	router := mux.NewRouter()
	router.HandleFunc("/reviews/{id}", handler.UpdateReview)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_DeleteReview_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.
		On("DeleteReview", 7).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/reviews/7", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/reviews/{id}", handler.DeleteReview)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetSitterRating_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.
		On("GetSitterRating", 3).
		Return(4.5, 10, nil)

	req := httptest.NewRequest(http.MethodGet, "/sitters/3/rating", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/sitters/{sitter_id}/rating", handler.GetSitterRating)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}
