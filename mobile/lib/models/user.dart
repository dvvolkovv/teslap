import 'package:freezed_annotation/freezed_annotation.dart';

part 'user.freezed.dart';
part 'user.g.dart';

@freezed
class User with _$User {
  const factory User({
    required String id,
    @JsonKey(name: 'external_id') String? externalId,
    @JsonKey(name: 'first_name') required String firstName,
    @JsonKey(name: 'last_name') required String lastName,
    required String email,
    String? phone,
    @JsonKey(name: 'date_of_birth') String? dateOfBirth,
    String? nationality,
    UserAddress? address,
    @JsonKey(name: 'kyc_status') @Default('pending') String kycStatus,
    @JsonKey(name: 'kyc_level') @Default(0) int kycLevel,
    UserTier? tier,
    @Default('en') String language,
    @JsonKey(name: 'created_at') String? createdAt,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

@freezed
class UserAddress with _$UserAddress {
  const factory UserAddress({
    String? line1,
    String? city,
    @JsonKey(name: 'postal_code') String? postalCode,
    String? country,
  }) = _UserAddress;

  factory UserAddress.fromJson(Map<String, dynamic> json) =>
      _$UserAddressFromJson(json);
}

@freezed
class UserTier with _$UserTier {
  const factory UserTier({
    required String name,
    UserTierLimits? limits,
  }) = _UserTier;

  factory UserTier.fromJson(Map<String, dynamic> json) =>
      _$UserTierFromJson(json);
}

@freezed
class UserTierLimits with _$UserTierLimits {
  const factory UserTierLimits({
    @JsonKey(name: 'daily_transfer') String? dailyTransfer,
    @JsonKey(name: 'monthly_transfer') String? monthlyTransfer,
    @JsonKey(name: 'daily_card') String? dailyCard,
  }) = _UserTierLimits;

  factory UserTierLimits.fromJson(Map<String, dynamic> json) =>
      _$UserTierLimitsFromJson(json);
}
