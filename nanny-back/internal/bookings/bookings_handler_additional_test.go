package bookings

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"nanny-backend/internal/common/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetOwnerBookings_Success(t *testing.T) {
	mockSvc := &mockBookingService{
		getOwnerBookingsFunc: func(ownerID int) ([]models.Booking, error) {
			return []models.Booking{
				{BookingID: 1, OwnerID: ownerID, Status: "pending"},
				{BookingID: 2, OwnerID: ownerID, Status: "confirmed"},
			}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/1/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetOwnerBookings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []models.Booking
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp))
	assert.Equal(t, "pending", resp[0].Status)
}

func TestHandler_GetOwnerBookings_InvalidID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/invalid/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "invalid"})
	rec := httptest.NewRecorder()

	handler.GetOwnerBookings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID владельца")
}

func TestHandler_GetOwnerBookings_ZeroID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/0/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "0"})
	rec := httptest.NewRecorder()

	handler.GetOwnerBookings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "ID владельца должен быть положительным числом")
}

func TestHandler_GetOwnerBookings_NegativeID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/-1/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "-1"})
	rec := httptest.NewRecorder()

	handler.GetOwnerBookings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_GetOwnerBookings_ServiceError(t *testing.T) {
	mockSvc := &mockBookingService{
		getOwnerBookingsFunc: func(ownerID int) ([]models.Booking, error) {
			return nil, errors.New("database error")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/owners/1/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"owner_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetOwnerBookings(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "database error")
}

func TestHandler_GetSitterBookings_Success(t *testing.T) {
	mockSvc := &mockBookingService{
		getSitterBookingsFunc: func(sitterID int) ([]models.Booking, error) {
			return []models.Booking{
				{BookingID: 1, SitterID: sitterID, Status: "pending"},
				{BookingID: 2, SitterID: sitterID, Status: "confirmed"},
			}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/1/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetSitterBookings(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []models.Booking
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp))
}

func TestHandler_GetSitterBookings_InvalidID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/invalid/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "invalid"})
	rec := httptest.NewRecorder()

	handler.GetSitterBookings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID няни")
}

func TestHandler_GetSitterBookings_ZeroID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/0/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "0"})
	rec := httptest.NewRecorder()

	handler.GetSitterBookings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "ID няни должен быть положительным числом")
}

func TestHandler_GetSitterBookings_NegativeID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/-5/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "-5"})
	rec := httptest.NewRecorder()

	handler.GetSitterBookings(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_GetSitterBookings_ServiceError(t *testing.T) {
	mockSvc := &mockBookingService{
		getSitterBookingsFunc: func(sitterID int) ([]models.Booking, error) {
			return nil, errors.New("database error")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/1/bookings", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetSitterBookings(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "database error")
}

func TestHandler_CreateBooking_ValidationErrors(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	reqBody := map[string]interface{}{
		"owner_id":   1,
		"sitter_id":  1,
		"pet_id":     1,
		"service_id": 1,
		"start_date": "invalid-date",
		"end_date":   "2024-12-31",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateBooking(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CreateBooking_InvalidBody(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/bookings", bytes.NewBuffer([]byte("invalid json")))
	rec := httptest.NewRecorder()

	handler.CreateBooking(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_GetBooking_InvalidID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/bookings/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.GetBooking(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_ConfirmBooking_ZeroID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/bookings/0/confirm", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	rec := httptest.NewRecorder()

	handler.ConfirmBooking(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CancelBooking_InvalidID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/bookings/abc/cancel", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rec := httptest.NewRecorder()

	handler.CancelBooking(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_CompleteBooking_ZeroID(t *testing.T) {
	mockSvc := &mockBookingService{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/bookings/0/complete", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "0"})
	rec := httptest.NewRecorder()

	handler.CompleteBooking(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

type mockBookingService struct {
	createBookingFunc     func(int, int, int, int, time.Time, time.Time) (int, error)
	getBookingByIDFunc    func(int) (*models.Booking, error)
	getOwnerBookingsFunc  func(int) ([]models.Booking, error)
	getSitterBookingsFunc func(int) ([]models.Booking, error)
	confirmBookingFunc    func(int) error
	cancelBookingFunc     func(int) error
	completeBookingFunc   func(int) error
}

func (m *mockBookingService) CreateBooking(ownerID, sitterID, petID, serviceID int, startDate, endDate time.Time) (int, error) {
	if m.createBookingFunc != nil {
		return m.createBookingFunc(ownerID, sitterID, petID, serviceID, startDate, endDate)
	}
	return 1, nil
}

func (m *mockBookingService) GetBookingByID(bookingID int) (*models.Booking, error) {
	if m.getBookingByIDFunc != nil {
		return m.getBookingByIDFunc(bookingID)
	}
	return &models.Booking{BookingID: bookingID}, nil
}

func (m *mockBookingService) GetOwnerBookings(ownerID int) ([]models.Booking, error) {
	if m.getOwnerBookingsFunc != nil {
		return m.getOwnerBookingsFunc(ownerID)
	}
	return []models.Booking{}, nil
}

func (m *mockBookingService) GetSitterBookings(sitterID int) ([]models.Booking, error) {
	if m.getSitterBookingsFunc != nil {
		return m.getSitterBookingsFunc(sitterID)
	}
	return []models.Booking{}, nil
}

func (m *mockBookingService) ConfirmBooking(bookingID int) error {
	if m.confirmBookingFunc != nil {
		return m.confirmBookingFunc(bookingID)
	}
	return nil
}

func (m *mockBookingService) CancelBooking(bookingID int) error {
	if m.cancelBookingFunc != nil {
		return m.cancelBookingFunc(bookingID)
	}
	return nil
}

func (m *mockBookingService) CompleteBooking(bookingID int) error {
	if m.completeBookingFunc != nil {
		return m.completeBookingFunc(bookingID)
	}
	return nil
}
