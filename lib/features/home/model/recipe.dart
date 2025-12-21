import 'package:appwrite/models.dart';
import 'package:flutter/foundation.dart';
import 'package:recipe_organizer/features/home/model/nutrition_info.dart';

/// Recipe model representing a cooking recipe.
@immutable
class Recipe {
  /// Creates a new [Recipe] instance.
  const Recipe({
    required this.id,
    required this.name,
    this.description,
    this.images = const [],
    this.prepTime,
    this.cookTime,
    this.totalTime,
    this.recipeYield,
    this.ingredients = const [],
    this.instructions = const [],
    this.authorName,
    this.authorUrl,
    this.recipeCategory,
    this.recipeCuisine,
    this.keywords,
    this.datePublished,
    this.dateModified,
    this.nutrition,
    this.sourceUrl,
    this.recipeRequestId,
    this.userId,
  });

  /// Creates a [Recipe] from a JSON map (legacy format).
  factory Recipe.fromJson(Map<String, dynamic> json) {
    return Recipe(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String?,
      images: json['imageUrl'] != null ? [json['imageUrl'] as String] : [],
      prepTime: json['prepTime'] as String?,
      cookTime: json['cookTime'] as String?,
      totalTime: json['totalTime'] as String?,
      authorName: json['author'] as String?,
      ingredients: (json['ingredients'] as List<dynamic>?)?.cast<String>() ?? const [],
      instructions: (json['instructions'] as List<dynamic>?)?.cast<String>() ?? const [],
    );
  }

  /// Creates a [Recipe] from an Appwrite Document.
  factory Recipe.fromDocument(Document doc) {
    return Recipe(
      id: doc.$id,
      name: doc.data['name'] as String,
      description: doc.data['description'] as String?,
      images: (doc.data['image'] as List<dynamic>?)?.cast<String>() ?? [],
      prepTime: doc.data['prep_time'] as String?,
      cookTime: doc.data['cook_time'] as String?,
      totalTime: doc.data['total_time'] as String?,
      recipeYield: (doc.data['recipe_yield'] as List<dynamic>?)?.cast<String>(),
      ingredients: (doc.data['ingredients'] as List<dynamic>?)?.cast<String>() ?? [],
      instructions: (doc.data['instructions'] as List<dynamic>?)?.cast<String>() ?? [],
      authorName: doc.data['author_name'] as String?,
      authorUrl: doc.data['author_url'] as String?,
      recipeCategory: (doc.data['recipe_category'] as List<dynamic>?)?.cast<String>(),
      recipeCuisine: (doc.data['recipe_cuisine'] as List<dynamic>?)?.cast<String>(),
      keywords: doc.data['keywords'] as String?,
      datePublished: doc.data['date_published'] as String?,
      dateModified: doc.data['date_modified'] as String?,
      nutrition: NutritionInfo.fromMap(doc.data),
      sourceUrl: doc.data['source_url'] as String?,
      recipeRequestId: doc.data['fk_recipe_request'] as String?,
      userId: doc.data['user_id'] as String?,
    );
  }

  /// Unique identifier for the recipe.
  final String id;

  /// Name of the recipe.
  final String name;

  /// Description of the recipe.
  final String? description;

  /// Image URLs.
  final List<String> images;

  /// First image URL for convenience.
  String? get imageUrl => images.isNotEmpty ? images.first : null;

  /// Preparation time in ISO 8601 duration format.
  final String? prepTime;

  /// Cooking time in ISO 8601 duration format.
  final String? cookTime;

  /// Total time in ISO 8601 duration format.
  final String? totalTime;

  /// Recipe yield/servings.
  final List<String>? recipeYield;

  /// First yield value for convenience.
  String? get servings => recipeYield?.isNotEmpty ?? false ? recipeYield!.first : null;

  /// List of ingredients.
  final List<String> ingredients;

  /// List of cooking instructions.
  final List<String> instructions;

  /// Recipe author name.
  final String? authorName;

  /// Recipe author URL.
  final String? authorUrl;

  /// Recipe categories.
  final List<String>? recipeCategory;

  /// Recipe cuisines.
  final List<String>? recipeCuisine;

  /// Keywords (comma-separated).
  final String? keywords;

  /// Original publish date.
  final String? datePublished;

  /// Original modification date.
  final String? dateModified;

  /// Nutritional information.
  final NutritionInfo? nutrition;

  /// Original source URL.
  final String? sourceUrl;

  /// Foreign key to recipe request.
  final String? recipeRequestId;

  /// User ID who owns this recipe.
  final String? userId;

  /// Converts this [Recipe] to a JSON map.
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'image': images,
      'prep_time': prepTime,
      'cook_time': cookTime,
      'total_time': totalTime,
      'recipe_yield': recipeYield,
      'ingredients': ingredients,
      'instructions': instructions,
      'author_name': authorName,
      'author_url': authorUrl,
      'recipe_category': recipeCategory,
      'recipe_cuisine': recipeCuisine,
      'keywords': keywords,
      'date_published': datePublished,
      'date_modified': dateModified,
      'source_url': sourceUrl,
      'fk_recipe_request': recipeRequestId,
      'user_id': userId,
      ...?nutrition?.toJson(),
    };
  }

  /// Creates a copy of this [Recipe] with the given fields replaced.
  Recipe copyWith({
    String? id,
    String? name,
    String? description,
    List<String>? images,
    String? prepTime,
    String? cookTime,
    String? totalTime,
    List<String>? recipeYield,
    List<String>? ingredients,
    List<String>? instructions,
    String? authorName,
    String? authorUrl,
    List<String>? recipeCategory,
    List<String>? recipeCuisine,
    String? keywords,
    String? datePublished,
    String? dateModified,
    NutritionInfo? nutrition,
    String? sourceUrl,
    String? recipeRequestId,
    String? userId,
  }) {
    return Recipe(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      images: images ?? this.images,
      prepTime: prepTime ?? this.prepTime,
      cookTime: cookTime ?? this.cookTime,
      totalTime: totalTime ?? this.totalTime,
      recipeYield: recipeYield ?? this.recipeYield,
      ingredients: ingredients ?? this.ingredients,
      instructions: instructions ?? this.instructions,
      authorName: authorName ?? this.authorName,
      authorUrl: authorUrl ?? this.authorUrl,
      recipeCategory: recipeCategory ?? this.recipeCategory,
      recipeCuisine: recipeCuisine ?? this.recipeCuisine,
      keywords: keywords ?? this.keywords,
      datePublished: datePublished ?? this.datePublished,
      dateModified: dateModified ?? this.dateModified,
      nutrition: nutrition ?? this.nutrition,
      sourceUrl: sourceUrl ?? this.sourceUrl,
      recipeRequestId: recipeRequestId ?? this.recipeRequestId,
      userId: userId ?? this.userId,
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
