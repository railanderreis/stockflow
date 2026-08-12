package sales_test

import (
	"testing"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/sales"
)

func TestSalesOrder_Lifecycle(t *testing.T) {
	so := &sales.SalesOrder{
		ID:     "so-1",
		Code:   "SO-2026-001",
		Status: sales.StatusDraft,
		Items: []*sales.SalesOrderItem{
			{
				ProductID:      "p1",
				Quantity:       10,
				UnitPriceCents: 1500,
				TotalCents:     15000,
			},
		},
	}

	// 1. Direct ship from DRAFT must fail
	if err := so.Ship(); err != sales.ErrInvalidStatusTransition {
		t.Errorf("expected ErrInvalidStatusTransition on direct ship, got %v", err)
	}

	// 2. Confirm Order
	if err := so.Confirm(); err != nil {
		t.Fatalf("unexpected error confirming order: %v", err)
	}
	if so.Status != sales.StatusConfirmed {
		t.Errorf("expected status CONFIRMED, got %s", so.Status)
	}

	// 3. Ship Confirmed Order
	if err := so.Ship(); err != nil {
		t.Fatalf("unexpected error shipping order: %v", err)
	}
	if so.Status != sales.StatusShipped {
		t.Errorf("expected status SHIPPED, got %s", so.Status)
	}

	// 4. Cancel Shipped Order must fail
	if err := so.Cancel(); err != sales.ErrOrderAlreadyShipped {
		t.Errorf("expected ErrOrderAlreadyShipped, got %v", err)
	}
}
