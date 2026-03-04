import 'package:flutter/material.dart';

import '../../../core/theme/app_theme.dart';

/// Visual representation of a TeslaPay Mastercard.
///
/// Maintains the standard card aspect ratio (1.586:1). Displays TeslaPay logo,
/// Mastercard logo, masked card number, cardholder name, and expiry.
/// Full details are revealed only after biometric authentication.
class CardWidget extends StatelessWidget {
  const CardWidget({
    required this.lastFour,
    required this.cardholderName,
    required this.expiry,
    this.isVirtual = true,
    this.isFrozen = false,
    this.showFullNumber = false,
    this.fullNumber,
    this.cvv,
    this.onTap,
    super.key,
  });

  final String lastFour;
  final String cardholderName;
  final String expiry;
  final bool isVirtual;
  final bool isFrozen;
  final bool showFullNumber;
  final String? fullNumber;
  final String? cvv;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    return Semantics(
      label:
          '${isVirtual ? 'Virtual' : 'Physical'} Mastercard ending in $lastFour, '
          'cardholder $cardholderName, expires $expiry'
          '${isFrozen ? ', card is frozen' : ''}',
      child: GestureDetector(
        onTap: onTap,
        child: AspectRatio(
          aspectRatio: 1.586,
          child: Container(
            decoration: BoxDecoration(
              gradient: isFrozen ? null : AppColors.cardGradient,
              color: isFrozen ? AppColors.neutral400 : null,
              borderRadius: BorderRadius.circular(AppRadius.lg),
              boxShadow: const [
                BoxShadow(
                  color: Color(0x200D1B2A),
                  blurRadius: 12,
                  offset: Offset(0, 4),
                ),
              ],
            ),
            child: Stack(
              children: [
                // Decorative subtle circles
                Positioned(
                  right: -30,
                  top: -30,
                  child: Container(
                    width: 120,
                    height: 120,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: AppColors.white.withValues(alpha: 0.04),
                    ),
                  ),
                ),
                Positioned(
                  right: 20,
                  top: 20,
                  child: Container(
                    width: 80,
                    height: 80,
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      color: AppColors.white.withValues(alpha: 0.04),
                    ),
                  ),
                ),

                // Frost overlay
                if (isFrozen)
                  Positioned.fill(
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(AppRadius.lg),
                        color: AppColors.white.withValues(alpha: 0.1),
                      ),
                      child: Center(
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            const Icon(
                              Icons.ac_unit,
                              color: AppColors.white,
                              size: 32,
                            ),
                            const SizedBox(height: AppSpacing.xs),
                            Text(
                              'FROZEN',
                              style: AppTypography.overline.copyWith(
                                color: AppColors.white,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),

                // Card content
                Padding(
                  padding: const EdgeInsets.all(AppSpacing.lg),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Top row: TeslaPay logo + Mastercard logo
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            'TeslaPay',
                            style: AppTypography.body1.copyWith(
                              color: AppColors.white,
                              fontWeight: FontWeight.w700,
                              letterSpacing: 1,
                            ),
                          ),
                          _MastercardLogo(),
                        ],
                      ),

                      const Spacer(),

                      // Card number
                      Text(
                        showFullNumber && fullNumber != null
                            ? _formatCardNumber(fullNumber!)
                            : '**** **** **** $lastFour',
                        style: AppTypography.mono.copyWith(
                          color: AppColors.white,
                          fontSize: 18,
                          letterSpacing: 2,
                        ),
                      ),

                      if (showFullNumber && cvv != null) ...[
                        const SizedBox(height: AppSpacing.xs),
                        Text(
                          'CVV: $cvv',
                          style: AppTypography.mono.copyWith(
                            color:
                                AppColors.white.withValues(alpha: 0.8),
                            fontSize: 14,
                          ),
                        ),
                      ],

                      const Spacer(),

                      // Bottom row: name + expiry
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                'CARDHOLDER',
                                style: AppTypography.overline.copyWith(
                                  color: AppColors.white
                                      .withValues(alpha: 0.6),
                                ),
                              ),
                              const SizedBox(height: AppSpacing.xxs),
                              Text(
                                cardholderName.toUpperCase(),
                                style: AppTypography.body2.copyWith(
                                  color: AppColors.white,
                                  fontWeight: FontWeight.w500,
                                  letterSpacing: 0.5,
                                ),
                              ),
                            ],
                          ),
                          Column(
                            crossAxisAlignment: CrossAxisAlignment.end,
                            children: [
                              Text(
                                'EXPIRES',
                                style: AppTypography.overline.copyWith(
                                  color: AppColors.white
                                      .withValues(alpha: 0.6),
                                ),
                              ),
                              const SizedBox(height: AppSpacing.xxs),
                              Text(
                                expiry,
                                style: AppTypography.body2.copyWith(
                                  color: AppColors.white,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ],
                          ),
                        ],
                      ),

                      // Virtual badge
                      if (isVirtual) ...[
                        const SizedBox(height: AppSpacing.sm),
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: AppSpacing.sm,
                            vertical: AppSpacing.xxs,
                          ),
                          decoration: BoxDecoration(
                            color: AppColors.white.withValues(alpha: 0.15),
                            borderRadius:
                                BorderRadius.circular(AppRadius.xs),
                          ),
                          child: Text(
                            'VIRTUAL',
                            style: AppTypography.overline.copyWith(
                              color: AppColors.white,
                            ),
                          ),
                        ),
                      ],
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  String _formatCardNumber(String number) {
    final clean = number.replaceAll(RegExp(r'\s'), '');
    final buffer = StringBuffer();
    for (var i = 0; i < clean.length; i++) {
      if (i > 0 && i % 4 == 0) buffer.write(' ');
      buffer.write(clean[i]);
    }
    return buffer.toString();
  }
}

/// Simple Mastercard logo using overlapping circles.
class _MastercardLogo extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    const size = 28.0;
    return SizedBox(
      width: size * 1.6,
      height: size,
      child: Stack(
        children: [
          Positioned(
            left: 0,
            child: Container(
              width: size,
              height: size,
              decoration: BoxDecoration(
                color: const Color(0xFFEB001B).withValues(alpha: 0.9),
                shape: BoxShape.circle,
              ),
            ),
          ),
          Positioned(
            left: size * 0.6,
            child: Container(
              width: size,
              height: size,
              decoration: BoxDecoration(
                color: const Color(0xFFF79E1B).withValues(alpha: 0.9),
                shape: BoxShape.circle,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
