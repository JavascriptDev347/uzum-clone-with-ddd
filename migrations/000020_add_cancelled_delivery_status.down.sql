ALTER TABLE orders ALTER COLUMN delivery_status DROP DEFAULT;

ALTER TYPE order_delivery_status RENAME TO order_delivery_status_old;

CREATE TYPE order_delivery_status AS ENUM ('preparing', 'handed_to_courier', 'delivered');

ALTER TABLE orders
    ALTER COLUMN delivery_status TYPE order_delivery_status
    USING delivery_status::text::order_delivery_status;

ALTER TABLE orders ALTER COLUMN delivery_status SET DEFAULT 'preparing';

DROP TYPE order_delivery_status_old;
