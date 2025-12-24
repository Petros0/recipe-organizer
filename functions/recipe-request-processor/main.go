package handler

import (
	"encoding/json"
	"errors"
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
	// Parse event payload - contains the full document data from the database event
	var payload DocumentEventPayload
	if bodyText := Context.Req.BodyText(); bodyText != "" {
		if err := json.Unmarshal([]byte(bodyText), &payload); err != nil {
			return Context.Res.Json(ErrorResponse{
				Error: fmt.Sprintf("Invalid event payload: %v", err),
			}, Context.Res.WithStatusCode(http.StatusBadRequest))
		}
	}

	if err := ValidatePayload(payload); err != nil {
		return Context.Res.Json(ErrorResponse{
			Error: err.Error(),
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

	// Create structured logger with request context
	logger := NewLogger(Context, payload.ID, payload.URL, payload.UserID)
	logger.Info("main", "Processing recipe request")

	// Initialize recipe request client for tracking
	requestClient := NewRecipeRequestClient()

	// Update status to IN_PROGRESS before fetching
	if err := requestClient.UpdateStatus(payload.ID, StatusInProgress); err != nil {
		logger.Error("main", "Error updating status to IN_PROGRESS", map[string]interface{}{
			"error": err.Error(),
		})
		return Context.Res.Json(ErrorResponse{
			Error: "Error updating status to IN_PROGRESS",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	logger.Info("main", "Status updated to IN_PROGRESS")

	// Create strategy executor with HTTP client first, then Firecrawl as fallback
	// Firecrawl handles bot protection and can use LLM extraction if no JSON-LD is found
	executor := NewStrategyExecutor(
		NewHTTPClientStrategy(logger),
		NewFirecrawlStrategy(logger),
	)

	// Fetch recipe using the strategy executor
	recipe, err := executor.Execute(payload.URL, logger)

	if err != nil {
		logger.Error("main", "Error fetching recipe", map[string]interface{}{
			"error": err.Error(),
		})
		// Update status to FAILED
		if updateErr := requestClient.UpdateStatus(payload.ID, StatusFailed); updateErr != nil {
			logger.Error("main", "Error updating status to FAILED", map[string]interface{}{
				"error": updateErr.Error(),
			})
			return Context.Res.Json(ErrorResponse{
				Error: "Error updating status to FAILED",
			}, Context.Res.WithStatusCode(http.StatusInternalServerError))
		}

		return Context.Res.Json(ErrorResponse{
			Error: "Failed to fetch recipe",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	if recipe == nil {
		logger.Error("main", "No recipe structured data found on page")
		// Update status to FAILED when no recipe found
		if updateErr := requestClient.UpdateStatus(payload.ID, StatusFailed); updateErr != nil {
			logger.Error("main", "Error updating status to FAILED", map[string]interface{}{
				"error": updateErr.Error(),
			})
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
		logger.Error("main", "Error saving recipe to database", map[string]interface{}{
			"error": err.Error(),
		})
		// Update status to FAILED since we couldn't save
		if updateErr := requestClient.UpdateStatus(payload.ID, StatusFailed); updateErr != nil {
			logger.Error("main", "Error updating status to FAILED", map[string]interface{}{
				"error": updateErr.Error(),
			})
		}
		return Context.Res.Json(ErrorResponse{
			Error: "Failed to save recipe to database",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	logger.Info("main", "Recipe saved to database", map[string]interface{}{
		"recipe_id": recipeID,
	})

	// Update status to COMPLETED on success
	if err := requestClient.UpdateStatus(payload.ID, StatusCompleted); err != nil {
		logger.Error("main", "Error updating status to COMPLETED", map[string]interface{}{
			"error": err.Error(),
		})
		return Context.Res.Json(ErrorResponse{
			Error: "Error updating status to COMPLETED",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	logger.WithDuration("main", "Recipe processing completed", map[string]interface{}{
		"recipe_id": recipeID,
	})

	return Context.Res.Json(toRecipeResponse(payload.URL, recipe))
}

func ValidatePayload(payload DocumentEventPayload) error {
	if payload.ID == "" {
		return errors.New("$id is required in event payload")
	}
	if payload.URL == "" {
		return errors.New("url is required in event payload")
	}
	if payload.UserID == "" {
		return errors.New("user_id is required in event payload")
	}

	return nil
}