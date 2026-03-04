import 'package:flutter/material.dart';

import '../../../core/theme/app_theme.dart';

/// A single transaction row matching the wireframe specification.
///
/// Height: 72 px. Left circle icon with merchant initial, centre name + meta,
/// right-aligned amount with status indicator.
class TransactionListTile extends StatelessWidget {
  const TransactionListTile({
    required this.title,
    required this.subtitle,
    required this.amount,
    this.isPositive = false,
    this.status,
    this.iconLetter,
    this.iconColor,
    this.onTap,
    super.key,
  });

  final String title;
  final String subtitle;
  final String amount;
  final bool isPositive;
  final String? status;
  final String? iconLetter;
  final Color? iconColor;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final letter =
        iconLetter ?? (title.isNotEmpty ? title[0].toUpperCase() : '?');

    return Semantics(
      label: '$title, $amount, $subtitle${status != null ? ', $status' : ''}',
      button: onTap != null,
      child: InkWell(
        onTap: onTap,
        child: Container(
          height: 72,
          padding:
              const EdgeInsets.symmetric(horizontal: AppSpacing.md),
          child: Row(
            children: [
              // Merchant avatar
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: iconColor ?? AppColors.primary100,
                  shape: BoxShape.circle,
                ),
                child: Center(
                  child: Text(
                    letter,
                    style: AppTypography.body1.copyWith(
                      fontWeight: FontWeight.w600,
                      color: iconColor != null
                          ? AppColors.white
                          : AppColors.primary500,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: AppSpacing.md),

              // Name and subtitle
              Expanded(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: AppTypography.body1.copyWith(
                        fontWeight: FontWeight.w500,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: AppSpacing.xxs),
                    Text(
                      subtitle,
                      style: AppTypography.caption.copyWith(
                        color: isDark
                            ? AppColors.darkTextSecondary
                            : AppColors.neutral500,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),

              // Amount
              Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    amount,
                    style: AppTypography.body1.copyWith(
                      fontWeight: FontWeight.w600,
                      color: isPositive
                          ? AppColors.success500
                          : null,
                      fontFeatures: const [
                        FontFeature.tabularFigures(),
                      ],
                    ),
                  ),
                  if (status != null) ...[
                    const SizedBox(height: AppSpacing.xxs),
                    _StatusBadge(status: status!),
                  ],
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _StatusBadge extends StatelessWidget {
  const _StatusBadge({required this.status});

  final String status;

  Color get _color {
    return switch (status.toLowerCase()) {
      'completed' || 'settled' => AppColors.success500,
      'pending' || 'processing' => AppColors.warning500,
      'failed' || 'declined' => AppColors.error500,
      _ => AppColors.info500,
    };
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.sm,
        vertical: AppSpacing.xxs,
      ),
      decoration: BoxDecoration(
        color: _color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(AppRadius.xs),
      ),
      child: Text(
        status,
        style: AppTypography.caption.copyWith(
          color: _color,
          fontWeight: FontWeight.w500,
          fontSize: 10,
        ),
      ),
    );
  }
}
