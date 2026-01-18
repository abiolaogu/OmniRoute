/// OmniRoute Ecosystem - Supermarket Chain Dashboard Screen
/// Comprehensive dashboard for supermarket chains and retail operations
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class SupermarketDashboardScreen extends ConsumerStatefulWidget {
  const SupermarketDashboardScreen({super.key});
  @override
  ConsumerState<SupermarketDashboardScreen> createState() => _SupermarketDashboardScreenState();
}

class _SupermarketDashboardScreenState extends ConsumerState<SupermarketDashboardScreen> {
  int _selectedStore = 0;

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
                _buildStoreSelector(),
                const SizedBox(height: 20),
                _buildSalesOverview(),
                const SizedBox(height: 20),
                _buildCategoryPerformance(),
                const SizedBox(height: 20),
                _buildShrinkageAlerts(),
                const SizedBox(height: 20),
                _buildPlanogramCompliance(),
                const SizedBox(height: 20),
                _buildCheckoutStats(),
                const SizedBox(height: 20),
                _buildLowStockAlerts(),
                const SizedBox(height: 20),
                _buildPromotionPerformance(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.supermarketColor,
        icon: const Icon(Icons.qr_code_scanner),
        label: const Text('Scan Product'),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 140,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.supermarketColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('SUPERMARKET CHAIN', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            Text('FreshMart Nigeria', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.supermarketColor, AppColors.supermarketColor.withValues(alpha: 0.8)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Icon(Icons.shopping_cart, size: 80, color: Colors.white.withValues(alpha: 0.2)),
          ),
        ),
      ),
    );
  }

  Widget _buildStoreSelector() {
    final stores = ['All Stores (12)', 'Lagos Flagship', 'Abuja Central', 'Port Harcourt'];
    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(color: AppColors.grey100, borderRadius: AppRadius.borderRadiusFull),
      child: Row(
        children: stores.asMap().entries.map((e) {
          final isSelected = _selectedStore == e.key;
          return Expanded(
            child: GestureDetector(
              onTap: () => setState(() => _selectedStore = e.key),
              child: Container(
                padding: const EdgeInsets.symmetric(vertical: 10),
                decoration: BoxDecoration(
                  color: isSelected ? Colors.white : Colors.transparent,
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Center(
                  child: Text(
                    e.value,
                    style: AppTypography.labelSmall.copyWith(
                      color: isSelected ? AppColors.supermarketColor : AppColors.textSecondary,
                      fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                    ),
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildSalesOverview() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.supermarketColor, AppColors.supermarketColor.withValues(alpha: 0.85)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Today\'s Sales', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: AppColors.success, borderRadius: AppRadius.borderRadiusSm),
                child: Row(children: [
                  const Icon(Icons.trending_up, color: Colors.white, size: 14),
                  const SizedBox(width: 4),
                  Text('+18.3%', style: AppTypography.labelSmall.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
                ]),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text('₦24,567,890', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildSalesMetric('Transactions', '3,456'),
              _buildSalesMetric('Avg. Basket', '₦7,110'),
              _buildSalesMetric('Items Sold', '28,450'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSalesMetric(String label, String value) {
    return Expanded(
      child: Column(
        children: [
          Text(value, style: AppTypography.titleMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
          Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
        ],
      ),
    );
  }

  Widget _buildCategoryPerformance() {
    final categories = [
      {'name': 'Fresh Produce', 'sales': '₦5.2M', 'percentage': 0.21, 'color': AppColors.success},
      {'name': 'Beverages', 'sales': '₦4.8M', 'percentage': 0.19, 'color': AppColors.info},
      {'name': 'Dairy & Frozen', 'sales': '₦3.9M', 'percentage': 0.16, 'color': AppColors.warning},
      {'name': 'Household', 'sales': '₦3.2M', 'percentage': 0.13, 'color': AppColors.primary},
      {'name': 'Personal Care', 'sales': '₦2.8M', 'percentage': 0.11, 'color': AppColors.error},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Category Performance'),
        const SizedBox(height: 12),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
          child: Column(
            children: categories.map((c) => Padding(
              padding: const EdgeInsets.symmetric(vertical: 8),
              child: Row(
                children: [
                  Container(width: 4, height: 32, decoration: BoxDecoration(color: c['color'] as Color, borderRadius: AppRadius.borderRadiusSm)),
                  const SizedBox(width: 12),
                  Expanded(child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(c['name'] as String, style: AppTypography.titleSmall),
                      const SizedBox(height: 4),
                      LinearProgressIndicator(
                        value: c['percentage'] as double,
                        backgroundColor: (c['color'] as Color).withValues(alpha: 0.2),
                        valueColor: AlwaysStoppedAnimation(c['color'] as Color),
                      ),
                    ],
                  )),
                  const SizedBox(width: 12),
                  Text(c['sales'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
                ],
              ),
            )).toList(),
          ),
        ),
      ],
    );
  }

  Widget _buildShrinkageAlerts() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.error.withValues(alpha: 0.05),
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.error.withValues(alpha: 0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.warning_amber, color: AppColors.error),
              const SizedBox(width: 8),
              Text('Shrinkage Alerts', style: AppTypography.titleMedium.copyWith(color: AppColors.error)),
              const Spacer(),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: AppColors.error, borderRadius: AppRadius.borderRadiusSm),
                child: Text('₦1.2M this week', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
              ),
            ],
          ),
          const SizedBox(height: 12),
          _buildShrinkageItem('Fresh Meat Section', 'High variance - 8.5% loss', '₦450,000'),
          _buildShrinkageItem('Electronics Aisle 3', 'Missing inventory detected', '₦380,000'),
          _buildShrinkageItem('Dairy Cold Storage', 'Expired products - 23 items', '₦120,000'),
        ],
      ),
    );
  }

  Widget _buildShrinkageItem(String location, String issue, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Row(
        children: [
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(location, style: AppTypography.titleSmall),
              Text(issue, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            ],
          )),
          Text(value, style: AppTypography.titleSmall.copyWith(color: AppColors.error, fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }

  Widget _buildPlanogramCompliance() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Planogram Compliance', style: AppTypography.titleMedium),
              Text('87%', style: AppTypography.headlineSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w700)),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildComplianceCard('Aisle 1-5', 92, AppColors.success),
              const SizedBox(width: 8),
              _buildComplianceCard('Aisle 6-10', 78, AppColors.warning),
              const SizedBox(width: 8),
              _buildComplianceCard('Fresh Section', 95, AppColors.success),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildComplianceCard(String section, int compliance, Color color) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
        child: Column(
          children: [
            Text('$compliance%', style: AppTypography.titleMedium.copyWith(color: color, fontWeight: FontWeight.w700)),
            Text(section, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary), textAlign: TextAlign.center),
          ],
        ),
      ),
    );
  }

  Widget _buildCheckoutStats() {
    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2,
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.6,
      children: [
        StatCard(title: 'Active Lanes', value: '18/24', icon: Icons.point_of_sale, iconColor: AppColors.success),
        StatCard(title: 'Avg. Wait Time', value: '2.4 min', icon: Icons.timer, iconColor: AppColors.warning),
        StatCard(title: 'Self-Checkout', value: '34%', icon: Icons.touch_app, iconColor: AppColors.info),
        StatCard(title: 'Peak Hour', value: '12-2 PM', icon: Icons.schedule, iconColor: AppColors.primary),
      ],
    );
  }

  Widget _buildLowStockAlerts() {
    final items = [
      {'name': 'Peak Milk 400g', 'stock': 12, 'reorder': 50},
      {'name': 'Indomie Chicken 70g', 'stock': 24, 'reorder': 100},
      {'name': 'Golden Penny Semovita', 'stock': 8, 'reorder': 40},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Low Stock Alerts', onViewAll: () {}),
        const SizedBox(height: 12),
        ...items.map((item) => Container(
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(color: AppColors.warning.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
                child: const Icon(Icons.inventory_2, color: AppColors.warning, size: 20),
              ),
              const SizedBox(width: 12),
              Expanded(child: Text(item['name'] as String, style: AppTypography.titleSmall)),
              Text('${item['stock']}/${item['reorder']}', style: AppTypography.titleSmall.copyWith(color: AppColors.warning, fontWeight: FontWeight.w600)),
              const SizedBox(width: 8),
              TextButton(onPressed: () {}, child: const Text('Reorder')),
            ],
          ),
        )),
      ],
    );
  }

  Widget _buildPromotionPerformance() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Active Promotions', style: AppTypography.titleMedium),
              Text('8 campaigns', style: AppTypography.labelMedium.copyWith(color: AppColors.textSecondary)),
            ],
          ),
          const SizedBox(height: 16),
          _buildPromoItem('Buy 3 Get 1 Free - Beverages', '+45% uplift', AppColors.success),
          _buildPromoItem('Weekend Fresh Deals', '+28% uplift', AppColors.success),
          _buildPromoItem('Loyalty Double Points', '+12% transactions', AppColors.info),
        ],
      ),
    );
  }

  Widget _buildPromoItem(String name, String impact, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          const Icon(Icons.local_offer, color: AppColors.supermarketColor, size: 20),
          const SizedBox(width: 12),
          Expanded(child: Text(name, style: AppTypography.titleSmall)),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
            child: Text(impact, style: AppTypography.labelSmall.copyWith(color: color, fontWeight: FontWeight.w600)),
          ),
        ],
      ),
    );
  }
}
