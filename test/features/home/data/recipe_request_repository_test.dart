import 'dart:convert';

import 'package:appwrite/appwrite.dart';
import 'package:appwrite/enums.dart' as enums;
import 'package:appwrite/models.dart' as models;
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:recipe_organizer/features/home/data/recipe_request_repository.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';

class MockFunctions extends Mock implements Functions {}

models.Execution _createFakeExecution(String responseBody) {
  return models.Execution(
    $id: 'test-execution-id',
    $createdAt: '2025-12-17T12:00:00.000Z',
    $updatedAt: '2025-12-17T12:00:00.000Z',
    $permissions: [],
    functionId: 'recipe-request',
    deploymentId: 'test-deployment-id',
    trigger: enums.ExecutionTrigger.http,
    status: enums.ExecutionStatus.completed,
    requestMethod: 'POST',
    requestPath: '/',
    requestHeaders: [],
    responseStatusCode: 200,
    responseBody: responseBody,
    responseHeaders: [],
    logs: '',
    errors: '',
    duration: 1,
    scheduledAt: '',
  );
}

void main() {
  late RecipeRequestRepository repository;
  late MockFunctions mockFunctions;

  setUp(() {
    mockFunctions = MockFunctions();
    repository = RecipeRequestRepository(mockFunctions);
  });

  group('RecipeRequestRepository', () {
    test('createRecipeRequest returns RecipeRequest on success', () async {
      // Arrange
      const inputUrl = 'https://example.com/recipe';
      const inputUserId = 'user-123';
      final responseBody = jsonEncode({
        'documentId': 'test-doc-id',
        'status': 'REQUESTED',
        'url': inputUrl,
      });

      when(
        () => mockFunctions.createExecution(
          functionId: any(named: 'functionId'),
          body: any(named: 'body'),
          method: any(named: 'method'),
          headers: any(named: 'headers'),
        ),
      ).thenAnswer((_) async => _createFakeExecution(responseBody));

      // Act
      final result = await repository.createRecipeRequest(
        url: inputUrl,
        userId: inputUserId,
      );

      // Assert
      expect(result.id, 'test-doc-id');
      expect(result.url, inputUrl);
      expect(result.userId, inputUserId);
      expect(result.status, RecipeRequestStatus.requested);
    });

    test('createRecipeRequest throws on error response', () async {
      // Arrange
      const inputUrl = 'https://example.com/recipe';
      const inputUserId = 'user-123';
      final responseBody = jsonEncode({
        'error': 'Invalid URL format',
      });

      when(
        () => mockFunctions.createExecution(
          functionId: any(named: 'functionId'),
          body: any(named: 'body'),
          method: any(named: 'method'),
          headers: any(named: 'headers'),
        ),
      ).thenAnswer((_) async => _createFakeExecution(responseBody));

      // Act & Assert
      expect(
        () => repository.createRecipeRequest(
          url: inputUrl,
          userId: inputUserId,
        ),
        throwsException,
      );
    });

    test('createRecipeRequest uses correct function ID', () async {
      // Arrange
      const inputUrl = 'https://example.com/recipe';
      const inputUserId = 'user-123';
      final responseBody = jsonEncode({
        'documentId': 'test-doc-id',
        'status': 'REQUESTED',
      });

      when(
        () => mockFunctions.createExecution(
          functionId: any(named: 'functionId'),
          body: any(named: 'body'),
          method: any(named: 'method'),
          headers: any(named: 'headers'),
        ),
      ).thenAnswer((_) async => _createFakeExecution(responseBody));

      // Act
      await repository.createRecipeRequest(url: inputUrl, userId: inputUserId);

      // Assert
      verify(
        () => mockFunctions.createExecution(
          functionId: 'recipe-request',
          body: any(named: 'body'),
          method: any(named: 'method'),
          headers: {'x-user-id': inputUserId},
        ),
      ).called(1);
    });
  });
}
