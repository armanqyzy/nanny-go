package bookings

import (
	"context"
	"database/sql"
	"fmt"
	"nanny-backend/internal/common/models"
	"time"
)

type Repository interface {
	Create(booking *models.Booking) (int, error)
	GetByID(bookingID int) (*models.Booking, error)
	GetByOwnerID(ownerID int) ([]models.Booking, error)
	GetBySitterID(sitterID int) ([]models.Booking, error)
	UpdateStatus(bookingID int, status string) error
	Delete(bookingID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(booking *models.Booking) (int, error) {
	var bookingID int
	err := r.db.QueryRow(`
		INSERT INTO bookings (owner_id, sitter_id, pet_id, service_id, start_time, end_time, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING booking_id
	`, booking.OwnerID, booking.SitterID, booking.PetID, booking.ServiceID,
		booking.StartTime, booking.EndTime, booking.Status).Scan(&bookingID)

	if err != nil {
		return 0, fmt.Errorf("could not create booking: %w", err)
	}

	return bookingID, nil
}

func (r *repository) GetByID(bookingID int) (*models.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	booking := &models.Booking{}
	err := r.db.QueryRowContext(ctx, `
		SELECT booking_id, owner_id, sitter_id, pet_id, service_id, start_time, end_time, status
		FROM bookings
		WHERE booking_id = $1
	`, bookingID).Scan(
		&booking.BookingID,
		&booking.OwnerID,
		&booking.SitterID,
		&booking.PetID,
		&booking.ServiceID,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("бронирование не найдено")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения бронирования: %w", err)
	}

	return booking, nil
}

func (r *repository) GetByOwnerID(ownerID int) ([]models.Booking, error) {
	rows, err := r.db.Query(`
		SELECT booking_id, owner_id, sitter_id, pet_id, service_id, start_time, end_time, status
		FROM bookings
		WHERE owner_id = $1
		ORDER BY start_time DESC
	`, ownerID)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения бронирований: %w", err)
	}
	defer rows.Close()

	return scanBookings(rows)
}

func (r *repository) GetBySitterID(sitterID int) ([]models.Booking, error) {
	rows, err := r.db.Query(`
		SELECT booking_id, owner_id, sitter_id, pet_id, service_id, start_time, end_time, status
		FROM bookings
		WHERE sitter_id = $1
		ORDER BY start_time DESC
	`, sitterID)

	if err != nil {
		return nil, fmt.Errorf("ошибка получения бронирований: %w", err)
	}
	defer rows.Close()

	return scanBookings(rows)
}

func (r *repository) UpdateStatus(bookingID int, status string) error {
	_, err := r.db.Exec(`
		UPDATE bookings
		SET status = $1
		WHERE booking_id = $2
	`, status, bookingID)

	if err != nil {
		return fmt.Errorf("не удалось обновить статус бронирования: %w", err)
	}

	return nil
}

func (r *repository) Delete(bookingID int) error {
	_, err := r.db.Exec(`DELETE FROM bookings WHERE booking_id = $1`, bookingID)
	if err != nil {
		return fmt.Errorf("не удалось удалить бронирование: %w", err)
	}
	return nil
}

func scanBookings(rows *sql.Rows) ([]models.Booking, error) {
	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		var startTime, endTime time.Time

		err := rows.Scan(
			&booking.BookingID,
			&booking.OwnerID,
			&booking.SitterID,
			&booking.PetID,
			&booking.ServiceID,
			&startTime,
			&endTime,
			&booking.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования бронирования: %w", err)
		}

		booking.StartTime = startTime
		booking.EndTime = endTime
		bookings = append(bookings, booking)
	}

	return bookings, nil
}
