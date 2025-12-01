package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// This Appwrite function will be executed every time your function is triggered
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

	// Fetch HTML content - try HTTP client first, fallback to headless browser if needed
	Context.Log(fmt.Sprintf("Fetching recipe from: %s", targetURL))
	recipe, err := fetchRecipeFromURL(targetURL)

	// If HTTP client fails with 403/429 (bot protection), try headless browser
	if err != nil && (strings.Contains(err.Error(), "403") || strings.Contains(err.Error(), "429")) {
		Context.Log("HTTP request blocked, attempting with headless browser...")
		recipe, err = fetchRecipeFromURLWithBrowser(targetURL)
	}

	if err != nil {
		Context.Error(fmt.Sprintf("Error fetching recipe: %v", err))
		return Context.Res.Json(ErrorResponse{
			Error: fmt.Sprintf("Failed to fetch recipe: %v", err),
		}, Context.Res.WithStatusCode(http.StatusInternalServerError))
	}

	if recipe == nil {
		return Context.Res.Json(ErrorResponse{
			Error: "No Recipe structured data found on the page",
		}, Context.Res.WithStatusCode(http.StatusNotFound))
	}

	return Context.Res.Json(toRecipeResponse(targetURL, recipe))
}
