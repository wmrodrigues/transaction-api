package service

import (
	"context"
	"errors"
	"testing"
	"transaction-api/internal/domain"
)

type mockUserRepository struct {
	user *domain.User
	err  error
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return m.user, m.err
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.err
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.user, m.err
}

func TestGetUser(t *testing.T) {
	expectedUser := &domain.User{
		ID:    "123",
		Name:  "Wash",
		Email: "wash@example.com",
	}
	repository := &mockUserRepository{user: expectedUser}
	userService := NewUserService(repository)
	ctx := context.Background()
	user, err := userService.GetUserById(ctx, "123")
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
	userService := NewUserService(repository)
	_, err := userService.GetUserById(ctx, "123")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"expected context.Canceled, got %v",
			err,
		)
	}
}
