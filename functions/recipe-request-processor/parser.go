package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParserLogFunc is a function type for parser logging
type ParserLogFunc func(string)

// parserLogger is used for internal logging within parser
var parserLogger ParserLogFunc

// SetParserLogger sets the logger for parser functions
func SetParserLogger(logger ParserLogFunc) {
	parserLogger = logger
}

func logParser(msg string) {
	if parserLogger != nil {
		parserLogger("[Parser] " + msg)
	}
}

// extractRecipeFromHTML extracts Recipe JSON-LD from HTML content
func extractRecipeFromHTML(htmlContent string) (*Recipe, error) {
	// Log HTML content info for debugging
	logParser(fmt.Sprintf("HTML content length: %d bytes", len(htmlContent)))

	// Check if HTML contains script tags at all
	if strings.Contains(htmlContent, "<script") {
		logParser("HTML contains <script> tags")
	} else {
		logParser("WARNING: HTML does NOT contain any <script> tags")
	}

	// Check if HTML contains JSON-LD specifically
	if strings.Contains(htmlContent, "application/ld+json") {
		logParser("HTML contains application/ld+json reference")
	} else {
		logParser("WARNING: HTML does NOT contain application/ld+json")
	}

	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Count all script tags
	allScripts := doc.Find("script").Length()
	jsonLDScripts := doc.Find("script[type='application/ld+json']").Length()
	logParser(fmt.Sprintf("Found %d total script tags, %d JSON-LD scripts", allScripts, jsonLDScripts))

	// Find all JSON-LD script tags
	var recipe *Recipe
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		if recipe != nil {
			return // Already found a recipe
		}

		jsonLD := s.Text()
		if jsonLD == "" {
			logParser(fmt.Sprintf("JSON-LD script #%d is empty", i))
			return
		}

		logParser(fmt.Sprintf("JSON-LD script #%d length: %d bytes", i, len(jsonLD)))

		// Try to parse as single object or array
		var data interface{}
		if err := json.Unmarshal([]byte(jsonLD), &data); err != nil {
			logParser(fmt.Sprintf("JSON-LD script #%d parse error: %v", i, err))
			return // Skip invalid JSON
		}

		// Handle different JSON-LD formats
		recipe = extractRecipeFromJSONLD(data)
		if recipe != nil {
			logParser(fmt.Sprintf("Found Recipe in JSON-LD script #%d: %s", i, recipe.Name))
		}
	})

	return recipe, nil
}

// extractRecipeFromJSONLD extracts Recipe from various JSON-LD formats
func extractRecipeFromJSONLD(data interface{}) *Recipe {
	switch v := data.(type) {
	case map[string]interface{}:
		// Single object - check if it's a Recipe or has @graph
		if recipe := extractRecipeFromObject(v); recipe != nil {
			return recipe
		}
		// Check for @graph array
		if graph, ok := v["@graph"].([]interface{}); ok {
			return extractRecipeFromArray(graph)
		}
	case []interface{}:
		// Array of objects
		return extractRecipeFromArray(v)
	}
	return nil
}

// extractRecipeFromArray extracts Recipe from an array of objects
func extractRecipeFromArray(arr []interface{}) *Recipe {
	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			if recipe := extractRecipeFromObject(obj); recipe != nil {
				return recipe
			}
		}
	}
	return nil
}

// extractRecipeFromObject extracts Recipe from a single object
func extractRecipeFromObject(obj map[string]interface{}) *Recipe {
	// Check if it's a Recipe type
	typeVal, ok := obj["@type"].(string)
	if !ok {
		return nil
	}

	// Handle both "Recipe" and "https://schema.org/Recipe"
	if typeVal != "Recipe" && !strings.Contains(typeVal, "Recipe") {
		return nil
	}

	// Parse Recipe
	recipe := &Recipe{}

	// Required: name
	if name, ok := obj["name"].(string); ok && name != "" {
		recipe.Name = name
	} else {
		// Name is required, skip if not present
		return nil
	}

	// Required: image (can be string or array)
	if imageVal, ok := obj["image"]; ok {
		recipe.Image = parseImage(imageVal)
		if len(recipe.Image) == 0 {
			// Image is required, skip if empty
			return nil
		}
	} else {
		// Image is required, skip if not present
		return nil
	}

	// Optional fields
	recipe.Context = getString(obj, "@context")
	recipe.Type = getString(obj, "@type")

	if desc := getStringPtr(obj, "description"); desc != nil {
		recipe.Description = desc
	}
	if prepTime := getStringPtr(obj, "prepTime"); prepTime != nil {
		recipe.PrepTime = prepTime
	}
	if cookTime := getStringPtr(obj, "cookTime"); cookTime != nil {
		recipe.CookTime = cookTime
	}
	if totalTime := getStringPtr(obj, "totalTime"); totalTime != nil {
		recipe.TotalTime = totalTime
	}
	if yield := getStringPtr(obj, "recipeYield"); yield != nil {
		recipe.RecipeYield = yield
	}
	if category := getStringPtr(obj, "recipeCategory"); category != nil {
		recipe.RecipeCategory = category
	}
	if cuisine := getStringPtr(obj, "recipeCuisine"); cuisine != nil {
		recipe.RecipeCuisine = cuisine
	}
	if keywords := getStringPtr(obj, "keywords"); keywords != nil {
		recipe.Keywords = keywords
	}
	if datePublished := getStringPtr(obj, "datePublished"); datePublished != nil {
		recipe.DatePublished = datePublished
	}
	if dateModified := getStringPtr(obj, "dateModified"); dateModified != nil {
		recipe.DateModified = dateModified
	}

	// Recipe ingredients
	if ingredients := getStringArray(obj, "recipeIngredient"); len(ingredients) > 0 {
		recipe.RecipeIngredient = ingredients
	}

	// Recipe instructions
	if instructions := parseInstructions(obj["recipeInstructions"]); len(instructions) > 0 {
		recipe.RecipeInstructions = instructions
	}

	// Author
	if author := parseAuthor(obj["author"]); author != nil {
		recipe.Author = author
	}

	// Nutrition
	if nutrition := parseNutrition(obj["nutrition"]); nutrition != nil {
		recipe.Nutrition = nutrition
	}

	return recipe
}
