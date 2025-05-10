package main

import (
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
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	id, err := InsertRecipe(imported, importRequest.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	go func() {
		advancedInfo, err := advancedInfoFromRecipe(imported)
		if err != nil {
			logger.Error("failed to do advanced import", "id", id, "error", err)
			return
		}
		err = AddAdvancedInfo(id, advancedInfo)
		if err != nil {
			logger.Error("failed to insert advanced info", "id", id, "error", err)
			return
		}
	}()

	ctx.JSON(http.StatusOK, imported)
}
