import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_input.dart';

/// Registration step 1: email, phone, and consent.
class RegisterScreen extends StatefulWidget {
  const RegisterScreen({super.key});

  @override
  State<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends State<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  bool _termsAccepted = false;
  bool _isLoading = false;

  @override
  void dispose() {
    _emailController.dispose();
    _phoneController.dispose();
    super.dispose();
  }

  bool get _isFormValid =>
      _emailController.text.isNotEmpty &&
      _phoneController.text.isNotEmpty &&
      _termsAccepted;

  Future<void> _onContinue() async {
    if (!_formKey.currentState!.validate() || !_termsAccepted) return;

    setState(() => _isLoading = true);

    // TODO: Call POST /auth/register and navigate on success.
    await Future<void>.delayed(const Duration(seconds: 1));

    if (mounted) {
      setState(() => _isLoading = false);
      context.go(AppRoutes.pinSetup);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: BackButton(onPressed: () => context.go(AppRoutes.welcome)),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding:
              const EdgeInsets.symmetric(horizontal: AppSpacing.screenMargin),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const SizedBox(height: AppSpacing.md),
                Text('Create Your Account', style: AppTypography.h1),
                const SizedBox(height: AppSpacing.sm),
                Text(
                  'Step 1 of 3',
                  style:
                      AppTypography.body2.copyWith(color: AppColors.neutral500),
                ),
                const SizedBox(height: AppSpacing.sm),
                // Progress bar
                ClipRRect(
                  borderRadius: BorderRadius.circular(AppRadius.full),
                  child: LinearProgressIndicator(
                    value: 0.33,
                    backgroundColor: AppColors.neutral100,
                    valueColor: const AlwaysStoppedAnimation<Color>(
                      AppColors.primary500,
                    ),
                    minHeight: 4,
                  ),
                ),

                const SizedBox(height: AppSpacing.xl),
                AppInput(
                  label: 'Email address',
                  hint: 'your@email.com',
                  controller: _emailController,
                  keyboardType: TextInputType.emailAddress,
                  textInputAction: TextInputAction.next,
                  onChanged: (_) => setState(() {}),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter your email';
                    }
                    final emailRegex = RegExp(r'^[^@]+@[^@]+\.[^@]+$');
                    if (!emailRegex.hasMatch(value)) {
                      return 'Please enter a valid email';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: AppSpacing.md),
                AppInput(
                  label: 'Phone number',
                  hint: '+370 612 345 67',
                  controller: _phoneController,
                  keyboardType: TextInputType.phone,
                  textInputAction: TextInputAction.done,
                  onChanged: (_) => setState(() {}),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter your phone number';
                    }
                    return null;
                  },
                ),

                const SizedBox(height: AppSpacing.lg),
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Semantics(
                      label: 'Accept terms of service and privacy policy',
                      child: Checkbox(
                        value: _termsAccepted,
                        onChanged: (v) =>
                            setState(() => _termsAccepted = v ?? false),
                        activeColor: AppColors.primary500,
                      ),
                    ),
                    Expanded(
                      child: Padding(
                        padding: const EdgeInsets.only(top: 12),
                        child: RichText(
                          text: TextSpan(
                            style: AppTypography.body2.copyWith(
                              color: AppColors.neutral700,
                            ),
                            children: [
                              const TextSpan(
                                text: 'I agree to the ',
                              ),
                              TextSpan(
                                text: 'Terms of Service',
                                style: const TextStyle(
                                  color: AppColors.primary500,
                                ),
                                recognizer: TapGestureRecognizer()
                                  ..onTap = () {
                                    // TODO: open Terms of Service
                                  },
                              ),
                              const TextSpan(text: ' and '),
                              TextSpan(
                                text: 'Privacy Policy',
                                style: const TextStyle(
                                  color: AppColors.primary500,
                                ),
                                recognizer: TapGestureRecognizer()
                                  ..onTap = () {
                                    // TODO: open Privacy Policy
                                  },
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ],
                ),

                const SizedBox(height: AppSpacing.xl),
                AppButton(
                  label: 'Continue',
                  onPressed: _isFormValid ? _onContinue : null,
                  isLoading: _isLoading,
                ),
                const SizedBox(height: AppSpacing.xl),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
