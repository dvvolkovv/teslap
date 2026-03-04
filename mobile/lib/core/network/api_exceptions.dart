/// Typed API exceptions following the RFC 7807 error format from TeslaPay's
/// API contract.
class ApiException implements Exception {
  const ApiException({
    required this.statusCode,
    required this.errorCode,
    required this.title,
    required this.detail,
    this.traceId,
  });

  factory ApiException.fromJson(Map<String, dynamic> json) {
    return ApiException(
      statusCode: json['status'] as int? ?? 0,
      errorCode: json['error_code'] as String? ?? 'UNKNOWN',
      title: json['title'] as String? ?? 'Error',
      detail: json['detail'] as String? ?? 'An unexpected error occurred.',
      traceId: json['trace_id'] as String?,
    );
  }

  final int statusCode;
  final String errorCode;
  final String title;
  final String detail;
  final String? traceId;

  bool get isAuthError =>
      errorCode == 'AUTH_001' ||
      errorCode == 'AUTH_002' ||
      errorCode == 'AUTH_003';

  bool get isScaRequired => errorCode == 'SCA_001';
  bool get isInsufficientFunds => errorCode == 'PAY_001';
  bool get isLimitExceeded => errorCode == 'PAY_002';
  bool get isInvalidIban => errorCode == 'PAY_003';
  bool get isKycRequired => errorCode == 'KYC_001';
  bool get isCardFrozen => errorCode == 'CARD_001';
  bool get isRateLimited => errorCode == 'RATE_LIMITED';

  @override
  String toString() => 'ApiException($errorCode, $statusCode): $detail';
}

/// Thrown when a network request fails due to connectivity issues.
class NetworkException implements Exception {
  const NetworkException([this.message = 'No internet connection.']);

  final String message;

  @override
  String toString() => 'NetworkException: $message';
}

/// Thrown when the server response could not be parsed.
class ParseException implements Exception {
  const ParseException([this.message = 'Failed to parse server response.']);

  final String message;

  @override
  String toString() => 'ParseException: $message';
}

/// Thrown when the user's session has expired and they must re-authenticate.
class SessionExpiredException implements Exception {
  const SessionExpiredException();

  @override
  String toString() => 'SessionExpiredException: User must re-authenticate.';
}
