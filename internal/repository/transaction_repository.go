package repository

import (
	"context"
	"transaction-api/internal/domain"
)

type TransactionRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByUserId(ctx context.Context, userId string, pagination domain.Pagination) (*domain.Page[domain.Transaction], error)
}
