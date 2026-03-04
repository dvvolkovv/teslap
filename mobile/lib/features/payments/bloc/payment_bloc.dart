import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/repositories/payment_repository.dart';

// =============================================================================
// Events
// =============================================================================

abstract class PaymentEvent extends Equatable {
  const PaymentEvent();

  @override
  List<Object?> get props => [];
}

class PaymentSepaRequested extends PaymentEvent {
  const PaymentSepaRequested({
    required this.senderAccountId,
    required this.recipientIban,
    required this.recipientName,
    required this.amount,
    this.currency = 'EUR',
    this.reference,
  });

  final String senderAccountId;
  final String recipientIban;
  final String recipientName;
  final String amount;
  final String currency;
  final String? reference;

  @override
  List<Object?> get props =>
      [senderAccountId, recipientIban, recipientName, amount, currency];
}

class PaymentInternalRequested extends PaymentEvent {
  const PaymentInternalRequested({
    required this.senderAccountId,
    required this.recipientAccountId,
    required this.amount,
    this.currency = 'EUR',
    this.reference,
  });

  final String senderAccountId;
  final String recipientAccountId;
  final String amount;
  final String currency;
  final String? reference;

  @override
  List<Object?> get props =>
      [senderAccountId, recipientAccountId, amount, currency];
}

class PaymentListRequested extends PaymentEvent {
  const PaymentListRequested({this.accountId});

  final String? accountId;

  @override
  List<Object?> get props => [accountId];
}

// =============================================================================
// States
// =============================================================================

abstract class PaymentState extends Equatable {
  const PaymentState();

  @override
  List<Object?> get props => [];
}

class PaymentInitial extends PaymentState {
  const PaymentInitial();
}

class PaymentLoading extends PaymentState {
  const PaymentLoading();
}

class PaymentSuccess extends PaymentState {
  const PaymentSuccess({required this.payment});

  final Map<String, dynamic> payment;

  @override
  List<Object?> get props => [payment];
}

class PaymentListLoaded extends PaymentState {
  const PaymentListLoaded({required this.payments});

  final List<dynamic> payments;

  @override
  List<Object?> get props => [payments];
}

class PaymentError extends PaymentState {
  const PaymentError({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

// =============================================================================
// BLoC
// =============================================================================

class PaymentBloc extends Bloc<PaymentEvent, PaymentState> {
  PaymentBloc({required PaymentRepository paymentRepository})
      : _paymentRepo = paymentRepository,
        super(const PaymentInitial()) {
    on<PaymentSepaRequested>(_onSepa);
    on<PaymentInternalRequested>(_onInternal);
    on<PaymentListRequested>(_onList);
  }

  final PaymentRepository _paymentRepo;

  Future<void> _onSepa(
    PaymentSepaRequested event,
    Emitter<PaymentState> emit,
  ) async {
    emit(const PaymentLoading());
    try {
      final result = await _paymentRepo.createSepaPayment(
        senderAccountId: event.senderAccountId,
        recipientIban: event.recipientIban,
        recipientName: event.recipientName,
        amount: event.amount,
        currency: event.currency,
        reference: event.reference,
      );
      emit(PaymentSuccess(payment: result));
    } catch (e) {
      emit(PaymentError(message: e.toString()));
    }
  }

  Future<void> _onInternal(
    PaymentInternalRequested event,
    Emitter<PaymentState> emit,
  ) async {
    emit(const PaymentLoading());
    try {
      final result = await _paymentRepo.createInternalPayment(
        senderAccountId: event.senderAccountId,
        recipientAccountId: event.recipientAccountId,
        amount: event.amount,
        currency: event.currency,
        reference: event.reference,
      );
      emit(PaymentSuccess(payment: result));
    } catch (e) {
      emit(PaymentError(message: e.toString()));
    }
  }

  Future<void> _onList(
    PaymentListRequested event,
    Emitter<PaymentState> emit,
  ) async {
    emit(const PaymentLoading());
    try {
      final payments =
          await _paymentRepo.listPayments(accountId: event.accountId);
      emit(PaymentListLoaded(payments: payments));
    } catch (e) {
      emit(PaymentError(message: e.toString()));
    }
  }
}
