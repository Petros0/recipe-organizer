import 'package:flutter_test/flutter_test.dart';
import 'package:recipe_organizer/features/home/model/nutrition_info.dart';

void main() {
  group('NutritionInfo', () {
    test('fromMap creates instance with all fields', () {
      // Arrange
      final inputMap = <String, dynamic>{
        'nutrition_calories': '250 kcal',
        'nutrition_fat_content': '12 g',
        'nutrition_saturated_fat_content': '4 g',
        'nutrition_cholesterol_content': '50 mg',
        'nutrition_sodium_content': '500 mg',
        'nutrition_carbohydrate_content': '30 g',
        'nutrition_fiber_content': '5 g',
        'nutrition_sugar_content': '10 g',
        'nutrition_protein_content': '20 g',
      };

      // Act
      final result = NutritionInfo.fromMap(inputMap);

      // Assert
      expect(result.calories, '250 kcal');
      expect(result.fatContent, '12 g');
      expect(result.saturatedFatContent, '4 g');
      expect(result.cholesterolContent, '50 mg');
      expect(result.sodiumContent, '500 mg');
      expect(result.carbohydrateContent, '30 g');
      expect(result.fiberContent, '5 g');
      expect(result.sugarContent, '10 g');
      expect(result.proteinContent, '20 g');
    });

    test('fromMap handles missing fields', () {
      // Arrange
      final inputMap = <String, dynamic>{
        'nutrition_calories': '250 kcal',
      };

      // Act
      final result = NutritionInfo.fromMap(inputMap);

      // Assert
      expect(result.calories, '250 kcal');
      expect(result.fatContent, isNull);
      expect(result.proteinContent, isNull);
    });

    test('hasData returns true when calories present', () {
      // Arrange
      const nutrition = NutritionInfo(calories: '250 kcal');

      // Assert
      expect(nutrition.hasData, isTrue);
    });

    test('hasData returns true when protein present', () {
      // Arrange
      const nutrition = NutritionInfo(proteinContent: '20 g');

      // Assert
      expect(nutrition.hasData, isTrue);
    });

    test('hasData returns false when no data', () {
      // Arrange
      const nutrition = NutritionInfo();

      // Assert
      expect(nutrition.hasData, isFalse);
    });

    test('toJson includes all fields', () {
      // Arrange
      const nutrition = NutritionInfo(
        calories: '250 kcal',
        fatContent: '12 g',
        proteinContent: '20 g',
      );

      // Act
      final result = nutrition.toJson();

      // Assert
      expect(result['nutrition_calories'], '250 kcal');
      expect(result['nutrition_fat_content'], '12 g');
      expect(result['nutrition_protein_content'], '20 g');
    });

    test('equality is based on calories and protein', () {
      // Arrange
      const nutrition1 = NutritionInfo(
        calories: '250 kcal',
        proteinContent: '20 g',
        fatContent: '10 g',
      );
      const nutrition2 = NutritionInfo(
        calories: '250 kcal',
        proteinContent: '20 g',
        fatContent: '15 g', // Different fat content
      );

      // Assert
      expect(nutrition1, nutrition2);
    });
  });
}
