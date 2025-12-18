package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestNew_Success(t *testing.T) {
	_, err := New("invalid_connection_string")
	if err == nil {
		t.Error("expected error with invalid connection string")
	}
}

func TestDatabase_Close(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	database := &Database{DB: db}

	mock.ExpectClose()

	err = database.Close()
	if err != nil {
		t.Errorf("unexpected error closing database: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDatabase_CloseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}

	database := &Database{DB: db}

	mock.ExpectClose().WillReturnError(sqlmock.ErrCancelled)

	err = database.Close()
	if err == nil {
		t.Error("expected error on close")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDatabaseStruct(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	database := &Database{DB: db}

	if database.DB == nil {
		t.Error("expected DB to be set")
	}
}
