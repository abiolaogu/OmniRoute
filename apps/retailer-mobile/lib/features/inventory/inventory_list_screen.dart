/// OmniRoute Ecosystem - Inventory Screen
/// Comprehensive inventory management with stock tracking, alerts, and analytics

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class InventoryListScreen extends ConsumerStatefulWidget {
  const InventoryListScreen({super.key});

  @override
  ConsumerState<InventoryListScreen> createState() => _InventoryListScreenState();
}

class _InventoryListScreenState extends ConsumerState<InventoryListScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final _searchController = TextEditingController();
  String _sortBy = 'name';
  bool _isGridView = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: _buildAppBar(),
      body: Column(
        children: [
          _buildStatsRow(),
          _buildSearchAndFilter(),
          _buildCategoryTabs(),
          Expanded(child: _buildInventoryList()),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showAddProductSheet(context),
        backgroundColor: AppColors.primary,
        icon: const Icon(Icons.add),
        label: const Text('Add Product'),
      ),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return AppBar(
      backgroundColor: AppColors.white,
      elevation: 0,
      title: Text(
        'Inventory',
        style: AppTypography.titleLarge.copyWith(color: AppColors.grey900),
      ),
      actions: [
        IconButton(
          icon: Icon(
            _isGridView ? Icons.view_list : Icons.grid_view,
            color: AppColors.grey800,
          ),
          onPressed: () => setState(() => _isGridView = !_isGridView),
        ),
        IconButton(
          icon: const Icon(Icons.qr_code_scanner, color: AppColors.grey800),
          onPressed: () {},
        ),
        PopupMenuButton<String>(
          icon: const Icon(Icons.more_vert, color: AppColors.grey800),
          onSelected: (value) {},
          itemBuilder: (context) => [
            const PopupMenuItem(value: 'export', child: Text('Export Inventory')),
            const PopupMenuItem(value: 'import', child: Text('Import Products')),
            const PopupMenuItem(value: 'bulk', child: Text('Bulk Update')),
            const PopupMenuItem(value: 'reports', child: Text('Inventory Reports')),
          ],
        ),
      ],
    );
  }

  Widget _buildStatsRow() {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Row(
        children: [
          Expanded(
            child: _StockStatCard(
              label: 'Total Products',
              value: '1,248',
              icon: Icons.inventory_2,
              color: AppColors.primary,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: _StockStatCard(
              label: 'Low Stock',
              value: '23',
              icon: Icons.warning_amber,
              color: AppColors.warning,
              onTap: () {},
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: _StockStatCard(
              label: 'Out of Stock',
              value: '8',
              icon: Icons.error_outline,
              color: AppColors.error,
              onTap: () {},
            ),
          ),
        ],
      ).animate().fadeIn().slideY(begin: 0.1, end: 0),
    );
  }

  Widget _buildSearchAndFilter() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Row(
        children: [
          Expanded(
            child: Container(
              height: 48,
              decoration: BoxDecoration(
                color: AppColors.white,
                borderRadius: AppRadius.borderRadiusMd,
                border: Border.all(color: AppColors.grey200),
              ),
              child: Row(
                children: [
                  const SizedBox(width: 12),
                  const Icon(Icons.search, color: AppColors.grey500, size: 20),
                  const SizedBox(width: 8),
                  Expanded(
                    child: TextField(
                      controller: _searchController,
                      decoration: InputDecoration(
                        hintText: 'Search products...',
                        border: InputBorder.none,
                        hintStyle: AppTypography.bodyMedium.copyWith(
                          color: AppColors.grey500,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(width: 12),
          Container(
            height: 48,
            decoration: BoxDecoration(
              color: AppColors.white,
              borderRadius: AppRadius.borderRadiusMd,
              border: Border.all(color: AppColors.grey200),
            ),
            child: PopupMenuButton<String>(
              initialValue: _sortBy,
              onSelected: (value) => setState(() => _sortBy = value),
              offset: const Offset(0, 48),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Row(
                  children: [
                    const Icon(Icons.sort, color: AppColors.grey700, size: 20),
                    const SizedBox(width: 4),
                    Text('Sort', style: AppTypography.labelMedium),
                    const Icon(Icons.arrow_drop_down, color: AppColors.grey700),
                  ],
                ),
              ),
              itemBuilder: (context) => [
                const PopupMenuItem(value: 'name', child: Text('Name (A-Z)')),
                const PopupMenuItem(value: 'stock_low', child: Text('Stock (Low to High)')),
                const PopupMenuItem(value: 'stock_high', child: Text('Stock (High to Low)')),
                const PopupMenuItem(value: 'recent', child: Text('Recently Added')),
                const PopupMenuItem(value: 'price', child: Text('Price')),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCategoryTabs() {
    return Container(
      margin: const EdgeInsets.only(top: 16),
      child: TabBar(
        controller: _tabController,
        isScrollable: true,
        labelColor: AppColors.primary,
        unselectedLabelColor: AppColors.grey600,
        indicatorColor: AppColors.primary,
        indicatorWeight: 3,
        labelPadding: const EdgeInsets.symmetric(horizontal: 20),
        tabs: const [
          Tab(text: 'All'),
          Tab(text: 'Groceries'),
          Tab(text: 'Beverages'),
          Tab(text: 'Electronics'),
        ],
      ),
    );
  }

  Widget _buildInventoryList() {
    // Mock inventory data
    final items = List.generate(20, (index) => {
      'id': 'PRD-${1000 + index}',
      'name': [
        'Golden Penny Semovita',
        'Indomie Chicken Noodles',
        'Peak Milk Powder',
        'Coca-Cola 50cl',
        'Nestle Milo 400g',
        'Dangote Sugar 1kg',
        'Kings Oil 5L',
        'Dano Milk 900g',
      ][index % 8],
      'sku': 'SKU-${10000 + index}',
      'stock': [45, 12, 5, 0, 78, 156, 23, 8][index % 8],
      'price': [2500.0, 180.0, 4500.0, 200.0, 2800.0, 1200.0, 8500.0, 5200.0][index % 8],
      'reorderLevel': 20,
      'category': ['Food', 'Food', 'Dairy', 'Beverages', 'Beverages', 'Food', 'Food', 'Dairy'][index % 8],
    });

    if (_isGridView) {
      return GridView.builder(
        padding: const EdgeInsets.all(16),
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          childAspectRatio: 0.75,
        ),
        itemCount: items.length,
        itemBuilder: (context, index) {
          return _InventoryGridCard(item: items[index])
              .animate(delay: Duration(milliseconds: 50 * (index % 10)))
              .fadeIn()
              .scale(begin: const Offset(0.95, 0.95));
        },
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: items.length,
      itemBuilder: (context, index) {
        return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: _InventoryListCard(item: items[index])
              .animate(delay: Duration(milliseconds: 50 * (index % 10)))
              .fadeIn()
              .slideX(begin: 0.05, end: 0),
        );
      },
    );
  }

  void _showAddProductSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.85,
        decoration: const BoxDecoration(
          color: AppColors.white,
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        child: Column(
          children: [
            Container(
              margin: const EdgeInsets.only(top: 12),
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppColors.grey300,
                borderRadius: AppRadius.borderRadiusFull,
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(20),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Add New Product', style: AppTypography.titleLarge),
                  IconButton(
                    icon: const Icon(Icons.close),
                    onPressed: () => Navigator.pop(context),
                  ),
                ],
              ),
            ),
            const Divider(height: 1),
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Product image upload
                    Center(
                      child: Container(
                        width: 120,
                        height: 120,
                        decoration: BoxDecoration(
                          color: AppColors.grey100,
                          borderRadius: AppRadius.borderRadiusMd,
                          border: Border.all(
                            color: AppColors.grey300,
                            style: BorderStyle.solid,
                          ),
                        ),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            const Icon(Icons.add_photo_alternate, color: AppColors.grey500, size: 32),
                            const SizedBox(height: 4),
                            Text('Add Photo', style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(height: 24),
                    _buildFormField('Product Name', 'Enter product name'),
                    const SizedBox(height: 16),
                    _buildFormField('SKU', 'Enter SKU or barcode'),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(child: _buildFormField('Price', '₦0.00', prefix: '₦')),
                        const SizedBox(width: 16),
                        Expanded(child: _buildFormField('Cost', '₦0.00', prefix: '₦')),
                      ],
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(child: _buildFormField('Current Stock', '0')),
                        const SizedBox(width: 16),
                        Expanded(child: _buildFormField('Reorder Level', '20')),
                      ],
                    ),
                    const SizedBox(height: 16),
                    _buildDropdownField('Category', ['Food', 'Beverages', 'Dairy', 'Electronics', 'Other']),
                    const SizedBox(height: 16),
                    _buildFormField('Description', 'Enter product description', maxLines: 3),
                  ],
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(20),
              child: Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () => Navigator.pop(context),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: const Text('Cancel'),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    flex: 2,
                    child: ElevatedButton(
                      onPressed: () => Navigator.pop(context),
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: const Text('Add Product'),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFormField(String label, String hint, {String? prefix, int maxLines = 1}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: AppTypography.labelMedium.copyWith(color: AppColors.grey700)),
        const SizedBox(height: 8),
        TextField(
          maxLines: maxLines,
          decoration: InputDecoration(
            hintText: hint,
            prefixText: prefix,
            border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
          ),
        ),
      ],
    );
  }

  Widget _buildDropdownField(String label, List<String> options) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: AppTypography.labelMedium.copyWith(color: AppColors.grey700)),
        const SizedBox(height: 8),
        DropdownButtonFormField<String>(
          decoration: InputDecoration(
            border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
          ),
          items: options.map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
          onChanged: (value) {},
          hint: const Text('Select category'),
        ),
      ],
    );
  }
}

class _StockStatCard extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;
  final Color color;
  final VoidCallback? onTap;

  const _StockStatCard({
    required this.label,
    required this.value,
    required this.icon,
    required this.color,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: AppColors.white,
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: AppColors.cardBorder),
        ),
        child: Column(
          children: [
            Icon(icon, color: color, size: 24),
            const SizedBox(height: 8),
            Text(
              value,
              style: AppTypography.headlineSmall.copyWith(
                color: AppColors.grey900,
                fontWeight: FontWeight.w700,
              ),
            ),
            Text(
              label,
              style: AppTypography.labelSmall.copyWith(color: AppColors.grey600),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}

class _InventoryListCard extends StatelessWidget {
  final Map<String, dynamic> item;

  const _InventoryListCard({required this.item});

  @override
  Widget build(BuildContext context) {
    final stock = item['stock'] as int;
    final reorderLevel = item['reorderLevel'] as int;
    final isLowStock = stock > 0 && stock <= reorderLevel;
    final isOutOfStock = stock == 0;

    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(
          color: isOutOfStock
              ? AppColors.error.withValues(alpha: 0.3)
              : isLowStock
                  ? AppColors.warning.withValues(alpha: 0.3)
                  : AppColors.cardBorder,
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 56,
            height: 56,
            decoration: BoxDecoration(
              color: AppColors.grey100,
              borderRadius: AppRadius.borderRadiusSm,
            ),
            child: const Icon(Icons.inventory_2, color: AppColors.grey500),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(item['name'] as String, style: AppTypography.titleSmall),
                const SizedBox(height: 2),
                Text(
                  '${item['sku']} • ${item['category']}',
                  style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
                ),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                formatCurrency(item['price'] as double),
                style: AppTypography.titleSmall,
              ),
              const SizedBox(height: 4),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: isOutOfStock
                      ? AppColors.errorBg
                      : isLowStock
                          ? AppColors.warningBg
                          : AppColors.successBg,
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Text(
                  isOutOfStock ? 'Out of Stock' : '$stock in stock',
                  style: AppTypography.labelSmall.copyWith(
                    color: isOutOfStock
                        ? AppColors.error
                        : isLowStock
                            ? AppColors.warning
                            : AppColors.success,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _InventoryGridCard extends StatelessWidget {
  final Map<String, dynamic> item;

  const _InventoryGridCard({required this.item});

  @override
  Widget build(BuildContext context) {
    final stock = item['stock'] as int;
    final reorderLevel = item['reorderLevel'] as int;
    final isLowStock = stock > 0 && stock <= reorderLevel;
    final isOutOfStock = stock == 0;

    return Container(
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Stack(
            children: [
              Container(
                height: 100,
                width: double.infinity,
                decoration: BoxDecoration(
                  color: AppColors.grey100,
                  borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
                ),
                child: const Icon(Icons.inventory_2, size: 40, color: AppColors.grey400),
              ),
              if (isOutOfStock || isLowStock)
                Positioned(
                  top: 8,
                  right: 8,
                  child: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: isOutOfStock ? AppColors.error : AppColors.warning,
                      borderRadius: AppRadius.borderRadiusFull,
                    ),
                    child: Text(
                      isOutOfStock ? 'Out' : 'Low',
                      style: AppTypography.labelSmall.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
            ],
          ),
          Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  item['name'] as String,
                  style: AppTypography.titleSmall,
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 4),
                Text(
                  item['sku'] as String,
                  style: AppTypography.bodySmall.copyWith(color: AppColors.grey500),
                ),
                const SizedBox(height: 8),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      formatCurrency(item['price'] as double),
                      style: AppTypography.titleSmall.copyWith(color: AppColors.success),
                    ),
                    Text(
                      '$stock pcs',
                      style: AppTypography.labelSmall.copyWith(
                        color: isOutOfStock
                            ? AppColors.error
                            : isLowStock
                                ? AppColors.warning
                                : AppColors.grey600,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
