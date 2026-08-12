package product

import (
	"errors"
	"math"
	"time"
)

var (
	ErrSKUAlreadyExists   = errors.New("product SKU already exists")
	ErrProductNotFound    = errors.New("product not found")
	ErrInvalidSKU         = errors.New("invalid or empty SKU")
	ErrInvalidName        = errors.New("invalid or empty product name")
	ErrFractionNotAllowed = errors.New("fractional quantity is not allowed for this unit of measure")
	ErrInvalidSafetyStock = errors.New("min safety stock must be non-negative")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrUnitNotFound       = errors.New("unit of measure not found")
)

type Category struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Brand struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UnitOfMeasure struct {
	ID              string    `json:"id"`
	Code            string    `json:"code"` // UN, CX, KG, LT, M
	Name            string    `json:"name"`
	AllowFractional bool      `json:"allow_fractional"`
	CreatedAt       time.Time `json:"created_at"`
}

type Product struct {
	ID             string         `json:"id"`
	SKU            string         `json:"sku"`
	Name           string         `json:"name"`
	CategoryID     string         `json:"category_id"`
	CategoryName   string         `json:"category_name,omitempty"`
	BrandID        *string        `json:"brand_id,omitempty"`
	BrandName      *string        `json:"brand_name,omitempty"`
	UnitID         string         `json:"unit_id"`
	Unit           *UnitOfMeasure `json:"unit,omitempty"`
	MinSafetyStock float64        `json:"min_safety_stock"`
	IsActive       bool           `json:"is_active"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

func (p *Product) ValidateQuantity(qty float64) error {
	if qty <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	if p.Unit != nil && !p.Unit.AllowFractional {
		if math.Floor(qty) != qty {
			return ErrFractionNotAllowed
		}
	}
	return nil
}

type ProductRepository interface {
	Create(p *Product) error
	Update(p *Product) error
	GetByID(id string) (*Product, error)
	GetBySKU(sku string) (*Product, error)
	List(filter ProductFilter) ([]*Product, int, error)
	SoftDelete(id string) error

	// Category, Brand & Unit methods
	CreateCategory(cat *Category) error
	ListCategories() ([]*Category, error)
	GetCategoryByID(id string) (*Category, error)

	CreateBrand(brand *Brand) error
	ListBrands() ([]*Brand, error)

	GetUnitByID(id string) (*UnitOfMeasure, error)
	ListUnits() ([]*UnitOfMeasure, error)
}

type ProductFilter struct {
	Search     string
	CategoryID string
	BrandID    string
	IsActive   *bool
	Page       int
	PageSize   int
}
