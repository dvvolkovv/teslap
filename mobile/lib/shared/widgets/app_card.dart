import 'package:flutter/material.dart';

import '../../core/theme/app_theme.dart';

/// A styled card container following TeslaPay design tokens.
///
/// Uses 12 px border radius, 16 px internal padding, and a subtle shadow
/// (light mode) or 1 px border (dark mode).
class AppCard extends StatelessWidget {
  const AppCard({
    required this.child,
    this.padding,
    this.margin,
    this.color,
    this.gradient,
    this.borderRadius,
    this.onTap,
    super.key,
  });

  final Widget child;
  final EdgeInsetsGeometry? padding;
  final EdgeInsetsGeometry? margin;
  final Color? color;
  final Gradient? gradient;
  final BorderRadius? borderRadius;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final effectiveRadius =
        borderRadius ?? BorderRadius.circular(AppRadius.md);

    final decoration = BoxDecoration(
      color: gradient == null
          ? (color ??
              (isDark ? AppColors.darkSurface : AppColors.white))
          : null,
      gradient: gradient,
      borderRadius: effectiveRadius,
      border: isDark && gradient == null
          ? Border.all(color: AppColors.darkBorder)
          : null,
      boxShadow: isDark || gradient != null
          ? null
          : const [
              BoxShadow(
                color: Color(0x140D1B2A), // neutral900 at 8%
                blurRadius: 3,
                offset: Offset(0, 1),
              ),
            ],
    );

    final content = Container(
      margin: margin,
      padding: padding ?? const EdgeInsets.all(AppSpacing.cardPadding),
      decoration: decoration,
      child: child,
    );

    if (onTap != null) {
      return Material(
        color: Colors.transparent,
        borderRadius: effectiveRadius,
        clipBehavior: Clip.antiAlias,
        child: InkWell(
          onTap: onTap,
          borderRadius: effectiveRadius,
          child: content,
        ),
      );
    }

    return content;
  }
}
