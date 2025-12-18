package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"nanny-backend/internal/common/models"

	"github.com/gorilla/mux"
)

type Service interface {
	CreateService(sitterID int, serviceType string, pricePerHour float64, description string) (int, error)
	GetService(serviceID int) (*models.Service, error)
	GetSitterServices(sitterID int) ([]models.Service, error)
	UpdateService(serviceID int, serviceType string, pricePerHour float64, description string) error
	DeleteService(serviceID int) error
	SearchServices(serviceType, location string) ([]ServiceWithSitter, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateService(sitterID int, serviceType string, pricePerHour float64, description string) (int, error) {
	validTypes := map[string]bool{"walking": true, "boarding": true, "home-care": true}
	if !validTypes[serviceType] {
		return 0, fmt.Errorf("incorrect type of service. Allowed: walking, boarding, home-care")
	}

	if pricePerHour <= 0 {
		return 0, fmt.Errorf("price must be more than 0")
	}

	srv := &models.Service{
		SitterID:     sitterID,
		Type:         serviceType,
		PricePerHour: pricePerHour,
		Description:  description,
	}

	serviceID, err := s.repo.Create(srv)
	if err != nil {
		return 0, fmt.Errorf("error creating service: %w", err)
	}

	return serviceID, nil
}

func (s *service) GetService(serviceID int) (*models.Service, error) {
	return s.repo.GetByID(serviceID)
}

func (s *service) GetSitterServices(sitterID int) ([]models.Service, error) {
	return s.repo.GetBySitterID(sitterID)
}

func (s *service) UpdateService(serviceID int, serviceType string, pricePerHour float64, description string) error {
	validTypes := map[string]bool{"walking": true, "boarding": true, "home-care": true}
	if !validTypes[serviceType] {
		return fmt.Errorf("incorrect type of service. Allowed: walking, boarding, home-care")
	}

	if pricePerHour <= 0 {
		return fmt.Errorf("price must be more than 0")
	}

	srv := &models.Service{
		ServiceID:    serviceID,
		Type:         serviceType,
		PricePerHour: pricePerHour,
		Description:  description,
	}

	return s.repo.Update(srv)
}

func (s *service) DeleteService(serviceID int) error {
	return s.repo.Delete(serviceID)
}

func (s *service) SearchServices(serviceType, location string) ([]ServiceWithSitter, error) {
	return s.repo.SearchServices(serviceType, location)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateServiceRequest struct {
	SitterID     int     `json:"sitter_id"`
	Type         string  `json:"type"`
	PricePerHour float64 `json:"price_per_hour"`
	Description  string  `json:"description,omitempty"`
}

type UpdateServiceRequest struct {
	Type         string  `json:"type"`
	PricePerHour float64 `json:"price_per_hour"`
	Description  string  `json:"description,omitempty"`
}

func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request) {
	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect data")
		return
	}

	serviceID, err := h.service.CreateService(req.SitterID, req.Type, req.PricePerHour, req.Description)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message":    "service created succesfully",
		"service_id": serviceID,
	})
}

func (h *Handler) GetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID service")
		return
	}

	service, err := h.service.GetService(serviceID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, service)
}

func (h *Handler) GetSitterServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sitterID, err := strconv.Atoi(vars["sitter_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID nanny")
		return
	}

	services, err := h.service.GetSitterServices(sitterID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, services)
}

func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID service")
		return
	}

	var req UpdateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect data")
		return
	}

	err = h.service.UpdateService(serviceID, req.Type, req.PricePerHour, req.Description)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "service updated succesfully",
	})
}

func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "incorrect ID service")
		return
	}

	err = h.service.DeleteService(serviceID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"message": "service deleted succesfully",
	})
}

func (h *Handler) SearchServices(w http.ResponseWriter, r *http.Request) {
	serviceType := r.URL.Query().Get("type")
	location := r.URL.Query().Get("location")

	services, err := h.service.SearchServices(serviceType, location)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, services)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
