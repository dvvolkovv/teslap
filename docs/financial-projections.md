# TeslaPay 3-Year Financial Projections

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Financial Analyst / CFO, Dream Team
**Status:** Draft for Review

---

## 1. Key Assumptions

### 1.1 General Assumptions

| Assumption | Value | Basis |
|-----------|-------|-------|
| Currency | EUR | Primary operating currency |
| Projection period | 36 months (3 years) | Industry standard for seed/Series A planning |
| MVP development phase | Months 1-9 | Per PRD release plan |
| Soft launch (beta) | Month 8-9 | 500 beta users |
| Public launch | Month 10 | Full marketing begins |
| Tax rate | 15% | Lithuanian corporate tax rate |
| Inflation adjustment | Not applied | Short projection window |
| ECB deposit facility rate | 3.25% (Y1), 2.75% (Y2), 2.50% (Y3) | Conservative declining rate assumption |

### 1.2 User Growth Assumptions

| Parameter | Value | Basis |
|-----------|-------|-------|
| Beta users (Month 8-9) | 500 | Soft launch target |
| Launch month registrations (M10) | 2,000 | Marketing push + PR |
| Monthly organic growth rate (M10-18) | 25% MoM | Aggressive early growth, crypto community |
| Monthly organic growth rate (M19-24) | 15% MoM | Growth deceleration |
| Monthly organic growth rate (M25-36) | 8% MoM | Maturation |
| Activation rate (registered to active) | 65% | Industry benchmark 60-70% |
| Monthly churn (active users) | 3.5% (Y1), 3.0% (Y2), 2.5% (Y3) | Improving with product maturity |
| Tier distribution (Y1): Free/Premium/Metal | 78% / 18% / 4% | Conservative; target 75/20/5 by M12 |
| Tier distribution (Y2) | 72% / 22% / 6% | Improving conversion |
| Tier distribution (Y3) | 68% / 25% / 7% | Mature conversion |

### 1.3 Revenue Assumptions

| Parameter | Value | Basis |
|-----------|-------|-------|
| Card adoption rate (% of active users) | 60% (Y1), 70% (Y2), 75% (Y3) | Growing as card features mature |
| Avg monthly card spend per card user | EUR 500 (Y1), EUR 650 (Y2), EUR 750 (Y3) | Growing with user trust |
| Interchange rate (blended) | 0.18% | Slightly below 0.20% cap due to merchant mix |
| FX active users (% of active) | 25% (Y1), 30% (Y2), 35% (Y3) | Multi-currency is a key draw |
| Avg FX volume per FX user/month | EUR 400 (Y1), EUR 500 (Y2), EUR 600 (Y3) | |
| Blended FX markup | 0.35% | Weighted across tiers |
| Crypto active users (% of active) | 15% (Y1), 20% (Y2), 25% (Y3) | Crypto adoption grows with features |
| Avg crypto volume per crypto user/month | EUR 200 (Y1), EUR 300 (Y2), EUR 400 (Y3) | |
| Blended crypto fee | 1.15% | Weighted across tiers |
| Avg deposit per active user | EUR 800 (Y1), EUR 1,200 (Y2), EUR 1,500 (Y3) | |

### 1.4 Cost Assumptions

| Parameter | Value | Basis |
|-----------|-------|-------|
| Development cost (M1-9, pre-launch) | EUR 1,700,000 total | 18-person team x 9 months + infra |
| AWS infrastructure base | EUR 32,500/mo (Y1) | See unit economics doc |
| AWS scaling factor | +15% per 2x users beyond 50K | |
| Enfuce per-transaction cost | EUR 0.03 | Conservative estimate |
| Enfuce monthly per-card fee | EUR 0.15 | |
| Banking Circle per SEPA | EUR 0.20 (outgoing), EUR 0.10 (incoming) | |
| Sumsub per verification | EUR 1.50 | Including ongoing monitoring |
| Team growth | +4 heads in Y2, +6 heads in Y3 | Support + engineering scaling |
| Marketing budget (post-launch) | 20% of revenue (Y1), 15% (Y2), 12% (Y3) | Front-loaded for growth |

---

## 2. User Growth Model

### 2.1 Monthly Active Users (MAU) -- Year 1

| Month | New Registrations | Cumulative Registered | Active Users (65%) | Churned (3.5%/mo) | Net Active MAU |
|-------|------------------|----------------------|-------------------|--------------------|---------------|
| M1-7 | 0 (development) | 0 | 0 | 0 | 0 |
| M8 | 300 (beta) | 300 | 195 | 0 | 195 |
| M9 | 200 (beta) | 500 | 325 | 7 | 318 |
| M10 | 2,000 (launch) | 2,500 | 1,625 | 11 | 1,614 |
| M11 | 2,500 | 5,000 | 3,250 | 56 | 3,194 |
| M12 | 3,125 | 8,125 | 5,281 | 112 | 5,170 |
| M13 | 3,900 | 12,025 | 7,816 | 181 | 7,635 |
| M14 | 4,875 | 16,900 | 10,985 | 267 | 10,718 |
| M15 | 6,094 | 22,994 | 14,946 | 375 | 14,571 |
| M16 | 7,617 | 30,611 | 19,897 | 510 | 19,387 |
| M17 | 9,522 | 40,133 | 26,086 | 678 | 25,408 |
| M18 | 11,902 | 52,035 | 33,823 | 889 | 32,934 |
| M19 | 13,688 | 65,723 | 42,720 | 1,153 | 41,567 |
| M20 | 15,741 | 81,464 | 52,952 | 1,455 | 51,497 |
| M21 | 18,102 | 99,566 | 64,718 | 1,802 | 62,916 |

**Year 1 End (Month 21, 12 months post-launch):** ~50,000 MAU (aligned with PRD target G1 of 50,000 MAU).

**Note:** Months 1-9 are development; public launch is Month 10. The "12 months post-launch" PRD target aligns with Month 21 in absolute timeline.

### 2.2 Quarterly Active Users Summary

| Quarter | End MAU | Registered Users | Premium Users | Metal Users |
|---------|---------|-----------------|---------------|-------------|
| Q1 (M1-3) | 0 | 0 | 0 | 0 |
| Q2 (M4-6) | 0 | 0 | 0 | 0 |
| Q3 (M7-9) | 318 | 500 | 57 | 13 |
| Q4 (M10-12) | 5,170 | 8,125 | 931 | 207 |
| Q5 (M13-15) | 14,571 | 22,994 | 2,623 | 583 |
| Q6 (M16-18) | 32,934 | 52,035 | 5,928 | 1,317 |
| Q7 (M19-21) | 51,497 | 81,464 | 11,329 | 3,090 |
| Q8 (M22-24) | 72,000 | 120,000 | 15,840 | 4,320 |
| Q9 (M25-27) | 90,000 | 155,000 | 22,500 | 6,300 |
| Q10 (M28-30) | 110,000 | 195,000 | 27,500 | 7,700 |
| Q11 (M31-33) | 130,000 | 240,000 | 32,500 | 9,100 |
| Q12 (M34-36) | 150,000 | 290,000 | 37,500 | 10,500 |

**Year-End Summary:**

| Metric | Y1 End (M21) | Y2 End (M33) | Y3 End (M36+) |
|--------|-------------|-------------|---------------|
| Total Registered | ~82,000 | ~240,000 | ~290,000+ |
| Monthly Active Users | ~51,000 | ~130,000 | ~150,000 |
| Premium subscribers | ~11,300 | ~32,500 | ~37,500 |
| Metal subscribers | ~3,100 | ~9,100 | ~10,500 |
| Active cards | ~35,700 | ~97,500 | ~112,500 |
| Crypto wallets active | ~7,700 | ~26,000 | ~37,500 |

---

## 3. Revenue Projections

### 3.1 Monthly Revenue -- Year 1 (Post-Launch Months Only)

| Revenue Stream | M10 | M11 | M12 | M13 | M14 | M15 | M16 | M17 | M18 | M19 | M20 | M21 |
|----------------|------|------|------|------|------|------|------|------|------|------|------|------|
| Interchange | 0.9K | 1.7K | 2.8K | 4.1K | 5.8K | 7.9K | 10.5K | 13.7K | 17.8K | 22.5K | 27.8K | 34.0K |
| FX markup | 0.6K | 1.1K | 1.8K | 2.7K | 3.8K | 5.1K | 6.8K | 8.9K | 11.5K | 14.5K | 17.9K | 22.0K |
| Crypto fees | 0.3K | 0.6K | 0.9K | 1.3K | 1.8K | 2.5K | 3.4K | 4.4K | 5.7K | 7.2K | 8.9K | 10.9K |
| Subscriptions | 1.5K | 3.0K | 4.9K | 7.3K | 10.3K | 14.0K | 18.7K | 24.4K | 31.6K | 41.9K | 52.7K | 66.2K |
| ATM fees | 0.1K | 0.2K | 0.4K | 0.5K | 0.7K | 1.0K | 1.3K | 1.7K | 2.2K | 2.8K | 3.5K | 4.3K |
| Card/other fees | 0.2K | 0.3K | 0.4K | 0.5K | 0.7K | 0.9K | 1.2K | 1.5K | 1.8K | 2.3K | 2.7K | 3.2K |
| Interest on funds | 0.4K | 0.8K | 1.3K | 2.0K | 2.8K | 3.9K | 5.2K | 6.8K | 8.8K | 11.1K | 13.7K | 16.7K |
| **Total/Month** | **4.0K** | **7.7K** | **12.5K** | **18.4K** | **25.9K** | **35.3K** | **47.1K** | **61.4K** | **79.4K** | **102.3K** | **127.2K** | **157.3K** |

**Year 1 Total Revenue (M10-M21):** EUR 678,500

### 3.2 Quarterly Revenue -- Years 1-3

| Quarter | Interchange | FX | Crypto | Subscriptions | ATM/Card/Other | Interest | **Total** |
|---------|-------------|-----|--------|---------------|---------------|----------|-----------|
| Q3 (M7-9) | 0.0K | 0.0K | 0.0K | 0.0K | 0.0K | 0.0K | **0.0K** |
| Q4 (M10-12) | 5.4K | 3.5K | 1.8K | 9.4K | 1.6K | 2.5K | **24.2K** |
| Q5 (M13-15) | 17.8K | 11.6K | 5.6K | 31.6K | 5.2K | 8.7K | **80.5K** |
| Q6 (M16-18) | 42.0K | 27.2K | 13.5K | 74.7K | 12.2K | 20.8K | **190.4K** |
| Q7 (M19-21) | 84.3K | 54.4K | 27.0K | 160.8K | 24.6K | 41.5K | **392.6K** |
| **Y1 Total** | | | | | | | **687.7K** |
| Q8 (M22-24) | 120K | 84K | 45K | 260K | 38K | 60K | **607K** |
| Q9 (M25-27) | 162K | 119K | 68K | 380K | 52K | 82K | **863K** |
| Q10 (M28-30) | 198K | 149K | 89K | 470K | 65K | 102K | **1,073K** |
| Q11 (M31-33) | 234K | 177K | 110K | 560K | 78K | 122K | **1,281K** |
| **Y2 Total** | | | | | | | **3,824K** |
| Q12 (M34-36) | 270K | 210K | 135K | 650K | 90K | 140K | **1,495K** |
| Q13-Q14 (M37-42) | 580K | 450K | 300K | 1,400K | 195K | 290K | **3,215K** |
| Q15 (M43-45) est. | 300K | 240K | 165K | 740K | 100K | 155K | **1,700K** |
| **Y3 Total** | | | | | | | **6,410K** |

### 3.3 Annual Revenue Summary

| Year | Revenue | YoY Growth | PRD Target | vs Target |
|------|---------|-----------|------------|-----------|
| Year 1 (M1-M21) | EUR 688K | -- | EUR 8M (aspirational) | Below (pre-scale) |
| Year 2 (M22-M33) | EUR 3,824K | +456% | -- | Strong growth trajectory |
| Year 3 (M34-M45) | EUR 6,410K | +68% | -- | Approaching profitability |
| **3-Year Cumulative** | **EUR 10,922K** | | | |

**Note on PRD Revenue Target:** The PRD target of EUR 8M annual revenue 12 months post-launch is extremely aggressive. Our base-case projection shows EUR 688K in the first 12 months of operation (including 9 months of development with zero revenue). The EUR 8M target would require approximately 180,000 MAU with high ARPU -- achievable in Year 3 but not Year 1. We recommend revising the PRD target to EUR 2-3M for the first 12 months of commercial operation (M10-M21) in the optimistic scenario.

---

## 4. Cost Projections

### 4.1 Monthly Costs -- Year 1

| Cost Category | M1-9 (Dev/Mo) | M10 | M12 | M15 | M18 | M21 |
|---------------|---------------|------|------|------|------|------|
| Team salaries | 135.0K | 135.0K | 135.0K | 142.0K | 142.0K | 150.0K |
| AWS infrastructure | 25.0K | 32.5K | 32.5K | 33.5K | 35.0K | 37.5K |
| Third-party SaaS | 3.0K | 4.0K | 4.5K | 5.0K | 5.5K | 6.0K |
| Enfuce (card processing) | 0.0K | 0.5K | 1.5K | 4.2K | 9.5K | 16.5K |
| Banking Circle (SEPA) | 0.0K | 0.3K | 0.8K | 2.3K | 5.2K | 9.0K |
| Sumsub (KYC) | 1.5K | 3.0K | 4.7K | 9.1K | 17.9K | 20.6K |
| Fuse.io / blockchain | 0.5K | 0.5K | 0.8K | 1.2K | 2.0K | 3.0K |
| Marketing | 5.0K | 15.0K | 20.0K | 30.0K | 40.0K | 35.0K |
| Office / admin | 10.0K | 10.0K | 10.0K | 11.0K | 11.0K | 12.0K |
| Legal / compliance | 8.0K | 6.0K | 6.0K | 6.0K | 6.0K | 7.0K |
| Fraud / chargebacks | 0.0K | 0.1K | 0.3K | 0.8K | 1.8K | 3.2K |
| Contingency (10%) | 18.8K | 20.7K | 21.6K | 24.5K | 27.6K | 30.0K |
| **Total/Month** | **206.8K** | **227.6K** | **237.7K** | **269.6K** | **303.5K** | **329.8K** |

### 4.2 Quarterly Cost Summary -- Years 1-3

| Quarter | Team | Infrastructure | Variable (Enfuce/BC/Sumsub) | Marketing | Other | **Total** |
|---------|------|---------------|----------------------------|-----------|-------|-----------|
| Q1 (M1-3) | 405K | 84K | 5K | 15K | 98K | **607K** |
| Q2 (M4-6) | 405K | 84K | 5K | 15K | 98K | **607K** |
| Q3 (M7-9) | 405K | 84K | 8K | 15K | 98K | **610K** |
| Q4 (M10-12) | 405K | 99K | 15K | 55K | 96K | **670K** |
| Q5 (M13-15) | 420K | 101K | 35K | 80K | 100K | **736K** |
| Q6 (M16-18) | 426K | 106K | 75K | 110K | 105K | **822K** |
| Q7 (M19-21) | 445K | 112K | 120K | 100K | 110K | **887K** |
| **Y1 Total** | | | | | | **4,939K** |
| Q8 (M22-24) | 480K | 125K | 160K | 91K | 118K | **974K** |
| Q9 (M25-27) | 495K | 135K | 195K | 130K | 125K | **1,080K** |
| Q10 (M28-30) | 510K | 145K | 230K | 161K | 132K | **1,178K** |
| Q11 (M31-33) | 530K | 155K | 270K | 192K | 140K | **1,287K** |
| **Y2 Total** | | | | | | **4,519K** |
| Q12 (M34-36) | 550K | 162K | 300K | 180K | 148K | **1,340K** |
| Q13-Q14 (M37-42) | 1,140K | 340K | 640K | 386K | 310K | **2,816K** |
| Q15 (M43-45) est. | 580K | 175K | 340K | 204K | 160K | **1,459K** |
| **Y3 Total** | | | | | | **5,615K** |

### 4.3 Annual Cost Summary

| Category | Year 1 | Year 2 | Year 3 |
|----------|--------|--------|--------|
| Team (salaries + benefits) | 2,911K | 2,015K | 2,270K |
| Infrastructure (AWS + SaaS) | 570K | 560K | 677K |
| Variable processing costs | 263K | 855K | 1,280K |
| Marketing | 390K | 574K | 770K |
| Office / admin / legal | 597K | 515K | 618K |
| Contingency | 208K | -- | -- |
| **Total** | **4,939K** | **4,519K** | **5,615K** |

---

## 5. Profit and Loss Summary

### 5.1 Monthly P&L -- Year 1 (Key Months)

| Line Item | M1-9 (Avg/Mo) | M10 | M12 | M15 | M18 | M21 |
|-----------|---------------|------|------|------|------|------|
| **Revenue** | 0K | 4.0K | 12.5K | 35.3K | 79.4K | 157.3K |
| **COGS (variable)** | 2K | 4.4K | 8.1K | 17.6K | 36.4K | 52.3K |
| **Gross Profit** | (2K) | (0.4K) | 4.4K | 17.7K | 43.0K | 105.0K |
| Gross Margin | -- | (10%) | 35% | 50% | 54% | 67% |
| Team costs | 135.0K | 135.0K | 135.0K | 142.0K | 142.0K | 150.0K |
| Infrastructure (fixed) | 28.0K | 36.5K | 37.0K | 38.5K | 40.5K | 43.5K |
| Marketing | 5.0K | 15.0K | 20.0K | 30.0K | 40.0K | 35.0K |
| Other opex | 18.0K | 16.0K | 16.0K | 17.0K | 17.0K | 19.0K |
| Contingency | 18.8K | 20.7K | 21.6K | 24.5K | 27.6K | 30.0K |
| **Total OpEx** | 204.8K | 223.2K | 229.6K | 252.0K | 267.1K | 277.5K |
| **EBITDA** | **(206.8K)** | **(223.6K)** | **(225.2K)** | **(234.3K)** | **(224.1K)** | **(172.5K)** |
| Cumulative EBITDA | (206.8K) | (2,088K) | (2,536K) | (3,240K) | (3,956K) | (4,537K) |

### 5.2 Quarterly P&L -- Years 2-3

| Line Item | Q8 | Q9 | Q10 | Q11 | Q12 | Q13-14 | Q15 |
|-----------|-----|-----|------|------|------|--------|------|
| **Revenue** | 607K | 863K | 1,073K | 1,281K | 1,495K | 3,215K | 1,700K |
| **COGS** | 175K | 215K | 255K | 295K | 330K | 700K | 375K |
| **Gross Profit** | 432K | 648K | 818K | 986K | 1,165K | 2,515K | 1,325K |
| Gross Margin | 71% | 75% | 76% | 77% | 78% | 78% | 78% |
| Operating Expenses | 799K | 865K | 923K | 992K | 1,010K | 2,116K | 1,084K |
| **EBITDA** | **(367K)** | **(217K)** | **(105K)** | **(6K)** | **155K** | **399K** | **241K** |
| Cumulative EBITDA | (4,904K) | (5,121K) | (5,226K) | (5,232K) | (5,077K) | (4,678K) | (4,437K) |

### 5.3 Annual P&L Summary

| Line Item | Year 1 | Year 2 | Year 3 |
|-----------|--------|--------|--------|
| **Revenue** | 688K | 3,824K | 6,410K |
| **COGS (variable processing)** | 263K | 940K | 1,405K |
| **Gross Profit** | 425K | 2,884K | 5,005K |
| **Gross Margin** | 62% | 75% | 78% |
| Operating Expenses | 4,676K | 3,579K | 4,210K |
| **EBITDA** | **(4,251K)** | **(695K)** | **795K** |
| **EBITDA Margin** | -- | (18%) | 12% |
| Depreciation / Amortization | 50K | 80K | 100K |
| **EBIT** | **(4,301K)** | **(775K)** | **695K** |
| Interest expense (debt) | 0K | 0K | 0K |
| **EBT** | **(4,301K)** | **(775K)** | **695K** |
| Tax (15%, only if profitable) | 0K | 0K | 104K |
| **Net Income** | **(4,301K)** | **(775K)** | **591K** |

---

## 6. Cash Flow Analysis

### 6.1 Operating Cash Flow

| Period | EBITDA | Working Capital Changes | Operating Cash Flow |
|--------|--------|----------------------|-------------------|
| Year 1 | (4,251K) | (200K) | (4,451K) |
| Year 2 | (695K) | (150K) | (845K) |
| Year 3 | 795K | (100K) | 695K |

### 6.2 Capital Expenditure

| Period | Amount | Purpose |
|--------|--------|---------|
| Year 1 | 150K | Development tooling, security hardware, Mastercard certification fees |
| Year 2 | 100K | Additional infrastructure, compliance certifications |
| Year 3 | 120K | Infrastructure expansion, new market licenses |

### 6.3 Cash Flow Summary

| Period | Operating CF | CapEx | Free Cash Flow | Cumulative FCF |
|--------|-------------|-------|---------------|---------------|
| Year 1 | (4,451K) | (150K) | (4,601K) | (4,601K) |
| Year 2 | (845K) | (100K) | (945K) | (5,546K) |
| Year 3 | 695K | (120K) | 575K | (4,971K) |

**Total cash requirement through profitability: EUR 5.5M**

---

## 7. Funding Requirements

### 7.1 Funding Need Assessment

| Milestone | Timing | Cumulative Cash Need | Buffer (20%) | **Total Funding Need** |
|-----------|--------|---------------------|-------------|----------------------|
| End of MVP development | Month 9 | EUR 1,860K | EUR 370K | EUR 2,230K |
| End of Year 1 | Month 21 | EUR 4,600K | EUR 920K | EUR 5,520K |
| EBITDA breakeven | Month 33 | EUR 5,230K | EUR 1,050K | EUR 6,280K |
| Cash flow positive | Month 38 | EUR 5,550K | EUR 1,110K | EUR 6,660K |

### 7.2 Recommended Funding Strategy

| Round | Timing | Amount | Purpose | Runway Created |
|-------|--------|--------|---------|---------------|
| **Pre-Seed / Existing Capital** | M1 | EUR 2,500K | MVP development + 3 months post-launch buffer | Through Month 12 |
| **Seed Round** | M12-15 | EUR 4,000K | Scale to 50K users, marketing, team expansion | Through Month 30 |
| **Series A (if needed)** | M24-27 | EUR 5,000-8,000K | Accelerate growth to 150K+ users, Phase 2 features, new markets | Through profitability |

**Note:** TeslaPay's existing revenue is EUR 1.69M (2024) with EUR 9.85M in safeguarded funds. Assuming existing operations continue generating EUR 150K/month in revenue, the actual external funding need is reduced. However, the new platform investment is incremental, so we model it as a separate capital requirement.

### 7.3 Use of Funds (Seed Round, EUR 4,000K)

| Category | Amount | % |
|----------|--------|---|
| Engineering team (6 months) | 1,200K | 30% |
| Infrastructure and third-party costs | 600K | 15% |
| Marketing and user acquisition | 800K | 20% |
| Compliance and legal | 400K | 10% |
| Working capital / operations | 600K | 15% |
| Reserve / contingency | 400K | 10% |
| **Total** | **4,000K** | **100%** |

---

## 8. Key Metrics Timeline

| Metric | M12 | M18 | M24 (Y2) | M30 | M36 (Y3) |
|--------|------|------|----------|------|----------|
| MAU | 5,170 | 32,934 | 72,000 | 110,000 | 150,000 |
| Monthly Revenue | 12.5K | 79.4K | 202K | 358K | 500K+ |
| Annual Revenue Run Rate | 150K | 953K | 2,424K | 4,296K | 6,000K+ |
| Monthly Burn Rate | 225K | 224K | 325K | 393K | 423K |
| Monthly Net Burn | 213K | 145K | 123K | 35K | (77K) positive |
| Cash Remaining (from 6.5M raised) | 4,280K | 2,410K | 1,350K | 630K | 955K |
| Gross Margin | 35% | 54% | 71% | 76% | 78% |
| EBITDA Margin | -- | -- | (16%) | 2% | 12% |
| LTV:CAC (cohort) | 2.0:1 | 2.8:1 | 3.2:1 | 3.8:1 | 4.2:1 |
| Blended CAC | 28 | 25 | 23 | 21 | 19 |

---

## 9. Sensitivity Analysis

### 9.1 Three Scenarios

| Parameter | Pessimistic | Base | Optimistic |
|-----------|------------|------|-----------|
| MAU at M21 (12 mo post-launch) | 25,000 | 51,000 | 80,000 |
| MAU at M36 | 80,000 | 150,000 | 250,000 |
| Blended monthly ARPU | EUR 3.50 | EUR 4.51 | EUR 5.50 |
| Premium conversion | 12% | 20% | 28% |
| Monthly churn (Y1) | 5.0% | 3.5% | 2.5% |
| CAC (blended) | EUR 35 | EUR 24 | EUR 18 |
| Fixed cost growth | +25% | Base | -10% |

### 9.2 Scenario Outcomes -- Year 3

| Metric | Pessimistic | Base | Optimistic |
|--------|------------|------|-----------|
| Y3 Revenue | EUR 3,200K | EUR 6,410K | EUR 10,500K |
| Y3 EBITDA | (EUR 1,200K) | EUR 795K | EUR 3,150K |
| EBITDA Breakeven | Month 42+ | Month 33 | Month 24 |
| Total funding required | EUR 9,000K | EUR 6,500K | EUR 4,500K |
| LTV:CAC at scale | 1.8:1 | 3.5:1 | 5.5:1 |
| Users at breakeven | 95,000 | 59,000 | 38,000 |

### 9.3 Key Risk Scenarios

**Scenario A: Crypto Winter**
- Crypto trading revenue drops 60%
- Crypto user adoption drops from 15% to 8%
- Impact: Revenue reduced by ~12% (EUR 770K over 3 years)
- Mitigation: Banking core sustains business; crypto is <15% of revenue

**Scenario B: Low Subscription Conversion**
- Premium conversion stays at 12% (vs. 20% target)
- Metal conversion stays at 2% (vs. 5%)
- Impact: Revenue reduced by ~22% (EUR 2.4M over 3 years)
- Mitigation: This is the highest-risk scenario. Must invest in product differentiation for paid tiers. Consider price reduction to EUR 5.99/EUR 9.99 to boost conversion.

**Scenario C: Regulatory Cost Increase**
- MiCA compliance requires additional EUR 200K/year
- PSD3 transition costs EUR 150K one-time
- Impact: Breakeven delayed 3-4 months
- Mitigation: Budget contingency; phased compliance investment

**Scenario D: Competitor Response**
- Revolut launches DeFi features, eroding TeslaPay's differentiation
- Impact: CAC increases 40% (from EUR 24 to EUR 34); growth slows
- Mitigation: Deepen Fuse.io integration; focus on underserved markets (Baltics, CEE)

---

## 10. Breakeven Analysis

### 10.1 EBITDA Breakeven

| Scenario | Month | MAU at Breakeven | Monthly Revenue | Monthly Costs |
|----------|-------|-----------------|-----------------|---------------|
| Optimistic | Month 24 | 72,000 | EUR 325K | EUR 315K |
| Base | Month 33 | 130,000 | EUR 430K | EUR 430K |
| Pessimistic | Month 42+ | 95,000+ | EUR 320K+ | EUR 350K+ |

### 10.2 Cash Flow Breakeven

Cash flow breakeven (cumulative positive FCF) is projected at Month 38-42 in the base case, requiring a total of approximately EUR 5.5-6.5M in funding.

### 10.3 Path to EUR 8M Revenue (PRD Target)

The PRD target of EUR 8M annual revenue requires:
- ~180,000 MAU with current ARPU
- OR ~130,000 MAU with ARPU increase to EUR 5.50/mo (driven by business accounts in Phase 2)
- **Achievable timeline:** Month 30-36 (Year 3) in the base case

---

## 11. Safeguarded Funds Projection

TeslaPay is required under EMD2 to safeguard customer funds. The interest earned on these pooled funds is a meaningful revenue stream.

| Metric | M12 | M18 | M24 | M30 | M36 |
|--------|------|------|------|------|------|
| MAU | 5,170 | 32,934 | 72,000 | 110,000 | 150,000 |
| Avg deposit/user | EUR 600 | EUR 800 | EUR 1,000 | EUR 1,200 | EUR 1,500 |
| Total safeguarded funds | EUR 3.1M | EUR 26.3M | EUR 72M | EUR 132M | EUR 225M |
| ECB rate assumption | 3.25% | 3.25% | 2.75% | 2.75% | 2.50% |
| Annual interest income | EUR 101K | EUR 855K | EUR 1,980K | EUR 3,630K | EUR 5,625K |
| Monthly interest income | EUR 8K | EUR 71K | EUR 165K | EUR 303K | EUR 469K |

**Note:** Interest on safeguarded funds could become TeslaPay's largest revenue stream by Year 3 if deposit growth materializes. This is highly dependent on user trust and the ECB rate environment. If ECB rates fall to 1.5%, this revenue halves.

---

## 12. Valuation Indicators

For investor discussions, comparable multiples suggest:

| Metric | TeslaPay Y3 | Multiple Range | Implied Valuation |
|--------|------------|---------------|-------------------|
| Revenue (EUR 6.4M) | 6.4M | 8-15x (neobank revenue) | EUR 51-96M |
| MAU (150K) | 150K | EUR 200-500/user | EUR 30-75M |
| Safeguarded funds (EUR 225M) | 225M | 0.5-1.5% of AUM | EUR 1.1-3.4M (low) |

These are indicative ranges only. Actual valuation depends on growth rate, unit economics proof, and market conditions at the time of fundraising.

---

## 13. Recommendations

1. **Secure EUR 2.5M pre-seed/seed immediately** to fund the 9-month MVP development phase with 3 months post-launch buffer.

2. **Plan seed round (EUR 4M) at Month 12-15** when early traction data (5K+ MAU, subscription conversion rates, retention) is available to support valuation.

3. **Revise PRD revenue target.** The EUR 8M target for 12 months post-launch is not achievable under any realistic scenario. Recommend EUR 2-3M for the first 12 months of commercial operation, with EUR 8M as a Year 3 target.

4. **Monitor subscription conversion weekly.** This is the single largest revenue lever and risk factor. If conversion is below 15% by Month 15, consider pricing or feature adjustments.

5. **Maintain EUR 500K minimum cash reserve** at all times as an operating buffer for regulatory requirements and unexpected costs.

6. **Explore revenue from existing operations.** TeslaPay already generates EUR 1.69M/year. Transitioning existing customers to the new platform can provide baseline revenue during the growth phase.

7. **Consider strategic partnership with Fuse.io.** A co-marketing arrangement or grant from the Fuse ecosystem fund could reduce marketing costs by EUR 100-200K in Year 1.

---

**Document Approval:**

| Role | Name | Date | Status |
|------|------|------|--------|
| CEO | TBD | | Pending |
| CFO | Dream Team Financial Analyst | 2026-03-03 | Submitted |
| Board | TBD | | Pending |
| Investor Relations | TBD | | Pending |
