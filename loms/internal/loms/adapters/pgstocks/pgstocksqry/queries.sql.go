// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: queries.sql

package pgstocksqry

import (
	"context"
)

const cancelReservation = `-- name: CancelReservation :exec
UPDATE stocks.stocks
SET reserved = reserved - $1
WHERE id = $2
`

type CancelReservationParams struct {
	Reserved int64
	ID       int64
}

func (q *Queries) CancelReservation(ctx context.Context, arg CancelReservationParams) error {
	_, err := q.db.Exec(ctx, cancelReservation, arg.Reserved, arg.ID)
	return err
}

const getStockBySKU = `-- name: GetStockBySKU :one
SELECT available, reserved
FROM stocks.stocks
WHERE id = $1
`

type GetStockBySKURow struct {
	Available int64
	Reserved  int64
}

func (q *Queries) GetStockBySKU(ctx context.Context, id int64) (GetStockBySKURow, error) {
	row := q.db.QueryRow(ctx, getStockBySKU, id)
	var i GetStockBySKURow
	err := row.Scan(&i.Available, &i.Reserved)
	return i, err
}

const removePayedReservation = `-- name: RemovePayedReservation :exec
UPDATE stocks.stocks
SET available = available - $1,
    reserved = reserved - $1
WHERE id = $2
`

type RemovePayedReservationParams struct {
	Available int64
	ID        int64
}

func (q *Queries) RemovePayedReservation(ctx context.Context, arg RemovePayedReservationParams) error {
	_, err := q.db.Exec(ctx, removePayedReservation, arg.Available, arg.ID)
	return err
}

const reserveStocks = `-- name: ReserveStocks :exec
UPDATE stocks.stocks
SET reserved = reserved + $1
WHERE id = $2
`

type ReserveStocksParams struct {
	Reserved int64
	ID       int64
}

func (q *Queries) ReserveStocks(ctx context.Context, arg ReserveStocksParams) error {
	_, err := q.db.Exec(ctx, reserveStocks, arg.Reserved, arg.ID)
	return err
}

const retrieveStockForUpdate = `-- name: RetrieveStockForUpdate :one
SELECT available, reserved
FROM stocks.stocks
WHERE id = $1 FOR UPDATE
`

type RetrieveStockForUpdateRow struct {
	Available int64
	Reserved  int64
}

func (q *Queries) RetrieveStockForUpdate(ctx context.Context, id int64) (RetrieveStockForUpdateRow, error) {
	row := q.db.QueryRow(ctx, retrieveStockForUpdate, id)
	var i RetrieveStockForUpdateRow
	err := row.Scan(&i.Available, &i.Reserved)
	return i, err
}