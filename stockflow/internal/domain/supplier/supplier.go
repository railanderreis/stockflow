package supplier

import (
	"errors"
	"time"
)

var (
	ErrSupplierNotFound      = errors.New("supplier not found")
	ErrDocumentAlreadyExists = errors.New("supplier document already registered")
	ErrInvalidDocument       = errors.New("invalid supplier document (CNPJ/CPF)")
	ErrInvalidUnitCost       = errors.New("unit cost cents must be greater than or equal to zero")
	ErrInvalidLeadTime       = errors.New("lead time days must be at least 1 day")
)

type Supplier struct {
	ID            string     `json:"id"`
	Document      string     `json:"document"`
	CorporateName string     `json:"corporate_name"`
	TradeName     string     `json:"trade_name"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	IsActive      bool       `json:"is_active"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type SupplierProduct struct {
	SupplierID          string    `json:"supplier_id"`
	ProductID           string    `json:"product_id"`
	SupplierProductCode string    `json:"supplier_product_code"`
	UnitCostCents       int64     `json:"unit_cost_cents"`
	LeadTimeDays        int       `json:"lead_time_days"`
	IsPreferred         bool      `json:"is_preferred"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type SupplierRepository interface {
	Create(s *Supplier) error
	Update(s *Supplier) error
	GetByID(id string) (*Supplier, error)
	GetByDocument(doc string) (*Supplier, error)
	List(page, pageSize int, search string) ([]*Supplier, int, error)

	// Association
	UpsertSupplierProduct(sp *SupplierProduct) error
	GetSupplierProductsByProductID(productID string) ([]*SupplierProduct, error)
}
