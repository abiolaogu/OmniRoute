/// OmniRoute Ecosystem - Sales Agent Dashboard
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class AgentDashboardScreen extends ConsumerStatefulWidget {
  const AgentDashboardScreen({super.key});
  @override ConsumerState<AgentDashboardScreen> createState() => _AgentDashboardScreenState();
}

class _AgentDashboardScreenState extends ConsumerState<AgentDashboardScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('Sales Agent'), backgroundColor: AppColors.agentColor, foregroundColor: Colors.white,
        actions: [
          IconButton(icon: const Icon(Icons.leaderboard), onPressed: () {}),
          IconButton(icon: const Icon(Icons.notifications_outlined), onPressed: () {}),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildEarningsCard(),
            const SizedBox(height: 20),
            _buildPerformanceStats(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Today\'s Tasks'),
            const SizedBox(height: 12),
            _buildTodaysTasks(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Route Planner'),
            const SizedBox(height: 12),
            _buildRoutePlanner(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Commission Breakdown'),
            const SizedBox(height: 12),
            _buildCommissionBreakdown(),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Recent Activities'),
            const SizedBox(height: 12),
            _buildRecentActivities(),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {}, backgroundColor: AppColors.agentColor,
        icon: const Icon(Icons.add_business, color: Colors.white),
        label: const Text('Register Retailer', style: TextStyle(color: Colors.white)),
      ),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildEarningsCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.agentColor, AppColors.agentColor.withValues(alpha: 0.85)]),
        borderRadius: AppRadius.borderRadiusLg,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('This Month\'s Earnings', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              const SizedBox(height: 8),
              Text('₦287,500', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
            ]),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusMd),
              child: Column(children: [
                const Icon(Icons.star, color: Colors.amber, size: 24),
                const SizedBox(height: 4),
                Text('Gold', style: AppTypography.labelSmall.copyWith(color: Colors.white, fontWeight: FontWeight.bold)),
              ]),
            ),
          ]),
          const SizedBox(height: 20),
          Row(children: [
            Expanded(child: _buildEarningItem('Base Salary', '₦150,000')),
            Expanded(child: _buildEarningItem('Commission', '₦112,500')),
            Expanded(child: _buildEarningItem('Bonus', '₦25,000')),
          ]),
        ],
      ),
    );
  }

  Widget _buildEarningItem(String label, String value) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white60)),
      const SizedBox(height: 4),
      Text(value, style: AppTypography.titleMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
    ]);
  }

  Widget _buildPerformanceStats() {
    return Row(children: [
      Expanded(child: _buildStatCard('Visits Today', '12/15', Icons.store, AppColors.primary, 0.8)),
      const SizedBox(width: 12),
      Expanded(child: _buildStatCard('Orders Taken', '8', Icons.shopping_cart, AppColors.success, null)),
      const SizedBox(width: 12),
      Expanded(child: _buildStatCard('New Retailers', '2', Icons.person_add, AppColors.info, null)),
    ]);
  }

  Widget _buildStatCard(String label, String value, IconData icon, Color color, double? progress) {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Icon(icon, color: color, size: 24),
        const SizedBox(height: 8),
        Text(value, style: AppTypography.titleLarge.copyWith(fontWeight: FontWeight.w700)),
        Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        if (progress != null) ...[
          const SizedBox(height: 8),
          LinearProgressIndicator(value: progress, backgroundColor: color.withValues(alpha: 0.2), valueColor: AlwaysStoppedAnimation(color)),
        ],
      ]),
    );
  }

  Widget _buildTodaysTasks() {
    final tasks = [
      {'retailer': 'Mama Ngozi Store', 'address': '45 Awolowo Rd, Ikeja', 'type': 'Restock Visit', 'time': '10:00 AM', 'status': 'completed'},
      {'retailer': 'Iya Bose Shop', 'address': '12 Market St, Yaba', 'type': 'New Registration', 'time': '11:30 AM', 'status': 'current'},
      {'retailer': 'Chukwu Provisions', 'address': '78 Herbert Macaulay', 'type': 'Order Follow-up', 'time': '1:00 PM', 'status': 'pending'},
      {'retailer': 'Blessed Mart', 'address': '23 Ojuelegba Rd', 'type': 'Payment Collection', 'time': '2:30 PM', 'status': 'pending'},
    ];
    return Column(children: tasks.map((t) {
      final isCompleted = t['status'] == 'completed';
      final isCurrent = t['status'] == 'current';
      return Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: isCurrent ? AppColors.agentColor.withValues(alpha: 0.05) : Colors.white,
          borderRadius: AppRadius.borderRadiusMd,
          border: isCurrent ? Border.all(color: AppColors.agentColor, width: 2) : null,
          boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)],
        ),
        child: Row(children: [
          Container(
            width: 48, height: 48,
            decoration: BoxDecoration(
              color: isCompleted ? AppColors.success : isCurrent ? AppColors.agentColor : AppColors.borderColor,
              shape: BoxShape.circle,
            ),
            child: Icon(isCompleted ? Icons.check : isCurrent ? Icons.directions_walk : Icons.schedule, color: Colors.white, size: 24),
          ),
          const SizedBox(width: 12),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text(t['retailer']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600, decoration: isCompleted ? TextDecoration.lineThrough : null)),
              Text(t['time']!, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            ]),
            const SizedBox(height: 4),
            Text(t['address']!, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            const SizedBox(height: 4),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(color: _getTaskTypeColor(t['type']!).withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
              child: Text(t['type']!, style: AppTypography.labelSmall.copyWith(color: _getTaskTypeColor(t['type']!))),
            ),
          ])),
          if (isCurrent) IconButton(icon: const Icon(Icons.navigation, color: AppColors.agentColor), onPressed: () {}),
        ]),
      );
    }).toList());
  }

  Color _getTaskTypeColor(String type) {
    switch (type) {
      case 'Restock Visit': return AppColors.primary;
      case 'New Registration': return AppColors.success;
      case 'Order Follow-up': return AppColors.info;
      case 'Payment Collection': return AppColors.warning;
      default: return AppColors.textSecondary;
    }
  }

  Widget _buildRoutePlanner() {
    return Container(
      height: 200,
      decoration: BoxDecoration(color: AppColors.borderColor, borderRadius: AppRadius.borderRadiusMd),
      child: Stack(children: [
        Center(child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
          Icon(Icons.map, size: 48, color: AppColors.textSecondary.withValues(alpha: 0.5)),
          const SizedBox(height: 8),
          Text('Route Map', style: AppTypography.titleMedium.copyWith(color: AppColors.textSecondary)),
          const SizedBox(height: 4),
          Text('15 stops • 23km total', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ])),
        Positioned(bottom: 16, right: 16, child: ElevatedButton.icon(
          onPressed: () {}, style: ElevatedButton.styleFrom(backgroundColor: AppColors.agentColor),
          icon: const Icon(Icons.navigation, color: Colors.white, size: 18),
          label: const Text('Start Navigation', style: TextStyle(color: Colors.white)),
        )),
      ]),
    );
  }

  Widget _buildCommissionBreakdown() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(children: [
        _buildCommissionRow('Order Commission (2%)', '₦67,500', '₦3.375M orders'),
        const Divider(height: 24),
        _buildCommissionRow('New Retailer Bonus', '₦15,000', '3 new retailers'),
        const Divider(height: 24),
        _buildCommissionRow('Collection Bonus', '₦20,000', '₦2M collected'),
        const Divider(height: 24),
        _buildCommissionRow('Performance Bonus', '₦10,000', 'Top 10% agents'),
        const SizedBox(height: 16),
        Container(
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(color: AppColors.success.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
          child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Text('Total Commission', style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.w600)),
            Text('₦112,500', style: AppTypography.titleLarge.copyWith(color: AppColors.success, fontWeight: FontWeight.w700)),
          ]),
        ),
      ]),
    );
  }

  Widget _buildCommissionRow(String label, String amount, String detail) {
    return Row(children: [
      Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(label, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
        Text(detail, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      ])),
      Text(amount, style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.w600)),
    ]);
  }

  Widget _buildRecentActivities() {
    final activities = [
      {'action': 'Completed visit', 'retailer': 'Mama Ngozi Store', 'time': '10:45 AM', 'icon': Icons.check_circle},
      {'action': 'Order placed', 'retailer': '₦45,000 - Iya Bose Shop', 'time': '11:20 AM', 'icon': Icons.shopping_cart},
      {'action': 'Payment collected', 'retailer': '₦120,000 - Chukwu Provisions', 'time': '12:30 PM', 'icon': Icons.payments},
    ];
    return Column(children: activities.map((a) => ListTile(
      contentPadding: EdgeInsets.zero,
      leading: CircleAvatar(backgroundColor: AppColors.agentColor.withValues(alpha: 0.1), child: Icon(a['icon'] as IconData, color: AppColors.agentColor, size: 20)),
      title: Text(a['action']!, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
      subtitle: Text(a['retailer']!, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      trailing: Text(a['time']!, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
    )).toList());
  }

  Widget _buildBottomNav() {
    return BottomNavigationBar(
      currentIndex: 0, type: BottomNavigationBarType.fixed, selectedItemColor: AppColors.agentColor,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.dashboard), label: 'Dashboard'),
        BottomNavigationBarItem(icon: Icon(Icons.route), label: 'Routes'),
        BottomNavigationBarItem(icon: Icon(Icons.store), label: 'Retailers'),
        BottomNavigationBarItem(icon: Icon(Icons.leaderboard), label: 'Rankings'),
        BottomNavigationBarItem(icon: Icon(Icons.account_circle), label: 'Profile'),
      ],
    );
  }
}
