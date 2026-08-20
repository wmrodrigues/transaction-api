package repository

import (
	"context"
	"transaction-api/internal/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
