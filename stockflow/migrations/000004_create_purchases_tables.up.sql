-- Status Enums
CREATE TYPE requisition_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED', 'CANCELLED');
CREATE TYPE purchase_order_status AS ENUM ('DRAFT', 'ISSUED', 'PARTIALLY_RECEIVED', 'RECEIVED', 'CANCELLED');

-- Purchase Requisitions
CREATE TABLE IF NOT EXISTS purchase_requisitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    requester_id UUID NOT NULL,
    approver_id UUID,
    status requisition_status NOT NULL DEFAULT 'DRAFT',
    total_cents BIGINT NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchase_requisition_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    requisition_id UUID NOT NULL REFERENCES purchase_requisitions(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity NUMERIC(15, 4) NOT NULL CHECK (quantity > 0),
    estimated_unit_cost_cents BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Purchase Orders
CREATE TABLE IF NOT EXISTS purchase_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    requisition_id UUID REFERENCES purchase_requisitions(id),
    supplier_id UUID NOT NULL REFERENCES suppliers(id),
    target_warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    status purchase_order_status NOT NULL DEFAULT 'DRAFT',
    total_cents BIGINT NOT NULL DEFAULT 0 CHECK (total_cents >= 0),
    buyer_id UUID NOT NULL,
    approved_by_id UUID,
    expected_delivery_at TIMESTAMPTZ,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS purchase_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES purchase_orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity_ordered NUMERIC(15, 4) NOT NULL CHECK (quantity_ordered > 0),
    quantity_received NUMERIC(15, 4) NOT NULL DEFAULT 0 CHECK (quantity_received >= 0),
    unit_cost_cents BIGINT NOT NULL CHECK (unit_cost_cents >= 0),
    total_cost_cents BIGINT NOT NULL CHECK (total_cost_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Goods Receipts
CREATE TABLE IF NOT EXISTS goods_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    purchase_order_id UUID NOT NULL REFERENCES purchase_orders(id),
    received_by_id UUID NOT NULL,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    invoice_number VARCHAR(100),
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes TEXT
);

CREATE TABLE IF NOT EXISTS goods_receipt_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goods_receipt_id UUID NOT NULL REFERENCES goods_receipts(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity_received NUMERIC(15, 4) NOT NULL CHECK (quantity_received > 0),
    unit_cost_cents BIGINT NOT NULL CHECK (unit_cost_cents >= 0)
);

CREATE INDEX IF NOT EXISTS idx_po_supplier ON purchase_orders(supplier_id);
CREATE INDEX IF NOT EXISTS idx_po_status ON purchase_orders(status);
CREATE INDEX IF NOT EXISTS idx_goods_receipts_po ON goods_receipts(purchase_order_id);