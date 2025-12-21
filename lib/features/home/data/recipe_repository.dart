import 'package:appwrite/appwrite.dart';
import 'package:recipe_organizer/core/appwrite/appwrite_constants.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';

/// Repository for recipe CRUD operations using Appwrite Databases.
class RecipeRepository {
  /// Creates a new [RecipeRepository] instance.
  RecipeRepository(this._databases);

  final Databases _databases;

  /// Gets a recipe by its associated recipe request ID.
  ///
  /// Returns null if no recipe is found.
  Future<Recipe?> getRecipeByRequestId(String requestId) async {
    final response = await _databases.listDocuments(
      databaseId: AppwriteConstants.databaseId,
      collectionId: AppwriteConstants.recipeCollectionId,
      queries: [
        Query.equal('fk_recipe_request', requestId),
        Query.limit(1),
      ],
    );

    if (response.documents.isEmpty) return null;
    return Recipe.fromDocument(response.documents.first);
  }

  /// Lists all recipes, ordered by creation date (newest first).
  Future<List<Recipe>> listRecipes({int limit = 25, int offset = 0}) async {
    final response = await _databases.listDocuments(
      databaseId: AppwriteConstants.databaseId,
      collectionId: AppwriteConstants.recipeCollectionId,
      queries: [
        Query.limit(limit),
        Query.offset(offset),
        Query.orderDesc(r'$createdAt'),
      ],
    );

    return response.documents.map(Recipe.fromDocument).toList();
  }

  /// Gets a recipe by its ID.
  Future<Recipe?> getRecipeById(String id) async {
    try {
      final doc = await _databases.getDocument(
        databaseId: AppwriteConstants.databaseId,
        collectionId: AppwriteConstants.recipeCollectionId,
        documentId: id,
      );
      return Recipe.fromDocument(doc);
    } on AppwriteException {
      return null;
    }
  }
}
