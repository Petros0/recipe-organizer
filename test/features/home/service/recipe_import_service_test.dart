import 'dart:async';

import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:recipe_organizer/core/appwrite/realtime_service.dart';
import 'package:recipe_organizer/features/home/data/recipe_repository.dart';
import 'package:recipe_organizer/features/home/data/recipe_request_repository.dart';
import 'package:recipe_organizer/features/home/model/recipe.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';
import 'package:recipe_organizer/features/home/service/recipe_import_service.dart';

class MockRecipeRequestRepository extends Mock implements RecipeRequestRepository {}

class MockRecipeRepository extends Mock implements RecipeRepository {}

class MockRealtimeService extends Mock implements RealtimeService {}

void main() {
  late RecipeImportService service;
  late MockRecipeRequestRepository mockRequestRepo;
  late MockRecipeRepository mockRecipeRepo;
  late MockRealtimeService mockRealtimeService;

  setUp(() {
    mockRequestRepo = MockRecipeRequestRepository();
    mockRecipeRepo = MockRecipeRepository();
    mockRealtimeService = MockRealtimeService();
    service = RecipeImportService(
      requestRepository: mockRequestRepo,
      recipeRepository: mockRecipeRepo,
      realtimeService: mockRealtimeService,
    );
  });

  group('RecipeImportService', () {
    group('importFromUrl', () {
      test('calls requestRepository.createRecipeRequest', () async {
        // Arrange
        const inputUrl = 'https://example.com/recipe';
        const inputUserId = 'user-123';
        final expectedRequest = RecipeRequest(
          id: 'test-id',
          url: inputUrl,
          status: RecipeRequestStatus.requested,
          userId: inputUserId,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        when(
          () => mockRequestRepo.createRecipeRequest(
            url: inputUrl,
            userId: inputUserId,
          ),
        ).thenAnswer((_) async => expectedRequest);

        // Act
        final result = await service.importFromUrl(
          url: inputUrl,
          userId: inputUserId,
        );

        // Assert
        expect(result, expectedRequest);
        verify(
          () => mockRequestRepo.createRecipeRequest(
            url: inputUrl,
            userId: inputUserId,
          ),
        ).called(1);
      });
    });

    group('subscribeToRequest', () {
      test('returns stream of RecipeRequest updates', () async {
        // Arrange
        const requestId = 'test-request-id';
        final payloadStream = Stream.fromIterable([
          <String, dynamic>{
            r'$id': requestId,
            'url': 'https://example.com/recipe',
            'status': 'IN_PROGRESS',
            'user_id': 'user-123',
            r'$createdAt': '2025-12-17T12:00:00.000Z',
            r'$updatedAt': '2025-12-17T12:05:00.000Z',
          },
        ]);

        when(
          () => mockRealtimeService.subscribeToRecipeRequest(requestId),
        ).thenAnswer((_) => payloadStream);

        // Act
        final result = service.subscribeToRequest(requestId);
        final requests = await result.toList();

        // Assert
        expect(requests, hasLength(1));
        expect(requests.first.id, requestId);
        expect(requests.first.status, RecipeRequestStatus.inProgress);
      });
    });

    group('getExtractedRecipe', () {
      test('calls recipeRepository.getRecipeByRequestId', () async {
        // Arrange
        const requestId = 'test-request-id';
        const expectedRecipe = Recipe(
          id: 'recipe-id',
          name: 'Test Recipe',
        );

        when(
          () => mockRecipeRepo.getRecipeByRequestId(requestId),
        ).thenAnswer((_) async => expectedRecipe);

        // Act
        final result = await service.getExtractedRecipe(requestId);

        // Assert
        expect(result, expectedRecipe);
        verify(() => mockRecipeRepo.getRecipeByRequestId(requestId)).called(1);
      });
    });

    group('listRecipes', () {
      test('calls recipeRepository.listRecipes with defaults', () async {
        // Arrange
        const expectedRecipes = [
          Recipe(id: '1', name: 'Recipe 1'),
          Recipe(id: '2', name: 'Recipe 2'),
        ];

        when(
          () => mockRecipeRepo.listRecipes(),
        ).thenAnswer((_) async => expectedRecipes);

        // Act
        final result = await service.listRecipes();

        // Assert
        expect(result, expectedRecipes);
        verify(() => mockRecipeRepo.listRecipes()).called(1);
      });
    });
  });
}
