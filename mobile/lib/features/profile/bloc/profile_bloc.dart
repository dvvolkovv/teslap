import 'package:equatable/equatable.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

import '../../../core/repositories/account_repository.dart';

// Events
abstract class ProfileEvent extends Equatable {
  const ProfileEvent();
  @override
  List<Object?> get props => [];
}

class ProfileLoadRequested extends ProfileEvent {
  const ProfileLoadRequested();
}

// States
abstract class ProfileState extends Equatable {
  const ProfileState();
  @override
  List<Object?> get props => [];
}

class ProfileInitial extends ProfileState {}

class ProfileLoading extends ProfileState {}

class ProfileLoaded extends ProfileState {
  const ProfileLoaded({
    required this.firstName,
    required this.lastName,
    required this.email,
    required this.phone,
    required this.tier,
    required this.kycStatus,
  });

  final String firstName;
  final String lastName;
  final String email;
  final String phone;
  final String tier;
  final String kycStatus;

  String get fullName =>
      '$firstName $lastName'.trim().isEmpty ? 'User' : '$firstName $lastName'.trim();

  String get initials {
    if (firstName.isNotEmpty && lastName.isNotEmpty) {
      return '${firstName[0]}${lastName[0]}'.toUpperCase();
    }
    if (email.isNotEmpty) return email[0].toUpperCase();
    return 'U';
  }

  @override
  List<Object?> get props =>
      [firstName, lastName, email, phone, tier, kycStatus];
}

class ProfileError extends ProfileState {
  const ProfileError(this.message);
  final String message;
  @override
  List<Object?> get props => [message];
}

// BLoC
class ProfileBloc extends Bloc<ProfileEvent, ProfileState> {
  ProfileBloc({required AccountRepository accountRepository})
      : _repo = accountRepository,
        super(ProfileInitial()) {
    on<ProfileLoadRequested>(_onLoad);
  }

  final AccountRepository _repo;

  Future<void> _onLoad(
    ProfileLoadRequested event,
    Emitter<ProfileState> emit,
  ) async {
    emit(ProfileLoading());
    try {
      final profile = await _repo.getProfile();
      emit(ProfileLoaded(
        firstName: profile['first_name']?.toString() ?? '',
        lastName: profile['last_name']?.toString() ?? '',
        email: profile['email']?.toString() ?? '',
        phone: profile['phone']?.toString() ?? '',
        tier: profile['tier']?.toString() ?? 'standard',
        kycStatus: profile['kyc_status']?.toString() ?? 'pending',
      ));
    } catch (e) {
      emit(ProfileError(e.toString()));
    }
  }
}
