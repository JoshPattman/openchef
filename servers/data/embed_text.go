package main

import (
	"errors"
	"os"
	"utils"

	"github.com/JoshPattman/jpf"
)

func embedText(text utils.EmbedStringRequest) ([]float64, error) {
	embedder := jpf.NewOpenAIEmbedder(os.Getenv("OPENAI_KEY"), "text-embedding-3-small")
	vec, err := embedder.Embed(text.String)
	if err != nil {
		return nil, errors.Join(errors.New("failed to vectorise"), err)
	}
	return vec, nil
}
