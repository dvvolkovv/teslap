import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/repositories/account_repository.dart';

// =============================================================================
// Events
// =============================================================================

abstract class HomeEvent extends Equatable {
  const HomeEvent();

  @override
  List<Object?> get props => [];
}

class HomeLoadRequested extends HomeEvent {
  const HomeLoadRequested();
}

class HomeRefreshRequested extends HomeEvent {
  const HomeRefreshRequested();
}

// =============================================================================
// States
// =============================================================================

abstract class HomeState extends Equatable {
  const HomeState();

  @override
  List<Object?> get props => [];
}

class HomeInitial extends HomeState {
  const HomeInitial();
}

class HomeLoading extends HomeState {
  const HomeLoading();
}

class HomeLoaded extends HomeState {
  const HomeLoaded({
    required this.profile,
    required this.accounts,
    required this.transactions,
    required this.totalBalanceEur,
  });

  final Map<String, dynamic> profile;
  final List<dynamic> accounts;
  final List<dynamic> transactions;
  final String totalBalanceEur;

  @override
  List<Object?> get props => [profile, accounts, transactions, totalBalanceEur];
}

class HomeError extends HomeState {
  const HomeError({required this.message});

  final String message;

  @override
  List<Object?> get props => [message];
}

// =============================================================================
// BLoC
// =============================================================================

class HomeBloc extends Bloc<HomeEvent, HomeState> {
  HomeBloc({required AccountRepository accountRepository})
      : _accountRepo = accountRepository,
        super(const HomeInitial()) {
    on<HomeLoadRequested>(_onLoad);
    on<HomeRefreshRequested>(_onRefresh);
  }

  final AccountRepository _accountRepo;

  Future<void> _onLoad(
    HomeLoadRequested event,
    Emitter<HomeState> emit,
  ) async {
    emit(const HomeLoading());
    await _loadData(emit);
  }

  Future<void> _onRefresh(
    HomeRefreshRequested event,
    Emitter<HomeState> emit,
  ) async {
    await _loadData(emit);
  }

  Future<void> _loadData(Emitter<HomeState> emit) async {
    try {
      final profile = await _accountRepo.getProfile();
      final accounts = await _accountRepo.getAccounts();

      // Load transactions from the first account if available.
      List<dynamic> transactions = [];
      if (accounts.isNotEmpty) {
        final firstAccount = accounts[0] as Map<String, dynamic>;
        final accountId = firstAccount['id'] as String;
        transactions =
            await _accountRepo.getTransactions(accountId, limit: 5);
      }

      // Calculate total balance from sub-accounts.
      double total = 0;
      for (final acc in accounts) {
        final map = acc as Map<String, dynamic>;
        final subAccounts = map['sub_accounts'] as List<dynamic>? ?? [];
        for (final sub in subAccounts) {
          final subMap = sub as Map<String, dynamic>;
          final balance = subMap['balance'] as Map<String, dynamic>?;
          if (balance != null) {
            final available = double.tryParse(
                    balance['available']?.toString() ?? '0') ??
                0;
            total += available;
          }
        }
      }

      emit(HomeLoaded(
        profile: profile,
        accounts: accounts,
        transactions: transactions,
        totalBalanceEur: total.toStringAsFixed(2),
      ));
    } catch (e) {
      emit(HomeError(message: e.toString()));
    }
  }
}
