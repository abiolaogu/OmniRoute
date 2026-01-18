/// OmniRoute Ecosystem - Splash Screen
import 'package:flutter/material.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:omniroute_ecosystem/core/theme/app_theme.dart';
import 'package:omniroute_ecosystem/core/router/app_router.dart';
import 'package:omniroute_ecosystem/providers/app_providers.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});
  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _navigateAfterDelay();
  }

  Future<void> _navigateAfterDelay() async {
    await Future.delayed(const Duration(milliseconds: 2500));
    if (!mounted) return;
    final authState = ref.read(authProvider);
    if (authState.isAuthenticated) {
      context.go(authState.isOnboardingComplete ? RoutePaths.dashboard : RoutePaths.kyc);
    } else {
      context.go(RoutePaths.welcome);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(begin: Alignment.topLeft, end: Alignment.bottomRight, colors: [Color(0xFF0D47A1), Color(0xFF1565C0), Color(0xFF1976D2)]),
        ),
        child: SafeArea(
          child: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: 120, height: 120,
                  decoration: BoxDecoration(color: Colors.white, shape: BoxShape.circle, boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.2), blurRadius: 30, offset: const Offset(0, 10))]),
                  child: Center(
                    child: Stack(
                      alignment: Alignment.center,
                      children: [
                        SizedBox(width: 100, height: 100, child: CircularProgressIndicator(strokeWidth: 3, valueColor: AlwaysStoppedAnimation<Color>(AppColors.primary.withValues(alpha: 0.3)))).animate(onPlay: (c) => c.repeat()).rotate(duration: const Duration(seconds: 3)),
                        const Text('OR', style: TextStyle(fontSize: 40, fontWeight: FontWeight.w800, color: Color(0xFF0D47A1), fontFamily: 'SpaceGrotesk')),
                      ],
                    ),
                  ),
                ).animate().scale(begin: const Offset(0.5, 0.5), end: const Offset(1, 1), duration: const Duration(milliseconds: 600), curve: Curves.elasticOut).fadeIn(duration: const Duration(milliseconds: 400)),
                const SizedBox(height: 24),
                const Text('OmniRoute', style: TextStyle(fontSize: 36, fontWeight: FontWeight.w700, color: Colors.white, fontFamily: 'SpaceGrotesk', letterSpacing: -0.5)).animate(delay: const Duration(milliseconds: 400)).fadeIn(duration: const Duration(milliseconds: 500)).slideY(begin: 0.3, end: 0, curve: Curves.easeOutCubic),
                const SizedBox(height: 8),
                Text('Commerce Ecosystem Platform', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w400, color: Colors.white.withValues(alpha: 0.85), letterSpacing: 0.5)).animate(delay: const Duration(milliseconds: 600)).fadeIn(duration: const Duration(milliseconds: 500)).slideY(begin: 0.3, end: 0, curve: Curves.easeOutCubic),
                const SizedBox(height: 48),
                SizedBox(width: 24, height: 24, child: CircularProgressIndicator(strokeWidth: 2.5, valueColor: AlwaysStoppedAnimation<Color>(Colors.white.withValues(alpha: 0.8)))).animate(delay: const Duration(milliseconds: 800)).fadeIn(duration: const Duration(milliseconds: 400)),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
