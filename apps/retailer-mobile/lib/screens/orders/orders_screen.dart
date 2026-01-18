/// OmniRoute Ecosystem - Orders Screen
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class OrdersScreen extends ConsumerStatefulWidget {
  const OrdersScreen({super.key});
  @override ConsumerState<OrdersScreen> createState() => _OrdersScreenState();
}

class _OrdersScreenState extends ConsumerState<OrdersScreen> with SingleTickerProviderStateMixin {
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
        title: const Text('Orders'),
        bottom: TabBar(
          controller: _tabController,
          labelColor: AppColors.primary,
          unselectedLabelColor: AppColors.textSecondary,
          indicatorColor: AppColors.primary,
          tabs: const [
            Tab(text: 'All'),
            Tab(text: 'Pending'),
            Tab(text: 'Processing'),
            Tab(text: 'Completed'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildOrderList('all'),
          _buildOrderList('pending'),
          _buildOrderList('processing'),
          _buildOrderList('completed'),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.primary,
        icon: const Icon(Icons.add, color: Colors.white),
        label: const Text('New Order', style: TextStyle(color: Colors.white)),
      ),
    );
  }

  Widget _buildOrderList(String filter) {
    final orders = _getFilteredOrders(filter);
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: orders.length,
      itemBuilder: (context, index) => _buildOrderCard(orders[index]),
    );
  }

  Widget _buildOrderCard(Map<String, dynamic> order) {
    final statusColor = _getStatusColor(order['status']);
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: AppRadius.borderRadiusMd,
        boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)],
      ),
      child: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: statusColor.withValues(alpha: 0.05),
              borderRadius: const BorderRadius.vertical(top: Radius.circular(12)),
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(children: [
                  Icon(Icons.receipt_long, color: statusColor, size: 20),
                  const SizedBox(width: 8),
                  Text(order['id'], style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
                ]),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                  decoration: BoxDecoration(color: statusColor, borderRadius: AppRadius.borderRadiusSm),
                  child: Text(order['status'], style: const TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w600)),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                Row(
                  children: [
                    CircleAvatar(
                      backgroundColor: AppColors.primary.withValues(alpha: 0.1),
                      child: Text(order['customer'][0], style: TextStyle(color: AppColors.primary, fontWeight: FontWeight.bold)),
                    ),
                    const SizedBox(width: 12),
                    Expanded(child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(order['customer'], style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
                        Text(order['address'], style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary), maxLines: 1, overflow: TextOverflow.ellipsis),
                      ],
                    )),
                  ],
                ),
                const SizedBox(height: 16),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    _buildOrderInfo(Icons.inventory_2, '${order['items']} items'),
                    _buildOrderInfo(Icons.calendar_today, order['date']),
                    _buildOrderInfo(Icons.payments, order['total']),
                  ],
                ),
                const Divider(height: 32),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    TextButton.icon(onPressed: () {}, icon: const Icon(Icons.visibility, size: 18), label: const Text('View')),
                    TextButton.icon(onPressed: () {}, icon: const Icon(Icons.call, size: 18), label: const Text('Call')),
                    if (order['status'] == 'Pending')
                      TextButton.icon(onPressed: () {}, icon: Icon(Icons.check, size: 18, color: AppColors.success), label: Text('Confirm', style: TextStyle(color: AppColors.success))),
                    if (order['status'] == 'Processing')
                      TextButton.icon(onPressed: () {}, icon: Icon(Icons.local_shipping, size: 18, color: AppColors.primary), label: Text('Ship', style: TextStyle(color: AppColors.primary))),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildOrderInfo(IconData icon, String text) {
    return Row(children: [
      Icon(icon, size: 16, color: AppColors.textSecondary),
      const SizedBox(width: 4),
      Text(text, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
    ]);
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case 'Pending': return AppColors.warning;
      case 'Processing': return AppColors.info;
      case 'Shipped': return AppColors.primary;
      case 'Delivered': return AppColors.success;
      case 'Cancelled': return AppColors.error;
      default: return AppColors.textSecondary;
    }
  }

  List<Map<String, dynamic>> _getFilteredOrders(String filter) {
    final allOrders = [
      {'id': 'ORD-7842', 'customer': 'Shoprite Mall', 'address': '45 Awolowo Rd, Ikeja', 'items': 12, 'total': '₦145,000', 'date': 'Today', 'status': 'Pending'},
      {'id': 'ORD-7841', 'customer': 'SPAR Nigeria', 'address': '12 Lekki Phase 1', 'items': 8, 'total': '₦89,500', 'date': 'Today', 'status': 'Processing'},
      {'id': 'ORD-7840', 'customer': 'Justrite Stores', 'address': '78 Herbert Macaulay', 'items': 25, 'total': '₦312,000', 'date': 'Yesterday', 'status': 'Shipped'},
      {'id': 'ORD-7839', 'customer': 'Mama Ngozi Store', 'address': '23 Ojuelegba Rd', 'items': 5, 'total': '₦45,000', 'date': 'Yesterday', 'status': 'Delivered'},
      {'id': 'ORD-7838', 'customer': 'Blessed Mart', 'address': '56 Yaba Market', 'items': 15, 'total': '₦178,000', 'date': '2 days ago', 'status': 'Delivered'},
    ];
    if (filter == 'all') return allOrders;
    return allOrders.where((o) => o['status'].toString().toLowerCase() == filter.toLowerCase()).toList();
  }
}
