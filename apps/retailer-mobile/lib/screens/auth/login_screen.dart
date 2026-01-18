/// OmniRoute Ecosystem - Login Screen
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});
  @override ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  bool _obscurePassword = true;
  bool _isLoading = false;
  String _selectedRole = 'retailer';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const SizedBox(height: 40),
              _buildLogo(),
              const SizedBox(height: 48),
              _buildWelcomeText(),
              const SizedBox(height: 32),
              _buildRoleSelector(),
              const SizedBox(height: 32),
              _buildLoginForm(),
              const SizedBox(height: 24),
              _buildLoginButton(),
              const SizedBox(height: 16),
              _buildForgotPassword(),
              const SizedBox(height: 32),
              _buildDivider(),
              const SizedBox(height: 24),
              _buildSocialLogins(),
              const SizedBox(height: 32),
              _buildSignUpLink(),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildLogo() {
    return Column(children: [
      Container(
        width: 80, height: 80,
        decoration: BoxDecoration(gradient: LinearGradient(colors: [AppColors.primary, AppColors.primary.withValues(alpha: 0.8)]), borderRadius: AppRadius.borderRadiusLg),
        child: const Icon(Icons.route, color: Colors.white, size: 48),
      ),
      const SizedBox(height: 16),
      Text('OmniRoute', style: AppTypography.headlineMedium.copyWith(fontWeight: FontWeight.w700, color: AppColors.primary)),
      Text('Commerce Ecosystem', style: AppTypography.labelMedium.copyWith(color: AppColors.textSecondary)),
    ]);
  }

  Widget _buildWelcomeText() {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Text('Welcome back!', style: AppTypography.headlineSmall.copyWith(fontWeight: FontWeight.w600)),
      const SizedBox(height: 4),
      Text('Sign in to continue to your dashboard', style: AppTypography.bodyMedium.copyWith(color: AppColors.textSecondary)),
    ]);
  }

  Widget _buildRoleSelector() {
    final roles = [
      {'id': 'manufacturer', 'label': 'Manufacturer', 'icon': Icons.factory, 'color': AppColors.manufacturerColor},
      {'id': 'distributor', 'label': 'Distributor', 'icon': Icons.local_shipping, 'color': AppColors.distributorColor},
      {'id': 'retailer', 'label': 'Retailer', 'icon': Icons.store, 'color': AppColors.retailerColor},
      {'id': 'agent', 'label': 'Agent', 'icon': Icons.support_agent, 'color': AppColors.agentColor},
      {'id': 'driver', 'label': 'Driver', 'icon': Icons.delivery_dining, 'color': AppColors.driverColor},
    ];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('I am a...', style: AppTypography.labelMedium.copyWith(fontWeight: FontWeight.w600)),
        const SizedBox(height: 12),
        Wrap(
          spacing: 8, runSpacing: 8,
          children: roles.map((r) {
            final isSelected = _selectedRole == r['id'];
            return GestureDetector(
              onTap: () => setState(() => _selectedRole = r['id'] as String),
              child: Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                decoration: BoxDecoration(
                  color: isSelected ? (r['color'] as Color) : Colors.white,
                  borderRadius: AppRadius.borderRadiusSm,
                  border: Border.all(color: isSelected ? (r['color'] as Color) : AppColors.borderColor),
                ),
                child: Row(mainAxisSize: MainAxisSize.min, children: [
                  Icon(r['icon'] as IconData, size: 16, color: isSelected ? Colors.white : (r['color'] as Color)),
                  const SizedBox(width: 4),
                  Text(r['label'] as String, style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: isSelected ? Colors.white : AppColors.textPrimary)),
                ]),
              ),
            );
          }).toList(),
        ),
      ],
    );
  }

  Widget _buildLoginForm() {
    return Form(
      key: _formKey,
      child: Column(children: [
        TextFormField(
          controller: _emailController,
          keyboardType: TextInputType.emailAddress,
          decoration: InputDecoration(
            labelText: 'Email or Phone',
            prefixIcon: const Icon(Icons.email_outlined),
            border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
            filled: true, fillColor: Colors.white,
          ),
          validator: (v) => v!.isEmpty ? 'Please enter your email' : null,
        ),
        const SizedBox(height: 16),
        TextFormField(
          controller: _passwordController,
          obscureText: _obscurePassword,
          decoration: InputDecoration(
            labelText: 'Password',
            prefixIcon: const Icon(Icons.lock_outline),
            suffixIcon: IconButton(icon: Icon(_obscurePassword ? Icons.visibility_off : Icons.visibility), onPressed: () => setState(() => _obscurePassword = !_obscurePassword)),
            border: OutlineInputBorder(borderRadius: AppRadius.borderRadiusMd),
            filled: true, fillColor: Colors.white,
          ),
          validator: (v) => v!.isEmpty ? 'Please enter your password' : null,
        ),
      ]),
    );
  }

  Widget _buildLoginButton() {
    return SizedBox(
      height: 52,
      child: ElevatedButton(
        onPressed: _isLoading ? null : _handleLogin,
        style: ElevatedButton.styleFrom(backgroundColor: AppColors.primary, shape: RoundedRectangleBorder(borderRadius: AppRadius.borderRadiusMd)),
        child: _isLoading ? const CircularProgressIndicator(color: Colors.white) : const Text('Sign In', style: TextStyle(color: Colors.white, fontSize: 16, fontWeight: FontWeight.w600)),
      ),
    );
  }

  Widget _buildForgotPassword() {
    return Center(child: TextButton(onPressed: () {}, child: const Text('Forgot Password?')));
  }

  Widget _buildDivider() {
    return Row(children: [
      const Expanded(child: Divider()),
      Padding(padding: const EdgeInsets.symmetric(horizontal: 16), child: Text('or continue with', style: AppTypography.labelSmall.copyWith(color: AppColors.textSecondary))),
      const Expanded(child: Divider()),
    ]);
  }

  Widget _buildSocialLogins() {
    return Row(mainAxisAlignment: MainAxisAlignment.center, children: [
      _buildSocialButton(Icons.g_mobiledata, 'Google'),
      const SizedBox(width: 16),
      _buildSocialButton(Icons.phone, 'Phone'),
    ]);
  }

  Widget _buildSocialButton(IconData icon, String label) {
    return OutlinedButton.icon(
      onPressed: () {},
      style: OutlinedButton.styleFrom(padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12)),
      icon: Icon(icon),
      label: Text(label),
    );
  }

  Widget _buildSignUpLink() {
    return Row(mainAxisAlignment: MainAxisAlignment.center, children: [
      Text('Don\'t have an account?', style: AppTypography.bodyMedium.copyWith(color: AppColors.textSecondary)),
      TextButton(onPressed: () {}, child: const Text('Sign Up')),
    ]);
  }

  void _handleLogin() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isLoading = true);
    await Future.delayed(const Duration(seconds: 2));
    setState(() => _isLoading = false);
    // Navigate based on role
  }
}
