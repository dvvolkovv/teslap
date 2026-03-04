import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/numpad_widget.dart';

/// PIN setup screen where the user chooses a 6-digit PIN and then confirms it.
class PinSetupScreen extends StatefulWidget {
  const PinSetupScreen({super.key});

  @override
  State<PinSetupScreen> createState() => _PinSetupScreenState();
}

class _PinSetupScreenState extends State<PinSetupScreen> {
  String _pin = '';
  String? _firstPin;
  bool _isConfirming = false;
  String? _error;

  static const int _pinLength = 6;

  void _onDigit(int digit) {
    if (_pin.length >= _pinLength) return;
    setState(() {
      _pin += digit.toString();
      _error = null;
    });

    if (_pin.length == _pinLength) {
      _onPinComplete();
    }
  }

  void _onBackspace() {
    if (_pin.isEmpty) return;
    setState(() {
      _pin = _pin.substring(0, _pin.length - 1);
      _error = null;
    });
  }

  void _onPinComplete() {
    if (!_isConfirming) {
      // First entry -- store and ask for confirmation.
      setState(() {
        _firstPin = _pin;
        _pin = '';
        _isConfirming = true;
      });
    } else {
      // Second entry -- check match.
      if (_pin == _firstPin) {
        // TODO: Hash and persist PIN via AuthManager, then navigate.
        context.go('/login');
      } else {
        setState(() {
          _pin = '';
          _firstPin = null;
          _isConfirming = false;
          _error = 'PINs do not match. Please try again.';
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final title = _isConfirming ? 'Confirm Your PIN' : 'Create Your PIN';
    final subtitle = _isConfirming
        ? 'Re-enter your 6-digit PIN to confirm.'
        : 'Choose a 6-digit PIN for quick access to your account.';

    return Scaffold(
      appBar: AppBar(
        leading: BackButton(
          onPressed: () {
            if (_isConfirming) {
              setState(() {
                _isConfirming = false;
                _pin = '';
                _firstPin = null;
                _error = null;
              });
            } else {
              context.pop();
            }
          },
        ),
      ),
      body: SafeArea(
        child: Padding(
          padding:
              const EdgeInsets.symmetric(horizontal: AppSpacing.screenMargin),
          child: Column(
            children: [
              const Spacer(),
              Text(title, style: AppTypography.h1),
              const SizedBox(height: AppSpacing.sm),
              Text(
                subtitle,
                style:
                    AppTypography.body2.copyWith(color: AppColors.neutral500),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: AppSpacing.xl),

              // PIN dots
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(
                  _pinLength,
                  (index) => Padding(
                    padding:
                        const EdgeInsets.symmetric(horizontal: AppSpacing.sm),
                    child: AnimatedContainer(
                      duration: const Duration(milliseconds: 150),
                      width: 12,
                      height: 12,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: index < _pin.length
                            ? AppColors.primary500
                            : AppColors.neutral200,
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
              const SizedBox(height: AppSpacing.lg),
            ],
          ),
        ),
      ),
    );
  }
}
