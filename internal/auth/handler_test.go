package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"nanny-backend/internal/common/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) RegisterOwner(fullName, email, phone, password string) error {
	args := m.Called(fullName, email, phone, password)
	return args.Error(0)
}

func (m *MockService) RegisterSitter(
	fullName, email, phone, password string,
	experienceYears int,
	certificates, preferences, location string,
) error {
	args := m.Called(
		fullName,
		email,
		phone,
		password,
		experienceYears,
		certificates,
		preferences,
		location,
	)
	return args.Error(0)
}

func (m *MockService) Login(email, password string) (*models.User, string, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*models.User), args.String(1), args.Error(2)
}

func TestHandler_RegisterOwner_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := RegisterOwnerRequest{
		FullName: "Nuray Alim",
		Email:    "nuray@test.com",
		Phone:    "+77001234567",
		Password: "strongpassword",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register/owner", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On(
			"RegisterOwner",
			reqBody.FullName,
			reqBody.Email,
			reqBody.Phone,
			reqBody.Password,
		).
		Return(nil)

	handler.RegisterOwner(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "владелец зарегистрирован успешно", resp["message"])

	mockService.AssertExpectations(t)
}

func TestHandler_RegisterOwner_InvalidBody(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/auth/register/owner", bytes.NewBuffer([]byte("{bad json")))
	rec := httptest.NewRecorder()

	handler.RegisterOwner(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandler_RegisterSitter_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := RegisterSitterRequest{
		FullName:        "Test Sitter",
		Email:           "sitter@test.com",
		Phone:           "+77001234568",
		Password:        "strongpassword",
		ExperienceYears: 3,
		Certificates:    "CPR",
		Preferences:     "Dogs",
		Location:        "Almaty",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register/sitter", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On(
			"RegisterSitter",
			reqBody.FullName,
			reqBody.Email,
			reqBody.Phone,
			reqBody.Password,
			reqBody.ExperienceYears,
			reqBody.Certificates,
			reqBody.Preferences,
			reqBody.Location,
		).
		Return(nil)

	handler.RegisterSitter(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "няня зарегистрирована, ожидает подтверждения", resp["message"])

	mockService.AssertExpectations(t)
}

func TestHandler_Login_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := LoginRequest{
		Email:    "user@test.com",
		Password: "password",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	user := &models.User{
		UserID:   1,
		Email:    "user@test.com",
		FullName: "Test User",
		Role:     "owner",
	}

	mockService.
		On("Login", reqBody.Email, reqBody.Password).
		Return(user, "jwt-token", nil)

	handler.Login(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resp)

	assert.Equal(t, "вход выполнен", resp["message"])
	assert.Equal(t, float64(1), resp["user_id"])
	assert.Equal(t, "owner", resp["role"])
	assert.Equal(t, "jwt-token", resp["token"])

	mockService.AssertExpectations(t)
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	reqBody := LoginRequest{
		Email:    "wrong@test.com",
		Password: "wrong",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	mockService.
		On("Login", reqBody.Email, reqBody.Password).
		Return(nil, "", assert.AnError)

	handler.Login(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
