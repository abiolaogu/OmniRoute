// OmniRoute Shared Types
// Nigerian Market B2B FMCG Commerce Platform

// =============================================================================
// CORE TYPES
// =============================================================================

export interface Money {
    amount: number;
    currency: 'NGN' | 'USD' | 'GBP' | 'EUR';
    formatted: string;
}

export interface Address {
    id: string;
    line1: string;
    line2?: string;
    city: string;
    state: NigerianState;
    lga: string;
    postalCode?: string;
    country: string;
    lat?: number;
    lng?: number;
    isDefault: boolean;
}

export type NigerianState =
    | 'Abia' | 'Adamawa' | 'Akwa Ibom' | 'Anambra' | 'Bauchi' | 'Bayelsa'
    | 'Benue' | 'Borno' | 'Cross River' | 'Delta' | 'Ebonyi' | 'Edo'
    | 'Ekiti' | 'Enugu' | 'FCT' | 'Gombe' | 'Imo' | 'Jigawa' | 'Kaduna'
    | 'Kano' | 'Katsina' | 'Kebbi' | 'Kogi' | 'Kwara' | 'Lagos' | 'Nasarawa'
    | 'Niger' | 'Ogun' | 'Ondo' | 'Osun' | 'Oyo' | 'Plateau' | 'Rivers'
    | 'Sokoto' | 'Taraba' | 'Yobe' | 'Zamfara';

export interface Phone {
    countryCode: string;
    number: string;
    formatted: string;
    isVerified: boolean;
}

export interface Pagination {
    page: number;
    pageSize: number;
    total: number;
    totalPages: number;
}

export interface ListResponse<T> {
    data: T[];
    pagination: Pagination;
}

// =============================================================================
// USER & AUTH TYPES
// =============================================================================

export type UserRole =
    | 'super_admin' | 'platform_admin' | 'bank_admin' | 'bank_analyst'
    | 'manufacturer' | 'distributor' | 'wholesaler' | 'retailer'
    | 'logistics_provider' | 'warehouse_operator' | 'consumer';

export interface User {
    id: string;
    email: string;
    phone: Phone;
    firstName: string;
    lastName: string;
    fullName: string;
    avatar?: string;
    role: UserRole;
    tenantId?: string;
    permissions: string[];
    isActive: boolean;
    emailVerified: boolean;
    phoneVerified: boolean;
    createdAt: Date;
    lastLoginAt?: Date;
}

export interface Session {
    user: User;
    accessToken: string;
    refreshToken: string;
    expiresAt: Date;
}

export interface HasuraClaims {
    'x-hasura-user-id': string;
    'x-hasura-default-role': UserRole;
    'x-hasura-allowed-roles': UserRole[];
    'x-hasura-tenant-id'?: string;
    name: string;
    email: string;
    exp: number;
    iat: number;
}

export interface Tenant {
    id: string;
    name: string;
    slug: string;
    type: 'manufacturer' | 'distributor' | 'wholesaler' | 'retailer' | 'bank';
    logo?: string;
    address: Address;
    phone: Phone;
    email: string;
    settings: TenantSettings;
    isActive: boolean;
    createdAt: Date;
}

export interface TenantSettings {
    currency: 'NGN';
    timezone: 'Africa/Lagos';
    taxRate: number;
    creditEnabled: boolean;
    maxCreditLimit: number;
    paymentTerms: number; // days
}

// =============================================================================
// COMMERCE TYPES
// =============================================================================

export type OrderStatus =
    | 'draft' | 'pending' | 'confirmed' | 'processing'
    | 'ready' | 'shipped' | 'delivered' | 'cancelled' | 'refunded';

export type OrderSource = 'web' | 'mobile' | 'ussd' | 'whatsapp' | 'api';

export interface Order {
    id: string;
    orderNumber: string;
    customerId: string;
    customer: Customer;
    tenantId: string;
    status: OrderStatus;
    source: OrderSource;
    items: OrderItem[];
    subtotal: Money;
    tax: Money;
    discount: Money;
    shippingCost: Money;
    total: Money;
    shippingAddress: Address;
    billingAddress: Address;
    notes?: string;
    metadata: Record<string, unknown>;
    createdAt: Date;
    updatedAt: Date;
    confirmedAt?: Date;
    shippedAt?: Date;
    deliveredAt?: Date;
}

export interface OrderItem {
    id: string;
    orderId: string;
    productId: string;
    product: Product;
    variantId?: string;
    variant?: ProductVariant;
    quantity: number;
    unitPrice: Money;
    total: Money;
    discount: Money;
}

export interface Product {
    id: string;
    sku: string;
    name: string;
    slug: string;
    description: string;
    shortDescription?: string;
    categoryId: string;
    category: Category;
    brandId?: string;
    brand?: Brand;
    images: string[];
    basePrice: Money;
    comparePrice?: Money;
    costPrice: Money;
    variants: ProductVariant[];
    options: ProductOption[];
    inventory: InventoryLevel[];
    isActive: boolean;
    isFeatured: boolean;
    tags: string[];
    createdAt: Date;
    updatedAt: Date;
}

export interface ProductVariant {
    id: string;
    productId: string;
    sku: string;
    name: string;
    options: Record<string, string>; // e.g., { size: 'Large', color: 'Red' }
    price: Money;
    costPrice: Money;
    image?: string;
    isDefault: boolean;
}

export interface ProductOption {
    id: string;
    name: string;
    values: string[];
    position: number;
}

export interface Category {
    id: string;
    name: string;
    slug: string;
    description?: string;
    image?: string;
    parentId?: string;
    children?: Category[];
    position: number;
    isActive: boolean;
}

export interface Brand {
    id: string;
    name: string;
    slug: string;
    logo?: string;
    description?: string;
}

// =============================================================================
// CUSTOMER TYPES
// =============================================================================

export type CustomerType = 'b2b' | 'b2c';
export type CustomerTier = 'standard' | 'silver' | 'gold' | 'platinum';

export interface Customer {
    id: string;
    type: CustomerType;
    email?: string;
    phone: Phone;
    firstName: string;
    lastName: string;
    fullName: string;
    companyName?: string;
    tier: CustomerTier;
    addresses: Address[];
    defaultAddressId?: string;
    creditFacility?: CreditFacility;
    totalOrders: number;
    totalSpent: Money;
    lastOrderAt?: Date;
    notes?: string;
    tags: string[];
    isActive: boolean;
    createdAt: Date;
}

export interface CreditFacility {
    id: string;
    customerId: string;
    creditLimit: Money;
    availableCredit: Money;
    usedCredit: Money;
    paymentTerms: number; // days
    interestRate: number;
    isApproved: boolean;
    approvedAt?: Date;
    approvedBy?: string;
}

// =============================================================================
// INVENTORY TYPES
// =============================================================================

export interface InventoryLevel {
    id: string;
    productId: string;
    variantId?: string;
    locationId: string;
    location: InventoryLocation;
    quantity: number;
    reserved: number;
    available: number;
    reorderPoint: number;
    reorderQuantity: number;
}

export interface InventoryLocation {
    id: string;
    name: string;
    type: 'warehouse' | 'store' | 'transit';
    address: Address;
    isActive: boolean;
}

export interface InventoryMovement {
    id: string;
    productId: string;
    variantId?: string;
    fromLocationId?: string;
    toLocationId?: string;
    type: 'in' | 'out' | 'transfer' | 'adjustment';
    quantity: number;
    reason: string;
    reference?: string;
    createdAt: Date;
    createdBy: string;
}

// =============================================================================
// FINANCIAL TYPES
// =============================================================================

export type PaymentStatus = 'pending' | 'processing' | 'completed' | 'failed' | 'refunded';
export type PaymentMethodType = 'card' | 'bank_transfer' | 'ussd' | 'mobile_money' | 'credit';

export interface Payment {
    id: string;
    orderId?: string;
    customerId: string;
    amount: Money;
    fee: Money;
    net: Money;
    method: PaymentMethodType;
    status: PaymentStatus;
    reference: string;
    providerReference?: string;
    metadata: Record<string, unknown>;
    paidAt?: Date;
    createdAt: Date;
}

export interface Settlement {
    id: string;
    tenantId: string;
    amount: Money;
    fee: Money;
    net: Money;
    status: 'pending' | 'processing' | 'completed' | 'failed';
    bankName: string;
    accountNumber: string;
    accountName: string;
    reference: string;
    settledAt?: Date;
    createdAt: Date;
}

// =============================================================================
// DELIVERY TYPES
// =============================================================================

export type FulfillmentStatus = 'pending' | 'assigned' | 'picking' | 'packed' | 'shipped' | 'delivered' | 'failed';

export interface Fulfillment {
    id: string;
    orderId: string;
    status: FulfillmentStatus;
    items: FulfillmentItem[];
    trackingNumber?: string;
    carrier?: string;
    estimatedDelivery?: Date;
    actualDelivery?: Date;
    proofOfDelivery?: string;
    driverId?: string;
    driver?: GigWorker;
    createdAt: Date;
}

export interface FulfillmentItem {
    id: string;
    fulfillmentId: string;
    orderItemId: string;
    quantity: number;
    pickedAt?: Date;
    packedAt?: Date;
}

export interface GigWorker {
    id: string;
    userId: string;
    firstName: string;
    lastName: string;
    phone: Phone;
    avatar?: string;
    vehicleType: 'motorcycle' | 'car' | 'van' | 'truck';
    vehiclePlate?: string;
    rating: number;
    totalDeliveries: number;
    isOnline: boolean;
    currentLocation?: { lat: number; lng: number };
}

// =============================================================================
// ANALYTICS TYPES
// =============================================================================

export interface DashboardMetrics {
    revenue: Money;
    revenueChange: number;
    orders: number;
    ordersChange: number;
    customers: number;
    customersChange: number;
    avgOrderValue: Money;
    avgOrderValueChange: number;
}

export interface ChartDataPoint {
    date: string;
    value: number;
    label?: string;
}

export interface TopProduct {
    productId: string;
    productName: string;
    quantity: number;
    revenue: Money;
    image?: string;
}

export interface RecentActivity {
    id: string;
    type: 'order' | 'payment' | 'customer' | 'inventory';
    message: string;
    metadata: Record<string, unknown>;
    createdAt: Date;
}
