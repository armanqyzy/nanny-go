package validator

import (
	"testing"
)

func TestValidateKazakhPhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{
			name:    "valid phone number",
			phone:   "+77771234567",
			wantErr: false,
		},
		{
			name:    "valid phone number with different operator",
			phone:   "+77012345678",
			wantErr: false,
		},
		{
			name:    "invalid - no plus sign",
			phone:   "77771234567",
			wantErr: true,
		},
		{
			name:    "invalid - starts with 8",
			phone:   "87771234567",
			wantErr: true,
		},
		{
			name:    "invalid - too short",
			phone:   "+7777123456",
			wantErr: true,
		},
		{
			name:    "invalid - too long",
			phone:   "+777712345678",
			wantErr: true,
		},
		{
			name:    "invalid - has spaces",
			phone:   "+7 777 123 45 67",
			wantErr: true,
		},
		{
			name:    "invalid - has dashes",
			phone:   "+7-777-123-45-67",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Phone string `validate:"phone_kz"`
			}

			data := TestStruct{Phone: tt.phone}
			err := Validate(&data)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateKazakhPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePetType(t *testing.T) {
	tests := []struct {
		name    string
		petType string
		wantErr bool
	}{
		{
			name:    "valid - собака",
			petType: "собака",
			wantErr: false,
		},
		{
			name:    "valid - кошка",
			petType: "кошка",
			wantErr: false,
		},
		{
			name:    "valid - птица",
			petType: "птица",
			wantErr: false,
		},
		{
			name:    "valid - грызун",
			petType: "грызун",
			wantErr: false,
		},
		{
			name:    "valid - рептилия",
			petType: "рептилия",
			wantErr: false,
		},
		{
			name:    "valid - другое",
			petType: "другое",
			wantErr: false,
		},
		{
			name:    "valid - case insensitive",
			petType: "СОБАКА",
			wantErr: false,
		},
		{
			name:    "invalid - unknown type",
			petType: "дракон",
			wantErr: true,
		},
		{
			name:    "invalid - english",
			petType: "dog",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Type string `validate:"pet_type"`
			}

			data := TestStruct{Type: tt.petType}
			err := Validate(&data)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePetType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "user@mail.example.com",
			wantErr: false,
		},
		{
			name:    "valid email with plus",
			email:   "user+tag@example.com",
			wantErr: false,
		},
		{
			name:    "invalid - no @",
			email:   "testexample.com",
			wantErr: true,
		},
		{
			name:    "invalid - no domain",
			email:   "test@",
			wantErr: true,
		},
		{
			name:    "invalid - no local part",
			email:   "@example.com",
			wantErr: true,
		},
		{
			name:    "invalid - spaces",
			email:   "test @example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type TestStruct struct {
				Email string `validate:"email"`
			}

			data := TestStruct{Email: tt.email}
			err := Validate(&data)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRegisterOwnerRequest(t *testing.T) {
	type RegisterOwnerRequest struct {
		FullName string `validate:"required,min=2,max=100"`
		Email    string `validate:"required,email"`
		Phone    string `validate:"required,phone_kz"`
		Password string `validate:"required,min=8,max=72"`
	}

	tests := []struct {
		name    string
		request RegisterOwnerRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: RegisterOwnerRequest{
				FullName: "Иван Иванов",
				Email:    "ivan@example.com",
				Phone:    "+77771234567",
				Password: "securepassword123",
			},
			wantErr: false,
		},
		{
			name: "invalid - missing full name",
			request: RegisterOwnerRequest{
				Email:    "ivan@example.com",
				Phone:    "+77771234567",
				Password: "securepassword123",
			},
			wantErr: true,
		},
		{
			name: "invalid - short full name",
			request: RegisterOwnerRequest{
				FullName: "И",
				Email:    "ivan@example.com",
				Phone:    "+77771234567",
				Password: "securepassword123",
			},
			wantErr: true,
		},
		{
			name: "invalid - bad email",
			request: RegisterOwnerRequest{
				FullName: "Иван Иванов",
				Email:    "not-an-email",
				Phone:    "+77771234567",
				Password: "securepassword123",
			},
			wantErr: true,
		},
		{
			name: "invalid - bad phone",
			request: RegisterOwnerRequest{
				FullName: "Иван Иванов",
				Email:    "ivan@example.com",
				Phone:    "87771234567",
				Password: "securepassword123",
			},
			wantErr: true,
		},
		{
			name: "invalid - short password",
			request: RegisterOwnerRequest{
				FullName: "Иван Иванов",
				Email:    "ivan@example.com",
				Phone:    "+77771234567",
				Password: "short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCreatePetRequest(t *testing.T) {
	type CreatePetRequest struct {
		OwnerID int    `validate:"required,gt=0"`
		Name    string `validate:"required,min=1,max=100"`
		Type    string `validate:"required,pet_type"`
		Age     int    `validate:"required,gte=0,lte=30"`
		Notes   string `validate:"max=500"`
	}

	tests := []struct {
		name    string
		request CreatePetRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: CreatePetRequest{
				OwnerID: 1,
				Name:    "Барсик",
				Type:    "кошка",
				Age:     3,
				Notes:   "Очень дружелюбный",
			},
			wantErr: false,
		},
		{
			name: "valid request without notes",
			request: CreatePetRequest{
				OwnerID: 1,
				Name:    "Рекс",
				Type:    "собака",
				Age:     5,
			},
			wantErr: false,
		},
		{
			name: "invalid - zero owner id",
			request: CreatePetRequest{
				OwnerID: 0,
				Name:    "Барсик",
				Type:    "кошка",
				Age:     3,
			},
			wantErr: true,
		},
		{
			name: "invalid - negative age",
			request: CreatePetRequest{
				OwnerID: 1,
				Name:    "Барсик",
				Type:    "кошка",
				Age:     -1,
			},
			wantErr: true,
		},
		{
			name: "invalid - too old",
			request: CreatePetRequest{
				OwnerID: 1,
				Name:    "Барсик",
				Type:    "кошка",
				Age:     31,
			},
			wantErr: true,
		},
		{
			name: "invalid - wrong pet type",
			request: CreatePetRequest{
				OwnerID: 1,
				Name:    "Барсик",
				Type:    "динозавр",
				Age:     3,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateCreateReviewRequest(t *testing.T) {
	type CreateReviewRequest struct {
		BookingID int    `validate:"required,gt=0"`
		OwnerID   int    `validate:"required,gt=0"`
		SitterID  int    `validate:"required,gt=0"`
		Rating    int    `validate:"required,gte=1,lte=5"`
		Comment   string `validate:"max=1000"`
	}

	tests := []struct {
		name    string
		request CreateReviewRequest
		wantErr bool
	}{
		{
			name: "valid request with comment",
			request: CreateReviewRequest{
				BookingID: 1,
				OwnerID:   1,
				SitterID:  1,
				Rating:    5,
				Comment:   "Отличная работа!",
			},
			wantErr: false,
		},
		{
			name: "valid request without comment",
			request: CreateReviewRequest{
				BookingID: 1,
				OwnerID:   1,
				SitterID:  1,
				Rating:    4,
			},
			wantErr: false,
		},
		{
			name: "invalid - rating too low",
			request: CreateReviewRequest{
				BookingID: 1,
				OwnerID:   1,
				SitterID:  1,
				Rating:    0,
			},
			wantErr: true,
		},
		{
			name: "invalid - rating too high",
			request: CreateReviewRequest{
				BookingID: 1,
				OwnerID:   1,
				SitterID:  1,
				Rating:    6,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero booking id",
			request: CreateReviewRequest{
				BookingID: 0,
				OwnerID:   1,
				SitterID:  1,
				Rating:    5,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
