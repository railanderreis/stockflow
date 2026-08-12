CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID NOT NULL REFERENCES roles(id),
    name VARCHAR(150) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    before_state JSONB,
    after_state JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Inserção de Roles Iniciais do Bounded Context Identity
INSERT INTO roles (id, name, description) VALUES
    ('00000000-0000-0000-0000-000000000001', 'ADMIN', 'Administrador do sistema'),
    ('00000000-0000-0000-0000-000000000002', 'MANAGER', 'Gestor operacional e financeiro'),
    ('00000000-0000-0000-0000-000000000003', 'BUYER', 'Comprador de suprimentos'),
    ('00000000-0000-0000-0000-000000000004', 'WAREHOUSE_OPERATOR', 'Operador de depósito'),
    ('00000000-0000-0000-0000-000000000005', 'FINANCE', 'Analista financeiro'),
    ('00000000-0000-0000-0000-000000000006', 'VIEWER', 'Acesso apenas leitura')
ON CONFLICT (name) DO NOTHING;

-- Inserção de Permissões Básicas
INSERT INTO permissions (id, code, description) VALUES
    ('10000000-0000-0000-0000-000000000001', 'product:read', 'Visualizar produtos'),
    ('10000000-0000-0000-0000-000000000002', 'product:create', 'Cadastrar produtos'),
    ('10000000-0000-0000-0000-000000000003', 'product:update', 'Atualizar produtos'),
    ('10000000-0000-0000-0000-000000000004', 'inventory:read', 'Visualizar estoque'),
    ('10000000-0000-0000-0000-000000000005', 'inventory:adjust', 'Ajustar estoque'),
    ('10000000-0000-0000-0000-000000000006', 'inventory:transfer', 'Transferir estoque entre depósitos'),
    ('10000000-0000-0000-0000-000000000007', 'purchase:read', 'Visualizar pedidos de compra'),
    ('10000000-0000-0000-0000-000000000008', 'purchase:create', 'Criar pedidos de compra'),
    ('10000000-0000-0000-0000-000000000009', 'purchase:approve', 'Aprovar pedidos de compra'),
    ('10000000-0000-0000-0000-000000000010', 'purchase:receive', 'Receber mercadorias de pedidos de compra')
ON CONFLICT (code) DO NOTHING;

-- Associação ADMIN com todas as permissões
INSERT INTO role_permissions (role_id, permission_id)
SELECT '00000000-0000-0000-0000-000000000001', id FROM permissions
ON CONFLICT DO NOTHING;