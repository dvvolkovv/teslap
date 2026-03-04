import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';

/// Button variant determines visual style.
enum AppButtonVariant { primary, secondary, text, danger, outline }

/// A design-system compliant button with loading state support.
///
/// Follows TeslaPay design tokens: 56 px height, 8 px radius,
/// SemiBold 600 label, full-width by default.
class AppButton extends StatelessWidget {
  const AppButton({
    required this.label,
    required this.onPressed,
    this.variant = AppButtonVariant.primary,
    this.isLoading = false,
    this.isFullWidth = true,
    this.icon,
    super.key,
  });

  final String label;
  final VoidCallback? onPressed;
  final AppButtonVariant variant;
  final bool isLoading;
  final bool isFullWidth;
  final IconData? icon;

  @override
  Widget build(BuildContext context) {
    final child = isLoading
        ? const SizedBox(
            width: 24,
            height: 24,
            child: CircularProgressIndicator(
              strokeWidth: 2.5,
              valueColor: AlwaysStoppedAnimation<Color>(AppColors.white),
            ),
          )
        : _buildLabel();

    final effectiveOnPressed = isLoading ? null : onPressed;

    Widget button;
    switch (variant) {
      case AppButtonVariant.primary:
        button = ElevatedButton(
          onPressed: effectiveOnPressed,
          child: child,
        );
      case AppButtonVariant.secondary:
        button = OutlinedButton(
          onPressed: effectiveOnPressed,
          child: child,
        );
      case AppButtonVariant.text:
        button = TextButton(
          onPressed: effectiveOnPressed,
          child: child,
        );
      case AppButtonVariant.danger:
        button = ElevatedButton(
          onPressed: effectiveOnPressed,
          style: ElevatedButton.styleFrom(
            backgroundColor: AppColors.error500,
            foregroundColor: AppColors.white,
            minimumSize: isFullWidth
                ? const Size(double.infinity, 56)
                : const Size(0, 56),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.sm),
            ),
            textStyle: AppTypography.button,
            elevation: 0,
          ),
          child: child,
        );
      case AppButtonVariant.outline:
        button = OutlinedButton(
          onPressed: effectiveOnPressed,
          style: OutlinedButton.styleFrom(
            foregroundColor: AppColors.neutral700,
            minimumSize: isFullWidth
                ? const Size(double.infinity, 56)
                : const Size(0, 56),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(AppRadius.sm),
            ),
            side: const BorderSide(color: AppColors.neutral200, width: 1.5),
            textStyle: AppTypography.button,
          ),
          child: child,
        );
    }

    if (!isFullWidth) {
      return button;
    }

    return SizedBox(width: double.infinity, child: button);
  }

  Widget _buildLabel() {
    if (icon != null) {
      return Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 20),
          const SizedBox(width: AppSpacing.sm),
          Text(label),
        ],
      );
    }
    return Text(label);
  }
}
