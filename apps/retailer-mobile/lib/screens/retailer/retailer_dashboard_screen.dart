/// OmniRoute Ecosystem - Retailer Dashboard
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class RetailerDashboardScreen extends ConsumerStatefulWidget {
  const RetailerDashboardScreen({super.key});
  @override ConsumerState<RetailerDashboardScreen> createState() => _RetailerDashboardScreenState();
}

class _RetailerDashboardScreenState extends ConsumerState<RetailerDashboardScreen> {
  int _selectedIndex = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('My Store'), backgroundColor: AppColors.retailerColor, foregroundColor: Colors.white,
        actions: [
          IconButton(icon: const Icon(Icons.shopping_cart_outlined), onPressed: () {}),
          IconButton(icon: const Icon(Icons.notifications_outlined), onPressed: () {}),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildCreditCard(),
            const SizedBox(height: 20),
            _buildStatsRow(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Order From Distributors'),
            const SizedBox(height: 12),
            _buildDistributorsList(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Low Stock Alert'),
            const SizedBox(height: 12),
            _buildLowStockItems(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Today\'s Sales'),
            const SizedBox(height: 12),
            _buildTodaySales(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Recent Transactions'),
            const SizedBox(height: 12),
            _buildRecentTransactions(),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {}, backgroundColor: AppColors.retailerColor,
        icon: const Icon(Icons.add_shopping_cart, color: Colors.white),
        label: const Text('New Order', style: TextStyle(color: Colors.white)),
      ),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildCreditCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.retailerColor, AppColors.retailerColor.withValues(alpha: 0.85)], begin: Alignment.topLeft, end: Alignment.bottomRight),
        borderRadius: AppRadius.borderRadiusLg,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Text('Trade Credit', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
            Container(padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4), decoration: BoxDecoration(color: Colors.green, borderRadius: AppRadius.borderRadiusSm),
              child: const Text('Active', style: TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w600))),
          ]),
          const SizedBox(height: 12),
          Text('₦2,500,000', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          Text('Available Credit', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
          const SizedBox(height: 20),
          LinearProgressIndicator(value: 0.35, backgroundColor: Colors.white24, valueColor: const AlwaysStoppedAnimation<Color>(Colors.white)),
          const SizedBox(height: 8),
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Text('Used: ₦1,350,000', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            Text('Limit: ₦3,850,000', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
          ]),
        ],
      ),
    );
  }

  Widget _buildStatsRow() {
    return Row(children: [
      Expanded(child: _buildStatItem('Today\'s Sales', '₦125,400', Icons.point_of_sale, AppColors.success)),
      const SizedBox(width: 12),
      Expanded(child: _buildStatItem('Pending Orders', '3', Icons.hourglass_empty, AppColors.warning)),
      const SizedBox(width: 12),
      Expanded(child: _buildStatItem('Low Stock', '8', Icons.warning_amber, AppColors.error)),
    ]);
  }

  Widget _buildStatItem(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Icon(icon, color: color, size: 24),
        const SizedBox(height: 8),
        Text(value, style: AppTypography.titleLarge.copyWith(fontWeight: FontWeight.w700)),
        Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      ]),
    );
  }

  Widget _buildDistributorsList() {
    final distributors = [
      {'name': 'Dangote Distributors', 'products': '124 products', 'delivery': '24hr delivery', 'logo': 'D'},
      {'name': 'SPAR Wholesale', 'products': '256 products', 'delivery': 'Same day', 'logo': 'S'},
      {'name': 'Flour Mills Direct', 'products': '89 products', 'delivery': '48hr delivery', 'logo': 'F'},
    ];
    return SizedBox(
      height: 120,
      child: ListView.separated(
        scrollDirection: Axis.horizontal, itemCount: distributors.length, separatorBuilder: (_, __) => const SizedBox(width: 12),
        itemBuilder: (context, index) {
          final d = distributors[index];
          return Container(
            width: 160, padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, border: Border.all(color: AppColors.borderColor)),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                CircleAvatar(backgroundColor: AppColors.retailerColor.withValues(alpha: 0.1), radius: 16, child: Text(d['logo']!, style: TextStyle(color: AppColors.retailerColor, fontWeight: FontWeight.bold, fontSize: 12))),
                const SizedBox(width: 8),
                Expanded(child: Text(d['name']!, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600), maxLines: 1, overflow: TextOverflow.ellipsis)),
              ]),
              const Spacer(),
              Text(d['products']!, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
              Text(d['delivery']!, style: AppTypography.labelSmall.copyWith(color: AppColors.success)),
            ]),
          );
        },
      ),
    );
  }

  Widget _buildLowStockItems() {
    final items = [
      {'name': 'Golden Penny Semovita 5kg', 'stock': '5', 'reorder': '20'},
      {'name': 'Peak Milk 400g', 'stock': '12', 'reorder': '50'},
      {'name': 'Indomie Chicken 70g (carton)', 'stock': '3', 'reorder': '15'},
    ];
    return Container(
      decoration: BoxDecoration(color: AppColors.error.withValues(alpha: 0.05), borderRadius: AppRadius.borderRadiusMd, border: Border.all(color: AppColors.error.withValues(alpha: 0.2))),
      child: Column(children: items.asMap().entries.map((e) => Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(border: e.key < items.length - 1 ? Border(bottom: BorderSide(color: AppColors.error.withValues(alpha: 0.1))) : null),
        child: Row(children: [
          Icon(Icons.warning_amber, color: AppColors.error, size: 20),
          const SizedBox(width: 12),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(e.value['name']!, style: AppTypography.labelMedium),
            Text('${e.value['stock']} left (reorder at ${e.value['reorder']})', style: AppTypography.labelSmall.copyWith(color: AppColors.error)),
          ])),
          TextButton(onPressed: () {}, child: const Text('Reorder')),
        ]),
      )).toList()),
    );
  }

  Widget _buildTodaySales() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(children: [
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text('Total Sales', style: AppTypography.labelMedium.copyWith(color: AppColors.textSecondary)),
            Text('₦125,400', style: AppTypography.headlineMedium.copyWith(fontWeight: FontWeight.w700)),
          ]),
          Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
            Text('Transactions', style: AppTypography.labelMedium.copyWith(color: AppColors.textSecondary)),
            Text('24', style: AppTypography.headlineMedium.copyWith(fontWeight: FontWeight.w700)),
          ]),
        ]),
        const SizedBox(height: 16),
        Row(children: [
          Expanded(child: _buildSalesChannel('Cash', '₦78,200', 0.62)),
          const SizedBox(width: 16),
          Expanded(child: _buildSalesChannel('Transfer', '₦35,700', 0.28)),
          const SizedBox(width: 16),
          Expanded(child: _buildSalesChannel('Credit', '₦11,500', 0.10)),
        ]),
      ]),
    );
  }

  Widget _buildSalesChannel(String label, String amount, double percentage) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      const SizedBox(height: 4),
      Text(amount, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
      const SizedBox(height: 4),
      LinearProgressIndicator(value: percentage, backgroundColor: AppColors.borderColor, valueColor: AlwaysStoppedAnimation<Color>(AppColors.retailerColor)),
    ]);
  }

  Widget _buildRecentTransactions() {
    final transactions = [
      {'customer': 'Walk-in Customer', 'items': '5 items', 'amount': '₦4,500', 'time': '10:24 AM', 'type': 'cash'},
      {'customer': 'Mrs. Adebayo', 'items': '3 items', 'amount': '₦12,800', 'time': '09:45 AM', 'type': 'transfer'},
      {'customer': 'Chief Okonkwo', 'items': '8 items', 'amount': '₦28,000', 'time': '09:12 AM', 'type': 'credit'},
    ];
    return Column(children: transactions.map((t) => Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
      child: Row(children: [
        CircleAvatar(backgroundColor: AppColors.retailerColor.withValues(alpha: 0.1), radius: 20, child: Icon(t['type'] == 'cash' ? Icons.payments : t['type'] == 'transfer' ? Icons.send : Icons.credit_card, color: AppColors.retailerColor, size: 20)),
        const SizedBox(width: 12),
        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(t['customer']!, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
          Text('${t['items']} • ${t['time']}', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ])),
        Text(t['amount']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600, color: AppColors.success)),
      ]),
    )).toList());
  }

  Widget _buildBottomNav() {
    return BottomNavigationBar(
      currentIndex: _selectedIndex, onTap: (i) => setState(() => _selectedIndex = i), type: BottomNavigationBarType.fixed, selectedItemColor: AppColors.retailerColor,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.dashboard), label: 'Home'),
        BottomNavigationBarItem(icon: Icon(Icons.shopping_bag), label: 'Order'),
        BottomNavigationBarItem(icon: Icon(Icons.point_of_sale), label: 'POS'),
        BottomNavigationBarItem(icon: Icon(Icons.inventory_2), label: 'Inventory'),
        BottomNavigationBarItem(icon: Icon(Icons.account_circle), label: 'Account'),
      ],
    );
  }
}
