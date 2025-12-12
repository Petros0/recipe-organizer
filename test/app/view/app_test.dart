// Ignore for testing purposes
// ignore_for_file: prefer_const_constructors

import 'package:flutter_test/flutter_test.dart';
import 'package:recipe_organizer/app/app.dart';
import 'package:recipe_organizer/core/di.dart';
import 'package:recipe_organizer/features/home/home.dart';

void main() {
  group('App', () {
    setUpAll(configureDependencies);

    testWidgets('renders HomePage', (tester) async {
      await tester.pumpWidget(App());
      expect(find.byType(HomePage), findsOneWidget);
    });
  });
}
