/// OmniRoute Ecosystem - Data Models
/// Type-safe models using freezed for immutability

import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:omniroute_ecosystem/core/constants/app_constants.dart';

part 'models.freezed.dart';
part 'models.g.dart';

// ============================================================================
// USER & AUTH MODELS
// ============================================================================

@freezed
class User with _$User {
  const factory User({
    required String id,
    required String email,
    required String phone,
    required String fullName,
    required ParticipantType participantType,
    String? profileImageUrl,
    String? businessName,
    @Default(false) bool isVerified,
    @Default(false) bool isOnboardingComplete,
    @Default('active') String status,
    DateTime? createdAt,
    DateTime? lastLoginAt,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

@freezed
class AuthResponse with _$AuthResponse {
  const factory AuthResponse({
    required String accessToken,
    required String refreshToken,
    required User user,
    required int expiresIn,
  }) = _AuthResponse;

  factory AuthResponse.fromJson(Map<String, dynamic> json) =>
      _$AuthResponseFromJson(json);
}

// ============================================================================
// KYC & ONBOARDING MODELS
// ============================================================================

@freezed
class KycDocument with _$KycDocument {
  const factory KycDocument({
    required String id,
    required String documentType,
    required String documentUrl,
    @Default('pending') String status,
    String? rejectionReason,
    DateTime? uploadedAt,
    DateTime? verifiedAt,
  }) = _KycDocument;

  factory KycDocument.fromJson(Map<String, dynamic> json) =>
      _$KycDocumentFromJson(json);
}

@freezed
class KycSubmission with _$KycSubmission {
  const factory KycSubmission({
    required String id,
    required String userId,
    required ParticipantType participantType,
    required List<KycDocument> documents,
    @Default('pending') String status,
    String? reviewNotes,
    DateTime? submittedAt,
    DateTime? reviewedAt,
  }) = _KycSubmission;

  factory KycSubmission.fromJson(Map<String, dynamic> json) =>
      _$KycSubmissionFromJson(json);
}

// ============================================================================
// BUSINESS PROFILE MODELS
// ============================================================================

@freezed
class BusinessProfile with _$BusinessProfile {
  const factory BusinessProfile({
    required String id,
    required String userId,
    required String businessName,
    required ParticipantType participantType,
    String? registrationNumber,
    String? taxId,
    String? address,
    String? city,
    String? state,
    String? country,
    String? postalCode,
    double? latitude,
    double? longitude,
    String? website,
    String? description,
    String? logoUrl,
    @Default([]) List<String> categories,
    @Default({}) Map<String, dynamic> metadata,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _BusinessProfile;

  factory BusinessProfile.fromJson(Map<String, dynamic> json) =>
      _$BusinessProfileFromJson(json);
}

// ============================================================================
// DASHBOARD STATS MODELS
// ============================================================================

@freezed
class DashboardStats with _$DashboardStats {
  const factory DashboardStats({
    required double totalRevenue,
    required double revenueGrowth,
    required int totalOrders,
    required int orderGrowth,
    required int totalCustomers,
    required int customerGrowth,
    required double averageOrderValue,
    required List<ChartDataPoint> revenueChart,
    required List<ChartDataPoint> ordersChart,
  }) = _DashboardStats;

  factory DashboardStats.fromJson(Map<String, dynamic> json) =>
      _$DashboardStatsFromJson(json);
}

@freezed
class ChartDataPoint with _$ChartDataPoint {
  const factory ChartDataPoint({
    required String label,
    required double value,
    String? secondaryValue,
  }) = _ChartDataPoint;

  factory ChartDataPoint.fromJson(Map<String, dynamic> json) =>
      _$ChartDataPointFromJson(json);
}

// ============================================================================
// ORDER MODELS
// ============================================================================

@freezed
class Order with _$Order {
  const factory Order({
    required String id,
    required String orderNumber,
    required String customerId,
    required String customerName,
    required List<OrderItem> items,
    required double subtotal,
    required double tax,
    required double shipping,
    required double total,
    @Default('pending') String status,
    @Default('pending') String paymentStatus,
    String? deliveryAddress,
    String? notes,
    DateTime? createdAt,
    DateTime? estimatedDelivery,
    DateTime? deliveredAt,
  }) = _Order;

  factory Order.fromJson(Map<String, dynamic> json) => _$OrderFromJson(json);
}

@freezed
class OrderItem with _$OrderItem {
  const factory OrderItem({
    required String productId,
    required String productName,
    required String sku,
    required int quantity,
    required double unitPrice,
    required double total,
    String? imageUrl,
  }) = _OrderItem;

  factory OrderItem.fromJson(Map<String, dynamic> json) =>
      _$OrderItemFromJson(json);
}

// ============================================================================
// PRODUCT MODELS
// ============================================================================

@freezed
class Product with _$Product {
  const factory Product({
    required String id,
    required String name,
    required String sku,
    required String categoryId,
    required double price,
    double? compareAtPrice,
    String? description,
    @Default([]) List<String> images,
    @Default(0) int stockQuantity,
    @Default(true) bool isActive,
    @Default({}) Map<String, dynamic> attributes,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) = _Product;

  factory Product.fromJson(Map<String, dynamic> json) => _$ProductFromJson(json);
}

@freezed
class Category with _$Category {
  const factory Category({
    required String id,
    required String name,
    String? parentId,
    String? imageUrl,
    @Default(0) int productCount,
    @Default(true) bool isActive,
  }) = _Category;

  factory Category.fromJson(Map<String, dynamic> json) =>
      _$CategoryFromJson(json);
}

// ============================================================================
// INVENTORY MODELS
// ============================================================================

@freezed
class InventoryItem with _$InventoryItem {
  const factory InventoryItem({
    required String id,
    required String productId,
    required String productName,
    required String sku,
    required int currentStock,
    required int reservedStock,
    required int availableStock,
    required int reorderLevel,
    required int reorderQuantity,
    String? warehouseId,
    String? warehouseName,
    DateTime? lastRestocked,
    DateTime? updatedAt,
  }) = _InventoryItem;

  factory InventoryItem.fromJson(Map<String, dynamic> json) =>
      _$InventoryItemFromJson(json);
}

@freezed
class StockAlert with _$StockAlert {
  const factory StockAlert({
    required String id,
    required String productId,
    required String productName,
    required String alertType, // 'low_stock', 'out_of_stock', 'overstock'
    required int currentStock,
    required int threshold,
    @Default(false) bool isResolved,
    DateTime? createdAt,
  }) = _StockAlert;

  factory StockAlert.fromJson(Map<String, dynamic> json) =>
      _$StockAlertFromJson(json);
}

// ============================================================================
// LOGISTICS MODELS
// ============================================================================

@freezed
class Delivery with _$Delivery {
  const factory Delivery({
    required String id,
    required String orderId,
    required String orderNumber,
    required String pickupAddress,
    required String deliveryAddress,
    required double pickupLat,
    required double pickupLng,
    required double deliveryLat,
    required double deliveryLng,
    String? driverId,
    String? driverName,
    String? vehicleNumber,
    @Default('pending') String status,
    double? estimatedDistance,
    int? estimatedDuration,
    DateTime? scheduledAt,
    DateTime? pickedUpAt,
    DateTime? deliveredAt,
    @Default([]) List<DeliveryUpdate> updates,
  }) = _Delivery;

  factory Delivery.fromJson(Map<String, dynamic> json) =>
      _$DeliveryFromJson(json);
}

@freezed
class DeliveryUpdate with _$DeliveryUpdate {
  const factory DeliveryUpdate({
    required String status,
    required String message,
    double? latitude,
    double? longitude,
    required DateTime timestamp,
  }) = _DeliveryUpdate;

  factory DeliveryUpdate.fromJson(Map<String, dynamic> json) =>
      _$DeliveryUpdateFromJson(json);
}

@freezed
class Vehicle with _$Vehicle {
  const factory Vehicle({
    required String id,
    required String vehicleNumber,
    required String vehicleType, // 'truck', 'van', 'bike', 'car'
    required String make,
    required String model,
    required int year,
    String? driverId,
    String? driverName,
    @Default('available') String status,
    double? currentLat,
    double? currentLng,
    double? fuelLevel,
    int? odometer,
    DateTime? lastServiceDate,
    DateTime? nextServiceDue,
  }) = _Vehicle;

  factory Vehicle.fromJson(Map<String, dynamic> json) => _$VehicleFromJson(json);
}

// ============================================================================
// FINANCIAL MODELS
// ============================================================================

@freezed
class Wallet with _$Wallet {
  const factory Wallet({
    required String id,
    required String userId,
    required double balance,
    required double pendingBalance,
    required String currency,
    @Default(true) bool isActive,
    DateTime? lastTransactionAt,
  }) = _Wallet;

  factory Wallet.fromJson(Map<String, dynamic> json) => _$WalletFromJson(json);
}

@freezed
class Transaction with _$Transaction {
  const factory Transaction({
    required String id,
    required String reference,
    required String type, // 'credit', 'debit'
    required String category, // 'order_payment', 'settlement', 'withdrawal', etc.
    required double amount,
    required String currency,
    @Default('pending') String status,
    String? description,
    String? counterpartyId,
    String? counterpartyName,
    @Default({}) Map<String, dynamic> metadata,
    DateTime? createdAt,
    DateTime? completedAt,
  }) = _Transaction;

  factory Transaction.fromJson(Map<String, dynamic> json) =>
      _$TransactionFromJson(json);
}

@freezed
class Settlement with _$Settlement {
  const factory Settlement({
    required String id,
    required String reference,
    required double amount,
    required String currency,
    required String bankName,
    required String accountNumber,
    required String accountName,
    @Default('pending') String status,
    DateTime? initiatedAt,
    DateTime? completedAt,
  }) = _Settlement;

  factory Settlement.fromJson(Map<String, dynamic> json) =>
      _$SettlementFromJson(json);
}

@freezed
class LoanApplication with _$LoanApplication {
  const factory LoanApplication({
    required String id,
    required String userId,
    required double requestedAmount,
    required double approvedAmount,
    required double interestRate,
    required int tenorDays,
    @Default('pending') String status,
    String? purpose,
    double? monthlyPayment,
    DateTime? appliedAt,
    DateTime? approvedAt,
    DateTime? disbursedAt,
    DateTime? dueDate,
  }) = _LoanApplication;

  factory LoanApplication.fromJson(Map<String, dynamic> json) =>
      _$LoanApplicationFromJson(json);
}

// ============================================================================
// WAREHOUSE MODELS
// ============================================================================

@freezed
class Warehouse with _$Warehouse {
  const factory Warehouse({
    required String id,
    required String name,
    required String address,
    required String city,
    required String state,
    required double latitude,
    required double longitude,
    required int totalCapacity,
    required int usedCapacity,
    @Default(true) bool isActive,
    String? managerId,
    String? managerName,
    String? phone,
    String? email,
    @Default([]) List<String> capabilities, // 'cold_storage', 'hazmat', etc.
    @Default({}) Map<String, dynamic> operatingHours,
  }) = _Warehouse;

  factory Warehouse.fromJson(Map<String, dynamic> json) =>
      _$WarehouseFromJson(json);
}

@freezed
class StorageUnit with _$StorageUnit {
  const factory StorageUnit({
    required String id,
    required String warehouseId,
    required String unitNumber,
    required String unitType, // 'rack', 'bin', 'shelf', 'zone'
    required double capacity,
    required double usedCapacity,
    String? zone,
    String? aisle,
    String? level,
    @Default(true) bool isActive,
    @Default([]) List<String> productIds,
  }) = _StorageUnit;

  factory StorageUnit.fromJson(Map<String, dynamic> json) =>
      _$StorageUnitFromJson(json);
}

// ============================================================================
// NOTIFICATION MODELS
// ============================================================================

@freezed
class AppNotification with _$AppNotification {
  const factory AppNotification({
    required String id,
    required String title,
    required String body,
    required String type,
    @Default(false) bool isRead,
    String? actionUrl,
    @Default({}) Map<String, dynamic> data,
    DateTime? createdAt,
  }) = _AppNotification;

  factory AppNotification.fromJson(Map<String, dynamic> json) =>
      _$AppNotificationFromJson(json);
}

// ============================================================================
// ANALYTICS MODELS
// ============================================================================

@freezed
class SalesAnalytics with _$SalesAnalytics {
  const factory SalesAnalytics({
    required double totalSales,
    required double salesGrowth,
    required int totalTransactions,
    required double averageTransactionValue,
    required List<ChartDataPoint> salesByPeriod,
    required List<ProductSales> topProducts,
    required List<CustomerSales> topCustomers,
  }) = _SalesAnalytics;

  factory SalesAnalytics.fromJson(Map<String, dynamic> json) =>
      _$SalesAnalyticsFromJson(json);
}

@freezed
class ProductSales with _$ProductSales {
  const factory ProductSales({
    required String productId,
    required String productName,
    required int quantitySold,
    required double revenue,
  }) = _ProductSales;

  factory ProductSales.fromJson(Map<String, dynamic> json) =>
      _$ProductSalesFromJson(json);
}

@freezed
class CustomerSales with _$CustomerSales {
  const factory CustomerSales({
    required String customerId,
    required String customerName,
    required int orderCount,
    required double totalSpent,
  }) = _CustomerSales;

  factory CustomerSales.fromJson(Map<String, dynamic> json) =>
      _$CustomerSalesFromJson(json);
}
