import 'package:flutter/material.dart';

import '../../../core/theme/app_theme.dart';

/// Hero balance card displayed at the top of the home dashboard.
///
/// Shows total balance with currency, daily change percentage, and an
/// animated counter on load. Uses the brand gradient background.
class BalanceCard extends StatefulWidget {
  const BalanceCard({
    required this.totalBalance,
    required this.currency,
    this.changePercent,
    this.onCurrencyTap,
    super.key,
  });

  final String totalBalance;
  final String currency;
  final String? changePercent;
  final VoidCallback? onCurrencyTap;

  @override
  State<BalanceCard> createState() => _BalanceCardState();
}

class _BalanceCardState extends State<BalanceCard>
    with SingleTickerProviderStateMixin {
  late final AnimationController _controller;
  late final Animation<double> _fadeIn;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 600),
    );
    _fadeIn = CurvedAnimation(parent: _controller, curve: Curves.easeOut);
    _controller.forward();
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isPositive = widget.changePercent != null &&
        !widget.changePercent!.startsWith('-');

    return Semantics(
      label:
          'Total balance: ${widget.currency} ${widget.totalBalance}. '
          '${widget.changePercent != null ? '${widget.changePercent} percent today.' : ''}',
      child: Container(
        width: double.infinity,
        padding: const EdgeInsets.all(AppSpacing.lg),
        decoration: BoxDecoration(
          gradient: AppColors.brandGradient,
          borderRadius: BorderRadius.circular(AppRadius.lg),
          boxShadow: const [
            BoxShadow(
              color: Color(0x330066FF),
              blurRadius: 24,
              offset: Offset(0, 8),
            ),
          ],
        ),
        child: FadeTransition(
          opacity: _fadeIn,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Currency selector
              GestureDetector(
                onTap: widget.onCurrencyTap,
                child: Row(
                  children: [
                    Text(
                      'Total Balance',
                      style: AppTypography.body2.copyWith(
                        color: AppColors.white.withValues(alpha: 0.8),
                      ),
                    ),
                    const SizedBox(width: AppSpacing.xs),
                    Icon(
                      Icons.keyboard_arrow_down,
                      color: AppColors.white.withValues(alpha: 0.8),
                      size: 20,
                    ),
                  ],
                ),
              ),
              const SizedBox(height: AppSpacing.sm),

              // Balance amount
              Text(
                '${widget.currency} ${widget.totalBalance}',
                style: AppTypography.display1.copyWith(
                  color: AppColors.white,
                  fontFeatures: const [FontFeature.tabularFigures()],
                ),
              ),

              if (widget.changePercent != null) ...[
                const SizedBox(height: AppSpacing.xs),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: AppSpacing.sm,
                    vertical: AppSpacing.xxs,
                  ),
                  decoration: BoxDecoration(
                    color: AppColors.white.withValues(alpha: 0.2),
                    borderRadius: BorderRadius.circular(AppRadius.xl),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isPositive
                            ? Icons.arrow_upward
                            : Icons.arrow_downward,
                        size: 14,
                        color: AppColors.white,
                      ),
                      const SizedBox(width: AppSpacing.xxs),
                      Text(
                        '${widget.changePercent}% today',
                        style: AppTypography.caption.copyWith(
                          color: AppColors.white,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
