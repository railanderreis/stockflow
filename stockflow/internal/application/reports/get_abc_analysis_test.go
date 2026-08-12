package reports_test

import (
	"context"
	"testing"
	"time"

	appReports "github.com/railanderreis/stockflow/stockflow/internal/application/reports"
	"github.com/railanderreis/stockflow/stockflow/internal/domain/reports"
)

type mockABCRepo struct {
	items []*reports.ABCReportItem
}

func (m *mockABCRepo) GetProductSalesSummary(ctx context.Context, start, end time.Time) ([]*reports.ABCReportItem, error) {
	return m.items, nil
}

func TestGetABCAnalysis_CalculationAndClassification(t *testing.T) {
	mockItems := []*reports.ABCReportItem{
		{ProductID: "p1", ProductCode: "A01", TotalRevenueCents: 700000}, // 70% -> Class A
		{ProductID: "p2", ProductCode: "B01", TotalRevenueCents: 150000}, // 15% (cum 85%) -> Class B
		{ProductID: "p3", ProductCode: "C01", TotalRevenueCents: 100000}, // 10% (cum 95%) -> Class B
		{ProductID: "p4", ProductCode: "D01", TotalRevenueCents: 50000},  // 5% (cum 100%) -> Class C
	}

	repo := &mockABCRepo{items: mockItems}
	useCase := appReports.NewGetABCAnalysisUseCase(repo)

	report, err := useCase.Execute(context.Background(), time.Now().AddDate(0, -1, 0), time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.TotalVolume != 1000000 {
		t.Errorf("expected total volume 1000000, got %d", report.TotalVolume)
	}

	expectedClasses := map[string]reports.ABCClass{
		"p1": reports.ClassA,
		"p2": reports.ClassB,
		"p3": reports.ClassB,
		"p4": reports.ClassC,
	}

	for _, item := range report.Items {
		if item.Class != expectedClasses[item.ProductID] {
			t.Errorf("product %s: expected class %s, got %s (cum: %.2f%%)",
				item.ProductID, expectedClasses[item.ProductID], item.Class, item.CumulativePercentage)
		}
	}
}
