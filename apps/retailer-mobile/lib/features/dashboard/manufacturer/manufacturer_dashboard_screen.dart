/// OmniRoute Ecosystem - Manufacturer Dashboard
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class ManufacturerDashboardScreen extends ConsumerWidget {
  const ManufacturerDashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(title: const Text('Manufacturing Hub'), backgroundColor: AppColors.manufacturerColor, foregroundColor: Colors.white),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          children: [
            _buildRevenueCard(),
            const SizedBox(height: 20),
            GridView.count(
              shrinkWrap: true, physics: const NeverScrollableScrollPhysics(),
              crossAxisCount: 2, crossAxisSpacing: 12, mainAxisSpacing: 12, childAspectRatio: 1.5,
              children: [
                StatCard(title: 'Products', value: '48', icon: Icons.category, iconColor: AppColors.manufacturerColor),
                StatCard(title: 'Active Orders', value: '156', icon: Icons.receipt_long, iconColor: AppColors.primary),
                StatCard(title: 'Distributors', value: '34', icon: Icons.people, iconColor: AppColors.info),
                StatCard(title: 'Production Rate', value: '94%', icon: Icons.speed, iconColor: AppColors.success),
              ],
            ),
            const SizedBox(height: 20),
            const SectionHeader(title: 'Quick Actions'),
            const SizedBox(height: 12),
            QuickActionGrid(actions: [
              QuickActionItem(icon: Icons.add_box, label: 'Add Product', color: AppColors.success, onTap: () {}),
              QuickActionItem(icon: Icons.local_shipping, label: 'Ship Order', color: AppColors.manufacturerColor, onTap: () {}),
              QuickActionItem(icon: Icons.analytics, label: 'Analytics', color: AppColors.primary, onTap: () {}),
              QuickActionItem(icon: Icons.inventory, label: 'Inventory', color: AppColors.info, onTap: () {}),
            ]),
          ],
        ),
      ),
    );
  }

  Widget _buildRevenueCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(gradient: LinearGradient(colors: [AppColors.manufacturerColor, AppColors.manufacturerColor.withValues(alpha: 0.85)]), borderRadius: AppRadius.borderRadiusLg),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Monthly Revenue', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
          const SizedBox(height: 8),
          Text('₦45,890,000', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 12),
          Row(children: [
            const Icon(Icons.trending_up, color: Colors.white70, size: 16),
            const SizedBox(width: 4),
            Text('+18.5% vs last month', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
          ]),
        ],
      ),
    );
  }
}
