package main

import (
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
)

func importHandler(ctx *gin.Context) {
	var recipie utils.Recipie
	err := ctx.BindJSON(&recipie)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	imported, err := importRecipie(recipie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, imported)
}

func importRecipie(utils.Recipie) (utils.RecipieImportInfo, error) {
	return utils.RecipieImportInfo{
		Summary: "A recipie",
		Vector:  []float64{0},
	}, nil
}
