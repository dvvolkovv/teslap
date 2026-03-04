import 'package:flutter/material.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';

/// Crypto tab showing wallet balance, token cards, quick actions,
/// and DeFi yield section (Phase 2 placeholder).
class CryptoScreen extends StatelessWidget {
  const CryptoScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: CustomScrollView(
          slivers: [
            SliverAppBar(
              floating: true,
              snap: true,
              backgroundColor: Theme.of(context).scaffoldBackgroundColor,
              title: Text(
                'Crypto',
                style: AppTypography.h2.copyWith(fontWeight: FontWeight.w700),
              ),
              centerTitle: false,
            ),
            SliverPadding(
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.screenMargin,
              ),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  // Total crypto balance hero
                  AppCard(
                    gradient: AppColors.cryptoGradient,
                    borderRadius: BorderRadius.circular(AppRadius.lg),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Total Crypto Balance',
                          style: AppTypography.body2.copyWith(
                            color: AppColors.white.withValues(alpha: 0.8),
                          ),
                        ),
                        const SizedBox(height: AppSpacing.sm),
                        Text(
                          'EUR 467.53',
                          style: AppTypography.display2.copyWith(
                            color: AppColors.white,
                          ),
                        ),
                        const SizedBox(height: AppSpacing.xs),
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: AppSpacing.sm,
                            vertical: AppSpacing.xxs,
                          ),
                          decoration: BoxDecoration(
                            color: AppColors.white.withValues(alpha: 0.2),
                            borderRadius:
                                BorderRadius.circular(AppRadius.xl),
                          ),
                          child: Text(
                            '-0.85% 24h',
                            style: AppTypography.caption.copyWith(
                              color: AppColors.white,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Token balances
                  Text('Tokens', style: AppTypography.h3),
                  const SizedBox(height: AppSpacing.sm),
                  _TokenCard(
                    symbol: 'FUSE',
                    name: 'Fuse Token',
                    balance: '150.50',
                    eurValue: 'EUR 7.53',
                    change: '-2.50',
                    isNegative: true,
                    onTap: () {},
                  ),
                  const SizedBox(height: AppSpacing.sm),
                  _TokenCard(
                    symbol: 'USDC',
                    name: 'USD Coin',
                    balance: '500.00',
                    eurValue: 'EUR 460.00',
                    change: '+0.01',
                    isNegative: false,
                    onTap: () {},
                  ),
                  const SizedBox(height: AppSpacing.sm),
                  _TokenCard(
                    symbol: 'USDT',
                    name: 'Tether',
                    balance: '0.00',
                    eurValue: 'EUR 0.00',
                    change: '-0.02',
                    isNegative: true,
                    onTap: () {},
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Quick actions
                  Row(
                    children: [
                      Expanded(
                        child: AppButton(
                          label: 'Buy',
                          onPressed: () {},
                          isFullWidth: true,
                        ),
                      ),
                      const SizedBox(width: AppSpacing.md),
                      Expanded(
                        child: AppButton(
                          label: 'Sell',
                          variant: AppButtonVariant.secondary,
                          onPressed: () {},
                          isFullWidth: true,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSpacing.sm),
                  Row(
                    children: [
                      Expanded(
                        child: AppButton(
                          label: 'Send',
                          variant: AppButtonVariant.outline,
                          onPressed: () {},
                          isFullWidth: true,
                        ),
                      ),
                      const SizedBox(width: AppSpacing.md),
                      Expanded(
                        child: AppButton(
                          label: 'Receive',
                          variant: AppButtonVariant.outline,
                          onPressed: () {},
                          isFullWidth: true,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // DeFi Yield (Phase 2)
                  Text('DeFi Yield', style: AppTypography.h3),
                  const SizedBox(height: AppSpacing.sm),
                  AppCard(
                    child: Row(
                      children: [
                        Container(
                          width: 48,
                          height: 48,
                          decoration: const BoxDecoration(
                            color: AppColors.accent100,
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(
                            PhosphorIconsRegular.trendUp,
                            color: AppColors.accent500,
                            size: 24,
                          ),
                        ),
                        const SizedBox(width: AppSpacing.md),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Earn with Solid soUSD',
                                style: AppTypography.body1.copyWith(
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                              Text(
                                'Current APY: 4.8%',
                                style: AppTypography.body2.copyWith(
                                  color: AppColors.accent500,
                                ),
                              ),
                            ],
                          ),
                        ),
                        const Icon(
                          Icons.chevron_right,
                          color: AppColors.neutral400,
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Wallet info
                  Text('Wallet', style: AppTypography.h3),
                  const SizedBox(height: AppSpacing.sm),
                  AppCard(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Wallet Address',
                          style: AppTypography.caption.copyWith(
                            color: AppColors.neutral500,
                          ),
                        ),
                        const SizedBox(height: AppSpacing.xs),
                        Row(
                          children: [
                            Expanded(
                              child: Text(
                                '0x742d35Cc...f2BD25',
                                style: AppTypography.mono,
                              ),
                            ),
                            IconButton(
                              icon: const Icon(
                                PhosphorIconsRegular.copy,
                                size: 20,
                              ),
                              onPressed: () {},
                              tooltip: 'Copy address',
                            ),
                            IconButton(
                              icon: const Icon(
                                PhosphorIconsRegular.qrCode,
                                size: 20,
                              ),
                              onPressed: () {},
                              tooltip: 'Show QR code',
                            ),
                          ],
                        ),
                        const SizedBox(height: AppSpacing.sm),
                        Row(
                          children: [
                            Container(
                              width: 8,
                              height: 8,
                              decoration: const BoxDecoration(
                                shape: BoxShape.circle,
                                color: AppColors.success500,
                              ),
                            ),
                            const SizedBox(width: AppSpacing.xs),
                            Text(
                              'Fuse Network - Connected',
                              style: AppTypography.caption.copyWith(
                                color: AppColors.success500,
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSpacing.xxl),
                ]),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _TokenCard extends StatelessWidget {
  const _TokenCard({
    required this.symbol,
    required this.name,
    required this.balance,
    required this.eurValue,
    required this.change,
    required this.isNegative,
    this.onTap,
  });

  final String symbol;
  final String name;
  final String balance;
  final String eurValue;
  final String change;
  final bool isNegative;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return AppCard(
      onTap: onTap,
      child: Row(
        children: [
          // Token icon placeholder
          Container(
            width: 40,
            height: 40,
            decoration: const BoxDecoration(
              color: AppColors.accent100,
              shape: BoxShape.circle,
            ),
            child: Center(
              child: Text(
                symbol[0],
                style: AppTypography.body1.copyWith(
                  fontWeight: FontWeight.w700,
                  color: AppColors.accent500,
                ),
              ),
            ),
          ),
          const SizedBox(width: AppSpacing.md),

          // Name and balance in token
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  name,
                  style: AppTypography.body1.copyWith(
                    fontWeight: FontWeight.w500,
                  ),
                ),
                Text(
                  '$balance $symbol',
                  style: AppTypography.h3,
                ),
              ],
            ),
          ),

          // EUR value and change
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                eurValue,
                style: AppTypography.body2.copyWith(
                  color: AppColors.neutral500,
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSpacing.xs,
                  vertical: AppSpacing.xxs,
                ),
                decoration: BoxDecoration(
                  color: (isNegative
                          ? AppColors.error500
                          : AppColors.success500)
                      .withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(AppRadius.xs),
                ),
                child: Text(
                  '${change}%',
                  style: AppTypography.caption.copyWith(
                    color: isNegative
                        ? AppColors.error500
                        : AppColors.success500,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
