import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/theme/app_theme.dart';
import '../../../models/crypto.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../bloc/crypto_bloc.dart';

/// Crypto tab — loads wallet, balances, and prices from backend via [CryptoBloc].
class CryptoScreen extends StatefulWidget {
  const CryptoScreen({super.key});

  @override
  State<CryptoScreen> createState() => _CryptoScreenState();
}

class _CryptoScreenState extends State<CryptoScreen> {
  @override
  void initState() {
    super.initState();
    context.read<CryptoBloc>().add(const CryptoLoadRequested());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: BlocConsumer<CryptoBloc, CryptoState>(
          listener: (context, state) {
            if (state is CryptoActionSuccess) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(state.message),
                  backgroundColor: AppColors.success500,
                ),
              );
            } else if (state is CryptoError) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(state.message),
                  backgroundColor: AppColors.error500,
                ),
              );
            }
          },
          builder: (context, state) {
            return CustomScrollView(
              slivers: [
                SliverAppBar(
                  floating: true,
                  snap: true,
                  backgroundColor:
                      Theme.of(context).scaffoldBackgroundColor,
                  title: Text('Crypto',
                      style: AppTypography.h2
                          .copyWith(fontWeight: FontWeight.w700)),
                  centerTitle: false,
                ),
                if (state is CryptoLoading)
                  const SliverFillRemaining(
                    child: Center(child: CircularProgressIndicator()),
                  )
                else if (state is CryptoLoaded)
                  _buildLoadedContent(context, state)
                else if (state is CryptoError)
                  SliverFillRemaining(
                    child: Center(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text('Error loading crypto',
                              style: AppTypography.h3),
                          const SizedBox(height: AppSpacing.sm),
                          TextButton(
                            onPressed: () => context
                                .read<CryptoBloc>()
                                .add(const CryptoLoadRequested()),
                            child: const Text('Retry'),
                          ),
                        ],
                      ),
                    ),
                  )
                else
                  const SliverFillRemaining(
                    child: Center(child: CircularProgressIndicator()),
                  ),
              ],
            );
          },
        ),
      ),
    );
  }

  Widget _buildLoadedContent(
      BuildContext context, CryptoLoaded state) {
    final wallet = state.wallet;
    final prices = state.prices;

    return SliverPadding(
      padding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.screenMargin),
      sliver: SliverList(
        delegate: SliverChildListDelegate([
          // Total crypto balance hero
          AppCard(
            gradient: AppColors.cryptoGradient,
            borderRadius: BorderRadius.circular(AppRadius.lg),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Total Crypto Balance',
                    style: AppTypography.body2.copyWith(
                        color:
                            AppColors.white.withValues(alpha: 0.8))),
                const SizedBox(height: AppSpacing.sm),
                Text('EUR ${state.totalValueEur}',
                    style: AppTypography.display2
                        .copyWith(color: AppColors.white)),
              ],
            ),
          ),
          const SizedBox(height: AppSpacing.lg),

          // Token balances
          Text('Tokens', style: AppTypography.h3),
          const SizedBox(height: AppSpacing.sm),
          for (final balance in wallet.balances) ...[
            _TokenCard(
              symbol: balance.tokenSymbol,
              name: balance.tokenName,
              balance: balance.balance,
              eurValue: 'EUR ${balance.valueEur ?? '0.00'}',
              change: _getPriceChange(prices, balance.tokenSymbol),
              isNegative: _isPriceNegative(
                  prices, balance.tokenSymbol),
              onTap: () {},
            ),
            const SizedBox(height: AppSpacing.sm),
          ],
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

          // Wallet info
          Text('Wallet', style: AppTypography.h3),
          const SizedBox(height: AppSpacing.sm),
          AppCard(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Wallet Address',
                    style: AppTypography.caption
                        .copyWith(color: AppColors.neutral500)),
                const SizedBox(height: AppSpacing.xs),
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        wallet.address.length > 20
                            ? '${wallet.address.substring(0, 10)}...${wallet.address.substring(wallet.address.length - 6)}'
                            : wallet.address,
                        style: AppTypography.mono,
                      ),
                    ),
                    IconButton(
                      icon: const Icon(PhosphorIconsRegular.copy,
                          size: 20),
                      onPressed: () {
                        Clipboard.setData(
                            ClipboardData(text: wallet.address));
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(
                              content: Text('Address copied')),
                        );
                      },
                      tooltip: 'Copy address',
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
                      '${wallet.network.toUpperCase()} Network - Connected',
                      style: AppTypography.caption
                          .copyWith(color: AppColors.success500),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: AppSpacing.xxl),
        ]),
      ),
    );
  }

  String _getPriceChange(List<CryptoPrice> prices, String symbol) {
    for (final p in prices) {
      if (p.symbol == symbol) return p.change24h;
    }
    return '0';
  }

  bool _isPriceNegative(List<CryptoPrice> prices, String symbol) {
    final change = _getPriceChange(prices, symbol);
    return change.startsWith('-');
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
          Container(
            width: 40,
            height: 40,
            decoration: const BoxDecoration(
              color: AppColors.accent100,
              shape: BoxShape.circle,
            ),
            child: Center(
              child: Text(
                symbol.isNotEmpty ? symbol[0] : '?',
                style: AppTypography.body1.copyWith(
                    fontWeight: FontWeight.w700,
                    color: AppColors.accent500),
              ),
            ),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(name,
                    style: AppTypography.body1
                        .copyWith(fontWeight: FontWeight.w500)),
                Text('$balance $symbol', style: AppTypography.h3),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(eurValue,
                  style: AppTypography.body2
                      .copyWith(color: AppColors.neutral500)),
              Container(
                padding: const EdgeInsets.symmetric(
                    horizontal: AppSpacing.xs,
                    vertical: AppSpacing.xxs),
                decoration: BoxDecoration(
                  color: (isNegative
                          ? AppColors.error500
                          : AppColors.success500)
                      .withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(AppRadius.xs),
                ),
                child: Text('$change%',
                    style: AppTypography.caption.copyWith(
                        color: isNegative
                            ? AppColors.error500
                            : AppColors.success500,
                        fontWeight: FontWeight.w500)),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
