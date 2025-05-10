package main

import (
	"errors"
	"net/http"
	"os"
	"utils"

	"github.com/JoshPattman/jpf"
	"github.com/gin-gonic/gin"
)

func embedTextHandler(ctx *gin.Context) {
	var req utils.EmbedStringRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	embedded, err := embedText(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, embedded)
}

func embedText(text utils.EmbedStringRequest) ([]float64, error) {
	embedder := jpf.NewOpenAIEmbedder(os.Getenv("OPENAI_KEY"), "text-embedding-3-small")
	vec, err := embedder.Embed(text.String)
	if err != nil {
		return nil, errors.Join(errors.New("failed to vectorise"), err)
	}
	return vec, nil
}
