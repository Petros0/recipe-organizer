/// State of the recipe import flow.
enum ImportState {
  /// No active import.
  idle,

  /// Creating request document.
  submitting,

  /// Showing skeleton, waiting for backend extraction.
  extracting,

  /// Recipe extracted, showing preview.
  preview,

  /// Import failed.
  error,

  /// Saving recipe to collection.
  saving,
}
