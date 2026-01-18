/// OmniRoute Ecosystem - Main Dashboard Shell
/// Adaptive shell with bottom navigation that changes based on participant type

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/core/router/app_router.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class MainDashboardShell extends ConsumerStatefulWidget {
  final Widget child;

  const MainDashboardShell({super.key, required this.child});

  @override
  ConsumerState<MainDashboardShell> createState() => _MainDashboardShellState();
}

class _MainDashboardShellState extends ConsumerState<MainDashboardShell> {
  @override
  Widget build(BuildContext context) {
    final user = ref.watch(authProvider).user;
    final currentIndex = ref.watch(bottomNavIndexProvider);

    if (user == null) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }

    final navItems = _getNavItemsForParticipant(user.participantType);

    return Scaffold(
      body: widget.child,
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          color: AppColors.white,
          boxShadow: AppShadows.bottomNav,
        ),
        child: SafeArea(
          top: false,
          child: SizedBox(
            height: 64,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: navItems.asMap().entries.map((entry) {
                final index = entry.key;
                final item = entry.value;
                final isSelected = index == currentIndex;

                return _NavItem(
                  icon: item.icon,
                  activeIcon: item.activeIcon,
                  label: item.label,
                  isSelected: isSelected,
                  onTap: () {
                    ref.read(bottomNavIndexProvider.notifier).state = index;
                    context.go(item.route);
                  },
                );
              }).toList(),
            ),
          ),
        ),
      ),
      drawer: _buildDrawer(context, user.participantType),
    );
  }

  List<_NavItemData> _getNavItemsForParticipant(ParticipantType type) {
    switch (type) {
      case ParticipantType.bank:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.account_balance_outlined, Icons.account_balance, 'Loans', '/loans'),
          _NavItemData(Icons.receipt_long_outlined, Icons.receipt_long, 'Settlements', '/settlements'),
          _NavItemData(Icons.analytics_outlined, Icons.analytics, 'Analytics', RoutePaths.analytics),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.logistics:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.local_shipping_outlined, Icons.local_shipping, 'Deliveries', RoutePaths.deliveries),
          _NavItemData(Icons.directions_car_outlined, Icons.directions_car, 'Fleet', RoutePaths.fleet),
          _NavItemData(Icons.route_outlined, Icons.route, 'Routes', '/routes'),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.warehouse:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.inventory_2_outlined, Icons.inventory_2, 'Inventory', RoutePaths.inventory),
          _NavItemData(Icons.move_to_inbox_outlined, Icons.move_to_inbox, 'Inbound', '/inbound'),
          _NavItemData(Icons.outbox_outlined, Icons.outbox, 'Outbound', '/outbound'),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.manufacturer:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.category_outlined, Icons.category, 'Products', RoutePaths.products),
          _NavItemData(Icons.receipt_long_outlined, Icons.receipt_long, 'Orders', RoutePaths.orders),
          _NavItemData(Icons.analytics_outlined, Icons.analytics, 'Analytics', RoutePaths.analytics),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.retailer:
      case ParticipantType.wholesaler:
      case ParticipantType.distributor:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.receipt_long_outlined, Icons.receipt_long, 'Orders', RoutePaths.orders),
          _NavItemData(Icons.inventory_2_outlined, Icons.inventory_2, 'Inventory', RoutePaths.inventory),
          _NavItemData(Icons.account_balance_wallet_outlined, Icons.account_balance_wallet, 'Wallet', RoutePaths.wallet),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.ecommerce:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.shopping_bag_outlined, Icons.shopping_bag, 'Products', RoutePaths.products),
          _NavItemData(Icons.receipt_long_outlined, Icons.receipt_long, 'Orders', RoutePaths.orders),
          _NavItemData(Icons.local_shipping_outlined, Icons.local_shipping, 'Shipping', RoutePaths.deliveries),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.investor:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.pie_chart_outlined, Icons.pie_chart, 'Portfolio', '/portfolio'),
          _NavItemData(Icons.trending_up_outlined, Icons.trending_up, 'Opportunities', '/opportunities'),
          _NavItemData(Icons.analytics_outlined, Icons.analytics, 'Reports', RoutePaths.analytics),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.entrepreneur:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.lightbulb_outlined, Icons.lightbulb, 'Ideas', '/ideas'),
          _NavItemData(Icons.school_outlined, Icons.school, 'Learn', '/learn'),
          _NavItemData(Icons.handshake_outlined, Icons.handshake, 'Network', '/network'),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];

      case ParticipantType.agent:
      case ParticipantType.driver:
        return [
          _NavItemData(Icons.dashboard_outlined, Icons.dashboard, 'Home', RoutePaths.dashboard),
          _NavItemData(Icons.assignment_outlined, Icons.assignment, 'Tasks', '/tasks'),
          _NavItemData(Icons.account_balance_wallet_outlined, Icons.account_balance_wallet, 'Earnings', '/earnings'),
          _NavItemData(Icons.leaderboard_outlined, Icons.leaderboard, 'Performance', '/performance'),
          _NavItemData(Icons.person_outline, Icons.person, 'Profile', RoutePaths.profile),
        ];
    }
  }

  Widget _buildDrawer(BuildContext context, ParticipantType type) {
    return Drawer(
      child: ListView(
        padding: EdgeInsets.zero,
        children: [
          DrawerHeader(
            decoration: const BoxDecoration(
              gradient: AppColors.primaryGradient,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                const CircleAvatar(
                  radius: 32,
                  backgroundColor: Colors.white,
                  child: Icon(Icons.person, size: 32, color: AppColors.primary),
                ),
                const SizedBox(height: 12),
                Text(ref.read(authProvider).user?.fullName ?? '',
                    style: AppTypography.titleMedium.copyWith(color: Colors.white)),
                Text(type.displayName,
                    style: AppTypography.bodySmall.copyWith(color: Colors.white70)),
              ],
            ),
          ),
          _DrawerItem(icon: Icons.settings, label: 'Settings', onTap: () => context.push(RoutePaths.settings)),
          _DrawerItem(icon: Icons.notifications, label: 'Notifications', onTap: () => context.push(RoutePaths.notifications)),
          _DrawerItem(icon: Icons.help_outline, label: 'Help & Support', onTap: () {}),
          _DrawerItem(icon: Icons.info_outline, label: 'About', onTap: () {}),
          const Divider(),
          _DrawerItem(
            icon: Icons.logout,
            label: 'Logout',
            isDestructive: true,
            onTap: () {
              ref.read(authProvider.notifier).logout();
              context.go(RoutePaths.welcome);
            },
          ),
          const SizedBox(height: 16),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: Text('OmniRoute v1.0.0', style: AppTypography.bodySmall.copyWith(color: AppColors.grey500)),
          ),
        ],
      ),
    );
  }
}

class _NavItemData {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final String route;

  const _NavItemData(this.icon, this.activeIcon, this.label, this.route);
}

class _NavItem extends StatelessWidget {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final bool isSelected;
  final VoidCallback onTap;

  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.isSelected,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      borderRadius: AppRadius.borderRadiusMd,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 200),
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: isSelected ? AppColors.primary.withValues(alpha: 0.1) : Colors.transparent,
          borderRadius: AppRadius.borderRadiusMd,
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(isSelected ? activeIcon : icon, color: isSelected ? AppColors.primary : AppColors.grey500, size: 24),
            const SizedBox(height: 4),
            Text(label,
                style: AppTypography.labelSmall.copyWith(
                    color: isSelected ? AppColors.primary : AppColors.grey500,
                    fontWeight: isSelected ? FontWeight.w600 : FontWeight.w500)),
          ],
        ),
      ),
    );
  }
}

class _DrawerItem extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;
  final bool isDestructive;

  const _DrawerItem({required this.icon, required this.label, required this.onTap, this.isDestructive = false});

  @override
  Widget build(BuildContext context) {
    final color = isDestructive ? AppColors.error : AppColors.grey700;
    return ListTile(
      leading: Icon(icon, color: color),
      title: Text(label, style: AppTypography.bodyMedium.copyWith(color: color, fontWeight: FontWeight.w500)),
      onTap: onTap,
    );
  }
}
