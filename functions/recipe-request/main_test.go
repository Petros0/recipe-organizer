package handler

import (
	"os"
	"testing"
)

// TestCreateRequest_Integration tests creating a recipe request record
func TestCreateRequest_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("APPWRITE_API_KEY")
	if apiKey == "" {
		t.Skip("APPWRITE_API_KEY not set, skipping integration test")
	}

	client := NewRecipeRequestClient()

	// Test creating a request
	testURL := "https://example.com/test-recipe"
	t.Logf("Creating request record for: %s", testURL)

	docID, err := client.CreateRequest(testURL)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	if docID == "" {
		t.Fatal("Expected non-empty document ID")
	}

	t.Logf("Created document with ID: %s", docID)
	t.Log("Integration test completed successfully")
}
