import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import 'core/auth/auth_manager.dart';
import 'core/auth/biometric_auth.dart';
import 'core/di/injection.dart';
import 'core/network/api_client.dart';
import 'core/routing/app_router.dart';
import 'core/theme/app_theme.dart';
import 'features/auth/bloc/auth_bloc.dart';

/// Root widget for the TeslaPay application.
///
/// Provides:
/// - Material 3 theming (light + dark from design system)
/// - GoRouter navigation
/// - Global BLoC providers (auth)
class TeslaPayApp extends StatelessWidget {
  const TeslaPayApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider<AuthBloc>(
          create: (_) => AuthBloc(
            authManager: getIt<AuthManager>(),
            biometricAuth: getIt<BiometricAuth>(),
            apiClient: getIt<ApiClient>(),
          ),
        ),
      ],
      child: MaterialApp.router(
        title: 'TeslaPay',
        debugShowCheckedModeBanner: false,

        // Theme
        theme: AppTheme.light,
        darkTheme: AppTheme.dark,
        themeMode: ThemeMode.system,

        // Navigation
        routerConfig: appRouter,
      ),
    );
  }
}
