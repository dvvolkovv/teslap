import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_card.dart';
import '../widgets/balance_card.dart';
import '../widgets/quick_actions.dart';
import '../widgets/transaction_list_tile.dart';

/// Home dashboard matching the wireframe from the design system.
///
/// Sections:
/// 1. Balance card (hero)
/// 2. Quick actions row
/// 3. Accounts overview
/// 4. Recent transactions (last 5)
/// 5. Crypto portfolio summary (if activated)
class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: RefreshIndicator(
          color: AppColors.primary500,
          onRefresh: () async {
            // TODO: reload accounts, transactions, crypto balance
            await Future<void>.delayed(const Duration(milliseconds: 800));
          },
          child: CustomScrollView(
            slivers: [
              // App bar
              SliverAppBar(
                floating: true,
                snap: true,
                backgroundColor: Theme.of(context).scaffoldBackgroundColor,
                title: Text(
                  'TeslaPay',
                  style: AppTypography.h2.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
                ),
                centerTitle: false,
                actions: [
                  IconButton(
                    icon: const Icon(PhosphorIconsRegular.bellSimple),
                    onPressed: () {
                      // TODO: navigate to notifications
                    },
                    tooltip: 'Notifications',
                  ),
                ],
              ),

              SliverPadding(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSpacing.screenMargin,
                ),
                sliver: SliverList(
                  delegate: SliverChildListDelegate([
                    // -------------------------------------------------------
                    // Balance Card
                    // -------------------------------------------------------
                    const BalanceCard(
                      totalBalance: '3,245.67',
                      currency: 'EUR',
                      changePercent: '+0.42',
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // -------------------------------------------------------
                    // Quick Actions
                    // -------------------------------------------------------
                    QuickActions(
                      onSend: () => context.push(AppRoutes.sendMoney),
                      onRequest: () {
                        // TODO: navigate to request money
                      },
                      onExchange: () {
                        // TODO: navigate to exchange
                      },
                      onTopUp: () {
                        // TODO: show top-up bottom sheet
                      },
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // -------------------------------------------------------
                    // Accounts Overview
                    // -------------------------------------------------------
                    _SectionHeader(
                      title: 'Accounts',
                      actionLabel: 'See all',
                      onAction: () {
                        // TODO: navigate to accounts list
                      },
                    ),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: const EdgeInsets.symmetric(
                        vertical: AppSpacing.sm,
                      ),
                      child: Column(
                        children: [
                          _AccountRow(
                            flag: 'EU',
                            currency: 'EUR',
                            balance: 'EUR 2,845.67',
                            onTap: () {
                              // TODO: navigate to EUR account detail
                            },
                          ),
                          const Divider(
                            indent: AppSpacing.xxl + AppSpacing.md,
                            height: 1,
                          ),
                          _AccountRow(
                            flag: 'US',
                            currency: 'USD',
                            balance: 'USD 420.00',
                            onTap: () {
                              // TODO: navigate to USD account detail
                            },
                          ),
                          const Divider(
                            indent: AppSpacing.xxl + AppSpacing.md,
                            height: 1,
                          ),
                          ListTile(
                            leading: Container(
                              width: 32,
                              height: 32,
                              decoration: const BoxDecoration(
                                color: AppColors.primary50,
                                shape: BoxShape.circle,
                              ),
                              child: const Icon(
                                PhosphorIconsRegular.plus,
                                size: 16,
                                color: AppColors.primary500,
                              ),
                            ),
                            title: Text(
                              'Add currency',
                              style: AppTypography.body2.copyWith(
                                color: AppColors.primary500,
                              ),
                            ),
                            dense: true,
                            onTap: () {
                              // TODO: show add currency bottom sheet
                            },
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // -------------------------------------------------------
                    // Recent Transactions
                    // -------------------------------------------------------
                    _SectionHeader(
                      title: 'Recent Transactions',
                      actionLabel: 'See all',
                      onAction: () {
                        // TODO: navigate to full transaction history
                      },
                    ),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          TransactionListTile(
                            title: 'Starbucks',
                            subtitle: 'Card payment  Today 09:15',
                            amount: '-EUR 4.50',
                            onTap: () {},
                          ),
                          const Divider(
                            indent: AppSpacing.xxl + AppSpacing.lg,
                            height: 1,
                          ),
                          TransactionListTile(
                            title: 'Anna Kowalski',
                            subtitle: 'SEPA transfer  Yesterday',
                            amount: '-EUR 200.00',
                            onTap: () {},
                          ),
                          const Divider(
                            indent: AppSpacing.xxl + AppSpacing.lg,
                            height: 1,
                          ),
                          TransactionListTile(
                            title: 'Monthly Salary',
                            subtitle: 'SEPA received  28 Feb',
                            amount: '+EUR 3,200.00',
                            isPositive: true,
                            onTap: () {},
                          ),
                          const Divider(
                            indent: AppSpacing.xxl + AppSpacing.lg,
                            height: 1,
                          ),
                          TransactionListTile(
                            title: 'FUSE Purchase',
                            subtitle: 'Crypto buy  28 Feb',
                            amount: '-EUR 20.00',
                            onTap: () {},
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // -------------------------------------------------------
                    // Crypto Portfolio Summary
                    // -------------------------------------------------------
                    _SectionHeader(
                      title: 'Crypto Portfolio',
                      actionLabel: 'View',
                      onAction: () => context.go(AppRoutes.crypto),
                    ),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      gradient: AppColors.cryptoGradient,
                      onTap: () => context.go(AppRoutes.crypto),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'Total',
                                style: AppTypography.caption.copyWith(
                                  color:
                                      AppColors.white.withValues(alpha: 0.8),
                                ),
                              ),
                              const SizedBox(height: AppSpacing.xxs),
                              Text(
                                'EUR 19.70',
                                style: AppTypography.h3.copyWith(
                                  color: AppColors.white,
                                ),
                              ),
                            ],
                          ),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: AppSpacing.sm,
                              vertical: AppSpacing.xs,
                            ),
                            decoration: BoxDecoration(
                              color: AppColors.white.withValues(alpha: 0.2),
                              borderRadius:
                                  BorderRadius.circular(AppRadius.xl),
                            ),
                            child: Text(
                              '+3.2% 24h',
                              style: AppTypography.caption.copyWith(
                                color: AppColors.white,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
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
      ),
    );
  }
}

// ---------------------------------------------------------------------------
// Private helpers
// ---------------------------------------------------------------------------

class _SectionHeader extends StatelessWidget {
  const _SectionHeader({
    required this.title,
    this.actionLabel,
    this.onAction,
  });

  final String title;
  final String? actionLabel;
  final VoidCallback? onAction;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(title, style: AppTypography.h3),
        if (actionLabel != null)
          GestureDetector(
            onTap: onAction,
            child: Text(
              '$actionLabel >',
              style: AppTypography.body2.copyWith(
                color: AppColors.primary500,
              ),
            ),
          ),
      ],
    );
  }
}

class _AccountRow extends StatelessWidget {
  const _AccountRow({
    required this.flag,
    required this.currency,
    required this.balance,
    this.onTap,
  });

  final String flag;
  final String currency;
  final String balance;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      onTap: onTap,
      leading: Container(
        width: 32,
        height: 32,
        decoration: const BoxDecoration(
          color: AppColors.neutral100,
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Text(
            flag,
            style: AppTypography.caption.copyWith(
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
      title: Text(currency, style: AppTypography.body1),
      trailing: Text(
        balance,
        style: AppTypography.body1.copyWith(
          fontWeight: FontWeight.w600,
          fontFeatures: const [FontFeature.tabularFigures()],
        ),
      ),
      dense: true,
    );
  }
}
