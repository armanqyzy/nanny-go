package reviews

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"nanny-backend/internal/common/models"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(review *models.Review) (int, error) {
	args := m.Called(review)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetByID(reviewID int) (*models.Review, error) {
	args := m.Called(reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Review), args.Error(1)
}

func (m *MockRepository) GetBySitterID(sitterID int) ([]models.Review, error) {
	args := m.Called(sitterID)
	return args.Get(0).([]models.Review), args.Error(1)
}

func (m *MockRepository) GetByBookingID(bookingID int) (*models.Review, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Review), args.Error(1)
}

func (m *MockRepository) Update(review *models.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockRepository) Delete(reviewID int) error {
	args := m.Called(reviewID)
	return args.Error(0)
}

func (m *MockRepository) GetSitterRating(sitterID int) (float64, int, error) {
	args := m.Called(sitterID)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func TestCreateReview_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("GetByBookingID", 1).
		Return(nil, errors.New("not found"))

	mockRepo.
		On("Create", mock.Anything).
		Return(10, nil)

	reviewID, err := service.CreateReview(
		1,
		2,
		3,
		5,
		"Great service",
	)

	assert.NoError(t, err)
	assert.Equal(t, 10, reviewID)
	mockRepo.AssertExpectations(t)
}

func TestCreateReview_InvalidRating(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	reviewID, err := service.CreateReview(
		1,
		2,
		3,
		6, // invalid rating
		"Bad",
	)

	assert.Error(t, err)
	assert.Equal(t, 0, reviewID)
}

func TestCreateReview_AlreadyExists(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existing := &models.Review{
		ReviewID:  1,
		BookingID: 1,
	}

	mockRepo.
		On("GetByBookingID", 1).
		Return(existing, nil)

	reviewID, err := service.CreateReview(
		1,
		2,
		3,
		5,
		"Duplicate",
	)

	assert.Error(t, err)
	assert.Equal(t, 0, reviewID)
}

func TestGetReview_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := &models.Review{
		ReviewID:  1,
		BookingID: 1,
		OwnerID:   2,
		SitterID:  3,
		Rating:    4,
		Comment:   "Nice",
		CreatedAt: time.Now(),
	}

	mockRepo.
		On("GetByID", 1).
		Return(expected, nil)

	review, err := service.GetReview(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, review)
	mockRepo.AssertExpectations(t)
}

func TestGetSitterReviews_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := []models.Review{
		{ReviewID: 1, SitterID: 3, Rating: 5},
		{ReviewID: 2, SitterID: 3, Rating: 4},
	}

	mockRepo.
		On("GetBySitterID", 3).
		Return(expected, nil)

	reviews, err := service.GetSitterReviews(3)

	assert.NoError(t, err)
	assert.Len(t, reviews, 2)
	assert.Equal(t, expected, reviews)
	mockRepo.AssertExpectations(t)
}

func TestGetBookingReview_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	expected := &models.Review{
		ReviewID:  1,
		BookingID: 5,
	}

	mockRepo.
		On("GetByBookingID", 5).
		Return(expected, nil)

	review, err := service.GetBookingReview(5)

	assert.NoError(t, err)
	assert.Equal(t, expected, review)
	mockRepo.AssertExpectations(t)
}

func TestUpdateReview_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	existing := &models.Review{
		ReviewID: 1,
		Rating:   3,
		Comment:  "Ok",
	}

	mockRepo.
		On("GetByID", 1).
		Return(existing, nil)

	mockRepo.
		On("Update", mock.Anything).
		Return(nil)

	err := service.UpdateReview(1, 5, "Excellent")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateReview_InvalidRating(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	err := service.UpdateReview(1, 0, "Bad")

	assert.Error(t, err)
}

func TestDeleteReview_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("Delete", 1).
		Return(nil)

	err := service.DeleteReview(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetSitterRating_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("GetSitterRating", 3).
		Return(4.5, 10, nil)

	rating, count, err := service.GetSitterRating(3)

	assert.NoError(t, err)
	assert.Equal(t, 4.5, rating)
	assert.Equal(t, 10, count)
	mockRepo.AssertExpectations(t)
}
