package pets

import (
	"database/sql"
	"fmt"

	"nanny-backend/internal/common/models"
)

type Repository interface {
	Create(pet *models.Pet) (int, error)
	GetByID(petID int) (*models.Pet, error)
	GetByOwnerID(ownerID int) ([]models.Pet, error)
	Update(pet *models.Pet) error
	Delete(petID int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(pet *models.Pet) (int, error) {
	var petID int
	err := r.db.QueryRow(`
		INSERT INTO pets (owner_id, name, type, age, notes)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING pet_id
	`, pet.OwnerID, pet.Name, pet.Type, pet.Age, pet.Notes).Scan(&petID)
	
	if err != nil {
		return 0, fmt.Errorf("не удалось создать питомца: %w", err)
	}
	
	return petID, nil
}

func (r *repository) GetByID(petID int) (*models.Pet, error) {
	pet := &models.Pet{}
	err := r.db.QueryRow(`
		SELECT pet_id, owner_id, name, type, age, notes
		FROM pets
		WHERE pet_id = $1
	`, petID).Scan(
		&pet.PetID,
		&pet.OwnerID,
		&pet.Name,
		&pet.Type,
		&pet.Age,
		&pet.Notes,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("питомец не найден")
	}
	if err != nil {
		return nil, fmt.Errorf("ошибка получения питомца: %w", err)
	}
	
	return pet, nil
}

func (r *repository) GetByOwnerID(ownerID int) ([]models.Pet, error) {
	rows, err := r.db.Query(`
		SELECT pet_id, owner_id, name, type, age, notes
		FROM pets
		WHERE owner_id = $1
	`, ownerID)
	
	if err != nil {
		return nil, fmt.Errorf("ошибка получения питомцев: %w", err)
	}
	defer rows.Close()
	
	var pets []models.Pet
	for rows.Next() {
		var pet models.Pet
		err := rows.Scan(
			&pet.PetID,
			&pet.OwnerID,
			&pet.Name,
			&pet.Type,
			&pet.Age,
			&pet.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования питомца: %w", err)
		}
		pets = append(pets, pet)
	}
	
	return pets, nil
}

func (r *repository) Update(pet *models.Pet) error {
	_, err := r.db.Exec(`
		UPDATE pets
		SET name = $1, type = $2, age = $3, notes = $4
		WHERE pet_id = $5
	`, pet.Name, pet.Type, pet.Age, pet.Notes, pet.PetID)
	
	if err != nil {
		return fmt.Errorf("не удалось обновить питомца: %w", err)
	}
	
	return nil
}

func (r *repository) Delete(petID int) error {
	_, err := r.db.Exec(`DELETE FROM pets WHERE pet_id = $1`, petID)
	if err != nil {
		return fmt.Errorf("не удалось удалить питомца: %w", err)
	}
	return nil
}
