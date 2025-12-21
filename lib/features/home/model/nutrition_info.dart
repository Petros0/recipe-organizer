import 'package:flutter/foundation.dart';

/// Nutritional information for a recipe.
@immutable
class NutritionInfo {
  /// Creates a new [NutritionInfo] instance.
  const NutritionInfo({
    this.calories,
    this.fatContent,
    this.saturatedFatContent,
    this.cholesterolContent,
    this.sodiumContent,
    this.carbohydrateContent,
    this.fiberContent,
    this.sugarContent,
    this.proteinContent,
  });

  /// Creates a [NutritionInfo] from a map.
  factory NutritionInfo.fromMap(Map<String, dynamic> map) {
    return NutritionInfo(
      calories: map['nutrition_calories'] as String?,
      fatContent: map['nutrition_fat_content'] as String?,
      saturatedFatContent: map['nutrition_saturated_fat_content'] as String?,
      cholesterolContent: map['nutrition_cholesterol_content'] as String?,
      sodiumContent: map['nutrition_sodium_content'] as String?,
      carbohydrateContent: map['nutrition_carbohydrate_content'] as String?,
      fiberContent: map['nutrition_fiber_content'] as String?,
      sugarContent: map['nutrition_sugar_content'] as String?,
      proteinContent: map['nutrition_protein_content'] as String?,
    );
  }

  /// Calorie content (e.g., "250 calories").
  final String? calories;

  /// Fat content (e.g., "12 g").
  final String? fatContent;

  /// Saturated fat content.
  final String? saturatedFatContent;

  /// Cholesterol content.
  final String? cholesterolContent;

  /// Sodium content.
  final String? sodiumContent;

  /// Carbohydrate content.
  final String? carbohydrateContent;

  /// Fiber content.
  final String? fiberContent;

  /// Sugar content.
  final String? sugarContent;

  /// Protein content.
  final String? proteinContent;

  /// Whether any nutrition data is available.
  bool get hasData => calories != null || fatContent != null || proteinContent != null || carbohydrateContent != null;

  /// Converts to a JSON map.
  Map<String, dynamic> toJson() {
    return {
      'nutrition_calories': calories,
      'nutrition_fat_content': fatContent,
      'nutrition_saturated_fat_content': saturatedFatContent,
      'nutrition_cholesterol_content': cholesterolContent,
      'nutrition_sodium_content': sodiumContent,
      'nutrition_carbohydrate_content': carbohydrateContent,
      'nutrition_fiber_content': fiberContent,
      'nutrition_sugar_content': sugarContent,
      'nutrition_protein_content': proteinContent,
    };
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is NutritionInfo && other.calories == calories && other.proteinContent == proteinContent;
  }

  @override
  int get hashCode => Object.hash(calories, proteinContent);
}
