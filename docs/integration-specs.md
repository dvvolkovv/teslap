# TeslaPay Third-Party Integration Specifications

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Principal Software Architect, Dream Team

---

## 1. Integration Architecture Overview

```
+------------------+     +------------------+     +------------------+
|  Mobile App      |     |  TeslaPay        |     |  External        |
|  (Flutter)       |     |  Backend         |     |  Systems         |
+--------+---------+     +--------+---------+     +--------+---------+
         |                         |                        |
         |  Sumsub SDK            |  Sumsub REST API       |
         +---------- ----------->|<----------------------->| Sumsub Cloud
         |  (in-app verification) |  (webhooks, admin)      |
         |                        |                         |
         |  FuseBox Dart SDK      |  FuseBox TS SDK         |
         +---------------------->|<----------------------->| Fuse Network
         |  (wallet, sign tx)     |  (bundler, paymaster)   | (blockchain)
         |                        |                         |
         |                        |  Enfuce REST API        |
         |                        |<----------------------->| Mastercard
         |                        |  (card ops, webhooks)   | Network
         |                        |                         |
         |                        |  Banking Circle API     |
         |                        |<----------------------->| SEPA Network
         |                        |  (SCT, SCT Inst, SDD)   | (EBA STEP2)
         |                        |                         |
         |  Apple Pay SDK         |  MDES API               |
         +---------------------->|<----------------------->| Mastercard
         |  Google Pay SDK        |  (tokenization)         | MDES
         |                        |                         |
         |                        |  APNs / FCM             |
         |                        +----------------------->| Apple/Google
         |                        |  (push notifications)   | Push Services
         +------------------------+-------------------------+
```

---

## 2. Sumsub Integration

### 2.1 Overview

Sumsub provides identity verification (KYC), Anti-Money Laundering (AML) screening, and ongoing monitoring. Integration uses both mobile SDK (in-app verification UI) and server-side REST API (webhooks, administrative operations).

### 2.2 Architecture

```
Mobile App (Flutter)                TeslaPay Backend               Sumsub Cloud
      |                                  |                              |
      |  1. POST /kyc/verify             |                              |
      +--------------------------------->|                              |
      |                                  |  2. Create applicant          |
      |                                  +----------------------------->|
      |                                  |  3. Generate SDK token        |
      |                                  +----------------------------->|
      |  4. Return SDK token             |<-----------------------------+
      |<---------------------------------+                              |
      |                                  |                              |
      |  5. Initialize Sumsub SDK        |                              |
      |  6. Document capture             |                              |
      |  7. Liveness check               |                              |
      |  8. Upload docs + selfie -------->------------------------------>|
      |                                  |                              |
      |  9. SDK callback (completed)     |                              |
      |<---------------------------------|                              |
      |                                  |                              |
      |                                  | 10. Webhook: applicantReviewed
      |                                  |<-----------------------------+
      |                                  |                              |
      |                                  | 11. Process result            |
      |                                  |   - GREEN: activate account   |
      |                                  |   - RED/RETRY: notify user    |
      |                                  |   - RED/FINAL: block          |
      |                                  |                              |
      | 12. Push notification            |                              |
      |<---------------------------------+                              |
```

### 2.3 Sumsub Configuration

**Verification Levels (mapped to Sumsub flows):**

| TeslaPay Level | Sumsub Flow Name | Documents Required | Triggers |
|----------------|------------------|--------------------|----------|
| Basic (Tier 1) | `teslapay-basic` | ID document + liveness | Account registration |
| Enhanced (Tier 2) | `teslapay-enhanced` | ID + proof of address | Tier upgrade, cumulative > EUR 15K |
| Full (Tier 3) | `teslapay-full` | ID + POA + source of funds | Premium tier, high-value transactions |

**Supported Document Types:**
- Passport (all countries)
- National ID card (front + back)
- Driver's license (front + back)
- Residence permit (front + back)

### 2.4 Server-Side API Integration

```
Base URL: https://api.sumsub.com

Authentication:
  - App Token + Secret Key (HMAC-SHA256 signed requests)
  - Signature: HMAC-SHA256(ts + method + path + body, secret)
  - Headers:
      X-App-Token: <app_token>
      X-App-Access-Ts: <unix_timestamp>
      X-App-Access-Sig: <hmac_signature>
```

**Key API Endpoints Used:**

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/resources/applicants` | POST | Create applicant |
| `/resources/accessTokens?userId={id}&levelName={level}` | POST | Generate SDK token |
| `/resources/applicants/{id}/status` | GET | Check verification status |
| `/resources/applicants/{id}/one` | GET | Get applicant data |
| `/resources/applicants/{id}/requiredIdDocs/status` | GET | Get document check status |
| `/resources/checks/latest?type=FACE_COMPARE` | GET | Get liveness result |

### 2.5 Webhook Processing

**Webhook Events Consumed:**

| Event Type | Action |
|------------|--------|
| `applicantReviewed` | Process verification decision (approve/reject user) |
| `applicantOnHold` | Queue for manual review |
| `applicantActionOnHold` | Additional verification needed |
| `applicantReset` | Verification reset, re-initiate |
| `applicantPending` | Verification in progress (informational) |
| `applicantPersonalInfoChanged` | Update user profile if applicable |

**Webhook Verification:**
```go
// Verify webhook signature
func verifyWebhook(body []byte, receivedDigest string, webhookSecret string) bool {
    mac := hmac.New(sha256.New, []byte(webhookSecret))
    mac.Write(body)
    expectedDigest := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(expectedDigest), []byte(receivedDigest))
}
```

**Decision Mapping:**

| Sumsub Result | TeslaPay Action |
|---------------|-----------------|
| `reviewAnswer: GREEN` | Set `kyc_status = verified`, activate account, create Fuse wallet, send welcome notification |
| `reviewAnswer: RED`, `rejectType: RETRY` | Set `kyc_status = rejected`, notify user to retry (max 3 attempts) |
| `reviewAnswer: RED`, `rejectType: FINAL` | Set `kyc_status = rejected`, block account, escalate to compliance |

### 2.6 Ongoing AML Monitoring

- Enabled for all verified users via Sumsub Ongoing Monitoring API
- Sumsub continuously screens users against sanctions, PEP, and adverse media databases
- New matches trigger `applicantReviewed` webhook with updated screening results
- TeslaPay KYC Service processes these alerts and routes them to the compliance review queue
- Alert SLA: reviewed within 24 hours

### 2.7 Flutter SDK Integration

```dart
// pubspec.yaml
dependencies:
  sumsub_kyc: ^1.x.x

// Initialization
final onTokenExpired = () async {
  // Fetch new token from TeslaPay backend
  final response = await api.refreshSumsubToken();
  return response.accessToken;
};

final result = await SumsubKyc.launch(
  accessToken: sumsubAccessToken,
  onTokenExpired: onTokenExpired,
  locale: 'en',
  theme: SumsubTheme(
    primaryColor: Color(0xFF1A73E8),  // TeslaPay brand color
  ),
);

// Handle result
if (result.status == SumsubStatus.completed) {
  // Verification submitted; wait for webhook
} else if (result.status == SumsubStatus.cancelled) {
  // User cancelled; show prompt to continue later
}
```

### 2.8 Error Handling and Resilience

| Scenario | Handling |
|----------|----------|
| Sumsub API timeout | Retry with exponential backoff (1s, 2s, 4s), max 3 retries |
| Sumsub SDK crash | Catch exception, log, show user-friendly error with retry option |
| Webhook delivery failure | Sumsub retries automatically; TeslaPay also polls applicant status every 5 min for pending verifications |
| Sumsub outage | Queue new verification requests, process when service recovers; existing users unaffected |
| Duplicate webhook | Idempotent processing based on `applicantId` + `inspectionId` |

---

## 3. Fuse.io Blockchain Integration

### 3.1 Overview

Fuse.io is an EVM-compatible blockchain (Chain ID: 122) with a Wallet-as-a-Service platform called FuseBox. TeslaPay uses FuseBox for Smart Wallet creation and management, leveraging ERC-4337 Account Abstraction for gasless transactions.

### 3.2 Architecture

```
+-------------------+         +----------------------+        +------------------+
| Flutter App       |         | Crypto Service       |        | Fuse Network     |
| (fuse_wallet_sdk) |         | (NestJS + FuseBox TS)|        |                  |
+--------+----------+         +----------+-----------+        +--------+---------+
         |                               |                             |
         | Local EOA key generation      |                             |
         | (device Keystore/Keychain)    |                             |
         |                               |                             |
         | Sign UserOperations           | Submit UserOps to Bundler   |
         +------------------------------>+----------------------------->|
         |                               |                             |
         |                               | Paymaster sponsors gas      |
         |                               +----------------------------->|
         |                               |                             |
         |                               | Query balances via RPC      |
         |                               +----------------------------->|
         |                               |                             |
         |                               | Smart Wallets API           |
         |                               | (create wallet, get history)|
         |                               +----------------------------->| FuseBox Backend
         |                               |                             | (NestJS)
         |                               | Notifications API           |
         |                               | (WebSocket/webhook)         |
         |                               |<----------------------------+
         |                               |                             |
```

### 3.3 FuseBox SDK Configuration

**Backend (TypeScript):**
```typescript
import { FuseSDK } from "@fuseio/fusebox-web-sdk";

// Initialize SDK with API key
const fuseSDK = await FuseSDK.init(
  process.env.FUSE_API_KEY,
  {
    withPaymaster: true,  // Enable gasless transactions
  }
);

// Create Smart Wallet for user
const smartWallet = await fuseSDK.createSmartWallet(
  userEOAAddress,  // User's EOA from mobile device
);

// Execute gasless token transfer
const transferOp = await fuseSDK.transferToken(
  tokenAddress,     // ERC-20 contract address
  recipientAddress,
  amount,
  { withPaymaster: true }
);
```

**Mobile (Dart):**
```dart
// pubspec.yaml
dependencies:
  fuse_wallet_sdk: ^2.x.x

// Initialize
final fuseSDK = await FuseSDK.init(
  apiKey: environment.fuseApiKey,
  withPaymaster: true,
);

// Create smart wallet
final credentials = EthPrivateKey.fromHex(privateKeyHex);
await fuseSDK.authenticate(credentials);
final walletAddress = fuseSDK.wallet.getSender();
```

### 3.4 Smart Wallet Lifecycle

| Step | Action | Component |
|------|--------|-----------|
| 1 | User completes KYC | KYC Service emits `kyc.approved` event |
| 2 | Generate EOA keypair | Mobile app generates keypair in secure enclave (Keychain/Keystore) |
| 3 | Send EOA public address to backend | Mobile -> Crypto Service |
| 4 | Create Smart Wallet via FuseBox | Crypto Service calls FuseBox Smart Wallets API |
| 5 | Store wallet metadata | Crypto Service stores `(user_id, smart_wallet_address, eoa_address)` |
| 6 | Return wallet info to mobile | Wallet address displayed in app |

**Key Security Property:** The EOA private key never leaves the user's device. The Smart Wallet is a contract wallet that recognizes the EOA as its owner. TeslaPay backend never has custody of user keys.

### 3.5 Gasless Transactions (ERC-4337)

```
User Intent: Send 50 USDC to 0x1234...

1. Mobile App: Create UserOperation
   {
     sender: smartWalletAddress,
     callData: encode(transfer(to, amount)),
     paymasterAndData: paymasterAddress + paymasterSignature
   }

2. Mobile App: Sign UserOperation with EOA key

3. Crypto Service: Submit to FuseBox Bundler
   - Bundler validates UserOperation
   - Paymaster validates sponsorship (fee deducted from USDC balance)
   - Bundler submits to Fuse mempool

4. Fuse Network: Execute
   - Entry Point contract processes UserOperation
   - Smart Wallet executes the USDC transfer
   - Paymaster deducts gas fee equivalent in USDC

5. Crypto Service: Monitor transaction
   - Listen for transaction receipt
   - Update blockchain_transactions table
   - Emit crypto.events -> transfer.completed
```

### 3.6 On-Ramp / Off-Ramp (Buy/Sell Crypto)

**Buy Flow (EUR -> USDC):**
```
1. User requests buy quote
2. Crypto Service fetches current USDC/EUR price from price aggregator
3. Quote presented to user (valid 30 seconds)
4. User confirms (SCA required)
5. Crypto Service -> Ledger Service: Debit user EUR account, credit crypto settlement
6. Crypto Service -> FuseBox: Mint/transfer USDC from TeslaPay treasury wallet to user's Smart Wallet
7. Confirm transaction on-chain
8. Update order status to completed
9. Emit events for notification and audit
```

**Sell Flow (USDC -> EUR):**
```
1. User requests sell quote
2. Crypto Service fetches current USDC/EUR price
3. Quote presented, user confirms (SCA required)
4. Crypto Service -> FuseBox: Transfer USDC from user's Smart Wallet to TeslaPay treasury
5. Confirm transaction on-chain
6. Crypto Service -> Ledger Service: Debit crypto settlement, credit user EUR account
7. Update order status to completed
```

**Treasury Wallet:** TeslaPay maintains a treasury Smart Wallet pre-funded with FUSE, USDC, and USDT to facilitate instant buy/sell operations. Treasury liquidity is monitored with alerts when reserves drop below thresholds.

### 3.7 Blockchain Event Monitoring

The Crypto Service runs a background worker that:
1. Subscribes to Fuse RPC WebSocket for new blocks
2. Scans blocks for transactions involving TeslaPay user wallets
3. Detects incoming token transfers (deposits)
4. Updates `blockchain_transactions` table
5. Emits notification events for user

```typescript
// Block scanner pseudocode
const provider = new ethers.WebSocketProvider(FUSE_WS_RPC_URL);

provider.on("block", async (blockNumber) => {
  const block = await provider.getBlockWithTransactions(blockNumber);
  for (const tx of block.transactions) {
    if (isMonitoredAddress(tx.to) || isMonitoredAddress(tx.from)) {
      await processTransaction(tx);
    }
  }
});
```

### 3.8 Circuit Breaker for Fuse Network

| Parameter | Value |
|-----------|-------|
| Failure threshold | 5 failures in 30 seconds |
| Circuit open duration | 60 seconds |
| Half-open test requests | 1 per 15 seconds |
| Close threshold | 3 consecutive successes |
| Fallback behavior | Return cached balances, queue outgoing transactions |

When the circuit is open:
- Balance queries return last cached values with a "stale data" indicator
- Buy/sell operations are temporarily unavailable (user sees "Crypto services temporarily unavailable")
- Send/receive operations are queued for retry when circuit closes
- Fiat operations (SEPA, cards) are completely unaffected

---

## 4. Mastercard Card Integration (via Enfuce)

### 4.1 Overview

Enfuce is a Finnish EMI and Mastercard principal member providing BIN sponsorship, issuer processing, and card management APIs. Enfuce handles PCI DSS compliance for card data, Mastercard scheme certification, and physical card personalization/production.

### 4.2 Architecture

```
+------------------+      +------------------+      +------------------+
| TeslaPay         |      | Enfuce           |      | Mastercard       |
| Card Service     |      | Issuer Processor |      | Network          |
+--------+---------+      +--------+---------+      +--------+---------+
         |                          |                         |
         | REST API (card mgmt)     |                         |
         +------------------------->|                         |
         |                          |                         |
         |                          | Mastercard Messages     |
         |                          |<----------------------->|
         |                          |                         |
         | Webhook: authorization   |                         |
         |<-------------------------+ (real-time auth)        |
         |                          |                         |
         | Response: approve/decline|                         |
         +------------------------->|                         |
         |                          |                         |
         | Webhook: settlement      |                         |
         |<-------------------------+ (batch, T+1)            |
         |                          |                         |
         |                          | MDES (tokenization)     |
         |                          |<----------------------->| Apple Pay
         |                          |                         | Google Pay
         |                          |                         |
         |                          | Card Production         |
         |                          +----------------------->| Personalization
         |                          |                         | Bureau
```

### 4.3 Card Lifecycle

| Operation | Enfuce API | TeslaPay Action |
|-----------|-----------|-----------------|
| Create virtual card | `POST /cards` | Store card metadata, return last4 + expiry |
| Order physical card | `POST /cards/{id}/physical` | Track delivery, set status to `inactive` |
| Activate card | `PUT /cards/{id}/activate` | Update status to `active` |
| Freeze card | `PUT /cards/{id}/freeze` | Update status, decline new authorizations |
| Unfreeze card | `PUT /cards/{id}/unfreeze` | Restore to `active` |
| Block card (permanent) | `PUT /cards/{id}/block` | Update status, no recovery |
| Replace card | `POST /cards/{id}/replace` | New card issued, old card deactivated |
| Set PIN | `PUT /cards/{id}/pin` | PIN set at Enfuce, never stored at TeslaPay |
| Get PIN | `GET /cards/{id}/pin` | PIN retrieved from Enfuce, displayed briefly |
| Update spending limits | `PUT /cards/{id}/limits` | Limits enforced at Enfuce during authorization |

### 4.4 Authorization Flow (Real-Time)

```
Timeline (total < 100ms budget):

T+0ms    Merchant POS -> Mastercard -> Enfuce
T+10ms   Enfuce -> TeslaPay webhook: POST /internal/webhooks/enfuce/authorization
T+15ms   Card Service: Validate card status, check controls
T+20ms   Card Service -> Fraud Detection: Score (gRPC, <10ms)
T+30ms   Card Service -> Account Service: Check balance
T+40ms   Card Service -> Ledger Service: Post authorization hold
T+60ms   Card Service -> Enfuce: Approve response
T+70ms   Enfuce -> Mastercard: Approve
T+80ms   [Kafka: card.events -> authorization.approved]
T+90ms   Notification Service: Push to user device
T+3000ms Push notification arrives on user's phone
```

**Authorization Webhook Payload (from Enfuce):**
```json
{
  "event_type": "authorization_request",
  "card_id": "enfuce-card-id",
  "authorization": {
    "id": "auth-id",
    "amount": 4599,               // Minor units (cents)
    "currency": "EUR",
    "merchant": {
      "name": "Amazon EU S.a r.l.",
      "mcc": "5411",
      "country": "LU",
      "city": "Luxembourg"
    },
    "pos_entry_mode": "chip_contactless",
    "is_recurring": false,
    "three_ds": {
      "version": "2.2.0",
      "status": "authenticated"
    }
  },
  "timestamp": "2026-03-03T15:00:00.123Z"
}
```

**Authorization Response:**
```json
{
  "approved": true,
  // OR
  "approved": false,
  "decline_reason": "insufficient_funds"
}
```

### 4.5 Settlement Processing

Card settlements arrive daily (T+1) via batch webhook from Enfuce. The Card Service processes each settled authorization:
1. Match settlement to original authorization
2. Convert authorization hold to final posting in Ledger
3. Handle partial settlements and refunds
4. Update authorization status to `settled`

### 4.6 3D Secure 2.0

```
Online Purchase -> Merchant -> Mastercard 3DS Server -> Enfuce -> TeslaPay

1. Enfuce sends 3DS challenge request webhook
2. Card Service creates 3DS challenge record
3. Notification Service sends push to user: "Approve payment of EUR 45.99 at Amazon?"
4. User opens app, authenticates with biometric
5. Card Service responds to Enfuce: challenge_result = "authenticated"
6. Enfuce responds to Mastercard: authentication successful
7. Authorization proceeds

Timeout: 5 minutes. If no response, declined.
Fallback: SMS OTP if push not delivered within 30 seconds.
```

### 4.7 Apple Pay / Google Pay Tokenization

**Apple Pay Provisioning:**
1. User taps "Add to Apple Wallet" in TeslaPay app
2. Flutter calls Apple Pay SDK (via platform channel)
3. Apple Pay SDK generates certificates and nonce
4. Card Service sends provisioning request to Enfuce with Apple Pay data
5. Enfuce interacts with Mastercard MDES for tokenization
6. Enfuce returns encrypted pass data
7. Card Service passes data back to Apple Pay SDK
8. Card appears in Apple Wallet

**Google Pay Provisioning:**
Similar flow using Google Pay SDK and Mastercard MDES.

### 4.8 PCI DSS Scope Minimization

TeslaPay does not store, process, or transmit cardholder data (PAN, CVV, PIN). All sensitive card data resides at Enfuce. TeslaPay only handles:
- Tokenized card references (processor_card_id)
- Last 4 digits of PAN
- Card expiry date
- Cardholder name

This classifies TeslaPay as a **SAQ-A merchant equivalent** for PCI DSS purposes. Enfuce maintains PCI DSS Level 1 certification.

PIN viewing and card number display are handled by Enfuce's secure display API, which returns time-limited, encrypted data decrypted only on the user's device.

---

## 5. SEPA Payment Integration (via Banking Circle)

### 5.1 Overview

Banking Circle provides indirect SEPA scheme access via API, supporting SEPA Credit Transfer (SCT), SEPA Instant Credit Transfer (SCT Inst), and SEPA Direct Debit (SDD). Banking Circle is a licensed bank providing payment infrastructure to EMIs and payment institutions.

### 5.2 Architecture

```
+------------------+      +------------------+      +------------------+
| TeslaPay         |      | Banking Circle   |      | SEPA Scheme      |
| Payment Service  |      | (Licensed Bank)  |      | (EBA STEP2 /     |
+--------+---------+      +--------+---------+      | RT1 / TIPS)      |
         |                          |                +--------+---------+
         | REST API                 |                         |
         | Submit SCT               |                         |
         +------------------------->| Route via STEP2         |
         |                          +----------------------->|
         |                          |                         |
         | Webhook: status update   |                         |
         |<-------------------------+                         |
         |                          |                         |
         | Submit SCT Inst          |                         |
         +------------------------->| Route via RT1/TIPS      |
         |                          +----------------------->|
         |                          |                         |
         | Response: settled        |   (<10 seconds)         |
         |<-------------------------+<------------------------+
         |                          |                         |
         | Incoming payment webhook |                         |
         |<-------------------------+ Receive from STEP2      |
         |                          |<------------------------+
```

### 5.3 SEPA Payment Types

| Type | Method | SLA | Availability |
|------|--------|-----|--------------|
| SEPA Credit Transfer (SCT) | Batch via API | T+1 business day | Business days |
| SEPA Instant (SCT Inst) | Real-time API | < 10 seconds | 24/7/365 |
| SEPA Direct Debit (SDD) | Batch via API | D-2 submission | Business days |

### 5.4 API Integration

```
Base URL: https://api.bankingcircle.com/api/v1

Authentication:
  - OAuth2 client credentials flow
  - Certificate-based mTLS for API access
  - Separate credentials for production and sandbox

Key Endpoints:
  POST /payments/singles         -- Submit single SCT/SCT Inst
  POST /payments/batches         -- Submit batch payments
  GET  /payments/{id}            -- Get payment status
  GET  /payments/incoming        -- List incoming payments
  POST /directdebits             -- Submit SDD collection
  GET  /accounts/{id}/statement  -- Account statement
```

**Submit SEPA Instant:**
```json
POST /payments/singles
{
  "debtorAccount": {
    "account": "LT123456789012345678",
    "financialInstitution": "TESLLT21"
  },
  "creditorAccount": {
    "account": "DE89370400440532013000",
    "financialInstitution": "COBADEFFXXX"
  },
  "instructedAmount": {
    "currency": "EUR",
    "amount": 100.00
  },
  "paymentScheme": "SCTInst",
  "remittanceInformation": "Invoice #1234",
  "endToEndId": "TP2026030300001",
  "requestedExecutionDate": "2026-03-03"
}
```

### 5.5 Incoming Payment Processing

Banking Circle notifies TeslaPay of incoming SEPA payments via webhook:

```json
{
  "event_type": "incoming_payment",
  "payment": {
    "id": "bc-payment-id",
    "debtorName": "Jane Smith",
    "debtorAccount": "DE89370400440532013000",
    "creditorAccount": "LT123456789012345678",
    "amount": {
      "currency": "EUR",
      "amount": 500.00
    },
    "remittanceInformation": "Salary March 2026",
    "bookingDate": "2026-03-03",
    "valueDate": "2026-03-03"
  }
}
```

**Processing Steps:**
1. Validate creditor IBAN belongs to a TeslaPay user
2. Identify target sub-account
3. Post credit entry in Ledger (debit safeguarded funds, credit customer account)
4. Update account balance
5. Send push notification to user
6. Log in audit trail

### 5.6 IBAN and BIC

TeslaPay IBANs follow the Lithuanian format:
```
LT xx TESL XXXX XXXX XXXX
|  |  |    |
|  |  |    +-- Account number (12 digits)
|  |  +------- Bank code (TESL = TeslaPay)
|  +---------- Check digits (mod 97)
+------------- Country code (LT = Lithuania)

BIC: TESLLT21
```

Note: The actual bank code and BIC must be registered with the Bank of Lithuania and EBA. If TeslaPay uses Banking Circle's BIC, IBANs will reflect Banking Circle's scheme.

### 5.7 Reconciliation

Daily reconciliation between TeslaPay ledger and Banking Circle statements:
1. Download daily statement from Banking Circle API
2. Compare each transaction against TeslaPay payment orders and ledger entries
3. Flag discrepancies in reconciliation table
4. Auto-resolve minor timing differences (T+1 settlements appearing on different dates)
5. Escalate unresolved discrepancies to operations team

---

## 6. FX Rate Provider

### 6.1 Rate Sources

| Source | Purpose | Update Frequency |
|--------|---------|-----------------|
| ECB Reference Rates | Base rates for EUR pairs | Daily (14:15 CET) |
| Rate Aggregator (e.g., CurrencyLayer/XE) | Real-time mid-market rates | Every 60 seconds |
| TeslaPay markup engine | Apply tier-specific markup | Per request |

### 6.2 Rate Calculation

```
user_rate = mid_market_rate * (1 + tier_markup)

Example: EUR/USD mid-market = 1.0900, Standard tier markup = 0.30%
  buy_rate = 1.0900 * 1.003 = 1.0933 (user buys USD, pays more EUR)
  sell_rate = 1.0900 * 0.997 = 1.0867 (user sells USD, receives less EUR)
```

### 6.3 Rate Locking

When a user views an FX quote:
1. Rate locked for 30 seconds
2. Quote ID returned with expiry timestamp
3. If user confirms within 30 seconds, locked rate is used
4. If expired, user must request a new quote

---

## 7. Push Notification Integration

### 7.1 Apple Push Notification Service (APNs)

| Parameter | Value |
|-----------|-------|
| Authentication | Token-based (JWT with .p8 key) |
| Environment | Production: `api.push.apple.com`, Sandbox: `api.sandbox.push.apple.com` |
| Priority | 10 (immediate) for transaction alerts, 5 (throttled) for marketing |
| Expiry | 0 (immediate delivery only) for time-sensitive, 86400 for others |

### 7.2 Firebase Cloud Messaging (FCM)

| Parameter | Value |
|-----------|-------|
| Authentication | Service account JWT |
| API | HTTP v1 (`fcm.googleapis.com/v1/projects/{id}/messages:send`) |
| Priority | HIGH for transaction alerts, NORMAL for marketing |
| TTL | 0 for time-sensitive, 86400 for others |

### 7.3 Notification Types

| Type | Priority | Channel | Mandatory |
|------|----------|---------|-----------|
| Card authorization | HIGH | Push | Yes (cannot disable) |
| Payment received | HIGH | Push | Yes |
| Payment sent | HIGH | Push | Yes |
| 3DS challenge | HIGH | Push + SMS fallback | Yes |
| KYC status change | NORMAL | Push + Email | Yes |
| Security alert | HIGH | Push + Email + SMS | Yes |
| Marketing | NORMAL | Push/Email (opt-in) | No |

---

## 8. Integration Environment Strategy

### 8.1 Sandbox/Test Environments

| Integration | Sandbox URL | Test Credentials |
|-------------|-------------|-----------------|
| Sumsub | `https://test-api.sumsub.com` | Test app token + secret |
| Fuse.io | Fuse Spark testnet (Chain ID: 123) | Test API key |
| Enfuce | `https://sandbox.enfuce.com` | Test BIN range |
| Banking Circle | `https://sandbox.bankingcircle.com` | Test IBAN range |
| APNs | `api.sandbox.push.apple.com` | Development certificate |
| FCM | Same endpoint | Test device tokens |

### 8.2 Integration Testing Strategy

1. **Contract tests:** Verify webhook payload schemas against provider documentation
2. **Mock services:** WireMock stubs for each integration in CI pipeline
3. **Sandbox tests:** Nightly test suite runs against sandbox environments
4. **Certification tests:** Enfuce and Banking Circle require certification testing before production go-live
5. **Load tests:** Simulate 10x expected traffic against all integrations (in sandbox)

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CTO | TBD | | Pending |
| Integration Lead | TBD | | Pending |
| Principal Architect | Dream Team Architect | 2026-03-03 | Submitted |
