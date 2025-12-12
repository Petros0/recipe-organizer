import 'package:appwrite/appwrite.dart';
import 'package:appwrite/models.dart';

/// Repository for authentication operations using Appwrite Account service.
class AuthRepository {
  AuthRepository({required Account account}) : _account = account;

  final Account _account;

  /// Retrieves the currently authenticated user.
  ///
  /// Returns the [User] if authenticated, or throws an [AppwriteException]
  /// if no session exists.
  Future<User> getCurrentUser() async {
    return _account.get();
  }

  /// Creates an anonymous session for unauthenticated users.
  ///
  /// Returns a [Session] representing the new anonymous session.
  Future<Session> createAnonymousSession() async {
    return _account.createAnonymousSession();
  }
}
