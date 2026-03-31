-- +goose Down
-- +goose StatementBegin
ALTER TABLE warehouses DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
