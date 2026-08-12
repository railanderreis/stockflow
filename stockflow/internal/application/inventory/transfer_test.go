package inventory_test

import (
	"context"
	"testing"

	appInventory "github.com/railanderreis/stockflow/stockflow/application/inventory"
	"github.com/railanderreis/stockflow/stockflow/domain/inventory"
)

// Mock repositories and TxManager setup for isolation testing
type mockTxManager struct{}

func (m *mockTxManager) WithTransaction(ctx context.Context, fn func(tx context.Context) error) error {
	return fn(ctx)
}

func TestTransferStockUseCase_Validation(t *testing.T) {
	txMgr := &mockTxManager{}

	// Test 1: Same warehouse transfer should fail
	uc := appInventory.NewTransferStockUseCase(txMgr, nil, nil, nil)
	input := appInventory.TransferInput{
		SourceWarehouseID: "wh-1",
		TargetWarehouseID: "wh-1",
		ProductID:         "prod-1",
		Quantity:          10,
	}

	err := uc.Execute(context.Background(), input)
	if err != inventory.ErrSameWarehouseTransfer {
		t.Errorf("expected ErrSameWarehouseTransfer, got %v", err)
	}

	// Test 2: Invalid non-positive quantity should fail
	input.TargetWarehouseID = "wh-2"
	input.Quantity = -5

	err = uc.Execute(context.Background(), input)
	if err != inventory.ErrInvalidQuantity {
		t.Errorf("expected ErrInvalidQuantity, got %v", err)
	}
}
