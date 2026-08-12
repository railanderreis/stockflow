package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/inventory"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/product"
)

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx context.Context) error) error
}

type TransferInput struct {
	SourceWarehouseID string  `json:"source_warehouse_id"`
	TargetWarehouseID string  `json:"target_warehouse_id"`
	ProductID         string  `json:"product_id"`
	Quantity          float64 `json:"quantity"`
	UserID            string  `json:"user_id"`
	Notes             string  `json:"notes"`
}

type TransferStockUseCase struct {
	txManager     TxManager
	inventoryRepo inventory.InventoryRepository
	warehouseRepo inventory.WarehouseRepository
	productRepo   product.ProductRepository
}

func NewTransferStockUseCase(
	txManager TxManager,
	inventoryRepo inventory.InventoryRepository,
	warehouseRepo inventory.WarehouseRepository,
	productRepo product.ProductRepository,
) *TransferStockUseCase {
	return &TransferStockUseCase{
		txManager:     txManager,
		inventoryRepo: inventoryRepo,
		warehouseRepo: warehouseRepo,
		productRepo:   productRepo,
	}
}

func (uc *TransferStockUseCase) Execute(ctx context.Context, input TransferInput) error {
	if input.SourceWarehouseID == input.TargetWarehouseID {
		return inventory.ErrSameWarehouseTransfer
	}

	if input.Quantity <= 0 {
		return inventory.ErrInvalidQuantity
	}

	// Fetch product and validate fractional quantity rules
	prod, err := uc.productRepo.GetByID(input.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if err := prod.ValidateQuantity(input.Quantity); err != nil {
		return err
	}

	// Verify source and target warehouses
	srcWh, err := uc.warehouseRepo.GetByID(input.SourceWarehouseID)
	if err != nil || !srcWh.IsActive {
		return inventory.ErrWarehouseNotFound
	}

	tgtWh, err := uc.warehouseRepo.GetByID(input.TargetWarehouseID)
	if err != nil || !tgtWh.IsActive {
		return inventory.ErrWarehouseNotFound
	}

	// Execute within a database transaction with SELECT FOR UPDATE
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		// 1. Lock and validate source stock
		srcItem, err := uc.inventoryRepo.GetStockItemForUpdate(txCtx, input.SourceWarehouseID, input.ProductID)
		if err != nil {
			return err
		}

		if srcItem == nil || srcItem.CalculateAvailable() < input.Quantity {
			return inventory.ErrInsufficientStock
		}

		// 2. Lock target stock
		tgtItem, err := uc.inventoryRepo.GetStockItemForUpdate(txCtx, input.TargetWarehouseID, input.ProductID)
		if err != nil {
			return err
		}
		if tgtItem == nil {
			tgtItem = &inventory.StockItem{
				WarehouseID:      input.TargetWarehouseID,
				ProductID:        input.ProductID,
				Quantity:         0,
				ReservedQuantity: 0,
			}
		}

		// 3. Update source stock
		srcItem.Quantity -= input.Quantity
		srcItem.UpdatedAt = time.Now()
		if err := uc.inventoryRepo.UpsertStockItem(txCtx, srcItem); err != nil {
			return err
		}

		// 4. Update target stock
		tgtItem.Quantity += input.Quantity
		tgtItem.UpdatedAt = time.Now()
		if err := uc.inventoryRepo.UpsertStockItem(txCtx, tgtItem); err != nil {
			return err
		}

		// 5. Record movement out from source
		moveOut := &inventory.StockMovement{
			WarehouseID:   input.SourceWarehouseID,
			ProductID:     input.ProductID,
			Type:          inventory.MovementTransferOut,
			Quantity:      input.Quantity,
			ReferenceType: "TRANSFER",
			ReferenceID:   fmt.Sprintf("%s->%s", input.SourceWarehouseID, input.TargetWarehouseID),
			UserID:        input.UserID,
			Notes:         input.Notes,
			CreatedAt:     time.Now(),
		}
		if err := uc.inventoryRepo.CreateMovement(txCtx, moveOut); err != nil {
			return err
		}

		// 6. Record movement in to target
		moveIn := &inventory.StockMovement{
			WarehouseID:   input.TargetWarehouseID,
			ProductID:     input.ProductID,
			Type:          inventory.MovementTransferIn,
			Quantity:      input.Quantity,
			ReferenceType: "TRANSFER",
			ReferenceID:   fmt.Sprintf("%s->%s", input.SourceWarehouseID, input.TargetWarehouseID),
			UserID:        input.UserID,
			Notes:         input.Notes,
			CreatedAt:     time.Now(),
		}
		return uc.inventoryRepo.CreateMovement(txCtx, moveIn)
	})
}
