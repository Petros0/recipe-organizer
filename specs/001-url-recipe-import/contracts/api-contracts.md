# API Contracts: URL-Based Recipe Import

**Feature**: 001-url-recipe-import  
**Date**: 2025-12-17

## Overview

This feature uses Appwrite's built-in APIs rather than custom REST endpoints. The contracts below document the Appwrite SDK usage patterns.

---

## 1. Create Recipe Request

**Purpose**: Submit a URL to begin recipe extraction

### Appwrite Function Call

```dart
// Function: recipe-request
Future<Execution> createRecipeRequest(String url) async {
  return await functions.createExecution(
    functionId: 'recipe-request',
    body: jsonEncode({'url': url}),
    xasync: true,  // Asynchronous execution
    method: ExecutionMethod.POST,
  );
}
```

### Request Body

```json
{
  "url": "https://example.com/recipe/chocolate-cake"
}
```

### Response (Execution)

```json
{
  "documentId": "abc123...",
  "status": "REQUESTED",
  "url": "https://example.com/recipe/chocolate-cake"
}
```

### Error Response

```json
{
  "error": "Invalid URL format"
}
```

---

## 2. Subscribe to Request Status

**Purpose**: Real-time updates on extraction progress

### Appwrite Realtime Subscription

```dart
final subscription = realtime.subscribe([
  'databases.6930a343001607ad7cbd.collections.6930a34300165ad1d129.documents.$documentId'
]);

subscription.stream.listen((RealtimeMessage event) {
  if (event.events.contains('databases.*.collections.*.documents.*.update')) {
    final status = event.payload['status'] as String;
    // Handle status change: REQUESTED → IN_PROGRESS → COMPLETED/FAILED
  }
});
```

### Event Payload

```json
{
  "$id": "abc123...",
  "url": "https://example.com/recipe/chocolate-cake",
  "status": "IN_PROGRESS",
  "$createdAt": "2025-12-17T12:00:00.000Z",
  "$updatedAt": "2025-12-17T12:00:05.000Z"
}
```

### Status Values

| Status | Description |
|--------|-------------|
| `REQUESTED` | Initial state, awaiting processor |
| `IN_PROGRESS` | Processor fetching and parsing |
| `COMPLETED` | Recipe successfully extracted |
| `FAILED` | Extraction failed |

---

## 3. Fetch Recipe by Request ID

**Purpose**: Get extracted recipe after COMPLETED status

### Appwrite Database Query

```dart
Future<Recipe?> getRecipeByRequestId(String requestId) async {
  final response = await databases.listDocuments(
    databaseId: '6930a343001607ad7cbd',
    collectionId: 'recipe',
    queries: [
      Query.equal('fk_recipe_request', requestId),
      Query.limit(1),
    ],
  );
  
  if (response.documents.isEmpty) return null;
  return Recipe.fromDocument(response.documents.first);
}
```

### Response Document

```json
{
  "$id": "recipe_xyz...",
  "$createdAt": "2025-12-17T12:00:10.000Z",
  "$updatedAt": "2025-12-17T12:00:10.000Z",
  "name": "Chocolate Cake",
  "description": "A rich, moist chocolate cake...",
  "image": ["https://example.com/images/cake.jpg"],
  "prep_time": "PT20M",
  "cook_time": "PT35M",
  "total_time": "PT55M",
  "recipe_yield": ["12 servings"],
  "ingredients": [
    "2 cups all-purpose flour",
    "1 3/4 cups sugar",
    "3/4 cup cocoa powder"
  ],
  "instructions": [
    "Preheat oven to 350°F",
    "Mix dry ingredients",
    "Add wet ingredients and combine"
  ],
  "author_name": "Chef Example",
  "author_url": "https://example.com/chef",
  "fk_recipe_request": "abc123..."
}
```

---

## 4. List User Recipes

**Purpose**: Fetch all saved recipes for home grid

### Appwrite Database Query

```dart
Future<List<Recipe>> listRecipes({int limit = 25, int offset = 0}) async {
  final response = await databases.listDocuments(
    databaseId: '6930a343001607ad7cbd',
    collectionId: 'recipe',
    queries: [
      Query.limit(limit),
      Query.offset(offset),
      Query.orderDesc('\$createdAt'),
    ],
  );
  
  return response.documents
      .map((doc) => Recipe.fromDocument(doc))
      .toList();
}
```

---

## Appwrite Configuration

| Resource | ID |
|----------|-----|
| Database | `6930a343001607ad7cbd` |
| recipe_request Collection | `6930a34300165ad1d129` |
| recipe Collection | `recipe` |
| recipe-request Function | `recipe-request` |
| recipe-request-processor Function | `recipe-request-processor` |

---

## Error Codes

| Code | Description | User Message |
|------|-------------|--------------|
| `invalid_url` | URL format validation failed | "Please enter a valid URL" |
| `fetch_failed` | Could not fetch page | "Unable to access this website" |
| `no_recipe_found` | No recipe data on page | "No recipe found on this page" |
| `parse_error` | Failed to parse recipe data | "Could not extract recipe details" |
| `timeout` | Request timed out | "Request timed out. Please try again" |
