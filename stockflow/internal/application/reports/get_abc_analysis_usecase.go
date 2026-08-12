package reports

import (
	"context"
	"sort"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/reports"
)

type ABCReportRepository interface {
	GetProductSalesSummary(ctx context.Context, start, end time.Time) ([]*reports.ABCReportItem, error)
}

type GetABCAnalysisUseCase struct {
	repo ABCReportRepository
}

func NewGetABCAnalysisUseCase(repo ABCReportRepository) *GetABCAnalysisUseCase {
	return &GetABCAnalysisUseCase{repo: repo}
}

func (uc *GetABCAnalysisUseCase) Execute(ctx context.Context, start, end time.Time) (*reports.ABCReport, error) {
	rawItems, err := uc.repo.GetProductSalesSummary(ctx, start, end)
	if err != nil {
		return nil, err
	}

	var totalRevenue int64
	for _, item := range rawItems {
		totalRevenue += item.TotalRevenueCents
	}

	if totalRevenue == 0 {
		return &reports.ABCReport{
			PeriodStart: start,
			PeriodEnd:   end,
			TotalVolume: 0,
			Items:       []*reports.ABCReportItem{},
		}, nil
	}

	// Sort descending by revenue
	sort.Slice(rawItems, func(i, j int) bool {
		return rawItems[i].TotalRevenueCents > rawItems[j].TotalRevenueCents
	})

	var cumulativeRevenue int64
	for _, item := range rawItems {
		item.RevenuePercentage = (float64(item.TotalRevenueCents) / float64(totalRevenue)) * 100.0
		cumulativeRevenue += item.TotalRevenueCents
		item.CumulativePercentage = (float64(cumulativeRevenue) / float64(totalRevenue)) * 100.0

		// Categorize based on cumulative percentage thresholds
		switch {
		case item.CumulativePercentage <= 80.0 || (len(rawItems) == 1):
			item.Class = reports.ClassA
		case item.CumulativePercentage <= 95.0:
			item.Class = reports.ClassB
		default:
			item.Class = reports.ClassC
		}
	}

	return &reports.ABCReport{
		PeriodStart: start,
		PeriodEnd:   end,
		TotalVolume: totalRevenue,
		Items:       rawItems,
	}, nil
}
