<!--
SYNC IMPACT REPORT
==================
Version change: 0.0.0 → 1.0.0 (MAJOR: Initial constitution adoption)

Modified principles: N/A (initial version)
Added sections:
  - I. Feature-First Architecture
  - II. Signals State Management
  - III. forUI Components
  - IV. Code Quality & Linting
  - V. Test Discipline
  - VI. Localization
  - Technology Stack
  - Development Workflow

Removed sections: N/A (initial version)

Templates requiring updates:
  - .specify/templates/plan-template.md ✅ (aligned with constitution)
  - .specify/templates/spec-template.md ✅ (aligned with constitution)
  - .specify/templates/tasks-template.md ✅ (aligned with constitution)
  - No commands directory found - N/A

Follow-up TODOs: None
-->

# Recipify Constitution

## Core Principles

### I. Feature-First Architecture

Every feature MUST follow the layered architecture pattern:

1. **View (Page)** - UI widgets with minimal logic (<200 lines per widget)
2. **Controller** - Manages UI state using Signals (SignalsMixin pattern)
3. **Service** - Implements business logic, called by controllers
4. **Repository** - Handles data persistence and Appwrite API interactions
5. **Cache** (optional) - Caches repository results for performance

**Rationale**: Separation of concerns enables independent testing, maintainability, and clear
responsibility boundaries. Each layer can be modified without affecting others.

**Non-negotiables**:
- One export per file with barrel file exports (`feature.dart`)
- Feature directories MUST reside under `lib/features/<feature_name>/`
- Widgets MUST use `const` constructors wherever possible
- Composition MUST be preferred over inheritance

### II. Signals State Management

All new state management MUST use Signals, not Bloc.

**Rationale**: Signals provides reactive, fine-grained state management with minimal boilerplate,
aligning with modern Flutter patterns and reducing widget rebuild overhead.

**Non-negotiables**:
- Controllers MUST use `SignalsMixin` for state management
- State changes MUST flow through signals, not direct widget setState
- Existing Bloc code MAY remain but new features MUST use Signals

### III. forUI Components

UI components MUST use forUI library before implementing custom widgets.

**Rationale**: Consistent design language, reduced code duplication, and faster development through
pre-built, tested components.

**Non-negotiables**:
- Check [forUI docs](https://forui.dev/docs) before implementing any UI component
- Use `FThemes` for theming (currently `FThemes.zinc.dark`)
- Custom widgets are permitted ONLY when forUI lacks the required component

### IV. Code Quality & Linting

All code MUST pass `very_good_analysis` linting rules.

**Rationale**: Consistent code style, early error detection, and enforced best practices across
the entire codebase.

**Non-negotiables**:
- Run `dart fix --apply && dart format . --line-length 120` before commits
- Zero linting warnings in production code
- Naming conventions enforced:
  - Classes: `PascalCase`
  - Variables/Functions: `camelCase`
  - Files/Directories: `snake_case`
  - Booleans: Verb prefixes (`isLoading`, `hasError`, `canDelete`)
  - Functions: Start with verbs (`fetchRecipes`, `saveIngredient`)

### V. Test Discipline

Tests MUST follow Arrange-Act-Assert pattern with proper coverage.

**Rationale**: Reliable tests ensure feature stability and enable confident refactoring.

**Non-negotiables**:
- Test variable naming: `inputX`, `mockX`, `actualX`, `expectedX`
- Use `mocktail` for mocking dependencies
- Use `pump_app.dart` helper for widget tests
- Run tests with: `very_good test --coverage --test-randomize-ordering-seed random`
- Test categories:
  - Unit tests: Controllers, services, repositories
  - Widget tests: UI components in isolation
  - Integration tests: API modules and flows

### VI. Localization

All user-facing strings MUST be localized via ARB files.

**Rationale**: International accessibility and maintainable string management.

**Non-negotiables**:
- Strings reside in `lib/l10n/arb/app_<locale>.arb`
- Access via `context.l10n.<key>`
- Include `@<key>` descriptions for translators
- Run `flutter gen-l10n --arb-dir="lib/l10n/arb"` after ARB changes

## Technology Stack

The following technology choices are binding for this project:

| Layer          | Technology                           | Version   |
| -------------- | ------------------------------------ | --------- |
| Frontend       | Flutter / Dart                       | 3.35+ / 3.9+ |
| State Mgmt     | Signals                              | ^6.2.0    |
| UI Components  | forUI                                | ^0.16.0   |
| Backend        | Appwrite (BaaS)                      | Current   |
| Functions      | Appwrite Functions (Go)              | 1.23      |
| Linting        | very_good_analysis                   | ^10.0.0   |
| Testing        | flutter_test, mocktail               | Latest    |
| Serialization  | dart_mappable                        | Latest    |

**Platforms**: iOS, Android (primary), Web (secondary)

## Development Workflow

### Environment Flavors

| Flavor      | Entry Point              | Use Case           |
| ----------- | ------------------------ | ------------------ |
| development | `main_development.dart`  | Local development  |
| staging     | `main_staging.dart`      | Testing/QA         |
| production  | `main_production.dart`   | Production release |

### Commands

```bash
# Development
just dev  # or: flutter run --flavor development --target lib/main_development.dart

# Testing
very_good test --coverage --test-randomize-ordering-seed random

# Formatting
just format  # or: dart fix --apply && dart format . --line-length 120

# Localization
just l10n  # or: flutter gen-l10n --arb-dir="lib/l10n/arb"
```

### Code Review Requirements

- All changes MUST pass linting and tests before merge
- Feature changes MUST include appropriate test coverage
- UI changes MUST use forUI components where applicable
- New features MUST follow the feature-first architecture

## Governance

This constitution supersedes all conflicting practices. Amendments require:

1. Documentation of the proposed change with rationale
2. Update to this constitution file with version increment
3. Propagation of changes to dependent templates and documentation

**Versioning Policy**:
- MAJOR: Principle removals, backward-incompatible governance changes
- MINOR: New principles, materially expanded guidance
- PATCH: Clarifications, wording refinements, typo fixes

**Compliance**: All code contributions MUST verify alignment with these principles. Reference
`AGENTS.md` for detailed runtime development guidance.

**Version**: 1.0.0 | **Ratified**: 2025-12-17 | **Last Amended**: 2025-12-17
