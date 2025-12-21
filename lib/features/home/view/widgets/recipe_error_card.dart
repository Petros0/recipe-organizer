import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// A card displaying an error state for failed recipe imports.
class RecipeErrorCard extends StatelessWidget {
  /// Creates a new [RecipeErrorCard] instance.
  const RecipeErrorCard({
    required this.errorMessage,
    this.onRetry,
    this.onDismiss,
    super.key,
  });

  /// The error message to display.
  final String errorMessage;

  /// Callback when retry is pressed.
  final VoidCallback? onRetry;

  /// Callback when dismiss is pressed.
  final VoidCallback? onDismiss;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    return FCard(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // Error icon area
          Container(
            height: 120,
            width: double.infinity,
            decoration: BoxDecoration(
              color: theme.colors.destructive.withValues(alpha: 0.1),
              borderRadius: BorderRadius.circular(8),
            ),
            child: Center(
              child: Icon(
                Icons.error_outline,
                size: 48,
                color: theme.colors.destructive,
              ),
            ),
          ),
          const SizedBox(height: 16),
          // Error title
          Text(
            l10n.importFailed,
            style: theme.typography.lg.copyWith(
              fontWeight: FontWeight.w600,
              color: theme.colors.foreground,
            ),
          ),
          const SizedBox(height: 8),
          // Error message
          Text(
            errorMessage,
            style: theme.typography.sm.copyWith(
              color: theme.colors.mutedForeground,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 16),
          // Action buttons
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              if (onDismiss != null) ...[
                FButton(
                  style: FButtonStyle.outline(),
                  onPress: onDismiss,
                  child: Text(l10n.cancelButton),
                ),
                const SizedBox(width: 12),
              ],
              if (onRetry != null)
                FButton(
                  onPress: onRetry,
                  child: Text(l10n.retryButton),
                ),
            ],
          ),
        ],
      ),
    );
  }
}
