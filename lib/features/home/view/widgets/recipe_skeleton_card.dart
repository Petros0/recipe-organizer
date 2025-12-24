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
      child: ClipRRect(
        borderRadius: BorderRadius.circular(8),
        child: AspectRatio(
          aspectRatio: 1,
          child: Stack(
            fit: StackFit.expand,
            children: [
              // Shimmer background
              Shimmer.fromColors(
                baseColor: theme.colors.secondary,
                highlightColor: theme.colors.background,
                child: ColoredBox(color: theme.colors.secondary),
              ),
              // Gradient overlay (same as RecipeCard)
              Positioned.fill(
                child: DecoratedBox(
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      begin: Alignment.topCenter,
                      end: Alignment.bottomCenter,
                      colors: [
                        Colors.transparent,
                        Colors.black.withValues(alpha: 0.7),
                      ],
                      stops: const [0.5, 1.0],
                    ),
                  ),
                ),
              ),
              // Loading info at bottom
              Positioned(
                left: 12,
                right: 12,
                bottom: 12,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // Title skeleton
                    Shimmer.fromColors(
                      baseColor: Colors.white24,
                      highlightColor: Colors.white54,
                      child: Container(
                        height: 16,
                        width: double.infinity,
                        decoration: BoxDecoration(
                          color: Colors.white24,
                          borderRadius: BorderRadius.circular(4),
                        ),
                      ),
                    ),
                    const SizedBox(height: 8),
                    // Status row
                    Row(
                      children: [
                        const SizedBox(
                          width: 12,
                          height: 12,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            color: Colors.white70,
                          ),
                        ),
                        const SizedBox(width: 8),
                        Flexible(
                          child: Text(
                            showLlmMessage ? l10n.extractingRecipeLlm : l10n.extractingRecipe,
                            style: theme.typography.xs.copyWith(
                              color: Colors.white70,
                            ),
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
