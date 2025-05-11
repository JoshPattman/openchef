package datadb

import (
	"database/sql"
	"errors"
)

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func InitDB(db DB) error {
	q := `
		CREATE TABLE IF NOT EXISTS recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			source TEXT,
			state TEXT
		);
		CREATE TABLE IF NOT EXISTS ai_recipes (
			recipe_id INTEGER PRIMARY KEY,
			summary TEXT,
			summary_embedding BLOB,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id)
		);
		CREATE TABLE IF NOT EXISTS steps (
			id INTEGER PRIMARY KEY,
			recipe_id INTEGER,
			place INTEGER,
			step TEXT NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id)
		);
		CREATE TABLE IF NOT EXISTS ingredients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT
		);
		CREATE TABLE IF NOT EXISTS units (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			shorthand TEXT
		);
		CREATE TABLE IF NOT EXISTS recipe_ingredients (
			recipe_id INTEGER,
			ingredient_id INTEGER,
			unit_id INTEGER NOT NULL,
			amount REAL,
			PRIMARY KEY (recipe_id, ingredient_id),
			FOREIGN KEY (recipe_id) REFERENCES recipes(id),
			FOREIGN KEY (ingredient_id) REFERENCES ingredients(id),
			FOREIGN KEY (unit_id) REFERENCES units(id)
		);
	`
	_, err := db.Exec(q)
	if err != nil {
		return errors.Join(errors.New("failed to init db"), err)
	}
	return nil
}
