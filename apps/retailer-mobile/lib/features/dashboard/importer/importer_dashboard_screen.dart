/// OmniRoute Ecosystem - Importer Dashboard Screen
/// Comprehensive dashboard for importers and trade companies
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class ImporterDashboardScreen extends ConsumerStatefulWidget {
  const ImporterDashboardScreen({super.key});
  @override
  ConsumerState<ImporterDashboardScreen> createState() => _ImporterDashboardScreenState();
}

class _ImporterDashboardScreenState extends ConsumerState<ImporterDashboardScreen> {
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
                _buildTradeFinanceCard(),
                const SizedBox(height: 20),
                _buildShipmentTracker(),
                const SizedBox(height: 20),
                _buildCustomsClearance(),
                const SizedBox(height: 20),
                _buildForexExposure(),
                const SizedBox(height: 20),
                _buildLetterOfCreditStatus(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildSupplierPayments(),
                const SizedBox(height: 20),
                _buildIncomingContainers(),
                const SizedBox(height: 100),
              ]),
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () {},
        backgroundColor: AppColors.importerColor,
        icon: const Icon(Icons.add),
        label: const Text('New Import Order'),
      ),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      expandedHeight: 140,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.importerColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Column(
          mainAxisAlignment: MainAxisAlignment.end,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('IMPORT OPERATIONS', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
            Text('Global Trade Corp', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
          ],
        ),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.importerColor, AppColors.importerColor.withValues(alpha: 0.8)],
            ),
          ),
          child: Align(
            alignment: Alignment.centerRight,
            child: Icon(Icons.public, size: 80, color: Colors.white.withValues(alpha: 0.2)),
          ),
        ),
      ),
    );
  }

  Widget _buildTradeFinanceCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.importerColor, AppColors.importerColor.withValues(alpha: 0.85)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Trade Finance Facility', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: AppColors.success, borderRadius: AppRadius.borderRadiusSm),
                child: Text('Active', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text('\$2.5M', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
              Text(' / \$5M limit', style: AppTypography.titleMedium.copyWith(color: Colors.white70)),
            ],
          ),
          const SizedBox(height: 12),
          LinearProgressIndicator(
            value: 0.5,
            backgroundColor: Colors.white.withValues(alpha: 0.3),
            valueColor: const AlwaysStoppedAnimation(Colors.white),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildFinanceMetric('Open LCs', '\$1.8M'),
              _buildFinanceMetric('In Transit', '\$450K'),
              _buildFinanceMetric('Cleared', '\$250K'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildFinanceMetric(String label, String value) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(12),
        margin: const EdgeInsets.only(right: 8),
        decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.15), borderRadius: AppRadius.borderRadiusSm),
        child: Column(
          children: [
            Text(value, style: AppTypography.titleMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
            Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
          ],
        ),
      ),
    );
  }

  Widget _buildShipmentTracker() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Active Shipments'),
        const SizedBox(height: 12),
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
          child: Column(
            children: [
              _buildShipmentItem('MV Atlantic Star', 'Shanghai → Lagos', 'At Sea - 12 days left', 0.6, AppColors.info),
              const Divider(),
              _buildShipmentItem('MV Ocean Pride', 'Dubai → Apapa', 'Customs Hold', 0.85, AppColors.warning),
              const Divider(),
              _buildShipmentItem('MV Global Trader', 'Rotterdam → Tincan', 'Arrived - Pending Clearance', 1.0, AppColors.success),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildShipmentItem(String vessel, String route, String status, double progress, Color statusColor) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(children: [
                const Icon(Icons.directions_boat, size: 20, color: AppColors.importerColor),
                const SizedBox(width: 8),
                Text(vessel, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
              ]),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: statusColor.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
                child: Text(status, style: AppTypography.labelSmall.copyWith(color: statusColor, fontWeight: FontWeight.w600)),
              ),
            ],
          ),
          const SizedBox(height: 4),
          Text(route, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
          const SizedBox(height: 8),
          LinearProgressIndicator(
            value: progress,
            backgroundColor: statusColor.withValues(alpha: 0.2),
            valueColor: AlwaysStoppedAnimation(statusColor),
          ),
        ],
      ),
    );
  }

  Widget _buildCustomsClearance() {
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
              const Icon(Icons.gavel, color: AppColors.warning),
              const SizedBox(width: 8),
              Text('Customs Clearance', style: AppTypography.titleMedium),
              const Spacer(),
              Text('3 pending', style: AppTypography.labelMedium.copyWith(color: AppColors.warning, fontWeight: FontWeight.w600)),
            ],
          ),
          const SizedBox(height: 16),
          _buildCustomsItem('Container MSKU-4521879', 'Duty Assessment', '₦4.2M', 'Pay Now'),
          _buildCustomsItem('Container TCLU-8976543', 'Document Review', '-', 'View'),
          _buildCustomsItem('Container OOLU-3345678', 'Physical Inspection', '-', 'Track'),
        ],
      ),
    );
  }

  Widget _buildCustomsItem(String container, String stage, String duty, String action) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(container, style: AppTypography.titleSmall),
              Text(stage, style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            ],
          )),
          if (duty != '-') Text(duty, style: AppTypography.titleSmall.copyWith(color: AppColors.warning, fontWeight: FontWeight.w600)),
          const SizedBox(width: 8),
          TextButton(onPressed: () {}, child: Text(action)),
        ],
      ),
    );
  }

  Widget _buildForexExposure() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Forex Exposure', style: AppTypography.titleMedium),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildForexCard('USD', '\$1,245,000', '₦1,931/\$', true),
              const SizedBox(width: 12),
              _buildForexCard('EUR', '€320,000', '₦2,105/€', false),
              const SizedBox(width: 12),
              _buildForexCard('CNY', '¥2,800,000', '₦265/¥', true),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildForexCard(String currency, String exposure, String rate, bool favorable) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: favorable ? AppColors.success.withValues(alpha: 0.1) : AppColors.error.withValues(alpha: 0.1),
          borderRadius: AppRadius.borderRadiusSm,
        ),
        child: Column(
          children: [
            Text(currency, style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.w700)),
            const SizedBox(height: 4),
            Text(exposure, style: AppTypography.labelSmall),
            const SizedBox(height: 4),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(favorable ? Icons.arrow_downward : Icons.arrow_upward, size: 12, color: favorable ? AppColors.success : AppColors.error),
                Text(rate, style: AppTypography.labelSmall.copyWith(color: favorable ? AppColors.success : AppColors.error)),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildLetterOfCreditStatus() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Letters of Credit', onViewAll: () {}),
        const SizedBox(height: 12),
        GridView.count(
          shrinkWrap: true,
          physics: const NeverScrollableScrollPhysics(),
          crossAxisCount: 2,
          crossAxisSpacing: 12,
          mainAxisSpacing: 12,
          childAspectRatio: 1.6,
          children: [
            StatCard(title: 'Active LCs', value: '8', icon: Icons.description, iconColor: AppColors.success),
            StatCard(title: 'Pending Issuance', value: '2', icon: Icons.pending, iconColor: AppColors.warning),
            StatCard(title: 'Expiring Soon', value: '3', icon: Icons.timer_off, iconColor: AppColors.error),
            StatCard(title: 'This Month Value', value: '\$850K', icon: Icons.attach_money, iconColor: AppColors.info),
          ],
        ),
      ],
    );
  }

  Widget _buildQuickActions() {
    return QuickActionGrid(
      actions: [
        QuickActionItem(icon: Icons.add_circle, label: 'Open LC', color: AppColors.success, onTap: () {}),
        QuickActionItem(icon: Icons.calculate, label: 'Duty Calc', color: AppColors.warning, onTap: () {}),
        QuickActionItem(icon: Icons.track_changes, label: 'Track Vessel', color: AppColors.info, onTap: () {}),
        QuickActionItem(icon: Icons.currency_exchange, label: 'Forex Rate', color: AppColors.primary, onTap: () {}),
        QuickActionItem(icon: Icons.receipt_long, label: 'Documents', color: AppColors.importerColor, onTap: () {}),
        QuickActionItem(icon: Icons.support_agent, label: 'Agent', color: AppColors.error, onTap: () {}),
      ],
    );
  }

  Widget _buildSupplierPayments() {
    final payments = [
      {'supplier': 'Shenzhen Electronics Co.', 'amount': '\$125,000', 'due': '5 days', 'status': 'pending'},
      {'supplier': 'Dubai Trading FZE', 'amount': '\$89,500', 'due': '12 days', 'status': 'scheduled'},
      {'supplier': 'Rotterdam Commodities', 'amount': '€45,000', 'due': 'Paid', 'status': 'completed'},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Supplier Payments', onViewAll: () {}),
        const SizedBox(height: 12),
        ...payments.map((p) {
          final isPaid = p['status'] == 'completed';
          return Container(
            margin: const EdgeInsets.only(bottom: 8),
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
            child: Row(
              children: [
                CircleAvatar(
                  backgroundColor: isPaid ? AppColors.success.withValues(alpha: 0.1) : AppColors.warning.withValues(alpha: 0.1),
                  child: Icon(isPaid ? Icons.check : Icons.schedule, color: isPaid ? AppColors.success : AppColors.warning, size: 20),
                ),
                const SizedBox(width: 12),
                Expanded(child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(p['supplier']!, style: AppTypography.titleSmall),
                    Text('Due: ${p['due']}', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                  ],
                )),
                Text(p['amount']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
              ],
            ),
          );
        }),
      ],
    );
  }

  Widget _buildIncomingContainers() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Incoming Containers', style: AppTypography.titleMedium),
              Text('Next 30 days', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _buildContainerStat('20ft', '12'),
              _buildContainerStat('40ft', '8'),
              _buildContainerStat('40ft HC', '5'),
              _buildContainerStat('Reefer', '2'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildContainerStat(String type, String count) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(color: AppColors.importerColor.withValues(alpha: 0.1), borderRadius: AppRadius.borderRadiusSm),
          child: Text(count, style: AppTypography.headlineSmall.copyWith(color: AppColors.importerColor, fontWeight: FontWeight.w700)),
        ),
        const SizedBox(height: 4),
        Text(type, style: AppTypography.labelSmall),
      ],
    );
  }
}
