# User Stories: TeslaPay Neobank

**Version:** 1.0
**Date:** 2026-03-03

Story points: S (1-2 days), M (3-5 days), L (1-2 weeks), XL (2-4 weeks)
Priority: Must (MVP), Should (MVP if time), Could (Phase 2), Won't (Out of scope)

---

## Epic 1: User Registration and Onboarding

### US-1.1: Account Registration
**As a** new user,
**I want to** register for a TeslaPay account using my email and phone number,
**so that** I can start the onboarding process.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User provides email address, phone number, and creates a password
- Email verified via OTP link (expires in 15 minutes)
- Phone verified via SMS OTP (expires in 5 minutes)
- Password meets minimum requirements: 8+ characters, 1 uppercase, 1 number, 1 special character
- Duplicate email/phone rejected with clear message
- Registration available in 5+ languages

### US-1.2: KYC Document Verification
**As a** registered user,
**I want to** verify my identity by submitting an ID document via Sumsub,
**so that** I can activate my account and access banking features.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- User can select document type: passport, national ID card, driver's license, residence permit
- Camera-based document capture with real-time quality feedback (blur, glare, cropping)
- Front and back capture for two-sided documents
- Document submitted to Sumsub API; result returned within 60 seconds for 90%+ of cases
- Approved users proceed to account activation
- Rejected users see clear reason and can retry (max 3 attempts before manual review)

### US-1.3: Liveness Check
**As a** registered user,
**I want to** complete a liveness check during onboarding,
**so that** TeslaPay can confirm I am a real person and prevent identity fraud.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Liveness check initiated after document submission
- Sumsub liveness SDK embedded in native app
- Detects and rejects photos, videos, masks, and deepfakes
- Completes in under 30 seconds for successful attempts
- Failure prompts retry with guidance; max 3 attempts

### US-1.4: NFC Document Reading
**As a** user with an ePassport,
**I want to** verify my identity by scanning the NFC chip in my passport,
**so that** I get faster approval with higher trust level.

**Priority:** Should | **Size:** M
**Acceptance Criteria:**
- NFC read option presented for passport holders
- App reads chip data (MRZ, photo, fingerprint hash)
- Data cross-referenced with submitted document for consistency
- Successful NFC read upgrades trust level; may skip additional verification steps
- Graceful fallback if NFC read fails (continue with photo verification)

### US-1.5: AML Screening at Onboarding
**As a** compliance officer,
**I want** every new user automatically screened against sanctions, PEP, and adverse media databases,
**so that** TeslaPay does not onboard prohibited persons.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Sumsub AML screening runs automatically after identity verification
- Checks: global sanctions lists, PEP databases, adverse media
- Clean results: user proceeds to account activation automatically
- Match found: user placed in manual review queue; not activated until cleared
- Screening result and decision logged in audit trail

### US-1.6: Account Tier Assignment
**As a** newly verified user,
**I want to** be assigned an account tier (Basic, Standard, or Premium),
**so that** I understand my transaction limits and available features.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- Basic tier assigned by default after KYC approval
- Tier details displayed: transaction limits, card eligibility, FX allowance, fees
- User can view upgrade path and requirements for higher tiers
- Upgrade requires additional verification (proof of address, source of funds)

### US-1.7: Data Migration for Existing Users
**As an** existing TeslaPay customer,
**I want to** migrate my account to the new platform without losing my data or funds,
**so that** I have a seamless transition experience.

**Priority:** Must | **Size:** XL
**Acceptance Criteria:**
- Existing IBAN preserved (if technically feasible) or new IBAN with redirect
- Balance transferred accurately to new ledger
- Transaction history from legacy system viewable
- Existing KYC data re-used where possible (no re-verification for compliant users)
- Migration can be done in batches; rollback plan available

---

## Epic 2: Account Management

### US-2.1: View Account Dashboard
**As a** user,
**I want to** see my account balances, recent transactions, and quick actions on the home screen,
**so that** I have an immediate overview of my finances.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Dashboard shows total balance across all currencies (EUR equivalent)
- Individual currency balances displayed with currency flag/icon
- Last 5 transactions shown with merchant/payee name, amount, and date
- Quick action buttons: Send, Request, Card, Exchange
- Dashboard loads in under 3 seconds on 4G connection

### US-2.2: Multi-Currency Sub-Accounts
**As a** user,
**I want to** open sub-accounts in different currencies (EUR, USD, GBP, PLN, CHF),
**so that** I can hold and manage multiple currency balances.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User can open a new currency sub-account from the app in 2 taps
- Each sub-account has its own balance and transaction history
- EUR sub-account created by default at onboarding
- Sub-account shows dedicated IBAN where applicable
- Sub-accounts can be closed if balance is zero

### US-2.3: View Transaction History
**As a** user,
**I want to** view my complete transaction history with search and filter options,
**so that** I can track my spending and find specific transactions.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- All transactions listed in reverse chronological order
- Filter by: date range, amount range, transaction type, currency, status
- Search by: merchant name, payee name, reference, amount
- Transaction detail shows: amount, currency, FX rate (if applicable), fee, timestamp, status, category
- Export to CSV and PDF
- Pagination with infinite scroll; loads 50 transactions per page

### US-2.4: Update Personal Information
**As a** user,
**I want to** update my personal details (address, phone, email),
**so that** my account information stays current.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User can update: email, phone, address, tax residency
- Email/phone changes require OTP verification on both old and new
- Address changes require proof of address document upload (verified via Sumsub)
- Name changes require new ID document verification
- All changes logged in audit trail with before/after values

### US-2.5: Account Closure
**As a** user,
**I want to** close my TeslaPay account,
**so that** I can leave the service and receive my remaining funds.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User initiates closure from Settings
- System checks for pending transactions, active cards, active mandates
- User must resolve blockers before closure proceeds
- Remaining balance transferred to user-specified external IBAN
- Account deactivated after fund transfer confirmed
- Data retained per regulatory requirements (5 years minimum)

### US-2.6: Notification Preferences
**As a** user,
**I want to** configure which notifications I receive and through which channels,
**so that** I only get alerts that matter to me.

**Priority:** Should | **Size:** S
**Acceptance Criteria:**
- Configurable categories: transactions, security, marketing, product updates
- Channels: push notifications, email, SMS
- Transaction notifications enabled by default (cannot be disabled for security)
- Marketing notifications opt-in only (GDPR compliant)
- Changes take effect immediately

---

## Epic 3: Payments and Transfers

### US-3.1: Send SEPA Credit Transfer
**As a** user,
**I want to** send a SEPA Credit Transfer to any IBAN in the SEPA zone,
**so that** I can pay bills and transfer money to other people.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- User enters: recipient IBAN, recipient name, amount, reference (optional)
- IBAN validated (format + checksum) before submission
- Fee and estimated delivery time shown before confirmation
- Confirmation requires biometric auth or PIN
- Transaction visible in history immediately with "Processing" status
- Status updated to "Completed" when settlement confirmed
- Settlement within 1 business day

### US-3.2: Send SEPA Instant Transfer
**As a** user,
**I want to** send money instantly to another SEPA account,
**so that** the recipient receives funds within seconds.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- User can toggle "Instant" when sending SEPA transfer
- Instant fee (if any) displayed clearly before confirmation
- Transfer settles within 10 seconds
- Recipient bank must support SCT Inst; fallback to regular SCT if not supported
- Available 24/7/365
- Confirmation with settlement timestamp shown to user

### US-3.3: Internal Transfer (TeslaPay to TeslaPay)
**As a** user,
**I want to** send money instantly to another TeslaPay user,
**so that** I can split bills and pay friends without waiting.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User can send by: IBAN, phone number (if recipient has linked phone), or username
- Transfer is instant (under 2 seconds)
- No fee for internal transfers
- Both sender and recipient see the transaction immediately
- Optional message/note attached to transfer

### US-3.4: Currency Exchange
**As a** user,
**I want to** exchange between currencies in my account,
**so that** I can manage my multi-currency balances.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- User selects source currency, target currency, and amount
- Live exchange rate displayed with TeslaPay markup (max 0.5%)
- Rate locked for 30 seconds after display
- Confirmation shows: amount in, rate, fee/markup, amount out
- Exchange executes instantly after confirmation
- Transaction recorded in both currency sub-accounts

### US-3.5: Receive SEPA Payments
**As a** user,
**I want to** receive SEPA payments into my TeslaPay IBAN,
**so that** I can receive salary, payments, and transfers.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Each user has a dedicated Lithuanian IBAN
- Incoming SEPA payments credited within standard settlement time
- Incoming SEPA Instant payments credited within 10 seconds
- Push notification sent upon receipt
- Transaction shows sender name, IBAN, reference, amount

### US-3.6: Schedule Recurring Payment
**As a** user,
**I want to** set up recurring payments (weekly, monthly, custom),
**so that** I can automate regular bills and transfers.

**Priority:** Should | **Size:** M
**Acceptance Criteria:**
- User sets: payee, amount, frequency (weekly/bi-weekly/monthly/custom), start date, end date (optional)
- Recurring payment executes automatically on schedule
- Push notification sent before and after each execution
- User can pause, resume, modify, or cancel recurring payment
- Failed payment (insufficient funds) triggers notification with retry option

### US-3.7: Manage Saved Payees
**As a** user,
**I want to** save frequently used payees,
**so that** I can send money quickly without re-entering details.

**Priority:** Should | **Size:** S
**Acceptance Criteria:**
- User can save payee with: name, IBAN, default reference
- Saved payees shown in transfer flow for quick selection
- User can edit or delete saved payees
- IBAN validated at save time

### US-3.8: SEPA Direct Debit Management
**As a** user,
**I want to** view and manage SEPA Direct Debit mandates on my account,
**so that** I can control which companies can debit my account.

**Priority:** Should | **Size:** M
**Acceptance Criteria:**
- User can view all active SDD mandates
- User can cancel/revoke a mandate
- Upcoming direct debits shown with expected date and amount
- User can dispute a direct debit within 8 weeks (SEPA rules)

---

## Epic 4: Mastercard Card Management

### US-4.1: Request Virtual Card
**As a** user,
**I want to** request a virtual Mastercard debit card instantly,
**so that** I can start making online and contactless purchases immediately.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Virtual card issued within 30 seconds of request
- Card number, expiry, and CVV displayed in-app (behind biometric auth)
- Card linked to user's EUR sub-account by default
- Card can be used for online purchases and added to mobile wallets immediately

### US-4.2: Request Physical Card
**As a** user,
**I want to** order a physical Mastercard debit card,
**so that** I can make in-store purchases and ATM withdrawals.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User selects delivery address (default: registered address)
- Delivery estimate displayed (5-10 business days within EU)
- Card shipped within 3 business days of order
- Tracking number provided via push notification
- Card requires activation in-app before first use
- Physical card shares the same PAN as virtual card, or separate PAN per user preference

### US-4.3: Activate Physical Card
**As a** user who received a physical card,
**I want to** activate it through the app,
**so that** I can start using it at POS terminals and ATMs.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- User scans card or enters last 4 digits to confirm receipt
- Card activated within 5 seconds
- Activation confirmation shown in app and via push notification

### US-4.4: Freeze and Unfreeze Card
**As a** user,
**I want to** freeze my card temporarily,
**so that** I can protect my card if I suspect misuse or cannot find it.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- Single toggle in app to freeze/unfreeze
- Freeze takes effect within 5 seconds
- All new authorizations declined while frozen
- Existing recurring payments may still process (configurable)
- Unfreeze restores full functionality immediately

### US-4.5: Set Spending Controls
**As a** user,
**I want to** set spending limits and merchant category restrictions on my card,
**so that** I can control my spending and reduce fraud risk.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User can set: per-transaction limit, daily limit, monthly limit
- User can block merchant categories (e.g., gambling, adult content)
- User can restrict geographic regions (e.g., only EEA, only specific countries)
- User can enable/disable: online payments, contactless, ATM, magnetic stripe
- Changes take effect within 10 seconds

### US-4.6: View and Change PIN
**As a** user,
**I want to** view my card PIN in the app and change it if needed,
**so that** I can use my card at ATMs and POS terminals that require PIN.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- View PIN: requires biometric authentication; PIN shown for 10 seconds then hidden
- Change PIN: user sets new 4-digit PIN; takes effect within 30 seconds
- PIN change confirmed via push notification
- Maximum 3 failed PIN attempts at ATM/POS before card blocked (standard Mastercard rules)

### US-4.7: Add Card to Apple Pay
**As an** iOS user,
**I want to** add my TeslaPay card to Apple Pay,
**so that** I can make contactless payments with my iPhone and Apple Watch.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- "Add to Apple Wallet" button visible on card screen
- Provisioning flow follows Apple Pay guidelines
- Card tokenized via Mastercard MDES
- Tokenized card appears in Apple Wallet within 60 seconds
- Transactions made via Apple Pay appear in TeslaPay app with Apple Pay indicator

### US-4.8: Add Card to Google Pay
**As an** Android user,
**I want to** add my TeslaPay card to Google Pay,
**so that** I can make contactless payments with my Android phone.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- "Add to Google Wallet" button visible on card screen
- Provisioning flow follows Google Pay guidelines
- Card tokenized via Mastercard MDES
- Tokenized card appears in Google Wallet within 60 seconds
- Transactions made via Google Pay appear in TeslaPay app with Google Pay indicator

### US-4.9: Receive Real-Time Transaction Notifications
**As a** card holder,
**I want to** receive instant push notifications for every card transaction,
**so that** I can monitor spending and detect unauthorized use immediately.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- Push notification sent within 3 seconds of authorization
- Notification includes: merchant name, amount, currency, location (if available)
- Declined transactions also trigger notification with decline reason
- Tapping notification opens transaction detail in app

### US-4.10: Report Card Lost or Stolen
**As a** user,
**I want to** report my card as lost or stolen and request a replacement,
**so that** I can prevent fraudulent use and get a new card.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- "Report Lost/Stolen" option in card settings
- Card blocked immediately upon report
- Reason captured: lost, stolen, damaged
- Replacement card (virtual) issued instantly
- Replacement physical card ordered automatically; delivery in 5-10 business days
- Old card number permanently deactivated

### US-4.11: 3D Secure Authentication
**As a** card holder making an online purchase,
**I want to** authenticate via 3D Secure through my TeslaPay app,
**so that** my online payments are secure.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- 3DS 2.0 challenge delivered via push notification
- User approves/declines in-app with biometric auth
- Challenge timeout: 5 minutes
- Fallback to SMS OTP if push delivery fails
- Transaction approved/declined based on user response

### US-4.12: ATM Withdrawal
**As a** card holder,
**I want to** withdraw cash at ATMs using my TeslaPay card,
**so that** I have access to cash when needed.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- ATM withdrawal supported at any Mastercard-accepting ATM
- Free withdrawal allowance per month based on account tier
- Fee displayed in-app for withdrawals exceeding free allowance
- ATM withdrawal limit configurable within tier maximum
- Real-time notification for each ATM transaction

---

## Epic 5: Fuse.io Crypto Integration

### US-5.1: Create Crypto Wallet
**As a** verified user,
**I want to** have a Fuse blockchain wallet created automatically,
**so that** I can access crypto features within TeslaPay.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- Fuse Smart Wallet created automatically after KYC approval
- Wallet address displayed in app with QR code and copy button
- Wallet creation does not require separate onboarding
- User informed about crypto features via brief in-app tutorial
- Wallet is self-custodial (keys derived from device, not stored on server)

### US-5.2: View Crypto Balances
**As a** user with a Fuse wallet,
**I want to** view my crypto token balances alongside my fiat balances,
**so that** I have a complete picture of my finances.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Crypto section on dashboard shows: FUSE, USDC, USDT balances
- Each token shows: balance in token units and EUR equivalent
- EUR equivalent updates with market price (max 60-second delay)
- Total portfolio value includes both fiat and crypto
- Price change indicator (24h %) shown per token

### US-5.3: Buy Crypto with Fiat
**As a** user,
**I want to** buy crypto tokens (FUSE, USDC, USDT) using my EUR balance,
**so that** I can invest in crypto without leaving the app.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- User selects: token to buy, amount in EUR (or token units)
- Exchange rate and fee displayed before confirmation
- Minimum purchase: EUR 5
- Confirmation requires biometric auth
- Tokens credited to Fuse wallet within 60 seconds
- Transaction recorded in both fiat and crypto history

### US-5.4: Sell Crypto to Fiat
**As a** user,
**I want to** sell my crypto tokens back to EUR,
**so that** I can realize gains or access funds in my bank account.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- User selects: token to sell, amount in token units (or EUR target)
- Exchange rate and fee displayed before confirmation
- EUR credited to account within 60 seconds
- Transaction recorded in both fiat and crypto history
- Sell all option available for full balance liquidation

### US-5.5: Send Crypto to External Wallet
**As a** user,
**I want to** send crypto tokens to an external Fuse network address,
**so that** I can transfer tokens to my other wallets or pay others.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- User enters: destination address (paste or QR scan), amount, token
- Address validated (correct format, not blacklisted)
- Network fee estimate shown before confirmation
- Confirmation requires biometric auth
- Transaction submitted to Fuse network; hash displayed
- Status tracked: pending, confirmed, failed

### US-5.6: Receive Crypto from External Wallet
**As a** user,
**I want to** receive crypto tokens from an external sender,
**so that** I can consolidate my crypto in TeslaPay.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Deposit address shown with QR code and copy button
- Incoming transactions detected within 1 block confirmation
- Push notification sent upon receipt
- Received tokens reflected in balance immediately after confirmation

### US-5.7: View Crypto Transaction History
**As a** user,
**I want to** view my crypto transaction history,
**so that** I can track my crypto activity.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- All crypto transactions listed: buys, sells, sends, receives
- Each entry shows: type, token, amount, EUR value at time, date, status, tx hash
- Tx hash links to Fuse block explorer
- Filter by: token, type, date range

### US-5.8: Earn Yield on Stablecoins (Phase 2)
**As a** user holding stablecoins,
**I want to** deposit them into Solid soUSD to earn yield,
**so that** my stablecoins grow over time.

**Priority:** Could | **Size:** XL
**Acceptance Criteria:**
- User can deposit USDC/USDT into Solid soUSD vault
- Current APY displayed transparently
- Yield accrued daily; visible in app
- User can withdraw at any time (no lock-up)
- Risk disclosure presented before first deposit
- Yield reported separately for tax purposes

### US-5.9: Gasless Transactions
**As a** user sending crypto on Fuse network,
**I want** my transactions to be gasless (fees deducted from my token balance),
**so that** I do not need to hold FUSE tokens just to pay fees.

**Priority:** Should | **Size:** L
**Acceptance Criteria:**
- Account abstraction (ERC-4337) used for all user transactions
- Gas fees deducted from the token being sent (e.g., USDC)
- Fee amount shown in token terms before confirmation
- User never sees "insufficient gas" error if they have token balance

---

## Epic 6: Security and Authentication

### US-6.1: Biometric Login
**As a** user,
**I want to** log in using Face ID, Touch ID, or fingerprint,
**so that** I can access my account quickly and securely.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Biometric auth offered at login after initial password setup
- Supports: Face ID (iOS), Touch ID (iOS), fingerprint (Android), face unlock (Android)
- Login completes in under 2 seconds
- Fallback to PIN if biometric fails 3 times
- Fallback to password if PIN fails 3 times

### US-6.2: PIN Setup and Management
**As a** user,
**I want to** set a 6-digit app PIN as a secondary authentication method,
**so that** I have an alternative to biometrics.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- 6-digit PIN required during initial setup
- PIN can be changed in settings (requires current PIN or biometric)
- 5 failed PIN attempts locks account for 30 minutes
- PIN stored securely (hashed, never in plaintext)

### US-6.3: Transaction Confirmation
**As a** user initiating a financial transaction,
**I want to** confirm with biometric auth or PIN,
**so that** unauthorized transactions are prevented even if my phone is unlocked.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- All outgoing payments, card actions, and crypto sends require confirmation
- Confirmation method: biometric (primary) or PIN (fallback)
- Confirmation screen shows transaction summary before auth
- Timeout: 60 seconds to complete confirmation

### US-6.4: Session Management
**As a** user,
**I want to** view and terminate active sessions on my account,
**so that** I can ensure no unauthorized access.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Settings shows all active sessions: device name, OS, location, last active time
- User can terminate any individual session or all other sessions
- Terminated session is invalidated immediately
- Push notification sent when new session is created on a new device

### US-6.5: Two-Factor Authentication
**As a** user,
**I want to** enable 2FA for high-risk actions (password change, beneficiary addition),
**so that** my account has an extra layer of security.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- 2FA enforced for: password change, email change, phone change, new payee, large transfers
- Methods: push notification approval, SMS OTP
- SMS OTP: 6 digits, expires in 5 minutes, max 3 attempts
- Cannot disable 2FA for security-critical actions

### US-6.6: Suspicious Activity Alert
**As a** user,
**I want to** receive alerts for unusual account activity,
**so that** I can react quickly to potential fraud.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Alerts triggered by: login from new device/location, large transaction, rapid successive transactions, international card use
- Alert delivered via push notification and email
- User can confirm or report as fraud directly from alert
- Reporting fraud freezes card and escalates to support

---

## Epic 7: Customer Support

### US-7.1: In-App Chat Support
**As a** user,
**I want to** contact customer support via in-app chat,
**so that** I can get help without leaving the app.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Chat accessible from any screen via floating button or support tab
- First response from agent within 5 minutes during business hours
- Chatbot handles FAQs; escalates to human agent when needed
- Chat history preserved across sessions
- Support available in English and Lithuanian minimum

### US-7.2: FAQ and Help Center
**As a** user,
**I want to** browse a searchable FAQ and help center,
**so that** I can find answers to common questions independently.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Help center accessible from settings and contextual help buttons
- Categories: account, cards, payments, crypto, security, fees
- Search function with keyword matching
- Articles available in 5+ languages
- Content updated at least monthly

### US-7.3: Transaction Dispute
**As a** user,
**I want to** dispute a card transaction I do not recognize,
**so that** I can get my money back for unauthorized charges.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- "Dispute" button on transaction detail screen
- User selects reason: unauthorized, duplicate, goods not received, amount incorrect
- Supporting evidence upload (optional)
- Dispute submitted to card processor; reference number provided
- Status updates via push notification
- Resolution within 15 business days (Mastercard chargeback rules)

---

## Epic 8: Compliance and Reporting

### US-8.1: Ongoing AML Monitoring
**As a** compliance officer,
**I want** all users continuously monitored against sanctions and PEP lists,
**so that** TeslaPay detects changes in user risk profiles.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- Sumsub ongoing monitoring active for all verified users
- New matches generate alert within 24 hours
- Alerts reviewed by compliance team with approve/escalate/freeze options
- Automated account restrictions for high-confidence matches
- Monthly monitoring report generated

### US-8.2: Transaction Monitoring
**As a** compliance officer,
**I want** all transactions analyzed for suspicious patterns,
**so that** TeslaPay can detect and report money laundering.

**Priority:** Must | **Size:** XL
**Acceptance Criteria:**
- Rule-based monitoring: structuring, rapid movement, high-risk corridors
- Threshold alerts: single transaction > EUR 10,000, cumulative > EUR 15,000/month
- ML-based anomaly detection (Phase 2)
- Alerts queue for compliance review with case management
- SAR filing workflow with regulatory templates

### US-8.3: Regulatory Reporting
**As a** compliance officer,
**I want** to generate required regulatory reports for the Bank of Lithuania,
**so that** TeslaPay meets its reporting obligations.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- Reports generated in required format and schedule
- Includes: prudential reports, transaction statistics, complaint data
- Automated data extraction from CBS
- Reports reviewed and approved before submission
- Submission audit trail maintained

### US-8.4: GDPR Data Subject Requests
**As a** user,
**I want to** request access to, correction of, or deletion of my personal data,
**so that** I can exercise my GDPR rights.

**Priority:** Must | **Size:** M
**Acceptance Criteria:**
- Data access request: user receives full data export within 30 days
- Data correction: user can request corrections; verified and applied within 5 business days
- Data deletion: data deleted except where regulatory retention applies; user informed of exceptions
- All requests logged and tracked
- In-app submission form with request types

### US-8.5: Audit Trail Access
**As an** internal auditor,
**I want to** query the complete audit trail for any user or transaction,
**so that** I can investigate issues and prepare for regulatory audits.

**Priority:** Must | **Size:** L
**Acceptance Criteria:**
- Every state change logged: timestamp, actor (user/system/admin), action, before/after values
- Searchable by: user ID, transaction ID, date range, action type
- Logs immutable (append-only)
- Retention: minimum 7 years
- Export capability for auditors

---

## Epic 9: Settings and Preferences

### US-9.1: Language Selection
**As a** user,
**I want to** switch the app language,
**so that** I can use TeslaPay in my preferred language.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- Languages available: English, Lithuanian, Russian, German, Polish (minimum)
- Language switch applied immediately without restart
- Language preference persisted across sessions
- Default: device language if supported; otherwise English

### US-9.2: Dark Mode
**As a** user,
**I want to** switch between light and dark mode,
**so that** I can use the app comfortably in different lighting conditions.

**Priority:** Should | **Size:** S
**Acceptance Criteria:**
- Toggle: Light, Dark, System Default
- System Default follows device setting
- All screens render correctly in both modes
- No text readability issues in either mode

### US-9.3: Fee Schedule and Limits View
**As a** user,
**I want to** view the complete fee schedule and my current limits,
**so that** I know the cost of services and my remaining allowances.

**Priority:** Must | **Size:** S
**Acceptance Criteria:**
- Fee schedule accessible from settings
- Shows: card fees, ATM fees, FX markup, transfer fees, crypto fees
- Current usage vs. limits displayed (e.g., "EUR 200 / EUR 1,000 ATM limit used this month")
- Fees and limits reflect user's account tier

---

## Story Map Summary

| Epic | Must | Should | Could | Total |
|------|------|--------|-------|-------|
| 1. Registration & Onboarding | 6 | 1 | 0 | 7 |
| 2. Account Management | 4 | 2 | 0 | 6 |
| 3. Payments & Transfers | 5 | 3 | 0 | 8 |
| 4. Card Management | 11 | 0 | 0 | 11 |
| 5. Crypto (Fuse.io) | 7 | 1 | 2 | 10 |
| 6. Security & Auth | 6 | 0 | 0 | 6 |
| 7. Customer Support | 3 | 0 | 0 | 3 |
| 8. Compliance & Reporting | 5 | 0 | 0 | 5 |
| 9. Settings & Preferences | 2 | 1 | 0 | 3 |
| **Total** | **49** | **8** | **2** | **59** |
