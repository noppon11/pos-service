CREATE TABLE IF NOT EXISTS branches (
    id SERIAL PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    branch_id TEXT NOT NULL,
    branch_name TEXT NOT NULL,
    status TEXT NOT NULL,
    timezone TEXT NOT NULL,
    currency TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (tenant_id, branch_id)
);