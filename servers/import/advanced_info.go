package main

import (
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
)

func advancedInfoHandler(ctx *gin.Context) {
	var recipe utils.Recipe
	err := ctx.BindJSON(&recipe)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	imported, err := advancedInfoFromRecipe(recipe)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, imported)
}

func advancedInfoFromRecipe(utils.Recipe) (utils.RecipeImportInfo, error) {
	return utils.RecipeImportInfo{
		Summary: "A recipie",
		Vector:  []float64{0},
	}, nil
}
