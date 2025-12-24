package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// DocumentEventPayload represents the Appwrite document event payload
// This is triggered by database events when a document is created/updated
type DocumentEventPayload struct {
	ID           string `json:"$id"`
	DatabaseID   string `json:"$databaseId"`
	CollectionID string `json:"$collectionId"`
	CreatedAt    string `json:"$createdAt"`
	UpdatedAt    string `json:"$updatedAt"`
	URL          string `json:"url"`
	Status       string `json:"status"`
	UserID       string `json:"user_id"`
}

// This Appwrite function will be executed every time your function is triggered
func Main(Context openruntimes.Context) openruntimes.Response {
	// Handle ping endpoint
	if Context.Req.Path == "/ping" {
		return Context.Res.Text("Pong")
	}

	// Parse event payload - contains the full document data from the database event
	var payload DocumentEventPayload
	if bodyText := Context.Req.BodyText(); bodyText != "" {
		if err := json.Unmarshal([]byte(bodyText), &payload); err != nil {
			return Context.Res.Json(ErrorResponse{
				Error: fmt.Sprintf("Invalid event payload: %v", err),
			}, Context.Res.WithStatusCode(http.StatusBadRequest))
		}
	}

	// Validate required fields from event payload
	if payload.ID == "" {
		return Context.Res.Json(ErrorResponse{
			Error: "$id is required in event payload",
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	if payload.URL == "" {
		return Context.Res.Json(ErrorResponse{
			Error: "url is required in event payload",
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	if payload.UserID == "" {
		return Context.Res.Json(ErrorResponse{
			Error: "user_id is required in event payload",
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	// Only process documents with REQUESTED status to avoid infinite loops
	// when we update the status ourselves
	if payload.Status != StatusRequested {
		Context.Log(fmt.Sprintf("Skipping document %s with status %s (only processing REQUESTED)", payload.ID, payload.Status))
		return Context.Res.Json(map[string]string{
			"message": fmt.Sprintf("Skipped: document status is %s, not REQUESTED", payload.Status),
		})
	}

	// Initialize recipe request client for tracking
	requestClient := NewRecipeRequestClient()

	// Update status to IN_PROGRESS before fetching
	if err := requestClient.UpdateStatus(payload.ID, StatusInProgress); err != nil {
		Context.Error(fmt.Sprintf("Error updating status to IN_PROGRESS: %v", err))
		return Context.Res.Json(ErrorResponse{
			Error: "Error updating status to IN_PROGRESS",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	Context.Log(fmt.Sprintf("Processing request %s for URL: %s", payload.ID, payload.URL))

	// Set up loggers for detailed extraction logging
	SetFirecrawlLogger(func(msgs ...interface{}) {
		Context.Log(msgs...)
	})
	SetParserLogger(func(msgs ...interface{}) {
		Context.Log(msgs...)
	})

	// Create strategy executor with HTTP client first, then Firecrawl as fallback
	// Firecrawl handles bot protection and can use LLM extraction if no JSON-LD is found
	executor := NewStrategyExecutor(
		&HTTPClientStrategy{},
		NewFirecrawlStrategy(),
	)

	// Fetch recipe using the strategy executor
	Context.Log(fmt.Sprintf("Fetching recipe from: %s", payload.URL))
	recipe, err := executor.Execute(payload.URL, Context.Log)

	if err != nil {
		Context.Error(fmt.Sprintf("Error fetching recipe: %v", err))
		// Update status to FAILED
		if updateErr := requestClient.UpdateStatus(payload.ID, StatusFailed); updateErr != nil {
			Context.Error(fmt.Sprintf("Error updating status to FAILED: %v", updateErr))
			return Context.Res.Json(ErrorResponse{
				Error: "Error updating status to FAILED",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError))
		}

		return Context.Res.Json(ErrorResponse{
			Error: "Failed to fetch recipe",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	if recipe == nil {
		// Update status to FAILED when no recipe found
		if updateErr := requestClient.UpdateStatus(payload.ID, StatusFailed); updateErr != nil {
			Context.Error(fmt.Sprintf("Error updating status to FAILED: %v", updateErr))
			return Context.Res.Json(ErrorResponse{
				Error: "Error updating status to FAILED",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError))
		}
		return Context.Res.Json(ErrorResponse{
			Error: "No Recipe structured data found on the page",
		}, Context.Res.WithStatusCode(http.StatusNotFound))
	}

	// Save recipe to database
	recipeID, err := requestClient.CreateRecipe(payload.ID, payload.UserID, recipe)
	if err != nil {
		Context.Error(fmt.Sprintf("Error saving recipe to database: %v", err))
		// Update status to FAILED since we couldn't save
		if updateErr := requestClient.UpdateStatus(payload.ID, StatusFailed); updateErr != nil {
			Context.Error(fmt.Sprintf("Error updating status to FAILED: %v", updateErr))
		}
		return Context.Res.Json(ErrorResponse{
			Error: "Failed to save recipe to database",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}
	Context.Log(fmt.Sprintf("Recipe saved with ID: %s", recipeID))

	// Update status to COMPLETED on success
	if err := requestClient.UpdateStatus(payload.ID, StatusCompleted); err != nil {
		Context.Error(fmt.Sprintf("Error updating status to COMPLETED: %v", err))
		return Context.Res.Json(ErrorResponse{
			Error: "Error updating status to COMPLETED",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	return Context.Res.Json(toRecipeResponse(payload.URL, recipe))
}
