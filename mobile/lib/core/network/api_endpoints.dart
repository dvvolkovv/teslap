/// Centralised API endpoint constants derived from the TeslaPay API contract.
abstract final class ApiEndpoints {
  static const String baseUrl = 'https://api.teslapay.eu/api/v1';
  static const String stagingBaseUrl =
      'https://api.staging.teslapay.eu/api/v1';

  // ---------------------------------------------------------------------------
  // Auth
  // ---------------------------------------------------------------------------
  static const String register = '/auth/register';
  static const String verifyEmail = '/auth/verify-email';
  static const String verifyPhone = '/auth/verify-phone';
  static const String login = '/auth/login';
  static const String biometricLogin = '/auth/biometric';
  static const String refreshToken = '/auth/refresh';
  static const String logout = '/auth/logout';
  static const String sessions = '/auth/sessions';
  static const String scaInitiate = '/auth/sca/initiate';
  static const String scaVerify = '/auth/sca/verify';

  // ---------------------------------------------------------------------------
  // User / Account
  // ---------------------------------------------------------------------------
  static const String userProfile = '/users/me';
  static const String accounts = '/accounts';

  static String subAccounts(String accountId) =>
      '/accounts/$accountId/sub-accounts';

  static String transactions(String accountId) =>
      '/accounts/$accountId/transactions';

  static String exportTransactions(String accountId) =>
      '/accounts/$accountId/transactions/export';

  // ---------------------------------------------------------------------------
  // Beneficiaries
  // ---------------------------------------------------------------------------
  static const String beneficiaries = '/beneficiaries';

  // ---------------------------------------------------------------------------
  // Payments
  // ---------------------------------------------------------------------------
  static const String paymentSepa = '/payments/sepa';
  static const String paymentInternal = '/payments/internal';
  static const String fxQuote = '/payments/fx/quote';
  static const String fxExecute = '/payments/fx/execute';
  static const String scheduledPayments = '/payments/scheduled';

  static String paymentStatus(String paymentId) => '/payments/$paymentId';

  // ---------------------------------------------------------------------------
  // Cards
  // ---------------------------------------------------------------------------
  static const String cardVirtual = '/cards/virtual';
  static const String cardPhysical = '/cards/physical';

  static String cardDetails(String cardId) => '/cards/$cardId/details';
  static String cardFreeze(String cardId) => '/cards/$cardId/freeze';
  static String cardUnfreeze(String cardId) => '/cards/$cardId/unfreeze';
  static String cardControls(String cardId) => '/cards/$cardId/controls';
  static String cardPinView(String cardId) => '/cards/$cardId/pin/view';
  static String cardPinChange(String cardId) => '/cards/$cardId/pin/change';
  static String cardActivate(String cardId) => '/cards/$cardId/activate';
  static String cardReport(String cardId) => '/cards/$cardId/report';
  static String cardTokenize(String cardId) => '/cards/$cardId/tokenize';
  static String cardTransactions(String cardId) =>
      '/cards/$cardId/transactions';

  // ---------------------------------------------------------------------------
  // Crypto
  // ---------------------------------------------------------------------------
  static const String cryptoWallet = '/crypto/wallet';
  static const String cryptoDepositAddress = '/crypto/wallet/deposit-address';
  static const String cryptoQuote = '/crypto/quote';
  static const String cryptoBuy = '/crypto/buy';
  static const String cryptoSell = '/crypto/sell';
  static const String cryptoSend = '/crypto/send';
  static const String cryptoTransactions = '/crypto/transactions';
  static const String cryptoPrices = '/crypto/prices';

  // ---------------------------------------------------------------------------
  // KYC
  // ---------------------------------------------------------------------------
  static const String kycVerify = '/kyc/verify';
  static const String kycStatus = '/kyc/status';
  static const String kycUpgrade = '/kyc/upgrade';

  // ---------------------------------------------------------------------------
  // Notifications
  // ---------------------------------------------------------------------------
  static const String notificationPreferences = '/notifications/preferences';
  static const String notifications = '/notifications';

  static String notificationRead(String id) => '/notifications/$id/read';
}
