package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// Main is the Appwrite function entry point
func Main(Context openruntimes.Context) openruntimes.Response {
	// Handle ping endpoint
	if Context.Req.Path == "/ping" {
		return Context.Res.Text("Pong")
	}

	// Extract URL from request
	var targetURL string

	// Try to get URL from query parameter first
	if urlParam, ok := Context.Req.Query["url"]; ok && urlParam != "" {
		targetURL = urlParam
	} else if bodyText := Context.Req.BodyText(); bodyText != "" {
		// Try to parse JSON body
		var body RequestBody
		if err := json.Unmarshal([]byte(bodyText), &body); err == nil && body.URL != "" {
			targetURL = body.URL
		}
	}

	// Validate URL
	if targetURL == "" {
		return Context.Res.Json(ErrorResponse{
			Error: "URL parameter is required. Provide 'url' as query parameter or in JSON body.",
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	// Validate URL format
	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return Context.Res.Json(ErrorResponse{
			Error: fmt.Sprintf("Invalid URL format: %s", targetURL),
		}, Context.Res.WithStatusCode(http.StatusBadRequest))
	}

	// Initialize recipe request client
	requestClient := NewRecipeRequestClient()

	// Create request record with REQUESTED status
	documentID, err := requestClient.CreateRequest(targetURL)
	if err != nil {
		Context.Error(fmt.Sprintf("Error creating request record: %v", err))
		return Context.Res.Json(ErrorResponse{
			Error: "Error creating request record",
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	Context.Log(fmt.Sprintf("Created request record: %s for URL: %s", documentID, targetURL))

	// Return success response with document ID
	return Context.Res.Json(SuccessResponse{
		DocumentID: documentID,
		Status:     StatusRequested,
		URL:        targetURL,
	})
}
