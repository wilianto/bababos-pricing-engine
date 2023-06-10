CREATE TABLE price (
    id SERIAL PRIMARY KEY,
    sku_id VARCHAR(16) NOT NULL REFERENCES item(sku_id), 
    supplier_id INT NOT NULL REFERENCES supplier(id),
    price DECIMAL(10,2),
    unit VARCHAR(32),
    stock INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
)