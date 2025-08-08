package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func Open(dsn string, migrationsDir string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	log.Println("Running database migrations...")
	if err := goose.SetDialect("postgres"); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := goose.Up(db, migrationsDir); err != nil {
		_ = db.Close()
		return nil, err
	}
	log.Println("Database migrations completed successfully")

	return db, nil
}
