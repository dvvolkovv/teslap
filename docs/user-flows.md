# TeslaPay User Flows

**Version:** 1.0
**Date:** 2026-03-03
**Author:** Senior UI/UX Designer, Dream Team

---

## Flow 1: Onboarding (Registration to First Deposit)

**Related Stories:** US-1.1, US-1.2, US-1.3, US-1.4, US-1.5, US-1.6
**Personas:** All (Sofia's scenario: opens account during lecture break in 3 minutes)

```
[App Launch]
    |
    v
[Splash Screen] ---- 2 seconds auto-advance
    |
    v
[Onboarding Carousel] ---- 3 slides, skip available
    |  Slide 1: Multi-currency illustration + "Banking without borders"
    |  Slide 2: Card illustration + "Your card, your rules"
    |  Slide 3: Crypto illustration + "Crypto made simple"
    |
    +-- [Get Started] button
    +-- [I already have an account] link --> Login Screen
    |
    v
[Registration: Email + Phone]
    |  - Enter email address
    |  - Enter phone number (with country code selector)
    |  - Accept Terms of Service + Privacy Policy (checkboxes + links)
    |  - [Continue] button
    |
    |-- VALIDATION: email format, phone format, not already registered
    |-- ERROR: "This email is already registered" --> [Log in instead] link
    |
    v
[Email Verification]
    |  - "We sent a 6-digit code to your@email.com"
    |  - OTP input (6 boxes, auto-focus next)
    |  - Auto-submit on 6th digit
    |  - [Resend code] link (enabled after 60s countdown)
    |  - [Change email] link --> back to previous
    |
    |-- ERROR: Wrong code --> shake animation + "Incorrect code. Try again."
    |-- ERROR: Expired --> "Code expired. Tap Resend."
    |-- 3 failed attempts --> "Too many attempts. Please wait 5 minutes."
    |
    v
[Phone Verification]
    |  - "We sent an SMS to +370 *** **89"
    |  - Same OTP pattern as email
    |  - Same error handling
    |
    v
[Create Password]
    |  - Password field (show/hide toggle)
    |  - Confirm password field
    |  - Real-time strength indicator:
    |    [ ] 8+ characters
    |    [ ] 1 uppercase letter
    |    [ ] 1 number
    |    [ ] 1 special character
    |  - All checks green --> [Continue] enabled
    |
    v
[Set App PIN]
    |  - "Create a 6-digit PIN for quick access"
    |  - Custom numpad with 6 dot indicators
    |  - On 6th digit --> [Confirm PIN] screen (re-enter)
    |  - Mismatch --> "PINs do not match. Try again." --> reset
    |
    v
[Enable Biometric Authentication]
    |  - Device-specific illustration (Face ID / Fingerprint)
    |  - "Use [Face ID] to log in faster?"
    |  - [Enable] button --> system biometric prompt
    |  - [Skip for now] link (can enable later in settings)
    |
    v
[KYC Introduction]
    |  - "Let's verify your identity"
    |  - "You'll need: a valid ID document, a well-lit space, 2 minutes"
    |  - Document type selection:
    |    ( ) Passport
    |    ( ) National ID Card
    |    ( ) Driver's License
    |    ( ) Residence Permit
    |  - [Start Verification] button
    |
    v
[Sumsub SDK: Document Capture]
    |  - Camera viewfinder with document outline overlay
    |  - Real-time feedback: "Hold steady", "Too blurry", "Move closer"
    |  - Auto-capture or manual shutter
    |  - Front captured --> review screen ("Looks good? / Retake")
    |  - Back captured (if two-sided document)
    |
    +-- [Optional: NFC Passport Read]  (US-1.4)
    |   |  - "Speed up verification by scanning your passport chip"
    |   |  - Illustration: hold phone to passport
    |   |  - NFC reading progress indicator
    |   |  - Success: "Chip read successfully" --> skip some steps
    |   |  - Failure: "Could not read chip. Continuing with photo verification."
    |
    v
[Sumsub SDK: Liveness Check]
    |  - Front camera activates
    |  - "Look at the camera and follow the instructions"
    |  - Liveness challenge (head turn, blink, etc.)
    |  - Processing indicator
    |
    |-- FAILURE: "Verification failed. Please try again in better lighting."
    |   Max 3 attempts --> "We need to review manually. You'll hear from us in 24h."
    |
    v
[KYC Processing Screen]
    |  - Animated progress indicator
    |  - "Verifying your identity..."
    |  - Typical wait: 15-60 seconds
    |
    +-- APPROVED (85%+ of cases)
    |   |
    |   v
    |   [Account Created!]
    |       - Confetti / check animation
    |       - "Welcome to TeslaPay, [Name]!"
    |       - Account tier displayed: "Basic Account"
    |       - Lithuanian IBAN shown (copy button)
    |       - Two CTAs:
    |         [Add Funds] --> Top Up flow
    |         [Get Your Card] --> Card request flow
    |       - [Explore the App] link --> Home Dashboard
    |
    +-- PENDING REVIEW (10% of cases)
    |   |
    |   v
    |   [Verification Under Review]
    |       - "We're reviewing your documents. This usually takes up to 24 hours."
    |       - "We'll notify you via push notification and email."
    |       - [Go to App] --> Home (limited: view-only, no transactions)
    |
    +-- REJECTED (5% of cases)
        |
        v
        [Verification Failed]
            - Clear reason: "Document expired", "Photo too blurry", "Name mismatch"
            - Attempt count: "Attempt 1 of 3"
            - [Try Again] --> Back to KYC Introduction
            - After 3 failures: [Contact Support] --> Chat
```

**Happy Path Duration Target:** Under 5 minutes (under 3 minutes for ePassport NFC users)

---

## Flow 2: Send Payment (SEPA + Internal)

**Related Stories:** US-3.1, US-3.2, US-3.3, US-6.3
**Personas:** Eva (sends EUR 200 to mother in Tallinn), Sofia (sends EUR 50 to roommate)

```
[Entry Points]
    +-- Home > Quick Actions > Send
    +-- Payments Tab > Send Money
    +-- Transaction Detail > Repeat Payment
    +-- Deep link: teslapay://payments/send
    |
    v
[Select Recipient]
    |  +-----------------------------------------+
    |  | [Search payees, IBAN, phone, name...]   |
    |  +-----------------------------------------+
    |  |                                         |
    |  | RECENT                                  |
    |  | (@) Anna K.  TeslaPay  LT12 3456...    |
    |  | (@) Mom      SEPA      EE12 7890...    |
    |  |                                         |
    |  | SAVED PAYEES                            |
    |  | (@) Landlord  DE89 3704...              |
    |  | (@) Electric  LT60 1010...              |
    |  |                                         |
    |  | [+ New Recipient]                       |
    |  | [Scan QR Code]                          |
    |  +-----------------------------------------+
    |
    +-- [Select saved/recent payee] --> Skip to Amount
    +-- [New Recipient] --> Enter details manually
    +-- [Scan QR] --> Camera --> Parse IBAN/amount --> pre-fill
    |
    v
[Enter Recipient Details] (if new)
    |  - Recipient name (text input)
    |  - IBAN (formatted input: LT12 3456 7890 1234 5678)
    |    Real-time validation: format check, checksum, country flag
    |  - OR Phone number (for TeslaPay internal lookup)
    |  - [Save this payee] toggle
    |  - [Continue]
    |
    |-- IBAN VALIDATION:
    |   +-- Valid IBAN, same TeslaPay BIN --> [INTERNAL TRANSFER detected]
    |   |   "This is a TeslaPay account. Transfer will be instant and free!"
    |   +-- Valid IBAN, external --> [SEPA TRANSFER]
    |   +-- Invalid IBAN --> "Please check the IBAN. It doesn't look right."
    |
    |-- PHONE LOOKUP:
    |   +-- Found TeslaPay user --> show name (masked), confirm
    |   +-- Not found --> "No TeslaPay account found. Try using their IBAN."
    |
    v
[Enter Amount]
    |  +-------------------------------------------+
    |  | From: EUR Account       EUR 3,245.67      |
    |  +-------------------------------------------+
    |  |                                           |
    |  |           EUR 200.00                      |
    |  |         [Custom numpad]                   |
    |  |                                           |
    |  +-------------------------------------------+
    |  | Reference: [optional text field]          |
    |  +-------------------------------------------+
    |  | Transfer speed:                           |
    |  |   (o) Standard SEPA -- 1 business day     |
    |  |   ( ) Instant SEPA -- 10 seconds   +EUR 0 |
    |  +-------------------------------------------+
    |  | Fee: EUR 0.00                             |
    |  | Recipient gets: EUR 200.00                |
    |  +-------------------------------------------+
    |  | [Continue]                                |
    |  +-------------------------------------------+
    |
    |-- INTERNAL TRANSFER: speed selection hidden (always instant, always free)
    |-- AMOUNT > BALANCE: [Continue] disabled, "Insufficient funds" in red
    |-- AMOUNT > TIER LIMIT: "This exceeds your daily limit of EUR X,XXX"
    |
    v
[Review and Confirm]
    |  +-------------------------------------------+
    |  |         Review Your Transfer              |
    |  +-------------------------------------------+
    |  | From        EUR Account (LT12...)         |
    |  | To          Anna Kowalski                 |
    |  |             EE12 7890 1234 5678 90        |
    |  | Amount      EUR 200.00                    |
    |  | Fee         EUR 0.00                      |
    |  | Total       EUR 200.00                    |
    |  | Delivery    Instant (SEPA Inst)            |
    |  | Reference   March rent                    |
    |  +-------------------------------------------+
    |  |                                           |
    |  |     [Confirm with Face ID]                |
    |  |                                           |
    |  +-------------------------------------------+
    |
    |-- Biometric auth prompt (Face ID / fingerprint)
    |   +-- SUCCESS --> Process payment
    |   +-- FAIL x3 --> Fall back to PIN
    |   +-- PIN FAIL x5 --> Account locked 30 min
    |
    v
[Processing]
    |  - Brief loading state (< 2s internal, < 10s SEPA Inst)
    |
    +-- SUCCESS
    |   |
    |   v
    |   [Transfer Successful!]
    |       - Green check animation
    |       - "EUR 200.00 sent to Anna Kowalski"
    |       - [Share Receipt] button
    |       - [Send Another] button
    |       - [Back to Home] button
    |       - Push notification sent to recipient (if TeslaPay user)
    |
    +-- FAILURE
        |
        v
        [Transfer Failed]
            - Error reason: "Recipient bank declined", "Network error"
            - [Try Again] button
            - [Contact Support] link
            - Funds not debited (or auto-reversed)
```

---

## Flow 3: Card Ordering and Activation

**Related Stories:** US-4.1, US-4.2, US-4.3, US-4.7, US-4.8
**Personas:** Sofia (instant virtual + Apple Pay), Eva (physical card for travel)

### 3A: Virtual Card (Instant)

```
[Entry Point: Card Tab -- "Get Your Card" OR Onboarding prompt]
    |
    v
[Card Introduction]
    |  - Card visual preview (TeslaPay Mastercard design)
    |  - "Get your free virtual Mastercard instantly"
    |  - Features: Online payments, Apple/Google Pay, contactless
    |  - [Get Virtual Card] primary button
    |  - [Order Physical Card] secondary link
    |
    v
[Issuing Card...]
    |  - Card flip animation (5-15 seconds)
    |  - "Creating your Mastercard..."
    |
    v
[Virtual Card Ready!]
    |  +-------------------------------------------+
    |  |  [TeslaPay Mastercard card visual]        |
    |  |  **** **** **** 4521                      |
    |  |  SOFIA PETROV         12/29               |
    |  +-------------------------------------------+
    |  |                                           |
    |  | Your virtual card is ready!               |
    |  |                                           |
    |  | [Show Card Details] -- biometric gated    |
    |  | [Add to Apple Pay]                        |
    |  | [Add to Google Pay]                       |
    |  | [Order Physical Card]                     |
    |  | [Go to Card]                              |
    |  +-------------------------------------------+
    |
    +-- [Add to Apple Pay] (US-4.7)
    |   |  - System Apple Pay provisioning sheet
    |   |  - Terms acceptance
    |   |  - Card verification (SMS OTP or in-app)
    |   |  - "Card added to Apple Wallet" confirmation
    |   |  - Duration: < 60 seconds
    |
    +-- [Add to Google Pay] (US-4.8)
        |  - Google Pay provisioning flow
        |  - Same pattern as Apple Pay
```

### 3B: Physical Card

```
[Entry: Card Tab > Order Physical Card]
    |
    v
[Delivery Address]
    |  - Pre-filled from profile registration address
    |  - [Edit Address] option
    |  - Address fields: street, city, postal code, country
    |  - Country restricted to EEA
    |  - [Continue]
    |
    v
[Card Design Selection] (if multiple options available)
    |  - Standard (black, free)
    |  - Premium (gradient, EUR 10)  -- future
    |  - Swipeable card previews
    |  - [Continue]
    |
    v
[Order Summary]
    |  - Card type: TeslaPay Mastercard Debit
    |  - Delivery to: [address]
    |  - Estimated delivery: 5-10 business days
    |  - Cost: Free (first card) / EUR 10 (replacement)
    |  - [Confirm Order]
    |
    v
[Order Confirmed]
    |  - "Your card is on its way!"
    |  - Estimated arrival date
    |  - "We'll send you a notification when it ships."
    |  - [Track Order] link
    |  - [Back to Card] button
    |
    ... [Days pass, tracking notifications sent] ...
    |
    v
[Push Notification: "Your card has arrived! Activate it now."]
    |  - Tap --> Deep link to activation
    |
    v
[Activate Physical Card]
    |  +-- Option A: NFC tap (hold card to phone)
    |  +-- Option B: Enter last 4 digits of card number
    |  |
    |  v
    |  [Enter Last 4 Digits]
    |  |  - 4 input boxes
    |  |  - [Activate]
    |  |
    |  v
    |  [Activating...]
    |  |  - 3-5 seconds
    |  |
    |  v
    |  [Card Activated!]
    |      - "Your physical card is ready to use."
    |      - "Your PIN for ATM and POS: [View PIN]" (biometric gated)
    |      - [Set ATM Limits]
    |      - [Done]
```

---

## Flow 4: Crypto Purchase and DeFi Staking

**Related Stories:** US-5.1, US-5.2, US-5.3, US-5.8
**Personas:** Marco (buys USDC, deposits to Solid), Sofia (buys EUR 20 FUSE out of curiosity)

### 4A: First Crypto Purchase

```
[Entry: Crypto Tab (first visit after KYC)]
    |
    v
[Crypto Wallet Introduction]
    |  - "Your Fuse blockchain wallet is ready"
    |  - Brief explanation: "Buy, sell, send crypto -- all from your TeslaPay account"
    |  - Educational carousel (3 micro-slides):
    |    1. "What is Fuse?" -- fast, low-cost blockchain
    |    2. "Stablecoins" -- USDC/USDT pegged to USD
    |    3. "Self-custody" -- your keys, your coins
    |  - [Start Exploring] button
    |  - Risk disclosure link at bottom
    |
    v
[Crypto Dashboard] (empty state)
    |  - Total balance: EUR 0.00
    |  - Token list:
    |    FUSE    0.00     EUR 0.00
    |    USDC    0.00     EUR 0.00
    |    USDT    0.00     EUR 0.00
    |  - [Buy Crypto] prominent CTA
    |  - Wallet address shown (collapsed, expandable)
    |
    v
[Buy Crypto]
    |
    v
[Select Token]
    |  +-------------------------------------------+
    |  | What would you like to buy?               |
    |  +-------------------------------------------+
    |  | (FUSE)  FUSE       EUR 0.042  +3.2%      |
    |  | (USDC)  USD Coin   EUR 0.92   +0.01%     |
    |  | (USDT)  Tether     EUR 0.92   -0.02%     |
    |  +-------------------------------------------+
    |
    v
[Enter Purchase Amount]
    |  +-------------------------------------------+
    |  | Buy FUSE                                  |
    |  +-------------------------------------------+
    |  | Pay with: EUR Account  (EUR 3,245.67)     |
    |  +-------------------------------------------+
    |  |                                           |
    |  |          EUR 20.00                        |
    |  |     ~ 476.19 FUSE                        |
    |  |                                           |
    |  | [EUR 5] [EUR 20] [EUR 50] [EUR 100]      |  <-- quick amount chips
    |  |                                           |
    |  +-------------------------------------------+
    |  | Rate:  1 FUSE = EUR 0.042                 |
    |  | Fee:   EUR 0.30 (1.5%)                    |
    |  | You get: ~469.05 FUSE                     |
    |  +-------------------------------------------+
    |  | [Continue]                                |
    |  +-------------------------------------------+
    |
    |-- FIRST TIME: Risk disclosure bottom sheet before Continue
    |   "Crypto assets are volatile. You may lose some or all of your investment.
    |    TeslaPay does not provide investment advice."
    |   [I understand the risks] checkbox + [Continue]
    |
    |-- AMOUNT < EUR 5: "Minimum purchase is EUR 5"
    |-- AMOUNT > BALANCE: "Insufficient EUR balance"
    |
    v
[Review Purchase]
    |  +-------------------------------------------+
    |  |         Review Your Purchase              |
    |  +-------------------------------------------+
    |  | Buy           469.05 FUSE                 |
    |  | Pay           EUR 20.00                   |
    |  | Fee           EUR 0.30                    |
    |  | Total debit   EUR 20.30                   |
    |  | Rate          1 FUSE = EUR 0.042          |
    |  | Rate valid for: 28s (countdown)           |
    |  +-------------------------------------------+
    |  |                                           |
    |  |     [Confirm with Face ID]                |
    |  |                                           |
    |  +-------------------------------------------+
    |
    |-- Rate expired (30s) --> "Rate expired. [Get new rate]"
    |-- Biometric confirm --> process
    |
    v
[Processing Purchase]
    |  - "Buying FUSE..."
    |  - Duration: < 60 seconds
    |
    +-- SUCCESS
    |   |
    |   v
    |   [Purchase Complete!]
    |       - Token logo + animated balance counter
    |       - "You now own 469.05 FUSE"
    |       - "Worth approximately EUR 19.70"
    |       - [View in Portfolio]
    |       - [Buy More]
    |       - [Back to Crypto]
    |
    +-- FAILURE
        |
        v
        [Purchase Failed]
            - Error reason
            - [Try Again] / [Contact Support]
            - EUR not debited
```

### 4B: DeFi Staking (Phase 2 -- Solid soUSD)

```
[Entry: Crypto Tab > DeFi Yield section > "Earn" CTA]
    |
    v
[Earn with Solid soUSD]
    |  +-------------------------------------------+
    |  | Earn yield on your stablecoins            |
    |  +-------------------------------------------+
    |  | Current APY:  4.8%                        |
    |  |                                           |
    |  | How it works:                             |
    |  | 1. Deposit USDC or USDT                   |
    |  | 2. Receive soUSD (yield-bearing)           |
    |  | 3. Withdraw anytime -- no lock-up          |
    |  |                                           |
    |  | [Deposit Now]                             |
    |  | [Learn More] --> educational bottom sheet  |
    |  +-------------------------------------------+
    |
    |-- FIRST TIME: Risk disclosure modal
    |   "DeFi yield involves smart contract risk. Funds are not protected
    |    by deposit guarantee schemes. Past returns do not guarantee future results."
    |   [I understand and accept the risks] --> Continue
    |
    v
[Deposit to Solid]
    |  - Select token: USDC (balance: 500.00) / USDT (balance: 0.00)
    |  - Enter amount (or "Max")
    |  - Fee estimate
    |  - [Review Deposit]
    |
    v
[Review Deposit]
    |  - Deposit: 500.00 USDC
    |  - You receive: ~500.00 soUSD
    |  - Estimated annual yield: ~24.00 USDC (at 4.8% APY)
    |  - Network fee: 0.00 (gasless)
    |  - [Confirm with Face ID]
    |
    v
[Deposit Processing]
    |  - On-chain transaction
    |  - Duration: 10-30 seconds on Fuse
    |
    v
[Deposit Successful!]
    |  - "You're now earning yield!"
    |  - soUSD balance: 500.00
    |  - Estimated daily yield: ~0.066 USDC
    |  - [View Earnings]
    |  - [Back to Crypto]
    |
    v
[Ongoing: Yield Dashboard]
    |  +-------------------------------------------+
    |  | Your Yield                                |
    |  +-------------------------------------------+
    |  | Deposited:     500.00 soUSD               |
    |  | Current value: 501.34 USDC equivalent     |
    |  | Earned:        1.34 USDC                  |
    |  | APY:           4.8%                       |
    |  | Duration:      20 days                    |
    |  +-------------------------------------------+
    |  | [Deposit More] [Withdraw]                 |
    |  +-------------------------------------------+
    |
    +-- [Withdraw]
        |  - Enter amount or "Withdraw All"
        |  - Review: soUSD returned, USDC received
        |  - Confirm with biometric
        |  - Processing
        |  - USDC returned to Fuse wallet balance
```

---

## Flow 5: Biometric Login

**Related Stories:** US-6.1, US-6.2
**Personas:** All (Eva checking spending on the train -- must be fast)

```
[App Opened (returning user)]
    |
    +-- App was backgrounded < 5 minutes
    |   |  --> Resume without auth (show last screen)
    |
    +-- App was backgrounded > 5 minutes OR cold start
        |
        v
        [Login Screen]
            |  +-------------------------------------------+
            |  |         [TeslaPay Logo]                   |
            |  |                                           |
            |  |      Welcome back, Eva                    |
            |  |                                           |
            |  |          (Face ID icon)                   |
            |  |      Tap to unlock                        |
            |  |                                           |
            |  |    [Use PIN instead]                      |
            |  |    [Use password]                         |
            |  |    [Not Eva? Switch account]              |
            |  +-------------------------------------------+
            |
            +-- AUTO: System biometric prompt fires automatically on screen load
            |
            +-- BIOMETRIC SUCCESS
            |   |  - Haptic feedback (success)
            |   |  - Fade transition to Home Dashboard
            |   |  - Duration: < 2 seconds total
            |   v
            |   [Home Dashboard]
            |
            +-- BIOMETRIC FAIL (attempt 1-2)
            |   |  - "Face not recognized. Try again."
            |   |  - [Try Again] --> re-prompt biometric
            |
            +-- BIOMETRIC FAIL (attempt 3)
            |   |  - "Biometric login disabled. Use your PIN."
            |   v
            |   [PIN Entry Screen]
            |       |  - 6 dot indicators
            |       |  - Custom numpad
            |       |  - [Forgot PIN?] link
            |       |
            |       +-- PIN CORRECT --> Home Dashboard
            |       |
            |       +-- PIN WRONG (attempt 1-4)
            |       |   - Dots shake + turn red
            |       |   - "Incorrect PIN. [X] attempts remaining."
            |       |
            |       +-- PIN WRONG (attempt 5)
            |           |  - "Account locked for 30 minutes."
            |           |  - [Use Password] link
            |           v
            |           [Password Entry]
            |               |  - Email (pre-filled, read-only)
            |               |  - Password field
            |               |  - [Log In]
            |               |
            |               +-- CORRECT --> Home + prompt to reset PIN
            |               +-- WRONG x3 --> "Account locked. Contact support."
            |               +-- [Forgot Password] --> Reset flow
            |
            +-- BIOMETRIC NOT ENROLLED
            |   |  - Skip biometric, show PIN entry directly
            |
            +-- BIOMETRIC CHANGED (new face/finger enrolled on device)
                |  - "Device biometrics have changed. Please verify with your PIN."
                |  - PIN entry required
                |  - After PIN success: "Re-enable biometric login?"
                |  - [Enable] / [Not now]
```

**New Device Login:**
```
[Login on New Device]
    |
    v
[Email + Password Entry]
    |
    v
[2FA Challenge]
    |  - Push notification to existing device: "New login from [Device], [Location]"
    |  - OR SMS OTP to registered phone
    |  - User approves on existing device OR enters OTP
    |
    v
[Set PIN on New Device]
    |
    v
[Enable Biometric on New Device]
    |
    v
[Home Dashboard]
    |
    +-- Old device receives push: "New login detected on [Device Name]"
        "If this wasn't you, [Secure Account] link"
```

---

## Flow 6: Card Dispute

**Related Stories:** US-7.3, US-4.9
**Personas:** Eva (unrecognized charge while traveling)

```
[Entry Points]
    +-- Transaction Detail > [Dispute] button
    +-- Card Tab > Recent Transactions > tap transaction > [Dispute]
    +-- Push notification ("Charge of EUR 49.99 at Unknown Merchant") > [Not you?] action
    |
    v
[Dispute: Select Transaction] (if not entered from specific transaction)
    |  - List of card transactions from last 120 days
    |  - User selects the disputed transaction
    |
    v
[Dispute: Select Reason]
    |  +-------------------------------------------+
    |  | Why are you disputing this transaction?   |
    |  +-------------------------------------------+
    |  | ( ) I don't recognize this transaction    |
    |  | ( ) I was charged the wrong amount        |
    |  | ( ) I was charged twice                   |
    |  | ( ) I returned the item / cancelled       |
    |  | ( ) I didn't receive the goods/service    |
    |  | ( ) Other                                 |
    |  +-------------------------------------------+
    |  | [Continue]                                |
    |  +-------------------------------------------+
    |
    v
[Dispute: Additional Details]
    |  +-------------------------------------------+
    |  | Help us understand what happened          |
    |  +-------------------------------------------+
    |  | Transaction: EUR 49.99 at UNKNOWN MERCHANT|
    |  | Date: 2 March 2026                        |
    |  +-------------------------------------------+
    |  | Description: [text area, optional]        |
    |  |                                           |
    |  | "Please describe what happened..."        |
    |  +-------------------------------------------+
    |  | Attach evidence (optional)                |
    |  | [+ Add Photo] [+ Add Document]            |
    |  +-------------------------------------------+
    |  | [Continue]                                |
    |  +-------------------------------------------+
    |
    |-- If reason is "I don't recognize":
    |   Extra question: "Did you have your card with you at the time?"
    |   ( ) Yes  ( ) No  ( ) I'm not sure
    |
    |-- If "No" or "Not sure":
    |   Immediate prompt: "We recommend freezing your card to prevent further unauthorized use."
    |   [Freeze Card Now] [Keep Active]
    |
    v
[Dispute: Review and Submit]
    |  +-------------------------------------------+
    |  |        Review Your Dispute                |
    |  +-------------------------------------------+
    |  | Transaction  EUR 49.99, UNKNOWN MERCHANT  |
    |  | Date         2 March 2026                 |
    |  | Reason       Unrecognized transaction     |
    |  | Card frozen  Yes                          |
    |  | Evidence     1 photo attached             |
    |  +-------------------------------------------+
    |  | By submitting, you confirm the information|
    |  | provided is accurate and complete.        |
    |  +-------------------------------------------+
    |  | [Submit Dispute]                          |
    |  +-------------------------------------------+
    |
    v
[Dispute Submitted]
    |  - "Your dispute has been submitted."
    |  - Reference number: DSP-2026030-A1B2
    |  - "We'll investigate and update you within 15 business days."
    |  - "You'll receive updates via push notification and email."
    |  - [View Dispute Status]
    |  - [Back to Card]
    |
    v
[Ongoing: Dispute Status] (accessible from Transaction Detail or Card > Disputes)
    |
    |  Status timeline:
    |  [x] Submitted -- 2 Mar 2026
    |  [x] Under investigation -- 3 Mar 2026
    |  [ ] Resolution pending
    |  [ ] Resolved
    |
    +-- RESOLVED: Approved
    |   - "Dispute resolved in your favor. EUR 49.99 refunded."
    |   - Refund visible in transaction history
    |   - Push notification + email sent
    |
    +-- RESOLVED: Denied
    |   - "After investigation, the charge appears valid."
    |   - Detailed reason provided
    |   - [Appeal] option --> Contact support
    |
    +-- RESOLVED: Partial
        - "Partial refund of EUR [X] issued."
        - Explanation provided
```

---

## Flow 7: Currency Exchange

**Related Stories:** US-3.4
**Personas:** Eva (converts EUR to PLN for Warsaw trip)

```
[Entry: Home Quick Actions > Exchange OR Payments Tab > Exchange]
    |
    v
[Exchange Screen]
    |  +-------------------------------------------+
    |  | Exchange                                  |
    |  +-------------------------------------------+
    |  | From                                      |
    |  | [EUR flag] EUR    EUR 3,245.67            |
    |  |                                           |
    |  |            EUR 500.00                     |
    |  +-------------------------------------------+
    |  |           [  Swap  ]                      |
    |  +-------------------------------------------+
    |  | To                                        |
    |  | [PLN flag] PLN    PLN 0.00                |
    |  |                                           |
    |  |          ~ PLN 2,147.50                   |
    |  +-------------------------------------------+
    |  |                                           |
    |  | Rate: 1 EUR = 4.2950 PLN                  |
    |  | Mid-market: 4.2850                        |
    |  | TeslaPay markup: 0.23%                    |
    |  | Rate updates in: 28s                      |
    |  |                                           |
    |  | [Exchange]                                |
    |  +-------------------------------------------+
    |
    |-- Tap currency selector --> bottom sheet with currency list
    |   (only currencies user has opened sub-accounts for,
    |    plus "Open [currency] account" option)
    |
    |-- Rate refreshes every 30 seconds
    |-- If rate expires during review: "Rate updated" notification inline
    |
    v
[Confirm Exchange]
    |  - Summary: EUR 500.00 --> PLN 2,147.50
    |  - Rate locked for 30 seconds
    |  - Biometric / PIN confirmation
    |
    v
[Exchange Complete]
    |  - "Exchanged EUR 500.00 to PLN 2,147.50"
    |  - Both balances updated
    |  - [Done]
```

---

## Flow 8: Request Money

**Related Stories:** US-3.3 (related)
**Personas:** Sofia (requesting EUR 50 from roommate for bills)

```
[Entry: Home Quick Actions > Request OR Payments Tab > Request]
    |
    v
[Request Money]
    |  +-------------------------------------------+
    |  | Request Money                             |
    |  +-------------------------------------------+
    |  | From: [Search contact / enter phone]      |
    |  +-------------------------------------------+
    |  |                                           |
    |  |          EUR 50.00                        |
    |  |                                           |
    |  | Note: "March utilities"                   |
    |  +-------------------------------------------+
    |  | [Send Request]                            |
    |  +-------------------------------------------+
    |
    v
[Request Sent]
    |  - Push notification sent to recipient (if TeslaPay user)
    |  - Share link generated (for non-TeslaPay users)
    |  - [Share via...] (WhatsApp, Telegram, SMS, etc.)
    |  - Request appears in sender's activity as "Pending"
    |
    +-- Recipient receives notification:
        "Sofia requested EUR 50.00 for 'March utilities'"
        [Pay Now] [Decline]
        |
        +-- [Pay Now] --> Pre-filled Send Money flow
        +-- [Decline] --> Notification to requester
```
