package reports

import (
	"context"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/reports"
)

type ReorderAlertsRepository interface {
	GetItemsBelowReorderPoint(ctx context.Context, warehouseID string) ([]*reports.ReorderAlertItem, error)
}

type GetReorderAlertsUseCase struct {
	repo ReorderAlertsRepository
}

func NewGetReorderAlertsUseCase(repo ReorderAlertsRepository) *GetReorderAlertsUseCase {
	return &GetReorderAlertsUseCase{repo: repo}
}

func (uc *GetReorderAlertsUseCase) Execute(ctx context.Context, warehouseID string) ([]*reports.ReorderAlertItem, error) {
	items, err := uc.repo.GetItemsBelowReorderPoint(ctx, warehouseID)
	if err != nil {
		return nil, err
	}

	// Calculate suggested purchase reorder quantities
	for _, item := range items {
		deficit := item.ReorderPoint - item.AvailableStock
		if deficit > 0 {
			item.SuggestedReorderQty = deficit + item.MinStock
		} else {
			item.SuggestedReorderQty = item.MinStock
		}
	}

	return items, nil
}
