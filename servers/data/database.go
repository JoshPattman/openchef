package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"utils"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/create_tables.sql
var createTableSql string

//go:embed sql/insert_recipe.sql
var insertRecipeSql string

//go:embed sql/update_advanced_info.sql
var updateAdvancedInfoSql string

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path.Join(os.Getenv("DATA_PERSIST_PATH"), "web.db"))
	if err != nil {
		return nil, errors.Join(fmt.Errorf("Error opening database"), err)
	}
	logger.Info("Connected to db")
	return db, nil
}

func InitDB(db *sql.DB) error {
	_, err := db.Exec(createTableSql)
	if err != nil {
		return errors.Join(fmt.Errorf("Error initializing database"), err)
	}
	logger.Info("Inited db")
	return nil
}

func InsertRecipe(rec utils.Recipe, url string) (int, error) {
	stepsJson, err := json.Marshal(rec.Ingredients)
	if err != nil {
		return -1, err
	}
	ingredsJson, err := json.Marshal(rec.Steps)
	if err != nil {
		return -1, err
	}
	resp := db.QueryRow(
		insertRecipeSql,
		rec.Name,
		url,
		string(stepsJson),
		string(ingredsJson),
	)
	var id int
	err = resp.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AddAdvancedInfo(id int, advancedInfo utils.RecipeImportInfo) error {
	vectorJson, err := json.Marshal(advancedInfo.Vector)
	if err != nil {
		return err
	}
	_, err = db.Exec(
		updateAdvancedInfoSql,
		advancedInfo.Summary,
		vectorJson,
		id,
	)
	return err
}
