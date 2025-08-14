package repository

import (
	"database/sql"
)

type BassineRepository struct {
	DB *sql.DB
}

func NewBassineRepository(db *sql.DB) *BassineRepository {
	return &BassineRepository{
		DB: db,
	}
}
