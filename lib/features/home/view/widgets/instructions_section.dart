import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/l10n/l10n.dart';

/// Widget displaying the instructions list section.
class InstructionsSection extends StatelessWidget {
  /// Creates a new [InstructionsSection] instance.
  const InstructionsSection({
    required this.instructions,
    super.key,
  });

  /// List of cooking instructions.
  final List<String> instructions;

  @override
  Widget build(BuildContext context) {
    final theme = context.theme;
    final l10n = context.l10n;

    if (instructions.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          l10n.instructions,
          style: theme.typography.lg.copyWith(
            fontWeight: FontWeight.w600,
            color: theme.colors.foreground,
          ),
        ),
        const SizedBox(height: 12),
        ...instructions.asMap().entries.map((entry) {
          final index = entry.key + 1;
          final instruction = entry.value;
          return Padding(
            padding: const EdgeInsets.only(bottom: 16),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  width: 28,
                  height: 28,
                  margin: const EdgeInsets.only(right: 12),
                  decoration: BoxDecoration(
                    color: theme.colors.primary,
                    shape: BoxShape.circle,
                  ),
                  child: Center(
                    child: Text(
                      '$index',
                      style: theme.typography.sm.copyWith(
                        color: theme.colors.primaryForeground,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
                Expanded(
                  child: Padding(
                    padding: const EdgeInsets.only(top: 4),
                    child: Text(
                      instruction,
                      style: theme.typography.base.copyWith(
                        color: theme.colors.foreground,
                        height: 1.5,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          );
        }),
      ],
    );
  }
}
