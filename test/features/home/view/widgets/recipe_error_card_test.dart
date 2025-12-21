import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_error_card.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

void main() {
  group('RecipeErrorCard', () {
    Widget buildTestWidget({
      required String errorMessage,
      VoidCallback? onRetry,
      VoidCallback? onDismiss,
    }) {
      return MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: FTheme(
          data: FThemes.zinc.dark,
          child: Scaffold(
            body: RecipeErrorCard(
              errorMessage: errorMessage,
              onRetry: onRetry,
              onDismiss: onDismiss,
            ),
          ),
        ),
      );
    }

    testWidgets('renders error message', (tester) async {
      // Arrange
      const errorMessage = 'Test error message';

      // Act
      await tester.pumpWidget(buildTestWidget(errorMessage: errorMessage));
      await tester.pump();

      // Assert
      expect(find.text(errorMessage), findsOneWidget);
      expect(find.text('Import failed'), findsOneWidget);
    });

    testWidgets('renders retry button when onRetry provided', (tester) async {
      // Arrange
      var retryPressed = false;

      // Act
      await tester.pumpWidget(
        buildTestWidget(
          errorMessage: 'Error',
          onRetry: () => retryPressed = true,
        ),
      );
      await tester.pump();

      // Assert
      expect(find.text('Retry'), findsOneWidget);

      // Act - tap retry
      await tester.tap(find.text('Retry'));
      await tester.pump();

      // Assert
      expect(retryPressed, isTrue);
    });

    testWidgets('renders cancel button when onDismiss provided', (
      tester,
    ) async {
      // Arrange
      var dismissPressed = false;

      // Act
      await tester.pumpWidget(
        buildTestWidget(
          errorMessage: 'Error',
          onDismiss: () => dismissPressed = true,
        ),
      );
      await tester.pump();

      // Assert
      expect(find.text('Cancel'), findsOneWidget);

      // Act - tap cancel
      await tester.tap(find.text('Cancel'));
      await tester.pump();

      // Assert
      expect(dismissPressed, isTrue);
    });

    testWidgets('displays error icon', (tester) async {
      // Arrange & Act
      await tester.pumpWidget(buildTestWidget(errorMessage: 'Error'));
      await tester.pump();

      // Assert
      expect(find.byIcon(Icons.error_outline), findsOneWidget);
    });
  });
}
