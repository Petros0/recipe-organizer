# recipe-request

A Go-based Appwrite Cloud Function that creates recipe request records for async processing.

## Overview

This function accepts a URL and creates a recipe request record in Appwrite with `REQUESTED` status. The actual recipe fetching is handled asynchronously by the `recipe-request-processor` function.

## Usage

### GET /ping

Health check endpoint.

**Response**

Sample `200` Response:

```text
Pong
```

### GET, POST /

Create a recipe request record.

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
  "documentId": "abc123xyz",
  "status": "REQUESTED",
  "url": "https://example.com/recipe"
}
```

**Error Responses**

| Status | Condition          |
| ------ | ------------------ |
| 400    | Missing or invalid URL |
| 500    | Failed to create request record |

## Configuration

| Setting           | Value     |
| ----------------- | --------- |
| Runtime           | Go (1.23) |
| Entrypoint        | `main.go` |
| Permissions       | `any`     |
| Timeout (Seconds) | 15        |

## Environment Variables

| Variable             | Required | Description                    |
| -------------------- | -------- | ------------------------------ |
| `APPWRITE_API_KEY`   | Yes      | Appwrite API key for database  |
| `APPWRITE_ENDPOINT`  | No       | Custom Appwrite endpoint       |
| `APPWRITE_PROJECT_ID`| No       | Custom Appwrite project ID     |

## Architecture

This function is the first step in a two-step async workflow:

1. **recipe-request** (this function) - Validates URL, creates request record, returns document ID immediately
2. **recipe-request-processor** - Triggered by database event, fetches recipe and updates status

### Request Flow

```
Request → URL Validation → Create Request Record (REQUESTED) → Return Document ID
```

## Testing

```bash
# Run unit tests only
go test -short -v ./...

# Run all tests including integration
go test -v ./...
```

## Dependencies

- `github.com/open-runtimes/types-for-go/v4` - Appwrite runtime types
- `github.com/appwrite/sdk-for-go` - Appwrite SDK
