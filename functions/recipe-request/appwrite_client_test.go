package handler

import (
	"os"
	"testing"
)

// MockRecipeRequestStore is a mock implementation for testing
type MockRecipeRequestStore struct {
	CreateRequestFunc  func(url string) (string, error)
	CreatedDocumentID  string
	CreateRequestCalls int
}

func (m *MockRecipeRequestStore) CreateRequest(url string) (string, error) {
	m.CreateRequestCalls++
	if m.CreateRequestFunc != nil {
		return m.CreateRequestFunc(url)
	}
	m.CreatedDocumentID = "mock-doc-id"
	return m.CreatedDocumentID, nil
}

func TestStatusConstants(t *testing.T) {
	if StatusRequested != "REQUESTED" {
		t.Errorf("got %s, want REQUESTED", StatusRequested)
	}
}

func TestDefaultConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "DefaultEndpoint",
			constant: DefaultEndpoint,
			expected: "https://fra.cloud.appwrite.io/v1",
		},
		{
			name:     "DefaultProjectID",
			constant: DefaultProjectID,
			expected: "691f8b990030db50617a",
		},
		{
			name:     "DatabaseID",
			constant: DatabaseID,
			expected: "6930a343001607ad7cbd",
		},
		{
			name:     "CollectionID",
			constant: CollectionID,
			expected: "6930a34300165ad1d129",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("got %s, want %s", tt.constant, tt.expected)
			}
		})
	}
}

func TestMockRecipeRequestStore_CreateRequest(t *testing.T) {
	mock := &MockRecipeRequestStore{}

	docID, err := mock.CreateRequest("https://example.com/recipe")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if docID != "mock-doc-id" {
		t.Errorf("got docID %s, want mock-doc-id", docID)
	}
	if mock.CreateRequestCalls != 1 {
		t.Errorf("got %d calls, want 1", mock.CreateRequestCalls)
	}
}

func TestNewRecipeRequestClient(t *testing.T) {
	// Test that client is created without panic
	client := NewRecipeRequestClient()

	if client == nil {
		t.Error("expected non-nil client")
	}
	if client.databases == nil {
		t.Error("expected non-nil databases service")
	}
}

func TestNewRecipeRequestClient_WithEnvVars(t *testing.T) {
	// Save original env vars
	origEndpoint := os.Getenv("APPWRITE_ENDPOINT")
	origProjectID := os.Getenv("APPWRITE_PROJECT_ID")
	origAPIKey := os.Getenv("APPWRITE_API_KEY")

	// Set test env vars
	os.Setenv("APPWRITE_ENDPOINT", "https://test.appwrite.io/v1")
	os.Setenv("APPWRITE_PROJECT_ID", "test-project")
	os.Setenv("APPWRITE_API_KEY", "test-api-key")

	defer func() {
		// Restore original env vars
		if origEndpoint != "" {
			os.Setenv("APPWRITE_ENDPOINT", origEndpoint)
		} else {
			os.Unsetenv("APPWRITE_ENDPOINT")
		}
		if origProjectID != "" {
			os.Setenv("APPWRITE_PROJECT_ID", origProjectID)
		} else {
			os.Unsetenv("APPWRITE_PROJECT_ID")
		}
		if origAPIKey != "" {
			os.Setenv("APPWRITE_API_KEY", origAPIKey)
		} else {
			os.Unsetenv("APPWRITE_API_KEY")
		}
	}()

	client := NewRecipeRequestClient()

	if client == nil {
		t.Error("expected non-nil client")
	}
}

// Integration tests - only run when APPWRITE_API_KEY is set
func TestRecipeRequestClient_CreateRequest_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("APPWRITE_API_KEY")
	if apiKey == "" {
		t.Skip("APPWRITE_API_KEY not set, skipping integration test")
	}

	client := NewRecipeRequestClient()

	// Test creating a request
	testURL := "https://example.com/test-recipe-integration"
	testUserID := "test-user-123"
	t.Logf("Creating request record for: %s", testURL)

	docID, err := client.CreateRequest(testURL, testUserID)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	t.Logf("Created document with ID: %s", docID)

	if docID == "" {
		t.Fatal("Expected non-empty document ID")
	}

	t.Log("Integration test completed successfully")
}
