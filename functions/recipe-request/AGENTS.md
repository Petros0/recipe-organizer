# AGENTS.md - recipe-request

This document provides guidance for AI agents working on this Appwrite Cloud Function.

## Overview

**recipe-request** is a Go-based Appwrite Cloud Function that creates recipe request records for async processing. It accepts a URL, validates it, creates a record in Appwrite with `REQUESTED` status, and returns the document ID immediately.

The actual recipe fetching is handled by the `recipe-request-processor` function, which is triggered by a database event.

## Tech Stack

| Layer         | Technology                         |
| ------------- | ---------------------------------- |
| Language      | Go 1.23+                           |
| Runtime       | Appwrite Functions (Go 1.23)       |
| Types         | open-runtimes/types-for-go v4      |
| Backend       | [Appwrite](https://appwrite.io) (BaaS) |

## Project Structure

```
functions/recipe-request/
├── main.go               # Handler entry point
├── types.go              # Request/Response types
├── appwrite_client.go    # Appwrite database client
├── main_test.go          # Integration tests
├── appwrite_client_test.go # Client tests
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── build/                # Build artifacts (gitignored)
├── README.md             # API documentation
└── AGENTS.md             # This file
```

## Data Models

### Request/Response Types

```go
// RequestBody represents the JSON request body
type RequestBody struct {
    URL string `json:"url"`
}

// SuccessResponse represents a successful request creation response
type SuccessResponse struct {
    DocumentID string `json:"documentId"`
    Status     string `json:"status"`
    URL        string `json:"url"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error string `json:"error"`
}
```

## API Endpoints

### GET /ping

Health check endpoint.

**Response:** `200 OK` - `"Pong"`

### GET, POST /

Create a recipe request record.

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
  "documentId": "abc123xyz",
  "status": "REQUESTED",
  "url": "https://example.com/recipe"
}
```

**Error Responses:**

| Status | Condition                |
| ------ | ------------------------ |
| 400    | Missing or invalid URL   |
| 500    | Failed to create record  |

## Architecture

### Two-Step Async Workflow

This function is part of a two-step async workflow:

1. **recipe-request** (this function)
   - Validates URL format
   - Creates request record with REQUESTED status
   - Returns document ID immediately

2. **recipe-request-processor**
   - Triggered by database create event
   - Fetches recipe using HTTP/Firecrawl
   - Updates status to IN_PROGRESS → COMPLETED/FAILED

### Request Flow

```
Request → URL Validation → Create Request Record (REQUESTED) → Return Document ID
```

### Key Functions

| Function                       | Purpose                                      |
| ------------------------------ | -------------------------------------------- |
| `Main`                         | Entry point, request handling                |
| `NewRecipeRequestClient`       | Create Appwrite client for database ops      |
| `CreateRequest`                | Create request record with REQUESTED status  |

## Testing

### Test Structure

| Test                                    | Type        | Description                |
| --------------------------------------- | ----------- | -------------------------- |
| `TestStatusConstants`                   | Unit        | Status constant validation |
| `TestDefaultConstants`                  | Unit        | Appwrite config constants  |
| `TestMockRecipeRequestStore_CreateRequest` | Unit     | Mock store testing         |
| `TestNewRecipeRequestClient`            | Unit        | Client creation            |
| `TestCreateRequest_Integration`         | Integration | Live request creation      |

### Running Tests

```bash
# Run unit tests only
go test -short -v ./...

# Run all tests including integration
go test -v ./...

# Run specific test
go test -v -run TestCreateRequest_Integration ./...
```

### Test Conventions

- Use table-driven tests with `tests := []struct{}`
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

| Variable             | Required | Description                    |
| -------------------- | -------- | ------------------------------ |
| `APPWRITE_API_KEY`   | Yes      | Appwrite API key for database  |
| `APPWRITE_ENDPOINT`  | No       | Custom Appwrite endpoint       |
| `APPWRITE_PROJECT_ID`| No       | Custom Appwrite project ID     |

## Dependencies

```go
require (
    github.com/open-runtimes/types-for-go/v4 v4.0.6  // Appwrite runtime types
    github.com/appwrite/sdk-for-go v0.15.0           // Appwrite SDK
)
```

## Important Notes for Agents

1. **Single Responsibility** - This function only creates request records, not fetch recipes
2. **URL Validation** - Validates URL format before creating record
3. **Status** - Always creates records with `REQUESTED` status
4. **Async Processing** - Recipe fetching is handled by `recipe-request-processor`
5. **Error Handling** - Return descriptive JSON errors with appropriate HTTP status codes
6. **Environment Variables** - Requires `APPWRITE_API_KEY` for database operations
