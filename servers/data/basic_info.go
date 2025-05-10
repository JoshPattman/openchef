package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"utils"

	"github.com/gin-gonic/gin"
)

// JSON represents a simplified JSON structure
type JSON map[string]any

func basicInfoHandler(ctx *gin.Context) {
	var importRequest utils.ImportFromURLRequest
	err := ctx.BindJSON(&importRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	imported, err := basicInfoFromRequest(importRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, imported)
}

func basicInfoFromRequest(urlImport utils.ImportFromURLRequest) (utils.Recipe, error) {
	logger.Info("Fetching URL", "url", urlImport)

	// Make HTTP GET request to fetch the HTML content
	resp, err := http.Get(urlImport.URL)
	if err != nil {
		return utils.Recipe{}, fmt.Errorf("failed to fetch URL", "error", err)
	}
	defer resp.Body.Close()

	// Check if the response status code is OK
	if resp.StatusCode != http.StatusOK {
		return utils.Recipe{}, fmt.Errorf("received non-OK response", "code", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.Recipe{}, fmt.Errorf("failed to read response body", "error", err)
	}

	logger.Info("Fetched bytes", "bytes", len(body), "url", urlImport.URL)

	// Extract structured recipe data from the HTML
	recipeSchema, err := ExtractRecipeSchema(string(body))
	if err != nil {
		logger.Warn("Could not extract recipe schema", "error", err)
		// Fall back to a default recipe if we can't extract structured data
		return utils.Recipe{
			Name: "Recipe from " + urlImport.URL,
		}, nil
	}

	logger.Debug("Found Schema", "schema", recipeSchema)

	// Convert the schema to our Recipe type
	recipe := utils.Recipe{
		Name: urlImport.URL,
	}

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

		return jsonContent, nil
	}

	return "", fmt.Errorf("no valid recipe JSON-LD found in any script tag")
}
