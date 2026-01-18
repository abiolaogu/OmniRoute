/// OmniRoute Ecosystem - Logistics Dashboard Screen
/// Dashboard for logistics companies with fleet, deliveries, and route management

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class LogisticsDashboardScreen extends ConsumerStatefulWidget {
  const LogisticsDashboardScreen({super.key});

  @override
  ConsumerState<LogisticsDashboardScreen> createState() => _LogisticsDashboardScreenState();
}

class _LogisticsDashboardScreenState extends ConsumerState<LogisticsDashboardScreen> {
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
                _buildEarningsCard(),
                const SizedBox(height: 20),
                _buildStatsGrid(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildActiveDeliveries(),
                const SizedBox(height: 20),
                _buildFleetStatus(),
                const SizedBox(height: 20),
                _buildDriverPerformance(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 120,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.logisticsColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Text('Logistics Hub', style: AppTypography.titleLarge.copyWith(color: Colors.white, fontSize: 18)),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.logisticsColor, AppColors.logisticsColor.withValues(alpha: 0.85)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Padding(
              padding: const EdgeInsets.only(right: 20),
              child: Icon(Icons.local_shipping, size: 80, color: Colors.white.withValues(alpha: 0.2)),
            ),
          ),
        ),
      ),
      actions: [
        IconButton(icon: const Icon(Icons.map, color: Colors.white), onPressed: () {}),
        IconButton(icon: const Icon(Icons.notifications_outlined, color: Colors.white), onPressed: () {}),
      ],
    );
  }

  Widget _buildEarningsCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.logisticsColor, AppColors.logisticsColor.withValues(alpha: 0.85)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: AppRadius.borderRadiusLg,
        boxShadow: AppShadows.lg,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text("Today's Earnings", style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.2),
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Row(
                  children: [
                    const Icon(Icons.trending_up, color: Colors.white, size: 14),
                    const SizedBox(width: 4),
                    Text('+23%', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text('₦487,500', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 20),
          Row(
            children: [
              _EarningMetric(label: 'Deliveries', value: '47'),
              const SizedBox(width: 32),
              _EarningMetric(label: 'Avg. Distance', value: '12km'),
              const SizedBox(width: 32),
              _EarningMetric(label: 'Active Drivers', value: '23'),
            ],
          ),
        ],
      ),
    ).animate().fadeIn().slideY(begin: 0.2, end: 0);
  }

  Widget _buildStatsGrid() {
    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2,
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.5,
      children: [
        StatCard(
          title: 'Pending Deliveries',
          value: '34',
          icon: Icons.pending_actions,
          iconColor: AppColors.warning,
        ),
        StatCard(
          title: 'In Transit',
          value: '18',
          icon: Icons.local_shipping,
          iconColor: AppColors.info,
        ),
        StatCard(
          title: 'Completed Today',
          value: '47',
          icon: Icons.check_circle,
          iconColor: AppColors.success,
          growthPercentage: 15.2,
        ),
        StatCard(
          title: 'Fleet Utilization',
          value: '78%',
          icon: Icons.directions_car,
          iconColor: AppColors.logisticsColor,
        ),
      ],
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
            QuickActionItem(icon: Icons.add_location, label: 'New Delivery', color: AppColors.success, onTap: () {}),
            QuickActionItem(icon: Icons.route, label: 'Optimize Routes', color: AppColors.primary, onTap: () {}),
            QuickActionItem(icon: Icons.person_add, label: 'Add Driver', color: AppColors.logisticsColor, onTap: () {}),
            QuickActionItem(icon: Icons.map, label: 'Live Map', color: AppColors.info, onTap: () {}, badge: '18'),
          ],
        ),
      ],
    );
  }

  Widget _buildActiveDeliveries() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Active Deliveries', actionText: 'View Map', onAction: () {}),
        const SizedBox(height: 12),
        ...List.generate(3, (i) => Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: _DeliveryCard(
            orderNumber: 'DEL-${5001 + i}',
            origin: ['Ikeja Warehouse', 'Apapa Port', 'Lekki Hub'][i],
            destination: ['Victoria Island', 'Surulere', 'Yaba Market'][i],
            driver: ['John Okonkwo', 'Ahmed Musa', 'Chidi Eze'][i],
            status: i == 0 ? 'in_transit' : i == 1 ? 'picking_up' : 'pending',
            eta: '${15 + i * 10} mins',
            progress: [0.7, 0.3, 0.0][i],
          ),
        )),
      ],
    );
  }

  Widget _buildFleetStatus() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Fleet Status', actionText: 'Manage', onAction: () {}),
        const SizedBox(height: 12),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: Column(
            children: [
              Row(
                children: [
                  _FleetStatusItem(label: 'Available', count: 15, color: AppColors.success),
                  _FleetStatusItem(label: 'On Trip', count: 23, color: AppColors.info),
                  _FleetStatusItem(label: 'Maintenance', count: 4, color: AppColors.warning),
                  _FleetStatusItem(label: 'Offline', count: 3, color: AppColors.error),
                ],
              ),
              const SizedBox(height: 16),
              ClipRRect(
                borderRadius: AppRadius.borderRadiusFull,
                child: Row(
                  children: [
                    Expanded(flex: 15, child: Container(height: 8, color: AppColors.success)),
                    Expanded(flex: 23, child: Container(height: 8, color: AppColors.info)),
                    Expanded(flex: 4, child: Container(height: 8, color: AppColors.warning)),
                    Expanded(flex: 3, child: Container(height: 8, color: AppColors.error)),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildDriverPerformance() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Top Drivers Today', actionText: 'View All', onAction: () {}),
        const SizedBox(height: 12),
        ...List.generate(3, (i) => Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: Row(
            children: [
              Stack(
                children: [
                  CircleAvatar(
                    radius: 24,
                    backgroundColor: [AppColors.success, AppColors.info, AppColors.warning][i].withValues(alpha: 0.2),
                    child: Text('${i + 1}', style: AppTypography.titleMedium.copyWith(color: [AppColors.success, AppColors.info, AppColors.warning][i])),
                  ),
                  if (i == 0)
                    Positioned(
                      right: 0, bottom: 0,
                      child: Container(
                        padding: const EdgeInsets.all(2),
                        decoration: const BoxDecoration(color: AppColors.warning, shape: BoxShape.circle),
                        child: const Icon(Icons.star, color: Colors.white, size: 12),
                      ),
                    ),
                ],
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(['John Okonkwo', 'Ahmed Musa', 'Chidi Eze'][i], style: AppTypography.titleSmall),
                    Text('${[12, 10, 9][i]} deliveries • ${[4.9, 4.8, 4.7][i]}★', style: AppTypography.bodySmall),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text('₦${[45000, 38000, 35000][i]}', style: AppTypography.titleSmall.copyWith(color: AppColors.success)),
                  Text('Earned today', style: AppTypography.labelSmall.copyWith(color: AppColors.grey500)),
                ],
              ),
            ],
          ),
        ).animate(delay: Duration(milliseconds: 100 * i)).fadeIn().slideX(begin: 0.1, end: 0)),
      ],
    );
  }
}

class _EarningMetric extends StatelessWidget {
  final String label;
  final String value;

  const _EarningMetric({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(value, style: AppTypography.titleMedium.copyWith(color: Colors.white)),
        Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white60)),
      ],
    );
  }
}

class _DeliveryCard extends StatelessWidget {
  final String orderNumber;
  final String origin;
  final String destination;
  final String driver;
  final String status;
  final String eta;
  final double progress;

  const _DeliveryCard({
    required this.orderNumber,
    required this.origin,
    required this.destination,
    required this.driver,
    required this.status,
    required this.eta,
    required this.progress,
  });

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
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(orderNumber, style: AppTypography.titleSmall),
              StatusChip(status: status),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _LocationRow(icon: Icons.circle, color: AppColors.success, label: origin),
                    Container(
                      margin: const EdgeInsets.only(left: 4),
                      height: 20,
                      width: 2,
                      color: AppColors.grey300,
                    ),
                    _LocationRow(icon: Icons.location_on, color: AppColors.error, label: destination),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: AppColors.grey50,
                  borderRadius: AppRadius.borderRadiusSm,
                ),
                child: Column(
                  children: [
                    const Icon(Icons.access_time, color: AppColors.grey600, size: 20),
                    const SizedBox(height: 4),
                    Text(eta, style: AppTypography.labelSmall.copyWith(fontWeight: FontWeight.w600)),
                    Text('ETA', style: AppTypography.labelSmall.copyWith(color: AppColors.grey500, fontSize: 10)),
                  ],
                ),
              ),
            ],
          ),
          if (progress > 0) ...[
            const SizedBox(height: 12),
            ClipRRect(
              borderRadius: AppRadius.borderRadiusFull,
              child: LinearProgressIndicator(
                value: progress,
                backgroundColor: AppColors.grey200,
                valueColor: const AlwaysStoppedAnimation(AppColors.logisticsColor),
                minHeight: 6,
              ),
            ),
          ],
          const SizedBox(height: 12),
          Row(
            children: [
              CircleAvatar(
                radius: 14,
                backgroundColor: AppColors.grey200,
                child: const Icon(Icons.person, size: 16, color: AppColors.grey600),
              ),
              const SizedBox(width: 8),
              Text(driver, style: AppTypography.bodySmall),
              const Spacer(),
              IconButton(
                icon: const Icon(Icons.phone, color: AppColors.primary, size: 20),
                onPressed: () {},
                constraints: const BoxConstraints(),
                padding: EdgeInsets.zero,
              ),
              const SizedBox(width: 16),
              IconButton(
                icon: const Icon(Icons.message, color: AppColors.info, size: 20),
                onPressed: () {},
                constraints: const BoxConstraints(),
                padding: EdgeInsets.zero,
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _LocationRow extends StatelessWidget {
  final IconData icon;
  final Color color;
  final String label;

  const _LocationRow({required this.icon, required this.color, required this.label});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, size: 10, color: color),
        const SizedBox(width: 8),
        Expanded(child: Text(label, style: AppTypography.bodySmall, maxLines: 1, overflow: TextOverflow.ellipsis)),
      ],
    );
  }
}

class _FleetStatusItem extends StatelessWidget {
  final String label;
  final int count;
  final Color color;

  const _FleetStatusItem({required this.label, required this.count, required this.color});

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Column(
        children: [
          Container(
            width: 12, height: 12,
            decoration: BoxDecoration(color: color, shape: BoxShape.circle),
          ),
          const SizedBox(height: 4),
          Text('$count', style: AppTypography.titleSmall),
          Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
        ],
      ),
    );
  }
}
