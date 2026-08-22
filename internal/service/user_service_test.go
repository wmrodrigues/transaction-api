package service

import (
	"context"
	"errors"
	"testing"
	"transaction-api/internal/domain"

	"gorm.io/gorm"
)

type mockUserRepository struct {
	user          *domain.User
	err           error
	getByEmailErr error
	createErr     error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.user, m.err
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	return m.err
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailErr != nil {
		return nil, m.getByEmailErr
	}
	return m.user, m.err
}

type mockTransactionManager struct{}

func (m *mockTransactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestGetUser(t *testing.T) {
	expectedUser := &domain.User{
		ID:    "123",
		Name:  "Wash",
		Email: "wash@example.com",
	}
	repository := &mockUserRepository{user: expectedUser}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	ctx := context.Background()
	user, err := userService.GetById(ctx, "123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Errorf("expected user ID to be %s, got %s", expectedUser.ID, user.ID)
	}
}

func TestGetUser_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	repository := &mockUserRepository{err: context.Canceled}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	_, err := userService.GetById(ctx, "123")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestGetUser_UserNotFound(t *testing.T) {
	repository := &mockUserRepository{err: errors.New("user not found")}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	_, err := userService.GetById(context.Background(), "123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUser_ValidateCredentialsSuccess(t *testing.T) {
	expectedUser := &domain.User{
		ID:       "123",
		Name:     "Wash",
		Email:    "wash@example.com",
		Password: "123456",
	}
	repository := &mockUserRepository{user: expectedUser}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	ctx := context.Background()
	err, _ := userService.ValidateCredentials(ctx, expectedUser.Email, expectedUser.Password)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetUser_ValidateCredentialsFailed(t *testing.T) {
	expectedUser := &domain.User{
		ID:       "123",
		Name:     "Wash",
		Email:    "wash@example.com",
		Password: "123456",
	}
	repository := &mockUserRepository{user: expectedUser}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	ctx := context.Background()
	_, err := userService.ValidateCredentials(ctx, expectedUser.Email, "wrong-password")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetUser_EmptyID(t *testing.T) {
	repository := &mockUserRepository{}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	_, err := userService.GetById(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID, got nil")
	}
	if err.Error() != "user id cannot be empty and must be a valid UUID" {
		t.Fatalf("expected empty ID error message, got %v", err)
	}
}

func TestGetUser_GormRecordNotFound(t *testing.T) {
	repository := &mockUserRepository{err: gorm.ErrRecordNotFound}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	_, err := userService.GetById(context.Background(), "123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "user not found" {
		t.Fatalf("expected 'user not found', got %v", err)
	}
}

func TestCreateUser_Success(t *testing.T) {
	repository := &mockUserRepository{
		getByEmailErr: errors.New("not found"),
	}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	user := &domain.User{
		Name:     "Wash",
		Email:    "wash@example.com",
		Password: "123456",
	}
	err := userService.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID == "" {
		t.Fatal("expected user ID to be assigned")
	}
}

func TestCreateUser_ValidationFailure(t *testing.T) {
	repository := &mockUserRepository{}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	tests := []struct {
		name string
		user domain.User
	}{
		{"missing name", domain.User{Email: "wash@example.com", Password: "123456"}},
		{"missing email", domain.User{Name: "Wash", Password: "123456"}},
		{"missing password", domain.User{Name: "Wash", Email: "wash@example.com"}},
	}
	// here we just make sure the domain validations work
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user := tc.user
			err := userService.Create(context.Background(), &user)
			if err == nil {
				t.Fatalf("expected validation error for %s, got nil", tc.name)
			}
		})
	}
}

func TestCreateUser_DuplicateUser(t *testing.T) {
	existingUser := &domain.User{
		ID:    "123",
		Name:  "Wash",
		Email: "wash@example.com",
	}
	repository := &mockUserRepository{user: existingUser}
	userService := NewUserService(repository, &mockTransactionRepository{}, &mockTransactionManager{})
	user := &domain.User{
		Name:     "Wash",
		Email:    "wash@example.com",
		Password: "123456",
	}
	err := userService.Create(context.Background(), user)
	if err == nil {
		t.Fatal("expected error for duplicate user, got nil")
	}
	if err.Error() != "user already exists" {
		t.Fatalf("expected 'user already exists', got %v", err)
	}
}
