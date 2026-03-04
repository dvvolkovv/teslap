import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '../../features/auth/screens/login_screen.dart';
import '../../features/auth/screens/pin_entry_screen.dart';
import '../../features/cards/screens/card_screen.dart';
import '../../features/crypto/screens/crypto_screen.dart';
import '../../features/home/screens/home_screen.dart';
import '../../features/onboarding/screens/pin_setup_screen.dart';
import '../../features/onboarding/screens/register_screen.dart';
import '../../features/onboarding/screens/splash_screen.dart';
import '../../features/onboarding/screens/welcome_screen.dart';
import '../../features/payments/screens/payments_screen.dart';
import '../../features/payments/screens/send_money_screen.dart';
import '../../features/profile/screens/profile_screen.dart';
import '../../shared/widgets/app_bottom_nav.dart';

/// Named route constants for type-safe navigation.
abstract final class AppRoutes {
  static const String splash = '/splash';
  static const String welcome = '/welcome';
  static const String register = '/register';
  static const String pinSetup = '/pin-setup';
  static const String login = '/login';
  static const String pinEntry = '/pin-entry';

  // Tabs
  static const String home = '/home';
  static const String payments = '/payments';
  static const String card = '/card';
  static const String crypto = '/crypto';
  static const String profile = '/profile';

  // Sub-routes
  static const String sendMoney = '/payments/send';
}

final GlobalKey<NavigatorState> _rootNavigatorKey =
    GlobalKey<NavigatorState>(debugLabel: 'root');

final GlobalKey<NavigatorState> _shellNavigatorKey =
    GlobalKey<NavigatorState>(debugLabel: 'shell');

/// Central GoRouter configuration for TeslaPay.
///
/// Pre-auth routes (splash, welcome, register, login) sit at root level.
/// Post-auth routes are nested inside a [ShellRoute] that provides the
/// bottom navigation bar.
final GoRouter appRouter = GoRouter(
  navigatorKey: _rootNavigatorKey,
  initialLocation: AppRoutes.splash,
  debugLogDiagnostics: true,
  routes: [
    // -----------------------------------------------------------------------
    // Pre-authentication routes
    // -----------------------------------------------------------------------
    GoRoute(
      path: AppRoutes.splash,
      name: 'splash',
      builder: (context, state) => const SplashScreen(),
    ),
    GoRoute(
      path: AppRoutes.welcome,
      name: 'welcome',
      builder: (context, state) => const WelcomeScreen(),
    ),
    GoRoute(
      path: AppRoutes.register,
      name: 'register',
      builder: (context, state) => const RegisterScreen(),
    ),
    GoRoute(
      path: AppRoutes.pinSetup,
      name: 'pinSetup',
      builder: (context, state) => const PinSetupScreen(),
    ),
    GoRoute(
      path: AppRoutes.login,
      name: 'login',
      builder: (context, state) => const LoginScreen(),
    ),
    GoRoute(
      path: AppRoutes.pinEntry,
      name: 'pinEntry',
      builder: (context, state) => const PinEntryScreen(),
    ),

    // -----------------------------------------------------------------------
    // Main application shell (authenticated)
    // -----------------------------------------------------------------------
    ShellRoute(
      navigatorKey: _shellNavigatorKey,
      builder: (context, state, child) => AppScaffold(child: child),
      routes: [
        GoRoute(
          path: AppRoutes.home,
          name: 'home',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: HomeScreen(),
          ),
        ),
        GoRoute(
          path: AppRoutes.payments,
          name: 'payments',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: PaymentsScreen(),
          ),
          routes: [
            GoRoute(
              path: 'send',
              name: 'sendMoney',
              parentNavigatorKey: _rootNavigatorKey,
              builder: (context, state) => const SendMoneyScreen(),
            ),
          ],
        ),
        GoRoute(
          path: AppRoutes.card,
          name: 'card',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: CardScreen(),
          ),
        ),
        GoRoute(
          path: AppRoutes.crypto,
          name: 'crypto',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: CryptoScreen(),
          ),
        ),
        GoRoute(
          path: AppRoutes.profile,
          name: 'profile',
          pageBuilder: (context, state) => const NoTransitionPage(
            child: ProfileScreen(),
          ),
        ),
      ],
    ),
  ],
);

/// Shell scaffold providing the bottom navigation bar for all authenticated
/// tab screens.
class AppScaffold extends StatelessWidget {
  const AppScaffold({required this.child, super.key});

  final Widget child;

  int _currentIndex(BuildContext context) {
    final location = GoRouterState.of(context).uri.toString();
    if (location.startsWith(AppRoutes.payments)) return 1;
    if (location.startsWith(AppRoutes.card)) return 2;
    if (location.startsWith(AppRoutes.crypto)) return 3;
    if (location.startsWith(AppRoutes.profile)) return 4;
    return 0;
  }

  void _onItemTapped(BuildContext context, int index) {
    switch (index) {
      case 0:
        context.go(AppRoutes.home);
      case 1:
        context.go(AppRoutes.payments);
      case 2:
        context.go(AppRoutes.card);
      case 3:
        context.go(AppRoutes.crypto);
      case 4:
        context.go(AppRoutes.profile);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: child,
      bottomNavigationBar: AppBottomNav(
        currentIndex: _currentIndex(context),
        onTap: (index) => _onItemTapped(context, index),
      ),
    );
  }
}
