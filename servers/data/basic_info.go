package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"utils"

	"github.com/JoshPattman/jpf"
)

type RecipeSchema struct {
	Context        string   `json:"@context"`
	Type           string   `json:"@type"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Image          any      `json:"image"`  // Can be string or ImageObject
	Author         any      `json:"author"` // Can be string or Person/Organization
	DatePublished  string   `json:"datePublished"`
	DateModified   string   `json:"dateModified"`
	PrepTime       string   `json:"prepTime"`  // ISO 8601 duration format
	CookTime       string   `json:"cookTime"`  // ISO 8601 duration format
	TotalTime      string   `json:"totalTime"` // ISO 8601 duration format
	Keywords       string   `json:"keywords"`
	RecipeCategory string   `json:"recipeCategory"`
	RecipeCuisine  string   `json:"recipeCuisine"`
	RecipeYield    string   `json:"recipeYield"` // Can be string or number
	Ingredients    []string `json:"recipeIngredient"`
	Instructions   any      `json:"recipeInstructions"` // Can be string[] or HowToStep[]
	//SuitableForDiet []string `json:"suitableForDiet"`
	Nutrition       any `json:"nutrition"`       // NutritionInformation
	AggregateRating any `json:"aggregateRating"` // AggregateRating
	Video           any `json:"video"`           // VideoObject
}

func basicInfoFromRequest(urlImport utils.ImportFromURLRequest) (utils.Recipe, error) {
	logger.Info("Fetching URL", "url", urlImport)

	// Make HTTP GET request to fetch the HTML content
	resp, err := http.Get(urlImport.URL)
	if err != nil {
		return utils.Recipe{}, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return utils.Recipe{}, fmt.Errorf("received non-OK response: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.Recipe{}, fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Info("Fetched bytes", "bytes", len(body), "url", urlImport.URL)

	// Extract structured recipe data from the HTML
	recipeSchemaString, err := ExtractRecipeSchema(string(body))
	if err != nil {
		logger.Warn("Could not extract recipe schema", "error", err)
		// Fall back to a default recipe if we can't extract structured data
		return utils.Recipe{
			Name: "Recipe from " + urlImport.URL,
		}, nil
	}

	logger.Debug("Found Schema", "schema", recipeSchemaString)

	// // Convert the recipe schema string into a recipe schema
	// var recipeSchema RecipeSchema
	// err = json.Unmarshal([]byte(recipeSchemaString), &recipeSchema)
	// if err != nil {
	// 	logger.Warn("Failed to parse recipe schema", "error", err)
	// 	return utils.Recipe{
	// 		Name: "Recipe from " + urlImport.URL,
	// 	}, nil
	// }

	// logger.Debug("Parsed schema", "parsed", recipeSchema)

	// // Convert the schema to our Recipe type
	// // Cursed but works for now for llm
	// steps, err := json.Marshal(recipeSchema.Instructions)
	// if err != nil {
	// 	return utils.Recipe{}, err
	// }
	// ingreds, err := json.Marshal(recipeSchema.Ingredients)
	// if err != nil {
	// 	return utils.Recipe{}, err
	// }
	// recipe := utils.Recipe{
	// 	Name:        recipeSchema.Name,
	// 	Steps:       []string{string(steps)},
	// 	Ingredients: []utils.Ingredient{{Name: string(ingreds)}},
	// }
	recipe, _, err := jpf.RunWithRetries(
		jpf.NewStandardOpenAIModel(os.Getenv("OPENAI_KEY"), "gpt-4o-mini", 0, 0, 0.0),
		NewRecipeNormaliser(),
		5,
		jpf.UserRole,
		recipeSchemaString,
	)
	if err != nil {
		return utils.Recipe{}, err
	}

	logger.Debug("Found basic recipe", "recipe", recipe)

	return recipe, nil
}

// ExtractRecipeSchema attempts to extract recipe schema from HTML content
func ExtractRecipeSchema(htmlContent string) (string, error) {
	// Look for JSON-LD script tags with more flexible pattern matching
	// This will match both <script type="application/ld+json"> and variants like
	// <script data-rh="true" type="application/ld+json">
	scriptPattern := `<script[^>]*type="application/ld\+json"[^>]*>`
	jsonldEnd := `</script>`

	// Find all script tags that might contain JSON-LD
	scriptStartRegex := regexp.MustCompile(scriptPattern)
	matches := scriptStartRegex.FindAllStringIndex(htmlContent, -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("no JSON-LD script tags found in the HTML content")
	}

	// Try each JSON-LD script tag until we find a valid recipe
	for _, match := range matches {
		start := match[1] // End index of the opening script tag

		// Find the end of this script tag
		endTagIndex := strings.Index(htmlContent[start:], jsonldEnd)
		if endTagIndex == -1 {
			continue // Skip this match if no closing tag is found
		}

		// Extract the JSON content
		jsonContent := htmlContent[start : start+endTagIndex]
		jsonContent = strings.TrimSpace(jsonContent)

		// Parse the JSON content
		var jsonData map[string]any
		if err := json.Unmarshal([]byte(jsonContent), &jsonData); err != nil {
			logger.Debug("Failed to parse JSON-LD", "error", err)
			continue
		}

		// Look for recipe objects in the JSON
		recipes := findRecipeObjects(jsonData)
		if len(recipes) > 0 {
			// Return the first recipe found
			recipeJSON, err := json.Marshal(recipes[0])
			if err != nil {
				continue
			}
			return string(recipeJSON), nil
		}
	}

	return "", fmt.Errorf("no valid recipe JSON-LD found in any script tag")
}

// findRecipeObjects recursively searches for objects with @type = "Recipe" in JSON data
func findRecipeObjects(data any) []map[string]any {
	var recipes []map[string]any

	switch v := data.(type) {
	case map[string]any:
		// Check if this object is a recipe
		if typeValue, ok := v["@type"].(string); ok {
			if typeValue == "Recipe" || strings.Contains(strings.ToLower(typeValue), "recipe") {
				recipes = append(recipes, v)
				return recipes
			}
		}

		// Check if there's a @graph property that might contain recipes
		if graph, ok := v["@graph"].([]any); ok {
			for _, item := range graph {
				recipes = append(recipes, findRecipeObjects(item)...)
			}
		} else {
			// Recursively search all properties
			for _, propValue := range v {
				recipes = append(recipes, findRecipeObjects(propValue)...)
			}
		}

	case []any:
		// Search through array elements
		for _, item := range v {
			recipes = append(recipes, findRecipeObjects(item)...)
		}
	}

	return recipes
}

var rschemaPrompt = `
- Convert the provided recipe into a JSON object that conforms to this schema:
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "The name of the recipe."
    },
    "ingredients": {
      "type": "array",
      "description": "A list of ingredients used in the recipe.",
      "items": {
        "type": "object",
        "properties": {
          "name": {
            "type": "string",
            "description": "The name of the ingredient."
          },
          "quantity": {
            "type": "number",
            "description": "The quantity of the ingredient."
          },
          "metric": {
            "type": "string",
            "description": "The unit of measurement for the ingredient."
          }
        },
        "required": ["name", "quantity", "metric"]
      }
    },
    "steps": {
      "type": "array",
      "description": "The steps to prepare the recipe.",
      "items": {
        "type": "string"
      }
    },
    "yield": {
      "type": "integer",
      "description": "The number of servings the recipe produces."
    }
  },
  "required": ["name", "ingredients", "steps", "yield"]
}
- The provided schema may provide a lot of information, you only need to include enough information to fill out the schema.
- Try not to re-word any sentences.
- Don't include any extra text other than the json object.
`

// Creates a function that attemps to convert a string into a recipe schema
func NewRecipeNormaliser() jpf.RetryFunction[string, utils.Recipe] {
	return &recipeNormaliser{}
}

type recipeNormaliser struct{}

// BuildInputMessages implements jpf.RetryFunction.
func (r *recipeNormaliser) BuildInputMessages(s string) ([]jpf.Message, error) {
	return []jpf.Message{
		{
			Role:    jpf.SystemRole,
			Content: rschemaPrompt,
		},
		{
			Role:    jpf.UserRole,
			Content: s,
		},
	}, nil
}

// FormatFeedback implements jpf.RetryFunction.
func (r *recipeNormaliser) FormatFeedback(err *jpf.ParseError) string {
	logger.Warn("Recipe Normaliser responded with invalid format", "error", err.Err, "format", err.Response)
	return fmt.Sprintf("An error occured with your last response: %v", err.Err)
}

// ParseResponseText implements jpf.RetryFunction.
func (r *recipeNormaliser) ParseResponseText(resp string) (utils.Recipe, error) {
	var rec utils.Recipe
	err := json.Unmarshal([]byte(resp), &rec)
	if err != nil {
		return utils.Recipe{}, err
	}
	if rec.Steps == nil {
		return utils.Recipe{}, fmt.Errorf("the steps object was nil, was the response in the wrong format?")
	}
	logger.Debug("Recipe normaliser response", "response", resp)
	return rec, nil
}
