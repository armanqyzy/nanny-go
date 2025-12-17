package models

import "time"

type User struct {
	UserID       int       `json:"user_id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Pet struct {
	PetID   int    `json:"pet_id"`
	OwnerID int    `json:"owner_id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Age     int    `json:"age"`
	Notes   string `json:"notes,omitempty"`
}

type Sitter struct {
	SitterID        int    `json:"sitter_id"`
	ExperienceYears int    `json:"experience_years"`
	Certificates    string `json:"certificates,omitempty"`
	Preferences     string `json:"preferences,omitempty"`
	Location        string `json:"location"`
	Status          string `json:"status"`
}

type Service struct {
	ServiceID    int     `json:"service_id"`
	SitterID     int     `json:"sitter_id"`
	Type         string  `json:"type"`
	PricePerHour float64 `json:"price_per_hour"`
	Description  string  `json:"description,omitempty"`
}

type Booking struct {
	BookingID int       `json:"booking_id"`
	OwnerID   int       `json:"owner_id"`
	SitterID  int       `json:"sitter_id"`
	PetID     int       `json:"pet_id"`
	ServiceID int       `json:"service_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
}

type Payment struct {
	PaymentID int       `json:"payment_id"`
	BookingID int       `json:"booking_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ReviewID  int       `json:"review_id"`
	BookingID int       `json:"booking_id"`
	OwnerID   int       `json:"owner_id"`
	SitterID  int       `json:"sitter_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	MessageID int       `json:"message_id"`
	ChatID    int       `json:"chat_id"`
	SenderID  int       `json:"sender_id"`
	Content   string    `json:"content"`
	SentAt    time.Time `json:"sent_at"`
}

type Chat struct {
	ChatID    int       `json:"chat_id"`
	BookingID int       `json:"booking_id"`
	CreatedAt time.Time `json:"created_at"`
}
