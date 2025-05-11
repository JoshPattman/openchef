package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path.Join(os.Getenv("DATA_PERSIST_PATH"), "web.db"))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error opening database"), err)
	}
	logger.Info("Connected to db")
	return db, nil
}
