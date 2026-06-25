CREATE TABLE IF NOT EXISTS products (
    sku VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    current_stock INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS movements (
    event_id VARCHAR(50) PRIMARY KEY,
    sku VARCHAR(50) NOT NULL REFERENCES products(sku),
    type VARCHAR(10) NOT NULL CHECK (type IN ('IN', 'OUT')),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Index for the API to query history efficiently
CREATE INDEX IF NOT EXISTS idx_movements_sku_time ON movements(sku, occurred_at DESC);
