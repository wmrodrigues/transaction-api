package postgres

import (
	"time"
	"transaction-api/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionModel struct {
	ID         string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     string         `gorm:"type:uuid" json:"user_id"`
	FromUserID *string        `gorm:"type:uuid" json:"from_user_id"`
	ToUserID   *string        `gorm:"type:uuid" json:"to_user_id"`
	Currency   string         `gorm:"type:varchar(3)" json:"currency"`
	Amount     int64          `gorm:"type:bigint" json:"amount"`
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
		ToUserID:   t.ToUserID,
		Currency:   t.Currency,
		Amount:     t.Amount,
		CreatedAt:  t.CreatedAt,
	}
}

func TransactionToModel(transaction *domain.Transaction) TransactionModel {
	id := transaction.ID
	if transaction.ID == "" {
		id = uuid.New().String()
	}
	// default currency to SGD, but I'd handle this in the service using a localization header
	currency := transaction.Currency
	if currency == "" {
		currency = string(domain.SGD)
	}
	return TransactionModel{
		ID:         id,
		UserID:     transaction.UserID,
		FromUserID: transaction.FromUserID,
		ToUserID:   transaction.ToUserID,
		Currency:   currency,
		Amount:     transaction.Amount,
	}
}
