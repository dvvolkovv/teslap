import 'package:flutter/material.dart';

/// TeslaPay typography scale.
///
/// All styles use the Inter font family (with fallback to the platform default).
/// Monospace text (IBANs, card numbers, tx hashes) uses JetBrains Mono.
abstract final class AppTypography {
  static const String _fontFamily = 'Inter';
  static const String _monoFontFamily = 'JetBrainsMono';

  // ---------------------------------------------------------------------------
  // Display
  // ---------------------------------------------------------------------------

  static const TextStyle display1 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 40,
    fontWeight: FontWeight.w700,
    height: 1.2, // 48px
    letterSpacing: -0.5,
  );

  static const TextStyle display2 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 32,
    fontWeight: FontWeight.w700,
    height: 1.25, // 40px
    letterSpacing: -0.25,
  );

  // ---------------------------------------------------------------------------
  // Headings
  // ---------------------------------------------------------------------------

  static const TextStyle h1 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 24,
    fontWeight: FontWeight.w600,
    height: 1.333, // 32px
    letterSpacing: 0,
  );

  static const TextStyle h2 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 20,
    fontWeight: FontWeight.w600,
    height: 1.4, // 28px
    letterSpacing: 0,
  );

  static const TextStyle h3 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 18,
    fontWeight: FontWeight.w500,
    height: 1.333, // 24px
    letterSpacing: 0,
  );

  // ---------------------------------------------------------------------------
  // Body
  // ---------------------------------------------------------------------------

  static const TextStyle body1 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 16,
    fontWeight: FontWeight.w400,
    height: 1.5, // 24px
    letterSpacing: 0,
  );

  static const TextStyle body2 = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 14,
    fontWeight: FontWeight.w400,
    height: 1.429, // 20px
    letterSpacing: 0.1,
  );

  // ---------------------------------------------------------------------------
  // Supporting
  // ---------------------------------------------------------------------------

  static const TextStyle caption = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 12,
    fontWeight: FontWeight.w400,
    height: 1.333, // 16px
    letterSpacing: 0.2,
  );

  static const TextStyle overline = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 10,
    fontWeight: FontWeight.w600,
    height: 1.6, // 16px
    letterSpacing: 1.5,
  );

  static const TextStyle button = TextStyle(
    fontFamily: _fontFamily,
    fontSize: 16,
    fontWeight: FontWeight.w600,
    height: 1.5, // 24px
    letterSpacing: 0.5,
  );

  // ---------------------------------------------------------------------------
  // Monospace
  // ---------------------------------------------------------------------------

  static const TextStyle mono = TextStyle(
    fontFamily: _monoFontFamily,
    fontSize: 14,
    fontWeight: FontWeight.w400,
    height: 1.429, // 20px
    letterSpacing: 0,
  );
}
