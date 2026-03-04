import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_button.dart';

/// Three-slide onboarding carousel shown on first launch.
///
/// Slides:
/// 1. Multi-currency banking
/// 2. Mastercard
/// 3. Crypto on Fuse
class WelcomeScreen extends StatefulWidget {
  const WelcomeScreen({super.key});

  @override
  State<WelcomeScreen> createState() => _WelcomeScreenState();
}

class _WelcomeScreenState extends State<WelcomeScreen> {
  final PageController _pageController = PageController();
  int _currentPage = 0;

  static const _slides = [
    _SlideData(
      icon: PhosphorIconsRegular.globe,
      title: 'Banking without borders',
      description:
          'Multi-currency accounts, instant SEPA transfers, and a Mastercard that works everywhere in Europe.',
    ),
    _SlideData(
      icon: PhosphorIconsRegular.creditCard,
      title: 'Your card, your rules',
      description:
          'A virtual Mastercard in seconds. Freeze, set limits, and add to Apple Pay or Google Pay instantly.',
    ),
    _SlideData(
      icon: PhosphorIconsRegular.currencyBtc,
      title: 'Crypto made simple',
      description:
          'Buy, sell, and send crypto from your bank account. Earn yield on stablecoins with Fuse network.',
    ),
  ];

  @override
  void dispose() {
    _pageController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            // Skip button
            Align(
              alignment: Alignment.centerRight,
              child: TextButton(
                onPressed: () => context.go(AppRoutes.register),
                child: Text(
                  'Skip',
                  style: AppTypography.body2.copyWith(
                    color: AppColors.neutral500,
                  ),
                ),
              ),
            ),

            // Page view
            Expanded(
              child: PageView.builder(
                controller: _pageController,
                itemCount: _slides.length,
                onPageChanged: (index) =>
                    setState(() => _currentPage = index),
                itemBuilder: (context, index) {
                  final slide = _slides[index];
                  return Padding(
                    padding: const EdgeInsets.symmetric(
                      horizontal: AppSpacing.lg,
                    ),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Container(
                          width: 120,
                          height: 120,
                          decoration: BoxDecoration(
                            color: AppColors.primary50,
                            borderRadius:
                                BorderRadius.circular(AppRadius.xl),
                          ),
                          child: Icon(
                            slide.icon,
                            size: 48,
                            color: AppColors.primary500,
                          ),
                        ),
                        const SizedBox(height: AppSpacing.xl),
                        Text(
                          slide.title,
                          style: AppTypography.h1,
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: AppSpacing.md),
                        Text(
                          slide.description,
                          style: AppTypography.body1.copyWith(
                            color: AppColors.neutral500,
                          ),
                          textAlign: TextAlign.center,
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),

            // Dot indicators
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: List.generate(
                _slides.length,
                (index) => AnimatedContainer(
                  duration: const Duration(milliseconds: 250),
                  margin:
                      const EdgeInsets.symmetric(horizontal: AppSpacing.xs),
                  width: index == _currentPage ? 24 : 8,
                  height: 8,
                  decoration: BoxDecoration(
                    color: index == _currentPage
                        ? AppColors.primary500
                        : AppColors.neutral200,
                    borderRadius: BorderRadius.circular(AppRadius.full),
                  ),
                ),
              ),
            ),
            const SizedBox(height: AppSpacing.xl),

            // CTA buttons
            Padding(
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.screenMargin,
              ),
              child: AppButton(
                label: 'Get Started',
                onPressed: () => context.go(AppRoutes.register),
              ),
            ),
            const SizedBox(height: AppSpacing.md),
            TextButton(
              onPressed: () => context.go(AppRoutes.login),
              child: Text(
                'I already have an account',
                style: AppTypography.body2.copyWith(
                  color: AppColors.primary500,
                ),
              ),
            ),
            const SizedBox(height: AppSpacing.lg),
          ],
        ),
      ),
    );
  }
}

class _SlideData {
  const _SlideData({
    required this.icon,
    required this.title,
    required this.description,
  });

  final IconData icon;
  final String title;
  final String description;
}
