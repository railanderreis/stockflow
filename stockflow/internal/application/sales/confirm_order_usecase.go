package sales

import (
	"context"
	"fmt"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/inventory"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/product"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/sales"
)

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(tx context.Context) error) error
}

type ConfirmOrderUseCase struct {
	txManager     TxManager
	salesRepo     sales.SalesRepository
	inventoryRepo inventory.InventoryRepository
	productRepo   product.ProductRepository
}

func NewConfirmOrderUseCase(
	txManager TxManager,
	salesRepo sales.SalesRepository,
	inventoryRepo inventory.InventoryRepository,
	productRepo product.ProductRepository,
) *ConfirmOrderUseCase {
	return &ConfirmOrderUseCase{
		txManager:     txManager,
		salesRepo:     salesRepo,
		inventoryRepo: inventoryRepo,
		productRepo:   productRepo,
	}
}

func (uc *ConfirmOrderUseCase) Execute(ctx context.Context, orderID string) error {
	so, err := uc.salesRepo.GetOrderByID(orderID)
	if err != nil {
		return fmt.Errorf("sales order not found: %w", err)
	}

	if err := so.Confirm(); err != nil {
		return err
	}

	// Transactional reservation of stock
	return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		for _, item := range so.Items {
			// Validate fractional rules
			prod, err := uc.productRepo.GetByID(item.ProductID)
			if err != nil {
				return fmt.Errorf("product %s not found: %w", item.ProductID, err)
			}
			if err := prod.ValidateQuantity(item.Quantity); err != nil {
				return err
			}

			// Lock row for pessimistic update
			stockItem, err := uc.inventoryRepo.GetStockItemForUpdate(txCtx, so.WarehouseID, item.ProductID)
			if err != nil {
				return err
			}

			if stockItem == nil || stockItem.CalculateAvailable() < item.Quantity {
				return fmt.Errorf("%w for product %s (requested: %.4f, available: %.4f)",
					sales.ErrInsufficientAvailable, item.ProductID, item.Quantity, stockItem.CalculateAvailable())
			}

			// Reserve stock
			stockItem.ReservedQuantity += item.Quantity
			if err := uc.inventoryRepo.UpsertStockItem(txCtx, stockItem); err != nil {
				return err
			}
		}

		return uc.salesRepo.UpdateOrder(so)
	})
}
