package datadb

import (
	"database/sql"
	"errors"
)

// Adds a unit to database and returns ID
func AddUnit(db DB, name string, shorthand string) (int, error) {
	q := `INSERT INTO units (name, shorthand) VALUES (?, ?) RETURNING id;`
	row := db.QueryRow(q, name, shorthand)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.Join(errors.New("failed to insert unit"), err)
	}
	return id, nil
}

// Looks up a unit by name, returning ID, if it was found, and an error
func LookupUnitByName(db DB, name string) (int, bool, error) {
	q := `SELECT id FROM units WHERE name=?;`
	row := db.QueryRow(q, name)
	var id int
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return -1, false, nil
	} else if err != nil {
		return -1, false, errors.Join(errors.New("failed to look up unit"), err)
	}
	return id, true, nil
}
