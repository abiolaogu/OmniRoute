/// OmniRoute Ecosystem - Router Configuration
/// GoRouter setup for multi-participant navigation with auth guards

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

// Feature imports - these would be the actual screen files
import 'package:omniroute_ecosystem/features/splash/splash_screen.dart';
import 'package:omniroute_ecosystem/features/onboarding/screens/welcome_screen.dart';
import 'package:omniroute_ecosystem/features/onboarding/screens/participant_selection_screen.dart';
import 'package:omniroute_ecosystem/features/auth/screens/login_screen.dart';
import 'package:omniroute_ecosystem/features/auth/screens/register_screen.dart';
import 'package:omniroute_ecosystem/features/auth/screens/otp_verification_screen.dart';
import 'package:omniroute_ecosystem/features/onboarding/screens/kyc_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/common/main_dashboard_shell.dart';
import 'package:omniroute_ecosystem/features/dashboard/bank/bank_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/logistics/logistics_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/warehouse/warehouse_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/manufacturer/manufacturer_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/retailer/retailer_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/wholesaler/wholesaler_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/investor/investor_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/entrepreneur/entrepreneur_dashboard_screen.dart';
import 'package:omniroute_ecosystem/features/dashboard/ecommerce/ecommerce_dashboard_screen.dart';

// ============================================================================
// ROUTE NAMES
// ============================================================================

class RouteNames {
  // Auth & Onboarding
  static const String splash = 'splash';
  static const String welcome = 'welcome';
  static const String participantSelection = 'participant-selection';
  static const String login = 'login';
  static const String register = 'register';
  static const String otpVerification = 'otp-verification';
  static const String kyc = 'kyc';

  // Dashboard Routes
  static const String dashboard = 'dashboard';
  static const String bankDashboard = 'bank-dashboard';
  static const String logisticsDashboard = 'logistics-dashboard';
  static const String warehouseDashboard = 'warehouse-dashboard';
  static const String manufacturerDashboard = 'manufacturer-dashboard';
  static const String retailerDashboard = 'retailer-dashboard';
  static const String wholesalerDashboard = 'wholesaler-dashboard';
  static const String investorDashboard = 'investor-dashboard';
  static const String entrepreneurDashboard = 'entrepreneur-dashboard';
  static const String ecommerceDashboard = 'ecommerce-dashboard';

  // Feature Routes
  static const String orders = 'orders';
  static const String orderDetails = 'order-details';
  static const String inventory = 'inventory';
  static const String products = 'products';
  static const String productDetails = 'product-details';
  static const String analytics = 'analytics';
  static const String wallet = 'wallet';
  static const String transactions = 'transactions';
  static const String deliveries = 'deliveries';
  static const String deliveryTracking = 'delivery-tracking';
  static const String fleet = 'fleet';
  static const String warehouses = 'warehouses';
  static const String settings = 'settings';
  static const String profile = 'profile';
  static const String notifications = 'notifications';
}

// ============================================================================
// ROUTE PATHS
// ============================================================================

class RoutePaths {
  static const String splash = '/';
  static const String welcome = '/welcome';
  static const String participantSelection = '/participant-selection';
  static const String login = '/login';
  static const String register = '/register';
  static const String otpVerification = '/otp-verification';
  static const String kyc = '/kyc';

  static const String dashboard = '/dashboard';
  static const String orders = '/orders';
  static const String orderDetails = '/orders/:id';
  static const String inventory = '/inventory';
  static const String products = '/products';
  static const String productDetails = '/products/:id';
  static const String analytics = '/analytics';
  static const String wallet = '/wallet';
  static const String transactions = '/transactions';
  static const String deliveries = '/deliveries';
  static const String deliveryTracking = '/deliveries/:id/tracking';
  static const String fleet = '/fleet';
  static const String warehouses = '/warehouses';
  static const String settings = '/settings';
  static const String profile = '/profile';
  static const String notifications = '/notifications';
}

// ============================================================================
// ROUTER PROVIDER
// ============================================================================

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authProvider);

  return GoRouter(
    initialLocation: RoutePaths.splash,
    debugLogDiagnostics: true,
    refreshListenable: _RouterRefreshStream(ref),
    redirect: (context, state) {
      final isAuthenticated = authState.isAuthenticated;
      final isOnboardingComplete = authState.isOnboardingComplete;
      final currentPath = state.matchedLocation;

      // Public routes that don't require auth
      final publicRoutes = [
        RoutePaths.splash,
        RoutePaths.welcome,
        RoutePaths.participantSelection,
        RoutePaths.login,
        RoutePaths.register,
        RoutePaths.otpVerification,
      ];

      final isPublicRoute = publicRoutes.contains(currentPath);

      // If loading, stay on splash
      if (authState.isLoading && currentPath == RoutePaths.splash) {
        return null;
      }

      // If not authenticated and trying to access protected route
      if (!isAuthenticated && !isPublicRoute) {
        return RoutePaths.welcome;
      }

      // If authenticated but not onboarded
      if (isAuthenticated && !isOnboardingComplete && currentPath != RoutePaths.kyc) {
        return RoutePaths.kyc;
      }

      // If authenticated and onboarded, but on public route
      if (isAuthenticated && isOnboardingComplete && isPublicRoute) {
        return RoutePaths.dashboard;
      }

      return null;
    },
    routes: [
      // ========================================================================
      // AUTH & ONBOARDING ROUTES
      // ========================================================================
      GoRoute(
        path: RoutePaths.splash,
        name: RouteNames.splash,
        builder: (context, state) => const SplashScreen(),
      ),
      GoRoute(
        path: RoutePaths.welcome,
        name: RouteNames.welcome,
        builder: (context, state) => const WelcomeScreen(),
      ),
      GoRoute(
        path: RoutePaths.participantSelection,
        name: RouteNames.participantSelection,
        builder: (context, state) => const ParticipantSelectionScreen(),
      ),
      GoRoute(
        path: RoutePaths.login,
        name: RouteNames.login,
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: RoutePaths.register,
        name: RouteNames.register,
        builder: (context, state) {
          final participantType = state.extra as ParticipantType?;
          return RegisterScreen(participantType: participantType);
        },
      ),
      GoRoute(
        path: RoutePaths.otpVerification,
        name: RouteNames.otpVerification,
        builder: (context, state) {
          final phone = state.extra as String? ?? '';
          return OtpVerificationScreen(phone: phone);
        },
      ),
      GoRoute(
        path: RoutePaths.kyc,
        name: RouteNames.kyc,
        builder: (context, state) => const KycScreen(),
      ),

      // ========================================================================
      // MAIN DASHBOARD SHELL WITH NESTED ROUTES
      // ========================================================================
      ShellRoute(
        builder: (context, state, child) {
          return MainDashboardShell(child: child);
        },
        routes: [
          // Main dashboard - routes to participant-specific dashboard
          GoRoute(
            path: RoutePaths.dashboard,
            name: RouteNames.dashboard,
            builder: (context, state) => const _ParticipantDashboardRouter(),
            routes: [
              // Participant-specific dashboard sub-routes
              GoRoute(
                path: 'bank',
                name: RouteNames.bankDashboard,
                builder: (context, state) => const BankDashboardScreen(),
              ),
              GoRoute(
                path: 'logistics',
                name: RouteNames.logisticsDashboard,
                builder: (context, state) => const LogisticsDashboardScreen(),
              ),
              GoRoute(
                path: 'warehouse',
                name: RouteNames.warehouseDashboard,
                builder: (context, state) => const WarehouseDashboardScreen(),
              ),
              GoRoute(
                path: 'manufacturer',
                name: RouteNames.manufacturerDashboard,
                builder: (context, state) => const ManufacturerDashboardScreen(),
              ),
              GoRoute(
                path: 'retailer',
                name: RouteNames.retailerDashboard,
                builder: (context, state) => const RetailerDashboardScreen(),
              ),
              GoRoute(
                path: 'wholesaler',
                name: RouteNames.wholesalerDashboard,
                builder: (context, state) => const WholesalerDashboardScreen(),
              ),
              GoRoute(
                path: 'investor',
                name: RouteNames.investorDashboard,
                builder: (context, state) => const InvestorDashboardScreen(),
              ),
              GoRoute(
                path: 'entrepreneur',
                name: RouteNames.entrepreneurDashboard,
                builder: (context, state) => const EntrepreneurDashboardScreen(),
              ),
              GoRoute(
                path: 'ecommerce',
                name: RouteNames.ecommerceDashboard,
                builder: (context, state) => const EcommerceDashboardScreen(),
              ),
            ],
          ),

          // Orders
          GoRoute(
            path: RoutePaths.orders,
            name: RouteNames.orders,
            builder: (context, state) => const OrdersScreen(),
            routes: [
              GoRoute(
                path: ':id',
                name: RouteNames.orderDetails,
                builder: (context, state) {
                  final orderId = state.pathParameters['id']!;
                  return OrderDetailsScreen(orderId: orderId);
                },
              ),
            ],
          ),

          // Inventory
          GoRoute(
            path: RoutePaths.inventory,
            name: RouteNames.inventory,
            builder: (context, state) => const InventoryScreen(),
          ),

          // Products
          GoRoute(
            path: RoutePaths.products,
            name: RouteNames.products,
            builder: (context, state) => const ProductsScreen(),
            routes: [
              GoRoute(
                path: ':id',
                name: RouteNames.productDetails,
                builder: (context, state) {
                  final productId = state.pathParameters['id']!;
                  return ProductDetailsScreen(productId: productId);
                },
              ),
            ],
          ),

          // Analytics
          GoRoute(
            path: RoutePaths.analytics,
            name: RouteNames.analytics,
            builder: (context, state) => const AnalyticsScreen(),
          ),

          // Wallet & Transactions
          GoRoute(
            path: RoutePaths.wallet,
            name: RouteNames.wallet,
            builder: (context, state) => const WalletScreen(),
          ),
          GoRoute(
            path: RoutePaths.transactions,
            name: RouteNames.transactions,
            builder: (context, state) => const TransactionsScreen(),
          ),

          // Deliveries
          GoRoute(
            path: RoutePaths.deliveries,
            name: RouteNames.deliveries,
            builder: (context, state) => const DeliveriesScreen(),
            routes: [
              GoRoute(
                path: ':id/tracking',
                name: RouteNames.deliveryTracking,
                builder: (context, state) {
                  final deliveryId = state.pathParameters['id']!;
                  return DeliveryTrackingScreen(deliveryId: deliveryId);
                },
              ),
            ],
          ),

          // Fleet Management
          GoRoute(
            path: RoutePaths.fleet,
            name: RouteNames.fleet,
            builder: (context, state) => const FleetScreen(),
          ),

          // Warehouses
          GoRoute(
            path: RoutePaths.warehouses,
            name: RouteNames.warehouses,
            builder: (context, state) => const WarehousesScreen(),
          ),

          // Settings & Profile
          GoRoute(
            path: RoutePaths.settings,
            name: RouteNames.settings,
            builder: (context, state) => const SettingsScreen(),
          ),
          GoRoute(
            path: RoutePaths.profile,
            name: RouteNames.profile,
            builder: (context, state) => const ProfileScreen(),
          ),

          // Notifications
          GoRoute(
            path: RoutePaths.notifications,
            name: RouteNames.notifications,
            builder: (context, state) => const NotificationsScreen(),
          ),
        ],
      ),
    ],
    errorBuilder: (context, state) => ErrorScreen(error: state.error),
  );
});

// ============================================================================
// HELPER WIDGETS
// ============================================================================

/// Routes to the appropriate dashboard based on participant type
class _ParticipantDashboardRouter extends ConsumerWidget {
  const _ParticipantDashboardRouter();

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final user = ref.watch(authProvider).user;

    if (user == null) {
      return const Center(child: CircularProgressIndicator());
    }

    // Route based on participant type
    switch (user.participantType) {
      case ParticipantType.bank:
        return const BankDashboardScreen();
      case ParticipantType.logistics:
        return const LogisticsDashboardScreen();
      case ParticipantType.warehouse:
        return const WarehouseDashboardScreen();
      case ParticipantType.manufacturer:
        return const ManufacturerDashboardScreen();
      case ParticipantType.retailer:
        return const RetailerDashboardScreen();
      case ParticipantType.wholesaler:
        return const WholesalerDashboardScreen();
      case ParticipantType.investor:
        return const InvestorDashboardScreen();
      case ParticipantType.entrepreneur:
        return const EntrepreneurDashboardScreen();
      case ParticipantType.ecommerce:
        return const EcommerceDashboardScreen();
      case ParticipantType.distributor:
        return const WholesalerDashboardScreen(); // Use wholesaler for distributors
      case ParticipantType.agent:
      case ParticipantType.driver:
        return const LogisticsDashboardScreen(); // Use logistics for agents/drivers
    }
  }
}

/// Refresh stream for router updates
class _RouterRefreshStream extends ChangeNotifier {
  _RouterRefreshStream(this._ref) {
    _ref.listen(authProvider, (_, __) => notifyListeners());
  }

  final Ref _ref;
}

// ============================================================================
// PLACEHOLDER SCREENS (to be implemented)
// ============================================================================

class OrdersScreen extends StatelessWidget {
  const OrdersScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class OrderDetailsScreen extends StatelessWidget {
  final String orderId;
  const OrderDetailsScreen({super.key, required this.orderId});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class InventoryScreen extends StatelessWidget {
  const InventoryScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class ProductsScreen extends StatelessWidget {
  const ProductsScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class ProductDetailsScreen extends StatelessWidget {
  final String productId;
  const ProductDetailsScreen({super.key, required this.productId});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class AnalyticsScreen extends StatelessWidget {
  const AnalyticsScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class WalletScreen extends StatelessWidget {
  const WalletScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class TransactionsScreen extends StatelessWidget {
  const TransactionsScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class DeliveriesScreen extends StatelessWidget {
  const DeliveriesScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class DeliveryTrackingScreen extends StatelessWidget {
  final String deliveryId;
  const DeliveryTrackingScreen({super.key, required this.deliveryId});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class FleetScreen extends StatelessWidget {
  const FleetScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class WarehousesScreen extends StatelessWidget {
  const WarehousesScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class NotificationsScreen extends StatelessWidget {
  const NotificationsScreen({super.key});
  @override
  Widget build(BuildContext context) => const Placeholder();
}

class ErrorScreen extends StatelessWidget {
  final Exception? error;
  const ErrorScreen({super.key, this.error});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 64, color: Colors.red),
            const SizedBox(height: 16),
            Text(
              'Page not found',
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 8),
            TextButton(
              onPressed: () => context.go(RoutePaths.dashboard),
              child: const Text('Go to Dashboard'),
            ),
          ],
        ),
      ),
    );
  }
}
