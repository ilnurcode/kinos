CREATE TABLE IF NOT EXISTS inventory_reservations (
    id BIGSERIAL PRIMARY KEY,
    reservation_id VARCHAR(255) NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (reservation_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_reservations_reservation_id
    ON inventory_reservations(reservation_id);
