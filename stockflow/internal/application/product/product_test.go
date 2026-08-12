package product_test

import (
	"testing"

	"github.com/railanderreis/stockflow/stockflow/internal/domain/product"
)

func TestProduct_ValidateQuantity(t *testing.T) {
	unitInteger := &product.UnitOfMeasure{
		ID:              "u1",
		Code:            "UN",
		Name:            "Unidade",
		AllowFractional: false,
	}

	unitFractional := &product.UnitOfMeasure{
		ID:              "u2",
		Code:            "KG",
		Name:            "Quilograma",
		AllowFractional: true,
	}

	prodInteger := &product.Product{
		ID:   "p1",
		SKU:  "SKU-001",
		Unit: unitInteger,
	}

	prodFractional := &product.Product{
		ID:   "p2",
		SKU:  "SKU-002",
		Unit: unitFractional,
	}

	// 1. Integer unit with integer quantity -> Valid
	if err := prodInteger.ValidateQuantity(5.0); err != nil {
		t.Errorf("expected no error for integer quantity on UN, got: %v", err)
	}

	// 2. Integer unit with fractional quantity -> Should fail
	if err := prodInteger.ValidateQuantity(5.5); err != product.ErrFractionNotAllowed {
		t.Errorf("expected ErrFractionNotAllowed, got: %v", err)
	}

	// 3. Fractional unit with fractional quantity -> Valid
	if err := prodFractional.ValidateQuantity(2.375); err != nil {
		t.Errorf("expected no error for fractional quantity on KG, got: %v", err)
	}

	// 4. Zero or negative quantity -> Should fail
	if err := prodFractional.ValidateQuantity(0); err == nil {
		t.Error("expected error for zero quantity, got nil")
	}
}
