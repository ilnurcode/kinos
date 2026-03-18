-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS inventory (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    available_quantity INT NOT NULL DEFAULT 0,
    warehouse_location VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_inventory_product_id ON inventory(product_id);
CREATE INDEX idx_inventory_warehouse_location ON inventory(warehouse_location);
CREATE INDEX idx_inventory_available_quantity ON inventory(available_quantity);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_inventory_available_quantity;
DROP INDEX IF EXISTS idx_inventory_warehouse_location;
DROP INDEX IF EXISTS idx_inventory_product_id;
DROP TABLE IF EXISTS inventory;
-- +goose StatementEnd
