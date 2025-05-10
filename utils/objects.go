package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type EmbedStringRequest struct {
	String string `json:"string"`
}

type ImportFromURLRequest struct {
	URL string `json:"url"`
}

type Recipe struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Metric   string  `json:"metric"`
}

type RecipeImportInfo struct {
	Summary string    `json:"summary"`
	Vector  []float64 `json:"vector"`
}

type SemanticSearchRequest struct {
	Search string `json:"search"`
	MaxN   int    `json:"max_n"`
}

func ToJSONReader(v interface{}) *bytes.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}
	return bytes.NewReader(b)
}

func MustReadEnvInt(name string) int {
	valStr := os.Getenv(name)
	if valStr == "" {
		fmt.Println("must specify env var")
		os.Exit(1)
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		fmt.Println("bad int")
		os.Exit(1)
	}
	return val
}
