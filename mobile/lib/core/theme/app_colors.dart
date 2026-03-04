import 'package:flutter/material.dart';

/// TeslaPay color palette derived from the design system.
///
/// Primary blue conveys trust and stability. Accent teal-green
/// differentiates crypto features from traditional banking.
abstract final class AppColors {
  // ---------------------------------------------------------------------------
  // Primary
  // ---------------------------------------------------------------------------
  static const Color primary50 = Color(0xFFF0F6FF);
  static const Color primary100 = Color(0xFFE6F0FF);
  static const Color primary400 = Color(0xFF3385FF);
  static const Color primary500 = Color(0xFF0066FF);
  static const Color primary600 = Color(0xFF0052CC);

  // ---------------------------------------------------------------------------
  // Accent (Crypto)
  // ---------------------------------------------------------------------------
  static const Color accent100 = Color(0xFFE6FBF5);
  static const Color accent500 = Color(0xFF00D4AA);
  static const Color accent600 = Color(0xFF00B892);

  // ---------------------------------------------------------------------------
  // Semantic
  // ---------------------------------------------------------------------------
  static const Color success100 = Color(0xFFE8F5E9);
  static const Color success500 = Color(0xFF00C853);

  static const Color warning100 = Color(0xFFFFF3E0);
  static const Color warning500 = Color(0xFFFF9800);

  static const Color error100 = Color(0xFFFFEBEE);
  static const Color error500 = Color(0xFFF44336);

  static const Color info100 = Color(0xFFE3F2FD);
  static const Color info500 = Color(0xFF2196F3);

  // ---------------------------------------------------------------------------
  // Neutral (Light mode)
  // ---------------------------------------------------------------------------
  static const Color neutral900 = Color(0xFF0D1B2A);
  static const Color neutral800 = Color(0xFF1B2838);
  static const Color neutral700 = Color(0xFF2E3A4D);
  static const Color neutral500 = Color(0xFF6B7B8F);
  static const Color neutral400 = Color(0xFF94A3B8);
  static const Color neutral300 = Color(0xFFB0BEC5);
  static const Color neutral200 = Color(0xFFD1D9E6);
  static const Color neutral100 = Color(0xFFEDF1F7);
  static const Color neutral50 = Color(0xFFF7F9FC);
  static const Color white = Color(0xFFFFFFFF);

  // ---------------------------------------------------------------------------
  // Dark mode
  // ---------------------------------------------------------------------------
  static const Color darkBackground = Color(0xFF0D1117);
  static const Color darkSurface = Color(0xFF161B22);
  static const Color darkSurfaceElevated = Color(0xFF1C2333);
  static const Color darkBorder = Color(0xFF30363D);
  static const Color darkTextPrimary = Color(0xFFF0F6FC);
  static const Color darkTextSecondary = Color(0xFF8B949E);
  static const Color darkTextTertiary = Color(0xFF6E7681);

  // ---------------------------------------------------------------------------
  // Gradients
  // ---------------------------------------------------------------------------
  static const LinearGradient brandGradient = LinearGradient(
    begin: Alignment(-0.7, -0.7),
    end: Alignment(0.7, 0.7),
    colors: [primary500, accent500],
  );

  static const LinearGradient cryptoGradient = LinearGradient(
    begin: Alignment(-0.7, -0.7),
    end: Alignment(0.7, 0.7),
    colors: [accent500, Color(0xFF00B4D8)],
  );

  static const LinearGradient cardGradient = LinearGradient(
    begin: Alignment(-0.5, -0.8),
    end: Alignment(0.5, 0.8),
    colors: [neutral800, neutral900],
  );
}
