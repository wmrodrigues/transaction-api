package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"transaction-api/internal/config"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

const usageText = `This program runs command on the db. Supported commands are:
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
  go run *.go <command> [args]
`

var CommitVersion string

func main() {
	log.Printf("starting migration service, commit: %s\n", CommitVersion)
	flag.Usage = printUsage
	flag.Parse()

	db, err := connectDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	migTableExists, err := migrationsTableExists(db)
	if err != nil {
		log.Fatalf("failed to check if migrations table exists: %v", err)
	}
	if !migTableExists {
		if _, _, err := migrations.Run(db, "init"); err != nil {
			log.Fatalf("failed to create migrations table: %v", err)
		}
	}

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if newVersion != oldVersion {
		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		log.Printf("version is %d. Nothing to do\n", oldVersion)
	}
}

func printUsage() {
	fmt.Print(usageText)
	flag.PrintDefaults()
	os.Exit(2)
}

func connectDatabase() (*pg.DB, error) {
	cfg := config.LoadConfigs()
	opts := &pg.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.PostgresAddr, cfg.PostgresPort),
		User:     cfg.PostgresUser,
		Password: cfg.PostgresPass,
		Database: cfg.PostgresDB,
	}
	db := pg.Connect(opts)

	//attempt to ping the database to verify the connection 3 times
	for i := 0; i < 3; i++ {
		if err := db.Ping(context.Background()); err == nil {
			log.Println("connected to database")
			return db, nil
		} else {
			log.Printf("failed to connect to database with url %s: %s, retrying in 5s", opts.ToURL(), err.Error())
			time.Sleep(5 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to connect to database")
}

func migrationsTableExists(db *pg.DB) (bool, error) {
	resp, err := db.Exec(`SELECT 1 FROM pg_tables WHERE tablename = 'gopg_migrations'`)
	if err != nil {
		return false, err
	}
	return resp.RowsReturned() == 1, nil
}
