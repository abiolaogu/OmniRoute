/// OmniRoute Ecosystem - Register Screen
/// Multi-step registration with participant-specific requirements

import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/core/router/app_router.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class RegisterScreen extends ConsumerStatefulWidget {
  final ParticipantType? participantType;

  const RegisterScreen({super.key, this.participantType});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _fullNameController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  final _businessNameController = TextEditingController();

  bool _obscurePassword = true;
  bool _obscureConfirmPassword = true;
  bool _isLoading = false;
  bool _acceptedTerms = false;
  late ParticipantType _selectedType;

  @override
  void initState() {
    super.initState();
    _selectedType = widget.participantType ??
        ref.read(selectedParticipantTypeProvider) ??
        ParticipantType.retailer;
  }

  @override
  void dispose() {
    _fullNameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    _businessNameController.dispose();
    super.dispose();
  }

  bool get _requiresBusinessName {
    return [
      ParticipantType.bank,
      ParticipantType.logistics,
      ParticipantType.warehouse,
      ParticipantType.manufacturer,
      ParticipantType.distributor,
      ParticipantType.wholesaler,
      ParticipantType.ecommerce,
    ].contains(_selectedType);
  }

  Future<void> _handleRegister() async {
    if (!_formKey.currentState!.validate()) return;
    if (!_acceptedTerms) {
      _showSnackBar('Please accept the terms and conditions');
      return;
    }

    setState(() => _isLoading = true);

    final success = await ref.read(authProvider.notifier).register(
      fullName: _fullNameController.text.trim(),
      email: _emailController.text.trim(),
      phone: _phoneController.text.trim(),
      password: _passwordController.text,
      participantType: _selectedType,
    );

    setState(() => _isLoading = false);

    if (success && mounted) {
      context.push(RoutePaths.otpVerification, extra: _phoneController.text.trim());
    } else if (mounted) {
      final error = ref.read(authProvider).error;
      _showSnackBar(error ?? 'Registration failed. Please try again.');
    }
  }

  void _showSnackBar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: AppColors.error,
        behavior: SnackBarBehavior.floating,
        margin: const EdgeInsets.all(16),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back, color: AppColors.grey800),
          onPressed: () => context.pop(),
        ),
        title: Text(
          'Create Account',
          style: AppTypography.titleLarge.copyWith(color: AppColors.grey900),
        ),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Participant type badge
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                decoration: BoxDecoration(
                  color: _getParticipantColor().withValues(alpha: 0.1),
                  borderRadius: AppRadius.borderRadiusMd,
                  border: Border.all(
                    color: _getParticipantColor().withValues(alpha: 0.3),
                  ),
                ),
                child: Row(
                  children: [
                    Icon(
                      _getParticipantIcon(),
                      color: _getParticipantColor(),
                      size: 20,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            'Registering as',
                            style: AppTypography.labelSmall.copyWith(
                              color: AppColors.grey600,
                            ),
                          ),
                          Text(
                            _selectedType.displayName,
                            style: AppTypography.titleSmall.copyWith(
                              color: _getParticipantColor(),
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                    TextButton(
                      onPressed: () => context.pop(),
                      child: Text(
                        'Change',
                        style: AppTypography.labelMedium.copyWith(
                          color: AppColors.primary,
                        ),
                      ),
                    ),
                  ],
                ),
              )
                  .animate()
                  .fadeIn()
                  .slideY(begin: -0.1, end: 0),

              const SizedBox(height: 24),

              // Full Name
              _buildLabel('Full Name'),
              const SizedBox(height: 8),
              _buildTextField(
                controller: _fullNameController,
                hint: 'Enter your full name',
                prefixIcon: Icons.person_outline,
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your full name';
                  }
                  if (value.split(' ').length < 2) {
                    return 'Please enter your full name (first and last)';
                  }
                  return null;
                },
              )
                  .animate(delay: const Duration(milliseconds: 100))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 20),

              // Business Name (conditional)
              if (_requiresBusinessName) ...[
                _buildLabel('Business Name'),
                const SizedBox(height: 8),
                _buildTextField(
                  controller: _businessNameController,
                  hint: 'Enter your business name',
                  prefixIcon: Icons.business,
                  textInputAction: TextInputAction.next,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter your business name';
                    }
                    return null;
                  },
                )
                    .animate(delay: const Duration(milliseconds: 150))
                    .fadeIn()
                    .slideY(begin: 0.1, end: 0),
                const SizedBox(height: 20),
              ],

              // Email
              _buildLabel('Email Address'),
              const SizedBox(height: 8),
              _buildTextField(
                controller: _emailController,
                hint: 'Enter your email address',
                prefixIcon: Icons.email_outlined,
                keyboardType: TextInputType.emailAddress,
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your email';
                  }
                  if (!ValidationPatterns.email.hasMatch(value)) {
                    return 'Please enter a valid email address';
                  }
                  return null;
                },
              )
                  .animate(delay: const Duration(milliseconds: 200))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 20),

              // Phone
              _buildLabel('Phone Number'),
              const SizedBox(height: 8),
              _buildTextField(
                controller: _phoneController,
                hint: '08012345678',
                prefixIcon: Icons.phone_outlined,
                keyboardType: TextInputType.phone,
                textInputAction: TextInputAction.next,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter your phone number';
                  }
                  if (!ValidationPatterns.phone.hasMatch(value)) {
                    return 'Please enter a valid Nigerian phone number';
                  }
                  return null;
                },
              )
                  .animate(delay: const Duration(milliseconds: 300))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 20),

              // Password
              _buildLabel('Password'),
              const SizedBox(height: 8),
              _buildTextField(
                controller: _passwordController,
                hint: 'Create a strong password',
                prefixIcon: Icons.lock_outline,
                obscureText: _obscurePassword,
                textInputAction: TextInputAction.next,
                suffixIcon: IconButton(
                  icon: Icon(
                    _obscurePassword
                        ? Icons.visibility_outlined
                        : Icons.visibility_off_outlined,
                    size: 20,
                    color: AppColors.grey500,
                  ),
                  onPressed: () {
                    setState(() => _obscurePassword = !_obscurePassword);
                  },
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter a password';
                  }
                  if (value.length < 8) {
                    return 'Password must be at least 8 characters';
                  }
                  if (!RegExp(r'[A-Z]').hasMatch(value)) {
                    return 'Password must contain at least one uppercase letter';
                  }
                  if (!RegExp(r'[0-9]').hasMatch(value)) {
                    return 'Password must contain at least one number';
                  }
                  return null;
                },
              )
                  .animate(delay: const Duration(milliseconds: 400))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 20),

              // Confirm Password
              _buildLabel('Confirm Password'),
              const SizedBox(height: 8),
              _buildTextField(
                controller: _confirmPasswordController,
                hint: 'Confirm your password',
                prefixIcon: Icons.lock_outline,
                obscureText: _obscureConfirmPassword,
                textInputAction: TextInputAction.done,
                onFieldSubmitted: (_) => _handleRegister(),
                suffixIcon: IconButton(
                  icon: Icon(
                    _obscureConfirmPassword
                        ? Icons.visibility_outlined
                        : Icons.visibility_off_outlined,
                    size: 20,
                    color: AppColors.grey500,
                  ),
                  onPressed: () {
                    setState(() =>
                        _obscureConfirmPassword = !_obscureConfirmPassword);
                  },
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please confirm your password';
                  }
                  if (value != _passwordController.text) {
                    return 'Passwords do not match';
                  }
                  return null;
                },
              )
                  .animate(delay: const Duration(milliseconds: 500))
                  .fadeIn()
                  .slideY(begin: 0.1, end: 0),

              const SizedBox(height: 24),

              // Terms checkbox
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(
                    width: 24,
                    height: 24,
                    child: Checkbox(
                      value: _acceptedTerms,
                      onChanged: (value) {
                        setState(() => _acceptedTerms = value ?? false);
                      },
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(4),
                      ),
                      activeColor: AppColors.primary,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Text.rich(
                      TextSpan(
                        text: 'I agree to the ',
                        style: AppTypography.bodySmall.copyWith(
                          color: AppColors.grey700,
                        ),
                        children: [
                          TextSpan(
                            text: 'Terms of Service',
                            style: AppTypography.labelSmall.copyWith(
                              color: AppColors.primary,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                          const TextSpan(text: ' and '),
                          TextSpan(
                            text: 'Privacy Policy',
                            style: AppTypography.labelSmall.copyWith(
                              color: AppColors.primary,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              )
                  .animate(delay: const Duration(milliseconds: 600))
                  .fadeIn(),

              const SizedBox(height: 32),

              // Register button
              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: _isLoading ? null : _handleRegister,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: _getParticipantColor(),
                    foregroundColor: Colors.white,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    elevation: 0,
                    disabledBackgroundColor: AppColors.grey300,
                  ),
                  child: _isLoading
                      ? const SizedBox(
                          width: 24,
                          height: 24,
                          child: CircularProgressIndicator(
                            strokeWidth: 2.5,
                            valueColor: AlwaysStoppedAnimation<Color>(
                              Colors.white,
                            ),
                          ),
                        )
                      : const Text(
                          'Create Account',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                ),
              )
                  .animate(delay: const Duration(milliseconds: 700))
                  .fadeIn()
                  .slideY(begin: 0.2, end: 0),

              const SizedBox(height: 24),

              // Sign in link
              Center(
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      'Already have an account? ',
                      style: AppTypography.bodyMedium.copyWith(
                        color: AppColors.grey600,
                      ),
                    ),
                    GestureDetector(
                      onTap: () => context.push(RoutePaths.login),
                      child: Text(
                        'Sign In',
                        style: AppTypography.labelLarge.copyWith(
                          color: AppColors.primary,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ],
                ),
              )
                  .animate(delay: const Duration(milliseconds: 800))
                  .fadeIn(),

              const SizedBox(height: 32),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildLabel(String text) {
    return Text(
      text,
      style: AppTypography.labelMedium.copyWith(
        color: AppColors.grey700,
        fontWeight: FontWeight.w600,
      ),
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String hint,
    required IconData prefixIcon,
    TextInputType? keyboardType,
    TextInputAction? textInputAction,
    bool obscureText = false,
    Widget? suffixIcon,
    String? Function(String?)? validator,
    void Function(String)? onFieldSubmitted,
  }) {
    return TextFormField(
      controller: controller,
      keyboardType: keyboardType,
      textInputAction: textInputAction,
      obscureText: obscureText,
      onFieldSubmitted: onFieldSubmitted,
      decoration: InputDecoration(
        hintText: hint,
        prefixIcon: Icon(prefixIcon, size: 20),
        suffixIcon: suffixIcon,
        filled: true,
        fillColor: AppColors.white,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: AppColors.grey300),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: AppColors.grey300),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: AppColors.primary, width: 2),
        ),
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(12),
          borderSide: const BorderSide(color: AppColors.error),
        ),
      ),
      validator: validator,
    );
  }

  Color _getParticipantColor() {
    switch (_selectedType) {
      case ParticipantType.bank:
        return AppColors.bankColor;
      case ParticipantType.logistics:
        return AppColors.logisticsColor;
      case ParticipantType.warehouse:
        return AppColors.warehouseColor;
      case ParticipantType.manufacturer:
        return AppColors.manufacturerColor;
      case ParticipantType.distributor:
        return AppColors.distributorColor;
      case ParticipantType.wholesaler:
        return AppColors.wholesalerColor;
      case ParticipantType.retailer:
        return AppColors.retailerColor;
      case ParticipantType.ecommerce:
        return AppColors.ecommerceColor;
      case ParticipantType.entrepreneur:
        return AppColors.entrepreneurColor;
      case ParticipantType.investor:
        return AppColors.investorColor;
      case ParticipantType.agent:
        return AppColors.agentColor;
      case ParticipantType.driver:
        return AppColors.driverColor;
    }
  }

  IconData _getParticipantIcon() {
    switch (_selectedType) {
      case ParticipantType.bank:
        return Icons.account_balance;
      case ParticipantType.logistics:
        return Icons.local_shipping;
      case ParticipantType.warehouse:
        return Icons.warehouse;
      case ParticipantType.manufacturer:
        return Icons.factory;
      case ParticipantType.distributor:
        return Icons.inventory_2;
      case ParticipantType.wholesaler:
        return Icons.storefront;
      case ParticipantType.retailer:
        return Icons.store;
      case ParticipantType.ecommerce:
        return Icons.shopping_bag;
      case ParticipantType.entrepreneur:
        return Icons.lightbulb;
      case ParticipantType.investor:
        return Icons.trending_up;
      case ParticipantType.agent:
        return Icons.person_pin;
      case ParticipantType.driver:
        return Icons.directions_car;
    }
  }
}
