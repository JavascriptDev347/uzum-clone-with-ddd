CREATE TYPE order_payment_status AS ENUM ('unpaid', 'paid');
CREATE TYPE order_delivery_status AS ENUM ('preparing', 'handed_to_courier', 'delivered');

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    user_id UUID NULL REFERENCES users(id),
    customer_address TEXT NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    customer_note TEXT NULL,
    payment_status order_payment_status NOT NULL DEFAULT 'unpaid',
    delivery_status order_delivery_status NOT NULL DEFAULT 'preparing',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_items (
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    product_name VARCHAR(255) NOT NULL,
    unit_price_amount NUMERIC(15,2) NOT NULL,
    unit_price_currency VARCHAR(3) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, product_id)
);

CREATE INDEX ON orders(user_id);
CREATE INDEX ON order_items(order_id);
CREATE INDEX ON order_items(product_id);
