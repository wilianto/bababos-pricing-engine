CREATE TABLE po_item (
    id SERIAL PRIMARY KEY,
    po_id INT NOT NULL REFERENCES po(id),
    sku_id VARCHAR(16) NOT NULL REFERENCES item(sku_id),
    qty INT,
    unit VARCHAR(32),
    price DECIMAL(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)