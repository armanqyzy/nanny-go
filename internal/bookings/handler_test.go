package bookings

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"nanny-backend/internal/common/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateBooking(
	ownerID, sitterID, petID, serviceID int,
	startTime, endTime time.Time,
) (int, error) {
	args := m.Called(ownerID, sitterID, petID, serviceID, startTime, endTime)
	return args.Int(0), args.Error(1)
}

func (m *MockService) GetBookingByID(bookingID int) (*models.Booking, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockService) GetOwnerBookings(ownerID int) ([]models.Booking, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockService) GetSitterBookings(sitterID int) ([]models.Booking, error) {
	args := m.Called(sitterID)
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockService) ConfirmBooking(bookingID int) error {
	return m.Called(bookingID).Error(0)
}

func (m *MockService) CancelBooking(bookingID int) error {
	return m.Called(bookingID).Error(0)
}

func (m *MockService) CompleteBooking(bookingID int) error {
	return m.Called(bookingID).Error(0)
}

func TestHandler_CreateBooking_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	reqBody := map[string]interface{}{
		"owner_id":   1,
		"sitter_id":  2,
		"pet_id":     3,
		"service_id": 4,
		"start_time": startTime.Format(time.RFC3339),
		"end_time":   endTime.Format(time.RFC3339),
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On(
			"CreateBooking",
			1, 2, 3, 4,
			mock.AnythingOfType("time.Time"),
			mock.AnythingOfType("time.Time"),
		).
		Return(42, nil)

	handler.CreateBooking(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, float64(42), resp["booking_id"])

	mockService.AssertExpectations(t)
}

func TestHandler_GetBookingByID_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	booking := &models.Booking{
		BookingID: 1,
		OwnerID:   2,
		SitterID:  3,
		Status:    "pending",
	}

	mockService.
		On("GetBookingByID", 1).
		Return(booking, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/bookings/1", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/bookings/{id}", handler.GetBooking)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_ConfirmBooking_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.
		On("ConfirmBooking", 10).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/bookings/10/confirm", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/bookings/{id}/confirm", handler.ConfirmBooking)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

func TestHandler_CancelBooking_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.
		On("CancelBooking", 10).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/bookings/10/cancel", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/bookings/{id}/cancel", handler.CancelBooking)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}

func TestHandler_CompleteBooking_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	mockService.
		On("CompleteBooking", 10).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/bookings/10/complete", nil)
	rec := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/bookings/{id}/complete", handler.CompleteBooking)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	mockService.AssertExpectations(t)
}
