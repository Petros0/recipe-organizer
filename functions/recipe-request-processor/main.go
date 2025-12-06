package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// ProcessorRequestBody represents the request body for the processor
// This is typically triggered by an Appwrite event when a recipe request is created
type ProcessorRequestBody struct {
	DocumentID string `json:"documentId"`
	URL        string `json:"url"`
}

// This Appwrite function will be executed every time your function is triggered
func Main(Context openruntimes.Context) openruntimes.Response {
	// Handle ping endpoint
	if Context.Req.Path == "/ping" {
		return Context.Res.Text("Pong")
	}

	// Parse request body - expects documentId and url from event trigger
	var body ProcessorRequestBody
	if bodyText := Context.Req.BodyText(); bodyText != "" {
		if err := json.Unmarshal([]byte(bodyText), &body); err != nil {
			return Context.Res.Json(ErrorResponse{
				Error: fmt.Sprintf("Invalid request body: %v", err),
			}, Context.Res.WithStatusCode(http.StatusBadRequest))
		}
	}

	// Validate required fields
	if body.DocumentID == "" {
		return Context.Res.Json(ErrorResponse{
			Error: "documentId is required in request body",
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	if body.URL == "" {
		return Context.Res.Json(ErrorResponse{
			Error: "url is required in request body",
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	// Initialize recipe request client for tracking
	requestClient := NewRecipeRequestClient()

	// Update status to IN_PROGRESS before fetching
	if err := requestClient.UpdateStatus(body.DocumentID, StatusInProgress); err != nil {
		Context.Error(fmt.Sprintf("Error updating status to IN_PROGRESS: %v", err))
		return Context.Res.Json(ErrorResponse{
			Error: fmt.Sprintf("Error updating status to IN_PROGRESS: %v", err),
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	Context.Log(fmt.Sprintf("Processing request %s for URL: %s", body.DocumentID, body.URL))

	// Create strategy executor with HTTP client first, then Firecrawl as fallback
	// Firecrawl handles bot protection and can use LLM extraction if no JSON-LD is found
	executor := NewStrategyExecutor(
		&HTTPClientStrategy{},
		NewFirecrawlStrategy(),
	)

	// Fetch recipe using the strategy executor
	Context.Log(fmt.Sprintf("Fetching recipe from: %s", body.URL))
	recipe, err := executor.Execute(body.URL, Context.Log)

	if err != nil {
		Context.Error(fmt.Sprintf("Error fetching recipe: %v", err))
		// Update status to FAILED
		if updateErr := requestClient.UpdateStatus(body.DocumentID, StatusFailed); updateErr != nil {
			Context.Error(fmt.Sprintf("Error updating status to FAILED: %v", updateErr))
			return Context.Res.Json(ErrorResponse{
				Error: fmt.Sprintf("Error updating status to FAILED: %v", updateErr),
			}, Context.Res.WithStatusCode(http.StatusInternalServerError))
		}

		return Context.Res.Json(ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch recipe: %v", err),
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	if recipe == nil {
		// Update status to FAILED when no recipe found
		if updateErr := requestClient.UpdateStatus(body.DocumentID, StatusFailed); updateErr != nil {
			Context.Error(fmt.Sprintf("Error updating status to FAILED: %v", updateErr))
			return Context.Res.Json(ErrorResponse{
				Error: fmt.Sprintf("Error updating status to FAILED: %v", updateErr),
			}, Context.Res.WithStatusCode(http.StatusInternalServerError))
		}
		return Context.Res.Json(ErrorResponse{
			Error: "No Recipe structured data found on the page",
		}, Context.Res.WithStatusCode(http.StatusNotFound))
	}

	// Update status to COMPLETED on success
	if err := requestClient.UpdateStatus(body.DocumentID, StatusCompleted); err != nil {
		Context.Error(fmt.Sprintf("Error updating status to COMPLETED: %v", err))
		return Context.Res.Json(ErrorResponse{
			Error: fmt.Sprintf("Error updating status to COMPLETED: %v", err),
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	return Context.Res.Json(toRecipeResponse(body.URL, recipe))
}
