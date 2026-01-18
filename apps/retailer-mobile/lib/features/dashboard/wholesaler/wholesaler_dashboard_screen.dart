/// OmniRoute Ecosystem - Wholesaler Dashboard
/// B2B sales dashboard with bulk orders, inventory management, and customer tracking

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class WholesalerDashboardScreen extends ConsumerStatefulWidget {
  const WholesalerDashboardScreen({super.key});

  @override
  ConsumerState<WholesalerDashboardScreen> createState() =>
      _WholesalerDashboardScreenState();
}

class _WholesalerDashboardScreenState
    extends ConsumerState<WholesalerDashboardScreen> {
  @override
  Widget build(BuildContext context) {
    final user = ref.watch(authProvider).user;

    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        backgroundColor: AppColors.white,
        elevation: 0,
        leading: Builder(
          builder: (context) => IconButton(
            icon: const Icon(Icons.menu, color: AppColors.grey800),
            onPressed: () => Scaffold.of(context).openDrawer(),
          ),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Good morning,',
              style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
            ),
            Text(
              user?.businessName ?? user?.fullName ?? 'Wholesaler',
              style: AppTypography.titleMedium.copyWith(color: AppColors.grey900),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.search, color: AppColors.grey800),
            onPressed: () {},
          ),
          IconButton(
            icon: Stack(
              children: [
                const Icon(Icons.notifications_outlined, color: AppColors.grey800),
                Positioned(
                  right: 0,
                  top: 0,
                  child: Container(
                    width: 8,
                    height: 8,
                    decoration: const BoxDecoration(
                      color: AppColors.error,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
              ],
            ),
            onPressed: () {},
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {},
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Wallet Card
              WalletCard(
                balance: 4850000,
                pendingBalance: 320000,
                onTopUp: () {},
                onWithdraw: () {},
                onTransfer: () {},
              )
                  .animate()
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Key Metrics
              _buildKeyMetrics()
                  .animate(delay: const Duration(milliseconds: 100))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Quick Actions
              SectionHeader(title: 'Quick Actions'),
              const SizedBox(height: 12),
              QuickActionGrid(
                actions: [
                  QuickActionItem(
                    icon: Icons.add_shopping_cart,
                    label: 'New Order',
                    color: AppColors.wholesalerColor,
                    onTap: () {},
                  ),
                  QuickActionItem(
                    icon: Icons.inventory_2,
                    label: 'Stock',
                    color: AppColors.success,
                    onTap: () {},
                    badge: '5',
                  ),
                  QuickActionItem(
                    icon: Icons.people,
                    label: 'Customers',
                    color: AppColors.info,
                    onTap: () {},
                  ),
                  QuickActionItem(
                    icon: Icons.local_shipping,
                    label: 'Deliveries',
                    color: AppColors.warning,
                    onTap: () {},
                  ),
                ],
              )
                  .animate(delay: const Duration(milliseconds: 200))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Sales Chart
              SectionHeader(
                title: 'Sales Overview',
                actionText: 'View Report',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildSalesChart()
                  .animate(delay: const Duration(milliseconds: 300))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Low Stock Alerts
              SectionHeader(
                title: 'Low Stock Alerts',
                actionText: 'Reorder',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildLowStockAlerts()
                  .animate(delay: const Duration(milliseconds: 400))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Recent Orders
              SectionHeader(
                title: 'Recent Orders',
                actionText: 'View All',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildRecentOrders()
                  .animate(delay: const Duration(milliseconds: 500))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Top Customers
              SectionHeader(
                title: 'Top Customers',
                actionText: 'View All',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildTopCustomers()
                  .animate(delay: const Duration(milliseconds: 600))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.wholesalerColor,
        icon: const Icon(Icons.add),
        label: const Text('New Order'),
      ),
    );
  }

  Widget _buildKeyMetrics() {
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.3,
      children: [
        StatCard(
          title: 'Today\'s Sales',
          value: '₦1.2M',
          icon: Icons.trending_up,
          iconColor: AppColors.wholesalerColor,
          growthPercentage: 15.3,
        ),
        StatCard(
          title: 'Orders',
          value: '47',
          icon: Icons.receipt_long,
          iconColor: AppColors.success,
          growthPercentage: 8.7,
        ),
        StatCard(
          title: 'Customers',
          value: '156',
          icon: Icons.people,
          iconColor: AppColors.info,
          growthPercentage: 5.2,
        ),
        StatCard(
          title: 'Low Stock Items',
          value: '12',
          icon: Icons.warning_amber,
          iconColor: AppColors.warning,
          growthPercentage: -3.1,
        ),
      ],
    );
  }

  Widget _buildSalesChart() {
    return Container(
      height: 200,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: BarChart(
        BarChartData(
          alignment: BarChartAlignment.spaceAround,
          maxY: 20,
          barTouchData: BarTouchData(enabled: false),
          titlesData: FlTitlesData(
            show: true,
            bottomTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                getTitlesWidget: (value, meta) {
                  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
                  return Text(
                    days[value.toInt()],
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.grey500,
                    ),
                  );
                },
                reservedSize: 30,
              ),
            ),
            leftTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
            topTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
            rightTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
          ),
          gridData: const FlGridData(show: false),
          borderData: FlBorderData(show: false),
          barGroups: [
            _makeBarGroup(0, 12),
            _makeBarGroup(1, 15),
            _makeBarGroup(2, 10),
            _makeBarGroup(3, 18),
            _makeBarGroup(4, 14),
            _makeBarGroup(5, 8),
            _makeBarGroup(6, 16),
          ],
        ),
      ),
    );
  }

  BarChartGroupData _makeBarGroup(int x, double y) {
    return BarChartGroupData(
      x: x,
      barRods: [
        BarChartRodData(
          toY: y,
          color: AppColors.wholesalerColor,
          width: 20,
          borderRadius: const BorderRadius.vertical(top: Radius.circular(4)),
        ),
      ],
    );
  }

  Widget _buildLowStockAlerts() {
    final alerts = [
      {'name': 'Golden Penny Semovita 10kg', 'stock': 15, 'reorder': 50},
      {'name': 'Peak Milk Tin 400g', 'stock': 23, 'reorder': 100},
      {'name': 'Indomie Noodles Carton', 'stock': 8, 'reorder': 30},
    ];

    return Column(
      children: alerts.map((item) {
        final stockPercent = (item['stock'] as int) / (item['reorder'] as int);
        return Container(
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusSm,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: AppColors.errorBg,
                  borderRadius: AppRadius.borderRadiusXs,
                ),
                child: const Icon(
                  Icons.inventory_2_outlined,
                  color: AppColors.error,
                  size: 20,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      item['name'] as String,
                      style: AppTypography.labelMedium,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    LinearProgressIndicator(
                      value: stockPercent,
                      backgroundColor: AppColors.grey200,
                      valueColor: AlwaysStoppedAnimation<Color>(
                        stockPercent < 0.3 ? AppColors.error : AppColors.warning,
                      ),
                      minHeight: 4,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 12),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${item['stock']} left',
                    style: AppTypography.labelMedium.copyWith(
                      color: AppColors.error,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    'Min: ${item['reorder']}',
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.grey500,
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Widget _buildRecentOrders() {
    return Column(
      children: [
        OrderListTile(
          orderNumber: 'ORD-2024-0847',
          customerName: 'Bello Supermarket',
          amount: 185000,
          status: 'processing',
          date: DateTime.now().subtract(const Duration(hours: 2)),
          onTap: () {},
        ),
        const SizedBox(height: 8),
        OrderListTile(
          orderNumber: 'ORD-2024-0846',
          customerName: 'Chika Retail Store',
          amount: 92500,
          status: 'shipped',
          date: DateTime.now().subtract(const Duration(hours: 5)),
          onTap: () {},
        ),
        const SizedBox(height: 8),
        OrderListTile(
          orderNumber: 'ORD-2024-0845',
          customerName: 'Mama Ngozi Shop',
          amount: 45000,
          status: 'delivered',
          date: DateTime.now().subtract(const Duration(hours: 8)),
          onTap: () {},
        ),
      ],
    );
  }

  Widget _buildTopCustomers() {
    final customers = [
      {'name': 'Bello Supermarket', 'orders': 45, 'total': 2850000},
      {'name': 'Chika Retail Store', 'orders': 38, 'total': 2100000},
      {'name': 'Mama Ngozi Shop', 'orders': 32, 'total': 1650000},
    ];

    return Column(
      children: customers.asMap().entries.map((entry) {
        final index = entry.key;
        final customer = entry.value;
        return Container(
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusSm,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: Row(
            children: [
              Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: _getRankColor(index),
                  shape: BoxShape.circle,
                ),
                child: Center(
                  child: Text(
                    '${index + 1}',
                    style: AppTypography.labelMedium.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      customer['name'] as String,
                      style: AppTypography.labelMedium,
                    ),
                    Text(
                      '${customer['orders']} orders',
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.grey500,
                      ),
                    ),
                  ],
                ),
              ),
              Text(
                formatCurrency((customer['total'] as int).toDouble()),
                style: AppTypography.titleSmall.copyWith(
                  color: AppColors.success,
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Color _getRankColor(int rank) {
    switch (rank) {
      case 0:
        return const Color(0xFFFFD700); // Gold
      case 1:
        return const Color(0xFFC0C0C0); // Silver
      case 2:
        return const Color(0xFFCD7F32); // Bronze
      default:
        return AppColors.grey500;
    }
  }
}
