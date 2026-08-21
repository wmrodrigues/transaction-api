package postgres

import (
	"context"
	"fmt"
	"math"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) repository.TransactionRepository {
	return &TransactionRepository{db: db}
}

func (t *TransactionRepository) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var model TransactionModel
	err := t.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("error getting transaction by id: %w", err)
	}
	transaction := model.toDomain()
	return &transaction, nil
}

func (t *TransactionRepository) GetBalanceByUserId(ctx context.Context, id string) (*domain.Transaction, error) {
	var model TransactionModel
	err := t.db.WithContext(ctx).Where("user_id = ?", id).Order("created_at DESC").First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("error getting balance by user id: %w", err)
	}
	transaction := model.toDomain()
	return &transaction, nil
}

func (t *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	model := TransactionToModel(transaction)
	err := t.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return fmt.Errorf("error creating transaction: %w", err)
	}
	return err
}

func (t *TransactionRepository) GetByUserId(ctx context.Context, userId string, pagination domain.Pagination) (*domain.Page[domain.Transaction], error) {
	if pagination.Page <= 0 {
		pagination.Page = 0
	}

	if pagination.PageSize <= 0 {
		pagination.PageSize = 10
	}

	offset := pagination.Page * pagination.PageSize

	var total int64

	err := t.db.
		WithContext(ctx).
		Model(&TransactionModel{}).
		Count(&total).
		Error

	if err != nil {
		return nil, err
	}

	var models []TransactionModel
	err = t.db.
		WithContext(ctx).
		Where("user_id = ?", userId).
		Order("created_at").
		Limit(pagination.PageSize).
		Offset(offset).
		Find(&models).
		Error

	if err != nil {
		return nil, err
	}

	transactions := make([]domain.Transaction, 0, len(models))

	for _, model := range models {
		transactions = append(transactions, model.toDomain())
	}

	totalPages := int(math.Ceil(float64(total) / float64(pagination.PageSize)))

	return &domain.Page[domain.Transaction]{
		Items:      transactions,
		Total:      total,
		Page:       pagination.Page,
		PageSize:   pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}
