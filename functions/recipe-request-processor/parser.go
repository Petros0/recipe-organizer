package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// parserLogger is used for logging within parser functions
var parserLogger LogFunc

// SetParserLogger sets the logger for parser functions
func SetParserLogger(logger LogFunc) {
	parserLogger = logger
}

func logParser(msg string) {
	if parserLogger != nil {
		parserLogger("[Parser] " + msg)
	}
}

// extractRecipeFromHTML extracts Recipe JSON-LD from HTML content
func extractRecipeFromHTML(htmlContent string) (*Recipe, error) {
	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find all JSON-LD script tags
	var recipe *Recipe
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		if recipe != nil {
			return // Already found a recipe
		}

		jsonLD := s.Text()
		if jsonLD == "" {
			return
		}

		// Log raw JSON-LD for debugging
		logParser(fmt.Sprintf("Found JSON-LD script #%d (length: %d chars)", i+1, len(jsonLD)))
		// Log a snippet around recipeYield if present
		if idx := strings.Index(jsonLD, "recipeYield"); idx >= 0 {
			start := idx
			end := idx + 50
			if end > len(jsonLD) {
				end = len(jsonLD)
			}
			logParser(fmt.Sprintf("Raw recipeYield snippet: %s", jsonLD[start:end]))
		}

		// Try to parse as single object or array
		var data interface{}
		if err := json.Unmarshal([]byte(jsonLD), &data); err != nil {
			logParser(fmt.Sprintf("Failed to parse JSON-LD: %v", err))
			return // Skip invalid JSON
		}

		// Handle different JSON-LD formats
		recipe = extractRecipeFromJSONLD(data)
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
	// Log recipeYield before parsing
	logParser(fmt.Sprintf("recipeYield raw value: %v (type: %T)", obj["recipeYield"], obj["recipeYield"]))
	if yield := parseStringOrArray(obj["recipeYield"]); len(yield) > 0 {
		recipe.RecipeYield = yield
		logParser(fmt.Sprintf("recipeYield parsed: %v", yield))
	} else {
		logParser("recipeYield parsed to empty array")
	}
	if category := parseStringOrArray(obj["recipeCategory"]); len(category) > 0 {
		recipe.RecipeCategory = category
	}
	if cuisine := parseStringOrArray(obj["recipeCuisine"]); len(cuisine) > 0 {
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
