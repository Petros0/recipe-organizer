package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

// Recipe represents a schema.org Recipe structured data
type Recipe struct {
	Context            string              `json:"@context,omitempty"`
	Type               string              `json:"@type,omitempty"`
	Name               string              `json:"name"`
	Image              []string            `json:"image,omitempty"`
	Author             *Person             `json:"author,omitempty"`
	Description        *string             `json:"description,omitempty"`
	PrepTime           *string             `json:"prepTime,omitempty"`
	CookTime           *string             `json:"cookTime,omitempty"`
	TotalTime          *string             `json:"totalTime,omitempty"`
	RecipeYield        *string             `json:"recipeYield,omitempty"`
	RecipeIngredient   []string            `json:"recipeIngredient,omitempty"`
	RecipeInstructions []RecipeInstruction `json:"recipeInstructions,omitempty"`
	RecipeCategory     *string             `json:"recipeCategory,omitempty"`
	RecipeCuisine      *string             `json:"recipeCuisine,omitempty"`
	Nutrition          *Nutrition          `json:"nutrition,omitempty"`
	Keywords           *string             `json:"keywords,omitempty"`
	DatePublished      *string             `json:"datePublished,omitempty"`
	DateModified       *string             `json:"dateModified,omitempty"`
}

// Person represents a schema.org Person
type Person struct {
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// Nutrition represents schema.org NutritionInformation
type Nutrition struct {
	Type                string  `json:"@type,omitempty"`
	Calories            *string `json:"calories,omitempty"`
	FatContent          *string `json:"fatContent,omitempty"`
	SaturatedFatContent *string `json:"saturatedFatContent,omitempty"`
	CholesterolContent  *string `json:"cholesterolContent,omitempty"`
	SodiumContent       *string `json:"sodiumContent,omitempty"`
	CarbohydrateContent *string `json:"carbohydrateContent,omitempty"`
	FiberContent        *string `json:"fiberContent,omitempty"`
	SugarContent        *string `json:"sugarContent,omitempty"`
	ProteinContent      *string `json:"proteinContent,omitempty"`
}

// RecipeInstruction represents a schema.org HowToStep or HowToSection
type RecipeInstruction struct {
	Type            string              `json:"@type,omitempty"`
	Text            string              `json:"text,omitempty"`
	Name            string              `json:"name,omitempty"`
	URL             string              `json:"url,omitempty"`
	ItemListElement []RecipeInstruction `json:"itemListElement,omitempty"` // For HowToSection
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RequestBody represents the JSON request body
type RequestBody struct {
	URL string `json:"url"`
}

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

	return Context.Res.Json(recipe)
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
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
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

// fetchRecipeFromURL fetches HTML from URL and extracts Recipe JSON-LD
func fetchRecipeFromURL(urlStr string) (*Recipe, error) {
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
	return extractRecipeFromHTML(string(body))
}

// extractRecipeFromHTML extracts Recipe JSON-LD from HTML content
func extractRecipeFromHTML(htmlContent string) (*Recipe, error) {
	// Parse HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Find all JSON-LD script tags
	var recipe *Recipe
	doc.Find("script[type='application/ld+json']").Each(func(i int, s *goquery.Selection) {
		if recipe != nil {
			return // Already found a recipe
		}

		jsonLD := s.Text()
		if jsonLD == "" {
			return
		}

		// Try to parse as single object or array
		var data interface{}
		if err := json.Unmarshal([]byte(jsonLD), &data); err != nil {
			return // Skip invalid JSON
		}

		// Handle different JSON-LD formats
		recipe = extractRecipeFromJSONLD(data)
	})

	return recipe, nil
}

// fetchRecipeFromURLWithBrowser uses a headless browser to fetch the page
// This is a fallback when HTTP client is blocked by bot protection
func fetchRecipeFromURLWithBrowser(urlStr string) (*Recipe, error) {
	// Launch browser with stealth options to avoid detection
	l := launcher.New().
		Headless(true).
		Set("disable-blink-features", "AutomationControlled").
		Set("excludeSwitches", "enable-automation").
		NoSandbox(true).
		Set("disable-dev-shm-usage", "true")

	// Launch browser and get control URL
	controlURL, err := l.Launch()
	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	// Connect to browser with timeout
	browser := rod.New().
		ControlURL(controlURL).
		Timeout(20 * time.Second)

	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}
	defer browser.MustClose()

	// Create page
	page := browser.MustPage("")
	defer page.MustClose()

	// Set realistic browser properties to avoid detection before navigating
	// Set extra headers (key-value pairs)
	page.MustSetExtraHeaders(
		"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"Accept-Language", "en-US,en;q=0.9",
	)
	page.MustSetViewport(1920, 1080, 1, false)

	// Navigate to the URL
	page.MustNavigate(urlStr)

	// Wait for page to load
	page.MustWaitLoad()

	// Wait a bit for any JavaScript to execute and render content
	time.Sleep(1 * time.Second)

	// Get the HTML content
	html, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("failed to get page HTML: %w", err)
	}

	// Extract recipe from HTML
	return extractRecipeFromHTML(html)
}

// extractRecipeFromJSONLD extracts Recipe from various JSON-LD formats
func extractRecipeFromJSONLD(data interface{}) *Recipe {
	switch v := data.(type) {
	case map[string]interface{}:
		// Single object - check if it's a Recipe or has @graph
		if recipe := extractRecipeFromObject(v); recipe != nil {
			return recipe
		}
		// Check for @graph array
		if graph, ok := v["@graph"].([]interface{}); ok {
			return extractRecipeFromArray(graph)
		}
	case []interface{}:
		// Array of objects
		return extractRecipeFromArray(v)
	}
	return nil
}

// extractRecipeFromArray extracts Recipe from an array of objects
func extractRecipeFromArray(arr []interface{}) *Recipe {
	for _, item := range arr {
		if obj, ok := item.(map[string]interface{}); ok {
			if recipe := extractRecipeFromObject(obj); recipe != nil {
				return recipe
			}
		}
	}
	return nil
}

// extractRecipeFromObject extracts Recipe from a single object
func extractRecipeFromObject(obj map[string]interface{}) *Recipe {
	// Check if it's a Recipe type
	typeVal, ok := obj["@type"].(string)
	if !ok {
		return nil
	}

	// Handle both "Recipe" and "https://schema.org/Recipe"
	if typeVal != "Recipe" && !strings.Contains(typeVal, "Recipe") {
		return nil
	}

	// Parse Recipe
	recipe := &Recipe{}

	// Required: name
	if name, ok := obj["name"].(string); ok && name != "" {
		recipe.Name = name
	} else {
		// Name is required, skip if not present
		return nil
	}

	// Required: image (can be string or array)
	if imageVal, ok := obj["image"]; ok {
		recipe.Image = parseImage(imageVal)
		if len(recipe.Image) == 0 {
			// Image is required, skip if empty
			return nil
		}
	} else {
		// Image is required, skip if not present
		return nil
	}

	// Optional fields
	recipe.Context = getString(obj, "@context")
	recipe.Type = getString(obj, "@type")

	if desc := getStringPtr(obj, "description"); desc != nil {
		recipe.Description = desc
	}
	if prepTime := getStringPtr(obj, "prepTime"); prepTime != nil {
		recipe.PrepTime = prepTime
	}
	if cookTime := getStringPtr(obj, "cookTime"); cookTime != nil {
		recipe.CookTime = cookTime
	}
	if totalTime := getStringPtr(obj, "totalTime"); totalTime != nil {
		recipe.TotalTime = totalTime
	}
	if yield := getStringPtr(obj, "recipeYield"); yield != nil {
		recipe.RecipeYield = yield
	}
	if category := getStringPtr(obj, "recipeCategory"); category != nil {
		recipe.RecipeCategory = category
	}
	if cuisine := getStringPtr(obj, "recipeCuisine"); cuisine != nil {
		recipe.RecipeCuisine = cuisine
	}
	if keywords := getStringPtr(obj, "keywords"); keywords != nil {
		recipe.Keywords = keywords
	}
	if datePublished := getStringPtr(obj, "datePublished"); datePublished != nil {
		recipe.DatePublished = datePublished
	}
	if dateModified := getStringPtr(obj, "dateModified"); dateModified != nil {
		recipe.DateModified = dateModified
	}

	// Recipe ingredients
	if ingredients := getStringArray(obj, "recipeIngredient"); len(ingredients) > 0 {
		recipe.RecipeIngredient = ingredients
	}

	// Recipe instructions
	if instructions := parseInstructions(obj["recipeInstructions"]); len(instructions) > 0 {
		recipe.RecipeInstructions = instructions
	}

	// Author
	if author := parseAuthor(obj["author"]); author != nil {
		recipe.Author = author
	}

	// Nutrition
	if nutrition := parseNutrition(obj["nutrition"]); nutrition != nil {
		recipe.Nutrition = nutrition
	}

	return recipe
}

// parseImage parses image field which can be string or array
func parseImage(imageVal interface{}) []string {
	var images []string

	switch v := imageVal.(type) {
	case string:
		if v != "" {
			images = append(images, v)
		}
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				images = append(images, str)
			} else if obj, ok := item.(map[string]interface{}); ok {
				// Could be ImageObject with url property
				if url := getString(obj, "url"); url != "" {
					images = append(images, url)
				}
			}
		}
	case map[string]interface{}:
		// ImageObject
		if url := getString(v, "url"); url != "" {
			images = append(images, url)
		}
	}

	return images
}

// parseInstructions parses recipeInstructions which can be HowToStep, HowToSection, or array
func parseInstructions(instructionsVal interface{}) []RecipeInstruction {
	var instructions []RecipeInstruction

	if instructionsVal == nil {
		return instructions
	}

	switch v := instructionsVal.(type) {
	case []interface{}:
		for _, item := range v {
			if inst := parseInstructionItem(item); inst != nil {
				instructions = append(instructions, *inst)
			}
		}
	default:
		if inst := parseInstructionItem(v); inst != nil {
			instructions = append(instructions, *inst)
		}
	}

	return instructions
}

// parseInstructionItem parses a single instruction item
func parseInstructionItem(item interface{}) *RecipeInstruction {
	obj, ok := item.(map[string]interface{})
	if !ok {
		return nil
	}

	inst := &RecipeInstruction{}
	inst.Type = getString(obj, "@type")
	inst.Text = getString(obj, "text")
	inst.Name = getString(obj, "name")
	inst.URL = getString(obj, "url")

	// Handle HowToSection with itemListElement
	if itemList, ok := obj["itemListElement"].([]interface{}); ok {
		for _, item := range itemList {
			if subInst := parseInstructionItem(item); subInst != nil {
				inst.ItemListElement = append(inst.ItemListElement, *subInst)
			}
		}
	}

	return inst
}

// parseAuthor parses author field
func parseAuthor(authorVal interface{}) *Person {
	if authorVal == nil {
		return nil
	}

	author := &Person{}

	switch v := authorVal.(type) {
	case string:
		author.Name = v
	case map[string]interface{}:
		author.Type = getString(v, "@type")
		author.Name = getString(v, "name")
		author.URL = getString(v, "url")
	case []interface{}:
		// Array of authors - take first
		if len(v) > 0 {
			return parseAuthor(v[0])
		}
	}

	if author.Name == "" {
		return nil
	}

	return author
}

// parseNutrition parses nutrition field
func parseNutrition(nutritionVal interface{}) *Nutrition {
	if nutritionVal == nil {
		return nil
	}

	obj, ok := nutritionVal.(map[string]interface{})
	if !ok {
		return nil
	}

	nutrition := &Nutrition{}
	nutrition.Type = getString(obj, "@type")
	nutrition.Calories = getStringPtr(obj, "calories")
	nutrition.FatContent = getStringPtr(obj, "fatContent")
	nutrition.SaturatedFatContent = getStringPtr(obj, "saturatedFatContent")
	nutrition.CholesterolContent = getStringPtr(obj, "cholesterolContent")
	nutrition.SodiumContent = getStringPtr(obj, "sodiumContent")
	nutrition.CarbohydrateContent = getStringPtr(obj, "carbohydrateContent")
	nutrition.FiberContent = getStringPtr(obj, "fiberContent")
	nutrition.SugarContent = getStringPtr(obj, "sugarContent")
	nutrition.ProteinContent = getStringPtr(obj, "proteinContent")

	return nutrition
}

// Helper functions for safe type conversion

func getString(obj map[string]interface{}, key string) string {
	val, ok := obj[key]
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func getStringPtr(obj map[string]interface{}, key string) *string {
	val, ok := obj[key]
	if !ok {
		return nil
	}
	if str, ok := val.(string); ok && str != "" {
		return &str
	}
	return nil
}

func getStringArray(obj map[string]interface{}, key string) []string {
	val, ok := obj[key]
	if !ok {
		return nil
	}

	var result []string
	switch v := val.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
	case string:
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}
