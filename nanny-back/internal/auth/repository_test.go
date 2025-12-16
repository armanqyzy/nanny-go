package auth

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail_Success(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

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
		now, // time.Time, а не string
	)

	mock.ExpectQuery(`FROM users WHERE email = \$1`).
		WithArgs("test@mail.com").
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail("test@mail.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@mail.com", user.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}
