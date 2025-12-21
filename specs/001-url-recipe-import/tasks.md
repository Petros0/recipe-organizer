# Tasks: URL-Based Recipe Import

**Input**: Design documents from `/specs/001-url-recipe-import/`
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅, quickstart.md ✅

**Tests**: Test tasks included per constitution Test Discipline (Section V). Tests follow the Arrange-Act-Assert pattern using mocktail for mocking.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Mobile app**: `lib/` for source, `test/` for tests
- Following feature-first architecture from AGENTS.md

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and dependency configuration

- [X] T001 Add shimmer package dependency to pubspec.yaml
- [X] T002 [P] Create Appwrite constants file at lib/core/appwrite/appwrite_constants.dart
- [X] T003 [P] Create RealtimeService wrapper at lib/core/appwrite/realtime_service.dart
- [X] T004 Register new services in dependency injection at lib/core/di.dart

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core models and data layer that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T005 Create RecipeRequest model with status enum at lib/features/home/model/recipe_request.dart
- [X] T006 [P] Update Recipe model with fromDocument factory and all Appwrite fields at lib/features/home/model/recipe.dart
- [X] T007 [P] Create NutritionInfo model at lib/features/home/model/nutrition_info.dart
- [X] T008 Implement RecipeRequestRepository with createRequest and subscribeToUpdates at lib/features/home/data/recipe_request_repository.dart
- [X] T009 [P] Create RecipeRepository with getRecipeByRequestId and listRecipes at lib/features/home/data/recipe_repository.dart
- [X] T010 Create RecipeImportService with importFromUrl and subscribeToRequest at lib/features/home/service/recipe_import_service.dart
- [X] T011 Add localization keys for import feature to lib/l10n/arb/app_en.arb

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Import Recipe from URL with Structured Data (Priority: P1) 🎯 MVP

**Goal**: Enable users to paste a recipe URL and have the app extract all recipe details via schema.org JSON-LD for preview and saving

**Independent Test**: Paste a URL from AllRecipes/Food Network, verify skeleton shows during extraction, recipe preview displays with all fields, and save adds to collection

### Implementation for User Story 1

- [X] T012 [US1] Create ImportState enum for state machine at lib/features/home/state/import_state.dart
- [X] T013 [P] [US1] Create RecipeSkeletonCard widget with shimmer effect at lib/features/home/view/widgets/recipe_skeleton_card.dart
- [X] T014 [US1] Update HomeController with importRecipe method, signals for activeRequest/importState, and real-time subscription at lib/features/home/state/home_controller.dart
- [X] T015 [US1] Create RecipePreviewPage with SliverAppBar hero image and scrollable content sections at lib/features/home/view/recipe_preview_page.dart
- [X] T016 [P] [US1] Create metadata row widget for prep/cook time and servings at lib/features/home/view/widgets/recipe_metadata_row.dart
- [X] T017 [P] [US1] Create ingredients section widget at lib/features/home/view/widgets/ingredients_section.dart
- [X] T018 [P] [US1] Create instructions section widget at lib/features/home/view/widgets/instructions_section.dart
- [X] T019 [US1] Update HomePage to show skeleton during extraction and navigate to preview on completion at lib/features/home/view/home_page.dart
- [X] T020 [US1] Add URL input dialog/modal for recipe import trigger at lib/features/home/view/widgets/import_recipe_dialog.dart
- [X] T020a [US1] Add client-side URL validation (well-formed HTTP/HTTPS check) in import dialog at lib/features/home/view/widgets/import_recipe_dialog.dart
- [X] T021 [US1] Implement save to collection flow in HomeController at lib/features/home/state/home_controller.dart

**Checkpoint**: User Story 1 complete - users can import recipes from sites with structured data

---

## Phase 4: User Story 2 - Import Recipe with LLM Fallback (Priority: P2)

**Goal**: Enable recipe import from websites without schema.org markup via LLM extraction

**Independent Test**: Submit a URL from a blog without JSON-LD markup, verify LLM extraction triggers with progress indication

**Note**: Backend LLM fallback is already implemented in Go functions. Frontend changes are minimal.

### Implementation for User Story 2

- [X] T022 [US2] Add LLM extraction status indication in RecipeSkeletonCard (extended timing message) at lib/features/home/view/widgets/recipe_skeleton_card.dart
- [X] T023 [US2] Add localization keys for LLM extraction progress at lib/l10n/arb/app_en.arb
- [PARKED] T024 [US2] Update RecipePreviewPage to support editable fields before save at lib/features/home/view/recipe_preview_page.dart

**Checkpoint**: User Story 2 complete - users can import from any recipe website

---

## Phase 5: User Story 3 - Handle Failed or Partial Imports (Priority: P3)

**Goal**: Provide clear feedback for extraction failures with retry and manual fill options

**Independent Test**: Submit a URL that returns 404 or bot protection, verify inline error banner appears with retry action

### Implementation for User Story 3

- [X] T025 [US3] Create RecipeErrorCard with inline error banner and retry button at lib/features/home/view/widgets/recipe_error_card.dart
- [X] T026 [US3] Add error state handling in HomeController with retry capability at lib/features/home/state/home_controller.dart
- [X] T027 [US3] Update HomePage to display RecipeErrorCard on failed import at lib/features/home/view/home_page.dart
- [PARKED] T028 [P] [US3] Add partial data indicators in RecipePreviewPage for missing fields at lib/features/home/view/recipe_preview_page.dart
- [X] T029 [P] [US3] Add localization keys for error messages at lib/l10n/arb/app_en.arb
- [PARKED] T030 [US3] Implement manual field editing for missing data in preview at lib/features/home/view/recipe_preview_page.dart

**Checkpoint**: User Story 3 complete - graceful error handling with recovery options

---

## Phase 6: User Story 4 - Source Attribution (Priority: P4)

**Goal**: Display and preserve original source URL and author information

**Independent Test**: Import a recipe and verify source URL and author name are visible in saved recipe

### Implementation for User Story 4

- [X] T031 [US4] Create SourceAttribution widget for recipe preview and detail view at lib/features/home/view/widgets/source_attribution.dart
- [X] T032 [US4] Add source attribution section to RecipePreviewPage with tappable link at lib/features/home/view/recipe_preview_page.dart
- [X] T033 [US4] Add localization key for source label at lib/l10n/arb/app_en.arb

**Checkpoint**: User Story 4 complete - proper attribution for imported recipes

---

## Phase 7: Testing

**Purpose**: Unit and widget tests for core functionality per constitution Test Discipline

- [X] T034 [P] Create RecipeRequest model unit tests at test/features/home/model/recipe_request_test.dart
- [X] T035 [P] Create Recipe model unit tests at test/features/home/model/recipe_test.dart
- [X] T036 [P] Create NutritionInfo model unit tests at test/features/home/model/nutrition_info_test.dart
- [X] T037 Create RecipeRequestRepository unit tests with mocktail at test/features/home/data/recipe_request_repository_test.dart
- [X] T038 [P] Create RecipeRepository unit tests at test/features/home/data/recipe_repository_test.dart
- [X] T039 Create RecipeImportService unit tests at test/features/home/service/recipe_import_service_test.dart
- [X] T040 Create HomeController unit tests at test/features/home/state/home_controller_test.dart
- [X] T041 [P] Create RecipeSkeletonCard widget test at test/features/home/view/widgets/recipe_skeleton_card_test.dart
- [X] T042 [P] Create RecipeErrorCard widget test at test/features/home/view/widgets/recipe_error_card_test.dart
- [X] T043 Create ImportRecipeDialog widget test at test/features/home/view/widgets/import_recipe_dialog_test.dart

**Checkpoint**: Core test coverage complete

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements and validation

- [X] T044 [P] Run dart fix --apply and dart format on all modified files
- [X] T045 [P] Run very_good_analysis linter and fix any issues
- [X] T046 Generate localizations with flutter gen-l10n
- [PARKED] T047 Run quickstart.md validation scenarios manually
- [PARKED] T048 Verify 60 fps UI performance during skeleton animation

---

## Phase 9: User-Scoped Data (Row-Level Security)

**Purpose**: Restrict recipe_request and recipe documents to the authenticated user who created them

### Implementation

- [X] T049 [P] Update Flutter to pass `x-appwrite-user-id` header when calling recipe-request function at lib/features/home/data/recipe_request_repository.dart
- [X] T050 [P] Update recipe-request Go function to read `x-appwrite-user-id` header and store userId in document at functions/recipe-request/main.go and functions/recipe-request/appwrite_client.go
- [X] T051 [P] Update recipe-request-processor Go function to read userId from recipe_request and store in recipe document at functions/recipe-request-processor/appwrite_client.go
- [X] T052 [P] Set user-specific permissions (RLS) on recipe_request documents using `role.User(userId)` at functions/recipe-request/appwrite_client.go
- [X] T053 [P] Set user-specific permissions (RLS) on recipe documents using `role.User(userId)` at functions/recipe-request-processor/appwrite_client.go
- [ ] T054 Update Appwrite recipe_request collection schema to include userId field (string, required)
- [ ] T055 Update Appwrite recipe collection schema to include userId field (string, required)
- [X] T056 Update Flutter RecipeRequest model to include userId field at lib/features/home/model/recipe_request.dart
- [X] T057 Update Flutter Recipe model to include userId field at lib/features/home/model/recipe.dart

**Checkpoint**: User-scoped data complete - users can only access their own recipes

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - US1 can proceed immediately after Phase 2
  - US2 can start after Phase 2 (or after US1 if preferred)
  - US3 can start after Phase 2 (or after US1 if preferred)
  - US4 can start after Phase 2 (or after US1 if preferred)
- **Testing (Phase 7)**: Can start after Phase 2 for models; controller/widget tests after respective implementation
- **Polish (Phase 8)**: Depends on all user stories and tests being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Builds on US1 skeleton/preview components
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Adds error handling to US1 flow
- **User Story 4 (P4)**: Can start after Foundational (Phase 2) - Adds attribution to US1 preview

### Within Each User Story

- Models before services
- Services before controllers
- Controllers before views
- Core widgets before pages
- Commit after each task or logical group

### Parallel Opportunities

- T002, T003 can run in parallel (different files)
- T006, T007, T009 can run in parallel (different models/repos)
- T013, T016, T017, T018 can run in parallel (independent widgets)
- T028, T029 can run in parallel (different files)
- T034, T035 can run in parallel (independent tools)

---

## Parallel Example: User Story 1 Widgets

```bash
# Launch all independent widgets together:
Task T013: "Create RecipeSkeletonCard widget at lib/features/home/view/widgets/recipe_skeleton_card.dart"
Task T016: "Create metadata row widget at lib/features/home/view/widgets/recipe_metadata_row.dart"
Task T017: "Create ingredients section widget at lib/features/home/view/widgets/ingredients_section.dart"
Task T018: "Create instructions section widget at lib/features/home/view/widgets/instructions_section.dart"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational (T005-T011)
3. Complete Phase 3: User Story 1 (T012-T021)
4. **STOP and VALIDATE**: Test import from major recipe website
5. Deploy/demo if ready

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test LLM fallback → Deploy/Demo
4. Add User Story 3 → Test error handling → Deploy/Demo
5. Add User Story 4 → Test attribution → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (core flow)
   - Developer B: User Story 3 (error handling) - can stub US1 UI
   - Developer C: User Story 4 (attribution widget)
3. After US1 complete: Developer B integrates error handling
4. User Story 2 last (minimal frontend changes, depends on backend)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Backend Go functions are COMPLETE - only Flutter frontend implementation needed
- forUI components should be used where available (FAlert, FButton, FCard)
- shimmer package for skeleton loading effect (forUI lacks built-in shimmer)
- All strings must use localization (lib/l10n/arb/app_en.arb)
- Follow Signals pattern for state management (not Bloc)
- Commit after each task or logical group
