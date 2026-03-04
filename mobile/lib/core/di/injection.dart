import 'package:get_it/get_it.dart';

import '../auth/auth_manager.dart';
import '../auth/biometric_auth.dart';
import '../network/api_client.dart';

/// Global service locator.
final GetIt getIt = GetIt.instance;

/// Registers all singleton services used across the application.
///
/// Must be called before [runApp] in `main.dart`.
Future<void> configureDependencies() async {
  // Core services -----------------------------------------------------------
  getIt.registerLazySingleton<AuthManager>(() => AuthManager());
  getIt.registerLazySingleton<BiometricAuth>(() => BiometricAuth());
  getIt.registerLazySingleton<ApiClient>(
    () => ApiClient(authManager: getIt<AuthManager>()),
  );
}
