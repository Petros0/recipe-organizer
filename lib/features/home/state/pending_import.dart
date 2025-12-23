import 'package:flutter/foundation.dart';
import 'package:recipe_organizer/features/home/model/recipe_request.dart';

/// Represents a pending recipe import with its current state.
@immutable
class PendingImport {
  /// Creates a new [PendingImport] instance.
  const PendingImport({
    required this.request,
    required this.userId,
    this.errorMessage,
  });

  /// The underlying recipe request.
  final RecipeRequest request;

  /// Error message if import failed, null otherwise.
  final String? errorMessage;

  /// The user ID who initiated the import.
  final String userId;

  /// Whether this import has failed.
  bool get hasError => errorMessage != null;

  /// Whether this import is still loading.
  bool get isLoading => !hasError;

  /// Creates a copy with the given fields replaced.
  PendingImport copyWith({
    RecipeRequest? request,
    String? errorMessage,
    String? userId,
  }) {
    return PendingImport(
      request: request ?? this.request,
      errorMessage: errorMessage ?? this.errorMessage,
      userId: userId ?? this.userId,
    );
  }

  /// Creates a copy with error cleared.
  PendingImport clearError() {
    return PendingImport(
      request: request,
      userId: userId,
    );
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is PendingImport && other.request.id == request.id;
  }

  @override
  int get hashCode => request.id.hashCode;
}
