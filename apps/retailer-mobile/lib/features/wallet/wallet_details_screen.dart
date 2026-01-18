/// OmniRoute Ecosystem - Wallet Screen
/// Comprehensive wallet management with transactions, settlements, and loans

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class WalletDetailsScreen extends ConsumerStatefulWidget {
  const WalletDetailsScreen({super.key});

  @override
  ConsumerState<WalletDetailsScreen> createState() => _WalletDetailsScreenState();
}

class _WalletDetailsScreenState extends ConsumerState<WalletDetailsScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  String _selectedPeriod = 'This Week';

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

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
                _buildBalanceCard(),
                const SizedBox(height: 20),
                _buildQuickActions(),
                const SizedBox(height: 20),
                _buildCashFlowChart(),
                const SizedBox(height: 20),
                _buildFinancialServices(),
                const SizedBox(height: 20),
                _buildRecentTransactions(),
                const SizedBox(height: 20),
                _buildPendingSettlements(),
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
      expandedHeight: 60,
      floating: false,
      pinned: true,
      backgroundColor: AppColors.white,
      elevation: 0,
      title: Text(
        'Wallet',
        style: AppTypography.titleLarge.copyWith(color: AppColors.grey900),
      ),
      actions: [
        IconButton(
          icon: const Icon(Icons.history, color: AppColors.grey800),
          onPressed: () {},
        ),
        IconButton(
          icon: const Icon(Icons.more_vert, color: AppColors.grey800),
          onPressed: () {},
        ),
      ],
    );
  }

  Widget _buildBalanceCard() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF1A237E), Color(0xFF283593)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: AppRadius.borderRadiusLg,
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF1A237E).withValues(alpha: 0.3),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    'Available Balance',
                    style: AppTypography.labelMedium.copyWith(
                      color: Colors.white.withValues(alpha: 0.7),
                    ),
                  ),
                  const SizedBox(height: 4),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text(
                        '₦3,456,780',
                        style: AppTypography.displaySmall.copyWith(
                          color: Colors.white,
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        '.50',
                        style: AppTypography.titleLarge.copyWith(
                          color: Colors.white.withValues(alpha: 0.7),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.15),
                  shape: BoxShape.circle,
                ),
                child: const Icon(
                  Icons.account_balance_wallet,
                  color: Colors.white,
                  size: 28,
                ),
              ),
            ],
          ),
          const SizedBox(height: 20),
          Row(
            children: [
              _BalanceInfo(
                label: 'Pending',
                value: '₦245,000',
                icon: Icons.pending_actions,
              ),
              const SizedBox(width: 24),
              _BalanceInfo(
                label: 'Credit Limit',
                value: '₦500,000',
                icon: Icons.credit_score,
              ),
            ],
          ),
          const SizedBox(height: 20),
          Row(
            children: [
              Expanded(
                child: _WalletButton(
                  icon: Icons.add,
                  label: 'Add Money',
                  onTap: () => _showAddMoneySheet(context),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _WalletButton(
                  icon: Icons.send,
                  label: 'Transfer',
                  onTap: () => _showTransferSheet(context),
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: _WalletButton(
                  icon: Icons.account_balance,
                  label: 'Withdraw',
                  onTap: () => _showWithdrawSheet(context),
                ),
              ),
            ],
          ),
        ],
      ),
    ).animate().fadeIn().slideY(begin: 0.2, end: 0);
  }

  Widget _buildQuickActions() {
    return Row(
      children: [
        Expanded(
          child: _ActionCard(
            icon: Icons.receipt_long,
            label: 'Pay Invoice',
            color: AppColors.primary,
            onTap: () {},
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _ActionCard(
            icon: Icons.qr_code_scanner,
            label: 'Scan to Pay',
            color: AppColors.success,
            onTap: () {},
          ),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: _ActionCard(
            icon: Icons.schedule,
            label: 'Schedule',
            color: AppColors.info,
            onTap: () {},
          ),
        ),
      ],
    ).animate(delay: const Duration(milliseconds: 200)).fadeIn().slideY(begin: 0.1, end: 0);
  }

  Widget _buildCashFlowChart() {
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
              Text('Cash Flow', style: AppTypography.titleMedium),
              DropdownButton<String>(
                value: _selectedPeriod,
                underline: const SizedBox(),
                items: ['Today', 'This Week', 'This Month', 'This Year']
                    .map((e) => DropdownMenuItem(value: e, child: Text(e)))
                    .toList(),
                onChanged: (value) => setState(() => _selectedPeriod = value!),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              _CashFlowIndicator(
                label: 'Income',
                value: '₦5.2M',
                color: AppColors.success,
                icon: Icons.trending_up,
              ),
              const SizedBox(width: 24),
              _CashFlowIndicator(
                label: 'Expenses',
                value: '₦3.1M',
                color: AppColors.error,
                icon: Icons.trending_down,
              ),
            ],
          ),
          const SizedBox(height: 24),
          SizedBox(
            height: 200,
            child: BarChart(
              BarChartData(
                barGroups: _getCashFlowData(),
                gridData: FlGridData(show: false),
                borderData: FlBorderData(show: false),
                titlesData: FlTitlesData(
                  leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                  bottomTitles: AxisTitles(
                    sideTitles: SideTitles(
                      showTitles: true,
                      getTitlesWidget: (value, meta) {
                        const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
                        return Padding(
                          padding: const EdgeInsets.only(top: 8),
                          child: Text(days[value.toInt()], style: AppTypography.labelSmall),
                        );
                      },
                    ),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  List<BarChartGroupData> _getCashFlowData() {
    return List.generate(7, (index) {
      return BarChartGroupData(
        x: index,
        barRods: [
          BarChartRodData(
            toY: [120, 150, 80, 180, 200, 160, 90][index].toDouble(),
            color: AppColors.success,
            width: 12,
            borderRadius: const BorderRadius.vertical(top: Radius.circular(4)),
          ),
          BarChartRodData(
            toY: [80, 90, 60, 100, 120, 80, 70][index].toDouble(),
            color: AppColors.error.withValues(alpha: 0.6),
            width: 12,
            borderRadius: const BorderRadius.vertical(top: Radius.circular(4)),
          ),
        ],
        barsSpace: 4,
      );
    });
  }

  Widget _buildFinancialServices() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Financial Services'),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: _ServiceCard(
                icon: Icons.credit_card,
                title: 'Apply for Loan',
                subtitle: 'Up to ₦5M at 2.5%',
                color: AppColors.primary,
                onTap: () => _showLoanApplicationSheet(context),
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _ServiceCard(
                icon: Icons.calendar_today,
                title: 'Buy Now Pay Later',
                subtitle: 'Split payments',
                color: AppColors.secondary,
                onTap: () {},
              ),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Row(
          children: [
            Expanded(
              child: _ServiceCard(
                icon: Icons.savings,
                title: 'Savings',
                subtitle: 'Earn 12% p.a.',
                color: AppColors.success,
                onTap: () {},
              ),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: _ServiceCard(
                icon: Icons.shield,
                title: 'Insurance',
                subtitle: 'Protect your business',
                color: AppColors.warning,
                onTap: () {},
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildRecentTransactions() {
    final transactions = [
      {'type': 'credit', 'title': 'Order Payment Received', 'ref': 'ORD-10045', 'amount': 125000.0, 'time': '2 hours ago'},
      {'type': 'debit', 'title': 'Supplier Payment', 'ref': 'SUP-7823', 'amount': 450000.0, 'time': '5 hours ago'},
      {'type': 'credit', 'title': 'Settlement from Jumia', 'ref': 'SET-4521', 'amount': 890000.0, 'time': 'Yesterday'},
      {'type': 'debit', 'title': 'Delivery Fee', 'ref': 'DEL-1234', 'amount': 15000.0, 'time': 'Yesterday'},
      {'type': 'credit', 'title': 'Order Payment Received', 'ref': 'ORD-10044', 'amount': 78500.0, 'time': '2 days ago'},
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionHeader(title: 'Recent Transactions', actionText: 'View All', onAction: () {}),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: transactions.length,
            separatorBuilder: (_, __) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final tx = transactions[index];
              final isCredit = tx['type'] == 'credit';
              return ListTile(
                leading: Container(
                  width: 44,
                  height: 44,
                  decoration: BoxDecoration(
                    color: isCredit ? AppColors.successBg : AppColors.errorBg,
                    shape: BoxShape.circle,
                  ),
                  child: Icon(
                    isCredit ? Icons.arrow_downward : Icons.arrow_upward,
                    color: isCredit ? AppColors.success : AppColors.error,
                    size: 20,
                  ),
                ),
                title: Text(tx['title'] as String, style: AppTypography.titleSmall),
                subtitle: Text(
                  '${tx['ref']} • ${tx['time']}',
                  style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
                ),
                trailing: Text(
                  '${isCredit ? '+' : '-'}${formatCurrency(tx['amount'] as double)}',
                  style: AppTypography.titleSmall.copyWith(
                    color: isCredit ? AppColors.success : AppColors.error,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildPendingSettlements() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionHeader(title: 'Pending Settlements'),
        const SizedBox(height: 12),
        _SettlementCard(
          source: 'Jumia Marketplace',
          amount: 1250000,
          expectedDate: DateTime.now().add(const Duration(days: 2)),
          orders: 45,
        ),
        const SizedBox(height: 12),
        _SettlementCard(
          source: 'Konga Sales',
          amount: 680000,
          expectedDate: DateTime.now().add(const Duration(days: 3)),
          orders: 23,
        ),
      ],
    );
  }

  void _showAddMoneySheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => _AddMoneySheet(),
    );
  }

  void _showTransferSheet(BuildContext context) {
    // Implement transfer sheet
  }

  void _showWithdrawSheet(BuildContext context) {
    // Implement withdraw sheet
  }

  void _showLoanApplicationSheet(BuildContext context) {
    // Implement loan application sheet
  }
}

class _BalanceInfo extends StatelessWidget {
  final String label;
  final String value;
  final IconData icon;

  const _BalanceInfo({
    required this.label,
    required this.value,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Icon(icon, color: Colors.white.withValues(alpha: 0.7), size: 18),
        const SizedBox(width: 8),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              label,
              style: AppTypography.labelSmall.copyWith(
                color: Colors.white.withValues(alpha: 0.6),
              ),
            ),
            Text(
              value,
              style: AppTypography.titleSmall.copyWith(color: Colors.white),
            ),
          ],
        ),
      ],
    );
  }
}

class _WalletButton extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _WalletButton({
    required this.icon,
    required this.label,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: Colors.white.withValues(alpha: 0.15),
          borderRadius: AppRadius.borderRadiusSm,
        ),
        child: Column(
          children: [
            Icon(icon, color: Colors.white, size: 22),
            const SizedBox(height: 4),
            Text(
              label,
              style: AppTypography.labelSmall.copyWith(color: Colors.white),
            ),
          ],
        ),
      ),
    );
  }
}

class _ActionCard extends StatelessWidget {
  final IconData icon;
  final String label;
  final Color color;
  final VoidCallback onTap;

  const _ActionCard({
    required this.icon,
    required this.label,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: AppColors.white,
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: AppColors.cardBorder),
        ),
        child: Column(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(icon, color: color, size: 22),
            ),
            const SizedBox(height: 8),
            Text(
              label,
              style: AppTypography.labelSmall.copyWith(
                color: AppColors.grey800,
                fontWeight: FontWeight.w500,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }
}

class _CashFlowIndicator extends StatelessWidget {
  final String label;
  final String value;
  final Color color;
  final IconData icon;

  const _CashFlowIndicator({
    required this.label,
    required this.value,
    required this.color,
    required this.icon,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(
            color: color,
            shape: BoxShape.circle,
          ),
        ),
        const SizedBox(width: 8),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(label, style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
            Text(value, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
          ],
        ),
      ],
    );
  }
}

class _ServiceCard extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final Color color;
  final VoidCallback onTap;

  const _ServiceCard({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.color,
    required this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: color.withValues(alpha: 0.05),
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: color.withValues(alpha: 0.2)),
        ),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.15),
                borderRadius: AppRadius.borderRadiusSm,
              ),
              child: Icon(icon, color: color, size: 22),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: AppTypography.titleSmall),
                  Text(
                    subtitle,
                    style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
                  ),
                ],
              ),
            ),
            Icon(Icons.chevron_right, color: color),
          ],
        ),
      ),
    );
  }
}

class _SettlementCard extends StatelessWidget {
  final String source;
  final double amount;
  final DateTime expectedDate;
  final int orders;

  const _SettlementCard({
    required this.source,
    required this.amount,
    required this.expectedDate,
    required this.orders,
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
      child: Row(
        children: [
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: AppColors.primaryLight.withValues(alpha: 0.15),
              borderRadius: AppRadius.borderRadiusSm,
            ),
            child: const Icon(Icons.account_balance, color: AppColors.primary),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(source, style: AppTypography.titleSmall),
                Text(
                  '$orders orders • Expected ${formatDate(expectedDate)}',
                  style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
                ),
              ],
            ),
          ),
          Text(
            formatCurrency(amount),
            style: AppTypography.titleMedium.copyWith(
              color: AppColors.success,
              fontWeight: FontWeight.w700,
            ),
          ),
        ],
      ),
    );
  }
}

class _AddMoneySheet extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      height: MediaQuery.of(context).size.height * 0.6,
      decoration: const BoxDecoration(
        color: AppColors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
      ),
      child: Column(
        children: [
          Container(
            margin: const EdgeInsets.only(top: 12),
            width: 40,
            height: 4,
            decoration: BoxDecoration(
              color: AppColors.grey300,
              borderRadius: AppRadius.borderRadiusFull,
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(20),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Add Money', style: AppTypography.titleLarge),
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
          ),
          const Divider(height: 1),
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('Amount', style: AppTypography.labelMedium),
                  const SizedBox(height: 8),
                  TextField(
                    keyboardType: TextInputType.number,
                    style: AppTypography.headlineMedium,
                    decoration: InputDecoration(
                      hintText: '0.00',
                      prefixText: '₦ ',
                      border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Wrap(
                    spacing: 8,
                    children: [5000, 10000, 25000, 50000, 100000].map((amount) {
                      return ActionChip(
                        label: Text('₦${formatNumber(amount)}'),
                        onPressed: () {},
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 24),
                  Text('Payment Method', style: AppTypography.labelMedium),
                  const SizedBox(height: 12),
                  _PaymentMethodOption(
                    icon: Icons.credit_card,
                    title: 'Debit/Credit Card',
                    subtitle: 'Instant • No fee',
                    isSelected: true,
                  ),
                  const SizedBox(height: 8),
                  _PaymentMethodOption(
                    icon: Icons.account_balance,
                    title: 'Bank Transfer',
                    subtitle: '1-3 hours • No fee',
                    isSelected: false,
                  ),
                  const SizedBox(height: 8),
                  _PaymentMethodOption(
                    icon: Icons.qr_code,
                    title: 'USSD',
                    subtitle: 'Dial *737# or *901#',
                    isSelected: false,
                  ),
                ],
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(20),
            child: SizedBox(
              width: double.infinity,
              height: 56,
              child: ElevatedButton(
                onPressed: () => Navigator.pop(context),
                child: const Text('Continue'),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _PaymentMethodOption extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final bool isSelected;

  const _PaymentMethodOption({
    required this.icon,
    required this.title,
    required this.subtitle,
    required this.isSelected,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: isSelected ? AppColors.primary.withValues(alpha: 0.05) : AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(
          color: isSelected ? AppColors.primary : AppColors.grey300,
          width: isSelected ? 2 : 1,
        ),
      ),
      child: Row(
        children: [
          Icon(icon, color: isSelected ? AppColors.primary : AppColors.grey600),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: AppTypography.titleSmall),
                Text(
                  subtitle,
                  style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
                ),
              ],
            ),
          ),
          if (isSelected)
            const Icon(Icons.check_circle, color: AppColors.primary),
        ],
      ),
    );
  }
}
