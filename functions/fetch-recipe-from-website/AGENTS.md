# AGENTS.md - fetch-recipe-from-website

This document provides guidance for AI agents working on this Appwrite Cloud Function.

## Overview

**fetch-recipe-from-website** is a Go-based Appwrite Cloud Function that extracts structured recipe data from any website by parsing schema.org Recipe JSON-LD markup. It provides automatic fallback to headless browser rendering when standard HTTP requests are blocked by bot protection.

## Tech Stack

| Layer        | Technology                         |
| ------------ | ---------------------------------- |
| Language     | Go 1.23+                           |
| Runtime      | Appwrite Functions (Go 1.23)       |
| HTML Parsing | [goquery](https://github.com/PuerkitoBio/goquery) |
| Headless     | [go-rod](https://go-rod.github.io) |
| Types        | open-runtimes/types-for-go v4      |

## Project Structure

```
functions/fetch-recipe-from-website/
├── main.go          # Handler and all extraction logic
├── main_test.go     # Unit and integration tests
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── build/           # Build artifacts (gitignored)
├── README.md        # API documentation
└── AGENTS.md        # This file
```

## Data Models

### Core Types

```go
// Recipe - Main response type (schema.org/Recipe)
type Recipe struct {
    Name               string              // Required
    Image              []string            // Required
    Author             *Person
    Description        *string
    PrepTime           *string             // ISO 8601 duration
    CookTime           *string             // ISO 8601 duration
    TotalTime          *string             // ISO 8601 duration
    RecipeYield        *string
    RecipeIngredient   []string
    RecipeInstructions []RecipeInstruction
    RecipeCategory     *string
    RecipeCuisine      *string
    Nutrition          *Nutrition
    Keywords           *string
    DatePublished      *string
    DateModified       *string
}

// Person - Author information
type Person struct {
    Type string
    Name string
    URL  string
}

// Nutrition - Nutritional information
type Nutrition struct {
    Calories            *string
    FatContent          *string
    SaturatedFatContent *string
    CholesterolContent  *string
    SodiumContent       *string
    CarbohydrateContent *string
    FiberContent        *string
    SugarContent        *string
    ProteinContent      *string
}

// RecipeInstruction - HowToStep or HowToSection
type RecipeInstruction struct {
    Type            string
    Text            string
    Name            string
    URL             string
    ItemListElement []RecipeInstruction // Nested for HowToSection
}
```

## API Endpoints

### GET /ping

Health check endpoint.

**Response:** `200 OK` - `"Pong"`

### GET, POST /

Extract recipe from URL.

**Input (Query Parameter):**
```
?url=https://example.com/recipe
```

**Input (JSON Body):**
```json
{
  "url": "https://example.com/recipe"
}
```

**Success Response:** `200 OK`
```json
{
  "name": "Recipe Name",
  "image": ["https://example.com/image.jpg"],
  "recipeIngredient": ["1 cup flour", "2 eggs"],
  "recipeInstructions": [{"@type": "HowToStep", "text": "Mix ingredients"}]
}
```

**Error Responses:**

| Status | Condition                     |
| ------ | ----------------------------- |
| 400    | Missing or invalid URL        |
| 404    | No Recipe JSON-LD found       |
| 500    | Failed to fetch/parse recipe  |

## Architecture

### Request Flow

```
Request → URL Validation → HTTP Fetch → Extract JSON-LD → Parse Recipe → Response
                              ↓ (403/429)
                     Headless Browser Fallback
```

### Key Functions

| Function                       | Purpose                                      |
| ------------------------------ | -------------------------------------------- |
| `Main`                         | Entry point, request handling                |
| `fetchRecipeFromURL`           | HTTP client with browser headers             |
| `fetchRecipeFromURLWithBrowser`| Headless browser fallback (go-rod)           |
| `extractRecipeFromHTML`        | Parse HTML, find JSON-LD script tags         |
| `extractRecipeFromJSONLD`      | Handle various JSON-LD formats               |
| `extractRecipeFromObject`      | Parse single Recipe object                   |
| `parseImage`                   | Handle string/array/ImageObject formats      |
| `parseInstructions`            | Handle HowToStep/HowToSection formats        |
| `parseAuthor`                  | Handle string/Person/array formats           |
| `parseNutrition`               | Parse NutritionInformation                   |

### Bot Protection Bypass

The function uses two strategies:

1. **HTTP Client with Browser Headers** - Sets realistic User-Agent, Accept, and Sec-Fetch headers
2. **Headless Browser Fallback** - Uses go-rod with stealth options when HTTP returns 403/429

## Testing

### Test Structure

| Test                               | Type        | Description                          |
| ---------------------------------- | ----------- | ------------------------------------ |
| `TestValidateURL`                  | Unit        | URL format validation                |
| `TestParseURL`                     | Unit        | URL parsing components               |
| `TestFetchRecipeFromURL_Integration`| Integration | Live recipe fetching (skipped short) |

### Running Tests

```bash
# Run unit tests only
go test -short -v ./...

# Run all tests including integration
go test -v ./...

# Run specific test
go test -v -run TestValidateURL ./...
```

### Test Conventions

- Use table-driven tests with `tests := []struct{}`
- Name test cases descriptively: `"valid URL - akispetretzikis.com recipe"`
- Skip integration tests in short mode: `testing.Short()`
- Log informative output for debugging: `t.Logf()`

## Common Commands

```bash
# Run the function locally (requires Appwrite CLI)
appwrite functions createExecution \
  --functionId=<function-id> \
  --body='{"url":"https://example.com/recipe"}'

# Build
go build -o build/handler .

# Test
go test -v ./...

# Format
go fmt ./...

# Lint (if golangci-lint installed)
golangci-lint run

# Update dependencies
go mod tidy
```

## Configuration

| Setting     | Value          |
| ----------- | -------------- |
| Runtime     | Go 1.23        |
| Entrypoint  | `main.go`      |
| Package     | `handler`      |
| Permissions | `any`          |
| Timeout     | 15 seconds     |

## Dependencies

```go
require (
    github.com/open-runtimes/types-for-go/v4 v4.0.6  // Appwrite runtime types
    github.com/PuerkitoBio/goquery v1.11.0           // HTML parsing
    github.com/go-rod/rod v0.116.2                   // Headless browser
)
```

## Important Notes for Agents

1. **Required Fields** - Recipe must have `name` and `image` to be valid
2. **JSON-LD Formats** - Handle single object, arrays, and `@graph` containers
3. **Type Variations** - Check both `"Recipe"` and URLs like `"https://schema.org/Recipe"`
4. **Image Formats** - Can be string, array of strings, or ImageObject with `url` property
5. **Instructions Formats** - Handle both HowToStep and nested HowToSection
6. **Author Formats** - Can be string, Person object, or array (takes first)
7. **Bot Protection** - HTTP 403/429 triggers automatic headless browser fallback
8. **Timeout** - HTTP client: 30s, Browser: 20s, Function total: 15s (increase if needed)
9. **Pointer Types** - Optional fields use `*string` to distinguish empty from missing
10. **Error Handling** - Return descriptive JSON errors with appropriate HTTP status codes

## Limitations

- Only extracts first Recipe found on page
- Requires schema.org JSON-LD (not microdata or RDFa)
- Headless browser may fail on sites with advanced bot protection
- Some recipe fields may be empty if not provided by source site

