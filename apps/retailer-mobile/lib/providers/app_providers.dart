/// OmniRoute Ecosystem - State Management
/// Riverpod providers for authentication, user state, and dashboard data

import 'dart:async';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';
import 'package:omniroute_ecosystem/core/network/api_client.dart';
import 'package:omniroute_ecosystem/models/models.dart';

// ============================================================================
// AUTH STATE
// ============================================================================

enum AuthStatus { initial, authenticated, unauthenticated, loading }

class AuthState {
  final AuthStatus status;
  final User? user;
  final String? error;
  final bool isOnboardingComplete;

  const AuthState({
    this.status = AuthStatus.initial,
    this.user,
    this.error,
    this.isOnboardingComplete = false,
  });

  AuthState copyWith({
    AuthStatus? status,
    User? user,
    String? error,
    bool? isOnboardingComplete,
  }) {
    return AuthState(
      status: status ?? this.status,
      user: user ?? this.user,
      error: error,
      isOnboardingComplete: isOnboardingComplete ?? this.isOnboardingComplete,
    );
  }

  bool get isAuthenticated => status == AuthStatus.authenticated;
  bool get isLoading => status == AuthStatus.loading;
}

class AuthNotifier extends StateNotifier<AuthState> {
  final ApiClient _apiClient;
  final FlutterSecureStorage _secureStorage;

  AuthNotifier(this._apiClient, this._secureStorage) : super(const AuthState()) {
    _checkAuthStatus();
  }

  Future<void> _checkAuthStatus() async {
    state = state.copyWith(status: AuthStatus.loading);

    try {
      final token = await _secureStorage.read(key: StorageKeys.accessToken);
      if (token != null) {
        // Validate token and get user
        final response = await _apiClient.get<User>(
          '/auth/me',
          fromJson: (json) => User.fromJson(json['user']),
        );

        if (response.isSuccess && response.data != null) {
          state = AuthState(
            status: AuthStatus.authenticated,
            user: response.data,
            isOnboardingComplete: response.data!.isOnboardingComplete,
          );
          return;
        }
      }
      state = state.copyWith(status: AuthStatus.unauthenticated);
    } catch (e) {
      state = state.copyWith(status: AuthStatus.unauthenticated);
    }
  }

  Future<bool> login({
    required String email,
    required String password,
  }) async {
    state = state.copyWith(status: AuthStatus.loading, error: null);

    final response = await _apiClient.post<AuthResponse>(
      ApiEndpoints.login,
      data: {'email': email, 'password': password},
      fromJson: (json) => AuthResponse.fromJson(json),
    );

    return response.when(
      success: (data) async {
        await _secureStorage.write(
          key: StorageKeys.accessToken,
          value: data.accessToken,
        );
        await _secureStorage.write(
          key: StorageKeys.refreshToken,
          value: data.refreshToken,
        );
        await _secureStorage.write(
          key: StorageKeys.userId,
          value: data.user.id,
        );

        state = AuthState(
          status: AuthStatus.authenticated,
          user: data.user,
          isOnboardingComplete: data.user.isOnboardingComplete,
        );
        return true;
      },
      error: (message) {
        state = state.copyWith(
          status: AuthStatus.unauthenticated,
          error: message,
        );
        return false;
      },
    );
  }

  Future<bool> register({
    required String fullName,
    required String email,
    required String phone,
    required String password,
    required ParticipantType participantType,
  }) async {
    state = state.copyWith(status: AuthStatus.loading, error: null);

    final response = await _apiClient.post<AuthResponse>(
      ApiEndpoints.register,
      data: {
        'full_name': fullName,
        'email': email,
        'phone': phone,
        'password': password,
        'participant_type': participantType.code,
      },
      fromJson: (json) => AuthResponse.fromJson(json),
    );

    return response.when(
      success: (data) async {
        await _secureStorage.write(
          key: StorageKeys.accessToken,
          value: data.accessToken,
        );
        await _secureStorage.write(
          key: StorageKeys.refreshToken,
          value: data.refreshToken,
        );
        await _secureStorage.write(
          key: StorageKeys.userId,
          value: data.user.id,
        );

        state = AuthState(
          status: AuthStatus.authenticated,
          user: data.user,
          isOnboardingComplete: false,
        );
        return true;
      },
      error: (message) {
        state = state.copyWith(
          status: AuthStatus.unauthenticated,
          error: message,
        );
        return false;
      },
    );
  }

  Future<bool> verifyOtp(String otp) async {
    state = state.copyWith(status: AuthStatus.loading, error: null);

    final response = await _apiClient.post(
      ApiEndpoints.verifyOtp,
      data: {'otp': otp},
    );

    return response.when(
      success: (_) {
        if (state.user != null) {
          state = state.copyWith(
            user: state.user!.copyWith(isVerified: true),
          );
        }
        return true;
      },
      error: (message) {
        state = state.copyWith(error: message);
        return false;
      },
    );
  }

  Future<void> logout() async {
    await _secureStorage.deleteAll();
    state = const AuthState(status: AuthStatus.unauthenticated);
  }

  void updateUser(User user) {
    state = state.copyWith(user: user);
  }

  void setOnboardingComplete() {
    state = state.copyWith(
      isOnboardingComplete: true,
      user: state.user?.copyWith(isOnboardingComplete: true),
    );
  }
}

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final secureStorage = ref.watch(secureStorageProvider);
  return AuthNotifier(apiClient, secureStorage);
});

// ============================================================================
// DASHBOARD PROVIDERS
// ============================================================================

final dashboardStatsProvider = FutureProvider.autoDispose<DashboardStats>((ref) async {
  final apiClient = ref.watch(apiClientProvider);
  
  final response = await apiClient.get<DashboardStats>(
    ApiEndpoints.dashboardStats,
    fromJson: (json) => DashboardStats.fromJson(json),
  );

  return response.when(
    success: (data) => data,
    error: (message) => throw Exception(message),
  );
});

// ============================================================================
// ORDERS PROVIDERS
// ============================================================================

class OrdersState {
  final List<Order> orders;
  final bool isLoading;
  final String? error;
  final bool hasMore;
  final int currentPage;

  const OrdersState({
    this.orders = const [],
    this.isLoading = false,
    this.error,
    this.hasMore = true,
    this.currentPage = 1,
  });

  OrdersState copyWith({
    List<Order>? orders,
    bool? isLoading,
    String? error,
    bool? hasMore,
    int? currentPage,
  }) {
    return OrdersState(
      orders: orders ?? this.orders,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
    );
  }
}

class OrdersNotifier extends StateNotifier<OrdersState> {
  final ApiClient _apiClient;

  OrdersNotifier(this._apiClient) : super(const OrdersState());

  Future<void> loadOrders({bool refresh = false}) async {
    if (state.isLoading) return;
    if (!refresh && !state.hasMore) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    final response = await _apiClient.get(
      ApiEndpoints.orders,
      queryParameters: {
        'page': state.currentPage,
        'limit': 20,
      },
    );

    response.when(
      success: (data) {
        final orders = (data['orders'] as List)
            .map((e) => Order.fromJson(e))
            .toList();

        state = state.copyWith(
          orders: refresh ? orders : [...state.orders, ...orders],
          isLoading: false,
          hasMore: orders.length >= 20,
          currentPage: state.currentPage + 1,
        );
      },
      error: (message) {
        state = state.copyWith(
          isLoading: false,
          error: message,
        );
      },
    );
  }

  Future<Order?> getOrderDetails(String orderId) async {
    final response = await _apiClient.get<Order>(
      ApiEndpoints.orderDetails.replaceAll('{id}', orderId),
      fromJson: (json) => Order.fromJson(json['order']),
    );

    return response.when(
      success: (data) => data,
      error: (_) => null,
    );
  }

  Future<bool> updateOrderStatus(String orderId, String status) async {
    final response = await _apiClient.put(
      ApiEndpoints.updateOrderStatus.replaceAll('{id}', orderId),
      data: {'status': status},
    );

    if (response.isSuccess) {
      // Update local state
      state = state.copyWith(
        orders: state.orders.map((o) {
          if (o.id == orderId) {
            return o.copyWith(status: status);
          }
          return o;
        }).toList(),
      );
    }

    return response.isSuccess;
  }
}

final ordersProvider = StateNotifierProvider<OrdersNotifier, OrdersState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return OrdersNotifier(apiClient);
});

// ============================================================================
// INVENTORY PROVIDERS
// ============================================================================

final inventoryProvider = FutureProvider.autoDispose<List<InventoryItem>>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get(
    ApiEndpoints.inventory,
  );

  return response.when(
    success: (data) {
      return (data['items'] as List)
          .map((e) => InventoryItem.fromJson(e))
          .toList();
    },
    error: (message) => throw Exception(message),
  );
});

final stockAlertsProvider = FutureProvider.autoDispose<List<StockAlert>>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get(
    ApiEndpoints.stockAlerts,
  );

  return response.when(
    success: (data) {
      return (data['alerts'] as List)
          .map((e) => StockAlert.fromJson(e))
          .toList();
    },
    error: (message) => throw Exception(message),
  );
});

// ============================================================================
// WALLET PROVIDERS
// ============================================================================

final walletProvider = FutureProvider.autoDispose<Wallet>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<Wallet>(
    ApiEndpoints.wallet,
    fromJson: (json) => Wallet.fromJson(json['wallet']),
  );

  return response.when(
    success: (data) => data,
    error: (message) => throw Exception(message),
  );
});

final transactionsProvider = FutureProvider.autoDispose<List<Transaction>>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get(
    ApiEndpoints.transactions,
    queryParameters: {'limit': 50},
  );

  return response.when(
    success: (data) {
      return (data['transactions'] as List)
          .map((e) => Transaction.fromJson(e))
          .toList();
    },
    error: (message) => throw Exception(message),
  );
});

// ============================================================================
// DELIVERIES PROVIDERS
// ============================================================================

final deliveriesProvider = FutureProvider.autoDispose<List<Delivery>>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get(
    ApiEndpoints.deliveries,
  );

  return response.when(
    success: (data) {
      return (data['deliveries'] as List)
          .map((e) => Delivery.fromJson(e))
          .toList();
    },
    error: (message) => throw Exception(message),
  );
});

// ============================================================================
// NOTIFICATIONS PROVIDERS
// ============================================================================

final notificationsProvider = FutureProvider.autoDispose<List<AppNotification>>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get(
    ApiEndpoints.notifications,
  );

  return response.when(
    success: (data) {
      return (data['notifications'] as List)
          .map((e) => AppNotification.fromJson(e))
          .toList();
    },
    error: (message) => throw Exception(message),
  );
});

final unreadNotificationCountProvider = Provider<int>((ref) {
  final notifications = ref.watch(notificationsProvider);
  return notifications.when(
    data: (data) => data.where((n) => !n.isRead).length,
    loading: () => 0,
    error: (_, __) => 0,
  );
});

// ============================================================================
// ANALYTICS PROVIDERS
// ============================================================================

final salesAnalyticsProvider = FutureProvider.autoDispose<SalesAnalytics>((ref) async {
  final apiClient = ref.watch(apiClientProvider);

  final response = await apiClient.get<SalesAnalytics>(
    ApiEndpoints.salesAnalytics,
    fromJson: (json) => SalesAnalytics.fromJson(json),
  );

  return response.when(
    success: (data) => data,
    error: (message) => throw Exception(message),
  );
});

// ============================================================================
// SELECTED PARTICIPANT TYPE
// ============================================================================

final selectedParticipantTypeProvider = StateProvider<ParticipantType?>((ref) => null);

// ============================================================================
// NAVIGATION INDEX
// ============================================================================

final bottomNavIndexProvider = StateProvider<int>((ref) => 0);
