CREATE TABLE IF NOT EXISTS products (
                                     product_id BIGSERIAL PRIMARY KEY,
                                     product_name VARCHAR(255) NOT NULL,
                                     manufacturers_id INTEGER NOT NULL REFERENCES manufacturers(manufacturers_id) ON DELETE CASCADE ,
                                     category_id INTEGER NOT NULL REFERENCES category(category_id) ON DELETE CASCADE ,
                                     price FLOAT NOT NULL
);