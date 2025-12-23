import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/view/widgets/ingredients_section.dart';
import 'package:recipe_organizer/features/home/view/widgets/instructions_section.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_metadata_row.dart';
import 'package:recipe_organizer/features/home/view/widgets/source_attribution.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// Page for previewing an extracted recipe before saving.
class RecipePreviewPage extends StatefulWidget {
  /// Creates a new [RecipePreviewPage] instance.
  const RecipePreviewPage({
    required this.recipe,
    required this.onSaveOrDelete,
    this.onCancel,
    this.isProcessing = false,
    this.isEditable = false,
    super.key,
  });

  /// The recipe to preview.
  final Recipe recipe;

  /// Callback when delete is pressed.
  final VoidCallback onSaveOrDelete;

  /// Callback when cancel is pressed.
  final VoidCallback? onCancel;

  /// Whether delete is in progress.
  final bool isProcessing;

  /// Whether fields are editable (for partial data).
  final bool isEditable;

  @override
  State<RecipePreviewPage> createState() => _RecipePreviewPageState();
}

class _RecipePreviewPageState extends State<RecipePreviewPage> {
  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    return Scaffold(
      backgroundColor: theme.colors.background,
      body: CustomScrollView(
        slivers: [
          _buildSliverAppBar(theme, l10n),
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildTitle(theme),
                  if (widget.recipe.description != null) ...[
                    const SizedBox(height: 8),
                    _buildDescription(theme),
                  ],
                  const SizedBox(height: 16),
                  RecipeMetadataRow(
                    prepTime: widget.recipe.prepTime,
                    cookTime: widget.recipe.cookTime,
                    servings: widget.recipe.servings,
                  ),
                  const SizedBox(height: 24),
                  IngredientsSection(ingredients: widget.recipe.ingredients),
                  const SizedBox(height: 24),
                  InstructionsSection(instructions: widget.recipe.instructions),
                  const SizedBox(height: 24),
                  if (widget.recipe.sourceUrl != null || widget.recipe.authorName != null)
                    SourceAttribution(
                      sourceUrl: widget.recipe.sourceUrl,
                      authorName: widget.recipe.authorName,
                      authorUrl: widget.recipe.authorUrl,
                    ),
                  const SizedBox(height: 100), // Space for bottom button
                ],
              ),
            ),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomBar(theme, l10n),
    );
  }

  Widget _buildSliverAppBar(FThemeData theme, AppLocalizations l10n) {
    final imageUrl = widget.recipe.imageUrl;

    return SliverAppBar(
      expandedHeight: 300,
      pinned: true,
      backgroundColor: theme.colors.background,
      foregroundColor: theme.colors.foreground,
      leading: IconButton(
        icon: Container(
          padding: const EdgeInsets.all(8),
          decoration: BoxDecoration(
            color: theme.colors.background.withValues(alpha: 0.8),
            shape: BoxShape.circle,
          ),
          child: Icon(
            Icons.arrow_back,
            color: theme.colors.foreground,
          ),
        ),
        onPressed: widget.onCancel ?? () => Navigator.of(context).pop(),
      ),
      flexibleSpace: FlexibleSpaceBar(
        background: imageUrl != null
            ? Image.network(
                imageUrl,
                fit: BoxFit.cover,
                errorBuilder: (_, _, _) => _buildPlaceholderImage(theme),
              )
            : _buildPlaceholderImage(theme),
      ),
    );
  }

  Widget _buildPlaceholderImage(FThemeData theme) {
    return ColoredBox(
      color: theme.colors.secondary,
      child: Center(
        child: Icon(
          Icons.restaurant,
          size: 64,
          color: theme.colors.secondaryForeground,
        ),
      ),
    );
  }

  Widget _buildTitle(FThemeData theme) {
    return Text(
      widget.recipe.name,
      style: theme.typography.xl2.copyWith(
        fontWeight: FontWeight.bold,
        color: theme.colors.foreground,
      ),
    );
  }

  Widget _buildDescription(FThemeData theme) {
    return Text(
      widget.recipe.description!,
      style: theme.typography.base.copyWith(
        color: theme.colors.mutedForeground,
        height: 1.5,
      ),
    );
  }

  Widget _buildBottomBar(FThemeData theme, AppLocalizations l10n) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: theme.colors.background,
        border: Border(
          top: BorderSide(color: theme.colors.border),
        ),
      ),
      child: SafeArea(
        child: FButton(
          style: FButtonStyle.destructive(),
          onPress: widget.isProcessing ? null : widget.onSaveOrDelete,
          child: widget.isProcessing
              ? SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    color: theme.colors.primaryForeground,
                  ),
                )
              : Text(l10n.deleteRecipe),
        ),
      ),
    );
  }
}
