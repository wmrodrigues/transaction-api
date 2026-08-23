package main

import (
	"fmt"
	"log"

	"github.com/go-pg/migrations/v8"
)

func init() {
	err := migrations.Register(

		func(db migrations.DB) error {
			fmt.Println("Creating initial data...")
			_, err := db.Exec(`	
			INSERT INTO users(id, name, email, password, active, created_at, updated_at, deleted_at)
			VALUES ('7f37e40d-ea0b-4cf0-9104-701f1737d145', 'wash1', 'wash1@gmail.com', '$2a$10$K4PpEiB8zNLEVaoWgOIEA.Gs4nf/cFSARt7P2amHE0lvYqAyGd7Gq', true, now(), now(), null);
			
			INSERT INTO users(id, name, email, password, active, created_at, updated_at, deleted_at)
			VALUES ('87a2b0f5-37a0-410d-ab23-59a3cb4fcf25', 'wash2', 'wash2@gmail.com', '$2a$10$dYXWBUqXDR948RhByX4Oe.xtd2pP7.EOrzsm.vpQuH8Pd2k3Vzij.', true, now(), now(), null);


			INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at, deleted_at)
			VALUES('b75a085d-9b34-4ade-afbc-d4049737f2f6', '7f37e40d-ea0b-4cf0-9104-701f1737d145', 'SGD', 1000, now(), now(), null);
			
			INSERT INTO accounts (id, user_id, currency, balance, created_at, updated_at, deleted_at)
			VALUES('def7bb60-ffd9-4bf6-8b54-af5032f2ed57', '87a2b0f5-37a0-410d-ab23-59a3cb4fcf25', 'SGD', 0, now(), now(), null);
			`)
			return err
		},

		func(db migrations.DB) error {
			_, err := db.Exec(`
				DELETE FROM accounts;
				DELETE FROM users;
			`)
			return err
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}
