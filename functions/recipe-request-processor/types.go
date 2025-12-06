package handler

// Recipe represents a schema.org Recipe structured data
type Recipe struct {
	Context            string              `json:"@context,omitempty"`
	Type               string              `json:"@type,omitempty"`
	Name               string              `json:"name"`
	Image              []string            `json:"image,omitempty"`
	Author             *Person             `json:"author,omitempty"`
	Description        *string             `json:"description,omitempty"`
	PrepTime           *string             `json:"prepTime,omitempty"`
	CookTime           *string             `json:"cookTime,omitempty"`
	TotalTime          *string             `json:"totalTime,omitempty"`
	RecipeYield        *string             `json:"recipeYield,omitempty"`
	RecipeIngredient   []string            `json:"recipeIngredient,omitempty"`
	RecipeInstructions []RecipeInstruction `json:"recipeInstructions,omitempty"`
	RecipeCategory     *string             `json:"recipeCategory,omitempty"`
	RecipeCuisine      *string             `json:"recipeCuisine,omitempty"`
	Nutrition          *Nutrition          `json:"nutrition,omitempty"`
	Keywords           *string             `json:"keywords,omitempty"`
	DatePublished      *string             `json:"datePublished,omitempty"`
	DateModified       *string             `json:"dateModified,omitempty"`
}

// Person represents a schema.org Person
type Person struct {
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Nutrition represents schema.org NutritionInformation
type Nutrition struct {
	Type                string  `json:"@type,omitempty"`
	Calories            *string `json:"calories,omitempty"`
	FatContent          *string `json:"fatContent,omitempty"`
	SaturatedFatContent *string `json:"saturatedFatContent,omitempty"`
	CholesterolContent  *string `json:"cholesterolContent,omitempty"`
	SodiumContent       *string `json:"sodiumContent,omitempty"`
	CarbohydrateContent *string `json:"carbohydrateContent,omitempty"`
	FiberContent        *string `json:"fiberContent,omitempty"`
	SugarContent        *string `json:"sugarContent,omitempty"`
	ProteinContent      *string `json:"proteinContent,omitempty"`
}

// RecipeInstruction represents a schema.org HowToStep or HowToSection
type RecipeInstruction struct {
	Type            string              `json:"@type,omitempty"`
	Text            string              `json:"text,omitempty"`
	Name            string              `json:"name,omitempty"`
	URL             string              `json:"url,omitempty"`
	ItemListElement []RecipeInstruction `json:"itemListElement,omitempty"` // For HowToSection
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RecipeResponse is the custom API response format
type RecipeResponse struct {
	URL          string        `json:"url"`
	Recipe       RecipeDetails `json:"recipe"`
	Instructions []string      `json:"instructions"`
	Ingredients  []string      `json:"ingredients"`
}

// RecipeDetails contains the flattened recipe metadata
type RecipeDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	PrepTime    string `json:"prepTime"`
	CookTime    string `json:"cookTime"`
	TotalTime   string `json:"totalTime"`
	Author      string `json:"author"`
}

// RequestBody represents the JSON request body
type RequestBody struct {
	URL string `json:"url"`
}
