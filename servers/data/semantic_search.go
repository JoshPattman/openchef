package main

import (
	"net/http"
	"os"
	"sort"
	"utils"

	"github.com/JoshPattman/jpf"
	"github.com/gin-gonic/gin"
)

func semanticSearchHandler(ctx *gin.Context) {
	req := utils.SemanticSearchRequest{}
	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	recs, err := semanticSearch(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, recs)
}

func semanticSearch(req utils.SemanticSearchRequest) ([]utils.Recipe, error) {
	recs, err := GetAllRecipes()
	if err != nil {
		panic(err)
	}
	embedder := jpf.NewOpenAIEmbedder(os.Getenv("OPENAI_KEY"), "text-embedding-3-small")
	query, err := embedder.Embed(req.Search)
	if err != nil {
		return nil, err
	}
	// TODO very inefficient
	sort.Slice(recs, func(i, j int) bool {
		c1, _ := jpf.CosineSimilarity(query, recs[i].Emb)
		c2, _ := jpf.CosineSimilarity(query, recs[j].Emb)
		return c1 > c2
	})

	result := make([]utils.Recipe, 0)
	for i := range recs {
		if i >= req.MaxN {
			break
		}
		result = append(result, recs[i].Rec)
	}
	return result, nil
}
