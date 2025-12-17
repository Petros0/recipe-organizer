import 'package:appwrite/appwrite.dart';
import 'package:appwrite/enums.dart';
import 'package:appwrite/models.dart';

/// Repository for recipe request operations using Appwrite Functions service.
class RecipeRequestRepository {
  RecipeRequestRepository(this._functions);

  final Functions _functions;

  final String _functionId = 'recipe-request';

  Future<void> createRecipeRequest(String url) async {
    final execution = await _functions.createExecution(
      functionId: _functionId,
      path: '?url=$url',
    );
    
    //   Future result = functions.createExecution(
    //   functionId: '<FUNCTION_ID>',
    //   body: '<BODY>', // optional
    //   xasync: false, // optional
    //   path: '<PATH>', // optional
    //   method: 'GET', // optional
    //   headers: {}, // optional
    // );

    // result
    //     .then((response) {
    //       print(response); // Success
    //     })
    //     .catchError((error) {
    //       print(error.response); // Failure
    //     });
    await _functions.createRecipeRequest(url);
  }
}