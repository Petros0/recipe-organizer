import 'package:appwrite/models.dart';
import 'package:recipe_organizer/features/auth/service/auth_service.dart';
import 'package:signals/signals.dart';

/// Controller for managing authentication state using Signals.
class AuthController {
  AuthController({required AuthService service}) : _service = service {
    isAuthenticated = computed(() => _currentUser.value != null);
    isAnonymous = computed(() {
      final user = _currentUser.value;
      if (user == null) return false;
      // Anonymous users have no email set
      return user.email.isEmpty;
    });
  }

  final AuthService _service;
  final Signal<User?> _currentUser = signal<User?>(null);

  /// The currently authenticated user, or null if not authenticated.
  ReadonlySignal<User?> get currentUser => _currentUser;

  /// Whether the user is authenticated.
  late final Computed<bool> isAuthenticated;

  /// Whether the current user is anonymous.
  late final Computed<bool> isAnonymous;

  /// Ensures the user is authenticated.
  ///
  /// Checks for an existing session or creates an anonymous one.
  Future<void> ensureAuthenticated() async {
    final user = await _service.ensureAuthenticated();
    _currentUser.value = user;
  }

  /// Clears the current user (logout).
  void clearUser() {
    _currentUser.value = null;
  }
}
