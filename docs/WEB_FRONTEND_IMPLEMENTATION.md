# OmniRoute Web Frontend - Complete Implementation Prompt

## Copy this entire prompt into Claude Code CLI to build the complete web frontend

---

```
Build the complete OmniRoute Multi-Portal Web Platform - a production-grade B2B/B2C FMCG commerce system for African markets.

## CRITICAL REQUIREMENTS

1. This is a REAL production application, not a demo
2. Generate ALL files with complete, working code
3. Follow enterprise patterns: Clean Architecture, DDD, SOLID
4. Optimize for Nigerian market: Naira currency, Nigerian phone formats, local addresses
5. Support offline-first where possible
6. Mobile-responsive (60%+ users on mobile)

## PROJECT STRUCTURE

Create this exact Turborepo monorepo structure:

```
omniroute-web/
├── package.json
├── pnpm-workspace.yaml
├── turbo.json
├── .env.example
├── .gitignore
├── README.md
│
├── apps/
│   ├── admin/                      # Platform Admin Portal
│   │   ├── app/
│   │   │   ├── (auth)/
│   │   │   │   ├── login/page.tsx
│   │   │   │   ├── forgot-password/page.tsx
│   │   │   │   └── layout.tsx
│   │   │   ├── (dashboard)/
│   │   │   │   ├── layout.tsx
│   │   │   │   ├── page.tsx                    # Dashboard
│   │   │   │   ├── orders/
│   │   │   │   │   ├── page.tsx                # Orders list
│   │   │   │   │   ├── [id]/page.tsx           # Order detail
│   │   │   │   │   └── create/page.tsx
│   │   │   │   ├── products/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── [id]/page.tsx
│   │   │   │   │   └── create/page.tsx
│   │   │   │   ├── customers/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── [id]/page.tsx
│   │   │   │   │   └── create/page.tsx
│   │   │   │   ├── inventory/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── transfers/page.tsx
│   │   │   │   ├── users/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── [id]/page.tsx
│   │   │   │   ├── tenants/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── [id]/page.tsx
│   │   │   │   ├── analytics/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── sales/page.tsx
│   │   │   │   │   └── customers/page.tsx
│   │   │   │   ├── finance/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── settlements/page.tsx
│   │   │   │   │   └── payouts/page.tsx
│   │   │   │   └── settings/
│   │   │   │       ├── page.tsx
│   │   │   │       ├── profile/page.tsx
│   │   │   │       └── integrations/page.tsx
│   │   │   ├── api/
│   │   │   │   └── auth/[...nextauth]/route.ts
│   │   │   ├── layout.tsx
│   │   │   └── providers.tsx
│   │   ├── components/
│   │   │   ├── layout/
│   │   │   │   ├── Sidebar.tsx
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Breadcrumb.tsx
│   │   │   │   └── UserMenu.tsx
│   │   │   ├── dashboard/
│   │   │   │   ├── StatCard.tsx
│   │   │   │   ├── RevenueChart.tsx
│   │   │   │   ├── OrdersChart.tsx
│   │   │   │   ├── TopProducts.tsx
│   │   │   │   └── ActivityFeed.tsx
│   │   │   ├── orders/
│   │   │   │   ├── OrdersTable.tsx
│   │   │   │   ├── OrderDetail.tsx
│   │   │   │   ├── OrderTimeline.tsx
│   │   │   │   ├── OrderItems.tsx
│   │   │   │   └── OrderActions.tsx
│   │   │   ├── products/
│   │   │   │   ├── ProductsTable.tsx
│   │   │   │   ├── ProductForm.tsx
│   │   │   │   ├── ProductVariants.tsx
│   │   │   │   └── InventoryLevels.tsx
│   │   │   └── customers/
│   │   │       ├── CustomersTable.tsx
│   │   │       ├── CustomerDetail.tsx
│   │   │       └── CustomerOrders.tsx
│   │   ├── hooks/
│   │   │   ├── useOrders.ts
│   │   │   ├── useProducts.ts
│   │   │   ├── useCustomers.ts
│   │   │   ├── useAnalytics.ts
│   │   │   └── useRealtime.ts
│   │   ├── lib/
│   │   │   ├── auth-provider.ts
│   │   │   ├── access-control.ts
│   │   │   └── utils.ts
│   │   ├── styles/
│   │   │   └── globals.css
│   │   ├── public/
│   │   ├── package.json
│   │   ├── next.config.js
│   │   ├── tailwind.config.js
│   │   └── tsconfig.json
│   │
│   ├── bank/                       # Bank/Financial Institution Portal
│   │   ├── app/
│   │   │   ├── (auth)/
│   │   │   │   └── login/page.tsx
│   │   │   ├── (dashboard)/
│   │   │   │   ├── layout.tsx
│   │   │   │   ├── page.tsx                    # Bank Dashboard
│   │   │   │   ├── loans/
│   │   │   │   │   ├── page.tsx                # Loan applications
│   │   │   │   │   ├── [id]/page.tsx           # Loan detail
│   │   │   │   │   └── disbursements/page.tsx
│   │   │   │   ├── atc/
│   │   │   │   │   ├── page.tsx                # ATC mandates
│   │   │   │   │   ├── collections/page.tsx
│   │   │   │   │   └── disputes/page.tsx
│   │   │   │   ├── settlements/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── reconciliation/page.tsx
│   │   │   │   ├── customers/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── [id]/page.tsx
│   │   │   │   ├── compliance/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── kyc/page.tsx
│   │   │   │   │   └── aml/page.tsx
│   │   │   │   └── reports/
│   │   │   │       ├── page.tsx
│   │   │   │       └── regulatory/page.tsx
│   │   │   └── ...
│   │   ├── components/
│   │   │   ├── loans/
│   │   │   │   ├── LoanApplicationsTable.tsx
│   │   │   │   ├── LoanDetail.tsx
│   │   │   │   ├── CreditAssessment.tsx
│   │   │   │   └── RepaymentSchedule.tsx
│   │   │   ├── atc/
│   │   │   │   ├── MandatesTable.tsx
│   │   │   │   ├── CollectionCalendar.tsx
│   │   │   │   └── FailedCollections.tsx
│   │   │   └── settlements/
│   │   │       ├── SettlementQueue.tsx
│   │   │       └── ReconciliationView.tsx
│   │   └── ...
│   │
│   ├── partner/                    # B2B Partner Portal
│   │   ├── app/
│   │   │   ├── (auth)/
│   │   │   │   └── login/page.tsx
│   │   │   ├── (dashboard)/
│   │   │   │   ├── layout.tsx
│   │   │   │   ├── page.tsx                    # Role-specific dashboard
│   │   │   │   ├── orders/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── [id]/page.tsx
│   │   │   │   │   ├── create/page.tsx
│   │   │   │   │   └── received/page.tsx
│   │   │   │   ├── products/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── [id]/page.tsx
│   │   │   │   │   └── pricing/page.tsx
│   │   │   │   ├── inventory/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── transfers/page.tsx
│   │   │   │   │   └── alerts/page.tsx
│   │   │   │   ├── customers/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── [id]/page.tsx
│   │   │   │   ├── suppliers/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── [id]/page.tsx
│   │   │   │   ├── finance/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── invoices/page.tsx
│   │   │   │   │   ├── payments/page.tsx
│   │   │   │   │   └── credit/page.tsx
│   │   │   │   ├── logistics/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── shipments/page.tsx
│   │   │   │   │   └── tracking/page.tsx
│   │   │   │   └── analytics/
│   │   │   │       ├── page.tsx
│   │   │   │       └── reports/page.tsx
│   │   │   └── ...
│   │   ├── components/
│   │   │   ├── dashboard/
│   │   │   │   ├── ManufacturerDashboard.tsx
│   │   │   │   ├── DistributorDashboard.tsx
│   │   │   │   ├── WholesalerDashboard.tsx
│   │   │   │   ├── RetailerDashboard.tsx
│   │   │   │   ├── LogisticsDashboard.tsx
│   │   │   │   └── WarehouseDashboard.tsx
│   │   │   └── ...
│   │   └── ...
│   │
│   └── shop/                       # B2C Consumer Shop
│       ├── app/
│       │   ├── (shop)/
│       │   │   ├── layout.tsx
│       │   │   ├── page.tsx                    # Homepage
│       │   │   ├── products/
│       │   │   │   ├── page.tsx                # Product listing
│       │   │   │   └── [slug]/page.tsx         # Product detail
│       │   │   ├── categories/
│       │   │   │   └── [slug]/page.tsx
│       │   │   ├── cart/page.tsx
│       │   │   ├── checkout/
│       │   │   │   ├── page.tsx
│       │   │   │   └── success/page.tsx
│       │   │   └── search/page.tsx
│       │   ├── (account)/
│       │   │   ├── layout.tsx
│       │   │   ├── orders/
│       │   │   │   ├── page.tsx
│       │   │   │   └── [id]/page.tsx
│       │   │   ├── profile/page.tsx
│       │   │   ├── addresses/page.tsx
│       │   │   └── wishlist/page.tsx
│       │   └── ...
│       ├── components/
│       │   ├── shop/
│       │   │   ├── Header.tsx
│       │   │   ├── Footer.tsx
│       │   │   ├── ProductCard.tsx
│       │   │   ├── ProductGrid.tsx
│       │   │   ├── CategoryNav.tsx
│       │   │   ├── SearchBar.tsx
│       │   │   ├── CartDrawer.tsx
│       │   │   └── CartItem.tsx
│       │   ├── checkout/
│       │   │   ├── CheckoutForm.tsx
│       │   │   ├── AddressForm.tsx
│       │   │   ├── PaymentMethods.tsx
│       │   │   └── OrderSummary.tsx
│       │   └── account/
│       │       ├── OrderHistory.tsx
│       │       └── AddressBook.tsx
│       └── ...
│
├── packages/
│   ├── ui/                         # Shared UI Components
│   │   ├── src/
│   │   │   ├── components/
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Card.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Select.tsx
│   │   │   │   ├── Table.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   ├── Drawer.tsx
│   │   │   │   ├── Tabs.tsx
│   │   │   │   ├── Badge.tsx
│   │   │   │   ├── Avatar.tsx
│   │   │   │   ├── Skeleton.tsx
│   │   │   │   ├── EmptyState.tsx
│   │   │   │   ├── ErrorBoundary.tsx
│   │   │   │   ├── DataTable/
│   │   │   │   │   ├── DataTable.tsx
│   │   │   │   │   ├── Pagination.tsx
│   │   │   │   │   ├── Filters.tsx
│   │   │   │   │   └── ColumnHeader.tsx
│   │   │   │   ├── Charts/
│   │   │   │   │   ├── AreaChart.tsx
│   │   │   │   │   ├── BarChart.tsx
│   │   │   │   │   ├── PieChart.tsx
│   │   │   │   │   └── LineChart.tsx
│   │   │   │   ├── Forms/
│   │   │   │   │   ├── FormField.tsx
│   │   │   │   │   ├── FormSection.tsx
│   │   │   │   │   ├── MoneyInput.tsx
│   │   │   │   │   ├── PhoneInput.tsx
│   │   │   │   │   ├── AddressInput.tsx
│   │   │   │   │   └── FileUpload.tsx
│   │   │   │   └── Display/
│   │   │   │       ├── MoneyDisplay.tsx
│   │   │   │       ├── DateDisplay.tsx
│   │   │   │       ├── StatusBadge.tsx
│   │   │   │       └── UserAvatar.tsx
│   │   │   └── index.ts
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   ├── api/                        # API Client & Data Provider
│   │   ├── src/
│   │   │   ├── client.ts           # Hasura GraphQL client
│   │   │   ├── data-provider.ts    # Refine data provider
│   │   │   ├── live-provider.ts    # Real-time subscriptions
│   │   │   ├── queries/
│   │   │   │   ├── orders.ts
│   │   │   │   ├── products.ts
│   │   │   │   ├── customers.ts
│   │   │   │   ├── inventory.ts
│   │   │   │   └── analytics.ts
│   │   │   ├── mutations/
│   │   │   │   ├── orders.ts
│   │   │   │   ├── products.ts
│   │   │   │   └── customers.ts
│   │   │   ├── subscriptions/
│   │   │   │   ├── orders.ts
│   │   │   │   └── inventory.ts
│   │   │   └── index.ts
│   │   ├── codegen.yml             # GraphQL codegen config
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   ├── auth/                       # Authentication Package
│   │   ├── src/
│   │   │   ├── auth-provider.ts    # Refine auth provider
│   │   │   ├── access-control.ts   # RBAC/ABAC
│   │   │   ├── session.ts
│   │   │   ├── jwt.ts
│   │   │   └── index.ts
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   ├── types/                      # Shared TypeScript Types
│   │   ├── src/
│   │   │   ├── index.ts
│   │   │   ├── user.ts
│   │   │   ├── order.ts
│   │   │   ├── product.ts
│   │   │   ├── customer.ts
│   │   │   ├── inventory.ts
│   │   │   ├── payment.ts
│   │   │   ├── delivery.ts
│   │   │   └── analytics.ts
│   │   ├── package.json
│   │   └── tsconfig.json
│   │
│   └── config/                     # Shared Configurations
│       ├── eslint/
│       │   └── index.js
│       ├── typescript/
│       │   └── base.json
│       ├── tailwind/
│       │   └── preset.js
│       └── package.json
│
└── docker/
    └── docker-compose.yml          # Local dev stack
```

## TECHNOLOGY STACK

### Core Framework
- Next.js 14.2+ (App Router, Server Components, Server Actions)
- TypeScript 5.4+ (strict mode)
- Turborepo 2.0+ for monorepo management
- pnpm 8.15+ for package management

### Admin Framework (Admin, Bank, Partner portals)
- Refine.dev v4.52+ (headless admin framework)
- @refinedev/nextjs-router
- @refinedev/antd (Ant Design integration)
- @refinedev/kbar (command palette)

### UI Libraries
- Ant Design 5.15+ (Admin/Bank/Partner)
- shadcn/ui + Radix UI (Shop)
- Tailwind CSS 3.4+
- Framer Motion 11+ (animations)

### Data Layer
- GraphQL with graphql-request
- @tanstack/react-query v5
- graphql-ws for subscriptions
- Hasura as GraphQL engine

### Forms & Validation
- React Hook Form 7.50+
- Zod 3.22+ for validation
- @hookform/resolvers

### Charts & Visualization
- Recharts 2.12+ (primary)
- Apache ECharts (complex visualizations)

### State Management
- Zustand 4.5+ (client state)
- TanStack Query (server state)
- React Context (theme, auth)

### Authentication
- NextAuth.js 4.24+ (auth framework)
- JWT with Hasura claims

### Testing
- Vitest (unit tests)
- Playwright (E2E tests)
- MSW (API mocking)

## DESIGN SYSTEM

### Brand Colors
```css
:root {
  /* Primary - Deep Navy (Trust, Enterprise) */
  --primary-50: #E6EBF4;
  --primary-100: #C2D1E8;
  --primary-200: #9AB4DA;
  --primary-300: #7297CC;
  --primary-400: #5482C2;
  --primary-500: #366DB8;
  --primary-600: #3065B1;
  --primary-700: #295AA8;
  --primary-800: #2250A0;
  --primary-900: #1A365D;  /* Main brand color */
  
  /* Secondary - Success Green (Growth, Money) */
  --secondary-50: #E8F5E9;
  --secondary-100: #C8E6C9;
  --secondary-200: #A5D6A7;
  --secondary-300: #81C784;
  --secondary-400: #66BB6A;
  --secondary-500: #4CAF50;
  --secondary-600: #43A047;
  --secondary-700: #388E3C;
  --secondary-800: #2E7D32;
  --secondary-900: #1B5E20;
  
  /* Accent - Gold (Premium, Value) */
  --accent-50: #FFF8E1;
  --accent-100: #FFECB3;
  --accent-200: #FFE082;
  --accent-300: #FFD54F;
  --accent-400: #FFCA28;
  --accent-500: #FFC107;
  --accent-600: #FFB300;
  --accent-700: #FFA000;
  --accent-800: #FF8F00;
  --accent-900: #D69E2E;
  
  /* Semantic Colors */
  --success: #38A169;
  --warning: #D69E2E;
  --error: #E53E3E;
  --info: #3182CE;
  
  /* Neutrals */
  --gray-50: #F7FAFC;
  --gray-100: #EDF2F7;
  --gray-200: #E2E8F0;
  --gray-300: #CBD5E0;
  --gray-400: #A0AEC0;
  --gray-500: #718096;
  --gray-600: #4A5568;
  --gray-700: #2D3748;
  --gray-800: #1A202C;
  --gray-900: #171923;
}
```

### Typography
```css
:root {
  /* Display Font - Headers, Titles */
  --font-display: 'Plus Jakarta Sans', -apple-system, sans-serif;
  
  /* Body Font - Content, UI */
  --font-body: 'DM Sans', -apple-system, sans-serif;
  
  /* Mono Font - Numbers, Code */
  --font-mono: 'JetBrains Mono', 'SF Mono', monospace;
  
  /* Font Sizes */
  --text-xs: 0.75rem;    /* 12px */
  --text-sm: 0.875rem;   /* 14px */
  --text-base: 1rem;     /* 16px */
  --text-lg: 1.125rem;   /* 18px */
  --text-xl: 1.25rem;    /* 20px */
  --text-2xl: 1.5rem;    /* 24px */
  --text-3xl: 1.875rem;  /* 30px */
  --text-4xl: 2.25rem;   /* 36px */
  --text-5xl: 3rem;      /* 48px */
}
```

### Spacing & Layout
```css
:root {
  /* Spacing Scale */
  --space-1: 0.25rem;   /* 4px */
  --space-2: 0.5rem;    /* 8px */
  --space-3: 0.75rem;   /* 12px */
  --space-4: 1rem;      /* 16px */
  --space-5: 1.25rem;   /* 20px */
  --space-6: 1.5rem;    /* 24px */
  --space-8: 2rem;      /* 32px */
  --space-10: 2.5rem;   /* 40px */
  --space-12: 3rem;     /* 48px */
  --space-16: 4rem;     /* 64px */
  
  /* Border Radius */
  --radius-sm: 4px;
  --radius-md: 6px;
  --radius-lg: 8px;
  --radius-xl: 12px;
  --radius-2xl: 16px;
  --radius-full: 9999px;
  
  /* Shadows */
  --shadow-xs: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06);
  --shadow-md: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  --shadow-lg: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}
```

## SHARED TYPES (packages/types/src/index.ts)

Generate comprehensive TypeScript types including:

### Core Types
- Money (amount, currency, formatted)
- Address (Nigerian format with LGA support)
- Phone (Nigerian format validation)
- Pagination, Filters, Sorting

### User & Auth Types
- User, UserRole, Permission
- Session, HasuraClaims
- Tenant, TenantSettings

### Commerce Types
- Order, OrderItem, OrderStatus, OrderSource
- Product, ProductVariant, ProductOption
- Category, Brand
- Inventory, InventoryLocation, InventoryMovement

### Customer Types
- Customer, CustomerType, CustomerTier
- Address (shipping/billing)
- CreditFacility, PaymentMethod

### Financial Types
- Payment, PaymentStatus, PaymentMethod
- Transaction, TransactionType
- Settlement, Payout

### Delivery Types
- Fulfillment, FulfillmentItem
- Shipment, ShipmentStatus
- GigWorker, Trip

## PORTAL-SPECIFIC REQUIREMENTS

### 1. ADMIN PORTAL (apps/admin/)

Purpose: Platform administration, multi-tenant management, system oversight

Key Features:
- Dashboard with GMV, orders, users, revenue metrics
- Real-time activity feed and system health
- Tenant CRUD with subscription management
- User management with role assignment
- Order management across all tenants
- Product catalog (master products)
- Financial operations (settlements, payouts)
- Analytics & custom reports
- Audit logs & compliance

Tech Stack:
- Refine.dev with Ant Design
- Full CRUD for all resources
- Real-time updates via subscriptions
- Export functionality (CSV, Excel, PDF)

### 2. BANK PORTAL (apps/bank/)

Purpose: Financial institution view for lending, collections, settlements

Key Features:
- Loan portfolio dashboard
- Loan application queue (Kanban + Table view)
- Credit assessment integration
- Disbursement tracking
- ATC (Authority to Collect) mandate management
- Collection calendar with scheduling
- Failed collection queue with retry
- Settlement queue and reconciliation
- KYC verification workflow
- AML alerts management
- Regulatory compliance reports

Specific Components:
- LoanApplicationsKanban with drag-drop status changes
- CreditScoreCard showing risk assessment
- RepaymentScheduleTimeline
- CollectionCalendar with daily/weekly views
- ReconciliationDiffView

### 3. PARTNER PORTAL (apps/partner/)

Purpose: B2B portal for supply chain participants

Role-Specific Dashboards:
1. Manufacturer: Production orders, inventory, distribution
2. Distributor: Route sales, credit management, deliveries
3. Wholesaler: Bulk orders, supplier management
4. Retailer: Ordering, POS integration, restock alerts
5. Logistics: Fleet management, trips, POD
6. Warehouse: Inventory, inbound/outbound, 3PL

Shared Features:
- Order management (create, receive, track)
- Inventory management with transfers
- Product catalog with custom pricing
- Customer/Supplier management
- Financial dashboard (invoices, payments, credit)
- Analytics & reporting

### 4. SHOP PORTAL (apps/shop/)

Purpose: B2C consumer-facing e-commerce

Key Features:
- Homepage with featured products, categories, promotions
- Product listing with filters (category, brand, price)
- Product detail with variants, images, reviews
- Shopping cart with persistent state
- Checkout flow (address, payment, confirmation)
- Order tracking
- Account management (orders, addresses, wishlist)
- Search with autocomplete

Design Requirements:
- Mobile-first (primary mobile users)
- Fast loading (< 2s LCP)
- Offline cart persistence
- PWA-ready

## KEY IMPLEMENTATIONS

### 1. Authentication Flow

```typescript
// packages/auth/src/auth-provider.ts
export const authProvider: AuthProvider = {
  login: async ({ email, password }) => {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    const { user, accessToken, refreshToken } = await response.json();
    
    // Store tokens
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('refreshToken', refreshToken);
    
    return { success: true, redirectTo: '/' };
  },
  
  logout: async () => {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    return { success: true, redirectTo: '/login' };
  },
  
  check: async () => {
    const token = localStorage.getItem('accessToken');
    if (!token) return { authenticated: false };
    
    // Verify token validity
    try {
      const decoded = jwtDecode(token);
      if (decoded.exp * 1000 < Date.now()) {
        // Token expired, try refresh
        await refreshAccessToken();
      }
      return { authenticated: true };
    } catch {
      return { authenticated: false };
    }
  },
  
  getIdentity: async () => {
    const token = localStorage.getItem('accessToken');
    if (!token) return null;
    
    const decoded = jwtDecode<HasuraClaims>(token);
    return {
      id: decoded['x-hasura-user-id'],
      name: decoded.name,
      email: decoded.email,
      role: decoded['x-hasura-default-role'],
    };
  },
};
```

### 2. Hasura Data Provider

```typescript
// packages/api/src/data-provider.ts
export const createDataProvider = (client: GraphQLClient): DataProvider => ({
  getList: async ({ resource, pagination, filters, sorters, meta }) => {
    const { current = 1, pageSize = 10 } = pagination ?? {};
    const offset = (current - 1) * pageSize;
    
    const where = buildWhereClause(filters);
    const orderBy = buildOrderByClause(sorters);
    const fields = meta?.fields ?? 'id';
    
    const query = gql`
      query GetList($limit: Int!, $offset: Int!, $where: ${resource}_bool_exp, $orderBy: [${resource}_order_by!]) {
        ${resource}(limit: $limit, offset: $offset, where: $where, order_by: $orderBy) {
          ${fields}
        }
        ${resource}_aggregate(where: $where) {
          aggregate { count }
        }
      }
    `;
    
    const result = await client.request(query, {
      limit: pageSize,
      offset,
      where,
      orderBy,
    });
    
    return {
      data: result[resource],
      total: result[`${resource}_aggregate`].aggregate.count,
    };
  },
  
  // Implement all other methods: getOne, create, update, deleteOne, getMany, custom
});
```

### 3. Real-time Subscriptions

```typescript
// packages/api/src/live-provider.ts
export const createLiveProvider = (wsClient: Client): LiveProvider => ({
  subscribe: ({ channel, types, callback, params }) => {
    const query = gql`
      subscription Subscribe${channel} {
        ${channel}(where: ${JSON.stringify(params?.where ?? {})}) {
          id
          __typename
        }
      }
    `;
    
    const unsubscribe = wsClient.subscribe({ query }, {
      next: ({ data }) => {
        if (data) {
          callback({
            type: types[0],
            channel,
            date: new Date(),
            payload: data[channel],
          });
        }
      },
      error: console.error,
      complete: () => {},
    });
    
    return unsubscribe;
  },
  
  unsubscribe: (unsubscribe) => {
    unsubscribe();
  },
});
```

### 4. Access Control (RBAC)

```typescript
// packages/auth/src/access-control.ts
const rolePermissions: Record<UserRole, Permission[]> = {
  super_admin: ['*'],
  platform_admin: [
    'orders:*',
    'products:*',
    'customers:*',
    'users:read',
    'users:create',
    'analytics:*',
  ],
  bank_admin: [
    'loans:*',
    'atc:*',
    'settlements:*',
    'customers:read',
    'compliance:*',
  ],
  manufacturer: [
    'orders:read',
    'orders:create',
    'products:*',
    'inventory:*',
    'customers:read',
  ],
  retailer: [
    'orders:*',
    'products:read',
    'inventory:read',
    'customers:read',
  ],
  // ... other roles
};

export const accessControlProvider: AccessControlProvider = {
  can: async ({ resource, action, params }) => {
    const user = await getUser();
    if (!user) return { can: false };
    
    const permissions = rolePermissions[user.role] || [];
    
    // Check wildcard
    if (permissions.includes('*')) return { can: true };
    
    // Check resource:action
    const permission = `${resource}:${action}`;
    const resourceWildcard = `${resource}:*`;
    
    const can = permissions.includes(permission) || 
                permissions.includes(resourceWildcard);
    
    return { can };
  },
};
```

## COMPONENTS TO IMPLEMENT

### Dashboard Components
- StatCard (with trend indicator, icon, sparkline)
- RevenueChart (area chart with period selector)
- OrdersChart (bar chart by status)
- TopProductsTable
- ActivityFeed (real-time updates)
- SystemHealthIndicator

### Data Table Components
- DataTable (with sorting, filtering, pagination)
- FilterBar (with saved filters)
- BulkActions (for selected rows)
- ExportButton (CSV, Excel, PDF)
- ColumnCustomizer

### Form Components
- MoneyInput (with currency selector, formatting)
- PhoneInput (with country code, validation)
- AddressInput (with Nigerian states/LGAs)
- ProductSelector (with search, variants)
- CustomerSelector (with search, create inline)
- DateRangePicker (with presets)

### Status & Display Components
- StatusBadge (order, payment, fulfillment status)
- PriorityIndicator
- ProgressBar (for fulfillment, payment)
- Timeline (order events, activity)
- AvatarGroup (for team members)

## IMPLEMENTATION CHECKLIST

Phase 1: Foundation
- [ ] Turborepo monorepo setup
- [ ] Shared packages (types, ui, api, auth, config)
- [ ] Design system with CSS variables
- [ ] Hasura client configuration
- [ ] Authentication setup

Phase 2: Admin Portal
- [ ] Layout (sidebar, header, breadcrumb)
- [ ] Dashboard with all widgets
- [ ] Orders CRUD with filters
- [ ] Products CRUD with variants
- [ ] Customers CRUD
- [ ] Users management
- [ ] Settings pages

Phase 3: Bank Portal
- [ ] Loan management
- [ ] ATC/Collections
- [ ] Settlements
- [ ] Compliance

Phase 4: Partner Portal
- [ ] Role-based routing
- [ ] 6 dashboard variants
- [ ] Orders (buy/sell)
- [ ] Inventory
- [ ] Finance

Phase 5: Shop Portal
- [ ] Homepage
- [ ] Product listing/detail
- [ ] Cart & checkout
- [ ] Account pages

Phase 6: Polish
- [ ] Loading states
- [ ] Error boundaries
- [ ] Empty states
- [ ] Animations
- [ ] Mobile optimization
- [ ] PWA features

## QUALITY REQUIREMENTS

1. Performance:
   - Lighthouse score > 90
   - LCP < 2.5s
   - FID < 100ms
   - CLS < 0.1

2. Accessibility:
   - WCAG 2.1 AA compliance
   - Keyboard navigation
   - Screen reader support

3. Code Quality:
   - TypeScript strict mode
   - ESLint + Prettier
   - No any types
   - Comprehensive error handling

4. Testing:
   - Unit tests for utilities
   - Component tests
   - E2E for critical paths

Generate all files with complete, production-ready code. Include all imports, types, and implementations. This is not a prototype - it's the real application.
```

---

## QUICK START AFTER GENERATION

```bash
# Navigate to project
cd omniroute-web

# Install dependencies
pnpm install

# Set up environment
cp .env.example .env.local
# Edit .env.local with your Hasura endpoint, auth secrets, etc.

# Start development
pnpm dev

# Individual portals
pnpm dev --filter=admin
pnpm dev --filter=bank
pnpm dev --filter=partner
pnpm dev --filter=shop

# Build
pnpm build

# Test
pnpm test
```

## EXPECTED OUTPUT

This prompt should generate:
- ~200+ React components
- ~80+ pages across 4 portals
- ~50+ custom hooks
- Complete GraphQL operations
- Full TypeScript types
- ~35,000+ lines of production code

---

## NOTES FOR CLAUDE CODE

1. Generate ALL files - don't skip any
2. Use real implementations, not placeholders
3. Include proper error handling everywhere
4. Add loading and empty states
5. Make everything responsive
6. Follow Nigerian market conventions (₦, phone format, states)
7. Include proper TypeScript types for everything
