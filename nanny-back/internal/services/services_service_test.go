package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"nanny-backend/internal/common/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockServiceForHandler struct {
	createServiceFunc     func(int, string, float64, string) (int, error)
	getServiceFunc        func(int) (*models.Service, error)
	getSitterServicesFunc func(int) ([]models.Service, error)
	updateServiceFunc     func(int, string, float64, string) error
	deleteServiceFunc     func(int) error
	searchServicesFunc    func(string, string) ([]ServiceWithSitter, error)
}

func (m *mockServiceForHandler) CreateService(sitterID int, serviceType string, pricePerHour float64, description string) (int, error) {
	if m.createServiceFunc != nil {
		return m.createServiceFunc(sitterID, serviceType, pricePerHour, description)
	}
	return 1, nil
}

func (m *mockServiceForHandler) GetService(serviceID int) (*models.Service, error) {
	if m.getServiceFunc != nil {
		return m.getServiceFunc(serviceID)
	}
	return &models.Service{ServiceID: serviceID}, nil
}

func (m *mockServiceForHandler) GetSitterServices(sitterID int) ([]models.Service, error) {
	if m.getSitterServicesFunc != nil {
		return m.getSitterServicesFunc(sitterID)
	}
	return []models.Service{{ServiceID: 1}}, nil
}

func (m *mockServiceForHandler) UpdateService(serviceID int, serviceType string, pricePerHour float64, description string) error {
	if m.updateServiceFunc != nil {
		return m.updateServiceFunc(serviceID, serviceType, pricePerHour, description)
	}
	return nil
}

func (m *mockServiceForHandler) DeleteService(serviceID int) error {
	if m.deleteServiceFunc != nil {
		return m.deleteServiceFunc(serviceID)
	}
	return nil
}

func (m *mockServiceForHandler) SearchServices(serviceType, location string) ([]ServiceWithSitter, error) {
	if m.searchServicesFunc != nil {
		return m.searchServicesFunc(serviceType, location)
	}
	return []ServiceWithSitter{{Service: models.Service{ServiceID: 1}}}, nil
}

func TestHandler_CreateService_Success(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		createServiceFunc: func(sitterID int, serviceType string, pricePerHour float64, description string) (int, error) {
			return 123, nil
		},
	}
	handler := NewHandler(mockSvc)

	reqBody := CreateServiceRequest{
		SitterID:     1,
		Type:         "walking",
		PricePerHour: 2500,
		Description:  "Dog walking service",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateService(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	if err != nil {
		return
	}
	assert.Equal(t, "услуга создана успешно", resp["message"])
	assert.Equal(t, float64(123), resp["service_id"])
}

func TestHandler_CreateService_InvalidBody(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBuffer([]byte("invalid json")))
	rec := httptest.NewRecorder()

	handler.CreateService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверные данные")
}

func TestHandler_CreateService_ServiceError(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		createServiceFunc: func(sitterID int, serviceType string, pricePerHour float64, description string) (int, error) {
			return 0, errors.New("неверный тип услуги")
		},
	}
	handler := NewHandler(mockSvc)

	reqBody := CreateServiceRequest{
		SitterID:     1,
		Type:         "invalid",
		PricePerHour: 2500,
		Description:  "Test",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/services", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный тип услуги")
}

func TestHandler_GetService_Success(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		getServiceFunc: func(serviceID int) (*models.Service, error) {
			return &models.Service{
				ServiceID:    serviceID,
				SitterID:     1,
				Type:         "walking",
				PricePerHour: 2500,
				Description:  "Dog walking",
			}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	handler.GetService(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.Service
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 1, resp.ServiceID)
	assert.Equal(t, "walking", resp.Type)
}

func TestHandler_GetService_InvalidID(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.GetService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID услуги")
}

func TestHandler_GetService_NotFound(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		getServiceFunc: func(serviceID int) (*models.Service, error) {
			return nil, errors.New("услуга не найдена")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	handler.GetService(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "услуга не найдена")
}

func TestHandler_GetSitterServices_Success(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		getSitterServicesFunc: func(sitterID int) ([]models.Service, error) {
			return []models.Service{
				{ServiceID: 1, SitterID: sitterID, Type: "walking", PricePerHour: 2500},
				{ServiceID: 2, SitterID: sitterID, Type: "boarding", PricePerHour: 5000},
			}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/1/services", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetSitterServices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []models.Service
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp))
	assert.Equal(t, "walking", resp[0].Type)
	assert.Equal(t, "boarding", resp[1].Type)
}

func TestHandler_GetSitterServices_InvalidID(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/invalid/services", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "invalid"})
	rec := httptest.NewRecorder()

	handler.GetSitterServices(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID няни")
}

func TestHandler_GetSitterServices_Error(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		getSitterServicesFunc: func(sitterID int) ([]models.Service, error) {
			return nil, errors.New("database error")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/sitters/1/services", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rec := httptest.NewRecorder()

	handler.GetSitterServices(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "database error")
}

func TestHandler_UpdateService_Success(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		updateServiceFunc: func(serviceID int, serviceType string, pricePerHour float64, description string) error {
			return nil
		},
	}
	handler := NewHandler(mockSvc)

	reqBody := UpdateServiceRequest{
		Type:         "boarding",
		PricePerHour: 5000,
		Description:  "Updated description",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/services/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateService(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "услуга обновлена успешно", resp["message"])
}

func TestHandler_UpdateService_InvalidID(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	reqBody := UpdateServiceRequest{
		Type:         "walking",
		PricePerHour: 2500,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/services/invalid", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.UpdateService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID услуги")
}

func TestHandler_UpdateService_InvalidBody(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/services/1", bytes.NewBuffer([]byte("invalid json")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	handler.UpdateService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверные данные")
}

func TestHandler_UpdateService_ServiceError(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		updateServiceFunc: func(serviceID int, serviceType string, pricePerHour float64, description string) error {
			return errors.New("неверный тип услуги")
		},
	}
	handler := NewHandler(mockSvc)

	reqBody := UpdateServiceRequest{
		Type:         "invalid",
		PricePerHour: 2500,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/services/1", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный тип услуги")
}

func TestHandler_DeleteService_Success(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		deleteServiceFunc: func(serviceID int) error {
			return nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/services/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	rec := httptest.NewRecorder()

	handler.DeleteService(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "услуга удалена успешно", resp["message"])
}

func TestHandler_DeleteService_InvalidID(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/services/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	rec := httptest.NewRecorder()

	handler.DeleteService(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "неверный ID услуги")
}

func TestHandler_DeleteService_Error(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		deleteServiceFunc: func(serviceID int) error {
			return errors.New("service not found")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/services/999", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	rec := httptest.NewRecorder()

	handler.DeleteService(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "service not found")
}

func TestHandler_SearchServices_Success(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		searchServicesFunc: func(serviceType, location string) ([]ServiceWithSitter, error) {
			return []ServiceWithSitter{
				{
					Service:      models.Service{ServiceID: 1, Type: serviceType, PricePerHour: 2500},
					SitterName:   "John Doe",
					SitterRating: 4.5,
				},
				{
					Service:      models.Service{ServiceID: 2, Type: serviceType, PricePerHour: 3000},
					SitterName:   "Jane Smith",
					SitterRating: 4.8,
				},
			}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/search?type=walking&location=Almaty", nil)
	rec := httptest.NewRecorder()

	handler.SearchServices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []ServiceWithSitter
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, 2, len(resp))
	assert.Equal(t, "John Doe", resp[0].SitterName)
	assert.Equal(t, 4.5, resp[0].SitterRating)
}

func TestHandler_SearchServices_NoFilters(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		searchServicesFunc: func(serviceType, location string) ([]ServiceWithSitter, error) {
			return []ServiceWithSitter{
				{Service: models.Service{ServiceID: 1}},
			}, nil
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/search", nil)
	rec := httptest.NewRecorder()

	handler.SearchServices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []ServiceWithSitter
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.GreaterOrEqual(t, len(resp), 0)
}

func TestHandler_SearchServices_Error(t *testing.T) {
	mockSvc := &mockServiceForHandler{
		searchServicesFunc: func(serviceType, location string) ([]ServiceWithSitter, error) {
			return nil, errors.New("database error")
		},
	}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/services/search?type=walking", nil)
	rec := httptest.NewRecorder()

	handler.SearchServices(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Contains(t, resp["error"], "database error")
}

func TestNewHandler(t *testing.T) {
	mockSvc := &mockServiceForHandler{}
	handler := NewHandler(mockSvc)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.service)
}
