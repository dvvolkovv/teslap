import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:phosphor_flutter/phosphor_flutter.dart';

import '../../../core/theme/app_theme.dart';
import '../../../shared/widgets/app_button.dart';
import '../../../shared/widgets/app_input.dart';
import '../bloc/payment_bloc.dart';

/// Multi-step send money flow.
/// Steps: 1) Select recipient  2) Enter amount  3) Review and confirm
class SendMoneyScreen extends StatefulWidget {
  const SendMoneyScreen({super.key});

  @override
  State<SendMoneyScreen> createState() => _SendMoneyScreenState();
}

class _SendMoneyScreenState extends State<SendMoneyScreen> {
  int _step = 0;
  final _ibanController = TextEditingController();
  final _nameController = TextEditingController();
  final _amountController = TextEditingController();
  final _referenceController = TextEditingController();
  bool _isInstant = false;

  @override
  void dispose() {
    _ibanController.dispose();
    _nameController.dispose();
    _amountController.dispose();
    _referenceController.dispose();
    super.dispose();
  }

  void _nextStep() {
    if (_step < 2) {
      setState(() => _step++);
    } else {
      _send();
    }
  }

  void _previousStep() {
    if (_step > 0) {
      setState(() => _step--);
    } else {
      context.pop();
    }
  }

  void _send() {
    context.read<PaymentBloc>().add(PaymentSepaRequested(
          senderAccountId: '', // Will use default account on backend
          recipientIban: _ibanController.text,
          recipientName: _nameController.text,
          amount: _amountController.text,
          reference: _referenceController.text.isNotEmpty
              ? _referenceController.text
              : null,
        ));
  }

  @override
  Widget build(BuildContext context) {
    return BlocListener<PaymentBloc, PaymentState>(
      listener: (context, state) {
        if (state is PaymentSuccess) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(
                'EUR ${_amountController.text} sent to ${_nameController.text}',
              ),
              backgroundColor: AppColors.success500,
            ),
          );
          context.pop();
        } else if (state is PaymentError) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(state.message),
              backgroundColor: AppColors.error500,
            ),
          );
        }
      },
      child: Scaffold(
        appBar: AppBar(
          leading: BackButton(onPressed: _previousStep),
          title: const Text('Send Money'),
        ),
        body: SafeArea(
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: AppSpacing.screenMargin,
                ),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(AppRadius.full),
                  child: LinearProgressIndicator(
                    value: (_step + 1) / 3,
                    backgroundColor: AppColors.neutral100,
                    valueColor: const AlwaysStoppedAnimation<Color>(
                      AppColors.primary500,
                    ),
                    minHeight: 4,
                  ),
                ),
              ),
              Expanded(
                child: AnimatedSwitcher(
                  duration: const Duration(milliseconds: 250),
                  child: _buildStep(),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildStep() {
    return switch (_step) {
      0 => _RecipientStep(
          key: const ValueKey('recipient'),
          ibanController: _ibanController,
          nameController: _nameController,
          onContinue: _nextStep,
        ),
      1 => _AmountStep(
          key: const ValueKey('amount'),
          amountController: _amountController,
          referenceController: _referenceController,
          isInstant: _isInstant,
          onInstantChanged: (v) => setState(() => _isInstant = v),
          onContinue: _nextStep,
        ),
      _ => BlocBuilder<PaymentBloc, PaymentState>(
          builder: (context, state) {
            return _ReviewStep(
              key: const ValueKey('review'),
              name: _nameController.text,
              iban: _ibanController.text,
              amount: _amountController.text,
              reference: _referenceController.text,
              isInstant: _isInstant,
              isSending: state is PaymentLoading,
              onConfirm: _nextStep,
            );
          },
        ),
    };
  }
}

// ---------------------------------------------------------------------------
// Step widgets
// ---------------------------------------------------------------------------

class _RecipientStep extends StatelessWidget {
  const _RecipientStep({
    super.key,
    required this.ibanController,
    required this.nameController,
    required this.onContinue,
  });

  final TextEditingController ibanController;
  final TextEditingController nameController;
  final VoidCallback onContinue;

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.screenMargin),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: AppSpacing.md),
          Text('Select Recipient', style: AppTypography.h1),
          const SizedBox(height: AppSpacing.lg),
          AppInput(
            label: 'Recipient name',
            hint: 'John Smith',
            controller: nameController,
            textInputAction: TextInputAction.next,
          ),
          const SizedBox(height: AppSpacing.md),
          AppInput(
            label: 'IBAN',
            hint: 'LT12 3456 7890 1234 5678',
            controller: ibanController,
            variant: AppInputVariant.iban,
            textInputAction: TextInputAction.done,
          ),
          const SizedBox(height: AppSpacing.xl),
          AppButton(label: 'Continue', onPressed: onContinue),
        ],
      ),
    );
  }
}

class _AmountStep extends StatelessWidget {
  const _AmountStep({
    super.key,
    required this.amountController,
    required this.referenceController,
    required this.isInstant,
    required this.onInstantChanged,
    required this.onContinue,
  });

  final TextEditingController amountController;
  final TextEditingController referenceController;
  final bool isInstant;
  final ValueChanged<bool> onInstantChanged;
  final VoidCallback onContinue;

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.screenMargin),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: AppSpacing.md),
          Text('Enter Amount', style: AppTypography.h1),
          const SizedBox(height: AppSpacing.lg),
          AppInput(
            label: 'Amount',
            hint: '0.00',
            controller: amountController,
            variant: AppInputVariant.amount,
            prefixIcon: const Padding(
              padding: EdgeInsets.only(left: AppSpacing.md),
              child: Text('EUR',
                  style: TextStyle(
                      fontWeight: FontWeight.w600, fontSize: 16)),
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          AppInput(
            label: 'Reference (optional)',
            hint: 'What is this for?',
            controller: referenceController,
            textInputAction: TextInputAction.done,
          ),
          const SizedBox(height: AppSpacing.lg),
          Row(
            children: [
              Expanded(
                  child: Text('SEPA Instant',
                      style: AppTypography.body1)),
              Switch(
                value: isInstant,
                onChanged: onInstantChanged,
                activeColor: AppColors.primary500,
              ),
            ],
          ),
          Text(
            isInstant
                ? 'Arrives in 10 seconds. Fee: EUR 0.50'
                : 'Arrives in 1 business day. Fee: EUR 0.00',
            style: AppTypography.caption
                .copyWith(color: AppColors.neutral500),
          ),
          const SizedBox(height: AppSpacing.xl),
          AppButton(label: 'Continue', onPressed: onContinue),
        ],
      ),
    );
  }
}

class _ReviewStep extends StatelessWidget {
  const _ReviewStep({
    super.key,
    required this.name,
    required this.iban,
    required this.amount,
    required this.reference,
    required this.isInstant,
    required this.isSending,
    required this.onConfirm,
  });

  final String name;
  final String iban;
  final String amount;
  final String reference;
  final bool isInstant;
  final bool isSending;
  final VoidCallback onConfirm;

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.screenMargin),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const SizedBox(height: AppSpacing.md),
          Text('Review Your Transfer', style: AppTypography.h1),
          const SizedBox(height: AppSpacing.lg),
          _ReviewRow(label: 'From', value: 'EUR Account'),
          _ReviewRow(label: 'To', value: name),
          _ReviewRow(label: 'IBAN', value: iban),
          _ReviewRow(label: 'Amount', value: 'EUR $amount'),
          _ReviewRow(
              label: 'Fee',
              value: isInstant ? 'EUR 0.50' : 'EUR 0.00'),
          _ReviewRow(
            label: 'Total',
            value:
                'EUR ${(double.tryParse(amount) ?? 0) + (isInstant ? 0.50 : 0)}',
            isBold: true,
          ),
          _ReviewRow(
            label: 'Delivery',
            value:
                isInstant ? 'Instant (SEPA Inst)' : '1 business day',
          ),
          if (reference.isNotEmpty)
            _ReviewRow(label: 'Reference', value: reference),
          const SizedBox(height: AppSpacing.xl),
          AppButton(
            label: 'Confirm with Face ID',
            icon: PhosphorIconsRegular.fingerprint,
            onPressed: onConfirm,
            isLoading: isSending,
          ),
        ],
      ),
    );
  }
}

class _ReviewRow extends StatelessWidget {
  const _ReviewRow({
    required this.label,
    required this.value,
    this.isBold = false,
  });

  final String label;
  final String value;
  final bool isBold;

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: AppSpacing.sm),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 100,
            child: Text(label,
                style: AppTypography.body2
                    .copyWith(color: AppColors.neutral500)),
          ),
          Expanded(
            child: Text(
              value,
              style: isBold
                  ? AppTypography.body1
                      .copyWith(fontWeight: FontWeight.w600)
                  : AppTypography.body1,
              textAlign: TextAlign.right,
            ),
          ),
        ],
      ),
    );
  }
}
