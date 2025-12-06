package handler

import (
	"os"
	"testing"
)

func TestFetchRecipe_Integration(t *testing.T) {
	// Skip integration tests if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	tests := []struct {
		name          string
		urlStr        string
		wantName      bool // Whether we expect a recipe name to be present
		wantImage     bool // Whether we expect an image to be present
		requireRecipe bool // Whether recipe JSON-LD is required (some sites don't have it)
	}{
		{
			name:          "akispetretzikis.com recipe",
			urlStr:        "https://akispetretzikis.com/recipe/175/rizoto-me-manitaria",
			wantName:      true,
			wantImage:     true,
			requireRecipe: true,
		},
		{
			name:          "ohmyveggies.com recipe",
			urlStr:        "https://ohmyveggies.com/french-bread-pizza-with-pesto-and-sun-dried-tomatoes",
			wantName:      true,
			wantImage:     true,
			requireRecipe: true,
		},
	}

	// Create strategy executor with HTTP client first, then Firecrawl as fallback
	executor := NewStrategyExecutor(
		&HTTPClientStrategy{},
		NewFirecrawlStrategy(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := executor.Execute(tt.urlStr, func(msgs ...interface{}) {
				t.Log(msgs...)
			})

			if err != nil {
				t.Logf("Could not fetch recipe from %s: %v", tt.urlStr, err)
				return
			}

			if recipe == nil {
				if tt.requireRecipe {
					t.Fatal("Expected recipe but got nil")
				}
				t.Log("No Recipe JSON-LD found (expected)")
				return
			}

			// Validate required fields
			if tt.wantName && recipe.Name == "" {
				t.Error("Expected recipe name")
			} else {
				t.Logf("Recipe: %s", recipe.Name)
			}

			if tt.wantImage && len(recipe.Image) == 0 {
				t.Error("Expected recipe image")
			}

			// Log additional info
			if len(recipe.RecipeIngredient) > 0 {
				t.Logf("Ingredients: %d", len(recipe.RecipeIngredient))
			}
			if len(recipe.RecipeInstructions) > 0 {
				t.Logf("Instructions: %d", len(recipe.RecipeInstructions))
			}
		})
	}
}

// TestFirecrawlStrategy_HTMLWithJSONLD tests Firecrawl's ability to fetch HTML and parse JSON-LD
func TestFirecrawlStrategy_HTMLWithJSONLD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey == "" {
		t.Skip("FIRECRAWL_API_KEY not set, skipping Firecrawl integration test")
	}

	strategy := NewFirecrawlStrategy()

	// Test with a site that has JSON-LD structured data
	url := "https://akispetretzikis.com/recipe/175/rizoto-me-manitaria"
	t.Logf("Testing Firecrawl HTML extraction with JSON-LD: %s", url)

	recipe, err := strategy.Fetch(url)
	if err != nil {
		t.Fatalf("Failed to fetch recipe: %v", err)
	}

	if recipe == nil {
		t.Fatal("Expected recipe but got nil")
	}

	// Validate required fields
	if recipe.Name == "" {
		t.Error("Expected recipe name")
	} else {
		t.Logf("Recipe Name: %s", recipe.Name)
	}

	if len(recipe.Image) == 0 {
		t.Error("Expected recipe image")
	} else {
		t.Logf("Recipe Image: %s", recipe.Image[0])
	}

	// Check for ingredients
	if len(recipe.RecipeIngredient) == 0 {
		t.Error("Expected recipe ingredients")
	} else {
		t.Logf("Ingredients count: %d", len(recipe.RecipeIngredient))
		for i, ing := range recipe.RecipeIngredient {
			if i < 3 { // Log first 3 ingredients
				t.Logf("  - %s", ing)
			}
		}
	}

	// Check for instructions
	if len(recipe.RecipeInstructions) == 0 {
		t.Error("Expected recipe instructions")
	} else {
		t.Logf("Instructions count: %d", len(recipe.RecipeInstructions))
	}
}

// TestFirecrawlStrategy_LLMExtraction tests Firecrawl's LLM extraction for sites without JSON-LD
func TestFirecrawlStrategy_LLMExtraction(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey == "" {
		t.Skip("FIRECRAWL_API_KEY not set, skipping Firecrawl integration test")
	}

	strategy := NewFirecrawlStrategy()

	// Test with a site that does NOT have JSON-LD structured data
	// This Substack newsletter has a recipe but no schema.org markup
	url := "https://alfiecooks.substack.com/p/caramelised-onion-sun-dried-tomato"
	t.Logf("Testing Firecrawl LLM extraction (no JSON-LD): %s", url)

	recipe, err := strategy.Fetch(url)
	if err != nil {
		t.Fatalf("Failed to fetch recipe: %v", err)
	}

	if recipe == nil {
		t.Fatal("Expected recipe but got nil - LLM extraction should have found the recipe")
	}

	// Validate extracted fields
	if recipe.Name == "" {
		t.Error("Expected recipe name from LLM extraction")
	} else {
		t.Logf("Recipe Name: %s", recipe.Name)
	}

	// Check for ingredients - the Substack post has ingredients
	if len(recipe.RecipeIngredient) == 0 {
		t.Error("Expected recipe ingredients from LLM extraction")
	} else {
		t.Logf("Ingredients count: %d", len(recipe.RecipeIngredient))
		for i, ing := range recipe.RecipeIngredient {
			if i < 5 { // Log first 5 ingredients
				t.Logf("  - %s", ing)
			}
		}
	}

	// Check for instructions - the Substack post has method steps
	if len(recipe.RecipeInstructions) == 0 {
		t.Error("Expected recipe instructions from LLM extraction")
	} else {
		t.Logf("Instructions count: %d", len(recipe.RecipeInstructions))
		for i, inst := range recipe.RecipeInstructions {
			if i < 3 { // Log first 3 instructions
				t.Logf("  Step %d: %s...", i+1, truncate(inst.Text, 80))
			}
		}
	}

	// Log other extracted fields
	if recipe.Description != nil && *recipe.Description != "" {
		t.Logf("Description: %s...", truncate(*recipe.Description, 100))
	}
	if recipe.Author != nil {
		t.Logf("Author: %s", recipe.Author.Name)
	}
}

// truncate truncates a string to maxLen characters and adds "..." if truncated
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

