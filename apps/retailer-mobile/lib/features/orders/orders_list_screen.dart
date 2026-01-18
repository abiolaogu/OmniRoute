/// OmniRoute Ecosystem - Orders Screen
/// Comprehensive order management with filtering, search, and status updates

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/core/router/app_router.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class OrdersListScreen extends ConsumerStatefulWidget {
  const OrdersListScreen({super.key});

  @override
  ConsumerState<OrdersListScreen> createState() => _OrdersListScreenState();
}

class _OrdersListScreenState extends ConsumerState<OrdersListScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final _searchController = TextEditingController();
  String _selectedFilter = 'all';
  bool _isSearching = false;

  final _orderStatuses = [
    {'key': 'all', 'label': 'All Orders', 'count': 156},
    {'key': 'pending', 'label': 'Pending', 'count': 23},
    {'key': 'processing', 'label': 'Processing', 'count': 45},
    {'key': 'shipped', 'label': 'Shipped', 'count': 32},
    {'key': 'delivered', 'label': 'Delivered', 'count': 48},
    {'key': 'cancelled', 'label': 'Cancelled', 'count': 8},
  ];

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: _orderStatuses.length, vsync: this);
    // Load orders on init
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(ordersProvider.notifier).loadOrders(refresh: true);
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final ordersState = ref.watch(ordersProvider);

    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: _buildAppBar(),
      body: Column(
        children: [
          // Filter Tabs
          _buildFilterTabs(),
          // Search and Filter Bar
          if (_isSearching) _buildSearchBar(),
          // Orders List
          Expanded(
            child: ordersState.isLoading && ordersState.orders.isEmpty
                ? _buildLoadingState()
                : ordersState.orders.isEmpty
                    ? _buildEmptyState()
                    : _buildOrdersList(ordersState),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.primary,
        icon: const Icon(Icons.add),
        label: const Text('New Order'),
      ),
    );
  }

  PreferredSizeWidget _buildAppBar() {
    return AppBar(
      backgroundColor: AppColors.white,
      elevation: 0,
      title: _isSearching
          ? TextField(
              controller: _searchController,
              autofocus: true,
              decoration: InputDecoration(
                hintText: 'Search orders...',
                border: InputBorder.none,
                hintStyle: AppTypography.bodyMedium.copyWith(
                  color: AppColors.grey500,
                ),
              ),
              onChanged: (value) {
                // Implement search
              },
            )
          : Text(
              'Orders',
              style: AppTypography.titleLarge.copyWith(
                color: AppColors.grey900,
              ),
            ),
      leading: IconButton(
        icon: Icon(
          _isSearching ? Icons.arrow_back : Icons.menu,
          color: AppColors.grey800,
        ),
        onPressed: () {
          if (_isSearching) {
            setState(() {
              _isSearching = false;
              _searchController.clear();
            });
          } else {
            Scaffold.of(context).openDrawer();
          }
        },
      ),
      actions: [
        IconButton(
          icon: Icon(
            _isSearching ? Icons.close : Icons.search,
            color: AppColors.grey800,
          ),
          onPressed: () {
            setState(() {
              _isSearching = !_isSearching;
              if (!_isSearching) {
                _searchController.clear();
              }
            });
          },
        ),
        IconButton(
          icon: const Icon(Icons.filter_list, color: AppColors.grey800),
          onPressed: () => _showFilterBottomSheet(context),
        ),
      ],
    );
  }

  Widget _buildFilterTabs() {
    return Container(
      color: AppColors.white,
      child: TabBar(
        controller: _tabController,
        isScrollable: true,
        labelColor: AppColors.primary,
        unselectedLabelColor: AppColors.grey600,
        indicatorColor: AppColors.primary,
        indicatorWeight: 3,
        labelStyle: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600),
        unselectedLabelStyle: AppTypography.labelMedium,
        tabs: _orderStatuses.map((status) {
          return Tab(
            child: Row(
              children: [
                Text(status['label'] as String),
                const SizedBox(width: 6),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: _tabController.index ==
                            _orderStatuses.indexOf(status)
                        ? AppColors.primary.withValues(alpha: 0.1)
                        : AppColors.grey200,
                    borderRadius: AppRadius.borderRadiusFull,
                  ),
                  child: Text(
                    '${status['count']}',
                    style: AppTypography.labelSmall.copyWith(
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
          );
        }).toList(),
        onTap: (index) {
          setState(() {
            _selectedFilter = _orderStatuses[index]['key'] as String;
          });
        },
      ),
    );
  }

  Widget _buildSearchBar() {
    return Container(
      padding: const EdgeInsets.all(16),
      color: AppColors.white,
      child: Row(
        children: [
          Expanded(
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              decoration: BoxDecoration(
                color: AppColors.grey50,
                borderRadius: AppRadius.borderRadiusMd,
                border: Border.all(color: AppColors.grey200),
              ),
              child: Row(
                children: [
                  const Icon(Icons.search, color: AppColors.grey500, size: 20),
                  const SizedBox(width: 12),
                  Expanded(
                    child: TextField(
                      controller: _searchController,
                      decoration: InputDecoration(
                        hintText: 'Search by order ID, customer...',
                        border: InputBorder.none,
                        hintStyle: AppTypography.bodyMedium.copyWith(
                          color: AppColors.grey500,
                        ),
                      ),
                    ),
                  ),
                  if (_searchController.text.isNotEmpty)
                    GestureDetector(
                      onTap: () => _searchController.clear(),
                      child: const Icon(Icons.close, color: AppColors.grey500, size: 20),
                    ),
                ],
              ),
            ),
          ),
          const SizedBox(width: 12),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: AppColors.primary,
              borderRadius: AppRadius.borderRadiusMd,
            ),
            child: const Icon(Icons.qr_code_scanner, color: Colors.white, size: 24),
          ),
        ],
      ),
    ).animate().fadeIn().slideY(begin: -0.2, end: 0);
  }

  Widget _buildLoadingState() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 5,
      itemBuilder: (context, index) => Padding(
        padding: const EdgeInsets.only(bottom: 12),
        child: const _OrderCardShimmer(),
      ),
    );
  }

  Widget _buildEmptyState() {
    return EmptyState(
      icon: Icons.receipt_long,
      title: 'No orders yet',
      subtitle: 'Your orders will appear here once you start receiving them.',
      actionText: 'Create Order',
      onAction: () {},
    );
  }

  Widget _buildOrdersList(OrdersState ordersState) {
    return RefreshIndicator(
      onRefresh: () => ref.read(ordersProvider.notifier).loadOrders(refresh: true),
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: ordersState.orders.length + (ordersState.hasMore ? 1 : 0),
        itemBuilder: (context, index) {
          if (index == ordersState.orders.length) {
            // Load more indicator
            if (!ordersState.isLoading) {
              ref.read(ordersProvider.notifier).loadOrders();
            }
            return const Center(
              child: Padding(
                padding: EdgeInsets.all(16),
                child: CircularProgressIndicator(),
              ),
            );
          }

          final order = ordersState.orders[index];
          return Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: _OrderCard(
              order: order,
              onTap: () => context.push('/orders/${order.id}'),
            ).animate(delay: Duration(milliseconds: 50 * (index % 10))).fadeIn().slideX(begin: 0.05, end: 0),
          );
        },
      ),
    );
  }

  void _showFilterBottomSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.6,
        decoration: const BoxDecoration(
          color: AppColors.white,
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        child: Column(
          children: [
            // Handle
            Container(
              margin: const EdgeInsets.only(top: 12),
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: AppColors.grey300,
                borderRadius: AppRadius.borderRadiusFull,
              ),
            ),
            // Header
            Padding(
              padding: const EdgeInsets.all(20),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Filter Orders', style: AppTypography.titleLarge),
                  TextButton(
                    onPressed: () {
                      // Reset filters
                      Navigator.pop(context);
                    },
                    child: Text(
                      'Reset',
                      style: AppTypography.labelLarge.copyWith(
                        color: AppColors.primary,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const Divider(height: 1),
            // Filter options
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Date Range', style: AppTypography.titleSmall),
                    const SizedBox(height: 12),
                    Wrap(
                      spacing: 8,
                      runSpacing: 8,
                      children: [
                        _FilterChip(label: 'Today', isSelected: true),
                        _FilterChip(label: 'Yesterday'),
                        _FilterChip(label: 'Last 7 days'),
                        _FilterChip(label: 'Last 30 days'),
                        _FilterChip(label: 'Custom'),
                      ],
                    ),
                    const SizedBox(height: 24),
                    Text('Payment Status', style: AppTypography.titleSmall),
                    const SizedBox(height: 12),
                    Wrap(
                      spacing: 8,
                      runSpacing: 8,
                      children: [
                        _FilterChip(label: 'All', isSelected: true),
                        _FilterChip(label: 'Paid'),
                        _FilterChip(label: 'Pending'),
                        _FilterChip(label: 'Failed'),
                      ],
                    ),
                    const SizedBox(height: 24),
                    Text('Amount Range', style: AppTypography.titleSmall),
                    const SizedBox(height: 12),
                    Row(
                      children: [
                        Expanded(
                          child: TextField(
                            decoration: InputDecoration(
                              labelText: 'Min',
                              prefixText: '₦',
                              border: OutlineInputBorder(
                                borderRadius: AppRadius.borderRadiusMd,
                              ),
                            ),
                            keyboardType: TextInputType.number,
                          ),
                        ),
                        const Padding(
                          padding: EdgeInsets.symmetric(horizontal: 12),
                          child: Text('-'),
                        ),
                        Expanded(
                          child: TextField(
                            decoration: InputDecoration(
                              labelText: 'Max',
                              prefixText: '₦',
                              border: OutlineInputBorder(
                                borderRadius: AppRadius.borderRadiusMd,
                              ),
                            ),
                            keyboardType: TextInputType.number,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            // Apply button
            Padding(
              padding: const EdgeInsets.all(20),
              child: SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: () => Navigator.pop(context),
                  child: const Text('Apply Filters'),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class _OrderCard extends StatelessWidget {
  final dynamic order; // Would be Order type
  final VoidCallback onTap;

  const _OrderCard({required this.order, required this.onTap});

  @override
  Widget build(BuildContext context) {
    // Mock data since we're using dynamic
    final orderNumber = order.orderNumber ?? 'ORD-10001';
    final customerName = order.customerName ?? 'Customer Name';
    final total = order.total ?? 125000.0;
    final status = order.status ?? 'pending';
    final itemCount = order.items?.length ?? 3;
    final createdAt = order.createdAt ?? DateTime.now();

    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: AppColors.white,
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: AppColors.cardBorder),
          boxShadow: AppShadows.sm,
        ),
        child: Column(
          children: [
            // Header row
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      orderNumber,
                      style: AppTypography.titleSmall.copyWith(
                        color: AppColors.primary,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      timeAgo(createdAt),
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.grey500,
                      ),
                    ),
                  ],
                ),
                StatusChip(status: status),
              ],
            ),
            const SizedBox(height: 12),
            const Divider(height: 1),
            const SizedBox(height: 12),
            // Customer info
            Row(
              children: [
                Container(
                  width: 40,
                  height: 40,
                  decoration: BoxDecoration(
                    color: AppColors.grey100,
                    borderRadius: AppRadius.borderRadiusSm,
                  ),
                  child: const Center(
                    child: Icon(Icons.person, color: AppColors.grey500, size: 20),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(customerName, style: AppTypography.titleSmall),
                      Text(
                        '$itemCount items',
                        style: AppTypography.bodySmall.copyWith(
                          color: AppColors.grey600,
                        ),
                      ),
                    ],
                  ),
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      formatCurrency(total),
                      style: AppTypography.titleMedium.copyWith(
                        color: AppColors.grey900,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    Text(
                      'Total',
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.grey500,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 12),
            // Action buttons
            Row(
              children: [
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: () {},
                    icon: const Icon(Icons.visibility, size: 18),
                    label: const Text('View'),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: AppColors.grey700,
                      side: const BorderSide(color: AppColors.grey300),
                      padding: const EdgeInsets.symmetric(vertical: 8),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: () {},
                    icon: const Icon(Icons.print, size: 18),
                    label: const Text('Print'),
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppColors.primary,
                      foregroundColor: Colors.white,
                      padding: const EdgeInsets.symmetric(vertical: 8),
                      elevation: 0,
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _OrderCardShimmer extends StatelessWidget {
  const _OrderCardShimmer();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: const [
              ShimmerLoading(width: 100, height: 20),
              ShimmerLoading(width: 80, height: 24, borderRadius: AppRadius.borderRadiusFull),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: const [
              ShimmerLoading(width: 40, height: 40),
              SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ShimmerLoading(width: 120, height: 16),
                    SizedBox(height: 4),
                    ShimmerLoading(width: 80, height: 12),
                  ],
                ),
              ),
              ShimmerLoading(width: 80, height: 20),
            ],
          ),
        ],
      ),
    );
  }
}

class _FilterChip extends StatelessWidget {
  final String label;
  final bool isSelected;

  const _FilterChip({required this.label, this.isSelected = false});

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
      decoration: BoxDecoration(
        color: isSelected ? AppColors.primary : AppColors.white,
        borderRadius: AppRadius.borderRadiusFull,
        border: Border.all(
          color: isSelected ? AppColors.primary : AppColors.grey300,
        ),
      ),
      child: Text(
        label,
        style: AppTypography.labelMedium.copyWith(
          color: isSelected ? Colors.white : AppColors.grey700,
          fontWeight: isSelected ? FontWeight.w600 : FontWeight.w500,
        ),
      ),
    );
  }
}
