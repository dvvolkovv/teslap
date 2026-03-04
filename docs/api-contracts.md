# TeslaPay API Contracts

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Principal Software Architect, Dream Team

---

## 1. API Design Principles

| Principle | Implementation |
|-----------|---------------|
| RESTful | Resources as nouns, HTTP verbs for operations, proper status codes |
| Versioned | URL path versioning: `/api/v1/...` |
| Idempotent | All POST/PUT operations accept `Idempotency-Key` header |
| Paginated | Cursor-based pagination for all list endpoints |
| Consistent Errors | RFC 7807 Problem Details for all error responses |
| Rate Limited | Per-user limits enforced at API Gateway |
| Authenticated | Bearer JWT for all endpoints except auth |
| SCA Enforced | PSD2 Strong Customer Authentication for financial operations |

### 1.1 Base URL

```
Production:  https://api.teslapay.eu/api/v1
Staging:     https://api.staging.teslapay.eu/api/v1
```

### 1.2 Common Headers

```
Request:
  Authorization: Bearer <jwt_access_token>
  Idempotency-Key: <uuid>                   -- Required for POST/PUT
  X-Request-ID: <uuid>                      -- Client-generated correlation ID
  X-Device-ID: <uuid>                       -- Registered device identifier
  Accept-Language: en                        -- Content language preference
  Content-Type: application/json

Response:
  X-Request-ID: <uuid>                      -- Echoed back
  X-RateLimit-Limit: 100
  X-RateLimit-Remaining: 87
  X-RateLimit-Reset: 1709478400
  Content-Type: application/json
```

### 1.3 Error Response Format (RFC 7807)

```json
{
  "type": "https://api.teslapay.eu/errors/insufficient-funds",
  "title": "Insufficient Funds",
  "status": 422,
  "detail": "Account EUR balance is 45.00, but payment requires 100.00",
  "instance": "/api/v1/payments/sepa",
  "error_code": "PAY_001",
  "trace_id": "abc123-def456"
}
```

### 1.4 Standard Error Codes

| HTTP Status | Error Code | Meaning |
|-------------|-----------|---------|
| 400 | `VALIDATION_ERROR` | Request body validation failed |
| 401 | `AUTH_001` | Missing or invalid access token |
| 401 | `AUTH_002` | Access token expired |
| 403 | `AUTH_003` | Insufficient permissions |
| 403 | `SCA_001` | Strong Customer Authentication required |
| 404 | `NOT_FOUND` | Resource not found |
| 409 | `CONFLICT` | Idempotency key already used with different payload |
| 422 | `PAY_001` | Insufficient funds |
| 422 | `PAY_002` | Transfer limit exceeded |
| 422 | `PAY_003` | Invalid IBAN |
| 422 | `KYC_001` | KYC verification required |
| 422 | `CARD_001` | Card is frozen |
| 429 | `RATE_LIMITED` | Rate limit exceeded |
| 500 | `INTERNAL_ERROR` | Internal server error |
| 503 | `SERVICE_UNAVAILABLE` | Service temporarily unavailable |

### 1.5 Pagination

All list endpoints use cursor-based pagination:

```json
// Request
GET /api/v1/transactions?cursor=eyJpZCI6...&limit=20

// Response
{
  "data": [...],
  "pagination": {
    "has_more": true,
    "next_cursor": "eyJpZCI6...",
    "total_count": 1234
  }
}
```

### 1.6 Versioning Strategy

- URL path versioning (`/api/v1/`, `/api/v2/`)
- Breaking changes require new version
- Old version supported for minimum 12 months after deprecation notice
- Deprecation communicated via `Sunset` response header and developer portal
- Non-breaking additions (new fields, new optional parameters) do NOT require new version

### 1.7 Rate Limiting

| Tier | Requests/Minute | Burst | Notes |
|------|-----------------|-------|-------|
| Basic | 60 | 10 | Default for authenticated users |
| Standard | 100 | 20 | Standard tier accounts |
| Premium | 300 | 50 | Premium tier accounts |
| Unauthenticated | 20 | 5 | Registration/login only |
| Webhook Callbacks | Unlimited | -- | Internal webhook processing |

---

## 2. Auth API

### 2.1 Register

```
POST /api/v1/auth/register

Request:
{
  "email": "user@example.com",
  "phone": "+37060012345",
  "password": "SecureP@ss123",
  "language": "en",
  "device": {
    "device_id": "uuid",
    "device_name": "iPhone 15 Pro",
    "device_type": "ios",
    "push_token": "APNs-token-here"
  },
  "consent": {
    "terms_accepted": true,
    "privacy_accepted": true,
    "marketing_opt_in": false
  }
}

Response: 201 Created
{
  "user_id": "uuid",
  "email_verification_sent": true,
  "phone_verification_sent": true
}
```

### 2.2 Verify Email

```
POST /api/v1/auth/verify-email

Request:
{
  "user_id": "uuid",
  "code": "123456"
}

Response: 200 OK
{
  "email_verified": true
}
```

### 2.3 Verify Phone

```
POST /api/v1/auth/verify-phone

Request:
{
  "user_id": "uuid",
  "code": "123456"
}

Response: 200 OK
{
  "phone_verified": true
}
```

### 2.4 Login

```
POST /api/v1/auth/login

Request:
{
  "email": "user@example.com",
  "password": "SecureP@ss123",
  "device": {
    "device_id": "uuid",
    "device_name": "iPhone 15 Pro",
    "device_type": "ios"
  }
}

Response: 200 OK
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "refresh_token": "dGVzbGFwYXktcmVm...",
  "token_type": "Bearer",
  "expires_in": 900,
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "kyc_status": "verified",
    "tier": "standard"
  }
}

Response: 403 (new device detected)
{
  "type": "https://api.teslapay.eu/errors/new-device",
  "title": "New Device Detected",
  "status": 403,
  "detail": "Please verify this device via the OTP sent to your phone",
  "error_code": "AUTH_004",
  "challenge_id": "uuid"
}
```

### 2.5 Biometric Login

```
POST /api/v1/auth/biometric

Request:
{
  "device_id": "uuid",
  "biometric_signature": "base64-encoded-signature",
  "biometric_type": "face_id"
}

Response: 200 OK
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "refresh_token": "dGVzbGFwYXktcmVm...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### 2.6 Refresh Token

```
POST /api/v1/auth/refresh

Request:
{
  "refresh_token": "dGVzbGFwYXktcmVm..."
}

Response: 200 OK
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "refresh_token": "bmV3LXJlZnJlc2g...",
  "token_type": "Bearer",
  "expires_in": 900
}
```

### 2.7 Logout

```
POST /api/v1/auth/logout

Request:
{
  "all_sessions": false      // true = logout everywhere
}

Response: 204 No Content
```

### 2.8 List Active Sessions

```
GET /api/v1/auth/sessions

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "device_name": "iPhone 15 Pro",
      "device_type": "ios",
      "ip_address": "192.168.1.1",
      "location": "Vilnius, LT",
      "last_active_at": "2026-03-03T14:30:00Z",
      "is_current": true
    }
  ]
}
```

### 2.9 Terminate Session

```
DELETE /api/v1/auth/sessions/{session_id}

Response: 204 No Content
```

### 2.10 SCA Challenge (PSD2)

```
POST /api/v1/auth/sca/initiate

Request:
{
  "action": "payment",
  "action_id": "uuid",           // payment order ID
  "method": "biometric"          // or "push", "sms_otp"
}

Response: 200 OK
{
  "challenge_id": "uuid",
  "method": "biometric",
  "expires_at": "2026-03-03T15:01:00Z"
}
```

```
POST /api/v1/auth/sca/verify

Request:
{
  "challenge_id": "uuid",
  "biometric_signature": "base64-encoded",
  // OR
  "otp_code": "123456"
}

Response: 200 OK
{
  "sca_token": "sca-token-uuid",
  "expires_at": "2026-03-03T15:05:00Z"
}
```

---

## 3. Account API

### 3.1 Get User Profile

```
GET /api/v1/users/me

Response: 200 OK
{
  "id": "uuid",
  "external_id": "TP-12345678",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "phone": "+37060012345",
  "date_of_birth": "1990-05-15",
  "nationality": "LTU",
  "address": {
    "line1": "Gedimino pr. 1",
    "city": "Vilnius",
    "postal_code": "01103",
    "country": "LTU"
  },
  "kyc_status": "verified",
  "kyc_level": 2,
  "tier": {
    "name": "standard",
    "limits": {
      "daily_transfer": "10000.00",
      "monthly_transfer": "50000.00",
      "daily_card": "5000.00"
    }
  },
  "language": "en",
  "created_at": "2026-01-15T10:30:00Z"
}
```

### 3.2 Update Profile

```
PATCH /api/v1/users/me

Request:
{
  "address": {
    "line1": "Konstitucijos pr. 7",
    "city": "Vilnius",
    "postal_code": "09308",
    "country": "LTU"
  },
  "language": "lt"
}

Response: 200 OK
{
  // Updated user profile
  "re_verification_required": true,
  "re_verification_reason": "address_change"
}
```

### 3.3 List Accounts

```
GET /api/v1/accounts

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "account_number": "TP1234567890",
      "status": "active",
      "sub_accounts": [
        {
          "id": "uuid",
          "currency": "EUR",
          "iban": "LT123456789012345678",
          "bic": "TESLLT21",
          "balance": {
            "available": "1250.50",
            "pending": "45.00",
            "total": "1295.50"
          },
          "is_default": true
        },
        {
          "id": "uuid",
          "currency": "USD",
          "iban": null,
          "balance": {
            "available": "500.00",
            "pending": "0.00",
            "total": "500.00"
          },
          "is_default": false
        }
      ],
      "total_balance_eur": "1715.50"
    }
  ]
}
```

### 3.4 Create Sub-Account

```
POST /api/v1/accounts/{account_id}/sub-accounts

Request:
{
  "currency": "GBP"
}

Response: 201 Created
{
  "id": "uuid",
  "currency": "GBP",
  "iban": null,
  "balance": {
    "available": "0.00",
    "pending": "0.00",
    "total": "0.00"
  }
}
```

### 3.5 Close Sub-Account

```
POST /api/v1/accounts/{account_id}/sub-accounts/{sub_account_id}/close

Request:
{
  "transfer_remaining_to": "uuid"    // Target sub-account for remaining balance
}

Response: 200 OK
{
  "status": "closed",
  "remaining_transferred": "0.00"
}
```

### 3.6 Get Transaction History

```
GET /api/v1/accounts/{account_id}/transactions?currency=EUR&type=payment&from=2026-01-01&to=2026-03-01&limit=20

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "type": "sepa_sct",
      "direction": "outgoing",
      "status": "completed",
      "amount": "-500.00",
      "currency": "EUR",
      "counterparty": {
        "name": "Jane Smith",
        "iban": "DE89370400440532013000"
      },
      "reference": "Invoice #1234",
      "fee": "0.20",
      "fx_rate": null,
      "category": "transfer",
      "created_at": "2026-02-28T14:30:00Z",
      "settled_at": "2026-03-01T09:00:00Z"
    }
  ],
  "pagination": {
    "has_more": true,
    "next_cursor": "eyJpZCI6...",
    "total_count": 156
  }
}
```

### 3.7 Export Transactions

```
POST /api/v1/accounts/{account_id}/transactions/export

Request:
{
  "format": "csv",            // or "pdf"
  "from": "2026-01-01",
  "to": "2026-03-01",
  "currency": "EUR"
}

Response: 202 Accepted
{
  "export_id": "uuid",
  "status": "processing",
  "download_url": null,
  "estimated_ready_at": "2026-03-03T15:05:00Z"
}
```

```
GET /api/v1/accounts/{account_id}/transactions/export/{export_id}

Response: 200 OK
{
  "export_id": "uuid",
  "status": "completed",
  "download_url": "https://api.teslapay.eu/downloads/uuid?token=signed-url",
  "expires_at": "2026-03-04T15:00:00Z"
}
```

### 3.8 Manage Beneficiaries

```
GET /api/v1/beneficiaries
POST /api/v1/beneficiaries
PUT /api/v1/beneficiaries/{id}
DELETE /api/v1/beneficiaries/{id}

POST /api/v1/beneficiaries
Request:
{
  "name": "Jane Smith",
  "iban": "DE89370400440532013000",
  "bic": "COBADEFFXXX",
  "default_reference": "Monthly rent"
}

Response: 201 Created
{
  "id": "uuid",
  "name": "Jane Smith",
  "iban": "DE89370400440532013000",
  "bic": "COBADEFFXXX",
  "bank_name": "Commerzbank",
  "is_internal": false,
  "created_at": "2026-03-03T15:00:00Z"
}
```

---

## 4. Payment API

### 4.1 Initiate SEPA Credit Transfer

```
POST /api/v1/payments/sepa
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid     -- Required for amounts > 30 EUR

Request:
{
  "source_sub_account_id": "uuid",
  "destination": {
    "iban": "DE89370400440532013000",
    "name": "Jane Smith",
    "bic": "COBADEFFXXX"
  },
  "amount": "500.00",
  "currency": "EUR",
  "reference": "Invoice #1234",
  "instant": false,
  "save_beneficiary": true
}

Response: 201 Created
{
  "id": "uuid",
  "type": "sepa_sct",
  "status": "authorized",
  "amount": "500.00",
  "currency": "EUR",
  "fee": "0.20",
  "estimated_arrival": "2026-03-04T17:00:00Z",
  "end_to_end_id": "TP2026030300001",
  "created_at": "2026-03-03T15:00:00Z"
}
```

### 4.2 Initiate SEPA Instant Transfer

```
POST /api/v1/payments/sepa
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid

Request:
{
  "source_sub_account_id": "uuid",
  "destination": {
    "iban": "DE89370400440532013000",
    "name": "Jane Smith"
  },
  "amount": "100.00",
  "currency": "EUR",
  "reference": "Dinner split",
  "instant": true
}

Response: 201 Created
{
  "id": "uuid",
  "type": "sepa_sct_inst",
  "status": "completed",
  "amount": "100.00",
  "currency": "EUR",
  "fee": "0.50",
  "settled_at": "2026-03-03T15:00:08Z",
  "created_at": "2026-03-03T15:00:00Z"
}
```

### 4.3 Internal Transfer

```
POST /api/v1/payments/internal
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid    -- Required for amounts > 30 EUR

Request:
{
  "source_sub_account_id": "uuid",
  "destination": {
    "iban": "LT987654321098765432",
    // OR
    "phone": "+37060099999",
    // OR
    "username": "TP-87654321"
  },
  "amount": "25.00",
  "currency": "EUR",
  "note": "Coffee money"
}

Response: 201 Created
{
  "id": "uuid",
  "type": "internal",
  "status": "completed",
  "amount": "25.00",
  "currency": "EUR",
  "fee": "0.00",
  "counterparty": {
    "name": "Bob Jones",
    "username": "TP-87654321"
  },
  "created_at": "2026-03-03T15:00:00Z"
}
```

### 4.4 Currency Exchange

```
GET /api/v1/payments/fx/quote?source=EUR&target=USD&amount=100

Response: 200 OK
{
  "quote_id": "uuid",
  "source_currency": "EUR",
  "target_currency": "USD",
  "source_amount": "100.00",
  "target_amount": "108.50",
  "rate": "1.08500",
  "mid_market_rate": "1.09000",
  "markup": "0.46%",
  "fee": "0.00",
  "expires_at": "2026-03-03T15:00:30Z"
}
```

```
POST /api/v1/payments/fx/execute
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid

Request:
{
  "quote_id": "uuid",
  "source_sub_account_id": "uuid",
  "target_sub_account_id": "uuid"
}

Response: 201 Created
{
  "id": "uuid",
  "type": "fx_exchange",
  "status": "completed",
  "source": {
    "currency": "EUR",
    "amount": "100.00"
  },
  "target": {
    "currency": "USD",
    "amount": "108.50"
  },
  "rate": "1.08500",
  "created_at": "2026-03-03T15:00:05Z"
}
```

### 4.5 Get Payment Status

```
GET /api/v1/payments/{payment_id}

Response: 200 OK
{
  "id": "uuid",
  "type": "sepa_sct",
  "status": "completed",
  "status_history": [
    {"status": "created", "at": "2026-03-03T15:00:00Z"},
    {"status": "authorized", "at": "2026-03-03T15:00:01Z"},
    {"status": "submitted", "at": "2026-03-03T15:00:02Z"},
    {"status": "completed", "at": "2026-03-04T09:00:00Z"}
  ],
  "amount": "500.00",
  "currency": "EUR",
  "fee": "0.20",
  "destination": {
    "name": "Jane Smith",
    "iban": "DE89370400440532013000"
  },
  "reference": "Invoice #1234",
  "end_to_end_id": "TP2026030300001"
}
```

### 4.6 Scheduled Payments

```
POST /api/v1/payments/scheduled

Request:
{
  "source_sub_account_id": "uuid",
  "destination": {
    "iban": "DE89370400440532013000",
    "name": "Landlord LLC"
  },
  "amount": "850.00",
  "currency": "EUR",
  "reference": "Rent",
  "frequency": "monthly",
  "start_date": "2026-04-01",
  "end_date": null
}

Response: 201 Created
{
  "id": "uuid",
  "status": "active",
  "next_execution": "2026-04-01",
  "frequency": "monthly"
}
```

```
GET /api/v1/payments/scheduled
PATCH /api/v1/payments/scheduled/{id}
DELETE /api/v1/payments/scheduled/{id}

PATCH /api/v1/payments/scheduled/{id}
Request:
{
  "status": "paused"
}
```

---

## 5. Card API

### 5.1 Request Virtual Card

```
POST /api/v1/cards/virtual
Headers:
  Idempotency-Key: uuid

Request:
{
  "sub_account_id": "uuid",
  "cardholder_name": "JOHN DOE"
}

Response: 201 Created
{
  "id": "uuid",
  "type": "virtual",
  "brand": "mastercard",
  "last_four": "4321",
  "expiry": "03/29",
  "cardholder_name": "JOHN DOE",
  "status": "active",
  "linked_currency": "EUR",
  "created_at": "2026-03-03T15:00:00Z"
}
```

### 5.2 Get Card Details (Sensitive)

```
POST /api/v1/cards/{card_id}/details
Headers:
  X-SCA-Token: sca-token-uuid    -- Biometric required

Response: 200 OK
{
  "card_number": "5412345678904321",
  "expiry": "03/29",
  "cvv": "123",
  "display_timeout": 10           -- Hide after 10 seconds in UI
}
```

### 5.3 Request Physical Card

```
POST /api/v1/cards/physical
Headers:
  Idempotency-Key: uuid

Request:
{
  "sub_account_id": "uuid",
  "cardholder_name": "JOHN DOE",
  "delivery_address": {
    "line1": "Gedimino pr. 1",
    "city": "Vilnius",
    "postal_code": "01103",
    "country": "LT"
  }
}

Response: 201 Created
{
  "id": "uuid",
  "type": "physical",
  "brand": "mastercard",
  "last_four": "4321",
  "status": "inactive",
  "estimated_delivery": "2026-03-13",
  "created_at": "2026-03-03T15:00:00Z"
}
```

### 5.4 Activate Physical Card

```
POST /api/v1/cards/{card_id}/activate

Request:
{
  "last_four": "4321"
}

Response: 200 OK
{
  "id": "uuid",
  "status": "active",
  "activated_at": "2026-03-13T12:00:00Z"
}
```

### 5.5 Freeze/Unfreeze Card

```
POST /api/v1/cards/{card_id}/freeze

Response: 200 OK
{
  "id": "uuid",
  "status": "frozen",
  "frozen_at": "2026-03-03T15:00:00Z"
}
```

```
POST /api/v1/cards/{card_id}/unfreeze

Response: 200 OK
{
  "id": "uuid",
  "status": "active"
}
```

### 5.6 Update Spending Controls

```
PUT /api/v1/cards/{card_id}/controls

Request:
{
  "per_transaction_limit": "500.00",
  "daily_limit": "2000.00",
  "monthly_limit": "10000.00",
  "atm_daily_limit": "500.00",
  "online_enabled": true,
  "contactless_enabled": true,
  "atm_enabled": true,
  "magstripe_enabled": false,
  "blocked_mcc_codes": ["7995"],        // Gambling
  "allowed_countries": []                // Empty = all countries
}

Response: 200 OK
{
  // Updated controls echoed back
}
```

### 5.7 View/Change PIN

```
POST /api/v1/cards/{card_id}/pin/view
Headers:
  X-SCA-Token: sca-token-uuid

Response: 200 OK
{
  "pin": "1234",
  "display_timeout": 10
}
```

```
POST /api/v1/cards/{card_id}/pin/change
Headers:
  X-SCA-Token: sca-token-uuid

Request:
{
  "new_pin": "5678"
}

Response: 200 OK
{
  "pin_changed": true,
  "effective_at": "2026-03-03T15:00:30Z"
}
```

### 5.8 Report Lost/Stolen

```
POST /api/v1/cards/{card_id}/report

Request:
{
  "reason": "stolen",
  "request_replacement": true
}

Response: 200 OK
{
  "blocked_card_id": "uuid",
  "blocked_at": "2026-03-03T15:00:00Z",
  "replacement_card": {
    "id": "uuid-new",
    "type": "virtual",
    "status": "active",
    "last_four": "8765"
  }
}
```

### 5.9 Add to Apple Pay / Google Pay

```
POST /api/v1/cards/{card_id}/tokenize

Request:
{
  "wallet_type": "apple_pay",
  "device_certificates": ["base64-cert"],     // Apple Pay provisioning data
  "nonce": "base64-nonce",
  "nonce_signature": "base64-sig"
}

Response: 200 OK
{
  "activation_data": "base64-encrypted-data",  // Pass to Apple/Google SDK
  "encrypted_pass_data": "base64-data",
  "ephemeral_public_key": "base64-key"
}
```

### 5.10 List Card Transactions

```
GET /api/v1/cards/{card_id}/transactions?limit=20

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "type": "purchase",
      "status": "settled",
      "amount": "-45.99",
      "currency": "EUR",
      "merchant": {
        "name": "Amazon.de",
        "category": "Online Shopping",
        "mcc": "5411",
        "country": "DE"
      },
      "is_contactless": false,
      "is_apple_pay": false,
      "authorized_at": "2026-03-03T14:30:00Z",
      "settled_at": "2026-03-04T09:00:00Z"
    }
  ],
  "pagination": {
    "has_more": true,
    "next_cursor": "eyJpZCI6..."
  }
}
```

### 5.11 Submit Dispute

```
POST /api/v1/cards/{card_id}/transactions/{transaction_id}/dispute

Request:
{
  "reason": "unauthorized",
  "description": "I did not make this purchase",
  "evidence_files": []
}

Response: 201 Created
{
  "dispute_id": "uuid",
  "status": "submitted",
  "expected_resolution_date": "2026-03-24"
}
```

---

## 6. Crypto API

### 6.1 Get Wallet

```
GET /api/v1/crypto/wallet

Response: 200 OK
{
  "id": "uuid",
  "smart_wallet_address": "0x742d35Cc6634C0532925a3b844Bc9e7595f2BD25",
  "network": "fuse",
  "chain_id": 122,
  "status": "active",
  "balances": [
    {
      "token": "FUSE",
      "symbol": "FUSE",
      "balance": "150.500000000000000000",
      "balance_eur": "7.53",
      "price_eur": "0.0500",
      "change_24h": "-2.50"
    },
    {
      "token": "USDC",
      "symbol": "USDC",
      "contract_address": "0x620fd5fa44BE6af63d68D7C73F8e17CEd3bc2FC8",
      "balance": "500.000000",
      "balance_eur": "460.00",
      "price_eur": "0.9200",
      "change_24h": "0.01"
    },
    {
      "token": "USDT",
      "symbol": "USDT",
      "contract_address": "0xFaDbBF8Ce7D5b7A1ba969a6e10C6c31E5d790F8B",
      "balance": "0.000000",
      "balance_eur": "0.00",
      "price_eur": "0.9200",
      "change_24h": "-0.02"
    }
  ],
  "total_balance_eur": "467.53",
  "created_at": "2026-01-15T10:30:00Z"
}
```

### 6.2 Get Deposit Address

```
GET /api/v1/crypto/wallet/deposit-address?token=USDC

Response: 200 OK
{
  "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f2BD25",
  "network": "Fuse Network",
  "chain_id": 122,
  "token": "USDC",
  "qr_code_url": "https://api.teslapay.eu/qr/0x742d35...",
  "warning": "Only send USDC on Fuse Network. Tokens sent on other networks will be lost."
}
```

### 6.3 Buy Crypto

```
GET /api/v1/crypto/quote?action=buy&token=USDC&fiat_amount=200&fiat_currency=EUR

Response: 200 OK
{
  "quote_id": "uuid",
  "action": "buy",
  "token": "USDC",
  "fiat_amount": "200.00",
  "fiat_currency": "EUR",
  "token_amount": "215.68",
  "exchange_rate": "1.0784",
  "fee": "2.00",
  "fee_currency": "EUR",
  "net_fiat_amount": "198.00",
  "expires_at": "2026-03-03T15:00:30Z"
}
```

```
POST /api/v1/crypto/buy
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid

Request:
{
  "quote_id": "uuid",
  "source_sub_account_id": "uuid"
}

Response: 201 Created
{
  "order_id": "uuid",
  "status": "executing",
  "token": "USDC",
  "token_amount": "215.68",
  "fiat_amount": "200.00",
  "fiat_currency": "EUR",
  "fee": "2.00"
}
```

### 6.4 Sell Crypto

```
POST /api/v1/crypto/sell
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid

Request:
{
  "quote_id": "uuid",
  "target_sub_account_id": "uuid"
}

Response: 201 Created
{
  "order_id": "uuid",
  "status": "executing",
  "token": "USDC",
  "token_amount": "100.00",
  "fiat_amount": "92.00",
  "fiat_currency": "EUR",
  "fee": "0.92"
}
```

### 6.5 Send Crypto

```
POST /api/v1/crypto/send
Headers:
  Idempotency-Key: uuid
  X-SCA-Token: sca-token-uuid

Request:
{
  "token": "USDC",
  "amount": "50.000000",
  "to_address": "0x1234567890abcdef1234567890abcdef12345678",
  "note": "Payment for services"
}

Response: 201 Created
{
  "id": "uuid",
  "tx_hash": "0xabc123...",
  "status": "pending",
  "token": "USDC",
  "amount": "50.000000",
  "to_address": "0x1234...",
  "gas_fee": "0.50",
  "gas_fee_token": "USDC",          // Gasless: fee in same token
  "created_at": "2026-03-03T15:00:00Z"
}
```

### 6.6 Crypto Transaction History

```
GET /api/v1/crypto/transactions?token=USDC&type=send&limit=20

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "type": "send",
      "status": "confirmed",
      "token": "USDC",
      "amount": "50.000000",
      "fiat_value_eur": "46.00",
      "from_address": "0x742d35...",
      "to_address": "0x1234...",
      "tx_hash": "0xabc123...",
      "block_number": 28456789,
      "confirmations": 12,
      "gas_fee": "0.50",
      "gas_fee_token": "USDC",
      "explorer_url": "https://explorer.fuse.io/tx/0xabc123...",
      "created_at": "2026-03-03T15:00:00Z",
      "confirmed_at": "2026-03-03T15:00:15Z"
    }
  ],
  "pagination": {
    "has_more": false,
    "next_cursor": null
  }
}
```

### 6.7 Get Price Feed

```
GET /api/v1/crypto/prices

Response: 200 OK
{
  "prices": [
    {
      "token": "FUSE",
      "price_eur": "0.0500",
      "price_usd": "0.0543",
      "change_24h_pct": "-2.50",
      "volume_24h_eur": "125000.00"
    },
    {
      "token": "USDC",
      "price_eur": "0.9200",
      "price_usd": "1.0000",
      "change_24h_pct": "0.01",
      "volume_24h_eur": "5000000.00"
    },
    {
      "token": "USDT",
      "price_eur": "0.9200",
      "price_usd": "1.0000",
      "change_24h_pct": "-0.02",
      "volume_24h_eur": "4800000.00"
    }
  ],
  "updated_at": "2026-03-03T14:59:45Z"
}
```

---

## 7. KYC API

### 7.1 Initiate Verification

```
POST /api/v1/kyc/verify

Request:
{
  "level": "basic"
}

Response: 200 OK
{
  "verification_id": "uuid",
  "sumsub_access_token": "sb-token-xxxxx",
  "sumsub_flow_name": "teslapay-basic",
  "status": "initiated"
}
```

The `sumsub_access_token` is used by the Sumsub Flutter SDK to render the verification flow natively in the app.

### 7.2 Get Verification Status

```
GET /api/v1/kyc/status

Response: 200 OK
{
  "verification_id": "uuid",
  "level": "basic",
  "status": "approved",
  "completed_at": "2026-01-15T10:35:00Z",
  "checks": {
    "document": "passed",
    "liveness": "passed",
    "aml_screening": "clear",
    "nfc_verification": "not_performed"
  },
  "next_review_date": "2027-01-15",
  "upgrade_available": true,
  "upgrade_requirements": [
    "proof_of_address",
    "source_of_funds"
  ]
}
```

### 7.3 Request Tier Upgrade

```
POST /api/v1/kyc/upgrade

Request:
{
  "target_tier": "standard"
}

Response: 200 OK
{
  "verification_id": "uuid",
  "sumsub_access_token": "sb-token-xxxxx",
  "sumsub_flow_name": "teslapay-enhanced",
  "required_documents": ["proof_of_address"],
  "status": "initiated"
}
```

### 7.4 KYC Webhook (Internal -- from Sumsub)

```
POST /internal/webhooks/sumsub

Headers:
  X-Payload-Digest: hmac-sha256-signature

Request (from Sumsub):
{
  "type": "applicantReviewed",
  "applicantId": "sumsub-applicant-id",
  "inspectionId": "sumsub-inspection-id",
  "correlationId": "teslapay-user-id",
  "externalUserId": "TP-12345678",
  "reviewResult": {
    "reviewAnswer": "GREEN",
    "rejectLabels": [],
    "reviewRejectType": null
  },
  "createdAt": "2026-01-15T10:35:00Z"
}

Response: 200 OK
```

---

## 8. Notification API

### 8.1 Get Notification Preferences

```
GET /api/v1/notifications/preferences

Response: 200 OK
{
  "push_transactions": true,
  "push_security": true,
  "push_marketing": false,
  "email_transactions": true,
  "email_security": true,
  "email_marketing": false,
  "sms_security": true
}
```

### 8.2 Update Notification Preferences

```
PUT /api/v1/notifications/preferences

Request:
{
  "push_marketing": true,
  "email_marketing": true
}

Response: 200 OK
```

### 8.3 List In-App Notifications

```
GET /api/v1/notifications?unread_only=true&limit=20

Response: 200 OK
{
  "data": [
    {
      "id": "uuid",
      "type": "transaction",
      "title": "Payment Received",
      "body": "You received EUR 500.00 from Jane Smith",
      "action_url": "/transactions/uuid",
      "is_read": false,
      "created_at": "2026-03-03T14:30:00Z"
    }
  ],
  "unread_count": 3,
  "pagination": {
    "has_more": true,
    "next_cursor": "eyJpZCI6..."
  }
}
```

### 8.4 Mark Notification Read

```
POST /api/v1/notifications/{id}/read

Response: 204 No Content
```

---

## 9. Webhook Specifications (Outgoing)

TeslaPay exposes webhook delivery for business partners and Phase 2 Open Banking integrations.

### 9.1 Webhook Format

```json
{
  "id": "uuid",
  "type": "payment.completed",
  "api_version": "2026-03-01",
  "created_at": "2026-03-03T15:00:00Z",
  "data": {
    "payment_id": "uuid",
    "amount": "500.00",
    "currency": "EUR",
    "status": "completed"
  }
}
```

### 9.2 Webhook Security

- HMAC-SHA256 signature in `X-TeslaPay-Signature` header
- Timestamp in `X-TeslaPay-Timestamp` header (reject if > 5 min old)
- Signature computed over: `timestamp.body`
- Webhook retries: 3 attempts with exponential backoff (1 min, 5 min, 30 min)
- Webhook endpoint must respond 2xx within 15 seconds

### 9.3 Webhook Event Types

| Event | Description |
|-------|-------------|
| `payment.created` | Payment order created |
| `payment.completed` | Payment settled successfully |
| `payment.failed` | Payment failed |
| `payment.returned` | SEPA return received |
| `card.authorized` | Card authorization approved |
| `card.declined` | Card authorization declined |
| `card.settled` | Card transaction settled |
| `kyc.approved` | KYC verification approved |
| `kyc.rejected` | KYC verification rejected |
| `crypto.deposit.received` | Incoming crypto deposit confirmed |
| `account.updated` | Account status changed |

---

## 10. Admin API (Internal)

Secured by API key + VPN allowlisting. Used by the compliance dashboard and internal tools.

### 10.1 User Management

```
GET    /admin/v1/users?status=active&kyc=verified&page=1
GET    /admin/v1/users/{id}
PATCH  /admin/v1/users/{id}         -- Suspend, close, update tier
GET    /admin/v1/users/{id}/audit   -- Full audit trail
```

### 10.2 KYC Review Queue

```
GET    /admin/v1/kyc/queue?status=open&priority=high
POST   /admin/v1/kyc/queue/{id}/assign
POST   /admin/v1/kyc/queue/{id}/decide

Request:
{
  "decision": "approved",
  "reason": "Documents verified manually"
}
```

### 10.3 Transaction Monitoring

```
GET    /admin/v1/monitoring/alerts?status=open
GET    /admin/v1/monitoring/alerts/{id}
POST   /admin/v1/monitoring/alerts/{id}/review

Request:
{
  "action": "clear",
  "notes": "False positive - verified customer"
}
```

### 10.4 SAR Filing

```
POST   /admin/v1/compliance/sar

Request:
{
  "user_id": "uuid",
  "alert_ids": ["uuid", "uuid"],
  "narrative": "Suspicious structuring pattern detected...",
  "report_type": "SAR"
}
```

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CTO | TBD | | Pending |
| API Lead | TBD | | Pending |
| Principal Architect | Dream Team Architect | 2026-03-03 | Submitted |
