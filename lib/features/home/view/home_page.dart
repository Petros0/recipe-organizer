import 'dart:async';

import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/core/di.dart';
import 'package:recipe_organizer/features/auth/state/auth_controller.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/state/home_controller.dart';
import 'package:recipe_organizer/features/home/state/import_state.dart';
import 'package:recipe_organizer/features/home/state/pending_import.dart';
import 'package:recipe_organizer/features/home/view/recipe_preview_page.dart';
import 'package:recipe_organizer/features/home/view/widgets/empty_state.dart';
import 'package:recipe_organizer/features/home/view/widgets/import_recipe_dialog.dart';
import 'package:recipe_organizer/features/home/view/widgets/pending_import_card.dart';
import 'package:recipe_organizer/features/home/view/widgets/recipe_card.dart';
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

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final theme = context.theme;
    final importState = _controller.importState.watch(context);

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
    final pendingImports = _controller.pendingImports.watch(context);

    if (isEmpty && pendingImports.isEmpty) {
      return EmptyState(onImportPressed: _showImportDialog);
    }

    return Stack(
      children: [
        _buildRecipeList(recipes, pendingImports),
        _buildFloatingActionButton(theme),
      ],
    );
  }

  Widget _buildRecipeList(List<Recipe> recipes, List<PendingImport> pending) {
    final totalItems = pending.length + recipes.length;

    return Padding(
      padding: const EdgeInsets.all(16),
      child: GridView.builder(
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          crossAxisSpacing: 16,
          mainAxisSpacing: 16,
        ),
        itemCount: totalItems,
        itemBuilder: (context, index) {
          // Show pending imports first
          if (index < pending.length) {
            final pendingImport = pending[index];
            return PendingImportCard(
              pendingImport: pendingImport,
              onRetry: () => _controller.retryPendingImport(
                pendingImport.request.id,
              ),
              onDismiss: () => _controller.dismissPendingImport(
                pendingImport.request.id,
              ),
            );
          }

          // Then show saved recipes
          final recipeIndex = index - pending.length;
          final recipe = recipes[recipeIndex];
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
        builder: (_) => RecipePreviewPage(
          recipe: recipe,
          onSaveOrDelete: () async {
            await _controller.deleteRecipe(recipe.id);
            if (mounted) Navigator.of(context).pop();
          },
        ),
      ),
    );
  }
}
