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

type IDVectorPair struct {
	ID     int
	Vector []float64
}

// Gets every id that has a vector
func GetAllIDVectorPairs(db DB) ([]IDVectorPair, error) {
	q := `SELECT recipe_id, summary_embedding FROM ai_recipes;`
	pairs := make([]IDVectorPair, 0)
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		var vecBlob []byte
		err := rows.Scan(&id, &vecBlob)
		if err != nil {
			return nil, err
		}
		var vec []float64
		err = json.Unmarshal(vecBlob, &vec)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, IDVectorPair{ID: id, Vector: vec})
	}
	return pairs, nil
}

func GetRecipeNameAndSummary(db DB, id int) (string, string, error) {
	q := `SELECT recipes.name, ai_recipes.summary FROM recipes JOIN ai_recipes ON recipes.id=ai_recipes.recipe_id WHERE recipes.id=?;`
	row := db.QueryRow(q, id)
	var name, summary string
	err := row.Scan(&name, &summary)
	if err != nil {
		return "", "", err
	}
	return name, summary, nil
}
