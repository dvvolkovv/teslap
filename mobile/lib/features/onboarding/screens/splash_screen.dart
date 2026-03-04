import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../auth/bloc/auth_bloc.dart';

/// Full-screen splash with the TeslaPay logo.
///
/// Auto-advances after 2 seconds based on current auth state:
/// - Authenticated and session valid -> Home
/// - Has device ID but unauthenticated -> Login
/// - Fresh install -> Welcome carousel
class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen>
    with SingleTickerProviderStateMixin {
  late final AnimationController _pulseController;
  late final Animation<double> _pulseAnimation;

  @override
  void initState() {
    super.initState();

    _pulseController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1200),
    )..repeat(reverse: true);

    _pulseAnimation = Tween<double>(begin: 0.95, end: 1.05).animate(
      CurvedAnimation(parent: _pulseController, curve: Curves.easeInOut),
    );

    // Dispatch auth check and wait for result.
    context.read<AuthBloc>().add(const AuthCheckRequested());
  }

  @override
  void dispose() {
    _pulseController.dispose();
    super.dispose();
  }

  void _navigate(AuthState state) {
    if (!mounted) return;

    // Delay for branding visibility.
    Future.delayed(const Duration(seconds: 2), () {
      if (!mounted) return;
      if (state is AuthAuthenticated) {
        context.go(AppRoutes.home);
      } else if (state is AuthNeedsOnboarding) {
        context.go(AppRoutes.welcome);
      } else {
        context.go(AppRoutes.login);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state is! AuthInitial && state is! AuthLoading) {
          _navigate(state);
        }
      },
      child: Scaffold(
        body: Container(
          width: double.infinity,
          height: double.infinity,
          decoration: const BoxDecoration(
            gradient: AppColors.brandGradient,
          ),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              ScaleTransition(
                scale: _pulseAnimation,
                child: Container(
                  width: 80,
                  height: 80,
                  decoration: BoxDecoration(
                    color: AppColors.white.withValues(alpha: 0.2),
                    shape: BoxShape.circle,
                  ),
                  child: const Center(
                    child: Text(
                      'T',
                      style: TextStyle(
                        fontFamily: 'Inter',
                        fontSize: 40,
                        fontWeight: FontWeight.w700,
                        color: AppColors.white,
                      ),
                    ),
                  ),
                ),
              ),
              const SizedBox(height: AppSpacing.md),
              Text(
                'TeslaPay',
                style: AppTypography.h1.copyWith(color: AppColors.white),
              ),
              const SizedBox(height: AppSpacing.xxl),
              const SizedBox(
                width: 24,
                height: 24,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  valueColor:
                      AlwaysStoppedAnimation<Color>(AppColors.white),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
