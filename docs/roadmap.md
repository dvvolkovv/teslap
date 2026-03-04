# TeslaPay Development Roadmap

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Senior Technical Project Manager, Dream Team
**Timeline:** 9 months (18 sprints, 2 weeks each)
**Target Launch:** 2026-12-01

---

## 1. Executive Summary

This roadmap covers the TeslaPay neobank MVP (Phase 1), transforming the existing EMI into a modern neobank with Core Banking System, Mastercard cards, Fuse.io crypto integration, and native mobile applications. The project spans 18 two-week sprints (9 months) organized into 5 phases with clearly defined milestones, dependencies, and a critical path.

**Start Date:** 2026-03-09 (Sprint 1)
**Soft Launch (Beta):** 2026-11-02 (Sprint 17)
**Public Launch:** 2026-12-01 (Sprint 18 end)

---

## 2. Team Composition

### Recommended Team: 18-20 People

| Role | Count | Responsibilities |
|------|-------|-----------------|
| Engineering Manager / Tech Lead | 1 | Technical leadership, architecture oversight, code reviews |
| Senior Backend Engineers (Go) | 4 | CBS, Auth, Payment, Card, KYC, Audit, Fraud services |
| Backend Engineer (TypeScript/NestJS) | 1 | Crypto Service (Fuse.io integration) |
| Senior Flutter Engineers | 2 | iOS and Android mobile application |
| Frontend Engineer (React/Next.js) | 1 | Admin dashboard and compliance tools |
| DevOps / Platform Engineer | 2 | AWS infrastructure, Kubernetes, CI/CD, monitoring |
| QA Engineer | 2 | Test automation, E2E testing, security testing |
| Product Manager | 1 | Backlog management, stakeholder communication, acceptance |
| UI/UX Designer | 1 | Ongoing design refinement, user testing |
| Compliance / Regulatory Specialist | 1 | Bank of Lithuania reporting, AML/KYC process design |
| Security Engineer | 1 | Penetration testing, security reviews, PCI DSS coordination |
| Project Manager | 1 | This role -- coordination, risk management, delivery tracking |

### Team Ramp-Up Schedule

| Period | Team Size | Notes |
|--------|-----------|-------|
| Sprints 1-2 | 10-12 | Core backend, DevOps, Flutter lead, PM, designer |
| Sprints 3-4 | 14-16 | Add remaining backend, second Flutter, QA, compliance |
| Sprints 5-18 | 18-20 | Full team including security engineer, second QA |

---

## 3. Phase Breakdown

### Phase 1: Foundation (Sprints 1-4, Weeks 1-8)
**Dates:** 2026-03-09 to 2026-05-01

| Sprint | Dates | Focus |
|--------|-------|-------|
| Sprint 1 | Mar 9 - Mar 20 | Infrastructure setup, CI/CD, project scaffolding, protobuf definitions |
| Sprint 2 | Mar 23 - Apr 3 | Auth Service, Account Service core, Flutter app shell, admin dashboard shell |
| Sprint 3 | Apr 6 - Apr 17 | Ledger Service (event sourcing, double-entry), KYC Service + Sumsub integration |
| Sprint 4 | Apr 20 - May 1 | Onboarding flow E2E (registration through KYC to account activation), mobile auth |

**Milestone M1: Foundation Complete (May 1)**
- Acceptance Criteria:
  - AWS infrastructure provisioned (EKS, RDS, MSK, ElastiCache) in dev + staging
  - Auth Service operational with JWT issuance, session management, biometric support
  - Account Service creates accounts with Lithuanian IBAN generation
  - Ledger Service posts double-entry journal entries with event sourcing
  - KYC Service processes Sumsub webhooks and manages verification state machine
  - Flutter app shell with authentication flow (login, registration, biometric)
  - CI/CD pipeline deploys to dev and staging automatically
  - All services have >80% unit test coverage

---

### Phase 2: Core Banking (Sprints 5-8, Weeks 9-16)
**Dates:** 2026-05-04 to 2026-06-26

| Sprint | Dates | Focus |
|--------|-------|-------|
| Sprint 5 | May 4 - May 15 | SEPA payments (SCT, SCT Inst) via Banking Circle, internal transfers |
| Sprint 6 | May 18 - May 29 | Multi-currency accounts, FX engine, fee management |
| Sprint 7 | Jun 1 - Jun 12 | Account dashboard, transaction history, notification service |
| Sprint 8 | Jun 15 - Jun 26 | Fraud detection (rules engine), recurring payments, reconciliation engine |

**Milestone M2: Core Banking Operational (June 26)**
- Acceptance Criteria:
  - SEPA Credit Transfer (SCT) and SEPA Instant (SCT Inst) working end-to-end
  - Internal TeslaPay-to-TeslaPay transfers instant and free
  - Multi-currency sub-accounts (EUR, USD, GBP, PLN, CHF) with FX exchange
  - Account dashboard shows balances and transactions on mobile
  - Push notifications delivered for transactions within 5 seconds
  - Fraud detection rules engine scoring transactions in real-time
  - Daily automated reconciliation running
  - Banking Circle sandbox integration tested and validated

---

### Phase 3: Cards and Payments (Sprints 9-12, Weeks 17-24)
**Dates:** 2026-06-29 to 2026-08-21

| Sprint | Dates | Focus |
|--------|-------|-------|
| Sprint 9 | Jun 29 - Jul 10 | Card Service core: virtual card issuance, card lifecycle (freeze/unfreeze) |
| Sprint 10 | Jul 13 - Jul 24 | Physical card ordering, PIN management, spending controls, 3D Secure |
| Sprint 11 | Jul 27 - Aug 7 | Apple Pay + Google Pay tokenization (Mastercard MDES), ATM withdrawal support |
| Sprint 12 | Aug 10 - Aug 21 | Card authorization flow E2E, real-time notifications, dispute management |

**Milestone M3: Card Program Live (August 21)**
- Acceptance Criteria:
  - Virtual Mastercard issued within 30 seconds
  - Physical card ordering with delivery tracking
  - Card freeze/unfreeze within 5 seconds
  - 3D Secure 2.0 authentication via push notification
  - Apple Pay and Google Pay provisioning working
  - Spending controls (per-transaction, daily, category, geographic)
  - Real-time transaction notifications within 3 seconds
  - Card authorization processing under 100ms
  - Enfuce sandbox fully integrated and tested

---

### Phase 4: Crypto Integration (Sprints 13-16, Weeks 25-32)
**Dates:** 2026-08-24 to 2026-10-16

| Sprint | Dates | Focus |
|--------|-------|-------|
| Sprint 13 | Aug 24 - Sep 4 | Fuse Smart Wallet creation, crypto balance display, wallet UI |
| Sprint 14 | Sep 7 - Sep 18 | Buy crypto (fiat-to-crypto on-ramp), sell crypto (crypto-to-fiat off-ramp) |
| Sprint 15 | Sep 21 - Oct 2 | Send/receive crypto, crypto transaction history, MiCA compliance review |
| Sprint 16 | Oct 5 - Oct 16 | Gasless transactions (ERC-4337), crypto-fiat bridge polish, compliance sign-off |

**Milestone M4: Crypto Features Complete (October 16)**
- Acceptance Criteria:
  - Fuse Smart Wallet created automatically after KYC approval
  - Crypto balances (FUSE, USDC, USDT) displayed with EUR equivalent
  - Buy/sell crypto with EUR balance, rate shown before confirmation
  - Send crypto to external Fuse addresses with QR scan
  - Receive crypto with deposit address and push notifications
  - Crypto transaction history with block explorer links
  - Gasless transactions via account abstraction (if "Should" priority met)
  - MiCA compliance review passed for all crypto features
  - Circuit breaker protects platform from Fuse network instability

---

### Phase 5: Polish and Launch (Sprints 17-18, Weeks 33-36)
**Dates:** 2026-10-19 to 2026-12-01

| Sprint | Dates | Focus |
|--------|-------|-------|
| Sprint 17 | Oct 19 - Oct 30 | UAT, penetration testing, performance testing, compliance audit prep, beta onboarding |
| Sprint 18 | Nov 2 - Nov 13 | Beta feedback fixes, data migration tooling, regulatory submission, launch preparation |
| Buffer | Nov 16 - Dec 1 | Final fixes, App Store/Play Store submission, marketing coordination, public launch |

**Milestone M5: Soft Launch / Beta (November 2)**
- Acceptance Criteria:
  - All "Must" priority user stories implemented and tested
  - Penetration test completed with zero critical/high findings
  - Performance test passed at 3x projected load (300 TPS)
  - Bank of Lithuania regulatory reporting operational
  - 500 beta users onboarded
  - App crash rate below 0.1%
  - System uptime above 99.95% during beta period

**Milestone M6: Public Launch (December 1)**
- Acceptance Criteria:
  - All beta feedback critical issues resolved
  - App Store and Google Play Store approval received
  - Data migration from legacy system tested and validated
  - Marketing campaign assets ready
  - Customer support team trained and operational
  - Monitoring and alerting fully operational with PagerDuty integration
  - DR failover tested successfully
  - Go/no-go checklist completed with all stakeholders

---

## 4. Workstream Dependencies

### Dependency Graph

```
Infrastructure (DevOps)
    |
    +---> Auth Service ----+
    |                      |
    +---> Account Service -+---> Payment Service --+---> Card Service
    |          |           |          |             |         |
    |          v           |          v             |         v
    |     KYC Service      |    Ledger Service      |   Apple/Google Pay
    |          |           |          |             |
    |          v           |          v             |
    |   Sumsub Integration |   Banking Circle       |
    |                      |   Integration          |
    +---> Notification Svc +                        |
    |                                               |
    +---> Fraud Detection Svc <---------------------+
    |
    +---> Crypto Service (independent until fiat bridge)
    |          |
    |          v
    |     Fuse.io Integration
    |
    +---> Audit Service (consumes events from all services)
    |
    +---> Flutter Mobile App (parallel, integrates with backend services incrementally)
    |
    +---> Admin Dashboard (parallel, integrates with backend services incrementally)
```

### Critical Dependencies (External)

| Dependency | Required By | Lead Time | Risk |
|-----------|-------------|-----------|------|
| Enfuce card processor contract + sandbox | Sprint 9 (Jun 29) | 4-8 weeks negotiation | HIGH -- must initiate by Sprint 3 |
| Banking Circle contract + sandbox | Sprint 5 (May 4) | 4-6 weeks negotiation | HIGH -- must initiate by Sprint 1 |
| Sumsub contract + API keys | Sprint 3 (Apr 6) | 2-3 weeks | MEDIUM -- must initiate immediately |
| Fuse.io partnership + SDK access | Sprint 13 (Aug 24) | 2-4 weeks | MEDIUM -- initiate by Sprint 8 |
| Apple Developer Enterprise account | Sprint 11 (Jul 27) | 1-2 weeks | LOW |
| Google Play Developer account | Sprint 11 (Jul 27) | 1-2 days | LOW |
| Mastercard BIN allocation via Enfuce | Sprint 9 (Jun 29) | 6-12 weeks | HIGH -- driven by Enfuce relationship |
| MiCA legal review for crypto features | Sprint 15 (Sep 21) | 4-6 weeks | MEDIUM -- engage counsel by Sprint 10 |
| Bank of Lithuania pre-notification | Sprint 16 (Oct 5) | 4-8 weeks | MEDIUM -- engage by Sprint 12 |

### Internal Dependencies (Service-to-Service)

| Dependent Service | Depends On | Blocking? | Notes |
|------------------|-----------|-----------|-------|
| Account Service | Auth Service | Yes | Account creation requires authenticated user |
| Ledger Service | None (core) | N/A | Foundation service, no upstream dependencies |
| Payment Service | Ledger Service, Account Service | Yes | Payments post to ledger, validate against accounts |
| Card Service | Account Service, Ledger Service, Fraud Detection | Yes | Cards linked to accounts, authorizations post to ledger |
| Crypto Service | Account Service, Ledger Service, KYC Service | Yes | Wallet creation after KYC, fiat bridge via ledger |
| KYC Service | Account Service | Yes | KYC results update account status |
| Notification Service | Kafka cluster | Yes | Consumes events from all services |
| Audit Service | Kafka cluster | Yes | Consumes events from all services |
| Fraud Detection | Ledger Service, Card Service | Partially | Sync scoring for card auth, async for transaction analysis |
| Flutter App | All backend services | Rolling | Integrates incrementally as services become available |

---

## 5. Critical Path Analysis

The critical path determines the minimum project duration. Any delay on critical path items delays the launch.

### Critical Path

```
Infrastructure Setup (S1)
  -> Auth + Account Services (S2)
    -> Ledger Service + KYC/Sumsub (S3)
      -> Onboarding E2E (S4)
        -> SEPA Payments / Banking Circle (S5)
          -> Multi-currency + FX (S6)
            -> Dashboard + Notifications (S7)
              -> Fraud + Reconciliation (S8)
                -> Card Issuance / Enfuce (S9)
                  -> Physical Cards + 3DS (S10)
                    -> Apple/Google Pay (S11)
                      -> Card Auth E2E (S12)
                        -> UAT + Pen Test (S17)
                          -> Beta + Launch Prep (S18)
                            -> PUBLIC LAUNCH
```

**Critical Path Duration:** 18 sprints (36 weeks)
**Float:** Zero on critical path items; 4 sprints of float on crypto workstream (Sprints 13-16 run parallel to Cards Sprints 9-12 completion)

### Near-Critical Paths

1. **Crypto Path** (4 sprints float): Crypto features (S13-16) run after Cards phase but are not gated by card completion. If cards slip, crypto can still proceed. However, crypto MUST complete before Sprint 17 UAT.

2. **Admin Dashboard** (6+ sprints float): Admin dashboard development runs in parallel and is not on the critical path. Compliance dashboards must be ready by Sprint 17 for regulatory review.

3. **Data Migration** (2 sprints float): Legacy system migration tooling in Sprint 18 has minimal float before launch.

---

## 6. Milestone Summary

| ID | Milestone | Target Date | Sprint | Gate |
|----|-----------|-------------|--------|------|
| M0 | Project Kickoff | 2026-03-09 | S1 | Team assembled, contracts initiated |
| M1 | Foundation Complete | 2026-05-01 | S4 end | Infrastructure, auth, accounts, ledger, KYC operational |
| M2 | Core Banking Operational | 2026-06-26 | S8 end | SEPA payments, FX, dashboard, notifications, fraud rules |
| M3 | Card Program Live | 2026-08-21 | S12 end | Virtual/physical cards, Apple/Google Pay, 3DS, ATM |
| M4 | Crypto Features Complete | 2026-10-16 | S16 end | Fuse wallet, buy/sell, send/receive, MiCA compliance |
| M5 | Soft Launch (Beta) | 2026-11-02 | S17 end | 500 beta users, pen test passed, compliance review |
| M6 | Public Launch | 2026-12-01 | S18 + buffer | App stores approved, marketing live, full operation |

---

## 7. Key Assumptions

1. Team hiring completes by Sprint 3 (full team of 18-20 by end of April 2026).
2. External partner contracts (Enfuce, Banking Circle, Sumsub) are initiated immediately and signed within 6 weeks.
3. TeslaPay's existing EMI license covers all planned MVP features without amendment.
4. Sandbox/test environments from all external partners are available within 2 weeks of contract signing.
5. No major regulatory changes (PSD3, MiCA amendments) during the development period that would require rework.
6. Budget is approved for the recommended team size and AWS infrastructure costs (estimated EUR 15K-25K/month for infrastructure).
7. Product Owner is available for backlog refinement and acceptance testing throughout the project.
8. Legacy system data migration requirements are documented by Sprint 12.

---

## 8. Budget Estimates (High-Level)

| Category | Monthly Cost (EUR) | 9-Month Total (EUR) |
|----------|-------------------|---------------------|
| Engineering team (18 people avg) | 135,000 - 180,000 | 1,215,000 - 1,620,000 |
| AWS infrastructure (dev+staging+prod) | 15,000 - 25,000 | 135,000 - 225,000 |
| External services (Sumsub, Enfuce, Banking Circle) | 5,000 - 15,000 | 45,000 - 135,000 |
| Tooling (GitHub, SonarQube, PagerDuty, etc.) | 3,000 - 5,000 | 27,000 - 45,000 |
| Security (pen testing, SOC 2 prep) | 5,000 - 10,000 | 45,000 - 90,000 |
| Contingency (15%) | -- | 220,000 - 320,000 |
| **Total** | **~185,000** | **~1,690,000 - 2,435,000** |

*Note: These are planning estimates. Detailed budgeting should be performed during Sprint 0 / project initiation.*

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| Product Owner | TBD | | Pending |
| CTO | TBD | | Pending |
| Engineering Manager | TBD | | Pending |
| Project Manager | Dream Team PM | 2026-03-03 | Submitted |
