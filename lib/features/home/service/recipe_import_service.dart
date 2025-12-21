import 'dart:async';

import 'package:recipe_organizer/core/appwrite/realtime_service.dart';
import 'package:recipe_organizer/features/home/data/recipe_repository.dart';
import 'package:recipe_organizer/features/home/data/recipe_request_repository.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';

/// Service for managing recipe import operations.
class RecipeImportService {
  /// Creates a new [RecipeImportService] instance.
  RecipeImportService({
    required RecipeRequestRepository requestRepository,
    required RecipeRepository recipeRepository,
    required RealtimeService realtimeService,
  }) : _requestRepository = requestRepository,
       _recipeRepository = recipeRepository,
       _realtimeService = realtimeService;

  final RecipeRequestRepository _requestRepository;
  final RecipeRepository _recipeRepository;
  final RealtimeService _realtimeService;

  /// Initiates a recipe import from the given URL.
  ///
  /// The [userId] identifies the authenticated user making the request.
  /// Returns the created [RecipeRequest].
  Future<RecipeRequest> importFromUrl({
    required String url,
    required String userId,
  }) async {
    return _requestRepository.createRecipeRequest(url: url, userId: userId);
  }

  /// Subscribes to status updates for a recipe request.
  ///
  /// Returns a stream that emits [RecipeRequest] updates.
  Stream<RecipeRequest> subscribeToRequest(String requestId) {
    return _realtimeService.subscribeToRecipeRequest(requestId).map(RecipeRequest.fromMap);
  }

  /// Gets the extracted recipe for a completed request.
  Future<Recipe?> getExtractedRecipe(String requestId) async {
    return _recipeRepository.getRecipeByRequestId(requestId);
  }

  /// Lists all saved recipes.
  Future<List<Recipe>> listRecipes({int limit = 25, int offset = 0}) async {
    return _recipeRepository.listRecipes(limit: limit, offset: offset);
  }
}
