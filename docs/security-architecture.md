# TeslaPay Security Architecture

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Principal Software Architect, Dream Team

---

## 1. Security Architecture Overview

TeslaPay operates under a **defense-in-depth** model with multiple security layers. As a regulated EMI handling customer funds and card data, the security architecture must satisfy PSD2/SCA, PCI DSS, GDPR, and AMLD6 requirements simultaneously.

### 1.1 Security Principles

| Principle | Application |
|-----------|-------------|
| Zero Trust | Every request authenticated and authorized, regardless of network position |
| Least Privilege | Services and users have minimum necessary permissions |
| Defense in Depth | Multiple independent security controls at each layer |
| Secure by Default | All defaults are restrictive; access is explicitly granted |
| Fail Secure | System failures result in denied access, not open access |
| Separation of Duties | No single person can complete a high-risk operation alone |
| Immutable Audit | All security-relevant events logged immutably |

### 1.2 Security Layers

```
Layer 1: Edge Security
  +-------------------------------------------------------+
  |  AWS WAF + CloudFront + AWS Shield (DDoS)              |
  |  - OWASP Top 10 rules                                  |
  |  - Bot detection                                        |
  |  - Geo-blocking (sanctioned countries)                  |
  |  - Rate limiting (coarse)                               |
  +-------------------------------------------------------+
                            |
Layer 2: API Gateway Security
  +-------------------------------------------------------+
  |  Kong API Gateway                                       |
  |  - JWT validation                                       |
  |  - Fine-grained rate limiting (per user/tier)           |
  |  - Request schema validation                            |
  |  - Certificate pinning enforcement                      |
  |  - IP allowlisting (admin API)                          |
  +-------------------------------------------------------+
                            |
Layer 3: Service Mesh Security
  +-------------------------------------------------------+
  |  Istio Service Mesh                                     |
  |  - mTLS between all services (automatic)                |
  |  - Authorization policies (service-to-service)          |
  |  - Network segmentation via Kubernetes NetworkPolicies  |
  +-------------------------------------------------------+
                            |
Layer 4: Application Security
  +-------------------------------------------------------+
  |  Service-Level Security                                 |
  |  - Input validation and sanitization                    |
  |  - Business logic authorization (RBAC + ABAC)           |
  |  - SCA enforcement for financial operations             |
  |  - PII field-level encryption                           |
  |  - Idempotency key validation                           |
  +-------------------------------------------------------+
                            |
Layer 5: Data Security
  +-------------------------------------------------------+
  |  Data-at-Rest and Data-in-Transit                       |
  |  - AES-256 encryption at rest (AWS KMS)                 |
  |  - TLS 1.3 in transit                                   |
  |  - Field-level encryption for PII (Vault Transit)       |
  |  - Database-level RLS policies                          |
  +-------------------------------------------------------+
```

---

## 2. Authentication Architecture

### 2.1 Authentication Flow

```
+-------------+     +----------+     +-------------+     +-------------+
| Mobile App  |     | API GW   |     | Auth Service|     | Redis       |
+------+------+     +----+-----+     +------+------+     +------+------+
       |                  |                  |                   |
       | 1. Login (email + password)         |                   |
       +----------------->+----------------->|                   |
       |                  |                  |                   |
       |                  |   2. Validate credentials            |
       |                  |   3. Check device fingerprint        |
       |                  |   4. Generate JWT (RS256)            |
       |                  |   5. Generate refresh token          |
       |                  |                  |                   |
       |                  |                  | 6. Store session  |
       |                  |                  +------------------>|
       |                  |                  |                   |
       | 7. Return {access_token, refresh_token}                |
       |<-----------------+<-----------------+                   |
       |                  |                  |                   |
       | 8. Subsequent API calls             |                   |
       | Authorization: Bearer <JWT>         |                   |
       +----------------->| 9. Validate JWT  |                   |
       |                  | (local, no call) |                   |
       |                  +------>           |                   |
       |                  |                  |                   |
       | 10. JWT expired  |                  |                   |
       | Refresh flow     |                  |                   |
       +----------------->+----------------->| 11. Validate      |
       |                  |                  |     refresh token  |
       |                  |                  +------------------>|
       |                  |                  |<-----------------+|
       | 12. New tokens   |                  | 13. Rotate refresh|
       |<-----------------+<-----------------+     token         |
       |                  |                  |                   |
```

### 2.2 JWT Token Structure

```json
{
  "header": {
    "alg": "RS256",
    "typ": "JWT",
    "kid": "key-id-2026-03"
  },
  "payload": {
    "iss": "https://auth.teslapay.eu",
    "sub": "user-uuid",
    "aud": "https://api.teslapay.eu",
    "exp": 1709478400,
    "iat": 1709477500,
    "jti": "unique-token-id",
    "scope": "user:read user:write payment:write card:read card:write crypto:read crypto:write",
    "tier": "standard",
    "kyc": "verified",
    "device_id": "device-uuid",
    "sca_level": 0
  }
}
```

**Token Specifications:**

| Token Type | Algorithm | Lifetime | Storage |
|-----------|-----------|----------|---------|
| Access Token | RS256 (asymmetric) | 15 minutes | In-memory only (mobile) |
| Refresh Token | Opaque (random 256-bit) | 30 days | flutter_secure_storage |
| SCA Token | RS256 | 5 minutes | In-memory only |

**Key Rotation:**
- RSA key pairs rotated every 90 days
- Old keys remain valid for token verification until all issued tokens expire
- JWKS endpoint at `/.well-known/jwks.json` publishes current and previous public keys

### 2.3 Biometric Authentication

```
Device Registration:
  1. User enables biometric login
  2. App generates asymmetric key pair in device secure enclave
     (Keychain with kSecAttrAccessibleWhenPasscodeSetThisDeviceOnly on iOS,
      Android Keystore with setUserAuthenticationRequired)
  3. Public key sent to Auth Service and stored in devices table
  4. Private key never leaves secure enclave

Biometric Login:
  1. App creates challenge payload (timestamp + device_id + nonce)
  2. User authenticates with Face ID / fingerprint
  3. Secure enclave signs challenge with private key
  4. Signature + payload sent to Auth Service
  5. Auth Service verifies signature against stored public key
  6. If valid, issue JWT tokens
```

### 2.4 Strong Customer Authentication (PSD2 SCA)

PSD2 requires SCA for:
- Electronic payment transactions
- Actions with fraud risk (adding new beneficiary, changing phone/email)
- Accessing payment account data (first time or after 90+ days)

**SCA Implementation:**

Two of three factors required:
1. **Knowledge:** App PIN or password
2. **Possession:** Registered device (device_id + push notification)
3. **Inherence:** Biometric (Face ID, fingerprint)

**SCA Flow:**

```
1. User initiates SCA-required action (e.g., payment > EUR 30)
2. API returns 403 with SCA challenge requirement
3. App shows confirmation screen with transaction details
4. User authenticates with biometric (inherence) on registered device (possession)
5. App sends signed SCA verification to Auth Service
6. Auth Service issues short-lived SCA token (5 min TTL)
7. App retries original request with X-SCA-Token header
8. API Gateway validates SCA token and proceeds
```

**SCA Exemptions (per PSD2):**
- Low-value transactions (< EUR 30, cumulative < EUR 100)
- Trusted beneficiaries (user-whitelisted payees)
- Recurring transactions with same amount and payee
- Contactless payments (< EUR 50)

---

## 3. Encryption Architecture

### 3.1 Encryption at Rest

| Data Store | Encryption Method | Key Management |
|-----------|-------------------|---------------|
| RDS PostgreSQL | AES-256 (AWS managed) | AWS KMS CMK per database |
| S3 Buckets | AES-256-GCM (SSE-KMS) | AWS KMS CMK per bucket |
| ElastiCache Redis | AES-256 | AWS KMS |
| EBS Volumes | AES-256 | AWS KMS |
| MSK Kafka | AES-256 | AWS KMS |
| OpenSearch | AES-256 | AWS KMS |
| Backup Snapshots | AES-256 | AWS KMS (separate backup key) |

### 3.2 Encryption in Transit

| Channel | Protocol | Notes |
|---------|----------|-------|
| Client to API Gateway | TLS 1.3 | Certificate pinning enforced in mobile app |
| Service to Service | mTLS (Istio) | Automatic certificate rotation |
| Service to Database | TLS 1.3 | PostgreSQL `sslmode=verify-full` |
| Service to Redis | TLS 1.3 | ElastiCache in-transit encryption |
| Service to Kafka | TLS 1.3 | MSK TLS listeners only |
| Service to External APIs | TLS 1.3 | Certificate validation enforced |

### 3.3 Field-Level Encryption (PII)

Sensitive PII fields are encrypted at the application level using HashiCorp Vault Transit engine before storage in PostgreSQL. This provides an additional layer beyond database-level encryption.

**Encrypted Fields:**

| Table | Field | Rationale |
|-------|-------|-----------|
| `users` | `tax_id` | Tax identification number -- highly sensitive |
| `users` | `date_of_birth` | Personal identifier |
| `verifications` | `document_number_masked` | Even masked doc numbers are protected |
| `audit_events.changes` | PII field values | Old/new values of PII changes |

**Vault Transit Engine Configuration:**
```
Key:       teslapay-pii-key
Type:      aes256-gcm96
Rotation:  Every 90 days (automatic)
Convergent: No (non-deterministic encryption -- same plaintext produces different ciphertext)
```

**Encryption Flow:**
```go
// Encrypt before storage
func encryptPII(ctx context.Context, plaintext string) (string, error) {
    secret, err := vaultClient.Logical().Write("transit/encrypt/teslapay-pii-key", map[string]interface{}{
        "plaintext": base64.StdEncoding.EncodeToString([]byte(plaintext)),
    })
    return secret.Data["ciphertext"].(string), err
}

// Decrypt when reading
func decryptPII(ctx context.Context, ciphertext string) (string, error) {
    secret, err := vaultClient.Logical().Write("transit/decrypt/teslapay-pii-key", map[string]interface{}{
        "ciphertext": ciphertext,
    })
    decoded, _ := base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
    return string(decoded), err
}
```

### 3.4 Certificate Pinning (Mobile App)

```dart
// Flutter certificate pinning configuration
class CertificatePinningInterceptor extends Interceptor {
  static const pins = {
    'api.teslapay.eu': [
      'sha256/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=', // Primary
      'sha256/BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=', // Backup
    ],
  };

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    // Pin validation happens in SecurityContext
    // If pin mismatch: throw CertificatePinningException
    // App shows "Secure connection could not be established"
  }
}
```

Pin rotation strategy:
- Always include at least 2 pins (current + next)
- App update published with new pin before old certificate expires
- Emergency: remote config flag can temporarily disable pinning (with increased monitoring)

---

## 4. PCI DSS Compliance Architecture

### 4.1 Scope Minimization Strategy

TeslaPay minimizes PCI DSS scope by never storing, processing, or transmitting cardholder data:

```
+------------------------------------------------------------------+
|                        OUT OF SCOPE                                |
|                   (TeslaPay Infrastructure)                        |
|                                                                    |
|  +-------------------+    +-------------------+                    |
|  | API Gateway       |    | Card Service      |                    |
|  | (no card data)    |    | (tokenized refs   |                    |
|  |                   |    |  only)             |                    |
|  +-------------------+    +-------------------+                    |
|                                                                    |
+------------------------------------------------------------------+
                              |
                              | Tokenized references only
                              | (processor_card_id, last_four)
                              |
+------------------------------------------------------------------+
|                        IN SCOPE (Enfuce)                           |
|                                                                    |
|  +-------------------+    +-------------------+                    |
|  | Card Processing   |    | HSM               |                    |
|  | (PAN, CVV, PIN    |    | (Key Storage)     |                    |
|  |  stored here)     |    |                   |                    |
|  +-------------------+    +-------------------+                    |
|                                                                    |
|  PCI DSS Level 1 Certified                                        |
+------------------------------------------------------------------+
```

### 4.2 TeslaPay PCI DSS Obligations

Even with scope minimization, TeslaPay must:

| Requirement | Implementation |
|-------------|---------------|
| SAQ-A compliance | Annual self-assessment questionnaire |
| Secure communication with processor | mTLS + TLS 1.3 to Enfuce API |
| No card data in logs | Log redaction filters for any potential PAN patterns |
| Secure card display | Card details (PAN, CVV) displayed via Enfuce secure display SDK, never passed through TeslaPay API |
| PIN management | PIN set/view operations proxied directly to Enfuce; PIN never in TeslaPay memory |
| Employee training | Annual PCI DSS awareness training |
| Incident response | Card data breach notification procedure to Mastercard within 24 hours |

### 4.3 Card Data Redaction

All logging infrastructure includes PAN detection and redaction:

```go
// PAN redaction regex applied to all log output
var panRegex = regexp.MustCompile(`\b(?:\d[ -]*?){13,19}\b`)

func redactPAN(input string) string {
    return panRegex.ReplaceAllStringFunc(input, func(match string) string {
        digits := strings.Map(func(r rune) rune {
            if r >= '0' && r <= '9' { return r }
            return -1
        }, match)
        if luhnCheck(digits) {
            return "****REDACTED****"
        }
        return match
    })
}
```

---

## 5. Key Management

### 5.1 Key Hierarchy

```
+-------------------------------------------------------+
|  AWS CloudHSM (FIPS 140-2 Level 3)                    |
|  +---------------------------------------------------+|
|  | Master Key (HSM Root)                              ||
|  | - Never exported                                   ||
|  | - Used to wrap all other keys                      ||
|  +---------------------------------------------------+|
+-------------------------------------------------------+
                            |
                            | Wraps
                            v
+-------------------------------------------------------+
|  AWS KMS (Customer Managed Keys)                       |
|  +-------------------+  +---------------------------+  |
|  | Database CMK      |  | Storage CMK               |  |
|  | (RDS encryption)  |  | (S3, EBS encryption)      |  |
|  +-------------------+  +---------------------------+  |
|  +-------------------+  +---------------------------+  |
|  | Kafka CMK         |  | Backup CMK                |  |
|  | (MSK encryption)  |  | (snapshot encryption)     |  |
|  +-------------------+  +---------------------------+  |
+-------------------------------------------------------+
                            |
                            | Wraps
                            v
+-------------------------------------------------------+
|  HashiCorp Vault                                       |
|  +-------------------+  +---------------------------+  |
|  | Transit Engine    |  | PKI Engine                |  |
|  | (PII encryption)  |  | (mTLS certificates)      |  |
|  +-------------------+  +---------------------------+  |
|  +-------------------+  +---------------------------+  |
|  | KV Engine         |  | JWT Signing Keys          |  |
|  | (API keys, secrets)|  | (RS256 key pairs)        |  |
|  +-------------------+  +---------------------------+  |
+-------------------------------------------------------+
```

### 5.2 Key Rotation Schedule

| Key Type | Rotation Period | Method |
|----------|----------------|--------|
| HSM Master Key | Annually | Manual ceremony with dual control |
| KMS CMKs | Annually (automatic) | AWS KMS automatic rotation |
| Vault Transit Key | 90 days (automatic) | Vault auto-rotation |
| JWT Signing Keys | 90 days | Manual rotation with JWKS overlap |
| mTLS Certificates | 24 hours | Istio automatic rotation |
| Database Passwords | 30 days | Vault dynamic secrets |
| API Keys (external) | 180 days | Manual rotation with partner coordination |

### 5.3 Secrets Management with Vault

All application secrets are managed by HashiCorp Vault. No secrets in environment variables, config files, or code.

```
Vault Configuration:
  Storage Backend:  Consul (HA)
  Seal Type:        AWS KMS auto-unseal
  Auth Methods:     Kubernetes (pod identity), AppRole (CI/CD)
  Audit:            File + syslog (all operations logged)

Secret Paths:
  secret/data/teslapay/auth-service/db      -- Auth DB credentials
  secret/data/teslapay/payment-service/db   -- Payment DB credentials
  secret/data/teslapay/integrations/sumsub  -- Sumsub API keys
  secret/data/teslapay/integrations/enfuce  -- Enfuce API keys
  secret/data/teslapay/integrations/banking-circle -- Banking Circle creds
  secret/data/teslapay/integrations/fuse    -- Fuse API key
  transit/keys/teslapay-pii-key             -- PII encryption key
  pki/issue/service-mesh                    -- Dynamic TLS certificates
```

**Dynamic Database Credentials:**
```go
// Vault generates short-lived PostgreSQL credentials
// Lease duration: 1 hour, renewable
secret, err := vaultClient.Logical().Read("database/creds/ledger-service")
username := secret.Data["username"].(string)
password := secret.Data["password"].(string)
// Connect to PostgreSQL with dynamic credentials
// Vault automatically revokes credentials when lease expires
```

---

## 6. Fraud Detection System

### 6.1 Architecture

```
+------------------+     +------------------+     +------------------+
| Transaction      |     | Fraud Detection  |     | Alert            |
| Sources          |     | Service          |     | Consumers        |
+--------+---------+     +--------+---------+     +--------+---------+
         |                         |                        |
         | card.events            |                        |
         | payment.events         |                        |
         +----------------------->| 1. Rule Engine         |
         |                        |    (sync + async)      |
         |                        |                        |
         | gRPC (auth scoring)    | 2. Velocity Engine     |
         +----------------------->|    (Redis counters)    |
         |  <10ms response        |                        |
         |                        | 3. Pattern Matching    +-> KYC Service
         |                        |    (historical data)   |   (freeze account)
         |                        |                        |
         |                        | 4. ML Scoring (P2)     +-> Card Service
         |                        |    (anomaly detection) |   (block card)
         |                        |                        |
         |                        +----------------------->+-> Notification
         |                        |  fraud.signals         |   (alert user)
         |                        |                        |
         |                        +----------------------->+-> Audit Service
         |                        |  audit.events          |   (log decision)
         |                        |                        |
         |                        +----------------------->+-> Compliance
         |                        |  review queue          |   (SAR review)
```

### 6.2 Rule Engine (Phase 1)

| Rule Category | Rules | Action |
|---------------|-------|--------|
| **Velocity** | > 5 card transactions in 5 minutes | Flag for review |
| **Velocity** | > 10 failed PIN attempts in 1 hour | Block card |
| **Amount** | Single transaction > EUR 10,000 | Flag + SCA |
| **Amount** | Cumulative > EUR 15,000 in 30 days | Flag for AML review |
| **Geographic** | Card used in 2+ countries within 1 hour | Decline + alert user |
| **Geographic** | Card used in high-risk country | Flag + SCA |
| **Pattern** | Multiple small transactions (structuring) | Flag for AML review |
| **Pattern** | Rapid account funding + immediate withdrawal | Block + review |
| **Blacklist** | Merchant on blacklist | Decline |
| **Time** | Transaction at unusual hour for user | Increase risk score |
| **Channel** | Online transaction without 3DS | Increase risk score |

### 6.3 Velocity Counters (Redis)

```
Key patterns:
  fraud:velocity:card:{card_id}:txn_count:5m    -- Sorted set, last 5 min
  fraud:velocity:card:{card_id}:amount:24h       -- Sorted set, last 24h
  fraud:velocity:user:{user_id}:login_fail:1h    -- Counter, last 1h
  fraud:velocity:user:{user_id}:country_set:1h   -- Set of countries, last 1h
  fraud:velocity:user:{user_id}:cumulative:30d   -- Total amount, rolling 30d

Operations:
  - ZADD with timestamp score on each transaction
  - ZRANGEBYSCORE to count within window
  - ZREMRANGEBYSCORE to clean expired entries
  - EXPIRE keys after maximum window + buffer
```

### 6.4 Risk Scoring

Each transaction receives a risk score (0-100):

```
Base score: 0

Factors:
  +10: New device
  +15: New country
  +20: Amount > 90th percentile for user
  +10: Nighttime transaction (user's timezone)
  +25: Failed 3DS attempt in last hour
  +30: Velocity threshold exceeded
  +15: High-risk MCC code
  -10: Known merchant (previously used)
  -15: Trusted device
  -20: Normal amount range for user

Action thresholds:
  0-30:  Approve (normal)
  31-60: Approve with enhanced monitoring
  61-80: Require additional SCA
  81-100: Decline and alert
```

---

## 7. GDPR Data Handling

### 7.1 Data Classification

| Category | Examples | Retention | Encryption | Access |
|----------|----------|-----------|------------|--------|
| **Public** | Company info, product descriptions | N/A | No | Everyone |
| **Internal** | Employee names, internal docs | Employment + 1 year | Standard | Staff |
| **Confidential** | User profiles, transaction data | 5 years (regulatory) | AES-256 + field-level | Service accounts |
| **Restricted** | Tax IDs, KYC documents, biometrics | 5 years (regulatory) | AES-256 + Vault Transit | Dedicated service accounts |

### 7.2 Data Subject Rights Implementation

| Right | Implementation | SLA |
|-------|---------------|-----|
| **Right of Access** | Export all user data from all services via Audit Service API; delivered as JSON/PDF | 30 days |
| **Right to Rectification** | User updates via Account Service; compliance review for regulated fields | 5 business days |
| **Right to Erasure** | Anonymization of personal data; transaction data retained per regulatory requirement; user notified of exceptions | 30 days |
| **Right to Portability** | Machine-readable export (JSON) of all user-provided data | 30 days |
| **Right to Restriction** | Account frozen, data retained but not processed | Immediate |

### 7.3 Data Erasure Process

When a user requests deletion (right to erasure):

```
1. Audit Service receives data_request(type=deletion)
2. For each service:
   a. Auth Service: Delete credentials, sessions, devices
   b. Account Service: Anonymize user profile:
      - first_name -> "DELETED"
      - last_name -> "DELETED"
      - email -> hash(email) + "@deleted.teslapay.eu"
      - phone -> hash(phone)
      - address -> NULL
      - tax_id -> NULL
   c. Ledger Service: RETAIN transaction data (5-year regulatory requirement)
      - Remove user name from transaction references
   d. Payment Service: RETAIN payment orders (5-year regulatory requirement)
      - Anonymize payee names in beneficiary lists
   e. Card Service: Delete card data, cancel active cards
   f. Crypto Service: Delete wallet metadata (blockchain data is immutable)
   g. KYC Service: RETAIN verification records (5-year AMLD6 requirement)
      - Mark as "data_subject_erasure_requested"
      - Delete from Sumsub after retention period
   h. Notification Service: Delete delivery logs, preferences
3. Log erasure event in Audit Service (retained for accountability)
4. Notify user of completion and retention exceptions
```

### 7.4 Data Processing Records (GDPR Art. 30)

| Processing Activity | Legal Basis | Data Categories | Recipients |
|---------------------|-------------|-----------------|------------|
| Account opening | Contract performance | Identity, contact | Sumsub (KYC) |
| Payment processing | Contract performance | Transaction data | Banking Circle, Enfuce |
| KYC/AML screening | Legal obligation | Identity, docs | Sumsub, regulators |
| Card issuance | Contract performance | Identity, card data | Enfuce, Mastercard |
| Crypto wallet | Contract performance | Identity, wallet | Fuse.io |
| Marketing | Consent (opt-in) | Contact, preferences | Email/push providers |
| Fraud detection | Legitimate interest | Transaction patterns | Internal only |
| Regulatory reporting | Legal obligation | Transaction, identity | Bank of Lithuania |

---

## 8. Network Security

### 8.1 Network Architecture

```
+-------------------------------------------------------+
|  AWS VPC (10.0.0.0/16)                                |
|                                                        |
|  Public Subnets (10.0.1.0/24, 10.0.2.0/24)           |
|  +--------------------+  +--------------------+       |
|  | NLB                |  | NAT Gateway        |       |
|  | (ingress only)     |  | (egress only)      |       |
|  +--------------------+  +--------------------+       |
|                                                        |
|  Private Subnets - Application (10.0.10.0/24, ...)    |
|  +--------------------+  +--------------------+       |
|  | EKS Worker Nodes   |  | EKS Worker Nodes   |       |
|  | (teslapay-core)    |  | (teslapay-crypto)  |       |
|  +--------------------+  +--------------------+       |
|                                                        |
|  Private Subnets - PCI (10.0.20.0/24, ...)            |
|  +--------------------+                               |
|  | EKS Worker Nodes   |  <-- Extra NetworkPolicies    |
|  | (teslapay-card)    |  <-- Restricted egress         |
|  +--------------------+                               |
|                                                        |
|  Private Subnets - Data (10.0.30.0/24, ...)           |
|  +--------------------+  +--------------------+       |
|  | RDS PostgreSQL     |  | ElastiCache Redis  |       |
|  +--------------------+  +--------------------+       |
|  +--------------------+  +--------------------+       |
|  | MSK Kafka          |  | OpenSearch          |       |
|  +--------------------+  +--------------------+       |
|                                                        |
|  No direct internet access from private subnets        |
|  All egress via NAT Gateway with VPC Flow Logs         |
+-------------------------------------------------------+
```

### 8.2 AWS WAF Rules

| Rule Group | Purpose | Action |
|-----------|---------|--------|
| AWS Managed - Common | OWASP Top 10 (SQLi, XSS, LFI, RFI) | Block |
| AWS Managed - Known Bad Inputs | Exploit payloads | Block |
| AWS Managed - Bot Control | Automated request detection | Challenge/Block |
| Custom - Geo Blocking | Block sanctioned countries (OFAC list) | Block |
| Custom - Rate Limiting | > 1000 requests/min per IP | Block |
| Custom - Login Protection | > 10 failed logins/min per IP | Block |
| Custom - API Abuse | Unusual request patterns | Challenge |

### 8.3 DDoS Protection

| Layer | Protection |
|-------|-----------|
| Layer 3/4 | AWS Shield Standard (automatic, free) |
| Layer 7 | AWS WAF rate limiting + AWS Shield Advanced |
| Application | Kong rate limiting per user |
| Origin | CloudFront + NLB absorb traffic before reaching EKS |

### 8.4 Kubernetes Network Policies

```yaml
# Example: Card Service can only communicate with specific services
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: card-service-policy
  namespace: teslapay-card
spec:
  podSelector:
    matchLabels:
      app: card-service
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: teslapay-infra    # API Gateway only
          podSelector:
            matchLabels:
              app: kong
      ports:
        - port: 8080
  egress:
    - to:
        - namespaceSelector:
            matchLabels:
              name: teslapay-core     # Ledger, Account services
      ports:
        - port: 9090                  # gRPC
    - to:
        - ipBlock:
            cidr: 0.0.0.0/0          # Enfuce API (external)
      ports:
        - port: 443
    - to:
        - namespaceSelector:
            matchLabels:
              name: teslapay-data     # Kafka, Redis, PostgreSQL
```

---

## 9. Security Monitoring and Incident Response

### 9.1 Security Monitoring Stack

| Tool | Purpose | Alert Targets |
|------|---------|--------------|
| AWS GuardDuty | AWS account threat detection | Security Slack + PagerDuty |
| AWS Security Hub | Centralized security findings | Security dashboard |
| Falco | Container runtime anomaly detection | Security Slack + PagerDuty |
| OWASP ZAP | Continuous DAST scanning | Security Jira |
| Trivy | Container vulnerability scanning | CI/CD pipeline (block deploy) |
| AWS CloudTrail | AWS API audit trail | S3 + OpenSearch |
| VPC Flow Logs | Network traffic analysis | S3 + OpenSearch |
| Vault Audit Log | Secrets access audit | OpenSearch |

### 9.2 Security Alert Categories

| Category | Examples | Response SLA |
|----------|---------|--------------|
| Critical | Unauthorized data access, active breach, card data exposure | 15 minutes |
| High | Brute force attack, privilege escalation, unusual admin activity | 1 hour |
| Medium | Failed authentication spike, new IAM role, unusual API patterns | 4 hours |
| Low | Vulnerability scan triggered, minor policy violation | 24 hours |

### 9.3 Incident Response Plan

```
Phase 1: Detection and Triage (0-15 min)
  - Automated alert triggers PagerDuty
  - On-call engineer assesses severity
  - If Critical: Assemble incident team immediately

Phase 2: Containment (15-60 min)
  - Isolate affected systems (network policy, pod kill)
  - Revoke compromised credentials
  - Enable enhanced logging on affected services
  - If card data: Notify Enfuce for card blocking

Phase 3: Investigation (1-24 hours)
  - Analyze audit trails and access logs
  - Determine attack vector and scope
  - Identify all affected users and data

Phase 4: Remediation (24-72 hours)
  - Patch vulnerability
  - Rotate all potentially compromised credentials
  - Restore from clean backups if needed

Phase 5: Notification (per regulatory requirements)
  - GDPR: Notify DPA within 72 hours of personal data breach
  - Notify affected users "without undue delay"
  - Mastercard: Notify within 24 hours for card data incidents
  - Bank of Lithuania: Notify per EMI incident reporting requirements

Phase 6: Post-Incident (1-2 weeks)
  - Root cause analysis
  - Update security controls
  - Document lessons learned
  - Update incident response playbook
```

---

## 10. Security Testing Program

### 10.1 Testing Schedule

| Test Type | Frequency | Performed By |
|-----------|-----------|-------------|
| SAST (SonarQube) | Every PR | Automated (CI/CD) |
| DAST (OWASP ZAP) | Weekly (staging) | Automated |
| Container Scanning (Trivy) | Every build | Automated (CI/CD) |
| Dependency Scanning | Daily | Automated (Dependabot) |
| Penetration Testing | Quarterly + after major releases | External firm (CREST certified) |
| Red Team Exercise | Annually | External firm |
| Social Engineering Test | Bi-annually | External firm |
| PCI DSS Assessment | Annually | QSA (Qualified Security Assessor) |
| SOC 2 Type II Audit | Annually | External auditor |

### 10.2 Secure Development Lifecycle

```
1. Threat Modeling (design phase)
   - STRIDE analysis for each new feature
   - Data flow diagrams reviewed by security team

2. Secure Coding Standards
   - OWASP Secure Coding Practices checklist
   - Go-specific: no use of unsafe package, no raw SQL, parameterized queries only
   - Input validation on all API endpoints (go-playground/validator)

3. Code Review (development phase)
   - 2 reviewer minimum, 1 must be security-aware
   - Security-critical code (auth, crypto, payment) requires security team review

4. Automated Testing (CI/CD)
   - SAST, DAST, container scanning (see above)
   - Security unit tests for auth, encryption, validation logic

5. Pre-Production (staging)
   - Security regression tests
   - Compliance check (PCI DSS, GDPR controls)

6. Production
   - Runtime protection (Falco)
   - Continuous monitoring (GuardDuty, WAF logs)
   - Bug bounty program (Phase 2)
```

---

## 11. Compliance Matrix

| Requirement | Standard | Implementation | Status |
|-------------|----------|---------------|--------|
| Data encryption at rest | PCI DSS 3.4, GDPR Art. 32 | AES-256 via AWS KMS | Architecture defined |
| Data encryption in transit | PCI DSS 4.1, GDPR Art. 32 | TLS 1.3, mTLS | Architecture defined |
| Access control | PCI DSS 7, GDPR Art. 25 | RBAC + ABAC, least privilege | Architecture defined |
| Audit logging | PCI DSS 10, AMLD6 | Immutable audit trail, 7-year retention | Architecture defined |
| Key management | PCI DSS 3.5, 3.6 | CloudHSM + Vault | Architecture defined |
| Network segmentation | PCI DSS 1.3 | VPC, subnets, NetworkPolicies | Architecture defined |
| Vulnerability management | PCI DSS 6.1, 6.2 | Trivy, SonarQube, ZAP, pen testing | Process defined |
| Strong authentication | PSD2 SCA, PCI DSS 8 | OAuth2 + biometric + device binding | Architecture defined |
| Data minimization | GDPR Art. 5 | Collect only necessary data, PII encryption | Architecture defined |
| Right to erasure | GDPR Art. 17 | Anonymization process with regulatory exceptions | Process defined |
| AML monitoring | AMLD6 | Sumsub ongoing monitoring + transaction rules | Architecture defined |
| Incident response | PCI DSS 12.10, GDPR Art. 33 | IR plan with 72-hour DPA notification | Process defined |

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CTO | TBD | | Pending |
| CISO / Security Lead | TBD | | Pending |
| DPO | TBD | | Pending |
| Principal Architect | Dream Team Architect | 2026-03-03 | Submitted |
