import 'package:freezed_annotation/freezed_annotation.dart';

part 'card.freezed.dart';
part 'card.g.dart';

@freezed
class PaymentCard with _$PaymentCard {
  const factory PaymentCard({
    required String id,
    required String type,
    @Default('mastercard') String brand,
    @JsonKey(name: 'last_four') required String lastFour,
    required String expiry,
    @JsonKey(name: 'cardholder_name') required String cardholderName,
    required String status,
    @JsonKey(name: 'linked_currency') String? linkedCurrency,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'estimated_delivery') String? estimatedDelivery,
  }) = _PaymentCard;

  factory PaymentCard.fromJson(Map<String, dynamic> json) =>
      _$PaymentCardFromJson(json);
}

@freezed
class CardDetails with _$CardDetails {
  const factory CardDetails({
    @JsonKey(name: 'card_number') required String cardNumber,
    required String expiry,
    required String cvv,
    @JsonKey(name: 'display_timeout') @Default(10) int displayTimeout,
  }) = _CardDetails;

  factory CardDetails.fromJson(Map<String, dynamic> json) =>
      _$CardDetailsFromJson(json);
}

@freezed
class CardControls with _$CardControls {
  const factory CardControls({
    @JsonKey(name: 'per_transaction_limit') String? perTransactionLimit,
    @JsonKey(name: 'daily_limit') String? dailyLimit,
    @JsonKey(name: 'monthly_limit') String? monthlyLimit,
    @JsonKey(name: 'atm_daily_limit') String? atmDailyLimit,
    @JsonKey(name: 'online_enabled') @Default(true) bool onlineEnabled,
    @JsonKey(name: 'contactless_enabled') @Default(true) bool contactlessEnabled,
    @JsonKey(name: 'atm_enabled') @Default(true) bool atmEnabled,
    @JsonKey(name: 'magstripe_enabled') @Default(false) bool magstripeEnabled,
  }) = _CardControls;

  factory CardControls.fromJson(Map<String, dynamic> json) =>
      _$CardControlsFromJson(json);
}
