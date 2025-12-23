# Feature Specification: Async Recipe Import UX

**Feature ID**: `002-async-recipe-import-ux`  
**Created**: 2025-12-23  
**Status**: Draft  
**Depends On**: `001-url-recipe-import`  
**Input**: Improve recipe import UX by making it non-blocking - navigate to home immediately and show loading state there

## Problem Statement

Currently when the user imports a recipe, the app waits for the recipe to be loaded before proceeding. This creates two UX issues:

1. **Blocking behavior**: Users cannot check their other recipes while waiting
2. **Perceived hang**: The application may appear stuck during extraction (especially for LLM fallback which takes longer)

## Solution Overview

When the user clicks "import" on a recipe URL:

1. **Immediate navigation**: Navigate to the home screen right away instead of waiting
2. **Optimistic UI**: Show the recipe in a loading/skeleton state on the home screen
3. **Seamless transition**: Once loaded, replace the skeleton with the actual recipe content

This provides a non-blocking, optimistic UI experience where users can continue browsing their recipes while new ones are being imported in the background.

## User Scenarios & Testing

### User Story 1 - Non-Blocking Import Flow (Priority: P1)

As a user, I want to continue browsing my recipes while a new recipe is being imported, so I don't have to wait and stare at a loading screen.

**Why this priority**: This is the core improvement. Users should never feel blocked by background operations.

**Independent Test**: Paste a recipe URL, verify immediate navigation to home, see skeleton card while loading, then see actual recipe when ready.

**Acceptance Scenarios**:

1. **Given** a user submits a recipe URL in the import dialog, **When** they tap "Import", **Then** they are immediately navigated to the home screen (within 300ms).
2. **Given** a recipe import is in progress, **When** the user is on the home screen, **Then** they see a skeleton card at the top of their recipe list indicating the import is loading.
3. **Given** a recipe import completes successfully, **When** the extraction finishes, **Then** the skeleton card is replaced with the actual recipe card (no navigation required).
4. **Given** a user has an import in progress, **When** they interact with their existing recipes, **Then** they can browse, open, and view them without interruption.

---

### User Story 2 - Import Error on Home Screen (Priority: P2)

As a user, when a background import fails, I want to see the error on my home screen so I can retry without losing context.

**Why this priority**: Error handling must work within the new flow.

**Independent Test**: Import a URL that fails, verify error card appears on home screen with retry option.

**Acceptance Scenarios**:

1. **Given** a recipe import fails (network error, invalid page, etc.), **When** the error occurs, **Then** the skeleton card transforms into an error card on the home screen.
2. **Given** an error card is shown, **When** the user taps "Retry", **Then** the import is attempted again and the card returns to skeleton state.
3. **Given** an error card is shown, **When** the user taps "Dismiss", **Then** the error card is removed from the list.

---

### User Story 3 - Multiple Concurrent Imports (Priority: P3)

As a user, I want to import multiple recipes at once, so I can batch-save recipes from my browser tabs.

**Why this priority**: Natural extension once imports are non-blocking.

**Independent Test**: Import 3 URLs in quick succession, verify 3 skeleton cards appear, each resolves independently.

**Acceptance Scenarios**:

1. **Given** the user has one import in progress, **When** they start another import, **Then** a second skeleton card appears on the home screen.
2. **Given** multiple imports are in progress, **When** one completes, **Then** only that skeleton is replaced (others continue loading).
3. **Given** multiple imports are in progress, **When** one fails, **Then** only that skeleton shows an error (others continue loading).

---

## Requirements

### Functional Requirements

- **FR-001**: System MUST navigate to home screen immediately after user confirms import (before extraction completes).
- **FR-002**: System MUST display a skeleton/loading card on the home screen for each in-progress import.
- **FR-003**: System MUST support multiple concurrent recipe imports.
- **FR-004**: System MUST update the skeleton card to show the actual recipe when extraction completes.
- **FR-005**: System MUST transform the skeleton card into an error card when extraction fails.
- **FR-006**: System MUST provide a retry action on error cards.
- **FR-007**: System MUST provide a dismiss action to remove error cards.
- **FR-008**: System MUST persist the pending import state so it survives app restart (via Appwrite realtime subscription reconnection).
- **FR-009**: System MUST position loading/error cards prominently (top of recipe list).
- **FR-010**: System MUST allow users to interact with existing recipes while imports are in progress.

### Non-Functional Requirements

- **NFR-001**: Navigation to home screen MUST occur within 300ms of tapping "Import".
- **NFR-002**: Skeleton-to-recipe transition MUST be smooth (no jarring layout shifts).
- **NFR-003**: The UI MUST remain responsive (60fps) while imports are processing.

## Key Changes from 001-url-recipe-import

| Aspect | Before (001) | After (002) |
|--------|--------------|-------------|
| Import flow | Wait on import dialog → Show preview → Navigate | Navigate immediately → Show skeleton on home → Preview on tap |
| Skeleton location | Import dialog/modal | Home screen recipe list |
| Error display | Modal/dialog | Inline card on home screen |
| Multiple imports | Not supported (blocking) | Supported (non-blocking) |
| User context | Lost during import wait | Preserved - can browse recipes |

## Technical Approach

### State Management Changes

1. **HomeController** needs to track multiple pending imports (list of `RecipeRequest` objects)
2. Each pending import has its own Appwrite realtime subscription
3. On import start: add to pending list, navigate to home
4. On import complete: remove from pending, add recipe to list
5. On import error: update pending item with error state

### UI Changes

1. **HomePage**: Render pending imports as skeleton/error cards at top of list
2. **ImportRecipeDialog**: Remove waiting state, just trigger import and close
3. **RecipeSkeletonCard**: Already exists, reuse on home screen
4. **RecipeErrorCard**: Already exists, reuse on home screen

## Success Criteria

- **SC-001**: Zero perceived blocking - users can browse recipes within 300ms of starting import
- **SC-002**: 100% of imports (success/failure) resolve correctly on home screen
- **SC-003**: Users can successfully import 3+ recipes concurrently
- **SC-004**: No increase in import failure rate compared to blocking flow

## Assumptions

- Appwrite realtime subscriptions handle multiple concurrent requests reliably
- Users understand that skeleton cards represent pending imports
- The existing skeleton and error card components are reusable

## Out of Scope

- Import queue/history page (could be future enhancement)
- Import progress percentage (skeleton shimmer is sufficient)
- Notification when import completes (user sees it on home screen)
- Import cancellation (can dismiss on error only)
