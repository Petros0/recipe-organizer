import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:signals/signals.dart';

/// Controller for managing the home page state using Signals.
class HomeController {
  /// Creates a new [HomeController] instance.
  HomeController() {
    isEmpty = computed(() => _recipes.value.isEmpty);
  }

  final Signal<List<Recipe>> _recipes = signal<List<Recipe>>([]);
  final Signal<bool> _isLoading = signal<bool>(false);
  final Signal<String?> _error = signal<String?>(null);

  /// The list of recipes.
  ReadonlySignal<List<Recipe>> get recipes => _recipes;

  /// Whether recipes are currently being loaded.
  ReadonlySignal<bool> get isLoading => _isLoading;

  /// The current error message, if any.
  ReadonlySignal<String?> get error => _error;

  /// Whether the recipe list is empty.
  late final Computed<bool> isEmpty;

  /// Loads mock recipes for development.
  void loadMockRecipes() {
    _recipes.value = _mockRecipes;
  }

  /// Clears all recipes (for testing empty state).
  void clearRecipes() {
    _recipes.value = [];
  }

  /// Simulates importing a recipe from a URL.
  Future<void> importRecipe(String url) async {
    if (url.isEmpty) {
      _error.value = 'Please enter a valid URL';
      return;
    }

    _isLoading.value = true;
    _error.value = null;

    try {
      // Simulate network delay
      await Future<void>.delayed(const Duration(seconds: 1));

      // For now, add a mock recipe based on the URL
      final newRecipe = Recipe(
        id: DateTime.now().millisecondsSinceEpoch.toString(),
        name: 'Imported Recipe',
        imageUrl:
            'https://images.unsplash.com/photo-1546069901-ba9599a7e63c?w=400',
        description: 'A delicious recipe imported from $url',
        author: 'Chef',
        ingredients: const ['Ingredient 1', 'Ingredient 2', 'Ingredient 3'],
        instructions: const ['Step 1', 'Step 2', 'Step 3'],
      );

      _recipes.value = [..._recipes.value, newRecipe];
    } on Exception catch (e) {
      _error.value = 'Failed to import recipe: $e';
    } finally {
      _isLoading.value = false;
    }
  }

  /// Clears the current error.
  void clearError() {
    _error.value = null;
  }
}

/// Mock recipes for development and testing.
const List<Recipe> _mockRecipes = [
  Recipe(
    id: '1',
    name: 'Spaghetti Carbonara',
    imageUrl:
        'https://images.unsplash.com/photo-1612874742237-6526221588e3?w=400',
    description: 'Classic Italian pasta dish with eggs, cheese, and pancetta',
    author: 'Italian Chef',
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
    imageUrl:
        'https://images.unsplash.com/photo-1540189549336-e6e99c3679fe?w=400',
    description: 'Fresh Mediterranean salad with feta cheese and olives',
    author: 'Mediterranean Kitchen',
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
    imageUrl:
        'https://images.unsplash.com/photo-1624353365286-3f8d62daad51?w=400',
    description: 'Decadent chocolate dessert with molten center',
    author: 'Pastry Chef',
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
