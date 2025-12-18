package admin

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetPendingSittersRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"sitter_id", "experience_years", "certificates", "preferences", "location", "status"}).
		AddRow(1, 5, "cert1", "dogs", "Almaty", "pending").
		AddRow(2, 3, "cert2", "cats", "Astana", "pending")

	mock.ExpectQuery("SELECT (.+) FROM sitters WHERE status").
		WillReturnRows(rows)

	sitters, err := repo.GetPendingSitters()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(sitters) != 2 {
		t.Errorf("expected 2 sitters, got %d", len(sitters))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetPendingSittersRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM sitters WHERE status").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetPendingSitters()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetPendingSittersRepository_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"sitter_id", "experience_years", "certificates", "preferences", "location", "status"}).
		AddRow("invalid", 5, "cert1", "dogs", "Almaty", "pending")

	mock.ExpectQuery("SELECT (.+) FROM sitters WHERE status").
		WillReturnRows(rows)

	_, err = repo.GetPendingSitters()
	if err == nil {
		t.Error("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestApproveSitterRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("UPDATE sitters SET status").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.ApproveSitter(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestApproveSitterRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("UPDATE sitters SET status").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err = repo.ApproveSitter(1)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestRejectSitterRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("UPDATE sitters SET status").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.RejectSitter(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestRejectSitterRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("UPDATE sitters SET status").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err = repo.RejectSitter(1)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetAllUsersRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// ВАЖНО: используем time.Time вместо строки
	now := time.Now()
	rows := sqlmock.NewRows([]string{"user_id", "full_name", "email", "phone", "role", "created_at"}).
		AddRow(1, "User One", "user1@test.com", "123456", "owner", now).
		AddRow(2, "User Two", "user2@test.com", "654321", "sitter", now)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	users, err := repo.GetAllUsers()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetAllUsersRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetAllUsers()
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetAllUsersRepository_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"user_id", "full_name", "email", "phone", "role", "created_at"}).
		AddRow("invalid", "User One", "user1@test.com", "123456", "owner", time.Now())

	mock.ExpectQuery("SELECT (.+) FROM users").
		WillReturnRows(rows)

	_, err = repo.GetAllUsers()
	if err == nil {
		t.Error("expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetUserByIDRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	// ВАЖНО: используем time.Time вместо строки
	now := time.Now()
	rows := sqlmock.NewRows([]string{"user_id", "full_name", "email", "phone", "role", "created_at"}).
		AddRow(1, "Test User", "test@example.com", "1234567890", "owner", now)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", user.UserID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetUserByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByID(999)
	if err == nil {
		t.Error("expected error for non-existent user")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetUserByIDRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM users WHERE user_id").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	_, err = repo.GetUserByID(1)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteUserRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("DELETE FROM users WHERE user_id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteUser(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteUserRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("DELETE FROM users WHERE user_id").
		WithArgs(1).
		WillReturnError(errors.New("db error"))

	err = repo.DeleteUser(1)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetSitterDetailsRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"sitter_id", "experience_years", "certificates", "preferences", "location", "status",
		"full_name", "email", "phone", "rating", "reviews",
	}).AddRow(1, 5, "cert", "dogs", "Almaty", "approved", "John Doe", "john@test.com", "123", 4.5, 10)

	mock.ExpectQuery("SELECT (.+) FROM sitters").
		WithArgs(1).
		WillReturnRows(rows)

	details, err := repo.GetSitterDetails(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if details.SitterID != 1 {
		t.Errorf("expected sitter ID 1, got %d", details.SitterID)
	}
	if details.FullName != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", details.FullName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetSitterDetailsRepository_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM sitters").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetSitterDetails(999)
	if err == nil {
		t.Error("expected error for non-existent sitter")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetSitterDetailsRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM sitters").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	_, err = repo.GetSitterDetails(1)
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateSitterStatusRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("UPDATE sitters SET status").
		WithArgs("approved", 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateSitterStatus(1, "approved")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateSitterStatusRepository_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	mock.ExpectExec("UPDATE sitters SET status").
		WithArgs("approved", 1).
		WillReturnError(errors.New("db error"))

	err = repo.UpdateSitterStatus(1, "approved")
	if err == nil {
		t.Error("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}
