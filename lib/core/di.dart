import 'package:appwrite/appwrite.dart';
import 'package:get_it/get_it.dart';
import 'package:recipe_organizer/features/auth/auth.dart';
import 'package:recipe_organizer/features/home/home.dart';

/// Global service locator instance for dependency injection.
final GetIt getIt = GetIt.instance;

void configureDependencies() {
  // Core
  getIt
    ..registerLazySingleton<Client>(
      () => Client()
        ..setProject('691f8b990030db50617a')
        ..setEndpoint('https://fra.cloud.appwrite.io/v1'),
    )
    ..registerLazySingleton<Functions>(() {
      final client = getIt<Client>();
      return Functions(client);
    })
    // Auth feature
    ..registerLazySingleton<AuthRepository>(() {
      final client = getIt<Client>();
      return AuthRepository(account: Account(client));
    })
    ..registerLazySingleton<AuthService>(() {
      final authRepository = getIt<AuthRepository>();
      return AuthService(repository: authRepository);
    })
    ..registerLazySingleton<AuthController>(() {
      final authService = getIt<AuthService>();
      return AuthController(service: authService);
    })
    // Home feature
    ..registerLazySingleton<HomeController>(HomeController.new)
    ..registerLazySingleton<RecipeRequestRepository>(() {
      final functions = getIt<Functions>();
      return RecipeRequestRepository(functions);
    });
}
