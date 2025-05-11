package main

import (
	"datadb"
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
)

func importURLHandler(ctx *gin.Context) {
	var importRequest utils.ImportFromURLRequest
	err := ctx.BindJSON(&importRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	imported, err := basicInfoFromRequest(importRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	id, err := AddRecipeToDB(imported, importRequest.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	go func() {
		advancedInfo, err := advancedInfoFromRecipe(imported)
		if err != nil {
			logger.Error("failed to do advanced import", "id", id, "error", err)
			return
		}
		err = AddAIRecipeToDB(id, advancedInfo)
		if err != nil {
			logger.Error("failed to insert advanced info", "id", id, "error", err)
			return
		}
	}()

	ctx.JSON(http.StatusOK, imported)
}

func AddRecipeToDB(rec utils.Recipe, url string) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	succsess := false
	defer func() {
		if !succsess {
			tx.Rollback()
		}
	}()

	recipeID, err := datadb.AddRecipe(tx, rec.Name, url, "import_pending")
	if err != nil {
		return -1, err
	}
	for i, step := range rec.Steps {
		_, err := datadb.AddStep(tx, recipeID, i, step)
		if err != nil {
			return -1, err
		}
	}
	// TODO should try to match these rather than creating duplicates
	ingredientIDs := make([]int, 0)
	unitIDs := make([]int, 0)
	for _, ing := range rec.Ingredients {
		ingID, err := datadb.AddIngredient(tx, ing.Name)
		if err != nil {
			return -1, err
		}
		ingredientIDs = append(ingredientIDs, ingID)
		unitID, err := datadb.AddUnit(tx, ing.Metric, ing.Metric)
		if err != nil {
			return -1, err
		}
		unitIDs = append(unitIDs, unitID)
	}
	for i, ing := range rec.Ingredients {
		err := datadb.AddRecipeIngredient(tx, recipeID, ingredientIDs[i], unitIDs[i], ing.Quantity)
		if err != nil {
			return -1, err
		}
	}

	succsess = true
	err = tx.Commit()
	if err != nil {
		succsess = false
		return -1, err
	}
	return recipeID, nil
}

func AddAIRecipeToDB(recID int, aiRec utils.RecipeImportInfo) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	succsess := false
	defer func() {
		if !succsess {
			tx.Rollback()
		}
	}()
	err = datadb.AddAIRecipeInfo(tx, recID, aiRec.Summary, aiRec.Vector)
	if err != nil {
		return err
	}
	// TODO maybe want to make state go to import failed, but not today
	err = datadb.UpdateRecipeState(tx, recID, "import_sucsess")
	if err != nil {
		return err
	}
	succsess = true
	err = tx.Commit()
	if err != nil {
		succsess = false
		return err
	}
	return nil
}
