CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_id TEXT NOT NULL,
    tenant_id TEXT NOT NULL,
    branch_id TEXT NOT NULL,
    name TEXT NOT NULL,
    sku TEXT NOT NULL,
    price NUMERIC NOT NULL,
    category_id TEXT NOT NULL,
    unit TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (tenant_id, branch_id, product_id),
    UNIQUE (tenant_id, branch_id, sku)
);

CREATE INDEX IF NOT EXISTS idx_products_tenant_branch
ON products (tenant_id, branch_id);