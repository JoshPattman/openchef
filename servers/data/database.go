package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/create_tables.sql
var createTableSql string

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path.Join(os.Getenv("DATA_PERSIST_PATH"), "web.db"))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error opening database"), err)
	}
	return db, nil
}

func InitDB(db *sql.DB) error {
	_, err := db.Exec(createTableSql)
	if err != nil {
		return errors.Join(fmt.Errorf("Error initializing database"), err)
	}
	return nil
}
