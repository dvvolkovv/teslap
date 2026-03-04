# TeslaPay Sprint Plan -- Sprints 1-6

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Senior Technical Project Manager, Dream Team
**Sprint Duration:** 2 weeks
**Sprint Ceremonies:** Planning (Monday S1), Daily Standup (15 min), Review (Friday S2), Retro (Friday S2)

---

## Sprint 1: Infrastructure and Project Scaffolding

**Dates:** 2026-03-09 to 2026-03-20
**Sprint Goal:** Establish cloud infrastructure, CI/CD pipeline, service scaffolding, and shared libraries so that engineering can begin feature development from Sprint 2 onward.

### User Stories

No user stories in Sprint 1 -- this is a pure technical foundation sprint.

### Technical Tasks

| ID | Task | Owner | Size | Notes |
|----|------|-------|------|-------|
| T1.01 | Provision AWS VPC, subnets, security groups (Terraform) | DevOps | L | eu-central-1, 3 AZs, public/private subnets |
| T1.02 | Deploy EKS cluster with Karpenter autoscaler | DevOps | L | Kubernetes 1.29+, managed node groups |
| T1.03 | Deploy RDS PostgreSQL instances (dev environment) | DevOps | M | 10 databases per architecture spec |
| T1.04 | Deploy MSK Kafka cluster (dev environment) | DevOps | M | 3 brokers, topics per architecture spec |
| T1.05 | Deploy ElastiCache Redis cluster (dev environment) | DevOps | M | Cluster mode enabled |
| T1.06 | Set up GitHub organization, repos (monorepo or multi-repo decision), branch protection | DevOps | S | Monorepo recommended for initial velocity |
| T1.07 | Configure GitHub Actions CI pipeline (lint, test, build, scan) | DevOps | L | SonarQube, Trivy integrated |
| T1.08 | Set up ArgoCD for GitOps deployment to dev/staging | DevOps | M | Canary deployment config for production later |
| T1.09 | Scaffold Go service template (project structure, Makefile, Dockerfile) | Tech Lead | M | Per tech-stack.md structure |
| T1.10 | Define protobuf/gRPC contracts for Auth, Account, Ledger services | Tech Lead + Backend | L | Shared proto repo |
| T1.11 | Scaffold Flutter app project with BLoC architecture | Flutter Lead | M | Feature-based folder structure per tech-stack.md |
| T1.12 | Set up Istio service mesh on EKS | DevOps | M | mTLS, traffic management |
| T1.13 | Deploy Kong API Gateway | DevOps | M | JWT validation, rate limiting plugins |
| T1.14 | Set up HashiCorp Vault for secrets management | DevOps | M | Dynamic database credentials |
| T1.15 | Deploy Prometheus + Grafana observability stack | DevOps | M | Base dashboards, alerting to Slack initially |
| T1.16 | Deploy Fluent Bit -> OpenSearch for logging | DevOps | M | PII redaction filters |
| T1.17 | Initiate Sumsub contract and sandbox access | PM | S | Business dependency -- must start immediately |
| T1.18 | Initiate Banking Circle contract negotiation | PM | S | Business dependency -- critical path for Sprint 5 |
| T1.19 | Initiate Enfuce card processor discussions | PM | S | Business dependency -- critical path for Sprint 9 |

### Definition of Done

- [ ] All dev environment infrastructure accessible by engineering team
- [ ] CI pipeline runs on every PR: lint, unit test, build, container scan
- [ ] At least one Go service (empty template) deploys successfully to dev via ArgoCD
- [ ] Flutter app builds and runs on iOS simulator and Android emulator
- [ ] All protobuf contracts for Sprint 2 services defined and code-generated
- [ ] Observability stack operational (logs, metrics, traces visible in Grafana)
- [ ] All external partner discussions initiated with documented timelines

### Sprint Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| AWS account setup delays (compliance approval) | Medium | High | Use existing TeslaPay AWS account or pre-approved sandbox |
| EKS cluster provisioning issues | Low | Medium | Fallback to docker-compose for local dev while resolving |
| Team members not yet onboarded | Medium | Medium | Sprint 1 scoped for 10-12 people; full team by Sprint 3 |

---

## Sprint 2: Auth Service, Account Service Core, App Shell

**Dates:** 2026-03-23 to 2026-04-03
**Sprint Goal:** Deliver working Auth and Account services so users can register and authenticate, and a Flutter app shell with login/registration screens.

### User Stories

| Story ID | Story | Priority | Size | Notes |
|----------|-------|----------|------|-------|
| US-1.1 | Account Registration (email, phone, password) | Must | M | Backend registration endpoint + Flutter registration screen |
| US-6.1 | Biometric Login | Must | M | Flutter biometric integration; backend session token flow |
| US-6.2 | PIN Setup and Management | Must | S | 6-digit PIN setup during onboarding |
| US-9.1 | Language Selection | Must | S | i18n framework setup; 5 languages |

### Technical Tasks

| ID | Task | Owner | Size | Notes |
|----|------|-------|------|-------|
| T2.01 | Implement Auth Service: user registration, login, JWT issuance | Backend | L | Password hashing (argon2), email/phone OTP |
| T2.02 | Implement Auth Service: session management, refresh tokens | Backend | M | Redis-backed sessions, 15min/30day TTL |
| T2.03 | Implement Auth Service: device binding, biometric credential storage | Backend | M | Public key registration per device |
| T2.04 | Implement Account Service: user profile CRUD | Backend | M | PostgreSQL, basic profile fields |
| T2.05 | Implement Account Service: Lithuanian IBAN generation | Backend | M | LT format, uniqueness guarantee |
| T2.06 | Implement Account Service: account tier assignment (default Basic) | Backend | S | Tier config, feature flags per tier |
| T2.07 | Flutter: Registration flow (email, phone, OTP verification) | Flutter | M | Multi-step form, OTP input |
| T2.08 | Flutter: Login screen with biometric + PIN fallback | Flutter | M | Face ID, Touch ID, fingerprint, PIN |
| T2.09 | Flutter: App shell with navigation (bottom tab bar) | Flutter | M | Dashboard, Payments, Cards, Crypto, Settings tabs |
| T2.10 | Flutter: i18n setup with 5 languages | Flutter | M | EN, LT, RU, DE, PL |
| T2.11 | Flutter: Secure storage for tokens and credentials | Flutter | S | flutter_secure_storage |
| T2.12 | Set up OpenAPI specs for Auth and Account services | Backend | S | Auto-generated from code annotations |
| T2.13 | Admin dashboard: project scaffolding (Next.js) | Frontend | M | Authentication against Auth Service |
| T2.14 | Set up integration test framework (testcontainers-go) | QA | M | Template for all backend services |
| T2.15 | Staging environment provisioning (mirror of dev, multi-AZ) | DevOps | L | Production-like topology |

### Definition of Done

- [ ] User can register with email + phone via the Flutter app
- [ ] OTP verification works for both email and phone
- [ ] User can log in with password, biometrics, or PIN
- [ ] JWT tokens issued and validated by API Gateway
- [ ] Account created with Lithuanian IBAN upon registration
- [ ] App displays in all 5 languages with working language switcher
- [ ] Auth Service: >80% unit test coverage, integration tests passing
- [ ] Account Service: >80% unit test coverage, integration tests passing
- [ ] Staging environment provisioned and accessible

### Sprint Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| IBAN generation algorithm complexity | Low | Medium | Research BIC/IBAN libraries (iban4j); validate with Bank of Lithuania format |
| Biometric SDK differences across iOS/Android | Medium | Low | Flutter platform channels abstract differences; test on physical devices |
| Sumsub contract not signed yet | Medium | Low | Sprint 2 does not require Sumsub; Sprint 3 does |

---

## Sprint 3: Ledger Service, KYC/Sumsub Integration

**Dates:** 2026-04-06 to 2026-04-17
**Sprint Goal:** Build the core financial ledger with event sourcing and integrate Sumsub for KYC verification, enabling the full onboarding pipeline.

### User Stories

| Story ID | Story | Priority | Size | Notes |
|----------|-------|----------|------|-------|
| US-1.2 | KYC Document Verification (Sumsub) | Must | L | Sumsub SDK in Flutter, webhook processing on backend |
| US-1.3 | Liveness Check | Must | M | Sumsub liveness SDK |
| US-1.5 | AML Screening at Onboarding | Must | M | Automatic post-verification screening |
| US-1.6 | Account Tier Assignment | Must | S | Auto-assign Basic tier after KYC approval |

### Technical Tasks

| ID | Task | Owner | Size | Notes |
|----|------|-------|------|-------|
| T3.01 | Implement Ledger Service: chart of accounts, journal entry posting | Backend | XL | Double-entry, ACID, event sourcing |
| T3.02 | Implement Ledger Service: event store with cryptographic chaining | Backend | L | Append-only, SHA-256 checksums |
| T3.03 | Implement Ledger Service: balance projection (CQRS read model) | Backend | M | Materialized balances from event stream |
| T3.04 | Implement Ledger Service: idempotency via posting IDs | Backend | M | Duplicate detection |
| T3.05 | Implement KYC Service: Sumsub webhook receiver | Backend | M | Signature validation, state machine |
| T3.06 | Implement KYC Service: verification state machine | Backend | M | initiated -> doc_submitted -> liveness -> aml_screened -> approved/rejected |
| T3.07 | Implement KYC Service: AML screening result processing | Backend | M | Green/Red/Review paths |
| T3.08 | Implement KYC Service: manual review queue (admin API) | Backend | M | For compliance team |
| T3.09 | Flutter: Sumsub SDK integration (document capture + liveness) | Flutter | L | Camera permissions, error handling, retry logic |
| T3.10 | Flutter: KYC status display and retry flow | Flutter | S | Status screen: pending, approved, rejected, retry |
| T3.11 | Kafka topic setup for ledger.events, kyc.events | DevOps | S | Replication factor 3, exactly-once |
| T3.12 | Ledger Service integration tests (posting, balance, idempotency) | QA | L | testcontainers with PostgreSQL + Kafka |
| T3.13 | KYC Service integration tests (webhook processing, state transitions) | QA | M | Mock Sumsub webhooks |
| T3.14 | Admin dashboard: KYC review queue UI | Frontend | M | List pending reviews, approve/reject actions |

### Definition of Done

- [ ] Ledger Service posts double-entry journal entries; every transaction balances to zero
- [ ] Event sourcing operational: events appended, balances projected from events
- [ ] Idempotent posting: duplicate posting IDs rejected gracefully
- [ ] Sumsub SDK integrated in Flutter app: document capture and liveness check functional
- [ ] KYC Service processes Sumsub webhooks and transitions verification state correctly
- [ ] AML screening results processed: clean users auto-approved, matches queued for review
- [ ] Admin dashboard shows KYC review queue
- [ ] All services: >80% unit test coverage, integration tests passing

### Sprint Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Sumsub contract/sandbox not ready | Medium | High | Develop against Sumsub sandbox docs; mock webhooks for testing |
| Event sourcing complexity (first implementation) | Medium | Medium | Tech Lead reviews all Ledger PRs; spike completed before sprint |
| Sumsub Flutter SDK compatibility issues | Low | Medium | Test on physical devices early; have native bridge fallback |

---

## Sprint 4: Onboarding E2E, Mobile Auth Polish

**Dates:** 2026-04-20 to 2026-05-01
**Sprint Goal:** Complete the end-to-end onboarding journey from registration through KYC to account activation, and polish the authentication experience.

### User Stories

| Story ID | Story | Priority | Size | Notes |
|----------|-------|----------|------|-------|
| US-2.1 | View Account Dashboard | Must | M | Basic dashboard with balance display (EUR only initially) |
| US-6.3 | Transaction Confirmation (biometric/PIN for actions) | Must | M | SCA for sensitive operations |
| US-6.4 | Session Management | Must | M | View/terminate active sessions |
| US-6.5 | Two-Factor Authentication | Must | M | 2FA for high-risk actions |
| US-2.6 | Notification Preferences | Should | S | Push notification opt-in/out |

### Technical Tasks

| ID | Task | Owner | Size | Notes |
|----|------|-------|------|-------|
| T4.01 | E2E onboarding flow: registration -> OTP -> KYC -> account activation | Backend + Flutter | L | Full path with error handling |
| T4.02 | Account Service: trigger Fuse wallet placeholder after KYC (deferred creation) | Backend | S | Flag for Sprint 13 crypto integration |
| T4.03 | Auth Service: SCA step-up flow for PSD2 compliance | Backend | M | Biometric/PIN challenge before sensitive ops |
| T4.04 | Auth Service: active session listing and remote termination | Backend | M | Device name, OS, last active |
| T4.05 | Auth Service: 2FA via push notification and SMS OTP | Backend | M | For password change, payee addition, etc. |
| T4.06 | Notification Service: push notification via APNs + FCM | Backend | L | Template engine, i18n, device token management |
| T4.07 | Notification Service: email notifications (transactional) | Backend | M | SES integration, templates |
| T4.08 | Notification Service: SMS via provider (Twilio or similar) | Backend | M | OTP delivery, alerts |
| T4.09 | Flutter: Account dashboard screen (balance, quick actions) | Flutter | M | EUR balance, last 5 transactions placeholder |
| T4.10 | Flutter: Session management screen | Flutter | S | List sessions, terminate button |
| T4.11 | Flutter: Notification preferences screen | Flutter | S | Toggle categories, channels |
| T4.12 | Flutter: SCA confirmation bottom sheet (biometric/PIN) | Flutter | M | Reusable component for all confirmations |
| T4.13 | E2E test suite for onboarding flow | QA | L | Happy path + error paths (KYC fail, retry, etc.) |
| T4.14 | Performance baseline: Auth + Account + Ledger services | QA | M | Response time benchmarks under load |
| T4.15 | Certificate pinning in Flutter app | Flutter | S | Per MOB-009 requirement |

### Definition of Done

- [ ] Complete onboarding: user registers, verifies email/phone, completes KYC, gets activated account with IBAN
- [ ] Account dashboard shows EUR balance (zero for new accounts)
- [ ] SCA step-up works for sensitive operations (biometric or PIN required)
- [ ] Session management: user can view and terminate other sessions
- [ ] 2FA works for password change and other high-risk actions
- [ ] Push notifications delivered via APNs (iOS) and FCM (Android)
- [ ] E2E onboarding tests passing in staging environment
- [ ] Performance baseline documented: API response times under 200ms (p95)

### Sprint Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Push notification delivery reliability | Medium | Medium | Test on physical devices; implement retry logic |
| Onboarding flow edge cases (partial completion, timeout) | Medium | Medium | Implement resume capability; save progress server-side |
| Banking Circle contract still unsigned | Medium | High | No impact on Sprint 4; must be signed by Sprint 5 start |

---

## MILESTONE M1 GATE REVIEW (May 1, 2026)

Checklist for Phase 1 (Foundation) completion:
- [ ] All dev + staging infrastructure operational
- [ ] Auth, Account, Ledger, KYC, Notification services deployed and functional
- [ ] Flutter app: registration, login, KYC, basic dashboard working
- [ ] CI/CD pipeline fully operational (build, test, scan, deploy)
- [ ] Observability operational (logs, metrics, traces, alerts)
- [ ] External partner status: Sumsub signed, Banking Circle in progress, Enfuce in progress
- [ ] Team fully ramped (18-20 members)
- [ ] All Sprint 1-4 DoD items completed

**Go/No-Go Decision:** Proceed to Phase 2 (Core Banking) only if all critical path items are complete.

---

## Sprint 5: SEPA Payments, Internal Transfers

**Dates:** 2026-05-04 to 2026-05-15
**Sprint Goal:** Enable sending and receiving money via SEPA transfers and instant TeslaPay-to-TeslaPay payments.

### User Stories

| Story ID | Story | Priority | Size | Notes |
|----------|-------|----------|------|-------|
| US-3.1 | Send SEPA Credit Transfer | Must | L | Banking Circle integration for SCT |
| US-3.2 | Send SEPA Instant Transfer | Must | L | Banking Circle SCT Inst |
| US-3.3 | Internal Transfer (TeslaPay to TeslaPay) | Must | M | Instant, free, by IBAN/phone |
| US-3.5 | Receive SEPA Payments | Must | M | Incoming payment processing |

### Technical Tasks

| ID | Task | Owner | Size | Notes |
|----|------|-------|------|-------|
| T5.01 | Implement Payment Service: payment order creation and validation | Backend | L | IBAN validation, limit checks, SCA |
| T5.02 | Implement Payment Service: SEPA SCT submission via Banking Circle API | Backend | L | Outbox pattern, saga orchestration |
| T5.03 | Implement Payment Service: SEPA SCT Inst submission | Backend | L | 10-second settlement target |
| T5.04 | Implement Payment Service: internal transfer (debit/credit via Ledger) | Backend | M | Atomic, sub-2-second |
| T5.05 | Implement Payment Service: incoming SEPA payment processing (webhook) | Backend | M | Banking Circle settlement webhooks |
| T5.06 | Payment Service <-> Ledger Service gRPC integration (debit/credit posting) | Backend | M | Hold -> settle -> release pattern |
| T5.07 | Payment Service <-> Account Service gRPC (account validation, limits) | Backend | M | Balance check, tier limits |
| T5.08 | Kafka topics: payment.events (initiated, completed, failed) | DevOps | S | Exactly-once semantics |
| T5.09 | Flutter: Send money screen (SEPA + internal) | Flutter | L | IBAN input, amount, reference, confirmation |
| T5.10 | Flutter: Payment confirmation with SCA | Flutter | M | Summary screen -> biometric -> success |
| T5.11 | Flutter: Payment status tracking (processing, completed, failed) | Flutter | S | Real-time status via polling or WebSocket |
| T5.12 | Notification: payment sent/received push notifications | Backend | S | Templates for payment events |
| T5.13 | Integration tests: SEPA payment flow E2E (with Banking Circle sandbox) | QA | L | Happy path + failure scenarios |
| T5.14 | Admin dashboard: payment monitoring screen | Frontend | M | Payment list, status, manual actions |

### Definition of Done

- [ ] User can send SEPA Credit Transfer from Flutter app; payment submitted to Banking Circle
- [ ] SEPA Instant Transfer settles within 10 seconds (in sandbox)
- [ ] Internal TeslaPay-to-TeslaPay transfer completes in under 2 seconds
- [ ] Incoming SEPA payments credited to user account with push notification
- [ ] Payment creates correct double-entry ledger postings (debit sender, credit recipient or hold)
- [ ] Payment saga handles failures gracefully (reversal on settlement failure)
- [ ] All payment events published to Kafka and consumed by Notification + Audit services

### Sprint Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| Banking Circle contract/sandbox not ready | Medium | Critical | This blocks SEPA payments entirely; escalate to executive level |
| SEPA Instant complexity (24/7 availability) | Medium | Medium | Start with SCT; add SCT Inst as second deliverable |
| Saga pattern implementation complexity | Medium | Medium | Use well-tested saga library or pattern from tech lead spike |

---

## Sprint 6: Multi-Currency, FX Engine, Fee Management

**Dates:** 2026-05-18 to 2026-05-29
**Sprint Goal:** Enable multi-currency accounts and currency exchange, and implement the fee management engine.

### User Stories

| Story ID | Story | Priority | Size | Notes |
|----------|-------|----------|------|-------|
| US-2.2 | Multi-Currency Sub-Accounts | Must | M | Open USD, GBP, PLN, CHF accounts |
| US-3.4 | Currency Exchange | Must | L | Buy/sell currencies with live FX rates |
| US-9.3 | Fee Schedule and Limits View | Must | S | Display fees and current usage vs limits |
| US-3.7 | Manage Saved Payees | Should | S | Save, edit, delete payees |

### Technical Tasks

| ID | Task | Owner | Size | Notes |
|----|------|-------|------|-------|
| T6.01 | Account Service: multi-currency sub-account creation | Backend | M | IBAN per sub-account where applicable |
| T6.02 | Payment Service: FX engine (rate fetching, markup, locking) | Backend | L | ECB rates + aggregator, max 0.5% markup, 30s lock |
| T6.03 | Payment Service: currency exchange execution via Ledger | Backend | M | Debit source currency, credit target currency atomically |
| T6.04 | Ledger Service: multi-currency posting support | Backend | M | Separate entries per currency |
| T6.05 | Fee management engine: configurable fee rules per tier, operation type | Backend | L | Fee applied as ledger entries |
| T6.06 | Account Service: transaction limits per tier (velocity checks) | Backend | M | Daily, monthly, per-transaction limits |
| T6.07 | Account Service: saved payee CRUD | Backend | S | IBAN validation, name, default reference |
| T6.08 | FX rate caching in Redis (30s TTL) | Backend | S | Rate provider integration |
| T6.09 | Flutter: Multi-currency account management screen | Flutter | M | Open new currency, view balances per currency |
| T6.10 | Flutter: Currency exchange screen (source, target, rate, confirm) | Flutter | M | Live rate display, 30s countdown, SCA |
| T6.11 | Flutter: Fee schedule and limits screen | Flutter | S | Tier-based fee table, usage progress bars |
| T6.12 | Flutter: Saved payees management | Flutter | S | List, add, edit, delete |
| T6.13 | Integration tests: FX engine (rate lock, expiry, execution) | QA | M | Edge cases: expired rate, insufficient balance |
| T6.14 | Integration tests: multi-currency postings | QA | M | Cross-currency ledger balancing |

### Definition of Done

- [ ] User can open sub-accounts in USD, GBP, PLN, CHF from the app
- [ ] Currency exchange shows live rate with max 0.5% markup; rate locked for 30 seconds
- [ ] Exchange executes instantly after confirmation; both sub-accounts updated
- [ ] Fee management applies correct fees based on account tier and operation type
- [ ] Transaction limits enforced per tier (daily, monthly, per-transaction)
- [ ] Saved payees: user can save, edit, delete, and use in payment flow
- [ ] All FX transactions create correct multi-currency ledger entries
- [ ] Fee schedule visible in app, reflecting user's current tier

### Sprint Risks

| Risk | Probability | Impact | Mitigation |
|------|------------|--------|------------|
| FX rate provider reliability | Medium | Medium | Cache rates; use ECB reference rates as fallback |
| Multi-currency ledger complexity | Medium | Medium | Thorough testing; tech lead review of all ledger changes |
| Fee calculation edge cases | Low | Medium | Comprehensive test suite for fee scenarios per tier |

---

## MILESTONE M2 PROGRESS CHECK (End of Sprint 6)

Mid-phase checkpoint (M2 full gate is at Sprint 8):
- [ ] SEPA SCT and SCT Inst functional in sandbox
- [ ] Internal transfers working end-to-end
- [ ] Multi-currency accounts operational
- [ ] FX engine with rate locking functional
- [ ] Fee management engine operational
- [ ] All Sprint 5-6 DoD items completed
- [ ] Banking Circle integration on track for production readiness by Sprint 8
- [ ] Enfuce contract status: must be signed by end of Sprint 6 for Sprint 9 card development

---

## Sprint Velocity Tracking Template

| Sprint | Planned Points | Completed Points | Velocity | Notes |
|--------|---------------|-----------------|----------|-------|
| Sprint 1 | -- | -- | -- | Technical setup, no story points |
| Sprint 2 | TBD | -- | -- | First measured sprint |
| Sprint 3 | TBD | -- | -- | |
| Sprint 4 | TBD | -- | -- | |
| Sprint 5 | TBD | -- | -- | |
| Sprint 6 | TBD | -- | -- | |

*Story point estimation to be calibrated during Sprint 2 planning based on team composition.*

---

## Cross-Sprint Dependency Matrix (Sprints 1-6)

```
S1 Infrastructure -------> S2 Auth + Account -------> S3 Ledger + KYC
        |                         |                         |
        |                         v                         v
        |                   S4 Onboarding E2E -------> S5 Payments
        |                         |                         |
        |                         v                         v
        +---------------------------------------------> S6 Multi-Ccy + FX
```

| From | To | Dependency | Risk if Late |
|------|----|-----------|-------------|
| S1 | S2 | Infrastructure must be operational | All development blocked |
| S2 | S3 | Auth + Account needed for KYC flow | KYC cannot process users |
| S2 | S4 | Auth needed for onboarding E2E | Cannot test full journey |
| S3 | S4 | Ledger + KYC needed for account activation | No active accounts |
| S3 | S5 | Ledger needed for payment posting | Payments cannot process |
| S4 | S5 | Active accounts needed for payments | No one to send/receive |
| S5 | S6 | Payment Service needed for FX execution | FX engine has no execution path |

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| Product Owner | TBD | | Pending |
| Engineering Manager | TBD | | Pending |
| Project Manager | Dream Team PM | 2026-03-03 | Submitted |
