CREATE TABLE rfq_item (
    id SERIAL PRIMARY KEY,
    rfq_id INT NOT NULL REFERENCES rfq(id),
    sku_id VARCHAR(16) NOT NULL REFERENCES item(sku_id),
    qty INT,
    unit VARCHAR(32),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)