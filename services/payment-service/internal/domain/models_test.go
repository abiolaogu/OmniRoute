// Package domain_test contains unit tests for the Payment Service domain models
package domain_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/omniroute/payment-service/internal/domain"
)

func TestMoney_Operations(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		a := domain.Money{Amount: decimal.NewFromFloat(100.00), Currency: "NGN"}
		b := domain.Money{Amount: decimal.NewFromFloat(50.00), Currency: "NGN"}

		result := a.Add(b)
		expected := decimal.NewFromFloat(150.00)

		if !result.Amount.Equal(expected) {
			t.Errorf("Money.Add() = %v, want %v", result.Amount, expected)
		}
	})

	t.Run("Sub", func(t *testing.T) {
		a := domain.Money{Amount: decimal.NewFromFloat(100.00), Currency: "NGN"}
		b := domain.Money{Amount: decimal.NewFromFloat(30.00), Currency: "NGN"}

		result := a.Sub(b)
		expected := decimal.NewFromFloat(70.00)

		if !result.Amount.Equal(expected) {
			t.Errorf("Money.Sub() = %v, want %v", result.Amount, expected)
		}
	})
}

func TestCreditTier_FromScore(t *testing.T) {
	tests := []struct {
		score    int
		expected domain.CreditTier
	}{
		{950, domain.CreditTierPremium},
		{800, domain.CreditTierPremium},
		{750, domain.CreditTierStandard},
		{600, domain.CreditTierStandard},
		{500, domain.CreditTierLimited},
		{400, domain.CreditTierLimited},
		{300, domain.CreditTierRestricted},
		{200, domain.CreditTierRestricted},
		{100, domain.CreditTierNoCredit},
		{0, domain.CreditTierNoCredit},
	}

	for _, tt := range tests {
		t.Run(string(tt.expected), func(t *testing.T) {
			tier := domain.CreditTierFromScore(tt.score)
			if tier != tt.expected {
				t.Errorf("CreditTierFromScore(%d) = %v, want %v", tt.score, tier, tt.expected)
			}
		})
	}
}

func TestCreditLimit_AvailableAmount(t *testing.T) {
	limit := &domain.CreditLimit{
		Amount:         decimal.NewFromFloat(1000000.00),
		UtilizedAmount: decimal.NewFromFloat(250000.00),
	}

	expected := decimal.NewFromFloat(750000.00)
	available := limit.CalculateAvailable()

	if !available.Equal(expected) {
		t.Errorf("CreditLimit.CalculateAvailable() = %v, want %v", available, expected)
	}
}

func TestCreditLimit_CanUtilize(t *testing.T) {
	limit := &domain.CreditLimit{
		Amount:         decimal.NewFromFloat(100000.00),
		UtilizedAmount: decimal.NewFromFloat(80000.00),
		IsActive:       true,
		IsFrozen:       false,
		ValidFrom:      time.Now().Add(-24 * time.Hour),
		ValidTo:        time.Now().Add(24 * time.Hour),
	}

	tests := []struct {
		name     string
		amount   decimal.Decimal
		expected bool
	}{
		{"within limit", decimal.NewFromFloat(15000.00), true},
		{"at limit", decimal.NewFromFloat(20000.00), true},
		{"exceeds limit", decimal.NewFromFloat(25000.00), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := limit.CanUtilize(tt.amount); got != tt.expected {
				t.Errorf("CreditLimit.CanUtilize(%v) = %v, want %v", tt.amount, got, tt.expected)
			}
		})
	}
}

func TestCreditLimit_Frozen(t *testing.T) {
	limit := &domain.CreditLimit{
		Amount:         decimal.NewFromFloat(100000.00),
		UtilizedAmount: decimal.NewFromFloat(0),
		IsActive:       true,
		IsFrozen:       true,
		FreezeReason:   "Payment overdue",
		ValidFrom:      time.Now().Add(-24 * time.Hour),
		ValidTo:        time.Now().Add(24 * time.Hour),
	}

	if limit.CanUtilize(decimal.NewFromFloat(1000.00)) {
		t.Error("Frozen CreditLimit.CanUtilize() should return false")
	}
}

func TestCreditScore_Calculation(t *testing.T) {
	score := &domain.CreditScore{
		TransactionScore: 300, // max 350
		PaymentScore:     320, // max 350
		BusinessScore:    180, // max 200
		ExternalScore:    80,  // max 100
	}

	score.CalculateTotal()

	expectedTotal := 880
	if score.TotalScore != expectedTotal {
		t.Errorf("CreditScore.TotalScore = %v, want %v", score.TotalScore, expectedTotal)
	}

	if score.Tier != domain.CreditTierPremium {
		t.Errorf("CreditScore.Tier = %v, want %v", score.Tier, domain.CreditTierPremium)
	}
}

func TestWallet_AvailableBalance(t *testing.T) {
	wallet := &domain.Wallet{
		Balance:     decimal.NewFromFloat(10000.00),
		HeldBalance: decimal.NewFromFloat(2500.00),
		IsActive:    true,
		IsFrozen:    false,
	}

	expected := decimal.NewFromFloat(7500.00)
	available := wallet.CalculateAvailable()

	if !available.Equal(expected) {
		t.Errorf("Wallet.CalculateAvailable() = %v, want %v", available, expected)
	}
}

func TestWallet_CanDebit(t *testing.T) {
	wallet := &domain.Wallet{
		Balance:     decimal.NewFromFloat(5000.00),
		HeldBalance: decimal.NewFromFloat(1000.00),
		IsActive:    true,
		IsFrozen:    false,
	}

	tests := []struct {
		name     string
		amount   decimal.Decimal
		expected bool
	}{
		{"within balance", decimal.NewFromFloat(3000.00), true},
		{"at limit", decimal.NewFromFloat(4000.00), true},
		{"exceeds available", decimal.NewFromFloat(4500.00), false},
		{"negative amount", decimal.NewFromFloat(-100.00), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wallet.CanDebit(tt.amount); got != tt.expected {
				t.Errorf("Wallet.CanDebit(%v) = %v, want %v", tt.amount, got, tt.expected)
			}
		})
	}
}

func TestPayment_StatusTransitions(t *testing.T) {
	tests := []struct {
		from    domain.PaymentStatus
		to      domain.PaymentStatus
		allowed bool
	}{
		{domain.PaymentStatusPending, domain.PaymentStatusProcessing, true},
		{domain.PaymentStatusProcessing, domain.PaymentStatusCompleted, true},
		{domain.PaymentStatusProcessing, domain.PaymentStatusFailed, true},
		{domain.PaymentStatusCompleted, domain.PaymentStatusRefunded, true},
		{domain.PaymentStatusFailed, domain.PaymentStatusPending, true}, // retry
		{domain.PaymentStatusCompleted, domain.PaymentStatusPending, false},
		{domain.PaymentStatusRefunded, domain.PaymentStatusCompleted, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"_to_"+string(tt.to), func(t *testing.T) {
			payment := &domain.Payment{Status: tt.from}
			if got := payment.CanTransitionTo(tt.to); got != tt.allowed {
				t.Errorf("Payment.CanTransitionTo(%v, %v) = %v, want %v",
					tt.from, tt.to, got, tt.allowed)
			}
		})
	}
}

func TestInvoice_BalanceDue(t *testing.T) {
	invoice := &domain.Invoice{
		TotalAmount: decimal.NewFromFloat(50000.00),
		PaidAmount:  decimal.NewFromFloat(20000.00),
	}

	expected := decimal.NewFromFloat(30000.00)
	balance := invoice.CalculateBalanceDue()

	if !balance.Equal(expected) {
		t.Errorf("Invoice.CalculateBalanceDue() = %v, want %v", balance, expected)
	}
}

func TestInvoice_IsOverdue(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		dueDate  time.Time
		status   domain.InvoiceStatus
		expected bool
	}{
		{
			name:     "not due yet",
			dueDate:  now.Add(7 * 24 * time.Hour),
			status:   domain.InvoiceStatusIssued,
			expected: false,
		},
		{
			name:     "past due",
			dueDate:  now.Add(-7 * 24 * time.Hour),
			status:   domain.InvoiceStatusIssued,
			expected: true,
		},
		{
			name:     "paid invoice",
			dueDate:  now.Add(-7 * 24 * time.Hour),
			status:   domain.InvoiceStatusPaid,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoice := &domain.Invoice{
				DueDate: tt.dueDate,
				Status:  tt.status,
			}
			if got := invoice.IsOverdue(); got != tt.expected {
				t.Errorf("Invoice.IsOverdue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCollectionTask_Completion(t *testing.T) {
	task := &domain.CollectionTask{
		AmountDue:       decimal.NewFromFloat(10000.00),
		AmountCollected: decimal.NewFromFloat(10000.00),
		Status:          domain.CollectionCompleted,
	}

	if !task.IsFullyCollected() {
		t.Error("CollectionTask.IsFullyCollected() should return true")
	}

	// Partial collection
	task.AmountCollected = decimal.NewFromFloat(5000.00)
	if task.IsFullyCollected() {
		t.Error("CollectionTask.IsFullyCollected() should return false for partial")
	}
}

func TestInstallment_Status(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		inst     domain.Installment
		expected string
	}{
		{
			name: "paid installment",
			inst: domain.Installment{
				Amount:     decimal.NewFromFloat(5000.00),
				PaidAmount: decimal.NewFromFloat(5000.00),
				DueDate:    now.Add(-7 * 24 * time.Hour),
				Status:     "paid",
			},
			expected: "paid",
		},
		{
			name: "overdue installment",
			inst: domain.Installment{
				Amount:     decimal.NewFromFloat(5000.00),
				PaidAmount: decimal.NewFromFloat(0),
				DueDate:    now.Add(-7 * 24 * time.Hour),
				Status:     "pending",
			},
			expected: "overdue",
		},
		{
			name: "pending installment",
			inst: domain.Installment{
				Amount:     decimal.NewFromFloat(5000.00),
				PaidAmount: decimal.NewFromFloat(0),
				DueDate:    now.Add(7 * 24 * time.Hour),
				Status:     "pending",
			},
			expected: "pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.inst.CalculateStatus()
			if status != tt.expected {
				t.Errorf("Installment.CalculateStatus() = %v, want %v", status, tt.expected)
			}
		})
	}
}

func TestSettlement_NetAmount(t *testing.T) {
	settlement := &domain.Settlement{
		GrossAmount: decimal.NewFromFloat(100000.00),
		Fees:        decimal.NewFromFloat(1500.00), // 1.5%
	}

	expected := decimal.NewFromFloat(98500.00)
	net := settlement.CalculateNetAmount()

	if !net.Equal(expected) {
		t.Errorf("Settlement.CalculateNetAmount() = %v, want %v", net, expected)
	}
}

func TestCustomerType_CreditLimitMultiplier(t *testing.T) {
	tests := []struct {
		ct       domain.CustomerType
		expected float64
	}{
		{domain.CustomerTypeConsumer, 0.5},
		{domain.CustomerTypeRetailer, 1.0},
		{domain.CustomerTypeWholesaler, 2.0},
		{domain.CustomerTypeDistributor, 3.0},
		{domain.CustomerTypeEnterprise, 5.0},
	}

	for _, tt := range tests {
		t.Run(string(tt.ct), func(t *testing.T) {
			multiplier := tt.ct.CreditLimitMultiplier()
			if multiplier != tt.expected {
				t.Errorf("CustomerType.CreditLimitMultiplier() = %v, want %v",
					multiplier, tt.expected)
			}
		})
	}
}
