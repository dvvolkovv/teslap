# TeslaPay System Architecture

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Principal Software Architect, Dream Team
**Status:** Draft for Review

---

## 1. Architecture Overview

TeslaPay is a microservices-based, event-driven neobank platform deployed on Kubernetes in AWS EU (Frankfurt). The architecture follows domain-driven design (DDD), with each bounded context owning its data, communicating via Apache Kafka for asynchronous events and gRPC for synchronous inter-service calls. The mobile client (Flutter) connects through an API Gateway (Kong) that handles authentication, rate limiting, and request routing.

### 1.1 Architecture Principles

| Principle | Description |
|-----------|-------------|
| Domain-Driven Design | Services aligned to business domains with clear bounded contexts |
| Event Sourcing | All financial state changes stored as immutable events in append-only logs |
| CQRS | Separate read/write models for ledger and transaction-heavy services |
| Zero Trust | Every service-to-service call authenticated and authorized; no implicit trust |
| Data Sovereignty | All PII and financial data resides in EU (Frankfurt); DR in EU (Ireland) |
| Defense in Depth | Multiple layers of security controls; no single point of failure |
| Idempotency | All write operations are idempotent with client-supplied idempotency keys |
| Observability First | Structured logging, distributed tracing, and metrics from day one |

---

## 2. High-Level Architecture Diagram

```
                                 +---------------------+
                                 |   CloudFront CDN    |
                                 |   + AWS WAF         |
                                 +----------+----------+
                                            |
                                            v
                                 +---------------------+
                                 |   AWS NLB / ALB     |
                                 +----------+----------+
                                            |
                                            v
+-------------------+           +---------------------+           +--------------------+
|                   |           |                     |           |                    |
|  Flutter Mobile   |<--------->|   API Gateway       |<--------->|  Admin Dashboard   |
|  (iOS / Android)  |   HTTPS   |   (Kong)            |   HTTPS   |  (React / Next.js) |
|                   |           |                     |           |                    |
+-------------------+           +----------+----------+           +--------------------+
                                           |
                          +----------------+----------------+
                          |                |                |
                    gRPC/REST        gRPC/REST        gRPC/REST
                          |                |                |
              +-----------v--+    +--------v------+   +----v-----------+
              | Auth Service |    | Account       |   | Payment        |
              | (IAM, OAuth2,|    | Service       |   | Service        |
              | Sessions,    |    | (User Accts,  |   | (SEPA, FX,     |
              | MFA, SCA)    |    | Multi-Ccy,    |   | Internal Xfer, |
              +---------+----+    | IBAN, Tiers)  |   | Scheduling)    |
                        |         +-------+-------+   +-------+--------+
                        |                 |                    |
                        |    +------------+----+     +---------v--------+
                        |    | Ledger Service  |     | Card Service     |
                        |    | (Double-Entry,  |     | (Issuing, 3DS,   |
                        |    | Event Sourcing, |     | Lifecycle, Apple/|
                        |    | Reconciliation) |     | Google Pay)      |
                        |    +--------+--------+     +--------+---------+
                        |             |                       |
               +--------v--------+   |              +--------v---------+
               | KYC Service     |   |              | Crypto Service   |
               | (Sumsub, Risk   |   |              | (Fuse.io, Smart  |
               | Scoring, AML    |   |              | Wallets, Buy/    |
               | Ongoing Monitor)|   |              | Sell, ERC-4337)  |
               +---------+-------+   |              +--------+---------+
                         |           |                       |
               +---------v-------+   |              +--------v---------+
               | Notification    |   |              | Audit Service    |
               | Service         |   |              | (Immutable Log,  |
               | (Push, Email,   |   |              | Compliance Trail,|
               | SMS, In-App)    |   |              | Regulatory Rpt)  |
               +-----------------+   |              +------------------+
                                     |
                         +-----------v-----------+
                         |   Fraud Detection     |
                         |   Service             |
                         |   (Rules Engine,      |
                         |   Velocity Checks,    |
                         |   ML Scoring Phase 2) |
                         +-----------------------+

            +========================================================+
            |                   Apache Kafka Cluster                  |
            |  Topics: ledger.events, payment.events, card.events,   |
            |  kyc.events, crypto.events, audit.events,              |
            |  notification.commands, fraud.signals                   |
            +========================================================+

            +========================================================+
            |                   Data Layer                            |
            |                                                        |
            |  +-------------+  +-----------+  +------------------+  |
            |  | PostgreSQL  |  | Redis     |  | Elasticsearch    |  |
            |  | (per-svc    |  | Cluster   |  | (Search, Audit   |  |
            |  | databases)  |  | (Cache,   |  | Log Indexing)    |  |
            |  |             |  | Sessions) |  |                  |  |
            |  +-------------+  +-----------+  +------------------+  |
            |                                                        |
            |  +-------------+  +-----------+                        |
            |  | S3           |  | Vault    |                        |
            |  | (Documents,  |  | (Secrets,|                        |
            |  | Exports,     |  | HSM Keys,|                        |
            |  | Backups)     |  | Certs)   |                        |
            |  +-------------+  +-----------+                        |
            +========================================================+

                     External Integrations
            +========================================================+
            |  Sumsub API    |  Mastercard Network  |  Fuse.io RPC   |
            |  (KYC/AML)     |  (via Enfuce/Marqeta)|  (Blockchain)  |
            |                |                      |                |
            |  SEPA Gateway  |  Apple/Google Push   |  FX Provider   |
            |  (Banking      |  (APNs / FCM)        |  (ECB rates +  |
            |   Circle)      |                      |   aggregator)  |
            +========================================================+
```

---

## 3. Microservices Breakdown

### 3.1 Auth Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | User authentication, session management, OAuth2/OIDC token issuance, MFA/SCA enforcement, device binding, biometric credential management |
| Tech Stack | Go, PostgreSQL, Redis |
| Data Ownership | Users credentials table, sessions, devices, MFA secrets, refresh tokens |
| API Style | REST (public), gRPC (internal) |
| Key Patterns | JWT access tokens (15 min TTL), opaque refresh tokens (30 day TTL), device fingerprinting, SCA step-up for PSD2 |

**Rationale:** Auth is a high-throughput, security-critical service. Go provides low latency and strong concurrency. Separated from Account Service to isolate the security boundary and allow independent scaling during login spikes.

### 3.2 Account Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | User profile management, IBAN generation, multi-currency sub-account lifecycle, account tiers, beneficiary/payee management, account closure |
| Tech Stack | Go, PostgreSQL |
| Data Ownership | User profiles, accounts, sub-accounts, IBANs, tiers, beneficiaries |
| API Style | REST (public), gRPC (internal) |
| Key Patterns | IBAN generation using Lithuanian format (LT + check digits + SWIFT + account number), tier-based feature flags |

**Rationale:** Core domain entity. Separated from Ledger because account lifecycle (opening, closing, KYC status) is a different concern from financial accounting. Account data is queried frequently and needs different caching strategies than ledger entries.

### 3.3 Ledger Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | Double-entry general ledger, posting journal entries, balance calculation, interest accrual, fee application, reconciliation, event sourcing for all financial state changes |
| Tech Stack | Go, PostgreSQL (with event store), Redis (balance cache) |
| Data Ownership | Chart of accounts, journal entries, ledger balances, event store, reconciliation records |
| API Style | gRPC (internal only -- never directly exposed to clients) |
| Key Patterns | Event sourcing (append-only journal), CQRS (write to event store, project to materialized balance views), optimistic locking on balance updates, idempotency via posting IDs |

**Rationale:** The Ledger is the core of the CBS. It must guarantee ACID properties for every posting. Event sourcing provides an immutable audit trail required by regulators and enables point-in-time balance reconstruction. CQRS separates the high-read balance queries from the write-heavy posting path. This service is internal-only to ensure no external actor can directly modify the ledger.

### 3.4 Payment Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | SEPA Credit Transfer (SCT), SEPA Instant (SCT Inst), internal transfers, SEPA Direct Debit management, recurring/scheduled payments, currency exchange, payment routing |
| Tech Stack | Go, PostgreSQL, Kafka (outbox pattern) |
| Data Ownership | Payment orders, payment routing rules, FX rates, recurring payment schedules, SDD mandates |
| API Style | REST (public), gRPC (internal) |
| Key Patterns | Saga pattern for multi-step payments (debit -> route -> settle -> credit), outbox pattern for reliable event publishing, idempotency keys, FX rate locking (30s window) |

**Rationale:** Payments are the highest-value flows. The saga pattern handles the distributed transaction nature of SEPA payments (local debit, external settlement, confirmation). Separated from Ledger because payment orchestration (routing, retries, external gateway communication) is a different concern from accounting.

### 3.5 Card Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | Virtual/physical card issuance, card lifecycle (activate, freeze, block, replace), PIN management, spending controls, 3D Secure authentication, Apple Pay / Google Pay tokenization, authorization processing, dispute management |
| Tech Stack | Go, PostgreSQL, Redis (authorization cache) |
| Data Ownership | Cards, card controls, 3DS challenges, disputes, tokenization records |
| API Style | REST (public), gRPC (internal), webhooks (from card processor) |
| Key Patterns | Card processor adapter pattern (abstract Enfuce/Marqeta behind interface), real-time authorization with sub-100ms response, PCI DSS scope minimization via tokenization |

**Rationale:** Card operations require real-time authorization decisions (sub-100ms). The adapter pattern allows switching card processors without affecting the rest of the system. PIN and card data are never stored in TeslaPay systems -- they remain with the PCI-certified processor. The Card Service only stores tokenized references.

### 3.6 Crypto Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | Fuse Smart Wallet creation and management, crypto balance queries, buy/sell (on-ramp/off-ramp), send/receive tokens, gasless transaction relay via ERC-4337, blockchain event monitoring, Solid soUSD integration (Phase 2) |
| Tech Stack | TypeScript/Node.js (NestJS), PostgreSQL, FuseBox SDK (Dart on mobile, TypeScript on backend) |
| Data Ownership | Wallet metadata (addresses, not keys), crypto orders, blockchain transaction records, price feeds |
| API Style | REST (public), gRPC (internal), WebSocket (price feeds) |
| Key Patterns | FuseBox Bundler for UserOperations, Paymaster for gas sponsorship, circuit breaker for blockchain RPC, on-ramp/off-ramp via Ledger Service integration |

**Rationale:** TypeScript/Node.js chosen because FuseBox backend SDK is TypeScript-native, and the Fuse ecosystem tooling is JavaScript/TypeScript. Running this in Go would require maintaining custom RPC bindings. The circuit breaker pattern protects the platform from Fuse network instability (PRD risk R03). Wallet keys are never stored server-side; the Smart Wallet is controlled by the user's device-derived EOA.

### 3.7 KYC Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | KYC/AML orchestration with Sumsub, verification status management, risk scoring, ongoing AML monitoring, manual review workflow, re-verification triggers, geographic restrictions |
| Tech Stack | Go, PostgreSQL |
| Data Ownership | Verification records, risk scores, AML screening results, review decisions (actual documents stored in Sumsub) |
| API Style | REST (public -- limited to initiating verification), gRPC (internal), webhooks (from Sumsub) |
| Key Patterns | State machine for verification workflow (initiated -> document_submitted -> liveness_passed -> aml_screened -> approved/rejected), webhook signature validation, retry with exponential backoff |

**Rationale:** KYC is a compliance-critical service. Documents and biometric data are stored by Sumsub (PCI-equivalent certified), not in TeslaPay systems, minimizing data liability. The state machine pattern ensures every verification follows the correct sequence and no step is skipped. Ongoing AML monitoring events arrive via Sumsub webhooks and are processed asynchronously.

### 3.8 Notification Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | Push notifications (APNs, FCM), email (transactional and marketing), SMS (OTP, alerts), in-app notifications, notification preferences, delivery tracking |
| Tech Stack | Go, PostgreSQL, Redis (deduplication), Kafka consumer |
| Data Ownership | Notification templates, delivery logs, user preferences, device tokens |
| API Style | gRPC (internal only -- notifications are triggered by events, not direct API calls) |
| Key Patterns | Fan-out consumer from Kafka topics, template engine with i18n, delivery channel selection based on user preferences, deduplication via Redis |

**Rationale:** Notification is a cross-cutting concern consumed by every other service. Making it event-driven via Kafka decouples producers from delivery mechanics. Push notification delivery must meet the 3-second SLA for card transaction alerts (MC-007).

### 3.9 Audit Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | Immutable audit trail, compliance event logging, regulatory report generation (Bank of Lithuania, CRS/FATCA), SAR filing support, GDPR data subject request processing, data retention enforcement |
| Tech Stack | Go, PostgreSQL (append-only tables), Elasticsearch (search/analysis) |
| Data Ownership | Audit events, compliance reports, data retention schedules |
| API Style | gRPC (internal), REST (admin dashboard) |
| Key Patterns | Append-only storage with cryptographic chaining (hash of previous entry), 7-year retention, Elasticsearch for full-text search across audit entries |

**Rationale:** Regulatory requirement. Every state change in the system produces an audit event consumed by this service. Cryptographic chaining ensures tamper evidence. Elasticsearch enables compliance officers to search and analyze audit trails efficiently. Separated from other services to provide a single, trusted source of truth for regulators.

### 3.10 Fraud Detection Service

| Attribute | Detail |
|-----------|--------|
| Responsibility | Real-time transaction risk scoring, velocity checks, rule-based fraud detection, ML-based anomaly detection (Phase 2), SAR alert generation, geographic anomaly detection |
| Tech Stack | Go (rules engine), Python (ML models Phase 2), Redis (velocity counters), Kafka consumer |
| Data Ownership | Fraud rules, risk scores, velocity counters, fraud cases |
| API Style | gRPC (synchronous scoring for authorization), Kafka consumer (async analysis) |
| Key Patterns | Two-phase evaluation: (1) synchronous lightweight rules during authorization (<10ms), (2) asynchronous deep analysis post-transaction. Velocity windows via Redis sorted sets |

**Rationale:** Must operate in two modes -- synchronous for real-time card authorization decisions and asynchronous for transaction pattern analysis. The rules engine handles Phase 1 requirements (structuring, high-risk corridors, threshold alerts). ML models in Phase 2 provide adaptive detection.

---

## 4. Communication Patterns

### 4.1 Synchronous Communication (gRPC)

Used for operations that require an immediate response:

| Caller | Callee | Use Case |
|--------|--------|----------|
| API Gateway | Auth Service | Token validation on every request |
| Payment Service | Ledger Service | Post debit/credit entries |
| Payment Service | Account Service | Validate account status and limits |
| Card Service | Fraud Detection | Real-time authorization scoring |
| Card Service | Ledger Service | Post card transaction entries |
| Crypto Service | Ledger Service | Post fiat leg of buy/sell |
| Account Service | KYC Service | Check verification status |

**Protocol:** gRPC with Protocol Buffers, mTLS between all services, circuit breakers (Hystrix pattern), retry with exponential backoff, deadline propagation.

### 4.2 Asynchronous Communication (Kafka)

Used for event-driven workflows where eventual consistency is acceptable:

| Topic | Producers | Consumers | Events |
|-------|-----------|-----------|--------|
| `ledger.events` | Ledger Service | Audit, Notification, Fraud Detection | `entry.posted`, `balance.updated`, `reconciliation.completed` |
| `payment.events` | Payment Service | Notification, Audit, Account | `payment.initiated`, `payment.completed`, `payment.failed`, `payment.scheduled` |
| `card.events` | Card Service | Notification, Audit, Fraud Detection | `card.issued`, `card.frozen`, `authorization.approved`, `authorization.declined` |
| `kyc.events` | KYC Service | Account, Notification, Audit, Crypto | `verification.completed`, `aml.alert`, `risk.score.updated` |
| `crypto.events` | Crypto Service | Notification, Audit, Fraud Detection | `wallet.created`, `crypto.bought`, `crypto.sold`, `transfer.sent` |
| `audit.events` | All Services | Audit Service | `state.changed` (generic envelope) |
| `notification.commands` | All Services | Notification Service | `send.push`, `send.email`, `send.sms` |
| `fraud.signals` | Fraud Detection | Card, Payment, Account | `risk.elevated`, `account.freeze.recommended` |

**Kafka Configuration:** 3 brokers minimum, replication factor 3, `acks=all` for financial topics, exactly-once semantics (EOS) enabled, 30-day retention for financial topics (90 days for audit).

### 4.3 Communication Pattern Decision Matrix

```
+---------------------+-------------------+------------------------------------------+
| Pattern             | When to Use       | Example                                  |
+---------------------+-------------------+------------------------------------------+
| Synchronous (gRPC)  | Real-time needed, | Balance check before payment,            |
|                     | strong consistency | card authorization scoring                |
+---------------------+-------------------+------------------------------------------+
| Async (Kafka event) | Eventual consistency| Send notification after payment,         |
|                     | acceptable, fan-out | update audit trail                       |
+---------------------+-------------------+------------------------------------------+
| Saga (orchestrated) | Multi-step business| SEPA payment: validate -> hold ->         |
|                     | transaction        | submit -> settle -> release              |
+---------------------+-------------------+------------------------------------------+
| Outbox Pattern      | Reliable event     | Payment Service writes to DB + outbox     |
|                     | publishing needed  | table atomically; relay publishes to Kafka|
+---------------------+-------------------+------------------------------------------+
```

---

## 5. Infrastructure Architecture

### 5.1 AWS Deployment (eu-central-1 Frankfurt, DR in eu-west-1 Ireland)

```
+==============================================================================+
|                          AWS eu-central-1 (Frankfurt) -- Primary             |
|                                                                              |
|  +---------------------------+  +---------------------------+                |
|  |      AZ-a                 |  |      AZ-b                 |               |
|  |                           |  |                           |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |  |  EKS Worker Nodes   |  |  |  |  EKS Worker Nodes   |  |               |
|  |  |  (m6i.2xlarge)      |  |  |  |  (m6i.2xlarge)      |  |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |                           |  |                           |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |  |  RDS PostgreSQL     |  |  |  |  RDS PostgreSQL     |  |               |
|  |  |  (Primary)          |  |  |  |  (Standby)          |  |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |                           |  |                           |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |  |  MSK Kafka Broker   |  |  |  |  MSK Kafka Broker   |  |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |                           |  |                           |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  |  |  ElastiCache Redis  |  |  |  |  ElastiCache Redis  |  |               |
|  |  |  (Primary)          |  |  |  |  (Replica)          |  |               |
|  |  +---------------------+  |  |  +---------------------+  |               |
|  +---------------------------+  +---------------------------+                |
|                                                                              |
|  +---------------------------+                                               |
|  |      AZ-c                 |                                               |
|  |  EKS Workers + MSK Broker |                                               |
|  |  + Redis Replica          |                                               |
|  +---------------------------+                                               |
|                                                                              |
|  Shared Services:                                                            |
|  +--------------------+  +------------------+  +------------------+          |
|  | ECR (Container     |  | S3 (Documents,   |  | CloudWatch +     |          |
|  |  Registry)         |  |  Backups)         |  |  X-Ray (Tracing) |          |
|  +--------------------+  +------------------+  +------------------+          |
|                                                                              |
|  +--------------------+  +------------------+  +------------------+          |
|  | Secrets Manager +  |  | KMS (Encryption  |  | OpenSearch       |          |
|  |  Vault (Secrets)   |  |  Keys)           |  | (Audit Logs)     |          |
|  +--------------------+  +------------------+  +------------------+          |
|                                                                              |
+==============================================================================+

+==============================================================================+
|                          AWS eu-west-1 (Ireland) -- DR                       |
|                                                                              |
|  +---------------------------+  +---------------------------+                |
|  |      AZ-a                 |  |      AZ-b                 |               |
|  |  EKS Workers (standby)    |  |  RDS Read Replica          |               |
|  |  S3 Cross-Region Repl.    |  |  MSK Mirror (async)        |               |
|  +---------------------------+  +---------------------------+                |
+==============================================================================+
```

### 5.2 Kubernetes Cluster Organization

```
Namespaces:
  teslapay-core        -- Auth, Account, Ledger, Payment services
  teslapay-card        -- Card Service (PCI DSS scoped)
  teslapay-crypto      -- Crypto Service
  teslapay-compliance  -- KYC, Audit, Fraud Detection services
  teslapay-infra       -- Notification, API Gateway, observability
  teslapay-data        -- Kafka Connect, schema registry, database operators
```

**Network Policies:** Each namespace has strict ingress/egress rules. The `teslapay-card` namespace has the most restrictive policies (PCI DSS CDE -- Cardholder Data Environment). Only the Card Service can communicate with the card processor, and only the API Gateway can reach the Card Service.

### 5.3 Scaling Strategy

| Service | Min Pods | Max Pods | Scaling Trigger | Notes |
|---------|----------|----------|-----------------|-------|
| API Gateway | 3 | 20 | CPU > 60%, RPS > 5000 | Front-door; absorbs traffic spikes |
| Auth Service | 3 | 15 | CPU > 70%, RPS > 3000 | Login spike handling |
| Account Service | 2 | 10 | CPU > 70% | Moderate, mostly reads |
| Ledger Service | 3 | 12 | CPU > 60%, Kafka lag > 1000 | Critical path; no scale-to-zero |
| Payment Service | 3 | 15 | CPU > 60%, queue depth | SEPA Instant requires always-on |
| Card Service | 3 | 15 | CPU > 60%, latency p95 > 50ms | Authorization latency critical |
| Crypto Service | 2 | 8 | CPU > 70% | Lower initial traffic |
| KYC Service | 2 | 8 | Pending verifications > 100 | Burst during marketing campaigns |
| Notification Service | 2 | 10 | Kafka consumer lag > 5000 | Bursty by nature |
| Audit Service | 2 | 6 | Kafka consumer lag > 10000 | Can tolerate slight lag |
| Fraud Detection | 3 | 10 | CPU > 60%, latency p95 > 5ms | Sync path must be fast |

---

## 6. Data Architecture

### 6.1 Database-per-Service Pattern

Each microservice owns its database. No direct cross-database queries. Data is shared via APIs and events.

| Service | Database | Type | Rationale |
|---------|----------|------|-----------|
| Auth | `auth_db` | PostgreSQL 16 | Credential storage, session management |
| Account | `account_db` | PostgreSQL 16 | User profiles, account data |
| Ledger | `ledger_db` | PostgreSQL 16 | Double-entry journal, event store |
| Payment | `payment_db` | PostgreSQL 16 | Payment orders, schedules |
| Card | `card_db` | PostgreSQL 16 | Card metadata (no PAN/CVV -- stored at processor) |
| Crypto | `crypto_db` | PostgreSQL 16 | Wallet metadata, crypto orders |
| KYC | `kyc_db` | PostgreSQL 16 | Verification records, risk scores |
| Notification | `notification_db` | PostgreSQL 16 | Templates, delivery logs, preferences |
| Audit | `audit_db` | PostgreSQL 16 + OpenSearch | Append-only events, indexed for search |
| Fraud | `fraud_db` | PostgreSQL 16 + Redis | Rules, velocity counters, cases |

### 6.2 Caching Strategy

| Cache Layer | Technology | Use Case | TTL |
|-------------|-----------|----------|-----|
| Session cache | Redis Cluster | Active sessions, device fingerprints | 15 min (access), 30 days (refresh) |
| Balance cache | Redis Cluster | Projected balances for quick reads | 5 sec (invalidated on ledger event) |
| FX rate cache | Redis | Current exchange rates | 30 sec |
| Authorization cache | Redis | Recent authorizations for dedup | 5 min |
| Account cache | Redis | Account status, tier, limits | 60 sec |

### 6.3 Event Store Design

The Ledger Service uses event sourcing. Every financial state change is an immutable event:

```
Event Store Entry:
{
  event_id:        UUID (globally unique)
  aggregate_id:    UUID (account or posting ID)
  aggregate_type:  "account" | "posting" | "journal"
  event_type:      "debit.posted" | "credit.posted" | "reversal.posted"
  event_data:      JSONB (full event payload)
  metadata:        JSONB (correlation_id, causation_id, actor)
  sequence_number: BIGINT (per-aggregate ordering)
  created_at:      TIMESTAMPTZ
  checksum:        TEXT (SHA-256 of event_data + previous checksum)
}
```

---

## 7. Cross-Cutting Concerns

### 7.1 Observability Stack

| Concern | Tool | Detail |
|---------|------|--------|
| Metrics | Prometheus + Grafana | RED metrics (Rate, Errors, Duration) per service, business metrics (TPS, payment volume) |
| Logging | Fluent Bit -> OpenSearch | Structured JSON logs, correlation IDs, PII redaction |
| Tracing | AWS X-Ray / OpenTelemetry | Distributed trace across all services, 100% sampling for errors, 10% for normal traffic |
| Alerting | Grafana Alertmanager -> PagerDuty | SLA-based alerts: p95 latency, error rate, Kafka lag, balance discrepancy |
| Dashboards | Grafana | Service health, business KPIs, compliance dashboards |

### 7.2 Service Mesh

Istio service mesh provides:
- mTLS between all services (zero-trust networking)
- Traffic management (canary deployments, circuit breaking)
- Observability (automatic metrics, tracing)
- Rate limiting at service level
- Retry policies with configurable backoff

### 7.3 API Gateway (Kong)

| Feature | Configuration |
|---------|--------------|
| Authentication | JWT validation, API key for admin |
| Rate Limiting | Per-user: 100 req/min (Standard), 300 req/min (Premium); global: 10,000 req/min |
| Request Transformation | Header injection (correlation ID, request ID) |
| Response Caching | GET endpoints with 5-10 sec TTL |
| IP Allowlisting | Admin API restricted to VPN IPs |
| Request Size Limit | 10MB (document uploads), 1MB (standard) |
| Logging | Full request/response logging (PII redacted) |

### 7.4 CI/CD Pipeline

```
Developer -> GitHub PR -> Lint + Unit Tests -> Code Review (2 approvals)
                                                    |
                                                    v
                              Integration Tests (Testcontainers)
                                                    |
                                                    v
                              Security Scan (SAST: SonarQube, DAST: OWASP ZAP)
                                                    |
                                                    v
                              Container Build -> ECR Push -> Trivy Scan
                                                    |
                                                    v
                              Deploy to Staging (ArgoCD) -> E2E Tests
                                                    |
                                                    v
                              Manual Approval Gate (for production)
                                                    |
                                                    v
                              Deploy to Production (ArgoCD, canary 10% -> 50% -> 100%)
                                                    |
                                                    v
                              Post-Deploy Smoke Tests -> Monitor for 30 min
```

---

## 8. Failure Modes and Resilience

### 8.1 Failure Scenarios

| Scenario | Impact | Mitigation |
|----------|--------|------------|
| Kafka broker failure | Event delivery delayed | 3-broker cluster, replication factor 3; ISR min 2 |
| PostgreSQL primary failure | Write unavailable | Multi-AZ RDS with automatic failover (< 60s) |
| Card processor timeout | Authorization delayed | 3-second timeout, decline with retry code |
| Sumsub outage | New onboarding blocked | Queue verification requests, retry when available |
| Fuse.io network down | Crypto operations unavailable | Circuit breaker, graceful degradation, fiat operations unaffected |
| Redis cluster failure | Cache miss, slower reads | Fall through to database; services remain functional |
| API Gateway failure | All external traffic blocked | 3+ replicas, NLB health checks, automatic replacement |
| Full AZ failure | Reduced capacity | Multi-AZ deployment, pods rescheduled to surviving AZs |
| Full region failure | Full outage | DR failover to eu-west-1, RTO < 15 min, RPO < 1 min |

### 8.2 Circuit Breaker Configuration

| External System | Open Threshold | Half-Open After | Close Threshold |
|-----------------|---------------|-----------------|-----------------|
| Sumsub API | 5 failures in 30s | 60 seconds | 3 successes |
| Card Processor | 3 failures in 10s | 30 seconds | 5 successes |
| Fuse.io RPC | 5 failures in 30s | 60 seconds | 3 successes |
| SEPA Gateway | 3 failures in 30s | 60 seconds | 5 successes |
| FX Rate Provider | 5 failures in 60s | 120 seconds | 3 successes |

---

## 9. Data Flow Diagrams

### 9.1 SEPA Payment Flow

```
User -> API GW -> Payment Service: Initiate SCT
                       |
                       v
               Account Service: Validate sender account, check limits
                       |
                       v
               Fraud Detection: Score transaction risk (sync, <10ms)
                       |
                       v
               Ledger Service: Post debit entry (hold)
                       |
                       v
               Payment Service: Submit to SEPA Gateway (Banking Circle)
                       |
                       v
               [Kafka: payment.events -> payment.initiated]
                       |
               Notification Service: Push "Payment sent" to user
                       |
               ... (T+0 to T+1 settlement) ...
                       |
               SEPA Gateway Webhook: Settlement confirmed
                       |
                       v
               Payment Service: Mark as settled
                       |
                       v
               Ledger Service: Release hold, post final debit
                       |
                       v
               [Kafka: payment.events -> payment.completed]
                       |
               Notification Service: Push "Payment completed" to user
               Audit Service: Log settlement event
```

### 9.2 Card Authorization Flow

```
Merchant POS -> Mastercard Network -> Card Processor (Enfuce)
                                            |
                                            v
                                     Webhook to Card Service (<100ms budget)
                                            |
                                            v
                                     Card Service: Validate card status, controls
                                            |
                                            v
                                     Fraud Detection: Score (sync, <10ms)
                                            |
                                            v
                                     Account Service: Check balance
                                            |
                                            v
                                     Ledger Service: Post authorization hold
                                            |
                                            v
                                     Card Service -> Card Processor: Approve/Decline
                                            |
                                     [Kafka: card.events -> authorization.approved]
                                            |
                                     Notification Service: Push to user (<3s total)
```

### 9.3 KYC Onboarding Flow

```
User (Mobile) -> Sumsub SDK (in-app): Capture document + liveness
                       |
                       v
                 Sumsub Cloud: Process verification
                       |
                       v
                 Sumsub Webhook -> KYC Service: applicantReviewed
                       |
                       v
                 KYC Service: Parse result (GREEN/RED)
                       |
                       +--- GREEN (approved) ---> Account Service: Activate account
                       |                                    |
                       |                                    v
                       |                          Crypto Service: Create Fuse Smart Wallet
                       |                                    |
                       |                                    v
                       |                          Notification: "Welcome to TeslaPay"
                       |
                       +--- RED (RETRY) -------> Notification: "Please resubmit docs"
                       |
                       +--- RED (FINAL) -------> Account Service: Block account
                                                 Notification: "Verification failed"
```

---

## 10. Technology Decision Records

### TDR-001: Custom CBS vs. Licensed Core Banking Platform

**Decision:** Build custom CBS (Ledger + Payment + Account services).

**Rationale:** Licensed platforms (Mambu, Thought Machine) cost EUR 50K-200K/month at scale and impose constraints on crypto integration. A custom CBS built with event sourcing provides full control over the ledger, which is the core differentiator when bridging traditional finance and blockchain. The team can optimize for TeslaPay's specific multi-currency + crypto requirements.

**Trade-off:** Higher initial development cost (estimated 4-6 developer-months) but lower long-term TCO and no vendor lock-in.

### TDR-002: Enfuce as Primary Card Processor

**Decision:** Use Enfuce as the Mastercard issuer processor with BIN sponsorship.

**Rationale:** Enfuce is a Finnish EMI with Mastercard principal membership, offering BIN sponsorship for EEA card programs. They handle PCI DSS certification, Mastercard scheme compliance, and card personalization. This avoids TeslaPay needing its own Mastercard principal membership (which requires significant capital and certification). Enfuce's API-first approach aligns with our microservices architecture.

**Trade-off:** Dependency on Enfuce for card operations. Mitigation: adapter pattern allows switching to Marqeta or GPS if needed.

### TDR-003: Banking Circle for SEPA Connectivity

**Decision:** Use Banking Circle for SEPA SCT, SCT Inst, and SDD connectivity.

**Rationale:** Banking Circle provides plug-and-play SEPA connectivity via API, supporting both standard and instant payments. As a licensed bank, they can provide TeslaPay with indirect SEPA scheme access without TeslaPay needing direct EBA STEP2 participation. They support both API and SWIFT connectivity.

**Trade-off:** Per-transaction fees reduce margin. At scale (>500K transactions/month), direct SEPA participation may become more economical.

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CTO | TBD | | Pending |
| Engineering Lead | TBD | | Pending |
| Principal Architect | Dream Team Architect | 2026-03-03 | Submitted |
