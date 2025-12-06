# recipe-request-processor

A Go-based Appwrite Cloud Function that extracts structured recipe data from any website using schema.org JSON-LD parsing with Firecrawl fallback for bot-protected sites.

## 🧰 Usage

### GET /ping

Health check endpoint.

**Response**

Sample `200` Response:

```text
Pong
```

### GET, POST /

Extract recipe data from a URL.

**Input (Query Parameter)**

```
?url=https://example.com/recipe
```

**Input (JSON Body)**

```json
{
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
| 400    | Missing or invalid URL       |
| 404    | No Recipe found on page      |
| 500    | Failed to fetch/parse recipe |

## ⚙️ Configuration

| Setting           | Value     |
| ----------------- | --------- |
| Runtime           | Go (1.23) |
| Entrypoint        | `main.go` |
| Permissions       | `any`     |
| Timeout (Seconds) | 15        |

## 🔒 Environment Variables

| Variable            | Required | Description                      |
| ------------------- | -------- | -------------------------------- |
| `FIRECRAWL_API_KEY` | Yes      | Firecrawl API key for fallback   |
| `APPWRITE_API_KEY`  | Yes      | Appwrite API key for database    |
| `APPWRITE_ENDPOINT` | No       | Custom Appwrite endpoint         |
| `APPWRITE_PROJECT_ID` | No     | Custom Appwrite project ID       |

## 🏗️ Architecture

The function uses a strategy pattern with automatic fallback:

1. **HTTP Client Strategy** - Fast, free HTTP requests with browser headers (~80% success rate)
2. **Firecrawl HTML Strategy** - Handles bot-protected sites, parses JSON-LD
3. **Firecrawl LLM Extraction** - AI-based extraction for sites without structured data

### Request Flow

```
Request → URL Validation → Create Request Record (REQUESTED)
    ↓
Update Status (IN_PROGRESS)
    ↓
HTTP Client → JSON-LD Parser
    ↓ (403/429 or no JSON-LD)
Firecrawl (HTML) → JSON-LD Parser
    ↓ (no JSON-LD found)
Firecrawl (LLM Extract) → Recipe Schema
    ↓
Update Status (COMPLETED/FAILED) → Response
```

## 🧪 Testing

```bash
# Run unit tests only
go test -short -v ./...

# Run all tests including integration
go test -v ./...

# Run specific test
go test -v -run TestExtractRecipeFromHTML ./...
```

## 📦 Dependencies

- `github.com/open-runtimes/types-for-go/v4` - Appwrite runtime types
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/appwrite/sdk-for-go` - Appwrite SDK
- `github.com/mendableai/firecrawl-go/v2` - Firecrawl web scraping API
