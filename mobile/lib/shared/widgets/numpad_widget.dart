import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../core/theme/app_theme.dart';

/// A custom number pad used for PIN entry and amount input.
///
/// Provides digits 0-9, a backspace key, and an optional biometric button
/// in the bottom-left position.
class NumpadWidget extends StatelessWidget {
  const NumpadWidget({
    required this.onDigit,
    required this.onBackspace,
    this.onBiometric,
    this.showBiometric = false,
    this.biometricIcon,
    super.key,
  });

  final ValueChanged<int> onDigit;
  final VoidCallback onBackspace;
  final VoidCallback? onBiometric;
  final bool showBiometric;
  final IconData? biometricIcon;

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        _buildRow([1, 2, 3]),
        const SizedBox(height: AppSpacing.sm),
        _buildRow([4, 5, 6]),
        const SizedBox(height: AppSpacing.sm),
        _buildRow([7, 8, 9]),
        const SizedBox(height: AppSpacing.sm),
        _buildBottomRow(context),
      ],
    );
  }

  Widget _buildRow(List<int> digits) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: digits
          .map(
            (d) => _NumpadKey(
              label: d.toString(),
              onTap: () {
                HapticFeedback.lightImpact();
                onDigit(d);
              },
            ),
          )
          .toList(),
    );
  }

  Widget _buildBottomRow(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        // Bottom-left: biometric or empty
        if (showBiometric && onBiometric != null)
          _NumpadKey(
            icon: biometricIcon ?? PhosphorIconsRegular.fingerprint,
            onTap: () {
              HapticFeedback.mediumImpact();
              onBiometric!();
            },
            semanticLabel: 'Use biometric authentication',
          )
        else
          const SizedBox(width: 80, height: 64),

        // Zero
        _NumpadKey(
          label: '0',
          onTap: () {
            HapticFeedback.lightImpact();
            onDigit(0);
          },
        ),

        // Backspace
        _NumpadKey(
          icon: PhosphorIconsRegular.backspace,
          onTap: () {
            HapticFeedback.lightImpact();
            onBackspace();
          },
          semanticLabel: 'Delete',
        ),
      ],
    );
  }
}

class _NumpadKey extends StatelessWidget {
  const _NumpadKey({
    this.label,
    this.icon,
    required this.onTap,
    this.semanticLabel,
  });

  final String? label;
  final IconData? icon;
  final VoidCallback onTap;
  final String? semanticLabel;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label: semanticLabel ?? label,
      button: true,
      child: Material(
        color: Colors.transparent,
        borderRadius: BorderRadius.circular(AppRadius.full),
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(AppRadius.full),
          splashColor: AppColors.primary100,
          child: SizedBox(
            width: 80,
            height: 64,
            child: Center(
              child: label != null
                  ? Text(
                      label!,
                      style: AppTypography.h1.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
                    )
                  : Icon(
                      icon,
                      size: 28,
                      color: AppColors.neutral700,
                    ),
            ),
          ),
        ),
      ),
    );
  }
}
