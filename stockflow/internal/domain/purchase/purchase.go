package purchase

import (
	"errors"
	"time"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition for purchase order")
	ErrRequisitionNotApproved  = errors.New("cannot create purchase order from non-approved requisition")
	ErrExceededOrderedQuantity = errors.New("received quantity exceeds ordered quantity")
	ErrOrderAlreadyClosed      = errors.New("purchase order is already fully received or cancelled")
	ErrEmptyOrderItems         = errors.New("purchase order must contain at least one item")
)

type RequisitionStatus string

const (
	ReqDraft     RequisitionStatus = "DRAFT"
	ReqSubmitted RequisitionStatus = "SUBMITTED"
	ReqApproved  RequisitionStatus = "APPROVED"
	ReqRejected  RequisitionStatus = "REJECTED"
	ReqCancelled RequisitionStatus = "CANCELLED"
)

type OrderStatus string

const (
	OrderDraft             OrderStatus = "DRAFT"
	OrderIssued            OrderStatus = "ISSUED"
	OrderPartiallyReceived OrderStatus = "PARTIALLY_RECEIVED"
	OrderReceived          OrderStatus = "RECEIVED"
	OrderCancelled         OrderStatus = "CANCELLED"
)

type PurchaseOrderItem struct {
	ID               string    `json:"id"`
	OrderID          string    `json:"order_id"`
	ProductID        string    `json:"product_id"`
	QuantityOrdered  float64   `json:"quantity_ordered"`
	QuantityReceived float64   `json:"quantity_received"`
	UnitCostCents    int64     `json:"unit_cost_cents"`
	TotalCostCents   int64     `json:"total_cost_cents"`
	CreatedAt        time.Time `json:"created_at"`
}

type PurchaseOrder struct {
	ID                 string               `json:"id"`
	Code               string               `json:"code"`
	RequisitionID      *string              `json:"requisition_id,omitempty"`
	SupplierID         string               `json:"supplier_id"`
	TargetWarehouseID  string               `json:"target_warehouse_id"`
	Status             OrderStatus          `json:"status"`
	TotalCents         int64                `json:"total_cents"`
	BuyerID            string               `json:"buyer_id"`
	ApprovedByID       *string              `json:"approved_by_id,omitempty"`
	ExpectedDeliveryAt *time.Time           `json:"expected_delivery_at,omitempty"`
	Notes              string               `json:"notes,omitempty"`
	Items              []*PurchaseOrderItem `json:"items"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

func (po *PurchaseOrder) Issue() error {
	if po.Status != OrderDraft {
		return ErrInvalidStatusTransition
	}
	if len(po.Items) == 0 {
		return ErrEmptyOrderItems
	}
	po.Status = OrderIssued
	po.UpdatedAt = time.Now()
	return nil
}

func (po *PurchaseOrder) ReceiveItem(productID string, qty float64) error {
	if po.Status != OrderIssued && po.Status != OrderPartiallyReceived {
		return ErrInvalidStatusTransition
	}

	var found *PurchaseOrderItem
	for _, item := range po.Items {
		if item.ProductID == productID {
			found = item
			break
		}
	}

	if found == nil {
		return errors.New("product not found in purchase order")
	}

	if found.QuantityReceived+qty > found.QuantityOrdered {
		return ErrExceededOrderedQuantity
	}

	found.QuantityReceived += qty

	// Evaluate overall PO status
	allReceived := true
	for _, item := range po.Items {
		if item.QuantityReceived < item.QuantityOrdered {
			allReceived = false
			break
		}
	}

	if allReceived {
		po.Status = OrderReceived
	} else {
		po.Status = OrderPartiallyReceived
	}

	po.UpdatedAt = time.Now()
	return nil
}

type GoodsReceiptItem struct {
	ProductID        string  `json:"product_id"`
	QuantityReceived float64 `json:"quantity_received"`
	UnitCostCents    int64   `json:"unit_cost_cents"`
}

type GoodsReceipt struct {
	ID              string              `json:"id"`
	Code            string              `json:"code"`
	PurchaseOrderID string              `json:"purchase_order_id"`
	ReceivedByID    string              `json:"received_by_id"`
	WarehouseID     string              `json:"warehouse_id"`
	InvoiceNumber   string              `json:"invoice_number,omitempty"`
	ReceivedAt      time.Time           `json:"received_at"`
	Notes           string              `json:"notes,omitempty"`
	Items           []*GoodsReceiptItem `json:"items"`
}

type PurchaseRepository interface {
	CreateOrder(po *PurchaseOrder) error
	UpdateOrder(po *PurchaseOrder) error
	GetOrderByID(id string) (*PurchaseOrder, error)
	CreateGoodsReceipt(tx any, gr *GoodsReceipt) error
}
