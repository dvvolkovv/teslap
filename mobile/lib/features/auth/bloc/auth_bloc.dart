import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/auth/auth_manager.dart';
import '../../../core/auth/biometric_auth.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/api_endpoints.dart';

// =============================================================================
// Events
// =============================================================================

abstract class AuthEvent extends Equatable {
  const AuthEvent();

  @override
  List<Object?> get props => [];
}

class AuthCheckRequested extends AuthEvent {
  const AuthCheckRequested();
}

class AuthLoginRequested extends AuthEvent {
  const AuthLoginRequested({
    required this.email,
    required this.password,
  });

  final String email;
  final String password;

  @override
  List<Object?> get props => [email, password];
}

class AuthBiometricRequested extends AuthEvent {
  const AuthBiometricRequested();
}

class AuthPinSubmitted extends AuthEvent {
  const AuthPinSubmitted({required this.pin});

  final String pin;

  @override
  List<Object?> get props => [pin];
}

class AuthLogoutRequested extends AuthEvent {
  const AuthLogoutRequested();
}

// =============================================================================
// States
// =============================================================================

abstract class AuthState extends Equatable {
  const AuthState();

  @override
  List<Object?> get props => [];
}

class AuthInitial extends AuthState {
  const AuthInitial();
}

class AuthLoading extends AuthState {
  const AuthLoading();
}

class AuthAuthenticated extends AuthState {
  const AuthAuthenticated();
}

class AuthUnauthenticated extends AuthState {
  const AuthUnauthenticated();
}

class AuthNeedsOnboarding extends AuthState {
  const AuthNeedsOnboarding();
}

class AuthFailure extends AuthState {
  const AuthFailure({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

// =============================================================================
// BLoC
// =============================================================================

class AuthBloc extends Bloc<AuthEvent, AuthState> {
  AuthBloc({
    required AuthManager authManager,
    required BiometricAuth biometricAuth,
    required ApiClient apiClient,
  })  : _authManager = authManager,
        _biometricAuth = biometricAuth,
        _apiClient = apiClient,
        super(const AuthInitial()) {
    on<AuthCheckRequested>(_onCheckRequested);
    on<AuthLoginRequested>(_onLoginRequested);
    on<AuthBiometricRequested>(_onBiometricRequested);
    on<AuthPinSubmitted>(_onPinSubmitted);
    on<AuthLogoutRequested>(_onLogoutRequested);
  }

  final AuthManager _authManager;
  final BiometricAuth _biometricAuth;
  final ApiClient _apiClient;

  Future<void> _onCheckRequested(
    AuthCheckRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoading());
    final isAuth = await _authManager.isAuthenticated;
    if (!isAuth) {
      // Check if this is a fresh install (no device ID).
      final deviceId = await _authManager.deviceId;
      if (deviceId == null) {
        emit(const AuthNeedsOnboarding());
      } else {
        emit(const AuthUnauthenticated());
      }
      return;
    }

    final isExpired = await _authManager.isSessionExpired;
    if (isExpired) {
      emit(const AuthUnauthenticated());
    } else {
      emit(const AuthAuthenticated());
    }
  }

  Future<void> _onLoginRequested(
    AuthLoginRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoading());
    try {
      final response = await _apiClient.post<Map<String, dynamic>>(
        ApiEndpoints.login,
        data: {
          'email': event.email,
          'password': event.password,
        },
      );

      final data = response.data;
      if (data != null) {
        await _authManager.saveTokens(
          accessToken: data['access_token'] as String,
          refreshToken: data['refresh_token'] as String,
        );
        await _authManager.updateLastActive();
        emit(const AuthAuthenticated());
      } else {
        emit(const AuthFailure(message: 'Invalid server response.'));
      }
    } catch (e) {
      emit(AuthFailure(message: e.toString()));
    }
  }

  Future<void> _onBiometricRequested(
    AuthBiometricRequested event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoading());
    final success = await _biometricAuth.authenticate(
      reason: 'Log in to TeslaPay',
    );

    if (success) {
      await _authManager.updateLastActive();
      emit(const AuthAuthenticated());
    } else {
      emit(const AuthUnauthenticated());
    }
  }

  Future<void> _onPinSubmitted(
    AuthPinSubmitted event,
    Emitter<AuthState> emit,
  ) async {
    emit(const AuthLoading());
    final valid = await _authManager.verifyPin(event.pin);
    if (valid) {
      await _authManager.updateLastActive();
      emit(const AuthAuthenticated());
    } else {
      emit(const AuthFailure(message: 'Incorrect PIN. Please try again.'));
    }
  }

  Future<void> _onLogoutRequested(
    AuthLogoutRequested event,
    Emitter<AuthState> emit,
  ) async {
    try {
      await _apiClient.post<void>(
        ApiEndpoints.logout,
        data: {'all_sessions': false},
      );
    } catch (_) {
      // Best-effort server logout; always clear local state.
    }
    await _authManager.clearSession();
    emit(const AuthUnauthenticated());
  }
}
