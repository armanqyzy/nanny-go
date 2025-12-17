package pets

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"nanny-backend/internal/common/models"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	pet := &models.Pet{
		OwnerID: 1,
		Name:    "Buddy",
		Type:    "dog",
		Age:     3,
		Notes:   "friendly",
	}

	mock.ExpectQuery(`INSERT INTO pets`).
		WithArgs(
			pet.OwnerID,
			pet.Name,
			pet.Type,
			pet.Age,
			pet.Notes,
		).
		WillReturnRows(sqlmock.NewRows([]string{"pet_id"}).AddRow(10))

	id, err := repo.Create(pet)

	assert.NoError(t, err)
	assert.Equal(t, 10, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{
		"pet_id",
		"owner_id",
		"name",
		"type",
		"age",
		"notes",
	}).AddRow(
		10,
		1,
		"Buddy",
		"dog",
		3,
		"friendly",
	)

	mock.ExpectQuery(`FROM pets WHERE pet_id = \$1`).
		WithArgs(10).
		WillReturnRows(rows)

	pet, err := repo.GetByID(10)

	assert.NoError(t, err)
	assert.NotNil(t, pet)
	assert.Equal(t, "Buddy", pet.Name)
	assert.Equal(t, "dog", pet.Type)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectQuery(`FROM pets WHERE pet_id = \$1`).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	pet, err := repo.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, pet)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByOwnerID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	rows := sqlmock.NewRows([]string{
		"pet_id",
		"owner_id",
		"name",
		"type",
		"age",
		"notes",
	}).
		AddRow(1, 5, "Buddy", "dog", 3, "friendly").
		AddRow(2, 5, "Murka", "cat", 2, "calm")

	mock.ExpectQuery(`FROM pets WHERE owner_id = \$1`).
		WithArgs(5).
		WillReturnRows(rows)

	pets, err := repo.GetByOwnerID(5)

	assert.NoError(t, err)
	assert.Len(t, pets, 2)
	assert.Equal(t, "Buddy", pets[0].Name)
	assert.Equal(t, "Murka", pets[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	pet := &models.Pet{
		PetID: 10,
		Name:  "Buddy Updated",
		Type:  "dog",
		Age:   4,
		Notes: "updated notes",
	}

	mock.ExpectExec(`UPDATE pets`).
		WithArgs(
			pet.Name,
			pet.Type,
			pet.Age,
			pet.Notes,
			pet.PetID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(pet)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &repository{db: db}

	mock.ExpectExec(`DELETE FROM pets`).
		WithArgs(10).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(10)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
