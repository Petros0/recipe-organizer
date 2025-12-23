import 'dart:async';

import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';
import 'package:recipe_organizer/features/home/service/recipe_import_service.dart';
import 'package:recipe_organizer/features/home/state/import_state.dart';
import 'package:signals/signals.dart';

/// Controller for managing the home page state using Signals.
class HomeController {
  /// Creates a new [HomeController] instance.
  HomeController({required RecipeImportService importService}) : _importService = importService {
    isEmpty = computed(() => _recipes.value.isEmpty);
    isExtracting = computed(
      () => _importState.value == ImportState.submitting || _importState.value == ImportState.extracting,
    );
  }

  final RecipeImportService _importService;
  StreamSubscription<RecipeRequest>? _requestSubscription;

  final Signal<List<Recipe>> _recipes = signal<List<Recipe>>([]);
  final Signal<bool> _isLoading = signal<bool>(false);
  final Signal<String?> _error = signal<String?>(null);
  final Signal<ImportState> _importState = signal<ImportState>(
    ImportState.idle,
  );
  final Signal<RecipeRequest?> _activeRequest = signal<RecipeRequest?>(null);
  final Signal<Recipe?> _previewRecipe = signal<Recipe?>(null);

  /// The list of recipes.
  ReadonlySignal<List<Recipe>> get recipes => _recipes;

  /// Whether recipes are currently being loaded.
  ReadonlySignal<bool> get isLoading => _isLoading;

  /// The current error message, if any.
  ReadonlySignal<String?> get error => _error;

  /// Whether the recipe list is empty.
  late final Computed<bool> isEmpty;

  /// Whether extraction is in progress.
  late final Computed<bool> isExtracting;

  /// Current import state.
  ReadonlySignal<ImportState> get importState => _importState;

  /// The active recipe request being processed.
  ReadonlySignal<RecipeRequest?> get activeRequest => _activeRequest;

  /// The recipe available for preview.
  ReadonlySignal<Recipe?> get previewRecipe => _previewRecipe;

  /// Loads recipes from the backend.
  Future<void> loadRecipes() async {
    _isLoading.value = true;
    _error.value = null;

    try {
      final recipes = await _importService.listRecipes();
      _recipes.value = recipes;
    } on Exception catch (e) {
      _error.value = 'Failed to load recipes: $e';
    } finally {
      _isLoading.value = false;
    }
  }

  /// Deletes a recipe by its ID.
  Future<void> deleteRecipe(String recipeId) async {
    try {
      await _importService.deleteRecipe(recipeId);
      _recipes.value = _recipes.value.where((r) => r.id != recipeId).toList();
    } on Exception catch (e) {
      _error.value = 'Failed to delete recipe: $e';
    }
  }

  /// Clears all recipes (for testing empty state).
  void clearRecipes() {
    _recipes.value = [];
  }

  /// Imports a recipe from a URL.
  ///
  /// The [userId] identifies the authenticated user making the request.
  Future<void> importRecipe({
    required String url,
    required String userId,
  }) async {
    if (url.isEmpty) {
      _error.value = 'Please enter a valid URL';
      return;
    }

    _importState.value = ImportState.submitting;
    _error.value = null;
    _lastImportUserId = userId;

    try {
      final request = await _importService.importFromUrl(
        url: url,
        userId: userId,
      );
      _activeRequest.value = request;
      _importState.value = ImportState.extracting;

      _subscribeToRequest(request.id);
    } on Exception catch (e) {
      _error.value = e.toString();
      _importState.value = ImportState.error;
    }
  }

  String? _lastImportUserId;

  void _subscribeToRequest(String requestId) {
    _requestSubscription?.cancel();
    _requestSubscription = _importService
        .subscribeToRequest(requestId)
        .listen(
          (request) {
            _activeRequest.value = request;

            switch (request.status) {
              case RecipeRequestStatus.completed:
                unawaited(_fetchExtractedRecipe(requestId));
              case RecipeRequestStatus.failed:
                _importState.value = ImportState.error;
                _error.value = 'Recipe extraction failed';
              case RecipeRequestStatus.requested:
              case RecipeRequestStatus.inProgress:
                // Stay in extracting state
                break;
            }
          },
          onError: (Object error) {
            _error.value = error.toString();
            _importState.value = ImportState.error;
          },
        );
  }

  Future<void> _fetchExtractedRecipe(String requestId) async {
    try {
      final recipe = await _importService.getExtractedRecipe(requestId);
      if (recipe != null) {
        _previewRecipe.value = recipe;
        _importState.value = ImportState.preview;
      } else {
        _error.value = 'Recipe not found';
        _importState.value = ImportState.error;
      }
    } on Exception catch (e) {
      _error.value = e.toString();
      _importState.value = ImportState.error;
    }
  }

  /// Retries the last failed import.
  Future<void> retryImport() async {
    final request = _activeRequest.value;
    final userId = _lastImportUserId;
    if (request != null && userId != null) {
      await importRecipe(url: request.url, userId: userId);
    }
  }

  /// Saves the preview recipe to the collection.
  Future<void> saveRecipe() async {
    final recipe = _previewRecipe.value;
    if (recipe == null) return;

    _importState.value = ImportState.saving;

    try {
      // Recipe is already saved by backend, just add to local list
      _recipes.value = [recipe, ..._recipes.value];
      _resetImportState();
    } on Exception catch (e) {
      _error.value = e.toString();
      _importState.value = ImportState.error;
    }
  }

  /// Cancels the current import and resets state.
  void cancelImport() {
    _resetImportState();
  }

  void _resetImportState() {
    _requestSubscription?.cancel();
    _requestSubscription = null;
    _importState.value = ImportState.idle;
    _activeRequest.value = null;
    _previewRecipe.value = null;
    _error.value = null;
  }

  /// Clears the current error.
  void clearError() {
    _error.value = null;
  }

  /// Disposes resources.
  void dispose() {
    _requestSubscription?.cancel();
  }
}
