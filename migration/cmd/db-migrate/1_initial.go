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

		_, err = db.Exec(`
			CREATE TABLE transactions (
				id UUID PRIMARY KEY,
				user_id UUID NOT NULL,
				from_user_id UUID,
				to_user_id UUID,
				currency VARCHAR(3) NOT NULL,
				amount BIGINT NOT NULL,
				balance BIGINT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				deleted_at TIMESTAMP,
				FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
				FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE
			);
		CREATE INDEX idx_transactions_user_id ON transactions (user_id);
		CREATE INDEX idx_transactions_deleted_at ON transactions (deleted_at);
		`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping tables...")
		_, err := db.Exec(`
		DROP TABLE IF EXISTS transactions;
		DROP TABLE IF EXISTS users;
		`)
		return err
	})
}
