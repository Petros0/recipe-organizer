package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	firecrawl "github.com/mendableai/firecrawl-go/v2"
)

// ErrNoJSONLD indicates that no JSON-LD structured data was found on the page
var ErrNoJSONLD = errors.New("no JSON-LD structured data found")

// FirecrawlStrategy implements FetchStrategy using Firecrawl API
// It uses a hybrid approach:
// 1. First tries to get HTML and parse JSON-LD (cheaper)
// 2. Falls back to LLM extraction if no JSON-LD found (more expensive but works for any page)
type FirecrawlStrategy struct {
	apiKey string
}

// NewFirecrawlStrategy creates a new FirecrawlStrategy with the API key from environment
func NewFirecrawlStrategy() *FirecrawlStrategy {
	return &FirecrawlStrategy{
		apiKey: os.Getenv("FIRECRAWL_API_KEY"),
	}
}

// Name returns the strategy name for logging
func (s *FirecrawlStrategy) Name() string {
	return "Firecrawl"
}

// CanRetry returns false - Firecrawl is typically the last resort
func (s *FirecrawlStrategy) CanRetry(err error) bool {
	return false
}

// Fetch uses Firecrawl API to fetch the page and extract recipe data
// It first tries HTML parsing for JSON-LD, then falls back to LLM extraction
func (s *FirecrawlStrategy) Fetch(url string) (*Recipe, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("FIRECRAWL_API_KEY environment variable is not set")
	}

	// Initialize Firecrawl client
	app, err := firecrawl.NewFirecrawlApp(s.apiKey, "")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firecrawl: %w", err)
	}

	// Step 1: Try to get HTML and parse JSON-LD (cheaper approach)
	recipe, err := s.fetchWithHTML(app, url)
	if err == nil && recipe != nil {
		return recipe, nil
	}

	// If we got a non-ErrNoJSONLD error, return it
	if err != nil && !errors.Is(err, ErrNoJSONLD) {
		return nil, err
	}

	// Step 2: Fall back to LLM extraction (for sites without JSON-LD)
	return s.fetchWithLLMExtraction(app, url)
}

// fetchWithHTML fetches the page HTML and parses JSON-LD
func (s *FirecrawlStrategy) fetchWithHTML(app *firecrawl.FirecrawlApp, url string) (*Recipe, error) {
	// Request HTML format
	params := &firecrawl.ScrapeParams{
		Formats: []string{"html"},
	}

	result, err := app.ScrapeURL(url, params)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape URL with Firecrawl: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("no result from Firecrawl")
	}

	// Extract recipe from HTML
	recipe, err := extractRecipeFromHTML(result.HTML)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	if recipe == nil {
		return nil, ErrNoJSONLD
	}

	return recipe, nil
}

// fetchWithLLMExtraction uses Firecrawl's LLM extraction to get structured recipe data
func (s *FirecrawlStrategy) fetchWithLLMExtraction(app *firecrawl.FirecrawlApp, url string) (*Recipe, error) {
	// Define JSON schema for recipe extraction
	jsonSchema := buildRecipeExtractionSchema()

	// Create a prompt to help the LLM extract recipe data
	prompt := `
Extract the complete recipe data from this page.
Focus on the main recipe content and ignore advertisements, related recipes, and sidebar content.
Extract all available fields including image, instructions, times, servings, nutrition, and author information when present.
`

	// Use scrape with JSON extraction options (v2 SDK feature)
	params := &firecrawl.ScrapeParams{
		Formats: []string{"json"},
		JsonOptions: &firecrawl.JsonOptions{
			Schema: jsonSchema,
			Prompt: &prompt,
		},
	}

	result, err := app.ScrapeURL(url, params)
	if err != nil {
		return nil, fmt.Errorf("failed to extract recipe with LLM: %w", err)
	}

	if result == nil || result.JSON == nil {
		return nil, fmt.Errorf("no data extracted from page")
	}

	// Parse the extracted data into Recipe struct
	return parseExtractedRecipe(result.JSON)
}

// buildRecipeExtractionSchema returns the JSON schema for recipe extraction
func buildRecipeExtractionSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "The name/title of the recipe",
			},
			"description": map[string]any{
				"type":        "string",
				"description": "A brief description of the recipe",
			},
			"image": map[string]any{
				"type":        "string",
				"description": "URL of the recipe image (usually it is the thumbnail image)",
			},
			"prepTime": map[string]any{
				"type":        "string",
				"description": "Preparation time (e.g., '15 minutes' or 'PT15M')",
			},
			"cookTime": map[string]any{
				"type":        "string",
				"description": "Cooking time (e.g., '30 minutes' or 'PT30M')",
			},
			"totalTime": map[string]any{
				"type":        "string",
				"description": "Total time (e.g., '45 minutes' or 'PT45M')",
			},
			"recipeYield": map[string]any{
				"type":        "string",
				"description": "Number of servings or yield (e.g., '4 servings')",
			},
			"recipeIngredient": map[string]any{
				"type":        "array",
				"items":       map[string]string{"type": "string"},
				"description": "List of ingredients with quantities",
			},
			"recipeInstructions": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type":        "string",
					"description": "A single step in the cooking instructions",
				},
				"description": "Step-by-step cooking instructions",
			},
			"author": map[string]any{
				"type":        "string",
				"description": "Author or creator of the recipe",
			},
			"recipeCategory": map[string]any{
				"type":        "string",
				"description": "Category (e.g., 'Dessert', 'Main Course')",
			},
			"recipeCuisine": map[string]any{
				"type":        "string",
				"description": "Cuisine type (e.g., 'Italian', 'Mexican')",
			},
			"nutrition": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"calories":            map[string]string{"type": "string"},
					"fatContent":          map[string]string{"type": "string"},
					"carbohydrateContent": map[string]string{"type": "string"},
					"proteinContent":      map[string]string{"type": "string"},
				},
				"description": "Nutritional information",
			},
		},
		"required": []string{"name", "recipeIngredient", "recipeInstructions"},
	}
}

// ExtractedRecipe represents the LLM-extracted recipe data
type ExtractedRecipe struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Image              string             `json:"image"`
	PrepTime           string             `json:"prepTime"`
	CookTime           string             `json:"cookTime"`
	TotalTime          string             `json:"totalTime"`
	RecipeYield        string             `json:"recipeYield"`
	RecipeIngredient   []string           `json:"recipeIngredient"`
	RecipeInstructions []string           `json:"recipeInstructions"`
	Author             string             `json:"author"`
	RecipeCategory     string             `json:"recipeCategory"`
	RecipeCuisine      string             `json:"recipeCuisine"`
	Nutrition          ExtractedNutrition `json:"nutrition"`
}

// ExtractedNutrition represents LLM-extracted nutrition data
type ExtractedNutrition struct {
	Calories            string `json:"calories"`
	FatContent          string `json:"fatContent"`
	CarbohydrateContent string `json:"carbohydrateContent"`
	ProteinContent      string `json:"proteinContent"`
}

// parseExtractedRecipe converts Firecrawl's JSON extraction result to a Recipe struct
func parseExtractedRecipe(data map[string]any) (*Recipe, error) {
	if data == nil {
		return nil, fmt.Errorf("no data extracted from page")
	}

	// Marshal the data back to JSON and unmarshal into our struct
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal extracted data: %w", err)
	}

	var extracted ExtractedRecipe
	if err := json.Unmarshal(dataBytes, &extracted); err != nil {
		return nil, fmt.Errorf("failed to parse extracted recipe: %w", err)
	}

	// Validate required fields
	if extracted.Name == "" {
		return nil, fmt.Errorf("extracted recipe missing required field: name")
	}

	// Convert to Recipe struct
	recipe := &Recipe{
		Context: "https://schema.org",
		Type:    "Recipe",
		Name:    extracted.Name,
	}

	// Handle image (convert single string to array)
	if extracted.Image != "" {
		recipe.Image = []string{extracted.Image}
	}

	// Set optional string fields
	if extracted.Description != "" {
		recipe.Description = &extracted.Description
	}
	if extracted.PrepTime != "" {
		recipe.PrepTime = &extracted.PrepTime
	}
	if extracted.CookTime != "" {
		recipe.CookTime = &extracted.CookTime
	}
	if extracted.TotalTime != "" {
		recipe.TotalTime = &extracted.TotalTime
	}
	if extracted.RecipeYield != "" {
		recipe.RecipeYield = &extracted.RecipeYield
	}
	if extracted.RecipeCategory != "" {
		recipe.RecipeCategory = &extracted.RecipeCategory
	}
	if extracted.RecipeCuisine != "" {
		recipe.RecipeCuisine = &extracted.RecipeCuisine
	}

	// Set ingredients
	recipe.RecipeIngredient = extracted.RecipeIngredient

	// Convert instructions from strings to RecipeInstruction
	if len(extracted.RecipeInstructions) > 0 {
		recipe.RecipeInstructions = make([]RecipeInstruction, len(extracted.RecipeInstructions))
		for i, step := range extracted.RecipeInstructions {
			recipe.RecipeInstructions[i] = RecipeInstruction{
				Type: "HowToStep",
				Text: step,
			}
		}
	}

	// Set author
	if extracted.Author != "" {
		recipe.Author = &Person{
			Type: "Person",
			Name: extracted.Author,
		}
	}

	// Set nutrition
	if extracted.Nutrition.Calories != "" || extracted.Nutrition.FatContent != "" ||
		extracted.Nutrition.CarbohydrateContent != "" || extracted.Nutrition.ProteinContent != "" {
		recipe.Nutrition = &Nutrition{
			Type: "NutritionInformation",
		}
		if extracted.Nutrition.Calories != "" {
			recipe.Nutrition.Calories = &extracted.Nutrition.Calories
		}
		if extracted.Nutrition.FatContent != "" {
			recipe.Nutrition.FatContent = &extracted.Nutrition.FatContent
		}
		if extracted.Nutrition.CarbohydrateContent != "" {
			recipe.Nutrition.CarbohydrateContent = &extracted.Nutrition.CarbohydrateContent
		}
		if extracted.Nutrition.ProteinContent != "" {
			recipe.Nutrition.ProteinContent = &extracted.Nutrition.ProteinContent
		}
	}

	return recipe, nil
}
