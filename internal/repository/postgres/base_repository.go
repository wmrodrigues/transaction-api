package postgres

import (
	"context"
	"transaction-api/internal/database"

	"gorm.io/gorm"
)

type BaseRepository struct {
	db *gorm.DB
}

func (r *BaseRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := database.TransactionFromContext(ctx); ok {
		return tx
	}
	return r.db
}
