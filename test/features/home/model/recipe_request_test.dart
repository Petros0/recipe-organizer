import 'package:flutter_test/flutter_test.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';

void main() {
  group('RecipeRequestStatus', () {
    test('fromString returns correct status for valid values', () {
      // Arrange & Act & Assert
      expect(
        RecipeRequestStatus.fromString('REQUESTED'),
        RecipeRequestStatus.requested,
      );
      expect(
        RecipeRequestStatus.fromString('IN_PROGRESS'),
        RecipeRequestStatus.inProgress,
      );
      expect(
        RecipeRequestStatus.fromString('COMPLETED'),
        RecipeRequestStatus.completed,
      );
      expect(
        RecipeRequestStatus.fromString('FAILED'),
        RecipeRequestStatus.failed,
      );
    });

    test('fromString returns requested for unknown value', () {
      // Arrange
      const unknownValue = 'UNKNOWN';

      // Act
      final result = RecipeRequestStatus.fromString(unknownValue);

      // Assert
      expect(result, RecipeRequestStatus.requested);
    });

    test('value property returns correct string', () {
      // Assert
      expect(RecipeRequestStatus.requested.value, 'REQUESTED');
      expect(RecipeRequestStatus.inProgress.value, 'IN_PROGRESS');
      expect(RecipeRequestStatus.completed.value, 'COMPLETED');
      expect(RecipeRequestStatus.failed.value, 'FAILED');
    });
  });

  group('RecipeRequest', () {
    test('fromMap creates instance correctly', () {
      // Arrange
      final inputMap = <String, dynamic>{
        r'$id': 'test-id',
        'url': 'https://example.com/recipe',
        'status': 'IN_PROGRESS',
        'user_id': 'user-123',
        r'$createdAt': '2025-12-17T12:00:00.000Z',
        r'$updatedAt': '2025-12-17T12:05:00.000Z',
      };

      // Act
      final result = RecipeRequest.fromMap(inputMap);

      // Assert
      expect(result.id, 'test-id');
      expect(result.url, 'https://example.com/recipe');
      expect(result.status, RecipeRequestStatus.inProgress);
      expect(result.userId, 'user-123');
      expect(result.createdAt, DateTime.utc(2025, 12, 17, 12));
      expect(result.updatedAt, DateTime.utc(2025, 12, 17, 12, 5));
    });

    test('copyWith creates copy with replaced fields', () {
      // Arrange
      final original = RecipeRequest(
        id: 'test-id',
        url: 'https://example.com/recipe',
        status: RecipeRequestStatus.requested,
        userId: 'user-123',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // Act
      final result = original.copyWith(status: RecipeRequestStatus.completed);

      // Assert
      expect(result.id, original.id);
      expect(result.url, original.url);
      expect(result.userId, original.userId);
      expect(result.status, RecipeRequestStatus.completed);
    });

    test('equality is based on id', () {
      // Arrange
      final request1 = RecipeRequest(
        id: 'test-id',
        url: 'https://example.com/recipe1',
        status: RecipeRequestStatus.requested,
        userId: 'user-123',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );
      final request2 = RecipeRequest(
        id: 'test-id',
        url: 'https://example.com/recipe2',
        status: RecipeRequestStatus.completed,
        userId: 'user-456',
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      );

      // Assert
      expect(request1, request2);
      expect(request1.hashCode, request2.hashCode);
    });
  });
}
