/// OmniRoute Ecosystem - Supermall Dashboard Screen
/// Comprehensive dashboard for supermall and shopping center operators
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class SupermallDashboardScreen extends ConsumerStatefulWidget {
  const SupermallDashboardScreen({super.key});
  @override
  ConsumerState<SupermallDashboardScreen> createState() => _SupermallDashboardScreenState();
}

class _SupermallDashboardScreenState extends ConsumerState<SupermallDashboardScreen> {
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
                _buildRevenueCard(),
                const SizedBox(height: 20),
                _buildOccupancyStats(),
                const SizedBox(height: 20),
                _buildFootTrafficChart(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildTenantPerformance(),
                const SizedBox(height: 20),
                _buildParkingStatus(),
                const SizedBox(height: 20),
                _buildUtilityConsumption(),
                const SizedBox(height: 20),
                _buildUpcomingEvents(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.supermallColor,
        icon: const Icon(Icons.add_business),
        label: const Text('Add Tenant'),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 140,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.supermallColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('SUPERMALL MANAGEMENT', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            Text('Victoria Island Mall', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.supermallColor, AppColors.supermallColor.withValues(alpha: 0.8)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Padding(
              padding: const EdgeInsets.only(right: 20),
              child: Icon(Icons.location_city, size: 80, color: Colors.white.withValues(alpha: 0.2)),
            ),
          ),
        ),
      ),
      actions: [
        IconButton(icon: const Icon(Icons.notifications_outlined, color: Colors.white), onPressed: () {}),
        IconButton(icon: const Icon(Icons.settings, color: Colors.white), onPressed: () {}),
      ],
    );
  }

  Widget _buildRevenueCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.supermallColor, AppColors.supermallColor.withValues(alpha: 0.8)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Total Revenue', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusSm),
                child: Text('This Month', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text('₦156,780,000', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildRevenueMetric('Tenant Rent', '₦98.5M', Icons.store),
              const SizedBox(width: 16),
              _buildRevenueMetric('Service Fees', '₦32.2M', Icons.miscellaneous_services),
              const SizedBox(width: 16),
              _buildRevenueMetric('Parking', '₦26.1M', Icons.local_parking),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildRevenueMetric(String label, String value, IconData icon) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: AppRadius.borderRadiusSm),
        child: Column(
          children: [
            Icon(icon, color: Colors.white, size: 20),
            const SizedBox(height: 4),
            Text(value, style: AppTypography.titleSmall.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
            Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white70), textAlign: TextAlign.center),
          ],
        ),
      ),
    );
  }

  Widget _buildOccupancyStats() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Mall Occupancy'),
        const SizedBox(height: 12),
        GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: 2,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          childAspectRatio: 1.6,
          children: [
            StatCard(title: 'Occupancy Rate', value: '94.2%', icon: Icons.business, iconColor: AppColors.success),
            StatCard(title: 'Active Tenants', value: '127', icon: Icons.storefront, iconColor: AppColors.primary),
            StatCard(title: 'Available Units', value: '8', icon: Icons.door_front_door, iconColor: AppColors.warning),
            StatCard(title: 'Pending Leases', value: '5', icon: Icons.pending_actions, iconColor: AppColors.info),
          ],
        ),
      ],
    );
  }

  Widget _buildFootTrafficChart() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Foot Traffic', style: AppTypography.titleMedium),
              Row(children: [
                const Icon(Icons.people, color: AppColors.success, size: 16),
                const SizedBox(width: 4),
                Text('42,350 today', style: AppTypography.labelMedium.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
              ]),
            ],
          ),
          const SizedBox(height: 16),
          SizedBox(
            height: 120,
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: List.generate(7, (i) {
                final heights = [0.4, 0.5, 0.6, 0.8, 0.9, 1.0, 0.7];
                final days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
                return Column(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    Container(
                      width: 32,
                      height: 100 * heights[i],
                      decoration: BoxDecoration(
                        color: i == 5 ? AppColors.supermallColor : AppColors.supermallColor.withValues(alpha: 0.3),
                        borderRadius: const BorderRadius.vertical(top: Radius.circular(4)),
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(days[i], style: AppTypography.labelSmall),
                  ],
                );
              }),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActions() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Quick Actions'),
        const SizedBox(height: 12),
        QuickActionGrid(
          actions: [
            QuickActionItem(icon: Icons.add_business, label: 'New Tenant', color: AppColors.success, onTap: () {}),
            QuickActionItem(icon: Icons.receipt_long, label: 'Billing', color: AppColors.primary, onTap: () {}),
            QuickActionItem(icon: Icons.local_parking, label: 'Parking', color: AppColors.warning, onTap: () {}),
            QuickActionItem(icon: Icons.event, label: 'Events', color: AppColors.info, onTap: () {}),
            QuickActionItem(icon: Icons.security, label: 'Security', color: AppColors.error, onTap: () {}),
            QuickActionItem(icon: Icons.analytics, label: 'Analytics', color: AppColors.supermallColor, onTap: () {}),
          ],
        ),
      ],
    );
  }

  Widget _buildTenantPerformance() {
    final tenants = [
      {'name': 'Shoprite', 'category': 'Grocery', 'revenue': '₦45.2M', 'trend': 12.5, 'rank': 1},
      {'name': 'Samsung Store', 'category': 'Electronics', 'revenue': '₦28.7M', 'trend': 8.3, 'rank': 2},
      {'name': 'Nike', 'category': 'Fashion', 'revenue': '₦18.9M', 'trend': -2.1, 'rank': 3},
      {'name': 'Filmhouse Cinemas', 'category': 'Entertainment', 'revenue': '₦15.4M', 'trend': 15.8, 'rank': 4},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Top Performing Tenants', onViewAll: () {}),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: tenants.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, i) {
              final t = tenants[i];
              final trend = t['trend'] as double;
              return ListTile(
                leading: CircleAvatar(backgroundColor: AppColors.supermallColor.withValues(alpha: 0.1), child: Text('${t['rank']}', style: TextStyle(color: AppColors.supermallColor, fontWeight: FontWeight.bold))),
                title: Text(t['name'] as String, style: AppTypography.titleSmall),
                subtitle: Text(t['category'] as String, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                trailing: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(t['revenue'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
                    Row(mainAxisSize: MainAxisSize.min, children: [
                      Icon(trend >= 0 ? Icons.trending_up : Icons.trending_down, size: 12, color: trend >= 0 ? AppColors.success : AppColors.error),
                      Text('${trend.abs()}%', style: AppTypography.labelSmall.copyWith(color: trend >= 0 ? AppColors.success : AppColors.error)),
                    ]),
                  ],
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildParkingStatus() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Parking Status', style: AppTypography.titleMedium),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: AppColors.success.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
                child: Text('654 Available', style: AppTypography.labelSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildParkingZone('Basement 1', 120, 95),
              const SizedBox(width: 12),
              _buildParkingZone('Basement 2', 150, 142),
              const SizedBox(width: 12),
              _buildParkingZone('Open Air', 200, 145),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildParkingZone(String zone, int total, int occupied) {
    final percentage = (occupied / total * 100).toInt();
    final color = percentage > 90 ? AppColors.error : percentage > 70 ? AppColors.warning : AppColors.success;
    return Expanded(
      child: Column(
        children: [
          Text(zone, style: AppTypography.labelSmall),
          const SizedBox(height: 8),
          Stack(
            alignment: Alignment.center,
            children: [
              SizedBox(
                width: 60, height: 60,
                child: CircularProgressIndicator(
                  value: occupied / total,
                  backgroundColor: color.withValues(alpha: 0.2),
                  valueColor: AlwaysStoppedAnimation(color),
                  strokeWidth: 6,
                ),
              ),
              Text('$percentage%', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
            ],
          ),
          const SizedBox(height: 4),
          Text('$occupied/$total', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ],
      ),
    );
  }

  Widget _buildUtilityConsumption() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Utility Consumption', style: AppTypography.titleMedium),
          const SizedBox(height: 16),
          _buildUtilityRow(Icons.electric_bolt, 'Electricity', '245,780 kWh', '₦18.4M', AppColors.warning),
          const Divider(),
          _buildUtilityRow(Icons.water_drop, 'Water', '12,450 m³', '₦2.1M', AppColors.info),
          const Divider(),
          _buildUtilityRow(Icons.ac_unit, 'HVAC', '₦8.5M', 'Monthly', AppColors.primary),
        ],
      ),
    );
  }

  Widget _buildUtilityRow(IconData icon, String name, String usage, String cost, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(8),
            decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
            child: Icon(icon, color: color, size: 20),
          ),
          const SizedBox(width: 12),
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(name, style: AppTypography.titleSmall),
              Text(usage, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            ],
          )),
          Text(cost, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }

  Widget _buildUpcomingEvents() {
    final events = [
      {'name': 'Black Friday Sale', 'date': 'Nov 29', 'tenants': 45},
      {'name': 'Christmas Carnival', 'date': 'Dec 15-25', 'tenants': 78},
      {'name': 'New Year Concert', 'date': 'Dec 31', 'tenants': 12},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Upcoming Events', onViewAll: () {}),
        const SizedBox(height: 12),
        ...events.map((e) => Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.supermallColor.withValues(alpha: 0.3)),
          ),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(color: AppColors.supermallColor.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
                child: const Icon(Icons.event, color: AppColors.supermallColor),
              ),
              const SizedBox(width: 12),
              Expanded(child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(e['name'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
                  Text('${e['tenants']} participating tenants', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                ],
              )),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(color: AppColors.supermallColor, borderRadius: AppRadius.borderRadiusSm),
                child: Text(e['date'] as String, style: const TextStyle(color: Colors.white, fontSize: 12, fontWeight: FontWeight.w600)),
              ),
            ],
          ),
        )),
      ],
    );
  }
}
