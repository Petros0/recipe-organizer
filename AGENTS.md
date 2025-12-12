# AGENTS.md - Recipe Organizer (Recipify)

This document provides guidance for AI agents working on this codebase.

## Project Overview

**Recipify** is a cross-platform mobile/web application that transforms content from any source into structured recipes. The app allows users to import recipes from websites by extracting schema.org Recipe structured data.

### Tech Stack

| Layer          | Technology                                                        |
| -------------- | ----------------------------------------------------------------- |
| Frontend       | Flutter 3.35+ / Dart 3.9+                                         |
| State Mgmt     | [Signals](https://dartsignals.dev/llms.txt)                       |
| UI Components  | [forUI](https://forui.dev/docs)                                   |
| Backend        | [Appwrite](https://appwrite.io) (BaaS)                            |
| Functions      | Appwrite Functions (Go 1.23)                                      |
| Linting        | [very_good_analysis](https://pub.dev/packages/very_good_analysis) |
| Testing        | flutter_test, bloc_test, mocktail                                 |
| Serialization  | dart_mappable (for toJson/fromJson)                               |

### Supported Platforms

iOS, Android

## Project Structure

```
lib/
  core/                     # Shared utilities and constants
    constants/              # App-wide constants
    theme/                  # Theme configuration (FThemeData)
    utils/                  # Utility functions
    widgets/                # Reusable widgets
  app/
    app.dart                # App barrel export
    view/
      app.dart              # Root MaterialApp widget
  features/
    <feature_name>/
      view/
        <feature>_page.dart       # UI page/screen
      state/
        <feature>_controller.dart # State management (Signals)
      service/
        <feature>_service.dart    # Business logic
      data/
        <feature>_repository.dart # Data persistence
      <feature>.dart              # Barrel export
  l10n/
    arb/                    # Translation files (app_en.arb, app_es.arb)
    gen/                    # Generated localization code
    l10n.dart               # Localization exports
  bootstrap.dart            # Cross-environment initialization
  main_development.dart     # Development entry point
  main_staging.dart         # Staging entry point
  main_production.dart      # Production entry point

test/
  app/                      # App-level tests
  <feature>/                # Feature tests (mirror lib structure)
  helpers/                  # Test utilities and pump_app

functions/
  # Go functions for recipe extraction
  recipe-request/
    main.go                   
  recipe-request-processor/
    main.go
```

## Architecture Patterns

### Feature Architecture

Each feature follows a layered architecture:

1. **View (Page)** - UI widgets, minimal logic
2. **Controller** - Manages UI state using Signals
3. **Service** - Implements business logic, called by controllers
4. **Repository** - Handles data persistence and external APIs
5. **Cache** (optional) - Caches repository results

### State Management with Signals

```dart
class _FeaturePageState extends State<FeaturePage> with SignalsMixin {
  late final counter = createSignal(0);

  @override
  Widget build(BuildContext context) {
    return Text('Value: $counter');
  }
}
```

### UI Components with forUI

IMPORTANT: Always first make sure that a component doesn't exist with [forUI](https://forui.dev/) before implementing it by yourself.

Use `FThemes` for theming and forUI widgets for consistent UI:

```dart
final theme = FThemes.zinc.dark;
return MaterialApp(
  theme: theme.toApproximateMaterialTheme(),
  builder: (_, child) => FAnimatedTheme(data: theme, child: child!),
);
```

## Code Conventions

### Naming

- **Classes**: `PascalCase`
- **Variables/Functions**: `camelCase`
- **Files/Directories**: `snake_case`
- **Environment Variables**: `UPPERCASE`
- **Booleans**: Verb prefixes (`isLoading`, `hasError`, `canDelete`)
- **Functions**: Start with verbs (`fetchRecipes`, `saveIngredient`)

### File Organization

- One export per file
- Use barrel files (`feature.dart`) to export all feature components
- Keep widgets small and focused (<200 lines)

### Flutter Best Practices

- Use `const` constructors wherever possible
- Avoid deeply nested widget trees - extract to smaller components
- Break down large widgets into focused, reusable widgets
- Prefer composition over inheritance

## Backend Services

### Appwrite Configuration

- **Project ID**: `691f8b990030db50617a`
- **Endpoint**: `https://fra.cloud.appwrite.io/v1`
- **Project Name**: `recipe-organizer`

### Cloud Functions

#### `fetch-recipe-from-website` (Go)

Extracts recipe data from URLs using schema.org JSON-LD:

- **Runtime**: Go 1.23
- **Endpoint**: `/ping` returns "Pong"
- **Input**: `url` (query param or JSON body)
- **Output**: Recipe JSON with name, image, ingredients, instructions, nutrition
- **Features**: HTTP client with browser headers, headless browser fallback for bot protection

## Common Commands

```bash
# Run app (development flavor)
flutter run --flavor development --target lib/main_development.dart

# Or use justfile
just dev

# Run tests
very_good test --coverage --test-randomize-ordering-seed random

# Generate localizations
just l10n
# or: flutter gen-l10n --arb-dir="lib/l10n/arb"

# Format code
just format
# or: dart fix --apply && dart format . --line-length 120

# List devices
just devices
```

## Testing Guidelines

### Test Structure

- **Unit tests**: Test controllers, services, repositories
- **Widget tests**: Test UI components in isolation
- **Integration tests**: Test API modules and flows

### Conventions

- Follow **Arrange-Act-Assert** pattern
- Name variables: `inputX`, `mockX`, `actualX`, `expectedX`
- Use `mocktail` for mocking
- Use `pump_app.dart` helper for widget tests

### Example Widget Test

```dart
testWidgets('renders counter value', (tester) async {
  await tester.pumpApp(const CounterPage());
  expect(find.text('0'), findsOneWidget);
});
```

## Localization

Translations are managed via ARB files in `lib/l10n/arb/`:

```json
{
  "@@locale": "en",
  "featureTitle": "Title",
  "@featureTitle": {
    "description": "Description for translators"
  }
}
```

Access in code:

```dart
final l10n = context.l10n;
Text(l10n.featureTitle);
```

## Environment Flavors

| Flavor      | Entry Point              | Use Case           |
| ----------- | ------------------------ | ------------------ |
| development | `main_development.dart`  | Local development  |
| staging     | `main_staging.dart`      | Testing/QA         |
| production  | `main_production.dart`   | Production release |

## Key Dependencies

```yaml
# State Management
signals: ^6.2.0

# UI Framework
forui: ^0.16.0
forui_assets: ^0.16.0

# Bloc (available but prefer Signals)
bloc: ^9.0.1
flutter_bloc: ^9.1.1

# Localization
flutter_localizations: sdk
intl: ^0.20.2

# Dev/Test
very_good_analysis: ^10.0.0
mocktail: ^1.0.4
bloc_test: ^10.0.0
```

## Important Notes for Agents

1. **Prefer Signals over Bloc** for new state management
2. **Use forUI components** for UI consistency
3. **Follow the feature structure** when adding new features
4. **Run `just format`** after making changes
5. **Keep widgets shallow** - extract components to avoid deep nesting
6. **Use const constructors** to minimize rebuilds
7. **Add translations** for all user-facing strings
8. **Write tests** following the Arrange-Act-Assert pattern
9. **Use dart_mappable** for JSON serialization
10. **Check `.cursor/rules/flutter.mdc`** for detailed Flutter guidelines

