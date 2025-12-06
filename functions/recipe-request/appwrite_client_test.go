package handler

import (
	"os"
	"strings"
	"testing"
)

// MockRecipeRequestStore is a mock implementation for testing
type MockRecipeRequestStore struct {
	CreateRequestFunc  func(url string) (string, error)
	UpdateStatusFunc   func(documentID, status string) error
	CreatedDocumentID  string
	LastStatus         string
	CreateRequestCalls int
	UpdateStatusCalls  int
}

func (m *MockRecipeRequestStore) CreateRequest(url string) (string, error) {
	m.CreateRequestCalls++
	if m.CreateRequestFunc != nil {
		return m.CreateRequestFunc(url)
	}
	m.CreatedDocumentID = "mock-doc-id"
	return m.CreatedDocumentID, nil
}

func (m *MockRecipeRequestStore) UpdateStatus(documentID, status string) error {
	m.UpdateStatusCalls++
	m.LastStatus = status
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(documentID, status)
	}
	return nil
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "StatusRequested",
			constant: StatusRequested,
			expected: "REQUESTED",
		},
		{
			name:     "StatusInProgress",
			constant: StatusInProgress,
			expected: "IN_PROGRESS",
		},
		{
			name:     "StatusCompleted",
			constant: StatusCompleted,
			expected: "COMPLETED",
		},
		{
			name:     "StatusFailed",
			constant: StatusFailed,
			expected: "FAILED",
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

func TestMockRecipeRequestStore_UpdateStatus(t *testing.T) {
	tests := []struct {
		name       string
		documentID string
		status     string
	}{
		{
			name:       "update to IN_PROGRESS",
			documentID: "doc-123",
			status:     StatusInProgress,
		},
		{
			name:       "update to COMPLETED",
			documentID: "doc-456",
			status:     StatusCompleted,
		},
		{
			name:       "update to FAILED",
			documentID: "doc-789",
			status:     StatusFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockRecipeRequestStore{}

			err := mock.UpdateStatus(tt.documentID, tt.status)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if mock.LastStatus != tt.status {
				t.Errorf("got status %s, want %s", mock.LastStatus, tt.status)
			}
			if mock.UpdateStatusCalls != 1 {
				t.Errorf("got %d calls, want 1", mock.UpdateStatusCalls)
			}
		})
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
func TestRecipeRequestClient_Integration(t *testing.T) {
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
	t.Logf("Creating request record for: %s", testURL)

	docID, err := client.CreateRequest(testURL)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	t.Logf("Created document with ID: %s", docID)

	if docID == "" {
		t.Fatal("Expected non-empty document ID")
	}

	// Test updating status to IN_PROGRESS
	t.Log("Updating status to IN_PROGRESS")
	err = client.UpdateStatus(docID, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to update status to IN_PROGRESS: %v", err)
	}

	// Test updating status to COMPLETED
	t.Log("Updating status to COMPLETED")
	err = client.UpdateStatus(docID, StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to update status to COMPLETED: %v", err)
	}

	t.Log("Integration test completed successfully")
}

func TestRecipeRequestClient_Integration_FailedStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("APPWRITE_API_KEY")
	if apiKey == "" {
		t.Skip("APPWRITE_API_KEY not set, skipping integration test")
	}

	client := NewRecipeRequestClient()

	// Test creating a request that will fail
	testURL := "https://example.com/test-recipe-failed"
	t.Logf("Creating request record for: %s", testURL)

	docID, err := client.CreateRequest(testURL)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	t.Logf("Created document with ID: %s", docID)

	// Update to IN_PROGRESS
	err = client.UpdateStatus(docID, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to update status to IN_PROGRESS: %v", err)
	}

	// Test updating status to FAILED
	// Note: This test requires the FAILED status to be deployed to Appwrite
	// Run `appwrite deploy` to update the schema if this test fails
	t.Log("Updating status to FAILED")
	err = client.UpdateStatus(docID, StatusFailed)
	if err != nil {
		// Check if it's a schema error (FAILED status not deployed yet)
		if strings.Contains(err.Error(), "Invalid document structure") ||
			strings.Contains(err.Error(), "must be one of") {
			t.Skipf("FAILED status not yet deployed to Appwrite schema. Run 'appwrite deploy' to update. Error: %v", err)
		}
		t.Fatalf("Failed to update status to FAILED: %v", err)
	}

	t.Log("Failed status integration test completed successfully")
}

