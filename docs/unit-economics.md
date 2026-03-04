# TeslaPay Unit Economics Analysis

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Financial Analyst / CFO, Dream Team
**Status:** Draft for Review

---

## 1. Executive Summary

This document models TeslaPay's unit economics across three account tiers (Free, Premium, Metal), projecting revenue per user, cost to serve, and lifetime value. All figures are in EUR unless stated otherwise. Conservative assumptions are applied throughout; optimistic scenarios are flagged separately.

**Key Findings:**
- Blended ARPU (Year 1): EUR 4.20/month (EUR 50.40/year)
- Blended cost to serve: EUR 2.05/month (EUR 24.60/year)
- Gross margin per user: EUR 2.15/month (51%)
- Target CAC: EUR 25 (blended)
- Estimated LTV (36-month horizon): EUR 62.10
- LTV:CAC ratio: 2.5:1 (Year 1 cohort), improving to 3.5:1+ at scale
- Payback period: ~12 months

---

## 2. Revenue Streams

### 2.1 Revenue Stream Breakdown

| # | Revenue Stream | Description | Estimated Contribution (Steady State) |
|---|----------------|-------------|--------------------------------------|
| 1 | Card interchange | 0.20% of consumer debit transactions (EU IFR cap) | 30-35% |
| 2 | FX markup | 0.30-0.50% spread on currency exchanges | 15-20% |
| 3 | Crypto trading fees | 1.0-1.5% on buy/sell transactions | 10-15% |
| 4 | Premium subscriptions | EUR 7.99/mo (Premium), EUR 14.99/mo (Metal) | 20-25% |
| 5 | ATM fees | 2% over free allowance | 3-5% |
| 6 | Interest on safeguarded funds | ECB deposit facility rate on pooled funds | 5-8% |
| 7 | Card issuance fees | Physical card replacements, express delivery | 2-3% |
| 8 | Other fees | Inactive account fees, expedited support | 1-2% |

### 2.2 Revenue Per User (ARPU) by Tier -- Monthly

**Assumptions:**
- Average card spend: Free EUR 400/mo, Premium EUR 900/mo, Metal EUR 1,500/mo
- FX usage: Free 1 exchange/mo (EUR 200), Premium 3/mo (EUR 500), Metal 5/mo (EUR 1,000)
- Crypto activity: Free 10% of users buy EUR 100/mo, Premium 25% buy EUR 300/mo, Metal 40% buy EUR 500/mo
- Average deposit balance: Free EUR 500, Premium EUR 2,000, Metal EUR 5,000
- ATM withdrawals: Free 2/mo, Premium 3/mo, Metal 4/mo (most within free allowance)

| Revenue Component | Free Tier | Premium (EUR 7.99/mo) | Metal (EUR 14.99/mo) |
|-------------------|-----------|----------------------|---------------------|
| Card interchange (0.20%) | EUR 0.80 | EUR 1.80 | EUR 3.00 |
| FX markup (avg 0.40%) | EUR 0.80 | EUR 2.00 | EUR 4.00 |
| Crypto fees (1.25% avg) | EUR 0.13 | EUR 0.94 | EUR 2.50 |
| Subscription fee | EUR 0.00 | EUR 7.99 | EUR 14.99 |
| ATM over-limit fees | EUR 0.10 | EUR 0.05 | EUR 0.02 |
| Interest on deposits (3.5% ECB rate) | EUR 0.15 | EUR 0.58 | EUR 1.46 |
| Card/other fees (avg) | EUR 0.05 | EUR 0.08 | EUR 0.10 |
| **Total Monthly ARPU** | **EUR 2.03** | **EUR 13.44** | **EUR 26.07** |
| **Annualized ARPU** | **EUR 24.36** | **EUR 161.28** | **EUR 312.84** |

### 2.3 Blended ARPU Calculation

**Assumed tier distribution (12 months post-launch):**
- Free: 75% of users
- Premium: 20% of users
- Metal: 5% of users

| Metric | Calculation | Result |
|--------|-------------|--------|
| Blended monthly ARPU | (0.75 x 2.03) + (0.20 x 13.44) + (0.05 x 26.07) | **EUR 4.51** |
| Blended annual ARPU | 4.51 x 12 | **EUR 54.12** |

**Benchmark Comparison:**

| Neobank | Estimated ARPU (Annual) | Source/Basis |
|---------|------------------------|--------------|
| Revolut | ~EUR 59 (GBP 59, 2024 annual report) | GBP 3.1B / 52.5M users |
| N26 | ~EUR 47 (estimated, 2024) | EUR 380M / 8M users |
| Wise | ~EUR 55 (estimated) | GBP 1.0B / 16M users |
| **TeslaPay (target)** | **EUR 54** | Model above |

TeslaPay's blended ARPU target of EUR 54 is achievable but requires the subscription conversion rates modeled above. Revolut's higher ARPU is driven by its wealth products (stocks/crypto), business accounts, and larger average card spend due to more mature user base.

---

## 3. Cost Structure

### 3.1 Variable Cost Per User -- Monthly

| Cost Component | Free Tier | Premium | Metal | Assumption |
|----------------|-----------|---------|-------|------------|
| **Card processing (Enfuce)** | | | | |
| - Per-transaction fee | EUR 0.12 | EUR 0.20 | EUR 0.30 | EUR 0.03/txn; Free 4 txns/mo avg, Premium 7, Metal 10 |
| - Monthly card maintenance | EUR 0.15 | EUR 0.15 | EUR 0.15 | Enfuce per-card fee ~EUR 0.15/mo estimated |
| - Virtual card issuance (amortized) | EUR 0.02 | EUR 0.02 | EUR 0.02 | ~EUR 0.50 one-time, amortized 24 months |
| **SEPA payments (Banking Circle)** | | | | |
| - Outgoing SCT | EUR 0.04 | EUR 0.08 | EUR 0.12 | EUR 0.20/txn; Free 0.2/mo, Premium 0.4/mo, Metal 0.6/mo |
| - Incoming SCT (credit) | EUR 0.02 | EUR 0.04 | EUR 0.06 | EUR 0.10/incoming; salary credits etc. |
| - SCT Instant surcharge | EUR 0.02 | EUR 0.05 | EUR 0.10 | EUR 0.50/instant; ~4% of payments are instant |
| **KYC (Sumsub)** | | | | |
| - Onboarding verification (amortized) | EUR 0.06 | EUR 0.06 | EUR 0.06 | EUR 1.50/check, amortized over 24 months |
| - Ongoing AML monitoring | EUR 0.05 | EUR 0.05 | EUR 0.05 | EUR 0.60/year per user |
| **Blockchain (Fuse.io)** | | | | |
| - Gas fees (paymaster) | EUR 0.01 | EUR 0.02 | EUR 0.05 | ~EUR 0.001/txn; minimal for most users |
| - FuseBox API | EUR 0.02 | EUR 0.03 | EUR 0.05 | Estimated API costs at scale |
| **Infrastructure (AWS, allocated)** | EUR 0.20 | EUR 0.20 | EUR 0.20 | Total infra cost / total users (see Section 3.2) |
| **Customer support** | EUR 0.25 | EUR 0.35 | EUR 0.50 | 0.5 tickets/mo free, 0.3 premium, 0.2 metal; EUR 3-5/ticket |
| **Fraud & chargebacks** | EUR 0.08 | EUR 0.10 | EUR 0.12 | 0.1% of card volume |
| **FX provider costs** | EUR 0.05 | EUR 0.12 | EUR 0.20 | Rate provider fees + spread costs |
| **Physical card (amortized)** | EUR 0.00 | EUR 0.15 | EUR 0.25 | Free: virtual only default; Premium: EUR 3.50 card / 24 mo; Metal: EUR 6.00 / 24 mo |
| **Compliance & reporting** | EUR 0.10 | EUR 0.10 | EUR 0.10 | Regulatory reporting tools, audit, legal allocated |
| **Push notifications (APNs/FCM)** | EUR 0.01 | EUR 0.02 | EUR 0.03 | Negligible at scale |
| **Total Variable Cost/Mo** | **EUR 1.20** | **EUR 1.74** | **EUR 2.36** |

### 3.2 Fixed Costs -- Monthly (Pre-Scale)

| Category | Monthly Cost (EUR) | Notes |
|----------|-------------------|-------|
| **Cloud Infrastructure (AWS)** | | |
| - EKS clusters (prod + staging + DR) | 4,500 | 3 clusters x ~EUR 1,500 |
| - RDS PostgreSQL instances (10 DBs) | 12,000 | Mix of xlarge/large Multi-AZ |
| - MSK Kafka cluster | 3,500 | 3 brokers, m5.2xlarge |
| - ElastiCache Redis | 2,500 | 3-shard cluster with replicas |
| - OpenSearch | 2,000 | 3-node cluster |
| - S3, CloudFront, Route53 | 1,000 | Storage + CDN |
| - NAT Gateway, data transfer | 2,000 | Cross-AZ and outbound |
| - WAF, GuardDuty, Security Hub | 1,500 | Security services |
| - Vault, KMS, CloudHSM | 2,500 | Secrets and key management |
| - Misc (ECR, CloudWatch, etc.) | 1,000 | Various support services |
| **Subtotal AWS** | **32,500** | |
| **Third-party SaaS** | | |
| - Sumsub (minimum commitment) | 2,000 | Monthly minimum + per-check above minimum |
| - FX rate provider | 500 | CurrencyLayer or similar |
| - SonarQube, monitoring tools | 500 | Code quality |
| - Intercom or similar (support) | 1,000 | Customer support platform |
| **Subtotal SaaS** | **4,000** | |
| **Team Costs** | | |
| - Engineering (8 engineers) | 64,000 | Avg EUR 8,000/mo gross per engineer (EU rates) |
| - DevOps / SRE (2) | 18,000 | EUR 9,000/mo avg |
| - Product / Design (2) | 14,000 | EUR 7,000/mo avg |
| - Compliance officer (1) | 8,000 | |
| - Customer support (3) | 12,000 | EUR 4,000/mo avg |
| - Finance / Operations (1) | 7,000 | |
| - Management / CEO (1) | 12,000 | |
| **Subtotal Team** | **135,000** | 18 headcount |
| **Other Fixed Costs** | | |
| - Office / coworking | 4,000 | Vilnius office |
| - Legal / regulatory | 5,000 | Ongoing legal counsel |
| - Insurance (D&O, cyber) | 2,000 | Monthly amortization |
| - Accounting / audit | 2,000 | External accounting |
| - Mastercard scheme fees | 3,000 | Annual BIN fees, scheme participation |
| - Bank of Lithuania reporting | 1,000 | Regulatory fees |
| **Subtotal Other** | **17,000** | |
| **Total Fixed Costs/Month** | **188,500** | |

### 3.3 Blended Cost to Serve Per User

| User Base Size | Fixed Cost/User/Mo | Variable Cost/User/Mo | Total Cost/User/Mo |
|----------------|-------------------|-----------------------|-------------------|
| 5,000 (launch) | EUR 37.70 | EUR 1.33 | EUR 39.03 |
| 10,000 | EUR 18.85 | EUR 1.33 | EUR 20.18 |
| 25,000 | EUR 7.54 | EUR 1.33 | EUR 8.87 |
| 50,000 (12-mo target) | EUR 3.77 | EUR 1.33 | EUR 5.10 |
| 100,000 | EUR 2.26 | EUR 1.33 | EUR 3.59 |
| 200,000 | EUR 1.32 | EUR 1.33 | EUR 2.65 |

**Note:** Fixed costs increase at ~15% for every doubling of users beyond 50K (more support staff, additional infra). The table above holds fixed costs constant for simplicity at the initial levels.

---

## 4. Customer Acquisition Cost (CAC)

### 4.1 CAC Model

| Channel | Cost/Acquisition | % of Users | Weighted CAC |
|---------|-----------------|------------|-------------|
| Organic / word-of-mouth | EUR 0 | 20% | EUR 0.00 |
| Referral program (EUR 10 reward each) | EUR 20 | 25% | EUR 5.00 |
| Social media ads (Meta, TikTok) | EUR 35 | 25% | EUR 8.75 |
| Google Ads (search) | EUR 50 | 15% | EUR 7.50 |
| Content marketing / SEO | EUR 15 | 10% | EUR 1.50 |
| Partnerships / crypto communities | EUR 20 | 5% | EUR 1.00 |
| **Blended CAC** | | **100%** | **EUR 23.75** |

**Target:** CAC < EUR 25 (PRD target: EUR 20-30 depending on maturity).

**Benchmark Comparison:**

| Neobank | CAC (Approximate) | Notes |
|---------|-------------------|-------|
| Revolut | EUR 35-40 | 2024; higher due to global expansion |
| N26 | EUR 40-50 | Higher due to brand campaigns |
| Monzo | EUR 5-10 | UK-focused, strong word-of-mouth |
| Wise | EUR 9-12 | Heavy word-of-mouth (~66%) |
| **TeslaPay (target)** | **EUR 24** | Crypto community + referrals |

TeslaPay's crypto-banking niche enables lower CAC than broad-market neobanks (N26, Revolut) because the crypto community is concentrated in online channels with lower CPA and higher virality.

### 4.2 CAC Payback Period

| Scenario | Blended ARPU/Mo | Blended Variable Cost/Mo | Contribution/Mo | CAC | Payback (Months) |
|----------|----------------|-------------------------|-----------------|-----|------------------|
| Conservative | EUR 4.00 | EUR 1.33 | EUR 2.67 | EUR 25 | 9.4 |
| Base | EUR 4.51 | EUR 1.33 | EUR 3.18 | EUR 24 | 7.5 |
| Optimistic | EUR 5.20 | EUR 1.33 | EUR 3.87 | EUR 22 | 5.7 |

---

## 5. Lifetime Value (LTV) Analysis

### 5.1 Assumptions

| Parameter | Value | Source/Rationale |
|-----------|-------|-----------------|
| Average customer lifespan | 36 months | Conservative; industry average 3-5 years |
| Monthly churn rate | 3.0% | Industry: 2-5% for neobanks |
| Annual churn rate | ~31% | Derived from monthly churn |
| Discount rate (monthly) | 0.83% | 10% annual WACC |
| Revenue growth per user/year | 10% | Users increase spend over time |

### 5.2 LTV Calculation

**LTV = ARPU x Gross Margin x (1 / churn rate)**

Using the contribution margin approach:

| Tier | Monthly Contribution | Monthly Churn | LTV | % of Users | Weighted LTV |
|------|---------------------|---------------|-----|-----------|-------------|
| Free | EUR 0.83 | 4.0% | EUR 20.75 | 75% | EUR 15.56 |
| Premium | EUR 11.70 | 2.0% | EUR 585.00 | 20% | EUR 117.00 |
| Metal | EUR 23.71 | 1.5% | EUR 1,580.67 | 5% | EUR 79.03 |
| **Blended** | | | | | **EUR 211.59** |

**Discounted LTV (36-month horizon, 10% annual discount):**

| Tier | Discounted LTV (36 mo) |
|------|----------------------|
| Free | EUR 17.80 |
| Premium | EUR 283.50 |
| Metal | EUR 595.20 |
| **Blended (weighted)** | **EUR 83.20** |

**Note:** The infinite-horizon LTV is significantly higher than the 36-month discounted version. For conservative planning, we use the 36-month discounted LTV.

### 5.3 LTV:CAC Ratio

| Scenario | LTV (36-mo discounted) | CAC | LTV:CAC |
|----------|----------------------|-----|---------|
| Conservative | EUR 62.10 | EUR 28 | 2.2:1 |
| Base | EUR 83.20 | EUR 24 | 3.5:1 |
| Optimistic | EUR 105.00 | EUR 20 | 5.3:1 |

**Industry benchmark:** 3.5:1 is considered healthy for neobanks. TeslaPay's base case meets this threshold.

**Risk:** If Premium/Metal conversion remains below 15% combined, the blended LTV drops below EUR 50, pushing LTV:CAC below 2:1. Subscription conversion is the single most important unit economics lever.

---

## 6. Breakeven Analysis

### 6.1 User-Level Breakeven

A user becomes profitable (recovers their CAC) at the following timelines:

| Tier | Monthly Contribution | CAC Recovery (Months) |
|------|---------------------|----------------------|
| Free | EUR 0.83 | 29 months (many never recover) |
| Premium | EUR 11.70 | 2.1 months |
| Metal | EUR 23.71 | 1.0 month |
| Blended | EUR 3.18 | 7.5 months |

**Critical insight:** Free-tier users are unprofitable in isolation. They must convert to paid tiers or generate sufficient interchange revenue to justify their CAC. The referral value of free users (bringing in paid users) is not captured in this direct calculation but is material.

### 6.2 Company-Level Breakeven

| Metric | Value | Assumptions |
|--------|-------|-------------|
| Monthly fixed costs | EUR 188,500 | See Section 3.2 |
| Blended contribution/user/mo | EUR 3.18 | See Section 4.2 |
| **Users needed for breakeven** | **59,280** | 188,500 / 3.18 |
| Target to reach 59K users | Month 14-16 | Based on growth model |

**Sensitivity to key variables:**

| Variable Change | New Breakeven Point |
|-----------------|-------------------|
| Premium conversion 25% (vs 20%) | 48,500 users |
| Premium conversion 15% (vs 20%) | 76,200 users |
| Fixed costs +20% | 71,100 users |
| ARPU +15% | 49,800 users |
| ARPU -15% | 73,500 users |

---

## 7. Industry Benchmark Comparison

| Metric | TeslaPay (Target) | Revolut (2024) | N26 (Est. 2024) | Industry Avg |
|--------|-------------------|----------------|------------------|-------------|
| Annual ARPU | EUR 54 | EUR 59 | EUR 47 | EUR 40-60 |
| Blended CAC | EUR 24 | EUR 35-40 | EUR 40-50 | EUR 20-50 |
| LTV:CAC | 3.5:1 | 3.3:1 | ~2.5:1 | 3.0-4.0:1 |
| Gross margin | 51% | ~55% | ~45% | 40-55% |
| Paid tier % | 25% | ~15% | ~12% | 10-20% |
| Monthly churn | 3.0% | ~2.5% | ~3.0% | 2-5% |
| Payback months | 7.5 | ~8 | ~12 | 6-18 |
| Interchange as % rev | 30-35% | ~22% | ~35% | 25-40% |

**Key observations:**
1. TeslaPay's target paid-tier conversion (25%) is aggressive vs. industry norms (10-20%). The crypto/DeFi value proposition must drive this.
2. Lower CAC than Revolut/N26 is achievable if crypto community targeting works but is not guaranteed.
3. Interchange-heavy revenue mix is normal for early-stage neobanks; diversification (subscriptions, crypto) is critical for long-term margin improvement.

---

## 8. Key Assumptions and Risks

### 8.1 Assumptions Register

| ID | Assumption | Sensitivity | Risk if Wrong |
|----|-----------|-------------|---------------|
| A1 | 25% of users convert to paid tiers within 12 months | HIGH | If only 15%, breakeven extends to 76K users |
| A2 | Average card spend: EUR 400/mo (Free), EUR 900/mo (Premium) | MEDIUM | Lower spend reduces interchange; 20% reduction = EUR 0.40 less ARPU |
| A3 | Crypto trading adoption: 10-40% by tier | HIGH | If crypto market enters bear cycle, adoption may drop to 5-15% |
| A4 | EU interchange cap remains at 0.20% (debit) | LOW | Regulatory risk; any reduction directly impacts largest revenue stream |
| A5 | ECB rate stays above 3.0% for interest income | MEDIUM | Rate cuts reduce interest on safeguarded funds |
| A6 | Monthly churn rate: 3.0% blended | HIGH | If churn exceeds 5%, LTV drops 40% |
| A7 | Enfuce per-transaction cost: EUR 0.03 avg | LOW | Contract-dependent; could be EUR 0.02-0.05 |
| A8 | Team size of 18 is sufficient for first 12 months | MEDIUM | If more hires needed, fixed costs rise 5-10% per head |
| A9 | CAC achievable at EUR 24 blended | MEDIUM | If organic/referral channels underperform, CAC could reach EUR 35+ |
| A10 | Fuse.io gas costs remain below EUR 0.001/tx | LOW | Fuse network is low-cost; minimal risk |

### 8.2 Financial Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|------------|------------|
| Low subscription conversion | Breakeven delayed 6-12 months | Medium | Test pricing, A/B onboarding, feature gating |
| Crypto bear market | 30-50% reduction in crypto revenue | Medium | Crypto revenue is 10-15% of total; banking core sustains business |
| Interchange rate reduction | 10-15% revenue reduction | Low | Diversify to subscriptions and FX; lobby via industry groups |
| ECB rate cuts | EUR 0.10-0.30 ARPU reduction | Medium | Build fee-based revenue to reduce interest dependency |
| Higher-than-expected churn | LTV:CAC falls below 2:1 | Medium | Invest in engagement features; loyalty programs in Phase 2 |
| Regulatory compliance costs | EUR 50-100K additional annual cost | Medium | Budget contingency; engage specialized compliance counsel |

---

## 9. Recommendations

1. **Prioritize subscription conversion.** This is the single most impactful lever. Free-to-Premium conversion should be the #1 product metric after MAU.

2. **Gate crypto features behind Premium/Metal tiers.** Allow basic buy/sell on Free with higher fees (1.5%); lower fees and DeFi yield for paid tiers (0.75-1.0%). This creates a natural upgrade path.

3. **Invest in referral program early.** EUR 20 cost per referred user (EUR 10 each) is the most cost-efficient acquisition channel outside organic. Target crypto communities on Twitter/X, Discord, and Telegram.

4. **Monitor CAC weekly.** If blended CAC exceeds EUR 30, pause paid acquisition and focus on organic/referral channels until product-market fit is validated.

5. **Negotiate volume-based Enfuce pricing.** Card processing is the largest variable cost. A 30% reduction in Enfuce fees (from EUR 0.03 to EUR 0.02/txn) improves gross margin by 2-3 percentage points.

6. **Build toward 50K users in 12 months.** This is the minimum scale where unit economics become sustainable. Below 25K users, the business burns cash at EUR 120-150K/month net.

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CEO | TBD | | Pending |
| CFO | Dream Team Financial Analyst | 2026-03-03 | Submitted |
| Board | TBD | | Pending |
