package utils

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
	Summary string `json:"summary"`
	Vector  []int  `json:"vector"`
}
