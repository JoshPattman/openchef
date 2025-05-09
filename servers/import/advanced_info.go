package main

import (
	"encoding/json"
	"net/http"
	"os"
	"utils"

	_ "embed"

	"github.com/JoshPattman/jpf"
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

func advancedInfoFromRecipe(recipie utils.Recipe) (utils.RecipeImportInfo, error) {
	model := jpf.NewStandardOpenAIModel(os.Getenv("OPENAI_KEY"), "gpt-4o-mini", 0, 0, 0.4)
	summariser := NewRecipieSummariser()
	summary, _, err := jpf.RunOneShot(model, summariser, recipie)
	if err != nil {
		return utils.RecipeImportInfo{}, err
	}
	embedder := jpf.NewOpenAIEmbedder(os.Getenv("OPENAI_KEY"), "text-embedding-3-small")
	vec, err := embedder.Embed(summary)
	if err != nil {
		return utils.RecipeImportInfo{}, err
	}
	return utils.RecipeImportInfo{
		Summary: summary,
		Vector:  vec,
	}, nil
}

var _ jpf.Function[utils.Recipe, string] = &recipieSummariser{}

func NewRecipieSummariser() jpf.Function[utils.Recipe, string] {
	return &recipieSummariser{}
}

//go:embed prompts/recipie_summarise.md
var recipieSummariserPrompt string

type recipieSummariser struct{}

// BuildInputMessages implements jpf.Function.
func (r *recipieSummariser) BuildInputMessages(recipie utils.Recipe) ([]jpf.Message, error) {
	content, err := json.MarshalIndent(recipie, "", "  ")
	if err != nil {
		return nil, err
	}
	return []jpf.Message{
		{
			Role:    jpf.SystemRole,
			Content: recipieSummariserPrompt,
		},
		{
			Role:    jpf.UserRole,
			Content: string(content),
		},
	}, nil
}

// ParseResponseText implements jpf.Function.
func (r *recipieSummariser) ParseResponseText(s string) (string, error) {
	return s, nil
}
