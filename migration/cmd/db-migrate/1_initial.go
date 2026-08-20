package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating all database initial tables...")
		_, err := db.Exec(`
			CREATE TABLE users (
			id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		);

		CREATE UNIQUE INDEX idx_users_email ON users (email);
		CREATE INDEX idx_users_deleted_at ON users (deleted_at);
		`)
		if err != nil {
			return fmt.Errorf("error creating users table: %v", err)
		}
		return nil
	}, func(db migrations.DB) error {
		fmt.Println("dropping tables...")
		_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
		`)
		return err
	})
}
