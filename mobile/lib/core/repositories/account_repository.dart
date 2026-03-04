import '../network/api_client.dart';
import '../network/api_endpoints.dart';

class AccountRepository {
  AccountRepository({required ApiClient apiClient}) : _api = apiClient;

  final ApiClient _api;

  Future<Map<String, dynamic>> getProfile() async {
    final response =
        await _api.get<Map<String, dynamic>>(ApiEndpoints.userProfile);
    return response.data!;
  }

  Future<Map<String, dynamic>> updateProfile(
      Map<String, dynamic> data) async {
    final response = await _api.patch<Map<String, dynamic>>(
      ApiEndpoints.userProfile,
      data: data,
    );
    return response.data!;
  }

  Future<List<dynamic>> getAccounts() async {
    final response =
        await _api.get<Map<String, dynamic>>(ApiEndpoints.accounts);
    final data = response.data!;
    return data['data'] as List<dynamic>? ?? [];
  }

  Future<List<dynamic>> getTransactions(
    String accountId, {
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.transactions(accountId),
      queryParameters: {'limit': limit, 'offset': offset},
    );
    final data = response.data!;
    return data['data'] as List<dynamic>? ?? [];
  }
}
