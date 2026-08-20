package main

import (
	"log"

	"github.com/go-pg/migrations/v8"
)

func init() {
	err := migrations.Register(func(db migrations.DB) error {
		_, err := db.Exec(`
			
		`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`

		`)
		return err
	})

	if err != nil {
		log.Fatal(err)
	}
}
