/// OmniRoute Ecosystem - Bank Dashboard Screen
/// Comprehensive dashboard for banks and financial institutions

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class BankDashboardScreen extends ConsumerStatefulWidget {
  const BankDashboardScreen({super.key});

  @override
  ConsumerState<BankDashboardScreen> createState() => _BankDashboardScreenState();
}

class _BankDashboardScreenState extends ConsumerState<BankDashboardScreen> {
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
                _buildPortfolioCard(),
                const SizedBox(height: 20),
                _buildStatsGrid(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildLoanPerformanceChart(),
                const SizedBox(height: 20),
                _buildRecentSettlements(),
                const SizedBox(height: 20),
                _buildPendingApprovals(),
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
      backgroundColor: AppColors.bankColor,
      flexibleSpace: FlexibleSpaceBar(
        titlePadding: const EdgeInsets.only(left: 16, bottom: 16),
        title: Text('Financial Dashboard',
            style: AppTypography.titleLarge.copyWith(color: Colors.white, fontSize: 18)),
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [AppColors.bankColor, AppColors.bankColor.withValues(alpha: 0.8)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
        ),
      ),
      actions: [
        IconButton(icon: const Icon(Icons.notifications_outlined, color: Colors.white), onPressed: () {}),
        IconButton(icon: const Icon(Icons.settings_outlined, color: Colors.white), onPressed: () {}),
      ],
    );
  }

  Widget _buildPortfolioCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [AppColors.bankColor, AppColors.bankColor.withValues(alpha: 0.85)],
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
              Text('Total Portfolio', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: AppColors.success.withValues(alpha: 0.2),
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Row(
                  children: [
                    const Icon(Icons.trending_up, color: AppColors.successLight, size: 14),
                    const SizedBox(width: 4),
                    Text('+12.5%', style: AppTypography.labelSmall.copyWith(color: AppColors.successLight)),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text('₦2,458,320,000', style: AppTypography.displaySmall.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
          const SizedBox(height: 24),
          Row(
            children: [
              _PortfolioMetric(label: 'Active Loans', value: '1,234', icon: Icons.receipt_long),
              const SizedBox(width: 24),
              _PortfolioMetric(label: 'Default Rate', value: '2.3%', icon: Icons.warning_amber),
              const SizedBox(width: 24),
              _PortfolioMetric(label: 'NPL Ratio', value: '1.8%', icon: Icons.pie_chart),
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
          title: 'Disbursed This Month',
          value: '₦458M',
          icon: Icons.payments,
          iconColor: AppColors.success,
          growthPercentage: 18.5,
        ),
        StatCard(
          title: 'Repayments Received',
          value: '₦312M',
          icon: Icons.account_balance_wallet,
          iconColor: AppColors.primary,
          growthPercentage: 8.2,
        ),
        StatCard(
          title: 'Pending Approvals',
          value: '47',
          icon: Icons.pending_actions,
          iconColor: AppColors.warning,
        ),
        StatCard(
          title: 'Settlement Due',
          value: '₦89M',
          icon: Icons.schedule,
          iconColor: AppColors.info,
        ),
      ].map((e) => e.animate().fadeIn().scale(begin: const Offset(0.95, 0.95))).toList(),
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
            QuickActionItem(icon: Icons.add_circle, label: 'New Loan', color: AppColors.success, onTap: () {}),
            QuickActionItem(icon: Icons.sync, label: 'Process Settlement', color: AppColors.primary, onTap: () {}),
            QuickActionItem(icon: Icons.assessment, label: 'Risk Report', color: AppColors.warning, onTap: () {}),
            QuickActionItem(icon: Icons.people, label: 'Borrowers', color: AppColors.info, onTap: () {}, badge: '12'),
          ],
        ),
      ],
    );
  }

  Widget _buildLoanPerformanceChart() {
    return Container(
      padding: const EdgeInsets.all(20),
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
              Text('Loan Performance', style: AppTypography.titleMedium),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: AppColors.grey100,
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Text('Last 6 months', style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
              ),
            ],
          ),
          const SizedBox(height: 24),
          SizedBox(
            height: 200,
            child: LineChart(
              LineChartData(
                gridData: FlGridData(show: true, drawVerticalLine: false),
                titlesData: FlTitlesData(
                  leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  bottomTitles: AxisTitles(
                    sideTitles: SideTitles(
                      showTitles: true,
                      getTitlesWidget: (value, meta) {
                        const months = ['Aug', 'Sep', 'Oct', 'Nov', 'Dec', 'Jan'];
                        return Text(months[value.toInt() % 6], style: AppTypography.labelSmall);
                      },
                    ),
                  ),
                ),
                borderData: FlBorderData(show: false),
                lineBarsData: [
                  LineChartBarData(
                    spots: const [
                      FlSpot(0, 3), FlSpot(1, 3.5), FlSpot(2, 4), FlSpot(3, 3.8), FlSpot(4, 4.2), FlSpot(5, 5),
                    ],
                    isCurved: true,
                    color: AppColors.bankColor,
                    barWidth: 3,
                    dotData: FlDotData(show: false),
                    belowBarData: BarAreaData(
                      show: true,
                      color: AppColors.bankColor.withValues(alpha: 0.1),
                    ),
                  ),
                  LineChartBarData(
                    spots: const [
                      FlSpot(0, 2), FlSpot(1, 2.2), FlSpot(2, 2.5), FlSpot(3, 2.3), FlSpot(4, 2.8), FlSpot(5, 3),
                    ],
                    isCurved: true,
                    color: AppColors.success,
                    barWidth: 3,
                    dotData: FlDotData(show: false),
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              _ChartLegend(color: AppColors.bankColor, label: 'Disbursed'),
              const SizedBox(width: 24),
              _ChartLegend(color: AppColors.success, label: 'Repaid'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildRecentSettlements() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Recent Settlements', actionText: 'View All', onAction: () {}),
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
              Container(
                width: 48, height: 48,
                decoration: BoxDecoration(color: AppColors.successBg, borderRadius: AppRadius.borderRadiusSm),
                child: const Icon(Icons.check_circle, color: AppColors.success),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text('Settlement #STL-${1000 + i}', style: AppTypography.titleSmall),
                    Text('Wholesaler - Dangote Foods Ltd', style: AppTypography.bodySmall),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text('₦${(12 + i * 3)}M', style: AppTypography.titleSmall.copyWith(color: AppColors.success)),
                  Text('${i + 1}h ago', style: AppTypography.labelSmall.copyWith(color: AppColors.grey500)),
                ],
              ),
            ],
          ),
        )),
      ],
    );
  }

  Widget _buildPendingApprovals() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Pending Approvals', actionText: 'View All', onAction: () {}),
        const SizedBox(height: 12),
        ...List.generate(2, (i) => Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.warning.withValues(alpha: 0.3)),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Loan Application #LA-${2000 + i}', style: AppTypography.titleSmall),
                  const StatusChip(status: 'pending'),
                ],
              ),
              const SizedBox(height: 8),
              Text('Retailer - Mama Ngozi Stores', style: AppTypography.bodySmall),
              const SizedBox(height: 12),
              Row(
                children: [
                  _LoanDetail(label: 'Amount', value: '₦${5 + i * 2}M'),
                  const SizedBox(width: 24),
                  _LoanDetail(label: 'Tenor', value: '${30 + i * 15} days'),
                  const SizedBox(width: 24),
                  _LoanDetail(label: 'Score', value: '${750 + i * 25}'),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(onPressed: () {}, child: const Text('Reject')),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: ElevatedButton(onPressed: () {}, child: const Text('Approve')),
                  ),
                ],
              ),
            ],
          ),
        )),
      ],
    );
  }
}

class _PortfolioMetric extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;

  const _PortfolioMetric({required this.label, required this.value, required this.icon});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: Colors.white70, size: 16),
        const SizedBox(width: 6),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(value, style: AppTypography.titleSmall.copyWith(color: Colors.white)),
            Text(label, style: AppTypography.labelSmall.copyWith(color: Colors.white60)),
          ],
        ),
      ],
    );
  }
}

class _ChartLegend extends StatelessWidget {
  final Color color;
  final String label;

  const _ChartLegend({required this.color, required this.label});

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Container(width: 12, height: 12, decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(3))),
        const SizedBox(width: 6),
        Text(label, style: AppTypography.labelSmall),
      ],
    );
  }
}

class _LoanDetail extends StatelessWidget {
  final String label;
  final String value;

  const _LoanDetail({required this.label, required this.value});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.grey500)),
        Text(value, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
      ],
    );
  }
}
