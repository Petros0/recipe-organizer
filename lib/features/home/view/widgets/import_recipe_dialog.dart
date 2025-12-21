import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// Dialog for importing a recipe from a URL.
class ImportRecipeDialog extends StatefulWidget {
  /// Creates a new [ImportRecipeDialog] instance.
  const ImportRecipeDialog({
    required this.onImport,
    this.isLoading = false,
    super.key,
  });

  /// Callback when a URL is submitted for import.
  final void Function(String url) onImport;

  /// Whether an import is currently in progress.
  final bool isLoading;

  @override
  State<ImportRecipeDialog> createState() => _ImportRecipeDialogState();
}

class _ImportRecipeDialogState extends State<ImportRecipeDialog> {
  final TextEditingController _urlController = TextEditingController();
  String? _validationError;

  @override
  void dispose() {
    _urlController.dispose();
    super.dispose();
  }

  bool _isValidUrl(String url) {
    if (url.isEmpty) return false;
    try {
      final uri = Uri.parse(url);
      return uri.hasScheme && (uri.scheme == 'http' || uri.scheme == 'https') && uri.host.isNotEmpty;
    } on FormatException {
      return false;
    }
  }

  void _handleImport() {
    final url = _urlController.text.trim();

    if (url.isEmpty) {
      setState(() {
        _validationError = context.l10n.invalidUrl;
      });
      return;
    }

    if (!_isValidUrl(url)) {
      setState(() {
        _validationError = context.l10n.invalidUrlFormat;
      });
      return;
    }

    setState(() {
      _validationError = null;
    });
    widget.onImport(url);
  }

  @override
  Widget build(BuildContext context) {
    final l10n = context.l10n;
    final theme = context.theme;

    return Container(
      constraints: const BoxConstraints(maxWidth: 400),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            l10n.importDialogTitle,
            style: theme.typography.xl.copyWith(
              fontWeight: FontWeight.w600,
              color: theme.colors.foreground,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            l10n.importDialogDescription,
            style: theme.typography.sm.copyWith(
              color: theme.colors.mutedForeground,
            ),
          ),
          const SizedBox(height: 24),
          FTextField(
            controller: _urlController,
            hint: 'https://example.com/recipe',
            label: Text(l10n.urlLabel),
            enabled: !widget.isLoading,
            keyboardType: TextInputType.url,
            textInputAction: TextInputAction.done,
            onSubmit: (_) => _handleImport(),
            onChange: (_) {
              if (_validationError != null) {
                setState(() {
                  _validationError = null;
                });
              }
            },
          ),
          if (_validationError != null) ...[
            const SizedBox(height: 8),
            Text(
              _validationError!,
              style: theme.typography.sm.copyWith(
                color: theme.colors.destructive,
              ),
            ),
          ],
          const SizedBox(height: 24),
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              FButton(
                style: FButtonStyle.outline(),
                onPress: widget.isLoading ? null : () => Navigator.of(context).pop(),
                child: Text(l10n.cancelButton),
              ),
              const SizedBox(width: 12),
              FButton(
                onPress: widget.isLoading ? null : _handleImport,
                child: widget.isLoading
                    ? SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          color: theme.colors.primaryForeground,
                        ),
                      )
                    : Text(l10n.importButton),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

/// Shows the import recipe dialog.
Future<void> showImportRecipeDialog({
  required BuildContext context,
  required void Function(String url) onImport,
  bool isLoading = false,
}) {
  return showFDialog(
    context: context,
    builder: (context, style, animation) => FDialog.raw(
      style: style.call,
      animation: animation,
      builder: (context, style) => Padding(
        padding: const EdgeInsets.all(24),
        child: ImportRecipeDialog(
          onImport: onImport,
          isLoading: isLoading,
        ),
      ),
    ),
  );
}
