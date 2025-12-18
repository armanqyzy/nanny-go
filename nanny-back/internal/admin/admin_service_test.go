package admin

import (
	"errors"
	"testing"

	"nanny-backend/internal/common/models"
)

type mockAdminRepository struct {
	getPendingSittersFunc  func() ([]models.Sitter, error)
	approveSitterFunc      func(int) error
	rejectSitterFunc       func(int) error
	getAllUsersFunc        func() ([]models.User, error)
	getUserByIDFunc        func(int) (*models.User, error)
	deleteUserFunc         func(int) error
	getSitterDetailsFunc   func(int) (*SitterDetails, error)
	updateSitterStatusFunc func(int, string) error
}

func (m *mockAdminRepository) GetPendingSitters() ([]models.Sitter, error) {
	if m.getPendingSittersFunc != nil {
		return m.getPendingSittersFunc()
	}
	return []models.Sitter{{SitterID: 1, Status: "pending"}}, nil
}

func (m *mockAdminRepository) ApproveSitter(sitterID int) error {
	if m.approveSitterFunc != nil {
		return m.approveSitterFunc(sitterID)
	}
	return nil
}

func (m *mockAdminRepository) RejectSitter(sitterID int) error {
	if m.rejectSitterFunc != nil {
		return m.rejectSitterFunc(sitterID)
	}
	return nil
}

func (m *mockAdminRepository) GetAllUsers() ([]models.User, error) {
	if m.getAllUsersFunc != nil {
		return m.getAllUsersFunc()
	}
	return []models.User{{UserID: 1, Email: "test@example.com"}}, nil
}

func (m *mockAdminRepository) GetUserByID(userID int) (*models.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(userID)
	}
	return &models.User{UserID: userID, Email: "test@example.com"}, nil
}

func (m *mockAdminRepository) DeleteUser(userID int) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(userID)
	}
	return nil
}

func (m *mockAdminRepository) GetSitterDetails(sitterID int) (*SitterDetails, error) {
	if m.getSitterDetailsFunc != nil {
		return m.getSitterDetailsFunc(sitterID)
	}
	return &SitterDetails{
		Sitter:   models.Sitter{SitterID: sitterID, Status: "pending"},
		FullName: "Test User",
		Email:    "test@example.com",
	}, nil
}

func (m *mockAdminRepository) UpdateSitterStatus(sitterID int, status string) error {
	if m.updateSitterStatusFunc != nil {
		return m.updateSitterStatusFunc(sitterID, status)
	}
	return nil
}

func TestGetPendingSitters(t *testing.T) {
	repo := &mockAdminRepository{
		getPendingSittersFunc: func() ([]models.Sitter, error) {
			return []models.Sitter{
				{SitterID: 1, Status: "pending"},
				{SitterID: 2, Status: "pending"},
			}, nil
		},
	}
	svc := NewService(repo)

	sitters, err := svc.GetPendingSitters()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(sitters) != 2 {
		t.Errorf("expected 2 sitters, got %d", len(sitters))
	}
}

func TestApproveSitter(t *testing.T) {
	tests := []struct {
		name              string
		sitterID          int
		mockGetDetails    func(int) (*SitterDetails, error)
		mockApproveSitter func(int) error
		expectError       bool
	}{
		{
			name:     "successful approval",
			sitterID: 1,
			mockGetDetails: func(id int) (*SitterDetails, error) {
				return &SitterDetails{
					Sitter: models.Sitter{SitterID: id, Status: "pending"},
				}, nil
			},
			mockApproveSitter: func(id int) error {
				return nil
			},
			expectError: false,
		},
		{
			name:     "sitter not found",
			sitterID: 999,
			mockGetDetails: func(id int) (*SitterDetails, error) {
				return nil, errors.New("sitter not found")
			},
			expectError: true,
		},
		{
			name:     "sitter not pending",
			sitterID: 1,
			mockGetDetails: func(id int) (*SitterDetails, error) {
				return &SitterDetails{
					Sitter: models.Sitter{SitterID: id, Status: "approved"},
				}, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepository{
				getSitterDetailsFunc: tt.mockGetDetails,
				approveSitterFunc:    tt.mockApproveSitter,
			}
			svc := NewService(repo)

			err := svc.ApproveSitter(tt.sitterID)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRejectSitter(t *testing.T) {
	tests := []struct {
		name             string
		sitterID         int
		mockGetDetails   func(int) (*SitterDetails, error)
		mockRejectSitter func(int) error
		expectError      bool
	}{
		{
			name:     "successful rejection",
			sitterID: 1,
			mockGetDetails: func(id int) (*SitterDetails, error) {
				return &SitterDetails{
					Sitter: models.Sitter{SitterID: id, Status: "pending"},
				}, nil
			},
			mockRejectSitter: func(id int) error {
				return nil
			},
			expectError: false,
		},
		{
			name:     "sitter already approved",
			sitterID: 1,
			mockGetDetails: func(id int) (*SitterDetails, error) {
				return &SitterDetails{
					Sitter: models.Sitter{SitterID: id, Status: "approved"},
				}, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepository{
				getSitterDetailsFunc: tt.mockGetDetails,
				rejectSitterFunc:     tt.mockRejectSitter,
			}
			svc := NewService(repo)

			err := svc.RejectSitter(tt.sitterID)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	repo := &mockAdminRepository{
		getAllUsersFunc: func() ([]models.User, error) {
			return []models.User{
				{UserID: 1, Email: "user1@example.com", Role: "owner"},
				{UserID: 2, Email: "user2@example.com", Role: "sitter"},
			}, nil
		},
	}
	svc := NewService(repo)

	users, err := svc.GetAllUsers()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestGetUser(t *testing.T) {
	repo := &mockAdminRepository{
		getUserByIDFunc: func(userID int) (*models.User, error) {
			return &models.User{UserID: userID, Email: "test@example.com"}, nil
		},
	}
	svc := NewService(repo)

	user, err := svc.GetUser(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.UserID != 1 {
		t.Errorf("expected user ID 1, got %d", user.UserID)
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		mockGetUser    func(int) (*models.User, error)
		mockDeleteUser func(int) error
		expectError    bool
	}{
		{
			name:   "successful deletion",
			userID: 1,
			mockGetUser: func(id int) (*models.User, error) {
				return &models.User{UserID: id}, nil
			},
			mockDeleteUser: func(id int) error {
				return nil
			},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: 999,
			mockGetUser: func(id int) (*models.User, error) {
				return nil, errors.New("user not found")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockAdminRepository{
				getUserByIDFunc: tt.mockGetUser,
				deleteUserFunc:  tt.mockDeleteUser,
			}
			svc := NewService(repo)

			err := svc.DeleteUser(tt.userID)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetSitterDetails(t *testing.T) {
	repo := &mockAdminRepository{
		getSitterDetailsFunc: func(sitterID int) (*SitterDetails, error) {
			return &SitterDetails{
				Sitter: models.Sitter{
					SitterID:        sitterID,
					ExperienceYears: 5,
					Status:          "approved",
				},
				FullName: "John Doe",
				Email:    "john@example.com",
			}, nil
		},
	}
	svc := NewService(repo)

	details, err := svc.GetSitterDetails(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if details.SitterID != 1 {
		t.Errorf("expected sitter ID 1, got %d", details.SitterID)
	}
	if details.FullName != "John Doe" {
		t.Errorf("expected name 'John Doe', got '%s'", details.FullName)
	}
}

func TestGetPendingSitters_Error(t *testing.T) {
	repo := &mockAdminRepository{
		getPendingSittersFunc: func() ([]models.Sitter, error) {
			return nil, errors.New("database error")
		},
	}
	svc := NewService(repo)

	_, err := svc.GetPendingSitters()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetAllUsers_Error(t *testing.T) {
	repo := &mockAdminRepository{
		getAllUsersFunc: func() ([]models.User, error) {
			return nil, errors.New("database error")
		},
	}
	svc := NewService(repo)

	_, err := svc.GetAllUsers()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetUser_Error(t *testing.T) {
	repo := &mockAdminRepository{
		getUserByIDFunc: func(userID int) (*models.User, error) {
			return nil, errors.New("user not found")
		},
	}
	svc := NewService(repo)

	_, err := svc.GetUser(999)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetSitterDetails_Error(t *testing.T) {
	repo := &mockAdminRepository{
		getSitterDetailsFunc: func(sitterID int) (*SitterDetails, error) {
			return nil, errors.New("sitter not found")
		},
	}
	svc := NewService(repo)

	_, err := svc.GetSitterDetails(999)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestApproveSitter_RepositoryError(t *testing.T) {
	repo := &mockAdminRepository{
		getSitterDetailsFunc: func(id int) (*SitterDetails, error) {
			return &SitterDetails{
				Sitter: models.Sitter{SitterID: id, Status: "pending"},
			}, nil
		},
		approveSitterFunc: func(id int) error {
			return errors.New("repository error")
		},
	}
	svc := NewService(repo)

	err := svc.ApproveSitter(1)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRejectSitter_RepositoryError(t *testing.T) {
	repo := &mockAdminRepository{
		getSitterDetailsFunc: func(id int) (*SitterDetails, error) {
			return &SitterDetails{
				Sitter: models.Sitter{SitterID: id, Status: "pending"},
			}, nil
		},
		rejectSitterFunc: func(id int) error {
			return errors.New("repository error")
		},
	}
	svc := NewService(repo)

	err := svc.RejectSitter(1)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestRejectSitter_GetDetailsError(t *testing.T) {
	repo := &mockAdminRepository{
		getSitterDetailsFunc: func(id int) (*SitterDetails, error) {
			return nil, errors.New("sitter not found")
		},
	}
	svc := NewService(repo)

	err := svc.RejectSitter(1)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDeleteUser_DeleteError(t *testing.T) {
	repo := &mockAdminRepository{
		getUserByIDFunc: func(id int) (*models.User, error) {
			return &models.User{UserID: id}, nil
		},
		deleteUserFunc: func(id int) error {
			return errors.New("delete failed")
		},
	}
	svc := NewService(repo)

	err := svc.DeleteUser(1)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
