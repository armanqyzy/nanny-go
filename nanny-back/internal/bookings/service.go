package bookings

import (
	"fmt"
	"time"

	"nanny-backend/internal/common/models"
)

type Service interface {
	CreateBooking(ownerID, sitterID, petID, serviceID int, startTime, endTime time.Time) (int, error)
	GetBookingByID(bookingID int) (*models.Booking, error)
	GetOwnerBookings(ownerID int) ([]models.Booking, error)
	GetSitterBookings(sitterID int) ([]models.Booking, error)
	ConfirmBooking(bookingID int) error
	CancelBooking(bookingID int) error
	CompleteBooking(bookingID int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateBooking(ownerID, sitterID, petID, serviceID int, startTime, endTime time.Time) (int, error) {

	if startTime.After(endTime) {
		return 0, fmt.Errorf("время начала не может быть позже времени окончания")
	}

	if startTime.Before(time.Now()) {
		return 0, fmt.Errorf("нельзя создать бронирование в прошлом")
	}

	booking := &models.Booking{
		OwnerID:   ownerID,
		SitterID:  sitterID,
		PetID:     petID,
		ServiceID: serviceID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    "pending",
	}

	bookingID, err := s.repo.Create(booking)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания бронирования: %w", err)
	}

	return bookingID, nil
}

func (s *service) GetBookingByID(bookingID int) (*models.Booking, error) {
	return s.repo.GetByID(bookingID)
}

func (s *service) GetOwnerBookings(ownerID int) ([]models.Booking, error) {
	return s.repo.GetByOwnerID(ownerID)
}

func (s *service) GetSitterBookings(sitterID int) ([]models.Booking, error) {
	return s.repo.GetBySitterID(sitterID)
}

func (s *service) ConfirmBooking(bookingID int) error {

	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != "pending" {
		return fmt.Errorf("можно подтвердить только бронирование со статусом 'pending'")
	}

	return s.repo.UpdateStatus(bookingID, "confirmed")
}

func (s *service) CancelBooking(bookingID int) error {
	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status == "completed" {
		return fmt.Errorf("нельзя отменить завершённое бронирование")
	}

	return s.repo.UpdateStatus(bookingID, "cancelled")
}

func (s *service) CompleteBooking(bookingID int) error {
	booking, err := s.repo.GetByID(bookingID)
	if err != nil {
		return err
	}

	if booking.Status != "confirmed" {
		return fmt.Errorf("можно завершить только подтверждённое бронирование")
	}

	return s.repo.UpdateStatus(bookingID, "completed")
}
