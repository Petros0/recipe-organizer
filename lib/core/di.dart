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
    // Auth feature
    ..registerLazySingleton<AuthRepository>(() => AuthRepository(account: Account(getIt<Client>())))
    ..registerLazySingleton<AuthService>(() => AuthService(repository: getIt<AuthRepository>()))
    ..registerLazySingleton<AuthController>(() => AuthController(service: getIt<AuthService>()))
    // Home feature
    ..registerLazySingleton<HomeController>(HomeController.new);
}
