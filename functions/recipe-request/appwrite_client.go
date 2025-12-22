package handler

import (
	"fmt"
	"os"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/permission"
	"github.com/appwrite/sdk-for-go/role"
)

// Status constants for recipe request tracking
const (
	StatusRequested = "REQUESTED"
)

// Default Appwrite configuration
const (
	DefaultEndpoint  = "https://fra.cloud.appwrite.io/v1"
	DefaultProjectID = "691f8b990030db50617a"
	DatabaseID       = "6930a343001607ad7cbd"
	CollectionID     = "6930a34300165ad1d129"
)

// RecipeRequestStore defines the interface for recipe request operations
type RecipeRequestStore interface {
	CreateRequest(url, userID string) (string, error)
}

// RecipeRequestClient handles database operations for recipe requests
type RecipeRequestClient struct {
	databases *databases.Databases
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

	appwriteClient := client.New()
	appwriteClient.Endpoint = endpoint
	appwriteClient.AddHeader("X-Appwrite-Project", projectID)
	appwriteClient.AddHeader("X-Appwrite-Key", apiKey)

	return &RecipeRequestClient{
		databases: databases.New(appwriteClient),
	}
}

// CreateRequest creates a new recipe request record with REQUESTED status
func (c *RecipeRequestClient) CreateRequest(url, userID string) (string, error) {
	data := map[string]interface{}{
		"url":     url,
		"status":  StatusRequested,
		"user_id": userID,
	}

	// Set permissions to allow only the user who created the request to read it
	permissions := []string{
		permission.Read(role.User(userID, "")),
		permission.Update(role.User(userID, "")),
		permission.Delete(role.User(userID, "")),
	}

	doc, err := c.databases.CreateDocument(
		DatabaseID,
		CollectionID,
		id.Unique(),
		data,
		c.databases.WithCreateDocumentPermissions(permissions),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create recipe request document: %w", err)
	}

	return doc.Id, nil
}
