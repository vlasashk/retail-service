-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS stocks;

CREATE TABLE IF NOT EXISTS stocks.stocks (
  id BIGINT PRIMARY KEY,
  available INTEGER NOT NULL,
  reserved INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stocks.stocks;
DROP SCHEMA IF EXISTS stocks;
-- +goose StatementEnd
