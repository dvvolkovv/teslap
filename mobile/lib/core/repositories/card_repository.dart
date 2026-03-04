import '../network/api_client.dart';
import '../network/api_endpoints.dart';

class CardRepository {
  CardRepository({required ApiClient apiClient}) : _api = apiClient;

  final ApiClient _api;

  Future<Map<String, dynamic>> issueVirtualCard({
    required String accountId,
    required String cardholderName,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.cardVirtual,
      data: {
        'account_id': accountId,
        'cardholder_name': cardholderName,
      },
    );
    return response.data!;
  }

  Future<Map<String, dynamic>> issuePhysicalCard({
    required String accountId,
    required String cardholderName,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.cardPhysical,
      data: {
        'account_id': accountId,
        'cardholder_name': cardholderName,
      },
    );
    return response.data!;
  }

  Future<List<dynamic>> listCards({String? accountId}) async {
    final response = await _api.get<Map<String, dynamic>>(
      '/cards',
      queryParameters: {
        if (accountId != null) 'account_id': accountId,
      },
    );
    final data = response.data!;
    return data['data'] as List<dynamic>? ?? [];
  }

  Future<Map<String, dynamic>> getCard(String cardId) async {
    final response = await _api.get<Map<String, dynamic>>(
      '/cards/$cardId',
    );
    return response.data!;
  }

  Future<void> freezeCard(String cardId) async {
    await _api.post<void>(ApiEndpoints.cardFreeze(cardId));
  }

  Future<void> unfreezeCard(String cardId) async {
    await _api.post<void>(ApiEndpoints.cardUnfreeze(cardId));
  }

  Future<void> blockCard(String cardId) async {
    await _api.post<void>('/cards/$cardId/block');
  }

  Future<void> activateCard(String cardId, String lastFour) async {
    await _api.post<void>(
      ApiEndpoints.cardActivate(cardId),
      data: {'last_four': lastFour},
    );
  }

  Future<Map<String, dynamic>> updateControls(
    String cardId,
    Map<String, dynamic> controls,
  ) async {
    final response = await _api.put<Map<String, dynamic>>(
      ApiEndpoints.cardControls(cardId),
      data: controls,
    );
    return response.data!;
  }

  Future<List<dynamic>> getCardTransactions(
    String cardId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.cardTransactions(cardId),
      queryParameters: {'limit': limit, 'offset': offset},
    );
    final data = response.data!;
    return data['data'] as List<dynamic>? ?? [];
  }
}
