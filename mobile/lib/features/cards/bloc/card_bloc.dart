import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/repositories/card_repository.dart';

// =============================================================================
// Events
// =============================================================================

abstract class CardEvent extends Equatable {
  const CardEvent();

  @override
  List<Object?> get props => [];
}

class CardListRequested extends CardEvent {
  const CardListRequested({this.accountId});

  final String? accountId;

  @override
  List<Object?> get props => [accountId];
}

class CardFreezeRequested extends CardEvent {
  const CardFreezeRequested({required this.cardId});

  final String cardId;

  @override
  List<Object?> get props => [cardId];
}

class CardUnfreezeRequested extends CardEvent {
  const CardUnfreezeRequested({required this.cardId});

  final String cardId;

  @override
  List<Object?> get props => [cardId];
}

class CardBlockRequested extends CardEvent {
  const CardBlockRequested({required this.cardId});

  final String cardId;

  @override
  List<Object?> get props => [cardId];
}

class CardIssueVirtualRequested extends CardEvent {
  const CardIssueVirtualRequested({
    required this.accountId,
    required this.cardholderName,
  });

  final String accountId;
  final String cardholderName;

  @override
  List<Object?> get props => [accountId, cardholderName];
}

class CardTransactionsRequested extends CardEvent {
  const CardTransactionsRequested({required this.cardId});

  final String cardId;

  @override
  List<Object?> get props => [cardId];
}

// =============================================================================
// States
// =============================================================================

abstract class CardState extends Equatable {
  const CardState();

  @override
  List<Object?> get props => [];
}

class CardInitial extends CardState {
  const CardInitial();
}

class CardLoading extends CardState {
  const CardLoading();
}

class CardListLoaded extends CardState {
  const CardListLoaded({required this.cards});

  final List<dynamic> cards;

  @override
  List<Object?> get props => [cards];
}

class CardActionSuccess extends CardState {
  const CardActionSuccess({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

class CardTransactionsLoaded extends CardState {
  const CardTransactionsLoaded({required this.transactions});

  final List<dynamic> transactions;

  @override
  List<Object?> get props => [transactions];
}

class CardError extends CardState {
  const CardError({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

// =============================================================================
// BLoC
// =============================================================================

class CardBloc extends Bloc<CardEvent, CardState> {
  CardBloc({required CardRepository cardRepository})
      : _cardRepo = cardRepository,
        super(const CardInitial()) {
    on<CardListRequested>(_onList);
    on<CardFreezeRequested>(_onFreeze);
    on<CardUnfreezeRequested>(_onUnfreeze);
    on<CardBlockRequested>(_onBlock);
    on<CardIssueVirtualRequested>(_onIssueVirtual);
    on<CardTransactionsRequested>(_onTransactions);
  }

  final CardRepository _cardRepo;

  Future<void> _onList(
    CardListRequested event,
    Emitter<CardState> emit,
  ) async {
    emit(const CardLoading());
    try {
      final cards = await _cardRepo.listCards(accountId: event.accountId);
      emit(CardListLoaded(cards: cards));
    } catch (e) {
      emit(CardError(message: e.toString()));
    }
  }

  Future<void> _onFreeze(
    CardFreezeRequested event,
    Emitter<CardState> emit,
  ) async {
    emit(const CardLoading());
    try {
      await _cardRepo.freezeCard(event.cardId);
      emit(const CardActionSuccess(message: 'Card frozen successfully'));
      // Reload card list.
      add(const CardListRequested());
    } catch (e) {
      emit(CardError(message: e.toString()));
    }
  }

  Future<void> _onUnfreeze(
    CardUnfreezeRequested event,
    Emitter<CardState> emit,
  ) async {
    emit(const CardLoading());
    try {
      await _cardRepo.unfreezeCard(event.cardId);
      emit(const CardActionSuccess(message: 'Card unfrozen successfully'));
      add(const CardListRequested());
    } catch (e) {
      emit(CardError(message: e.toString()));
    }
  }

  Future<void> _onBlock(
    CardBlockRequested event,
    Emitter<CardState> emit,
  ) async {
    emit(const CardLoading());
    try {
      await _cardRepo.blockCard(event.cardId);
      emit(const CardActionSuccess(message: 'Card blocked permanently'));
      add(const CardListRequested());
    } catch (e) {
      emit(CardError(message: e.toString()));
    }
  }

  Future<void> _onIssueVirtual(
    CardIssueVirtualRequested event,
    Emitter<CardState> emit,
  ) async {
    emit(const CardLoading());
    try {
      await _cardRepo.issueVirtualCard(
        accountId: event.accountId,
        cardholderName: event.cardholderName,
      );
      emit(const CardActionSuccess(message: 'Virtual card issued'));
      add(const CardListRequested());
    } catch (e) {
      emit(CardError(message: e.toString()));
    }
  }

  Future<void> _onTransactions(
    CardTransactionsRequested event,
    Emitter<CardState> emit,
  ) async {
    emit(const CardLoading());
    try {
      final txs = await _cardRepo.getCardTransactions(event.cardId);
      emit(CardTransactionsLoaded(transactions: txs));
    } catch (e) {
      emit(CardError(message: e.toString()));
    }
  }
}
