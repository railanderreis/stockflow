package reports

import "time"

type ABCClass string

const (
	ClassA ABCClass = "A" // Top 80% cumulative value
	ClassB ABCClass = "B" // Next 15% cumulative value (80% - 95%)
	ClassC ABCClass = "C" // Final 5% cumulative value (95% - 100%)
)

type ABCReportItem struct {
	ProductID            string   `json:"product_id"`
	ProductCode          string   `json:"product_code"`
	ProductName          string   `json:"product_name"`
	TotalQuantitySold    float64  `json:"total_quantity_sold"`
	TotalRevenueCents    int64    `json:"total_revenue_cents"`
	RevenuePercentage    float64  `json:"revenue_percentage"`
	CumulativePercentage float64  `json:"cumulative_percentage"`
	Class                ABCClass `json:"class"`
}

type ABCReport struct {
	PeriodStart time.Time        `json:"period_start"`
	PeriodEnd   time.Time        `json:"period_end"`
	TotalVolume int64            `json:"total_volume_cents"`
	Items       []*ABCReportItem `json:"items"`
}

type StockValuationItem struct {
	WarehouseID         string  `json:"warehouse_id"`
	WarehouseName       string  `json:"warehouse_name"`
	ProductID           string  `json:"product_id"`
	ProductCode         string  `json:"product_code"`
	ProductName         string  `json:"product_name"`
	CategoryName        string  `json:"category_name"`
	Quantity            float64 `json:"quantity"`
	ReservedQuantity    float64 `json:"reserved_quantity"`
	AvailableQuantity   float64 `json:"available_quantity"`
	UnitCostCents       int64   `json:"unit_cost_cents"`
	TotalValuationCents int64   `json:"total_valuation_cents"`
}

type ValuationSummary struct {
	TotalSKUs           int64                 `json:"total_skus"`
	TotalQuantity       float64               `json:"total_quantity"`
	TotalValuationCents int64                 `json:"total_valuation_cents"`
	Items               []*StockValuationItem `json:"items"`
}

type ReorderAlertItem struct {
	WarehouseID         string  `json:"warehouse_id"`
	WarehouseName       string  `json:"warehouse_name"`
	ProductID           string  `json:"product_id"`
	ProductCode         string  `json:"product_code"`
	ProductName         string  `json:"product_name"`
	CurrentStock        float64 `json:"current_stock"`
	ReservedStock       float64 `json:"reserved_stock"`
	AvailableStock      float64 `json:"available_stock"`
	MinStock            float64 `json:"min_stock"`
	ReorderPoint        float64 `json:"reorder_point"`
	SuggestedReorderQty float64 `json:"suggested_reorder_qty"`
	SupplierID          string  `json:"supplier_id,omitempty"`
	SupplierName        string  `json:"supplier_name,omitempty"`
}
