/// OmniRoute Ecosystem - Settings Screen
/// Comprehensive settings for profile, security, preferences, and app configuration

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/core/router/app_router.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class SettingsListScreen extends ConsumerStatefulWidget {
  const SettingsListScreen({super.key});

  @override
  ConsumerState<SettingsListScreen> createState() => _SettingsListScreenState();
}

class _SettingsListScreenState extends ConsumerState<SettingsListScreen> {
  bool _notificationsEnabled = true;
  bool _biometricsEnabled = false;
  bool _darkModeEnabled = false;
  String _selectedLanguage = 'English';
  String _selectedCurrency = 'NGN (₦)';

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authProvider);
    final user = authState.user;

    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        backgroundColor: AppColors.white,
        elevation: 0,
        title: Text(
          'Settings',
          style: AppTypography.titleLarge.copyWith(color: AppColors.grey900),
        ),
      ),
      body: SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Profile Section
            _buildProfileSection(user),
            const SizedBox(height: 8),
            
            // Account Settings
            _buildSection(
              title: 'Account',
              children: [
                _SettingsTile(
                  icon: Icons.person_outline,
                  title: 'Personal Information',
                  subtitle: 'Name, email, phone number',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.business,
                  title: 'Business Profile',
                  subtitle: 'Business details, documents',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.account_balance,
                  title: 'Bank Accounts',
                  subtitle: 'Manage linked bank accounts',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.location_on_outlined,
                  title: 'Addresses',
                  subtitle: 'Manage delivery addresses',
                  onTap: () {},
                ),
              ],
            ),
            
            // Security Settings
            _buildSection(
              title: 'Security',
              children: [
                _SettingsTile(
                  icon: Icons.lock_outline,
                  title: 'Change Password',
                  subtitle: 'Update your password',
                  onTap: () => _showChangePasswordSheet(context),
                ),
                _SettingsTile(
                  icon: Icons.fingerprint,
                  title: 'Biometric Login',
                  subtitle: 'Use fingerprint or face ID',
                  trailing: Switch(
                    value: _biometricsEnabled,
                    onChanged: (value) => setState(() => _biometricsEnabled = value),
                    activeColor: AppColors.primary,
                  ),
                ),
                _SettingsTile(
                  icon: Icons.security,
                  title: 'Two-Factor Authentication',
                  subtitle: 'Add extra security to your account',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.devices,
                  title: 'Active Sessions',
                  subtitle: 'Manage devices logged into your account',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.history,
                  title: 'Login History',
                  subtitle: 'View recent login activity',
                  onTap: () {},
                ),
              ],
            ),
            
            // Notifications
            _buildSection(
              title: 'Notifications',
              children: [
                _SettingsTile(
                  icon: Icons.notifications_outlined,
                  title: 'Push Notifications',
                  subtitle: 'Receive order and delivery updates',
                  trailing: Switch(
                    value: _notificationsEnabled,
                    onChanged: (value) => setState(() => _notificationsEnabled = value),
                    activeColor: AppColors.primary,
                  ),
                ),
                _SettingsTile(
                  icon: Icons.email_outlined,
                  title: 'Email Notifications',
                  subtitle: 'Receive email updates',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.sms_outlined,
                  title: 'SMS Notifications',
                  subtitle: 'Receive SMS alerts',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.tune,
                  title: 'Notification Preferences',
                  subtitle: 'Customize what notifications you receive',
                  onTap: () {},
                ),
              ],
            ),
            
            // App Preferences
            _buildSection(
              title: 'Preferences',
              children: [
                _SettingsTile(
                  icon: Icons.dark_mode_outlined,
                  title: 'Dark Mode',
                  subtitle: 'Switch to dark theme',
                  trailing: Switch(
                    value: _darkModeEnabled,
                    onChanged: (value) => setState(() => _darkModeEnabled = value),
                    activeColor: AppColors.primary,
                  ),
                ),
                _SettingsTile(
                  icon: Icons.language,
                  title: 'Language',
                  subtitle: _selectedLanguage,
                  onTap: () => _showLanguageSheet(context),
                ),
                _SettingsTile(
                  icon: Icons.attach_money,
                  title: 'Currency',
                  subtitle: _selectedCurrency,
                  onTap: () => _showCurrencySheet(context),
                ),
              ],
            ),
            
            // Business Settings
            _buildSection(
              title: 'Business',
              children: [
                _SettingsTile(
                  icon: Icons.receipt_long,
                  title: 'Invoice Settings',
                  subtitle: 'Customize invoice templates',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.local_shipping_outlined,
                  title: 'Delivery Settings',
                  subtitle: 'Manage delivery preferences',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.inventory_2_outlined,
                  title: 'Inventory Alerts',
                  subtitle: 'Set low stock thresholds',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.people_outline,
                  title: 'Team Management',
                  subtitle: 'Manage team members and permissions',
                  onTap: () {},
                ),
              ],
            ),
            
            // Integrations
            _buildSection(
              title: 'Integrations',
              children: [
                _SettingsTile(
                  icon: Icons.store,
                  title: 'Marketplace Connections',
                  subtitle: 'Jumia, Konga, Jiji integrations',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.account_balance_wallet_outlined,
                  title: 'Payment Gateways',
                  subtitle: 'Paystack, Flutterwave settings',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.analytics_outlined,
                  title: 'Accounting Software',
                  subtitle: 'Connect QuickBooks, Xero',
                  onTap: () {},
                ),
              ],
            ),
            
            // Support
            _buildSection(
              title: 'Support',
              children: [
                _SettingsTile(
                  icon: Icons.help_outline,
                  title: 'Help Center',
                  subtitle: 'FAQs and guides',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.chat_bubble_outline,
                  title: 'Contact Support',
                  subtitle: 'Chat with our support team',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.bug_report_outlined,
                  title: 'Report a Problem',
                  subtitle: 'Help us improve the app',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.feedback_outlined,
                  title: 'Give Feedback',
                  subtitle: 'Share your thoughts with us',
                  onTap: () {},
                ),
              ],
            ),
            
            // About
            _buildSection(
              title: 'About',
              children: [
                _SettingsTile(
                  icon: Icons.info_outline,
                  title: 'About OmniRoute',
                  subtitle: 'Version 1.0.0',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.description_outlined,
                  title: 'Terms of Service',
                  subtitle: 'Read our terms and conditions',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.privacy_tip_outlined,
                  title: 'Privacy Policy',
                  subtitle: 'How we handle your data',
                  onTap: () {},
                ),
                _SettingsTile(
                  icon: Icons.star_outline,
                  title: 'Rate Us',
                  subtitle: 'Love the app? Rate us on the store',
                  onTap: () {},
                ),
              ],
            ),
            
            // Danger Zone
            const SizedBox(height: 16),
            _buildSection(
              title: 'Danger Zone',
              titleColor: AppColors.error,
              children: [
                _SettingsTile(
                  icon: Icons.logout,
                  title: 'Sign Out',
                  subtitle: 'Sign out of your account',
                  iconColor: AppColors.error,
                  titleColor: AppColors.error,
                  onTap: () => _showLogoutConfirmation(context),
                ),
                _SettingsTile(
                  icon: Icons.delete_forever,
                  title: 'Delete Account',
                  subtitle: 'Permanently delete your account and data',
                  iconColor: AppColors.error,
                  titleColor: AppColors.error,
                  onTap: () => _showDeleteAccountConfirmation(context),
                ),
              ],
            ),
            const SizedBox(height: 32),
          ],
        ),
      ),
    );
  }

  Widget _buildProfileSection(dynamic user) {
    return Container(
      margin: const EdgeInsets.all(16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: AppColors.white,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.cardBorder),
      ),
      child: Row(
        children: [
          Container(
            width: 64,
            height: 64,
            decoration: BoxDecoration(
              color: AppColors.primary.withValues(alpha: 0.1),
              shape: BoxShape.circle,
            ),
            child: Center(
              child: Text(
                user?.fullName?.substring(0, 2).toUpperCase() ?? 'OR',
                style: AppTypography.titleLarge.copyWith(
                  color: AppColors.primary,
                  fontWeight: FontWeight.w700,
                ),
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  user?.fullName ?? 'User Name',
                  style: AppTypography.titleMedium,
                ),
                const SizedBox(height: 4),
                Text(
                  user?.email ?? 'user@email.com',
                  style: AppTypography.bodySmall.copyWith(
                    color: AppColors.grey600,
                  ),
                ),
                const SizedBox(height: 4),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: AppColors.successBg,
                    borderRadius: AppRadius.borderRadiusFull,
                  ),
                  child: Text(
                    user?.participantType.displayName ?? 'Participant',
                    style: AppTypography.labelSmall.copyWith(
                      color: AppColors.success,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.edit_outlined, color: AppColors.primary),
            onPressed: () {},
          ),
        ],
      ),
    ).animate().fadeIn().slideY(begin: 0.1, end: 0);
  }

  Widget _buildSection({
    required String title,
    required List<Widget> children,
    Color? titleColor,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
          child: Text(
            title.toUpperCase(),
            style: AppTypography.labelSmall.copyWith(
              color: titleColor ?? AppColors.grey600,
              fontWeight: FontWeight.w600,
              letterSpacing: 1.2,
            ),
          ),
        ),
        Container(
          margin: const EdgeInsets.symmetric(horizontal: 16),
          decoration: BoxDecoration(
            color: AppColors.white,
            borderRadius: AppRadius.borderRadiusMd,
            border: Border.all(color: AppColors.cardBorder),
          ),
          child: Column(
            children: children.asMap().entries.map((entry) {
              final isLast = entry.key == children.length - 1;
              return Column(
                children: [
                  entry.value,
                  if (!isLast) const Divider(height: 1, indent: 56),
                ],
              );
            }).toList(),
          ),
        ),
      ],
    );
  }

  void _showChangePasswordSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        height: MediaQuery.of(context).size.height * 0.5,
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
                  Text('Change Password', style: AppTypography.titleLarge),
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
                  children: [
                    TextField(
                      obscureText: true,
                      decoration: InputDecoration(
                        labelText: 'Current Password',
                        border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      obscureText: true,
                      decoration: InputDecoration(
                        labelText: 'New Password',
                        border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
                      ),
                    ),
                    const SizedBox(height: 16),
                    TextField(
                      obscureText: true,
                      decoration: InputDecoration(
                        labelText: 'Confirm New Password',
                        border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
                      ),
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
                  child: const Text('Update Password'),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _showLanguageSheet(BuildContext context) {
    final languages = ['English', 'Hausa', 'Yoruba', 'Igbo', 'Pidgin'];
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        decoration: const BoxDecoration(
          color: AppColors.white,
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
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
              child: Text('Select Language', style: AppTypography.titleLarge),
            ),
            ...languages.map((lang) => ListTile(
                  title: Text(lang),
                  trailing: _selectedLanguage == lang
                      ? const Icon(Icons.check, color: AppColors.primary)
                      : null,
                  onTap: () {
                    setState(() => _selectedLanguage = lang);
                    Navigator.pop(context);
                  },
                )),
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }

  void _showCurrencySheet(BuildContext context) {
    final currencies = ['NGN (₦)', 'USD (\$)', 'GBP (£)', 'EUR (€)'];
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        decoration: const BoxDecoration(
          color: AppColors.white,
          borderRadius: BorderRadius.vertical(top: Radius.circular(24)),
        ),
        child: Column(
          mainAxisSize: MainAxisSize.min,
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
              child: Text('Select Currency', style: AppTypography.titleLarge),
            ),
            ...currencies.map((curr) => ListTile(
                  title: Text(curr),
                  trailing: _selectedCurrency == curr
                      ? const Icon(Icons.check, color: AppColors.primary)
                      : null,
                  onTap: () {
                    setState(() => _selectedCurrency = curr);
                    Navigator.pop(context);
                  },
                )),
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }

  void _showLogoutConfirmation(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: AppRadius.borderRadiusMd),
        title: const Text('Sign Out'),
        content: const Text('Are you sure you want to sign out of your account?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              ref.read(authProvider.notifier).logout();
              context.go(RoutePaths.welcome);
            },
            style: ElevatedButton.styleFrom(backgroundColor: AppColors.error),
            child: const Text('Sign Out'),
          ),
        ],
      ),
    );
  }

  void _showDeleteAccountConfirmation(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: AppRadius.borderRadiusMd),
        title: Row(
          children: [
            const Icon(Icons.warning_amber, color: AppColors.error),
            const SizedBox(width: 8),
            const Text('Delete Account'),
          ],
        ),
        content: const Text(
          'This action is permanent and cannot be undone. All your data, orders, and settings will be permanently deleted.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              // Implement account deletion
              Navigator.pop(context);
            },
            style: ElevatedButton.styleFrom(backgroundColor: AppColors.error),
            child: const Text('Delete Forever'),
          ),
        ],
      ),
    );
  }
}

class _SettingsTile extends StatelessWidget {
  final IconData icon;
  final String title;
  final String subtitle;
  final Widget? trailing;
  final VoidCallback? onTap;
  final Color? iconColor;
  final Color? titleColor;

  const _SettingsTile({
    required this.icon,
    required this.title,
    required this.subtitle,
    this.trailing,
    this.onTap,
    this.iconColor,
    this.titleColor,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: Container(
        width: 40,
        height: 40,
        decoration: BoxDecoration(
          color: (iconColor ?? AppColors.primary).withValues(alpha: 0.1),
          borderRadius: AppRadius.borderRadiusSm,
        ),
        child: Icon(icon, color: iconColor ?? AppColors.primary, size: 20),
      ),
      title: Text(
        title,
        style: AppTypography.titleSmall.copyWith(color: titleColor),
      ),
      subtitle: Text(
        subtitle,
        style: AppTypography.bodySmall.copyWith(color: AppColors.grey600),
      ),
      trailing: trailing ?? const Icon(Icons.chevron_right, color: AppColors.grey400),
      onTap: onTap,
    );
  }
}
