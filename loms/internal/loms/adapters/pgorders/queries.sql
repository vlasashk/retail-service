-- name: CreateOrder :one
INSERT INTO orders.orders (user_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: InsertOrderItems :copyfrom
INSERT INTO orders.order_items (sku_id, order_id, count, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetOrderById :one
SELECT id, user_id, status
FROM orders.orders
WHERE id = $1;

-- name: GetOrderItems :many
SELECT sku_id, order_id, count
FROM orders.order_items
WHERE order_id = $1;
