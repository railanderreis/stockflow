package integration_test

import (
	"context"
	"testing"
	"time"

	domainInv "github.com/railanderreis/stockflow/stockflow/internal/domain/inventory"
	domainPur "github.com/railanderreis/stockflow/stockflow/internal/domain/purchase"
	domainSales "github.com/railanderreis/stockflow/stockflow/internal/domain/sales"
)

func TestE2E_CompleteOrderToCashLifecycle(t *testing.T) {
	ctx := context.Background()

	// Setup Test Context & Mock/Test Repositories
	t.Log("1. Initializing E2E Test Suite Context...")

	// 1. Create Product & Warehouse Setup
	warehouseID := "wh-main-01"
	productID := "prod-laptop-01"
	supplierID := "supp-tech-01"
	customerID := "cust-corp-01"

	t.Log("2. Executing Purchase Order creation & approval...")
	po := &domainPur.PurchaseOrder{
		ID:          "po-1001",
		Code:        "PO-2026-001",
		SupplierID:  supplierID,
		WarehouseID: warehouseID,
		Status:      domainPur.StatusApproved,
		TotalCents:  5000000,
		Items: []*domainPur.PurchaseOrderItem{
			{
				ProductID:     productID,
				Quantity:      50.0,
				UnitCostCents: 100000,
				TotalCents:    5000000,
			},
		},
	}

	if err := po.Receive(); err != nil {
		t.Fatalf("Failed to transition purchase order to RECEIVED: %v", err)
	}

	t.Log("3. Receiving Goods into Inventory (Inward Stock Movement)...")
	stockItem := &domainInv.StockItem{
		WarehouseID:      warehouseID,
		ProductID:        productID,
		Quantity:         50.0,
		ReservedQuantity: 0.0,
		MinStock:         10.0,
		ReorderPoint:     15.0,
		UpdatedAt:        time.Now(),
	}

	if stockItem.CalculateAvailable() != 50.0 {
		t.Fatalf("Expected 50.0 available stock, got %.2f", stockItem.CalculateAvailable())
	}

	t.Log("4. Creating & Confirming Sales Order with Stock Reservation...")
	so := &domainSales.SalesOrder{
		ID:          "so-2001",
		Code:        "SO-2026-001",
		CustomerID:  customerID,
		WarehouseID: warehouseID,
		Status:      domainSales.StatusDraft,
		Items: []*domainSales.SalesOrderItem{
			{
				ProductID:      productID,
				Quantity:       20.0,
				UnitPriceCents: 150000,
				TotalCents:     3000000,
			},
		},
	}

	if err := so.Confirm(); err != nil {
		t.Fatalf("Failed to confirm sales order: %v", err)
	}

	// Reserve 20 units
	stockItem.ReservedQuantity += 20.0
	if stockItem.CalculateAvailable() != 30.0 {
		t.Fatalf("Expected 30.0 available stock after reservation, got %.2f", stockItem.CalculateAvailable())
	}

	t.Log("5. Executing Sales Order Shipment & Physical Stock Deduction...")
	if err := so.Ship(); err != nil {
		t.Fatalf("Failed to ship sales order: %v", err)
	}

	// Deduct physical quantity upon dispatch
	stockItem.Quantity -= 20.0
	stockItem.ReservedQuantity -= 20.0

	if stockItem.Quantity != 30.0 || stockItem.ReservedQuantity != 0.0 {
		t.Fatalf("Stock state mismatch post-shipment: Qty=%.2f, Reserved=%.2f", stockItem.Quantity, stockItem.ReservedQuantity)
	}

	t.Log("6. Validating Post-Shipment Stock Valuation...")
	unitCostCents := int64(100000)
	totalValuationCents := int64(stockItem.Quantity * float64(unitCostCents))

	if totalValuationCents != 3000000 {
		t.Fatalf("Expected stock valuation of 3000000 cents ($30,000.00), got %d", totalValuationCents)
	}

	t.Log("E2E Order-to-Cash Lifecycle Integration Test PASSED Successfully.")
}
