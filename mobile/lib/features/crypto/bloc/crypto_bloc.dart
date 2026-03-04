import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/repositories/crypto_repository.dart';
import '../../../models/crypto.dart';

// =============================================================================
// Events
// =============================================================================

abstract class CryptoEvent extends Equatable {
  const CryptoEvent();

  @override
  List<Object?> get props => [];
}

class CryptoLoadRequested extends CryptoEvent {
  const CryptoLoadRequested();
}

class CryptoRefreshRequested extends CryptoEvent {
  const CryptoRefreshRequested();
}

class CryptoBuyRequested extends CryptoEvent {
  const CryptoBuyRequested({
    required this.symbol,
    required this.amount,
  });

  final String symbol;
  final String amount;

  @override
  List<Object?> get props => [symbol, amount];
}

class CryptoSellRequested extends CryptoEvent {
  const CryptoSellRequested({
    required this.symbol,
    required this.amount,
  });

  final String symbol;
  final String amount;

  @override
  List<Object?> get props => [symbol, amount];
}

class CryptoSendRequested extends CryptoEvent {
  const CryptoSendRequested({
    required this.symbol,
    required this.amount,
    required this.recipientAddress,
  });

  final String symbol;
  final String amount;
  final String recipientAddress;

  @override
  List<Object?> get props => [symbol, amount, recipientAddress];
}

// =============================================================================
// States
// =============================================================================

abstract class CryptoState extends Equatable {
  const CryptoState();

  @override
  List<Object?> get props => [];
}

class CryptoInitial extends CryptoState {
  const CryptoInitial();
}

class CryptoLoading extends CryptoState {
  const CryptoLoading();
}

class CryptoLoaded extends CryptoState {
  const CryptoLoaded({
    required this.wallet,
    required this.prices,
    required this.transactions,
    required this.totalValueEur,
  });

  final CryptoWallet wallet;
  final List<CryptoPrice> prices;
  final List<CryptoTransaction> transactions;
  final String totalValueEur;

  @override
  List<Object?> get props => [wallet, prices, transactions, totalValueEur];
}

class CryptoActionSuccess extends CryptoState {
  const CryptoActionSuccess({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

class CryptoError extends CryptoState {
  const CryptoError({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

// =============================================================================
// BLoC
// =============================================================================

class CryptoBloc extends Bloc<CryptoEvent, CryptoState> {
  CryptoBloc({required CryptoRepository cryptoRepository})
      : _cryptoRepo = cryptoRepository,
        super(const CryptoInitial()) {
    on<CryptoLoadRequested>(_onLoad);
    on<CryptoRefreshRequested>(_onRefresh);
    on<CryptoBuyRequested>(_onBuy);
    on<CryptoSellRequested>(_onSell);
    on<CryptoSendRequested>(_onSend);
  }

  final CryptoRepository _cryptoRepo;

  Future<void> _onLoad(
    CryptoLoadRequested event,
    Emitter<CryptoState> emit,
  ) async {
    emit(const CryptoLoading());
    await _loadData(emit);
  }

  Future<void> _onRefresh(
    CryptoRefreshRequested event,
    Emitter<CryptoState> emit,
  ) async {
    await _loadData(emit);
  }

  Future<void> _loadData(Emitter<CryptoState> emit) async {
    try {
      final wallet = await _cryptoRepo.getWallet();
      final prices = await _cryptoRepo.getPrices();
      final transactions = await _cryptoRepo.getTransactions(limit: 10);

      // Calculate total value in EUR from balances.
      double total = 0;
      for (final balance in wallet.balances) {
        final valueEur =
            double.tryParse(balance.valueEur ?? '0') ?? 0;
        total += valueEur;
      }

      emit(CryptoLoaded(
        wallet: wallet,
        prices: prices,
        transactions: transactions,
        totalValueEur: total.toStringAsFixed(2),
      ));
    } catch (e) {
      emit(CryptoError(message: e.toString()));
    }
  }

  Future<void> _onBuy(
    CryptoBuyRequested event,
    Emitter<CryptoState> emit,
  ) async {
    emit(const CryptoLoading());
    try {
      final quote = await _cryptoRepo.getQuote(
        action: 'buy',
        symbol: event.symbol,
        amount: event.amount,
      );
      await _cryptoRepo.buyCrypto(quote.id);
      emit(CryptoActionSuccess(
        message: 'Bought ${event.amount} EUR of ${event.symbol}',
      ));
      add(const CryptoLoadRequested());
    } catch (e) {
      emit(CryptoError(message: e.toString()));
    }
  }

  Future<void> _onSell(
    CryptoSellRequested event,
    Emitter<CryptoState> emit,
  ) async {
    emit(const CryptoLoading());
    try {
      final quote = await _cryptoRepo.getQuote(
        action: 'sell',
        symbol: event.symbol,
        amount: event.amount,
      );
      await _cryptoRepo.sellCrypto(quote.id);
      emit(CryptoActionSuccess(
        message: 'Sold ${event.amount} of ${event.symbol}',
      ));
      add(const CryptoLoadRequested());
    } catch (e) {
      emit(CryptoError(message: e.toString()));
    }
  }

  Future<void> _onSend(
    CryptoSendRequested event,
    Emitter<CryptoState> emit,
  ) async {
    emit(const CryptoLoading());
    try {
      await _cryptoRepo.sendCrypto(
        symbol: event.symbol,
        amount: event.amount,
        recipientAddress: event.recipientAddress,
      );
      emit(CryptoActionSuccess(
        message: 'Sent ${event.amount} ${event.symbol}',
      ));
      add(const CryptoLoadRequested());
    } catch (e) {
      emit(CryptoError(message: e.toString()));
    }
  }
}
