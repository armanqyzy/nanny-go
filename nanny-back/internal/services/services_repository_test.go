package services

import (
	"database/sql"
	"errors"
	"testing"

	"nanny-backend/internal/common/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateServiceRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	srv := &models.Service{
		SitterID:     1,
		Type:         "walking",
		PricePerHour: 2500,
		Description:  "Dog walking",
	}

	mock.ExpectQuery("INSERT INTO services").
		WithArgs(srv.SitterID, srv.Type, srv.PricePerHour, srv.Description).
		WillReturnRows(sqlmock.NewRows([]string{"service_id"}).AddRow(1))

	id, err := repo.Create(srv)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if id != 1 {
		t.Errorf("expected ID 1, got %d", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetByIDRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"service_id", "sitter_id", "type", "price_per_hour", "description"}).
		AddRow(1, 2, "walking", 2500.0, "Dog walking")

	mock.ExpectQuery("SELECT (.+) FROM services WHERE service_id").
		WithArgs(1).
		WillReturnRows(rows)

	srv, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if srv.ServiceID != 1 {
		t.Errorf("expected service ID 1, got %d", srv.ServiceID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	mock.ExpectQuery("SELECT (.+) FROM services WHERE service_id").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByID(999)
	if err == nil {
		t.Error("expected error for non-existent service")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestGetBySitterIDRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"service_id", "sitter_id", "type", "price_per_hour", "description"}).
		AddRow(1, 2, "walking", 2500.0, "Dog walking").
		AddRow(2, 2, "boarding", 5000.0, "Pet boarding")

	mock.ExpectQuery("SELECT (.+) FROM services WHERE sitter_id").
		WithArgs(2).
		WillReturnRows(rows)

	services, err := repo.GetBySitterID(2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestUpdateRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	srv := &models.Service{
		ServiceID:    1,
		SitterID:     2,
		Type:         "boarding",
		PricePerHour: 5000,
		Description:  "Updated",
	}

	mock.ExpectExec("UPDATE services SET").
		WithArgs(srv.Type, srv.PricePerHour, srv.Description, srv.ServiceID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(srv)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestDeleteRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	mock.ExpectExec("DELETE FROM services WHERE service_id").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSearchServicesRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"service_id", "sitter_id", "type", "price_per_hour", "description",
		"full_name", "rating",
	}).AddRow(1, 2, "walking", 2500.0, "Dog walking", "John Doe", 4.5)

	mock.ExpectQuery("SELECT (.+) FROM services").
		WillReturnRows(rows)

	services, err := repo.SearchServices("", "")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(services) == 0 {
		t.Error("expected at least one service")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestSearchWithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{
		"service_id", "sitter_id", "type", "price_per_hour", "description",
		"full_name", "rating",
	}).AddRow(1, 2, "walking", 2500.0, "Dog walking", "John Doe", 4.5)

	mock.ExpectQuery("SELECT (.+) FROM services").
		WithArgs("walking", "%Almaty%").
		WillReturnRows(rows)

	services, err := repo.SearchServices("walking", "Almaty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(services) != 1 {
		t.Errorf("expected 1 service, got %d", len(services))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestRepositoryErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := NewRepository(db)

	t.Run("create error", func(t *testing.T) {
		srv := &models.Service{SitterID: 1, Type: "walking", PricePerHour: 2500}
		mock.ExpectQuery("INSERT INTO services").
			WillReturnError(errors.New("db error"))

		_, err := repo.Create(srv)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("update error", func(t *testing.T) {
		srv := &models.Service{ServiceID: 1, Type: "walking", PricePerHour: 2500}
		mock.ExpectExec("UPDATE services").
			WillReturnError(errors.New("db error"))

		err := repo.Update(srv)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("delete error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM services").
			WillReturnError(errors.New("db error"))

		err := repo.Delete(1)
		if err == nil {
			t.Error("expected error")
		}
	})
}
