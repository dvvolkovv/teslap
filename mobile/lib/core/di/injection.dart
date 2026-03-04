import 'package:get_it/get_it.dart';

import '../auth/auth_manager.dart';
import '../auth/biometric_auth.dart';
import '../network/api_client.dart';
import '../repositories/account_repository.dart';
import '../repositories/card_repository.dart';
import '../repositories/crypto_repository.dart';
import '../repositories/notification_repository.dart';
import '../repositories/payment_repository.dart';

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

  // Repositories ------------------------------------------------------------
  getIt.registerLazySingleton<AccountRepository>(
    () => AccountRepository(apiClient: getIt<ApiClient>()),
  );
  getIt.registerLazySingleton<PaymentRepository>(
    () => PaymentRepository(apiClient: getIt<ApiClient>()),
  );
  getIt.registerLazySingleton<CardRepository>(
    () => CardRepository(apiClient: getIt<ApiClient>()),
  );
  getIt.registerLazySingleton<CryptoRepository>(
    () => CryptoRepository(apiClient: getIt<ApiClient>()),
  );
  getIt.registerLazySingleton<NotificationRepository>(
    () => NotificationRepository(apiClient: getIt<ApiClient>()),
  );
}
