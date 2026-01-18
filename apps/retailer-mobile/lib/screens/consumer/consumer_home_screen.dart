/// OmniRoute Ecosystem - Consumer Home Screen
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

class ConsumerHomeScreen extends ConsumerStatefulWidget {
  const ConsumerHomeScreen({super.key});
  @override ConsumerState<ConsumerHomeScreen> createState() => _ConsumerHomeScreenState();
}

class _ConsumerHomeScreenState extends ConsumerState<ConsumerHomeScreen> {
  final _searchController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      body: CustomScrollView(
        slivers: [
          _buildAppBar(),
          SliverPadding(
            padding: const EdgeInsets.all(16),
            sliver: SliverList(delegate: SliverChildListDelegate([
              _buildSearchBar(),
              const SizedBox(height: 20),
              _buildCategories(),
              const SizedBox(height: 20),
              _buildPromoBanner(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Popular Near You'),
              const SizedBox(height: 12),
              _buildPopularStores(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Best Deals'),
              const SizedBox(height: 12),
              _buildBestDeals(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Recently Ordered'),
              const SizedBox(height: 12),
              _buildRecentlyOrdered(),
              const SizedBox(height: 20),
              const SectionHeader(title: 'Grocery Essentials'),
              const SizedBox(height: 12),
              _buildEssentials(),
            ])),
          ),
        ],
      ),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildAppBar() {
    return SliverAppBar(
      floating: true, backgroundColor: AppColors.consumerColor, foregroundColor: Colors.white,
      title: Row(children: [
        const Icon(Icons.location_on, size: 20),
        const SizedBox(width: 8),
        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text('Deliver to', style: AppTypography.labelSmall.copyWith(color: Colors.white70)),
          Row(children: [
            Text('Victoria Island, Lagos', style: AppTypography.titleSmall.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
            const Icon(Icons.keyboard_arrow_down, color: Colors.white, size: 20),
          ]),
        ])),
      ]),
      actions: [
        IconButton(icon: const Icon(Icons.favorite_border), onPressed: () {}),
        Stack(children: [
          IconButton(icon: const Icon(Icons.shopping_cart_outlined), onPressed: () {}),
          Positioned(right: 8, top: 8, child: Container(
            padding: const EdgeInsets.all(4),
            decoration: const BoxDecoration(color: Colors.red, shape: BoxShape.circle),
            child: const Text('3', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold)),
          )),
        ]),
      ],
    );
  }

  Widget _buildSearchBar() {
    return Container(
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
      child: TextField(
        controller: _searchController,
        decoration: InputDecoration(
          hintText: 'Search for products, stores...',
          prefixIcon: const Icon(Icons.search, color: AppColors.textSecondary),
          suffixIcon: IconButton(icon: const Icon(Icons.mic, color: AppColors.consumerColor), onPressed: () {}),
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        ),
      ),
    );
  }

  Widget _buildCategories() {
    final categories = [
      {'name': 'Groceries', 'icon': Icons.shopping_basket, 'color': AppColors.success},
      {'name': 'Drinks', 'icon': Icons.local_drink, 'color': AppColors.info},
      {'name': 'Pharmacy', 'icon': Icons.medical_services, 'color': AppColors.error},
      {'name': 'Fresh', 'icon': Icons.eco, 'color': AppColors.warning},
      {'name': 'Electronics', 'icon': Icons.devices, 'color': AppColors.primary},
      {'name': 'More', 'icon': Icons.more_horiz, 'color': AppColors.textSecondary},
    ];
    return SizedBox(
      height: 90,
      child: ListView.separated(
        scrollDirection: Axis.horizontal, itemCount: categories.length, separatorBuilder: (_, __) => const SizedBox(width: 16),
        itemBuilder: (context, index) {
          final c = categories[index];
          return Column(mainAxisAlignment: MainAxisAlignment.center, children: [
            Container(
              width: 56, height: 56,
              decoration: BoxDecoration(color: (c['color'] as Color).withValues(alpha: 0.1), shape: BoxShape.circle),
              child: Icon(c['icon'] as IconData, color: c['color'] as Color, size: 28),
            ),
            const SizedBox(height: 8),
            Text(c['name']!.toString(), style: AppTypography.labelSmall),
          ]);
        },
      ),
    );
  }

  Widget _buildPromoBanner() {
    return Container(
      height: 140,
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: [AppColors.consumerColor, AppColors.consumerColor.withValues(alpha: 0.8)]),
        borderRadius: AppRadius.borderRadiusLg,
      ),
      child: Stack(children: [
        Positioned(right: -20, bottom: -20, child: Icon(Icons.local_offer, size: 140, color: Colors.white.withValues(alpha: 0.1))),
        Padding(
          padding: const EdgeInsets.all(20),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, mainAxisAlignment: MainAxisAlignment.center, children: [
            Container(padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4), decoration: BoxDecoration(color: Colors.amber, borderRadius: AppRadius.borderRadiusSm),
              child: const Text('NEW USER', style: TextStyle(fontSize: 10, fontWeight: FontWeight.bold))),
            const SizedBox(height: 8),
            Text('Get 20% OFF', style: AppTypography.headlineMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w700)),
            Text('on your first order!', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
            const SizedBox(height: 8),
            Text('Use code: WELCOME20', style: AppTypography.labelMedium.copyWith(color: Colors.white, fontWeight: FontWeight.w600)),
          ]),
        ),
      ]),
    );
  }

  Widget _buildPopularStores() {
    final stores = [
      {'name': 'Shoprite', 'delivery': '15-25 min', 'rating': '4.8', 'offer': '10% off'},
      {'name': 'SPAR', 'delivery': '20-30 min', 'rating': '4.7', 'offer': 'Free delivery'},
      {'name': 'Justrite', 'delivery': '25-35 min', 'rating': '4.6', 'offer': null},
    ];
    return SizedBox(
      height: 180,
      child: ListView.separated(
        scrollDirection: Axis.horizontal, itemCount: stores.length, separatorBuilder: (_, __) => const SizedBox(width: 16),
        itemBuilder: (context, index) {
          final s = stores[index];
          return Container(
            width: 160,
            decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 10)]),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Container(
                height: 80,
                decoration: BoxDecoration(color: AppColors.consumerColor.withValues(alpha: 0.1), borderRadius: const BorderRadius.vertical(top: Radius.circular(12))),
                child: Center(child: Icon(Icons.store, size: 40, color: AppColors.consumerColor)),
              ),
              Padding(padding: const EdgeInsets.all(12), child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(s['name']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 4),
                Row(children: [
                  const Icon(Icons.star, size: 14, color: Colors.amber),
                  Text(' ${s['rating']}', style: AppTypography.labelSmall),
                  Text(' • ${s['delivery']}', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
                ]),
                if (s['offer'] != null) ...[
                  const SizedBox(height: 4),
                  Text(s['offer']!, style: AppTypography.labelSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
                ],
              ])),
            ]),
          );
        },
      ),
    );
  }

  Widget _buildBestDeals() {
    final deals = [
      {'name': 'Golden Penny Spaghetti', 'price': '₦850', 'oldPrice': '₦1,200', 'discount': '29%'},
      {'name': 'Indomie Chicken (40pcs)', 'price': '₦3,500', 'oldPrice': '₦4,500', 'discount': '22%'},
      {'name': 'Peak Milk 400g', 'price': '₦2,100', 'oldPrice': '₦2,500', 'discount': '16%'},
    ];
    return SizedBox(
      height: 200,
      child: ListView.separated(
        scrollDirection: Axis.horizontal, itemCount: deals.length, separatorBuilder: (_, __) => const SizedBox(width: 12),
        itemBuilder: (context, index) {
          final d = deals[index];
          return Container(
            width: 140,
            decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, border: Border.all(color: AppColors.borderColor)),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Stack(children: [
                Container(height: 100, decoration: BoxDecoration(color: AppColors.borderColor, borderRadius: const BorderRadius.vertical(top: Radius.circular(12))),
                  child: const Center(child: Icon(Icons.image, size: 40, color: AppColors.textSecondary))),
                Positioned(top: 8, left: 8, child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                  decoration: BoxDecoration(color: AppColors.error, borderRadius: AppRadius.borderRadiusSm),
                  child: Text('-${d['discount']}', style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold)),
                )),
              ]),
              Padding(padding: const EdgeInsets.all(12), child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(d['name']!, style: AppTypography.labelMedium, maxLines: 2, overflow: TextOverflow.ellipsis),
                const SizedBox(height: 4),
                Row(children: [
                  Text(d['price']!, style: AppTypography.titleSmall.copyWith(fontWeight: FontWeight.w700, color: AppColors.consumerColor)),
                  const SizedBox(width: 4),
                  Text(d['oldPrice']!, style: AppTypography.labelSmall.copyWith(decoration: TextDecoration.lineThrough, color: AppColors.textSecondary)),
                ]),
              ])),
            ]),
          );
        },
      ),
    );
  }

  Widget _buildRecentlyOrdered() {
    final items = ['Rice 5kg', 'Cooking Oil 2L', 'Sugar 1kg', 'Milk Powder'];
    return SizedBox(
      height: 100,
      child: ListView.separated(
        scrollDirection: Axis.horizontal, itemCount: items.length, separatorBuilder: (_, __) => const SizedBox(width: 12),
        itemBuilder: (context, index) => Container(
          width: 80,
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, border: Border.all(color: AppColors.borderColor)),
          child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
            Container(width: 40, height: 40, decoration: BoxDecoration(color: AppColors.consumerColor.withValues(alpha: 0.1), shape: BoxShape.circle),
              child: const Icon(Icons.shopping_bag, color: AppColors.consumerColor, size: 20)),
            const SizedBox(height: 8),
            Text(items[index], style: AppTypography.labelSmall, textAlign: TextAlign.center, maxLines: 1, overflow: TextOverflow.ellipsis),
          ]),
        ),
      ),
    );
  }

  Widget _buildEssentials() {
    return GridView.count(
      shrinkWrap: true, physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 3, crossAxisSpacing: 12, mainAxisSpacing: 12, childAspectRatio: 0.8,
      children: ['Rice', 'Oil', 'Sugar', 'Salt', 'Flour', 'Beans'].map((name) => Container(
        decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 5)]),
        child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
          Container(width: 48, height: 48, decoration: BoxDecoration(color: AppColors.consumerColor.withValues(alpha: 0.1), shape: BoxShape.circle),
            child: const Icon(Icons.shopping_bag, color: AppColors.consumerColor)),
          const SizedBox(height: 8),
          Text(name, style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
          Text('From ₦500', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary)),
        ]),
      )).toList(),
    );
  }

  Widget _buildBottomNav() {
    return BottomNavigationBar(
      currentIndex: 0, type: BottomNavigationBarType.fixed, selectedItemColor: AppColors.consumerColor,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.home), label: 'Home'),
        BottomNavigationBarItem(icon: Icon(Icons.search), label: 'Browse'),
        BottomNavigationBarItem(icon: Icon(Icons.local_offer), label: 'Deals'),
        BottomNavigationBarItem(icon: Icon(Icons.receipt), label: 'Orders'),
        BottomNavigationBarItem(icon: Icon(Icons.account_circle), label: 'Account'),
      ],
    );
  }
}
