package reviews

import (
	"database/sql"
	"fmt"

	"nanny-backend/internal/common/models"
)

type Repository interface {
	Create(review *models.Review) (int, error)
	GetByID(reviewID int) (*models.Review, error)
	GetBySitterID(sitterID int) ([]models.Review, error)
	GetByBookingID(bookingID int) (*models.Review, error)
	Update(review *models.Review) error
	Delete(reviewID int) error
	GetSitterRating(sitterID int) (float64, int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(review *models.Review) (int, error) {
	var reviewID int
	err := r.db.QueryRow(`
		INSERT INTO reviews (booking_id, owner_id, sitter_id, rating, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING review_id
	`, review.BookingID, review.OwnerID, review.SitterID, review.Rating, review.Comment).Scan(&reviewID)

	if err != nil {
		return 0, fmt.Errorf("не удалось создать отзыв: %w", err)
	}

	return reviewID, nil
}

func (r *repository) GetByID(reviewID int) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.QueryRow(`
		SELECT review_id, booking_id, owner_id, sitter_id, rating, comment, created_at
		FROM reviews
		WHERE review_id = $1
	`, reviewID).Scan(
		&review.ReviewID,
		&review.BookingID,
		&review.OwnerID,
		&review.SitterID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("отзыв не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения отзыва: %w", err)
	}

	return review, nil
}

func (r *repository) GetBySitterID(sitterID int) ([]models.Review, error) {
	rows, err := r.db.Query(`
		SELECT review_id, booking_id, owner_id, sitter_id, rating, comment, created_at
		FROM reviews
		WHERE sitter_id = $1
		ORDER BY created_at DESC
	`, sitterID)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения отзывов: %w", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(
			&review.ReviewID,
			&review.BookingID,
			&review.OwnerID,
			&review.SitterID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования отзыва: %w", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *repository) GetByBookingID(bookingID int) (*models.Review, error) {
	review := &models.Review{}
	err := r.db.QueryRow(`
		SELECT review_id, booking_id, owner_id, sitter_id, rating, comment, created_at
		FROM reviews
		WHERE booking_id = $1
	`, bookingID).Scan(
		&review.ReviewID,
		&review.BookingID,
		&review.OwnerID,
		&review.SitterID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("отзыв не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения отзыва: %w", err)
	}

	return review, nil
}

func (r *repository) Update(review *models.Review) error {
	_, err := r.db.Exec(`
		UPDATE reviews
		SET rating = $1, comment = $2
		WHERE review_id = $3
	`, review.Rating, review.Comment, review.ReviewID)

	if err != nil {
		return fmt.Errorf("не удалось обновить отзыв: %w", err)
	}

	return nil
}

func (r *repository) Delete(reviewID int) error {
	_, err := r.db.Exec(`DELETE FROM reviews WHERE review_id = $1`, reviewID)
	if err != nil {
		return fmt.Errorf("не удалось удалить отзыв: %w", err)
	}
	return nil
}

func (r *repository) GetSitterRating(sitterID int) (float64, int, error) {
	var avgRating sql.NullFloat64
	var count int

	err := r.db.QueryRow(`
		SELECT AVG(rating), COUNT(*)
		FROM reviews
		WHERE sitter_id = $1
	`, sitterID).Scan(&avgRating, &count)

	if err != nil {
		return 0, 0, fmt.Errorf("ошибка расчёта рейтинга: %w", err)
	}

	if !avgRating.Valid {
		return 0, 0, nil
	}

	return avgRating.Float64, count, nil
}
