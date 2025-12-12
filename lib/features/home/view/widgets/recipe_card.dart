import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';

/// A card widget displaying a recipe with its image and name.
class RecipeCard extends StatelessWidget {
  /// Creates a new [RecipeCard] instance.
  const RecipeCard({
    required this.recipe,
    this.onTap,
    super.key,
  });

  /// The recipe to display.
  final Recipe recipe;

  /// Callback when the card is tapped.
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;

    return GestureDetector(
      onTap: onTap,
      child: FCard(
        child: ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: AspectRatio(
            aspectRatio: 1,
            child: Stack(
              fit: StackFit.expand,
              children: [
                _buildImage(theme),
                _buildGradientOverlay(),
                _buildRecipeInfo(theme),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildImage(FThemeData theme) {
    if (recipe.imageUrl != null && recipe.imageUrl!.isNotEmpty) {
      return Image.network(
        recipe.imageUrl!,
        fit: BoxFit.cover,
        errorBuilder: (context, error, stackTrace) => _buildPlaceholder(theme),
        loadingBuilder: (context, child, loadingProgress) {
          if (loadingProgress == null) return child;
          return _buildLoadingIndicator(theme);
        },
      );
    }
    return _buildPlaceholder(theme);
  }

  Widget _buildPlaceholder(FThemeData theme) {
    return ColoredBox(
      color: theme.colors.secondary,
      child: Center(
        child: Icon(
          Icons.restaurant,
          size: 48,
          color: theme.colors.secondaryForeground,
        ),
      ),
    );
  }

  Widget _buildLoadingIndicator(FThemeData theme) {
    return ColoredBox(
      color: theme.colors.secondary,
      child: Center(
        child: CircularProgressIndicator(
          color: theme.colors.primary,
          strokeWidth: 2,
        ),
      ),
    );
  }

  Widget _buildGradientOverlay() {
    return Positioned.fill(
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
    );
  }

  Widget _buildRecipeInfo(FThemeData theme) {
    return Positioned(
      left: 12,
      right: 12,
      bottom: 12,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            recipe.name,
            style: theme.typography.base.copyWith(
              color: Colors.white,
              fontWeight: FontWeight.w600,
            ),
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
          if (recipe.author != null) ...[
            const SizedBox(height: 4),
            Text(
              recipe.author!,
              style: theme.typography.sm.copyWith(
                color: Colors.white70,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ],
      ),
    );
  }
}
