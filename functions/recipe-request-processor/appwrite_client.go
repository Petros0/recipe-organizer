package handler

import (
	"fmt"
	"os"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/role"
	"github.com/appwrite/sdk-for-go/tablesdb"
)

// Status constants for recipe request tracking
const (
	StatusRequested  = "REQUESTED"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

// Default Appwrite configuration
const (
	DefaultEndpoint    = "https://fra.cloud.appwrite.io/v1"
	DefaultProjectID   = "691f8b990030db50617a"
	DatabaseID         = "6930a343001607ad7cbd"
	CollectionID       = "6930a34300165ad1d129"
	RecipeCollectionID = "recipe"
)

// RecipeRequestStore defines the interface for recipe request operations
type RecipeRequestStore interface {
	UpdateStatus(documentID, status string) error
	CreateRecipe(requestID, userID string, recipe *Recipe) (string, error)
}

// RecipeRequestClient handles database operations for recipe requests
type RecipeRequestClient struct {
	tablesdb *tablesdb.TablesDB
}

// NewRecipeRequestClient creates a new RecipeRequestClient with Appwrite configuration
func NewRecipeRequestClient() *RecipeRequestClient {
	endpoint := os.Getenv("APPWRITE_ENDPOINT")
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	projectID := os.Getenv("APPWRITE_PROJECT_ID")
	if projectID == "" {
		projectID = DefaultProjectID
	}

	apiKey := os.Getenv("APPWRITE_API_KEY")

	client := appwrite.NewClient(
		appwrite.WithEndpoint(endpoint),
		appwrite.WithProject(projectID),
		appwrite.WithKey(apiKey),
	)

	return &RecipeRequestClient{
		tablesdb: appwrite.NewTablesDB(client),
	}
}

// UpdateStatus updates the status of an existing recipe request
func (c *RecipeRequestClient) UpdateStatus(documentID, status string) error {
	data := map[string]interface{}{
		"status": status,
	}

	_, err := c.tablesdb.UpdateRow(
		DatabaseID,
		CollectionID,
		documentID,
		c.tablesdb.WithUpdateRowData(data),
	)
	if err != nil {
		return fmt.Errorf("failed to update recipe request status to %s: %w", status, err)
	}
	return nil
}

// CreateRecipe creates a new recipe document linked to a recipe request
func (c *RecipeRequestClient) CreateRecipe(requestID, userID string, recipe *Recipe) (string, error) {
	data := recipeToMap(requestID, userID, recipe)

	// Set permissions to allow only the user who owns the recipe to access it
	permissions := []string{
		permission.Read(role.User(userID, "")),
		permission.Update(role.User(userID, "")),
		permission.Delete(role.User(userID, "")),
	}

	doc, err := c.tablesdb.CreateRow(
		DatabaseID,
		RecipeCollectionID,
		id.Unique(),
		data,
		c.tablesdb.WithCreateRowPermissions(permissions),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create recipe document: %w", err)
	}

	return doc.Id, nil
}

// recipeToMap converts a Recipe struct to a map for Appwrite document creation
func recipeToMap(requestID, userID string, recipe *Recipe) map[string]interface{} {
	data := map[string]interface{}{
		"fk_recipe_request": requestID,
		"user_id":           userID,
		"name":              recipe.Name,
	}

	// Optional string fields
	if recipe.Description != nil {
		data["description"] = *recipe.Description
	}
	if recipe.PrepTime != nil {
		data["prep_time"] = *recipe.PrepTime
	}
	if recipe.CookTime != nil {
		data["cook_time"] = *recipe.CookTime
	}
	if recipe.TotalTime != nil {
		data["total_time"] = *recipe.TotalTime
	}
	if len(recipe.RecipeYield) > 0 {
		data["recipe_yield"] = recipe.RecipeYield
	}
	if len(recipe.RecipeCategory) > 0 {
		data["recipe_category"] = recipe.RecipeCategory
	}
	if len(recipe.RecipeCuisine) > 0 {
		data["recipe_cuisine"] = recipe.RecipeCuisine
	}
	if recipe.Keywords != nil {
		data["keywords"] = *recipe.Keywords
	}
	if recipe.DatePublished != nil {
		data["date_published"] = *recipe.DatePublished
	}
	if recipe.DateModified != nil {
		data["date_modified"] = *recipe.DateModified
	}

	// Image array
	if len(recipe.Image) > 0 {
		data["image"] = recipe.Image
	}

	// Ingredients array
	if len(recipe.RecipeIngredient) > 0 {
		data["ingredients"] = recipe.RecipeIngredient
	}

	// Instructions - flatten to string array
	if len(recipe.RecipeInstructions) > 0 {
		data["instructions"] = flattenInstructions(recipe.RecipeInstructions)
	}

	// Author fields (flattened)
	if recipe.Author != nil {
		if recipe.Author.Name != "" {
			data["author_name"] = recipe.Author.Name
		}
		if recipe.Author.URL != "" {
			data["author_url"] = recipe.Author.URL
		}
	}

	// Nutrition fields (flattened)
	if recipe.Nutrition != nil {
		if recipe.Nutrition.Calories != nil {
			data["nutrition_calories"] = *recipe.Nutrition.Calories
		}
		if recipe.Nutrition.FatContent != nil {
			data["nutrition_fat"] = *recipe.Nutrition.FatContent
		}
		if recipe.Nutrition.SaturatedFatContent != nil {
			data["nutrition_saturated_fat"] = *recipe.Nutrition.SaturatedFatContent
		}
		if recipe.Nutrition.CholesterolContent != nil {
			data["nutrition_cholesterol"] = *recipe.Nutrition.CholesterolContent
		}
		if recipe.Nutrition.SodiumContent != nil {
			data["nutrition_sodium"] = *recipe.Nutrition.SodiumContent
		}
		if recipe.Nutrition.CarbohydrateContent != nil {
			data["nutrition_carbohydrate"] = *recipe.Nutrition.CarbohydrateContent
		}
		if recipe.Nutrition.FiberContent != nil {
			data["nutrition_fiber"] = *recipe.Nutrition.FiberContent
		}
		if recipe.Nutrition.SugarContent != nil {
			data["nutrition_sugar"] = *recipe.Nutrition.SugarContent
		}
		if recipe.Nutrition.ProteinContent != nil {
			data["nutrition_protein"] = *recipe.Nutrition.ProteinContent
		}
	}

	return data
}
