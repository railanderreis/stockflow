package sales

import (
	"errors"
	"time"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition for sales order")
	ErrInsufficientAvailable   = errors.New("insufficient available stock for reservation")
	ErrOrderAlreadyShipped     = errors.New("cannot cancel or modify an order that has already been shipped")
	ErrEmptyOrderItems         = errors.New("sales order must contain at least one item")
	ErrCustomerInactive        = errors.New("customer account is inactive")
)

type OrderStatus string

const (
	StatusDraft     OrderStatus = "DRAFT"
	StatusQuoted    OrderStatus = "QUOTED"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusDelivered OrderStatus = "DELIVERED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type SalesOrderItem struct {
	ID             string    `json:"id"`
	SalesOrderID   string    `json:"sales_order_id"`
	ProductID      string    `json:"product_id"`
	Quantity       float64   `json:"quantity"`
	UnitPriceCents int64     `json:"unit_price_cents"`
	DiscountCents  int64     `json:"discount_cents"`
	TotalCents     int64     `json:"total_cents"`
	CreatedAt      time.Time `json:"created_at"`
}

type SalesOrder struct {
	ID            string            `json:"id"`
	Code          string            `json:"code"`
	CustomerID    string            `json:"customer_id"`
	WarehouseID   string            `json:"warehouse_id"`
	Status        OrderStatus       `json:"status"`
	SubtotalCents int64             `json:"subtotal_cents"`
	DiscountCents int64             `json:"discount_cents"`
	TotalCents    int64             `json:"total_cents"`
	SellerID      string            `json:"seller_id"`
	Notes         string            `json:"notes,omitempty"`
	ConfirmedAt   *time.Time        `json:"confirmed_at,omitempty"`
	ShippedAt     *time.Time        `json:"shipped_at,omitempty"`
	CancelledAt   *time.Time        `json:"cancelled_at,omitempty"`
	Items         []*SalesOrderItem `json:"items"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

func (so *SalesOrder) Confirm() error {
	if so.Status != StatusDraft && so.Status != StatusQuoted {
		return ErrInvalidStatusTransition
	}
	if len(so.Items) == 0 {
		return ErrEmptyOrderItems
	}
	now := time.Now()
	so.Status = StatusConfirmed
	so.ConfirmedAt = &now
	so.UpdatedAt = now
	return nil
}

func (so *SalesOrder) Ship() error {
	if so.Status != StatusConfirmed {
		return ErrInvalidStatusTransition
	}
	now := time.Now()
	so.Status = StatusShipped
	so.ShippedAt = &now
	so.UpdatedAt = now
	return nil
}

func (so *SalesOrder) Cancel() error {
	if so.Status == StatusShipped || so.Status == StatusDelivered {
		return ErrOrderAlreadyShipped
	}
	if so.Status == StatusCancelled {
		return ErrInvalidStatusTransition
	}
	now := time.Now()
	so.Status = StatusCancelled
	so.CancelledAt = &now
	so.UpdatedAt = now
	return nil
}

type SalesRepository interface {
	CreateOrder(so *SalesOrder) error
	UpdateOrder(so *SalesOrder) error
	GetOrderByID(id string) (*SalesOrder, error)
	CreateShipment(tx any, shipment *Shipment) error
}

type ShipmentItem struct {
	ProductID       string  `json:"product_id"`
	QuantityShipped float64 `json:"quantity_shipped"`
}

type Shipment struct {
	ID             string          `json:"id"`
	Code           string          `json:"code"`
	SalesOrderID   string          `json:"sales_order_id"`
	WarehouseID    string          `json:"warehouse_id"`
	InvoiceNumber  string          `json:"invoice_number,omitempty"`
	DispatchedByID string          `json:"dispatched_by_id"`
	ShippedAt      time.Time       `json:"shipped_at"`
	Notes          string          `json:"notes,omitempty"`
	Items          []*ShipmentItem `json:"items"`
}
