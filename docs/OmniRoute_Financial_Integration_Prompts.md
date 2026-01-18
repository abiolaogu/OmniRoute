# OmniRoute Financial Integration & Infrastructure Modernization Prompts

## Overview

This document contains implementation prompts for:
1. **Bank & Financial Institution Integration** via Hasura GraphQL API
2. **Authority to Collect (ATC) System** for B2B distribution
3. **Infrastructure Migration** to YugabyteDB, DragonflyDB, and Redpanda

---

# PART 1: BANK & FINANCIAL INSTITUTION INTEGRATION

## F-P01: Hasura GraphQL Bank Integration Gateway (Go + Hasura)

```
PROJECT: OmniRoute - Bank Integration Gateway with Hasura

CONTEXT:
Build a comprehensive Bank and Financial Institution Integration Gateway using Hasura GraphQL API to coordinate payments across the distribution ecosystem. This system enables:

- Real-time payment initiation and tracking
- Multi-bank connectivity (commercial banks, microfinance, mobile money)
- Payment orchestration across multiple financial rails
- Account validation and verification
- Balance inquiries and statement retrieval
- Bulk payment processing for supplier settlements
- Direct debit mandate management
- Standing order configuration
- Webhook-based payment notifications
- Reconciliation automation
- Regulatory reporting (CBN, NDIC compliance for Nigeria)

TECHNICAL REQUIREMENTS:
- API Layer: Hasura GraphQL Engine v2.x
- Backend: Go 1.22+ for custom business logic (Hasura Actions)
- Database: YugabyteDB (PostgreSQL-compatible, distributed)
- Cache: DragonflyDB (Redis-compatible)
- Events: Redpanda (Kafka-compatible)
- Security: mTLS for bank connections, HSM for key management

SUPPORTED FINANCIAL RAILS:
- NIBSS (Nigeria Inter-Bank Settlement System)
  - NIP: NIBSS Instant Payment
  - NEFT: NIBSS Electronic Fund Transfer
  - Direct Debit
- RTGS (Real-Time Gross Settlement)
- Mobile Money APIs (M-Pesa, MTN MoMo, Airtel Money)
- Card Networks (Visa, Mastercard via payment processors)
- Open Banking APIs (PSD2-style where available)

HASURA SCHEMA DESIGN:
```graphql
# Bank Connection Configuration
type BankConnection {
  id: uuid!
  tenant_id: uuid!
  bank_code: String!
  bank_name: String!
  connection_type: BankConnectionType!
  api_endpoint: String!
  credentials: jsonb  # Encrypted
  status: ConnectionStatus!
  last_health_check: timestamptz
  supported_operations: [BankOperation!]!
  rate_limits: jsonb
  created_at: timestamptz!
  updated_at: timestamptz!
}

enum BankConnectionType {
  NIBSS_NIP
  NIBSS_NEFT
  RTGS
  MOBILE_MONEY
  OPEN_BANKING
  CARD_PROCESSOR
}

# Virtual Account for Collections
type VirtualAccount {
  id: uuid!
  tenant_id: uuid!
  account_number: String! @unique
  bank_code: String!
  account_name: String!
  owner_id: uuid!
  owner_type: AccountOwnerType!
  currency: String!
  balance: numeric!
  status: AccountStatus!
  expiry_date: timestamptz
  collection_rules: jsonb
  metadata: jsonb
  created_at: timestamptz!
}

enum AccountOwnerType {
  CUSTOMER
  SUPPLIER
  WORKER
  GROUP
  ESCROW
}

# Payment Transaction
type PaymentTransaction {
  id: uuid!
  tenant_id: uuid!
  reference: String! @unique
  external_reference: String
  type: PaymentType!
  direction: PaymentDirection!
  rail: PaymentRail!
  source_account: AccountDetails!
  destination_account: AccountDetails!
  amount: numeric!
  currency: String!
  fee: numeric!
  status: PaymentStatus!
  status_history: [PaymentStatusChange!]!
  initiated_by: uuid!
  approved_by: uuid
  narration: String
  metadata: jsonb
  bank_response: jsonb
  created_at: timestamptz!
  completed_at: timestamptz
}

enum PaymentType {
  SUPPLIER_PAYMENT
  WORKER_PAYOUT
  CUSTOMER_REFUND
  COLLECTION
  INTERNAL_TRANSFER
  BULK_DISBURSEMENT
}

enum PaymentStatus {
  PENDING
  AWAITING_APPROVAL
  PROCESSING
  SENT_TO_BANK
  CONFIRMED
  FAILED
  REVERSED
  EXPIRED
}

# Bulk Payment Batch
type BulkPaymentBatch {
  id: uuid!
  tenant_id: uuid!
  reference: String!
  name: String!
  total_amount: numeric!
  total_count: Int!
  successful_count: Int!
  failed_count: Int!
  status: BatchStatus!
  payments: [PaymentTransaction!]! @relation
  initiated_by: uuid!
  approval_workflow: uuid
  created_at: timestamptz!
  completed_at: timestamptz
}

# Bank Statement
type BankStatement {
  id: uuid!
  tenant_id: uuid!
  account_id: uuid!
  statement_date: date!
  opening_balance: numeric!
  closing_balance: numeric!
  total_credits: numeric!
  total_debits: numeric!
  entries: [StatementEntry!]!
  reconciliation_status: ReconciliationStatus!
  created_at: timestamptz!
}
```

HASURA ACTIONS (Go Backend):
```go
// cmd/bank-gateway/main.go
package main

import (
    "github.com/omniroute/bank-gateway/internal/handler"
    "github.com/omniroute/bank-gateway/internal/service"
    "github.com/omniroute/bank-gateway/internal/bank"
)

func main() {
    // Initialize bank connectors
    nibssConnector := bank.NewNIBSSConnector(config.NIBSS)
    mobileMoneyConnector := bank.NewMobileMoneyConnector(config.MobileMoney)
    
    // Initialize services
    paymentService := service.NewPaymentService(
        nibssConnector,
        mobileMoneyConnector,
        yugabyteDB,
        dragonflyCache,
        redpandaProducer,
    )
    
    // Register Hasura Action handlers
    router := chi.NewRouter()
    router.Post("/actions/initiate-payment", handler.InitiatePayment(paymentService))
    router.Post("/actions/verify-account", handler.VerifyAccount(paymentService))
    router.Post("/actions/process-bulk-payment", handler.ProcessBulkPayment(paymentService))
    router.Post("/actions/create-virtual-account", handler.CreateVirtualAccount(paymentService))
    router.Post("/actions/fetch-statement", handler.FetchStatement(paymentService))
    
    // Webhook handlers for bank callbacks
    router.Post("/webhooks/nibss/payment-notification", handler.NIBSSPaymentWebhook(paymentService))
    router.Post("/webhooks/mobile-money/callback", handler.MobileMoneyWebhook(paymentService))
}

// internal/service/payment_service.go
type PaymentService struct {
    nibss         *bank.NIBSSConnector
    mobileMoney   *bank.MobileMoneyConnector
    db            *yugabyte.Client
    cache         *dragonfly.Client
    events        *redpanda.Producer
}

func (s *PaymentService) InitiatePayment(ctx context.Context, 
    req *InitiatePaymentRequest) (*PaymentTransaction, error) {
    // 1. Validate request and check limits
    if err := s.validatePaymentRequest(ctx, req); err != nil {
        return nil, err
    }
    
    // 2. Check for duplicates (idempotency)
    existing, _ := s.cache.Get(ctx, "payment:ref:"+req.Reference)
    if existing != "" {
        return s.getTransaction(ctx, existing)
    }
    
    // 3. Create transaction record
    tx := &PaymentTransaction{
        ID:            uuid.New(),
        TenantID:      req.TenantID,
        Reference:     req.Reference,
        Type:          req.Type,
        Direction:     PaymentDirectionOutbound,
        Rail:          s.selectOptimalRail(req),
        Amount:        req.Amount,
        Currency:      req.Currency,
        Status:        PaymentStatusPending,
        SourceAccount: req.SourceAccount,
        DestAccount:   req.DestinationAccount,
    }
    
    // 4. Check if approval required
    if s.requiresApproval(ctx, tx) {
        tx.Status = PaymentStatusAwaitingApproval
        return s.saveAndNotify(ctx, tx)
    }
    
    // 5. Send to bank
    return s.sendToBank(ctx, tx)
}

func (s *PaymentService) sendToBank(ctx context.Context, 
    tx *PaymentTransaction) (*PaymentTransaction, error) {
    tx.Status = PaymentStatusSentToBank
    s.db.Save(ctx, tx)
    
    var result *BankResponse
    var err error
    
    switch tx.Rail {
    case PaymentRailNIBSSNIP:
        result, err = s.nibss.SendNIPTransfer(ctx, &NIBSSTransferRequest{
            SessionID:        tx.ID.String(),
            Amount:           tx.Amount,
            SourceBankCode:   tx.SourceAccount.BankCode,
            SourceAccount:    tx.SourceAccount.AccountNumber,
            DestBankCode:     tx.DestAccount.BankCode,
            DestAccount:      tx.DestAccount.AccountNumber,
            Narration:        tx.Narration,
        })
    case PaymentRailMobileMoney:
        result, err = s.mobileMoney.Disburse(ctx, &MobileMoneyRequest{
            Reference:   tx.Reference,
            Amount:      tx.Amount,
            Phone:       tx.DestAccount.Phone,
            Provider:    tx.DestAccount.Provider,
        })
    }
    
    if err != nil {
        tx.Status = PaymentStatusFailed
        tx.BankResponse = map[string]interface{}{"error": err.Error()}
    } else {
        tx.ExternalReference = result.TransactionID
        tx.Status = PaymentStatusConfirmed
        tx.CompletedAt = time.Now()
    }
    
    s.db.Save(ctx, tx)
    s.events.Publish(ctx, "payments.completed", tx)
    
    return tx, nil
}

// internal/bank/nibss_connector.go
type NIBSSConnector struct {
    client     *http.Client
    baseURL    string
    bankCode   string
    secretKey  []byte
    aesKey     []byte
    ivKey      []byte
}

func (c *NIBSSConnector) SendNIPTransfer(ctx context.Context,
    req *NIBSSTransferRequest) (*BankResponse, error) {
    // Build NIP request XML
    // Encrypt with AES
    // Sign with RSA
    // Send to NIBSS endpoint
    // Parse and return response
}

func (c *NIBSSConnector) VerifyAccount(ctx context.Context,
    bankCode, accountNumber string) (*AccountVerification, error) {
    // Call NIBSS Name Enquiry API
    // Return account name and status
}
```

HASURA METADATA CONFIGURATION:
```yaml
# hasura/metadata/actions.yaml
actions:
  - name: initiatePayment
    definition:
      kind: synchronous
      handler: '{{BANK_GATEWAY_URL}}/actions/initiate-payment'
      timeout: 30
      request_transform:
        template_engine: Kriti
    permissions:
      - role: finance_admin
      - role: payment_initiator

  - name: verifyBankAccount
    definition:
      kind: synchronous
      handler: '{{BANK_GATEWAY_URL}}/actions/verify-account'
      timeout: 10
    permissions:
      - role: user

  - name: processBulkPayment
    definition:
      kind: asynchronous
      handler: '{{BANK_GATEWAY_URL}}/actions/process-bulk-payment'
    permissions:
      - role: finance_admin

  - name: createVirtualAccount
    definition:
      kind: synchronous
      handler: '{{BANK_GATEWAY_URL}}/actions/create-virtual-account'
    permissions:
      - role: finance_admin
      - role: customer_admin

# hasura/metadata/event_triggers.yaml
event_triggers:
  - name: payment_status_changed
    table:
      schema: public
      name: payment_transactions
    definition:
      update:
        columns:
          - status
    webhook: '{{NOTIFICATION_SERVICE_URL}}/webhooks/payment-status'
    
  - name: virtual_account_credited
    table:
      schema: public
      name: virtual_account_transactions
    definition:
      insert:
        columns: '*'
    webhook: '{{ORDER_SERVICE_URL}}/webhooks/payment-received'
```

GRAPHQL QUERIES AND MUTATIONS:
```graphql
# Payment Operations
mutation InitiatePayment($input: InitiatePaymentInput!) {
  initiatePayment(input: $input) {
    id
    reference
    status
    externalReference
  }
}

mutation ProcessBulkPayment($batchId: uuid!, $payments: [PaymentInput!]!) {
  processBulkPayment(batchId: $batchId, payments: $payments) {
    batchId
    totalProcessed
    successCount
    failedCount
    failures {
      reference
      error
    }
  }
}

query GetPaymentStatus($reference: String!) {
  payment_transactions(where: {reference: {_eq: $reference}}) {
    id
    reference
    status
    amount
    currency
    statusHistory {
      status
      timestamp
      details
    }
    bankResponse
  }
}

# Virtual Account Operations
mutation CreateVirtualAccount($input: CreateVirtualAccountInput!) {
  createVirtualAccount(input: $input) {
    accountNumber
    bankCode
    bankName
    accountName
    expiryDate
  }
}

subscription VirtualAccountTransactions($accountId: uuid!) {
  virtual_account_transactions(
    where: {account_id: {_eq: $accountId}}
    order_by: {created_at: desc}
    limit: 50
  ) {
    id
    amount
    type
    reference
    narration
    balanceAfter
    createdAt
  }
}

# Reconciliation
query GetReconciliationReport($startDate: date!, $endDate: date!) {
  reconciliation_summary(args: {start_date: $startDate, end_date: $endDate}) {
    totalTransactions
    totalAmount
    matchedCount
    unmatchedCount
    discrepancies {
      transactionId
      expectedAmount
      actualAmount
      difference
    }
  }
}
```

TEST CASES:
1. Single payment initiation and confirmation flow
2. Bulk payment batch processing (1000+ payments)
3. Virtual account creation and credit notification
4. NIBSS NIP transfer success and failure handling
5. Mobile money disbursement across providers
6. Payment approval workflow execution
7. Duplicate payment detection (idempotency)
8. Bank statement fetch and reconciliation
9. Webhook reliability and retry handling
10. Multi-tenant isolation verification

EXPECTED DELIVERABLES:
1. Hasura metadata and migrations
2. Go backend for Hasura Actions
3. NIBSS connector implementation
4. Mobile money connector (M-Pesa, MTN, Airtel)
5. Virtual account management
6. Reconciliation engine
7. GraphQL subscriptions for real-time updates
8. Comprehensive test suite
9. API documentation
10. Deployment manifests
```

---

## F-P02: Authority to Collect (ATC) System (Go + Hasura)

```
PROJECT: OmniRoute - Authority to Collect (ATC) System

CONTEXT:
Build an Authority to Collect (ATC) system that enables distributors, wholesalers, and manufacturers to delegate collection rights to downstream partners in the FMCG distribution chain. This follows standard distribution norms where:

- Manufacturers grant ATC to Distributors
- Distributors grant ATC to Sub-distributors or Wholesalers
- Wholesalers grant ATC to Retailers (in some cases)
- ATC enables automatic payment routing and settlement
- Supports tiered commission structures
- Handles payment splitting and cascading settlements

ATC USE CASES IN DISTRIBUTION:
1. **Manufacturer → Distributor ATC**: Distributor collects payments from retailers on behalf of manufacturer, keeps agreed margin
2. **Distributor → Van Sales ATC**: Van sales reps collect cash/mobile money, remit to distributor
3. **Wholesaler → Retailer Credit ATC**: Wholesaler authorizes collection from retailer's future sales
4. **Platform → All Parties**: OmniRoute acts as settlement orchestrator

TECHNICAL REQUIREMENTS:
- API Layer: Hasura GraphQL Engine v2.x
- Backend: Go 1.22+ for ATC business logic
- Database: YugabyteDB (distributed PostgreSQL)
- Cache: DragonflyDB
- Events: Redpanda
- Workflow: Temporal for complex settlement flows

DOMAIN MODEL:
```graphql
# Authority to Collect Grant
type ATCGrant {
  id: uuid!
  tenant_id: uuid!
  reference: String! @unique
  
  # Grantor (who gives authority)
  grantor_id: uuid!
  grantor_type: PartyType!
  grantor_name: String!
  
  # Grantee (who receives authority)
  grantee_id: uuid!
  grantee_type: PartyType!
  grantee_name: String!
  
  # Scope of authority
  scope: ATCScope!
  collection_type: CollectionType!
  
  # Financial terms
  currency: String!
  max_amount: numeric              # Per-transaction limit
  cumulative_limit: numeric        # Total collection limit
  cumulative_collected: numeric!   # Running total
  
  # Commission/Margin structure
  commission_type: CommissionType!
  commission_rate: numeric         # Percentage
  commission_flat: numeric         # Fixed amount
  min_commission: numeric
  max_commission: numeric
  
  # Settlement terms
  settlement_frequency: SettlementFrequency!
  settlement_day: Int              # Day of week/month
  settlement_delay_days: Int!      # Days after collection
  settlement_account: AccountDetails!
  
  # Validity
  status: ATCStatus!
  effective_from: timestamptz!
  effective_to: timestamptz
  
  # Hierarchy
  parent_atc_id: uuid              # For cascading ATCs
  child_atcs: [ATCGrant!]! @relation
  
  # Metadata
  terms_document_url: String
  approved_by: uuid
  approved_at: timestamptz
  created_at: timestamptz!
  updated_at: timestamptz!
}

enum PartyType {
  MANUFACTURER
  DISTRIBUTOR
  WHOLESALER
  RETAILER
  WORKER
  PLATFORM
}

enum ATCScope {
  ALL_PRODUCTS           # Can collect for any product
  PRODUCT_CATEGORY       # Specific categories only
  SPECIFIC_PRODUCTS      # Listed products only
  SPECIFIC_CUSTOMERS     # Listed customers only
  GEOGRAPHIC_AREA        # Within defined territory
}

enum CollectionType {
  CASH
  MOBILE_MONEY
  BANK_TRANSFER
  CHEQUE
  CREDIT_AGAINST_PURCHASE
  ALL_METHODS
}

enum CommissionType {
  PERCENTAGE
  FLAT_FEE
  TIERED_PERCENTAGE      # Different rates for volume bands
  HYBRID                 # Flat + Percentage
}

enum SettlementFrequency {
  INSTANT                # Real-time settlement
  DAILY
  WEEKLY
  BIWEEKLY
  MONTHLY
  ON_DEMAND
}

enum ATCStatus {
  DRAFT
  PENDING_APPROVAL
  ACTIVE
  SUSPENDED
  EXPIRED
  REVOKED
}

# Collection Record
type ATCCollection {
  id: uuid!
  tenant_id: uuid!
  atc_grant_id: uuid!
  atc_grant: ATCGrant! @relation
  
  # Collection details
  reference: String! @unique
  collected_from_id: uuid!
  collected_from_type: PartyType!
  collected_from_name: String!
  
  collection_method: CollectionType!
  payment_reference: String         # Link to payment transaction
  
  gross_amount: numeric!
  commission_amount: numeric!
  net_amount: numeric!              # Amount to settle to grantor
  
  # Related entities
  order_id: uuid
  invoice_id: uuid
  
  # Settlement tracking
  settlement_status: SettlementStatus!
  settlement_batch_id: uuid
  settled_at: timestamptz
  
  # Metadata
  collection_point: GeoPoint
  collected_by: uuid                # Worker who collected
  notes: String
  
  created_at: timestamptz!
}

enum SettlementStatus {
  PENDING
  SCHEDULED
  PROCESSING
  SETTLED
  FAILED
  DISPUTED
}

# Settlement Batch
type ATCSettlementBatch {
  id: uuid!
  tenant_id: uuid!
  reference: String!
  
  atc_grant_id: uuid!
  settlement_period_start: timestamptz!
  settlement_period_end: timestamptz!
  
  total_collections: Int!
  gross_amount: numeric!
  total_commission: numeric!
  net_settlement_amount: numeric!
  
  collections: [ATCCollection!]! @relation
  
  # Payment details
  payment_transaction_id: uuid
  payment_status: PaymentStatus!
  
  # Reconciliation
  reconciliation_status: ReconciliationStatus!
  discrepancy_amount: numeric
  discrepancy_notes: String
  
  created_at: timestamptz!
  settled_at: timestamptz
}

# ATC Scope Configuration
type ATCScopeConfig {
  id: uuid!
  atc_grant_id: uuid!
  scope_type: ATCScopeType!
  scope_values: [String!]!          # Product IDs, Category IDs, etc.
  created_at: timestamptz!
}

enum ATCScopeType {
  PRODUCT_ID
  PRODUCT_CATEGORY
  CUSTOMER_ID
  CUSTOMER_TIER
  GEOGRAPHIC_ZONE
  ORDER_TYPE
}

# Commission Tier (for tiered commission structures)
type ATCCommissionTier {
  id: uuid!
  atc_grant_id: uuid!
  tier_number: Int!
  min_amount: numeric!
  max_amount: numeric
  commission_rate: numeric!
  created_at: timestamptz!
}
```

GO BACKEND IMPLEMENTATION:
```go
// internal/domain/atc.go
package domain

type ATCGrant struct {
    ID                   uuid.UUID
    TenantID             uuid.UUID
    Reference            string
    GrantorID            uuid.UUID
    GrantorType          PartyType
    GranteeID            uuid.UUID
    GranteeType          PartyType
    Scope                ATCScope
    CollectionType       CollectionType
    Currency             string
    MaxAmount            decimal.Decimal
    CumulativeLimit      decimal.Decimal
    CumulativeCollected  decimal.Decimal
    CommissionType       CommissionType
    CommissionRate       decimal.Decimal
    CommissionTiers      []CommissionTier
    SettlementFrequency  SettlementFrequency
    SettlementDelayDays  int
    SettlementAccount    AccountDetails
    Status               ATCStatus
    EffectiveFrom        time.Time
    EffectiveTo          *time.Time
    ParentATCID          *uuid.UUID
}

// internal/service/atc_service.go
type ATCService struct {
    db              *yugabyte.Client
    cache           *dragonfly.Client
    events          *redpanda.Producer
    paymentService  *PaymentService
    workflowClient  client.Client  // Temporal
}

func (s *ATCService) CreateATCGrant(ctx context.Context, 
    req *CreateATCRequest) (*ATCGrant, error) {
    // 1. Validate grantor has authority to grant
    if err := s.validateGrantorAuthority(ctx, req); err != nil {
        return nil, fmt.Errorf("grantor validation failed: %w", err)
    }
    
    // 2. Check for conflicting ATCs
    if conflict := s.checkATCConflict(ctx, req); conflict != nil {
        return nil, fmt.Errorf("conflicting ATC exists: %s", conflict.Reference)
    }
    
    // 3. Validate scope configuration
    if err := s.validateATCScope(ctx, req.Scope, req.ScopeConfig); err != nil {
        return nil, fmt.Errorf("invalid scope: %w", err)
    }
    
    // 4. Create ATC grant
    atc := &ATCGrant{
        ID:                  uuid.New(),
        TenantID:            req.TenantID,
        Reference:           s.generateATCReference(),
        GrantorID:           req.GrantorID,
        GrantorType:         req.GrantorType,
        GranteeID:           req.GranteeID,
        GranteeType:         req.GranteeType,
        Scope:               req.Scope,
        CollectionType:      req.CollectionType,
        Currency:            req.Currency,
        MaxAmount:           req.MaxAmount,
        CumulativeLimit:     req.CumulativeLimit,
        CumulativeCollected: decimal.Zero,
        CommissionType:      req.CommissionType,
        CommissionRate:      req.CommissionRate,
        CommissionTiers:     req.CommissionTiers,
        SettlementFrequency: req.SettlementFrequency,
        SettlementDelayDays: req.SettlementDelayDays,
        SettlementAccount:   req.SettlementAccount,
        Status:              ATCStatusPendingApproval,
        EffectiveFrom:       req.EffectiveFrom,
        EffectiveTo:         req.EffectiveTo,
        ParentATCID:         req.ParentATCID,
    }
    
    // 5. Save to database
    if err := s.db.Create(ctx, atc); err != nil {
        return nil, err
    }
    
    // 6. Create scope configurations
    for _, scope := range req.ScopeConfig {
        if err := s.db.Create(ctx, &ATCScopeConfig{
            ID:         uuid.New(),
            ATCGrantID: atc.ID,
            ScopeType:  scope.Type,
            ScopeValues: scope.Values,
        }); err != nil {
            return nil, err
        }
    }
    
    // 7. Publish event
    s.events.Publish(ctx, "atc.created", atc)
    
    return atc, nil
}

func (s *ATCService) RecordCollection(ctx context.Context,
    req *RecordCollectionRequest) (*ATCCollection, error) {
    // 1. Find applicable ATC grant
    atc, err := s.findApplicableATC(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("no valid ATC found: %w", err)
    }
    
    // 2. Validate collection against ATC limits
    if err := s.validateCollectionLimits(ctx, atc, req.Amount); err != nil {
        return nil, err
    }
    
    // 3. Calculate commission
    commission := s.calculateCommission(atc, req.Amount)
    netAmount := req.Amount.Sub(commission)
    
    // 4. Create collection record
    collection := &ATCCollection{
        ID:                uuid.New(),
        TenantID:          req.TenantID,
        ATCGrantID:        atc.ID,
        Reference:         s.generateCollectionReference(),
        CollectedFromID:   req.CollectedFromID,
        CollectedFromType: req.CollectedFromType,
        CollectedFromName: req.CollectedFromName,
        CollectionMethod:  req.CollectionMethod,
        PaymentReference:  req.PaymentReference,
        GrossAmount:       req.Amount,
        CommissionAmount:  commission,
        NetAmount:         netAmount,
        OrderID:           req.OrderID,
        InvoiceID:         req.InvoiceID,
        SettlementStatus:  SettlementStatusPending,
        CollectionPoint:   req.CollectionPoint,
        CollectedBy:       req.CollectedBy,
    }
    
    // 5. Update cumulative collected in ATC
    atc.CumulativeCollected = atc.CumulativeCollected.Add(req.Amount)
    
    // 6. Save in transaction
    err = s.db.Transaction(ctx, func(tx *yugabyte.Tx) error {
        if err := tx.Create(ctx, collection); err != nil {
            return err
        }
        if err := tx.Update(ctx, atc); err != nil {
            return err
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    
    // 7. Handle instant settlement if configured
    if atc.SettlementFrequency == SettlementFrequencyInstant {
        go s.processInstantSettlement(context.Background(), collection, atc)
    }
    
    // 8. Publish event
    s.events.Publish(ctx, "atc.collection.recorded", collection)
    
    return collection, nil
}

func (s *ATCService) calculateCommission(atc *ATCGrant, 
    amount decimal.Decimal) decimal.Decimal {
    var commission decimal.Decimal
    
    switch atc.CommissionType {
    case CommissionTypePercentage:
        commission = amount.Mul(atc.CommissionRate).Div(decimal.NewFromInt(100))
        
    case CommissionTypeFlatFee:
        commission = atc.CommissionFlat
        
    case CommissionTypeTiered:
        for _, tier := range atc.CommissionTiers {
            if amount.GreaterThanOrEqual(tier.MinAmount) &&
               (tier.MaxAmount.IsZero() || amount.LessThanOrEqual(tier.MaxAmount)) {
                commission = amount.Mul(tier.CommissionRate).Div(decimal.NewFromInt(100))
                break
            }
        }
        
    case CommissionTypeHybrid:
        percentageComm := amount.Mul(atc.CommissionRate).Div(decimal.NewFromInt(100))
        commission = atc.CommissionFlat.Add(percentageComm)
    }
    
    // Apply min/max constraints
    if atc.MinCommission.GreaterThan(decimal.Zero) && commission.LessThan(atc.MinCommission) {
        commission = atc.MinCommission
    }
    if atc.MaxCommission.GreaterThan(decimal.Zero) && commission.GreaterThan(atc.MaxCommission) {
        commission = atc.MaxCommission
    }
    
    return commission
}

// internal/workflow/settlement_workflow.go
func ATCSettlementWorkflow(ctx workflow.Context, 
    batchID uuid.UUID) (*SettlementResult, error) {
    logger := workflow.GetLogger(ctx)
    
    // Activity: Fetch pending collections for batch
    var collections []ATCCollection
    err := workflow.ExecuteActivity(ctx, 
        activities.FetchPendingCollections, batchID).Get(ctx, &collections)
    if err != nil {
        return nil, err
    }
    
    // Activity: Create settlement batch record
    var batch ATCSettlementBatch
    err = workflow.ExecuteActivity(ctx,
        activities.CreateSettlementBatch, batchID, collections).Get(ctx, &batch)
    if err != nil {
        return nil, err
    }
    
    // Activity: Initiate payment to grantor
    var paymentResult PaymentResult
    err = workflow.ExecuteActivity(ctx,
        activities.InitiateSettlementPayment, batch).Get(ctx, &paymentResult)
    if err != nil {
        // Handle payment failure
        workflow.ExecuteActivity(ctx, activities.HandleSettlementFailure, batch, err)
        return nil, err
    }
    
    // Activity: Update collection records
    err = workflow.ExecuteActivity(ctx,
        activities.MarkCollectionsSettled, batch.ID).Get(ctx, nil)
    if err != nil {
        return nil, err
    }
    
    // Activity: Send notifications
    workflow.ExecuteActivity(ctx, activities.NotifySettlement, batch)
    
    return &SettlementResult{
        BatchID:    batch.ID,
        Amount:     batch.NetSettlementAmount,
        Status:     "SETTLED",
    }, nil
}
```

HASURA GRAPHQL OPERATIONS:
```graphql
# ATC Grant Operations
mutation CreateATCGrant($input: CreateATCGrantInput!) {
  createATCGrant(input: $input) {
    id
    reference
    status
    grantorName
    granteeName
    commissionRate
    effectiveFrom
    effectiveTo
  }
}

mutation ApproveATCGrant($grantId: uuid!, $approvedBy: uuid!) {
  update_atc_grants_by_pk(
    pk_columns: {id: $grantId}
    _set: {
      status: ACTIVE
      approved_by: $approvedBy
      approved_at: "now()"
    }
  ) {
    id
    status
  }
}

mutation RevokeATCGrant($grantId: uuid!, $reason: String!) {
  revokeATCGrant(grantId: $grantId, reason: $reason) {
    success
    effectiveDate
  }
}

# Collection Operations
mutation RecordCollection($input: RecordCollectionInput!) {
  recordCollection(input: $input) {
    id
    reference
    grossAmount
    commissionAmount
    netAmount
    settlementStatus
  }
}

query GetCollectionsByATC($atcId: uuid!, $startDate: timestamptz!, $endDate: timestamptz!) {
  atc_collections(
    where: {
      atc_grant_id: {_eq: $atcId}
      created_at: {_gte: $startDate, _lte: $endDate}
    }
    order_by: {created_at: desc}
  ) {
    id
    reference
    collectedFromName
    grossAmount
    commissionAmount
    netAmount
    settlementStatus
    createdAt
  }
}

# Settlement Operations
query GetPendingSettlements($granteeId: uuid!) {
  atc_collections_aggregate(
    where: {
      atc_grant: {grantee_id: {_eq: $granteeId}}
      settlement_status: {_eq: PENDING}
    }
  ) {
    aggregate {
      count
      sum {
        gross_amount
        commission_amount
        net_amount
      }
    }
  }
}

subscription SettlementUpdates($grantorId: uuid!) {
  atc_settlement_batches(
    where: {atc_grant: {grantor_id: {_eq: $grantorId}}}
    order_by: {created_at: desc}
    limit: 10
  ) {
    id
    reference
    grossAmount
    netSettlementAmount
    paymentStatus
    settledAt
  }
}

# Hierarchical ATC Query
query GetATCHierarchy($rootAtcId: uuid!) {
  atc_grants_by_pk(id: $rootAtcId) {
    id
    reference
    grantorName
    granteeName
    commissionRate
    childAtcs {
      id
      reference
      granteeName
      commissionRate
      childAtcs {
        id
        reference
        granteeName
        commissionRate
      }
    }
  }
}
```

TEST CASES:
1. ATC grant creation with various commission types
2. Cascading ATC hierarchy (3+ levels deep)
3. Collection recording with commission calculation
4. Tiered commission calculation accuracy
5. Cumulative limit enforcement
6. Settlement batch creation and processing
7. Instant settlement flow
8. ATC revocation and pending collection handling
9. Multi-tenant isolation
10. Concurrent collection recording

EXPECTED DELIVERABLES:
1. Complete Hasura schema and metadata
2. Go backend for ATC business logic
3. Temporal workflows for settlements
4. Commission calculation engine
5. Settlement reconciliation
6. Hierarchical ATC management
7. GraphQL subscriptions for real-time updates
8. Comprehensive test suite
9. API documentation
10. Deployment configuration
```

---

# PART 2: INFRASTRUCTURE MIGRATION PROMPTS

## I-P01: YugabyteDB Migration (PostgreSQL Replacement)

```
PROJECT: OmniRoute - YugabyteDB Migration

CONTEXT:
Migrate from PostgreSQL to YugabyteDB for distributed, horizontally scalable database layer. YugabyteDB provides:
- PostgreSQL compatibility (use existing queries)
- Automatic sharding and rebalancing
- Multi-region deployment
- High availability with Raft consensus
- Distributed transactions

TECHNICAL REQUIREMENTS:
- YugabyteDB: 2.20+ (latest stable)
- Compatibility: PostgreSQL 11.2 wire protocol
- Deployment: Kubernetes with yugabyte-operator
- Regions: Lagos (primary), London (secondary), Singapore (tertiary)

KEY DIFFERENCES FROM POSTGRESQL:
1. Primary keys should be hash-sharded by default
2. Use HASH partitioning for hot tables
3. Colocated tables for small reference data
4. Tablet splitting for large tables
5. Connection pooling via Odyssey/PgBouncer compatible

SCHEMA MODIFICATIONS:
```sql
-- Use HASH sharding for distributed tables
CREATE TABLE orders (
    id UUID DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_amount DECIMAL(15,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY ((tenant_id) HASH, id ASC)  -- Hash on tenant_id
) SPLIT INTO 16 TABLETS;

-- Colocated tables for reference data (small, frequently joined)
CREATE DATABASE omniroute WITH COLOCATION = true;

CREATE TABLE currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100),
    symbol VARCHAR(10)
) WITH (COLOCATED = true);

CREATE TABLE countries (
    code VARCHAR(2) PRIMARY KEY,
    name VARCHAR(100),
    currency_code VARCHAR(3)
) WITH (COLOCATED = true);

-- Geo-partitioned table for data residency
CREATE TABLE customer_data (
    id UUID,
    tenant_id UUID,
    region VARCHAR(10),
    data JSONB,
    PRIMARY KEY ((region) HASH, tenant_id, id)
) TABLESPACE nigeria_tablespace;  -- Data stays in Nigeria

-- Index strategies for YugabyteDB
CREATE INDEX NONCONCURRENTLY idx_orders_customer 
ON orders (customer_id ASC) 
SPLIT INTO 8 TABLETS;

-- Use covering indexes for common queries
CREATE INDEX idx_orders_status_covering 
ON orders (tenant_id HASH, status ASC) 
INCLUDE (total_amount, created_at);
```

GO DATABASE LAYER:
```go
// pkg/database/yugabyte.go
package database

import (
    "context"
    "database/sql"
    "github.com/jackc/pgx/v5/pgxpool"
)

type YugabyteConfig struct {
    Hosts           []string  // Multiple hosts for HA
    Port            int
    Database        string
    User            string
    Password        string
    SSLMode         string
    MaxConnections  int
    MinConnections  int
    LoadBalance     bool      // YugabyteDB load balancing
    TopologyKeys    string    // Prefer local region
}

func NewYugabytePool(cfg YugabyteConfig) (*pgxpool.Pool, error) {
    // Build connection string with YugabyteDB-specific options
    connStr := fmt.Sprintf(
        "host=%s port=%d dbname=%s user=%s password=%s sslmode=%s "+
        "load_balance=%t topology_keys=%s "+
        "pool_max_conns=%d pool_min_conns=%d",
        strings.Join(cfg.Hosts, ","),
        cfg.Port,
        cfg.Database,
        cfg.User,
        cfg.Password,
        cfg.SSLMode,
        cfg.LoadBalance,
        cfg.TopologyKeys,  // e.g., "aws.us-west-2.us-west-2a"
        cfg.MaxConnections,
        cfg.MinConnections,
    )
    
    poolConfig, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        return nil, err
    }
    
    // Configure for YugabyteDB
    poolConfig.ConnConfig.RuntimeParams["application_name"] = "omniroute"
    
    return pgxpool.NewWithConfig(context.Background(), poolConfig)
}

// Repository with YugabyteDB optimizations
type OrderRepository struct {
    pool *pgxpool.Pool
}

func (r *OrderRepository) Create(ctx context.Context, order *Order) error {
    // Use UPSERT for idempotency (YugabyteDB handles conflicts efficiently)
    query := `
        INSERT INTO orders (id, tenant_id, customer_id, status, total_amount, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (tenant_id, id) DO UPDATE SET
            status = EXCLUDED.status,
            total_amount = EXCLUDED.total_amount
    `
    _, err := r.pool.Exec(ctx, query, 
        order.ID, order.TenantID, order.CustomerID, 
        order.Status, order.TotalAmount, order.CreatedAt)
    return err
}

func (r *OrderRepository) GetByTenant(ctx context.Context, 
    tenantID uuid.UUID, limit int) ([]Order, error) {
    // Query optimized for hash-sharded tenant_id
    query := `
        SELECT id, tenant_id, customer_id, status, total_amount, created_at
        FROM orders
        WHERE tenant_id = $1
        ORDER BY created_at DESC
        LIMIT $2
    `
    rows, err := r.pool.Query(ctx, query, tenantID, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var orders []Order
    for rows.Next() {
        var o Order
        if err := rows.Scan(&o.ID, &o.TenantID, &o.CustomerID, 
            &o.Status, &o.TotalAmount, &o.CreatedAt); err != nil {
            return nil, err
        }
        orders = append(orders, o)
    }
    return orders, nil
}
```

KUBERNETES DEPLOYMENT:
```yaml
# yugabyte-cluster.yaml
apiVersion: yugabyte.com/v1alpha1
kind: YBCluster
metadata:
  name: omniroute-yb
  namespace: database
spec:
  image:
    repository: yugabytedb/yugabyte
    tag: 2.20.1.0-b97
  
  tserver:
    replicas: 3
    storage:
      count: 2
      size: 100Gi
      storageClass: ssd
    resources:
      requests:
        cpu: "4"
        memory: "16Gi"
      limits:
        cpu: "8"
        memory: "32Gi"
    gflags:
      # Performance tuning
      - key: ysql_max_connections
        value: "500"
      - key: ysql_enable_packed_row
        value: "true"
      - key: enable_automatic_tablet_splitting
        value: "true"
      - key: tablet_split_low_phase_size_threshold_bytes
        value: "536870912"  # 512MB
  
  master:
    replicas: 3
    storage:
      size: 10Gi
      storageClass: ssd
    resources:
      requests:
        cpu: "2"
        memory: "4Gi"
    gflags:
      - key: replication_factor
        value: "3"
      - key: load_balancer_max_concurrent_adds
        value: "5"
  
  enableTLS: true
  tlsSecret: yugabyte-tls
  
  # Multi-region configuration
  gflags:
    master:
      - key: placement_cloud
        value: "gcp"
      - key: placement_region
        value: "africa-south1"
      - key: placement_zone
        value: "africa-south1-a"
```

MIGRATION SCRIPT:
```go
// cmd/migrate-to-yugabyte/main.go
func main() {
    // 1. Create schema in YugabyteDB
    // 2. Migrate data in batches
    // 3. Verify row counts
    // 4. Switch application connection strings
    // 5. Monitor and rollback if needed
}

func migrateTable(ctx context.Context, srcPool, dstPool *pgxpool.Pool, 
    tableName string, batchSize int) error {
    // Get total count
    var count int64
    srcPool.QueryRow(ctx, "SELECT COUNT(*) FROM "+tableName).Scan(&count)
    
    // Migrate in batches
    for offset := int64(0); offset < count; offset += int64(batchSize) {
        rows, err := srcPool.Query(ctx, 
            fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT %d OFFSET %d", 
                tableName, batchSize, offset))
        if err != nil {
            return err
        }
        
        // Bulk insert to YugabyteDB
        batch := &pgx.Batch{}
        for rows.Next() {
            // Add to batch
        }
        rows.Close()
        
        results := dstPool.SendBatch(ctx, batch)
        results.Close()
        
        log.Printf("Migrated %d/%d rows of %s", offset+int64(batchSize), count, tableName)
    }
    
    return nil
}
```

TEST CASES:
1. Connection pool with load balancing
2. Multi-region read/write operations
3. Distributed transaction consistency
4. Tablet splitting under load
5. Failover and recovery
6. Query performance vs PostgreSQL
7. Index effectiveness
8. Data residency verification

EXPECTED DELIVERABLES:
1. YugabyteDB Kubernetes deployment
2. Schema migration scripts
3. Go database layer with YugabyteDB optimizations
4. Data migration tooling
5. Monitoring and alerting setup
6. Performance benchmarks
7. Runbooks for operations
8. Rollback procedures
```

---

## I-P02: DragonflyDB Migration (Redis Replacement)

```
PROJECT: OmniRoute - DragonflyDB Migration

CONTEXT:
Migrate from Redis to DragonflyDB for improved performance and efficiency. DragonflyDB provides:
- Full Redis API compatibility
- Multi-threaded architecture (uses all CPU cores)
- 25x more memory efficient than Redis
- Supports Redis Cluster protocol
- Native snapshotting with instant recovery

TECHNICAL REQUIREMENTS:
- DragonflyDB: v1.x (latest stable)
- Compatibility: Redis 7.x API
- Deployment: Kubernetes with StatefulSet
- Memory: 32GB per node
- Cluster: 3 nodes for HA

KEY DIFFERENCES FROM REDIS:
1. No need for Redis Cluster sharding (single instance handles more)
2. Better memory efficiency for large datasets
3. Built-in TLS without performance penalty
4. Faster snapshots and recovery

GO CLIENT CONFIGURATION:
```go
// pkg/cache/dragonfly.go
package cache

import (
    "context"
    "time"
    "github.com/redis/go-redis/v9"
)

type DragonflyConfig struct {
    Addresses    []string
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
    MaxRetries   int
    TLSEnabled   bool
}

func NewDragonflyClient(cfg DragonflyConfig) (*redis.ClusterClient, error) {
    opts := &redis.ClusterOptions{
        Addrs:        cfg.Addresses,
        Password:     cfg.Password,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        MaxRetries:   cfg.MaxRetries,
        
        // DragonflyDB specific optimizations
        ReadOnly:       false,
        RouteRandomly:  true,
        
        // Connection timeouts
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
        PoolTimeout:  4 * time.Second,
    }
    
    if cfg.TLSEnabled {
        opts.TLSConfig = &tls.Config{
            MinVersion: tls.VersionTLS12,
        }
    }
    
    client := redis.NewClusterClient(opts)
    
    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("dragonfly ping failed: %w", err)
    }
    
    return client, nil
}

// Cache service with DragonflyDB
type CacheService struct {
    client *redis.ClusterClient
}

func (c *CacheService) SetWithTTL(ctx context.Context, 
    key string, value interface{}, ttl time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }
    return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
    data, err := c.client.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dest)
}

// Use DragonflyDB's efficient hash operations for session storage
func (c *CacheService) SetSession(ctx context.Context, 
    sessionID string, session *Session, ttl time.Duration) error {
    key := fmt.Sprintf("session:%s", sessionID)
    
    pipe := c.client.Pipeline()
    pipe.HSet(ctx, key, map[string]interface{}{
        "user_id":    session.UserID.String(),
        "tenant_id":  session.TenantID.String(),
        "roles":      strings.Join(session.Roles, ","),
        "created_at": session.CreatedAt.Unix(),
    })
    pipe.Expire(ctx, key, ttl)
    
    _, err := pipe.Exec(ctx)
    return err
}

// Distributed locking with DragonflyDB
func (c *CacheService) AcquireLock(ctx context.Context, 
    key string, ttl time.Duration) (bool, error) {
    result, err := c.client.SetNX(ctx, "lock:"+key, "1", ttl).Result()
    return result, err
}

func (c *CacheService) ReleaseLock(ctx context.Context, key string) error {
    return c.client.Del(ctx, "lock:"+key).Err()
}

// Rate limiting using DragonflyDB's efficient INCR
func (c *CacheService) CheckRateLimit(ctx context.Context, 
    key string, limit int, window time.Duration) (bool, error) {
    pipe := c.client.Pipeline()
    
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    return incr.Val() <= int64(limit), nil
}
```

KUBERNETES DEPLOYMENT:
```yaml
# dragonfly-statefulset.yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: dragonfly
  namespace: cache
spec:
  serviceName: dragonfly
  replicas: 3
  selector:
    matchLabels:
      app: dragonfly
  template:
    metadata:
      labels:
        app: dragonfly
    spec:
      containers:
      - name: dragonfly
        image: docker.dragonflydb.io/dragonflydb/dragonfly:v1.13.0
        args:
          - "--cluster_mode=emulated"
          - "--maxmemory=28gb"
          - "--proactor_threads=8"
          - "--dbfilename=dump"
          - "--snapshot_cron=*/30 * * * *"
          - "--requirepass=$(DRAGONFLY_PASSWORD)"
          - "--tls"
          - "--tls_cert_file=/certs/tls.crt"
          - "--tls_key_file=/certs/tls.key"
        env:
          - name: DRAGONFLY_PASSWORD
            valueFrom:
              secretKeyRef:
                name: dragonfly-secret
                key: password
        ports:
          - containerPort: 6379
            name: redis
        resources:
          requests:
            cpu: "4"
            memory: "32Gi"
          limits:
            cpu: "8"
            memory: "32Gi"
        volumeMounts:
          - name: data
            mountPath: /data
          - name: certs
            mountPath: /certs
            readOnly: true
        livenessProbe:
          tcpSocket:
            port: 6379
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          exec:
            command:
              - sh
              - -c
              - 'redis-cli -a $DRAGONFLY_PASSWORD ping | grep -q PONG'
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
        - name: certs
          secret:
            secretName: dragonfly-tls
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        storageClassName: ssd
        resources:
          requests:
            storage: 100Gi
---
apiVersion: v1
kind: Service
metadata:
  name: dragonfly
  namespace: cache
spec:
  type: ClusterIP
  ports:
    - port: 6379
      targetPort: 6379
  selector:
    app: dragonfly
```

TEST CASES:
1. Redis API compatibility
2. Connection pooling performance
3. Hash operations for sessions
4. Pub/Sub functionality
5. Distributed locking
6. Rate limiting accuracy
7. Snapshot and recovery
8. Memory efficiency vs Redis

EXPECTED DELIVERABLES:
1. DragonflyDB Kubernetes deployment
2. Go client wrapper
3. Migration scripts from Redis
4. Monitoring with Prometheus
5. Performance benchmarks
6. Operations runbooks
```

---

## I-P03: Redpanda Migration (Kafka Replacement)

```
PROJECT: OmniRoute - Redpanda Migration

CONTEXT:
Migrate from Apache Kafka to Redpanda for simplified operations and better performance. Redpanda provides:
- Full Kafka API compatibility
- No ZooKeeper dependency
- Lower latency (p99 < 10ms)
- Built-in Schema Registry
- Simpler operations
- Better resource efficiency

TECHNICAL REQUIREMENTS:
- Redpanda: v23.x (latest stable)
- Compatibility: Kafka 3.x API
- Deployment: Kubernetes with redpanda-operator
- Cluster: 3 brokers for HA
- Storage: 500GB NVMe per broker

KEY DIFFERENCES FROM KAFKA:
1. No ZooKeeper (uses Raft internally)
2. Single binary deployment
3. Built-in Admin API and Console
4. Native Schema Registry
5. Better performance with fewer resources

GO CLIENT CONFIGURATION:
```go
// pkg/messaging/redpanda.go
package messaging

import (
    "context"
    "crypto/tls"
    "github.com/twmb/franz-go/pkg/kgo"
    "github.com/twmb/franz-go/pkg/sasl/scram"
)

type RedpandaConfig struct {
    Brokers        []string
    Username       string
    Password       string
    TLSEnabled     bool
    ConsumerGroup  string
    ProducerConfig ProducerConfig
    ConsumerConfig ConsumerConfig
}

type ProducerConfig struct {
    BatchMaxBytes int
    LingerMs      int
    Compression   string
    Acks          string
}

type ConsumerConfig struct {
    AutoOffsetReset string
    MaxPollRecords  int
}

func NewRedpandaClient(cfg RedpandaConfig) (*kgo.Client, error) {
    opts := []kgo.Opt{
        kgo.SeedBrokers(cfg.Brokers...),
        kgo.ConsumerGroup(cfg.ConsumerGroup),
        
        // Producer settings
        kgo.ProducerBatchMaxBytes(int32(cfg.ProducerConfig.BatchMaxBytes)),
        kgo.ProducerLinger(time.Duration(cfg.ProducerConfig.LingerMs) * time.Millisecond),
        kgo.RequiredAcks(parseAcks(cfg.ProducerConfig.Acks)),
        
        // Consumer settings
        kgo.ConsumeResetOffset(parseOffset(cfg.ConsumerConfig.AutoOffsetReset)),
        
        // Retry settings
        kgo.RetryBackoffFn(func(attempt int) time.Duration {
            return time.Duration(attempt*100) * time.Millisecond
        }),
        kgo.RequestRetries(5),
    }
    
    // Add SASL authentication
    if cfg.Username != "" {
        mechanism := scram.Auth{
            User: cfg.Username,
            Pass: cfg.Password,
        }.AsSha256Mechanism()
        opts = append(opts, kgo.SASL(mechanism))
    }
    
    // Add TLS if enabled
    if cfg.TLSEnabled {
        opts = append(opts, kgo.DialTLSConfig(&tls.Config{
            MinVersion: tls.VersionTLS12,
        }))
    }
    
    return kgo.NewClient(opts...)
}

// Producer service
type EventProducer struct {
    client *kgo.Client
}

func (p *EventProducer) Publish(ctx context.Context, 
    topic string, key string, event interface{}) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    record := &kgo.Record{
        Topic: topic,
        Key:   []byte(key),
        Value: data,
        Headers: []kgo.RecordHeader{
            {Key: "content-type", Value: []byte("application/json")},
            {Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
        },
    }
    
    // Synchronous produce with delivery confirmation
    results := p.client.ProduceSync(ctx, record)
    return results.FirstErr()
}

func (p *EventProducer) PublishBatch(ctx context.Context, 
    topic string, events []Event) error {
    records := make([]*kgo.Record, len(events))
    for i, event := range events {
        data, _ := json.Marshal(event)
        records[i] = &kgo.Record{
            Topic: topic,
            Key:   []byte(event.Key),
            Value: data,
        }
    }
    
    results := p.client.ProduceSync(ctx, records...)
    return results.FirstErr()
}

// Consumer service with proper offset management
type EventConsumer struct {
    client   *kgo.Client
    handlers map[string]EventHandler
}

func (c *EventConsumer) Subscribe(topics []string, handler EventHandler) {
    c.client.AddConsumeTopics(topics...)
    for _, topic := range topics {
        c.handlers[topic] = handler
    }
}

func (c *EventConsumer) Start(ctx context.Context) error {
    for {
        fetches := c.client.PollFetches(ctx)
        if fetches.IsClientClosed() {
            return nil
        }
        
        fetches.EachError(func(topic string, partition int32, err error) {
            log.Printf("fetch error topic %s partition %d: %v", topic, partition, err)
        })
        
        fetches.EachRecord(func(record *kgo.Record) {
            handler, ok := c.handlers[record.Topic]
            if !ok {
                return
            }
            
            event := Event{
                Topic:     record.Topic,
                Key:       string(record.Key),
                Value:     record.Value,
                Timestamp: record.Timestamp,
                Partition: record.Partition,
                Offset:    record.Offset,
            }
            
            if err := handler.Handle(ctx, event); err != nil {
                log.Printf("handler error: %v", err)
                // Handle error (DLQ, retry, etc.)
            }
        })
        
        // Commit offsets after processing
        if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
            log.Printf("commit error: %v", err)
        }
    }
}
```

KUBERNETES DEPLOYMENT:
```yaml
# redpanda-cluster.yaml
apiVersion: cluster.redpanda.com/v1alpha1
kind: Redpanda
metadata:
  name: omniroute
  namespace: messaging
spec:
  image: "docker.redpanda.com/redpandadata/redpanda:v23.3.5"
  version: v23.3.5
  
  statefulset:
    replicas: 3
    budget:
      maxUnavailable: 1
    
  resources:
    cpu:
      cores: 4
    memory:
      container:
        max: 16Gi
      redpanda:
        memory: 12Gi
        reserveMemory: 2Gi
  
  storage:
    capacity: 500Gi
    storageClassName: nvme-ssd
  
  configuration:
    developerMode: false
    rpcServer:
      port: 33145
    kafkaApi:
      port: 9092
      tls:
        enabled: true
        cert: redpanda-tls
    adminApi:
      port: 9644
      tls:
        enabled: true
        cert: redpanda-tls
    schemaRegistry:
      port: 8081
      tls:
        enabled: true
        cert: redpanda-tls
  
  auth:
    sasl:
      enabled: true
      users:
        - name: omniroute-producer
          password:
            valueFrom:
              secretKeyRef:
                name: redpanda-users
                key: producer-password
          mechanism: SCRAM-SHA-256
        - name: omniroute-consumer
          password:
            valueFrom:
              secretKeyRef:
                name: redpanda-users
                key: consumer-password
          mechanism: SCRAM-SHA-256
  
  # Topic configuration
  additionalConfiguration:
    log_retention_ms: "604800000"  # 7 days
    log_segment_size: "134217728"  # 128MB
    group_topic_partitions: "16"
    default_topic_replications: "3"
    default_topic_partitions: "6"
---
# Create topics via Kubernetes Job
apiVersion: batch/v1
kind: Job
metadata:
  name: create-topics
  namespace: messaging
spec:
  template:
    spec:
      containers:
      - name: rpk
        image: docker.redpanda.com/redpandadata/redpanda:v23.3.5
        command:
          - /bin/bash
          - -c
          - |
            rpk topic create orders.created --partitions 12 --replicas 3
            rpk topic create orders.updated --partitions 12 --replicas 3
            rpk topic create payments.completed --partitions 6 --replicas 3
            rpk topic create inventory.updated --partitions 12 --replicas 3
            rpk topic create workers.assigned --partitions 6 --replicas 3
            rpk topic create atc.collections --partitions 6 --replicas 3
            rpk topic create atc.settlements --partitions 3 --replicas 3
      restartPolicy: OnFailure
```

SCHEMA REGISTRY INTEGRATION:
```go
// pkg/messaging/schema_registry.go
package messaging

import (
    "github.com/twmb/franz-go/pkg/sr"
)

type SchemaRegistry struct {
    client *sr.Client
}

func NewSchemaRegistry(urls []string, opts ...sr.Opt) (*SchemaRegistry, error) {
    client, err := sr.NewClient(append(opts, sr.URLs(urls...))...)
    if err != nil {
        return nil, err
    }
    return &SchemaRegistry{client: client}, nil
}

func (r *SchemaRegistry) RegisterSchema(ctx context.Context, 
    subject string, schema string) (int, error) {
    ss, err := r.client.CreateSchema(ctx, subject, sr.Schema{
        Schema: schema,
        Type:   sr.TypeJSON,
    })
    if err != nil {
        return 0, err
    }
    return ss.ID, nil
}

// Serde for automatic serialization/deserialization
type OrderEventSerde struct {
    serde sr.Serde
}

func NewOrderEventSerde(registry *sr.Client, topic string) (*OrderEventSerde, error) {
    var serde sr.Serde
    serde.Register(
        1, // Schema ID
        OrderEvent{},
        sr.EncodeFn(json.Marshal),
        sr.DecodeFn(json.Unmarshal),
    )
    return &OrderEventSerde{serde: serde}, nil
}
```

TEST CASES:
1. Kafka API compatibility
2. Producer delivery guarantees
3. Consumer group rebalancing
4. Schema Registry operations
5. Topic creation and management
6. SASL authentication
7. TLS encryption
8. Latency benchmarks (p99 < 10ms)

EXPECTED DELIVERABLES:
1. Redpanda Kubernetes deployment
2. Go producer/consumer clients
3. Schema Registry integration
4. Migration scripts from Kafka
5. Topic configuration
6. Monitoring setup
7. Performance benchmarks
8. Operations runbooks
```

---

# PART 3: COMBINED INTEGRATION PROMPT

## MASTER-P01: Complete Financial Integration with Modern Infrastructure

```
PROJECT: OmniRoute - Complete Financial Integration Platform

CONTEXT:
Implement the complete financial integration layer for OmniRoute using:
- Hasura GraphQL API for bank integrations
- Authority to Collect (ATC) system for B2B distribution
- YugabyteDB for distributed database
- DragonflyDB for caching
- Redpanda for event streaming

This is a comprehensive prompt that ties all components together.

ARCHITECTURE OVERVIEW:
┌─────────────────────────────────────────────────────────────────┐
│                    OmniRoute Financial Platform                  │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   Hasura    │  │  Go Backend │  │   Temporal Workflows    │ │
│  │  GraphQL    │◄─│  (Actions)  │◄─│   (Settlements)         │ │
│  └──────┬──────┘  └──────┬──────┘  └───────────┬─────────────┘ │
│         │                │                      │               │
├─────────┼────────────────┼──────────────────────┼───────────────┤
│         ▼                ▼                      ▼               │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    Service Layer                            ││
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    ││
│  │  │ Payment  │  │   ATC    │  │  Bank    │  │ Settle-  │    ││
│  │  │ Service  │  │ Service  │  │ Connector│  │  ment    │    ││
│  │  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    ││
│  └───────┼─────────────┼─────────────┼─────────────┼───────────┘│
│          │             │             │             │            │
├──────────┼─────────────┼─────────────┼─────────────┼────────────┤
│          ▼             ▼             ▼             ▼            │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                   Data Layer                                ││
│  │  ┌────────────┐  ┌────────────┐  ┌────────────────────────┐││
│  │  │ YugabyteDB │  │ DragonflyDB│  │      Redpanda         │││
│  │  │ (Primary)  │  │  (Cache)   │  │  (Event Streaming)    │││
│  │  └────────────┘  └────────────┘  └────────────────────────┘││
│  └─────────────────────────────────────────────────────────────┘│
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                   External Integrations                         │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────────┐│
│  │  NIBSS  │  │ M-Pesa  │  │ MTN MoMo│  │   Other Banks       ││
│  └─────────┘  └─────────┘  └─────────┘  └─────────────────────┘│
└─────────────────────────────────────────────────────────────────┘

IMPLEMENTATION PHASES:

PHASE 1: Infrastructure Setup
- Deploy YugabyteDB cluster (3 nodes, multi-AZ)
- Deploy DragonflyDB cluster
- Deploy Redpanda cluster (3 brokers)
- Configure networking and security

PHASE 2: Database Schema
- Create all tables in YugabyteDB with proper sharding
- Set up Hasura metadata and relationships
- Configure event triggers

PHASE 3: Core Services
- Implement Payment Service
- Implement ATC Service
- Implement Bank Connectors (NIBSS, Mobile Money)
- Set up Temporal workflows for settlements

PHASE 4: Hasura Integration
- Configure Actions for all operations
- Set up GraphQL subscriptions
- Implement authorization rules

PHASE 5: Testing & Deployment
- Integration testing
- Performance testing
- Security testing
- Production deployment

COMPLETE DOCKER COMPOSE FOR LOCAL DEVELOPMENT:
```yaml
version: '3.8'
services:
  yugabytedb:
    image: yugabytedb/yugabyte:2.20.1.0-b97
    command: ["bin/yugabyted", "start", "--daemon=false"]
    ports:
      - "5433:5433"
      - "9000:9000"
      - "7000:7000"
    environment:
      - YUGABYTE_DB=omniroute
    volumes:
      - yugabyte-data:/home/yugabyte/yb_data

  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly:v1.13.0
    ports:
      - "6379:6379"
    command: ["--maxmemory", "2gb", "--proactor_threads", "2"]
    volumes:
      - dragonfly-data:/data

  redpanda:
    image: docker.redpanda.com/redpandadata/redpanda:v23.3.5
    command:
      - redpanda
      - start
      - --smp 1
      - --memory 1G
      - --overprovisioned
      - --kafka-addr PLAINTEXT://0.0.0.0:9092
      - --advertise-kafka-addr PLAINTEXT://redpanda:9092
      - --pandaproxy-addr 0.0.0.0:8082
      - --advertise-pandaproxy-addr redpanda:8082
      - --schema-registry-addr 0.0.0.0:8081
    ports:
      - "9092:9092"
      - "8081:8081"
      - "8082:8082"
      - "9644:9644"
    volumes:
      - redpanda-data:/var/lib/redpanda/data

  hasura:
    image: hasura/graphql-engine:v2.36.0
    ports:
      - "8080:8080"
    environment:
      HASURA_GRAPHQL_DATABASE_URL: "postgres://yugabyte:yugabyte@yugabytedb:5433/omniroute"
      HASURA_GRAPHQL_ENABLE_CONSOLE: "true"
      HASURA_GRAPHQL_ADMIN_SECRET: "admin-secret"
      HASURA_GRAPHQL_JWT_SECRET: '{"type":"HS256","key":"your-jwt-secret-key-min-32-chars"}'
      HASURA_GRAPHQL_UNAUTHORIZED_ROLE: "anonymous"
      BANK_GATEWAY_URL: "http://bank-gateway:8090"
    depends_on:
      - yugabytedb

  temporal:
    image: temporalio/auto-setup:1.22.0
    ports:
      - "7233:7233"
    environment:
      - DB=postgresql
      - DB_PORT=5433
      - POSTGRES_USER=yugabyte
      - POSTGRES_PWD=yugabyte
      - POSTGRES_SEEDS=yugabytedb

  bank-gateway:
    build: ./services/bank-gateway
    ports:
      - "8090:8090"
    environment:
      - YUGABYTE_HOST=yugabytedb
      - YUGABYTE_PORT=5433
      - DRAGONFLY_HOST=dragonfly
      - REDPANDA_BROKERS=redpanda:9092
      - TEMPORAL_HOST=temporal:7233
    depends_on:
      - yugabytedb
      - dragonfly
      - redpanda
      - temporal

volumes:
  yugabyte-data:
  dragonfly-data:
  redpanda-data:
```

EXPECTED DELIVERABLES:
1. Complete infrastructure deployment (YugabyteDB, DragonflyDB, Redpanda)
2. Hasura GraphQL API with all schemas
3. Go backend services (Payment, ATC, Bank Connectors)
4. Temporal workflows for settlements
5. Bank integration adapters (NIBSS, Mobile Money)
6. Comprehensive test suite
7. Docker Compose for local development
8. Kubernetes manifests for production
9. CI/CD pipeline configuration
10. API documentation and runbooks
```

---

*Document Version: 1.0 | January 2026 | BillyRonks Global Limited*
