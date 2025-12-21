// GENERATED CODE - DO NOT MODIFY BY HAND
// Consider adding this file to your .gitignore.

import 'dart:io';
import 'dart:typed_data';

import 'package:flutter_test/flutter_test.dart';


import 'app/view/app_test.dart' as _a;
import 'features/home/state/home_controller_test.dart' as _b;
import 'features/home/model/nutrition_info_test.dart' as _c;
import 'features/home/model/recipe_test.dart' as _d;
import 'features/home/model/recipe_request_test.dart' as _e;
import 'features/home/view/widgets/recipe_error_card_test.dart' as _f;
import 'features/home/view/widgets/recipe_skeleton_card_test.dart' as _g;
import 'features/home/view/widgets/import_recipe_dialog_test.dart' as _h;
import 'features/home/service/recipe_import_service_test.dart' as _i;
import 'features/home/data/recipe_repository_test.dart' as _j;
import 'features/home/data/recipe_request_repository_test.dart' as _k;

void main() {
  goldenFileComparator = _TestOptimizationAwareGoldenFileComparator(goldenFileComparator as LocalFileComparator);
  group('app/view/app_test.dart', () { _a.main(); });
  group('features/home/state/home_controller_test.dart', () { _b.main(); });
  group('features/home/model/nutrition_info_test.dart', () { _c.main(); });
  group('features/home/model/recipe_test.dart', () { _d.main(); });
  group('features/home/model/recipe_request_test.dart', () { _e.main(); });
  group('features/home/view/widgets/recipe_error_card_test.dart', () { _f.main(); });
  group('features/home/view/widgets/recipe_skeleton_card_test.dart', () { _g.main(); });
  group('features/home/view/widgets/import_recipe_dialog_test.dart', () { _h.main(); });
  group('features/home/service/recipe_import_service_test.dart', () { _i.main(); });
  group('features/home/data/recipe_repository_test.dart', () { _j.main(); });
  group('features/home/data/recipe_request_repository_test.dart', () { _k.main(); });
}


class _TestOptimizationAwareGoldenFileComparator extends GoldenFileComparator {
  final List<String> goldenFilePaths;
  final LocalFileComparator previousGoldenFileComparator;

  _TestOptimizationAwareGoldenFileComparator(this.previousGoldenFileComparator)
      : goldenFilePaths = _goldenFilePaths;

  static List<String> get _goldenFilePaths =>
      Directory.fromUri((goldenFileComparator as LocalFileComparator).basedir)
          .listSync(recursive: true, followLinks: true)
          .whereType<File>()
          .map((file) => file.path)
          .where((path) => path.endsWith('.png'))
          .toList();
  @override
  Future<bool> compare(Uint8List imageBytes, Uri golden)  => previousGoldenFileComparator.compare(imageBytes, golden);

  @override
  Uri getTestUri(Uri key, int? version) {
    final keyString = key.toFilePath();
    return Uri.parse(goldenFilePaths
        .singleWhere((goldenFilePath) => goldenFilePath.endsWith(keyString)));
  }

  @override
  Future<void> update(Uri golden, Uint8List imageBytes) => previousGoldenFileComparator.update(golden, imageBytes);

}
