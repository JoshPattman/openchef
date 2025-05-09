package utils

type Recipie struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
	Steps       []string `json:"steps"`
}

type RecipieImportInfo struct {
	Summary string `json:"summary"`
	Vector  []int  `json:"vector"`
}
