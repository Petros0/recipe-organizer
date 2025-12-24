import 'package:appwrite/appwrite.dart';
import 'package:get_it/get_it.dart';
import 'package:recipe_organizer/core/appwrite/appwrite_constants.dart';
import 'package:recipe_organizer/core/appwrite/realtime_service.dart';
import 'package:recipe_organizer/features/auth/auth.dart';
import 'package:recipe_organizer/features/home/home.dart';

/// Global service locator instance for dependency injection.
final GetIt getIt = GetIt.instance;

void configureDependencies() {
  // Core
  getIt
    ..registerLazySingleton<Client>(
      () => Client()
        ..setProject(AppwriteConstants.projectId)
        ..setEndpoint(AppwriteConstants.endpoint),
    )
    ..registerLazySingleton<Functions>(() {
      final client = getIt<Client>();
      return Functions(client);
    })
    ..registerLazySingleton<TablesDB>(() {
      final client = getIt<Client>();
      return TablesDB(client);
    })
    ..registerLazySingleton<Realtime>(() {
      final client = getIt<Client>();
      return Realtime(client);
    })
    ..registerLazySingleton<RealtimeService>(() {
      final realtime = getIt<Realtime>();
      return RealtimeService(realtime);
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
    ..registerLazySingleton<RecipeRequestRepository>(() {
      final functions = getIt<Functions>();
      return RecipeRequestRepository(functions);
    })
    ..registerLazySingleton<RecipeRepository>(() {
      final tablesDB = getIt<TablesDB>();
      return RecipeRepository(tablesDB);
    })
    ..registerLazySingleton<RecipeImportService>(() {
      final requestRepo = getIt<RecipeRequestRepository>();
      final recipeRepo = getIt<RecipeRepository>();
      final realtimeService = getIt<RealtimeService>();
      return RecipeImportService(
        requestRepository: requestRepo,
        recipeRepository: recipeRepo,
        realtimeService: realtimeService,
      );
    })
    ..registerLazySingleton<HomeController>(() {
      final importService = getIt<RecipeImportService>();
      return HomeController(importService: importService);
    });
}
