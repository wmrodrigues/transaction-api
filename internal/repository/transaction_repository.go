package repository

import (
	"context"
	"transaction-api/internal/domain"
)

type TransactionRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetBalanceByUserId(ctx context.Context, id string) (*domain.Transaction, error)
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetAll(ctx context.Context, pagination domain.Pagination) (*domain.Page[domain.Transaction], error)
}
