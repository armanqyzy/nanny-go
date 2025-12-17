package bookings

import (
	"errors"
	"testing"
	"time"

	"nanny-backend/internal/common/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of Repository interface
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

// TestCreateBooking_Success tests successful booking creation
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

// TestCreateBooking_EndTimeBeforeStartTime tests validation for invalid time range
func TestCreateBooking_EndTimeBeforeStartTime(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(-1 * time.Hour) // Invalid: end before start

	// Act
	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, bookingID)
	assert.Contains(t, err.Error(), "время начала не может быть позже")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreateBooking_StartTimeInPast tests validation for past booking
func TestCreateBooking_StartTimeInPast(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(-1 * time.Hour) // In the past
	endTime := time.Now().Add(1 * time.Hour)

	// Act
	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, bookingID)
	assert.Contains(t, err.Error(), "нельзя создать бронирование в прошлом")
	mockRepo.AssertNotCalled(t, "Create")
}

// TestCreateBooking_RepositoryError tests error handling from repository
func TestCreateBooking_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	startTime := time.Now().Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	mockRepo.On("Create", mock.Anything).Return(0, errors.New("database error"))

	// Act
	bookingID, err := service.CreateBooking(1, 2, 3, 4, startTime, endTime)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, 0, bookingID)
	assert.Contains(t, err.Error(), "ошибка создания бронирования")
	mockRepo.AssertExpectations(t)
}

// TestGetBookingByID_Success tests successful retrieval of booking
func TestGetBookingByID_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedBooking := &models.Booking{
		BookingID: 1,
		OwnerID:   2,
		SitterID:  3,
		Status:    "confirmed",
	}

	mockRepo.On("GetByID", 1).Return(expectedBooking, nil)

	// Act
	booking, err := service.GetBookingByID(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedBooking, booking)
	mockRepo.AssertExpectations(t)
}

// TestGetBookingByID_NotFound tests error when booking not found
func TestGetBookingByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("GetByID", 999).Return((*models.Booking)(nil), errors.New("booking not found"))

	// Act
	booking, err := service.GetBookingByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, booking)
	mockRepo.AssertExpectations(t)
}

// TestConfirmBooking_Success tests successful booking confirmation
func TestConfirmBooking_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "pending",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)
	mockRepo.On("UpdateStatus", 1, "confirmed").Return(nil)

	// Act
	err := service.ConfirmBooking(1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestConfirmBooking_InvalidStatus tests confirmation of non-pending booking
func TestConfirmBooking_InvalidStatus(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "confirmed", // Already confirmed
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)

	// Act
	err := service.ConfirmBooking(1)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "можно подтвердить только бронирование со статусом 'pending'")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

// TestCancelBooking_Success tests successful booking cancellation
func TestCancelBooking_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "pending",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)
	mockRepo.On("UpdateStatus", 1, "cancelled").Return(nil)

	// Act
	err := service.CancelBooking(1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCancelBooking_CompletedBooking tests cancellation of completed booking
func TestCancelBooking_CompletedBooking(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "completed",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)

	// Act
	err := service.CancelBooking(1)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "нельзя отменить завершённое бронирование")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

// TestCompleteBooking_Success tests successful booking completion
func TestCompleteBooking_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "confirmed",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)
	mockRepo.On("UpdateStatus", 1, "completed").Return(nil)

	// Act
	err := service.CompleteBooking(1)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCompleteBooking_NotConfirmed tests completion of non-confirmed booking
func TestCompleteBooking_NotConfirmed(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existingBooking := &models.Booking{
		BookingID: 1,
		Status:    "pending",
	}

	mockRepo.On("GetByID", 1).Return(existingBooking, nil)

	// Act
	err := service.CompleteBooking(1)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "можно завершить только подтверждённое бронирование")
	mockRepo.AssertNotCalled(t, "UpdateStatus")
}

// TestGetOwnerBookings_Success tests retrieval of owner's bookings
func TestGetOwnerBookings_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedBookings := []models.Booking{
		{BookingID: 1, OwnerID: 5},
		{BookingID: 2, OwnerID: 5},
	}

	mockRepo.On("GetByOwnerID", 5).Return(expectedBookings, nil)

	// Act
	bookings, err := service.GetOwnerBookings(5)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.Equal(t, expectedBookings, bookings)
	mockRepo.AssertExpectations(t)
}

// TestGetSitterBookings_Success tests retrieval of sitter's bookings
func TestGetSitterBookings_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expectedBookings := []models.Booking{
		{BookingID: 3, SitterID: 10},
		{BookingID: 4, SitterID: 10},
		{BookingID: 5, SitterID: 10},
	}

	mockRepo.On("GetBySitterID", 10).Return(expectedBookings, nil)

	// Act
	bookings, err := service.GetSitterBookings(10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, bookings, 3)
	assert.Equal(t, expectedBookings, bookings)
	mockRepo.AssertExpectations(t)
}
