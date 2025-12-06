# recipe-request-processor

A Go-based Appwrite Cloud Function that processes recipe requests by fetching and extracting structured recipe data from websites.

## Overview

This function is triggered by a database event when a recipe request is created. It receives the document ID and URL, fetches the recipe using schema.org JSON-LD parsing with Firecrawl fallback, and updates the request status.

## Usage

### GET /ping

Health check endpoint.

**Response**

Sample `200` Response:

```text
Pong
```

### POST /

Process a recipe request (typically triggered by database event).

**Input (JSON Body)**

```json
{
  "documentId": "abc123xyz",
  "url": "https://example.com/recipe"
}
```

**Response**

Sample `200` Response:

```json
{
  "url": "https://example.com/recipe",
  "recipe": {
    "name": "Chocolate Cake",
    "description": "A rich and moist chocolate cake",
    "image": "https://example.com/cake.jpg",
    "prepTime": "PT20M",
    "cookTime": "PT35M",
    "totalTime": "PT55M",
    "author": "Chef John"
  },
  "instructions": [
    "Preheat oven to 350°F",
    "Mix dry ingredients",
    "Add wet ingredients",
    "Bake for 35 minutes"
  ],
  "ingredients": [
    "2 cups flour",
    "1 cup sugar",
    "3/4 cup cocoa powder",
    "2 eggs",
    "1 cup milk"
  ]
}
```

**Error Responses**

| Status | Condition                    |
| ------ | ---------------------------- |
| 400    | Missing documentId or url    |
| 404    | No Recipe found on page      |
| 500    | Failed to fetch/parse recipe |

## Configuration

| Setting           | Value     |
| ----------------- | --------- |
| Runtime           | Go (1.23) |
| Entrypoint        | `main.go` |
| Permissions       | `any`     |
| Timeout (Seconds) | 15        |

## Environment Variables

| Variable            | Required | Description                      |
| ------------------- | -------- | -------------------------------- |
| `FIRECRAWL_API_KEY` | Yes      | Firecrawl API key for fallback   |
| `APPWRITE_API_KEY`  | Yes      | Appwrite API key for database    |
| `APPWRITE_ENDPOINT` | No       | Custom Appwrite endpoint         |
| `APPWRITE_PROJECT_ID` | No     | Custom Appwrite project ID       |

## Architecture

This function is the second step in a two-step async workflow:

1. **recipe-request** - Validates URL, creates request record, returns document ID immediately
2. **recipe-request-processor** (this function) - Fetches recipe and updates status

### Request Flow

```
Event Trigger → Update Status (IN_PROGRESS)
    ↓
HTTP Client → JSON-LD Parser
    ↓ (403/429 or no JSON-LD)
Firecrawl (HTML) → JSON-LD Parser
    ↓ (no JSON-LD found)
Firecrawl (LLM Extract) → Recipe Schema
    ↓
Update Status (COMPLETED/FAILED) → Response
```

### Extraction Strategies

The function uses a cost-optimized hybrid approach:

1. **HTTP Client** - Free, works for ~80% of recipe sites
2. **Firecrawl HTML** - Handles bot protection, parses JSON-LD (1 credit)
3. **Firecrawl LLM Extract** - AI extraction for sites without structured data (higher cost)

## Testing

```bash
# Run unit tests only
go test -short -v ./...

# Run all tests including integration
go test -v ./...

# Run specific test
go test -v -run TestExtractRecipeFromHTML ./...
```

## Dependencies

- `github.com/open-runtimes/types-for-go/v4` - Appwrite runtime types
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/appwrite/sdk-for-go` - Appwrite SDK
- `github.com/mendableai/firecrawl-go/v2` - Firecrawl web scraping API
