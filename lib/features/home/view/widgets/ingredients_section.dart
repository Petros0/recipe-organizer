import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// Widget displaying the ingredients list section.
class IngredientsSection extends StatelessWidget {
  /// Creates a new [IngredientsSection] instance.
  const IngredientsSection({
    required this.ingredients,
    super.key,
  });

  /// List of ingredients.
  final List<String> ingredients;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    if (ingredients.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.ingredients,
          style: theme.typography.lg.copyWith(
            fontWeight: FontWeight.w600,
            color: theme.colors.foreground,
          ),
        ),
        const SizedBox(height: 12),
        ...ingredients.map(
          (ingredient) => Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  width: 6,
                  height: 6,
                  margin: const EdgeInsets.only(top: 6, right: 12),
                  decoration: BoxDecoration(
                    color: theme.colors.primary,
                    shape: BoxShape.circle,
                  ),
                ),
                Expanded(
                  child: Text(
                    ingredient,
                    style: theme.typography.base.copyWith(
                      color: theme.colors.foreground,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
