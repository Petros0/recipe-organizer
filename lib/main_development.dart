import 'package:recipe_organizer/app/app.dart';
import 'package:recipe_organizer/bootstrap.dart';

Future<void> main() async {
  await bootstrap(() => const App());
}
