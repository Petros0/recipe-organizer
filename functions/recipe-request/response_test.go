package handler

import (
	"testing"
)

func TestToRecipeResponse(t *testing.T) {
	// Helper to create string pointers
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name                 string
		inputURL             string
		inputRecipe          *Recipe
		expectedURL          string
		expectedName         string
		expectedDesc         string
		expectedImage        string
		expectedPrepTime     string
		expectedCookTime     string
		expectedTotalTime    string
		expectedAuthor       string
		expectedIngredients  []string
		expectedInstructions []string
	}{
		{
			name:     "full recipe with all fields",
			inputURL: "https://example.com/recipe/1",
			inputRecipe: &Recipe{
				Name:             "Test Recipe",
				Description:      strPtr("A delicious test recipe"),
				Image:            []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
				PrepTime:         strPtr("PT15M"),
				CookTime:         strPtr("PT30M"),
				TotalTime:        strPtr("PT45M"),
				Author:           &Person{Name: "Test Chef"},
				RecipeIngredient: []string{"1 cup flour", "2 eggs", "1 tsp salt"},
				RecipeInstructions: []RecipeInstruction{
					{Type: "HowToStep", Text: "Mix flour and salt"},
					{Type: "HowToStep", Text: "Add eggs and stir"},
				},
			},
			expectedURL:          "https://example.com/recipe/1",
			expectedName:         "Test Recipe",
			expectedDesc:         "A delicious test recipe",
			expectedImage:        "https://example.com/image1.jpg",
			expectedPrepTime:     "PT15M",
			expectedCookTime:     "PT30M",
			expectedTotalTime:    "PT45M",
			expectedAuthor:       "Test Chef",
			expectedIngredients:  []string{"1 cup flour", "2 eggs", "1 tsp salt"},
			expectedInstructions: []string{"Mix flour and salt", "Add eggs and stir"},
		},
		{
			name:     "minimal recipe with only required fields",
			inputURL: "https://example.com/recipe/2",
			inputRecipe: &Recipe{
				Name:  "Minimal Recipe",
				Image: []string{"https://example.com/image.jpg"},
			},
			expectedURL:          "https://example.com/recipe/2",
			expectedName:         "Minimal Recipe",
			expectedDesc:         "",
			expectedImage:        "https://example.com/image.jpg",
			expectedPrepTime:     "",
			expectedCookTime:     "",
			expectedTotalTime:    "",
			expectedAuthor:       "",
			expectedIngredients:  []string{},
			expectedInstructions: []string{},
		},
		{
			name:     "recipe with nested HowToSection instructions",
			inputURL: "https://example.com/recipe/3",
			inputRecipe: &Recipe{
				Name:  "Recipe with Sections",
				Image: []string{"https://example.com/image.jpg"},
				RecipeInstructions: []RecipeInstruction{
					{
						Type: "HowToSection",
						Name: "Preparation",
						ItemListElement: []RecipeInstruction{
							{Type: "HowToStep", Text: "Prepare ingredients"},
							{Type: "HowToStep", Text: "Preheat oven"},
						},
					},
					{
						Type: "HowToSection",
						Name: "Cooking",
						ItemListElement: []RecipeInstruction{
							{Type: "HowToStep", Text: "Cook for 30 minutes"},
						},
					},
				},
			},
			expectedURL:          "https://example.com/recipe/3",
			expectedName:         "Recipe with Sections",
			expectedDesc:         "",
			expectedImage:        "https://example.com/image.jpg",
			expectedPrepTime:     "",
			expectedCookTime:     "",
			expectedTotalTime:    "",
			expectedAuthor:       "",
			expectedIngredients:  []string{},
			expectedInstructions: []string{"Prepare ingredients", "Preheat oven", "Cook for 30 minutes"},
		},
		{
			name:     "recipe with no images returns empty image",
			inputURL: "https://example.com/recipe/4",
			inputRecipe: &Recipe{
				Name:  "No Image Recipe",
				Image: []string{},
			},
			expectedURL:          "https://example.com/recipe/4",
			expectedName:         "No Image Recipe",
			expectedDesc:         "",
			expectedImage:        "",
			expectedPrepTime:     "",
			expectedCookTime:     "",
			expectedTotalTime:    "",
			expectedAuthor:       "",
			expectedIngredients:  []string{},
			expectedInstructions: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toRecipeResponse(tt.inputURL, tt.inputRecipe)

			// Check URL
			if result.URL != tt.expectedURL {
				t.Errorf("URL = %q, want %q", result.URL, tt.expectedURL)
			}

			// Check Recipe details
			if result.Recipe.Name != tt.expectedName {
				t.Errorf("Recipe.Name = %q, want %q", result.Recipe.Name, tt.expectedName)
			}
			if result.Recipe.Description != tt.expectedDesc {
				t.Errorf("Recipe.Description = %q, want %q", result.Recipe.Description, tt.expectedDesc)
			}
			if result.Recipe.Image != tt.expectedImage {
				t.Errorf("Recipe.Image = %q, want %q", result.Recipe.Image, tt.expectedImage)
			}
			if result.Recipe.PrepTime != tt.expectedPrepTime {
				t.Errorf("Recipe.PrepTime = %q, want %q", result.Recipe.PrepTime, tt.expectedPrepTime)
			}
			if result.Recipe.CookTime != tt.expectedCookTime {
				t.Errorf("Recipe.CookTime = %q, want %q", result.Recipe.CookTime, tt.expectedCookTime)
			}
			if result.Recipe.TotalTime != tt.expectedTotalTime {
				t.Errorf("Recipe.TotalTime = %q, want %q", result.Recipe.TotalTime, tt.expectedTotalTime)
			}
			if result.Recipe.Author != tt.expectedAuthor {
				t.Errorf("Recipe.Author = %q, want %q", result.Recipe.Author, tt.expectedAuthor)
			}

			// Check Ingredients
			if len(result.Ingredients) != len(tt.expectedIngredients) {
				t.Errorf("Ingredients count = %d, want %d", len(result.Ingredients), len(tt.expectedIngredients))
			} else {
				for i, ing := range result.Ingredients {
					if ing != tt.expectedIngredients[i] {
						t.Errorf("Ingredients[%d] = %q, want %q", i, ing, tt.expectedIngredients[i])
					}
				}
			}

			// Check Instructions
			if len(result.Instructions) != len(tt.expectedInstructions) {
				t.Errorf("Instructions count = %d, want %d", len(result.Instructions), len(tt.expectedInstructions))
			} else {
				for i, inst := range result.Instructions {
					if inst != tt.expectedInstructions[i] {
						t.Errorf("Instructions[%d] = %q, want %q", i, inst, tt.expectedInstructions[i])
					}
				}
			}
		})
	}
}

