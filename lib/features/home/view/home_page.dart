import 'dart:async';

import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/core/di.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/state/home_controller.dart';
import 'package:recipe_organizer/features/home/view/widgets/empty_state.dart';
import 'package:recipe_organizer/features/home/view/widgets/import_recipe_dialog.dart';
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

  @override
  void initState() {
    super.initState();
    _controller = getIt<HomeController>();
    // Load mock recipes on init for development
    _controller.loadMockRecipes();
  }

  Future<void> _showImportDialog() async {
    await showImportRecipeDialog(
      context: context,
      onImport: (url) {
        Navigator.of(context).pop();
        unawaited(_controller.importRecipe(url));
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final theme = context.theme;

    return FScaffold(
      header: FHeader(
        title: Text(l10n.homeTitle),
      ),
      child: _buildContent(theme),
    );
  }

  Widget _buildContent(FThemeData theme) {
    final isEmpty = _controller.isEmpty.watch(context);
    final recipes = _controller.recipes.watch(context);

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
    // TODO(petrosstergioulas): Navigate to recipe detail page
    debugPrint('Recipe tapped: ${recipe.name}');
  }
}
