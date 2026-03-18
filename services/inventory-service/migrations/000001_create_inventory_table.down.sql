-- +goose Down
DROP INDEX IF EXISTS idx_inventory_available_quantity;
DROP INDEX IF EXISTS idx_inventory_warehouse_location;
DROP INDEX IF EXISTS idx_inventory_product_id;
DROP TABLE IF EXISTS inventory;
