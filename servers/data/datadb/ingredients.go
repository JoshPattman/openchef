package datadb

import (
	"database/sql"
	"errors"
)

// Adds an ingredient to database and returns ID
func AddIngredient(db DB, name string) (int, error) {
	q := `INSERT INTO ingredients (name) VALUES (?) RETURNING id;`
	row := db.QueryRow(q, name)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.Join(errors.New("failed to insert ingredient"), err)
	}
	return id, nil
}

// Looks up an ingredient by name, returning ID, if it was found, and an error
func LookupIngredientByName(db DB, name string) (int, bool, error) {
	q := `SELECT id FROM ingredients WHERE name=?;`
	row := db.QueryRow(q, name)
	var id int
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, false, nil
	} else if err != nil {
		return -1, false, errors.Join(errors.New("failed to look up ingredient"), err)
	}
	return id, true, nil
}
