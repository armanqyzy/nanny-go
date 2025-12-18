package bookings

import (
	"errors"
	"testing"
	"time"

	"nanny-backend/internal/common/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(booking *models.Booking) (int, error) {
	args := m.Called(booking)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetByID(bookingID int) (*models.Booking, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *MockRepository) GetByOwnerID(ownerID int) ([]models.Booking, error) {
	args := m.Called(ownerID)
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockRepository) GetBySitterID(sitterID int) ([]models.Booking, error) {
	args := m.Called(sitterID)
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *MockRepository) UpdateStatus(bookingID int, status string) error {
	args := m.Called(bookingID, status)
	return args.Error(0)
}
func (m *MockRepository) Delete(bookingID int) error {
	args := m.Called(bookingID)
	return args.Error(0)
}

func TestCreateBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	mockRepo.On("Create", mock.MatchedBy(func(b *models.Booking) bool {
		return b.OwnerID == 1 &&
			b.SitterID == 2 &&
			b.PetID == 3 &&
			b.ServiceID == 4 &&
			b.Status == "pending"
	})).Return(42, nil)

	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	assert.NoError(t, err)
	assert.Equal(t, 42, bookingID)
	mockRepo.AssertExpectations(t)
}

func TestCreateBooking_EndTimeBeforeStartTime(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(-1 * time.Hour)

	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	assert.Error(t, err)
	assert.Equal(t, 0, bookingID)
	assert.Contains(t, err.Error(), "start date cannot be later")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateBooking_StartTimeInPast(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(1 * time.Hour)

	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	assert.Error(t, err)
	assert.Equal(t, 0, bookingID)
	assert.Contains(t, err.Error(), "cannot create a booking in the past time")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateBooking_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	mockRepo.On("Create", mock.Anything).Return(0, errors.New("database error"))

	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	assert.Error(t, err)
	assert.Equal(t, 0, bookingID)
	assert.Contains(t, err.Error(), "error creating booking")
	mockRepo.AssertExpectations(t)
}

func TestGetBookingByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedBooking := &models.Booking{
		BookingID: 1,
		OwnerID:   2,
		SitterID:  3,
		Status:    "confirmed",
	}

	mockRepo.On("GetByID", 1).Return(expectedBooking, nil)

	booking, err := service.GetBookingByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedBooking, booking)
	mockRepo.AssertExpectations(t)
}

func TestGetBookingByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("GetByID", 999).Return((*models.Booking)(nil), errors.New("booking not found"))

	booking, err := service.GetBookingByID(999)

	assert.Error(t, err)
	assert.Nil(t, booking)
	mockRepo.AssertExpectations(t)
}

func TestConfirmBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "pending",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)
	mockRepo.On("UpdateStatus", 1, "confirmed").Return(nil)

	err := service.ConfirmBooking(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConfirmBooking_InvalidStatus(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "confirmed",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)

	err := service.ConfirmBooking(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can approve only booking with status 'pending'")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

func TestCancelBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "pending",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)
	mockRepo.On("UpdateStatus", 1, "cancelled").Return(nil)

	err := service.CancelBooking(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCancelBooking_CompletedBooking(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "completed",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)

	err := service.CancelBooking(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot cancel completed booking")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

func TestCompleteBooking_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "confirmed",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)
	mockRepo.On("UpdateStatus", 1, "completed").Return(nil)

	err := service.CompleteBooking(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCompleteBooking_NotConfirmed(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "pending",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)

	err := service.CompleteBooking(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only finish accepted booking")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

func TestGetOwnerBookings_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedBookings := []models.Booking{
		{BookingID: 1, OwnerID: 5},
		{BookingID: 2, OwnerID: 5},
	}

	mockRepo.On("GetByOwnerID", 5).Return(expectedBookings, nil)

	bookings, err := service.GetOwnerBookings(5)

	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, expectedBookings, bookings)
	mockRepo.AssertExpectations(t)
}

func TestGetSitterBookings_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedBookings := []models.Booking{
		{BookingID: 3, SitterID: 10},
		{BookingID: 4, SitterID: 10},
		{BookingID: 5, SitterID: 10},
	}

	mockRepo.On("GetBySitterID", 10).Return(expectedBookings, nil)

	bookings, err := service.GetSitterBookings(10)

	assert.NoError(t, err)
	assert.Len(t, bookings, 3)
	assert.Equal(t, expectedBookings, bookings)
	mockRepo.AssertExpectations(t)
}
