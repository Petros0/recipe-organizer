import 'dart:async';
import 'dart:developer' as developer;

import 'package:appwrite/appwrite.dart';
import 'package:recipe_organizer/core/appwrite/appwrite_constants.dart';

/// Wrapper service for Appwrite Realtime subscriptions.
class RealtimeService {
  /// Creates a new [RealtimeService] instance.
  RealtimeService(this._realtime);

  final Realtime _realtime;

  /// Subscribes to document updates for a specific recipe request.
  ///
  /// Returns a stream that emits the document payload on each update.
  Stream<Map<String, dynamic>> subscribeToRecipeRequest(String documentId) {
    final channel =
        'databases.${AppwriteConstants.databaseId}'
        '.collections.${AppwriteConstants.recipeRequestCollectionId}'
        '.documents.$documentId';

    developer.log('Subscribing to channel: $channel');

    final subscription = _realtime.subscribe([channel]);

    return subscription.stream.map((message) {
      developer.log('Realtime message received: ${message.events}');
      developer.log('Realtime payload: ${message.payload}');
      return message.payload;
    });
  }

  /// Subscribes to all recipe request document events.
  Stream<Map<String, dynamic>> subscribeToRecipeRequests() {
    const channel =
        'databases.${AppwriteConstants.databaseId}'
        '.collections.${AppwriteConstants.recipeRequestCollectionId}'
        '.documents';

    final subscription = _realtime.subscribe([channel]);

    return subscription.stream.map((message) {
      return message.payload;
    });
  }
}
