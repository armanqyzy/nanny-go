package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"nanny-backend/internal/common/models"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(user *models.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) CreateSitter(sitter *models.Sitter) error {
	args := m.Called(sitter)
	return args.Error(0)
}

func TestRegisterOwner_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("CreateUser", mock.Anything).
		Return(1, nil)

	err := service.RegisterOwner(
		"Test User",
		"test@mail.com",
		"+77001234567",
		"password123",
	)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterOwner_EmailExists(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("CreateUser", mock.Anything).
		Return(0, errors.New("email already exists"))

	err := service.RegisterOwner(
		"Test User",
		"test@mail.com",
		"+77001234567",
		"password123",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error registration")
	mockRepo.AssertExpectations(t)
}

func TestRegisterSitter_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("CreateUser", mock.Anything).
		Return(1, nil)

	mockRepo.
		On("CreateSitter", mock.Anything).
		Return(nil)

	err := service.RegisterSitter(
		"Test Sitter",
		"sitter@mail.com",
		"+77001234567",
		"password123",
		5,
		"CPR Certified",
		"Dogs, Cats",
		"Almaty",
	)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterSitter_CreateUserError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("CreateUser", mock.Anything).
		Return(0, errors.New("database error"))

	err := service.RegisterSitter(
		"Test Sitter",
		"sitter@mail.com",
		"+77001234567",
		"password123",
		5,
		"CPR",
		"Dogs",
		"Almaty",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating user")
	mockRepo.AssertExpectations(t)
}

func TestRegisterSitter_CreateSitterError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("CreateUser", mock.Anything).
		Return(1, nil)

	mockRepo.
		On("CreateSitter", mock.Anything).
		Return(errors.New("database error"))

	err := service.RegisterSitter(
		"Test Sitter",
		"sitter@mail.com",
		"+77001234567",
		"password123",
		5,
		"CPR",
		"Dogs",
		"Almaty",
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating nanny prodile")
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)

	user := &models.User{
		UserID:       1,
		Email:        "test@mail.com",
		PasswordHash: string(hashedPassword),
		Role:         "owner",
		FullName:     "Test User",
	}

	mockRepo.
		On("GetUserByEmail", "test@mail.com").
		Return(user, nil)

	resultUser, token, err := service.Login("test@mail.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, "owner", resultUser.Role)
	assert.Equal(t, 1, resultUser.UserID)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("GetUserByEmail", "wrong@mail.com").
		Return(nil, errors.New("user not found"))

	user, token, err := service.Login("wrong@mail.com", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "incorrect email or password")
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("correctpassword"),
		bcrypt.DefaultCost,
	)

	user := &models.User{
		UserID:       1,
		Email:        "test@mail.com",
		PasswordHash: string(hashedPassword),
		Role:         "owner",
	}

	mockRepo.
		On("GetUserByEmail", "test@mail.com").
		Return(user, nil)

	resultUser, token, err := service.Login("test@mail.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, resultUser)
	assert.Empty(t, token)
	assert.Contains(t, err.Error(), "incorrect email or password")
	mockRepo.AssertExpectations(t)
}

func TestLogin_SitterRole(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)

	user := &models.User{
		UserID:       2,
		Email:        "sitter@mail.com",
		PasswordHash: string(hashedPassword),
		Role:         "sitter",
		FullName:     "Test Sitter",
	}

	mockRepo.
		On("GetUserByEmail", "sitter@mail.com").
		Return(user, nil)

	resultUser, token, err := service.Login("sitter@mail.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, "sitter", resultUser.Role)
	assert.Equal(t, 2, resultUser.UserID)
	mockRepo.AssertExpectations(t)
}
