-- name: RetrieveStockForUpdate :one
SELECT available, reserved
FROM stocks.stocks
WHERE id = $1 FOR UPDATE;

-- name: ReserveStocks :exec
UPDATE stocks.stocks
SET reserved = reserved + $1
WHERE id = $2;

-- name: RemovePayedReservation :exec
UPDATE stocks.stocks
SET available = available - $1,
    reserved = reserved - $1
WHERE id = $2;

-- name: CancelReservation :exec
UPDATE stocks.stocks
SET reserved = reserved - $1
WHERE id = $2;

-- name: GetStockBySKU :one
SELECT available, reserved
FROM stocks.stocks
WHERE id = $1;
