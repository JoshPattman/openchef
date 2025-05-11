package datadb

import (
	"encoding/json"
	"errors"
)

// Adds a recipe to the recipes table
func AddRecipe(db DB, name string, source string, state string) (int, error) {
	q := `INSERT INTO recipes (name, source, state) VALUES (?, ?, ?) RETURNING id;`
	row := db.QueryRow(q, name, source, state)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.Join(errors.New("failed to insert recipe"), err)
	}
	return id, nil
}

// Adds an AI recipe to the AI
func AddAIRecipeInfo(db DB, recipeID int, summary string, vector []float64) error {
	vectorBlob, err := json.Marshal(vector)
	if err != nil {
		return errors.Join(errors.New("failed to serialise vector"), err)
	}
	q := `INSERT INTO ai_recipes (summary, summary_embedding) VALUES (?, ?);`
	_, err = db.Exec(q, summary, vectorBlob)
	if err != nil {
		return errors.Join(errors.New("failed to insert ai recipe info"), err)
	}
	return nil
}

// Updates the state of a recipe id
func UpdateRecipeState(db DB, id int, state string) error {
	q := `UPDATE recipes SET state=? WHERE id=id;`
	_, err := db.Exec(q, state, id)
	if err != nil {
		return errors.Join(errors.New("failed to update recipe state"), err)
	}
	return nil
}

// Adds a single step
func AddStep(db DB, recipeID int, place int, step string) (int, error) {
	q := `INSERT INTO steps (recipe_id, place, step) VALUES (?, ?, ?) RETURNING id;`
	row := db.QueryRow(q, recipeID, place, step)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.Join(errors.New("failed to insert step"), err)
	}
	return id, nil
}

// Adds a single recipe ingredient
func AddRecipeIngredient(db DB, recipeID int, ingredientID int, unitID int, amount float64) error {
	q := `INSERT INTO recipe_ingredients (recipe_id, ingredient_id, unit_id, amount) VALUES (?, ?, ?, ?);`
	_, err := db.Exec(q, recipeID, ingredientID, unitID, amount)
	if err != nil {
		return errors.Join(errors.New("failed to insert recipe ingredient"), err)
	}
	return nil
}
