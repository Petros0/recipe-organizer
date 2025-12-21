import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';
import 'package:url_launcher/url_launcher.dart';

/// Widget displaying source attribution for a recipe.
class SourceAttribution extends StatelessWidget {
  /// Creates a new [SourceAttribution] instance.
  const SourceAttribution({
    this.sourceUrl,
    this.authorName,
    this.authorUrl,
    super.key,
  });

  /// The original source URL.
  final String? sourceUrl;

  /// The author name.
  final String? authorName;

  /// The author's URL.
  final String? authorUrl;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    if (sourceUrl == null && authorName == null) {
      return const SizedBox.shrink();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.source,
          style: theme.typography.sm.copyWith(
            fontWeight: FontWeight.w600,
            color: theme.colors.mutedForeground,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: theme.colors.secondary,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Row(
            children: [
              Icon(
                Icons.link,
                size: 20,
                color: theme.colors.mutedForeground,
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (authorName != null)
                      Text(
                        authorName!,
                        style: theme.typography.sm.copyWith(
                          fontWeight: FontWeight.w500,
                          color: theme.colors.foreground,
                        ),
                      ),
                    if (sourceUrl != null)
                      GestureDetector(
                        onTap: () => _launchUrl(sourceUrl!),
                        child: Text(
                          _formatUrl(sourceUrl!),
                          style: theme.typography.xs.copyWith(
                            color: theme.colors.primary,
                            decoration: TextDecoration.underline,
                          ),
                          maxLines: 1,
                          overflow: TextOverflow.ellipsis,
                        ),
                      ),
                  ],
                ),
              ),
              if (sourceUrl != null)
                IconButton(
                  icon: Icon(
                    Icons.open_in_new,
                    size: 18,
                    color: theme.colors.mutedForeground,
                  ),
                  onPressed: () => _launchUrl(sourceUrl!),
                ),
            ],
          ),
        ),
      ],
    );
  }

  String _formatUrl(String url) {
    try {
      final uri = Uri.parse(url);
      return uri.host;
    } on FormatException {
      return url;
    }
  }

  Future<void> _launchUrl(String url) async {
    final uri = Uri.parse(url);
    if (await canLaunchUrl(uri)) {
      await launchUrl(uri, mode: LaunchMode.externalApplication);
    }
  }
}
