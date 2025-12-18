package reviews

import (
	"database/sql"
	"testing"
	"time"

	"nanny-backend/internal/common/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	review := &models.Review{
		BookingID: 1,
		OwnerID:   2,
		SitterID:  3,
		Rating:    5,
		Comment:   "Great service",
	}

	mock.ExpectQuery(`INSERT INTO reviews`).
		WithArgs(
			review.BookingID,
			review.OwnerID,
			review.SitterID,
			review.Rating,
			review.Comment,
		).
		WillReturnRows(sqlmock.NewRows([]string{"review_id"}).AddRow(10))

	id, err := repo.Create(review)

	assert.NoError(t, err)
	assert.Equal(t, 10, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"review_id",
		"booking_id",
		"owner_id",
		"sitter_id",
		"rating",
		"comment",
		"created_at",
	}).AddRow(
		10,
		1,
		2,
		3,
		5,
		"Great service",
		now,
	)

	mock.ExpectQuery(`FROM reviews WHERE review_id = \$1`).
		WithArgs(10).
		WillReturnRows(rows)

	review, err := repo.GetByID(10)

	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, 5, review.Rating)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	mock.ExpectQuery(`FROM reviews WHERE review_id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	review, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, review)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBySitterID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"review_id",
		"booking_id",
		"owner_id",
		"sitter_id",
		"rating",
		"comment",
		"created_at",
	}).
		AddRow(1, 10, 2, 5, 4, "Good", now).
		AddRow(2, 11, 3, 5, 5, "Excellent", now.Add(time.Minute))

	mock.ExpectQuery(`FROM reviews WHERE sitter_id = \$1`).
		WithArgs(5).
		WillReturnRows(rows)

	reviews, err := repo.GetBySitterID(5)

	assert.NoError(t, err)
	assert.Len(t, reviews, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByBookingID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"review_id",
		"booking_id",
		"owner_id",
		"sitter_id",
		"rating",
		"comment",
		"created_at",
	}).AddRow(
		10,
		1,
		2,
		3,
		5,
		"Nice",
		now,
	)

	mock.ExpectQuery(`FROM reviews WHERE booking_id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	review, err := repo.GetByBookingID(1)

	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, "Nice", review.Comment)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	review := &models.Review{
		ReviewID: 10,
		Rating:   4,
		Comment:  "Updated comment",
	}

	mock.ExpectExec(`UPDATE reviews`).
		WithArgs(
			review.Rating,
			review.Comment,
			review.ReviewID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(review)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	mock.ExpectExec(`DELETE FROM reviews`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(10)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSitterRating_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{
		"avg",
		"count",
	}).AddRow(4.5, 6)

	mock.ExpectQuery(`SELECT AVG\(rating\), COUNT\(\*\)`).
		WithArgs(3).
		WillReturnRows(rows)

	avg, count, err := repo.GetSitterRating(3)

	assert.NoError(t, err)
	assert.Equal(t, 4.5, avg)
	assert.Equal(t, 6, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSitterRating_NoReviews(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{
		"avg",
		"count",
	}).AddRow(nil, 0)

	mock.ExpectQuery(`SELECT AVG\(rating\), COUNT\(\*\)`).
		WithArgs(3).
		WillReturnRows(rows)

	avg, count, err := repo.GetSitterRating(3)

	assert.NoError(t, err)
	assert.Equal(t, 0.0, avg)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
