import 'package:flutter/services.dart';
import 'package:local_auth/local_auth.dart';

/// Wraps the [LocalAuthentication] plugin to provide biometric
/// authentication (Face ID, fingerprint, iris) across iOS and Android.
class BiometricAuth {
  BiometricAuth({LocalAuthentication? auth})
      : _auth = auth ?? LocalAuthentication();

  final LocalAuthentication _auth;

  /// Returns `true` if the device supports biometric authentication
  /// and at least one biometric is enrolled.
  Future<bool> get isAvailable async {
    try {
      final canCheck = await _auth.canCheckBiometrics;
      final isSupported = await _auth.isDeviceSupported();
      return canCheck && isSupported;
    } on PlatformException {
      return false;
    }
  }

  /// Returns the list of enrolled biometric types.
  Future<List<BiometricType>> get enrolledBiometrics async {
    try {
      return await _auth.getAvailableBiometrics();
    } on PlatformException {
      return [];
    }
  }

  /// Prompts the user for biometric authentication.
  ///
  /// [reason] is displayed in the system prompt.
  /// Returns `true` if authentication succeeded.
  Future<bool> authenticate({
    String reason = 'Authenticate to access TeslaPay',
  }) async {
    try {
      return await _auth.authenticate(
        localizedReason: reason,
        options: const AuthenticationOptions(
          stickyAuth: true,
          biometricOnly: true,
        ),
      );
    } on PlatformException {
      return false;
    }
  }

  /// Returns a human-readable label for the primary biometric type.
  Future<String> get biometricLabel async {
    final biometrics = await enrolledBiometrics;
    if (biometrics.contains(BiometricType.face)) {
      return 'Face ID';
    } else if (biometrics.contains(BiometricType.fingerprint)) {
      return 'Fingerprint';
    } else if (biometrics.contains(BiometricType.iris)) {
      return 'Iris';
    }
    return 'Biometric';
  }
}
