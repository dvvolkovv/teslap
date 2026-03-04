class CryptoWallet {
  final String id;
  final String address;
  final String network;
  final String status;
  final List<CryptoBalance> balances;

  CryptoWallet({
    required this.id,
    required this.address,
    this.network = 'fuse',
    this.status = 'active',
    this.balances = const [],
  });

  factory CryptoWallet.fromJson(Map<String, dynamic> json) {
    return CryptoWallet(
      id: json['id'] as String,
      address: json['address'] as String,
      network: json['network'] as String? ?? 'fuse',
      status: json['status'] as String? ?? 'active',
      balances: (json['balances'] as List<dynamic>?)
              ?.map((e) => CryptoBalance.fromJson(e as Map<String, dynamic>))
              .toList() ??
          [],
    );
  }
}

class CryptoBalance {
  final String tokenSymbol;
  final String tokenName;
  final String balance;
  final String? valueEur;

  CryptoBalance({
    required this.tokenSymbol,
    required this.tokenName,
    required this.balance,
    this.valueEur,
  });

  factory CryptoBalance.fromJson(Map<String, dynamic> json) {
    return CryptoBalance(
      tokenSymbol: json['token_symbol'] as String,
      tokenName: json['token_name'] as String,
      balance: json['balance'] as String? ?? '0',
      valueEur: json['value_eur'] as String?,
    );
  }
}

class CryptoPrice {
  final String symbol;
  final String name;
  final String priceEur;
  final String priceUsd;
  final String change24h;

  CryptoPrice({
    required this.symbol,
    required this.name,
    required this.priceEur,
    required this.priceUsd,
    required this.change24h,
  });

  factory CryptoPrice.fromJson(Map<String, dynamic> json) {
    return CryptoPrice(
      symbol: json['symbol'] as String,
      name: json['name'] as String,
      priceEur: json['price_eur'] as String? ?? '0',
      priceUsd: json['price_usd'] as String? ?? '0',
      change24h: json['change_24h'] as String? ?? '0',
    );
  }
}

class CryptoQuote {
  final String id;
  final String action;
  final String tokenSymbol;
  final String fiatAmount;
  final String cryptoAmount;
  final String rate;
  final String feeAmount;
  final String feePct;
  final String expiresAt;

  CryptoQuote({
    required this.id,
    required this.action,
    required this.tokenSymbol,
    required this.fiatAmount,
    required this.cryptoAmount,
    required this.rate,
    required this.feeAmount,
    required this.feePct,
    required this.expiresAt,
  });

  factory CryptoQuote.fromJson(Map<String, dynamic> json) {
    return CryptoQuote(
      id: json['id'] as String,
      action: json['action'] as String,
      tokenSymbol: json['token_symbol'] as String,
      fiatAmount: json['fiat_amount'] as String? ?? '0',
      cryptoAmount: json['crypto_amount'] as String? ?? '0',
      rate: json['rate'] as String? ?? '0',
      feeAmount: json['fee_amount'] as String? ?? '0',
      feePct: json['fee_pct'] as String? ?? '0',
      expiresAt: json['expires_at'] as String,
    );
  }
}

class CryptoTransaction {
  final String id;
  final String type;
  final String tokenSymbol;
  final String amount;
  final String? fiatAmount;
  final String status;
  final String? recipientAddress;
  final String createdAt;

  CryptoTransaction({
    required this.id,
    required this.type,
    required this.tokenSymbol,
    required this.amount,
    this.fiatAmount,
    required this.status,
    this.recipientAddress,
    required this.createdAt,
  });

  factory CryptoTransaction.fromJson(Map<String, dynamic> json) {
    return CryptoTransaction(
      id: json['id'] as String,
      type: json['type'] as String,
      tokenSymbol: json['token_symbol'] as String,
      amount: json['amount'] as String? ?? '0',
      fiatAmount: json['fiat_amount'] as String?,
      status: json['status'] as String,
      recipientAddress: json['recipient_address'] as String?,
      createdAt: json['created_at'] as String,
    );
  }
}
