import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_skeleton_card.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

void main() {
  group('RecipeSkeletonCard', () {
    Widget buildTestWidget({bool showLlmMessage = false}) {
      return MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: FTheme(
          data: FThemes.zinc.dark,
          child: Scaffold(
            body: RecipeSkeletonCard(showLlmMessage: showLlmMessage),
          ),
        ),
      );
    }

    testWidgets('renders shimmer skeleton elements', (tester) async {
      // Arrange & Act
      await tester.pumpWidget(buildTestWidget());
      await tester.pump();

      // Assert
      expect(find.byType(RecipeSkeletonCard), findsOneWidget);
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('shows extracting message by default', (tester) async {
      // Arrange & Act
      await tester.pumpWidget(buildTestWidget());
      await tester.pump();

      // Assert
      expect(find.text('Extracting recipe...'), findsOneWidget);
    });

    testWidgets('shows LLM message when showLlmMessage is true', (
      tester,
    ) async {
      // Arrange & Act
      await tester.pumpWidget(buildTestWidget(showLlmMessage: true));
      await tester.pump();

      // Assert
      expect(find.text('Analyzing page content...'), findsOneWidget);
    });
  });
}
