import 'package:flutter_secure_storage/flutter_secure_storage.dart';

/// Manages JWT access and refresh tokens using platform secure storage.
///
/// All tokens are stored encrypted via [FlutterSecureStorage] (Keychain on iOS,
/// Encrypted SharedPreferences on Android).
class AuthManager {
  AuthManager({FlutterSecureStorage? storage})
      : _storage = storage ?? const FlutterSecureStorage();

  final FlutterSecureStorage _storage;

  static const String _accessTokenKey = 'teslapay_access_token';
  static const String _refreshTokenKey = 'teslapay_refresh_token';
  static const String _deviceIdKey = 'teslapay_device_id';
  static const String _pinHashKey = 'teslapay_pin_hash';
  static const String _biometricEnabledKey = 'teslapay_biometric_enabled';
  static const String _lastActiveKey = 'teslapay_last_active';

  /// Duration after which the user must re-authenticate on resume.
  static const Duration sessionTimeout = Duration(minutes: 5);

  // ---------------------------------------------------------------------------
  // Token management
  // ---------------------------------------------------------------------------

  Future<String?> get accessToken => _storage.read(key: _accessTokenKey);

  Future<String?> get refreshToken => _storage.read(key: _refreshTokenKey);

  Future<String?> get deviceId => _storage.read(key: _deviceIdKey);

  Future<void> saveTokens({
    required String accessToken,
    required String refreshToken,
  }) async {
    await Future.wait([
      _storage.write(key: _accessTokenKey, value: accessToken),
      _storage.write(key: _refreshTokenKey, value: refreshToken),
    ]);
  }

  Future<void> saveDeviceId(String deviceId) =>
      _storage.write(key: _deviceIdKey, value: deviceId);

  Future<bool> get isAuthenticated async {
    final token = await accessToken;
    return token != null && token.isNotEmpty;
  }

  Future<void> clearSession() async {
    await Future.wait([
      _storage.delete(key: _accessTokenKey),
      _storage.delete(key: _refreshTokenKey),
    ]);
  }

  // ---------------------------------------------------------------------------
  // PIN management
  // ---------------------------------------------------------------------------

  Future<void> savePin(String pinHash) =>
      _storage.write(key: _pinHashKey, value: pinHash);

  Future<String?> get storedPinHash => _storage.read(key: _pinHashKey);

  Future<bool> verifyPin(String pinHash) async {
    final stored = await storedPinHash;
    return stored != null && stored == pinHash;
  }

  // ---------------------------------------------------------------------------
  // Biometric preference
  // ---------------------------------------------------------------------------

  Future<void> setBiometricEnabled(bool enabled) =>
      _storage.write(key: _biometricEnabledKey, value: enabled.toString());

  Future<bool> get isBiometricEnabled async {
    final value = await _storage.read(key: _biometricEnabledKey);
    return value == 'true';
  }

  // ---------------------------------------------------------------------------
  // Session timeout tracking
  // ---------------------------------------------------------------------------

  Future<void> updateLastActive() => _storage.write(
        key: _lastActiveKey,
        value: DateTime.now().toIso8601String(),
      );

  /// Returns `true` if the session has been idle longer than [sessionTimeout].
  Future<bool> get isSessionExpired async {
    final lastActive = await _storage.read(key: _lastActiveKey);
    if (lastActive == null) return true;

    final lastActiveDate = DateTime.tryParse(lastActive);
    if (lastActiveDate == null) return true;

    return DateTime.now().difference(lastActiveDate) > sessionTimeout;
  }

  // ---------------------------------------------------------------------------
  // Full wipe (logout everywhere, account closure)
  // ---------------------------------------------------------------------------

  Future<void> clearAll() => _storage.deleteAll();
}
