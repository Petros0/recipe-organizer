import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';
import 'package:recipe_organizer/features/home/state/pending_import.dart';
import 'package:recipe_organizer/features/home/view/widgets/pending_import_card.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_error_card.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_skeleton_card.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

void main() {
  group('PendingImportCard', () {
    Widget buildTestWidget({
      required PendingImport pendingImport,
      VoidCallback? onRetry,
      VoidCallback? onDismiss,
    }) {
      return MaterialApp(
        localizationsDelegates: AppLocalizations.localizationsDelegates,
        supportedLocales: AppLocalizations.supportedLocales,
        home: FTheme(
          data: FThemes.zinc.dark,
          child: Scaffold(
            body: PendingImportCard(
              pendingImport: pendingImport,
              onRetry: onRetry,
              onDismiss: onDismiss,
            ),
          ),
        ),
      );
    }

    PendingImport createPendingImport({String? errorMessage}) {
      return PendingImport(
        request: RecipeRequest(
          id: 'test-id',
          url: 'https://example.com/recipe',
          status: RecipeRequestStatus.requested,
          userId: 'user-123',
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        userId: 'user-123',
        errorMessage: errorMessage,
      );
    }

    testWidgets('shows skeleton when import is loading', (tester) async {
      // Arrange
      final pendingImport = createPendingImport();

      // Act
      await tester.pumpWidget(buildTestWidget(pendingImport: pendingImport));
      await tester.pump();

      // Assert
      expect(find.byType(RecipeSkeletonCard), findsOneWidget);
      expect(find.byType(RecipeErrorCard), findsNothing);
    });

    testWidgets('shows error card when import has error', (tester) async {
      // Arrange
      final pendingImport = createPendingImport(
        errorMessage: 'Failed to extract recipe',
      );

      // Act
      await tester.pumpWidget(buildTestWidget(pendingImport: pendingImport));
      await tester.pump();

      // Assert
      expect(find.byType(RecipeErrorCard), findsOneWidget);
      expect(find.byType(RecipeSkeletonCard), findsNothing);
      expect(find.text('Failed to extract recipe'), findsOneWidget);
    });

    testWidgets('calls onRetry when retry is pressed', (tester) async {
      // Arrange
      var retryPressed = false;
      final pendingImport = createPendingImport(
        errorMessage: 'Import failed',
      );

      // Act
      await tester.pumpWidget(
        buildTestWidget(
          pendingImport: pendingImport,
          onRetry: () => retryPressed = true,
        ),
      );
      await tester.pump();
      await tester.tap(find.text('Retry'));
      await tester.pump();

      // Assert
      expect(retryPressed, isTrue);
    });

    testWidgets('calls onDismiss when dismiss is pressed', (tester) async {
      // Arrange
      var dismissPressed = false;
      final pendingImport = createPendingImport(
        errorMessage: 'Import failed',
      );

      // Act
      await tester.pumpWidget(
        buildTestWidget(
          pendingImport: pendingImport,
          onDismiss: () => dismissPressed = true,
        ),
      );
      await tester.pump();
      await tester.tap(find.text('Cancel'));
      await tester.pump();

      // Assert
      expect(dismissPressed, isTrue);
    });
  });
}
