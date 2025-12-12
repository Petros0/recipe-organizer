import 'dart:developer';

import 'package:appwrite/appwrite.dart';
import 'package:appwrite/models.dart';
import 'package:recipe_organizer/features/auth/data/auth_repository.dart';

/// Service for handling authentication business logic.
class AuthService {
  AuthService({required AuthRepository repository}) : _repository = repository;

  final AuthRepository _repository;

  /// Ensures the user is authenticated.
  ///
  /// Checks for an existing session first. If no session exists,
  /// creates an anonymous session automatically.
  ///
  /// Returns the authenticated [User].
  Future<User> ensureAuthenticated() async {
    try {
      final user = await _repository.getCurrentUser();
      log('User already authenticated: ${user.$id}');
      return user;
    } on AppwriteException catch (e) {
      if (e.code == 401) {
        log('No existing session, creating anonymous session...');
        await _repository.createAnonymousSession();
        final user = await _repository.getCurrentUser();
        log('Anonymous session created: ${user.$id}');
        return user;
      }
      rethrow;
    }
  }
}
