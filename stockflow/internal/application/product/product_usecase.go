package product

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/product"
)

type CreateProductInput struct {
	SKU            string  `json:"sku"`
	Name           string  `json:"name"`
	CategoryID     string  `json:"category_id"`
	BrandID        *string `json:"brand_id"`
	UnitID         string  `json:"unit_id"`
	MinSafetyStock float64 `json:"min_safety_stock"`
}

type CreateProductUseCase struct {
	repo product.ProductRepository
}

func NewCreateProductUseCase(repo product.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{repo: repo}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (*product.Product, error) {
	sku := strings.TrimSpace(input.SKU)
	if sku == "" {
		return nil, product.ErrInvalidSKU
	}

	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, product.ErrInvalidName
	}

	if input.MinSafetyStock < 0 {
		return nil, product.ErrInvalidSafetyStock
	}

	existing, err := uc.repo.GetBySKU(sku)
	if err == nil && existing != nil {
		return nil, product.ErrSKUAlreadyExists
	}

	unit, err := uc.repo.GetUnitByID(input.UnitID)
	if err != nil || unit == nil {
		return nil, product.ErrUnitNotFound
	}

	category, err := uc.repo.GetCategoryByID(input.CategoryID)
	if err != nil || category == nil {
		return nil, product.ErrCategoryNotFound
	}

	p := &product.Product{
		SKU:            sku,
		Name:           name,
		CategoryID:     input.CategoryID,
		BrandID:        input.BrandID,
		UnitID:         input.UnitID,
		Unit:           unit,
		MinSafetyStock: input.MinSafetyStock,
		IsActive:       true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Create(p); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return p, nil
}
