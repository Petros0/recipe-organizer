# Tasks: Async Recipe Import UX

**Input**: Design documents from `/specs/002-async-recipe-import-ux/`  
**Prerequisites**: spec.md ✅  
**Depends On**: `001-url-recipe-import` (all phases complete)

**Tests**: Test tasks included per constitution Test Discipline (Section V). Tests follow the Arrange-Act-Assert pattern using mocktail for mocking.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: State Management Refactor

**Purpose**: Update HomeController to support multiple concurrent pending imports

- [X] T001 [US1] Add `pendingImports` signal (list of RecipeRequest) to HomeController at lib/features/home/state/home_controller.dart
- [X] T002 [US1] Refactor `importRecipe` method to add to pendingImports list instead of blocking at lib/features/home/state/home_controller.dart
- [X] T003 [US1] Create `_subscribeToImport` method that manages realtime subscription per pending import at lib/features/home/state/home_controller.dart
- [X] T004 [US1] Add `_onImportComplete` callback to remove from pending and add recipe to list at lib/features/home/state/home_controller.dart
- [X] T005 [US1] Add `_onImportError` callback to update pending import with error state at lib/features/home/state/home_controller.dart
- [X] T006 [P] [US2] Add `retryImport(String requestId)` method to HomeController at lib/features/home/state/home_controller.dart
- [X] T007 [P] [US2] Add `dismissImport(String requestId)` method to HomeController at lib/features/home/state/home_controller.dart

**Checkpoint**: Controller supports multiple concurrent imports with retry/dismiss

---

## Phase 2: UI Updates

**Purpose**: Update UI to show pending imports on home screen and navigate immediately

- [X] T008 [US1] Update ImportRecipeDialog to close and navigate immediately after triggering import at lib/features/home/view/widgets/import_recipe_dialog.dart
- [X] T009 [US1] Update HomePage to render pendingImports as skeleton/error cards at top of recipe list at lib/features/home/view/home_page.dart
- [X] T010 [P] [US1] Create PendingImportCard widget that switches between skeleton and error states at lib/features/home/view/widgets/pending_import_card.dart
- [X] T011 [US2] Update RecipeErrorCard to accept onRetry and onDismiss callbacks at lib/features/home/view/widgets/recipe_error_card.dart
- [X] T012 [US2] Wire retry/dismiss actions from PendingImportCard to HomeController at lib/features/home/view/home_page.dart

**Checkpoint**: UI shows pending imports on home screen with full interaction

---

## Phase 3: Localization

**Purpose**: Add new localization strings

- [X] T013 [P] Add localization keys for pending import states at lib/l10n/arb/app_en.arb:
  - `importPending`: "Importing recipe..."
  - `importRetry`: "Retry"
  - `importDismiss`: "Dismiss"
  - `multipleImportsInProgress`: "{count} recipes importing"

**Checkpoint**: All strings localized

---

## Phase 4: Testing

**Purpose**: Unit and widget tests for new functionality

- [X] T014 [P] Add HomeController tests for pendingImports management at test/features/home/state/home_controller_test.dart
- [X] T015 [P] Add HomeController tests for concurrent import handling at test/features/home/state/home_controller_test.dart
- [X] T016 [P] Create PendingImportCard widget test at test/features/home/view/widgets/pending_import_card_test.dart
- [X] T017 Update ImportRecipeDialog test for immediate navigation at test/features/home/view/widgets/import_recipe_dialog_test.dart

**Checkpoint**: Core test coverage complete

---

## Phase 5: Polish

**Purpose**: Final improvements and validation

- [X] T018 [P] Run dart fix --apply and dart format on all modified files
- [X] T019 [P] Run very_good_analysis linter and fix any issues
- [X] T020 Generate localizations with flutter gen-l10n
- [ ] T021 Manual validation: import 3 recipes concurrently, verify all resolve correctly

---

## Phase 6: Bug Fixes (Post-Implementation)

**Purpose**: Fix issues found during manual testing

- [X] T022 Fix RecipeSkeletonCard overflow in grid - use Expanded instead of fixed height at lib/features/home/view/widgets/recipe_skeleton_card.dart
- [X] T023 Remove "Save to Collection" button from RecipePreviewPage - recipes are already saved by backend at lib/features/home/view/recipe_preview_page.dart
- [X] T024 Update HomePage _navigateToPreview to use delete action instead of save at lib/features/home/view/home_page.dart

**Checkpoint**: UI issues resolved

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: No dependencies - can start immediately (requires 001-url-recipe-import complete)
- **Phase 2**: Depends on Phase 1 (T001-T005) for controller signals
- **Phase 3**: No dependencies - can run in parallel with Phase 2
- **Phase 4**: Depends on Phase 1-2 for implementation
- **Phase 5**: Depends on all phases complete
- **Phase 6**: Post-implementation bug fixes

### Within Phases

- T001-T005 must be sequential (building up controller logic)
- T006-T007 can run in parallel after T005
- T008-T009 sequential; T010 can run in parallel
- All test tasks (T014-T017) can run in parallel

### Parallel Opportunities

- T006, T007 can run in parallel (different methods)
- T010, T013 can run in parallel (different files)
- T014, T015, T016 can run in parallel (independent tests)

---

## Notes

- Reuse existing `RecipeSkeletonCard` and `RecipeErrorCard` from 001-url-recipe-import
- The `PendingImportCard` is a wrapper that chooses skeleton vs error based on state
- Navigation happens via existing routing - no new routes needed
- Appwrite realtime subscriptions already work - just need to manage multiple
