import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';
import 'package:shimmer/shimmer.dart';

/// A skeleton loading card with shimmer effect for recipe import.
class RecipeSkeletonCard extends StatelessWidget {
  /// Creates a new [RecipeSkeletonCard] instance.
  const RecipeSkeletonCard({
    this.showLlmMessage = false,
    super.key,
  });

  /// Whether to show LLM extraction message (longer processing).
  final bool showLlmMessage;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    return FCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Hero image skeleton
          Shimmer.fromColors(
            baseColor: theme.colors.secondary,
            highlightColor: theme.colors.background,
            child: Container(
              height: 200,
              decoration: BoxDecoration(
                color: theme.colors.secondary,
                borderRadius: BorderRadius.circular(8),
              ),
            ),
          ),
          const SizedBox(height: 16),
          // Title skeleton
          Shimmer.fromColors(
            baseColor: theme.colors.secondary,
            highlightColor: theme.colors.background,
            child: Container(
              height: 24,
              width: double.infinity,
              decoration: BoxDecoration(
                color: theme.colors.secondary,
                borderRadius: BorderRadius.circular(4),
              ),
            ),
          ),
          const SizedBox(height: 12),
          // Metadata skeleton row
          Row(
            children: [
              _buildSkeletonChip(theme, 80),
              const SizedBox(width: 8),
              _buildSkeletonChip(theme, 80),
              const SizedBox(width: 8),
              _buildSkeletonChip(theme, 60),
            ],
          ),
          const SizedBox(height: 16),
          // Status message
          Center(
            child: Text(
              showLlmMessage ? l10n.extractingRecipeLlm : l10n.extractingRecipe,
              style: theme.typography.sm.copyWith(
                color: theme.colors.mutedForeground,
              ),
            ),
          ),
          const SizedBox(height: 8),
          // Loading indicator
          Center(
            child: SizedBox(
              width: 20,
              height: 20,
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: theme.colors.primary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSkeletonChip(FThemeData theme, double width) {
    return Shimmer.fromColors(
      baseColor: theme.colors.secondary,
      highlightColor: theme.colors.background,
      child: Container(
        height: 28,
        width: width,
        decoration: BoxDecoration(
          color: theme.colors.secondary,
          borderRadius: BorderRadius.circular(14),
        ),
      ),
    );
  }
}
