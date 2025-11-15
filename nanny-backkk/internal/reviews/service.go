package reviews

import (
	"fmt"

	"nanny-backend/internal/common/models"
)

type Service interface {
	CreateReview(bookingID, ownerID, sitterID, rating int, comment string) (int, error)
	GetReview(reviewID int) (*models.Review, error)
	GetSitterReviews(sitterID int) ([]models.Review, error)
	GetBookingReview(bookingID int) (*models.Review, error)
	UpdateReview(reviewID, rating int, comment string) error
	DeleteReview(reviewID int) error
	GetSitterRating(sitterID int) (float64, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateReview(bookingID, ownerID, sitterID, rating int, comment string) (int, error) {
	// Валидация рейтинга
	if rating < 1 || rating > 5 {
		return 0, fmt.Errorf("рейтинг должен быть от 1 до 5")
	}

	// Проверяем, нет ли уже отзыва на это бронирование
	existing, _ := s.repo.GetByBookingID(bookingID)
	if existing != nil {
		return 0, fmt.Errorf("отзыв на это бронирование уже существует")
	}

	review := &models.Review{
		BookingID: bookingID,
		OwnerID:   ownerID,
		SitterID:  sitterID,
		Rating:    rating,
		Comment:   comment,
	}

	reviewID, err := s.repo.Create(review)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания отзыва: %w", err)
	}

	return reviewID, nil
}

func (s *service) GetReview(reviewID int) (*models.Review, error) {
	return s.repo.GetByID(reviewID)
}

func (s *service) GetSitterReviews(sitterID int) ([]models.Review, error) {
	return s.repo.GetBySitterID(sitterID)
}

func (s *service) GetBookingReview(bookingID int) (*models.Review, error) {
	return s.repo.GetByBookingID(bookingID)
}

func (s *service) UpdateReview(reviewID, rating int, comment string) error {
	// Валидация рейтинга
	if rating < 1 || rating > 5 {
		return fmt.Errorf("рейтинг должен быть от 1 до 5")
	}

	// Проверяем существование отзыва
	review, err := s.repo.GetByID(reviewID)
	if err != nil {
		return err
	}

	review.Rating = rating
	review.Comment = comment

	return s.repo.Update(review)
}

func (s *service) DeleteReview(reviewID int) error {
	return s.repo.Delete(reviewID)
}

func (s *service) GetSitterRating(sitterID int) (float64, int, error) {
	return s.repo.GetSitterRating(sitterID)
}
