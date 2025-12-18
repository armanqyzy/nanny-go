package admin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"nanny-backend/internal/common/models"

	"github.com/gorilla/mux"
)

type mockAdminServiceForHandler struct {
	getPendingSittersFunc func() ([]models.Sitter, error)
	approveSitterFunc     func(int) error
	rejectSitterFunc      func(int) error
	getAllUsersFunc       func() ([]models.User, error)
	getUserFunc           func(int) (*models.User, error)
	deleteUserFunc        func(int) error
	getSitterDetailsFunc  func(int) (*SitterDetails, error)
}

func (m *mockAdminServiceForHandler) GetPendingSitters() ([]models.Sitter, error) {
	if m.getPendingSittersFunc != nil {
		return m.getPendingSittersFunc()
	}
	return []models.Sitter{{SitterID: 1}}, nil
}

func (m *mockAdminServiceForHandler) ApproveSitter(id int) error {
	if m.approveSitterFunc != nil {
		return m.approveSitterFunc(id)
	}
	return nil
}

func (m *mockAdminServiceForHandler) RejectSitter(id int) error {
	if m.rejectSitterFunc != nil {
		return m.rejectSitterFunc(id)
	}
	return nil
}

func (m *mockAdminServiceForHandler) GetAllUsers() ([]models.User, error) {
	if m.getAllUsersFunc != nil {
		return m.getAllUsersFunc()
	}
	return []models.User{{UserID: 1}}, nil
}

func (m *mockAdminServiceForHandler) GetUser(id int) (*models.User, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(id)
	}
	return &models.User{UserID: id}, nil
}

func (m *mockAdminServiceForHandler) DeleteUser(id int) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(id)
	}
	return nil
}

func (m *mockAdminServiceForHandler) GetSitterDetails(id int) (*SitterDetails, error) {
	if m.getSitterDetailsFunc != nil {
		return m.getSitterDetailsFunc(id)
	}
	return &SitterDetails{Sitter: models.Sitter{SitterID: id}}, nil
}

func TestGetPendingSittersHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getPendingSittersFunc: func() ([]models.Sitter, error) {
			return []models.Sitter{{SitterID: 1}, {SitterID: 2}}, nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/sitters/pending", nil)
	rr := httptest.NewRecorder()

	handler.GetPendingSitters(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestGetPendingSittersHandler_Error(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getPendingSittersFunc: func() ([]models.Sitter, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/sitters/pending", nil)
	rr := httptest.NewRecorder()

	handler.GetPendingSitters(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}

func TestApproveSitterHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		approveSitterFunc: func(id int) error {
			return nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/api/admin/sitters/1/approve", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rr := httptest.NewRecorder()

	handler.ApproveSitter(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestApproveSitterHandler_InvalidID(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/sitters/invalid/approve", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "invalid"})
	rr := httptest.NewRecorder()

	handler.ApproveSitter(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestApproveSitterHandler_ServiceError(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		approveSitterFunc: func(id int) error {
			return errors.New("already approved")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/api/admin/sitters/1/approve", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rr := httptest.NewRecorder()

	handler.ApproveSitter(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestRejectSitterHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		rejectSitterFunc: func(id int) error {
			return nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/api/admin/sitters/1/reject", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rr := httptest.NewRecorder()

	handler.RejectSitter(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestGetAllUsersHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getAllUsersFunc: func() ([]models.User, error) {
			return []models.User{{UserID: 1}, {UserID: 2}}, nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	rr := httptest.NewRecorder()

	handler.GetAllUsers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestGetUserHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getUserFunc: func(id int) (*models.User, error) {
			return &models.User{UserID: id, Email: "test@example.com"}, nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestGetUserHandler_NotFound(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getUserFunc: func(id int) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/999", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "999"})
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

func TestDeleteUserHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		deleteUserFunc: func(id int) error {
			return nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestGetSitterDetailsHandler(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getSitterDetailsFunc: func(id int) (*SitterDetails, error) {
			return &SitterDetails{
				Sitter:   models.Sitter{SitterID: id},
				FullName: "John Doe",
			}, nil
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/sitters/1", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rr := httptest.NewRecorder()

	handler.GetSitterDetails(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

func TestRejectSitterHandler_InvalidID(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodPut, "/api/admin/sitters/invalid/reject", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "invalid"})
	rr := httptest.NewRecorder()

	handler.RejectSitter(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestRejectSitterHandler_ServiceError(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		rejectSitterFunc: func(id int) error {
			return errors.New("rejection failed")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodPut, "/api/admin/sitters/1/reject", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "1"})
	rr := httptest.NewRecorder()

	handler.RejectSitter(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestGetAllUsersHandler_Error(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getAllUsersFunc: func() ([]models.User, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	rr := httptest.NewRecorder()

	handler.GetAllUsers(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}

func TestGetUserHandler_InvalidID(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/users/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "invalid"})
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestDeleteUserHandler_InvalidID(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "invalid"})
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestDeleteUserHandler_ServiceError(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		deleteUserFunc: func(id int) error {
			return errors.New("delete failed")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodDelete, "/api/admin/users/1", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
	rr := httptest.NewRecorder()

	handler.DeleteUser(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
}

func TestGetSitterDetailsHandler_InvalidID(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{}
	handler := NewHandler(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/admin/sitters/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "invalid"})
	rr := httptest.NewRecorder()

	handler.GetSitterDetails(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
}

func TestGetSitterDetailsHandler_NotFound(t *testing.T) {
	mockSvc := &mockAdminServiceForHandler{
		getSitterDetailsFunc: func(id int) (*SitterDetails, error) {
			return nil, errors.New("sitter not found")
		},
	}

	handler := NewHandler(mockSvc)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/sitters/999", nil)
	req = mux.SetURLVars(req, map[string]string{"sitter_id": "999"})
	rr := httptest.NewRecorder()

	handler.GetSitterDetails(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}
