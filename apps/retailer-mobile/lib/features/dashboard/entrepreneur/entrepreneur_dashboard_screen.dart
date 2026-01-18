/// OmniRoute Ecosystem - Entrepreneur Dashboard
/// New business dashboard with learning resources, mentorship, and growth tools

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/widgets/common_widgets.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class EntrepreneurDashboardScreen extends ConsumerStatefulWidget {
  const EntrepreneurDashboardScreen({super.key});

  @override
  ConsumerState<EntrepreneurDashboardScreen> createState() =>
      _EntrepreneurDashboardScreenState();
}

class _EntrepreneurDashboardScreenState
    extends ConsumerState<EntrepreneurDashboardScreen> {
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
              'Hello,',
              style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
            ),
            Text(
              user?.fullName ?? 'Entrepreneur',
              style: AppTypography.titleMedium.copyWith(color: AppColors.grey900),
            ),
          ],
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.school_outlined, color: AppColors.grey800),
            onPressed: () {},
          ),
          IconButton(
            icon: Stack(
              children: [
                const Icon(Icons.notifications_outlined, color: AppColors.grey800),
                Positioned(
                  right: 0,
                  top: 0,
                  child: Container(
                    width: 8,
                    height: 8,
                    decoration: const BoxDecoration(
                      color: AppColors.error,
                      shape: BoxShape.circle,
                    ),
                  ),
                ),
              ],
            ),
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
              // Welcome Card with Progress
              _buildWelcomeCard()
                  .animate()
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Business Setup Progress
              SectionHeader(
                title: 'Business Setup',
                actionText: 'Continue',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildSetupProgress()
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
                    icon: Icons.store,
                    label: 'My Store',
                    color: AppColors.entrepreneurColor,
                    onTap: () {},
                  ),
                  QuickActionItem(
                    icon: Icons.add_shopping_cart,
                    label: 'Add Product',
                    color: AppColors.success,
                    onTap: () {},
                  ),
                  QuickActionItem(
                    icon: Icons.school,
                    label: 'Learn',
                    color: AppColors.info,
                    onTap: () {},
                    badge: '3',
                  ),
                  QuickActionItem(
                    icon: Icons.people,
                    label: 'Mentors',
                    color: AppColors.warning,
                    onTap: () {},
                  ),
                ],
              )
                  .animate(delay: const Duration(milliseconds: 200))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Key Metrics
              _buildKeyMetrics()
                  .animate(delay: const Duration(milliseconds: 300))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Learning Path
              SectionHeader(
                title: 'Continue Learning',
                actionText: 'View All',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildLearningPath()
                  .animate(delay: const Duration(milliseconds: 400))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Funding Opportunities
              SectionHeader(
                title: 'Funding Opportunities',
                actionText: 'Explore',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildFundingOpportunities()
                  .animate(delay: const Duration(milliseconds: 500))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Mentorship
              SectionHeader(
                title: 'Connect with Mentors',
                actionText: 'Find More',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildMentors()
                  .animate(delay: const Duration(milliseconds: 600))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Community
              SectionHeader(
                title: 'Community Feed',
                actionText: 'View All',
                onAction: () {},
              ),
              const SizedBox(height: 12),
              _buildCommunityFeed()
                  .animate(delay: const Duration(milliseconds: 700))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildWelcomeCard() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFFF9A825), Color(0xFFFFB300)],
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
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '🚀 Your Journey Begins!',
                      style: AppTypography.titleLarge.copyWith(
                        color: Colors.white,
                        fontWeight: FontWeight.w700,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      'Complete your business setup to unlock all features and start selling.',
                      style: AppTypography.bodyMedium.copyWith(
                        color: Colors.white.withValues(alpha: 0.9),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 16),
              Container(
                width: 60,
                height: 60,
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.2),
                  shape: BoxShape.circle,
                ),
                child: Center(
                  child: Text(
                    '45%',
                    style: AppTypography.titleMedium.copyWith(
                      color: Colors.white,
                      fontWeight: FontWeight.w700,
                    ),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          LinearProgressIndicator(
            value: 0.45,
            backgroundColor: Colors.white.withValues(alpha: 0.3),
            valueColor: const AlwaysStoppedAnimation<Color>(Colors.white),
            minHeight: 8,
            borderRadius: BorderRadius.circular(4),
          ),
        ],
      ),
    );
  }

  Widget _buildSetupProgress() {
    final steps = [
      {'title': 'Create Account', 'done': true},
      {'title': 'Verify Identity', 'done': true},
      {'title': 'Add Products', 'done': false},
      {'title': 'Set Up Payments', 'done': false},
      {'title': 'Launch Store', 'done': false},
    ];

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: Column(
        children: steps.asMap().entries.map((entry) {
          final index = entry.key;
          final step = entry.value;
          final isDone = step['done'] as bool;
          final isLast = index == steps.length - 1;

          return Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Column(
                children: [
                  Container(
                    width: 28,
                    height: 28,
                    decoration: BoxDecoration(
                      color: isDone ? AppColors.success : AppColors.grey200,
                      shape: BoxShape.circle,
                    ),
                    child: Center(
                      child: isDone
                          ? const Icon(Icons.check, color: Colors.white, size: 16)
                          : Text(
                              '${index + 1}',
                              style: AppTypography.labelSmall.copyWith(
                                color: AppColors.grey600,
                              ),
                            ),
                    ),
                  ),
                  if (!isLast)
                    Container(
                      width: 2,
                      height: 24,
                      color: isDone ? AppColors.success : AppColors.grey200,
                    ),
                ],
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Padding(
                  padding: EdgeInsets.only(bottom: isLast ? 0 : 16),
                  child: Text(
                    step['title'] as String,
                    style: AppTypography.bodyMedium.copyWith(
                      color: isDone ? AppColors.grey500 : AppColors.grey900,
                      decoration: isDone ? TextDecoration.lineThrough : null,
                    ),
                  ),
                ),
              ),
            ],
          );
        }).toList(),
      ),
    );
  }

  Widget _buildKeyMetrics() {
    return GridView.count(
      crossAxisCount: 2,
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      crossAxisSpacing: 12,
      mainAxisSpacing: 12,
      childAspectRatio: 1.4,
      children: [
        StatCard(
          title: 'Products Listed',
          value: '0',
          icon: Icons.inventory_2,
          iconColor: AppColors.entrepreneurColor,
        ),
        StatCard(
          title: 'Total Views',
          value: '0',
          icon: Icons.visibility,
          iconColor: AppColors.info,
        ),
        StatCard(
          title: 'Orders',
          value: '0',
          icon: Icons.receipt_long,
          iconColor: AppColors.success,
        ),
        StatCard(
          title: 'Revenue',
          value: '₦0',
          icon: Icons.account_balance_wallet,
          iconColor: AppColors.warning,
        ),
      ],
    );
  }

  Widget _buildLearningPath() {
    final courses = [
      {
        'title': 'Starting Your FMCG Business',
        'progress': 0.6,
        'lessons': 12,
        'completed': 7,
        'duration': '2h 30m',
      },
      {
        'title': 'Pricing Strategies for Profit',
        'progress': 0.0,
        'lessons': 8,
        'completed': 0,
        'duration': '1h 45m',
      },
    ];

    return Column(
      children: courses.map((course) {
        return Container(
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
                width: 60,
                height: 60,
                decoration: BoxDecoration(
                  color: AppColors.info.withValues(alpha: 0.1),
                  borderRadius: AppRadius.borderRadiusSm,
                ),
                child: const Icon(Icons.play_circle_filled, color: AppColors.info, size: 32),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      course['title'] as String,
                      style: AppTypography.titleSmall,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Row(
                      children: [
                        Text(
                          '${course['completed']}/${course['lessons']} lessons',
                          style: AppTypography.labelSmall.copyWith(
                            color: AppColors.grey500,
                          ),
                        ),
                        const SizedBox(width: 8),
                        const Icon(Icons.schedule, size: 12, color: AppColors.grey500),
                        const SizedBox(width: 2),
                        Text(
                          course['duration'] as String,
                          style: AppTypography.labelSmall.copyWith(
                            color: AppColors.grey500,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    LinearProgressIndicator(
                      value: course['progress'] as double,
                      backgroundColor: AppColors.grey200,
                      valueColor: const AlwaysStoppedAnimation<Color>(AppColors.info),
                      minHeight: 4,
                      borderRadius: BorderRadius.circular(2),
                    ),
                  ],
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Widget _buildFundingOpportunities() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.success.withValues(alpha: 0.3)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(
                  color: AppColors.successBg,
                  borderRadius: AppRadius.borderRadiusSm,
                ),
                child: const Icon(Icons.savings, color: AppColors.success),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Starter Business Loan',
                      style: AppTypography.titleSmall,
                    ),
                    Text(
                      'Up to ₦500,000 at 2% monthly',
                      style: AppTypography.bodySmall.copyWith(
                        color: AppColors.grey600,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              _buildFundingTag('No Collateral'),
              const SizedBox(width: 8),
              _buildFundingTag('Quick Approval'),
              const SizedBox(width: 8),
              _buildFundingTag('Flexible Terms'),
            ],
          ),
          const SizedBox(height: 16),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton(
              onPressed: () {},
              style: OutlinedButton.styleFrom(
                foregroundColor: AppColors.success,
                side: const BorderSide(color: AppColors.success),
              ),
              child: const Text('Check Eligibility'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFundingTag(String text) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AppColors.grey100,
        borderRadius: AppRadius.borderRadiusFull,
      ),
      child: Text(
        text,
        style: AppTypography.labelSmall.copyWith(
          color: AppColors.grey700,
        ),
      ),
    );
  }

  Widget _buildMentors() {
    final mentors = [
      {'name': 'Chidi Okonkwo', 'expertise': 'Retail Strategy', 'rating': 4.9},
      {'name': 'Amina Yusuf', 'expertise': 'Supply Chain', 'rating': 4.8},
    ];

    return SizedBox(
      height: 140,
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: mentors.length,
        itemBuilder: (context, index) {
          final mentor = mentors[index];
          return Container(
            width: 200,
            margin: EdgeInsets.only(right: index < mentors.length - 1 ? 12 : 0),
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: AppColors.white,
              borderRadius: AppRadius.borderRadiusMd,
              border: Border.all(color: AppColors.cardBorder),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    CircleAvatar(
                      radius: 20,
                      backgroundColor: AppColors.entrepreneurColor.withValues(alpha: 0.2),
                      child: Text(
                        (mentor['name'] as String).substring(0, 1),
                        style: AppTypography.titleSmall.copyWith(
                          color: AppColors.entrepreneurColor,
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            mentor['name'] as String,
                            style: AppTypography.labelMedium,
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                          Text(
                            mentor['expertise'] as String,
                            style: AppTypography.labelSmall.copyWith(
                              color: AppColors.grey500,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                const Spacer(),
                Row(
                  children: [
                    const Icon(Icons.star, color: Color(0xFFFFB300), size: 16),
                    const SizedBox(width: 4),
                    Text(
                      '${mentor['rating']}',
                      style: AppTypography.labelSmall.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    const Spacer(),
                    TextButton(
                      onPressed: () {},
                      style: TextButton.styleFrom(
                        padding: const EdgeInsets.symmetric(horizontal: 12),
                        minimumSize: Size.zero,
                        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      ),
                      child: const Text('Connect'),
                    ),
                  ],
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildCommunityFeed() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              CircleAvatar(
                radius: 18,
                backgroundColor: AppColors.primary.withValues(alpha: 0.2),
                child: const Text('AO', style: TextStyle(fontSize: 12)),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Adeola Ogundimu',
                      style: AppTypography.labelMedium,
                    ),
                    Text(
                      '2 hours ago',
                      style: AppTypography.labelSmall.copyWith(
                        color: AppColors.grey500,
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Text(
            'Just made my first ₦50,000 sale through OmniRoute! The platform made it so easy to connect with wholesalers. 🎉',
            style: AppTypography.bodyMedium,
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              TextButton.icon(
                onPressed: () {},
                icon: const Icon(Icons.thumb_up_outlined, size: 18),
                label: const Text('24'),
                style: TextButton.styleFrom(
                  foregroundColor: AppColors.grey600,
                  padding: EdgeInsets.zero,
                ),
              ),
              const SizedBox(width: 16),
              TextButton.icon(
                onPressed: () {},
                icon: const Icon(Icons.chat_bubble_outline, size: 18),
                label: const Text('8'),
                style: TextButton.styleFrom(
                  foregroundColor: AppColors.grey600,
                  padding: EdgeInsets.zero,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
