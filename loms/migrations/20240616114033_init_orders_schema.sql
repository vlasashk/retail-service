-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS orders;

CREATE TYPE orders.order_status AS ENUM (
  'Unknown',
  'New',
  'AwaitingPayment',
  'Payed',
  'Cancelled',
  'Failed'
);

CREATE TABLE IF NOT EXISTS orders.orders (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  status orders.order_status NOT NULL DEFAULT 'Unknown',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders.order_items (
  sku_id BIGINT NOT NULL,
  order_id BIGINT REFERENCES orders.orders(id) ON DELETE CASCADE NOT NULL,
  count BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON orders.order_items(order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS orders.idx_order_items_order_id;

DROP TABLE IF EXISTS orders.order_items;
DROP TABLE IF EXISTS orders.orders;
DROP TYPE IF EXISTS orders.order_status;
DROP SCHEMA IF EXISTS orders;
-- +goose StatementEnd

