// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package pgordersqry

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders.orders (user_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4)
RETURNING id
`

type CreateOrderParams struct {
	UserID    int64
	Status    OrdersOrderStatus
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (int64, error) {
	row := q.db.QueryRow(ctx, createOrder,
		arg.UserID,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getOrderById = `-- name: GetOrderById :one
SELECT id, user_id, status
FROM orders.orders
WHERE id = $1
`

type GetOrderByIdRow struct {
	ID     int64
	UserID int64
	Status OrdersOrderStatus
}

func (q *Queries) GetOrderById(ctx context.Context, id int64) (GetOrderByIdRow, error) {
	row := q.db.QueryRow(ctx, getOrderById, id)
	var i GetOrderByIdRow
	err := row.Scan(&i.ID, &i.UserID, &i.Status)
	return i, err
}

const getOrderItems = `-- name: GetOrderItems :many
SELECT sku_id, order_id, count
FROM orders.order_items
WHERE order_id = $1
`

type GetOrderItemsRow struct {
	SkuID   int64
	OrderID int64
	Count   int64
}

func (q *Queries) GetOrderItems(ctx context.Context, orderID int64) ([]GetOrderItemsRow, error) {
	rows, err := q.db.Query(ctx, getOrderItems, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetOrderItemsRow
	for rows.Next() {
		var i GetOrderItemsRow
		if err := rows.Scan(&i.SkuID, &i.OrderID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type InsertOrderItemsParams struct {
	SkuID     int64
	OrderID   int64
	Count     int64
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}