-- Warehouses
CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Stock Items (Current balance per Warehouse and Product)
CREATE TABLE IF NOT EXISTS stock_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity NUMERIC(15, 4) NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    reserved_quantity NUMERIC(15, 4) NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uk_warehouse_product UNIQUE (warehouse_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_stock_items_warehouse ON stock_items(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_stock_items_product ON stock_items(product_id);

-- Movement Types Enum
CREATE TYPE movement_type AS ENUM ('IN', 'OUT', 'TRANSFER_IN', 'TRANSFER_OUT', 'ADJUSTMENT');

-- Stock Movements (Immutable Ledger)
CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    product_id UUID NOT NULL REFERENCES products(id),
    type movement_type NOT NULL,
    quantity NUMERIC(15, 4) NOT NULL CHECK (quantity > 0),
    unit_cost_cents BIGINT DEFAULT 0,
    reference_type VARCHAR(50), -- 'PURCHASE_ORDER', 'SALES_ORDER', 'MANUAL', 'TRANSFER'
    reference_id VARCHAR(100),
    user_id UUID NOT NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_stock_movements_warehouse_product ON stock_movements(warehouse_id, product_id);
CREATE INDEX IF NOT EXISTS idx_stock_movements_created_at ON stock_movements(created_at DESC);

-- Default Warehouse Seed
INSERT INTO warehouses (id, code, name, is_active) VALUES
    ('30000000-0000-0000-0000-000000000001', 'WH-MAIN', 'Depósito Principal', TRUE),
    ('30000000-0000-0000-0000-000000000002', 'WH-SEC', 'Depósito Secundário', TRUE)
ON CONFLICT (code) DO NOTHING;