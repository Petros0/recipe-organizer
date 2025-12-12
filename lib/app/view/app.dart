import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/features/home/home.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = FThemes.zinc.dark;
    return MaterialApp(
      theme: theme.toApproximateMaterialTheme(),
      builder: (_, child) => FAnimatedTheme(data: theme, child: child!),
      localizationsDelegates: AppLocalizations.localizationsDelegates,
      supportedLocales: AppLocalizations.supportedLocales,
      home: const HomePage(),
    );
  }
}
