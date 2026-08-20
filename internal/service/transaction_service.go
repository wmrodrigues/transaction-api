package service

import (
	"context"
	"errors"
	"fmt"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"gorm.io/gorm"
)

type TransactionService interface {
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetBalanceByUserId(ctx context.Context, userId string) (*domain.Transaction, error)
	Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	GetAll(ctx context.Context, pagination domain.Pagination) (*domain.Page[domain.Transaction], error)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(transactionRepository repository.TransactionRepository) TransactionService {
	return &transactionService{transactionRepository: transactionRepository}
}

func (t *transactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("transaction id cannot be empty and must be a valid UUID")
	}
	transaction, err := t.transactionRepository.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("transaction not found")
	}
	return transaction, err
}

func (t *transactionService) GetBalanceByUserId(ctx context.Context, userId string) (*domain.Transaction, error) {
	if userId == "" {
		return nil, errors.New("user id cannot be empty and must be a valid UUID")
	}
	transaction, err := t.transactionRepository.GetBalanceByUserId(ctx, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("balance not found for specified user")
	}
	return transaction, err
}

func (t *transactionService) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	if err := transaction.Validate(); err != nil {
		return nil, fmt.Errorf("error validating transaction data: %w", err)
	}

	return nil, nil

}

func (t *transactionService) GetAll(ctx context.Context, pagination domain.Pagination) (*domain.Page[domain.Transaction], error) {
	return t.transactionRepository.GetAll(ctx, pagination)
}
