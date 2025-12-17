# Data Model: URL-Based Recipe Import

**Feature**: 001-url-recipe-import  
**Date**: 2025-12-17

## Entities

### RecipeRequest

Represents a pending recipe import request. Created when user submits a URL.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | `String` | Yes | Appwrite document ID (`$id`) |
| `url` | `String` | Yes | Source URL submitted by user |
| `status` | `RecipeRequestStatus` | Yes | Current processing status |
| `createdAt` | `DateTime` | Yes | When request was created (`$createdAt`) |
| `updatedAt` | `DateTime` | Yes | Last status update (`$updatedAt`) |

**Status Enum** (`RecipeRequestStatus`):
- `requested` - Initial state after document creation
- `inProgress` - Backend function is processing
- `completed` - Recipe successfully extracted
- `failed` - Extraction failed

**Appwrite Collection**: `recipe_request` (ID: `6930a34300165ad1d129`)

---

### Recipe

Represents a fully extracted recipe stored in the database.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | `String` | Yes | Appwrite document ID (`$id`) |
| `name` | `String` | Yes | Recipe title |
| `description` | `String?` | No | Recipe description |
| `images` | `List<String>` | No | Image URLs |
| `prepTime` | `String?` | No | ISO 8601 duration |
| `cookTime` | `String?` | No | ISO 8601 duration |
| `totalTime` | `String?` | No | ISO 8601 duration |
| `recipeYield` | `List<String>?` | No | Servings/yield |
| `ingredients` | `List<String>` | No | Ingredient list |
| `instructions` | `List<String>` | No | Cooking steps |
| `authorName` | `String?` | No | Recipe author |
| `authorUrl` | `String?` | No | Author's URL |
| `recipeCategory` | `List<String>?` | No | Categories (e.g., "Dinner") |
| `recipeCuisine` | `List<String>?` | No | Cuisines (e.g., "Italian") |
| `keywords` | `String?` | No | Comma-separated keywords |
| `datePublished` | `String?` | No | Original publish date |
| `dateModified` | `String?` | No | Original modification date |
| `nutrition` | `NutritionInfo?` | No | Nutritional information |
| `sourceUrl` | `String?` | No | Original recipe URL (from request) |
| `recipeRequestId` | `String?` | No | FK to recipe_request |

**Appwrite Collection**: `recipe` (ID: `recipe`)

---

### NutritionInfo

Embedded value object for nutritional data.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `calories` | `String?` | No | e.g., "250 calories" |
| `fatContent` | `String?` | No | e.g., "12 g" |
| `saturatedFatContent` | `String?` | No | |
| `cholesterolContent` | `String?` | No | |
| `sodiumContent` | `String?` | No | |
| `carbohydrateContent` | `String?` | No | |
| `fiberContent` | `String?` | No | |
| `sugarContent` | `String?` | No | |
| `proteinContent` | `String?` | No | |

---

## Relationships

```
┌─────────────────┐         ┌─────────────────┐
│  RecipeRequest  │ 1 ──── 1│     Recipe      │
│                 │         │                 │
│  id             │         │  id             │
│  url ───────────┼────────►│  sourceUrl      │
│  status         │         │  recipeRequestId│
└─────────────────┘         └─────────────────┘
```

- **RecipeRequest → Recipe**: One-to-one relationship
- Recipe is created when RecipeRequest status becomes `COMPLETED`
- Recipe's `recipeRequestId` links back to the originating request

---

## State Transitions

### RecipeRequest Status Flow

```
                    ┌──────────────┐
User submits URL ──►│  REQUESTED   │
                    └──────┬───────┘
                           │ Function triggered
                           ▼
                    ┌──────────────┐
                    │ IN_PROGRESS  │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │                         │
              ▼                         ▼
       ┌──────────────┐         ┌──────────────┐
       │  COMPLETED   │         │    FAILED    │
       └──────────────┘         └──────────────┘
              │
              ▼
       Recipe document created
```

---

## Validation Rules

### RecipeRequest
- `url` must be a valid URL format (validated by Appwrite)
- `status` must be one of the enum values

### Recipe
- `name` is required and max 512 characters
- `description` max 10,000 characters
- `ingredients` items max 1,024 characters each
- `instructions` items max 4,096 characters each

---

## Appwrite Document Mapping

### RecipeRequest fromDocument

```dart
factory RecipeRequest.fromDocument(Document doc) {
  return RecipeRequest(
    id: doc.$id,
    url: doc.data['url'] as String,
    status: RecipeRequestStatus.fromString(doc.data['status'] as String),
    createdAt: DateTime.parse(doc.$createdAt),
    updatedAt: DateTime.parse(doc.$updatedAt),
  );
}
```

### Recipe fromDocument

```dart
factory Recipe.fromDocument(Document doc) {
  return Recipe(
    id: doc.$id,
    name: doc.data['name'] as String,
    description: doc.data['description'] as String?,
    images: (doc.data['image'] as List<dynamic>?)?.cast<String>() ?? [],
    prepTime: doc.data['prep_time'] as String?,
    cookTime: doc.data['cook_time'] as String?,
    totalTime: doc.data['total_time'] as String?,
    recipeYield: (doc.data['recipe_yield'] as List<dynamic>?)?.cast<String>(),
    ingredients: (doc.data['ingredients'] as List<dynamic>?)?.cast<String>() ?? [],
    instructions: (doc.data['instructions'] as List<dynamic>?)?.cast<String>() ?? [],
    authorName: doc.data['author_name'] as String?,
    authorUrl: doc.data['author_url'] as String?,
    recipeCategory: (doc.data['recipe_category'] as List<dynamic>?)?.cast<String>(),
    recipeCuisine: (doc.data['recipe_cuisine'] as List<dynamic>?)?.cast<String>(),
    keywords: doc.data['keywords'] as String?,
    datePublished: doc.data['date_published'] as String?,
    dateModified: doc.data['date_modified'] as String?,
    nutrition: _parseNutrition(doc.data),
    recipeRequestId: doc.data['fk_recipe_request'] as String?,
  );
}
```
