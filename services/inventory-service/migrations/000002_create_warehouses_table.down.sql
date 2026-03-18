-- +goose Down
DROP INDEX IF EXISTS idx_warehouses_name;
DROP TABLE IF EXISTS warehouses;
