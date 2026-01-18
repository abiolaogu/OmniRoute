/// OmniRoute Ecosystem - Warehouse Dashboard
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class WarehouseDashboardScreen extends ConsumerWidget {
  const WarehouseDashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('Warehouse Operations'),
        backgroundColor: AppColors.warehouseColor,
        foregroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildCapacityCard(),
            const SizedBox(height: 20),
            GridView.count(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 2,
              crossAxisSpacing: 12,
              mainAxisSpacing: 12,
              childAspectRatio: 1.5,
              children: [
                StatCard(title: 'Total SKUs', value: '2,456', icon: Icons.inventory_2, iconColor: AppColors.warehouseColor),
                StatCard(title: 'Inbound Today', value: '45', icon: Icons.move_to_inbox, iconColor: AppColors.success),
                StatCard(title: 'Outbound Today', value: '78', icon: Icons.outbox, iconColor: AppColors.info),
                StatCard(title: 'Pick Accuracy', value: '99.2%', icon: Icons.check_circle, iconColor: AppColors.success),
              ],
            ),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Quick Actions'),
            const SizedBox(height: 12),
            QuickActionGrid(
              actions: [
                QuickActionItem(icon: Icons.qr_code_scanner, label: 'Scan Item', color: AppColors.primary, onTap: () {}),
                QuickActionItem(icon: Icons.inventory, label: 'Stock Count', color: AppColors.warehouseColor, onTap: () {}),
                QuickActionItem(icon: Icons.local_shipping, label: 'New Shipment', color: AppColors.success, onTap: () {}),
                QuickActionItem(icon: Icons.assignment, label: 'Pick List', color: AppColors.info, onTap: () {}, badge: '12'),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildCapacityCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.warehouseColor, AppColors.warehouseColor.withValues(alpha: 0.85)],
        ),
        borderRadius: AppRadius.borderRadiusLg,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Warehouse Capacity', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
          const SizedBox(height: 8),
          Row(
            children: [
              Text('78%', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    ClipRRect(
                      borderRadius: AppRadius.borderRadiusFull,
                      child: LinearProgressIndicator(value: 0.78, backgroundColor: Colors.white24, valueColor: const AlwaysStoppedAnimation(Colors.white), minHeight: 8),
                    ),
                    const SizedBox(height: 4),
                    Text('15,600 / 20,000 sqft', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
