/// OmniRoute Ecosystem - Application Constants
/// Defines all participant types and core constants for the multi-tenant ecosystem

library;

// ============================================================================
// PARTICIPANT TYPES
// ============================================================================

/// All ecosystem participant types supported by the platform
enum ParticipantType {
  bank('Bank / Financial Institution', 'bank', '🏦'),
  logistics('Logistics Company', 'logistics', '🚚'),
  warehouse('Warehouse Operator', 'warehouse', '🏭'),
  manufacturer('Manufacturer', 'manufacturer', '⚙️'),
  distributor('Distributor', 'distributor', '📦'),
  wholesaler('Wholesaler', 'wholesaler', '🏪'),
  retailer('Retailer', 'retailer', '🛒'),
  ecommerce('E-commerce / Dropshipper', 'ecommerce', '🛍️'),
  entrepreneur('Entrepreneur', 'entrepreneur', '💡'),
  investor('Investor', 'investor', '💰'),
  agent('Field Agent', 'agent', '👤'),
  driver('Delivery Driver', 'driver', '🚗');

  final String displayName;
  final String code;
  final String emoji;

  const ParticipantType(this.displayName, this.code, this.emoji);

  /// Returns appropriate dashboard route for participant type
  String get dashboardRoute => '/dashboard/$code';

  /// Returns onboarding requirements for KYC
  List<String> get requiredDocuments {
    switch (this) {
      case ParticipantType.bank:
        return [
          'CBN License',
          'Certificate of Incorporation',
          'Board Resolution',
          'AML/CFT Policy',
          'Directors ID',
        ];
      case ParticipantType.logistics:
        return [
          'Business Registration (CAC)',
          'Fleet Documentation',
          'Insurance Certificates',
          'Driver Licenses',
          'Vehicle Particulars',
        ];
      case ParticipantType.warehouse:
        return [
          'Business Registration (CAC)',
          'Warehouse License',
          'Fire Safety Certificate',
          'Property Documents',
          'Insurance Certificate',
        ];
      case ParticipantType.manufacturer:
        return [
          'Business Registration (CAC)',
          'NAFDAC Registration',
          'SON Certification',
          'Environmental Compliance',
          'Factory License',
        ];
      case ParticipantType.distributor:
      case ParticipantType.wholesaler:
        return [
          'Business Registration (CAC)',
          'Tax Clearance',
          'Trade License',
          'Bank Statement (6 months)',
        ];
      case ParticipantType.retailer:
        return [
          'Business Name Registration',
          'Valid ID (NIN/Voter\'s Card)',
          'Utility Bill',
          'Shop Photo',
        ];
      case ParticipantType.ecommerce:
        return [
          'Business Registration',
          'Valid ID',
          'Bank Account Details',
          'Platform Links',
        ];
      case ParticipantType.entrepreneur:
        return [
          'Valid ID (NIN/Voter\'s Card)',
          'Proof of Address',
          'Business Plan (Optional)',
        ];
      case ParticipantType.investor:
        return [
          'Valid ID',
          'Proof of Address',
          'Source of Funds Declaration',
          'Investment Profile',
        ];
      case ParticipantType.agent:
      case ParticipantType.driver:
        return [
          'Valid ID (NIN/Voter\'s Card)',
          'Proof of Address',
          'Guarantor Details',
          'Bank Account',
        ];
    }
  }
}

// ============================================================================
// API ENDPOINTS
// ============================================================================

class ApiEndpoints {
  static const String baseUrl = 'https://api.omniroute.io/v1';
  static const String stagingUrl = 'https://staging-api.omniroute.io/v1';

  // Auth
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String verifyOtp = '/auth/verify-otp';
  static const String refreshToken = '/auth/refresh';
  static const String forgotPassword = '/auth/forgot-password';

  // Onboarding
  static const String participantTypes = '/onboarding/participant-types';
  static const String submitKyc = '/onboarding/kyc';
  static const String uploadDocument = '/onboarding/documents';
  static const String verifyBusiness = '/onboarding/verify-business';

  // Dashboard
  static const String dashboardStats = '/dashboard/stats';
  static const String notifications = '/dashboard/notifications';
  static const String activities = '/dashboard/activities';

  // Orders
  static const String orders = '/orders';
  static const String orderDetails = '/orders/{id}';
  static const String createOrder = '/orders/create';
  static const String updateOrderStatus = '/orders/{id}/status';

  // Products
  static const String products = '/products';
  static const String productCategories = '/products/categories';
  static const String productSearch = '/products/search';

  // Inventory
  static const String inventory = '/inventory';
  static const String stockLevels = '/inventory/stock-levels';
  static const String stockAlerts = '/inventory/alerts';

  // Logistics
  static const String deliveries = '/logistics/deliveries';
  static const String routes = '/logistics/routes';
  static const String tracking = '/logistics/tracking/{id}';
  static const String fleet = '/logistics/fleet';

  // Finance
  static const String wallet = '/finance/wallet';
  static const String transactions = '/finance/transactions';
  static const String settlements = '/finance/settlements';
  static const String loans = '/finance/loans';
  static const String invoices = '/finance/invoices';

  // Analytics
  static const String salesAnalytics = '/analytics/sales';
  static const String inventoryAnalytics = '/analytics/inventory';
  static const String performanceMetrics = '/analytics/performance';
}

// ============================================================================
// STORAGE KEYS
// ============================================================================

class StorageKeys {
  static const String accessToken = 'access_token';
  static const String refreshToken = 'refresh_token';
  static const String userId = 'user_id';
  static const String participantType = 'participant_type';
  static const String userProfile = 'user_profile';
  static const String onboardingComplete = 'onboarding_complete';
  static const String biometricsEnabled = 'biometrics_enabled';
  static const String themeMode = 'theme_mode';
  static const String locale = 'locale';
  static const String fcmToken = 'fcm_token';
  static const String lastSyncTime = 'last_sync_time';
}

// ============================================================================
// FEATURE FLAGS
// ============================================================================

class FeatureFlags {
  static const bool enableBiometrics = true;
  static const bool enablePushNotifications = true;
  static const bool enableOfflineMode = true;
  static const bool enableDarkMode = true;
  static const bool enableAnalytics = true;
  static const bool enableCrashReporting = true;
  static const bool enableInAppUpdates = true;
}

// ============================================================================
// VALIDATION PATTERNS
// ============================================================================

class ValidationPatterns {
  static final RegExp email = RegExp(
    r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
  );
  static final RegExp phone = RegExp(r'^(\+234|0)[789][01]\d{8}$');
  static final RegExp nin = RegExp(r'^\d{11}$');
  static final RegExp bvn = RegExp(r'^\d{11}$');
  static final RegExp cacNumber = RegExp(r'^RC\d+$|^BN\d+$');
  static final RegExp accountNumber = RegExp(r'^\d{10}$');
}

// ============================================================================
// DATE FORMATS
// ============================================================================

class DateFormats {
  static const String apiFormat = 'yyyy-MM-dd\'T\'HH:mm:ss.SSS\'Z\'';
  static const String displayDate = 'MMM dd, yyyy';
  static const String displayDateTime = 'MMM dd, yyyy • HH:mm';
  static const String shortDate = 'dd/MM/yyyy';
  static const String monthYear = 'MMMM yyyy';
  static const String timeOnly = 'HH:mm';
}

// ============================================================================
// CURRENCY
// ============================================================================

class CurrencyConfig {
  static const String defaultCurrency = 'NGN';
  static const String currencySymbol = '₦';
  static const int decimalPlaces = 2;
  static const String thousandSeparator = ',';
  static const String decimalSeparator = '.';
}
