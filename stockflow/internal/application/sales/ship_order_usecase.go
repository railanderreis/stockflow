package sales

import (
	"context"
	"fmt"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/inventory"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/sales"
)

type ShipOrderInput struct {
	SalesOrderID   string `json:"sales_order_id"`
	DispatchedByID string `json:"dispatched_by_id"`
	InvoiceNumber  string `json:"invoice_number"`
	Notes          string `json:"notes"`
}

type ShipOrderUseCase struct {
	txManager     TxManager
	salesRepo     sales.SalesRepository
	inventoryRepo inventory.InventoryRepository
}

func NewShipOrderUseCase(
	txManager TxManager,
	salesRepo sales.SalesRepository,
	inventoryRepo inventory.InventoryRepository,
) *ShipOrderUseCase {
	return &ShipOrderUseCase{
		txManager:     txManager,
		salesRepo:     salesRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (uc *ShipOrderUseCase) Execute(ctx context.Context, input ShipOrderInput) (*sales.Shipment, error) {
	so, err := uc.salesRepo.GetOrderByID(input.SalesOrderID)
	if err != nil {
		return nil, fmt.Errorf("sales order not found: %w", err)
	}

	if err := so.Ship(); err != nil {
		return nil, err
	}

	shipment := &sales.Shipment{
		Code:           fmt.Sprintf("SHIP-%d", time.Now().UnixNano()),
		SalesOrderID:   so.ID,
		WarehouseID:    so.WarehouseID,
		InvoiceNumber:  input.InvoiceNumber,
		DispatchedByID: input.DispatchedByID,
		ShippedAt:      time.Now(),
		Notes:          input.Notes,
	}

	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, item := range so.Items {
			// Lock stock row
			stockItem, err := uc.inventoryRepo.GetStockItemForUpdate(txCtx, so.WarehouseID, item.ProductID)
			if err != nil {
				return err
			}

			if stockItem == nil || stockItem.ReservedQuantity < item.Quantity {
				return fmt.Errorf("inconsistent reservation state for product %s", item.ProductID)
			}

			// Deduct physical quantity and clear reservation
			stockItem.Quantity -= item.Quantity
			stockItem.ReservedQuantity -= item.Quantity
			stockItem.UpdatedAt = time.Now()

			if err := uc.inventoryRepo.UpsertStockItem(txCtx, stockItem); err != nil {
				return err
			}

			// Register OUT stock movement ledger entry
			movement := &inventory.StockMovement{
				WarehouseID:   so.WarehouseID,
				ProductID:     item.ProductID,
				Type:          inventory.MovementOut,
				Quantity:      item.Quantity,
				UnitCostCents: item.UnitPriceCents,
				ReferenceType: "SALES_ORDER",
				ReferenceID:   so.Code,
				UserID:        input.DispatchedByID,
				Notes:         fmt.Sprintf("Invoice: %s", input.InvoiceNumber),
				CreatedAt:     time.Now(),
			}
			if err := uc.inventoryRepo.CreateMovement(txCtx, movement); err != nil {
				return err
			}

			shipment.Items = append(shipment.Items, &sales.ShipmentItem{
				ProductID:       item.ProductID,
				QuantityShipped: item.Quantity,
			})
		}

		if err := uc.salesRepo.UpdateOrder(so); err != nil {
			return err
		}

		return uc.salesRepo.CreateShipment(txCtx, shipment)
	})

	if err != nil {
		return nil, err
	}

	return shipment, nil
}
