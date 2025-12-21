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

  /// Loads mock recipes for development.
  void loadMockRecipes() {
    _recipes.value = _mockRecipes;
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

/// Mock recipes for development and testing.
const List<Recipe> _mockRecipes = [
  Recipe(
    id: '1',
    name: 'Spaghetti Carbonara',
    images: [
      'https://images.unsplash.com/photo-1612874742237-6526221588e3?w=400',
    ],
    description: 'Classic Italian pasta dish with eggs, cheese, and pancetta',
    authorName: 'Italian Chef',
    prepTime: 'PT15M',
    cookTime: 'PT20M',
    totalTime: 'PT35M',
    ingredients: [
      '400g spaghetti',
      '200g pancetta',
      '4 egg yolks',
      '100g pecorino cheese',
      'Black pepper',
    ],
    instructions: [
      'Cook the spaghetti in salted water',
      'Fry the pancetta until crispy',
      'Mix egg yolks with cheese',
      'Combine everything off heat',
      'Season with black pepper',
    ],
  ),
  Recipe(
    id: '2',
    name: 'Greek Salad',
    images: [
      'https://images.unsplash.com/photo-1540189549336-e6e99c3679fe?w=400',
    ],
    description: 'Fresh Mediterranean salad with feta cheese and olives',
    authorName: 'Mediterranean Kitchen',
    prepTime: 'PT10M',
    totalTime: 'PT10M',
    ingredients: [
      'Tomatoes',
      'Cucumber',
      'Red onion',
      'Feta cheese',
      'Kalamata olives',
      'Olive oil',
      'Oregano',
    ],
    instructions: [
      'Chop vegetables into chunks',
      'Add feta cheese and olives',
      'Drizzle with olive oil',
      'Sprinkle oregano',
    ],
  ),
  Recipe(
    id: '3',
    name: 'Chocolate Lava Cake',
    images: [
      'https://images.unsplash.com/photo-1624353365286-3f8d62daad51?w=400',
    ],
    description: 'Decadent chocolate dessert with molten center',
    authorName: 'Pastry Chef',
    prepTime: 'PT15M',
    cookTime: 'PT12M',
    totalTime: 'PT27M',
    ingredients: [
      '200g dark chocolate',
      '100g butter',
      '2 eggs',
      '2 egg yolks',
      '50g sugar',
      '25g flour',
    ],
    instructions: [
      'Melt chocolate and butter',
      'Whisk eggs with sugar',
      'Combine and add flour',
      'Bake at 200°C for 12 minutes',
      'Serve immediately',
    ],
  ),
];
