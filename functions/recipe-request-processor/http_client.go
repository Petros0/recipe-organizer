package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// HTTPClientStrategy implements FetchStrategy using standard HTTP client
type HTTPClientStrategy struct{}

// Name returns the strategy name for logging
func (s *HTTPClientStrategy) Name() string {
	return "HTTP Client"
}

// CanRetry returns true if the error indicates bot protection (403/429) or no JSON-LD found
func (s *HTTPClientStrategy) CanRetry(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	// Retry on bot protection (403/429) or when no JSON-LD is found
	return strings.Contains(errStr, "403") || strings.Contains(errStr, "429") || errors.Is(err, ErrNoJSONLD)
}

// Fetch fetches HTML from URL and extracts Recipe JSON-LD using HTTP client
func (s *HTTPClientStrategy) Fetch(urlStr string) (*Recipe, error) {
	// Create cookie jar to handle sessions and cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Create HTTP client with timeout and cookie jar to handle sessions
	client := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	// Create request with realistic browser headers to avoid bot detection
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set realistic browser headers to mimic a real browser request
	setBrowserHeaders(req, urlStr)

	// Fetch the page
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract recipe from HTML body
	recipe, err := extractRecipeFromHTML(string(body))
	if err != nil {
		return nil, err
	}

	// If no recipe found, return ErrNoJSONLD so we can retry with Firecrawl
	if recipe == nil {
		return nil, ErrNoJSONLD
	}

	return recipe, nil
}

// setBrowserHeaders sets realistic browser headers to avoid bot detection
func setBrowserHeaders(req *http.Request, targetURL string) {
	// Parse URL to get domain for Referer
	parsedURL, _ := url.Parse(targetURL)
	baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	// Set realistic browser headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	// Note: Don't set Accept-Encoding manually - Go's http.Transport handles
	// compression automatically and will decompress the response
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Cache-Control", "max-age=0")

	// Set Referer if we have a base URL
	if baseURL != "" {
		req.Header.Set("Referer", baseURL)
	}
}

