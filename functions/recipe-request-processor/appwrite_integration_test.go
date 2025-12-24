package handler

import (
	"os"
	"testing"
)

// Integration tests for Appwrite TablesDB operations
// Run with: APPWRITE_API_KEY=your_key go test -v -run Integration

func skipIfNoAPIKey(t *testing.T) {
	t.Helper()
	if os.Getenv("APPWRITE_API_KEY") == "" {
		t.Skip("APPWRITE_API_KEY not set, skipping integration test")
	}
}

func TestIntegration_RecipeRequest_CRUD(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewRecipeRequestClient()
	testUserID := "integration-test-user"

	// CREATE - Create a new recipe request
	t.Log("Creating recipe request...")
	requestID, err := createRecipeRequest(client, "https://example.com/integration-test-recipe", testUserID)
	if err != nil {
		t.Fatalf("Failed to create recipe request: %v", err)
	}
	t.Logf("Created recipe request with ID: %s", requestID)

	// Cleanup at the end
	defer func() {
		t.Log("Cleaning up - deleting recipe request...")
		if err := deleteRecipeRequest(client, requestID); err != nil {
			t.Logf("Warning: Failed to delete recipe request: %v", err)
		}
	}()

	// READ - Verify the request was created (via status update which requires the doc to exist)
	t.Log("Updating status to IN_PROGRESS...")
	err = client.UpdateStatus(requestID, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to update status to IN_PROGRESS: %v", err)
	}

	// UPDATE - Update status to COMPLETED
	t.Log("Updating status to COMPLETED...")
	err = client.UpdateStatus(requestID, StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to update status to COMPLETED: %v", err)
	}

	t.Log("Recipe request CRUD test completed successfully")
}

func TestIntegration_Recipe_CRUD(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewRecipeRequestClient()
	testUserID := "integration-test-user"

	// First create a recipe request (required for recipe FK)
	t.Log("Creating parent recipe request...")
	requestID, err := createRecipeRequest(client, "https://example.com/recipe-crud-test", testUserID)
	if err != nil {
		t.Fatalf("Failed to create recipe request: %v", err)
	}
	t.Logf("Created recipe request with ID: %s", requestID)

	// Cleanup recipe request at the end
	defer func() {
		t.Log("Cleaning up - deleting recipe request...")
		if err := deleteRecipeRequest(client, requestID); err != nil {
			t.Logf("Warning: Failed to delete recipe request: %v", err)
		}
	}()

	// CREATE - Create a recipe
	t.Log("Creating recipe...")
	recipe := createTestRecipe()
	recipeID, err := client.CreateRecipe(requestID, testUserID, recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}
	t.Logf("Created recipe with ID: %s", recipeID)

	// Cleanup recipe at the end
	defer func() {
		t.Log("Cleaning up - deleting recipe...")
		if err := deleteRecipe(client, recipeID); err != nil {
			t.Logf("Warning: Failed to delete recipe: %v", err)
		}
	}()

	// Verify recipe was created by checking the ID is not empty
	if recipeID == "" {
		t.Fatal("Expected non-empty recipe ID")
	}

	t.Log("Recipe CRUD test completed successfully")
}

func TestIntegration_FullWorkflow(t *testing.T) {
	skipIfNoAPIKey(t)

	client := NewRecipeRequestClient()
	testUserID := "integration-test-user"

	// Step 1: Create recipe request (simulates user submitting a URL)
	t.Log("Step 1: Creating recipe request...")
	requestID, err := createRecipeRequest(client, "https://example.com/full-workflow-test", testUserID)
	if err != nil {
		t.Fatalf("Failed to create recipe request: %v", err)
	}
	t.Logf("Created recipe request: %s", requestID)

	// Cleanup
	defer func() {
		deleteRecipeRequest(client, requestID)
	}()

	// Step 2: Update status to IN_PROGRESS (simulates processor picking up the request)
	t.Log("Step 2: Updating status to IN_PROGRESS...")
	err = client.UpdateStatus(requestID, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	// Step 3: Create the recipe (simulates successful extraction)
	t.Log("Step 3: Creating recipe from extracted data...")
	recipe := createTestRecipe()
	recipeID, err := client.CreateRecipe(requestID, testUserID, recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}
	t.Logf("Created recipe: %s", recipeID)

	// Cleanup recipe
	defer func() {
		deleteRecipe(client, recipeID)
	}()

	// Step 4: Update status to COMPLETED
	t.Log("Step 4: Updating status to COMPLETED...")
	err = client.UpdateStatus(requestID, StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}

	t.Log("Full workflow integration test completed successfully")
}

// Helper functions

func createRecipeRequest(client *RecipeRequestClient, url, userID string) (string, error) {
	data := map[string]interface{}{
		"url":     url,
		"status":  StatusRequested,
		"user_id": userID,
	}

	doc, err := client.tablesdb.CreateRow(
		DatabaseID,
		CollectionID,
		"unique()",
		data,
	)
	if err != nil {
		return "", err
	}
	return doc.Id, nil
}

func deleteRecipeRequest(client *RecipeRequestClient, documentID string) error {
	_, err := client.tablesdb.DeleteRow(
		DatabaseID,
		CollectionID,
		documentID,
	)
	return err
}

func deleteRecipe(client *RecipeRequestClient, recipeID string) error {
	_, err := client.tablesdb.DeleteRow(
		DatabaseID,
		RecipeCollectionID,
		recipeID,
	)
	return err
}

func createTestRecipe() *Recipe {
	description := "A delicious test recipe for integration testing"
	prepTime := "PT15M"
	cookTime := "PT30M"
	totalTime := "PT45M"
	calories := "350 kcal"

	return &Recipe{
		Context:     "https://schema.org",
		Type:        "Recipe",
		Name:        "Integration Test Recipe",
		Description: &description,
		Image:       []string{"https://example.com/test-recipe.jpg"},
		PrepTime:    &prepTime,
		CookTime:    &cookTime,
		TotalTime:   &totalTime,
		RecipeYield: []string{"4 servings"},
		RecipeIngredient: []string{
			"1 cup flour",
			"2 eggs",
			"1/2 cup sugar",
		},
		RecipeInstructions: []RecipeInstruction{
			{Type: "HowToStep", Text: "Mix all ingredients together"},
			{Type: "HowToStep", Text: "Bake at 350°F for 30 minutes"},
		},
		RecipeCategory: []string{"Dessert"},
		RecipeCuisine:  []string{"American"},
		Author: &Person{
			Type: "Person",
			Name: "Test Chef",
		},
		Nutrition: &Nutrition{
			Type:     "NutritionInformation",
			Calories: &calories,
		},
	}
}
