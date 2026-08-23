package postgres

import (
	"time"
	"transaction-api/internal/domain"

	"gorm.io/gorm"
)

type AccountModel struct {
	ID        string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string         `gorm:"type:uuid;not null" json:"user_id"`
	Currency  string         `gorm:"type:varchar(3);not null" json:"currency"`
	Balance   int64          `gorm:"type:bigint;not null;default:0" json:"balance"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (a *AccountModel) TableName() string {
	return "accounts"
}

func (a *AccountModel) toDomain() domain.Account {
	return domain.Account{
		ID:       a.ID,
		UserID:   a.UserID,
		Currency: a.Currency,
		Balance:  a.Balance,
	}
}

func AccountToModel(account *domain.Account) AccountModel {
	return AccountModel{
		ID:       account.ID,
		UserID:   account.UserID,
		Currency: account.Currency,
		Balance:  account.Balance,
	}
}
