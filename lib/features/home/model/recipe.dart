import 'package:flutter/foundation.dart';

/// Recipe model representing a cooking recipe.
@immutable
class Recipe {
  /// Creates a new [Recipe] instance.
  const Recipe({
    required this.id,
    required this.name,
    this.imageUrl,
    this.description,
    this.prepTime,
    this.cookTime,
    this.totalTime,
    this.author,
    this.ingredients = const [],
    this.instructions = const [],
  });

  /// Creates a [Recipe] from a JSON map.
  factory Recipe.fromJson(Map<String, dynamic> json) {
    return Recipe(
      id: json['id'] as String,
      name: json['name'] as String,
      imageUrl: json['imageUrl'] as String?,
      description: json['description'] as String?,
      prepTime: json['prepTime'] as String?,
      cookTime: json['cookTime'] as String?,
      totalTime: json['totalTime'] as String?,
      author: json['author'] as String?,
      ingredients:
          (json['ingredients'] as List<dynamic>?)?.cast<String>() ?? const [],
      instructions:
          (json['instructions'] as List<dynamic>?)?.cast<String>() ?? const [],
    );
  }

  /// Unique identifier for the recipe.
  final String id;

  /// Name of the recipe.
  final String name;

  /// URL to the recipe image.
  final String? imageUrl;

  /// Description of the recipe.
  final String? description;

  /// Preparation time in ISO 8601 duration format.
  final String? prepTime;

  /// Cooking time in ISO 8601 duration format.
  final String? cookTime;

  /// Total time in ISO 8601 duration format.
  final String? totalTime;

  /// Author of the recipe.
  final String? author;

  /// List of ingredients.
  final List<String> ingredients;

  /// List of cooking instructions.
  final List<String> instructions;

  /// Converts this [Recipe] to a JSON map.
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'imageUrl': imageUrl,
      'description': description,
      'prepTime': prepTime,
      'cookTime': cookTime,
      'totalTime': totalTime,
      'author': author,
      'ingredients': ingredients,
      'instructions': instructions,
    };
  }

  /// Creates a copy of this [Recipe] with the given fields replaced.
  Recipe copyWith({
    String? id,
    String? name,
    String? imageUrl,
    String? description,
    String? prepTime,
    String? cookTime,
    String? totalTime,
    String? author,
    List<String>? ingredients,
    List<String>? instructions,
  }) {
    return Recipe(
      id: id ?? this.id,
      name: name ?? this.name,
      imageUrl: imageUrl ?? this.imageUrl,
      description: description ?? this.description,
      prepTime: prepTime ?? this.prepTime,
      cookTime: cookTime ?? this.cookTime,
      totalTime: totalTime ?? this.totalTime,
      author: author ?? this.author,
      ingredients: ingredients ?? this.ingredients,
      instructions: instructions ?? this.instructions,
    );
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Recipe && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}
