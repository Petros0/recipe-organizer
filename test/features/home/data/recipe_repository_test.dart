import 'package:appwrite/appwrite.dart';
import 'package:appwrite/models.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:recipe_organizer/features/home/data/recipe_repository.dart';

class MockTablesDB extends Mock implements TablesDB {}

class FakeRow extends Fake implements Row {
  FakeRow({
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

class FakeRowList extends Fake implements RowList {
  FakeRowList(this._rows);

  final List<Row> _rows;

  @override
  List<Row> get rows => _rows;
}

void main() {
  late RecipeRepository repository;
  late MockTablesDB mockTablesDB;

  setUp(() {
    mockTablesDB = MockTablesDB();
    repository = RecipeRepository(mockTablesDB);
  });

  group('RecipeRepository', () {
    group('getRecipeByRequestId', () {
      test('returns recipe when found', () async {
        // Arrange
        const requestId = 'test-request-id';
        final fakeRow = FakeRow(
          id: 'recipe-id',
          data: {
            'name': 'Test Recipe',
            'description': 'A test recipe',
            'ingredients': <dynamic>['ingredient1'],
            'instructions': <dynamic>['step1'],
          },
        );

        when(
          () => mockTablesDB.listRows(
            databaseId: any(named: 'databaseId'),
            tableId: any(named: 'tableId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeRowList([fakeRow]));

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
          () => mockTablesDB.listRows(
            databaseId: any(named: 'databaseId'),
            tableId: any(named: 'tableId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeRowList([]));

        // Act
        final result = await repository.getRecipeByRequestId(requestId);

        // Assert
        expect(result, isNull);
      });
    });

    group('listRecipes', () {
      test('returns list of recipes', () async {
        // Arrange
        final fakeRow1 = FakeRow(
          id: 'recipe-1',
          data: {'name': 'Recipe 1'},
        );
        final fakeRow2 = FakeRow(
          id: 'recipe-2',
          data: {'name': 'Recipe 2'},
        );

        when(
          () => mockTablesDB.listRows(
            databaseId: any(named: 'databaseId'),
            tableId: any(named: 'tableId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeRowList([fakeRow1, fakeRow2]));

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
          () => mockTablesDB.listRows(
            databaseId: any(named: 'databaseId'),
            tableId: any(named: 'tableId'),
            queries: any(named: 'queries'),
          ),
        ).thenAnswer((_) async => FakeRowList([]));

        // Act
        final result = await repository.listRecipes();

        // Assert
        expect(result, isEmpty);
      });
    });
  });
}
