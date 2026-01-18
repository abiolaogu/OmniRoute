/// OmniRoute Ecosystem - B2B Manufacturer Dashboard Screen
/// Comprehensive dashboard for B2B/industrial manufacturers
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class B2BManufacturerDashboardScreen extends ConsumerStatefulWidget {
  const B2BManufacturerDashboardScreen({super.key});
  @override
  ConsumerState<B2BManufacturerDashboardScreen> createState() => _B2BManufacturerDashboardScreenState();
}

class _B2BManufacturerDashboardScreenState extends ConsumerState<B2BManufacturerDashboardScreen> {
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
                _buildContractOrdersCard(),
                const SizedBox(height: 20),
                _buildProductionPipeline(),
                const SizedBox(height: 20),
                _buildMOQAlerts(),
                const SizedBox(height: 20),
                _buildContractPricing(),
                const SizedBox(height: 20),
                _buildProductionCapacity(),
                const SizedBox(height: 20),
                _buildTopBuyers(),
                const SizedBox(height: 20),
                _buildQualityCertifications(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.b2bManufacturerColor,
        icon: const Icon(Icons.handshake),
        label: const Text('New Contract'),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 140,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.b2bManufacturerColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('B2B MANUFACTURING', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            Text('Industrial Components Ltd', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.b2bManufacturerColor, AppColors.b2bManufacturerColor.withValues(alpha: 0.8)],
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Icon(Icons.precision_manufacturing, size: 80, color: Colors.white.withValues(alpha: 0.2)),
          ),
        ),
      ),
    );
  }

  Widget _buildContractOrdersCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.b2bManufacturerColor, AppColors.b2bManufacturerColor.withValues(alpha: 0.85)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Contract Orders', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusSm),
                child: Text('Q1 2026', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text('₦458,750,000', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 4),
          Text('Confirmed order value', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildOrderMetric('Active Contracts', '24'),
              _buildOrderMetric('Pending RFQs', '12'),
              _buildOrderMetric('Renewals Due', '5'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildOrderMetric(String label, String value) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(12),
        margin: const EdgeInsets.only(right: 8),
        decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: AppRadius.borderRadiusSm),
        child: Column(
          children: [
            Text(value, style: AppTypography.headlineSmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
            Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white70), textAlign: TextAlign.center),
          ],
        ),
      ),
    );
  }

  Widget _buildProductionPipeline() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Production Pipeline'),
        const SizedBox(height: 12),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
          child: Column(
            children: [
              _buildPipelineStage('Scheduled', 15, '₦125M', AppColors.info),
              _buildPipelineStage('In Production', 8, '₦89M', AppColors.warning),
              _buildPipelineStage('Quality Check', 4, '₦45M', AppColors.primary),
              _buildPipelineStage('Ready to Ship', 6, '₦72M', AppColors.success),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildPipelineStage(String stage, int orders, String value, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Container(width: 4, height: 40, decoration: BoxDecoration(color: color, borderRadius: AppRadius.borderRadiusSm)),
          const SizedBox(width: 12),
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(stage, style: AppTypography.titleSmall),
              Text('$orders orders', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            ],
          )),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
            child: Text(value, style: AppTypography.titleSmall.copyWith(color: color, fontWeight: FontWeight.w600)),
          ),
        ],
      ),
    );
  }

  Widget _buildMOQAlerts() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.warning.withValues(alpha: 0.05),
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.warning.withValues(alpha: 0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.inventory_2, color: AppColors.warning),
              const SizedBox(width: 8),
              Text('MOQ Alerts', style: AppTypography.titleMedium),
              const Spacer(),
              Text('3 pending', style: AppTypography.labelMedium.copyWith(color: AppColors.warning, fontWeight: FontWeight.w600)),
            ],
          ),
          const SizedBox(height: 12),
          _buildMOQItem('PVC Conduit 25mm', 'Customer: ElectroCorp', 'Ordered: 800 | MOQ: 1000'),
          _buildMOQItem('Steel Brackets Type-B', 'Customer: BuildMart', 'Ordered: 450 | MOQ: 500'),
          _buildMOQItem('Industrial Valves 2"', 'Customer: PetroSupply', 'Ordered: 180 | MOQ: 200'),
        ],
      ),
    );
  }

  Widget _buildMOQItem(String product, String customer, String qty) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(product, style: AppTypography.titleSmall),
              Text(customer, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
              Text(qty, style: AppTypography.labelSmall.copyWith(color: AppColors.warning)),
            ],
          )),
          TextButton.icon(onPressed: () {}, icon: const Icon(Icons.edit, size: 16), label: const Text('Adjust')),
        ],
      ),
    );
  }

  Widget _buildContractPricing() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Contract Pricing Tiers', style: AppTypography.titleMedium),
          const SizedBox(height: 16),
          _buildPricingTier('Tier 1 (1-50 units)', '₦45,000/unit', AppColors.textSecondary),
          _buildPricingTier('Tier 2 (51-200 units)', '₦42,000/unit', AppColors.info),
          _buildPricingTier('Tier 3 (201-500 units)', '₦38,500/unit', AppColors.success),
          _buildPricingTier('Enterprise (500+ units)', 'Custom Quote', AppColors.b2bManufacturerColor),
        ],
      ),
    );
  }

  Widget _buildPricingTier(String tier, String price, Color color) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Icon(Icons.local_offer, color: color, size: 20),
          const SizedBox(width: 12),
          Expanded(child: Text(tier, style: AppTypography.titleSmall)),
          Text(price, style: AppTypography.titleSmall.copyWith(color: color, fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }

  Widget _buildProductionCapacity() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Production Capacity'),
        const SizedBox(height: 12),
        GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: 2,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          childAspectRatio: 1.6,
          children: [
            StatCard(title: 'Line 1 Utilization', value: '87%', icon: Icons.precision_manufacturing, iconColor: AppColors.success),
            StatCard(title: 'Line 2 Utilization', value: '72%', icon: Icons.precision_manufacturing, iconColor: AppColors.warning),
            StatCard(title: 'Lead Time', value: '14 days', icon: Icons.schedule, iconColor: AppColors.info),
            StatCard(title: 'Capacity Available', value: '28%', icon: Icons.speed, iconColor: AppColors.primary),
          ],
        ),
      ],
    );
  }

  Widget _buildTopBuyers() {
    final buyers = [
      {'name': 'Nigerian Breweries', 'volume': '₦85M', 'orders': 12, 'trend': 15.2},
      {'name': 'Dangote Industries', 'volume': '₦72M', 'orders': 8, 'trend': 8.5},
      {'name': 'Cadbury Nigeria', 'volume': '₦58M', 'orders': 15, 'trend': -3.2},
      {'name': 'Nestle Nigeria', 'volume': '₦45M', 'orders': 6, 'trend': 22.1},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Top B2B Buyers', onViewAll: () {}),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: buyers.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, i) {
              final b = buyers[i];
              final trend = b['trend'] as double;
              return ListTile(
                leading: CircleAvatar(
                  backgroundColor: AppColors.b2bManufacturerColor.withValues(alpha: 0.1),
                  child: Text('${i + 1}', style: TextStyle(color: AppColors.b2bManufacturerColor, fontWeight: FontWeight.bold)),
                ),
                title: Text(b['name'] as String, style: AppTypography.titleSmall),
                subtitle: Text('${b['orders']} orders this quarter', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                trailing: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(b['volume'] as String, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
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

  Widget _buildQualityCertifications() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Quality Certifications', style: AppTypography.titleMedium),
          const SizedBox(height: 16),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _buildCertBadge('ISO 9001:2015', AppColors.success),
              _buildCertBadge('ISO 14001', AppColors.success),
              _buildCertBadge('SON Certified', AppColors.success),
              _buildCertBadge('CE Mark', AppColors.success),
              _buildCertBadge('API Q1', AppColors.warning, dueIn: '45 days'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildCertBadge(String name, Color color, {String? dueIn}) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm, border: Border.all(color: color.withValues(alpha: 0.3))),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.verified, color: color, size: 16),
          const SizedBox(width: 4),
          Text(name, style: AppTypography.labelSmall.copyWith(color: color, fontWeight: FontWeight.w600)),
          if (dueIn != null) ...[
            const SizedBox(width: 8),
            Text('(Renewal: $dueIn)', style: AppTypography.labelSmall.copyWith(color: AppColors.warning)),
          ],
        ],
      ),
    );
  }

  Widget _buildQuickActions() {
    return QuickActionGrid(
      actions: [
        QuickActionItem(icon: Icons.request_quote, label: 'New Quote', color: AppColors.success, onTap: () {}),
        QuickActionItem(icon: Icons.production_quantity_limits, label: 'Schedule', color: AppColors.warning, onTap: () {}),
        QuickActionItem(icon: Icons.inventory, label: 'Raw Materials', color: AppColors.info, onTap: () {}),
        QuickActionItem(icon: Icons.qr_code, label: 'Batch Track', color: AppColors.primary, onTap: () {}),
        QuickActionItem(icon: Icons.analytics, label: 'Reports', color: AppColors.b2bManufacturerColor, onTap: () {}),
        QuickActionItem(icon: Icons.support, label: 'Support', color: AppColors.error, onTap: () {}),
      ],
    );
  }
}
