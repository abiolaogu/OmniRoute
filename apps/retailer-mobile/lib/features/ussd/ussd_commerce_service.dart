/// OmniRoute Ecosystem - USSD Commerce Gateway
/// Layer 4: Accessibility - USSD Interface for Feature Phone Users
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';

/// USSD Menu State Management
class USSDSession {
  final String sessionId;
  final String msisdn;
  final List<String> menuStack;
  final Map<String, dynamic> sessionData;
  final DateTime startTime;
  
  const USSDSession({
    required this.sessionId,
    required this.msisdn,
    this.menuStack = const ['main'],
    this.sessionData = const {},
    required this.startTime,
  });
  
  USSDSession copyWith({
    List<String>? menuStack,
    Map<String, dynamic>? sessionData,
  }) => USSDSession(
    sessionId: sessionId,
    msisdn: msisdn,
    menuStack: menuStack ?? this.menuStack,
    sessionData: {...this.sessionData, ...?sessionData},
    startTime: startTime,
  );
}

/// USSD Commerce Service
class USSDCommerceService {
  // Main Menu
  static const String mainMenu = '''
Welcome to OmniRoute
1. Place Order
2. Check Order Status
3. View Balance
4. Make Payment
5. My Account
0. Exit
''';

  // Order Menu
  static const String orderMenu = '''
Place Order
1. Quick Reorder (Last Order)
2. Browse Categories
3. Enter Product Code
4. View Cart
0. Back
''';

  // Categories
  static const String categoriesMenu = '''
Categories
1. Beverages
2. Food Items
3. Personal Care
4. Household
5. Electronics
0. Back
''';

  // Process USSD input
  String processInput(USSDSession session, String input) {
    final currentMenu = session.menuStack.last;
    
    switch (currentMenu) {
      case 'main':
        return _handleMainMenu(input);
      case 'order':
        return _handleOrderMenu(input);
      case 'categories':
        return _handleCategoriesMenu(input);
      case 'balance':
        return _handleBalanceCheck(session);
      case 'payment':
        return _handlePayment(session, input);
      default:
        return mainMenu;
    }
  }
  
  String _handleMainMenu(String input) {
    switch (input) {
      case '1': return orderMenu;
      case '2': return 'Enter Order ID:\n';
      case '3': return 'Your Balance:\nWallet: NGN 45,230\nCredit: NGN 150,000\n\n0. Back';
      case '4': return 'Payment Amount:\nEnter amount in Naira\n';
      case '5': return 'My Account\n1. Profile\n2. Credit Limit\n3. Transaction History\n0. Back';
      case '0': return 'END Thank you for using OmniRoute';
      default: return 'Invalid option. $mainMenu';
    }
  }
  
  String _handleOrderMenu(String input) {
    switch (input) {
      case '1': return 'Reorder Last:\n12x Peak Milk 400g\n5x Indomie Carton\nTotal: NGN 85,000\n\n1. Confirm\n0. Cancel';
      case '2': return categoriesMenu;
      case '3': return 'Enter Product Code:\n';
      case '4': return 'Your Cart:\n[Empty]\n0. Back';
      case '0': return mainMenu;
      default: return 'Invalid. $orderMenu';
    }
  }
  
  String _handleCategoriesMenu(String input) {
    final categories = {
      '1': 'Beverages:\n1. Soft Drinks\n2. Water\n3. Juices\n4. Alcohol\n0. Back',
      '2': 'Food Items:\n1. Rice/Grains\n2. Cooking Oil\n3. Noodles\n4. Canned\n0. Back',
      '3': 'Personal Care:\n1. Soap\n2. Toiletries\n3. Cosmetics\n0. Back',
      '4': 'Household:\n1. Cleaning\n2. Kitchen\n3. Electrical\n0. Back',
      '5': 'Electronics:\n1. Phones\n2. Accessories\n3. Appliances\n0. Back',
    };
    return categories[input] ?? 'Invalid. $categoriesMenu';
  }
  
  String _handleBalanceCheck(USSDSession session) {
    return '''
Account: ${session.msisdn}
Wallet Balance: NGN 45,230.00
Available Credit: NGN 150,000.00
Credit Used: NGN 35,000.00
Next Payment: Jan 25, 2026

0. Back
''';
  }
  
  String _handlePayment(USSDSession session, String input) {
    if (RegExp(r'^\d+$').hasMatch(input)) {
      final amount = int.tryParse(input) ?? 0;
      if (amount > 0) {
        return '''
Confirm Payment:
Amount: NGN ${amount.toString().replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (m) => '${m[1]},')}
From: Wallet

1. Confirm with PIN
0. Cancel
''';
      }
    }
    return 'Enter valid amount:\n';
  }
}

/// USSD Gateway Dashboard Widget (for Admin monitoring)
class USSDGatewayDashboard extends ConsumerWidget {
  const USSDGatewayDashboard({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('USSD Gateway Dashboard'),
        backgroundColor: AppColors.primary,
        foregroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildOverviewCard(),
            const SizedBox(height: 16),
            _buildActiveSessionsCard(),
            const SizedBox(height: 16),
            _buildMenuAnalytics(),
            const SizedBox(height: 16),
            _buildRecentTransactions(),
          ],
        ),
      ),
    );
  }

  Widget _buildOverviewCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('USSD Overview', style: AppTypography.titleMedium),
            const SizedBox(height: 16),
            Row(
              children: [
                _buildStatTile('Active Sessions', '1,247', Icons.phone_android, AppColors.success),
                _buildStatTile('Today Orders', '3,891', Icons.shopping_cart, AppColors.primary),
                _buildStatTile('Success Rate', '94.2%', Icons.check_circle, AppColors.success),
                _buildStatTile('Avg Duration', '45s', Icons.timer, AppColors.info),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatTile(String label, String value, IconData icon, Color color) {
    return Expanded(
      child: Column(
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(height: 8),
          Text(value, style: AppTypography.titleLarge.copyWith(fontWeight: FontWeight.bold)),
          Text(label, style: AppTypography.labelSmall),
        ],
      ),
    );
  }

  Widget _buildActiveSessionsCard() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Active Sessions by Carrier', style: AppTypography.titleMedium),
            const SizedBox(height: 16),
            _buildCarrierRow('MTN', 542, 0.43),
            _buildCarrierRow('Airtel', 312, 0.25),
            _buildCarrierRow('Glo', 245, 0.20),
            _buildCarrierRow('9mobile', 148, 0.12),
          ],
        ),
      ),
    );
  }

  Widget _buildCarrierRow(String carrier, int sessions, double percentage) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          SizedBox(width: 80, child: Text(carrier, style: AppTypography.titleSmall)),
          Expanded(
            child: LinearProgressIndicator(
              value: percentage,
              backgroundColor: AppColors.grey200,
              valueColor: const AlwaysStoppedAnimation(AppColors.primary),
            ),
          ),
          const SizedBox(width: 16),
          Text('$sessions', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }

  Widget _buildMenuAnalytics() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Menu Flow Analytics', style: AppTypography.titleMedium),
            const SizedBox(height: 16),
            _buildFlowRow('Main → Order', '45%'),
            _buildFlowRow('Order → Quick Reorder', '62%'),
            _buildFlowRow('Order → Categories', '28%'),
            _buildFlowRow('Main → Balance', '22%'),
            _buildFlowRow('Main → Payment', '18%'),
          ],
        ),
      ),
    );
  }

  Widget _buildFlowRow(String flow, String percentage) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(flow, style: AppTypography.bodyMedium),
          Text(percentage, style: AppTypography.labelMedium.copyWith(color: AppColors.primary, fontWeight: FontWeight.bold)),
        ],
      ),
    );
  }

  Widget _buildRecentTransactions() {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Recent USSD Transactions', style: AppTypography.titleMedium),
            const SizedBox(height: 16),
            _buildTransactionRow('+234 803 XXX 1234', 'Order Placed', 'NGN 45,000', '2 min ago'),
            _buildTransactionRow('+234 805 XXX 5678', 'Payment', 'NGN 12,500', '5 min ago'),
            _buildTransactionRow('+234 701 XXX 9012', 'Balance Check', '-', '7 min ago'),
            _buildTransactionRow('+234 809 XXX 3456', 'Order Placed', 'NGN 89,200', '12 min ago'),
          ],
        ),
      ),
    );
  }

  Widget _buildTransactionRow(String phone, String action, String amount, String time) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        children: [
          const Icon(Icons.phone_android, size: 20, color: AppColors.grey500),
          const SizedBox(width: 8),
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(phone, style: AppTypography.titleSmall),
              Text(action, style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
            ],
          )),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(amount, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.bold)),
              Text(time, style: AppTypography.labelSmall.copyWith(color: AppColors.grey500)),
            ],
          ),
        ],
      ),
    );
  }
}
