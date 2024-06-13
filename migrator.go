package main

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Migrator struct {
	db  *sql.DB
	dir string
}

func NewMigrator(db *sql.DB, dir string) (*Migrator, error) {
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return nil, err
	}

	return &Migrator{
		db:  db,
		dir: dir,
	}, nil
}

func NewMigratorEmbed(db *sql.DB) (*Migrator, error) {
	err := goose.SetDialect("sqlite3")
	goose.SetBaseFS(migrations)
	if err != nil {
		return nil, err
	}

	return &Migrator{
		db:  db,
		dir: "migrations",
	}, nil
}

func (m Migrator) Up() error {
	return goose.Up(m.db, m.dir)
}

func (m Migrator) Down() error {
	return goose.Down(m.db, m.dir)
}

func (m Migrator) Status() error {
	return goose.Status(m.db, m.dir)
}

func (m Migrator) Reset() error {
	return goose.Reset(m.db, m.dir)
}
