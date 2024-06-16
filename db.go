package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func StartDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	m, err := NewMigratorEmbed(db)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil {
		return nil, err
	}

	return db, nil
}
