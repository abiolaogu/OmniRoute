/// OmniRoute Ecosystem - Products Screen
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class ProductsScreen extends ConsumerStatefulWidget {
  const ProductsScreen({super.key});
  @override ConsumerState<ProductsScreen> createState() => _ProductsScreenState();
}

class _ProductsScreenState extends ConsumerState<ProductsScreen> {
  String _selectedCategory = 'All';
  bool _isGridView = true;

  final _categories = ['All', 'Food', 'Beverages', 'Personal Care', 'Household', 'Electronics'];

  final _products = [
    {'name': 'Golden Penny Semovita 5kg', 'sku': 'GPS-5KG', 'price': '₦4,200', 'stock': 45, 'category': 'Food', 'image': null},
    {'name': 'Peak Evaporated Milk 400g', 'sku': 'PEM-400', 'price': '₦850', 'stock': 120, 'category': 'Beverages', 'image': null},
    {'name': 'Indomie Chicken 70g (Carton)', 'sku': 'IND-70C', 'price': '₦6,500', 'stock': 8, 'category': 'Food', 'image': null},
    {'name': 'Dettol Soap 175g', 'sku': 'DET-175', 'price': '₦450', 'stock': 200, 'category': 'Personal Care', 'image': null},
    {'name': 'Kings Vegetable Oil 5L', 'sku': 'KVO-5L', 'price': '₦8,900', 'stock': 32, 'category': 'Food', 'image': null},
    {'name': 'Morning Fresh 750ml', 'sku': 'MF-750', 'price': '₦1,200', 'stock': 0, 'category': 'Household', 'image': null},
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('Products'),
        actions: [
          IconButton(icon: const Icon(Icons.search), onPressed: () => _showSearchDialog()),
          IconButton(icon: const Icon(Icons.filter_list), onPressed: () => _showFilterSheet()),
          IconButton(icon: Icon(_isGridView ? Icons.list : Icons.grid_view), onPressed: () => setState(() => _isGridView = !_isGridView)),
        ],
      ),
      body: Column(
        children: [
          _buildCategoryChips(),
          Expanded(child: _isGridView ? _buildProductGrid() : _buildProductList()),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
        backgroundColor: AppColors.primary,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildCategoryChips() {
    return Container(
      height: 56,
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: 16),
        itemCount: _categories.length,
        separatorBuilder: (_, __) => const SizedBox(width: 8),
        itemBuilder: (context, index) {
          final cat = _categories[index];
          final isSelected = cat == _selectedCategory;
          return ChoiceChip(
            label: Text(cat),
            selected: isSelected,
            onSelected: (selected) => setState(() => _selectedCategory = cat),
            selectedColor: AppColors.primary,
            labelStyle: TextStyle(color: isSelected ? Colors.white : AppColors.textPrimary),
          );
        },
      ),
    );
  }

  Widget _buildProductGrid() {
    final filtered = _getFilteredProducts();
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(crossAxisCount: 2, crossAxisSpacing: 12, mainAxisSpacing: 12, childAspectRatio: 0.75),
      itemCount: filtered.length,
      itemBuilder: (context, index) => _buildProductCard(filtered[index]),
    );
  }

  Widget _buildProductList() {
    final filtered = _getFilteredProducts();
    return ListView.separated(
      padding: const EdgeInsets.all(16),
      itemCount: filtered.length,
      separatorBuilder: (_, __) => const SizedBox(height: 12),
      itemBuilder: (context, index) => _buildProductListItem(filtered[index]),
    );
  }

  Widget _buildProductCard(Map<String, dynamic> product) {
    final isLowStock = (product['stock'] as int) < 10;
    final isOutOfStock = (product['stock'] as int) == 0;
    return Container(
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Stack(children: [
            Container(
              height: 100,
              decoration: BoxDecoration(color: AppColors.borderColor, borderRadius: const BorderRadius.vertical(top: Radius.circular(12))),
              child: const Center(child: Icon(Icons.image, size: 40, color: AppColors.textSecondary)),
            ),
            if (isOutOfStock)
              Positioned.fill(child: Container(
                decoration: BoxDecoration(color: Colors.black.withValues(alpha: 0.5), borderRadius: const BorderRadius.vertical(top: Radius.circular(12))),
                child: const Center(child: Text('OUT OF STOCK', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 12))),
              )),
            if (isLowStock && !isOutOfStock)
              Positioned(top: 8, left: 8, child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(color: AppColors.warning, borderRadius: AppRadius.borderRadiusSm),
                child: const Text('LOW STOCK', style: TextStyle(color: Colors.white, fontSize: 8, fontWeight: FontWeight.bold)),
              )),
          ]),
          Padding(
            padding: const EdgeInsets.all(12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(product['name'], style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600), maxLines: 2, overflow: TextOverflow.ellipsis),
                const SizedBox(height: 4),
                Text(product['sku'], style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                const SizedBox(height: 8),
                Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
                  Text(product['price'], style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w700, color: AppColors.primary)),
                  Text('${product['stock']} in stock', style: AppTypography.labelSmall.copyWith(color: isLowStock ? AppColors.warning : AppColors.success)),
                ]),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildProductListItem(Map<String, dynamic> product) {
    final isLowStock = (product['stock'] as int) < 10;
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Row(
        children: [
          Container(
            width: 64, height: 64,
            decoration: BoxDecoration(color: AppColors.borderColor, borderRadius: AppRadius.borderRadiusSm),
            child: const Icon(Icons.image, color: AppColors.textSecondary),
          ),
          const SizedBox(width: 12),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(product['name'], style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 4),
            Text('${product['sku']} • ${product['category']}', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
          ])),
          Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
            Text(product['price'], style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w700, color: AppColors.primary)),
            const SizedBox(height: 4),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(color: isLowStock ? AppColors.warning.withValues(alpha: 0.1) : AppColors.success.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
              child: Text('${product['stock']}', style: TextStyle(color: isLowStock ? AppColors.warning : AppColors.success, fontWeight: FontWeight.w600, fontSize: 12)),
            ),
          ]),
        ],
      ),
    );
  }

  List<Map<String, dynamic>> _getFilteredProducts() {
    if (_selectedCategory == 'All') return _products;
    return _products.where((p) => p['category'] == _selectedCategory).toList();
  }

  void _showSearchDialog() {
    showSearch(context: context, delegate: ProductSearchDelegate(products: _products));
  }

  void _showFilterSheet() {
    showModalBottomSheet(
      context: context,
      builder: (context) => Container(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Filter Products', style: AppTypography.titleLarge.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 24),
            Text('Stock Status', style: AppTypography.titleSmall),
            const SizedBox(height: 8),
            Wrap(spacing: 8, children: [
              FilterChip(label: const Text('In Stock'), selected: true, onSelected: (_) {}),
              FilterChip(label: const Text('Low Stock'), selected: false, onSelected: (_) {}),
              FilterChip(label: const Text('Out of Stock'), selected: false, onSelected: (_) {}),
            ]),
            const SizedBox(height: 16),
            Text('Price Range', style: AppTypography.titleSmall),
            const SizedBox(height: 8),
            RangeSlider(values: const RangeValues(0, 10000), min: 0, max: 20000, onChanged: (_) {}),
            const SizedBox(height: 24),
            SizedBox(width: double.infinity, child: ElevatedButton(onPressed: () => Navigator.pop(context), child: const Text('Apply Filters'))),
          ],
        ),
      ),
    );
  }
}

class ProductSearchDelegate extends SearchDelegate {
  final List<Map<String, dynamic>> products;
  ProductSearchDelegate({required this.products});

  @override List<Widget> buildActions(BuildContext context) => [IconButton(icon: const Icon(Icons.clear), onPressed: () => query = '')];
  @override Widget buildLeading(BuildContext context) => IconButton(icon: const Icon(Icons.arrow_back), onPressed: () => close(context, null));
  @override Widget buildResults(BuildContext context) => _buildSearchResults();
  @override Widget buildSuggestions(BuildContext context) => _buildSearchResults();

  Widget _buildSearchResults() {
    final results = products.where((p) => p['name'].toString().toLowerCase().contains(query.toLowerCase())).toList();
    return ListView.builder(
      itemCount: results.length,
      itemBuilder: (context, index) => ListTile(
        leading: Container(width: 48, height: 48, decoration: BoxDecoration(color: AppColors.borderColor, borderRadius: AppRadius.borderRadiusSm), child: const Icon(Icons.image)),
        title: Text(results[index]['name']),
        subtitle: Text('${results[index]['sku']} • ${results[index]['price']}'),
        trailing: Text('${results[index]['stock']} in stock'),
        onTap: () => close(context, results[index]),
      ),
    );
  }
}
