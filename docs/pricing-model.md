# TeslaPay Pricing Strategy and Model

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Financial Analyst / CFO, Dream Team
**Status:** Draft for Review

---

## 1. Pricing Philosophy

TeslaPay's pricing strategy is designed around three principles:

1. **Free tier must deliver real value.** Users get a fully functional bank account, virtual Mastercard, SEPA payments, and basic crypto access. This drives adoption and word-of-mouth.
2. **Paid tiers unlock frequency and depth.** Premium and Metal tiers reduce fees on high-frequency activities (FX, crypto, ATM) and unlock exclusive features (DeFi yield, metal card, priority support).
3. **Crypto is the upgrade trigger.** The unique crypto/DeFi features are the primary reason users upgrade. Traditional banking features establish trust; crypto features drive revenue.

---

## 2. Tier Structure

### 2.1 Overview

| Feature | Free | Premium (EUR 7.99/mo) | Metal (EUR 14.99/mo) |
|---------|------|----------------------|---------------------|
| **Annual equivalent** | EUR 0 | EUR 95.88 (or EUR 79.90/yr annual plan) | EUR 179.88 (or EUR 149.90/yr annual plan) |
| **Target segment** | Price-sensitive, evaluating | Active users, multi-currency | Power users, crypto-active, high-spend |
| **Target % of base** | 75% | 20% | 5% |

### 2.2 Account Features

| Feature | Free | Premium | Metal |
|---------|------|---------|-------|
| EUR current account with IBAN | Yes | Yes | Yes |
| Multi-currency accounts | EUR + 1 more | EUR + 5 currencies | EUR + all available |
| SEPA Credit Transfer | Free (unlimited) | Free (unlimited) | Free (unlimited) |
| SEPA Instant | EUR 0.50/transfer | Free (unlimited) | Free (unlimited) |
| Internal TeslaPay transfers | Free | Free | Free |
| Monthly transaction limit | EUR 15,000 | EUR 50,000 | EUR 150,000 |
| Scheduled/recurring payments | 3 active | Unlimited | Unlimited |
| Account statements (PDF) | Monthly | Monthly + on-demand | Monthly + on-demand |

### 2.3 Card Features

| Feature | Free | Premium | Metal |
|---------|------|---------|-------|
| Virtual Mastercard debit | 1 card | Up to 3 cards | Up to 5 cards |
| Physical card | EUR 5.00 (one-time) | Free (included) | Free metal card (included) |
| Card design | Standard | Premium design | Exclusive metal card |
| Apple Pay / Google Pay | Yes | Yes | Yes |
| 3D Secure (in-app) | Yes | Yes | Yes |
| Spending controls | Basic (limits only) | Full (MCC blocks, geo) | Full + auto-rules |
| Contactless (NFC) | Yes | Yes | Yes |
| Card replacement | EUR 10.00 | EUR 5.00 | Free |
| Express card delivery | EUR 25.00 | EUR 15.00 | Free |

### 2.4 FX Features

| Feature | Free | Premium | Metal |
|---------|------|---------|-------|
| FX markup (weekday, major pairs) | 0.50% | 0.25% | 0.10% |
| FX markup (weekday, exotic pairs) | 1.00% | 0.50% | 0.25% |
| FX markup (weekend/holiday) | +0.50% surcharge | +0.25% surcharge | No surcharge |
| Free FX exchange limit/month | EUR 1,000 | EUR 10,000 | Unlimited |
| FX over-limit markup | +0.50% additional | +0.25% additional | N/A |
| Supported currency pairs | 5 major (EUR, USD, GBP, PLN, CHF) | 5 major + 10 additional | All available |

**Major pairs:** EUR/USD, EUR/GBP, EUR/PLN, EUR/CHF, GBP/USD
**Weekend:** Friday 23:00 CET to Sunday 23:00 CET; public holidays per ECB calendar

### 2.5 ATM Features

| Feature | Free | Premium | Metal |
|---------|------|---------|-------|
| Free ATM withdrawals/month | 2 (up to EUR 200) | 5 (up to EUR 500) | Unlimited (up to EUR 1,500) |
| ATM fee (over allowance, domestic) | 2.0% (min EUR 1.00) | 1.5% (min EUR 1.00) | 1.0% (min EUR 1.00) |
| ATM fee (international) | 2.0% + EUR 1.50 | 1.5% + EUR 1.00 | 1.0% (no fixed fee) |
| ATM operator surcharge | Passed through | Passed through | Passed through |

### 2.6 Crypto Features

| Feature | Free | Premium | Metal |
|---------|------|---------|-------|
| Fuse Smart Wallet | Yes | Yes | Yes |
| Supported tokens | FUSE, USDC, USDT | FUSE, USDC, USDT + expanded list | All Fuse tokens |
| Crypto buy/sell fee | 1.50% | 1.00% | 0.75% |
| Crypto buy/sell limit (monthly) | EUR 1,000 | EUR 10,000 | EUR 50,000 |
| Crypto send/receive (on-chain) | Yes (max 3 tx/month) | Yes (unlimited) | Yes (unlimited) |
| Crypto withdrawal to external wallet | EUR 1.00 flat fee | Free (5/month) | Free (unlimited) |
| Gasless transactions (ERC-4337) | No (user pays gas) | Yes (up to 10/month) | Yes (unlimited) |
| DeFi yield (Solid soUSD) -- Phase 2 | Not available | Yes (up to EUR 5,000) | Yes (unlimited) |
| Token swaps -- Phase 2 | Not available | Yes | Yes |
| NFT display -- Phase 2 | Not available | Not available | Yes |

### 2.7 Support and Other

| Feature | Free | Premium | Metal |
|---------|------|---------|-------|
| In-app chat support | Standard (24h SLA) | Priority (2h SLA) | Priority (30min SLA) |
| Phone support | Not available | Callback within 4h | Direct line, 24/7 |
| Dedicated account manager | No | No | Yes (100+ Metal users, then scaled) |
| Financial insights / analytics | Basic spending summary | Full analytics + budgeting | Full + tax export |
| Dark mode | Yes | Yes | Yes |
| Multi-language | Yes (5 languages) | Yes (5 languages) | Yes (5 languages) |

---

## 3. Crypto Fee Structure (Detail)

### 3.1 Buy/Sell Spread

The crypto buy/sell fee is applied as a markup on the market price:

| Token Pair | Free Tier | Premium | Metal |
|-----------|-----------|---------|-------|
| EUR -> USDC | 1.50% | 1.00% | 0.75% |
| EUR -> USDT | 1.50% | 1.00% | 0.75% |
| EUR -> FUSE | 1.50% | 1.00% | 0.75% |
| USDC -> EUR | 1.50% | 1.00% | 0.75% |
| USDT -> EUR | 1.50% | 1.00% | 0.75% |
| FUSE -> EUR | 1.50% | 1.00% | 0.75% |

**Spread calculation:** TeslaPay adds the fee percentage to the mid-market rate sourced from CoinGecko/CoinMarketCap aggregated feeds. The effective user price is displayed before confirmation and locked for 30 seconds.

### 3.2 On-Chain Transaction Fees

| Operation | Free Tier | Premium | Metal |
|-----------|-----------|---------|-------|
| Send USDC/USDT on Fuse | EUR 0.50 flat | Free (gasless, 10/mo) | Free (gasless, unlimited) |
| Send FUSE on Fuse | Network gas (~EUR 0.001) | Network gas (~EUR 0.001) | Free (gasless, unlimited) |
| Withdraw to external wallet | EUR 1.00 flat | Free (5/mo), then EUR 0.50 | Free (unlimited) |
| Receive crypto | Free | Free | Free |

### 3.3 Phase 2: DeFi / Yield

| Feature | Premium | Metal |
|---------|---------|-------|
| Solid soUSD deposit | 0.50% entry fee | No entry fee |
| Solid soUSD withdrawal | Free (7-day lock) | Free (no lock) |
| Token swap fee | 0.75% per swap | 0.50% per swap |
| Maximum DeFi deposit | EUR 5,000 | Unlimited |

---

## 4. FX Fee Structure (Detail)

### 4.1 Fee Tiers

| Category | Definition | Free | Premium | Metal |
|----------|-----------|------|---------|-------|
| **Major pairs (weekday)** | EUR vs USD, GBP, CHF, PLN, SEK, NOK, DKK | 0.50% | 0.25% | 0.10% |
| **Minor pairs (weekday)** | EUR vs CZK, HUF, RON, BGN, HRK | 0.75% | 0.40% | 0.20% |
| **Exotic pairs (weekday)** | All other supported pairs | 1.00% | 0.50% | 0.25% |
| **Weekend surcharge** | Applied on top of weekday rate | +0.50% | +0.25% | +0.00% |
| **Free allowance** | No markup applied up to this amount/month | EUR 1,000 | EUR 10,000 | Unlimited |
| **Over-limit surcharge** | Applied after free allowance exhausted | +0.50% | +0.25% | N/A |

### 4.2 FX Revenue Example

A Premium user exchanging EUR 2,000 to USD on a Tuesday:
- First EUR 10,000 at 0.25% markup
- Markup revenue: EUR 2,000 x 0.25% = EUR 5.00
- TeslaPay cost (rate provider): ~EUR 0.50
- Net FX revenue: EUR 4.50

### 4.3 Weekend Rate Justification

Weekend and holiday surcharges exist because:
1. Interbank FX markets are closed; rates carry higher risk
2. TeslaPay takes on spread risk between Friday close and Monday open
3. Industry standard (Revolut, Wise apply similar surcharges)

Metal tier absorbs this risk as a premium benefit, which is a key upgrade driver for frequent FX users.

---

## 5. ATM Fee Structure (Detail)

### 5.1 Fee Schedule

| Scenario | Free | Premium | Metal |
|----------|------|---------|-------|
| Domestic, within allowance | Free | Free | Free |
| Domestic, over allowance | 2.0% (min EUR 1.00) | 1.5% (min EUR 1.00) | 1.0% (min EUR 1.00) |
| International (non-EUR), within allowance | Free (uses FX rate) | Free (uses FX rate) | Free (uses FX rate) |
| International, over allowance | 2.0% + EUR 1.50 fixed | 1.5% + EUR 1.00 fixed | 1.0% (no fixed fee) |
| Operator surcharge (any) | Passed to user | Passed to user | Passed to user |

### 5.2 Monthly Allowance

| Tier | Free Withdrawals | Free Amount | Resets |
|------|-----------------|-------------|--------|
| Free | 2 per month | Up to EUR 200 total | 1st of each month |
| Premium | 5 per month | Up to EUR 500 total | 1st of each month |
| Metal | Unlimited | Up to EUR 1,500 total | 1st of each month |

**Note:** "Unlimited" for Metal means no per-withdrawal count limit, but the EUR 1,500 monthly cap still applies to prevent ATM arbitrage.

---

## 6. Card Fee Schedule

| Fee | Free | Premium | Metal |
|-----|------|---------|-------|
| First virtual card | Free | Free | Free |
| Additional virtual cards | EUR 2.00 each | Free (up to 3) | Free (up to 5) |
| First physical card | EUR 5.00 | Free | Free (metal) |
| Physical card replacement (lost/stolen) | EUR 10.00 | EUR 5.00 | Free |
| Physical card replacement (damaged) | EUR 5.00 | Free | Free |
| Express delivery (2-day) | EUR 25.00 | EUR 15.00 | Free |
| Standard delivery (7-10 days) | Free (with card order) | Free | Free |
| Inactive card fee (no tx > 12 months) | EUR 3.00/month | No fee | No fee |

---

## 7. Account Fees

| Fee | Free | Premium | Metal |
|-----|------|---------|-------|
| Account opening | Free | Free | Free |
| Account maintenance | Free | EUR 7.99/mo | EUR 14.99/mo |
| Account closure | Free | Free | Free |
| Inactive account (no login > 12 months) | EUR 5.00/month | EUR 5.00/month | No fee |
| Paper statement (postal) | EUR 5.00 each | EUR 5.00 each | Free |
| SWIFT transfer (if supported, Phase 2) | EUR 5.00 + 0.3% | EUR 3.00 + 0.2% | EUR 1.00 + 0.1% |

---

## 8. Competitive Pricing Comparison

### 8.1 Subscription Pricing

| Neobank | Free Tier | Mid Tier | Top Tier |
|---------|-----------|----------|----------|
| **TeslaPay** | **EUR 0** | **EUR 7.99/mo** | **EUR 14.99/mo** |
| Revolut | EUR 0 | EUR 9.99/mo (Plus) | EUR 16.99/mo (Ultra) |
| N26 | EUR 0 | EUR 9.90/mo (Smart) | EUR 16.90/mo (Metal) |
| Crypto.com | USD 0 | -- | -- (staking tiers) |
| Wise | EUR 0 | -- (no tiers) | -- |

TeslaPay is priced 15-20% below Revolut and N26, reflecting its smaller brand and the need to drive initial adoption. Pricing can be raised after establishing product-market fit.

### 8.2 FX Fees

| Neobank | Free Tier Weekday | Mid Tier Weekday | Top Tier Weekday | Weekend Surcharge |
|---------|-------------------|------------------|------------------|-------------------|
| **TeslaPay** | **0.50%** | **0.25%** | **0.10%** | **+0.50% / +0.25% / +0.00%** |
| Revolut | 0.00% (up to EUR 1K) | 0.00% (up to EUR 5K) | 0.00% (unlimited) | +1.00% / +0.50% / +0.00% |
| N26 | 0.00% (Mastercard rate) | 0.00% | 0.00% | N/A (Mastercard rate) |
| Wise | 0.33-0.62% flat | -- | -- | No surcharge |

**Note:** TeslaPay's FX rates are competitive but not as aggressive as Revolut's free-tier offer. This is intentional -- FX is a significant revenue stream. N26 passes Mastercard rates which are typically 0.2-0.5% depending on the pair.

### 8.3 Crypto Fees

| Neobank | Buy/Sell Fee | Self-Custody | DeFi/Yield | Tokens Supported |
|---------|-------------|-------------|------------|-----------------|
| **TeslaPay** | **1.50% / 1.00% / 0.75%** | **Yes (Fuse wallet)** | **Yes (Solid soUSD)** | **3+ (Fuse network)** |
| Revolut | 1.49% (free) / 0.99% (premium) | No | No | 200+ |
| N26 | ~1.5% via partner | No | No | Limited |
| Crypto.com | 0-2.99% | Yes (DeFi wallet) | Yes (Earn) | 250+ |

TeslaPay's crypto fees are competitive with Revolut and significantly cheaper than Crypto.com's high-tier fees. The self-custody and DeFi yield differentiation justifies parity pricing.

### 8.4 ATM Fees

| Neobank | Free Tier Allowance | Over-Limit Fee | International |
|---------|--------------------|--------------|--------------|
| **TeslaPay** | **2x, up to EUR 200** | **2.0%** | **2.0% + EUR 1.50** |
| Revolut | 5x, up to EUR 200 | 2.0% | 2.0% |
| N26 | 3-5x/month | EUR 2.00 flat | EUR 2.00 + FX |
| Wise | 2x, up to EUR 200 | EUR 0.50 + 1.75% | EUR 0.50 + 1.75% |

TeslaPay's ATM allowance is slightly below Revolut's but comparable to Wise. The fee structure is standard for the industry.

---

## 9. Pricing Strategy Rationale

### 9.1 Why EUR 7.99 and EUR 14.99

| Factor | Reasoning |
|--------|-----------|
| **Psychological pricing** | Sub-EUR 8 and sub-EUR 15 price points reduce perceived cost |
| **Competitor discount** | 15-20% below Revolut/N26 compensates for lower brand awareness |
| **Margin positive at scale** | Premium at EUR 7.99 generates EUR 11.70 contribution/mo (cost EUR 1.74); Metal at EUR 14.99 generates EUR 23.71 (cost EUR 2.36) |
| **Annual discount incentive** | EUR 79.90/yr (Premium) and EUR 149.90/yr (Metal) = ~17% discount for annual commitment, reducing churn |

### 9.2 Free-to-Paid Conversion Strategy

| Trigger | Mechanism | Expected Impact |
|---------|-----------|----------------|
| FX limit hit (EUR 1,000/mo) | In-app upgrade prompt when approaching limit | 5-8% of free users convert |
| Crypto limit hit (EUR 1,000/mo) | In-app upsell showing fee savings | 3-5% of crypto-active free users |
| ATM limit hit | Prompt at next withdrawal | 2-3% of ATM-active free users |
| Physical card desire | Card order flow shows Premium includes free card | 4-6% of free users |
| DeFi yield availability (Phase 2) | Marketing push for yield access | 5-10% of crypto-active users |
| Trial offer | 1 month free Premium trial at sign-up | 15-20% trial-to-paid conversion |

### 9.3 Price Increase Path

TeslaPay should plan for price increases after establishing market position:

| Timeline | Action | Rationale |
|----------|--------|-----------|
| Launch to Month 12 | EUR 7.99 / EUR 14.99 | Establish base, acquire users |
| Month 12-18 | Introduce annual plans at discount | Lock in users, reduce churn |
| Month 18-24 | Evaluate increase to EUR 9.99 / EUR 16.99 | Match competitors if NPS > 40 and churn < 2.5% |
| Month 24+ | Introduce Business tier (EUR 19.99/mo) | Capture SME segment from Phase 2 |

---

## 10. Revenue Impact Modeling

### 10.1 Revenue by Fee Type (50,000 Users, 12-Month Target)

| Fee Type | Users Affected | Avg Revenue/User/Mo | Monthly Revenue | % of Total |
|----------|---------------|--------------------|-----------------|-----------|
| Interchange (cards) | 35,000 active card users | EUR 1.10 | EUR 38,500 | 28% |
| FX markup | 15,000 FX users | EUR 1.80 | EUR 27,000 | 20% |
| Crypto fees | 8,000 crypto users | EUR 1.50 | EUR 12,000 | 9% |
| Premium subscriptions | 10,000 | EUR 7.99 | EUR 79,900 | 28% |
| Metal subscriptions | 2,500 | EUR 14.99 | EUR 37,475 | 13% |
| ATM over-limit | 5,000 over-limit events | EUR 1.50 avg | EUR 7,500 | 5% |
| Card fees | 500 orders/replacements | EUR 8.00 avg | EUR 4,000 | 3% |
| Interest on safeguarded funds | 50,000 | EUR 0.35 | EUR 17,500 | -- |
| **Total Monthly Revenue** | | | **EUR 223,875** | |
| **Annualized** | | | **EUR 2,686,500** | |

**Note:** Interest on safeguarded funds (EUR 17,500/mo at EUR 50M deposits x 3.5% / 12) is technically treasury revenue and may be reported separately from user-facing revenue.

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CEO | TBD | | Pending |
| CFO | Dream Team Financial Analyst | 2026-03-03 | Submitted |
| Product Owner | TBD | | Pending |
