import 'package:freezed_annotation/freezed_annotation.dart';

part 'account.freezed.dart';
part 'account.g.dart';

@freezed
class Account with _$Account {
  const factory Account({
    required String id,
    @JsonKey(name: 'account_number') String? accountNumber,
    @Default('active') String status,
    @JsonKey(name: 'sub_accounts')
    @Default([])
    List<SubAccount> subAccounts,
    @JsonKey(name: 'total_balance_eur') String? totalBalanceEur,
  }) = _Account;

  factory Account.fromJson(Map<String, dynamic> json) =>
      _$AccountFromJson(json);
}

@freezed
class SubAccount with _$SubAccount {
  const factory SubAccount({
    required String id,
    required String currency,
    String? iban,
    String? bic,
    required AccountBalance balance,
    @JsonKey(name: 'is_default') @Default(false) bool isDefault,
  }) = _SubAccount;

  factory SubAccount.fromJson(Map<String, dynamic> json) =>
      _$SubAccountFromJson(json);
}

@freezed
class AccountBalance with _$AccountBalance {
  const factory AccountBalance({
    @Default('0.00') String available,
    @Default('0.00') String pending,
    @Default('0.00') String total,
  }) = _AccountBalance;

  factory AccountBalance.fromJson(Map<String, dynamic> json) =>
      _$AccountBalanceFromJson(json);
}
