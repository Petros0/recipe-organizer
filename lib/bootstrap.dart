import 'dart:async';
import 'dart:developer';

import 'package:flutter/widgets.dart';
import 'package:recipe_organizer/core/di.dart';
import 'package:recipe_organizer/features/auth/auth.dart';

Future<void> bootstrap(FutureOr<Widget> Function() builder) async {
  WidgetsFlutterBinding.ensureInitialized();

  FlutterError.onError = (details) {
    log(details.exceptionAsString(), stackTrace: details.stack);
  };

  // Configure dependencies
  configureDependencies();

  // Initialize authentication - auto-login as anonymous if needed
  await getIt<AuthController>().ensureAuthenticated();

  runApp(await builder());
}
