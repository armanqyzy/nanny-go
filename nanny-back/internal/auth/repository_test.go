package auth

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"nanny-backend/internal/common/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	user := &models.User{
		FullName:     "Test User",
		Email:        "test@mail.com",
		Phone:        "+77001234567",
		PasswordHash: "hashed_password",
		Role:         "owner",
	}

	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(1)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.FullName, user.Email, user.Phone, user.PasswordHash, user.Role).
		WillReturnRows(rows)

	userID, err := repo.CreateUser(user)

	assert.NoError(t, err)
	assert.Equal(t, 1, userID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	user := &models.User{
		FullName:     "Test User",
		Email:        "test@mail.com",
		Phone:        "+77001234567",
		PasswordHash: "hashed_password",
		Role:         "owner",
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.FullName, user.Email, user.Phone, user.PasswordHash, user.Role).
		WillReturnError(errors.New("duplicate email"))

	userID, err := repo.CreateUser(user)

	assert.Error(t, err)
	assert.Equal(t, 0, userID)
	assert.Contains(t, err.Error(), "could not create a user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"user_id",
		"full_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"created_at",
	}).AddRow(
		1,
		"Test User",
		"test@mail.com",
		"+77001234567",
		"hashed_password",
		"owner",
		now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
		WithArgs("test@mail.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("test@mail.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.UserID)
	assert.Equal(t, "test@mail.com", user.Email)
	assert.Equal(t, "Test User", user.FullName)
	assert.Equal(t, "owner", user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
		WithArgs("notfound@mail.com").
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByEmail("notfound@mail.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
		WithArgs("test@mail.com").
		WillReturnError(errors.New("database connection error"))

	user, err := repo.GetUserByEmail("test@mail.com")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "error getting user")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail_SitterRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	now := time.Now()

	rows := sqlmock.NewRows([]string{
		"user_id",
		"full_name",
		"email",
		"phone",
		"password_hash",
		"role",
		"created_at",
	}).AddRow(
		2,
		"Test Sitter",
		"sitter@mail.com",
		"+77001234568",
		"hashed_password",
		"sitter",
		now,
	)

	mock.ExpectQuery(`SELECT (.+) FROM users WHERE email`).
		WithArgs("sitter@mail.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("sitter@mail.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 2, user.UserID)
	assert.Equal(t, "sitter", user.Role)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSitter_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	sitter := &models.Sitter{
		SitterID:        1,
		ExperienceYears: 5,
		Certificates:    "CPR Certified",
		Preferences:     "Dogs, Cats",
		Location:        "Almaty",
		Status:          "pending",
	}

	mock.ExpectExec(`INSERT INTO sitters`).
		WithArgs(
			sitter.SitterID,
			sitter.ExperienceYears,
			sitter.Certificates,
			sitter.Preferences,
			sitter.Location,
			sitter.Status,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateSitter(sitter)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSitter_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	sitter := &models.Sitter{
		SitterID:        999,
		ExperienceYears: 5,
		Certificates:    "CPR",
		Preferences:     "Dogs",
		Location:        "Almaty",
		Status:          "pending",
	}

	mock.ExpectExec(`INSERT INTO sitters`).
		WithArgs(
			sitter.SitterID,
			sitter.ExperienceYears,
			sitter.Certificates,
			sitter.Preferences,
			sitter.Location,
			sitter.Status,
		).
		WillReturnError(errors.New("foreign key constraint failed"))

	err = repo.CreateSitter(sitter)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not create a nanny profile")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSitter_DuplicateID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	sitter := &models.Sitter{
		SitterID:        1,
		ExperienceYears: 3,
		Certificates:    "First Aid",
		Preferences:     "Cats",
		Location:        "Astana",
		Status:          "pending",
	}

	mock.ExpectExec(`INSERT INTO sitters`).
		WithArgs(
			sitter.SitterID,
			sitter.ExperienceYears,
			sitter.Certificates,
			sitter.Preferences,
			sitter.Location,
			sitter.Status,
		).
		WillReturnError(errors.New("duplicate key value"))

	err = repo.CreateSitter(sitter)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not create nanny profile")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestNewRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	assert.NotNil(t, repo)
}
