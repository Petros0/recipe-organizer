package handler

import (
	"testing"
)

func TestExtractRecipeFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantName string
		wantNil  bool
		wantErr  bool
	}{
		{
			name: "valid JSON-LD recipe",
			html: `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
{
  "@type": "Recipe",
  "name": "Chocolate Cake",
  "image": "https://example.com/cake.jpg"
}
</script>
</head>
<body></body>
</html>`,
			wantName: "Chocolate Cake",
		},
		{
			name: "JSON-LD in @graph format",
			html: `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
{
  "@context": "https://schema.org",
  "@graph": [
    {"@type": "WebPage", "name": "Page"},
    {"@type": "Recipe", "name": "Pasta", "image": "https://example.com/pasta.jpg"}
  ]
}
</script>
</head>
<body></body>
</html>`,
			wantName: "Pasta",
		},
		{
			name: "JSON-LD as array",
			html: `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
[
  {"@type": "Organization", "name": "Org"},
  {"@type": "Recipe", "name": "Salad", "image": "https://example.com/salad.jpg"}
]
</script>
</head>
<body></body>
</html>`,
			wantName: "Salad",
		},
		{
			name: "no JSON-LD script",
			html: `<!DOCTYPE html>
<html>
<head><title>No Recipe</title></head>
<body></body>
</html>`,
			wantNil: true,
		},
		{
			name: "JSON-LD without Recipe type",
			html: `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
{"@type": "Article", "name": "Blog Post"}
</script>
</head>
<body></body>
</html>`,
			wantNil: true,
		},
		{
			name: "invalid JSON in script",
			html: `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
{invalid json}
</script>
</head>
<body></body>
</html>`,
			wantNil: true,
		},
		{
			name: "multiple JSON-LD scripts, recipe in second",
			html: `<!DOCTYPE html>
<html>
<head>
<script type="application/ld+json">
{"@type": "WebSite", "name": "Site"}
</script>
<script type="application/ld+json">
{"@type": "Recipe", "name": "Soup", "image": "https://example.com/soup.jpg"}
</script>
</head>
<body></body>
</html>`,
			wantName: "Soup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := extractRecipeFromHTML(tt.html)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
				return
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.wantNil {
				if recipe != nil {
					t.Errorf("Expected nil recipe, got %+v", recipe)
				}
				return
			}

			if recipe == nil {
				t.Fatal("Expected recipe but got nil")
			}

			if recipe.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", recipe.Name, tt.wantName)
			}
		})
	}
}

func TestExtractRecipeFromJSONLD(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		wantName string
		wantNil  bool
	}{
		{
			name: "single Recipe object",
			data: map[string]interface{}{
				"@type": "Recipe",
				"name":  "Test Recipe",
				"image": "https://example.com/img.jpg",
			},
			wantName: "Test Recipe",
		},
		{
			name: "Recipe with full schema.org type URL",
			data: map[string]interface{}{
				"@type": "https://schema.org/Recipe",
				"name":  "Schema Recipe",
				"image": "https://example.com/img.jpg",
			},
			wantName: "Schema Recipe",
		},
		{
			name: "@graph with Recipe",
			data: map[string]interface{}{
				"@graph": []interface{}{
					map[string]interface{}{"@type": "WebPage"},
					map[string]interface{}{
						"@type": "Recipe",
						"name":  "Graph Recipe",
						"image": "https://example.com/img.jpg",
					},
				},
			},
			wantName: "Graph Recipe",
		},
		{
			name: "array with Recipe",
			data: []interface{}{
				map[string]interface{}{"@type": "Organization"},
				map[string]interface{}{
					"@type": "Recipe",
					"name":  "Array Recipe",
					"image": "https://example.com/img.jpg",
				},
			},
			wantName: "Array Recipe",
		},
		{
			name:    "nil data",
			data:    nil,
			wantNil: true,
		},
		{
			name:    "string data",
			data:    "not an object",
			wantNil: true,
		},
		{
			name: "object without @type",
			data: map[string]interface{}{
				"name":  "No Type",
				"image": "https://example.com/img.jpg",
			},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe := extractRecipeFromJSONLD(tt.data)

			if tt.wantNil {
				if recipe != nil {
					t.Errorf("Expected nil, got %+v", recipe)
				}
				return
			}

			if recipe == nil {
				t.Fatal("Expected recipe but got nil")
			}

			if recipe.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", recipe.Name, tt.wantName)
			}
		})
	}
}

func TestExtractRecipeFromObject(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name             string
		obj              map[string]interface{}
		wantNil          bool
		wantName         string
		wantImageCount   int
		wantDescription  *string
		wantIngredients  int
		wantInstructions int
		wantYield        []string
		wantCategory     []string
	}{
		{
			name: "full recipe",
			obj: map[string]interface{}{
				"@type":            "Recipe",
				"name":             "Full Recipe",
				"image":            []interface{}{"img1.jpg", "img2.jpg"},
				"description":      "A delicious recipe",
				"recipeIngredient": []interface{}{"flour", "sugar", "eggs"},
				"recipeInstructions": []interface{}{
					map[string]interface{}{"@type": "HowToStep", "text": "Step 1"},
					map[string]interface{}{"@type": "HowToStep", "text": "Step 2"},
				},
			},
			wantName:         "Full Recipe",
			wantImageCount:   2,
			wantDescription:  strPtr("A delicious recipe"),
			wantIngredients:  3,
			wantInstructions: 2,
		},
		{
			name: "minimal valid recipe",
			obj: map[string]interface{}{
				"@type": "Recipe",
				"name":  "Minimal",
				"image": "single-image.jpg",
			},
			wantName:       "Minimal",
			wantImageCount: 1,
		},
		{
			name: "missing name",
			obj: map[string]interface{}{
				"@type": "Recipe",
				"image": "img.jpg",
			},
			wantNil: true,
		},
		{
			name: "empty name",
			obj: map[string]interface{}{
				"@type": "Recipe",
				"name":  "",
				"image": "img.jpg",
			},
			wantNil: true,
		},
		{
			name: "missing image",
			obj: map[string]interface{}{
				"@type": "Recipe",
				"name":  "No Image",
			},
			wantNil: true,
		},
		{
			name: "wrong type",
			obj: map[string]interface{}{
				"@type": "Article",
				"name":  "Article",
				"image": "img.jpg",
			},
			wantNil: true,
		},
		{
			name: "type contains Recipe",
			obj: map[string]interface{}{
				"@type": "schema:Recipe",
				"name":  "Schema Recipe",
				"image": "img.jpg",
			},
			wantName:       "Schema Recipe",
			wantImageCount: 1,
		},
		{
			name: "recipeYield as string",
			obj: map[string]interface{}{
				"@type":       "Recipe",
				"name":        "Yield Test",
				"image":       "img.jpg",
				"recipeYield": "6",
			},
			wantName:       "Yield Test",
			wantImageCount: 1,
			wantYield:      []string{"6"},
		},
		{
			name: "recipeYield as number",
			obj: map[string]interface{}{
				"@type":       "Recipe",
				"name":        "Yield Number Test",
				"image":       "img.jpg",
				"recipeYield": float64(4),
			},
			wantName:       "Yield Number Test",
			wantImageCount: 1,
			wantYield:      []string{"4"},
		},
		{
			name: "recipeCategory as number",
			obj: map[string]interface{}{
				"@type":          "Recipe",
				"name":           "Category Number Test",
				"image":          "img.jpg",
				"recipeCategory": float64(19),
			},
			wantName:       "Category Number Test",
			wantImageCount: 1,
			wantCategory:   []string{"19"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe := extractRecipeFromObject(tt.obj)

			if tt.wantNil {
				if recipe != nil {
					t.Errorf("Expected nil, got %+v", recipe)
				}
				return
			}

			if recipe == nil {
				t.Fatal("Expected recipe but got nil")
			}

			if recipe.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", recipe.Name, tt.wantName)
			}

			if len(recipe.Image) != tt.wantImageCount {
				t.Errorf("Image count = %d, want %d", len(recipe.Image), tt.wantImageCount)
			}

			if tt.wantDescription != nil {
				if recipe.Description == nil {
					t.Error("Expected description but got nil")
				} else if *recipe.Description != *tt.wantDescription {
					t.Errorf("Description = %q, want %q", *recipe.Description, *tt.wantDescription)
				}
			}

			if len(recipe.RecipeIngredient) != tt.wantIngredients {
				t.Errorf("Ingredients count = %d, want %d", len(recipe.RecipeIngredient), tt.wantIngredients)
			}

			if len(recipe.RecipeInstructions) != tt.wantInstructions {
				t.Errorf("Instructions count = %d, want %d", len(recipe.RecipeInstructions), tt.wantInstructions)
			}

			if len(tt.wantYield) > 0 {
				if len(recipe.RecipeYield) != len(tt.wantYield) {
					t.Errorf("RecipeYield = %v, want %v", recipe.RecipeYield, tt.wantYield)
				} else {
					for i, v := range tt.wantYield {
						if recipe.RecipeYield[i] != v {
							t.Errorf("RecipeYield[%d] = %q, want %q", i, recipe.RecipeYield[i], v)
						}
					}
				}
			}

			if len(tt.wantCategory) > 0 {
				if len(recipe.RecipeCategory) != len(tt.wantCategory) {
					t.Errorf("RecipeCategory = %v, want %v", recipe.RecipeCategory, tt.wantCategory)
				} else {
					for i, v := range tt.wantCategory {
						if recipe.RecipeCategory[i] != v {
							t.Errorf("RecipeCategory[%d] = %q, want %q", i, recipe.RecipeCategory[i], v)
						}
					}
				}
			}
		})
	}
}
