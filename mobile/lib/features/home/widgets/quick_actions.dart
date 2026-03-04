import 'package:flutter/material.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/theme/app_theme.dart';

/// A row of circular quick-action buttons: Send, Request, Exchange, Top Up.
///
/// Each button has a 48 px circular icon container with primary-100 background
/// and a caption label below.
class QuickActions extends StatelessWidget {
  const QuickActions({
    this.onSend,
    this.onRequest,
    this.onExchange,
    this.onTopUp,
    super.key,
  });

  final VoidCallback? onSend;
  final VoidCallback? onRequest;
  final VoidCallback? onExchange;
  final VoidCallback? onTopUp;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        _QuickActionItem(
          icon: PhosphorIconsRegular.arrowUp,
          label: 'Send',
          onTap: onSend,
        ),
        _QuickActionItem(
          icon: PhosphorIconsRegular.arrowDown,
          label: 'Request',
          onTap: onRequest,
        ),
        _QuickActionItem(
          icon: PhosphorIconsRegular.arrowsLeftRight,
          label: 'Exchange',
          onTap: onExchange,
        ),
        _QuickActionItem(
          icon: PhosphorIconsRegular.plus,
          label: 'Top Up',
          onTap: onTopUp,
        ),
      ],
    );
  }
}

class _QuickActionItem extends StatelessWidget {
  const _QuickActionItem({
    required this.icon,
    required this.label,
    this.onTap,
  });

  final IconData icon;
  final String label;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label: label,
      button: true,
      child: GestureDetector(
        onTap: onTap,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 48,
              height: 48,
              decoration: const BoxDecoration(
                color: AppColors.primary100,
                shape: BoxShape.circle,
              ),
              child: Icon(icon, size: 24, color: AppColors.primary500),
            ),
            const SizedBox(height: AppSpacing.xs),
            Text(
              label,
              style: AppTypography.caption.copyWith(
                color: AppColors.neutral700,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
