-- Audit Logs Table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID NOT NULL,
    actor_email VARCHAR(255) NOT NULL,
    action VARCHAR(100) NOT NULL, -- e.g., "sales_order.confirm", "product.update", "stock.adjust"
    resource_type VARCHAR(100) NOT NULL, -- e.g., "SALES_ORDER", "PRODUCT", "STOCK_ITEM"
    resource_id VARCHAR(255) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_actor ON audit_logs(actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);

-- View for Inventory Valuation per Warehouse & Product
CREATE OR REPLACE VIEW v_stock_valuation AS
SELECT 
    si.warehouse_id,
    w.name AS warehouse_name,
    si.product_id,
    p.code AS product_code,
    p.name AS product_name,
    c.name AS category_name,
    si.quantity,
    si.reserved_quantity,
    (si.quantity - si.reserved_quantity) AS available_quantity,
    p.cost_price_cents,
    CAST(si.quantity * p.cost_price_cents AS BIGINT) AS total_valuation_cents,
    si.updated_at
FROM stock_items si
JOIN products p ON si.product_id = p.id
JOIN warehouses w ON si.warehouse_id = w.id
LEFT JOIN categories c ON p.category_id = c.id
WHERE p.is_active = TRUE AND si.quantity > 0;