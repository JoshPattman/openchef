package main

import (
	"datadb"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
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

func semanticSearch(req utils.SemanticSearchRequest) (utils.SemanticSearchResponse, error) {
	// TODO this all should probably use a single tx
	recs, err := datadb.GetAllIDVectorPairs(db)
	if err != nil {
		return utils.SemanticSearchResponse{}, err
	}
	embedder := jpf.NewOpenAIEmbedder(os.Getenv("OPENAI_KEY"), "text-embedding-3-small")
	query, err := embedder.Embed(req.Search)
	if err != nil {
		return utils.SemanticSearchResponse{}, err
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
			return utils.SemanticSearchResponse{}, err
		}
		results = append(results, utils.SemanticSearchResult{
			ID:      id,
			Name:    name,
			Summary: summary,
		})
	}

	sum := NewResultsSummariser()
	model := jpf.NewStandardOpenAIModel(os.Getenv("OPENAI_KEY"), "gpt-4o-mini", 0, 0, 0.5)
	summary, _, err := jpf.RunOneShot(model, sum, ResultSummariserInput{
		Query:   req.Search,
		Results: results,
	})
	if err != nil {
		return utils.SemanticSearchResponse{}, err
	}

	return utils.SemanticSearchResponse{
		Results: results,
		Summary: summary,
	}, nil
}

type ResultSummariserInput struct {
	Query   string
	Results []utils.SemanticSearchResult
}

func NewResultsSummariser() jpf.Function[ResultSummariserInput, string] {
	return &resultsSummariser{}
}

type resultsSummariser struct{}

// BuildInputMessages implements jpf.Function.
func (r *resultsSummariser) BuildInputMessages(input ResultSummariserInput) ([]jpf.Message, error) {
	sysPrompt := `
	- You are recipe summariser, an expert foodie.
	- Your role is to summarise some recipes for the user, given a query that they asked.
	- You should cite the IDs of the dishes you mention in your answer, using sqaure backets, e.g. "[5]" or "[1][3]".
	`
	userPrompt := fmt.Sprintf("I am looking for: %s", input.Query)
	resultStrs := []string{}
	for _, res := range input.Results {
		resultStrs = append(resultStrs, fmt.Sprintf("- **%s**: %s", res.Name, res.Summary))
	}
	resultPrompt := fmt.Sprintf("The following are relevant recipies to the users query:\n%s", strings.Join(resultStrs, "\n"))
	return []jpf.Message{
		{
			Role:    jpf.SystemRole,
			Content: sysPrompt,
		},
		{
			Role:    jpf.UserRole,
			Content: userPrompt,
		},
		{
			Role:    jpf.SystemRole,
			Content: resultPrompt,
		},
	}, nil
}

// ParseResponseText implements jpf.Function.
func (r *resultsSummariser) ParseResponseText(s string) (string, error) {
	return s, nil
}
