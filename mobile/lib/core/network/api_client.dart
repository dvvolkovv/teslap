import 'dart:io';

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:uuid/uuid.dart';

import '../auth/auth_manager.dart';
import 'api_endpoints.dart';
import 'api_exceptions.dart';

/// Dio-based HTTP client for all TeslaPay API communication.
///
/// Features:
/// - Automatic JWT injection via [AuthManager]
/// - Transparent token refresh on 401
/// - Idempotency-Key header injection for mutating requests
/// - Standard error mapping to typed exceptions
/// - Request / response logging in debug mode
class ApiClient {
  ApiClient({
    required AuthManager authManager,
    String? baseUrl,
  })  : _authManager = authManager,
        _dio = Dio(
          BaseOptions(
            baseUrl: baseUrl ?? ApiEndpoints.baseUrl,
            connectTimeout: const Duration(seconds: 15),
            receiveTimeout: const Duration(seconds: 30),
            headers: {
              HttpHeaders.acceptHeader: 'application/json',
              HttpHeaders.contentTypeHeader: 'application/json',
            },
          ),
        ) {
    _dio.interceptors.addAll([
      _AuthInterceptor(authManager: _authManager, dio: _dio),
      _IdempotencyInterceptor(),
      if (kDebugMode) _LoggingInterceptor(),
    ]);
  }

  final Dio _dio;
  final AuthManager _authManager;

  // ---------------------------------------------------------------------------
  // Public convenience methods
  // ---------------------------------------------------------------------------

  Future<Response<T>> get<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) {
    return _handleRequest(
      () => _dio.get<T>(
        path,
        queryParameters: queryParameters,
        options: options,
      ),
    );
  }

  Future<Response<T>> post<T>(
    String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) {
    return _handleRequest(
      () => _dio.post<T>(
        path,
        data: data,
        queryParameters: queryParameters,
        options: options,
      ),
    );
  }

  Future<Response<T>> put<T>(
    String path, {
    dynamic data,
    Options? options,
  }) {
    return _handleRequest(
      () => _dio.put<T>(path, data: data, options: options),
    );
  }

  Future<Response<T>> patch<T>(
    String path, {
    dynamic data,
    Options? options,
  }) {
    return _handleRequest(
      () => _dio.patch<T>(path, data: data, options: options),
    );
  }

  Future<Response<T>> delete<T>(
    String path, {
    dynamic data,
    Options? options,
  }) {
    return _handleRequest(
      () => _dio.delete<T>(path, data: data, options: options),
    );
  }

  // ---------------------------------------------------------------------------
  // Error mapping
  // ---------------------------------------------------------------------------

  Future<Response<T>> _handleRequest<T>(
    Future<Response<T>> Function() request,
  ) async {
    try {
      return await request();
    } on DioException catch (e) {
      throw _mapDioException(e);
    }
  }

  Exception _mapDioException(DioException e) {
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.connectionError:
        return const NetworkException();

      case DioExceptionType.badResponse:
        final data = e.response?.data;
        if (data is Map<String, dynamic>) {
          return ApiException.fromJson(data);
        }
        return ApiException(
          statusCode: e.response?.statusCode ?? 0,
          errorCode: 'UNKNOWN',
          title: 'Error',
          detail: e.message ?? 'An unexpected error occurred.',
        );

      case DioExceptionType.cancel:
        return const ApiException(
          statusCode: 0,
          errorCode: 'CANCELLED',
          title: 'Request Cancelled',
          detail: 'The request was cancelled.',
        );

      default:
        return const ApiException(
          statusCode: 0,
          errorCode: 'UNKNOWN',
          title: 'Unknown Error',
          detail: 'An unexpected error occurred.',
        );
    }
  }
}

// =============================================================================
// Interceptors
// =============================================================================

/// Injects the JWT access token and handles transparent refresh on 401.
class _AuthInterceptor extends Interceptor {
  _AuthInterceptor({
    required AuthManager authManager,
    required Dio dio,
  })  : _authManager = authManager,
        _dio = dio;

  final AuthManager _authManager;
  final Dio _dio;
  bool _isRefreshing = false;

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final token = await _authManager.accessToken;
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    final deviceId = await _authManager.deviceId;
    if (deviceId != null) {
      options.headers['X-Device-ID'] = deviceId;
    }

    options.headers['X-Request-ID'] = const Uuid().v4();

    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    if (err.response?.statusCode != 401 || _isRefreshing) {
      return handler.next(err);
    }

    _isRefreshing = true;
    try {
      final refreshToken = await _authManager.refreshToken;
      if (refreshToken == null) {
        await _authManager.clearSession();
        _isRefreshing = false;
        return handler.next(err);
      }

      final response = await _dio.post<Map<String, dynamic>>(
        ApiEndpoints.refreshToken,
        data: {'refresh_token': refreshToken},
      );

      final data = response.data;
      if (data != null) {
        await _authManager.saveTokens(
          accessToken: data['access_token'] as String,
          refreshToken: data['refresh_token'] as String,
        );

        // Retry the original request with the new token.
        final retryOptions = err.requestOptions;
        retryOptions.headers['Authorization'] =
            'Bearer ${data['access_token']}';

        final retryResponse = await _dio.fetch<dynamic>(retryOptions);
        _isRefreshing = false;
        return handler.resolve(retryResponse);
      }
    } catch (_) {
      await _authManager.clearSession();
    }

    _isRefreshing = false;
    handler.next(err);
  }
}

/// Injects a unique Idempotency-Key for all POST and PUT requests.
class _IdempotencyInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final method = options.method.toUpperCase();
    if (method == 'POST' || method == 'PUT') {
      options.headers['Idempotency-Key'] ??= const Uuid().v4();
    }
    handler.next(options);
  }
}

/// Debug-only request/response logger.
class _LoggingInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    debugPrint('--> ${options.method} ${options.uri}');
    handler.next(options);
  }

  @override
  void onResponse(
    Response<dynamic> response,
    ResponseInterceptorHandler handler,
  ) {
    debugPrint(
      '<-- ${response.statusCode} ${response.requestOptions.uri}',
    );
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    debugPrint(
      '<-- ERROR ${err.response?.statusCode} ${err.requestOptions.uri}',
    );
    handler.next(err);
  }
}
