import 'package:flutter/material.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../core/theme/app_theme.dart';

/// TeslaPay bottom navigation bar with 5 tabs matching the information
/// architecture: Home, Payments, Card, Crypto, Profile.
///
/// Active tab uses bold icon weight and primary color. Inactive tabs use
/// regular weight and neutral color.
class AppBottomNav extends StatelessWidget {
  const AppBottomNav({
    required this.currentIndex,
    required this.onTap,
    super.key,
  });

  final int currentIndex;
  final ValueChanged<int> onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Container(
      decoration: BoxDecoration(
        color: isDark ? AppColors.darkSurface : AppColors.white,
        border: Border(
          top: BorderSide(
            color: isDark ? AppColors.darkBorder : AppColors.neutral200,
          ),
        ),
      ),
      child: SafeArea(
        top: false,
        child: SizedBox(
          height: AppSpacing.bottomNavHeight,
          child: Row(
            children: List.generate(5, (index) {
              final isActive = index == currentIndex;
              return Expanded(
                child: _NavItem(
                  icon: _icon(index, active: false),
                  activeIcon: _icon(index, active: true),
                  label: _label(index),
                  isActive: isActive,
                  onTap: () => onTap(index),
                ),
              );
            }),
          ),
        ),
      ),
    );
  }

  IconData _icon(int index, {required bool active}) {
    return switch (index) {
      0 => active ? PhosphorIconsBold.house : PhosphorIconsRegular.house,
      1 => active
          ? PhosphorIconsBold.arrowsLeftRight
          : PhosphorIconsRegular.arrowsLeftRight,
      2 => active
          ? PhosphorIconsBold.creditCard
          : PhosphorIconsRegular.creditCard,
      3 => active
          ? PhosphorIconsBold.currencyBtc
          : PhosphorIconsRegular.currencyBtc,
      4 => active
          ? PhosphorIconsBold.userCircle
          : PhosphorIconsRegular.userCircle,
      _ => PhosphorIconsRegular.question,
    };
  }

  String _label(int index) {
    return switch (index) {
      0 => 'Home',
      1 => 'Payments',
      2 => 'Card',
      3 => 'Crypto',
      4 => 'Profile',
      _ => '',
    };
  }
}

class _NavItem extends StatelessWidget {
  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.isActive,
    required this.onTap,
  });

  final IconData icon;
  final IconData activeIcon;
  final String label;
  final bool isActive;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final color = isActive ? AppColors.primary500 : AppColors.neutral400;
    final textColor = isActive ? AppColors.primary500 : AppColors.neutral500;

    return Semantics(
      label: label,
      selected: isActive,
      button: true,
      child: InkResponse(
        onTap: onTap,
        highlightShape: BoxShape.rectangle,
        containedInkWell: true,
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              isActive ? activeIcon : icon,
              size: 24,
              color: color,
            ),
            const SizedBox(height: AppSpacing.xxs),
            Text(
              label,
              style: AppTypography.caption.copyWith(color: textColor),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
          ],
        ),
      ),
    );
  }
}
