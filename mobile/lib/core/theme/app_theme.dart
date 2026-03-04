import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app_colors.dart';
import 'app_spacing.dart';
import 'app_typography.dart';

export 'app_colors.dart';
export 'app_spacing.dart';
export 'app_typography.dart';

/// Provides the light and dark [ThemeData] for TeslaPay.
abstract final class AppTheme {
  // ---------------------------------------------------------------------------
  // Light theme
  // ---------------------------------------------------------------------------
  static ThemeData get light {
    final colorScheme = ColorScheme.light(
      primary: AppColors.primary500,
      onPrimary: AppColors.white,
      primaryContainer: AppColors.primary100,
      onPrimaryContainer: AppColors.primary600,
      secondary: AppColors.accent500,
      onSecondary: AppColors.white,
      secondaryContainer: AppColors.accent100,
      onSecondaryContainer: AppColors.accent600,
      error: AppColors.error500,
      onError: AppColors.white,
      errorContainer: AppColors.error100,
      surface: AppColors.white,
      onSurface: AppColors.neutral900,
      onSurfaceVariant: AppColors.neutral700,
      outline: AppColors.neutral200,
      outlineVariant: AppColors.neutral100,
      shadow: AppColors.neutral900.withValues(alpha: 0.08),
    );

    return ThemeData(
      useMaterial3: true,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: AppColors.neutral50,
      fontFamily: 'Inter',
      textTheme: _buildTextTheme(AppColors.neutral900),
      appBarTheme: AppBarTheme(
        backgroundColor: Colors.transparent,
        foregroundColor: AppColors.neutral900,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: true,
        titleTextStyle: AppTypography.h2.copyWith(color: AppColors.neutral900),
        systemOverlayStyle: SystemUiOverlayStyle.dark,
      ),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
        backgroundColor: AppColors.white,
        selectedItemColor: AppColors.primary500,
        unselectedItemColor: AppColors.neutral400,
        type: BottomNavigationBarType.fixed,
        selectedLabelStyle: AppTypography.caption.copyWith(
          color: AppColors.primary500,
        ),
        unselectedLabelStyle: AppTypography.caption.copyWith(
          color: AppColors.neutral500,
        ),
        elevation: 0,
      ),
      cardTheme: CardThemeData(
        color: AppColors.white,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppRadius.md),
          side: BorderSide(color: AppColors.neutral200.withValues(alpha: 0.5)),
        ),
        margin: EdgeInsets.zero,
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.neutral100,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.md,
          vertical: AppSpacing.md,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.neutral200),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.neutral200),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(
            color: AppColors.primary500,
            width: 2,
          ),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.error500),
        ),
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(
            color: AppColors.error500,
            width: 2,
          ),
        ),
        hintStyle: AppTypography.body1.copyWith(color: AppColors.neutral500),
        labelStyle: AppTypography.caption.copyWith(color: AppColors.neutral500),
        errorStyle: AppTypography.caption.copyWith(color: AppColors.error500),
        floatingLabelStyle: AppTypography.caption.copyWith(
          color: AppColors.primary500,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.primary500,
          foregroundColor: AppColors.white,
          minimumSize: const Size(double.infinity, 56),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppRadius.sm),
          ),
          textStyle: AppTypography.button,
          elevation: 0,
          disabledBackgroundColor: AppColors.primary500.withValues(alpha: 0.4),
          disabledForegroundColor: AppColors.white.withValues(alpha: 0.7),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.primary500,
          minimumSize: const Size(double.infinity, 56),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppRadius.sm),
          ),
          side: const BorderSide(color: AppColors.primary500, width: 1.5),
          textStyle: AppTypography.button,
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.primary500,
          textStyle: AppTypography.button,
        ),
      ),
      dividerTheme: const DividerThemeData(
        color: AppColors.neutral200,
        thickness: 1,
        space: 0,
      ),
      bottomSheetTheme: const BottomSheetThemeData(
        backgroundColor: AppColors.white,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(
            top: Radius.circular(AppRadius.lg),
          ),
        ),
        showDragHandle: true,
        dragHandleColor: AppColors.neutral200,
        dragHandleSize: Size(40, 4),
      ),
      switchTheme: SwitchThemeData(
        thumbColor: WidgetStateProperty.all(AppColors.white),
        trackColor: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return AppColors.primary500;
          }
          return AppColors.neutral200;
        }),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Dark theme
  // ---------------------------------------------------------------------------
  static ThemeData get dark {
    final colorScheme = ColorScheme.dark(
      primary: AppColors.primary400,
      onPrimary: AppColors.white,
      primaryContainer: AppColors.primary600,
      onPrimaryContainer: AppColors.primary100,
      secondary: AppColors.accent500,
      onSecondary: AppColors.neutral900,
      secondaryContainer: AppColors.accent600,
      onSecondaryContainer: AppColors.accent100,
      error: AppColors.error500,
      onError: AppColors.white,
      errorContainer: AppColors.error100,
      surface: AppColors.darkSurface,
      onSurface: AppColors.darkTextPrimary,
      onSurfaceVariant: AppColors.darkTextSecondary,
      outline: AppColors.darkBorder,
      outlineVariant: AppColors.darkSurfaceElevated,
      shadow: Colors.transparent,
    );

    return ThemeData(
      useMaterial3: true,
      colorScheme: colorScheme,
      scaffoldBackgroundColor: AppColors.darkBackground,
      fontFamily: 'Inter',
      textTheme: _buildTextTheme(AppColors.darkTextPrimary),
      appBarTheme: AppBarTheme(
        backgroundColor: Colors.transparent,
        foregroundColor: AppColors.darkTextPrimary,
        elevation: 0,
        scrolledUnderElevation: 0,
        centerTitle: true,
        titleTextStyle: AppTypography.h2.copyWith(
          color: AppColors.darkTextPrimary,
        ),
        systemOverlayStyle: SystemUiOverlayStyle.light,
      ),
      bottomNavigationBarTheme: BottomNavigationBarThemeData(
        backgroundColor: AppColors.darkSurface,
        selectedItemColor: AppColors.primary400,
        unselectedItemColor: AppColors.darkTextTertiary,
        type: BottomNavigationBarType.fixed,
        selectedLabelStyle: AppTypography.caption.copyWith(
          color: AppColors.primary400,
        ),
        unselectedLabelStyle: AppTypography.caption.copyWith(
          color: AppColors.darkTextTertiary,
        ),
        elevation: 0,
      ),
      cardTheme: CardThemeData(
        color: AppColors.darkSurface,
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(AppRadius.md),
          side: const BorderSide(color: AppColors.darkBorder),
        ),
        margin: EdgeInsets.zero,
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: AppColors.darkSurface,
        contentPadding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.md,
          vertical: AppSpacing.md,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.darkBorder),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.darkBorder),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(
            color: AppColors.primary400,
            width: 2,
          ),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.error500),
        ),
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(AppRadius.sm),
          borderSide: const BorderSide(color: AppColors.error500, width: 2),
        ),
        hintStyle: AppTypography.body1.copyWith(
          color: AppColors.darkTextTertiary,
        ),
        labelStyle: AppTypography.caption.copyWith(
          color: AppColors.darkTextTertiary,
        ),
        errorStyle: AppTypography.caption.copyWith(color: AppColors.error500),
        floatingLabelStyle: AppTypography.caption.copyWith(
          color: AppColors.primary400,
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: AppColors.primary500,
          foregroundColor: AppColors.white,
          minimumSize: const Size(double.infinity, 56),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppRadius.sm),
          ),
          textStyle: AppTypography.button,
          elevation: 0,
          disabledBackgroundColor: AppColors.primary500.withValues(alpha: 0.4),
          disabledForegroundColor: AppColors.white.withValues(alpha: 0.7),
        ),
      ),
      outlinedButtonTheme: OutlinedButtonThemeData(
        style: OutlinedButton.styleFrom(
          foregroundColor: AppColors.primary400,
          minimumSize: const Size(double.infinity, 56),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(AppRadius.sm),
          ),
          side: const BorderSide(color: AppColors.primary400, width: 1.5),
          textStyle: AppTypography.button,
        ),
      ),
      textButtonTheme: TextButtonThemeData(
        style: TextButton.styleFrom(
          foregroundColor: AppColors.primary400,
          textStyle: AppTypography.button,
        ),
      ),
      dividerTheme: const DividerThemeData(
        color: AppColors.darkBorder,
        thickness: 1,
        space: 0,
      ),
      bottomSheetTheme: const BottomSheetThemeData(
        backgroundColor: AppColors.darkSurface,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.vertical(
            top: Radius.circular(AppRadius.lg),
          ),
        ),
        showDragHandle: true,
        dragHandleColor: AppColors.darkBorder,
        dragHandleSize: Size(40, 4),
      ),
      switchTheme: SwitchThemeData(
        thumbColor: WidgetStateProperty.all(AppColors.white),
        trackColor: WidgetStateProperty.resolveWith((states) {
          if (states.contains(WidgetState.selected)) {
            return AppColors.primary400;
          }
          return AppColors.darkBorder;
        }),
      ),
    );
  }

  // ---------------------------------------------------------------------------
  // Helpers
  // ---------------------------------------------------------------------------
  static TextTheme _buildTextTheme(Color baseColor) {
    return TextTheme(
      displayLarge: AppTypography.display1.copyWith(color: baseColor),
      displayMedium: AppTypography.display2.copyWith(color: baseColor),
      headlineLarge: AppTypography.h1.copyWith(color: baseColor),
      headlineMedium: AppTypography.h2.copyWith(color: baseColor),
      headlineSmall: AppTypography.h3.copyWith(color: baseColor),
      bodyLarge: AppTypography.body1.copyWith(color: baseColor),
      bodyMedium: AppTypography.body2.copyWith(color: baseColor),
      labelLarge: AppTypography.button.copyWith(color: baseColor),
      labelMedium: AppTypography.caption.copyWith(color: baseColor),
      labelSmall: AppTypography.overline.copyWith(color: baseColor),
    );
  }
}
