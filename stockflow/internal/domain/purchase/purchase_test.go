package purchase_test

import (
	"testing"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/purchase"
)

func TestPurchaseOrder_Lifecycle(t *testing.T) {
	po := &purchase.PurchaseOrder{
		ID:     "po-1",
		Code:   "PO-2026-001",
		Status: purchase.OrderDraft,
		Items: []*purchase.PurchaseOrderItem{
			{
				ProductID:        "p1",
				QuantityOrdered:  100,
				QuantityReceived: 0,
			},
		},
	}

	// 1. Cannot receive DRAFT order
	if err := po.ReceiveItem("p1", 50); err != purchase.ErrInvalidStatusTransition {
		t.Errorf("expected ErrInvalidStatusTransition, got %v", err)
	}

	// 2. Issue Order
	if err := po.Issue(); err != nil {
		t.Fatalf("unexpected error issuing order: %v", err)
	}
	if po.Status != purchase.OrderIssued {
		t.Errorf("expected status OrderIssued, got %s", po.Status)
	}

	// 3. Partial Receive
	if err := po.ReceiveItem("p1", 40); err != nil {
		t.Fatalf("unexpected error receiving items: %v", err)
	}
	if po.Status != purchase.OrderPartiallyReceived {
		t.Errorf("expected status OrderPartiallyReceived, got %s", po.Status)
	}

	// 4. Over-receiving should fail
	if err := po.ReceiveItem("p1", 70); err != purchase.ErrExceededOrderedQuantity {
		t.Errorf("expected ErrExceededOrderedQuantity, got %v", err)
	}

	// 5. Complete Remaining Receive
	if err := po.ReceiveItem("p1", 60); err != nil {
		t.Fatalf("unexpected error completing receive: %v", err)
	}
	if po.Status != purchase.OrderReceived {
		t.Errorf("expected status OrderReceived, got %s", po.Status)
	}
}
