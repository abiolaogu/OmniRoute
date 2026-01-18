# OmniRoute Ecosystem Platform

A comprehensive B2B FMCG ecosystem onboarding platform built with Flutter that provides tailored experiences for all participants in the commerce ecosystem.

## 🎯 Overview

OmniRoute Ecosystem is a mobile-first platform designed to connect banks, logistics companies, warehouses, manufacturers, retailers, wholesalers, e-commerce businesses, entrepreneurs, investors, and other participants in the B2B FMCG supply chain.

## 🏗️ Architecture

### Tech Stack
- **Framework**: Flutter 3.2+
- **State Management**: Riverpod
- **Navigation**: GoRouter
- **Network**: Dio with retry interceptors
- **Local Storage**: Hive + Flutter Secure Storage
- **Charts**: FL Chart, Syncfusion Charts
- **Animations**: flutter_animate

### Project Structure

```
lib/
├── core/
│   ├── constants/        # App constants, participant types, API endpoints
│   ├── network/          # API client with auth interceptors
│   ├── router/           # GoRouter configuration
│   └── theme/            # Design system (colors, typography, spacing)
│
├── features/
│   ├── splash/           # Animated splash screen
│   ├── onboarding/       # Welcome, participant selection, KYC
│   ├── auth/             # Login, register, OTP verification
│   ├── dashboard/
│   │   ├── common/       # Main dashboard shell with adaptive navigation
│   │   ├── bank/         # Banking dashboard (loans, settlements)
│   │   ├── logistics/    # Fleet, deliveries, routes
│   │   ├── warehouse/    # Inventory, inbound/outbound
│   │   ├── manufacturer/ # Products, production
│   │   ├── retailer/     # Sales, orders, inventory
│   │   ├── wholesaler/   # Bulk orders, distribution
│   │   ├── ecommerce/    # Dropshipping, marketplace integrations
│   │   ├── investor/     # Portfolio, opportunities
│   │   └── entrepreneur/ # Ideas, learning, networking
│   ├── orders/           # Order management
│   ├── inventory/        # Stock management
│   ├── wallet/           # Finance, transactions
│   └── settings/         # App configuration
│
├── models/               # Data models (freezed)
├── providers/            # Riverpod state providers
└── widgets/              # Reusable UI components
```

## 👥 Participant Types

| Type | Features |
|------|----------|
| **Bank** | Loans, settlements, ATC (Authority to Collect), analytics |
| **Logistics** | Fleet management, delivery tracking, route optimization |
| **Warehouse** | Inventory management, inbound/outbound, storage |
| **Manufacturer** | Product catalog, order fulfillment, analytics |
| **Distributor** | Distribution management, territory |
| **Wholesaler** | Bulk orders, retailer management |
| **Retailer** | POS, inventory, reordering, BNPL |
| **E-commerce** | Dropshipping, marketplace integrations, multi-channel |
| **Entrepreneur** | Business ideas, learning resources, networking |
| **Investor** | Portfolio management, investment opportunities |
| **Agent** | Field tasks, commission tracking, performance |
| **Driver** | Deliveries, earnings, navigation |

## 🎨 Design System

### Colors
- **Primary**: #0D47A1 (Deep Blue)
- **Secondary**: #00BFA5 (Teal)
- **Accent**: #FF6D00 (Orange)
- Each participant type has a unique accent color

### Typography
- **Display Font**: Space Grotesk
- **Body Font**: Inter
- Responsive scaling with constraints

### Components
- StatCard - Dashboard statistics
- WalletCard - Balance display with actions
- OrderListTile - Order summary
- DeliveryListTile - Delivery tracking
- StatusChip - Status indicators
- QuickActionGrid - Action shortcuts

## 🚀 Getting Started

### Prerequisites
- Flutter 3.2+ 
- Dart 3.2+

### Installation

```bash
# Clone the repository
git clone https://github.com/your-org/omniroute_ecosystem.git

# Navigate to project
cd omniroute_ecosystem

# Install dependencies
flutter pub get

# Generate freezed models
flutter pub run build_runner build --delete-conflicting-outputs

# Run the app
flutter run
```

### Environment Setup

Create `.env` file:
```env
API_BASE_URL=https://api.omniroute.io/v1
PAYSTACK_PUBLIC_KEY=your_key
GOOGLE_MAPS_API_KEY=your_key
```

## 📱 Features

### Authentication
- Email/phone registration
- OTP verification
- Biometric login
- Social sign-in (Google, Apple)

### Onboarding
- Participant type selection
- KYC document upload
- Business verification

### Dashboard
- Role-specific dashboards
- Real-time statistics
- Quick actions
- Recent activity

### Orders
- Order creation & tracking
- Status updates
- Filtering & search
- Bulk actions

### Inventory
- Stock tracking
- Low stock alerts
- Barcode scanning
- Reorder management

### Wallet
- Balance management
- Transaction history
- Bank withdrawals
- Loan applications
- BNPL (Buy Now Pay Later)

### Settings
- Profile management
- Security settings
- Notification preferences
- Business configuration

## 📊 State Management

Using Riverpod for:
- Authentication state
- User profile
- Orders pagination
- Inventory data
- Wallet balance
- Notifications

```dart
// Example provider
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final secureStorage = ref.watch(secureStorageProvider);
  return AuthNotifier(apiClient, secureStorage);
});
```

## 🛣️ Navigation

GoRouter with:
- Auth-guarded routes
- Participant-specific routing
- Deep linking support
- Shell routes for bottom nav

```dart
// Route based on participant type
switch (user.participantType) {
  case ParticipantType.bank:
    return const BankDashboardScreen();
  case ParticipantType.logistics:
    return const LogisticsDashboardScreen();
  // ...
}
```

## 🔒 Security

- Secure token storage (FlutterSecureStorage)
- Token refresh with interceptors
- Biometric authentication
- SSL pinning (configure in production)

## 📈 Performance

- Lazy loading with pagination
- Image caching
- Shimmer loading states
- Optimistic UI updates

## 🧪 Testing

```bash
# Unit tests
flutter test

# Integration tests
flutter test integration_test/
```

## 📦 Build

```bash
# Android APK
flutter build apk --release

# Android App Bundle
flutter build appbundle --release

# iOS
flutter build ios --release
```

## 🔄 CI/CD

Configured for:
- GitHub Actions
- Firebase App Distribution
- Play Store deployment
- App Store Connect

## 📋 Dependencies

Key packages:
- `flutter_riverpod` - State management
- `go_router` - Navigation
- `dio` - HTTP client
- `hive_flutter` - Local database
- `freezed` - Data classes
- `fl_chart` - Charts
- `flutter_animate` - Animations
- `google_maps_flutter` - Maps

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

## 📄 License

Proprietary - BillyRonks Global Limited

## 📞 Support

- Email: support@omniroute.io
- Documentation: https://docs.omniroute.io
- Issues: GitHub Issues

---

**OmniRoute** - Connecting the Commerce Ecosystem 🚀
