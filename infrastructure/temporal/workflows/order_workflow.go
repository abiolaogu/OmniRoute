// Package workflows provides Temporal workflow definitions for OmniRoute.
package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// OrderWorkflowInput contains input for order processing workflow
type OrderWorkflowInput struct {
	OrderID       string
	CustomerID    string
	TenantID      string
	Items         []OrderItem
	TotalAmount   float64
	PaymentMethod string
}

type OrderItem struct {
	ProductID string
	VariantID string
	Quantity  int
	UnitPrice float64
}

// OrderWorkflowResult contains the result of order processing
type OrderWorkflowResult struct {
	OrderID      string
	Status       string
	ShipmentID   string
	WorkerID     string
	CompletedAt  time.Time
	ErrorMessage string
}

// OrderProcessingWorkflow orchestrates the complete order lifecycle
func OrderProcessingWorkflow(ctx workflow.Context, input OrderWorkflowInput) (*OrderWorkflowResult, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting order processing workflow", "orderID", input.OrderID)

	// Configure activity options
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	result := &OrderWorkflowResult{OrderID: input.OrderID}

	// Step 1: Validate and reserve inventory
	var inventoryReserved bool
	err := workflow.ExecuteActivity(ctx, "ReserveInventory", input.OrderID, input.Items).Get(ctx, &inventoryReserved)
	if err != nil {
		result.Status = "failed"
		result.ErrorMessage = "Failed to reserve inventory: " + err.Error()
		return result, nil
	}

	if !inventoryReserved {
		result.Status = "out_of_stock"
		return result, nil
	}

	// Step 2: Process payment
	var paymentSuccess bool
	err = workflow.ExecuteActivity(ctx, "ProcessPayment", input.OrderID, input.TotalAmount, input.PaymentMethod).Get(ctx, &paymentSuccess)
	if err != nil {
		// Release inventory on payment failure
		_ = workflow.ExecuteActivity(ctx, "ReleaseInventory", input.OrderID).Get(ctx, nil)
		result.Status = "payment_failed"
		result.ErrorMessage = "Payment processing failed: " + err.Error()
		return result, nil
	}

	// Step 3: Confirm order
	err = workflow.ExecuteActivity(ctx, "ConfirmOrder", input.OrderID).Get(ctx, nil)
	if err != nil {
		logger.Error("Failed to confirm order", "error", err)
	}

	// Step 4: Assign gig worker
	var workerID string
	err = workflow.ExecuteActivity(ctx, "AssignWorker", input.OrderID, input.TenantID).Get(ctx, &workerID)
	if err != nil {
		logger.Warn("No worker available, order queued", "orderID", input.OrderID)
	}
	result.WorkerID = workerID

	// Step 5: Create shipment
	var shipmentID string
	err = workflow.ExecuteActivity(ctx, "CreateShipment", input.OrderID).Get(ctx, &shipmentID)
	if err != nil {
		logger.Error("Failed to create shipment", "error", err)
	}
	result.ShipmentID = shipmentID

	// Step 6: Send notifications
	_ = workflow.ExecuteActivity(ctx, "SendOrderNotification", input.OrderID, input.CustomerID, "confirmed").Get(ctx, nil)

	// Step 7: Wait for delivery confirmation (with timeout)
	deliveryCtx, cancel := workflow.WithCancel(ctx)
	defer cancel()

	var delivered bool
	selector := workflow.NewSelector(ctx)

	// Signal channel for delivery confirmation
	deliverySignal := workflow.GetSignalChannel(ctx, "delivery_confirmed")
	selector.AddReceive(deliverySignal, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &delivered)
	})

	// Timeout after 48 hours
	timer := workflow.NewTimer(deliveryCtx, 48*time.Hour)
	selector.AddFuture(timer, func(f workflow.Future) {
		delivered = false
	})

	selector.Select(ctx)

	if delivered {
		result.Status = "completed"
		result.CompletedAt = workflow.Now(ctx)
		_ = workflow.ExecuteActivity(ctx, "SendOrderNotification", input.OrderID, input.CustomerID, "delivered").Get(ctx, nil)
	} else {
		result.Status = "pending_delivery"
	}

	return result, nil
}

// CreditReviewWorkflowInput contains input for credit review workflow
type CreditReviewWorkflowInput struct {
	CustomerID      string
	TenantID        string
	RequestedAmount float64
}

// CreditReviewWorkflow handles credit limit review and approval
func CreditReviewWorkflow(ctx workflow.Context, input CreditReviewWorkflowInput) (bool, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting credit review workflow", "customerID", input.CustomerID)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Step 1: Get credit score
	var creditScore int
	err := workflow.ExecuteActivity(ctx, "GetCreditScore", input.CustomerID).Get(ctx, &creditScore)
	if err != nil {
		return false, err
	}

	// Step 2: Auto-approve if score is high enough
	if creditScore >= 700 && input.RequestedAmount <= 1000000 {
		_ = workflow.ExecuteActivity(ctx, "ApproveCredit", input.CustomerID, input.RequestedAmount).Get(ctx, nil)
		_ = workflow.ExecuteActivity(ctx, "SendCreditNotification", input.CustomerID, "approved").Get(ctx, nil)
		return true, nil
	}

	// Step 3: Queue for manual review
	_ = workflow.ExecuteActivity(ctx, "QueueForReview", input.CustomerID, input.RequestedAmount).Get(ctx, nil)

	// Step 4: Wait for manual approval (with timeout)
	approvalSignal := workflow.GetSignalChannel(ctx, "credit_decision")
	var approved bool

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(approvalSignal, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &approved)
	})

	timer := workflow.NewTimer(ctx, 72*time.Hour)
	selector.AddFuture(timer, func(f workflow.Future) {
		approved = false
	})

	selector.Select(ctx)

	if approved {
		_ = workflow.ExecuteActivity(ctx, "ApproveCredit", input.CustomerID, input.RequestedAmount).Get(ctx, nil)
		_ = workflow.ExecuteActivity(ctx, "SendCreditNotification", input.CustomerID, "approved").Get(ctx, nil)
	} else {
		_ = workflow.ExecuteActivity(ctx, "SendCreditNotification", input.CustomerID, "rejected").Get(ctx, nil)
	}

	return approved, nil
}

// PayoutWorkflowInput contains input for worker payout workflow
type PayoutWorkflowInput struct {
	TenantID   string
	PayoutDate string
	Workers    []string
}

// WorkerPayoutWorkflow handles weekly worker payouts
func WorkerPayoutWorkflow(ctx workflow.Context, input PayoutWorkflowInput) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting worker payout workflow", "date", input.PayoutDate)

	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	// Process each worker's payout
	for _, workerID := range input.Workers {
		// Calculate earnings
		var earnings float64
		err := workflow.ExecuteActivity(ctx, "CalculateWorkerEarnings", workerID, input.PayoutDate).Get(ctx, &earnings)
		if err != nil {
			logger.Error("Failed to calculate earnings", "workerID", workerID, "error", err)
			continue
		}

		if earnings <= 0 {
			continue
		}

		// Process bank transfer
		var transferID string
		err = workflow.ExecuteActivity(ctx, "ProcessBankTransfer", workerID, earnings).Get(ctx, &transferID)
		if err != nil {
			logger.Error("Failed to process transfer", "workerID", workerID, "error", err)
			_ = workflow.ExecuteActivity(ctx, "SendPayoutNotification", workerID, "failed").Get(ctx, nil)
			continue
		}

		// Update ledger
		_ = workflow.ExecuteActivity(ctx, "UpdatePayoutLedger", workerID, transferID, earnings).Get(ctx, nil)

		// Send notification
		_ = workflow.ExecuteActivity(ctx, "SendPayoutNotification", workerID, "success").Get(ctx, nil)
	}

	return nil
}
