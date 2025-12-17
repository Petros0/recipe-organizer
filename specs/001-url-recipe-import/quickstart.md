# Quickstart: URL-Based Recipe Import

**Feature**: 001-url-recipe-import  
**Date**: 2025-12-17

## Overview

This guide provides the essential information for implementing the URL-based recipe import feature in Flutter.

## Architecture Summary

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Flutter App                                 │
├─────────────────────────────────────────────────────────────────────┤
│  View Layer                                                          │
│  ┌──────────────┐  ┌────────────────────┐  ┌────────────────────┐   │
│  │  HomePage    │  │ RecipePreviewPage  │  │ ImportRecipeDialog │   │
│  └──────────────┘  └────────────────────┘  └────────────────────┘   │
│         │                   │                        │               │
├─────────┼───────────────────┼────────────────────────┼───────────────┤
│  State  │                   │                        │               │
│         ▼                   ▼                        ▼               │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                    HomeController (Signals)                   │   │
│  │  recipes: Signal<List<Recipe>>                               │   │
│  │  activeRequest: Signal<RecipeRequest?>                       │   │
│  │  importState: Computed<ImportState>                          │   │
│  └──────────────────────────────────────────────────────────────┘   │
│         │                                                            │
├─────────┼────────────────────────────────────────────────────────────┤
│  Service│                                                            │
│         ▼                                                            │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                  RecipeImportService                          │   │
│  │  - importFromUrl(url) → Future<RecipeRequest>                │   │
│  │  - subscribeToRequest(id) → Stream<RecipeRequest>            │   │
│  └──────────────────────────────────────────────────────────────┘   │
│         │                                                            │
├─────────┼────────────────────────────────────────────────────────────┤
│  Data   │                                                            │
│         ▼                                                            │
│  ┌────────────────────────┐  ┌─────────────────────────────────┐    │
│  │ RecipeRequestRepository│  │      RecipeRepository           │    │
│  │ - createRequest()      │  │ - getRecipeByRequestId()        │    │
│  │ - subscribeToUpdates() │  │ - listRecipes()                 │    │
│  └────────────────────────┘  └─────────────────────────────────┘    │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Appwrite Backend                              │
│  ┌─────────────────┐   ┌────────────────┐   ┌──────────────────┐   │
│  │ Functions       │   │ Database       │   │ Realtime         │   │
│  │ - recipe-request│   │ - recipe_request   │ - Subscriptions  │   │
│  │ - processor     │   │ - recipe       │   │                  │   │
│  └─────────────────┘   └────────────────┘   └──────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Key Files to Create/Modify

### New Files

| File | Purpose |
|------|---------|
| `lib/core/appwrite/appwrite_constants.dart` | Database/collection IDs |
| `lib/core/appwrite/realtime_service.dart` | Realtime subscription wrapper |
| `lib/features/home/model/recipe_request.dart` | RecipeRequest model + enum |
| `lib/features/home/data/recipe_repository.dart` | Recipe CRUD operations |
| `lib/features/home/service/recipe_import_service.dart` | Import business logic |
| `lib/features/home/view/recipe_preview_page.dart` | Full-detail preview |
| `lib/features/home/view/widgets/recipe_skeleton_card.dart` | Shimmer loading |
| `lib/features/home/view/widgets/recipe_error_card.dart` | Inline error state |

### Modified Files

| File | Changes |
|------|---------|
| `lib/core/di.dart` | Register new services/repositories |
| `lib/features/home/model/recipe.dart` | Add fromDocument, nutrition fields |
| `lib/features/home/data/recipe_request_repository.dart` | Complete implementation |
| `lib/features/home/state/home_controller.dart` | Connect to real data |
| `lib/features/home/view/home_page.dart` | Add skeleton/error states |
| `pubspec.yaml` | Add shimmer package |

## Import Flow Implementation

### Step 1: User Submits URL

```dart
// In HomeController
Future<void> importRecipe(String url) async {
  _importState.value = ImportState.submitting;
  
  try {
    final request = await _importService.importFromUrl(url);
    _activeRequest.value = request;
    _importState.value = ImportState.extracting;
    
    // Start listening for updates
    _subscribeToRequest(request.id);
  } catch (e) {
    _error.value = e.toString();
    _importState.value = ImportState.error;
  }
}
```

### Step 2: Subscribe to Status Updates

```dart
// In RecipeImportService
Stream<RecipeRequest> subscribeToRequest(String requestId) {
  return _realtimeService.subscribe(
    channel: 'databases.$databaseId.collections.$collectionId.documents.$requestId',
    transform: (payload) => RecipeRequest.fromMap(payload),
  );
}
```

### Step 3: Handle Status Changes

```dart
// In HomeController
void _subscribeToRequest(String requestId) {
  _requestSubscription?.cancel();
  _requestSubscription = _importService
      .subscribeToRequest(requestId)
      .listen((request) {
        _activeRequest.value = request;
        
        switch (request.status) {
          case RecipeRequestStatus.completed:
            _fetchExtractedRecipe(requestId);
          case RecipeRequestStatus.failed:
            _importState.value = ImportState.error;
          default:
            break;
        }
      });
}
```

### Step 4: Fetch Completed Recipe

```dart
Future<void> _fetchExtractedRecipe(String requestId) async {
  final recipe = await _recipeRepository.getRecipeByRequestId(requestId);
  if (recipe != null) {
    _previewRecipe.value = recipe;
    _importState.value = ImportState.preview;
  }
}
```

## UI State Machine

```dart
enum ImportState {
  idle,        // No active import
  submitting,  // Creating request document
  extracting,  // Showing skeleton, waiting for completion
  preview,     // Recipe extracted, showing preview
  error,       // Import failed
  saving,      // Saving to collection
}
```

## Widget Mapping

| State | Widget |
|-------|--------|
| `idle` | `RecipeGrid` or `EmptyState` |
| `submitting` | `ImportRecipeDialog` with spinner |
| `extracting` | `RecipeSkeletonCard` (shimmer) |
| `preview` | `RecipePreviewPage` |
| `error` | `RecipeErrorCard` |
| `saving` | `RecipePreviewPage` with spinner |

## Testing Strategy

| Layer | Test Type | Key Tests |
|-------|-----------|-----------|
| Model | Unit | Recipe.fromDocument parsing |
| Repository | Unit (mocked) | Appwrite SDK calls |
| Service | Unit (mocked) | Import flow orchestration |
| Controller | Unit (mocked) | State transitions |
| Widgets | Widget | Skeleton render, error display |
| Integration | Integration | Full import flow |

## Dependencies

```yaml
# pubspec.yaml
dependencies:
  shimmer: ^3.0.0  # Add for skeleton effect
```

## Localization Keys to Add

```json
{
  "extractingRecipe": "Extracting recipe...",
  "importFailed": "Import failed",
  "importFailedRetry": "Tap to retry",
  "recipePreviewTitle": "Recipe Preview",
  "saveToCollection": "Save to Collection",
  "prepTime": "Prep",
  "cookTime": "Cook",
  "servings": "Servings",
  "ingredients": "Ingredients",
  "instructions": "Instructions",
  "source": "Source"
}
```
