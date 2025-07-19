package event

import (
	"database/sql"
)

type EventHandler struct {
	db *sql.DB
}

func NeweventHandler(db *sql.DB) *EventHandler {
	return &EventHandler{
		db: db,
	}
}
