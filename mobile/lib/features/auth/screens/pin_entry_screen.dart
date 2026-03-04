import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/numpad_widget.dart';
import '../bloc/auth_bloc.dart';

/// PIN entry screen used as a fallback when biometric auth fails.
class PinEntryScreen extends StatefulWidget {
  const PinEntryScreen({super.key});

  @override
  State<PinEntryScreen> createState() => _PinEntryScreenState();
}

class _PinEntryScreenState extends State<PinEntryScreen>
    with SingleTickerProviderStateMixin {
  String _pin = '';
  int _attempts = 0;
  String? _error;

  late final AnimationController _shakeController;

  static const int _pinLength = 6;
  static const int _maxAttempts = 5;

  @override
  void initState() {
    super.initState();
    _shakeController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 400),
    );
  }

  @override
  void dispose() {
    _shakeController.dispose();
    super.dispose();
  }

  void _onDigit(int digit) {
    if (_pin.length >= _pinLength) return;
    setState(() {
      _pin += digit.toString();
      _error = null;
    });

    if (_pin.length == _pinLength) {
      _submit();
    }
  }

  void _onBackspace() {
    if (_pin.isEmpty) return;
    setState(() {
      _pin = _pin.substring(0, _pin.length - 1);
      _error = null;
    });
  }

  void _submit() {
    context.read<AuthBloc>().add(AuthPinSubmitted(pin: _pin));
  }

  void _onAuthFailure(String message) {
    _attempts++;
    if (_attempts >= _maxAttempts) {
      setState(() {
        _error = 'Account locked for 30 minutes.';
        _pin = '';
      });
      return;
    }

    _shakeController.forward(from: 0);
    setState(() {
      _error = 'Incorrect PIN. ${_maxAttempts - _attempts} attempts remaining.';
      _pin = '';
    });
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<AuthBloc, AuthState>(
      listener: (context, state) {
        if (state is AuthAuthenticated) {
          context.go(AppRoutes.home);
        } else if (state is AuthFailure) {
          _onAuthFailure(state.message);
        }
      },
      child: Scaffold(
        appBar: AppBar(
          leading: BackButton(onPressed: () => context.go(AppRoutes.login)),
        ),
        body: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(
              horizontal: AppSpacing.screenMargin,
            ),
            child: Column(
              children: [
                const Spacer(),
                Text('Enter Your PIN', style: AppTypography.h1),
                const SizedBox(height: AppSpacing.xl),

                // Animated PIN dots
                AnimatedBuilder(
                  animation: _shakeController,
                  builder: (context, child) {
                    final sineValue =
                        ((_shakeController.value * 3 * 3.14159 * 2) % (3.14159 * 2));
                    final direction = sineValue < 3.14159 ? 1.0 : -1.0;
                    final offset = 10 * _shakeController.value * direction;
                    return Transform.translate(
                      offset: Offset(offset, 0),
                      child: child,
                    );
                  },
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: List.generate(
                      _pinLength,
                      (index) => Padding(
                        padding: const EdgeInsets.symmetric(
                          horizontal: AppSpacing.sm,
                        ),
                        child: AnimatedContainer(
                          duration: const Duration(milliseconds: 150),
                          width: 12,
                          height: 12,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: index < _pin.length
                                ? (_error != null
                                    ? AppColors.error500
                                    : AppColors.primary500)
                                : AppColors.neutral200,
                          ),
                        ),
                      ),
                    ),
                  ),
                ),

                if (_error != null) ...[
                  const SizedBox(height: AppSpacing.md),
                  Text(
                    _error!,
                    style:
                        AppTypography.body2.copyWith(color: AppColors.error500),
                    textAlign: TextAlign.center,
                  ),
                ],

                const Spacer(),
                NumpadWidget(
                  onDigit: _onDigit,
                  onBackspace: _onBackspace,
                ),
                const SizedBox(height: AppSpacing.md),
                TextButton(
                  onPressed: () {
                    // TODO: navigate to forgot PIN / password flow
                  },
                  child: Text(
                    'Forgot PIN?',
                    style: AppTypography.body2.copyWith(
                      color: AppColors.primary500,
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
