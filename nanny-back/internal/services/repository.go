package services

import (
	"database/sql"
	"fmt"

	"nanny-backend/internal/common/models"
)

type Repository interface {
	Create(service *models.Service) (int, error)
	GetByID(serviceID int) (*models.Service, error)
	GetBySitterID(sitterID int) ([]models.Service, error)
	Update(service *models.Service) error
	Delete(serviceID int) error
	SearchServices(serviceType, location string) ([]ServiceWithSitter, error)
}

type ServiceWithSitter struct {
	models.Service
	SitterName   string  `json:"sitter_name"`
	SitterRating float64 `json:"sitter_rating"`
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(service *models.Service) (int, error) {
	var serviceID int
	err := r.db.QueryRow(`
		INSERT INTO services (sitter_id, type, price_per_hour, description)
		VALUES ($1, $2, $3, $4)
		RETURNING service_id
	`, service.SitterID, service.Type, service.PricePerHour, service.Description).Scan(&serviceID)

	if err != nil {
		return 0, fmt.Errorf("coould not создать serviceу: %w", err)
	}

	return serviceID, nil
}

func (r *repository) GetByID(serviceID int) (*models.Service, error) {
	service := &models.Service{}
	err := r.db.QueryRow(`
		SELECT service_id, sitter_id, type, price_per_hour, description
		FROM services
		WHERE service_id = $1
	`, serviceID).Scan(
		&service.ServiceID,
		&service.SitterID,
		&service.Type,
		&service.PricePerHour,
		&service.Description,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("service not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting service: %w", err)
	}

	return service, nil
}

func (r *repository) GetBySitterID(sitterID int) ([]models.Service, error) {
	rows, err := r.db.Query(`
		SELECT service_id, sitter_id, type, price_per_hour, description
		FROM services
		WHERE sitter_id = $1
	`, sitterID)

	if err != nil {
		return nil, fmt.Errorf("error getting service: %w", err)
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var service models.Service
		err := rows.Scan(
			&service.ServiceID,
			&service.SitterID,
			&service.Type,
			&service.PricePerHour,
			&service.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *repository) Update(service *models.Service) error {
	_, err := r.db.Exec(`
		UPDATE services
		SET type = $1, price_per_hour = $2, description = $3
		WHERE service_id = $4
	`, service.Type, service.PricePerHour, service.Description, service.ServiceID)

	if err != nil {
		return fmt.Errorf("coould not update service: %w", err)
	}

	return nil
}

func (r *repository) Delete(serviceID int) error {
	_, err := r.db.Exec(`DELETE FROM services WHERE service_id = $1`, serviceID)
	if err != nil {
		return fmt.Errorf("coould not delete service: %w", err)
	}
	return nil
}

func (r *repository) SearchServices(serviceType, location string) ([]ServiceWithSitter, error) {
	query := `
		SELECT 
			s.service_id, s.sitter_id, s.type, s.price_per_hour, s.description,
			u.full_name as sitter_name,
			COALESCE(AVG(r.rating), 0) as sitter_rating
		FROM services s
		JOIN sitters st ON s.sitter_id = st.sitter_id
		JOIN users u ON st.sitter_id = u.user_id
		LEFT JOIN reviews r ON st.sitter_id = r.sitter_id
		WHERE st.status = 'approved'
	`

	args := []interface{}{}
	argCount := 1

	if serviceType != "" {
		query += fmt.Sprintf(" AND s.type = $%d", argCount)
		args = append(args, serviceType)
		argCount++
	}

	if location != "" {
		query += fmt.Sprintf(" AND st.location ILIKE $%d", argCount)
		args = append(args, "%"+location+"%")
	}

	query += " GROUP BY s.service_id, u.full_name ORDER BY sitter_rating DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error searching service: %w", err)
	}
	defer rows.Close()

	var services []ServiceWithSitter
	for rows.Next() {
		var service ServiceWithSitter
		err := rows.Scan(
			&service.ServiceID,
			&service.SitterID,
			&service.Type,
			&service.PricePerHour,
			&service.Description,
			&service.SitterName,
			&service.SitterRating,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning service: %w", err)
		}
		services = append(services, service)
	}

	return services, nil
}
