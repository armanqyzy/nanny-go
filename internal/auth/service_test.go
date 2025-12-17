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
	assert.Contains(t, err.Error(), "ошибка регистрации")
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
	}

	mockRepo.
		On("GetUserByEmail", "test@mail.com").
		Return(user, nil)

	resultUser, token, err := service.Login("test@mail.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, resultUser)
	assert.NotEmpty(t, token)
	assert.Equal(t, "owner", resultUser.Role)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.
		On("GetUserByEmail", "wrong@mail.com").
		Return(nil, errors.New("not found"))

	user, token, err := service.Login("wrong@mail.com", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Empty(t, token)
}
