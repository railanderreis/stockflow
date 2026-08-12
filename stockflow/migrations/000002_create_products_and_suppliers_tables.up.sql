-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Brands
CREATE TABLE IF NOT EXISTS brands (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Units of Measure
CREATE TABLE IF NOT EXISTS units_of_measure (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(10) NOT NULL UNIQUE, -- UN, CX, KG, LT, M
    name VARCHAR(50) NOT NULL,
    allow_fractional BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Products
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sku VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    brand_id UUID REFERENCES brands(id),
    unit_id UUID NOT NULL REFERENCES units_of_measure(id),
    min_safety_stock NUMERIC(15, 4) NOT NULL DEFAULT 0 CHECK (min_safety_stock >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products(deleted_at) WHERE deleted_at IS NULL;

-- Suppliers
CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document VARCHAR(20) NOT NULL UNIQUE, -- CNPJ / CPF
    corporate_name VARCHAR(200) NOT NULL,
    trade_name VARCHAR(200),
    email VARCHAR(255),
    phone VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_suppliers_document ON suppliers(document);

-- Supplier Products (Relation & Pricing)
CREATE TABLE IF NOT EXISTS supplier_products (
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    supplier_product_code VARCHAR(100),
    unit_cost_cents BIGINT NOT NULL CHECK (unit_cost_cents >= 0),
    lead_time_days INT NOT NULL DEFAULT 1 CHECK (lead_time_days >= 1),
    is_preferred BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (supplier_id, product_id)
);

-- Initial Units Seed
INSERT INTO units_of_measure (id, code, name, allow_fractional) VALUES
    ('20000000-0000-0000-0000-000000000001', 'UN', 'Unidade', FALSE),
    ('20000000-0000-0000-0000-000000000002', 'CX', 'Caixa', FALSE),
    ('20000000-0000-0000-0000-000000000003', 'KG', 'Quilograma', TRUE),
    ('20000000-0000-0000-0000-000000000004', 'LT', 'Litro', TRUE),
    ('20000000-0000-0000-0000-000000000005', 'M', 'Metro', TRUE)
ON CONFLICT (code) DO NOTHING;