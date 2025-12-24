import 'package:appwrite/appwrite.dart';
import 'package:recipe_organizer/core/appwrite/appwrite_constants.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';

/// Repository for recipe CRUD operations using Appwrite TablesDB.
class RecipeRepository {
  /// Creates a new [RecipeRepository] instance.
  RecipeRepository(this._tablesDB);

  final TablesDB _tablesDB;

  /// Gets a recipe by its associated recipe request ID.
  ///
  /// Returns null if no recipe is found.
  /// Also populates the sourceUrl from the recipe request.
  Future<Recipe?> getRecipeByRequestId(String requestId) async {
    final response = await _tablesDB.listRows(
      databaseId: AppwriteConstants.databaseId,
      tableId: AppwriteConstants.recipeCollectionId,
      queries: [
        Query.equal('fk_recipe_request', requestId),
        Query.limit(1),
      ],
    );

    if (response.rows.isEmpty) return null;

    var recipe = Recipe.fromRow(response.rows.first);

    // Fetch source URL from recipe request if not already set
    if (recipe.sourceUrl == null && recipe.recipeRequestId != null) {
      final sourceUrl = await _getSourceUrlFromRequest(recipe.recipeRequestId!);
      if (sourceUrl != null) {
        recipe = recipe.copyWith(sourceUrl: sourceUrl);
      }
    }

    return recipe;
  }

  /// Lists all recipes, ordered by creation date (newest first).
  /// Also populates the sourceUrl from the recipe requests.
  Future<List<Recipe>> listRecipes({int limit = 25, int offset = 0}) async {
    final response = await _tablesDB.listRows(
      databaseId: AppwriteConstants.databaseId,
      tableId: AppwriteConstants.recipeCollectionId,
      queries: [
        Query.limit(limit),
        Query.offset(offset),
        Query.orderDesc(r'$createdAt'),
      ],
    );

    final recipes = response.rows.map(Recipe.fromRow).toList();

    // Populate source URLs from recipe requests
    return Future.wait(recipes.map(_populateSourceUrl));
  }

  Future<Recipe> _populateSourceUrl(Recipe recipe) async {
    if (recipe.sourceUrl != null || recipe.recipeRequestId == null) {
      return recipe;
    }

    final sourceUrl = await _getSourceUrlFromRequest(recipe.recipeRequestId!);
    if (sourceUrl != null) {
      return recipe.copyWith(sourceUrl: sourceUrl);
    }
    return recipe;
  }

  Future<String?> _getSourceUrlFromRequest(String requestId) async {
    try {
      final row = await _tablesDB.getRow(
        databaseId: AppwriteConstants.databaseId,
        tableId: AppwriteConstants.recipeRequestCollectionId,
        rowId: requestId,
      );
      return row.data['url'] as String?;
    } on AppwriteException {
      return null;
    }
  }

  /// Gets a recipe by its ID.
  Future<Recipe?> getRecipeById(String id) async {
    try {
      final row = await _tablesDB.getRow(
        databaseId: AppwriteConstants.databaseId,
        tableId: AppwriteConstants.recipeCollectionId,
        rowId: id,
      );
      var recipe = Recipe.fromRow(row);

      // Fetch source URL from recipe request if not already set
      if (recipe.sourceUrl == null && recipe.recipeRequestId != null) {
        final sourceUrl = await _getSourceUrlFromRequest(recipe.recipeRequestId!);
        if (sourceUrl != null) {
          recipe = recipe.copyWith(sourceUrl: sourceUrl);
        }
      }

      return recipe;
    } on AppwriteException {
      return null;
    }
  }

  /// Deletes a recipe by its ID.
  Future<void> deleteRecipe(String id) async {
    await _tablesDB.deleteRow(
      databaseId: AppwriteConstants.databaseId,
      tableId: AppwriteConstants.recipeCollectionId,
      rowId: id,
    );
  }
}
