import 'package:flutter/widgets.dart';
import 'package:recipe_organizer/l10n/gen/app_localizations.dart';

export 'package:recipe_organizer/l10n/gen/app_localizations.dart';

extension AppLocalizationsX on BuildContext {
  AppLocalizations get l10n => AppLocalizations.of(this);
}
