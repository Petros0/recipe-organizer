import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_error_card.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

void main() {
  group('RecipeErrorCard', () {
    Widget buildTestWidget({
      required String errorMessage,
    }) {
      return MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: FTheme(
          data: FThemes.zinc.dark,
          child: Scaffold(
            body: RecipeErrorCard(
              errorMessage: errorMessage,
            ),
          ),
        ),
      );
    }

    testWidgets('renders error message and import failed text', (tester) async {
      // Arrange
      const errorMessage = 'Test error message';

      // Act
      await tester.pumpWidget(buildTestWidget(errorMessage: errorMessage));
      await tester.pump();

      // Assert
      expect(find.text('Import failed'), findsOneWidget);
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
