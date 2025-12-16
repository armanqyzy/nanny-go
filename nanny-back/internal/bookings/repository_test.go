package bookings

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"nanny-backend/internal/common/models"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(2 * time.Hour)

	booking := &models.Booking{
		OwnerID:   1,
		SitterID:  2,
		PetID:     3,
		ServiceID: 4,
		StartTime: start,
		EndTime:   end,
		Status:    "pending",
	}

	mock.ExpectQuery(`INSERT INTO bookings`).
		WithArgs(
			booking.OwnerID,
			booking.SitterID,
			booking.PetID,
			booking.ServiceID,
			booking.StartTime,
			booking.EndTime,
			booking.Status,
		).
		WillReturnRows(sqlmock.NewRows([]string{"booking_id"}).AddRow(10))

	id, err := repo.Create(booking)

	assert.NoError(t, err)
	assert.Equal(t, 10, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	start := time.Now()
	end := start.Add(2 * time.Hour)

	rows := sqlmock.NewRows([]string{
		"booking_id",
		"owner_id",
		"sitter_id",
		"pet_id",
		"service_id",
		"start_time",
		"end_time",
		"status",
	}).AddRow(
		10,
		1,
		2,
		3,
		4,
		start,
		end,
		"confirmed",
	)

	mock.ExpectQuery(`FROM bookings WHERE booking_id = \$1`).
		WithArgs(10).
		WillReturnRows(rows)

	booking, err := repo.GetByID(10)

	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, 10, booking.BookingID)
	assert.Equal(t, "confirmed", booking.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(`FROM bookings WHERE booking_id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	booking, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, booking)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByOwnerID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"booking_id",
		"owner_id",
		"sitter_id",
		"pet_id",
		"service_id",
		"start_time",
		"end_time",
		"status",
	}).
		AddRow(1, 5, 10, 3, 4, now, now.Add(1*time.Hour), "pending").
		AddRow(2, 5, 11, 4, 5, now, now.Add(2*time.Hour), "confirmed")

	mock.ExpectQuery(`FROM bookings WHERE owner_id = \$1`).
		WithArgs(5).
		WillReturnRows(rows)

	bookings, err := repo.GetByOwnerID(5)

	assert.NoError(t, err)
	assert.Len(t, bookings, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBySitterID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"booking_id",
		"owner_id",
		"sitter_id",
		"pet_id",
		"service_id",
		"start_time",
		"end_time",
		"status",
	}).AddRow(
		1, 2, 7, 3, 4, now, now.Add(time.Hour), "completed",
	)

	mock.ExpectQuery(`FROM bookings WHERE sitter_id = \$1`).
		WithArgs(7).
		WillReturnRows(rows)

	bookings, err := repo.GetBySitterID(7)

	assert.NoError(t, err)
	assert.Len(t, bookings, 1)
	assert.Equal(t, "completed", bookings[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectExec(`UPDATE bookings`).
		WithArgs("confirmed", 10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateStatus(10, "confirmed")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectExec(`DELETE FROM bookings WHERE booking_id = \$1`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(10)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
