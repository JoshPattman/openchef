package main

import (
	"datadb"
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

func semanticSearch(req utils.SemanticSearchRequest) ([]utils.SemanticSearchResult, error) {
	// TODO this all should probably use a single tx
	recs, err := datadb.GetAllIDVectorPairs(db)
	if err != nil {
		return nil, err
	}
	embedder := jpf.NewOpenAIEmbedder(os.Getenv("OPENAI_KEY"), "text-embedding-3-small")
	query, err := embedder.Embed(req.Search)
	if err != nil {
		return nil, err
	}
	// TODO very inefficient
	sort.Slice(recs, func(i, j int) bool {
		c1, _ := jpf.CosineSimilarity(query, recs[i].Vector)
		c2, _ := jpf.CosineSimilarity(query, recs[j].Vector)
		return c1 > c2
	})

	ids := make([]int, 0)
	for i := range recs {
		if i >= req.MaxN {
			break
		}
		ids = append(ids, recs[i].ID)
	}

	results := make([]utils.SemanticSearchResult, 0)
	for _, id := range ids {
		name, summary, err := datadb.GetRecipeNameAndSummary(db, id)
		if err != nil {
			return nil, err
		}
		results = append(results, utils.SemanticSearchResult{
			ID:      id,
			Name:    name,
			Summary: summary,
		})
	}

	return results, nil
}
