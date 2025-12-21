import 'package:appwrite/appwrite.dart';
import 'package:appwrite/models.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:recipe_organizer/features/home/data/recipe_repository.dart';

class MockDatabases extends Mock implements Databases {}

class FakeDocument extends Fake implements Document {
  FakeDocument({
    required String id,
    required Map<String, dynamic> data,
  }) : _id = id,
       _data = data;

  final String _id;
  final Map<String, dynamic> _data;

  @override
  String get $id => _id;

  @override
  Map<String, dynamic> get data => _data;
}

class FakeDocumentList extends Fake implements DocumentList {
  FakeDocumentList(this._documents);

  final List<Document> _documents;

  @override
  List<Document> get documents => _documents;
}

void main() {
  late RecipeRepository repository;
  late MockDatabases mockDatabases;

  setUp(() {
    mockDatabases = MockDatabases();
    repository = RecipeRepository(mockDatabases);
  });

  group('RecipeRepository', () {
    group('getRecipeByRequestId', () {
      test('returns recipe when found', () async {
        // Arrange
        const requestId = 'test-request-id';
        final fakeDoc = FakeDocument(
          id: 'recipe-id',
          data: {
            'name': 'Test Recipe',
            'description': 'A test recipe',
            'ingredients': <dynamic>['ingredient1'],
            'instructions': <dynamic>['step1'],
          },
        );

        when(
          () => mockDatabases.listDocuments(
            databaseId: any(named: 'databaseId'),
            collectionId: any(named: 'collectionId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeDocumentList([fakeDoc]));

        // Act
        final result = await repository.getRecipeByRequestId(requestId);

        // Assert
        expect(result, isNotNull);
        expect(result!.id, 'recipe-id');
        expect(result.name, 'Test Recipe');
      });

      test('returns null when not found', () async {
        // Arrange
        const requestId = 'test-request-id';

        when(
          () => mockDatabases.listDocuments(
            databaseId: any(named: 'databaseId'),
            collectionId: any(named: 'collectionId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeDocumentList([]));

        // Act
        final result = await repository.getRecipeByRequestId(requestId);

        // Assert
        expect(result, isNull);
      });
    });

    group('listRecipes', () {
      test('returns list of recipes', () async {
        // Arrange
        final fakeDoc1 = FakeDocument(
          id: 'recipe-1',
          data: {'name': 'Recipe 1'},
        );
        final fakeDoc2 = FakeDocument(
          id: 'recipe-2',
          data: {'name': 'Recipe 2'},
        );

        when(
          () => mockDatabases.listDocuments(
            databaseId: any(named: 'databaseId'),
            collectionId: any(named: 'collectionId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeDocumentList([fakeDoc1, fakeDoc2]));

        // Act
        final result = await repository.listRecipes();

        // Assert
        expect(result, hasLength(2));
        expect(result[0].name, 'Recipe 1');
        expect(result[1].name, 'Recipe 2');
      });

      test('returns empty list when no recipes', () async {
        // Arrange
        when(
          () => mockDatabases.listDocuments(
            databaseId: any(named: 'databaseId'),
            collectionId: any(named: 'collectionId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeDocumentList([]));

        // Act
        final result = await repository.listRecipes();

        // Assert
        expect(result, isEmpty);
      });
    });
  });
}
