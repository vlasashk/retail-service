-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS outbox;

CREATE TYPE outbox.order_status AS ENUM (
  'Unknown',
  'New',
  'AwaitingPayment',
  'Payed',
  'Cancelled',
  'Failed'
);

CREATE TABLE IF NOT EXISTS outbox.notifier (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT,
    status outbox.order_status NOT NULL,
    is_sent BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox.notifier;
DROP TYPE IF EXISTS outbox.order_status;
DROP SCHEMA IF EXISTS outbox;
-- +goose StatementEnd



