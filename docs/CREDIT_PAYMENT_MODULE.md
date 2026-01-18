# OmniRoute Commerce Platform
# Credit & Payment Management Module
# Technical Design Document

---

## Overview

The Credit & Payment Management Module is the financial backbone of OmniRoute, enabling:
- **Trade Credit**: Automated credit limits for B2B customers
- **Payment Processing**: Multi-method payment orchestration
- **Collections**: Intelligent collection management
- **Settlement**: Automated payouts to all stakeholders
- **Financing**: Integration with external finance partners

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    CREDIT & PAYMENT MANAGEMENT ARCHITECTURE                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                        PAYMENT ORCHESTRATION LAYER                       │  │
│   │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │  │
│   │  │   Cards     │ │   Bank      │ │   Mobile    │ │   Wallets   │       │  │
│   │  │  (Paystack) │ │  Transfer   │ │   Money     │ │  (Internal) │       │  │
│   │  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘       │  │
│   │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │  │
│   │  │    USSD     │ │    Cash     │ │   Credit    │ │    POS      │       │  │
│   │  │  (Direct)   │ │ Collection  │ │  (B2B Only) │ │  Terminal   │       │  │
│   │  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘       │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                        │                                        │
│                                        ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                          CORE FINANCIAL SERVICES                         │  │
│   │                                                                          │  │
│   │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐            │  │
│   │  │    CREDIT      │  │   INVOICING    │  │  COLLECTIONS   │            │  │
│   │  │    ENGINE      │  │   & BILLING    │  │   MANAGEMENT   │            │  │
│   │  │                │  │                │  │                │            │  │
│   │  │ • Scoring      │  │ • Generation   │  │ • Aging        │            │  │
│   │  │ • Limits       │  │ • Delivery     │  │ • Reminders    │            │  │
│   │  │ • Terms        │  │ • Tracking     │  │ • Escalation   │            │  │
│   │  │ • Monitoring   │  │ • Reconcile    │  │ • Write-off    │            │  │
│   │  └────────────────┘  └────────────────┘  └────────────────┘            │  │
│   │                                                                          │  │
│   │  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐            │  │
│   │  │    WALLET      │  │  SETTLEMENT    │  │   REPORTING    │            │  │
│   │  │   SERVICE      │  │    ENGINE      │  │  & ANALYTICS   │            │  │
│   │  │                │  │                │  │                │            │  │
│   │  │ • Balance      │  │ • Schedules    │  │ • Cash flow    │            │  │
│   │  │ • Transactions │  │ • Payouts      │  │ • Aging        │            │  │
│   │  │ • Holds        │  │ • Splits       │  │ • Performance  │            │  │
│   │  │ • Withdrawals  │  │ • Reconcile    │  │ • Compliance   │            │  │
│   │  └────────────────┘  └────────────────┘  └────────────────┘            │  │
│   │                                                                          │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                        │                                        │
│                                        ▼                                        │
│   ┌─────────────────────────────────────────────────────────────────────────┐  │
│   │                      EXTERNAL FINANCE INTEGRATIONS                       │  │
│   │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐       │  │
│   │  │   Invoice   │ │   Working   │ │    BNPL     │ │  Insurance  │       │  │
│   │  │  Financing  │ │  Capital    │ │  Partners   │ │  Partners   │       │  │
│   │  └─────────────┘ └─────────────┘ └─────────────┘ └─────────────┘       │  │
│   └─────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## Module 1: Credit Engine

### 1.1 Credit Scoring Model

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        CREDIT SCORING FRAMEWORK                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   SCORE COMPONENTS (Total: 1000 points)                                     │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 1. TRANSACTION HISTORY (350 points)                                 │   │
│   │    ├── Order frequency: 0-100 pts                                   │   │
│   │    │   • <5 orders/month: 20 pts                                    │   │
│   │    │   • 5-15 orders/month: 60 pts                                  │   │
│   │    │   • >15 orders/month: 100 pts                                  │   │
│   │    ├── Order value consistency: 0-100 pts                           │   │
│   │    │   • High variance: 30 pts                                      │   │
│   │    │   • Medium variance: 70 pts                                    │   │
│   │    │   • Low variance: 100 pts                                      │   │
│   │    ├── Platform tenure: 0-100 pts                                   │   │
│   │    │   • <3 months: 20 pts                                          │   │
│   │    │   • 3-12 months: 60 pts                                        │   │
│   │    │   • >12 months: 100 pts                                        │   │
│   │    └── Growth trend: 0-50 pts                                       │   │
│   │        • Declining: 0 pts                                           │   │
│   │        • Stable: 25 pts                                             │   │
│   │        • Growing: 50 pts                                            │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 2. PAYMENT BEHAVIOR (350 points)                                    │   │
│   │    ├── On-time payment rate: 0-150 pts                              │   │
│   │    │   • <70%: 0 pts                                                │   │
│   │    │   • 70-85%: 75 pts                                             │   │
│   │    │   • 85-95%: 120 pts                                            │   │
│   │    │   • >95%: 150 pts                                              │   │
│   │    ├── Days past due (average): 0-100 pts                           │   │
│   │    │   • >30 days: 0 pts                                            │   │
│   │    │   • 15-30 days: 40 pts                                         │   │
│   │    │   • 7-14 days: 70 pts                                          │   │
│   │    │   • <7 days: 100 pts                                           │   │
│   │    ├── Payment method reliability: 0-50 pts                         │   │
│   │    │   • Mostly cash: 30 pts                                        │   │
│   │    │   • Mixed: 40 pts                                              │   │
│   │    │   • Digital payments: 50 pts                                   │   │
│   │    └── Collection difficulty: 0-50 pts                              │   │
│   │        • High effort: 10 pts                                        │   │
│   │        • Medium effort: 30 pts                                      │   │
│   │        • Low effort: 50 pts                                         │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 3. BUSINESS PROFILE (200 points)                                    │   │
│   │    ├── Business registration: 0-50 pts                              │   │
│   │    │   • Unregistered: 10 pts                                       │   │
│   │    │   • Registered <2 years: 30 pts                                │   │
│   │    │   • Registered >2 years: 50 pts                                │   │
│   │    ├── Location quality: 0-50 pts                                   │   │
│   │    │   • High-risk area: 20 pts                                     │   │
│   │    │   • Medium-risk area: 35 pts                                   │   │
│   │    │   • Low-risk area: 50 pts                                      │   │
│   │    ├── Business type: 0-50 pts                                      │   │
│   │    │   • Kiosk: 25 pts                                              │   │
│   │    │   • Small shop: 35 pts                                         │   │
│   │    │   • Medium store: 45 pts                                       │   │
│   │    │   • Large retailer: 50 pts                                     │   │
│   │    └── Verified inventory: 0-50 pts                                 │   │
│   │        • Not verified: 0 pts                                        │   │
│   │        • Partially verified: 25 pts                                 │   │
│   │        • Fully verified: 50 pts                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │ 4. EXTERNAL SIGNALS (100 points)                                    │   │
│   │    ├── BVN/NIN verification: 0-30 pts                               │   │
│   │    ├── Bank statement analysis: 0-40 pts                            │   │
│   │    └── Social/alternative data: 0-30 pts                            │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   SCORE INTERPRETATION:                                                     │
│   • 800-1000: Premium (highest limits, best terms)                         │
│   • 600-799:  Standard (normal limits and terms)                           │
│   • 400-599:  Limited (reduced limits, shorter terms)                      │
│   • 200-399:  Restricted (minimal credit, COD preferred)                   │
│   • <200:     No credit (COD only)                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Credit Limit Calculation

```python
# Credit Limit Formula (Pseudocode)

def calculate_credit_limit(customer: Customer, score: CreditScore) -> CreditLimit:
    """
    Calculate credit limit based on multiple factors
    """
    
    # Base limit from transaction history
    avg_monthly_orders = customer.get_avg_monthly_order_value(months=6)
    base_limit = avg_monthly_orders * 0.5  # 50% of monthly volume
    
    # Score multiplier
    score_multipliers = {
        (800, 1000): 2.0,   # Premium: 2x base
        (600, 799): 1.5,    # Standard: 1.5x base
        (400, 599): 1.0,    # Limited: 1x base
        (200, 399): 0.5,    # Restricted: 0.5x base
        (0, 199): 0.0       # No credit
    }
    
    multiplier = get_multiplier_for_score(score.total, score_multipliers)
    
    # Tenure bonus
    tenure_months = customer.tenure_months
    tenure_bonus = min(tenure_months * 0.02, 0.5)  # Max 50% bonus
    
    # Apply multipliers
    calculated_limit = base_limit * multiplier * (1 + tenure_bonus)
    
    # Apply absolute caps by customer type
    max_limits = {
        CustomerType.RETAILER: 500_000,      # ₦500K
        CustomerType.WHOLESALER: 5_000_000,   # ₦5M
        CustomerType.DISTRIBUTOR: 50_000_000  # ₦50M
    }
    
    final_limit = min(calculated_limit, max_limits[customer.type])
    
    # Minimum limits
    min_limits = {
        CustomerType.RETAILER: 10_000,       # ₦10K
        CustomerType.WHOLESALER: 100_000,    # ₦100K
        CustomerType.DISTRIBUTOR: 500_000    # ₦500K
    }
    
    if final_limit < min_limits[customer.type] and score.total >= 400:
        final_limit = min_limits[customer.type]
    
    return CreditLimit(
        amount=final_limit,
        currency="NGN",
        payment_terms=get_payment_terms(score.total),
        valid_from=datetime.now(),
        valid_to=datetime.now() + timedelta(days=90),
        review_date=datetime.now() + timedelta(days=30)
    )


def get_payment_terms(score: int) -> int:
    """Returns payment terms in days"""
    if score >= 800:
        return 30  # Net 30
    elif score >= 600:
        return 14  # Net 14
    elif score >= 400:
        return 7   # Net 7
    else:
        return 0   # COD
```

### 1.3 Credit Monitoring & Alerts

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        CREDIT MONITORING SYSTEM                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   REAL-TIME MONITORING                                                      │
│   ────────────────────                                                      │
│   • Credit utilization tracking                                             │
│   • Payment due date monitoring                                             │
│   • Early warning indicators                                                │
│                                                                              │
│   ALERT TRIGGERS                                                            │
│   ────────────────                                                          │
│   │                                                                         │
│   ├── UTILIZATION ALERTS                                                   │
│   │   • 70% utilized → Notify customer                                     │
│   │   • 85% utilized → Notify customer + account manager                   │
│   │   • 95% utilized → Block new orders, escalate                          │
│   │                                                                         │
│   ├── PAYMENT ALERTS                                                       │
│   │   • 3 days before due → Reminder                                       │
│   │   • Due date → Payment reminder                                        │
│   │   • 1 day past due → Urgent reminder                                   │
│   │   • 7 days past due → Block credit, call                               │
│   │   • 14 days past due → Collection escalation                           │
│   │   • 30 days past due → External collection                             │
│   │                                                                         │
│   ├── BEHAVIORAL ALERTS                                                    │
│   │   • Sudden order volume increase (>200%)                               │
│   │   • Change in product mix (high-value shift)                           │
│   │   • New delivery addresses                                             │
│   │   • Multiple failed payments                                           │
│   │                                                                         │
│   └── RISK ALERTS                                                          │
│       • Score drop >100 points                                             │
│       • Negative external signals                                          │
│       • Customer disputes                                                  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Module 2: Payment Processing

### 2.1 Payment Method Support

| Method | Provider(s) | Use Case | Fees |
|--------|-------------|----------|------|
| **Card Payments** | Paystack, Flutterwave | Online orders, recurring | 1.5% + ₦100 |
| **Bank Transfer** | Direct, Paystack | Large B2B orders | 0.5% capped ₦2000 |
| **Mobile Money** | OPay, PalmPay, Paga | Retail, last-mile | 1.0% |
| **USSD** | Banks, Paystack | Feature phones | 1.5% capped ₦1000 |
| **Wallet** | Internal | Instant, no fees | 0% |
| **Cash Collection** | Gig workers | Retail COD | 2.0% (collection fee) |
| **Trade Credit** | Internal | B2B | Based on terms |
| **POS Terminal** | Bank POS, mPOS | In-store | 0.5% |

### 2.2 Payment Orchestration Logic

```go
// Payment Orchestration Service (Go)

package payment

import (
    "context"
    "fmt"
    "sort"
    
    "github.com/shopspring/decimal"
)

// PaymentOrchestrator handles payment method selection and execution
type PaymentOrchestrator struct {
    providers     map[string]PaymentProvider
    routingRules  []RoutingRule
    fallbackOrder []string
}

// RoutingRule defines conditions for payment routing
type RoutingRule struct {
    Priority      int
    Condition     func(*PaymentRequest) bool
    Provider      string
    SplitPercent  decimal.Decimal // For split payments
}

// ProcessPayment handles a payment request with intelligent routing
func (o *PaymentOrchestrator) ProcessPayment(ctx context.Context, req *PaymentRequest) (*PaymentResult, error) {
    // 1. Validate request
    if err := o.validateRequest(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    // 2. Check for wallet balance (always try wallet first for partial)
    if req.AllowWalletDeduction {
        walletAmount := o.getWalletBalance(ctx, req.CustomerID)
        if walletAmount.GreaterThan(decimal.Zero) {
            deductAmount := decimal.Min(walletAmount, req.Amount)
            if err := o.deductFromWallet(ctx, req.CustomerID, deductAmount); err == nil {
                req.Amount = req.Amount.Sub(deductAmount)
                req.WalletDeducted = deductAmount
                
                if req.Amount.LessThanOrEqual(decimal.Zero) {
                    return &PaymentResult{
                        Status:         PaymentStatusCompleted,
                        WalletDeducted: deductAmount,
                        Provider:       "wallet",
                    }, nil
                }
            }
        }
    }
    
    // 3. Select provider based on routing rules
    provider, err := o.selectProvider(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("no suitable provider: %w", err)
    }
    
    // 4. Execute payment with retry and fallback
    result, err := o.executeWithFallback(ctx, req, provider)
    if err != nil {
        return nil, fmt.Errorf("payment failed: %w", err)
    }
    
    // 5. Handle webhooks asynchronously
    go o.processWebhook(ctx, result)
    
    return result, nil
}

// selectProvider chooses the best payment provider
func (o *PaymentOrchestrator) selectProvider(ctx context.Context, req *PaymentRequest) (string, error) {
    // Sort rules by priority
    rules := make([]RoutingRule, len(o.routingRules))
    copy(rules, o.routingRules)
    sort.Slice(rules, func(i, j int) bool {
        return rules[i].Priority > rules[j].Priority
    })
    
    // Find first matching rule
    for _, rule := range rules {
        if rule.Condition(req) {
            // Check provider health
            if o.providers[rule.Provider].IsHealthy(ctx) {
                return rule.Provider, nil
            }
        }
    }
    
    // Fallback to first healthy provider
    for _, providerName := range o.fallbackOrder {
        if o.providers[providerName].IsHealthy(ctx) {
            return providerName, nil
        }
    }
    
    return "", fmt.Errorf("no healthy providers available")
}

// executeWithFallback attempts payment with fallback providers
func (o *PaymentOrchestrator) executeWithFallback(
    ctx context.Context, 
    req *PaymentRequest, 
    primaryProvider string,
) (*PaymentResult, error) {
    
    // Try primary provider
    result, err := o.providers[primaryProvider].Charge(ctx, req)
    if err == nil && result.Status == PaymentStatusCompleted {
        return result, nil
    }
    
    // Log primary failure
    o.logFailure(ctx, req, primaryProvider, err)
    
    // Try fallback providers if retriable error
    if isRetriableError(err) {
        for _, fallbackName := range o.fallbackOrder {
            if fallbackName == primaryProvider {
                continue
            }
            
            provider := o.providers[fallbackName]
            if !provider.IsHealthy(ctx) || !provider.Supports(req.Method) {
                continue
            }
            
            result, err = provider.Charge(ctx, req)
            if err == nil && result.Status == PaymentStatusCompleted {
                result.FallbackUsed = true
                result.OriginalProvider = primaryProvider
                return result, nil
            }
        }
    }
    
    return nil, fmt.Errorf("all providers failed: %w", err)
}

// Default routing rules
func DefaultRoutingRules() []RoutingRule {
    return []RoutingRule{
        // High-value B2B → Bank transfer (lower fees)
        {
            Priority: 100,
            Condition: func(r *PaymentRequest) bool {
                return r.CustomerType == CustomerTypeB2B && 
                       r.Amount.GreaterThan(decimal.NewFromInt(100000))
            },
            Provider: "bank_transfer",
        },
        // Card payments → Paystack (best success rates)
        {
            Priority: 90,
            Condition: func(r *PaymentRequest) bool {
                return r.Method == PaymentMethodCard
            },
            Provider: "paystack",
        },
        // Mobile money → OPay (best coverage)
        {
            Priority: 80,
            Condition: func(r *PaymentRequest) bool {
                return r.Method == PaymentMethodMobileMoney
            },
            Provider: "opay",
        },
        // USSD → Bank direct (most reliable)
        {
            Priority: 70,
            Condition: func(r *PaymentRequest) bool {
                return r.Method == PaymentMethodUSSD
            },
            Provider: "ussd_direct",
        },
        // Default → Paystack
        {
            Priority: 0,
            Condition: func(r *PaymentRequest) bool {
                return true
            },
            Provider: "paystack",
        },
    }
}
```

### 2.3 Cash Collection Workflow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         CASH COLLECTION WORKFLOW                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ┌─────────────┐                                                           │
│   │   ORDER     │                                                           │
│   │  DELIVERED  │                                                           │
│   └──────┬──────┘                                                           │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                    COLLECTION TASK CREATED                            │  │
│   │  • Invoice amount: ₦50,000                                           │  │
│   │  • Due date: 2025-01-25                                              │  │
│   │  • Customer: ABC Retail Store                                         │  │
│   │  • Location: 123 Market Street, Lagos                                │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                    TASK ASSIGNMENT                                    │  │
│   │                                                                       │  │
│   │  Priority Queue:                                                      │  │
│   │  1. Delivery driver who made delivery (same-trip collection)         │  │
│   │  2. Dedicated collection agent (Gold+ level)                         │  │
│   │  3. Any available gig worker in area                                 │  │
│   │                                                                       │  │
│   │  Assignment Algorithm:                                                │  │
│   │  • Distance from current location                                    │  │
│   │  • Worker rating and collection success rate                         │  │
│   │  • Current task load                                                 │  │
│   │  • Historical success with this customer                             │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                    COLLECTION EXECUTION                               │  │
│   │                                                                       │  │
│   │  Worker App Flow:                                                     │  │
│   │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐   │  │
│   │  │ Accept  │→ │ Navigate│→ │ Check-in│→ │ Collect │→ │ Confirm │   │  │
│   │  │  Task   │  │ to Site │  │ at Site │  │ Payment │  │ Amount  │   │  │
│   │  └─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘   │  │
│   │                                                                       │  │
│   │  Collection Methods:                                                  │  │
│   │  • Cash: Photo of cash, count verification                           │  │
│   │  • Transfer: Initiate transfer to platform account                   │  │
│   │  • POS: Use mobile POS for card payment                              │  │
│   │  • Mobile Money: Generate payment link                               │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                    VERIFICATION & SETTLEMENT                          │  │
│   │                                                                       │  │
│   │  Cash Handling:                                                       │  │
│   │  ┌─────────────────────────────────────────────────────────────────┐ │  │
│   │  │  Option A: Same-Day Bank Deposit                                │ │  │
│   │  │  • Worker deposits at nearest bank branch                       │ │  │
│   │  │  • Uploads deposit slip photo                                   │ │  │
│   │  │  • System matches to collection task                            │ │  │
│   │  │  • Worker paid collection fee upon confirmation                 │ │  │
│   │  ├─────────────────────────────────────────────────────────────────┤ │  │
│   │  │  Option B: Hub Drop-off                                         │ │  │
│   │  │  • Worker drops cash at designated hub                          │ │  │
│   │  │  • Hub manager verifies and receipts                            │ │  │
│   │  │  • Bulk deposit by hub daily                                    │ │  │
│   │  ├─────────────────────────────────────────────────────────────────┤ │  │
│   │  │  Option C: Agent Cash-out                                       │ │  │
│   │  │  • Worker uses OPay/PalmPay agent                               │ │  │
│   │  │  • Digital transfer to platform                                 │ │  │
│   │  │  • Instant reconciliation                                       │ │  │
│   │  └─────────────────────────────────────────────────────────────────┘ │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│          │                                                                   │
│          ▼                                                                   │
│   ┌──────────────────────────────────────────────────────────────────────┐  │
│   │                    RECONCILIATION                                     │  │
│   │                                                                       │  │
│   │  • Match collection to invoice                                       │  │
│   │  • Update customer balance                                           │  │
│   │  • Release credit utilization                                        │  │
│   │  • Calculate worker earnings                                         │  │
│   │  • Update customer credit score                                      │  │
│   └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Module 3: Collections Management

### 3.1 Automated Collection Workflow

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    INTELLIGENT COLLECTIONS WORKFLOW                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   TIME          ACTION                           CHANNEL                    │
│   ─────────────────────────────────────────────────────────────────────     │
│                                                                              │
│   D-3           Pre-due reminder                 SMS + App Push            │
│                 "Your payment of ₦50,000 is                                │
│                  due in 3 days"                                            │
│                                                                              │
│   D-1           Final reminder                   SMS + WhatsApp            │
│                 "Payment due tomorrow.                                      │
│                  Pay now: [link]"                                          │
│                                                                              │
│   D+0           Due date notification            SMS + Call (IVR)          │
│   (Due Day)     "Payment of ₦50,000 is due                                 │
│                  today. Pay now to maintain                                │
│                  your credit limit."                                       │
│                                                                              │
│   D+1           Gentle follow-up                 WhatsApp                  │
│                 "Hi [Name], we noticed your                                │
│                  payment is pending. Need                                  │
│                  help? Reply to chat."                                     │
│                                                                              │
│   D+3           Escalation Level 1              Human Call                 │
│                 Sales rep / account manager                                │
│                 calls to understand situation                              │
│                                                                              │
│   D+7           Credit Freeze                   SMS + App                  │
│                 New orders blocked                                         │
│                 "Your account is on hold due                               │
│                  to overdue payment."                                      │
│                                                                              │
│   D+14          Collection Task Assigned        Gig Worker                 │
│                 Physical visit scheduled                                   │
│                 "A representative will visit                               │
│                  to collect payment."                                      │
│                                                                              │
│   D+21          Escalation Level 2              Manager Call               │
│                 Payment plan discussion                                    │
│                 "Let's discuss a payment                                   │
│                  arrangement."                                             │
│                                                                              │
│   D+30          External Collection             Collection Agency          │
│                 Case transferred to agency                                 │
│                 Internal: Provision for bad                                │
│                 debt begins                                                │
│                                                                              │
│   D+60          Final Notice                    Legal Letter               │
│                 "Final notice before legal                                 │
│                  action"                                                   │
│                                                                              │
│   D+90          Write-off Review               Internal                    │
│                 Case reviewed for write-off                                │
│                 Customer blacklisted                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Payment Plan Management

```go
// Payment Plan Service

type PaymentPlan struct {
    ID              uuid.UUID
    CustomerID      uuid.UUID
    InvoiceIDs      []uuid.UUID
    TotalAmount     decimal.Decimal
    Installments    []Installment
    Status          PlanStatus
    CreatedAt       time.Time
    ApprovedBy      uuid.UUID
}

type Installment struct {
    Number          int
    Amount          decimal.Decimal
    DueDate         time.Time
    PaidAmount      decimal.Decimal
    PaidAt          *time.Time
    Status          InstallmentStatus
}

// CreatePaymentPlan creates a payment plan for overdue invoices
func (s *PaymentPlanService) CreatePaymentPlan(
    ctx context.Context,
    customerID uuid.UUID,
    invoiceIDs []uuid.UUID,
    numInstallments int,
    firstPaymentDate time.Time,
) (*PaymentPlan, error) {
    
    // Calculate total overdue amount
    totalAmount, err := s.calculateTotalOverdue(ctx, invoiceIDs)
    if err != nil {
        return nil, err
    }
    
    // Validate customer eligibility
    eligible, reason := s.checkEligibility(ctx, customerID, totalAmount)
    if !eligible {
        return nil, fmt.Errorf("customer not eligible: %s", reason)
    }
    
    // Generate installments
    installmentAmount := totalAmount.Div(decimal.NewFromInt(int64(numInstallments)))
    installmentAmount = installmentAmount.Round(2)
    
    installments := make([]Installment, numInstallments)
    runningTotal := decimal.Zero
    
    for i := 0; i < numInstallments; i++ {
        dueDate := firstPaymentDate.AddDate(0, 0, i*7) // Weekly installments
        
        amount := installmentAmount
        if i == numInstallments-1 {
            // Last installment gets remainder
            amount = totalAmount.Sub(runningTotal)
        }
        
        installments[i] = Installment{
            Number:  i + 1,
            Amount:  amount,
            DueDate: dueDate,
            Status:  InstallmentStatusPending,
        }
        
        runningTotal = runningTotal.Add(amount)
    }
    
    plan := &PaymentPlan{
        ID:           uuid.New(),
        CustomerID:   customerID,
        InvoiceIDs:   invoiceIDs,
        TotalAmount:  totalAmount,
        Installments: installments,
        Status:       PlanStatusPendingApproval,
        CreatedAt:    time.Now(),
    }
    
    // Save and return
    if err := s.repo.Save(ctx, plan); err != nil {
        return nil, err
    }
    
    // Notify for approval
    s.notifyForApproval(ctx, plan)
    
    return plan, nil
}
```

---

## Module 4: Settlement Engine

### 4.1 Multi-Party Settlement

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      SETTLEMENT DISTRIBUTION FLOW                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   ORDER VALUE: ₦100,000                                                     │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                    SETTLEMENT WATERFALL                              │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   STEP 1: PLATFORM FEES                                                     │
│   ─────────────────────────────────────────────────────────────────────     │
│   │ Payment processing fee (1.5%)          │ ₦1,500  → Payment Provider   │
│   │ Platform transaction fee (0.5%)        │ ₦500    → OmniRoute          │
│   │ Subtotal after platform fees           │ ₦98,000                      │
│   └─────────────────────────────────────────────────────────────────────     │
│                                                                              │
│   STEP 2: LOGISTICS (if applicable)                                         │
│   ─────────────────────────────────────────────────────────────────────     │
│   │ Delivery fee                           │ ₦2,000  → Logistics Partner  │
│   │ Subtotal after logistics               │ ₦96,000                      │
│   └─────────────────────────────────────────────────────────────────────     │
│                                                                              │
│   STEP 3: GIG WORKER PAYMENTS (if applicable)                               │
│   ─────────────────────────────────────────────────────────────────────     │
│   │ Delivery task                          │ ₦500    → Gig Worker         │
│   │ Collection task (2% of collected)      │ ₦2,000  → Gig Worker         │
│   │ Subtotal after gig payments            │ ₦93,500                      │
│   └─────────────────────────────────────────────────────────────────────     │
│                                                                              │
│   STEP 4: FINANCING COSTS (if credit used)                                  │
│   ─────────────────────────────────────────────────────────────────────     │
│   │ Credit facilitation fee (1%)           │ ₦1,000  → Finance Partner    │
│   │ Subtotal after financing               │ ₦92,500                      │
│   └─────────────────────────────────────────────────────────────────────     │
│                                                                              │
│   STEP 5: MERCHANT PAYOUT                                                   │
│   ─────────────────────────────────────────────────────────────────────     │
│   │ Net to Manufacturer/Seller             │ ₦92,500 → Merchant Wallet    │
│   └─────────────────────────────────────────────────────────────────────     │
│                                                                              │
│   SETTLEMENT SCHEDULE:                                                      │
│   • T+0: Platform fees deducted                                             │
│   • T+1: Gig workers paid (instant withdrawal available)                    │
│   • T+1: Logistics partners paid                                            │
│   • T+2: Merchants settled (next business day)                              │
│   • T+3: Finance partners reconciled                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Settlement Configuration

```go
// Settlement Configuration

type SettlementConfig struct {
    TenantID        uuid.UUID
    
    // Schedule
    SettlementFrequency  string // "daily", "weekly", "monthly"
    SettlementDayOfWeek  int    // For weekly: 1-7
    SettlementDayOfMonth int    // For monthly: 1-28
    
    // Thresholds
    MinimumPayout       decimal.Decimal
    HoldbackPercent     decimal.Decimal // Reserve for refunds
    HoldbackDays        int             // Days to hold reserve
    
    // Bank details
    BankName            string
    AccountNumber       string
    AccountName         string
    
    // Splits
    Splits              []SettlementSplit
}

type SettlementSplit struct {
    RecipientType   string  // "merchant", "platform", "partner"
    RecipientID     uuid.UUID
    Percentage      decimal.Decimal
    FixedAmount     decimal.Decimal
    Priority        int
}

// Example: Distributor with 3% platform fee and 2% to manufacturer
func ExampleDistributorConfig() *SettlementConfig {
    return &SettlementConfig{
        SettlementFrequency: "daily",
        MinimumPayout:       decimal.NewFromInt(10000),
        HoldbackPercent:     decimal.NewFromFloat(0.05), // 5% holdback
        HoldbackDays:        7,
        Splits: []SettlementSplit{
            {
                RecipientType: "platform",
                Percentage:    decimal.NewFromFloat(0.03), // 3% platform fee
                Priority:      1,
            },
            {
                RecipientType: "partner",
                RecipientID:   uuid.MustParse("..."), // Manufacturer ID
                Percentage:    decimal.NewFromFloat(0.02), // 2% to manufacturer
                Priority:      2,
            },
            {
                RecipientType: "merchant",
                Percentage:    decimal.NewFromFloat(0.95), // 95% to distributor
                Priority:      3,
            },
        },
    }
}
```

---

## Module 5: External Finance Integration

### 5.1 Finance Partner API

```go
// Finance Partner Integration Interface

type FinancePartner interface {
    // Credit Operations
    RequestCreditLine(ctx context.Context, req *CreditLineRequest) (*CreditLineResponse, error)
    DisburseLoan(ctx context.Context, req *DisbursementRequest) (*DisbursementResponse, error)
    
    // Invoice Financing
    SubmitInvoiceForFinancing(ctx context.Context, invoice *Invoice) (*FinancingOffer, error)
    AcceptFinancingOffer(ctx context.Context, offerID string) (*FinancingConfirmation, error)
    
    // BNPL
    CheckBNPLEligibility(ctx context.Context, customerID uuid.UUID, amount decimal.Decimal) (*BNPLEligibility, error)
    CreateBNPLOrder(ctx context.Context, req *BNPLRequest) (*BNPLResponse, error)
    
    // Reporting
    GetOutstandingBalance(ctx context.Context, customerID uuid.UUID) (decimal.Decimal, error)
    GetRepaymentSchedule(ctx context.Context, loanID string) ([]Repayment, error)
}

// Credit Line Request
type CreditLineRequest struct {
    CustomerID          uuid.UUID
    RequestedAmount     decimal.Decimal
    Purpose             string
    
    // Customer data for underwriting
    CustomerProfile     *CustomerProfile
    TransactionHistory  *TransactionSummary
    PaymentBehavior     *PaymentBehaviorSummary
    
    // Platform assessment
    PlatformCreditScore int
    RecommendedLimit    decimal.Decimal
}

// Invoice Financing Flow
type InvoiceFinancingService struct {
    partners map[string]FinancePartner
}

func (s *InvoiceFinancingService) GetBestOffer(
    ctx context.Context,
    invoice *Invoice,
) (*FinancingOffer, error) {
    
    var bestOffer *FinancingOffer
    
    // Query all partners in parallel
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for name, partner := range s.partners {
        wg.Add(1)
        go func(name string, p FinancePartner) {
            defer wg.Done()
            
            offer, err := p.SubmitInvoiceForFinancing(ctx, invoice)
            if err != nil {
                log.Printf("Partner %s error: %v", name, err)
                return
            }
            
            mu.Lock()
            defer mu.Unlock()
            
            if bestOffer == nil || offer.EffectiveRate.LessThan(bestOffer.EffectiveRate) {
                bestOffer = offer
            }
        }(name, partner)
    }
    
    wg.Wait()
    
    if bestOffer == nil {
        return nil, fmt.Errorf("no financing offers available")
    }
    
    return bestOffer, nil
}
```

### 5.2 BNPL Integration

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         BNPL CHECKOUT FLOW                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│   RETAILER CHECKOUT                                                         │
│   ─────────────────                                                         │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                                                                     │   │
│   │   Order Summary                                                     │   │
│   │   ───────────────                                                   │   │
│   │   Products:         ₦100,000                                        │   │
│   │   Delivery:         ₦2,000                                          │   │
│   │   ─────────────────────────                                         │   │
│   │   Total:            ₦102,000                                        │   │
│   │                                                                     │   │
│   │   Payment Options                                                   │   │
│   │   ───────────────                                                   │   │
│   │   ┌───────────────────────────────────────────────────────────┐    │   │
│   │   │ ○ Pay Now (Bank Transfer / Card)                          │    │   │
│   │   │   Pay ₦102,000 today                                      │    │   │
│   │   └───────────────────────────────────────────────────────────┘    │   │
│   │                                                                     │   │
│   │   ┌───────────────────────────────────────────────────────────┐    │   │
│   │   │ ○ Pay with Credit (Net 14)                                │    │   │
│   │   │   Pay ₦102,000 in 14 days                                 │    │   │
│   │   │   Available credit: ₦250,000                              │    │   │
│   │   └───────────────────────────────────────────────────────────┘    │   │
│   │                                                                     │   │
│   │   ┌───────────────────────────────────────────────────────────┐    │   │
│   │   │ ● Buy Now, Pay Later (4 installments)              ✓      │    │   │
│   │   │                                                           │    │   │
│   │   │   Pay ₦25,500 today                                       │    │   │
│   │   │   Then ₦25,500 weekly × 3                                 │    │   │
│   │   │                                                           │    │   │
│   │   │   Total: ₦102,000 (0% interest)                          │    │   │
│   │   │                                                           │    │   │
│   │   │   Powered by [Partner Logo]                               │    │   │
│   │   └───────────────────────────────────────────────────────────┘    │   │
│   │                                                                     │   │
│   │   [  Complete Order  ]                                             │   │
│   │                                                                     │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   BNPL ELIGIBILITY RULES:                                                   │
│   • Minimum platform tenure: 3 months                                       │
│   • Minimum orders: 10 completed                                            │
│   • Payment history: >80% on-time                                           │
│   • No outstanding defaults                                                 │
│   • Order value: ₦20,000 - ₦500,000                                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Database Schema (Credit & Payment Tables)

```sql
-- Credit Management Tables

CREATE TABLE credit_scores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    
    -- Score Components
    transaction_score INTEGER NOT NULL,
    payment_score INTEGER NOT NULL,
    business_score INTEGER NOT NULL,
    external_score INTEGER NOT NULL,
    total_score INTEGER NOT NULL,
    
    -- Calculation Details
    calculation_details JSONB NOT NULL,
    
    -- Validity
    calculated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until TIMESTAMPTZ NOT NULL,
    
    -- Previous
    previous_score INTEGER,
    score_change INTEGER,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(customer_id, calculated_at)
);

CREATE INDEX idx_credit_scores_customer ON credit_scores(customer_id);
CREATE INDEX idx_credit_scores_date ON credit_scores(calculated_at);

-- Payment Transactions
CREATE TABLE payment_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- References
    order_id UUID REFERENCES orders(id),
    invoice_id UUID REFERENCES invoices(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    
    -- Transaction Details
    transaction_ref VARCHAR(100) UNIQUE NOT NULL,
    amount DECIMAL(15,4) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    
    -- Method & Provider
    payment_method VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    provider_ref VARCHAR(255),
    provider_response JSONB,
    
    -- Status
    status VARCHAR(20) NOT NULL CHECK (status IN (
        'pending', 'processing', 'completed', 'failed', 'refunded', 'cancelled'
    )),
    failure_reason TEXT,
    
    -- Fees
    provider_fee DECIMAL(10,4),
    platform_fee DECIMAL(10,4),
    net_amount DECIMAL(15,4),
    
    -- Collection (for cash/field collection)
    collected_by_user_id UUID REFERENCES users(id),
    collected_by_gig_worker_id UUID REFERENCES gig_workers(id),
    collection_proof JSONB,
    
    -- Timestamps
    initiated_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payment_trans_customer ON payment_transactions(customer_id);
CREATE INDEX idx_payment_trans_order ON payment_transactions(order_id);
CREATE INDEX idx_payment_trans_status ON payment_transactions(status);
CREATE INDEX idx_payment_trans_date ON payment_transactions(initiated_at);

-- Collection Tasks
CREATE TABLE collection_tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Reference
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    
    -- Amount
    expected_amount DECIMAL(15,4) NOT NULL,
    collected_amount DECIMAL(15,4) DEFAULT 0,
    
    -- Assignment
    assigned_to_user_id UUID REFERENCES users(id),
    assigned_to_gig_worker_id UUID REFERENCES gig_workers(id),
    assigned_at TIMESTAMPTZ,
    
    -- Schedule
    scheduled_date DATE NOT NULL,
    scheduled_time_window JSONB,
    
    -- Execution
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN (
        'pending', 'assigned', 'in_progress', 'completed', 
        'partial', 'failed', 'cancelled'
    )),
    
    attempt_count INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMPTZ,
    last_attempt_notes TEXT,
    
    -- Completion
    completed_at TIMESTAMPTZ,
    completion_proof JSONB,
    payment_transaction_id UUID REFERENCES payment_transactions(id),
    
    -- Metadata
    priority INTEGER DEFAULT 5,
    notes TEXT,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_collection_tasks_invoice ON collection_tasks(invoice_id);
CREATE INDEX idx_collection_tasks_customer ON collection_tasks(customer_id);
CREATE INDEX idx_collection_tasks_status ON collection_tasks(status);
CREATE INDEX idx_collection_tasks_date ON collection_tasks(scheduled_date);

-- Settlement Batches
CREATE TABLE settlement_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Batch Details
    batch_number VARCHAR(50) UNIQUE NOT NULL,
    settlement_date DATE NOT NULL,
    
    -- Amounts
    gross_amount DECIMAL(15,4) NOT NULL,
    total_fees DECIMAL(15,4) NOT NULL,
    net_amount DECIMAL(15,4) NOT NULL,
    
    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN (
        'pending', 'processing', 'completed', 'failed', 'partial'
    )),
    
    -- Bank Details
    bank_name VARCHAR(100),
    account_number VARCHAR(50),
    account_name VARCHAR(255),
    
    -- Transfer Details
    transfer_reference VARCHAR(255),
    transfer_date TIMESTAMPTZ,
    
    -- Counts
    transaction_count INTEGER NOT NULL,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE INDEX idx_settlement_batches_date ON settlement_batches(settlement_date);
CREATE INDEX idx_settlement_batches_status ON settlement_batches(status);

-- Payment Plans
CREATE TABLE payment_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    customer_id UUID NOT NULL REFERENCES customers(id),
    
    -- Covered Invoices
    invoice_ids UUID[] NOT NULL,
    
    -- Plan Details
    total_amount DECIMAL(15,4) NOT NULL,
    num_installments INTEGER NOT NULL,
    installment_amount DECIMAL(15,4) NOT NULL,
    
    -- Terms
    first_payment_date DATE NOT NULL,
    frequency VARCHAR(20) NOT NULL, -- weekly, biweekly, monthly
    
    -- Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN (
        'pending', 'approved', 'active', 'completed', 
        'defaulted', 'cancelled'
    )),
    
    -- Approval
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMPTZ,
    
    -- Tracking
    paid_installments INTEGER DEFAULT 0,
    paid_amount DECIMAL(15,4) DEFAULT 0,
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Payment Plan Installments
CREATE TABLE payment_plan_installments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payment_plan_id UUID NOT NULL REFERENCES payment_plans(id),
    
    installment_number INTEGER NOT NULL,
    amount DECIMAL(15,4) NOT NULL,
    due_date DATE NOT NULL,
    
    -- Payment
    paid_amount DECIMAL(15,4) DEFAULT 0,
    paid_at TIMESTAMPTZ,
    payment_transaction_id UUID REFERENCES payment_transactions(id),
    
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN (
        'pending', 'paid', 'partial', 'overdue', 'waived'
    )),
    
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(payment_plan_id, installment_number)
);
```

---

## API Endpoints Summary

```yaml
# Credit & Payment API Endpoints

Credit:
  POST   /api/v1/credit/scores/{customer_id}/calculate  # Calculate credit score
  GET    /api/v1/credit/scores/{customer_id}            # Get current score
  GET    /api/v1/credit/scores/{customer_id}/history    # Score history
  POST   /api/v1/credit/limits/{customer_id}/request    # Request limit increase
  GET    /api/v1/credit/limits/{customer_id}            # Get credit limit
  PUT    /api/v1/credit/limits/{customer_id}            # Update limit (admin)

Payments:
  POST   /api/v1/payments/initiate                      # Start payment
  GET    /api/v1/payments/{payment_id}                  # Get payment status
  POST   /api/v1/payments/{payment_id}/verify           # Verify payment
  POST   /api/v1/payments/webhooks/{provider}           # Provider webhooks

Collections:
  GET    /api/v1/collections/tasks                      # List collection tasks
  POST   /api/v1/collections/tasks                      # Create collection task
  PUT    /api/v1/collections/tasks/{task_id}            # Update task
  POST   /api/v1/collections/tasks/{task_id}/complete   # Complete collection

Invoices:
  GET    /api/v1/invoices                               # List invoices
  GET    /api/v1/invoices/{invoice_id}                  # Get invoice
  POST   /api/v1/invoices/{invoice_id}/remind           # Send reminder
  GET    /api/v1/invoices/aging                         # Aging report

Payment Plans:
  POST   /api/v1/payment-plans                          # Create plan
  GET    /api/v1/payment-plans/{plan_id}                # Get plan
  PUT    /api/v1/payment-plans/{plan_id}/approve        # Approve plan
  POST   /api/v1/payment-plans/{plan_id}/pay            # Make installment payment

Settlements:
  GET    /api/v1/settlements/batches                    # List settlement batches
  GET    /api/v1/settlements/batches/{batch_id}         # Get batch details
  POST   /api/v1/settlements/trigger                    # Trigger settlement (admin)
  GET    /api/v1/settlements/report                     # Settlement report

Wallets:
  GET    /api/v1/wallets/{owner_type}/{owner_id}        # Get wallet
  POST   /api/v1/wallets/{wallet_id}/withdraw           # Request withdrawal
  GET    /api/v1/wallets/{wallet_id}/transactions       # Transaction history
```

---

## Conclusion

The Credit & Payment Management Module provides a comprehensive financial infrastructure that:

1. **Automates Credit Assessment** - Real-time scoring based on platform behavior
2. **Orchestrates Payments** - Intelligent routing across multiple providers
3. **Enables Collections** - Automated reminders with gig worker integration
4. **Handles Settlements** - Multi-party splits with configurable schedules
5. **Integrates Finance Partners** - BNPL, invoice financing, working capital

This module is the **revenue engine** of OmniRoute - every transaction flows through it, enabling monetization while providing value to all stakeholders.

---

*Document Version: 1.0*
*Last Updated: January 2025*
*Technical Owner: Platform Engineering*
