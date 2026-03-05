import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/auth/biometric_auth.dart';
import '../../../core/di/injection.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../bloc/auth_bloc.dart';

/// Biometric login screen shown to returning users.
///
/// Automatically triggers the system biometric prompt on load.
/// Falls back to PIN or password entry on failure.
class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  bool _biometricAvailable = false;
  String _biometricLabel = 'Biometric';

  @override
  void initState() {
    super.initState();
    _checkBiometricAvailability();
  }

  Future<void> _checkBiometricAvailability() async {
    final bioAuth = getIt<BiometricAuth>();
    final available = await bioAuth.isAvailable;
    final label = await bioAuth.biometricLabel;

    if (mounted) {
      setState(() {
        _biometricAvailable = available;
        _biometricLabel = label;
      });

      if (available) {
        // Auto-trigger biometric prompt.
        _requestBiometric();
      }
    }
  }

  void _requestBiometric() {
    context.read<AuthBloc>().add(const AuthBiometricRequested());
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state is AuthAuthenticated) {
          context.go(AppRoutes.home);
        }
      },
      child: Scaffold(
        body: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSpacing.screenMargin,
            ),
            child: Column(
              children: [
                const Spacer(flex: 2),

                // Logo
                Container(
                  width: 72,
                  height: 72,
                  decoration: const BoxDecoration(
                    gradient: AppColors.brandGradient,
                    shape: BoxShape.circle,
                  ),
                  child: const Center(
                    child: Text(
                      'T',
                      style: TextStyle(
                        fontFamily: 'Inter',
                        fontSize: 32,
                        fontWeight: FontWeight.w700,
                        color: AppColors.white,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: AppSpacing.lg),

                Text(
                  'Welcome back',
                  style: AppTypography.h1,
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: AppSpacing.xxl),

                // Biometric button
                if (_biometricAvailable) ...[
                  GestureDetector(
                    onTap: _requestBiometric,
                    child: Semantics(
                      label: 'Tap to unlock with $_biometricLabel',
                      button: true,
                      child: Column(
                        children: [
                          Container(
                            width: 72,
                            height: 72,
                            decoration: BoxDecoration(
                              color: AppColors.primary50,
                              shape: BoxShape.circle,
                              border: Border.all(
                                color: AppColors.primary100,
                                width: 2,
                              ),
                            ),
                            child: const Icon(
                              PhosphorIconsRegular.fingerprint,
                              size: 36,
                              color: AppColors.primary500,
                            ),
                          ),
                          const SizedBox(height: AppSpacing.md),
                          Text(
                            'Tap to unlock',
                            style: AppTypography.body2.copyWith(
                              color: AppColors.neutral500,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: AppSpacing.xl),
                ],

                const Spacer(flex: 3),

                // Fallback options
                TextButton(
                  onPressed: () => context.go(AppRoutes.pinEntry),
                  child: Text(
                    'Use PIN instead',
                    style: AppTypography.body1.copyWith(
                      color: AppColors.primary500,
                    ),
                  ),
                ),
                const SizedBox(height: AppSpacing.sm),
                TextButton(
                  onPressed: () => context.push(AppRoutes.emailLogin),
                  child: Text(
                    'Use password',
                    style: AppTypography.body2.copyWith(
                      color: AppColors.neutral500,
                    ),
                  ),
                ),
                const SizedBox(height: AppSpacing.lg),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
