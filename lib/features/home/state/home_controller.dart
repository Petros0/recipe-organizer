import 'dart:async';

import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';
import 'package:recipe_organizer/features/home/service/recipe_import_service.dart';
import 'package:recipe_organizer/features/home/state/import_state.dart';
import 'package:recipe_organizer/features/home/state/pending_import.dart';
import 'package:signals/signals.dart';

/// Controller for managing the home page state using Signals.
class HomeController {
  /// Creates a new [HomeController] instance.
  HomeController({required RecipeImportService importService}) : _importService = importService {
    isEmpty = computed(() => _recipes.value.isEmpty && _pendingImports.value.isEmpty);
    isExtracting = computed(
      () => _importState.value == ImportState.submitting || _importState.value == ImportState.extracting,
    );
  }

  final RecipeImportService _importService;
  StreamSubscription<RecipeRequest>? _requestSubscription;
  final Map<String, StreamSubscription<RecipeRequest>> _pendingSubscriptions = {};

  final Signal<List<Recipe>> _recipes = signal<List<Recipe>>([]);
  final Signal<bool> _isLoading = signal<bool>(false);
  final Signal<String?> _error = signal<String?>(null);
  final Signal<ImportState> _importState = signal<ImportState>(
    ImportState.idle,
  );
  final Signal<RecipeRequest?> _activeRequest = signal<RecipeRequest?>(null);
  final Signal<Recipe?> _previewRecipe = signal<Recipe?>(null);
  final Signal<List<PendingImport>> _pendingImports = signal<List<PendingImport>>([]);

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

  /// List of pending imports (loading or failed).
  ReadonlySignal<List<PendingImport>> get pendingImports => _pendingImports;

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
  /// This is a non-blocking operation - the import is added to pending list.
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

      // Add to pending imports list
      final pendingImport = PendingImport(request: request, userId: userId);
      _pendingImports.value = [pendingImport, ..._pendingImports.value];

      // Also set as active request for backward compatibility
      _activeRequest.value = request;
      _importState.value = ImportState.extracting;

      _subscribeToImport(request.id, userId);
    } on Exception catch (e) {
      _error.value = e.toString();
      _importState.value = ImportState.error;
    }
  }

  String? _lastImportUserId;

  void _subscribeToImport(String requestId, String userId) {
    // Cancel any existing subscription for this request
    unawaited(_pendingSubscriptions[requestId]?.cancel());

    _pendingSubscriptions[requestId] = _importService
        .subscribeToRequest(requestId)
        .listen(
          (request) {
            // Update the pending import with new request status
            _updatePendingImport(requestId, (pending) => pending.copyWith(request: request));

            // Also update active request for backward compatibility
            if (_activeRequest.value?.id == requestId) {
              _activeRequest.value = request;
            }

            switch (request.status) {
              case RecipeRequestStatus.completed:
                unawaited(_onImportComplete(requestId));
              case RecipeRequestStatus.failed:
                _onImportError(requestId, 'Recipe extraction failed');
              case RecipeRequestStatus.requested:
              case RecipeRequestStatus.inProgress:
                // Stay in loading state
                break;
            }
          },
          onError: (Object error) {
            _onImportError(requestId, error.toString());
          },
        );
  }

  void _updatePendingImport(String requestId, PendingImport Function(PendingImport) update) {
    _pendingImports.value = _pendingImports.value.map((pending) {
      if (pending.request.id == requestId) {
        return update(pending);
      }
      return pending;
    }).toList();
  }

  Future<void> _onImportComplete(String requestId) async {
    try {
      final recipe = await _importService.getExtractedRecipe(requestId);
      if (recipe != null) {
        // Remove from pending and add to recipes
        _pendingImports.value = _pendingImports.value.where((p) => p.request.id != requestId).toList();
        _recipes.value = [recipe, ..._recipes.value];

        // Update state for backward compatibility
        if (_activeRequest.value?.id == requestId) {
          _previewRecipe.value = recipe;
          _importState.value = ImportState.preview;
        }

        // Clean up subscription
        unawaited(_pendingSubscriptions[requestId]?.cancel());
        _pendingSubscriptions.remove(requestId);
      } else {
        _onImportError(requestId, 'Recipe not found');
      }
    } on Exception catch (e) {
      _onImportError(requestId, e.toString());
    }
  }

  void _onImportError(String requestId, String errorMessage) {
    // Update pending import with error
    _updatePendingImport(
      requestId,
      (pending) => PendingImport(
        request: pending.request,
        userId: pending.userId,
        errorMessage: errorMessage,
      ),
    );

    // Update state for backward compatibility
    if (_activeRequest.value?.id == requestId) {
      _error.value = errorMessage;
      _importState.value = ImportState.error;
    }
  }

  /// Retries a failed import by its request ID.
  Future<void> retryPendingImport(String requestId) async {
    final pendingIndex = _pendingImports.value.indexWhere((p) => p.request.id == requestId);
    if (pendingIndex == -1) return;

    final pending = _pendingImports.value[pendingIndex];
    if (!pending.hasError) return;

    // Clear error and restart subscription
    _updatePendingImport(requestId, (p) => p.clearError());
    _subscribeToImport(requestId, pending.userId);

    // Re-import with the original URL
    await importRecipe(url: pending.request.url, userId: pending.userId);

    // Remove the old failed import
    dismissPendingImport(requestId);
  }

  /// Dismisses a pending import (removes from list).
  void dismissPendingImport(String requestId) {
    unawaited(_pendingSubscriptions[requestId]?.cancel());
    _pendingSubscriptions.remove(requestId);
    _pendingImports.value = _pendingImports.value.where((p) => p.request.id != requestId).toList();
  }

  /// Retries the last failed import.
  @Deprecated('Use retryPendingImport instead')
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
    unawaited(_requestSubscription?.cancel());
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
    unawaited(_requestSubscription?.cancel());
    for (final subscription in _pendingSubscriptions.values) {
      unawaited(subscription.cancel());
    }
    _pendingSubscriptions.clear();
  }
}
