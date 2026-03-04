import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_card.dart';
import '../bloc/home_bloc.dart';
import '../widgets/balance_card.dart';
import '../widgets/quick_actions.dart';
import '../widgets/transaction_list_tile.dart';

/// Home dashboard — loads real data from the backend via [HomeBloc].
class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  void initState() {
    super.initState();
    context.read<HomeBloc>().add(const HomeLoadRequested());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: BlocBuilder<HomeBloc, HomeState>(
          builder: (context, state) {
            return RefreshIndicator(
              color: AppColors.primary500,
              onRefresh: () async {
                context.read<HomeBloc>().add(const HomeRefreshRequested());
                await context.read<HomeBloc>().stream.firstWhere(
                      (s) => s is! HomeLoading,
                    );
              },
              child: CustomScrollView(
                slivers: [
                  SliverAppBar(
                    floating: true,
                    snap: true,
                    backgroundColor:
                        Theme.of(context).scaffoldBackgroundColor,
                    title: Text(
                      'TeslaPay',
                      style: AppTypography.h2
                          .copyWith(fontWeight: FontWeight.w700),
                    ),
                    centerTitle: false,
                    actions: [
                      IconButton(
                        icon:
                            const Icon(PhosphorIconsRegular.bellSimple),
                        onPressed: () {},
                        tooltip: 'Notifications',
                      ),
                    ],
                  ),
                  if (state is HomeLoading)
                    const SliverFillRemaining(
                      child:
                          Center(child: CircularProgressIndicator()),
                    )
                  else if (state is HomeError)
                    SliverFillRemaining(
                      child: _ErrorView(
                        message: state.message,
                        onRetry: () => context
                            .read<HomeBloc>()
                            .add(const HomeLoadRequested()),
                      ),
                    )
                  else if (state is HomeLoaded)
                    _LoadedContent(state: state),
                ],
              ),
            );
          },
        ),
      ),
    );
  }
}

class _ErrorView extends StatelessWidget {
  const _ErrorView({required this.message, required this.onRetry});

  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.screenMargin),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(PhosphorIconsRegular.warning,
                size: 48, color: AppColors.error500),
            const SizedBox(height: AppSpacing.md),
            Text('Failed to load data', style: AppTypography.h3),
            const SizedBox(height: AppSpacing.sm),
            Text(message,
                style: AppTypography.body2
                    .copyWith(color: AppColors.neutral500),
                textAlign: TextAlign.center),
            const SizedBox(height: AppSpacing.lg),
            TextButton(onPressed: onRetry, child: const Text('Retry')),
          ],
        ),
      ),
    );
  }
}

class _LoadedContent extends StatelessWidget {
  const _LoadedContent({required this.state});

  final HomeLoaded state;

  @override
  Widget build(BuildContext context) {
    final accounts = state.accounts;
    final transactions = state.transactions;

    return SliverPadding(
      padding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.screenMargin),
      sliver: SliverList(
        delegate: SliverChildListDelegate([
          BalanceCard(
            totalBalance: state.totalBalanceEur,
            currency: 'EUR',
            changePercent: '',
          ),
          const SizedBox(height: AppSpacing.lg),
          QuickActions(
            onSend: () => context.push(AppRoutes.sendMoney),
            onRequest: () {},
            onExchange: () {},
            onTopUp: () {},
          ),
          const SizedBox(height: AppSpacing.lg),
          _SectionHeader(
              title: 'Accounts',
              actionLabel: 'See all',
              onAction: () {}),
          const SizedBox(height: AppSpacing.sm),
          _AccountsCard(accounts: accounts),
          const SizedBox(height: AppSpacing.lg),
          _SectionHeader(
              title: 'Recent Transactions',
              actionLabel: 'See all',
              onAction: () {}),
          const SizedBox(height: AppSpacing.sm),
          _TransactionsCard(transactions: transactions),
          const SizedBox(height: AppSpacing.lg),
          _SectionHeader(
              title: 'Crypto Portfolio',
              actionLabel: 'View',
              onAction: () => context.go(AppRoutes.crypto)),
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
                    Text('Total',
                        style: AppTypography.caption.copyWith(
                            color:
                                AppColors.white.withValues(alpha: 0.8))),
                    const SizedBox(height: AppSpacing.xxs),
                    Text('View portfolio',
                        style: AppTypography.h3
                            .copyWith(color: AppColors.white)),
                  ],
                ),
                const Icon(Icons.chevron_right, color: AppColors.white),
              ],
            ),
          ),
          const SizedBox(height: AppSpacing.xxl),
        ]),
      ),
    );
  }
}

class _AccountsCard extends StatelessWidget {
  const _AccountsCard({required this.accounts});

  final List<dynamic> accounts;

  @override
  Widget build(BuildContext context) {
    if (accounts.isEmpty) {
      return AppCard(
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.lg),
          child: Center(
            child: Text('No accounts yet',
                style: AppTypography.body2
                    .copyWith(color: AppColors.neutral500)),
          ),
        ),
      );
    }

    final widgets = <Widget>[];
    for (final acc in accounts) {
      final map = acc as Map<String, dynamic>;
      final subs = map['sub_accounts'] as List<dynamic>? ?? [];
      for (final sub in subs) {
        final s = sub as Map<String, dynamic>;
        final cur = s['currency'] as String? ?? 'EUR';
        final bal = s['balance'] as Map<String, dynamic>?;
        final avail = bal?['available']?.toString() ?? '0.00';
        final flag = cur == 'EUR'
            ? 'EU'
            : cur == 'USD'
                ? 'US'
                : cur == 'GBP'
                    ? 'GB'
                    : cur.substring(0, 2);
        if (widgets.isNotEmpty) {
          widgets.add(const Divider(
              indent: AppSpacing.xxl + AppSpacing.md, height: 1));
        }
        widgets.add(_AccountRow(
            flag: flag, currency: cur, balance: '$cur $avail'));
      }
    }

    return AppCard(
      padding: const EdgeInsets.symmetric(vertical: AppSpacing.sm),
      child: Column(children: widgets),
    );
  }
}

class _TransactionsCard extends StatelessWidget {
  const _TransactionsCard({required this.transactions});

  final List<dynamic> transactions;

  @override
  Widget build(BuildContext context) {
    if (transactions.isEmpty) {
      return AppCard(
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.lg),
          child: Center(
            child: Text('No transactions yet',
                style: AppTypography.body2
                    .copyWith(color: AppColors.neutral500)),
          ),
        ),
      );
    }

    return AppCard(
      padding: EdgeInsets.zero,
      child: Column(
        children: [
          for (int i = 0; i < transactions.length; i++) ...[
            if (i > 0)
              const Divider(
                  indent: AppSpacing.xxl + AppSpacing.lg, height: 1),
            Builder(builder: (context) {
              final tx = transactions[i] as Map<String, dynamic>;
              final type = tx['type']?.toString() ?? '';
              final amount = tx['amount']?.toString() ?? '0.00';
              final currency = tx['currency']?.toString() ?? 'EUR';
              final isPositive =
                  type == 'credit' || type == 'receive';
              final title = tx['recipient_name']?.toString() ??
                  tx['reference']?.toString() ??
                  type;
              final status = tx['status']?.toString() ?? '';
              return TransactionListTile(
                title: title,
                subtitle: status,
                amount:
                    '${isPositive ? '+' : '-'}$currency $amount',
                isPositive: isPositive,
                onTap: () {},
              );
            }),
          ],
        ],
      ),
    );
  }
}

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
            child: Text('$actionLabel >',
                style: AppTypography.body2
                    .copyWith(color: AppColors.primary500)),
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
  });

  final String flag;
  final String currency;
  final String balance;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Container(
        width: 32,
        height: 32,
        decoration: const BoxDecoration(
          color: AppColors.neutral100,
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Text(flag,
              style: AppTypography.caption
                  .copyWith(fontWeight: FontWeight.w600)),
        ),
      ),
      title: Text(currency, style: AppTypography.body1),
      trailing: Text(balance,
          style: AppTypography.body1.copyWith(
            fontWeight: FontWeight.w600,
            fontFeatures: const [FontFeature.tabularFigures()],
          )),
      dense: true,
    );
  }
}
