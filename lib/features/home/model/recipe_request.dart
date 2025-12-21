import 'package:appwrite/models.dart';
import 'package:flutter/foundation.dart';

/// Status of a recipe import request.
enum RecipeRequestStatus {
  /// Initial state after document creation.
  requested('REQUESTED'),

  /// Backend function is processing.
  inProgress('IN_PROGRESS'),

  /// Recipe successfully extracted.
  completed('COMPLETED'),

  /// Extraction failed.
  failed('FAILED');

  const RecipeRequestStatus(this.value);

  /// The string value as stored in Appwrite.
  final String value;

  /// Creates a status from its string representation.
  static RecipeRequestStatus fromString(String value) {
    return RecipeRequestStatus.values.firstWhere(
      (status) => status.value == value,
      orElse: () => RecipeRequestStatus.requested,
    );
  }
}

/// Represents a pending recipe import request.
@immutable
class RecipeRequest {
  /// Creates a new [RecipeRequest] instance.
  const RecipeRequest({
    required this.id,
    required this.url,
    required this.status,
    required this.userId,
    required this.createdAt,
    required this.updatedAt,
  });

  /// Creates a [RecipeRequest] from an Appwrite Document.
  factory RecipeRequest.fromDocument(Document doc) {
    return RecipeRequest(
      id: doc.$id,
      url: doc.data['url'] as String,
      status: RecipeRequestStatus.fromString(doc.data['status'] as String),
      userId: doc.data['user_id'] as String,
      createdAt: DateTime.parse(doc.$createdAt),
      updatedAt: DateTime.parse(doc.$updatedAt),
    );
  }

  /// Creates a [RecipeRequest] from a map payload (e.g., from Realtime).
  factory RecipeRequest.fromMap(Map<String, dynamic> map) {
    return RecipeRequest(
      id: map[r'$id'] as String,
      url: map['url'] as String,
      status: RecipeRequestStatus.fromString(map['status'] as String),
      userId: map['user_id'] as String,
      createdAt: DateTime.parse(map[r'$createdAt'] as String),
      updatedAt: DateTime.parse(map[r'$updatedAt'] as String),
    );
  }

  /// Unique identifier for the request.
  final String id;

  /// Source URL submitted by user.
  final String url;

  /// Current processing status.
  final RecipeRequestStatus status;

  /// User ID who created the request.
  final String userId;

  /// When the request was created.
  final DateTime createdAt;

  /// When the request was last updated.
  final DateTime updatedAt;

  /// Creates a copy with the given fields replaced.
  RecipeRequest copyWith({
    String? id,
    String? url,
    RecipeRequestStatus? status,
    String? userId,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return RecipeRequest(
      id: id ?? this.id,
      url: url ?? this.url,
      status: status ?? this.status,
      userId: userId ?? this.userId,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is RecipeRequest && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}
