import '../network/api_client.dart';
import '../network/api_endpoints.dart';
import '../../models/notification.dart';

class NotificationRepository {
  NotificationRepository({required ApiClient apiClient}) : _api = apiClient;

  final ApiClient _api;

  Future<List<AppNotification>> getNotifications({
    int limit = 20,
    int offset = 0,
  }) async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.notifications,
      queryParameters: {'limit': limit, 'offset': offset},
    );
    final data = response.data!;
    final items = data['data'] as List<dynamic>? ?? [];
    return items
        .map((e) => AppNotification.fromJson(e as Map<String, dynamic>))
        .toList();
  }

  Future<void> markAsRead(String notificationId) async {
    await _api.post<void>(ApiEndpoints.notificationRead(notificationId));
  }

  Future<NotificationPreferences> getPreferences() async {
    final response = await _api.get<Map<String, dynamic>>(
      ApiEndpoints.notificationPreferences,
    );
    return NotificationPreferences.fromJson(response.data!);
  }

  Future<NotificationPreferences> updatePreferences(
      NotificationPreferences prefs) async {
    final response = await _api.put<Map<String, dynamic>>(
      ApiEndpoints.notificationPreferences,
      data: prefs.toJson(),
    );
    return NotificationPreferences.fromJson(response.data!);
  }
}
