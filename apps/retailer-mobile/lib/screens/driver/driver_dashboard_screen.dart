/// OmniRoute Ecosystem - Driver/Gig Worker Dashboard
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class DriverDashboardScreen extends ConsumerStatefulWidget {
  const DriverDashboardScreen({super.key});
  @override ConsumerState<DriverDashboardScreen> createState() => _DriverDashboardScreenState();
}

class _DriverDashboardScreenState extends ConsumerState<DriverDashboardScreen> {
  bool _isOnline = true;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      body: CustomScrollView(
        slivers: [
          _buildAppBar(),
          SliverPadding(
            padding: const EdgeInsets.all(16),
            sliver: SliverList(delegate: SliverChildListDelegate([
              _buildOnlineToggle(),
              const SizedBox(height: 20),
              _buildEarningsCard(),
              const SizedBox(height: 20),
              _buildActiveDelivery(),
              const SizedBox(height: 20),
              _buildTodaysStats(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Upcoming Deliveries'),
              const SizedBox(height: 12),
              _buildDeliveryQueue(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Performance'),
              const SizedBox(height: 12),
              _buildPerformanceMetrics(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Earnings History'),
              const SizedBox(height: 12),
              _buildEarningsHistory(),
            ])),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 120, floating: false, pinned: true, backgroundColor: AppColors.driverColor, foregroundColor: Colors.white,
      flexibleSpace: FlexibleSpaceBar(
        background: Container(
          decoration: BoxDecoration(gradient: LinearGradient(colors: [AppColors.driverColor, AppColors.driverColor.withValues(alpha: 0.85)])),
          child: SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Row(children: [
                const CircleAvatar(radius: 28, backgroundColor: Colors.white, child: Icon(Icons.person, color: AppColors.driverColor, size: 32)),
                const SizedBox(width: 16),
                Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, mainAxisAlignment: MainAxisAlignment.center, children: [
                  Text('Welcome back!', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
                  Text('Emeka Johnson', style: AppTypography.titleLarge.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
                  Row(children: [
                    Icon(Icons.star, color: Colors.amber, size: 16),
                    const SizedBox(width: 4),
                    Text('4.92', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w600)),
                    const SizedBox(width: 8),
                    Text('• 847 deliveries', style: TextStyle(color: Colors.white70, fontSize: 12)),
                  ]),
                ])),
              ]),
            ),
          ),
        ),
      ),
      actions: [
        IconButton(icon: const Icon(Icons.notifications_outlined), onPressed: () {}),
        IconButton(icon: const Icon(Icons.help_outline), onPressed: () {}),
      ],
    );
  }

  Widget _buildOnlineToggle() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      decoration: BoxDecoration(
        color: _isOnline ? AppColors.success.withValues(alpha: 0.1) : Colors.grey.withValues(alpha: 0.1),
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: _isOnline ? AppColors.success : Colors.grey),
      ),
      child: Row(children: [
        Icon(_isOnline ? Icons.wifi : Icons.wifi_off, color: _isOnline ? AppColors.success : Colors.grey, size: 28),
        const SizedBox(width: 16),
        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(_isOnline ? 'You\'re Online' : 'You\'re Offline', style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.w600)),
          Text(_isOnline ? 'Ready to receive deliveries' : 'Go online to receive deliveries', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ])),
        Switch(value: _isOnline, onChanged: (v) => setState(() => _isOnline = v), activeColor: AppColors.success),
      ]),
    );
  }

  Widget _buildEarningsCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(gradient: LinearGradient(colors: [AppColors.driverColor, AppColors.driverColor.withValues(alpha: 0.85)]), borderRadius: AppRadius.borderRadiusLg),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Text('Today\'s Earnings', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
            Container(padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4), decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusSm),
              child: Row(children: [
                const Icon(Icons.trending_up, color: Colors.white, size: 14),
                const SizedBox(width: 4),
                Text('+23%', style: TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w600)),
              ])),
          ]),
          const SizedBox(height: 12),
          Text('₦18,750', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 20),
          Row(children: [
            Expanded(child: _buildEarningBreakdown('Deliveries', '₦15,000', '12 trips')),
            Container(width: 1, height: 40, color: Colors.white24),
            Expanded(child: _buildEarningBreakdown('Tips', '₦2,250', '8 tips')),
            Container(width: 1, height: 40, color: Colors.white24),
            Expanded(child: _buildEarningBreakdown('Bonus', '₦1,500', 'Peak hours')),
          ]),
        ],
      ),
    );
  }

  Widget _buildEarningBreakdown(String label, String amount, String detail) {
    return Column(children: [
      Text(amount, style: AppTypography.titleMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
      Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
      Text(detail, style: AppTypography.labelSmall.copyWith(color: Colors.white54, fontSize: 10)),
    ]);
  }

  Widget _buildActiveDelivery() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.primary.withValues(alpha: 0.05),
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.primary, width: 2),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Container(padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4), decoration: BoxDecoration(color: AppColors.primary, borderRadius: AppRadius.borderRadiusSm),
              child: const Text('ACTIVE DELIVERY', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold))),
            Text('ORD-7842', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
          ]),
          const SizedBox(height: 16),
          Row(children: [
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                Container(width: 12, height: 12, decoration: BoxDecoration(color: AppColors.success, shape: BoxShape.circle)),
                const SizedBox(width: 8),
                Text('Pickup', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
              ]),
              Padding(padding: const EdgeInsets.only(left: 5), child: Container(width: 2, height: 20, color: AppColors.borderColor)),
              Row(children: [
                Container(width: 12, height: 12, decoration: BoxDecoration(color: AppColors.error, shape: BoxShape.circle)),
                const SizedBox(width: 8),
                Text('Dropoff', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
              ]),
            ])),
            Expanded(flex: 3, child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('Dangote Warehouse, Apapa', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
              const SizedBox(height: 16),
              Text('Shoprite Mall, Ikeja', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
            ])),
          ]),
          const SizedBox(height: 16),
          Row(children: [
            Expanded(child: _buildDeliveryInfo(Icons.inventory_2, '8 items')),
            Expanded(child: _buildDeliveryInfo(Icons.route, '12.5 km')),
            Expanded(child: _buildDeliveryInfo(Icons.payments, '₦1,800')),
            Expanded(child: _buildDeliveryInfo(Icons.timer, '25 min')),
          ]),
          const SizedBox(height: 16),
          Row(children: [
            Expanded(child: OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.phone), label: const Text('Call'))),
            const SizedBox(width: 12),
            Expanded(flex: 2, child: ElevatedButton.icon(
              onPressed: () {}, style: ElevatedButton.styleFrom(backgroundColor: AppColors.driverColor),
              icon: const Icon(Icons.navigation, color: Colors.white),
              label: const Text('Navigate', style: TextStyle(color: Colors.white)),
            )),
          ]),
        ],
      ),
    );
  }

  Widget _buildDeliveryInfo(IconData icon, String value) {
    return Column(children: [
      Icon(icon, color: AppColors.textSecondary, size: 20),
      const SizedBox(height: 4),
      Text(value, style: AppTypography.labelSmall.copyWith(fontWeight: FontWeight.w600)),
    ]);
  }

  Widget _buildTodaysStats() {
    return Row(children: [
      Expanded(child: _buildStatTile('Completed', '12', Icons.check_circle, AppColors.success)),
      const SizedBox(width: 12),
      Expanded(child: _buildStatTile('Distance', '67 km', Icons.route, AppColors.primary)),
      const SizedBox(width: 12),
      Expanded(child: _buildStatTile('Online', '6.5 hrs', Icons.timer, AppColors.info)),
    ]);
  }

  Widget _buildStatTile(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(children: [
        Icon(icon, color: color, size: 28),
        const SizedBox(height: 8),
        Text(value, style: AppTypography.titleLarge.copyWith(fontWeight: FontWeight.w700)),
        Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      ]),
    );
  }

  Widget _buildDeliveryQueue() {
    final deliveries = [
      {'from': 'SPAR Warehouse', 'to': 'Justrite Stores', 'items': '5 items', 'distance': '8.2 km', 'payout': '₦1,200', 'eta': '15 min'},
      {'from': 'Flour Mills', 'to': 'Mama Ngozi Store', 'items': '3 items', 'distance': '4.5 km', 'payout': '₦800', 'eta': '25 min'},
    ];
    return Column(children: deliveries.map((d) => Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Row(children: [
            Icon(Icons.arrow_upward, color: AppColors.success, size: 16),
            const SizedBox(width: 4),
            Text(d['from']!, style: AppTypography.labelMedium),
          ]),
          Text(d['payout']!, style: AppTypography.titleSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
        ]),
        Row(children: [
          Icon(Icons.arrow_downward, color: AppColors.error, size: 16),
          const SizedBox(width: 4),
          Text(d['to']!, style: AppTypography.labelMedium),
        ]),
        const SizedBox(height: 8),
        Row(children: [
          _buildTag(d['items']!),
          const SizedBox(width: 8),
          _buildTag(d['distance']!),
          const SizedBox(width: 8),
          _buildTag('ETA ${d['eta']}'),
        ]),
      ]),
    )).toList());
  }

  Widget _buildTag(String text) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(color: AppColors.borderColor, borderRadius: AppRadius.borderRadiusSm),
      child: Text(text, style: AppTypography.labelSmall),
    );
  }

  Widget _buildPerformanceMetrics() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(children: [
        _buildMetricRow('Acceptance Rate', '94%', 0.94, AppColors.success),
        const SizedBox(height: 16),
        _buildMetricRow('Completion Rate', '98%', 0.98, AppColors.success),
        const SizedBox(height: 16),
        _buildMetricRow('On-Time Delivery', '91%', 0.91, AppColors.primary),
        const SizedBox(height: 16),
        _buildMetricRow('Customer Rating', '4.92', 0.984, AppColors.warning),
      ]),
    );
  }

  Widget _buildMetricRow(String label, String value, double progress, Color color) {
    return Row(children: [
      Expanded(flex: 2, child: Text(label, style: AppTypography.labelMedium)),
      Expanded(flex: 3, child: LinearProgressIndicator(value: progress, backgroundColor: color.withValues(alpha: 0.2), valueColor: AlwaysStoppedAnimation(color))),
      const SizedBox(width: 12),
      Text(value, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
    ]);
  }

  Widget _buildEarningsHistory() {
    final history = [
      {'day': 'Today', 'amount': '₦18,750', 'trips': '12 trips'},
      {'day': 'Yesterday', 'amount': '₦22,400', 'trips': '15 trips'},
      {'day': 'Mon, Jan 16', 'amount': '₦19,800', 'trips': '13 trips'},
    ];
    return Column(children: history.map((h) => ListTile(
      contentPadding: EdgeInsets.zero,
      title: Text(h['day']!, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
      subtitle: Text(h['trips']!, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
      trailing: Text(h['amount']!, style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.w600, color: AppColors.success)),
    )).toList());
  }

  Widget _buildBottomNav() {
    return BottomNavigationBar(
      currentIndex: 0, type: BottomNavigationBarType.fixed, selectedItemColor: AppColors.driverColor,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.dashboard), label: 'Home'),
        BottomNavigationBarItem(icon: Icon(Icons.local_shipping), label: 'Deliveries'),
        BottomNavigationBarItem(icon: Icon(Icons.account_balance_wallet), label: 'Earnings'),
        BottomNavigationBarItem(icon: Icon(Icons.bar_chart), label: 'Stats'),
        BottomNavigationBarItem(icon: Icon(Icons.account_circle), label: 'Profile'),
      ],
    );
  }
}
