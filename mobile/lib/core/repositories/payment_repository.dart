import '../network/api_client.dart';
import '../network/api_endpoints.dart';

class PaymentRepository {
  PaymentRepository({required ApiClient apiClient}) : _api = apiClient;

  final ApiClient _api;

  Future<Map<String, dynamic>> createSepaPayment({
    required String senderAccountId,
    required String recipientIban,
    required String recipientName,
    required String amount,
    String currency = 'EUR',
    String? reference,
    String? description,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.paymentSepa,
      data: {
        'sender_account_id': senderAccountId,
        'recipient_iban': recipientIban,
        'recipient_name': recipientName,
        'amount': amount,
        'currency': currency,
        if (reference != null) 'reference': reference,
        if (description != null) 'description': description,
      },
    );
    return response.data!;
  }

  Future<Map<String, dynamic>> createInternalPayment({
    required String senderAccountId,
    required String recipientAccountId,
    required String amount,
    String currency = 'EUR',
    String? reference,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.paymentInternal,
      data: {
        'sender_account_id': senderAccountId,
        'recipient_account_id': recipientAccountId,
        'amount': amount,
        'currency': currency,
        if (reference != null) 'reference': reference,
      },
    );
    return response.data!;
  }

  Future<Map<String, dynamic>> getPayment(String paymentId) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.paymentStatus(paymentId),
    );
    return response.data!;
  }

  Future<List<dynamic>> listPayments({
    String? accountId,
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      '/payments',
      queryParameters: {
        if (accountId != null) 'account_id': accountId,
        'limit': limit,
        'offset': offset,
      },
    );
    final data = response.data!;
    return data['data'] as List<dynamic>? ?? [];
  }

  Future<Map<String, dynamic>> getFxQuote({
    required String from,
    required String to,
    required String amount,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.fxQuote,
      queryParameters: {'from': from, 'to': to, 'amount': amount},
    );
    return response.data!;
  }

  Future<Map<String, dynamic>> executeFx({
    required String quoteId,
    required String accountId,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.fxExecute,
      data: {'quote_id': quoteId, 'account_id': accountId},
    );
    return response.data!;
  }
}
