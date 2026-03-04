class AppNotification {
  final String id;
  final String type;
  final String channel;
  final String title;
  final String? body;
  final String status;
  final String? readAt;
  final String createdAt;

  AppNotification({
    required this.id,
    required this.type,
    required this.channel,
    required this.title,
    this.body,
    required this.status,
    this.readAt,
    required this.createdAt,
  });

  bool get isRead => readAt != null;

  factory AppNotification.fromJson(Map<String, dynamic> json) {
    return AppNotification(
      id: json['id'] as String,
      type: json['type'] as String,
      channel: json['channel'] as String,
      title: json['title'] as String,
      body: json['body'] as String?,
      status: json['status'] as String,
      readAt: json['read_at'] as String?,
      createdAt: json['created_at'] as String,
    );
  }
}

class NotificationPreferences {
  final bool pushEnabled;
  final bool emailEnabled;
  final bool smsEnabled;
  final bool paymentAlerts;
  final bool cardAlerts;
  final bool kycAlerts;
  final bool marketing;

  NotificationPreferences({
    this.pushEnabled = true,
    this.emailEnabled = true,
    this.smsEnabled = false,
    this.paymentAlerts = true,
    this.cardAlerts = true,
    this.kycAlerts = true,
    this.marketing = false,
  });

  factory NotificationPreferences.fromJson(Map<String, dynamic> json) {
    return NotificationPreferences(
      pushEnabled: json['push_enabled'] as bool? ?? true,
      emailEnabled: json['email_enabled'] as bool? ?? true,
      smsEnabled: json['sms_enabled'] as bool? ?? false,
      paymentAlerts: json['payment_alerts'] as bool? ?? true,
      cardAlerts: json['card_alerts'] as bool? ?? true,
      kycAlerts: json['kyc_alerts'] as bool? ?? true,
      marketing: json['marketing'] as bool? ?? false,
    );
  }

  Map<String, dynamic> toJson() => {
        'push_enabled': pushEnabled,
        'email_enabled': emailEnabled,
        'sms_enabled': smsEnabled,
        'payment_alerts': paymentAlerts,
        'card_alerts': cardAlerts,
        'kyc_alerts': kycAlerts,
        'marketing': marketing,
      };
}
