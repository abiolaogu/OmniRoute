/// OmniRoute Ecosystem - Retailer Dashboard Screen
/// Comprehensive dashboard for retailers with orders, inventory, and sales insights

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class RetailerDashboardScreen extends ConsumerStatefulWidget {
  const RetailerDashboardScreen({super.key});

  @override
  ConsumerState<RetailerDashboardScreen> createState() => _RetailerDashboardScreenState();
}

class _RetailerDashboardScreenState extends ConsumerState<RetailerDashboardScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      body: CustomScrollView(
        slivers: [
          _buildAppBar(),
          SliverPadding(
            padding: const EdgeInsets.all(16),
            sliver: SliverList(
              delegate: SliverChildListDelegate([
                _buildWalletCard(),
                const SizedBox(height: 20),
                _buildStatsGrid(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildSalesChart(),
                const SizedBox(height: 20),
                _buildRecentOrders(),
                const SizedBox(height: 20),
                _buildLowStockAlerts(),
                const SizedBox(height: 20),
                _buildRecommendedProducts(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.retailerColor,
        icon: const Icon(Icons.add_shopping_cart),
        label: const Text('New Order'),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 140,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.retailerColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Good Morning! 👋', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            Text('Mama Ngozi Stores', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.retailerColor, AppColors.retailerColor.withValues(alpha: 0.8)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Padding(
              padding: const EdgeInsets.only(right: 20),
              child: Icon(Icons.store, size: 80, color: Colors.white.withValues(alpha: 0.2)),
            ),
          ),
        ),
      ),
      actions: [
        Stack(
          children: [
            IconButton(icon: const Icon(Icons.notifications_outlined, color: Colors.white), onPressed: () {}),
            Positioned(
              right: 8, top: 8,
              child: Container(
                width: 8, height: 8,
                decoration: const BoxDecoration(color: AppColors.warning, shape: BoxShape.circle),
              ),
            ),
          ],
        ),
        IconButton(icon: const Icon(Icons.qr_code_scanner, color: Colors.white), onPressed: () {}),
      ],
    );
  }

  Widget _buildWalletCard() {
    return WalletCard(
      balance: 245890.50,
      pendingBalance: 45000.00,
      onTopUp: () {},
      onWithdraw: () {},
      onTransfer: () {},
    ).animate().fadeIn().slideY(begin: 0.2, end: 0);
  }

  Widget _buildStatsGrid() {
    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2,
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.5,
      children: [
        StatCard(
          title: "Today's Sales",
          value: '₦89,450',
          icon: Icons.point_of_sale,
          iconColor: AppColors.success,
          growthPercentage: 12.5,
        ),
        StatCard(
          title: 'Active Orders',
          value: '23',
          icon: Icons.shopping_bag,
          iconColor: AppColors.primary,
        ),
        StatCard(
          title: 'Products',
          value: '156',
          icon: Icons.inventory_2,
          iconColor: AppColors.info,
        ),
        StatCard(
          title: 'Low Stock Items',
          value: '8',
          icon: Icons.warning_amber,
          iconColor: AppColors.warning,
        ),
      ].asMap().entries.map((e) => 
        e.value.animate(delay: Duration(milliseconds: 100 * e.key)).fadeIn().scale(begin: const Offset(0.95, 0.95))
      ).toList(),
    );
  }

  Widget _buildQuickActions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Quick Actions'),
        const SizedBox(height: 12),
        QuickActionGrid(
          actions: [
            QuickActionItem(icon: Icons.add_shopping_cart, label: 'Order Stock', color: AppColors.success, onTap: () {}),
            QuickActionItem(icon: Icons.qr_code_scanner, label: 'Scan Product', color: AppColors.primary, onTap: () {}),
            QuickActionItem(icon: Icons.payment, label: 'Request Loan', color: AppColors.retailerColor, onTap: () {}),
            QuickActionItem(icon: Icons.receipt_long, label: 'View Orders', color: AppColors.info, onTap: () {}, badge: '5'),
          ],
        ),
      ],
    );
  }

  Widget _buildSalesChart() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Sales This Week', style: AppTypography.titleMedium),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: AppColors.successBg,
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Row(
                  children: [
                    const Icon(Icons.trending_up, color: AppColors.success, size: 16),
                    const SizedBox(width: 4),
                    Text('+18.5%', style: AppTypography.labelSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          SizedBox(
            height: 180,
            child: BarChart(
              BarChartData(
                barGroups: _getBarGroups(),
                gridData: FlGridData(show: false),
                borderData: FlBorderData(show: false),
                titlesData: FlTitlesData(
                  leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  bottomTitles: AxisTitles(
                    sideTitles: SideTitles(
                      showTitles: true,
                      getTitlesWidget: (value, meta) {
                        const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
                        return Text(days[value.toInt()], style: AppTypography.labelSmall);
                      },
                    ),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  List<BarChartGroupData> _getBarGroups() {
    final values = [65, 78, 45, 89, 95, 120, 85];
    return values.asMap().entries.map((e) {
      return BarChartGroupData(
        x: e.key,
        barRods: [
          BarChartRodData(
            toY: e.value.toDouble(),
            color: e.key == 5 ? AppColors.success : AppColors.retailerColor,
            width: 24,
            borderRadius: const BorderRadius.vertical(top: Radius.circular(6)),
          ),
        ],
      );
    }).toList();
  }

  Widget _buildRecentOrders() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Recent Orders', actionText: 'View All', onAction: () {}),
        const SizedBox(height: 12),
        ...List.generate(3, (i) => Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: OrderListTile(
            orderNumber: 'ORD-${10045 + i}',
            customerName: 'Supplier: Dangote Foods',
            amount: 125000 + (i * 25000),
            status: i == 0 ? 'in_transit' : i == 1 ? 'delivered' : 'pending',
            date: DateTime.now().subtract(Duration(hours: i * 3)),
            onTap: () {},
          ),
        )),
      ],
    );
  }

  Widget _buildLowStockAlerts() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Low Stock Alerts', actionText: 'Restock All', onAction: () {}),
        const SizedBox(height: 12),
        SizedBox(
          height: 140,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: 5,
            itemBuilder: (context, i) => Container(
              width: 140,
              margin: const EdgeInsets.only(right: 12),
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: AppColors.white,
                borderRadius: AppRadius.borderRadiusMd,
                border: Border.all(color: AppColors.warning.withValues(alpha: 0.3)),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    width: 40, height: 40,
                    decoration: BoxDecoration(
                      color: AppColors.warningBg,
                      borderRadius: AppRadius.borderRadiusSm,
                    ),
                    child: const Icon(Icons.inventory, color: AppColors.warning, size: 20),
                  ),
                  const SizedBox(height: 8),
                  Text(['Indomie Noodles', 'Golden Penny Rice', 'Peak Milk', 'Bournvita', 'Milo'][i],
                      style: AppTypography.labelMedium, maxLines: 2, overflow: TextOverflow.ellipsis),
                  const Spacer(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text('${[5, 3, 8, 2, 6][i]} left', style: AppTypography.labelSmall.copyWith(color: AppColors.warning)),
                      const Icon(Icons.add_circle, color: AppColors.primary, size: 20),
                    ],
                  ),
                ],
              ),
            ).animate(delay: Duration(milliseconds: 100 * i)).fadeIn().slideX(begin: 0.2, end: 0),
          ),
        ),
      ],
    );
  }

  Widget _buildRecommendedProducts() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Recommended for You', actionText: 'Browse All', onAction: () {}),
        const SizedBox(height: 12),
        SizedBox(
          height: 180,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: 4,
            itemBuilder: (context, i) => Container(
              width: 150,
              margin: const EdgeInsets.only(right: 12),
              decoration: BoxDecoration(
                color: AppColors.white,
                borderRadius: AppRadius.borderRadiusMd,
                border: Border.all(color: AppColors.cardBorder),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    height: 80,
                    decoration: BoxDecoration(
                      color: AppColors.grey100,
                      borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
                    ),
                    child: Center(child: Icon(Icons.shopping_basket, size: 40, color: AppColors.grey400)),
                  ),
                  Padding(
                    padding: const EdgeInsets.all(12),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(['Nestle Maggi', 'Honeywell Flour', 'Dano Milk', 'Kings Oil'][i],
                            style: AppTypography.labelMedium, maxLines: 1, overflow: TextOverflow.ellipsis),
                        const SizedBox(height: 4),
                        Text('₦${[2500, 15000, 3500, 8500][i]}', style: AppTypography.titleSmall.copyWith(color: AppColors.success)),
                        Text('Per carton', style: AppTypography.labelSmall.copyWith(color: AppColors.grey500)),
                      ],
                    ),
                  ),
                ],
              ),
            ).animate(delay: Duration(milliseconds: 100 * i)).fadeIn().slideX(begin: 0.2, end: 0),
          ),
        ),
      ],
    );
  }
}
