package main

import (
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
)

func advancedInfoHandler(ctx *gin.Context) {
	var recipie utils.Recipie
	err := ctx.BindJSON(&recipie)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	imported, err := advancedInfoFromRecipie(recipie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, imported)
}

func advancedInfoFromRecipie(utils.Recipie) (utils.RecipieImportInfo, error) {
	return utils.RecipieImportInfo{
		Summary: "A recipie",
		Vector:  []float64{0},
	}, nil
}
