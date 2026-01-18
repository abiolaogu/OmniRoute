/// OmniRoute Ecosystem - Distributor Dashboard
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class DistributorDashboardScreen extends ConsumerWidget {
  const DistributorDashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(title: const Text('Distribution Hub'), backgroundColor: AppColors.distributorColor, foregroundColor: Colors.white,
        actions: [IconButton(icon: const Icon(Icons.notifications_outlined), onPressed: () {}), IconButton(icon: const Icon(Icons.qr_code_scanner), onPressed: () {})],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildWelcomeCard(),
            const SizedBox(height: 20),
            _buildStatsGrid(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Order Pipeline'),
            const SizedBox(height: 12),
            _buildOrderPipeline(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Quick Actions'),
            const SizedBox(height: 12),
            _buildQuickActions(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Top Retailers'),
            const SizedBox(height: 12),
            _buildTopRetailers(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Recent Orders'),
            const SizedBox(height: 12),
            _buildRecentOrdersList(),
          ],
        ),
      ),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildWelcomeCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.distributorColor, AppColors.distributorColor.withValues(alpha: 0.8)], begin: Alignment.topLeft, end: Alignment.bottomRight),
        borderRadius: AppRadius.borderRadiusLg,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('Good Morning!', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              const SizedBox(height: 4),
              Text('Dangote Distributors', style: AppTypography.headlineMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
            ]),
            Container(padding: const EdgeInsets.all(12), decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusMd),
              child: const Icon(Icons.storefront, color: Colors.white, size: 32)),
          ]),
          const SizedBox(height: 20),
          Row(children: [
            Expanded(child: _buildMetricItem('Today\'s Sales', '₦8.2M', Icons.trending_up)),
            Container(width: 1, height: 40, color: Colors.white24),
            Expanded(child: _buildMetricItem('Orders', '47', Icons.receipt)),
            Container(width: 1, height: 40, color: Colors.white24),
            Expanded(child: _buildMetricItem('Deliveries', '32', Icons.local_shipping)),
          ]),
        ],
      ),
    );
  }

  Widget _buildMetricItem(String label, String value, IconData icon) {
    return Column(children: [
      Icon(icon, color: Colors.white70, size: 20),
      const SizedBox(height: 4),
      Text(value, style: AppTypography.titleLarge.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
      Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
    ]);
  }

  Widget _buildStatsGrid() {
    return GridView.count(
      shrinkWrap: true, physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2, crossAxisSpacing: 12, mainAxisSpacing: 12, childAspectRatio: 1.6,
      children: [
        StatCard(title: 'Active Retailers', value: '156', icon: Icons.store, iconColor: AppColors.primary, trend: '+12', trendUp: true),
        StatCard(title: 'Pending Orders', value: '23', icon: Icons.hourglass_empty, iconColor: AppColors.warning),
        StatCard(title: 'Inventory Items', value: '847', icon: Icons.inventory_2, iconColor: AppColors.info),
        StatCard(title: 'Credit Outstanding', value: '₦12.4M', icon: Icons.account_balance_wallet, iconColor: AppColors.error),
      ],
    );
  }

  Widget _buildOrderPipeline() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: [
        _buildPipelineStep('New', '12', AppColors.info, true),
        _buildPipelineArrow(),
        _buildPipelineStep('Processing', '8', AppColors.warning, false),
        _buildPipelineArrow(),
        _buildPipelineStep('Shipped', '15', AppColors.primary, false),
        _buildPipelineArrow(),
        _buildPipelineStep('Delivered', '32', AppColors.success, false),
      ]),
    );
  }

  Widget _buildPipelineStep(String label, String count, Color color, bool isActive) {
    return Column(children: [
      Container(width: 48, height: 48, decoration: BoxDecoration(color: isActive ? color : color.withValues(alpha: 0.1), shape: BoxShape.circle, border: Border.all(color: color, width: 2)),
        child: Center(child: Text(count, style: TextStyle(color: isActive ? Colors.white : color, fontWeight: FontWeight.bold)))),
      const SizedBox(height: 8),
      Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
    ]);
  }

  Widget _buildPipelineArrow() => Icon(Icons.arrow_forward, color: AppColors.textSecondary.withValues(alpha: 0.5), size: 20);

  Widget _buildQuickActions() {
    return QuickActionGrid(actions: [
      QuickActionItem(icon: Icons.add_shopping_cart, label: 'New Order', color: AppColors.success, onTap: () {}),
      QuickActionItem(icon: Icons.qr_code, label: 'Scan Product', color: AppColors.distributorColor, onTap: () {}),
      QuickActionItem(icon: Icons.route, label: 'Track Delivery', color: AppColors.primary, onTap: () {}),
      QuickActionItem(icon: Icons.inventory, label: 'Stock Check', color: AppColors.info, onTap: () {}),
      QuickActionItem(icon: Icons.receipt_long, label: 'Invoices', color: AppColors.warning, onTap: () {}),
      QuickActionItem(icon: Icons.people, label: 'Retailers', color: AppColors.secondary, onTap: () {}),
    ]);
  }

  Widget _buildTopRetailers() {
    final retailers = [
      {'name': 'Shoprite Mall', 'orders': '124', 'value': '₦4.2M', 'growth': '+15%'},
      {'name': 'SPAR Nigeria', 'orders': '98', 'value': '₦3.1M', 'growth': '+8%'},
      {'name': 'Justrite Stores', 'orders': '87', 'value': '₦2.8M', 'growth': '+22%'},
    ];
    return Column(children: retailers.map((r) => Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Row(children: [
        CircleAvatar(backgroundColor: AppColors.distributorColor.withValues(alpha: 0.1), child: Text(r['name']![0], style: TextStyle(color: AppColors.distributorColor, fontWeight: FontWeight.bold))),
        const SizedBox(width: 12),
        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(r['name']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
          Text('${r['orders']} orders this month', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ])),
        Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
          Text(r['value']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
          Text(r['growth']!, style: AppTypography.labelSmall.copyWith(color: AppColors.success)),
        ]),
      ]),
    )).toList());
  }

  Widget _buildRecentOrdersList() {
    final orders = [
      {'id': 'ORD-7842', 'retailer': 'Mama Put Restaurant', 'items': '12 items', 'total': '₦145,000', 'status': 'Processing'},
      {'id': 'ORD-7841', 'retailer': 'Lagos Supermarket', 'items': '8 items', 'total': '₦89,500', 'status': 'Shipped'},
      {'id': 'ORD-7840', 'retailer': 'Ibadan Wholesale', 'items': '25 items', 'total': '₦312,000', 'status': 'Delivered'},
    ];
    return Column(children: orders.map((o) => OrderListItem(orderId: o['id']!, customer: o['retailer']!, items: o['items']!, total: o['total']!, status: o['status']!, onTap: () {})).toList());
  }

  Widget _buildBottomNav() {
    return BottomNavigationBar(
      currentIndex: 0, type: BottomNavigationBarType.fixed, selectedItemColor: AppColors.distributorColor,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.dashboard), label: 'Dashboard'),
        BottomNavigationBarItem(icon: Icon(Icons.receipt_long), label: 'Orders'),
        BottomNavigationBarItem(icon: Icon(Icons.inventory_2), label: 'Inventory'),
        BottomNavigationBarItem(icon: Icon(Icons.store), label: 'Retailers'),
        BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
      ],
    );
  }
}
