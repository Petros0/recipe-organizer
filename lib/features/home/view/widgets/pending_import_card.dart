import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/state/pending_import.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_error_card.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_skeleton_card.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// A card that displays a pending import, switching between skeleton and error.
class PendingImportCard extends StatelessWidget {
  /// Creates a new [PendingImportCard] instance.
  const PendingImportCard({
    required this.pendingImport,
    this.onRetry,
    this.onDismiss,
    super.key,
  });

  /// The pending import to display.
  final PendingImport pendingImport;

  /// Callback when retry is pressed on an error.
  final VoidCallback? onRetry;

  /// Callback when dismiss is pressed on an error.
  final VoidCallback? onDismiss;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    final child = pendingImport.hasError
        ? RecipeErrorCard(errorMessage: pendingImport.errorMessage ?? 'Import failed')
        : const RecipeSkeletonCard();

    // Only show menu for error state
    if (!pendingImport.hasError) {
      return child;
    }

    return Stack(
      children: [
        child,
        Positioned(
          top: 8,
          right: 8,
          child: PopupMenuButton<String>(
            icon: Container(
              padding: const EdgeInsets.all(4),
              decoration: BoxDecoration(
                color: theme.colors.background.withValues(alpha: 0.8),
                shape: BoxShape.circle,
              ),
              child: Icon(
                Icons.more_vert,
                color: theme.colors.foreground,
                size: 18,
              ),
            ),
            padding: EdgeInsets.zero,
            onSelected: (value) {
              if (value == 'retry') {
                onRetry?.call();
              } else if (value == 'dismiss') {
                onDismiss?.call();
              }
            },
            itemBuilder: (context) => [
              if (onRetry != null)
                PopupMenuItem(
                  value: 'retry',
                  child: Row(
                    children: [
                      Icon(Icons.refresh, size: 18, color: theme.colors.foreground),
                      const SizedBox(width: 8),
                      Text(l10n.retryButton),
                    ],
                  ),
                ),
              if (onDismiss != null)
                PopupMenuItem(
                  value: 'dismiss',
                  child: Row(
                    children: [
                      Icon(Icons.close, size: 18, color: theme.colors.destructive),
                      const SizedBox(width: 8),
                      Text(l10n.cancelButton),
                    ],
                  ),
                ),
            ],
          ),
        ),
      ],
    );
  }
}
