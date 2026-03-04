import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '../../core/theme/app_theme.dart';

/// Input field variant determining specialised formatting and validation.
enum AppInputVariant { standard, password, iban, amount, otp }

/// A design-system compliant text input with floating label, error states,
/// and variant-specific formatting.
class AppInput extends StatefulWidget {
  const AppInput({
    this.label,
    this.hint,
    this.variant = AppInputVariant.standard,
    this.controller,
    this.focusNode,
    this.errorText,
    this.prefixIcon,
    this.suffixIcon,
    this.onChanged,
    this.onSubmitted,
    this.validator,
    this.keyboardType,
    this.textInputAction,
    this.maxLength,
    this.enabled = true,
    this.autofocus = false,
    this.obscureText = false,
    super.key,
  });

  final String? label;
  final String? hint;
  final AppInputVariant variant;
  final TextEditingController? controller;
  final FocusNode? focusNode;
  final String? errorText;
  final Widget? prefixIcon;
  final Widget? suffixIcon;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onSubmitted;
  final FormFieldValidator<String>? validator;
  final TextInputType? keyboardType;
  final TextInputAction? textInputAction;
  final int? maxLength;
  final bool enabled;
  final bool autofocus;
  final bool obscureText;

  @override
  State<AppInput> createState() => _AppInputState();
}

class _AppInputState extends State<AppInput> {
  late bool _obscureText;

  @override
  void initState() {
    super.initState();
    _obscureText = widget.variant == AppInputVariant.password || widget.obscureText;
  }

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: widget.controller,
      focusNode: widget.focusNode,
      obscureText: _obscureText,
      enabled: widget.enabled,
      autofocus: widget.autofocus,
      keyboardType: _keyboardType,
      textInputAction: widget.textInputAction,
      maxLength: widget.maxLength,
      inputFormatters: _inputFormatters,
      style: widget.variant == AppInputVariant.iban
          ? AppTypography.mono
          : AppTypography.body1,
      onChanged: widget.onChanged,
      onFieldSubmitted: widget.onSubmitted,
      validator: widget.validator,
      decoration: InputDecoration(
        labelText: widget.label,
        hintText: widget.hint,
        errorText: widget.errorText,
        prefixIcon: widget.prefixIcon,
        suffixIcon: _buildSuffix(),
        counterText: '',
      ),
    );
  }

  TextInputType get _keyboardType {
    if (widget.keyboardType != null) return widget.keyboardType!;
    return switch (widget.variant) {
      AppInputVariant.amount => const TextInputType.numberWithOptions(
          decimal: true,
        ),
      AppInputVariant.otp => TextInputType.number,
      AppInputVariant.iban => TextInputType.text,
      AppInputVariant.password => TextInputType.visiblePassword,
      _ => TextInputType.text,
    };
  }

  List<TextInputFormatter>? get _inputFormatters {
    return switch (widget.variant) {
      AppInputVariant.amount => [
          FilteringTextInputFormatter.allow(RegExp(r'[\d.]')),
        ],
      AppInputVariant.otp => [
          FilteringTextInputFormatter.digitsOnly,
          LengthLimitingTextInputFormatter(6),
        ],
      AppInputVariant.iban => [
          FilteringTextInputFormatter.allow(RegExp(r'[A-Za-z0-9 ]')),
          _IbanFormatter(),
        ],
      _ => null,
    };
  }

  Widget? _buildSuffix() {
    if (widget.suffixIcon != null) return widget.suffixIcon;
    if (widget.variant == AppInputVariant.password) {
      return IconButton(
        icon: Icon(
          _obscureText ? Icons.visibility_off : Icons.visibility,
          color: AppColors.neutral400,
          size: 24,
        ),
        onPressed: () => setState(() => _obscureText = !_obscureText),
        splashRadius: 20,
        tooltip: _obscureText ? 'Show password' : 'Hide password',
      );
    }
    return null;
  }
}

/// Formats text in IBAN groups of 4 characters separated by spaces.
class _IbanFormatter extends TextInputFormatter {
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    final stripped = newValue.text.replaceAll(' ', '').toUpperCase();
    final buffer = StringBuffer();
    for (var i = 0; i < stripped.length; i++) {
      if (i > 0 && i % 4 == 0) buffer.write(' ');
      buffer.write(stripped[i]);
    }
    final formatted = buffer.toString();
    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(offset: formatted.length),
    );
  }
}
