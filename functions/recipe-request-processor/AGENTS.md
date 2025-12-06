# AGENTS.md - recipe-request-processor

This document provides guidance for AI agents working on this Appwrite Cloud Function.

## Overview

**recipe-request-processor** is a Go-based Appwrite Cloud Function that processes recipe requests by extracting structured recipe data from websites. It uses a multi-layered approach:
1. First tries parsing schema.org Recipe JSON-LD markup via HTTP
2. Falls back to Firecrawl API for bot-protected sites
3. Uses Firecrawl's LLM extraction for sites without JSON-LD structured data

This function is triggered by a database event when `recipe-request` creates a new request record.

## Tech Stack

| Layer         | Technology                         |
| ------------- | ---------------------------------- |
| Language      | Go 1.23+                           |
| Runtime       | Appwrite Functions (Go 1.23)       |
| HTML Parsing  | [goquery](https://github.com/PuerkitoBio/goquery) |
| Web Scraping  | [Firecrawl](https://firecrawl.dev) |
| Types         | open-runtimes/types-for-go v4      |
| Backend       | [Appwrite](https://appwrite.io) (BaaS) |

## Project Structure

```
functions/recipe-request-processor/
├── main.go               # Handler entry point
├── strategy.go           # Strategy pattern executor
├── http_client.go        # HTTP client strategy
├── firecrawl_client.go   # Firecrawl API strategy (fallback + LLM extraction)
├── parser.go             # JSON-LD extraction logic
├── field_parsers.go      # Field-specific parsers
├── types.go              # Data models
├── response.go           # Response formatting
├── utils.go              # Utility functions
├── appwrite_client.go    # Appwrite database client for status tracking
├── main_test.go          # Integration tests
├── parser_test.go        # JSON-LD parsing tests
├── field_parsers_test.go # Field parser unit tests
├── response_test.go      # Response transformation tests
├── appwrite_client_test.go # Appwrite client tests
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── build/                # Build artifacts (gitignored)
├── README.md             # API documentation
└── AGENTS.md             # This file
```

## Data Models

### Core Types

```go
// ProcessorRequestBody represents the request body for the processor
type ProcessorRequestBody struct {
    DocumentID string `json:"documentId"`
    URL        string `json:"url"`
}

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
```

## API Endpoints

### GET /ping

Health check endpoint.

**Response:** `200 OK` - `"Pong"`

### POST /

Process a recipe request (triggered by database event).

**Input (JSON Body):**
```json
{
  "documentId": "abc123xyz",
  "url": "https://example.com/recipe"
}
```

**Success Response:** `200 OK`
```json
{
  "url": "https://example.com/recipe",
  "recipe": {
    "name": "Recipe Name",
    "description": "A delicious recipe",
    "image": "https://example.com/image.jpg",
    "prepTime": "PT15M",
    "cookTime": "PT30M",
    "totalTime": "PT45M",
    "author": "Chef Name"
  },
  "instructions": ["Step 1", "Step 2"],
  "ingredients": ["1 cup flour", "2 eggs"]
}
```

**Error Responses:**

| Status | Condition                     |
| ------ | ----------------------------- |
| 400    | Missing documentId or url     |
| 404    | No Recipe found               |
| 500    | Failed to fetch/parse recipe  |

## Architecture

### Two-Step Async Workflow

This function is the second step in a two-step async workflow:

1. **recipe-request**
   - Validates URL format
   - Creates request record with REQUESTED status
   - Returns document ID immediately

2. **recipe-request-processor** (this function)
   - Triggered by database create event
   - Updates status to IN_PROGRESS
   - Fetches recipe using HTTP/Firecrawl
   - Updates status to COMPLETED/FAILED

### Request Flow

```
Event Trigger (documentId, url) → Update Status (IN_PROGRESS)
    ↓
HTTP Client → JSON-LD Parser
    ↓ (403/429 or no JSON-LD)
Firecrawl (HTML) → JSON-LD Parser
    ↓ (no JSON-LD found)
Firecrawl (LLM Extract) → Recipe Schema
    ↓
Update Status (COMPLETED/FAILED) → Response
```

### Strategy Pattern

The function uses a strategy pattern with automatic fallback:

1. **HTTPClientStrategy** - Fast, free HTTP requests with browser headers
2. **FirecrawlStrategy** - Managed API with bot protection bypass and LLM extraction

### Key Functions

| Function                       | Purpose                                      |
| ------------------------------ | -------------------------------------------- |
| `Main`                         | Entry point, request handling                |
| `StrategyExecutor.Execute`     | Executes strategies with fallback logic      |
| `HTTPClientStrategy.Fetch`     | HTTP client with browser headers             |
| `FirecrawlStrategy.Fetch`      | Firecrawl API (HTML + LLM fallback)          |
| `extractRecipeFromHTML`        | Parse HTML, find JSON-LD script tags         |
| `extractRecipeFromJSONLD`      | Handle various JSON-LD formats               |
| `extractRecipeFromObject`      | Parse single Recipe object                   |
| `parseImage`                   | Handle string/array/ImageObject formats      |
| `parseInstructions`            | Handle HowToStep/HowToSection formats        |
| `parseAuthor`                  | Handle string/Person/array formats           |
| `parseNutrition`               | Parse NutritionInformation                   |
| `fetchWithLLMExtraction`       | AI-based recipe extraction for sites without JSON-LD |
| `NewRecipeRequestClient`       | Create Appwrite client for status tracking   |
| `UpdateStatus`                 | Update request status in Appwrite            |
| `toRecipeResponse`             | Transform Recipe to API response format      |

### Extraction Strategies

The function uses a cost-optimized hybrid approach:

1. **HTTP Client** - Free, works for ~80% of recipe sites
2. **Firecrawl HTML** - Handles bot protection, parses JSON-LD (1 credit)
3. **Firecrawl LLM Extract** - AI extraction for sites without structured data (higher cost)

## Testing

### Test Structure

| Test                                | Type        | Description                          |
| ----------------------------------- | ----------- | ------------------------------------ |
| `TestExtractRecipeFromHTML`         | Unit        | HTML JSON-LD extraction              |
| `TestExtractRecipeFromJSONLD`       | Unit        | JSON-LD format handling              |
| `TestParseImage`                    | Unit        | Image field parsing                  |
| `TestParseInstructions`             | Unit        | Instruction parsing                  |
| `TestParseAuthor`                   | Unit        | Author field parsing                 |
| `TestParseNutrition`                | Unit        | Nutrition parsing                    |
| `TestToRecipeResponse`              | Unit        | Response transformation              |
| `TestStatusConstants`               | Unit        | Status constant validation           |
| `TestFetchRecipe_Integration`       | Integration | Live recipe fetching (skipped short) |
| `TestFirecrawlStrategy_HTMLWithJSONLD` | Integration | Firecrawl HTML parsing            |
| `TestFirecrawlStrategy_LLMExtraction`  | Integration | Firecrawl LLM extraction          |

### Running Tests

```bash
# Run unit tests only
go test -short -v ./...

# Run all tests including integration
go test -v ./...

# Run specific test
go test -v -run TestExtractRecipeFromHTML ./...
```

### Test Conventions

- Use table-driven tests with `tests := []struct{}`
- Name test cases descriptively: `"valid URL - akispetretzikis.com recipe"`
- Skip integration tests in short mode: `testing.Short()`
- Log informative output for debugging: `t.Logf()`

## Common Commands

```bash
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

## Environment Variables

| Variable           | Required | Description                        |
| ------------------ | -------- | ---------------------------------- |
| `FIRECRAWL_API_KEY`| Yes      | Firecrawl API key for fallback     |
| `APPWRITE_API_KEY` | Yes      | Appwrite API key for database ops  |
| `APPWRITE_ENDPOINT`| No       | Custom Appwrite endpoint           |
| `APPWRITE_PROJECT_ID` | No    | Custom Appwrite project ID         |

## Dependencies

```go
require (
    github.com/open-runtimes/types-for-go/v4 v4.0.6  // Appwrite runtime types
    github.com/PuerkitoBio/goquery v1.10.2           // HTML parsing
    github.com/appwrite/sdk-for-go v0.15.0           // Appwrite SDK
    github.com/mendableai/firecrawl-go/v2 v2.4.0     // Firecrawl web scraping API
)
```

## Important Notes for Agents

1. **Event-Triggered** - This function is triggered by database events, not direct API calls
2. **Input Format** - Expects `documentId` and `url` in request body
3. **No URL Validation** - URL validation is done by `recipe-request`, not here
4. **No Request Creation** - Request records are created by `recipe-request`, not here
5. **Status Tracking** - Updates status: IN_PROGRESS → COMPLETED/FAILED
6. **Required Fields** - Recipe must have `name` and `image` to be valid (relaxed for LLM extraction)
7. **JSON-LD Formats** - Handle single object, arrays, and `@graph` containers
8. **Type Variations** - Check both `"Recipe"` and URLs like `"https://schema.org/Recipe"`
9. **Image Formats** - Can be string, array of strings, or ImageObject with `url` property
10. **Instructions Formats** - Handle both HowToStep and nested HowToSection
11. **Author Formats** - Can be string, Person object, or array (takes first)
12. **Bot Protection** - HTTP 403/429 triggers automatic Firecrawl fallback
13. **LLM Extraction** - Used when JSON-LD is not found on the page
14. **Timeout** - HTTP client: 30s, Firecrawl: API default
15. **Pointer Types** - Optional fields use `*string` to distinguish empty from missing
16. **Error Handling** - Return descriptive JSON errors with appropriate HTTP status codes
17. **Cost Optimization** - HTTP client is tried first (free), Firecrawl only when needed

## Limitations

- Only extracts first Recipe found on page
- JSON-LD parsing doesn't support microdata or RDFa
- LLM extraction quality depends on page content structure
- Some recipe fields may be empty if not provided by source site
- Requires `FIRECRAWL_API_KEY` environment variable for fallback
- Requires `APPWRITE_API_KEY` for status tracking
