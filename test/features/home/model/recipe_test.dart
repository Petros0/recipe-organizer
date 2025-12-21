import 'package:flutter_test/flutter_test.dart';
import 'package:recipe_organizer/features/home/model/nutrition_info.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';

void main() {
  group('Recipe', () {
    test('fromJson creates instance with legacy format', () {
      // Arrange
      final inputJson = <String, dynamic>{
        'id': 'test-id',
        'name': 'Test Recipe',
        'imageUrl': 'https://example.com/image.jpg',
        'description': 'A test description',
        'prepTime': 'PT15M',
        'cookTime': 'PT30M',
        'author': 'Test Author',
        'ingredients': ['ingredient1', 'ingredient2'],
        'instructions': ['step1', 'step2'],
      };

      // Act
      final result = Recipe.fromJson(inputJson);

      // Assert
      expect(result.id, 'test-id');
      expect(result.name, 'Test Recipe');
      expect(result.images, ['https://example.com/image.jpg']);
      expect(result.description, 'A test description');
      expect(result.prepTime, 'PT15M');
      expect(result.cookTime, 'PT30M');
      expect(result.authorName, 'Test Author');
      expect(result.ingredients, ['ingredient1', 'ingredient2']);
      expect(result.instructions, ['step1', 'step2']);
    });

    test('imageUrl getter returns first image', () {
      // Arrange
      const recipe = Recipe(
        id: 'test-id',
        name: 'Test Recipe',
        images: ['image1.jpg', 'image2.jpg'],
      );

      // Assert
      expect(recipe.imageUrl, 'image1.jpg');
    });

    test('imageUrl getter returns null when no images', () {
      // Arrange
      const recipe = Recipe(
        id: 'test-id',
        name: 'Test Recipe',
      );

      // Assert
      expect(recipe.imageUrl, isNull);
    });

    test('servings getter returns first yield value', () {
      // Arrange
      const recipe = Recipe(
        id: 'test-id',
        name: 'Test Recipe',
        recipeYield: ['4 servings', '8 portions'],
      );

      // Assert
      expect(recipe.servings, '4 servings');
    });

    test('copyWith creates copy with replaced fields', () {
      // Arrange
      const original = Recipe(
        id: 'test-id',
        name: 'Original Name',
        description: 'Original description',
      );

      // Act
      final result = original.copyWith(name: 'New Name');

      // Assert
      expect(result.id, 'test-id');
      expect(result.name, 'New Name');
      expect(result.description, 'Original description');
    });

    test('toJson includes all fields', () {
      // Arrange
      const recipe = Recipe(
        id: 'test-id',
        name: 'Test Recipe',
        description: 'A test description',
        images: ['image.jpg'],
        prepTime: 'PT15M',
        cookTime: 'PT30M',
        ingredients: ['ingredient1'],
        instructions: ['step1'],
        authorName: 'Author',
        sourceUrl: 'https://example.com',
        nutrition: NutritionInfo(calories: '100 kcal'),
      );

      // Act
      final result = recipe.toJson();

      // Assert
      expect(result['name'], 'Test Recipe');
      expect(result['description'], 'A test description');
      expect(result['image'], ['image.jpg']);
      expect(result['prep_time'], 'PT15M');
      expect(result['cook_time'], 'PT30M');
      expect(result['author_name'], 'Author');
      expect(result['source_url'], 'https://example.com');
      expect(result['nutrition_calories'], '100 kcal');
    });

    test('equality is based on id', () {
      // Arrange
      const recipe1 = Recipe(id: 'test-id', name: 'Recipe 1');
      const recipe2 = Recipe(id: 'test-id', name: 'Recipe 2');

      // Assert
      expect(recipe1, recipe2);
      expect(recipe1.hashCode, recipe2.hashCode);
    });
  });
}
