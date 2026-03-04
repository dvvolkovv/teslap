import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../../home/widgets/transaction_list_tile.dart';
import '../bloc/card_bloc.dart';
import '../widgets/card_widget.dart';

/// Card tab — loads card data from the backend via [CardBloc].
class CardScreen extends StatefulWidget {
  const CardScreen({super.key});

  @override
  State<CardScreen> createState() => _CardScreenState();
}

class _CardScreenState extends State<CardScreen> {
  bool _showDetails = false;

  @override
  void initState() {
    super.initState();
    context.read<CardBloc>().add(const CardListRequested());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: BlocConsumer<CardBloc, CardState>(
          listener: (context, state) {
            if (state is CardActionSuccess) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(state.message),
                  backgroundColor: AppColors.success500,
                ),
              );
            } else if (state is CardError) {
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
                  title: Text('Card',
                      style: AppTypography.h2
                          .copyWith(fontWeight: FontWeight.w700)),
                  centerTitle: false,
                ),
                if (state is CardLoading)
                  const SliverFillRemaining(
                    child: Center(child: CircularProgressIndicator()),
                  )
                else if (state is CardListLoaded)
                  _buildCardContent(context, state.cards)
                else if (state is CardError)
                  SliverFillRemaining(
                    child: Center(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text('Error loading cards',
                              style: AppTypography.h3),
                          const SizedBox(height: AppSpacing.sm),
                          TextButton(
                            onPressed: () => context
                                .read<CardBloc>()
                                .add(const CardListRequested()),
                            child: const Text('Retry'),
                          ),
                        ],
                      ),
                    ),
                  )
                else
                  _buildNoCards(context),
              ],
            );
          },
        ),
      ),
    );
  }

  Widget _buildNoCards(BuildContext context) {
    return SliverFillRemaining(
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Icon(PhosphorIconsRegular.creditCard,
                size: 64, color: AppColors.neutral300),
            const SizedBox(height: AppSpacing.lg),
            Text('No cards yet', style: AppTypography.h3),
            const SizedBox(height: AppSpacing.md),
            AppButton(
              label: 'Issue Virtual Card',
              onPressed: () {
                context.read<CardBloc>().add(
                      const CardIssueVirtualRequested(
                        accountId: '',
                        cardholderName: 'John Doe',
                      ),
                    );
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCardContent(
      BuildContext context, List<dynamic> cards) {
    if (cards.isEmpty) return _buildNoCards(context);

    final card = cards[0] as Map<String, dynamic>;
    final lastFour = card['last_four']?.toString() ?? '0000';
    final cardholderName =
        card['cardholder_name']?.toString() ?? 'Cardholder';
    final expiryMonth = card['expiry_month']?.toString() ?? '01';
    final expiryYear = card['expiry_year']?.toString() ?? '29';
    final cardType = card['type']?.toString() ?? 'virtual';
    final cardStatus = card['status']?.toString() ?? 'active';
    final cardId = card['id']?.toString() ?? '';
    final isFrozen = cardStatus == 'frozen';

    return SliverPadding(
      padding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.screenMargin),
      sliver: SliverList(
        delegate: SliverChildListDelegate([
          CardWidget(
            lastFour: lastFour,
            cardholderName: cardholderName,
            expiry: '$expiryMonth/$expiryYear',
            isVirtual: cardType == 'virtual',
            isFrozen: isFrozen,
            showFullNumber: _showDetails,
          ),
          const SizedBox(height: AppSpacing.lg),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _CardAction(
                icon: _showDetails
                    ? PhosphorIconsRegular.eyeSlash
                    : PhosphorIconsRegular.eye,
                label: _showDetails ? 'Hide Details' : 'Show Details',
                onTap: () =>
                    setState(() => _showDetails = !_showDetails),
              ),
              _CardAction(
                icon: PhosphorIconsRegular.snowflake,
                label: isFrozen ? 'Unfreeze' : 'Freeze',
                onTap: () {
                  if (isFrozen) {
                    context.read<CardBloc>().add(
                        CardUnfreezeRequested(cardId: cardId));
                  } else {
                    context.read<CardBloc>().add(
                        CardFreezeRequested(cardId: cardId));
                  }
                },
                isDanger: isFrozen,
              ),
              _CardAction(
                icon: PhosphorIconsRegular.wallet,
                label: 'Apple Pay',
                onTap: () {},
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.lg),
          AppCard(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Card Settings', style: AppTypography.h3),
                const SizedBox(height: AppSpacing.md),
                _SettingsRow(
                    icon: PhosphorIconsRegular.lockKey,
                    label: 'View PIN',
                    onTap: () {}),
                _SettingsRow(
                    icon: PhosphorIconsRegular.gauge,
                    label: 'Spending Limits',
                    onTap: () {}),
                _SettingsRow(
                    icon: PhosphorIconsRegular.shieldCheck,
                    label: 'Security Controls',
                    onTap: () {}),
                _SettingsRow(
                    icon: PhosphorIconsRegular.link,
                    label: 'Linked Account',
                    trailing: 'EUR',
                    onTap: () {}),
              ],
            ),
          ),
          const SizedBox(height: AppSpacing.lg),
          Text('Recent Card Transactions', style: AppTypography.h3),
          const SizedBox(height: AppSpacing.sm),
          AppCard(
            padding: EdgeInsets.zero,
            child: Column(
              children: [
                TransactionListTile(
                  title: 'No transactions yet',
                  subtitle: 'Use your card to see transactions here',
                  amount: '',
                  onTap: () {},
                ),
              ],
            ),
          ),
          const SizedBox(height: AppSpacing.lg),
          AppButton(
            label: 'Order Physical Card',
            variant: AppButtonVariant.secondary,
            onPressed: () {},
          ),
          const SizedBox(height: AppSpacing.sm),
          AppButton(
            label: 'Report Lost or Stolen',
            variant: AppButtonVariant.danger,
            onPressed: () {
              context
                  .read<CardBloc>()
                  .add(CardBlockRequested(cardId: cardId));
            },
          ),
          const SizedBox(height: AppSpacing.xxl),
        ]),
      ),
    );
  }
}

class _CardAction extends StatelessWidget {
  const _CardAction({
    required this.icon,
    required this.label,
    required this.onTap,
    this.isDanger = false,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool isDanger;

  @override
  Widget build(BuildContext context) {
    final color = isDanger ? AppColors.warning500 : AppColors.primary500;
    return GestureDetector(
      onTap: onTap,
      child: Column(
        children: [
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: color.withValues(alpha: 0.1),
              shape: BoxShape.circle,
            ),
            child: Icon(icon, size: 24, color: color),
          ),
          const SizedBox(height: AppSpacing.xs),
          Text(label,
              style: AppTypography.caption
                  .copyWith(color: AppColors.neutral700)),
        ],
      ),
    );
  }
}

class _SettingsRow extends StatelessWidget {
  const _SettingsRow({
    required this.icon,
    required this.label,
    required this.onTap,
    this.trailing,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final String? trailing;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
        child: Row(
          children: [
            Icon(icon, size: 24, color: AppColors.neutral500),
            const SizedBox(width: AppSpacing.md),
            Expanded(child: Text(label, style: AppTypography.body1)),
            if (trailing != null)
              Text(trailing!,
                  style: AppTypography.body2
                      .copyWith(color: AppColors.neutral500)),
            const SizedBox(width: AppSpacing.xs),
            const Icon(Icons.chevron_right,
                size: 20, color: AppColors.neutral400),
          ],
        ),
      ),
    );
  }
}
