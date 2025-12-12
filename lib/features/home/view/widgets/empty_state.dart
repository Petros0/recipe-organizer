import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// Widget displayed when the recipe list is empty.
class EmptyState extends StatelessWidget {
  /// Creates a new [EmptyState] instance.
  const EmptyState({
    required this.onImportPressed,
    super.key,
  });

  /// Callback when the import button is pressed.
  final VoidCallback onImportPressed;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            _buildIcon(theme),
            const SizedBox(height: 24),
            _buildTitle(theme, l10n),
            const SizedBox(height: 12),
            _buildDescription(theme, l10n),
            const SizedBox(height: 32),
            _buildImportButton(l10n),
          ],
        ),
      ),
    );
  }

  Widget _buildIcon(FThemeData theme) {
    return Container(
      width: 120,
      height: 120,
      decoration: BoxDecoration(
        color: theme.colors.secondary,
        shape: BoxShape.circle,
      ),
      child: Icon(
        Icons.menu_book_rounded,
        size: 56,
        color: theme.colors.mutedForeground,
      ),
    );
  }

  Widget _buildTitle(FThemeData theme, AppLocalizations l10n) {
    return Text(
      l10n.emptyStateTitle,
      style: theme.typography.xl2.copyWith(
        fontWeight: FontWeight.w600,
        color: theme.colors.foreground,
      ),
      textAlign: TextAlign.center,
    );
  }

  Widget _buildDescription(FThemeData theme, AppLocalizations l10n) {
    return Text(
      l10n.emptyStateDescription,
      style: theme.typography.base.copyWith(
        color: theme.colors.mutedForeground,
      ),
      textAlign: TextAlign.center,
    );
  }

  Widget _buildImportButton(AppLocalizations l10n) {
    return FButton(
      onPress: onImportPressed,
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          const Icon(Icons.add),
          const SizedBox(width: 8),
          Text(l10n.importRecipeButton),
        ],
      ),
    );
  }
}
