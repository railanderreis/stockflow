package purchase

import (
	"context"
	"fmt"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/inventory"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/product"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/purchase"
)

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx context.Context) error) error
}

type ReceiveGoodsInputItem struct {
	ProductID        float64 `json:"-"`
	ProductIDStr     string  `json:"product_id"`
	QuantityReceived float64 `json:"quantity_received"`
	UnitCostCents    int64   `json:"unit_cost_cents"`
}

type ReceiveGoodsInput struct {
	PurchaseOrderID string                  `json:"purchase_order_id"`
	ReceivedByID    string                  `json:"received_by_id"`
	InvoiceNumber   string                  `json:"invoice_number"`
	Notes           string                  `json:"notes"`
	Items           []ReceiveGoodsInputItem `json:"items"`
}

type ReceiveGoodsUseCase struct {
	txManager     TxManager
	purchaseRepo  purchase.PurchaseRepository
	inventoryRepo inventory.InventoryRepository
	productRepo   product.ProductRepository
}

func NewReceiveGoodsUseCase(
	txManager TxManager,
	purchaseRepo purchase.PurchaseRepository,
	inventoryRepo inventory.InventoryRepository,
	productRepo product.ProductRepository,
) *ReceiveGoodsUseCase {
	return &ReceiveGoodsUseCase{
		txManager:     txManager,
		purchaseRepo:  purchaseRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

func (uc *ReceiveGoodsUseCase) Execute(ctx context.Context, input ReceiveGoodsInput) (*purchase.GoodsReceipt, error) {
	po, err := uc.purchaseRepo.GetOrderByID(input.PurchaseOrderID)
	if err != nil {
		return nil, fmt.Errorf("purchase order not found: %w", err)
	}

	if po.Status != purchase.OrderIssued && po.Status != purchase.OrderPartiallyReceived {
		return nil, purchase.ErrInvalidStatusTransition
	}

	receipt := &purchase.GoodsReceipt{
		Code:            fmt.Sprintf("REC-%d", time.Now().UnixNano()),
		PurchaseOrderID: po.ID,
		ReceivedByID:    input.ReceivedByID,
		WarehouseID:     po.TargetWarehouseID,
		InvoiceNumber:   input.InvoiceNumber,
		ReceivedAt:      time.Now(),
		Notes:           input.Notes,
	}

	// Transactional processing: updates PO, registers receipt, and performs stock entry
	err = uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, itemInput := range input.Items {
			if itemInput.QuantityReceived <= 0 {
				continue
			}

			// 1. Validate fractional quantity rules
			prod, err := uc.productRepo.GetByID(itemInput.ProductIDStr)
			if err != nil {
				return fmt.Errorf("product %s not found: %w", itemInput.ProductIDStr, err)
			}
			if err := prod.ValidateQuantity(itemInput.QuantityReceived); err != nil {
				return err
			}

			// 2. Update PO internal item progress
			if err := po.ReceiveItem(itemInput.ProductIDStr, itemInput.QuantityReceived); err != nil {
				return err
			}

			// 3. Lock and update stock level in destination warehouse
			stockItem, err := uc.inventoryRepo.GetStockItemForUpdate(txCtx, po.TargetWarehouseID, itemInput.ProductIDStr)
			if err != nil {
				return err
			}
			if stockItem == nil {
				stockItem = &inventory.StockItem{
					WarehouseID: po.TargetWarehouseID,
					ProductID:   itemInput.ProductIDStr,
					Quantity:    0,
				}
			}

			stockItem.Quantity += itemInput.QuantityReceived
			stockItem.UpdatedAt = time.Now()
			if err := uc.inventoryRepo.UpsertStockItem(txCtx, stockItem); err != nil {
				return err
			}

			// 4. Record IN stock movement ledger entry
			movement := &inventory.StockMovement{
				WarehouseID:   po.TargetWarehouseID,
				ProductID:     itemInput.ProductIDStr,
				Type:          inventory.MovementIn,
				Quantity:      itemInput.QuantityReceived,
				UnitCostCents: itemInput.UnitCostCents,
				ReferenceType: "PURCHASE_ORDER",
				ReferenceID:   po.Code,
				UserID:        input.ReceivedByID,
				Notes:         fmt.Sprintf("Invoice: %s", input.InvoiceNumber),
				CreatedAt:     time.Now(),
			}
			if err := uc.inventoryRepo.CreateMovement(txCtx, movement); err != nil {
				return err
			}

			receipt.Items = append(receipt.Items, &purchase.GoodsReceiptItem{
				ProductID:        itemInput.ProductIDStr,
				QuantityReceived: itemInput.QuantityReceived,
				UnitCostCents:    itemInput.UnitCostCents,
			})
		}

		// 5. Save updated Purchase Order status
		if err := uc.purchaseRepo.UpdateOrder(po); err != nil {
			return err
		}

		// 6. Save Goods Receipt
		return uc.purchaseRepo.CreateGoodsReceipt(txCtx, receipt)
	})

	if err != nil {
		return nil, err
	}

	return receipt, nil
}
