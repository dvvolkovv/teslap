import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_card.dart';

/// Root screen for the Payments tab.
///
/// Sections: Quick actions, Saved payees, Scheduled payments,
/// Direct debits.
class PaymentsScreen extends StatelessWidget {
  const PaymentsScreen({super.key});

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
                'Payments',
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
                  // Quick actions grid
                  Row(
                    children: [
                      _ActionTile(
                        icon: PhosphorIconsRegular.arrowUp,
                        label: 'Send',
                        onTap: () => context.push(AppRoutes.sendMoney),
                      ),
                      const SizedBox(width: AppSpacing.md),
                      _ActionTile(
                        icon: PhosphorIconsRegular.arrowDown,
                        label: 'Request',
                        onTap: () {},
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSpacing.md),
                  Row(
                    children: [
                      _ActionTile(
                        icon: PhosphorIconsRegular.arrowsLeftRight,
                        label: 'Exchange',
                        onTap: () {},
                      ),
                      const SizedBox(width: AppSpacing.md),
                      _ActionTile(
                        icon: PhosphorIconsRegular.qrCode,
                        label: 'Scan QR',
                        onTap: () {},
                      ),
                    ],
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Saved Payees
                  _SectionHeader(title: 'Saved Payees', action: 'Manage'),
                  const SizedBox(height: AppSpacing.sm),
                  AppCard(
                    padding: EdgeInsets.zero,
                    child: Column(
                      children: [
                        _PayeeRow(
                          name: 'Anna Kowalski',
                          subtitle: 'EE12 7890 ****',
                          onTap: () => context.push(AppRoutes.sendMoney),
                        ),
                        const Divider(
                          indent: AppSpacing.xxl + AppSpacing.lg,
                          height: 1,
                        ),
                        _PayeeRow(
                          name: 'Landlord LLC',
                          subtitle: 'DE89 3704 ****',
                          onTap: () => context.push(AppRoutes.sendMoney),
                        ),
                        const Divider(
                          indent: AppSpacing.xxl + AppSpacing.lg,
                          height: 1,
                        ),
                        ListTile(
                          leading: Container(
                            width: 40,
                            height: 40,
                            decoration: const BoxDecoration(
                              color: AppColors.primary50,
                              shape: BoxShape.circle,
                            ),
                            child: const Icon(
                              PhosphorIconsRegular.plus,
                              color: AppColors.primary500,
                              size: 20,
                            ),
                          ),
                          title: Text(
                            'Add Payee',
                            style: AppTypography.body1.copyWith(
                              color: AppColors.primary500,
                            ),
                          ),
                          onTap: () {},
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // Scheduled Payments
                  _SectionHeader(
                    title: 'Scheduled Payments',
                    action: 'See All',
                  ),
                  const SizedBox(height: AppSpacing.sm),
                  AppCard(
                    child: Column(
                      children: [
                        _ScheduledRow(
                          name: 'Rent - Landlord LLC',
                          amount: 'EUR 850.00',
                          schedule: 'Monthly, next: 1 Apr',
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

class _ActionTile extends StatelessWidget {
  const _ActionTile({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: AppCard(
        onTap: onTap,
        child: Column(
          children: [
            Icon(icon, size: 28, color: AppColors.primary500),
            const SizedBox(height: AppSpacing.sm),
            Text(label, style: AppTypography.body1),
          ],
        ),
      ),
    );
  }
}

class _SectionHeader extends StatelessWidget {
  const _SectionHeader({required this.title, this.action});

  final String title;
  final String? action;

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(title, style: AppTypography.h3),
        if (action != null)
          GestureDetector(
            onTap: () {},
            child: Text(
              '$action >',
              style: AppTypography.body2.copyWith(
                color: AppColors.primary500,
              ),
            ),
          ),
      ],
    );
  }
}

class _PayeeRow extends StatelessWidget {
  const _PayeeRow({
    required this.name,
    required this.subtitle,
    this.onTap,
  });

  final String name;
  final String subtitle;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return ListTile(
      onTap: onTap,
      leading: Container(
        width: 40,
        height: 40,
        decoration: const BoxDecoration(
          color: AppColors.neutral100,
          shape: BoxShape.circle,
        ),
        child: Center(
          child: Text(
            name.isNotEmpty ? name[0] : '?',
            style: AppTypography.body1.copyWith(fontWeight: FontWeight.w600),
          ),
        ),
      ),
      title: Text(name, style: AppTypography.body1),
      subtitle: Text(
        subtitle,
        style: AppTypography.caption.copyWith(color: AppColors.neutral500),
      ),
      trailing: const Icon(
        Icons.chevron_right,
        size: 20,
        color: AppColors.neutral400,
      ),
    );
  }
}

class _ScheduledRow extends StatelessWidget {
  const _ScheduledRow({
    required this.name,
    required this.amount,
    required this.schedule,
  });

  final String name;
  final String amount;
  final String schedule;

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Container(
          width: 40,
          height: 40,
          decoration: const BoxDecoration(
            color: AppColors.primary100,
            shape: BoxShape.circle,
          ),
          child: const Icon(
            PhosphorIconsRegular.calendarBlank,
            size: 20,
            color: AppColors.primary500,
          ),
        ),
        const SizedBox(width: AppSpacing.md),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(name, style: AppTypography.body1),
              Text(
                schedule,
                style: AppTypography.caption.copyWith(
                  color: AppColors.neutral500,
                ),
              ),
            ],
          ),
        ),
        Text(
          amount,
          style: AppTypography.body1.copyWith(
            fontWeight: FontWeight.w600,
            fontFeatures: const [FontFeature.tabularFigures()],
          ),
        ),
      ],
    );
  }
}
