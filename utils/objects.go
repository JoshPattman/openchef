package utils

import (
	"fmt"
	"os"
	"strconv"
)

type ImportFromURLRequest struct {
	URL string `json:"url"`
}

type Recipie struct {
	Name        string       `json:"name"`
	Ingredients []Ingredient `json:"ingredients"`
	Steps       []string     `json:"steps"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Metric   string  `json:"metric"`
}

type RecipieImportInfo struct {
	Summary string    `json:"summary"`
	Vector  []float64 `json:"vector"`
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
