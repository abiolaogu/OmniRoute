/// OmniRoute Ecosystem - E-commerce Dashboard Screen
/// Comprehensive dashboard for e-commerce businesses and dropshippers

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class EcommerceDashboardScreen extends ConsumerStatefulWidget {
  const EcommerceDashboardScreen({super.key});

  @override
  ConsumerState<EcommerceDashboardScreen> createState() => _EcommerceDashboardScreenState();
}

class _EcommerceDashboardScreenState extends ConsumerState<EcommerceDashboardScreen> {
  int _selectedPeriod = 0; // 0: Today, 1: Week, 2: Month

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
                _buildPeriodSelector(),
                const SizedBox(height: 16),
                _buildStatsGrid(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildSalesPerformance(),
                const SizedBox(height: 20),
                _buildTopSellingProducts(),
                const SizedBox(height: 20),
                _buildRecentOrders(),
                const SizedBox(height: 20),
                _buildDropshippingOpportunities(),
                const SizedBox(height: 20),
                _buildMarketplaceIntegrations(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.ecommerceColor,
        icon: const Icon(Icons.add_business),
        label: const Text('List Product'),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 140,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.ecommerceColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'E-Commerce Hub 🛍️',
              style: AppTypography.labelSmall.copyWith(color: Colors.white70),
            ),
            Text(
              'DropShip Express',
              style: AppTypography.titleMedium.copyWith(color: Colors.white),
            ),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [
                AppColors.ecommerceColor,
                AppColors.ecommerceColor.withValues(alpha: 0.8),
              ],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Padding(
              padding: const EdgeInsets.only(right: 20),
              child: Icon(
                Icons.shopping_bag,
                size: 80,
                color: Colors.white.withValues(alpha: 0.2),
              ),
            ),
          ),
        ),
      ),
      actions: [
        Stack(
          children: [
            IconButton(
              icon: const Icon(Icons.notifications_outlined, color: Colors.white),
              onPressed: () {},
            ),
            Positioned(
              right: 8,
              top: 8,
              child: Container(
                padding: const EdgeInsets.all(4),
                decoration: const BoxDecoration(
                  color: AppColors.warning,
                  shape: BoxShape.circle,
                ),
                child: Text(
                  '12',
                  style: AppTypography.labelSmall.copyWith(
                    color: Colors.white,
                    fontSize: 10,
                  ),
                ),
              ),
            ),
          ],
        ),
        IconButton(
          icon: const Icon(Icons.store, color: Colors.white),
          onPressed: () {},
        ),
      ],
    );
  }

  Widget _buildWalletCard() {
    return WalletCard(
      balance: 1845670.00,
      pendingBalance: 325000.00,
      onTopUp: () {},
      onWithdraw: () {},
      onTransfer: () {},
    ).animate().fadeIn().slideY(begin: 0.2, end: 0);
  }

  Widget _buildPeriodSelector() {
    final periods = ['Today', 'This Week', 'This Month'];
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(
        color: AppColors.grey100,
        borderRadius: AppRadius.borderRadiusFull,
      ),
      child: Row(
        children: periods.asMap().entries.map((entry) {
          final isSelected = _selectedPeriod == entry.key;
          return Expanded(
            child: GestureDetector(
              onTap: () => setState(() => _selectedPeriod = entry.key),
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                padding: const EdgeInsets.symmetric(vertical: 12),
                decoration: BoxDecoration(
                  color: isSelected ? AppColors.white : Colors.transparent,
                  borderRadius: AppRadius.borderRadiusFull,
                  boxShadow: isSelected ? AppShadows.sm : null,
                ),
                child: Center(
                  child: Text(
                    entry.value,
                    style: AppTypography.labelMedium.copyWith(
                      color: isSelected ? AppColors.ecommerceColor : AppColors.grey600,
                      fontWeight: isSelected ? FontWeight.w600 : FontWeight.w500,
                    ),
                  ),
                ),
              ),
            ),
          );
        }).toList(),
      ),
    );
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
          title: 'Total Sales',
          value: '₦2.4M',
          icon: Icons.attach_money,
          iconColor: AppColors.success,
          growthPercentage: 23.5,
        ),
        StatCard(
          title: 'Orders',
          value: '1,248',
          icon: Icons.shopping_cart,
          iconColor: AppColors.primary,
          growthPercentage: 15.2,
        ),
        StatCard(
          title: 'Products Listed',
          value: '342',
          icon: Icons.inventory_2,
          iconColor: AppColors.ecommerceColor,
        ),
        StatCard(
          title: 'Conversion Rate',
          value: '4.8%',
          icon: Icons.trending_up,
          iconColor: AppColors.info,
          growthPercentage: 0.8,
        ),
      ].asMap().entries.map((e) {
        return e.value
            .animate(delay: Duration(milliseconds: 100 * e.key))
            .fadeIn()
            .scale(begin: const Offset(0.95, 0.95));
      }).toList(),
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
            QuickActionItem(
              icon: Icons.add_box,
              label: 'Add Product',
              color: AppColors.success,
              onTap: () {},
            ),
            QuickActionItem(
              icon: Icons.local_shipping,
              label: 'Track Shipments',
              color: AppColors.primary,
              onTap: () {},
              badge: '8',
            ),
            QuickActionItem(
              icon: Icons.people,
              label: 'Find Suppliers',
              color: AppColors.ecommerceColor,
              onTap: () {},
            ),
            QuickActionItem(
              icon: Icons.analytics,
              label: 'Analytics',
              color: AppColors.info,
              onTap: () {},
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildSalesPerformance() {
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
              Text('Sales Performance', style: AppTypography.titleMedium),
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
                    Text(
                      '+23.5%',
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.success,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            '₦2,456,780 total revenue',
            style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
          ),
          const SizedBox(height: 24),
          SizedBox(
            height: 180,
            child: LineChart(
              LineChartData(
                gridData: FlGridData(
                  show: true,
                  drawVerticalLine: false,
                  horizontalInterval: 1,
                  getDrawingHorizontalLine: (value) => FlLine(
                    color: AppColors.grey200,
                    strokeWidth: 1,
                  ),
                ),
                titlesData: FlTitlesData(
                  leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  bottomTitles: AxisTitles(
                    sideTitles: SideTitles(
                      showTitles: true,
                      getTitlesWidget: (value, meta) {
                        const labels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
                        if (value.toInt() < labels.length) {
                          return Padding(
                            padding: const EdgeInsets.only(top: 8),
                            child: Text(
                              labels[value.toInt()],
                              style: AppTypography.labelSmall,
                            ),
                          );
                        }
                        return const SizedBox.shrink();
                      },
                    ),
                  ),
                ),
                borderData: FlBorderData(show: false),
                lineBarsData: [
                  LineChartBarData(
                    spots: const [
                      FlSpot(0, 3),
                      FlSpot(1, 3.5),
                      FlSpot(2, 2.8),
                      FlSpot(3, 4.2),
                      FlSpot(4, 4.8),
                      FlSpot(5, 5.5),
                      FlSpot(6, 4.5),
                    ],
                    isCurved: true,
                    color: AppColors.ecommerceColor,
                    barWidth: 3,
                    isStrokeCapRound: true,
                    dotData: FlDotData(show: false),
                    belowBarData: BarAreaData(
                      show: true,
                      color: AppColors.ecommerceColor.withValues(alpha: 0.1),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTopSellingProducts() {
    final products = [
      {'name': 'iPhone 15 Pro Max', 'sold': 156, 'revenue': 234560000, 'trend': 12.5},
      {'name': 'Samsung Galaxy S24', 'sold': 124, 'revenue': 148800000, 'trend': 8.2},
      {'name': 'AirPods Pro 2', 'sold': 289, 'revenue': 86700000, 'trend': -2.1},
      {'name': 'MacBook Air M3', 'sold': 67, 'revenue': 100500000, 'trend': 15.8},
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Top Selling Products', actionText: 'View All', onAction: () {}),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: products.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final product = products[index];
              final trend = product['trend'] as double;
              return ListTile(
                leading: Container(
                  width: 48,
                  height: 48,
                  decoration: BoxDecoration(
                    color: AppColors.grey100,
                    borderRadius: AppRadius.borderRadiusSm,
                  ),
                  child: const Icon(Icons.shopping_bag, color: AppColors.grey500),
                ),
                title: Text(
                  product['name'] as String,
                  style: AppTypography.titleSmall,
                ),
                subtitle: Text(
                  '${product['sold']} sold • ${formatCurrency((product['revenue'] as int).toDouble())}',
                  style: AppTypography.bodySmall,
                ),
                trailing: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: trend >= 0 ? AppColors.successBg : AppColors.errorBg,
                    borderRadius: AppRadius.borderRadiusFull,
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        trend >= 0 ? Icons.trending_up : Icons.trending_down,
                        size: 14,
                        color: trend >= 0 ? AppColors.success : AppColors.error,
                      ),
                      const SizedBox(width: 4),
                      Text(
                        '${trend.abs()}%',
                        style: AppTypography.labelSmall.copyWith(
                          color: trend >= 0 ? AppColors.success : AppColors.error,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildRecentOrders() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Recent Orders', actionText: 'View All', onAction: () {}),
        const SizedBox(height: 12),
        ...List.generate(
          3,
          (i) => Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: OrderListTile(
              orderNumber: 'ORD-EC${20045 + i}',
              customerName: ['John Doe - Lagos', 'Sarah M. - Abuja', 'Mike T. - PH'][i],
              amount: [185000, 245000, 89000][i].toDouble(),
              status: ['processing', 'shipped', 'delivered'][i],
              date: DateTime.now().subtract(Duration(hours: i * 2)),
              onTap: () {},
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildDropshippingOpportunities() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(
          title: 'Dropshipping Opportunities',
          actionText: 'Browse All',
          onAction: () {},
        ),
        const SizedBox(height: 12),
        SizedBox(
          height: 200,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: 4,
            itemBuilder: (context, i) {
              final opportunities = [
                {'title': 'Electronics Bundle', 'margin': '25-40%', 'supplier': 'TechHub NG', 'minOrder': 10},
                {'title': 'Fashion Collection', 'margin': '30-50%', 'supplier': 'StyleCity', 'minOrder': 20},
                {'title': 'Home Appliances', 'margin': '20-35%', 'supplier': 'HomeMax', 'minOrder': 5},
                {'title': 'Beauty Products', 'margin': '40-60%', 'supplier': 'GlowUp Co', 'minOrder': 25},
              ];
              final item = opportunities[i];
              return Container(
                width: 180,
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
                        gradient: LinearGradient(
                          colors: [
                            AppColors.ecommerceColor.withValues(alpha: 0.8),
                            AppColors.ecommerceColor,
                          ],
                        ),
                        borderRadius: const BorderRadius.vertical(
                          top: Radius.circular(12),
                        ),
                      ),
                      child: Center(
                        child: Icon(
                          [Icons.devices, Icons.checkroom, Icons.home, Icons.spa][i],
                          size: 36,
                          color: Colors.white,
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(12),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            item['title'] as String,
                            style: AppTypography.titleSmall,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 4),
                          Text(
                            'Margin: ${item['margin']}',
                            style: AppTypography.labelSmall.copyWith(
                              color: AppColors.success,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          const SizedBox(height: 4),
                          Text(
                            'By ${item['supplier']}',
                            style: AppTypography.bodySmall.copyWith(
                              color: AppColors.grey600,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'Min: ${item['minOrder']} pcs',
                                style: AppTypography.labelSmall.copyWith(
                                  color: AppColors.grey500,
                                ),
                              ),
                              Icon(
                                Icons.arrow_forward,
                                size: 16,
                                color: AppColors.ecommerceColor,
                              ),
                            ],
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ).animate(delay: Duration(milliseconds: 100 * i)).fadeIn().slideX(begin: 0.2, end: 0);
            },
          ),
        ),
      ],
    );
  }

  Widget _buildMarketplaceIntegrations() {
    final marketplaces = [
      {'name': 'Jumia', 'status': 'Connected', 'products': 124, 'color': Colors.orange},
      {'name': 'Konga', 'status': 'Connected', 'products': 98, 'color': Colors.red},
      {'name': 'Jiji', 'status': 'Pending', 'products': 0, 'color': Colors.green},
      {'name': 'Instagram Shop', 'status': 'Connected', 'products': 67, 'color': Colors.purple},
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(
          title: 'Marketplace Integrations',
          actionText: 'Manage',
          onAction: () {},
        ),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: marketplaces.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final marketplace = marketplaces[index];
              final isConnected = marketplace['status'] == 'Connected';
              return ListTile(
                leading: Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: (marketplace['color'] as Color).withValues(alpha: 0.1),
                    borderRadius: AppRadius.borderRadiusSm,
                  ),
                  child: Icon(
                    Icons.store,
                    color: marketplace['color'] as Color,
                    size: 20,
                  ),
                ),
                title: Text(
                  marketplace['name'] as String,
                  style: AppTypography.titleSmall,
                ),
                subtitle: Text(
                  isConnected
                      ? '${marketplace['products']} products synced'
                      : 'Click to connect',
                  style: AppTypography.bodySmall,
                ),
                trailing: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: isConnected ? AppColors.successBg : AppColors.warningBg,
                    borderRadius: AppRadius.borderRadiusFull,
                  ),
                  child: Text(
                    marketplace['status'] as String,
                    style: AppTypography.labelSmall.copyWith(
                      color: isConnected ? AppColors.success : AppColors.warning,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
