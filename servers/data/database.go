package main

import (
	"bytes"
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

//go:embed sql/get_all_recipes.sql
var getAllRecipesSql string

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

type RecipeEmbeddingPair struct {
	Rec utils.Recipe
	Emb []float64
}

func GetAllRecipes() ([]RecipeEmbeddingPair, error) {
	rows, err := db.Query(getAllRecipesSql)
	if err != nil {
		return nil, err
	}
	recs := make([]RecipeEmbeddingPair, 0)
	for rows.Next() {
		rec := utils.Recipe{}
		var ing, steps string
		var emBytes []byte
		var vec []float64
		err := rows.Scan(&rec.Name, &ing, &steps, &emBytes)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(bytes.NewBufferString(ing)).Decode(&rec.Ingredients)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(bytes.NewBufferString(steps)).Decode(&rec.Steps)
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(bytes.NewBuffer(emBytes)).Decode(&vec)
		if err != nil {
			continue // For now just skip as not imported
		}
		recs = append(recs, RecipeEmbeddingPair{rec, vec})
	}
	return recs, nil
}
