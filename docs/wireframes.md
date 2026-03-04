# TeslaPay Wireframes

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Senior UI/UX Designer, Dream Team

All wireframes are designed for a standard phone viewport (375x812px logical, iPhone 14 reference).
ASCII art represents layout structure; actual implementation follows the design system tokens.

---

## 1. Onboarding Flow

### 1.1 Splash Screen

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|              (( T ))                  |
|            TeslaPay                   |
|                                       |
|                                       |
|                                       |
|                                       |
|            [loading...]               |
|                                       |
|                                       |
+---------------------------------------+
```

Notes:
- TeslaPay logo centered vertically and horizontally
- Subtle pulse animation on logo
- Background: `color.primary.500` with subtle gradient
- Auto-advances after 2 seconds

### 1.2 Welcome Carousel

```
+---------------------------------------+
|            [status bar]               |
|                                   Skip|
|                                       |
|                                       |
|       +-------------------------+     |
|       |                         |     |
|       |    [Illustration:       |     |
|       |     multi-currency      |     |
|       |     globe with cards]   |     |
|       |                         |     |
|       +-------------------------+     |
|                                       |
|        Banking without borders        |
|                                       |
|    Multi-currency accounts, instant   |
|    SEPA transfers, and a Mastercard   |
|    that works everywhere in Europe.   |
|                                       |
|            o  .  .                    |
|                                       |
|   +-------------------------------+   |
|   |         Get Started           |   |
|   +-------------------------------+   |
|                                       |
|     I already have an account         |
|                                       |
+---------------------------------------+
```

Notes:
- Swipeable 3 slides. Dot indicators at bottom.
- Slide 2: Card illustration + "Your card, your rules"
- Slide 3: Crypto/blockchain illustration + "Crypto made simple"
- [Get Started] = primary button, [Already have account] = text link

### 1.3 Registration Screen

```
+---------------------------------------+
|            [status bar]               |
| <                                     |
|                                       |
|    Create Your Account                |
|                                       |
|    Step 1 of 3                        |
|    [====------] progress bar          |
|                                       |
|    Email address                      |
|    +-------------------------------+  |
|    | your@email.com                |  |
|    +-------------------------------+  |
|                                       |
|    Phone number                       |
|    +----+ +------------------------+  |
|    |+370| | 612 345 67             |  |
|    +----+ +------------------------+  |
|                                       |
|    [x] I agree to the Terms of        |
|        Service and Privacy Policy     |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |           Continue            |   |
|   +-------------------------------+   |
|                                       |
+---------------------------------------+
```

### 1.4 OTP Verification

```
+---------------------------------------+
|            [status bar]               |
| <                                     |
|                                       |
|    Verify Your Email                  |
|                                       |
|    We sent a code to                  |
|    eva@example.com                    |
|                                       |
|                                       |
|    +--+ +--+ +--+ +--+ +--+ +--+     |
|    | 4| | 7| | 2| |  | |  | |  |     |
|    +--+ +--+ +--+ +--+ +--+ +--+     |
|                                       |
|                                       |
|    Didn't receive it?                 |
|    Resend code (52s)                  |
|                                       |
|    Change email address               |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
+---------------------------------------+
```

Notes:
- 6 individual input boxes, auto-advance on digit entry
- Auto-submit when 6th digit entered
- Resend countdown timer, becomes tappable link at 0

### 1.5 Set PIN Screen

```
+---------------------------------------+
|            [status bar]               |
| <                                     |
|                                       |
|    Create Your PIN                    |
|                                       |
|    Choose a 6-digit PIN for           |
|    quick access to your account.      |
|                                       |
|                                       |
|         o  o  o  *  *  *             |
|                                       |
|    (3 digits entered)                 |
|                                       |
|                                       |
|   +-------+  +-------+  +-------+    |
|   |   1   |  |   2   |  |   3   |    |
|   +-------+  +-------+  +-------+    |
|   +-------+  +-------+  +-------+    |
|   |   4   |  |   5   |  |   6   |    |
|   +-------+  +-------+  +-------+    |
|   +-------+  +-------+  +-------+    |
|   |   7   |  |   8   |  |   9   |    |
|   +-------+  +-------+  +-------+    |
|   +-------+  +-------+  +-------+    |
|   |       |  |   0   |  |  <x   |    |
|   +-------+  +-------+  +-------+    |
+---------------------------------------+
```

### 1.6 Enable Biometric

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|                                       |
|                                       |
|            +--------+                 |
|            | (face) |                 |
|            +--------+                 |
|                                       |
|       Enable Face ID?                 |
|                                       |
|    Use Face ID to log in and          |
|    confirm transactions quickly       |
|    and securely.                      |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |        Enable Face ID         |   |
|   +-------------------------------+   |
|                                       |
|          Skip for now                 |
|                                       |
|                                       |
+---------------------------------------+
```

### 1.7 KYC Document Selection

```
+---------------------------------------+
|            [status bar]               |
| <                                     |
|                                       |
|    Verify Your Identity               |
|                                       |
|    Select your document type:         |
|                                       |
|    +-------------------------------+  |
|    | [icon] Passport               |  |
|    | Recommended -- fastest        |  |
|    +-------------------------------+  |
|                                       |
|    +-------------------------------+  |
|    | [icon] National ID Card       |  |
|    +-------------------------------+  |
|                                       |
|    +-------------------------------+  |
|    | [icon] Driver's License       |  |
|    +-------------------------------+  |
|                                       |
|    +-------------------------------+  |
|    | [icon] Residence Permit       |  |
|    +-------------------------------+  |
|                                       |
|    You'll need your document and      |
|    a well-lit space. Takes ~2 min.    |
|                                       |
+---------------------------------------+
```

### 1.8 KYC Processing

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|           [animated spinner]          |
|                                       |
|       Verifying your identity...      |
|                                       |
|       This usually takes less         |
|       than a minute.                  |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
|                                       |
+---------------------------------------+
```

### 1.9 Account Created Success

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|          [check animation]            |
|                                       |
|    Welcome to TeslaPay, Eva!          |
|                                       |
|    Your account is ready.             |
|                                       |
|    +-------------------------------+  |
|    | Your IBAN                     |  |
|    | LT12 3456 7890 1234 5678     |  |
|    |                     [Copy]    |  |
|    +-------------------------------+  |
|                                       |
|    Account tier: Basic                |
|    [View tier details]                |
|                                       |
|   +-------------------------------+   |
|   |         Add Funds             |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   |       Get Your Card           |   |
|   +-------------------------------+   |
|                                       |
|         Explore the App               |
|                                       |
+---------------------------------------+
```

---

## 2. Home / Dashboard

### 2.1 Dashboard (Main State)

```
+---------------------------------------+
|            [status bar]               |
| TeslaPay                    (!) [bell]|
|                                       |
| +-----------------------------------+ |
| |                                   | |
| | Total Balance                     | |
| |                                   | |
| |     EUR 3,245.67                  | |
| |     +0.42% today                  | |
| |                                   | |
| +-----------------------------------+ |
|                                       |
|  [Send]  [Request] [Exchange] [Top Up]|
|   (->)     (<-)      (<>)      (+)   |
|                                       |
| Accounts                    See all > |
| +-----------------------------------+ |
| | [EU] EUR      EUR 2,845.67       | |
| | [US] USD      USD 420.00         | |
| |      + Add currency               | |
| +-----------------------------------+ |
|                                       |
| Recent Transactions         See all > |
| +-----------------------------------+ |
| | [S] Starbucks      -EUR 4.50     | |
| |     Card payment    Today 09:15   | |
| +-----------------------------------+ |
| | [A] Anna Kowalski  -EUR 200.00   | |
| |     SEPA transfer   Yesterday    | |
| +-----------------------------------+ |
| | [M] Monthly Salary +EUR 3,200.00 | |
| |     SEPA received   28 Feb       | |
| +-----------------------------------+ |
| | [F] FUSE Purchase  -EUR 20.00    | |
| |     Crypto buy      28 Feb       | |
| +-----------------------------------+ |
|                                       |
| Crypto Portfolio              View >  |
| +-----------------------------------+ |
| | Total: EUR 19.70     +3.2% 24h   | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
|  [*]   |   [ ]    | [ ] |  [ ]   |[ ]|
+---------------------------------------+
```

Notes:
- Balance card uses `gradient.brand` background, white text
- Quick action icons are circular with `color.primary.100` background
- Transaction items: 72px height, left icon circle is merchant initial with category color
- Positive amounts in `color.success.500`, negative in `color.neutral.900`
- Pull-to-refresh updates all data
- Crypto section only shows if user has activated wallet

### 2.2 Dashboard (Empty State -- New User)

```
+---------------------------------------+
|            [status bar]               |
| TeslaPay                    (!) [bell]|
|                                       |
| +-----------------------------------+ |
| |                                   | |
| | Total Balance                     | |
| |                                   | |
| |       EUR 0.00                    | |
| |                                   | |
| +-----------------------------------+ |
|                                       |
|  [Send]  [Request] [Exchange] [Top Up]|
|   (->)     (<-)      (<>)      (+)   |
|                                       |
| +-----------------------------------+ |
| |     [illustration: empty wallet]  | |
| |                                   | |
| |   Add funds to get started        | |
| |                                   | |
| |   Share your IBAN to receive      | |
| |   transfers, or top up by card.   | |
| |                                   | |
| |   +---------------------------+   | |
| |   |    Add Funds              |   | |
| |   +---------------------------+   | |
| |                                   | |
| |   +---------------------------+   | |
| |   |    Share Your IBAN        |   | |
| |   +---------------------------+   | |
| +-----------------------------------+ |
|                                       |
| Getting Started                       |
| +-----------------------------------+ |
| | [x] Create account                | |
| | [x] Verify identity               | |
| | [ ] Add funds                     | |
| | [ ] Get your card                 | |
| | [ ] Try crypto                    | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

---

## 3. Accounts Screen

### 3.1 Account Detail (EUR)

```
+---------------------------------------+
|            [status bar]               |
| <  EUR Account                        |
|                                       |
|         EUR 2,845.67                  |
|                                       |
|    IBAN: LT12 3456 7890 1234 5678    |
|                            [Copy]     |
|                                       |
|    [Send]    [Receive]   [Exchange]   |
|                                       |
| +-----------------------------------+ |
| | [Filter v] [Search...]            | |
| +-----------------------------------+ |
|                                       |
| Today                                 |
| +-----------------------------------+ |
| | [S] Starbucks         -EUR 4.50  | |
| |     Card payment       09:15     | |
| +-----------------------------------+ |
|                                       |
| Yesterday                             |
| +-----------------------------------+ |
| | [A] Anna Kowalski   -EUR 200.00  | |
| |     SEPA transfer     14:32      | |
| +-----------------------------------+ |
| | [L] Landlord        -EUR 800.00  | |
| |     Recurring         08:00      | |
| +-----------------------------------+ |
|                                       |
| 28 February                           |
| +-----------------------------------+ |
| | [C] Company AG     +EUR 3,200.00 | |
| |     SEPA received     09:45      | |
| +-----------------------------------+ |
| | [E] Exchange        -EUR 500.00  | |
| |     EUR to PLN        11:20      | |
| +-----------------------------------+ |
|                                       |
|     [Export CSV]  [Export PDF]         |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

---

## 4. Payments

### 4.1 Payments Tab (Root)

```
+---------------------------------------+
|            [status bar]               |
|    Payments                           |
|                                       |
|  +--------+ +--------+ +--------+    |
|  |  Send  | |Request | |Exchange|    |
|  |  (->)  | |  (<-)  | |  (<>) |    |
|  +--------+ +--------+ +--------+    |
|                        +--------+    |
|                        |Scan QR |    |
|                        | [QR]   |    |
|                        +--------+    |
|                                       |
| Saved Payees                Manage >  |
| +-----------------------------------+ |
| | (@) Anna Kowalski                 | |
| |     EE12 **** **** **** 90       | |
| +-----------------------------------+ |
| | (@) Landlord Berlin               | |
| |     DE89 **** **** **** 00       | |
| +-----------------------------------+ |
| | (@) Mom                           | |
| |     LT60 **** **** **** 78      | |
| +-----------------------------------+ |
| |     + Add new payee              | |
| +-----------------------------------+ |
|                                       |
| Scheduled                   See all > |
| +-----------------------------------+ |
| | Rent -- EUR 800       Monthly     | |
| | Next: 1 Apr 2026                  | |
| +-----------------------------------+ |
| | Netflix -- EUR 13.99   Monthly    | |
| | Next: 15 Mar 2026                 | |
| +-----------------------------------+ |
|                                       |
| Direct Debits               Manage >  |
| +-----------------------------------+ |
| | Vodafone LT     Active           | |
| | Next: ~EUR 25   12 Mar 2026      | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

### 4.2 Send Money -- Amount Entry

```
+---------------------------------------+
|            [status bar]               |
| X  Send Money                         |
|                                       |
|    To: Anna Kowalski                  |
|        EE12 7890 1234 5678 90        |
|        [Change]                       |
|                                       |
|    From: [EU] EUR Account             |
|           EUR 2,845.67 available      |
|                                       |
|                                       |
|           EUR 200.00                  |
|                                       |
|    Reference (optional)               |
|    +-------------------------------+  |
|    | March rent                    |  |
|    +-------------------------------+  |
|                                       |
|    Transfer speed                     |
|    (o) Instant     Free        10s    |
|    ( ) Standard    Free       1 day   |
|                                       |
|    +-------------------------------+  |
|    | Fee:           EUR 0.00      |  |
|    | Recipient gets: EUR 200.00   |  |
|    +-------------------------------+  |
|                                       |
|   +-------------------------------+   |
|   |          Continue             |   |
|   +-------------------------------+   |
|                                       |
|   +---+ +---+ +---+ +---+ +---+      |
|   | 1 | | 2 | | 3 | | 4 | | 5 |     |
|   +---+ +---+ +---+ +---+ +---+      |
|   | 6 | | 7 | | 8 | | 9 | | 0 |     |
|   +---+ +---+ +---+ +---+ +---+      |
|               | . | |<x |            |
|               +---+ +---+            |
+---------------------------------------+
```

### 4.3 Send Money -- Review

```
+---------------------------------------+
|            [status bar]               |
| <  Review Transfer                    |
|                                       |
|   +-------------------------------+   |
|   |                               |   |
|   | From    EUR Account           |   |
|   |         LT12 3456 7890...     |   |
|   |                               |   |
|   | To      Anna Kowalski         |   |
|   |         EE12 7890 1234...     |   |
|   |                               |   |
|   | Amount  EUR 200.00            |   |
|   | Fee     EUR 0.00              |   |
|   | Total   EUR 200.00            |   |
|   |                               |   |
|   | Speed   Instant (SEPA Inst)   |   |
|   | Ref     March rent            |   |
|   |                               |   |
|   +-------------------------------+   |
|                                       |
|                                       |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |     Confirm with Face ID      |   |
|   |         [face icon]           |   |
|   +-------------------------------+   |
|                                       |
+---------------------------------------+
```

### 4.4 Transfer Success

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|                                       |
|                                       |
|          [animated check]             |
|                                       |
|        Transfer Successful            |
|                                       |
|    EUR 200.00 sent to                 |
|    Anna Kowalski                      |
|                                       |
|    Arrived instantly                  |
|    2 March 2026, 14:32 CET           |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |       Share Receipt           |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   |       Send Another            |   |
|   +-------------------------------+   |
|                                       |
|          Back to Home                 |
|                                       |
+---------------------------------------+
```

---

## 5. Card Management

### 5.1 Card Tab (Active Card)

```
+---------------------------------------+
|            [status bar]               |
|    Card                               |
|                                       |
|   +-------------------------------+   |
|   |  [TeslaPay]       [Mastercard]|   |
|   |                               |   |
|   |                               |   |
|   |  **** **** **** 4521          |   |
|   |                               |   |
|   |  EVA TAMM          12/29     |   |
|   +-------------------------------+   |
|   (o virtual)  ( physical)            |
|                                       |
|  +--------+ +---------+ +---------+  |
|  | Show   | | Freeze  | |Apple Pay|  |
|  |Details | |   [*]   | |  [+]   |  |
|  +--------+ +---------+ +---------+  |
|                                       |
| Card Settings                         |
| +-----------------------------------+ |
| | [icon] View PIN                >  | |
| +-----------------------------------+ |
| | [icon] Change PIN              >  | |
| +-----------------------------------+ |
| | [icon] Spending Limits         >  | |
| +-----------------------------------+ |
| | [icon] Security Controls       >  | |
| +-----------------------------------+ |
| | [icon] Linked Account         >  | |
| +-----------------------------------+ |
|                                       |
| Recent Card Transactions    See all > |
| +-----------------------------------+ |
| | [S] Starbucks         -EUR 4.50  | |
| |     Contactless        Today     | |
| +-----------------------------------+ |
| | [A] Amazon.de        -EUR 29.99  | |
| |     Online             Yesterday | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| | [icon] Order Physical Card     >  | |
| +-----------------------------------+ |
| | [icon] Report Lost / Stolen    >  | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

Notes:
- Card visual uses `gradient.card`, white text, Mastercard logo
- Card is swipeable if user has multiple cards (virtual + physical)
- Freeze toggle is a prominent switch with icon animation
- When frozen: card visual gets gray overlay with snowflake icon

### 5.2 Card Tab (No Card Yet)

```
+---------------------------------------+
|            [status bar]               |
|    Card                               |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |                               |   |
|   |    [card illustration]        |   |
|   |                               |   |
|   |  Get your TeslaPay            |   |
|   |  Mastercard                   |   |
|   |                               |   |
|   |  Pay online, in stores,       |   |
|   |  and withdraw cash at any     |   |
|   |  ATM worldwide.               |   |
|   |                               |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   |   Get Virtual Card (Instant)  |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   |   Order Physical Card         |   |
|   +-------------------------------+   |
|                                       |
|   Features:                           |
|   [check] Apple Pay & Google Pay      |
|   [check] Real-time notifications     |
|   [check] Spending controls           |
|   [check] Freeze / unfreeze          |
|   [check] 3D Secure protection        |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

### 5.3 Card Frozen State

```
+---------------------------------------+
|            [status bar]               |
|    Card                               |
|                                       |
|   +-------------------------------+   |
|   |  [TeslaPay]       [Mastercard]|   |
|   |        CARD FROZEN            |   |
|   |         [snowflake]           |   |
|   |  **** **** **** 4521          |   |
|   |                               |   |
|   |  EVA TAMM          12/29     |   |
|   +-------------------------------+   |
|    (gray overlay over entire card)    |
|                                       |
|  +--------+ +---------+ +---------+  |
|  | Show   | |Unfreeze | |Apple Pay|  |
|  |Details | |   [*]   | |  [+]   |  |
|  +--------+ +---------+ +---------+  |
|                                       |
| +-----------------------------------+ |
| | (!) Card is frozen. All new       | |
| |     transactions will be declined.| |
| |     Tap Unfreeze to restore.     | |
| +-----------------------------------+ |
|                                       |
```

### 5.4 Spending Limits

```
+---------------------------------------+
|            [status bar]               |
| <  Spending Limits                    |
|                                       |
| Per Transaction                       |
| +-----------------------------------+ |
| |  EUR 0 ----[====o----]--- EUR 5K  | |
| |  Current: EUR 2,500               | |
| +-----------------------------------+ |
|                                       |
| Daily Limit                           |
| +-----------------------------------+ |
| |  EUR 0 ---[======o---]--- EUR 10K | |
| |  Current: EUR 5,000               | |
| |  Used today: EUR 234.49           | |
| +-----------------------------------+ |
|                                       |
| Monthly Limit                         |
| +-----------------------------------+ |
| |  EUR 0 -[========o--]--- EUR 50K  | |
| |  Current: EUR 25,000              | |
| |  Used this month: EUR 1,834.49    | |
| +-----------------------------------+ |
|                                       |
| ATM Withdrawal                        |
| +-----------------------------------+ |
| |  EUR 0 --[===o------]--- EUR 2K   | |
| |  Current: EUR 500                  | |
| |  Free left: 3 of 5 withdrawals   | |
| +-----------------------------------+ |
|                                       |
|   +-------------------------------+   |
|   |         Save Changes          |   |
|   +-------------------------------+   |
|                                       |
| Limits depend on your account tier.   |
| Current tier: Basic  [Upgrade]        |
|                                       |
+---------------------------------------+
```

### 5.5 Security Controls

```
+---------------------------------------+
|            [status bar]               |
| <  Security Controls                  |
|                                       |
| Payment Methods                       |
| +-----------------------------------+ |
| | Online payments          [====]   | |
| +-----------------------------------+ |
| | Contactless (NFC)        [====]   | |
| +-----------------------------------+ |
| | ATM withdrawals          [====]   | |
| +-----------------------------------+ |
| | Magnetic stripe          [    ]   | |
| +-----------------------------------+ |
|                                       |
| Merchant Categories                   |
| +-----------------------------------+ |
| | Gambling                 Blocked  | |
| +-----------------------------------+ |
| | Adult content            Blocked  | |
| +-----------------------------------+ |
| | Cryptocurrency exchanges Allowed  | |
| +-----------------------------------+ |
| | [Manage all categories >]         | |
| +-----------------------------------+ |
|                                       |
| Geographic Restrictions               |
| +-----------------------------------+ |
| | Allowed regions: Europe (EEA)     | |
| | [Change regions >]                | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

---

## 6. Crypto Wallet

### 6.1 Crypto Tab (Active)

```
+---------------------------------------+
|            [status bar]               |
|    Crypto                      [?]    |
|                                       |
| +-----------------------------------+ |
| |  Total Crypto Balance             | |
| |                                   | |
| |      EUR 519.70                   | |
| |      +2.8% (24h)                 | |
| +-----------------------------------+ |
|                                       |
|  [Buy]   [Sell]   [Send]  [Receive]  |
|                                       |
| Your Tokens                           |
| +-----------------------------------+ |
| | (F) FUSE                          | |
| |     469.05 FUSE                   | |
| |     EUR 19.70           +3.2%    | |
| +-----------------------------------+ |
| | (U) USDC                          | |
| |     500.00 USDC                   | |
| |     EUR 460.00          +0.01%   | |
| +-----------------------------------+ |
| | (T) USDT                          | |
| |     43.50 USDT                    | |
| |     EUR 40.00           -0.02%   | |
| +-----------------------------------+ |
|                                       |
| Earn Yield                            |
| +-----------------------------------+ |
| | [star] Solid soUSD    APY: 4.8%  | |
| |        Earn yield on stablecoins  | |
| |        [Start Earning >]          | |
| +-----------------------------------+ |
|                                       |
| Recent Activity               All >   |
| +-----------------------------------+ |
| | [B] Bought FUSE     -EUR 20.00   | |
| |     469.05 FUSE      28 Feb      | |
| +-----------------------------------+ |
| | [D] Deposited USDC  +500 USDC    | |
| |     From 0x1a2b...   25 Feb      | |
| +-----------------------------------+ |
|                                       |
| Wallet                                |
| +-----------------------------------+ |
| | Your address: 0x1a2b...3c4d      | |
| |               [Copy] [QR]        | |
| | Network: Fuse   [green dot] Live  | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

### 6.2 Crypto Tab (First Visit / Empty)

```
+---------------------------------------+
|            [status bar]               |
|    Crypto                      [?]    |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |                               |   |
|   |    [blockchain illustration]  |   |
|   |                               |   |
|   |  Your crypto wallet is ready  |   |
|   |                               |   |
|   |  Buy, sell, and send crypto   |   |
|   |  right from your TeslaPay     |   |
|   |  account. Powered by Fuse     |   |
|   |  blockchain.                  |   |
|   |                               |   |
|   +-------------------------------+   |
|                                       |
|   What you can do:                    |
|                                       |
|   [check] Buy FUSE, USDC, and USDT   |
|           starting from EUR 5         |
|                                       |
|   [check] Send crypto globally        |
|           with near-zero fees         |
|                                       |
|   [check] Earn yield on stablecoins   |
|           (coming soon)               |
|                                       |
|   [check] Self-custody: your keys,    |
|           your coins                  |
|                                       |
|   +-------------------------------+   |
|   |       Buy Your First Crypto   |   |
|   +-------------------------------+   |
|                                       |
|       Learn about crypto risks        |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

### 6.3 Buy Crypto -- Amount

```
+---------------------------------------+
|            [status bar]               |
| X  Buy FUSE                          |
|                                       |
|    Pay with                           |
|    [EU] EUR Account   EUR 2,845.67   |
|                                       |
|                                       |
|             EUR 20.00                 |
|           ~ 469.05 FUSE              |
|                                       |
|    [EUR 5] [EUR 20] [EUR 50] [EUR100]|
|                                       |
|                                       |
|    +-------------------------------+  |
|    | Rate   1 FUSE = EUR 0.0426   |  |
|    | Fee    EUR 0.30 (1.5%)       |  |
|    | You get ~469.05 FUSE          |  |
|    +-------------------------------+  |
|                                       |
|   +-------------------------------+   |
|   |          Continue             |   |
|   +-------------------------------+   |
|                                       |
|   +---+ +---+ +---+ +---+ +---+      |
|   | 1 | | 2 | | 3 | | 4 | | 5 |     |
|   +---+ +---+ +---+ +---+ +---+      |
|   | 6 | | 7 | | 8 | | 9 | | 0 |     |
|   +---+ +---+ +---+ +---+ +---+      |
|               | . | |<x |            |
|               +---+ +---+            |
+---------------------------------------+
```

### 6.4 Receive Crypto

```
+---------------------------------------+
|            [status bar]               |
| <  Receive Crypto                     |
|                                       |
|    Select token to receive:           |
|    [FUSE v]                           |
|                                       |
|   +-------------------------------+   |
|   |                               |   |
|   |      +------------------+     |   |
|   |      |                  |     |   |
|   |      |   [QR CODE]      |     |   |
|   |      |                  |     |   |
|   |      |                  |     |   |
|   |      +------------------+     |   |
|   |                               |   |
|   +-------------------------------+   |
|                                       |
|    Your Fuse wallet address:          |
|    0x1a2b3c4d5e6f7890abcd            |
|    ef1234567890abcd3c4d              |
|                                       |
|    [Copy Address]  [Share]            |
|                                       |
| +-----------------------------------+ |
| | (!) Only send FUSE tokens on the  | |
| |     Fuse network to this address. | |
| |     Sending other tokens or using | |
| |     other networks may result in  | |
| |     permanent loss of funds.      | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

### 6.5 DeFi Yield Dashboard (Phase 2)

```
+---------------------------------------+
|            [status bar]               |
| <  Earn Yield                         |
|                                       |
| +-----------------------------------+ |
| |  Solid soUSD                      | |
| |                                   | |
| |  Deposited     500.00 soUSD       | |
| |  Current value EUR 461.34         | |
| |  Earned        1.34 USDC          | |
| |  APY           4.8%               | |
| |  Duration      20 days            | |
| +-----------------------------------+ |
|                                       |
|  +-------------+ +-------------+      |
|  | Deposit More| |  Withdraw   |      |
|  +-------------+ +-------------+      |
|                                       |
| Yield History                         |
| +-----------------------------------+ |
| | Today         +0.066 USDC        | |
| | Yesterday     +0.066 USDC        | |
| | 28 Feb        +0.066 USDC        | |
| | 27 Feb        +0.065 USDC        | |
| | ...                               | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| | (i) Yield is generated by the     | |
| |     Solid protocol on the Fuse    | |
| |     network. Funds are not        | |
| |     protected by deposit guarantee| |
| |     schemes. Learn more           | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

---

## 7. Transaction History and Details

### 7.1 Full Transaction History

```
+---------------------------------------+
|            [status bar]               |
| <  Transactions                       |
|                                       |
| +-----------------------------------+ |
| | [Search transactions...]          | |
| +-----------------------------------+ |
| | [All v] [Date v] [Amount v] [+]  | |
| +-----------------------------------+ |
|                                       |
| Today                                 |
| +-----------------------------------+ |
| | [S] Starbucks Berlin  -EUR 4.50  | |
| |     Card - Contactless  09:15     | |
| +-----------------------------------+ |
|                                       |
| Yesterday -- 1 March                  |
| +-----------------------------------+ |
| | [A] Anna Kowalski   -EUR 200.00  | |
| |     SEPA Instant Out   14:32      | |
| +-----------------------------------+ |
| | [L] Landlord Berlin  -EUR 800.00 | |
| |     SEPA Recurring     08:00      | |
| +-----------------------------------+ |
| | [R] Rewe Supermarket  -EUR 34.67 | |
| |     Card - Contactless  18:45     | |
| +-----------------------------------+ |
|                                       |
| 28 February                           |
| +-----------------------------------+ |
| | [C] Company AG     +EUR 3,200.00 | |
| |     SEPA In            09:45      | |
| +-----------------------------------+ |
| | [E] EUR > PLN       -EUR 500.00  | |
| |     Exchange           11:20      | |
| +-----------------------------------+ |
| | [F] Buy FUSE         -EUR 20.30  | |
| |     Crypto buy         15:10      | |
| +-----------------------------------+ |
|                                       |
| ... (infinite scroll) ...             |
|                                       |
|    [Export: CSV | PDF]                 |
|                                       |
+---------------------------------------+
```

### 7.2 Transaction Detail (Card Payment)

```
+---------------------------------------+
|            [status bar]               |
| <  Transaction                        |
|                                       |
|              [S]                      |
|           Starbucks                   |
|         Berlin, Germany               |
|                                       |
|          -EUR 4.50                    |
|        Completed                      |
|                                       |
| +-----------------------------------+ |
| |                                   | |
| | Type       Card payment           | |
| | Method     Contactless (NFC)      | |
| | Card       **** 4521              | |
| |                                   | |
| | Category   Food & Drink           | |
| |            [Change >]             | |
| |                                   | |
| | Date       2 March 2026           | |
| | Time       09:15 CET              | |
| |                                   | |
| | Fee        EUR 0.00               | |
| | FX rate    N/A                     | |
| |                                   | |
| | Status     Completed              | |
| | Ref        MC-2026030-XY12        | |
| |                                   | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| | [icon] Repeat this payment     >  | |
| +-----------------------------------+ |
| | [icon] Share receipt           >  | |
| +-----------------------------------+ |
| | [icon] Dispute transaction     >  | |
| +-----------------------------------+ |
| | [icon] Report a problem        >  | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

### 7.3 Transaction Detail (SEPA Transfer)

```
+---------------------------------------+
|            [status bar]               |
| <  Transaction                        |
|                                       |
|              [A]                      |
|         Anna Kowalski                 |
|     EE12 7890 1234 5678 90           |
|                                       |
|         -EUR 200.00                   |
|        Completed                      |
|                                       |
| +-----------------------------------+ |
| |                                   | |
| | Type       SEPA Instant Transfer  | |
| | Direction  Outgoing               | |
| |                                   | |
| | Recipient  Anna Kowalski          | |
| | IBAN       EE12 7890 1234 5678 90 | |
| | BIC        HABAEE2X               | |
| |                                   | |
| | Reference  March rent             | |
| |                                   | |
| | Date       1 March 2026           | |
| | Time       14:32:05 CET           | |
| |                                   | |
| | Fee        EUR 0.00               | |
| | Status     Settled                | |
| |                                   | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| | [icon] Send again              >  | |
| +-----------------------------------+ |
| | [icon] Share receipt           >  | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

---

## 8. Profile and Settings

### 8.1 Profile Tab (Root)

```
+---------------------------------------+
|            [status bar]               |
|    Profile                            |
|                                       |
|   +-------------------------------+   |
|   | (@)  Eva Tamm                 |   |
|   |      eva@example.com          |   |
|   |      [Basic] tier badge       |   |
|   +-------------------------------+   |
|                                       |
| ACCOUNT                               |
| +-----------------------------------+ |
| | [icon] Personal Information    >  | |
| +-----------------------------------+ |
| | [icon] Verification Status     >  | |
| |        [green] Verified           | |
| +-----------------------------------+ |
| | [icon] Account Tier            >  | |
| |        Basic - Upgrade available  | |
| +-----------------------------------+ |
| | [icon] Fees and Limits         >  | |
| +-----------------------------------+ |
|                                       |
| SECURITY                              |
| +-----------------------------------+ |
| | [icon] Face ID          [====]    | |
| +-----------------------------------+ |
| | [icon] Change PIN              >  | |
| +-----------------------------------+ |
| | [icon] Change Password         >  | |
| +-----------------------------------+ |
| | [icon] Active Sessions         >  | |
| +-----------------------------------+ |
|                                       |
| PREFERENCES                           |
| +-----------------------------------+ |
| | [icon] Notifications           >  | |
| +-----------------------------------+ |
| | [icon] Language         English>  | |
| +-----------------------------------+ |
| | [icon] Appearance       System >  | |
| +-----------------------------------+ |
|                                       |
| SUPPORT                               |
| +-----------------------------------+ |
| | [icon] Help Center             >  | |
| +-----------------------------------+ |
| | [icon] Chat with Us            >  | |
| +-----------------------------------+ |
|                                       |
| LEGAL                                 |
| +-----------------------------------+ |
| | [icon] Terms of Service        >  | |
| +-----------------------------------+ |
| | [icon] Privacy Policy          >  | |
| +-----------------------------------+ |
| | [icon] GDPR Data Requests      >  | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| | [icon] Close Account           >  | |
| +-----------------------------------+ |
|                                       |
|           [Log Out]                   |
|                                       |
|        Version 1.0.0 (42)            |
|                                       |
+---------------------------------------+
|  Home  | Payments | Card | Crypto | Me|
+---------------------------------------+
```

### 8.2 Fees and Limits Screen

```
+---------------------------------------+
|            [status bar]               |
| <  Fees and Limits                    |
|                                       |
|    Account tier: Basic                |
|    [Upgrade to Standard >]           |
|                                       |
| TRANSFERS                             |
| +-----------------------------------+ |
| | SEPA transfer        Free         | |
| | SEPA Instant         Free         | |
| | Internal transfer    Free         | |
| +-----------------------------------+ |
|                                       |
| CARD                                  |
| +-----------------------------------+ |
| | Virtual card         Free         | |
| | Physical card        Free (1st)   | |
| | Card replacement     EUR 10       | |
| | ATM (free)          5/month       | |
| |   Used: 2 of 5                    | |
| | ATM (over limit)    2%            | |
| +-----------------------------------+ |
|                                       |
| CURRENCY EXCHANGE                     |
| +-----------------------------------+ |
| | FX markup           0.5%          | |
| | Monthly FX free     EUR 1,000     | |
| |   Used: EUR 500 of EUR 1,000     | |
| |   [=======-----]                  | |
| | Over limit          1.0%          | |
| +-----------------------------------+ |
|                                       |
| CRYPTO                                |
| +-----------------------------------+ |
| | Buy/Sell fee         1.5%         | |
| | Send fee             Network fee  | |
| | Receive              Free         | |
| +-----------------------------------+ |
|                                       |
| LIMITS                                |
| +-----------------------------------+ |
| | Daily transfer      EUR 5,000     | |
| | Monthly transfer    EUR 25,000    | |
| | Daily card spend    EUR 5,000     | |
| | ATM per withdrawal  EUR 500       | |
| | Crypto buy/day      EUR 1,000     | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

### 8.3 Active Sessions

```
+---------------------------------------+
|            [status bar]               |
| <  Active Sessions                    |
|                                       |
| Current Session                       |
| +-----------------------------------+ |
| | [phone] iPhone 14 Pro             | |
| |   iOS 19.2                        | |
| |   Berlin, Germany                 | |
| |   Active now                      | |
| |   This device                     | |
| +-----------------------------------+ |
|                                       |
| Other Sessions                        |
| +-----------------------------------+ |
| | [tablet] iPad Air                 | |
| |   iPadOS 19.2                     | |
| |   Berlin, Germany                 | |
| |   Last active: 2h ago            | |
| |               [Terminate]         | |
| +-----------------------------------+ |
|                                       |
|                                       |
|   +-------------------------------+   |
|   |   Terminate All Other Sessions|   |
|   +-------------------------------+   |
|                                       |
+---------------------------------------+
```

### 8.4 Notification Preferences

```
+---------------------------------------+
|            [status bar]               |
| <  Notifications                      |
|                                       |
| TRANSACTIONS (always on)              |
| +-----------------------------------+ |
| | Card payments        [====]  [L]  | |
| | Incoming transfers   [====]  [L]  | |
| | Outgoing transfers   [====]  [L]  | |
| | Crypto activity      [====]       | |
| +-----------------------------------+ |
| [L] = locked, cannot disable         |
|                                       |
| SECURITY (always on)                  |
| +-----------------------------------+ |
| | New device login     [====]  [L]  | |
| | Suspicious activity  [====]  [L]  | |
| | 3D Secure            [====]  [L]  | |
| +-----------------------------------+ |
|                                       |
| ACCOUNT                               |
| +-----------------------------------+ |
| | KYC updates          [====]       | |
| | Card delivery        [====]       | |
| | Scheduled payments   [====]       | |
| +-----------------------------------+ |
|                                       |
| MARKETING                             |
| +-----------------------------------+ |
| | Product updates      [    ]       | |
| | Offers and rewards   [    ]       | |
| +-----------------------------------+ |
|                                       |
| Channels                              |
| +-----------------------------------+ |
| | Push notifications   [====]       | |
| | Email                [====]       | |
| | SMS                  [    ]       | |
| +-----------------------------------+ |
|                                       |
+---------------------------------------+
```

---

## 9. Push Notification Examples

### 9.1 Card Transaction (Lock Screen)

```
+---------------------------------------+
|                                       |
|  TeslaPay                     now     |
|  Card Payment: EUR 4.50 at           |
|  Starbucks, Berlin.                  |
|  Tap for details.                     |
|                                       |
+---------------------------------------+
```

### 9.2 Incoming SEPA Transfer

```
+---------------------------------------+
|                                       |
|  TeslaPay                    2m ago   |
|  Money Received: EUR 3,200.00 from   |
|  Company AG. Reference: Salary Mar.  |
|                                       |
+---------------------------------------+
```

### 9.3 Card Declined

```
+---------------------------------------+
|                                       |
|  TeslaPay                     now     |
|  Card Declined: EUR 150.00 at        |
|  Online Store. Reason: exceeds daily  |
|  limit. Tap to adjust limits.        |
|                                       |
+---------------------------------------+
```

### 9.4 3D Secure Challenge

```
+---------------------------------------+
|                                       |
|  TeslaPay                     now     |
|  Payment Approval Required:           |
|  EUR 89.99 at Amazon.de.             |
|  Tap to approve or decline.          |
|  [Approve]  [Decline]                |
|                                       |
+---------------------------------------+
```

### 9.5 Crypto Received

```
+---------------------------------------+
|                                       |
|  TeslaPay                     now     |
|  Crypto Received: 100.00 USDC        |
|  (~EUR 92.00) from 0x1a2b...3c4d.   |
|                                       |
+---------------------------------------+
```

### 9.6 Suspicious Activity

```
+---------------------------------------+
|                                       |
|  TeslaPay                     now     |
|  Security Alert: New login detected   |
|  from Samsung Galaxy, Warsaw, PL.     |
|  Was this you? Tap to review.        |
|  [Yes, it's me]  [Secure Account]    |
|                                       |
+---------------------------------------+
```

### 9.7 Card Shipped

```
+---------------------------------------+
|                                       |
|  TeslaPay                   1h ago    |
|  Your TeslaPay Mastercard has been    |
|  shipped! Expected delivery:          |
|  8-12 March 2026. Track: LP12345678. |
|                                       |
+---------------------------------------+
```

### 9.8 KYC Approved

```
+---------------------------------------+
|                                       |
|  TeslaPay                     now     |
|  Identity Verified! Your TeslaPay     |
|  account is now active. Add funds     |
|  or get your card to get started.    |
|                                       |
+---------------------------------------+
```

---

## 10. Login Screen

### 10.1 Biometric Login (Returning User)

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|                                       |
|              (( T ))                  |
|            TeslaPay                   |
|                                       |
|                                       |
|       Welcome back, Eva              |
|                                       |
|                                       |
|           +--------+                  |
|           | (face) |                  |
|           +--------+                  |
|        Tap to unlock                  |
|                                       |
|                                       |
|                                       |
|        Use PIN instead                |
|                                       |
|        Use password                   |
|                                       |
|     Not Eva? Switch account           |
|                                       |
|                                       |
+---------------------------------------+
```

### 10.2 PIN Login (Fallback)

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|              (( T ))                  |
|            TeslaPay                   |
|                                       |
|                                       |
|       Enter your PIN                  |
|                                       |
|                                       |
|         o  o  *  *  *  *             |
|                                       |
|                                       |
|   +-------+  +-------+  +-------+    |
|   |   1   |  |   2   |  |   3   |    |
|   +-------+  +-------+  +-------+    |
|   +-------+  +-------+  +-------+    |
|   |   4   |  |   5   |  |   6   |    |
|   +-------+  +-------+  +-------+    |
|   +-------+  +-------+  +-------+    |
|   |   7   |  |   8   |  |   9   |    |
|   +-------+  +-------+  +-------+    |
|   +-------+  +-------+  +-------+    |
|   |(face) |  |   0   |  |  <x   |    |
|   +-------+  +-------+  +-------+    |
|                                       |
|        Forgot PIN?                    |
+---------------------------------------+
```

Notes:
- Bottom-left numpad key shows biometric icon for switching back
- After 5 failed PIN attempts, locked for 30 minutes

---

## 11. 3D Secure In-App Challenge

```
+---------------------------------------+
|            [status bar]               |
|                                       |
|                                       |
|    Payment Approval                   |
|                                       |
|   +-------------------------------+   |
|   |                               |   |
|   | Merchant   Amazon.de          |   |
|   | Amount     EUR 89.99          |   |
|   | Card       **** 4521          |   |
|   | Date       2 March 2026       |   |
|   |                               |   |
|   +-------------------------------+   |
|                                       |
|    Do you approve this payment?       |
|                                       |
|    Expires in: 4:32                   |
|                                       |
|   +-------------------------------+   |
|   |    Approve with Face ID       |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   |         Decline               |   |
|   +-------------------------------+   |
|                                       |
|    Not your purchase?                 |
|    Report fraud and freeze card       |
|                                       |
+---------------------------------------+
```

---

## 12. Exchange Currency

```
+---------------------------------------+
|            [status bar]               |
| <  Exchange                           |
|                                       |
|    From                               |
|   +-------------------------------+   |
|   | [EU flag] EUR   EUR 2,845.67  |   |
|   |                               |   |
|   |         500.00                |   |
|   +-------------------------------+   |
|                                       |
|              [swap icon]              |
|                                       |
|    To                                 |
|   +-------------------------------+   |
|   | [PL flag] PLN   PLN 0.00     |   |
|   |                               |   |
|   |       2,147.50               |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   | Rate      1 EUR = 4.2950 PLN |   |
|   | Markup    0.23%              |   |
|   | Mid-rate  1 EUR = 4.2850 PLN |   |
|   | Refreshes in: 28s            |   |
|   +-------------------------------+   |
|                                       |
|   +-------------------------------+   |
|   |     Exchange EUR 500.00       |   |
|   +-------------------------------+   |
|                                       |
+---------------------------------------+
```

---

## 13. Support Chat

```
+---------------------------------------+
|            [status bar]               |
| <  Support                            |
|                                       |
| +-----------------------------------+ |
| |  [bot] TeslaPay Support  10:15   | |
| |  Hi Eva! How can I help you      | |
| |  today?                           | |
| |                                   | |
| |  Quick topics:                    | |
| |  [Card issue] [Payment help]     | |
| |  [Account] [Crypto] [Other]      | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| |                          [you]    | |
| |  I have a question about a       | |
| |  card payment I don't recognize.  | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| |  [bot] TeslaPay Support  10:16   | |
| |  I can help with that. Let me    | |
| |  connect you with a specialist.  | |
| |                                   | |
| |  Estimated wait: < 2 minutes     | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| |  [agent] Marta K.       10:17   | |
| |  Hi Eva, I'm Marta from the     | |
| |  support team. Let me look into  | |
| |  this for you. Can you tell me   | |
| |  which transaction?              | |
| +-----------------------------------+ |
|                                       |
| +-----------------------------------+ |
| | [Type a message...]     [send >] | |
| | [+] [photo]                      | |
| +-----------------------------------+ |
+---------------------------------------+
```

---

## 14. Dispute Flow

### 14.1 Select Reason

```
+---------------------------------------+
|            [status bar]               |
| X  Dispute Transaction                |
|                                       |
|    EUR 49.99 at UNKNOWN MERCHANT      |
|    2 March 2026                       |
|                                       |
|    Why are you disputing this?        |
|                                       |
| +-----------------------------------+ |
| | ( ) I don't recognize this        | |
| |     transaction                   | |
| +-----------------------------------+ |
| | ( ) I was charged the wrong       | |
| |     amount                        | |
| +-----------------------------------+ |
| | ( ) I was charged twice           | |
| +-----------------------------------+ |
| | ( ) I returned the item or        | |
| |     cancelled the service         | |
| +-----------------------------------+ |
| | ( ) I didn't receive the goods    | |
| |     or service                    | |
| +-----------------------------------+ |
| | ( ) Other reason                  | |
| +-----------------------------------+ |
|                                       |
|   +-------------------------------+   |
|   |          Continue             |   |
|   +-------------------------------+   |
|                                       |
+---------------------------------------+
```

### 14.2 Additional Details

```
+---------------------------------------+
|            [status bar]               |
| <  Dispute Details                    |
|                                       |
|    Did you have your card with        |
|    you at the time?                   |
|                                       |
|    ( ) Yes                            |
|    ( ) No                             |
|    (o) I'm not sure                   |
|                                       |
|   +-------------------------------+   |
|   | (!) We recommend freezing your|   |
|   |     card to prevent further   |   |
|   |     unauthorized charges.     |   |
|   |     [Freeze Card Now]         |   |
|   +-------------------------------+   |
|                                       |
|    Please describe what happened:     |
|    +-------------------------------+  |
|    | I don't recognize this charge |  |
|    | from this merchant. I have    |  |
|    | not purchased anything...     |  |
|    +-------------------------------+  |
|                                       |
|    Attach evidence (optional):        |
|    [+ Add Photo] [+ Add Document]     |
|                                       |
|   +-------------------------------+   |
|   |          Continue             |   |
|   +-------------------------------+   |
|                                       |
+---------------------------------------+
```

---

## Design Notes for Engineering

1. **All monetary values** use `type.body1` weight 600 (semibold) and `JetBrains Mono` for tabular alignment in lists. Display balances use `type.display1` or `type.display2` in `Inter`.

2. **Pull-to-refresh** is available on: Home dashboard, Transaction history, Crypto tab, Account detail.

3. **Skeleton loaders** replace all content areas during load. Match the layout shapes of the content being loaded.

4. **Empty states** always include: illustration, title, description, and a primary CTA button.

5. **Error states** always include: error icon (red), clear error message, a retry action, and a link to support.

6. **Offline banner** appears as a persistent top bar: "No internet connection. Some features may not be available." Background: `color.warning.500`, text: white.

7. **Transaction list grouping** is by date with sticky headers that collapse on scroll.

8. **Card visual** uses `CustomPainter` in Flutter for the gradient card with Mastercard and TeslaPay logos rendered as assets, not text.

9. **Numpad** is a custom Flutter widget (not system keyboard) for all financial input screens. This ensures consistent experience and prevents clipboard attacks on amount fields.

10. **Biometric prompts** use the platform-native API (`local_auth` package in Flutter) and must never show a custom biometric UI -- always delegate to the OS.
