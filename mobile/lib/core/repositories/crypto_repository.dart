import '../network/api_client.dart';
import '../network/api_endpoints.dart';
import '../../models/crypto.dart';

class CryptoRepository {
  CryptoRepository({required ApiClient apiClient}) : _api = apiClient;

  final ApiClient _api;

  Future<CryptoWallet> getWallet() async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.cryptoWallet,
    );
    return CryptoWallet.fromJson(response.data!);
  }

  Future<List<CryptoPrice>> getPrices() async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.cryptoPrices,
    );
    final data = response.data!;
    final prices = data['data'] as List<dynamic>? ?? [];
    return prices
        .map((e) => CryptoPrice.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<CryptoQuote> getQuote({
    required String action,
    required String symbol,
    required String amount,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.cryptoQuote,
      queryParameters: {
        'action': action,
        'symbol': symbol,
        'amount': amount,
      },
    );
    return CryptoQuote.fromJson(response.data!);
  }

  Future<Map<String, dynamic>> buyCrypto(String quoteId) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.cryptoBuy,
      data: {'quote_id': quoteId},
    );
    return response.data!;
  }

  Future<Map<String, dynamic>> sellCrypto(String quoteId) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.cryptoSell,
      data: {'quote_id': quoteId},
    );
    return response.data!;
  }

  Future<Map<String, dynamic>> sendCrypto({
    required String symbol,
    required String amount,
    required String recipientAddress,
  }) async {
    final response = await _api.post<Map<String, dynamic>>(
      ApiEndpoints.cryptoSend,
      data: {
        'symbol': symbol,
        'amount': amount,
        'recipient_address': recipientAddress,
      },
    );
    return response.data!;
  }

  Future<List<CryptoTransaction>> getTransactions({
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.cryptoTransactions,
      queryParameters: {'limit': limit, 'offset': offset},
    );
    final data = response.data!;
    final txs = data['data'] as List<dynamic>? ?? [];
    return txs
        .map((e) => CryptoTransaction.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}
