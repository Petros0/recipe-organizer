import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/view/widgets/import_recipe_dialog.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

void main() {
  group('ImportRecipeDialog', () {
    Widget buildTestWidget({
      required void Function(String) onImport,
      bool isLoading = false,
    }) {
      return MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: FTheme(
          data: FThemes.zinc.dark,
          child: Scaffold(
            body: ImportRecipeDialog(
              onImport: onImport,
              isLoading: isLoading,
            ),
          ),
        ),
      );
    }

    testWidgets('renders dialog with title and description', (tester) async {
      // Arrange & Act
      await tester.pumpWidget(buildTestWidget(onImport: (_) {}));
      await tester.pump();

      // Assert
      expect(find.text('Import Recipe'), findsOneWidget);
      expect(
        find.text('Paste the URL of a recipe page to import it'),
        findsOneWidget,
      );
    });

    testWidgets('renders URL input field', (tester) async {
      // Arrange & Act
      await tester.pumpWidget(buildTestWidget(onImport: (_) {}));
      await tester.pump();

      // Assert
      expect(find.text('Recipe URL'), findsOneWidget);
      expect(find.byType(FTextField), findsOneWidget);
    });

    testWidgets('calls onImport with valid URL', (tester) async {
      // Arrange
      String? submittedUrl;
      const testUrl = 'https://example.com/recipe';

      await tester.pumpWidget(
        buildTestWidget(
          onImport: (url) => submittedUrl = url,
        ),
      );
      await tester.pump();

      // Act
      await tester.enterText(find.byType(TextField), testUrl);
      await tester.tap(find.text('Import'));
      await tester.pump();

      // Assert
      expect(submittedUrl, testUrl);
    });

    testWidgets('shows validation error for empty URL', (tester) async {
      // Arrange
      String? submittedUrl;

      await tester.pumpWidget(
        buildTestWidget(
          onImport: (url) => submittedUrl = url,
        ),
      );
      await tester.pump();

      // Act - tap import without entering URL
      await tester.tap(find.text('Import'));
      await tester.pump();

      // Assert
      expect(submittedUrl, isNull);
      expect(find.text('Please enter a valid URL'), findsOneWidget);
    });

    testWidgets('shows validation error for invalid URL format', (
      tester,
    ) async {
      // Arrange
      String? submittedUrl;

      await tester.pumpWidget(
        buildTestWidget(
          onImport: (url) => submittedUrl = url,
        ),
      );
      await tester.pump();

      // Act - enter invalid URL
      await tester.enterText(find.byType(TextField), 'not-a-valid-url');
      await tester.tap(find.text('Import'));
      await tester.pump();

      // Assert
      expect(submittedUrl, isNull);
      expect(
        find.text('URL must start with http:// or https://'),
        findsOneWidget,
      );
    });

    testWidgets('shows loading indicator when isLoading is true', (
      tester,
    ) async {
      // Arrange & Act
      await tester.pumpWidget(
        buildTestWidget(
          onImport: (_) {},
          isLoading: true,
        ),
      );
      await tester.pump();

      // Assert
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });
  });
}
