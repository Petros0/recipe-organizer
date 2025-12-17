# Research: URL-Based Recipe Import

**Feature**: 001-url-recipe-import  
**Date**: 2025-12-17

## Research Tasks

### 1. Appwrite Real-time Subscription for Flutter

**Task**: Research best practices for Appwrite real-time subscriptions in Flutter with Signals

**Decision**: Use Appwrite `Realtime.subscribe()` with document channel subscription

**Rationale**:
- Appwrite SDK provides native Flutter support for real-time subscriptions
- Subscribe to `databases.<DB_ID>.collections.<COLLECTION_ID>.documents.<DOC_ID>` for specific document updates
- Integrates cleanly with Signals by updating signal values in subscription callback

**Implementation Pattern**:
```dart
final realtime = Realtime(client);
final subscription = realtime.subscribe([
  'databases.6930a343001607ad7cbd.collections.6930a34300165ad1d129.documents.$documentId'
]);
subscription.stream.listen((response) {
  // Update signal with new status
  final status = response.payload['status'] as String;
  _requestStatus.value = RecipeRequestStatus.fromString(status);
});
```

**Alternatives Considered**:
- Polling: Rejected - inefficient, poor UX, unnecessary load
- WebSockets custom: Rejected - Appwrite handles this natively

---

### 2. Shimmer/Skeleton Loading in forUI

**Task**: Check if forUI provides shimmer or skeleton loading components

**Decision**: Use custom skeleton widget with Flutter shimmer package or simple animation

**Rationale**:
- forUI does not include a built-in shimmer/skeleton component (checked docs)
- `shimmer` package (pub.dev) is lightweight and well-maintained
- Alternative: use `AnimatedContainer` with gradient animation for custom implementation

**Implementation Pattern**:
```dart
Shimmer.fromColors(
  baseColor: theme.colors.secondary,
  highlightColor: theme.colors.secondaryForeground,
  child: Container(
    decoration: BoxDecoration(
      color: theme.colors.secondary,
      borderRadius: BorderRadius.circular(8),
    ),
  ),
)
```

**Alternatives Considered**:
- forUI built-in: Not available
- Custom gradient animation: Viable but more code; shimmer package preferred

---

### 3. Recipe Model Alignment with Appwrite Schema

**Task**: Research how to align Flutter Recipe model with Appwrite database schema

**Decision**: Extend existing Recipe model with full Appwrite schema fields, use snake_case mapping

**Rationale**:
- Appwrite uses snake_case keys (e.g., `prep_time`, `author_name`)
- Flutter model uses camelCase (e.g., `prepTime`, `authorName`)
- Add `fromDocument` factory constructor for Appwrite Document mapping
- Consider dart_mappable for cleaner serialization (per constitution)

**Field Mapping** (from appwrite.config.json):
| Appwrite Field | Dart Field | Type |
|----------------|------------|------|
| `$id` | `id` | `String` |
| `name` | `name` | `String` |
| `description` | `description` | `String?` |
| `image` | `images` | `List<String>` |
| `prep_time` | `prepTime` | `String?` |
| `cook_time` | `cookTime` | `String?` |
| `total_time` | `totalTime` | `String?` |
| `recipe_yield` | `recipeYield` | `List<String>?` |
| `ingredients` | `ingredients` | `List<String>` |
| `instructions` | `instructions` | `List<String>` |
| `author_name` | `authorName` | `String?` |
| `author_url` | `authorUrl` | `String?` |
| `nutrition_*` | `nutrition` | `NutritionInfo?` (nested) |

**Alternatives Considered**:
- Manual JSON parsing: Current approach, works but verbose
- json_serializable: Alternative, but dart_mappable preferred per constitution

---

### 4. Error State UI Pattern

**Task**: Research inline error banner pattern for recipe card skeleton

**Decision**: Transform skeleton card to error state with FAlert-style banner inside card bounds

**Rationale**:
- forUI provides `FAlert` component for error/warning states
- Error state replaces shimmer animation while keeping card structure
- Include retry action and descriptive error message
- Maintain visual continuity with skeleton dimensions

**Implementation Pattern**:
```dart
FCard(
  child: Column(
    children: [
      // Error icon area (replaces hero image skeleton)
      Container(
        height: 200,
        child: Center(child: FAssets.icons.circleAlert(...)),
      ),
      // Error message area
      FAlert(
        icon: FAssets.icons.triangleAlert,
        title: Text('Import Failed'),
        subtitle: Text(errorMessage),
      ),
      // Retry button
      FButton(onPress: onRetry, child: Text('Retry')),
    ],
  ),
)
```

**Alternatives Considered**:
- Toast/Snackbar: Rejected - doesn't maintain context with specific import
- Full-screen error: Rejected - too disruptive for single-item failure
- Dialog: Rejected - blocks user from other actions

---

### 5. Recipe Preview Page Layout

**Task**: Research full-detail scrollable layout pattern for recipe preview

**Decision**: Use `CustomScrollView` with `SliverAppBar` for hero image and sliver list for content

**Rationale**:
- `SliverAppBar` with `expandedHeight` provides elegant hero image with scroll collapse
- Content sections (metadata, ingredients, instructions) as sliver children
- Matches iOS/Android native patterns for detail views
- forUI components for content cards and lists

**Layout Structure**:
```
┌─────────────────────────────────────┐
│          Hero Image                 │  SliverAppBar (expandedHeight: 300)
│      (collapses on scroll)          │
├─────────────────────────────────────┤
│  ⏱ 15 min prep  │  🍽 4 servings    │  Metadata row
├─────────────────────────────────────┤
│  Ingredients                        │  Section header
│  • 400g spaghetti                   │  
│  • 200g pancetta                    │  List items
│  • ...                              │
├─────────────────────────────────────┤
│  Instructions                       │  Section header
│  1. Cook the spaghetti...           │  
│  2. Fry the pancetta...             │  Numbered steps
│  • ...                              │
├─────────────────────────────────────┤
│  Source: allrecipes.com             │  Attribution footer
└─────────────────────────────────────┘
```

**Alternatives Considered**:
- Tabs (ingredients/instructions): More complex, less scannable
- Accordion sections: Hides content by default
- Cards layout: Visual fragmentation

---

### 6. State Machine for Recipe Import Flow

**Task**: Research state management pattern for multi-stage import flow

**Decision**: Use enum-based state with Signals computed for derived UI states

**Rationale**:
- Import flow has clear stages: idle → submitting → extracting → preview/error
- Appwrite provides status via real-time (REQUESTED → IN_PROGRESS → COMPLETED/FAILED)
- Computed signals derive UI state from request status

**State Flow**:
```
User submits URL
     │
     ▼
┌─────────────┐
│   IDLE      │ (initial state)
└─────────────┘
     │ importRecipe(url)
     ▼
┌─────────────┐
│  SUBMITTING │ (creating request document)
└─────────────┘
     │ document created
     ▼
┌─────────────┐
│  REQUESTED  │◀─── Appwrite status
└─────────────┘
     │ function triggered
     ▼
┌─────────────┐
│ IN_PROGRESS │◀─── Appwrite status (shimmer showing)
└─────────────┘
     │
     ├───────────────────┐
     ▼                   ▼
┌─────────────┐    ┌─────────────┐
│  COMPLETED  │    │   FAILED    │
└─────────────┘    └─────────────┘
     │                   │
     ▼                   ▼
 Preview page      Error card
```

**Alternatives Considered**:
- Bloc: Constitution prefers Signals for new code
- StreamController: More boilerplate than Signals

---

## Dependencies to Add

```yaml
dependencies:
  shimmer: ^3.0.0  # For skeleton loading effect
```

No other new dependencies required - Appwrite SDK already included.

## Resolved Clarifications

| Original Unknown | Resolution |
|-----------------|------------|
| Real-time mechanism | Appwrite Realtime subscription on recipe_request collection |
| Shimmer component | Custom widget using shimmer package (forUI lacks built-in) |
| Recipe model schema | Extend with full Appwrite fields + fromDocument factory |
| Error UI pattern | Inline error card replacing skeleton |
| Preview layout | CustomScrollView with SliverAppBar hero pattern |
