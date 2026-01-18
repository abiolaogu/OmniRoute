/// OmniRoute Ecosystem - Investor Dashboard
/// Investment portfolio tracking with opportunities and ROI analytics

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class InvestorDashboardScreen extends ConsumerStatefulWidget {
  const InvestorDashboardScreen({super.key});

  @override
  ConsumerState<InvestorDashboardScreen> createState() =>
      _InvestorDashboardScreenState();
}

class _InvestorDashboardScreenState
    extends ConsumerState<InvestorDashboardScreen> {
  @override
  Widget build(BuildContext context) {
    final user = ref.watch(authProvider).user;

    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        backgroundColor: AppColors.white,
        elevation: 0,
        leading: Builder(
          builder: (context) => IconButton(
            icon: const Icon(Icons.menu, color: AppColors.grey800),
            onPressed: () => Scaffold.of(context).openDrawer(),
          ),
        ),
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Welcome back,',
              style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
            ),
            Text(
              user?.fullName ?? 'Investor',
              style: AppTypography.titleMedium.copyWith(color: AppColors.grey900),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications_outlined, color: AppColors.grey800),
            onPressed: () {},
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {},
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Portfolio Summary Card
              _buildPortfolioCard()
                  .animate()
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Key Metrics
              _buildKeyMetrics()
                  .animate(delay: const Duration(milliseconds: 100))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Quick Actions
              SectionHeader(title: 'Quick Actions'),
              const SizedBox(height: 12),
              QuickActionGrid(
                actions: [
                  QuickActionItem(
                    icon: Icons.add_chart,
                    label: 'Invest',
                    color: AppColors.investorColor,
                    onTap: () {},
                  ),
                  QuickActionItem(
                    icon: Icons.trending_up,
                    label: 'Portfolio',
                    color: AppColors.success,
                    onTap: () {},
                  ),
                  QuickActionItem(
                    icon: Icons.explore,
                    label: 'Discover',
                    color: AppColors.info,
                    onTap: () {},
                    badge: '12',
                  ),
                  QuickActionItem(
                    icon: Icons.analytics,
                    label: 'Reports',
                    color: AppColors.warning,
                    onTap: () {},
                  ),
                ],
              )
                  .animate(delay: const Duration(milliseconds: 200))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Portfolio Performance Chart
              SectionHeader(
                title: 'Portfolio Performance',
                actionText: 'Details',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildPerformanceChart()
                  .animate(delay: const Duration(milliseconds: 300))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Active Investments
              SectionHeader(
                title: 'Active Investments',
                actionText: 'View All',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildActiveInvestments()
                  .animate(delay: const Duration(milliseconds: 400))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Investment Opportunities
              SectionHeader(
                title: 'New Opportunities',
                actionText: 'Explore',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildOpportunities()
                  .animate(delay: const Duration(milliseconds: 500))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Recent Returns
              SectionHeader(
                title: 'Recent Returns',
                actionText: 'History',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildRecentReturns()
                  .animate(delay: const Duration(milliseconds: 600))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildPortfolioCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF283593), Color(0xFF3949AB)],
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
              Text(
                'Total Portfolio Value',
                style: AppTypography.labelMedium.copyWith(
                  color: Colors.white.withValues(alpha: 0.8),
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                decoration: BoxDecoration(
                  color: AppColors.success.withValues(alpha: 0.2),
                  borderRadius: AppRadius.borderRadiusFull,
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const Icon(Icons.trending_up, color: Colors.greenAccent, size: 14),
                    const SizedBox(width: 4),
                    Text(
                      '+18.5%',
                      style: AppTypography.labelSmall.copyWith(
                        color: Colors.greenAccent,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            '₦45,750,000',
            style: AppTypography.displaySmall.copyWith(
              color: Colors.white,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            'Unrealized gains: ₦7,125,000',
            style: AppTypography.bodySmall.copyWith(
              color: Colors.white.withValues(alpha: 0.7),
            ),
          ),
          const SizedBox(height: 20),
          Row(
            children: [
              Expanded(
                child: _buildPortfolioStat('Invested', '₦38.6M'),
              ),
              Container(
                width: 1,
                height: 40,
                color: Colors.white.withValues(alpha: 0.2),
              ),
              Expanded(
                child: _buildPortfolioStat('Returns', '₦7.1M'),
              ),
              Container(
                width: 1,
                height: 40,
                color: Colors.white.withValues(alpha: 0.2),
              ),
              Expanded(
                child: _buildPortfolioStat('ROI', '18.5%'),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildPortfolioStat(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: AppTypography.titleMedium.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.w600,
          ),
        ),
        Text(
          label,
          style: AppTypography.labelSmall.copyWith(
            color: Colors.white.withValues(alpha: 0.7),
          ),
        ),
      ],
    );
  }

  Widget _buildKeyMetrics() {
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.3,
      children: [
        StatCard(
          title: 'Active Investments',
          value: '8',
          icon: Icons.business_center,
          iconColor: AppColors.investorColor,
        ),
        StatCard(
          title: 'Avg. Monthly Return',
          value: '2.3%',
          icon: Icons.show_chart,
          iconColor: AppColors.success,
          growthPercentage: 0.5,
        ),
        StatCard(
          title: 'This Month Returns',
          value: '₦892K',
          icon: Icons.account_balance_wallet,
          iconColor: AppColors.info,
          growthPercentage: 12.3,
        ),
        StatCard(
          title: 'Pending Payouts',
          value: '₦450K',
          icon: Icons.schedule,
          iconColor: AppColors.warning,
        ),
      ],
    );
  }

  Widget _buildPerformanceChart() {
    return Container(
      height: 200,
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: LineChart(
        LineChartData(
          gridData: FlGridData(
            show: true,
            drawVerticalLine: false,
            horizontalInterval: 10,
            getDrawingHorizontalLine: (value) {
              return FlLine(
                color: AppColors.grey200,
                strokeWidth: 1,
              );
            },
          ),
          titlesData: FlTitlesData(
            show: true,
            bottomTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                reservedSize: 30,
                getTitlesWidget: (value, meta) {
                  const months = ['J', 'F', 'M', 'A', 'M', 'J'];
                  if (value.toInt() >= 0 && value.toInt() < months.length) {
                    return Text(
                      months[value.toInt()],
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.grey500,
                      ),
                    );
                  }
                  return const SizedBox();
                },
              ),
            ),
            leftTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
            topTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
            rightTitles: const AxisTitles(
              sideTitles: SideTitles(showTitles: false),
            ),
          ),
          borderData: FlBorderData(show: false),
          lineBarsData: [
            LineChartBarData(
              spots: const [
                FlSpot(0, 30),
                FlSpot(1, 35),
                FlSpot(2, 32),
                FlSpot(3, 40),
                FlSpot(4, 38),
                FlSpot(5, 45),
              ],
              isCurved: true,
              color: AppColors.investorColor,
              barWidth: 3,
              isStrokeCapRound: true,
              dotData: const FlDotData(show: false),
              belowBarData: BarAreaData(
                show: true,
                color: AppColors.investorColor.withValues(alpha: 0.1),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildActiveInvestments() {
    final investments = [
      {
        'name': 'Lagos Wholesale Hub',
        'type': 'Revenue Share',
        'invested': 5000000,
        'returns': 875000,
        'roi': 17.5,
      },
      {
        'name': 'Quick Logistics Ltd',
        'type': 'Equity',
        'invested': 10000000,
        'returns': 2200000,
        'roi': 22.0,
      },
      {
        'name': 'FreshMart Chain',
        'type': 'Debt',
        'invested': 8000000,
        'returns': 1120000,
        'roi': 14.0,
      },
    ];

    return Column(
      children: investments.map((inv) {
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
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
                  Container(
                    width: 48,
                    height: 48,
                    decoration: BoxDecoration(
                      color: AppColors.investorColor.withValues(alpha: 0.1),
                      borderRadius: AppRadius.borderRadiusSm,
                    ),
                    child: Center(
                      child: Text(
                        (inv['name'] as String).substring(0, 2).toUpperCase(),
                        style: AppTypography.titleSmall.copyWith(
                          color: AppColors.investorColor,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          inv['name'] as String,
                          style: AppTypography.titleSmall,
                        ),
                        Text(
                          inv['type'] as String,
                          style: AppTypography.bodySmall.copyWith(
                            color: AppColors.grey500,
                          ),
                        ),
                      ],
                    ),
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text(
                        '+${(inv['roi'] as double).toStringAsFixed(1)}%',
                        style: AppTypography.titleSmall.copyWith(
                          color: AppColors.success,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      Text(
                        'ROI',
                        style: AppTypography.labelSmall.copyWith(
                          color: AppColors.grey500,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Invested',
                          style: AppTypography.labelSmall.copyWith(
                            color: AppColors.grey500,
                          ),
                        ),
                        Text(
                          formatCurrency((inv['invested'] as int).toDouble()),
                          style: AppTypography.labelMedium,
                        ),
                      ],
                    ),
                  ),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.end,
                      children: [
                        Text(
                          'Returns',
                          style: AppTypography.labelSmall.copyWith(
                            color: AppColors.grey500,
                          ),
                        ),
                        Text(
                          formatCurrency((inv['returns'] as int).toDouble()),
                          style: AppTypography.labelMedium.copyWith(
                            color: AppColors.success,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Widget _buildOpportunities() {
    final opportunities = [
      {
        'name': 'FMCG Distribution Network',
        'target': 25000000,
        'raised': 18000000,
        'roi': '18-22%',
        'term': '12 months',
      },
      {
        'name': 'Cold Chain Logistics',
        'target': 50000000,
        'raised': 35000000,
        'roi': '20-25%',
        'term': '18 months',
      },
    ];

    return Column(
      children: opportunities.map((opp) {
        final progress = (opp['raised'] as int) / (opp['target'] as int);
        return Container(
          margin: const EdgeInsets.only(bottom: 12),
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.investorColor.withValues(alpha: 0.3)),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      opp['name'] as String,
                      style: AppTypography.titleSmall,
                    ),
                  ),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: AppColors.successBg,
                      borderRadius: AppRadius.borderRadiusFull,
                    ),
                    child: Text(
                      opp['roi'] as String,
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.success,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              LinearProgressIndicator(
                value: progress,
                backgroundColor: AppColors.grey200,
                valueColor: const AlwaysStoppedAnimation<Color>(AppColors.investorColor),
                minHeight: 8,
                borderRadius: BorderRadius.circular(4),
              ),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '${formatCurrency((opp['raised'] as int).toDouble())} raised',
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.grey600,
                    ),
                  ),
                  Text(
                    'Target: ${formatCurrency((opp['target'] as int).toDouble())}',
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.grey600,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  const Icon(Icons.schedule, size: 14, color: AppColors.grey500),
                  const SizedBox(width: 4),
                  Text(
                    opp['term'] as String,
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.grey600,
                    ),
                  ),
                  const Spacer(),
                  TextButton(
                    onPressed: () {},
                    style: TextButton.styleFrom(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      backgroundColor: AppColors.investorColor,
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(8),
                      ),
                    ),
                    child: const Text('Invest Now'),
                  ),
                ],
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Widget _buildRecentReturns() {
    final returns = [
      {'source': 'Lagos Wholesale Hub', 'amount': 125000, 'date': 'Today'},
      {'source': 'Quick Logistics Ltd', 'amount': 220000, 'date': 'Yesterday'},
      {'source': 'FreshMart Chain', 'amount': 140000, 'date': '3 days ago'},
    ];

    return Column(
      children: returns.map((ret) {
        return Container(
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusSm,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: Row(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: AppColors.successBg,
                  shape: BoxShape.circle,
                ),
                child: const Icon(Icons.arrow_downward, color: AppColors.success, size: 20),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      ret['source'] as String,
                      style: AppTypography.labelMedium,
                    ),
                    Text(
                      ret['date'] as String,
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.grey500,
                      ),
                    ),
                  ],
                ),
              ),
              Text(
                '+${formatCurrency((ret['amount'] as int).toDouble())}',
                style: AppTypography.titleSmall.copyWith(
                  color: AppColors.success,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }
}
