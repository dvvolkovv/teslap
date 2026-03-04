import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_card.dart';
import '../../auth/bloc/auth_bloc.dart';

/// Profile tab with user info, account settings, security, preferences,
/// support, legal, and logout.
class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state is AuthUnauthenticated) {
          context.go(AppRoutes.login);
        }
      },
      child: Scaffold(
        body: SafeArea(
          child: CustomScrollView(
            slivers: [
              SliverAppBar(
                floating: true,
                snap: true,
                backgroundColor: Theme.of(context).scaffoldBackgroundColor,
                title: Text(
                  'Profile',
                  style:
                      AppTypography.h2.copyWith(fontWeight: FontWeight.w700),
                ),
                centerTitle: false,
              ),
              SliverPadding(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSpacing.screenMargin,
                ),
                sliver: SliverList(
                  delegate: SliverChildListDelegate([
                    // User header
                    AppCard(
                      child: Row(
                        children: [
                          Container(
                            width: 56,
                            height: 56,
                            decoration: const BoxDecoration(
                              gradient: AppColors.brandGradient,
                              shape: BoxShape.circle,
                            ),
                            child: const Center(
                              child: Text(
                                'JD',
                                style: TextStyle(
                                  color: AppColors.white,
                                  fontWeight: FontWeight.w700,
                                  fontSize: 20,
                                ),
                              ),
                            ),
                          ),
                          const SizedBox(width: AppSpacing.md),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  'John Doe',
                                  style: AppTypography.h3,
                                ),
                                Text(
                                  'john.doe@example.com',
                                  style: AppTypography.body2.copyWith(
                                    color: AppColors.neutral500,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: AppSpacing.sm,
                              vertical: AppSpacing.xs,
                            ),
                            decoration: BoxDecoration(
                              color: AppColors.primary100,
                              borderRadius:
                                  BorderRadius.circular(AppRadius.xs),
                            ),
                            child: Text(
                              'Standard',
                              style: AppTypography.caption.copyWith(
                                color: AppColors.primary500,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // Account section
                    _SectionTitle(title: 'Account'),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          _ProfileRow(
                            icon: PhosphorIconsRegular.user,
                            label: 'Personal Information',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.shieldCheck,
                            label: 'Verification Status',
                            trailing: 'Verified',
                            trailingColor: AppColors.success500,
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.medal,
                            label: 'Account Tier',
                            trailing: 'Standard',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.receipt,
                            label: 'Fees and Limits',
                            onTap: () {},
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // Security section
                    _SectionTitle(title: 'Security'),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          _ProfileRow(
                            icon: PhosphorIconsRegular.fingerprint,
                            label: 'Biometric Login',
                            trailing: 'Enabled',
                            trailingColor: AppColors.success500,
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.keyhole,
                            label: 'Change PIN',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.lockKey,
                            label: 'Change Password',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.devices,
                            label: 'Active Sessions',
                            onTap: () {},
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // Preferences section
                    _SectionTitle(title: 'Preferences'),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          _ProfileRow(
                            icon: PhosphorIconsRegular.bell,
                            label: 'Notifications',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.translate,
                            label: 'Language',
                            trailing: 'English',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.moonStars,
                            label: 'Appearance',
                            trailing: 'System',
                            onTap: () {},
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // Support section
                    _SectionTitle(title: 'Support'),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          _ProfileRow(
                            icon: PhosphorIconsRegular.question,
                            label: 'Help Center',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.chatCircle,
                            label: 'Chat with Us',
                            onTap: () {},
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // Legal section
                    _SectionTitle(title: 'Legal'),
                    const SizedBox(height: AppSpacing.sm),
                    AppCard(
                      padding: EdgeInsets.zero,
                      child: Column(
                        children: [
                          _ProfileRow(
                            icon: PhosphorIconsRegular.file,
                            label: 'Terms of Service',
                            onTap: () {},
                          ),
                          const Divider(height: 1),
                          _ProfileRow(
                            icon: PhosphorIconsRegular.shieldStar,
                            label: 'Privacy Policy',
                            onTap: () {},
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),

                    // Logout
                    SizedBox(
                      width: double.infinity,
                      child: TextButton(
                        onPressed: () {
                          context
                              .read<AuthBloc>()
                              .add(const AuthLogoutRequested());
                        },
                        child: Text(
                          'Log Out',
                          style: AppTypography.button.copyWith(
                            color: AppColors.error500,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(height: AppSpacing.sm),
                    Center(
                      child: Text(
                        'TeslaPay v1.0.0',
                        style: AppTypography.caption.copyWith(
                          color: AppColors.neutral400,
                        ),
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

class _SectionTitle extends StatelessWidget {
  const _SectionTitle({required this.title});

  final String title;

  @override
  Widget build(BuildContext context) {
    return Text(
      title.toUpperCase(),
      style: AppTypography.overline.copyWith(
        color: AppColors.neutral500,
      ),
    );
  }
}

class _ProfileRow extends StatelessWidget {
  const _ProfileRow({
    required this.icon,
    required this.label,
    required this.onTap,
    this.trailing,
    this.trailingColor,
  });

  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final String? trailing;
  final Color? trailingColor;

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Padding(
        padding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.md,
          vertical: AppSpacing.md,
        ),
        child: Row(
          children: [
            Icon(icon, size: 24, color: AppColors.neutral500),
            const SizedBox(width: AppSpacing.md),
            Expanded(child: Text(label, style: AppTypography.body1)),
            if (trailing != null)
              Text(
                trailing!,
                style: AppTypography.body2.copyWith(
                  color: trailingColor ?? AppColors.neutral500,
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
