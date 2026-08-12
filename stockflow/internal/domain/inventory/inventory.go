package inventory

import (
	"errors"
	"time"
)

var (
	ErrWarehouseNotFound     = errors.New("warehouse not found")
	ErrInsufficientStock     = errors.New("insufficient stock for requested operation")
	ErrInvalidQuantity       = errors.New("quantity must be greater than zero")
	ErrSameWarehouseTransfer = errors.New("source and destination warehouses must be different")
	ErrWarehouseInactive     = errors.New("warehouse is inactive")
)

type MovementType string

const (
	MovementIn          MovementType = "IN"
	MovementOut         MovementType = "OUT"
	MovementTransferIn  MovementType = "TRANSFER_IN"
	MovementTransferOut MovementType = "TRANSFER_OUT"
	MovementAdjustment  MovementType = "ADJUSTMENT"
)

type Warehouse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Address   string    `json:"address,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StockItem struct {
	ID               string    `json:"id"`
	WarehouseID      string    `json:"warehouse_id"`
	ProductID        string    `json:"product_id"`
	Quantity         float64   `json:"quantity"`
	ReservedQuantity float64   `json:"reserved_quantity"`
	AvailableQty     float64   `json:"available_quantity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (s *StockItem) CalculateAvailable() float64 {
	return s.Quantity - s.ReservedQuantity
}

type StockMovement struct {
	ID            string       `json:"id"`
	WarehouseID   string       `json:"warehouse_id"`
	ProductID     string       `json:"product_id"`
	Type          MovementType `json:"type"`
	Quantity      float64      `json:"quantity"`
	UnitCostCents int64        `json:"unit_cost_cents"`
	ReferenceType string       `json:"reference_type,omitempty"`
	ReferenceID   string       `json:"reference_id,omitempty"`
	UserID        string       `json:"user_id"`
	Notes         string       `json:"notes,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}

type WarehouseRepository interface {
	Create(w *Warehouse) error
	GetByID(id string) (*Warehouse, error)
	List() ([]*Warehouse, error)
}

type InventoryRepository interface {
	// Pessimistic Locking Query within a transaction
	GetStockItemForUpdate(ctx any, warehouseID, productID string) (*StockItem, error)
	UpsertStockItem(ctx any, item *StockItem) error
	CreateMovement(ctx any, movement *StockMovement) error
	ListStockItems(filter StockFilter) ([]*StockItem, int, error)
	ListMovements(filter MovementFilter) ([]*StockMovement, int, error)
}

type StockFilter struct {
	WarehouseID string
	ProductID   string
	Search      string
	Page        int
	PageSize    int
}

type MovementFilter struct {
	WarehouseID string
	ProductID   string
	Type        MovementType
	Page        int
	PageSize    int
}
