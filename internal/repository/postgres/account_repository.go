package postgres

import (
	"context"
	"fmt"
	"transaction-api/internal/domain"
	"transaction-api/internal/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	BaseRepository
}

func NewAccountRepository(db *gorm.DB) repository.AccountRepository {
	return &AccountRepository{BaseRepository: BaseRepository{db: db}}
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	model := AccountToModel(account)
	err := r.getDB(ctx).Create(&model).Error
	if err != nil {
		return fmt.Errorf("error creating account: %w", err)
	}
	return nil
}

func (r *AccountRepository) GetForUpdate(ctx context.Context, userID, currency string) (*domain.Account, error) {
	var model AccountModel
	// here we load the account and lock the row FOR UPDATE so concurrent transactions serialize on it
	err := r.getDB(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND currency = ?", userID, currency).
		First(&model).Error
	if err != nil {
		return nil, fmt.Errorf("error locking account for update: %w", err)
	}
	account := model.toDomain()
	return &account, nil
}

func (r *AccountRepository) GetByUserID(ctx context.Context, userID string) ([]domain.Account, error) {
	var models []AccountModel
	// here we should improve it once we make the api available for other currencies
	err := r.getDB(ctx).Where("user_id = ?", userID).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("error getting accounts by user id: %w", err)
	}
	accounts := make([]domain.Account, 0, len(models))
	for _, model := range models {
		accounts = append(accounts, model.toDomain())
	}
	return accounts, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, accountID string, newBalance int64) error {
	err := r.getDB(ctx).
		Model(&AccountModel{}).
		Where("id = ?", accountID).
		Update("balance", newBalance).Error
	if err != nil {
		return fmt.Errorf("error updating account balance: %w", err)
	}
	return nil
}
