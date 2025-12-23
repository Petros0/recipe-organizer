import 'dart:async';

import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';
import 'package:recipe_organizer/features/home/service/recipe_import_service.dart';
import 'package:recipe_organizer/features/home/state/home_controller.dart';
import 'package:recipe_organizer/features/home/state/import_state.dart';

class MockRecipeImportService extends Mock implements RecipeImportService {}

void main() {
  late HomeController controller;
  late MockRecipeImportService mockImportService;

  setUp(() {
    mockImportService = MockRecipeImportService();
    controller = HomeController(importService: mockImportService);
  });

  tearDown(() {
    controller.dispose();
  });

  group('HomeController', () {
    group('initial state', () {
      test('has empty recipes', () {
        expect(controller.recipes.value, isEmpty);
      });

      test('isEmpty is true', () {
        expect(controller.isEmpty.value, isTrue);
      });

      test('importState is idle', () {
        expect(controller.importState.value, ImportState.idle);
      });
    });

    group('loadRecipes', () {
      test('populates recipes list from service', () async {
        // Arrange
        const recipes = [
          Recipe(id: '1', name: 'Recipe 1'),
          Recipe(id: '2', name: 'Recipe 2'),
        ];
        when(() => mockImportService.listRecipes(limit: 25, offset: 0)).thenAnswer((_) async => recipes);

        // Act
        await controller.loadRecipes();

        // Assert
        expect(controller.recipes.value, equals(recipes));
        expect(controller.isEmpty.value, isFalse);
      });

      test('sets error on failure', () async {
        // Arrange
        when(() => mockImportService.listRecipes(limit: 25, offset: 0)).thenThrow(Exception('Network error'));

        // Act
        await controller.loadRecipes();

        // Assert
        expect(controller.error.value, isNotNull);
        expect(controller.error.value, contains('Failed to load recipes'));
      });
    });

    group('clearRecipes', () {
      test('empties recipes list', () async {
        // Arrange
        const recipes = [Recipe(id: '1', name: 'Recipe 1')];
        when(() => mockImportService.listRecipes(limit: 25, offset: 0)).thenAnswer((_) async => recipes);
        await controller.loadRecipes();

        // Act
        controller.clearRecipes();

        // Assert
        expect(controller.recipes.value, isEmpty);
        expect(controller.isEmpty.value, isTrue);
      });
    });

    group('deleteRecipe', () {
      test('removes recipe from list', () async {
        // Arrange
        const recipes = [
          Recipe(id: '1', name: 'Recipe 1'),
          Recipe(id: '2', name: 'Recipe 2'),
        ];
        when(() => mockImportService.listRecipes(limit: 25, offset: 0)).thenAnswer((_) async => recipes);
        when(() => mockImportService.deleteRecipe('1')).thenAnswer((_) async {});
        await controller.loadRecipes();

        // Act
        await controller.deleteRecipe('1');

        // Assert
        expect(controller.recipes.value.length, equals(1));
        expect(controller.recipes.value.first.id, equals('2'));
      });

      test('sets error on failure', () async {
        // Arrange
        when(() => mockImportService.deleteRecipe('1')).thenThrow(Exception('Delete failed'));

        // Act
        await controller.deleteRecipe('1');

        // Assert
        expect(controller.error.value, isNotNull);
        expect(controller.error.value, contains('Failed to delete recipe'));
      });
    });

    group('importRecipe', () {
      test('sets error for empty URL', () async {
        // Act
        await controller.importRecipe(url: '', userId: 'user-123');

        // Assert
        expect(controller.error.value, 'Please enter a valid URL');
      });

      test('sets state to submitting then extracting on success', () async {
        // Arrange
        const inputUrl = 'https://example.com/recipe';
        const inputUserId = 'user-123';
        final request = RecipeRequest(
          id: 'test-id',
          url: inputUrl,
          status: RecipeRequestStatus.requested,
          userId: inputUserId,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        when(
          () => mockImportService.importFromUrl(
            url: inputUrl,
            userId: inputUserId,
          ),
        ).thenAnswer((_) async => request);
        when(
          () => mockImportService.subscribeToRequest(request.id),
        ).thenAnswer((_) => const Stream.empty());

        // Act
        await controller.importRecipe(url: inputUrl, userId: inputUserId);

        // Assert
        expect(controller.importState.value, ImportState.extracting);
        expect(controller.activeRequest.value, request);
      });

      test('sets error state on exception', () async {
        // Arrange
        const inputUrl = 'https://example.com/recipe';
        const inputUserId = 'user-123';

        when(
          () => mockImportService.importFromUrl(
            url: inputUrl,
            userId: inputUserId,
          ),
        ).thenThrow(Exception('Network error'));

        // Act
        await controller.importRecipe(url: inputUrl, userId: inputUserId);

        // Assert
        expect(controller.importState.value, ImportState.error);
        expect(controller.error.value, isNotNull);
      });
    });

    group('saveRecipe', () {
      test('adds preview recipe to list', () async {
        // Arrange
        const inputUrl = 'https://example.com/recipe';
        const inputUserId = 'user-123';
        final request = RecipeRequest(
          id: 'test-id',
          url: inputUrl,
          status: RecipeRequestStatus.completed,
          userId: inputUserId,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        const recipe = Recipe(id: 'recipe-id', name: 'Test Recipe');

        // Setup the controller state
        when(
          () => mockImportService.importFromUrl(
            url: inputUrl,
            userId: inputUserId,
          ),
        ).thenAnswer((_) async => request);
        when(
          () => mockImportService.subscribeToRequest(request.id),
        ).thenAnswer((_) => Stream.value(request));
        when(
          () => mockImportService.getExtractedRecipe(request.id),
        ).thenAnswer((_) async => recipe);

        await controller.importRecipe(url: inputUrl, userId: inputUserId);

        // Wait for stream processing
        await Future<void>.delayed(const Duration(milliseconds: 100));

        // Act
        await controller.saveRecipe();

        // Assert
        expect(controller.recipes.value, contains(recipe));
        expect(controller.importState.value, ImportState.idle);
      });
    });

    group('cancelImport', () {
      test('resets import state', () async {
        // Arrange
        const inputUrl = 'https://example.com/recipe';
        const inputUserId = 'user-123';
        final request = RecipeRequest(
          id: 'test-id',
          url: inputUrl,
          status: RecipeRequestStatus.requested,
          userId: inputUserId,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        when(
          () => mockImportService.importFromUrl(
            url: inputUrl,
            userId: inputUserId,
          ),
        ).thenAnswer((_) async => request);
        when(
          () => mockImportService.subscribeToRequest(request.id),
        ).thenAnswer((_) => const Stream.empty());

        await controller.importRecipe(url: inputUrl, userId: inputUserId);

        // Act
        controller.cancelImport();

        // Assert
        expect(controller.importState.value, ImportState.idle);
        expect(controller.activeRequest.value, isNull);
        expect(controller.previewRecipe.value, isNull);
      });
    });
  });
}
