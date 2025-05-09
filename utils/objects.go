package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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

func ToJSONReader(v interface{}) *bytes.Reader {
	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal JSON: %v", err))
	}
	return bytes.NewReader(b)
}
