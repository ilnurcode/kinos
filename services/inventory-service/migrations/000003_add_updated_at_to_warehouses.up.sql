-- +goose Up
-- +goose StatementBegin
ALTER TABLE warehouses ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE warehouses DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
