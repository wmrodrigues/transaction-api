package database

import (
	"context"

	"gorm.io/gorm"
)

type transactionKey struct{}

func ContextWithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, transactionKey{}, tx)
}

func TransactionFromContext(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	return tx, ok
}
