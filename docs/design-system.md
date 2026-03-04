# TeslaPay Design System

**Version:** 1.0
**Date:** 2026-03-03
**Platform:** Flutter (iOS + Android)
**Author:** Senior UI/UX Designer, Dream Team

---

## 1. Design Principles

1. **Trust First** -- Every design decision must reinforce that TeslaPay is a regulated, safe financial institution. No gimmicks.
2. **Clarity Over Cleverness** -- Financial data must be instantly readable. When in doubt, simplify.
3. **Progressive Disclosure** -- Show what matters now; reveal complexity on demand. Especially critical for crypto features.
4. **One App, Two Worlds** -- Fiat and crypto must feel like a single, cohesive experience, never bolted-on.
5. **Speed Is UX** -- Perceived and actual performance. Skeleton loaders, optimistic UI, instant feedback.

---

## 2. Color Palette

### 2.1 Primary Colors

| Token                  | Hex       | Usage                                           |
|------------------------|-----------|------------------------------------------------|
| `color.primary.500`    | `#0066FF` | Primary actions, links, active states, brand accent |
| `color.primary.600`    | `#0052CC` | Primary button pressed state                    |
| `color.primary.400`    | `#3385FF` | Primary button hover/focus                      |
| `color.primary.100`    | `#E6F0FF` | Primary tint backgrounds, selected list items   |
| `color.primary.50`     | `#F0F6FF` | Subtle primary background                       |

Rationale: A confident, trustworthy blue anchors the brand. Distinctly different from Revolut (dark violet) and N26 (teal). Conveys stability and modernity.

### 2.2 Secondary / Accent Colors

| Token                     | Hex       | Usage                                        |
|---------------------------|-----------|----------------------------------------------|
| `color.accent.500`        | `#00D4AA` | Crypto/Web3 features, yield indicators, success accent |
| `color.accent.600`        | `#00B892` | Pressed state for crypto actions              |
| `color.accent.100`        | `#E6FBF5` | Crypto section backgrounds                    |

Rationale: A vibrant teal-green signals innovation and differentiates the crypto side of the app from traditional banking blue.

### 2.3 Semantic Colors

| Token                      | Hex       | Usage                                   |
|----------------------------|-----------|----------------------------------------|
| `color.success.500`        | `#00C853` | Successful transactions, positive amounts |
| `color.success.100`        | `#E8F5E9` | Success background                      |
| `color.warning.500`        | `#FF9800` | Pending states, warnings, KYC review    |
| `color.warning.100`        | `#FFF3E0` | Warning background                      |
| `color.error.500`          | `#F44336` | Errors, declined transactions, negative amounts |
| `color.error.100`          | `#FFEBEE` | Error background                        |
| `color.info.500`           | `#2196F3` | Informational notices, tips             |
| `color.info.100`           | `#E3F2FD` | Info background                         |

### 2.4 Neutral Colors

| Token                    | Hex       | Usage                                    |
|--------------------------|-----------|------------------------------------------|
| `color.neutral.900`      | `#0D1B2A` | Primary text (light mode)                |
| `color.neutral.800`      | `#1B2838` | Headlines                                |
| `color.neutral.700`      | `#2E3A4D` | Secondary text                           |
| `color.neutral.500`      | `#6B7B8F` | Placeholder text, disabled state         |
| `color.neutral.400`      | `#94A3B8` | Icons (inactive)                         |
| `color.neutral.200`      | `#D1D9E6` | Borders, dividers                        |
| `color.neutral.100`      | `#EDF1F7` | Card backgrounds, input fills            |
| `color.neutral.50`       | `#F7F9FC` | Page background (light mode)             |
| `color.neutral.white`    | `#FFFFFF` | Cards, surfaces (light mode)             |

### 2.5 Dark Mode Colors

| Token                           | Hex       | Usage                              |
|---------------------------------|-----------|------------------------------------|
| `color.dark.background`         | `#0D1117` | Page background                    |
| `color.dark.surface`            | `#161B22` | Cards, bottom sheets               |
| `color.dark.surfaceElevated`    | `#1C2333` | Elevated cards, modals             |
| `color.dark.border`             | `#30363D` | Borders, dividers                  |
| `color.dark.textPrimary`        | `#F0F6FC` | Primary text                       |
| `color.dark.textSecondary`      | `#8B949E` | Secondary text                     |
| `color.dark.textTertiary`       | `#6E7681` | Placeholder, disabled              |

Note: Primary and semantic colors remain the same in dark mode but may be slightly adjusted for contrast (e.g., `color.primary.400` used instead of `.500` for better readability on dark backgrounds).

### 2.6 Gradient

| Token                  | Value                                   | Usage                        |
|------------------------|-----------------------------------------|------------------------------|
| `gradient.brand`       | `linear(135deg, #0066FF, #00D4AA)`      | Card visuals, premium badges |
| `gradient.crypto`      | `linear(135deg, #00D4AA, #00B4D8)`      | Crypto section headers       |
| `gradient.card`        | `linear(160deg, #1B2838, #0D1B2A)`      | Virtual card display         |

---

## 3. Typography

### 3.1 Font Family

| Role       | Font                | Fallback         | Source       |
|------------|---------------------|------------------|--------------|
| Primary    | **Inter**           | SF Pro, Roboto   | Google Fonts |
| Monospace  | **JetBrains Mono**  | SF Mono, monospace | Google Fonts |

Rationale: Inter is optimized for screens, has excellent number readability (critical for a banking app), supports all required languages (Latin, Cyrillic for Russian/Lithuanian), and is widely used in fintech.

### 3.2 Type Scale (8px baseline, 1.5x for large text)

| Token              | Size  | Weight    | Line Height | Letter Spacing | Usage                        |
|--------------------|-------|-----------|-------------|----------------|------------------------------|
| `type.display1`    | 40px  | Bold 700  | 48px        | -0.5px         | Balance (large, hero)        |
| `type.display2`    | 32px  | Bold 700  | 40px        | -0.25px        | Section balance              |
| `type.h1`          | 24px  | SemiBold 600 | 32px     | 0              | Screen titles                |
| `type.h2`          | 20px  | SemiBold 600 | 28px     | 0              | Section titles               |
| `type.h3`          | 18px  | Medium 500 | 24px       | 0              | Card titles, subsection      |
| `type.body1`       | 16px  | Regular 400 | 24px      | 0              | Primary body text            |
| `type.body2`       | 14px  | Regular 400 | 20px      | 0.1px          | Secondary text, descriptions |
| `type.caption`     | 12px  | Regular 400 | 16px      | 0.2px          | Labels, timestamps, meta     |
| `type.overline`    | 10px  | SemiBold 600 | 16px     | 1.5px          | Section labels (UPPERCASE)   |
| `type.button`      | 16px  | SemiBold 600 | 24px     | 0.5px          | Button labels                |
| `type.mono`        | 14px  | Regular 400 | 20px      | 0              | IBAN, card numbers, tx hashes |

### 3.3 Number Display

Monetary amounts use **tabular (monospaced) figures** from Inter to ensure decimal alignment. Currency symbols precede the amount with a thin space: `EUR 1,234.56`.

---

## 4. Spacing and Grid

### 4.1 Spacing Scale (8px grid)

| Token       | Value | Usage                                      |
|-------------|-------|--------------------------------------------|
| `space.2xs` | 2px   | Tight inline spacing                       |
| `space.xs`  | 4px   | Icon-to-text gap, inline elements          |
| `space.sm`  | 8px   | Compact padding, list item internal spacing|
| `space.md`  | 16px  | Standard padding, card content padding     |
| `space.lg`  | 24px  | Section spacing, card-to-card gap          |
| `space.xl`  | 32px  | Major section separation                   |
| `space.2xl` | 48px  | Screen top padding, hero spacing           |
| `space.3xl` | 64px  | Very large separation                      |

### 4.2 Layout Grid

- **Screen margin:** 16px (compact) / 20px (standard) left and right
- **Max content width:** 428px (iPhone 14 Pro Max reference)
- **Column grid:** 4-column grid with 16px gutter for card layouts
- **Card padding:** 16px internal
- **List item height:** 64px minimum (touch target + content)

### 4.3 Safe Areas

- **Top:** Respect system status bar + 8px
- **Bottom:** Respect system home indicator + bottom navigation height (56px) + 8px

---

## 5. Border Radius

| Token            | Value | Usage                                 |
|------------------|-------|---------------------------------------|
| `radius.xs`      | 4px   | Badges, small chips                   |
| `radius.sm`      | 8px   | Input fields, small buttons           |
| `radius.md`      | 12px  | Cards, modals, bottom sheets          |
| `radius.lg`      | 16px  | Large cards, card visuals             |
| `radius.xl`      | 24px  | Pills, floating action buttons        |
| `radius.full`    | 9999px | Avatars, round buttons               |

---

## 6. Elevation / Shadows

| Token              | Value                                          | Usage                  |
|--------------------|------------------------------------------------|------------------------|
| `shadow.sm`        | `0 1px 3px rgba(13,27,42,0.08)`               | Cards at rest          |
| `shadow.md`        | `0 4px 12px rgba(13,27,42,0.12)`              | Elevated cards, FABs   |
| `shadow.lg`        | `0 8px 24px rgba(13,27,42,0.16)`              | Modals, bottom sheets  |
| `shadow.card`      | `0 2px 8px rgba(13,27,42,0.06)`               | Transaction list items |

Dark mode: shadows are replaced with border emphasis (`color.dark.border` at 1px) since shadows are not visible on dark backgrounds.

---

## 7. Component Library

### 7.1 Buttons

#### Primary Button
- Background: `color.primary.500`
- Text: `#FFFFFF`, `type.button`
- Height: 56px
- Border radius: `radius.sm` (8px)
- Full width on forms; inline width with 24px horizontal padding
- States: Default, Pressed (`color.primary.600`), Disabled (40% opacity), Loading (spinner replaces text)

#### Secondary Button
- Background: transparent
- Border: 1.5px `color.primary.500`
- Text: `color.primary.500`, `type.button`
- Height: 56px
- Same states as primary with appropriate color adjustments

#### Tertiary / Text Button
- Background: transparent
- Text: `color.primary.500`, `type.button`
- No border
- Underline on hover/focus
- Used for "Cancel", "Skip", "Learn more"

#### Danger Button
- Background: `color.error.500`
- Text: `#FFFFFF`
- Used for destructive actions: close account, block card, report fraud

#### Quick Action Button (Icon + Label)
- Circular icon container: 48px diameter, `color.primary.100` background
- Icon: 24px, `color.primary.500`
- Label below: `type.caption`, `color.neutral.700`
- Used on dashboard for Send, Request, Exchange, Card

### 7.2 Input Fields

#### Text Input
- Height: 56px
- Background: `color.neutral.100` (light) / `color.dark.surface` (dark)
- Border: 1px `color.neutral.200` (light) / `color.dark.border` (dark)
- Border radius: `radius.sm`
- Label: `type.caption`, floats above on focus
- Placeholder: `color.neutral.500`
- Focus state: border becomes `color.primary.500`, 2px
- Error state: border becomes `color.error.500`, error message below in `color.error.500` `type.caption`
- Padding: 16px horizontal

#### Amount Input
- Large centered display: `type.display1` for amount
- Currency selector chip to the left
- Numpad below (custom, not system keyboard)
- Shows converted equivalent in secondary currency below in `type.body2`

#### PIN Input
- 6 dots in a row, each 12px diameter
- Filled dot: `color.primary.500`
- Empty dot: `color.neutral.200`
- Custom numpad with biometric button in bottom-left

#### Search Input
- Height: 48px
- Leading search icon
- Trailing clear button (visible when text entered)
- Border radius: `radius.xl`
- Background: `color.neutral.100`

### 7.3 Cards

#### Balance Card (Dashboard Hero)
- Full width minus margins
- Height: ~160px
- Background: `gradient.brand` or solid `color.neutral.900`
- Content: Total balance (`type.display1`, white), currency toggle, percentage change
- Border radius: `radius.lg`

#### Transaction Item
- Height: 72px
- Left: Merchant/payee icon (40px circle)
- Center: Name (`type.body1`), category/date (`type.caption`)
- Right: Amount (`type.body1`, bold), status indicator
- Positive amounts: `color.success.500`; negative amounts: `color.neutral.900` (not red, to avoid alarm)
- Divider: `color.neutral.200`, 1px, indented from left icon

#### Currency Card
- Height: 80px
- Left: Currency flag (24px round)
- Center: Currency name/code (`type.body1`), balance in native (`type.h3`)
- Right: EUR equivalent (`type.body2`, `color.neutral.500`)
- Tappable -- navigates to currency detail

#### Virtual Card Display
- Aspect ratio: 1.586:1 (standard card ratio)
- Background: `gradient.card`
- Content: Mastercard logo (top right), TeslaPay logo (top left), card number (masked), cardholder name, expiry
- Card number revealed on biometric auth tap
- Shadow: `shadow.md`
- Border radius: `radius.lg`

#### Crypto Token Card
- Height: 80px
- Left: Token icon (40px)
- Center: Token name (`type.body1`), balance in tokens (`type.h3`)
- Right: EUR equivalent, 24h change badge (`color.success.500` or `color.error.500`)

#### Info Card
- Background: `color.info.100` (light) / tinted dark variant
- Left: Info icon
- Body: `type.body2`
- Optional dismiss button
- Border radius: `radius.md`

### 7.4 Bottom Navigation Bar

- Height: 56px (plus safe area)
- Background: `color.neutral.white` (light) / `color.dark.surface` (dark)
- Top border: 1px `color.neutral.200`
- 5 tabs with icon (24px) + label (`type.caption`)
- Active: `color.primary.500` (icon + text)
- Inactive: `color.neutral.400` (icon), `color.neutral.500` (text)
- Tab items: Home, Payments, Card, Crypto, Profile

### 7.5 Top App Bar

- Height: 56px
- Background: transparent (scrolls with content) or solid on scroll
- Left: Back arrow or hamburger (never hamburger in this app -- use back only)
- Center: Screen title (`type.h2`)
- Right: Contextual actions (search, filter, help)

### 7.6 Bottom Sheet

- Border radius: `radius.lg` (top corners only)
- Handle: 40px wide, 4px tall, centered, `color.neutral.200`
- Background: `color.neutral.white` (light) / `color.dark.surface` (dark)
- Padding: 24px top, 16px sides, safe area bottom
- Dismissible by drag-down or tapping scrim
- Scrim: `rgba(0,0,0,0.5)`

### 7.7 Badges and Chips

#### Status Badge
- Height: 24px
- Padding: 4px 8px
- Border radius: `radius.xs`
- Variants: Pending (warning), Completed (success), Failed (error), Processing (info)
- Text: `type.caption`, matching semantic color

#### Tier Badge
- Premium: `gradient.brand` background, white text
- Standard: `color.primary.100` background, `color.primary.500` text
- Basic: `color.neutral.100` background, `color.neutral.700` text

#### Currency Chip
- Height: 32px
- Border radius: `radius.xl`
- Flag + code (e.g., flag + "EUR")
- Background: `color.neutral.100`
- Selected: `color.primary.100`, border `color.primary.500`

### 7.8 Toggle / Switch
- Width: 51px, Height: 31px (iOS standard)
- Active: `color.primary.500`
- Inactive: `color.neutral.200`
- Thumb: white, `shadow.sm`

### 7.9 Notification Toast
- Position: top of screen, below status bar
- Height: auto (min 56px)
- Background: `color.neutral.900` (light) / `color.dark.surfaceElevated` (dark)
- Text: white, `type.body2`
- Left icon: semantic color circle
- Auto-dismiss after 4 seconds
- Swipe up to dismiss

### 7.10 Skeleton Loader
- Background: `color.neutral.100`
- Shimmer: left-to-right gradient animation
- Shapes match content layout (rounded rects for text, circles for avatars)
- Duration: 1.5s per cycle

### 7.11 Empty State
- Centered illustration (line art style, using primary + accent colors)
- Title: `type.h2`
- Description: `type.body2`, `color.neutral.500`
- CTA button below

---

## 8. Iconography

### 8.1 Icon Set

**Recommended:** Phosphor Icons (https://phosphoricons.com)

Rationale:
- Consistent line weight (1.5px stroke at 24px)
- Available in multiple weights: thin, light, regular, bold, fill
- Excellent coverage of finance, crypto, and general UI icons
- MIT licensed, suitable for commercial use
- Flutter package available (`phosphor_flutter`)

### 8.2 Icon Sizes

| Token        | Size | Usage                                  |
|--------------|------|----------------------------------------|
| `icon.sm`    | 16px | Inline text icons, badges              |
| `icon.md`    | 24px | Standard UI icons, nav bar, list items |
| `icon.lg`    | 32px | Quick action icons, card actions       |
| `icon.xl`    | 48px | Empty states, onboarding illustrations |

### 8.3 Icon Style Rules

- Use **Regular** weight for navigation and standard UI
- Use **Bold** weight for active navigation tab icon
- Use **Fill** weight sparingly, only for active/selected states
- Colors: `color.neutral.700` default, `color.primary.500` active, semantic colors for status
- Always pair icons with text labels for accessibility

### 8.4 Custom Icons Needed

- TeslaPay logo (wordmark + icon mark)
- Currency flags (round, 24px)
- Crypto token icons (FUSE, USDC, USDT) -- 40px
- Mastercard logo (card display)
- Apple Pay / Google Pay badges

---

## 9. Motion and Animation

### 9.1 Principles

- **Purposeful:** Every animation communicates something (state change, spatial relationship)
- **Fast:** Financial app users expect responsiveness. Keep transitions under 300ms
- **Consistent:** Same elements animate the same way everywhere

### 9.2 Timing

| Token              | Duration | Curve              | Usage                          |
|--------------------|----------|--------------------|--------------------------------|
| `motion.fast`      | 150ms    | easeOut            | Button press, toggle, fade     |
| `motion.normal`    | 250ms    | easeInOut          | Page transitions, card expand  |
| `motion.slow`      | 350ms    | easeInOut          | Bottom sheet open, modal       |
| `motion.spring`    | 400ms    | spring(0.7, 0.8)   | Pull-to-refresh, bounce        |

### 9.3 Specific Animations

- **Page transition:** Slide left (push), slide right (pop), fade for tab switches
- **Balance update:** Counter animation (numbers roll) on dashboard load
- **Card freeze/unfreeze:** Card visual grays out with frost overlay animation
- **Transaction notification:** Slide down from top, auto-dismiss slide up
- **Biometric prompt:** Subtle scale-up of fingerprint/face icon
- **Success confirmation:** Check mark draws itself in a circle (Lottie animation)
- **Loading states:** Skeleton shimmer, not spinners (except on buttons)

---

## 10. Accessibility (WCAG 2.1 AA)

### 10.1 Color Contrast

- Text on backgrounds: minimum 4.5:1 ratio for `type.body2` and smaller
- Large text (18px+ bold or 24px+ regular): minimum 3:1 ratio
- Interactive elements: minimum 3:1 against adjacent colors
- All color combinations validated with contrast checker

### 10.2 Touch Targets

- Minimum touch target: 48x48px (per WCAG and Material Design)
- Adequate spacing between targets: minimum 8px
- Bottom nav items span full width of their section

### 10.3 Screen Reader Support

- All images and icons have `semanticLabel` (Flutter)
- Interactive elements have meaningful labels (not "Button 1")
- Balance amounts read as: "Total balance: one thousand two hundred thirty-four euros and fifty-six cents"
- Card numbers read digit by digit with pauses
- Transaction lists announce: "Payment to [merchant], [amount], [date], [status]"
- Custom components implement `Semantics` widget in Flutter

### 10.4 Text Scaling

- Support system text scale up to 200%
- Layout must not break at 150% text scale (tested)
- Long text truncates with ellipsis and full text available on tap

### 10.5 Reduced Motion

- Respect `MediaQuery.of(context).disableAnimations`
- When reduced motion is on: replace animations with instant state changes
- No autoplay animations or looping animations without user control

### 10.6 Color Blindness

- Never use color alone to convey information
- Always pair with icons, text labels, or patterns
- Transaction amounts: positive uses UP arrow icon + green; negative uses DOWN arrow + amount
- Status badges: always include text label, not just colored dot

---

## 11. Dark Mode Implementation

### 11.1 Strategy

- Respect system setting by default (`MediaQuery.of(context).platformBrightness`)
- User can override in Settings: Light / Dark / System
- Preference persisted in local storage

### 11.2 Mapping Rules

| Light Token               | Dark Equivalent                    |
|---------------------------|------------------------------------|
| `color.neutral.white`     | `color.dark.background`            |
| `color.neutral.50`        | `color.dark.background`            |
| `color.neutral.100`       | `color.dark.surface`               |
| `color.neutral.200`       | `color.dark.border`                |
| `color.neutral.900`       | `color.dark.textPrimary`           |
| `color.neutral.700`       | `color.dark.textSecondary`         |
| `color.neutral.500`       | `color.dark.textTertiary`          |
| Shadows                   | 1px borders (`color.dark.border`)  |

### 11.3 Special Cases

- Virtual card display: same gradient in both modes (dark card on dark bg works)
- Charts and graphs: use slightly lighter variants in dark mode
- Illustrations: provide dark-mode optimized versions

---

## 12. Responsive Patterns

### 12.1 Device Targets

| Device Class   | Width Range  | Notes                           |
|----------------|-------------|--------------------------------|
| Small phone    | 320-375px   | iPhone SE, compact Androids    |
| Standard phone | 376-428px   | iPhone 14/15, most Androids    |
| Large phone    | 429-480px   | iPhone Pro Max, Samsung Ultra  |

### 12.2 Adaptation Rules

- Balance card: scales text size with screen width
- Transaction list: items are always full-width
- Card display: maintains aspect ratio, scales width to fit screen minus margins
- Bottom sheet: max height 90% of screen
- Quick actions: 4 items in row on 375px+, 3 items with scroll on smaller

---

## 13. Localization Considerations

### 13.1 Supported Languages (MVP)

English (en), Lithuanian (lt), Russian (ru), German (de), Polish (pl)

### 13.2 Layout Rules

- All languages are LTR (no RTL requirement for MVP)
- German text is typically 30% longer than English -- test all layouts with German strings
- Russian uses Cyrillic -- ensure Inter font Cyrillic subset is loaded
- Currency formatting respects locale: `1.234,56` (DE) vs `1,234.56` (EN)
- Date formatting respects locale: `03.03.2026` (DE) vs `3 Mar 2026` (EN)

---

## 14. Design Token Export Format

All tokens are exported as:
- **Flutter:** Dart constants in `lib/theme/tokens.dart`
- **Figma:** Figma Variables for design handoff
- **JSON:** For documentation and cross-platform reference

Example Flutter implementation:

```dart
class TeslaPayColors {
  static const primary500 = Color(0xFF0066FF);
  static const primary600 = Color(0xFF0052CC);
  static const primary400 = Color(0xFF3385FF);
  static const primary100 = Color(0xFFE6F0FF);
  static const accent500 = Color(0xFF00D4AA);
  static const success500 = Color(0xFF00C853);
  static const warning500 = Color(0xFFFF9800);
  static const error500 = Color(0xFFF44336);
  // ... etc
}

class TeslaPayTypography {
  static const display1 = TextStyle(
    fontFamily: 'Inter',
    fontSize: 40,
    fontWeight: FontWeight.w700,
    height: 1.2,
    letterSpacing: -0.5,
  );
  // ... etc
}

class TeslaPaySpacing {
  static const xs = 4.0;
  static const sm = 8.0;
  static const md = 16.0;
  static const lg = 24.0;
  static const xl = 32.0;
  // ... etc
}
```
