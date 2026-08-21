package database

import (
	"context"

	"gorm.io/gorm"
)

type GormManager struct {
	db *gorm.DB
}

func NewGormManager(db *gorm.DB) *GormManager {
	return &GormManager{db: db}
}

func (m *GormManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		transactionContext := ContextWithTransaction(ctx, tx)
		return fn(transactionContext)
	})
}
