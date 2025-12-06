package handler

import (
	"os"
	"testing"
)

// MockRecipeRequestStore is a mock implementation for testing
type MockRecipeRequestStore struct {
	UpdateStatusFunc  func(documentID, status string) error
	LastDocumentID    string
	LastStatus        string
	UpdateStatusCalls int
}

func (m *MockRecipeRequestStore) UpdateStatus(documentID, status string) error {
	m.UpdateStatusCalls++
	m.LastDocumentID = documentID
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
			if mock.LastDocumentID != tt.documentID {
				t.Errorf("got documentID %s, want %s", mock.LastDocumentID, tt.documentID)
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
// Note: These tests require an existing document ID to update status
func TestRecipeRequestClient_UpdateStatus_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	apiKey := os.Getenv("APPWRITE_API_KEY")
	if apiKey == "" {
		t.Skip("APPWRITE_API_KEY not set, skipping integration test")
	}

	testDocID := os.Getenv("TEST_DOCUMENT_ID")
	if testDocID == "" {
		t.Skip("TEST_DOCUMENT_ID not set, skipping integration test. Create a request via recipe-request first.")
	}

	client := NewRecipeRequestClient()

	// Test updating status to IN_PROGRESS
	t.Log("Updating status to IN_PROGRESS")
	err := client.UpdateStatus(testDocID, StatusInProgress)
	if err != nil {
		t.Fatalf("Failed to update status to IN_PROGRESS: %v", err)
	}

	// Test updating status to COMPLETED
	t.Log("Updating status to COMPLETED")
	err = client.UpdateStatus(testDocID, StatusCompleted)
	if err != nil {
		t.Fatalf("Failed to update status to COMPLETED: %v", err)
	}

	t.Log("Integration test completed successfully")
}
