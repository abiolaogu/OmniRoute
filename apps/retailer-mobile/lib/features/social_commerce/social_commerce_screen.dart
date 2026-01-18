/// OmniRoute Ecosystem - Social Commerce Module
/// Layer 5: Social Commerce - Group Buying, Referrals, Reputation
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';

// =============================================================================
// DOMAIN MODELS
// =============================================================================

/// Group Buy Campaign
class GroupBuyCampaign {
  final String id;
  final String productId;
  final String productName;
  final String productImage;
  final double originalPrice;
  final List<PriceTier> priceTiers;
  final int currentParticipants;
  final int maxParticipants;
  final DateTime endTime;
  final String organizerId;
  final CampaignStatus status;

  const GroupBuyCampaign({
    required this.id,
    required this.productId,
    required this.productName,
    required this.productImage,
    required this.originalPrice,
    required this.priceTiers,
    required this.currentParticipants,
    required this.maxParticipants,
    required this.endTime,
    required this.organizerId,
    required this.status,
  });

  double get currentPrice {
    for (final tier in priceTiers.reversed) {
      if (currentParticipants >= tier.minParticipants) {
        return tier.price;
      }
    }
    return originalPrice;
  }

  double get discount => ((originalPrice - currentPrice) / originalPrice * 100);
  int get participantsNeeded => priceTiers.last.minParticipants - currentParticipants;
  Duration get timeRemaining => endTime.difference(DateTime.now());
}

class PriceTier {
  final int minParticipants;
  final double price;
  final String label;

  const PriceTier({
    required this.minParticipants,
    required this.price,
    required this.label,
  });
}

enum CampaignStatus { active, successful, failed, cancelled }

/// Referral Program
class ReferralProgram {
  final String userId;
  final String referralCode;
  final int totalReferrals;
  final int successfulReferrals;
  final double totalEarnings;
  final double pendingEarnings;
  final ReferralTier currentTier;
  final List<Referral> recentReferrals;

  const ReferralProgram({
    required this.userId,
    required this.referralCode,
    required this.totalReferrals,
    required this.successfulReferrals,
    required this.totalEarnings,
    required this.pendingEarnings,
    required this.currentTier,
    required this.recentReferrals,
  });
}

class Referral {
  final String refereeId;
  final String refereeName;
  final DateTime joinDate;
  final double commission;
  final ReferralStatus status;

  const Referral({
    required this.refereeId,
    required this.refereeName,
    required this.joinDate,
    required this.commission,
    required this.status,
  });
}

enum ReferralStatus { pending, active, converted, churned }
enum ReferralTier { bronze, silver, gold, platinum, ambassador }

/// Reputation Passport
class ReputationPassport {
  final String userId;
  final int overallScore;
  final int transactionCount;
  final double paymentReliability;
  final double orderAccuracy;
  final double responseRate;
  final int reviewCount;
  final double averageRating;
  final List<Badge> badges;
  final List<Review> recentReviews;
  final DateTime memberSince;
  final VerificationLevel verification;

  const ReputationPassport({
    required this.userId,
    required this.overallScore,
    required this.transactionCount,
    required this.paymentReliability,
    required this.orderAccuracy,
    required this.responseRate,
    required this.reviewCount,
    required this.averageRating,
    required this.badges,
    required this.recentReviews,
    required this.memberSince,
    required this.verification,
  });

  String get trustLevel {
    if (overallScore >= 90) return 'Excellent';
    if (overallScore >= 75) return 'Good';
    if (overallScore >= 60) return 'Fair';
    return 'Building';
  }
}

class Badge {
  final String id;
  final String name;
  final String description;
  final String icon;
  final DateTime earnedAt;

  const Badge({
    required this.id,
    required this.name,
    required this.description,
    required this.icon,
    required this.earnedAt,
  });
}

class Review {
  final String reviewerId;
  final String reviewerName;
  final int rating;
  final String comment;
  final DateTime date;
  final String transactionType;

  const Review({
    required this.reviewerId,
    required this.reviewerName,
    required this.rating,
    required this.comment,
    required this.date,
    required this.transactionType,
  });
}

enum VerificationLevel { none, email, phone, id, business, premium }

// =============================================================================
// SOCIAL COMMERCE DASHBOARD
// =============================================================================

class SocialCommerceDashboardScreen extends ConsumerStatefulWidget {
  const SocialCommerceDashboardScreen({super.key});
  @override
  ConsumerState<SocialCommerceDashboardScreen> createState() => _SocialCommerceDashboardScreenState();
}

class _SocialCommerceDashboardScreenState extends ConsumerState<SocialCommerceDashboardScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('Social Commerce'),
        backgroundColor: AppColors.primary,
        foregroundColor: Colors.white,
        bottom: TabBar(
          controller: _tabController,
          indicatorColor: Colors.white,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white70,
          tabs: const [
            Tab(icon: Icon(Icons.groups), text: 'Group Buy'),
            Tab(icon: Icon(Icons.share), text: 'Referrals'),
            Tab(icon: Icon(Icons.verified), text: 'Reputation'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: const [
          _GroupBuyTab(),
          _ReferralTab(),
          _ReputationTab(),
        ],
      ),
    );
  }
}

class _GroupBuyTab extends StatelessWidget {
  const _GroupBuyTab();

  @override
  Widget build(BuildContext context) {
    final campaigns = [
      GroupBuyCampaign(
        id: '1',
        productId: 'prod-001',
        productName: 'Samsung Galaxy A54 (128GB)',
        productImage: 'assets/images/phone.png',
        originalPrice: 285000,
        priceTiers: const [
          PriceTier(minParticipants: 5, price: 270000, label: '5+ buyers'),
          PriceTier(minParticipants: 10, price: 255000, label: '10+ buyers'),
          PriceTier(minParticipants: 25, price: 240000, label: '25+ buyers'),
        ],
        currentParticipants: 18,
        maxParticipants: 50,
        endTime: DateTime.now().add(const Duration(hours: 23)),
        organizerId: 'user-001',
        status: CampaignStatus.active,
      ),
      GroupBuyCampaign(
        id: '2',
        productId: 'prod-002',
        productName: 'Indomie Chicken Flavour (Carton of 40)',
        productImage: 'assets/images/indomie.png',
        originalPrice: 12000,
        priceTiers: const [
          PriceTier(minParticipants: 10, price: 11000, label: '10+ buyers'),
          PriceTier(minParticipants: 50, price: 10000, label: '50+ buyers'),
          PriceTier(minParticipants: 100, price: 9000, label: '100+ buyers'),
        ],
        currentParticipants: 67,
        maxParticipants: 200,
        endTime: DateTime.now().add(const Duration(hours: 48)),
        organizerId: 'user-002',
        status: CampaignStatus.active,
      ),
    ];

    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildCreateGroupBuyCard(context),
          const SizedBox(height: 20),
          Text('Active Campaigns', style: AppTypography.titleMedium),
          const SizedBox(height: 12),
          ...campaigns.map((c) => _buildCampaignCard(c)),
        ],
      ),
    );
  }

  Widget _buildCreateGroupBuyCard(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(colors: [AppColors.primary, Color(0xFF1565C0)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Start a Group Buy', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
                const SizedBox(height: 4),
                Text('Invite friends to buy together and unlock bigger discounts!',
                    style: AppTypography.bodySmall.copyWith(color: Colors.white70)),
              ],
            ),
          ),
          ElevatedButton(
            onPressed: () {},
            style: ElevatedButton.styleFrom(backgroundColor: Colors.white, foregroundColor: AppColors.primary),
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }

  Widget _buildCampaignCard(GroupBuyCampaign campaign) {
    final hours = campaign.timeRemaining.inHours;
    final minutes = campaign.timeRemaining.inMinutes % 60;

    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                width: 60, height: 60,
                decoration: BoxDecoration(color: AppColors.grey100, borderRadius: AppRadius.borderRadiusSm),
                child: const Icon(Icons.shopping_bag, color: AppColors.grey500),
              ),
              const SizedBox(width: 12),
              Expanded(child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(campaign.productName, style: AppTypography.titleSmall, maxLines: 2),
                  const SizedBox(height: 4),
                  Row(children: [
                    Text('₦${campaign.currentPrice.toStringAsFixed(0)}', style: AppTypography.titleMedium.copyWith(color: AppColors.primary, fontWeight: FontWeight.bold)),
                    const SizedBox(width: 8),
                    Text('₦${campaign.originalPrice.toStringAsFixed(0)}', style: AppTypography.labelSmall.copyWith(decoration: TextDecoration.lineThrough, color: AppColors.grey500)),
                    const SizedBox(width: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(color: AppColors.success, borderRadius: AppRadius.borderRadiusSm),
                      child: Text('-${campaign.discount.toStringAsFixed(0)}%', style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold)),
                    ),
                  ]),
                ],
              )),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text('${campaign.currentParticipants} joined', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
                  const SizedBox(height: 4),
                  LinearProgressIndicator(
                    value: campaign.currentParticipants / campaign.maxParticipants,
                    backgroundColor: AppColors.grey200,
                    valueColor: const AlwaysStoppedAnimation(AppColors.success),
                  ),
                ],
              )),
              const SizedBox(width: 16),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Row(children: [
                    const Icon(Icons.timer, size: 14, color: AppColors.warning),
                    const SizedBox(width: 4),
                    Text('${hours}h ${minutes}m left', style: AppTypography.labelSmall.copyWith(color: AppColors.warning)),
                  ]),
                  const SizedBox(height: 4),
                  Text('${campaign.participantsNeeded} more for next tier', style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
                ],
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(child: OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.share, size: 16), label: const Text('Share'))),
              const SizedBox(width: 12),
              Expanded(child: ElevatedButton(onPressed: () {}, child: const Text('Join Now'))),
            ],
          ),
        ],
      ),
    );
  }
}

class _ReferralTab extends StatelessWidget {
  const _ReferralTab();

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          _buildReferralCodeCard(),
          const SizedBox(height: 20),
          _buildEarningsCard(),
          const SizedBox(height: 20),
          _buildTierProgress(),
          const SizedBox(height: 20),
          _buildRecentReferrals(),
        ],
      ),
    );
  }

  Widget _buildReferralCodeCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(colors: [Color(0xFF7B1FA2), Color(0xFF9C27B0)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Column(
        children: [
          Text('Your Referral Code', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
          const SizedBox(height: 8),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
            decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusSm),
            child: Text('OMNI-ABC123', style: AppTypography.headlineSmall.copyWith(color: Colors.white, fontWeight: FontWeight.bold, letterSpacing: 2)),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextButton.icon(onPressed: () {}, icon: const Icon(Icons.copy, color: Colors.white, size: 16), label: Text('Copy', style: AppTypography.labelMedium.copyWith(color: Colors.white))),
              const SizedBox(width: 16),
              TextButton.icon(onPressed: () {}, icon: const Icon(Icons.share, color: Colors.white, size: 16), label: Text('Share', style: AppTypography.labelMedium.copyWith(color: Colors.white))),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildEarningsCard() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Row(
        children: [
          Expanded(child: _buildEarningStat('Total Referrals', '47', AppColors.primary)),
          Container(width: 1, height: 40, color: AppColors.grey200),
          Expanded(child: _buildEarningStat('Active Users', '32', AppColors.success)),
          Container(width: 1, height: 40, color: AppColors.grey200),
          Expanded(child: _buildEarningStat('Earnings', '₦127K', AppColors.warning)),
        ],
      ),
    );
  }

  Widget _buildEarningStat(String label, String value, Color color) {
    return Column(
      children: [
        Text(value, style: AppTypography.titleLarge.copyWith(color: color, fontWeight: FontWeight.bold)),
        Text(label, style: AppTypography.labelSmall),
      ],
    );
  }

  Widget _buildTierProgress() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusMd),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('Your Tier: Gold', style: AppTypography.titleMedium),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: Colors.amber, borderRadius: AppRadius.borderRadiusSm),
                child: const Text('5% commission', style: TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.bold)),
              ),
            ],
          ),
          const SizedBox(height: 16),
          LinearProgressIndicator(value: 0.7, backgroundColor: AppColors.grey200, valueColor: const AlwaysStoppedAnimation(Colors.amber)),
          const SizedBox(height: 8),
          Text('18 more referrals to reach Platinum (7% commission)', style: AppTypography.labelSmall.copyWith(color: AppColors.grey600)),
        ],
      ),
    );
  }

  Widget _buildRecentReferrals() {
    final referrals = [
      const Referral(refereeId: '1', refereeName: 'John Okafor', joinDate: null, commission: 2500, status: ReferralStatus.converted),
      const Referral(refereeId: '2', refereeName: 'Mary Adebayo', joinDate: null, commission: 1500, status: ReferralStatus.active),
      const Referral(refereeId: '3', refereeName: 'Chidi Nnamdi', joinDate: null, commission: 0, status: ReferralStatus.pending),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Recent Referrals', style: AppTypography.titleMedium),
        const SizedBox(height: 12),
        ...referrals.map((r) => Container(
          margin: const EdgeInsets.only(bottom: 8),
          padding: const EdgeInsets.all(12),
          decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
          child: Row(
            children: [
              CircleAvatar(backgroundColor: AppColors.grey200, child: Text(r.refereeName[0])),
              const SizedBox(width: 12),
              Expanded(child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(r.refereeName, style: AppTypography.titleSmall),
                  Text(_getStatusText(r.status), style: AppTypography.labelSmall.copyWith(color: _getStatusColor(r.status))),
                ],
              )),
              if (r.commission > 0) Text('+₦${r.commission.toStringAsFixed(0)}', style: AppTypography.titleSmall.copyWith(color: AppColors.success, fontWeight: FontWeight.w600)),
            ],
          ),
        )),
      ],
    );
  }

  String _getStatusText(ReferralStatus status) {
    switch (status) {
      case ReferralStatus.pending: return 'Pending signup';
      case ReferralStatus.active: return 'Active user';
      case ReferralStatus.converted: return 'Converted - earned';
      case ReferralStatus.churned: return 'Inactive';
    }
  }

  Color _getStatusColor(ReferralStatus status) {
    switch (status) {
      case ReferralStatus.pending: return AppColors.warning;
      case ReferralStatus.active: return AppColors.info;
      case ReferralStatus.converted: return AppColors.success;
      case ReferralStatus.churned: return AppColors.error;
    }
  }
}

class _ReputationTab extends StatelessWidget {
  const _ReputationTab();

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          _buildReputationScoreCard(),
          const SizedBox(height: 20),
          _buildMetricsGrid(),
          const SizedBox(height: 20),
          _buildBadges(),
          const SizedBox(height: 20),
          _buildRecentReviews(),
        ],
      ),
    );
  }

  Widget _buildReputationScoreCard() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: const LinearGradient(colors: [AppColors.success, Color(0xFF43A047)]),
        borderRadius: AppRadius.borderRadiusMd,
      ),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.verified, color: Colors.white, size: 24),
              const SizedBox(width: 8),
              Text('Reputation Passport', style: AppTypography.titleMedium.copyWith(color: Colors.white)),
            ],
          ),
          const SizedBox(height: 16),
          Text('87', style: AppTypography.displayLarge.copyWith(color: Colors.white, fontWeight: FontWeight.bold)),
          Text('Excellent Trust Score', style: AppTypography.labelMedium.copyWith(color: Colors.white70)),
          const SizedBox(height: 16),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: AppRadius.borderRadiusFull),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.check_circle, color: Colors.white, size: 16),
                const SizedBox(width: 4),
                Text('Business Verified', style: AppTypography.labelSmall.copyWith(color: Colors.white)),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildMetricsGrid() {
    return GridView.count(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisCount: 2,
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.5,
      children: [
        _buildMetricCard('Payment Reliability', '98%', Icons.payments, AppColors.success),
        _buildMetricCard('Order Accuracy', '94%', Icons.check_circle, AppColors.primary),
        _buildMetricCard('Response Rate', '92%', Icons.chat, AppColors.info),
        _buildMetricCard('Avg Rating', '4.7★', Icons.star, AppColors.warning),
      ],
    );
  }

  Widget _buildMetricCard(String label, String value, IconData icon, Color color) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(height: 8),
          Text(value, style: AppTypography.titleMedium.copyWith(fontWeight: FontWeight.bold)),
          Text(label, style: AppTypography.labelSmall, textAlign: TextAlign.center),
        ],
      ),
    );
  }

  Widget _buildBadges() {
    final badges = [
      const Badge(id: '1', name: 'Top Buyer', description: '1000+ orders', icon: '🏆', earnedAt: null),
      const Badge(id: '2', name: 'Quick Payer', description: 'Always pays on time', icon: '⚡', earnedAt: null),
      const Badge(id: '3', name: 'Trusted', description: '2+ years member', icon: '✓', earnedAt: null),
      const Badge(id: '4', name: 'Community', description: 'Active reviewer', icon: '💬', earnedAt: null),
    ];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Badges Earned', style: AppTypography.titleMedium),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8,
          runSpacing: 8,
          children: badges.map((b) => Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
            decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusFull, border: Border.all(color: AppColors.grey200)),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(b.icon, style: const TextStyle(fontSize: 16)),
                const SizedBox(width: 6),
                Text(b.name, style: AppTypography.labelMedium),
              ],
            ),
          )).toList(),
        ),
      ],
    );
  }

  Widget _buildRecentReviews() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('Recent Reviews', style: AppTypography.titleMedium),
        const SizedBox(height: 12),
        _buildReviewCard('Mama Ngozi Stores', 5, 'Excellent business partner. Always reliable!', '2 days ago'),
        _buildReviewCard('Chidi Wholesale', 4, 'Good dealings, prompt response', '1 week ago'),
      ],
    );
  }

  Widget _buildReviewCard(String reviewer, int rating, String comment, String time) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(color: Colors.white, borderRadius: AppRadius.borderRadiusSm),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(reviewer, style: AppTypography.titleSmall),
              Row(children: List.generate(5, (i) => Icon(i < rating ? Icons.star : Icons.star_border, size: 14, color: Colors.amber))),
            ],
          ),
          const SizedBox(height: 8),
          Text(comment, style: AppTypography.bodySmall),
          const SizedBox(height: 4),
          Text(time, style: AppTypography.labelSmall.copyWith(color: AppColors.grey500)),
        ],
      ),
    );
  }
}
