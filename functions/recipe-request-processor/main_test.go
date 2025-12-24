package handler

import (
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
		NewHTTPClientStrategy(nil),
		NewFirecrawlStrategy(nil),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := executor.Execute(tt.urlStr, nil)

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
