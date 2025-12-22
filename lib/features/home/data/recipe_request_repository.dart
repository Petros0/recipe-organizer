import 'dart:convert';
import 'dart:developer' as developer;

import 'package:appwrite/appwrite.dart';
import 'package:appwrite/enums.dart';
import 'package:recipe_organizer/core/appwrite/appwrite_constants.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';

/// Repository for recipe request operations using Appwrite Functions service.
class RecipeRequestRepository {
  /// Creates a new [RecipeRequestRepository] instance.
  RecipeRequestRepository(this._functions);

  final Functions _functions;

  /// Creates a new recipe import request.
  ///
  /// The [userId] is passed to the function to scope the created documents
  /// to the authenticated user.
  /// Returns the [RecipeRequest] created from the function response.
  Future<RecipeRequest> createRecipeRequest({
    required String url,
    required String userId,
  }) async {
    developer.log('Calling recipe-request function with URL: $url');

    final execution = await _functions.createExecution(
      functionId: AppwriteConstants.recipeRequestFunctionId,
      body: jsonEncode({'url': url}),
      method: ExecutionMethod.pOST,
      headers: {'x-user-id': userId},
    );

    developer.log('Function response status: ${execution.status}');
    developer.log('Function response body: ${execution.responseBody}');

    final responseBody = jsonDecode(execution.responseBody) as Map<String, dynamic>;

    if (responseBody.containsKey('error')) {
      throw Exception(responseBody['error'] as String);
    }

    final documentId = responseBody['documentId'] as String;
    developer.log('Created document ID: $documentId');

    return RecipeRequest(
      id: documentId,
      url: url,
      status: RecipeRequestStatus.fromString(
        responseBody['status'] as String? ?? 'REQUESTED',
      ),
      userId: userId,
      createdAt: DateTime.now(),
      updatedAt: DateTime.now(),
    );
  }
}
