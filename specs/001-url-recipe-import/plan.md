# Implementation Plan: URL-Based Recipe Import

**Branch**: `001-url-recipe-import` | **Date**: 2025-12-17 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-url-recipe-import/spec.md`

## Summary

Implement the Flutter UI for URL-based recipe import. The Go backend (recipe-request, recipe-request-processor functions) is complete. This plan covers connecting the frontend to Appwrite, implementing real-time status updates via Appwrite subscriptions, and building the skeleton/shimmer loading states, recipe preview page, and error handling UI.

## Technical Context

**Language/Version**: Flutter 3.35+ / Dart 3.9+  
**Primary Dependencies**: Signals ^6.2.0, forUI ^0.16.0, Appwrite SDK  
**Storage**: Appwrite Database (recipe-organizer-db)  
**Testing**: flutter_test, mocktail, very_good_analysis  
**Target Platform**: iOS, Android (primary)  
**Project Type**: Mobile application  
**Performance Goals**: Recipe import under 30s for structured data sites, 60 fps UI  
**Constraints**: Responsive UI during extraction, graceful error handling  
**Scale/Scope**: Single-user local, recipes collection management

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Feature-First Architecture | ✅ PASS | Using view/controller/service/repository layers |
| II. Signals State Management | ✅ PASS | HomeController uses Signals, new controllers will too |
| III. forUI Components | ✅ PASS | Will check forUI for shimmer/skeleton before custom impl |
| IV. Code Quality & Linting | ✅ PASS | Will run very_good_analysis |
| V. Test Discipline | ✅ PASS | Will use mocktail, AAA pattern |
| VI. Localization | ✅ PASS | Will add ARB entries for new strings |

## Project Structure

### Documentation (this feature)

```text
specs/001-url-recipe-import/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
lib/
├── core/
│   ├── appwrite/              # Appwrite client configuration (NEW)
│   │   ├── appwrite_config.dart
│   │   └── realtime_service.dart
│   └── di.dart
├── features/
│   └── home/
│       ├── data/
│       │   ├── recipe_request_repository.dart  # Update: complete implementation
│       │   └── recipe_repository.dart          # NEW: fetch recipes from DB
│       ├── model/
│       │   ├── recipe.dart                     # Update: align with DB schema
│       │   └── recipe_request.dart             # NEW: request status model
│       ├── service/
│       │   └── recipe_import_service.dart      # NEW: business logic
│       ├── state/
│       │   └── home_controller.dart            # Update: connect to repositories
│       └── view/
│           ├── home_page.dart
│           ├── recipe_preview_page.dart        # NEW: full-detail preview
│           └── widgets/
│               ├── recipe_skeleton_card.dart   # NEW: shimmer loading
│               ├── recipe_error_card.dart      # NEW: inline error state
│               └── ...

test/
└── features/
    └── home/
        ├── data/
        ├── service/
        ├── state/
        └── view/

functions/                    # COMPLETE - Go backend
├── recipe-request/
└── recipe-request-processor/
```

**Structure Decision**: Mobile app structure following feature-first architecture. Backend functions are complete; only Flutter frontend implementation needed.

## Constitution Check (Post-Design)

*Re-evaluated after Phase 1 design completion.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Feature-First Architecture | ✅ PASS | Design follows view/controller/service/repository layers |
| II. Signals State Management | ✅ PASS | ImportState enum + computed signals for derived UI |
| III. forUI Components | ✅ PASS | Researched: forUI lacks shimmer; using shimmer package |
| IV. Code Quality & Linting | ✅ PASS | Will run very_good_analysis before commits |
| V. Test Discipline | ✅ PASS | Test strategy defined in quickstart.md |
| VI. Localization | ✅ PASS | Localization keys documented in quickstart.md |

## Complexity Tracking

No constitution violations. Implementation follows established patterns.

## Generated Artifacts

| Artifact | Path | Purpose |
|----------|------|---------|
| Plan | `specs/001-url-recipe-import/plan.md` | This implementation plan |
| Research | `specs/001-url-recipe-import/research.md` | Technical decisions and rationale |
| Data Model | `specs/001-url-recipe-import/data-model.md` | Entity definitions and mappings |
| Contracts | `specs/001-url-recipe-import/contracts/api-contracts.md` | Appwrite SDK usage patterns |
| Quickstart | `specs/001-url-recipe-import/quickstart.md` | Implementation guide |

## Next Steps

Phase 2 (Task Generation) should be executed via `/speckit.tasks` command to create:
- `specs/001-url-recipe-import/tasks.md` with granular implementation tasks
