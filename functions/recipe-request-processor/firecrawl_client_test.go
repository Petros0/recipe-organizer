package handler

import (
	"os"
	"strings"
	"testing"
)

// TestFirecrawlStrategy_RawHTMLWithJSONLD tests that rawHtml format preserves JSON-LD script tags
// and allows proper extraction of recipe data in the original language
func TestFirecrawlStrategy_RawHTMLWithJSONLD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey == "" {
		t.Skip("FIRECRAWL_API_KEY not set, skipping Firecrawl integration test")
	}

	tests := []struct {
		name              string
		url               string
		wantGreekName     bool     // Expect Greek characters in recipe name
		wantCorrectImage  bool     // Expect actual recipe image, not icons
		wantNutrition     bool     // Expect nutrition data to be present
		invalidImageParts []string // Substrings that should NOT appear in image URL
	}{
		{
			name:              "akispetretzikis.com recipe - Greek millefeuille",
			url:               "https://akispetretzikis.com/recipe/9418/milfeig-se-pothri",
			wantGreekName:     true,
			wantCorrectImage:  true,
			wantNutrition:     true,
			invalidImageParts: []string{"star", "icon", "logo", "avatar", ".svg"},
		},
		{
			name:              "akispetretzikis.com recipe - mushroom risotto",
			url:               "https://akispetretzikis.com/recipe/175/rizoto-me-manitaria",
			wantGreekName:     true,
			wantCorrectImage:  true,
			wantNutrition:     true,
			invalidImageParts: []string{"star", "icon", "logo", "avatar", ".svg"},
		},
	}

	strategy := NewFirecrawlStrategy()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing URL: %s", tt.url)

			recipe, err := strategy.Fetch(tt.url)
			if err != nil {
				t.Fatalf("Failed to fetch recipe: %v", err)
			}

			if recipe == nil {
				t.Fatal("Expected recipe but got nil")
			}

			// Check recipe name
			if recipe.Name == "" {
				t.Error("Recipe name is empty")
			} else {
				t.Logf("Recipe Name: %s", recipe.Name)
			}

			// Check for Greek characters (indicates JSON-LD parsing, not LLM translation)
			if tt.wantGreekName {
				hasGreek := containsGreekCharacters(recipe.Name)
				if !hasGreek {
					t.Errorf("Expected Greek characters in recipe name, got: %s (likely LLM extraction was used instead of JSON-LD)", recipe.Name)
				} else {
					t.Log("✓ Recipe name contains Greek characters (JSON-LD extraction confirmed)")
				}
			}

			// Check image URL
			if len(recipe.Image) == 0 {
				t.Error("Recipe has no images")
			} else {
				imageURL := recipe.Image[0]
				t.Logf("Recipe Image: %s", imageURL)

				if tt.wantCorrectImage {
					for _, invalidPart := range tt.invalidImageParts {
						if strings.Contains(strings.ToLower(imageURL), invalidPart) {
							t.Errorf("Image URL contains invalid part '%s': %s (likely LLM picked wrong image)", invalidPart, imageURL)
						}
					}
				}
			}

			// Check nutrition data
			if tt.wantNutrition {
				if recipe.Nutrition == nil {
					t.Error("Expected nutrition data but got nil")
				} else {
					t.Log("✓ Nutrition data present")
					if recipe.Nutrition.Calories != nil {
						t.Logf("  Calories: %s", *recipe.Nutrition.Calories)
					}
					if recipe.Nutrition.FatContent != nil {
						t.Logf("  Fat: %s", *recipe.Nutrition.FatContent)
					}
					if recipe.Nutrition.CarbohydrateContent != nil {
						t.Logf("  Carbs: %s", *recipe.Nutrition.CarbohydrateContent)
					}
					if recipe.Nutrition.ProteinContent != nil {
						t.Logf("  Protein: %s", *recipe.Nutrition.ProteinContent)
					}
				}
			}

			// Check ingredients
			if len(recipe.RecipeIngredient) == 0 {
				t.Error("Recipe has no ingredients")
			} else {
				t.Logf("Ingredients count: %d", len(recipe.RecipeIngredient))
				// Log first 3 ingredients
				for i, ing := range recipe.RecipeIngredient {
					if i < 3 {
						t.Logf("  - %s", ing)
					}
				}
			}

			// Check instructions
			if len(recipe.RecipeInstructions) == 0 {
				t.Error("Recipe has no instructions")
			} else {
				t.Logf("Instructions count: %d", len(recipe.RecipeInstructions))
			}

			// Check author
			if recipe.Author != nil {
				t.Logf("Author: %s", recipe.Author.Name)
				// Check if author name is in Greek (another indicator of JSON-LD vs LLM)
				if tt.wantGreekName && containsGreekCharacters(recipe.Author.Name) {
					t.Log("✓ Author name contains Greek characters")
				}
			}
		})
	}
}

// TestFirecrawlStrategy_LLMExtractionFallback tests that LLM extraction works for sites without JSON-LD
func TestFirecrawlStrategy_LLMExtractionFallback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("FIRECRAWL_API_KEY")
	if apiKey == "" {
		t.Skip("FIRECRAWL_API_KEY not set, skipping Firecrawl integration test")
	}

	// Test with a site that does NOT have JSON-LD structured data
	url := "https://alfiecooks.substack.com/p/caramelised-onion-sun-dried-tomato"
	t.Logf("Testing LLM extraction fallback with URL: %s", url)

	strategy := NewFirecrawlStrategy()

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

	// Check for ingredients
	if len(recipe.RecipeIngredient) == 0 {
		t.Error("Expected recipe ingredients from LLM extraction")
	} else {
		t.Logf("Ingredients count: %d", len(recipe.RecipeIngredient))
	}

	// Check for instructions
	if len(recipe.RecipeInstructions) == 0 {
		t.Error("Expected recipe instructions from LLM extraction")
	} else {
		t.Logf("Instructions count: %d", len(recipe.RecipeInstructions))
	}
}

// containsGreekCharacters checks if a string contains Greek Unicode characters
func containsGreekCharacters(s string) bool {
	for _, r := range s {
		// Greek and Coptic block: U+0370 to U+03FF
		// Greek Extended block: U+1F00 to U+1FFF
		if (r >= 0x0370 && r <= 0x03FF) || (r >= 0x1F00 && r <= 0x1FFF) {
			return true
		}
	}
	return false
}
