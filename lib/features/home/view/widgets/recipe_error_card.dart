import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// A card displaying an error state for failed recipe imports.
class RecipeErrorCard extends StatelessWidget {
  /// Creates a new [RecipeErrorCard] instance.
  const RecipeErrorCard({
    required this.errorMessage,
    super.key,
  });

  /// The error message to display.
  final String errorMessage;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    return FCard(
      child: ClipRRect(
        borderRadius: BorderRadius.circular(8),
        child: AspectRatio(
          aspectRatio: 1,
          child: ColoredBox(
            color: theme.colors.destructive.withValues(alpha: 0.1),
            child: Center(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 40,
                    color: theme.colors.destructive,
                  ),
                  const SizedBox(height: 8),
                  Text(
                    l10n.importFailed,
                    style: theme.typography.sm.copyWith(
                      fontWeight: FontWeight.w600,
                      color: theme.colors.foreground,
                    ),
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
