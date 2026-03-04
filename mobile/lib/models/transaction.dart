import 'package:freezed_annotation/freezed_annotation.dart';

part 'transaction.freezed.dart';
part 'transaction.g.dart';

@freezed
class Transaction with _$Transaction {
  const factory Transaction({
    required String id,
    required String type,
    required String direction,
    required String status,
    required String amount,
    required String currency,
    TransactionCounterparty? counterparty,
    TransactionMerchant? merchant,
    String? reference,
    String? fee,
    @JsonKey(name: 'fx_rate') String? fxRate,
    String? category,
    @JsonKey(name: 'created_at') required String createdAt,
    @JsonKey(name: 'settled_at') String? settledAt,
  }) = _Transaction;

  factory Transaction.fromJson(Map<String, dynamic> json) =>
      _$TransactionFromJson(json);
}

@freezed
class TransactionCounterparty with _$TransactionCounterparty {
  const factory TransactionCounterparty({
    String? name,
    String? iban,
    String? username,
  }) = _TransactionCounterparty;

  factory TransactionCounterparty.fromJson(Map<String, dynamic> json) =>
      _$TransactionCounterpartyFromJson(json);
}

@freezed
class TransactionMerchant with _$TransactionMerchant {
  const factory TransactionMerchant({
    String? name,
    String? category,
    String? mcc,
    String? country,
  }) = _TransactionMerchant;

  factory TransactionMerchant.fromJson(Map<String, dynamic> json) =>
      _$TransactionMerchantFromJson(json);
}

@freezed
class PaginatedTransactions with _$PaginatedTransactions {
  const factory PaginatedTransactions({
    required List<Transaction> data,
    required Pagination pagination,
  }) = _PaginatedTransactions;

  factory PaginatedTransactions.fromJson(Map<String, dynamic> json) =>
      _$PaginatedTransactionsFromJson(json);
}

@freezed
class Pagination with _$Pagination {
  const factory Pagination({
    @JsonKey(name: 'has_more') @Default(false) bool hasMore,
    @JsonKey(name: 'next_cursor') String? nextCursor,
    @JsonKey(name: 'total_count') int? totalCount,
  }) = _Pagination;

  factory Pagination.fromJson(Map<String, dynamic> json) =>
      _$PaginationFromJson(json);
}
