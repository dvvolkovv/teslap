# TeslaPay Information Architecture

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Senior UI/UX Designer, Dream Team

---

## 1. Bottom Tab Navigation (5 Tabs)

```
+----------+----------+----------+----------+----------+
|   Home   | Payments |   Card   |  Crypto  | Profile  |
|   (O)    |   (->)   |   [=]    |   {B}    |   (@)    |
+----------+----------+----------+----------+----------+
```

| Tab      | Icon (Phosphor)     | Label    | Badge Behavior                      |
|----------|---------------------|----------|-------------------------------------|
| Home     | House               | Home     | Dot badge on new notifications      |
| Payments | ArrowsLeftRight     | Payments | None                                |
| Card     | CreditCard          | Card     | Dot badge on pending activation     |
| Crypto   | CurrencyBtc         | Crypto   | Dot badge on received crypto        |
| Profile  | UserCircle          | Profile  | Dot badge on KYC action required    |

---

## 2. Complete Screen Hierarchy

### 2.1 Pre-Authentication (Unauthenticated)

```
App Launch
  |
  +-- Splash Screen
  |
  +-- Welcome / Onboarding Carousel (3 slides)
  |     Slide 1: "Banking without borders" (multi-currency)
  |     Slide 2: "Your card, your rules" (Mastercard)
  |     Slide 3: "Crypto made simple" (Fuse.io)
  |
  +-- Login Screen
  |     +-- Biometric Prompt (Face ID / Fingerprint)
  |     +-- PIN Entry (fallback)
  |     +-- Password Entry (fallback)
  |     +-- Forgot Password Flow
  |           +-- Enter Email
  |           +-- OTP Verification
  |           +-- Set New Password
  |           +-- Confirmation
  |
  +-- Registration Flow
        +-- Step 1: Email + Phone
        +-- Step 2: Email OTP Verification
        +-- Step 3: Phone SMS Verification
        +-- Step 4: Create Password
        +-- Step 5: Set App PIN (6-digit)
        +-- Step 6: Enable Biometric Auth
        +-- Step 7: KYC Start (Sumsub)
        |     +-- Select Document Type
        |     +-- Document Capture (Front)
        |     +-- Document Capture (Back, if applicable)
        |     +-- NFC Passport Read (optional)
        |     +-- Liveness Check
        |     +-- Processing Screen
        |     +-- KYC Result (Approved / Pending Review / Rejected)
        +-- Step 8: Account Created Confirmation
        +-- Step 9: Tier Assignment Display
        +-- Step 10: First Action Prompt (Add Funds / Get Card)
```

### 2.2 Tab 1: Home

```
Home (Dashboard)
  |
  +-- [Hero] Balance Card
  |     +-- Tap: Currency selector (switch between total/individual)
  |
  +-- [Section] Quick Actions Row
  |     +-- Send Money --> Payments/Send flow
  |     +-- Request Money --> Payments/Request flow
  |     +-- Exchange --> Payments/Exchange flow
  |     +-- Top Up --> Top Up bottom sheet
  |
  +-- [Section] Accounts Overview
  |     +-- Currency Row (EUR) --> Account Detail
  |     +-- Currency Row (USD) --> Account Detail
  |     +-- Currency Row (other) --> Account Detail
  |     +-- "+ Add Currency" --> Add Currency bottom sheet
  |
  +-- [Section] Recent Transactions (last 5)
  |     +-- Transaction Item --> Transaction Detail
  |     +-- "See All" --> Full Transaction History
  |
  +-- [Section] Crypto Portfolio Summary (if activated)
  |     +-- Mini balance card --> Crypto tab
  |
  +-- [FAB] Support Chat Button
  |
  +-- Notifications Bell (top right)
        +-- Notification Center
              +-- Notification Item --> Deep link target
              +-- Mark all as read
              +-- Notification Preferences link
```

**Account Detail Screen:**
```
Account Detail (e.g., EUR Account)
  |
  +-- Balance display
  +-- IBAN display (copy, share)
  +-- Account actions: Send, Receive, Exchange
  +-- Transaction list (filtered to this currency)
  |     +-- Filter bar (date, type, amount)
  |     +-- Transaction Item --> Transaction Detail
  +-- Export (CSV / PDF)
  +-- Account settings
        +-- Rename account
        +-- Close account (if balance zero)
```

**Transaction Detail Screen:**
```
Transaction Detail
  |
  +-- Amount and currency
  +-- Status badge (Completed / Pending / Failed / Declined)
  +-- Counterparty (name, IBAN)
  +-- Reference / Description
  +-- Category (auto-assigned, editable)
  +-- Date and time
  +-- FX rate (if applicable)
  +-- Fee breakdown
  +-- [Actions]
       +-- Repeat this payment
       +-- Share receipt
       +-- Dispute transaction (card transactions only)
       +-- Report a problem
```

### 2.3 Tab 2: Payments

```
Payments
  |
  +-- [Section] Quick Actions
  |     +-- Send Money
  |     +-- Request Money
  |     +-- Exchange
  |     +-- Scan QR
  |
  +-- [Section] Saved Payees
  |     +-- Payee Item --> Send pre-filled
  |     +-- "+ Add Payee" --> Add Payee flow
  |     +-- "Manage" --> Payee Management
  |
  +-- [Section] Scheduled Payments
  |     +-- Upcoming payment item --> Payment detail
  |     +-- "+ Schedule" --> Create recurring
  |     +-- "See All" --> Full scheduled list
  |
  +-- [Section] Direct Debits (SDD)
        +-- Active mandate item --> Mandate detail
        +-- "Manage" --> Mandate management
```

**Send Money Flow:**
```
Send Money
  |
  +-- Step 1: Select Recipient
  |     +-- Search saved payees
  |     +-- Enter IBAN manually
  |     +-- Enter phone/username (TeslaPay internal)
  |     +-- Scan QR code
  |     +-- Recent recipients list
  |
  +-- Step 2: Enter Amount
  |     +-- Amount input (numpad)
  |     +-- Currency selector (source account)
  |     +-- Add reference / note
  |     +-- Toggle: Regular SEPA / SEPA Instant
  |     +-- Fee and delivery time display
  |
  +-- Step 3: Review and Confirm
  |     +-- Full summary (from, to, amount, fee, arrival)
  |     +-- Biometric / PIN confirmation
  |
  +-- Step 4: Success / Failure Screen
        +-- Share receipt
        +-- Send another
        +-- Back to Home
```

**Exchange Flow:**
```
Exchange Currency
  |
  +-- Source currency selector + amount
  +-- Swap button (reverse direction)
  +-- Target currency selector + converted amount
  +-- Live rate display (30s lock)
  +-- Fee/markup disclosure
  +-- Confirm with biometric/PIN
  +-- Success screen
```

### 2.4 Tab 3: Card

```
Card
  |
  +-- [Hero] Card Visual (swipeable if multiple cards)
  |     +-- Virtual card display
  |     +-- Physical card display
  |     +-- "+ Order Card" (if no card yet)
  |
  +-- [Quick Actions]
  |     +-- Show Details (biometric required)
  |     +-- Freeze / Unfreeze toggle
  |     +-- Add to Apple Pay / Google Pay
  |
  +-- [Section] Card Settings
  |     +-- View PIN (biometric required)
  |     +-- Change PIN
  |     +-- Spending Limits
  |     |     +-- Per-transaction limit slider
  |     |     +-- Daily limit slider
  |     |     +-- Monthly limit slider
  |     |     +-- ATM withdrawal limit
  |     +-- Security Controls
  |     |     +-- Toggle: Online payments
  |     |     +-- Toggle: Contactless
  |     |     +-- Toggle: ATM withdrawals
  |     |     +-- Toggle: Magnetic stripe
  |     |     +-- Merchant category blocks
  |     |     +-- Geographic restrictions
  |     +-- Linked Account (change source account)
  |
  +-- [Section] Recent Card Transactions
  |     +-- Transaction Item --> Transaction Detail
  |     +-- "See All" --> Transaction History (card filter)
  |
  +-- [Section] Card Actions
        +-- Order Physical Card --> Order flow
        +-- Report Lost / Stolen --> Report flow
        +-- Replace Card --> Replacement flow
```

**Order Physical Card Flow:**
```
Order Physical Card
  |
  +-- Step 1: Confirm Delivery Address
  |     +-- Pre-filled from profile
  |     +-- Edit address option
  |
  +-- Step 2: Card Options
  |     +-- Standard (free) or Premium design (fee)
  |
  +-- Step 3: Fee Confirmation (if applicable)
  |
  +-- Step 4: Order Summary + Confirm
  |
  +-- Step 5: Order Confirmed
        +-- Estimated delivery date
        +-- Track order link
```

**Activate Physical Card Flow:**
```
Activate Card
  |
  +-- Option A: Scan card (NFC or camera)
  +-- Option B: Enter last 4 digits
  +-- Activation processing
  +-- Success screen
```

### 2.5 Tab 4: Crypto

```
Crypto
  |
  +-- [Hero] Total Crypto Balance (EUR equivalent)
  |     +-- 24h change indicator
  |
  +-- [Section] Token Balances
  |     +-- FUSE card --> Token Detail
  |     +-- USDC card --> Token Detail
  |     +-- USDT card --> Token Detail
  |
  +-- [Quick Actions]
  |     +-- Buy Crypto
  |     +-- Sell Crypto
  |     +-- Send
  |     +-- Receive
  |
  +-- [Section] DeFi Yield (Phase 2)
  |     +-- Solid soUSD card
  |     +-- Current APY
  |     +-- "Earn" CTA
  |
  +-- [Section] Recent Crypto Activity
  |     +-- Crypto transaction item --> Crypto Tx Detail
  |     +-- "See All" --> Crypto Transaction History
  |
  +-- [Section] Wallet Info
        +-- Wallet address (copy, QR)
        +-- Fuse network status indicator
        +-- Learn about Fuse (educational link)
```

**Token Detail Screen:**
```
Token Detail (e.g., FUSE)
  |
  +-- Token balance (token + EUR equivalent)
  +-- Price chart (24h / 7d / 30d / 1y)
  +-- 24h price change
  +-- Actions: Buy, Sell, Send, Receive
  +-- Transaction history for this token
```

**Buy Crypto Flow:**
```
Buy Crypto
  |
  +-- Step 1: Select Token (FUSE / USDC / USDT)
  +-- Step 2: Enter Amount (EUR or token amount)
  |     +-- Live rate display
  |     +-- Fee disclosure
  |     +-- Minimum: EUR 5
  +-- Step 3: Review Order
  |     +-- Rate, fee, you receive amount
  |     +-- Risk disclosure (first-time only)
  +-- Step 4: Confirm (biometric / PIN)
  +-- Step 5: Success (token balance updated)
```

**Send Crypto Flow:**
```
Send Crypto
  |
  +-- Step 1: Select Token
  +-- Step 2: Enter Recipient Address
  |     +-- Paste address
  |     +-- Scan QR code
  |     +-- Address validation
  +-- Step 3: Enter Amount
  |     +-- Network fee estimate
  +-- Step 4: Review and Confirm (biometric / PIN)
  +-- Step 5: Transaction Submitted
        +-- Tx hash
        +-- Link to Fuse explorer
```

**Receive Crypto Screen:**
```
Receive Crypto
  |
  +-- Select token
  +-- QR code (wallet address)
  +-- Address text (copy button)
  +-- Share address
  +-- "Only send [TOKEN] on the Fuse network" warning
```

**DeFi Yield Flow (Phase 2):**
```
Earn Yield
  |
  +-- Solid soUSD overview
  |     +-- Current APY
  |     +-- Total deposited
  |     +-- Earned yield
  +-- Deposit stablecoins
  |     +-- Select token (USDC / USDT)
  |     +-- Enter amount
  |     +-- Risk disclosure
  |     +-- Confirm
  +-- Withdraw
  |     +-- Enter amount or "Withdraw all"
  |     +-- Confirm
  +-- Yield history
```

### 2.6 Tab 5: Profile

```
Profile
  |
  +-- [Header] User name, email, tier badge
  |
  +-- [Section] Account
  |     +-- Personal Information --> Edit profile
  |     +-- Verification Status --> KYC detail
  |     +-- Account Tier --> Tier info + upgrade path
  |     +-- Fees and Limits --> Fee schedule
  |
  +-- [Section] Security
  |     +-- Biometric Login toggle
  |     +-- Change PIN
  |     +-- Change Password
  |     +-- Two-Factor Authentication
  |     +-- Active Sessions --> Session management
  |
  +-- [Section] Preferences
  |     +-- Notifications --> Notification settings
  |     +-- Language --> Language picker
  |     +-- Appearance (Light / Dark / System)
  |     +-- Currency Display (primary fiat reference)
  |
  +-- [Section] Support
  |     +-- Help Center --> FAQ / Articles
  |     +-- Chat with Us --> Support chat
  |     +-- Report a Problem
  |
  +-- [Section] Legal
  |     +-- Terms of Service
  |     +-- Privacy Policy
  |     +-- Licenses
  |     +-- GDPR Data Requests
  |
  +-- [Section] Account Actions
  |     +-- Export Data
  |     +-- Close Account --> Account closure flow
  |
  +-- Log Out
  +-- App Version
```

---

## 3. Navigation Patterns

### 3.1 Primary Navigation

| Pattern         | Usage                                        |
|----------------|----------------------------------------------|
| Bottom tabs     | Top-level sections (5 tabs)                  |
| Stack push      | Drilling into detail (e.g., transaction detail) |
| Bottom sheet    | Quick actions, confirmations, selectors      |
| Full-screen modal | Multi-step flows (send money, buy crypto, KYC) |

### 3.2 Navigation Rules

1. **Tab memory:** Each tab remembers its navigation state. Switching between tabs does not reset the stack.
2. **Tab re-tap:** Tapping the active tab scrolls to top and pops to root of that tab.
3. **Back button (Android):** Pops current screen. At tab root, switches to Home tab. At Home tab root, shows exit confirmation.
4. **iOS swipe-back gesture:** Supported on all pushed screens.
5. **Multi-step flows:** Use full-screen modal with progress indicator. Back navigates to previous step. Close (X) shows "Are you sure?" confirmation if data was entered.
6. **Bottom sheets:** Used for quick selections (currency picker, document type), confirmations, and filters. Never nested more than 1 level deep.

### 3.3 Auth-Gated Screens

These screens require step-up authentication (biometric or PIN) before display:

- Card details (full number, CVV)
- PIN view
- Send money confirmation
- Crypto send confirmation
- Exchange confirmation
- Security settings changes
- Session management
- Account closure

---

## 4. Deep Linking Structure

### 4.1 URI Scheme

Base: `teslapay://` (custom scheme) and `https://app.teslapay.eu/` (universal links)

### 4.2 Deep Link Routes

| Route                                  | Target Screen                  | Source                        |
|----------------------------------------|--------------------------------|-------------------------------|
| `teslapay://home`                      | Home dashboard                 | General                       |
| `teslapay://transactions/{id}`         | Transaction detail             | Push notification             |
| `teslapay://cards`                     | Card tab                       | Card notification             |
| `teslapay://cards/activate`            | Card activation flow           | Delivery notification         |
| `teslapay://cards/3ds/{challengeId}`   | 3DS challenge screen           | 3DS push notification         |
| `teslapay://payments/send`             | Send money flow                | Share/widget                  |
| `teslapay://payments/send?iban={iban}` | Send money pre-filled          | QR code scan                  |
| `teslapay://payments/request`          | Request money                  | Share link                    |
| `teslapay://crypto`                    | Crypto tab                     | Crypto notification           |
| `teslapay://crypto/receive`            | Receive crypto screen          | Share                         |
| `teslapay://crypto/tx/{hash}`          | Crypto transaction detail      | Push notification             |
| `teslapay://profile/kyc`              | KYC verification screen        | Re-verification notification  |
| `teslapay://profile/security`         | Security settings              | Security alert                |
| `teslapay://support`                   | Support chat                   | Any help link                 |
| `teslapay://exchange?from=EUR&to=USD`  | Exchange pre-filled            | Widget / shortcut             |

### 4.3 Deep Link Behavior

1. If user is not authenticated: show biometric/PIN login, then navigate to target
2. If target requires step-up auth: show auth prompt at target screen
3. If target does not exist (e.g., deleted transaction): show error with "Go to Home" fallback
4. Universal links verified via `.well-known/apple-app-site-association` and `assetlinks.json`

---

## 5. State Handling

### 5.1 App States

| State                | Behavior                                           |
|----------------------|---------------------------------------------------|
| Fresh install        | Show onboarding carousel, then registration/login  |
| Logged out           | Show login screen (biometric first, then PIN)      |
| Backgrounded < 5min  | Resume without auth                                |
| Backgrounded > 5min  | Require biometric/PIN on resume                    |
| No internet          | Show cached data with "Offline" banner. Disable sends/buys |
| KYC pending          | Limited dashboard: show KYC status, block financial features |
| KYC rejected         | Show rejection screen with retry option or support link |
| Account frozen       | Show frozen status with support contact            |

### 5.2 Onboarding State Machine

```
[Not Registered] --> [Email Verified] --> [Phone Verified] --> [Password Set]
--> [PIN Set] --> [Biometric Enabled] --> [KYC Submitted] --> [KYC Pending]
--> [KYC Approved] --> [Account Active] --> [First Deposit Prompted]
```

Users can resume onboarding from where they left off if they close the app mid-flow.

---

## 6. Notification Architecture

### 6.1 Push Notification Types

| Type              | Category     | Taps To                        | Priority |
|-------------------|--------------|-------------------------------|----------|
| Card transaction  | Transaction  | Transaction detail             | High     |
| SEPA received     | Transaction  | Transaction detail             | High     |
| Crypto received   | Transaction  | Crypto tx detail               | High     |
| Card declined     | Transaction  | Transaction detail             | High     |
| 3DS challenge     | Security     | 3DS auth screen                | Critical |
| New device login   | Security     | Security settings              | Critical |
| Suspicious activity | Security   | Alert detail screen            | Critical |
| KYC approved      | Account      | Home dashboard                 | Medium   |
| KYC action needed | Account      | KYC screen                     | Medium   |
| Card shipped      | Account      | Card tab                       | Medium   |
| Card delivered    | Account      | Card activation                | Medium   |
| Scheduled payment | Transaction  | Payment detail                 | Medium   |
| Price alert       | Crypto       | Token detail                   | Low      |
| Product update    | Marketing    | Feature screen / webview       | Low      |

### 6.2 Notification Channels (Android)

- `transactions` -- All financial transactions (cannot be muted)
- `security` -- Security alerts and 3DS (cannot be muted)
- `account` -- Account updates, KYC, card delivery
- `crypto` -- Crypto-specific notifications
- `marketing` -- Product updates, offers (opt-in only)

---

## 7. Search Architecture

### 7.1 Global Search (Home Tab)

Accessible via search icon on Home screen. Searches across:
- Transactions (merchant name, reference, amount)
- Payees (name, IBAN)
- Crypto transactions (token, hash)
- Help articles (title, keywords)

Results grouped by category with "See all" per group.

### 7.2 Contextual Search

- Transaction history: filter + search bar at top
- Payee list: search by name
- Help center: search by keyword
- Country/currency selectors: search by name or code
