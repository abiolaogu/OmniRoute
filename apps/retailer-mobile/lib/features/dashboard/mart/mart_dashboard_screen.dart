/// OmniRoute Ecosystem - Mart/Mini-Mart Dashboard Screen
/// Streamlined dashboard for small convenience stores and mini-marts
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class MartDashboardScreen extends ConsumerStatefulWidget {
  const MartDashboardScreen({super.key});
  @override
  ConsumerState<MartDashboardScreen> createState() => _MartDashboardScreenState();
}

class _MartDashboardScreenState extends ConsumerState<MartDashboardScreen> {
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
                _buildTodaySalesCard(),
                const SizedBox(height: 16),
                _buildQuickStats(),
                const SizedBox(height: 16),
                _buildQuickRestock(),
                const SizedBox(height: 16),
                _buildTopSellingToday(),
                const SizedBox(height: 16),
                _buildMobilePOS(),
                const SizedBox(height: 16),
                _buildCreditCustomers(),
                const SizedBox(height: 16),
                _buildQuickActions(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.large(
        onPressed: () {},
        backgroundColor: AppColors.martColor,
        child: const Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.point_of_sale, size: 28),
            Text('SELL', style: TextStyle(fontSize: 10, fontWeight: FontWeight.w600)),
          ],
        ),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 120,
      floating: true,
      pinned: true,
      backgroundColor: AppColors.martColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 12),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('CONVENIENCE STORE', style: AppTypography.labelSmall.copyWith(color: Colors.white70, fontSize: 10)),
            Text('Quick Mart Express', style: AppTypography.titleSmall.copyWith(color: Colors.white)),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(colors: [AppColors.martColor, AppColors.martColor.withValues(alpha: 0.85)]),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Icon(Icons.store, size: 60, color: Colors.white.withValues(alpha: 0.2)),
          ),
        ),
      ),
      actions: [
        IconButton(icon: const Icon(Icons.qr_code_scanner, color: Colors.white), onPressed: () {}),
        IconButton(icon: const Icon(Icons.notifications_outlined, color: Colors.white), onPressed: () {}),
      ],
    );
  }

  Widget _buildTodaySalesCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.martColor, AppColors.martColor.withValues(alpha: 0.85)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Today\'s Sales', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
                const SizedBox(height: 4),
                Text('₦127,450', style: AppTypography.headlineMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
                const SizedBox(height: 8),
                Row(children: [
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                    decoration: BoxDecoration(color: AppColors.success, borderRadius: AppRadius.borderRadiusSm),
                    child: Row(children: [
                      const Icon(Icons.trending_up, color: Colors.white, size: 12),
                      const SizedBox(width: 2),
                      Text('+23%', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
                    ]),
                  ),
                  const SizedBox(width: 8),
                  Text('vs yesterday', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
                ]),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), shape: BoxShape.circle),
            child: Column(
              children: [
                Text('156', style: AppTypography.titleLarge.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
                Text('sales', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickStats() {
    return Row(
      children: [
        Expanded(child: _buildStatTile('Avg. Sale', '₦817', Icons.shopping_basket, AppColors.info)),
        const SizedBox(width: 12),
        Expanded(child: _buildStatTile('Profit', '₦28,450', Icons.trending_up, AppColors.success)),
        const SizedBox(width: 12),
        Expanded(child: _buildStatTile('Credit', '₦15,200', Icons.credit_card, AppColors.warning)),
      ],
    );
  }

  Widget _buildStatTile(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
      child: Column(
        children: [
          Icon(icon, color: color, size: 20),
          const SizedBox(height: 4),
          Text(value, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
          Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ],
      ),
    );
  }

  Widget _buildQuickRestock() {
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
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(children: [
                const Icon(Icons.warning, color: AppColors.error, size: 20),
                const SizedBox(width: 8),
                Text('Quick Restock Needed', style: AppTypography.titleSmall.copyWith(color: AppColors.error)),
              ]),
              TextButton(onPressed: () {}, child: const Text('Order All')),
            ],
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _buildRestockChip('Peak Milk', 3),
              _buildRestockChip('Soft Drinks', 5),
              _buildRestockChip('Bread', 2),
              _buildRestockChip('Sugar', 4),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildRestockChip(String item, int qty) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm, border: Border.all(color: AppColors.error.withValues(alpha: 0.3))),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(item, style: AppTypography.labelMedium),
          const SizedBox(width: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
            decoration: BoxDecoration(color: AppColors.error, borderRadius: AppRadius.borderRadiusSm),
            child: Text('$qty left', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
          ),
        ],
      ),
    );
  }

  Widget _buildTopSellingToday() {
    final items = [
      {'name': 'Coca-Cola 50cl', 'qty': 45, 'revenue': '₦11,250'},
      {'name': 'MTN Airtime', 'qty': 38, 'revenue': '₦19,000'},
      {'name': 'Indomie Chicken', 'qty': 32, 'revenue': '₦6,400'},
      {'name': 'Pure Water (Bag)', 'qty': 28, 'revenue': '₦5,600'},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Top Selling Today'),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: items.length,
            separatorBuilder: (_, __) => const Divider(height: 1, indent: 56),
            itemBuilder: (context, i) {
              final item = items[i];
              return ListTile(
                dense: true,
                leading: CircleAvatar(
                  radius: 16,
                  backgroundColor: AppColors.martColor.withValues(alpha: 0.1),
                  child: Text('${i + 1}', style: TextStyle(color: AppColors.martColor, fontWeight: FontWeight.bold, fontSize: 12)),
                ),
                title: Text(item['name'] as String, style: AppTypography.titleSmall),
                subtitle: Text('${item['qty']} sold', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                trailing: Text(item['revenue'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildMobilePOS() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.success, AppColors.success.withValues(alpha: 0.85)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Row(
        children: [
          const Icon(Icons.phone_android, color: Colors.white, size: 40),
          const SizedBox(width: 16),
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('Mobile POS Active', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
              Text('Accept payments anywhere in your store', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            ],
          )),
          OutlinedButton(
            onPressed: () {},
            style: OutlinedButton.styleFrom(foregroundColor: Colors.white, side: const BorderSide(color: Colors.white)),
            child: const Text('Open POS'),
          ),
        ],
      ),
    );
  }

  Widget _buildCreditCustomers() {
    final customers = [
      {'name': 'Mama Ngozi', 'amount': '₦8,500', 'days': 5},
      {'name': 'Baba Tunde', 'amount': '₦4,200', 'days': 12},
      {'name': 'Iya Basira', 'amount': '₦2,500', 'days': 3},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Credit Customers', onViewAll: () {}),
        const SizedBox(height: 8),
        ...customers.map((c) {
          final days = c['days'] as int;
          final isOverdue = days > 7;
          return Container(
            margin: const EdgeInsets.only(bottom: 8),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: AppRadius.borderRadiusSm,
              border: isOverdue ? Border.all(color: AppColors.error.withValues(alpha: 0.3)) : null,
            ),
            child: Row(
              children: [
                CircleAvatar(radius: 18, backgroundColor: AppColors.grey200, child: const Icon(Icons.person, size: 18)),
                const SizedBox(width: 12),
                Expanded(child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(c['name'] as String, style: AppTypography.titleSmall),
                    Text('$days days ago', style: AppTypography.labelSmall.copyWith(color: isOverdue ? AppColors.error : AppColors.textSecondary)),
                  ],
                )),
                Text(c['amount'] as String, style: AppTypography.titleSmall.copyWith(color: isOverdue ? AppColors.error : AppColors.textPrimary, fontWeight: FontWeight.w600)),
                const SizedBox(width: 8),
                IconButton(icon: Icon(Icons.message, size: 18, color: AppColors.martColor), onPressed: () {}),
              ],
            ),
          );
        }),
      ],
    );
  }

  Widget _buildQuickActions() {
    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 4,
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 0.9,
      children: [
        _buildActionTile(Icons.qr_code_scanner, 'Scan', AppColors.primary),
        _buildActionTile(Icons.add_shopping_cart, 'Stock In', AppColors.success),
        _buildActionTile(Icons.calculate, 'Expenses', AppColors.warning),
        _buildActionTile(Icons.bar_chart, 'Reports', AppColors.info),
      ],
    );
  }

  Widget _buildActionTile(IconData icon, String label, Color color) {
    return GestureDetector(
      onTap: () {},
      child: Container(
        decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(color: color.withValues(alpha: 0.1), shape: BoxShape.circle),
              child: Icon(icon, color: color, size: 22),
            ),
            const SizedBox(height: 6),
            Text(label, style: AppTypography.labelSmall),
          ],
        ),
      ),
    );
  }
}
