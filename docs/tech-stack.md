# TeslaPay Technology Stack

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Principal Software Architect, Dream Team

---

## 1. Technology Decisions Summary

| Layer | Technology | Version | Rationale |
|-------|-----------|---------|-----------|
| **Backend (Core)** | Go | 1.22+ | Performance, concurrency, type safety, fintech industry standard |
| **Backend (Crypto)** | TypeScript / NestJS | Node 20 LTS | FuseBox SDK is TypeScript-native |
| **Mobile** | Flutter / Dart | 3.x | Cross-platform, Fuse Wallet SDK is Dart-native, single codebase |
| **Admin Dashboard** | React / Next.js | 14+ | Rich ecosystem, SSR for compliance dashboards |
| **Primary Database** | PostgreSQL | 16 | ACID compliance, JSON support, proven in financial systems |
| **Cache** | Redis Cluster | 7.x | Sub-millisecond latency, pub/sub, sorted sets for velocity |
| **Message Broker** | Apache Kafka (AWS MSK) | 3.7+ | Event sourcing, exactly-once semantics, high throughput |
| **Search / Audit** | OpenSearch | 2.x | Full-text search on audit logs, regulatory queries |
| **Object Storage** | AWS S3 | -- | Document storage, backups, exports |
| **Container Orchestration** | Kubernetes (AWS EKS) | 1.29+ | Industry standard, auto-scaling, multi-AZ |
| **Service Mesh** | Istio | 1.21+ | mTLS, traffic management, observability |
| **API Gateway** | Kong | 3.x | Plugin ecosystem, rate limiting, JWT validation |
| **IaC** | Terraform + Terragrunt | 1.7+ | Reproducible infrastructure, multi-environment |
| **CI/CD** | GitHub Actions + ArgoCD | -- | GitOps, canary deployments |
| **Secrets** | HashiCorp Vault + AWS KMS | -- | Dynamic secrets, HSM-backed encryption keys |
| **Monitoring** | Prometheus + Grafana | -- | Metrics collection, dashboarding, alerting |
| **Logging** | Fluent Bit -> OpenSearch | -- | Structured log aggregation with PII redaction |
| **Tracing** | OpenTelemetry + AWS X-Ray | -- | Distributed tracing across microservices |
| **Container Registry** | AWS ECR | -- | Private registry, vulnerability scanning |
| **DNS / CDN** | AWS Route53 + CloudFront | -- | Low-latency global distribution, DDoS protection |
| **WAF** | AWS WAF | -- | OWASP Top 10, bot detection, geo-blocking |

---

## 2. Backend Language: Go

### Why Go

| Factor | Assessment |
|--------|------------|
| **Performance** | Compiled, garbage-collected with sub-millisecond pauses. Card authorization path requires <100ms response; Go routinely achieves <5ms for business logic. |
| **Concurrency** | Goroutines and channels handle thousands of concurrent SEPA and card operations efficiently. 100 TPS target is trivially achievable. |
| **Type Safety** | Static typing catches errors at compile time -- critical for financial calculations where runtime errors mean money loss. |
| **Operational Simplicity** | Single binary deployment. No JVM tuning, no dependency hell. Container images are <20MB. |
| **Fintech Adoption** | Used by Monzo, Revolut, Square, Stripe. Proven at neobank scale. |
| **Ecosystem** | Strong gRPC support (protobuf), excellent PostgreSQL drivers (pgx), Kafka libraries (confluent-kafka-go). |
| **Hiring** | Growing talent pool in EU. Simpler language attracts broader hiring than Kotlin/Scala. |

### Why Not Kotlin/JVM

Kotlin was considered as an alternative. JVM provides excellent libraries for financial computation (BigDecimal) and Spring Boot has a mature ecosystem. However, JVM cold start times (3-5s) are problematic for Kubernetes autoscaling, memory footprint is 3-5x higher than Go, and the team would need JVM operational expertise. Go's `math/big` and `shopspring/decimal` libraries provide sufficient precision for financial calculations.

### Why Not Rust

Rust offers superior performance and memory safety but has a steeper learning curve, smaller hiring pool, and slower development velocity. For TeslaPay's 9-month MVP timeline, Go's simplicity and ecosystem maturity are more appropriate.

### Crypto Service Exception: TypeScript/NestJS

The Crypto Service uses TypeScript because:
1. FuseBox SDK (`@fuseio/fusebox-web-sdk`) is TypeScript-native
2. Ethereum/EVM tooling (ethers.js, viem) is JavaScript/TypeScript-first
3. NestJS provides a structured framework suitable for production backend services
4. The Fuse team maintains TypeScript examples and documentation

This is an intentional, bounded exception. The Crypto Service communicates with other services via gRPC (using auto-generated TypeScript stubs from proto files), so the language boundary is invisible to callers.

---

## 3. Mobile: Flutter/Dart

### Why Flutter

| Factor | Assessment |
|--------|------------|
| **Fuse SDK** | `fuse_wallet_sdk` is a Dart package published on pub.dev. Native Dart integration eliminates bridging overhead. |
| **Sumsub SDK** | Sumsub provides a Flutter SDK for document capture and liveness detection. |
| **Cross-Platform** | Single codebase for iOS and Android. Team of 2-3 Flutter engineers vs. 4-6 for native (Swift + Kotlin). |
| **Performance** | Dart compiles to native ARM code. UI renders at 60fps. Meets the 3-second cold start requirement. |
| **Fintech Adoption** | Used by Nubank (150M+ users), Google Pay, Grab. Proven at neobank scale. |
| **UI Consistency** | Pixel-perfect control. TeslaPay brand identity rendered identically on both platforms. |
| **Time to Market** | 40-50% faster development than dual-native approach. Critical for 9-month MVP. |

### Why Not Native (Swift + Kotlin)

The PRD initially mentioned native development. However, given (a) Fuse SDK is Dart-native, (b) Sumsub has a Flutter SDK, (c) the 9-month timeline, and (d) the budget reality of a EUR 1.69M revenue company, Flutter is the pragmatic choice. Platform-specific features (Apple Pay provisioning, biometrics, NFC) are accessible via Flutter platform channels.

### Flutter Architecture

```
lib/
  core/
    di/             -- Dependency injection (get_it + injectable)
    network/        -- Dio HTTP client, interceptors, certificate pinning
    security/       -- Biometric auth, secure storage (flutter_secure_storage)
    l10n/           -- Internationalization (5+ languages)
  features/
    auth/           -- Login, registration, PIN, biometric
    dashboard/      -- Home screen, balance overview
    accounts/       -- Sub-accounts, IBAN display
    payments/       -- SEPA, internal transfers, FX
    cards/          -- Card display, controls, Apple/Google Pay
    crypto/         -- Fuse wallet, buy/sell, send/receive
    kyc/            -- Sumsub SDK integration
    settings/       -- Preferences, language, dark mode
    support/        -- In-app chat, FAQ
  shared/
    widgets/        -- Reusable UI components
    models/         -- Shared domain models
    utils/          -- Formatters, validators
```

**State Management:** BLoC pattern (flutter_bloc). Chosen for its testability, separation of UI from business logic, and suitability for event-driven architectures (mirrors the backend event model).

---

## 4. Database: PostgreSQL 16

### Why PostgreSQL

| Factor | Assessment |
|--------|------------|
| **ACID Compliance** | Mandatory for double-entry bookkeeping. Every journal entry posting must be atomic. |
| **Financial Precision** | `NUMERIC(19,4)` type provides exact decimal arithmetic -- no floating-point errors in money calculations. |
| **JSON Support** | `JSONB` columns for event store payloads, flexible metadata, and audit event data. |
| **Partitioning** | Native table partitioning by date range for ledger_entries (millions of rows per month). |
| **Row-Level Security** | RLS policies enforce data isolation between services and tenants if multi-tenancy is added later. |
| **Maturity** | 35+ years of development. Used by Goldman Sachs, Revolut, Wise for core banking. |
| **AWS RDS** | Managed service with Multi-AZ, automated backups, point-in-time recovery. |
| **Logical Replication** | Enables CDC (Change Data Capture) for event-driven patterns without application changes. |

### PostgreSQL Configuration for Financial Workloads

```
-- Key settings for RDS PostgreSQL 16
max_connections = 200
shared_buffers = 8GB            -- 25% of instance RAM
effective_cache_size = 24GB     -- 75% of instance RAM
work_mem = 64MB
maintenance_work_mem = 2GB
wal_level = logical             -- For CDC / Debezium
max_wal_size = 4GB
checkpoint_timeout = 15min
synchronous_commit = on         -- Mandatory for financial data
default_transaction_isolation = 'read committed'
```

### Database-per-Service Sizing

| Database | Instance Class | Storage | Rationale |
|----------|---------------|---------|-----------|
| `ledger_db` | db.r6g.2xlarge | 500GB gp3 | Highest write volume, event store |
| `account_db` | db.r6g.xlarge | 100GB gp3 | Moderate reads, low writes |
| `payment_db` | db.r6g.xlarge | 200GB gp3 | Payment orders, scheduling |
| `card_db` | db.r6g.xlarge | 100GB gp3 | Card metadata, authorizations |
| `auth_db` | db.r6g.large | 50GB gp3 | Sessions, credentials |
| `crypto_db` | db.r6g.large | 100GB gp3 | Wallet metadata, crypto orders |
| `kyc_db` | db.r6g.large | 50GB gp3 | Verification records |
| `notification_db` | db.r6g.large | 100GB gp3 | Delivery logs |
| `audit_db` | db.r6g.xlarge | 1TB gp3 | Append-only, 7-year retention |
| `fraud_db` | db.r6g.large | 50GB gp3 | Rules, cases |

---

## 5. Message Broker: Apache Kafka (AWS MSK)

### Why Kafka

| Factor | Assessment |
|--------|------------|
| **Event Sourcing** | Kafka's immutable, append-only log is the natural storage for financial events. |
| **Exactly-Once Semantics** | Kafka EOS prevents duplicate financial event processing -- critical for double-entry posting. |
| **High Throughput** | Handles 100K+ events/second. TeslaPay's 100 TPS is well within capacity. |
| **Durability** | Replication factor 3, `acks=all` ensures no event loss. |
| **Consumer Groups** | Multiple services consume the same events independently (notification + audit + fraud). |
| **Compacted Topics** | Used for account state projection -- latest balance always available. |
| **Ecosystem** | Kafka Connect for CDC (Debezium), Schema Registry for event schema evolution. |

### Why Not RabbitMQ

RabbitMQ excels at task queues but lacks Kafka's log semantics needed for event sourcing and replay. Kafka's ability to replay events from any offset is essential for rebuilding read models and auditing.

### Why Not AWS SQS/SNS

SQS/SNS has simpler operations but lacks exactly-once delivery, event replay, and the rich consumer group semantics needed for CQRS projections.

### MSK Configuration

```
Cluster:
  Brokers: 3 (one per AZ)
  Instance: kafka.m5.2xlarge
  Storage: 1TB per broker (auto-expanding)
  Replication Factor: 3
  Min In-Sync Replicas: 2
  Retention:
    Financial topics: 30 days
    Audit topics: 90 days
    Compacted topics: infinite (state)
  Encryption: TLS in transit, KMS at rest
  Authentication: SASL/SCRAM
```

---

## 6. Cache: Redis Cluster 7.x (AWS ElastiCache)

### Use Cases

| Use Case | Data Structure | TTL | Justification |
|----------|---------------|-----|---------------|
| Session tokens | STRING | 15 min / 30 days | Fast session validation on every API call |
| Balance projections | HASH | 5 sec | Avoid hitting PostgreSQL for balance reads |
| FX rates | HASH | 30 sec | Frequently queried, updated from FX provider |
| Rate limiting counters | Sorted Set | 1 min window | Per-user request counting |
| Velocity checks | Sorted Set | 24h window | Transaction velocity (fraud detection) |
| Idempotency keys | STRING | 24h | Prevent duplicate payment processing |
| OTP codes | STRING | 5 min | SMS/email verification codes |
| Authorization dedup | STRING | 5 min | Prevent duplicate card authorization processing |

### Configuration

```
Cluster Mode: Enabled
Node Type: cache.r6g.xlarge
Shards: 3 (with replicas)
Encryption: In-transit (TLS), At-rest (KMS)
Eviction Policy: volatile-lru
Max Memory: 13GB per node
Persistence: None (cache-only; all data reconstructible from PostgreSQL)
```

---

## 7. Search and Analytics: OpenSearch 2.x

### Why OpenSearch (not Elasticsearch)

AWS OpenSearch Service is the managed offering. It provides:
- Full-text search across audit logs (compliance requirement ACC-004)
- Log aggregation and analysis
- No Elastic license concerns (OpenSearch is Apache 2.0)

### Use Cases

| Index | Source | Purpose |
|-------|--------|---------|
| `audit-events-*` | Audit Service | Regulatory audit trail search (7-year retention) |
| `application-logs-*` | Fluent Bit | Centralized application log analysis |
| `transaction-search-*` | Ledger Service CDC | Transaction history search with complex filters |
| `compliance-alerts-*` | KYC/Fraud Services | AML alert investigation |

---

## 8. API Gateway: Kong 3.x

### Why Kong

| Factor | Assessment |
|--------|------------|
| **Plugin Ecosystem** | JWT auth, rate limiting, request transformation, logging -- all built-in. |
| **Performance** | Nginx-based, handles 50K+ requests/second per instance. |
| **Kubernetes Native** | Kong Ingress Controller integrates with EKS. |
| **Open Source** | Kong Gateway OSS covers all TeslaPay requirements. |
| **PSD2 Compliance** | Custom plugins for SCA step-up enforcement. |

### Why Not AWS API Gateway

AWS API Gateway has a 29-second timeout limit and per-request pricing that becomes expensive at scale. Kong runs on our own infrastructure with predictable costs and full control over configuration.

### Why Not Envoy (standalone)

Envoy is excellent as a sidecar proxy (used via Istio) but lacks the API management features (developer portal, analytics) that Kong provides for external-facing APIs.

---

## 9. Infrastructure and DevOps

### 9.1 Container Orchestration: AWS EKS

| Decision | Detail |
|----------|--------|
| Managed Kubernetes | EKS reduces operational burden vs. self-managed k8s |
| Node Groups | Managed node groups with m6i.2xlarge instances (8 vCPU, 32GB RAM) |
| Cluster Autoscaler | Karpenter for fast node provisioning |
| Pod Disruption Budgets | Minimum 2 pods for all financial services |

### 9.2 Infrastructure as Code

| Tool | Purpose |
|------|---------|
| Terraform | AWS resource provisioning (VPC, EKS, RDS, MSK, ElastiCache) |
| Terragrunt | Multi-environment orchestration (dev, staging, production) |
| Helm | Kubernetes application packaging |
| ArgoCD | GitOps-based continuous delivery |
| Kustomize | Environment-specific configuration overlays |

### 9.3 Environments

| Environment | Purpose | Infrastructure |
|-------------|---------|---------------|
| `dev` | Developer integration testing | Single-AZ, small instances, shared databases |
| `staging` | Pre-production, E2E testing | Multi-AZ, production-like topology, test data |
| `production` | Live traffic | Multi-AZ (3 AZs), full HA, real data |
| `dr` | Disaster recovery | eu-west-1, warm standby, replicated data |

---

## 10. Security Tooling

| Tool | Purpose | Rationale |
|------|---------|-----------|
| HashiCorp Vault | Dynamic secrets, PKI, transit encryption | Industry standard for secrets management in financial services |
| AWS KMS | Encryption key management, envelope encryption | FIPS 140-2 Level 3 HSM backing for key storage |
| AWS CloudHSM | PCI DSS key management for card operations | Required for storing card encryption keys in FIPS 140-2 Level 3 HSM |
| Trivy | Container image vulnerability scanning | Shift-left security in CI/CD pipeline |
| SonarQube | Static Application Security Testing (SAST) | Code quality and vulnerability detection |
| OWASP ZAP | Dynamic Application Security Testing (DAST) | Runtime vulnerability scanning |
| Falco | Runtime security monitoring | Detect anomalous container behavior |
| AWS GuardDuty | Threat detection | Network and account-level threat detection |
| AWS Security Hub | Centralized security findings | Aggregate findings from all security tools |

---

## 11. Observability Stack Detail

### Prometheus + Grafana

```
Metrics Collection:
  - Service metrics: go_* (runtime), grpc_* (gRPC), http_* (REST)
  - Business metrics: teslapay_payment_total, teslapay_card_authorization_total
  - Infrastructure metrics: node_*, container_*, kubelet_*

Key Dashboards:
  1. Platform Overview: All services health, error rates, latency
  2. Payment Operations: SEPA volume, success rate, latency by type
  3. Card Operations: Authorization rate, decline reasons, 3DS success
  4. Compliance: KYC approval rate, AML alerts, audit volume
  5. SLA Monitoring: p95 latency vs. targets, availability %
```

### Alert Rules (Critical)

| Alert | Condition | Severity | Response |
|-------|-----------|----------|----------|
| HighErrorRate | Error rate > 5% for 5 min | Critical | Page on-call engineer |
| HighLatency | p95 > 500ms for 5 min | Warning | Investigate, scale if needed |
| LedgerImbalance | Debit != Credit sum | Critical | Immediate investigation, potential freeze |
| KafkaConsumerLag | Lag > 10000 for 10 min | Warning | Scale consumers |
| DatabaseConnectionPool | Usage > 80% | Warning | Check for connection leaks |
| CertificateExpiry | Expiry < 30 days | Warning | Rotate certificate |
| VaultSealed | Vault sealed | Critical | Unseal immediately |

---

## 12. Development Standards

### Go Services

```
Project Structure (per service):
  cmd/
    server/         -- Main entrypoint
  internal/
    domain/         -- Domain models, business rules (no external dependencies)
    application/    -- Use cases / application services
    infrastructure/
      postgres/     -- Repository implementations
      kafka/        -- Event publisher/consumer
      grpc/         -- gRPC server/client implementations
      http/         -- REST handlers
    config/         -- Configuration loading
  api/
    proto/          -- gRPC protobuf definitions
    openapi/        -- OpenAPI 3.0 specs
  migrations/       -- SQL migration files (golang-migrate)
  tests/
    integration/    -- Integration tests (testcontainers)
    e2e/            -- End-to-end tests

Libraries:
  - HTTP Router: chi (lightweight, idiomatic)
  - gRPC: google.golang.org/grpc
  - PostgreSQL: jackc/pgx/v5
  - Kafka: confluent-kafka-go
  - Validation: go-playground/validator
  - Decimal: shopspring/decimal
  - Testing: testify, testcontainers-go
  - Logging: slog (stdlib structured logging)
  - Config: viper
  - Migrations: golang-migrate
```

### Code Quality Gates

| Gate | Tool | Threshold |
|------|------|-----------|
| Unit Test Coverage | Go test | >= 80% |
| Lint | golangci-lint | Zero warnings |
| Cyclomatic Complexity | golangci-lint (gocyclo) | <= 15 per function |
| Dependency Vulnerabilities | govulncheck | Zero critical/high |
| Container Vulnerabilities | Trivy | Zero critical |
| SAST | SonarQube | Quality Gate pass |

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CTO | TBD | | Pending |
| Engineering Lead | TBD | | Pending |
| Principal Architect | Dream Team Architect | 2026-03-03 | Submitted |
