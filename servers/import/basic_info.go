package main

import (
	"net/http"
	"utils"

	"github.com/gin-gonic/gin"
)

func basicInfoHandler(ctx *gin.Context) {
	var importRequest utils.ImportFromURLRequest
	err := ctx.BindJSON(&importRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	imported, err := basicInfoFromRequest(importRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, imported)
}

func basicInfoFromRequest(utils.ImportFromURLRequest) (utils.Recipie, error) {
	return utils.Recipie{
		Name: "Good soup",
	}, nil
}
