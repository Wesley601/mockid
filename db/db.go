package db

import (
	"database/sql"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"

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

func EnsureDBFile(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		return nil
	}
	s := strings.Split(filePath, "/")
	if len(s) > 1 {
		p, file := path.Join(s[:len(s)-1]...), s[len(s)-1]
		if filepath.Ext(file) != ".db" {
			return errors.New("file must have the .db extension")
		}
		err := os.Mkdir(p, 0777)
		if !os.IsExist(err) {
			return err
		}
	}
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	return file.Close()
}
