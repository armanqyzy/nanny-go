package services

import (
	"errors"
	"testing"

	"nanny-backend/internal/common/models"
)

type mockServiceRepository struct {
	createFunc         func(*models.Service) (int, error)
	getByIDFunc        func(int) (*models.Service, error)
	getBySitterIDFunc  func(int) ([]models.Service, error)
	updateFunc         func(*models.Service) error
	deleteFunc         func(int) error
	searchServicesFunc func(string, string) ([]ServiceWithSitter, error)
}

func (m *mockServiceRepository) Create(srv *models.Service) (int, error) {
	if m.createFunc != nil {
		return m.createFunc(srv)
	}
	return 1, nil
}

func (m *mockServiceRepository) GetByID(id int) (*models.Service, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id)
	}
	return &models.Service{ServiceID: id}, nil
}

func (m *mockServiceRepository) GetBySitterID(sitterID int) ([]models.Service, error) {
	if m.getBySitterIDFunc != nil {
		return m.getBySitterIDFunc(sitterID)
	}
	return []models.Service{{ServiceID: 1}}, nil
}

func (m *mockServiceRepository) Update(srv *models.Service) error {
	if m.updateFunc != nil {
		return m.updateFunc(srv)
	}
	return nil
}

func (m *mockServiceRepository) Delete(id int) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockServiceRepository) SearchServices(serviceType, location string) ([]ServiceWithSitter, error) {
	if m.searchServicesFunc != nil {
		return m.searchServicesFunc(serviceType, location)
	}
	return []ServiceWithSitter{{Service: models.Service{ServiceID: 1}}}, nil
}

func TestCreateService(t *testing.T) {
	tests := []struct {
		name         string
		sitterID     int
		serviceType  string
		pricePerHour float64
		description  string
		mockCreate   func(*models.Service) (int, error)
		expectError  bool
	}{
		{
			name:         "successful creation",
			sitterID:     1,
			serviceType:  "walking",
			pricePerHour: 2500,
			description:  "Dog walking",
			mockCreate: func(srv *models.Service) (int, error) {
				return 1, nil
			},
			expectError: false,
		},
		{
			name:         "invalid service type",
			sitterID:     1,
			serviceType:  "invalid",
			pricePerHour: 2500,
			description:  "Test",
			expectError:  true,
		},
		{
			name:         "invalid price",
			sitterID:     1,
			serviceType:  "walking",
			pricePerHour: -100,
			description:  "Test",
			expectError:  true,
		},
		{
			name:         "repository error",
			sitterID:     1,
			serviceType:  "walking",
			pricePerHour: 2500,
			description:  "Test",
			mockCreate: func(srv *models.Service) (int, error) {
				return 0, errors.New("database error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockServiceRepository{createFunc: tt.mockCreate}
			svc := NewService(repo)

			id, err := svc.CreateService(tt.sitterID, tt.serviceType, tt.pricePerHour, tt.description)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && id == 0 {
				t.Error("expected valid service ID")
			}
		})
	}
}

func TestGetService(t *testing.T) {
	repo := &mockServiceRepository{
		getByIDFunc: func(id int) (*models.Service, error) {
			return &models.Service{ServiceID: id, Type: "walking"}, nil
		},
	}
	svc := NewService(repo)

	service, err := svc.GetService(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if service.ServiceID != 1 {
		t.Errorf("expected service ID 1, got %d", service.ServiceID)
	}
}

func TestGetSitterServices(t *testing.T) {
	repo := &mockServiceRepository{
		getBySitterIDFunc: func(sitterID int) ([]models.Service, error) {
			return []models.Service{
				{ServiceID: 1, SitterID: sitterID},
				{ServiceID: 2, SitterID: sitterID},
			}, nil
		},
	}
	svc := NewService(repo)

	services, err := svc.GetSitterServices(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(services) != 2 {
		t.Errorf("expected 2 services, got %d", len(services))
	}
}

func TestUpdateService(t *testing.T) {
	tests := []struct {
		name         string
		serviceID    int
		serviceType  string
		pricePerHour float64
		description  string
		expectError  bool
	}{
		{
			name:         "successful update",
			serviceID:    1,
			serviceType:  "boarding",
			pricePerHour: 5000,
			description:  "Updated",
			expectError:  false,
		},
		{
			name:         "invalid type",
			serviceID:    1,
			serviceType:  "invalid",
			pricePerHour: 5000,
			description:  "Test",
			expectError:  true,
		},
		{
			name:         "invalid price",
			serviceID:    1,
			serviceType:  "boarding",
			pricePerHour: 0,
			description:  "Test",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockServiceRepository{}
			svc := NewService(repo)

			err := svc.UpdateService(tt.serviceID, tt.serviceType, tt.pricePerHour, tt.description)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestDeleteService(t *testing.T) {
	repo := &mockServiceRepository{
		deleteFunc: func(id int) error {
			return nil
		},
	}
	svc := NewService(repo)

	err := svc.DeleteService(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSearchServices(t *testing.T) {
	repo := &mockServiceRepository{
		searchServicesFunc: func(serviceType, location string) ([]ServiceWithSitter, error) {
			return []ServiceWithSitter{
				{Service: models.Service{ServiceID: 1, Type: serviceType}},
			}, nil
		},
	}
	svc := NewService(repo)

	services, err := svc.SearchServices("walking", "Almaty")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(services) == 0 {
		t.Error("expected at least one service")
	}
}
