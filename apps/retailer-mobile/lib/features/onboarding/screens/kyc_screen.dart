/// OmniRoute Ecosystem - KYC Verification Screen
/// Multi-step document upload and verification process

import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:image_picker/image_picker.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/core/router/app_router.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class KycScreen extends ConsumerStatefulWidget {
  const KycScreen({super.key});

  @override
  ConsumerState<KycScreen> createState() => _KycScreenState();
}

class _KycScreenState extends ConsumerState<KycScreen> {
  final ImagePicker _picker = ImagePicker();
  final Map<String, File?> _uploadedDocuments = {};
  int _currentStep = 0;
  bool _isUploading = false;

  late ParticipantType _participantType;
  late List<String> _requiredDocuments;

  @override
  void initState() {
    super.initState();
    final user = ref.read(authProvider).user;
    _participantType = user?.participantType ?? ParticipantType.retailer;
    _requiredDocuments = _participantType.requiredDocuments;
  }

  double get _progress => _uploadedDocuments.values.where((f) => f != null).length / _requiredDocuments.length;

  Future<void> _pickDocument(String documentType, ImageSource source) async {
    final XFile? image = await _picker.pickImage(source: source, imageQuality: 80);
    if (image != null) {
      setState(() => _uploadedDocuments[documentType] = File(image.path));
    }
  }

  void _showPickerOptions(String documentType) {
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Upload $documentType', style: AppTypography.titleMedium),
            const SizedBox(height: 24),
            Row(
              children: [
                Expanded(
                  child: _PickerOption(
                    icon: Icons.camera_alt,
                    label: 'Camera',
                    onTap: () {
                      Navigator.pop(context);
                      _pickDocument(documentType, ImageSource.camera);
                    },
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _PickerOption(
                    icon: Icons.photo_library,
                    label: 'Gallery',
                    onTap: () {
                      Navigator.pop(context);
                      _pickDocument(documentType, ImageSource.gallery);
                    },
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  Future<void> _submitKyc() async {
    if (_uploadedDocuments.values.where((f) => f != null).length < _requiredDocuments.length) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please upload all required documents')),
      );
      return;
    }

    setState(() => _isUploading = true);

    // Simulate upload delay
    await Future.delayed(const Duration(seconds: 2));

    ref.read(authProvider.notifier).setOnboardingComplete();

    setState(() => _isUploading = false);

    if (mounted) {
      _showSuccessDialog();
    }
  }

  void _showSuccessDialog() {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 80,
              height: 80,
              decoration: BoxDecoration(
                color: AppColors.successBg,
                shape: BoxShape.circle,
              ),
              child: const Icon(Icons.check_circle, color: AppColors.success, size: 48),
            ).animate().scale(begin: const Offset(0.5, 0.5), curve: Curves.elasticOut),
            const SizedBox(height: 24),
            Text('Verification Submitted!', style: AppTypography.titleLarge),
            const SizedBox(height: 8),
            Text('Your documents are being reviewed. This usually takes 24-48 hours.',
                style: AppTypography.bodyMedium.copyWith(color: AppColors.grey600), textAlign: TextAlign.center),
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: () {
                  Navigator.pop(context);
                  context.go(RoutePaths.dashboard);
                },
                child: const Text('Go to Dashboard'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.scaffoldBackground,
      appBar: AppBar(
        title: const Text('Account Verification'),
        backgroundColor: AppColors.white,
        elevation: 0,
      ),
      body: Column(
        children: [
          _buildProgressHeader(),
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(24),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildStepIndicator(),
                  const SizedBox(height: 24),
                  ..._buildDocumentCards(),
                ],
              ),
            ),
          ),
          _buildBottomBar(),
        ],
      ),
    );
  }

  Widget _buildProgressHeader() {
    return Container(
      padding: const EdgeInsets.all(20),
      color: AppColors.white,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text('${(_progress * 100).toInt()}% Complete', style: AppTypography.labelMedium.copyWith(color: AppColors.primary)),
              Text('${_uploadedDocuments.values.where((f) => f != null).length}/${_requiredDocuments.length} documents',
                  style: AppTypography.bodySmall.copyWith(color: AppColors.grey600)),
            ],
          ),
          const SizedBox(height: 8),
          LinearProgressIndicator(
            value: _progress,
            backgroundColor: AppColors.grey200,
            valueColor: const AlwaysStoppedAnimation(AppColors.primary),
            minHeight: 8,
            borderRadius: AppRadius.borderRadiusFull,
          ),
        ],
      ),
    );
  }

  Widget _buildStepIndicator() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: AppColors.infoBg,
        borderRadius: AppRadius.borderRadiusMd,
        border: Border.all(color: AppColors.info.withValues(alpha: 0.2)),
      ),
      child: Row(
        children: [
          const Icon(Icons.info_outline, color: AppColors.info),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Verifying as ${_participantType.displayName}',
                    style: AppTypography.labelMedium.copyWith(color: AppColors.info)),
                Text('Please upload clear photos of all required documents',
                    style: AppTypography.bodySmall.copyWith(color: AppColors.grey600)),
              ],
            ),
          ),
        ],
      ),
    ).animate().fadeIn().slideY(begin: 0.2, end: 0);
  }

  List<Widget> _buildDocumentCards() {
    return _requiredDocuments.asMap().entries.map((entry) {
      final index = entry.key;
      final docType = entry.value;
      final isUploaded = _uploadedDocuments[docType] != null;

      return Container(
        margin: const EdgeInsets.only(bottom: 16),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color: AppColors.white,
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: isUploaded ? AppColors.success : AppColors.grey200),
        ),
        child: Row(
          children: [
            Container(
              width: 56,
              height: 56,
              decoration: BoxDecoration(
                color: isUploaded ? AppColors.successBg : AppColors.grey100,
                borderRadius: AppRadius.borderRadiusSm,
              ),
              child: isUploaded
                  ? ClipRRect(
                      borderRadius: AppRadius.borderRadiusSm,
                      child: Image.file(_uploadedDocuments[docType]!, fit: BoxFit.cover),
                    )
                  : Icon(Icons.description_outlined, color: AppColors.grey400),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(docType, style: AppTypography.titleSmall),
                  const SizedBox(height: 4),
                  Text(isUploaded ? 'Uploaded successfully' : 'Required',
                      style: AppTypography.bodySmall.copyWith(
                          color: isUploaded ? AppColors.success : AppColors.grey500)),
                ],
              ),
            ),
            if (isUploaded)
              IconButton(
                icon: const Icon(Icons.check_circle, color: AppColors.success),
                onPressed: () => _showPickerOptions(docType),
              )
            else
              ElevatedButton(
                onPressed: () => _showPickerOptions(docType),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.primary,
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                ),
                child: const Text('Upload'),
              ),
          ],
        ),
      ).animate(delay: Duration(milliseconds: 100 * index)).fadeIn().slideX(begin: 0.1, end: 0);
    }).toList();
  }

  Widget _buildBottomBar() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: AppColors.white,
        boxShadow: AppShadows.bottomNav,
      ),
      child: SafeArea(
        top: false,
        child: SizedBox(
          width: double.infinity,
          height: 56,
          child: ElevatedButton(
            onPressed: _isUploading ? null : _submitKyc,
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.primary,
              foregroundColor: Colors.white,
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
              elevation: 0,
            ),
            child: _isUploading
                ? const Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, valueColor: AlwaysStoppedAnimation(Colors.white))),
                      SizedBox(width: 12),
                      Text('Uploading...'),
                    ],
                  )
                : const Text('Submit for Verification', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
          ),
        ),
      ),
    );
  }
}

class _PickerOption extends StatelessWidget {
  final IconData icon;
  final String label;
  final VoidCallback onTap;

  const _PickerOption({required this.icon, required this.label, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(
          color: AppColors.grey50,
          borderRadius: AppRadius.borderRadiusMd,
          border: Border.all(color: AppColors.grey200),
        ),
        child: Column(
          children: [
            Icon(icon, size: 32, color: AppColors.primary),
            const SizedBox(height: 8),
            Text(label, style: AppTypography.labelMedium),
          ],
        ),
      ),
    );
  }
}
