package repository

import (
	"context"
	"transaction-api/internal/domain"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetForUpdate(ctx context.Context, userID, currency string) (*domain.Account, error)
	GetByUserID(ctx context.Context, userID string) ([]domain.Account, error)
	UpdateBalance(ctx context.Context, accountID string, newBalance int64) error
}
