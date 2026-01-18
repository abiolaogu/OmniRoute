/// OmniRoute Ecosystem - Inventory Screen
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class InventoryScreen extends ConsumerStatefulWidget {
  const InventoryScreen({super.key});
  @override ConsumerState<InventoryScreen> createState() => _InventoryScreenState();
}

class _InventoryScreenState extends ConsumerState<InventoryScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('Inventory'),
        actions: [
          IconButton(icon: const Icon(Icons.qr_code_scanner), onPressed: () {}),
          IconButton(icon: const Icon(Icons.add), onPressed: () {}),
        ],
        bottom: TabBar(
          controller: _tabController,
          isScrollable: true,
          tabs: const [Tab(text: 'Overview'), Tab(text: 'Stock Levels'), Tab(text: 'Movements'), Tab(text: 'Alerts')],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [_buildOverviewTab(), _buildStockLevelsTab(), _buildMovementsTab(), _buildAlertsTab()],
      ),
    );
  }

  Widget _buildOverviewTab() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildSummaryCards(),
          const SizedBox(height: 24),
          const SectionHeader(title: 'Stock Value by Category'),
          const SizedBox(height: 12),
          _buildCategoryBreakdown(),
          const SizedBox(height: 24),
          const SectionHeader(title: 'Warehouse Distribution'),
          const SizedBox(height: 12),
          _buildWarehouseCards(),
          const SizedBox(height: 24),
          const SectionHeader(title: 'Recent Activity'),
          const SizedBox(height: 12),
          _buildRecentActivity(),
        ],
      ),
    );
  }

  Widget _buildSummaryCards() {
    return GridView.count(
      shrinkWrap: true, physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2, crossAxisSpacing: 12, mainAxisSpacing: 12, childAspectRatio: 1.5,
      children: [
        _buildSummaryCard('Total SKUs', '847', Icons.inventory_2, AppColors.primary),
        _buildSummaryCard('Total Value', '₦125M', Icons.account_balance_wallet, AppColors.success),
        _buildSummaryCard('Low Stock', '23', Icons.warning, AppColors.warning),
        _buildSummaryCard('Out of Stock', '8', Icons.error, AppColors.error),
      ],
    );
  }

  Widget _buildSummaryCard(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
        Icon(icon, color: color, size: 28),
        Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(value, style: AppTypography.headlineSmall.copyWith(fontWeight: FontWeight.w700)),
          Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ]),
      ]),
    );
  }

  Widget _buildCategoryBreakdown() {
    final categories = [
      {'name': 'Food & Beverages', 'value': '₦45M', 'percentage': 0.36, 'color': AppColors.primary},
      {'name': 'Personal Care', 'value': '₦28M', 'percentage': 0.22, 'color': AppColors.success},
      {'name': 'Household', 'value': '₦22M', 'percentage': 0.18, 'color': AppColors.info},
      {'name': 'Electronics', 'value': '₦18M', 'percentage': 0.14, 'color': AppColors.warning},
      {'name': 'Others', 'value': '₦12M', 'percentage': 0.10, 'color': AppColors.textSecondary},
    ];
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(children: categories.map((c) => Padding(
        padding: const EdgeInsets.only(bottom: 12),
        child: Row(children: [
          Container(width: 12, height: 12, decoration: BoxDecoration(color: c['color'] as Color, shape: BoxShape.circle)),
          const SizedBox(width: 12),
          Expanded(child: Text(c['name'] as String, style: AppTypography.labelMedium)),
          Text(c['value'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
          const SizedBox(width: 12),
          SizedBox(width: 100, child: LinearProgressIndicator(value: c['percentage'] as double, backgroundColor: (c['color'] as Color).withValues(alpha: 0.2), valueColor: AlwaysStoppedAnimation(c['color'] as Color))),
        ]),
      )).toList()),
    );
  }

  Widget _buildWarehouseCards() {
    final warehouses = [
      {'name': 'Main Warehouse - Apapa', 'items': 456, 'value': '₦78M', 'capacity': 0.72},
      {'name': 'Distribution Center - Ikeja', 'items': 234, 'value': '₦32M', 'capacity': 0.54},
      {'name': 'Satellite Store - VI', 'items': 157, 'value': '₦15M', 'capacity': 0.89},
    ];
    return Column(children: warehouses.map((w) => Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Text(w['name'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
          Text(w['value'] as String, style: AppTypography.titleSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
        ]),
        const SizedBox(height: 8),
        Row(children: [
          Text('${w['items']} items', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
          const Spacer(),
          Text('${((w['capacity'] as double) * 100).toInt()}% capacity', style: AppTypography.labelSmall.copyWith(color: (w['capacity'] as double) > 0.8 ? AppColors.warning : AppColors.success)),
        ]),
        const SizedBox(height: 8),
        LinearProgressIndicator(value: w['capacity'] as double, backgroundColor: AppColors.borderColor, valueColor: AlwaysStoppedAnimation((w['capacity'] as double) > 0.8 ? AppColors.warning : AppColors.success)),
      ]),
    )).toList());
  }

  Widget _buildRecentActivity() {
    final activities = [
      {'action': 'Stock In', 'product': 'Golden Penny Semovita 5kg', 'qty': '+100', 'time': '2 hours ago'},
      {'action': 'Stock Out', 'product': 'Indomie Chicken 70g', 'qty': '-50', 'time': '3 hours ago'},
      {'action': 'Adjustment', 'product': 'Peak Milk 400g', 'qty': '-5', 'time': '5 hours ago'},
    ];
    return Column(children: activities.map((a) => ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(
        backgroundColor: a['action'] == 'Stock In' ? AppColors.success.withValues(alpha: 0.1) : a['action'] == 'Stock Out' ? AppColors.error.withValues(alpha: 0.1) : AppColors.warning.withValues(alpha: 0.1),
        child: Icon(a['action'] == 'Stock In' ? Icons.add : a['action'] == 'Stock Out' ? Icons.remove : Icons.edit, color: a['action'] == 'Stock In' ? AppColors.success : a['action'] == 'Stock Out' ? AppColors.error : AppColors.warning, size: 20),
      ),
      title: Text(a['product']!, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
      subtitle: Text('${a['action']} • ${a['time']}', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      trailing: Text(a['qty']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600, color: a['qty']!.startsWith('+') ? AppColors.success : AppColors.error)),
    )).toList());
  }

  Widget _buildStockLevelsTab() {
    final items = [
      {'name': 'Golden Penny Semovita 5kg', 'sku': 'GPS-5KG', 'onHand': 45, 'reserved': 10, 'available': 35, 'reorder': 20},
      {'name': 'Peak Milk 400g', 'sku': 'PEM-400', 'onHand': 120, 'reserved': 0, 'available': 120, 'reorder': 50},
      {'name': 'Indomie Chicken 70g', 'sku': 'IND-70C', 'onHand': 8, 'reserved': 5, 'available': 3, 'reorder': 25},
    ];
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: items.length,
      itemBuilder: (context, index) {
        final item = items[index];
        final isLow = (item['available'] as int) <= (item['reorder'] as int);
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, border: isLow ? Border.all(color: AppColors.warning) : null),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Expanded(child: Text(item['name'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600))),
              if (isLow) Container(padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2), decoration: BoxDecoration(color: AppColors.warning, borderRadius: AppRadius.borderRadiusSm), child: const Text('LOW', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold))),
            ]),
            Text(item['sku'] as String, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            const SizedBox(height: 12),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              _buildStockColumn('On Hand', '${item['onHand']}'),
              _buildStockColumn('Reserved', '${item['reserved']}'),
              _buildStockColumn('Available', '${item['available']}'),
              _buildStockColumn('Reorder At', '${item['reorder']}'),
            ]),
          ]),
        );
      },
    );
  }

  Widget _buildStockColumn(String label, String value) {
    return Column(children: [
      Text(value, style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.w700)),
      Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
    ]);
  }

  Widget _buildMovementsTab() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 10,
      itemBuilder: (context, index) => Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
        child: Row(children: [
          CircleAvatar(backgroundColor: index % 2 == 0 ? AppColors.success.withValues(alpha: 0.1) : AppColors.error.withValues(alpha: 0.1),
            child: Icon(index % 2 == 0 ? Icons.arrow_downward : Icons.arrow_upward, color: index % 2 == 0 ? AppColors.success : AppColors.error)),
          const SizedBox(width: 12),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text('Product ${index + 1}', style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
            Text(index % 2 == 0 ? 'Received from Supplier' : 'Dispatched to Customer', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
          ])),
          Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
            Text('${index % 2 == 0 ? '+' : '-'}${(index + 1) * 10}', style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600, color: index % 2 == 0 ? AppColors.success : AppColors.error)),
            Text('2 hours ago', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
          ]),
        ]),
      ),
    );
  }

  Widget _buildAlertsTab() {
    final alerts = [
      {'type': 'critical', 'title': 'Out of Stock', 'message': '8 products are out of stock', 'count': 8},
      {'type': 'warning', 'title': 'Low Stock', 'message': '23 products below reorder point', 'count': 23},
      {'type': 'info', 'title': 'Expiring Soon', 'message': '5 batches expiring in 30 days', 'count': 5},
    ];
    return ListView(
      padding: const EdgeInsets.all(16),
      children: alerts.map((a) => Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: a['type'] == 'critical' ? AppColors.error.withValues(alpha: 0.05) : a['type'] == 'warning' ? AppColors.warning.withValues(alpha: 0.05) : AppColors.info.withValues(alpha: 0.05),
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: a['type'] == 'critical' ? AppColors.error.withValues(alpha: 0.3) : a['type'] == 'warning' ? AppColors.warning.withValues(alpha: 0.3) : AppColors.info.withValues(alpha: 0.3)),
        ),
        child: Row(children: [
          Icon(a['type'] == 'critical' ? Icons.error : a['type'] == 'warning' ? Icons.warning : Icons.info,
            color: a['type'] == 'critical' ? AppColors.error : a['type'] == 'warning' ? AppColors.warning : AppColors.info, size: 32),
          const SizedBox(width: 16),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(a['title'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
            Text(a['message'] as String, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
          ])),
          TextButton(onPressed: () {}, child: const Text('View')),
        ]),
      )).toList(),
    );
  }
}
