import 'dart:async';

import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/core/di.dart';
import 'package:recipe_organizer/features/auth/state/auth_controller.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/state/home_controller.dart';
import 'package:recipe_organizer/features/home/state/import_state.dart';
import 'package:recipe_organizer/features/home/view/recipe_preview_page.dart';
import 'package:recipe_organizer/features/home/view/widgets/empty_state.dart';
import 'package:recipe_organizer/features/home/view/widgets/import_recipe_dialog.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_card.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_error_card.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_skeleton_card.dart';
import 'package:recipe_organizer/l10n/l10n.dart';
import 'package:signals/signals_flutter.dart';

/// The home page displaying a grid of imported recipes.
class HomePage extends StatefulWidget {
  /// Creates a new [HomePage] instance.
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> with SignalsMixin {
  late final HomeController _controller;
  late final AuthController _authController;

  @override
  void initState() {
    super.initState();
    _controller = getIt<HomeController>();
    _authController = getIt<AuthController>();
    // Load recipes from the database
    _controller.loadRecipes();
  }

  Future<void> _showImportDialog() async {
    await showImportRecipeDialog(
      context: context,
      onImport: (url) {
        Navigator.of(context).pop();
        final userId = _authController.currentUser.value?.$id;
        if (userId != null) {
          unawaited(_controller.importRecipe(url: url, userId: userId));
        }
      },
    );
  }

  void _navigateToPreview(Recipe recipe) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (context) => RecipePreviewPage(
          recipe: recipe,
          onSaveOrDelete: () {
            _controller.saveRecipe();
            Navigator.of(context).pop();
          },
          onCancel: () {
            _controller.cancelImport();
            Navigator.of(context).pop();
          },
          isProcessing: _controller.importState.value == ImportState.saving,
          isNewRecipe: true,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final theme = context.theme;
    final importState = _controller.importState.watch(context);
    final previewRecipe = _controller.previewRecipe.watch(context);

    // Navigate to preview when recipe is ready
    if (importState == ImportState.preview && previewRecipe != null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _navigateToPreview(previewRecipe);
      });
    }

    return FScaffold(
      header: FHeader(
        title: Text(l10n.homeTitle),
      ),
      child: _buildContent(theme, importState),
    );
  }

  Widget _buildContent(FThemeData theme, ImportState importState) {
    final isEmpty = _controller.isEmpty.watch(context);
    final recipes = _controller.recipes.watch(context);
    final error = _controller.error.watch(context);
    final isExtracting = _controller.isExtracting.watch(context);

    // Show error card if import failed
    if (importState == ImportState.error && error != null) {
      return Stack(
        children: [
          if (recipes.isNotEmpty) _buildRecipeGrid(recipes) else EmptyState(onImportPressed: _showImportDialog),
          Positioned(
            top: 16,
            left: 16,
            right: 16,
            child: RecipeErrorCard(
              errorMessage: error,
              onRetry: _controller.retryImport,
              onDismiss: _controller.clearError,
            ),
          ),
        ],
      );
    }

    // Show skeleton during extraction
    if (isExtracting) {
      return Stack(
        children: [
          if (recipes.isNotEmpty) _buildRecipeGrid(recipes) else const Center(child: RecipeSkeletonCard()),
          if (recipes.isNotEmpty)
            const Positioned(
              top: 16,
              left: 16,
              right: 16,
              child: RecipeSkeletonCard(),
            ),
          _buildFloatingActionButton(theme),
        ],
      );
    }

    if (isEmpty) {
      return EmptyState(onImportPressed: _showImportDialog);
    }

    return Stack(
      children: [
        _buildRecipeGrid(recipes),
        _buildFloatingActionButton(theme),
      ],
    );
  }

  Widget _buildRecipeGrid(List<Recipe> recipes) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: GridView.builder(
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          crossAxisSpacing: 16,
          mainAxisSpacing: 16,
        ),
        itemCount: recipes.length,
        itemBuilder: (context, index) {
          final recipe = recipes[index];
          return RecipeCard(
            recipe: recipe,
            onTap: () => _onRecipeTap(recipe),
          );
        },
      ),
    );
  }

  Widget _buildFloatingActionButton(FThemeData theme) {
    return Positioned(
      right: 16,
      bottom: 16,
      child: FloatingActionButton(
        onPressed: _showImportDialog,
        backgroundColor: theme.colors.primary,
        foregroundColor: theme.colors.primaryForeground,
        child: const Icon(Icons.add),
      ),
    );
  }

  void _onRecipeTap(Recipe recipe) {
    Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (context) => RecipePreviewPage(
          recipe: recipe,
          onSaveOrDelete: () async {
            await _controller.deleteRecipe(recipe.id);
            if (context.mounted) Navigator.of(context).pop();
          },
        ),
      ),
    );
  }
}
