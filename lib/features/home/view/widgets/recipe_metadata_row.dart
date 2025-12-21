import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// Widget displaying recipe metadata (prep time, cook time, servings).
class RecipeMetadataRow extends StatelessWidget {
  /// Creates a new [RecipeMetadataRow] instance.
  const RecipeMetadataRow({
    this.prepTime,
    this.cookTime,
    this.servings,
    super.key,
  });

  /// Preparation time in ISO 8601 duration format.
  final String? prepTime;

  /// Cooking time in ISO 8601 duration format.
  final String? cookTime;

  /// Servings/yield string.
  final String? servings;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;
    final items = <Widget>[];

    if (prepTime != null) {
      items.add(
        _buildChip(
          theme,
          Icons.timer_outlined,
          '${l10n.prepTime}: ${_formatDuration(prepTime!)}',
        ),
      );
    }

    if (cookTime != null) {
      items.add(
        _buildChip(
          theme,
          Icons.local_fire_department_outlined,
          '${l10n.cookTime}: ${_formatDuration(cookTime!)}',
        ),
      );
    }

    if (servings != null) {
      items.add(
        _buildChip(
          theme,
          Icons.restaurant_outlined,
          servings!,
        ),
      );
    }

    if (items.isEmpty) return const SizedBox.shrink();

    return Wrap(
      spacing: 8,
      runSpacing: 8,
      children: items,
    );
  }

  Widget _buildChip(FThemeData theme, IconData icon, String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: theme.colors.secondary,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            icon,
            size: 16,
            color: theme.colors.mutedForeground,
          ),
          const SizedBox(width: 4),
          Text(
            label,
            style: theme.typography.sm.copyWith(
              color: theme.colors.foreground,
            ),
          ),
        ],
      ),
    );
  }

  /// Formats ISO 8601 duration (e.g., PT15M) to human-readable format.
  String _formatDuration(String isoDuration) {
    final regex = RegExp(r'PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?');
    final match = regex.firstMatch(isoDuration);

    if (match == null) return isoDuration;

    final hours = match.group(1);
    final minutes = match.group(2);

    final parts = <String>[];
    if (hours != null) parts.add('${hours}h');
    if (minutes != null) parts.add('${minutes}m');

    return parts.isEmpty ? isoDuration : parts.join(' ');
  }
}
