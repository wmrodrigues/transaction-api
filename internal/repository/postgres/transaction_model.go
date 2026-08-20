package postgres

import (
	"time"
	"transaction-api/internal/domain"

	"gorm.io/gorm"
)

type TransactionModel struct {
	ID         string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string         `gorm:"type:uuid" json:"user_id"`
	FromUserID string         `gorm:"type:uuid" json:"from_user_id"`
	Currency   string         `gorm:"type:varchar(3)" json:"currency"`
	Amount     int64          `gorm:"type:bigint" json:"amount"`
	Balance    int64          `gorm:"type:bigint" json:"balance"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (t *TransactionModel) TableName() string {
	return "transactions"
}

func (t *TransactionModel) toDomain() domain.Transaction {
	return domain.Transaction{
		ID:         t.ID,
		UserID:     t.UserID,
		FromUserID: t.FromUserID,
		Currency:   t.Currency,
		Amount:     t.Amount,
		Balance:    t.Balance,
		CreatedAt:  t.CreatedAt,
	}
}

func TransactionToModel(transaction *domain.Transaction) TransactionModel {
	return TransactionModel{
		ID:         transaction.ID,
		UserID:     transaction.UserID,
		FromUserID: transaction.FromUserID,
		Currency:   transaction.Currency,
		Amount:     transaction.Amount,
		Balance:    transaction.Balance,
	}
}
