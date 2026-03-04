import 'package:flutter/material.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_card.dart';
import '../../home/widgets/transaction_list_tile.dart';
import '../widgets/card_widget.dart';

/// Card tab showing the user's virtual/physical Mastercard, quick actions,
/// card settings summary, and recent card transactions.
class CardScreen extends StatefulWidget {
  const CardScreen({super.key});

  @override
  State<CardScreen> createState() => _CardScreenState();
}

class _CardScreenState extends State<CardScreen> {
  bool _isFrozen = false;
  bool _showDetails = false;

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
                'Card',
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
                  // Card visual
                  CardWidget(
                    lastFour: '4521',
                    cardholderName: 'John Doe',
                    expiry: '03/29',
                    isVirtual: true,
                    isFrozen: _isFrozen,
                    showFullNumber: _showDetails,
                    fullNumber: _showDetails ? '5412345678904521' : null,
                    cvv: _showDetails ? '123' : null,
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Quick actions
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                    children: [
                      _CardAction(
                        icon: _showDetails
                            ? PhosphorIconsRegular.eyeSlash
                            : PhosphorIconsRegular.eye,
                        label:
                            _showDetails ? 'Hide Details' : 'Show Details',
                        onTap: () {
                          // TODO: require biometric before showing
                          setState(() => _showDetails = !_showDetails);
                        },
                      ),
                      _CardAction(
                        icon: _isFrozen
                            ? PhosphorIconsRegular.snowflake
                            : PhosphorIconsRegular.snowflake,
                        label: _isFrozen ? 'Unfreeze' : 'Freeze',
                        onTap: () => setState(() => _isFrozen = !_isFrozen),
                        isDanger: _isFrozen,
                      ),
                      _CardAction(
                        icon: PhosphorIconsRegular.wallet,
                        label: 'Apple Pay',
                        onTap: () {
                          // TODO: tokenize card
                        },
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Card Settings
                  AppCard(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text('Card Settings', style: AppTypography.h3),
                        const SizedBox(height: AppSpacing.md),
                        _SettingsRow(
                          icon: PhosphorIconsRegular.lockKey,
                          label: 'View PIN',
                          onTap: () {},
                        ),
                        _SettingsRow(
                          icon: PhosphorIconsRegular.gauge,
                          label: 'Spending Limits',
                          onTap: () {},
                        ),
                        _SettingsRow(
                          icon: PhosphorIconsRegular.shieldCheck,
                          label: 'Security Controls',
                          onTap: () {},
                        ),
                        _SettingsRow(
                          icon: PhosphorIconsRegular.link,
                          label: 'Linked Account',
                          trailing: 'EUR',
                          onTap: () {},
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Recent Card Transactions
                  Text('Recent Card Transactions', style: AppTypography.h3),
                  const SizedBox(height: AppSpacing.sm),
                  AppCard(
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        TransactionListTile(
                          title: 'Amazon.de',
                          subtitle: 'Online Shopping  Today 14:30',
                          amount: '-EUR 45.99',
                          status: 'Settled',
                          onTap: () {},
                        ),
                        const Divider(
                          indent: AppSpacing.xxl + AppSpacing.lg,
                          height: 1,
                        ),
                        TransactionListTile(
                          title: 'Bolt',
                          subtitle: 'Transport  Yesterday',
                          amount: '-EUR 8.50',
                          status: 'Completed',
                          onTap: () {},
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Card actions
                  AppButton(
                    label: 'Order Physical Card',
                    variant: AppButtonVariant.secondary,
                    onPressed: () {
                      // TODO: navigate to order physical card flow
                    },
                  ),
                  const SizedBox(height: AppSpacing.sm),
                  AppButton(
                    label: 'Report Lost or Stolen',
                    variant: AppButtonVariant.danger,
                    onPressed: () {
                      // TODO: navigate to report flow
                    },
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
    return Semantics(
      label: label,
      button: true,
      child: GestureDetector(
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
              Text(
                trailing!,
                style: AppTypography.body2.copyWith(
                  color: AppColors.neutral500,
                ),
              ),
            const SizedBox(width: AppSpacing.xs),
            const Icon(
              Icons.chevron_right,
              size: 20,
              color: AppColors.neutral400,
            ),
          ],
        ),
      ),
    );
  }
}
