# Product Requirements Document: TeslaPay Neobank Infrastructure

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Business Analyst, Dream Team
**Status:** Draft for Review
**Product:** TeslaPay Neobank Platform (https://www.teslapay.eu/)

---

## 1. Executive Summary

TeslaPay (UAB TeslaPay) is a Lithuanian-registered Electronic Money Institution (EMI) regulated by the Bank of Lithuania. The company has been operational since 2017, holding EUR 9.85M in safeguarded customer funds and generating EUR 1.69M in revenue (2024). This PRD defines requirements for a complete infrastructure rebuild -- a modern, scalable neobank platform comprising a Core Banking System (CBS), Mastercard card issuance, Fuse.io blockchain integration, Sumsub-powered KYC/AML, and native mobile applications for iOS and Android.

## 2. Vision and Strategic Goals

### 2.1 Product Vision

Transform TeslaPay from a basic EMI with IBAN accounts and prepaid cards into a full-featured European neobank that bridges traditional finance and Web3, offering multi-currency accounts, Mastercard debit cards, blockchain-powered payments, and DeFi yield -- all through a modern mobile-first experience.

### 2.2 Strategic Goals

| ID | Goal | Measure | Target (12 months post-launch) |
|----|------|---------|-------------------------------|
| G1 | Grow active user base | Monthly Active Users (MAU) | 50,000 |
| G2 | Increase customer funds | Total safeguarded funds | EUR 50M |
| G3 | Revenue growth | Annual revenue | EUR 8M |
| G4 | Card adoption | Active Mastercard holders | 25,000 |
| G5 | Crypto engagement | Users with Fuse wallet activated | 10,000 |
| G6 | Compliance excellence | Regulatory incidents | 0 critical findings |
| G7 | App store rating | Average rating (iOS + Android) | 4.5+ stars |

## 3. Scope

### 3.1 In Scope (MVP -- Phase 1, 0-9 months)

1. **Core Banking System (CBS)**
   - Double-entry general ledger
   - Customer account management (current accounts, EUR primary)
   - Multi-currency support (EUR, USD, GBP, PLN, CHF minimum)
   - SEPA Credit Transfers (SCT) and SEPA Direct Debits (SDD)
   - SEPA Instant Credit Transfers (SCT Inst)
   - Internal transfers between TeslaPay accounts
   - Interest calculation engine (savings accounts)
   - Fee management engine
   - Transaction limits and velocity controls

2. **User Accounts and Onboarding**
   - Mobile-first registration flow
   - Personal accounts with dedicated IBAN (Lithuanian IBAN)
   - KYC/AML verification via Sumsub integration
   - Multi-level account tiers (Basic, Standard, Premium)
   - Account settings and profile management

3. **Accounting and Regulatory Reporting**
   - General ledger with full audit trail
   - Automated reconciliation engine
   - Bank of Lithuania regulatory reporting
   - AML transaction monitoring and SAR filing support
   - GDPR-compliant data management
   - Tax reporting support (CRS/FATCA)

4. **Mastercard Card Issuance**
   - Virtual card issuance (instant)
   - Physical card issuance and delivery
   - Mastercard Debit card program
   - Card lifecycle management (activate, freeze, block, replace, PIN management)
   - 3D Secure 2.0 authentication
   - Contactless (NFC) payments
   - Apple Pay and Google Pay tokenization
   - Real-time transaction notifications
   - Spending controls (merchant category blocks, geographic restrictions, transaction limits)
   - ATM withdrawal support

5. **Sumsub Integration**
   - Document verification (passport, ID card, driver's license, residence permit)
   - Liveness detection with deepfake checks
   - NFC document reading (ePassport chip)
   - AML screening (sanctions, PEP, adverse media)
   - Ongoing monitoring (continuous AML screening)
   - Risk scoring
   - Address verification
   - Proof of funds / Source of income verification workflows

6. **Fuse.io Integration (Basic)**
   - Fuse Smart Wallet creation per user
   - FUSE token balance display
   - Stablecoin support (USDC, USDT on Fuse network)
   - Basic crypto-to-fiat and fiat-to-crypto conversion
   - Crypto deposit and withdrawal
   - Transaction history for blockchain operations

7. **Mobile Applications**
   - Native iOS application (Swift/SwiftUI, iOS 16+)
   - Native Android application (Kotlin/Jetpack Compose, Android 10+)
   - Biometric authentication (Face ID, Touch ID, fingerprint)
   - PIN/passcode fallback
   - Push notifications (transactional, marketing with opt-in)
   - In-app chat support
   - Multi-language support (English, Lithuanian, Russian, German, Polish minimum)
   - Dark mode
   - Accessibility (WCAG 2.1 AA)

### 3.2 In Scope (Phase 2, 9-18 months)

1. **Business Accounts**
   - KYB (Know Your Business) via Sumsub
   - Multi-user access with role-based permissions
   - Bulk payment processing (SEPA batch)
   - Invoice management
   - Accounting software integration (Xero, QuickBooks)

2. **Advanced Fuse.io / DeFi Features**
   - Solid integration (soUSD yield-bearing stablecoin)
   - DeFi yield (savings via Fuse DeFi protocols)
   - Token swap functionality
   - NFT display and management
   - Gasless transactions via account abstraction (ERC-4337)

3. **Open Banking (PSD2/PSD3)**
   - Account Information Service (AIS) -- aggregate external bank accounts
   - Payment Initiation Service (PIS)
   - Third-party API access (TPP)
   - Strong Customer Authentication (SCA) compliant flows

4. **Advanced Card Features**
   - Cashback / rewards program
   - Installment payments
   - Card-linked offers
   - Corporate cards with expense management

5. **Savings and Financial Products**
   - Term deposits
   - Savings goals (round-ups, scheduled deposits)
   - Budgeting tools with spending analytics
   - Recurring payments and standing orders

### 3.3 Out of Scope (Won't)

- Full banking license acquisition (TeslaPay operates under EMI license)
- Lending / credit products (requires separate license)
- Insurance products
- Stock or ETF trading
- Proprietary blockchain development
- Desktop/web banking application (Phase 2+ consideration)

## 4. Target Users

See `docs/personas.md` for detailed personas. Summary:

| Persona | Description | Primary Needs |
|---------|-------------|---------------|
| Digital Native Expat | 25-35, lives/works across EU borders | Multi-currency, low FX fees, fast onboarding |
| Crypto-Curious Professional | 28-40, interested in crypto but not expert | Easy crypto access, yield, traditional banking in one app |
| Small Business Owner | 30-50, freelancer or small EU business | Business IBAN, card expense management, multi-currency |
| Privacy-Conscious User | 30-45, values data sovereignty | Self-custodial crypto wallet, minimal data sharing, EU-regulated |
| Tech-Savvy Student/Young Professional | 18-25, first banking relationship | Low/no fees, modern UX, instant cards, peer payments |

## 5. Functional Requirements

### 5.1 Core Banking System (CBS)

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| CBS-001 | System shall maintain a double-entry general ledger with real-time posting | Must | Every transaction creates at least two ledger entries that balance to zero |
| CBS-002 | System shall support creation of current accounts with unique IBAN | Must | IBAN generated follows Lithuanian format (LTxx), unique per account |
| CBS-003 | System shall support multi-currency accounts (EUR, USD, GBP, PLN, CHF) | Must | User can hold balances in 5+ currencies; each currency has separate sub-account |
| CBS-004 | System shall process SEPA Credit Transfers within 1 business day | Must | SCT submitted before cutoff time settles T+1; status updates in real-time |
| CBS-005 | System shall process SEPA Instant Credit Transfers within 10 seconds | Must | SCT Inst settles within 10s; 99.5% availability 24/7/365 |
| CBS-006 | System shall support internal transfers (TeslaPay-to-TeslaPay) instantly | Must | Transfer completes in under 2 seconds; both accounts updated atomically |
| CBS-007 | System shall calculate and apply interest on savings accounts daily | Must | Interest accrued daily, paid monthly; rate configurable per tier/product |
| CBS-008 | System shall enforce transaction limits per account tier | Must | Limits checked before processing; declined transactions return clear error |
| CBS-009 | System shall support scheduled/recurring payments | Should | User can set frequency (daily/weekly/monthly/custom), amount, and payee |
| CBS-010 | System shall maintain complete transaction history with search/filter | Must | All transactions queryable by date, amount, type, status; export to CSV/PDF |
| CBS-011 | System shall support currency exchange at competitive FX rates | Must | FX rate displayed before confirmation; markup not exceeding 0.5% on mid-market |
| CBS-012 | System shall process SEPA Direct Debits (as payer) | Should | User can authorize and manage SDD mandates |

### 5.2 User Account Management

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| UAM-001 | System shall support self-service registration via mobile app | Must | Registration completes in under 5 minutes for standard flow |
| UAM-002 | System shall assign account tiers (Basic, Standard, Premium) | Must | Tier determines limits, features, and pricing; upgradeable |
| UAM-003 | System shall support multi-language interface | Must | Minimum 5 languages; language switchable without re-login |
| UAM-004 | System shall allow users to update personal information | Must | Changes to regulated fields (name, address) trigger re-verification |
| UAM-005 | System shall support account closure with fund disbursement | Must | Remaining funds transferred out; account archived per retention policy |
| UAM-006 | System shall enforce email and phone verification | Must | OTP-based verification; both required before account activation |
| UAM-007 | System shall support beneficiary/payee management | Should | User can save, edit, delete payees; payee validation against IBAN registry |

### 5.3 Mastercard Card Management

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| MC-001 | System shall issue virtual Mastercard debit cards instantly | Must | Card number, expiry, CVV available in-app within 30 seconds of request |
| MC-002 | System shall support physical card ordering with delivery tracking | Must | Card shipped within 3 business days; tracking number provided |
| MC-003 | System shall support card freeze/unfreeze | Must | Freeze takes effect within 5 seconds; all new authorizations declined |
| MC-004 | System shall implement 3D Secure 2.0 | Must | 3DS challenge via push notification or in-app; fallback to SMS OTP |
| MC-005 | System shall support Apple Pay provisioning | Must | Card added to Apple Wallet via in-app button; tokenization via Mastercard MDES |
| MC-006 | System shall support Google Pay provisioning | Must | Card added to Google Wallet; tokenization via Mastercard MDES |
| MC-007 | System shall provide real-time transaction notifications | Must | Push notification within 3 seconds of authorization |
| MC-008 | System shall support spending controls | Must | User can set per-transaction limit, daily limit, merchant category blocks |
| MC-009 | System shall support PIN management (set, change, view) | Must | PIN viewable in-app with biometric auth; changeable at ATM or in-app |
| MC-010 | System shall support ATM withdrawals with configurable limits | Must | ATM withdrawal limit per tier; fee-free allowance configurable |
| MC-011 | System shall support card replacement (lost/stolen/damaged) | Must | Old card blocked immediately; new card issued within 1 business day (virtual) |
| MC-012 | System shall display card details securely in-app | Must | Card number masked by default; full number shown after biometric auth |

### 5.4 Sumsub KYC/AML Integration

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| KYC-001 | System shall perform document verification via Sumsub | Must | Accepts passport, national ID, driver's license, residence permit from 200+ countries |
| KYC-002 | System shall perform liveness detection | Must | Liveness check passes only for live person; blocks photos, videos, deepfakes |
| KYC-003 | System shall perform AML screening at onboarding | Must | Checks against sanctions lists, PEP databases, adverse media; result in under 60s |
| KYC-004 | System shall perform ongoing AML monitoring | Must | Continuous screening; alerts generated within 24h of new matches |
| KYC-005 | System shall support NFC document reading | Should | ePassport chip read for enhanced verification; reduces fraud |
| KYC-006 | System shall implement risk-based verification tiers | Must | Low-risk: basic docs; Medium: docs + proof of address; High: docs + POA + source of funds |
| KYC-007 | System shall support re-verification triggers | Must | Triggered by profile changes, transaction anomalies, or regulatory requirement |
| KYC-008 | System shall maintain verification audit trail | Must | Full log of verification attempts, decisions, and manual reviews stored for 5+ years |
| KYC-009 | System shall support manual review workflow | Must | Compliance team can review, approve, or reject flagged verifications |
| KYC-010 | System shall enforce geographic restrictions | Must | Onboarding blocked for sanctioned countries; configurable country list |

### 5.5 Fuse.io Blockchain Integration

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| FUSE-001 | System shall create a Fuse Smart Wallet for each verified user | Must | Wallet created automatically post-KYC; address displayed in-app |
| FUSE-002 | System shall display crypto balances (FUSE, USDC, USDT) | Must | Balances refresh within 30 seconds; show fiat equivalent |
| FUSE-003 | System shall support crypto deposit (receive) to Fuse wallet | Must | Deposit address (QR + copy) shown; incoming transfers detected within 1 block confirmation |
| FUSE-004 | System shall support crypto withdrawal (send) from Fuse wallet | Must | User specifies address and amount; confirmation with fee estimate before sending |
| FUSE-005 | System shall support fiat-to-crypto conversion (buy) | Must | User can buy FUSE/USDC/USDT using EUR balance; rate shown before confirmation |
| FUSE-006 | System shall support crypto-to-fiat conversion (sell) | Must | User can sell FUSE/USDC/USDT to EUR balance; rate shown before confirmation |
| FUSE-007 | System shall provide crypto transaction history | Must | All blockchain transactions listed with hash, timestamp, amount, status |
| FUSE-008 | System shall support gasless transactions via account abstraction | Should | Users do not need to hold FUSE tokens to pay gas; fees deducted from stablecoin balance |
| FUSE-009 | System shall integrate Solid soUSD yield (Phase 2) | Could | Users can deposit stablecoins into Solid for yield; APY displayed transparently |
| FUSE-010 | System shall support token swap within app (Phase 2) | Could | Swap between supported tokens on Fuse DEX; slippage tolerance configurable |

### 5.6 Mobile Application

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| MOB-001 | App shall support biometric login (Face ID, Touch ID, fingerprint) | Must | Login completes in under 2 seconds; fallback to PIN available |
| MOB-002 | App shall display account dashboard with balances and recent activity | Must | Dashboard loads in under 3 seconds; shows all currency balances |
| MOB-003 | App shall support push notifications for transactions | Must | Notification delivered within 5 seconds of transaction event |
| MOB-004 | App shall provide in-app customer support chat | Must | Chat accessible from any screen; response time SLA < 5 minutes |
| MOB-005 | App shall support dark mode | Should | Toggle in settings; respects system preference by default |
| MOB-006 | App shall meet WCAG 2.1 AA accessibility standards | Must | VoiceOver/TalkBack tested; all interactive elements have labels |
| MOB-007 | App shall support deep linking for notifications | Should | Tapping notification opens relevant transaction detail screen |
| MOB-008 | App shall support offline mode for balance/history viewing | Could | Cached data available when offline; clear indicator of stale data |
| MOB-009 | App shall implement certificate pinning | Must | API calls fail gracefully if certificate mismatch detected |
| MOB-010 | App shall support remote session termination | Must | User can log out all other sessions from settings |

### 5.7 Accounting and Compliance

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|-------------------|
| ACC-001 | System shall maintain a general ledger with immutable entries | Must | No ledger entry can be deleted; corrections via reversal entries only |
| ACC-002 | System shall perform automated daily reconciliation | Must | Reconciliation runs end-of-day; discrepancies flagged automatically |
| ACC-003 | System shall generate Bank of Lithuania regulatory reports | Must | Reports generated in required format and frequency |
| ACC-004 | System shall maintain full audit trail for all operations | Must | Every state change logged with timestamp, actor, previous/new value |
| ACC-005 | System shall support GDPR data subject requests | Must | Right to access, erasure (with regulatory retention exceptions), portability |
| ACC-006 | System shall implement data retention policies | Must | Transaction data retained 5 years minimum; personal data per GDPR |
| ACC-007 | System shall support CRS/FATCA reporting | Should | Annual reporting data extracted automatically |
| ACC-008 | System shall generate SAR (Suspicious Activity Reports) | Must | Compliance officers can file SARs; auto-generation for high-risk patterns |

## 6. Non-Functional Requirements

### 6.1 Performance

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-P01 | API response time (95th percentile) | < 200ms |
| NFR-P02 | Mobile app cold start time | < 3 seconds |
| NFR-P03 | Concurrent users supported | 10,000 minimum |
| NFR-P04 | Transaction processing throughput | 100 TPS minimum |
| NFR-P05 | Database query response time (95th percentile) | < 50ms |

### 6.2 Availability and Reliability

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-A01 | System uptime (CBS, Card Processing) | 99.95% (26 min downtime/month max) |
| NFR-A02 | SEPA Instant availability | 99.5% (per ECB requirement) |
| NFR-A03 | Recovery Time Objective (RTO) | < 15 minutes |
| NFR-A04 | Recovery Point Objective (RPO) | < 1 minute |
| NFR-A05 | Disaster recovery | Active-passive in separate EU region |

### 6.3 Security

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-S01 | Data at rest encryption | AES-256 |
| NFR-S02 | Data in transit encryption | TLS 1.3 |
| NFR-S03 | PCI DSS compliance | Level 1 (via card processor) |
| NFR-S04 | Penetration testing frequency | Quarterly + after major releases |
| NFR-S05 | SOC 2 Type II audit | Annually |
| NFR-S06 | Session timeout | 5 minutes inactive (configurable) |
| NFR-S07 | Sensitive data masking | PAN, CVV, passwords never logged |

### 6.4 Scalability

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-SC01 | Horizontal scaling | Auto-scale based on load; handle 10x traffic spikes |
| NFR-SC02 | Database scaling | Read replicas; partitioning by account |
| NFR-SC03 | Target capacity (18 months) | 500,000 accounts, 5M transactions/month |

### 6.5 Regulatory and Compliance

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-R01 | EMD2/PSD2 compliance | Full compliance with Bank of Lithuania requirements |
| NFR-R02 | PSD3 readiness | Architecture supports PSD3 transition (expected 2027) |
| NFR-R03 | GDPR compliance | Full compliance; DPO appointed |
| NFR-R04 | AML Directive 6 (AMLD6) compliance | Transaction monitoring, SAR, record keeping |
| NFR-R05 | MiCA (Markets in Crypto-Assets) compliance | Crypto features comply with MiCA framework |

## 7. Technical Architecture Constraints

### 7.1 Architecture Principles

- **Microservices architecture** with event-driven communication
- **Cloud-native** deployment (AWS or GCP, EU region -- Frankfurt or Ireland)
- **API-first** design with OpenAPI 3.0 specifications
- **Event sourcing** for financial transactions (immutable event log)
- **CQRS** (Command Query Responsibility Segregation) for read/write optimization

### 7.2 Key Technology Decisions

| Component | Requirement | Rationale |
|-----------|-------------|-----------|
| CBS Engine | Build or license (Mambu, Thought Machine, or custom) | Core differentiator; must support multi-currency and crypto |
| Card Processor | Mastercard-certified issuer processor (Marqeta, GPS, Enfuce) | BIN sponsorship and Mastercard network access |
| Payment Gateway | SEPA connectivity via banking partner or direct EBA/STEP2 | Required for SEPA SCT/SDD/Inst |
| KYC/AML | Sumsub (mandated) | Client requirement; supports EU regulatory needs |
| Blockchain | Fuse.io SDK (Flutter + TypeScript) (mandated) | Client requirement; EVM-compatible, low fees |
| Mobile Framework | Native (Swift + Kotlin) or cross-platform (Flutter) | Flutter preferred given Fuse SDK Flutter support |
| Database | PostgreSQL (primary), Redis (cache), Kafka (events) | Proven stack for financial systems |
| Infrastructure | Kubernetes on cloud; Terraform IaC | Scalable, reproducible deployments |

### 7.3 Integration Points

```
+-------------------+       +-------------------+       +-------------------+
|   Mobile Apps     |<----->|   API Gateway     |<----->|   CBS Engine      |
|  (iOS/Android)    |       |  (Auth, Rate      |       |  (Ledger, Accts,  |
+-------------------+       |   Limiting)       |       |   Payments)       |
                            +-------------------+       +-------------------+
                                    |                           |
                    +---------------+---------------+           |
                    |               |               |           |
            +-------v---+   +------v----+   +------v------+    |
            |  Sumsub   |   | Mastercard|   |   Fuse.io   |    |
            |  KYC/AML  |   | Processor |   |  Blockchain |    |
            +-----------+   +-----------+   +-------------+    |
                                                                |
                                                    +-----------v---------+
                                                    | Accounting /        |
                                                    | Regulatory Reporting|
                                                    +---------------------+
```

### 7.4 Data Residency

- All personal data stored in EU (GDPR requirement)
- Primary data center: EU (Frankfurt preferred)
- Disaster recovery: EU (Ireland or Netherlands)
- No personal data transferred outside EEA without adequate safeguards

## 8. Risks and Mitigations

| ID | Risk | Likelihood | Impact | Mitigation |
|----|------|-----------|--------|------------|
| R01 | Mastercard BIN sponsorship delays | Medium | High | Engage multiple issuer processors early; have backup plan |
| R02 | Regulatory changes (PSD3 transition) | High | Medium | Design architecture for adaptability; monitor regulatory pipeline |
| R03 | Fuse.io network instability or low liquidity | Medium | Medium | Implement circuit breakers; fiat operations independent of blockchain |
| R04 | MiCA compliance complexity for crypto features | High | High | Engage crypto-regulatory counsel; phased crypto feature rollout |
| R05 | Sumsub service outage blocking onboarding | Low | High | Implement offline queue; retry mechanism; escalation path |
| R06 | Data breach / security incident | Low | Critical | Defense-in-depth; SOC monitoring; incident response plan; insurance |
| R07 | Mobile app rejection by App Store / Google Play | Low | Medium | Pre-submission review; compliance with platform guidelines |
| R08 | Scalability bottlenecks under growth | Medium | Medium | Load testing at 3x projected capacity; auto-scaling infrastructure |
| R09 | Key person dependency | Medium | High | Document all architecture decisions; cross-train team members |
| R10 | Bank of Lithuania audit findings | Medium | High | Internal compliance team; pre-audit readiness reviews quarterly |

## 9. Assumptions

| ID | Assumption |
|----|------------|
| A01 | TeslaPay's existing EMI license with Bank of Lithuania remains valid and covers planned services |
| A02 | TeslaPay will secure or already has a Mastercard BIN sponsorship arrangement |
| A03 | Fuse.io network will remain operational and maintain sufficient liquidity for planned features |
| A04 | Sumsub API and pricing remain stable throughout development and initial launch |
| A05 | Target markets are EEA countries (no UK initially) |
| A06 | EUR is the primary currency; other currencies are secondary |
| A07 | Development team has access to required sandbox/test environments for all integrations |
| A08 | Existing TeslaPay customer data can be migrated to the new platform |
| A09 | Budget supports native mobile development or Flutter cross-platform approach |
| A10 | MiCA regulatory framework is finalized and applicable guidance is available for crypto features |

## 10. Dependencies

| ID | Dependency | Owner | Impact if Delayed |
|----|-----------|-------|-------------------|
| D01 | Mastercard issuer processor contract | TeslaPay Biz Dev | No card issuance; blocks MC-* requirements |
| D02 | Sumsub contract and API keys | TeslaPay Biz Dev | No onboarding; blocks KYC-* requirements |
| D03 | Fuse.io partnership and SDK access | TeslaPay Biz Dev | No crypto features; blocks FUSE-* requirements |
| D04 | SEPA connectivity (direct or via banking partner) | TeslaPay Biz Dev | No EUR payments; blocks CBS-004/005 |
| D05 | Apple Developer Enterprise account | Engineering | No iOS distribution; blocks MOB-001 |
| D06 | Google Play Developer account | Engineering | No Android distribution; blocks MOB-001 |
| D07 | Cloud infrastructure provisioning | Engineering/DevOps | No deployment environment |
| D08 | CBS platform selection (build vs. buy) | Architecture Team | Blocks all CBS-* requirements |
| D09 | Data migration plan from legacy system | Engineering | Existing customers cannot use new platform |
| D10 | Compliance sign-off on crypto features | Compliance/Legal | Blocks FUSE-* go-live |

## 11. Success Metrics and KPIs

### 11.1 Product KPIs

| Metric | Target (MVP Launch) | Target (6 months) | Target (12 months) |
|--------|--------------------|--------------------|---------------------|
| Registered users | 5,000 | 20,000 | 50,000 |
| Monthly Active Users | 2,000 | 12,000 | 35,000 |
| Cards issued (virtual + physical) | 3,000 | 15,000 | 25,000 |
| Monthly transaction volume (EUR) | 5M | 25M | 100M |
| Fuse wallet activations | 500 | 3,000 | 10,000 |
| Customer acquisition cost (CAC) | < EUR 30 | < EUR 25 | < EUR 20 |
| Net Promoter Score (NPS) | > 30 | > 40 | > 50 |

### 11.2 Operational KPIs

| Metric | Target |
|--------|--------|
| KYC approval rate | > 85% (auto-approved) |
| KYC median completion time | < 3 minutes |
| SEPA payment success rate | > 99.5% |
| Card authorization success rate | > 98% |
| Customer support first response time | < 5 minutes |
| App crash rate | < 0.1% |
| Mean Time to Recovery (MTTR) | < 15 minutes |

## 12. Release Plan

### Phase 1: MVP (Months 1-9)

| Milestone | Target Date | Deliverables |
|-----------|-------------|-------------|
| M1: Architecture & Design | Month 1-2 | System architecture, API specs, DB schema, UI/UX wireframes |
| M2: CBS Core | Month 2-5 | Ledger, accounts, SEPA payments, FX engine |
| M3: KYC Integration | Month 3-4 | Sumsub integration, onboarding flow |
| M4: Card Program | Month 4-7 | Virtual/physical card issuance, Apple/Google Pay |
| M5: Fuse Basic | Month 5-7 | Wallet creation, crypto buy/sell, balance display |
| M6: Mobile Apps | Month 2-8 | iOS + Android apps, full feature integration |
| M7: Testing & Compliance | Month 7-8 | UAT, penetration testing, compliance review |
| M8: Soft Launch | Month 8-9 | Beta with 500 users; iterate based on feedback |
| M9: Public Launch | Month 9 | Full launch; marketing campaign |

### Phase 2: Growth (Months 9-18)

- Business accounts
- Open Banking (PSD2 APIs)
- Advanced DeFi/Solid integration
- Savings products and budgeting tools
- Rewards/cashback program

## 13. Glossary

| Term | Definition |
|------|-----------|
| ABS/CBS | Automated Banking System / Core Banking System |
| AML | Anti-Money Laundering |
| BIN | Bank Identification Number (first 6-8 digits of card number) |
| CRS | Common Reporting Standard (tax reporting) |
| CQRS | Command Query Responsibility Segregation |
| EEA | European Economic Area |
| EMI | Electronic Money Institution |
| FATCA | Foreign Account Tax Compliance Act |
| FATF | Financial Action Task Force |
| FX | Foreign Exchange |
| GDPR | General Data Protection Regulation |
| IBAN | International Bank Account Number |
| KYB | Know Your Business |
| KYC | Know Your Customer |
| MDES | Mastercard Digital Enablement Service |
| MiCA | Markets in Crypto-Assets (EU regulation) |
| NFC | Near-Field Communication |
| PEP | Politically Exposed Person |
| PSD2/PSD3 | Payment Services Directive 2/3 |
| SAR | Suspicious Activity Report |
| SCA | Strong Customer Authentication |
| SCT | SEPA Credit Transfer |
| SDD | SEPA Direct Debit |
| SEPA | Single Euro Payments Area |
| soUSD | Solid yield-bearing stablecoin on Fuse |
| TPS | Transactions Per Second |

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| Product Owner | TBD | | Pending |
| CTO | TBD | | Pending |
| Compliance Officer | TBD | | Pending |
| Business Analyst | Dream Team BA | 2026-03-03 | Submitted |
