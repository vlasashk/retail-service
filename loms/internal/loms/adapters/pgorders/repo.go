package pgorders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"route256/loms/config"
	"route256/loms/internal/loms/models"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type OrdersRepo struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, cfg config.OrdersRepoCfg, logger zerolog.Logger) (*OrdersRepo, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name)

	pool, err := pgconnect.Connect(ctx, url, logger)
	if err != nil {
		return nil, err
	}

	return &OrdersRepo{Pool: pool}, nil
}

func (or *OrdersRepo) Create(ctx context.Context, order models.Order) (int64, error) {
	qry := `INSERT INTO orders.orders (user_id, status, created_at, updated_at)
				VALUES ($1, $2, $3, $4)
				RETURNING id;`
	var orderID int64

	orderData := orderToDTO(order)

	tx, err := or.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("connection acquire fail: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				log.Error().Err(err).Caller().Send()
			}
		}
	}()

	if err = tx.QueryRow(ctx, qry, orderData.UserID, orderData.Status, orderData.CreatedAt, orderData.UpdatedAt).Scan(&orderID); err != nil {
		return 0, fmt.Errorf("query fail: %w", err)
	}

	itemsData := itemsToDTO(orderID, order.Items)

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"orders", "order_items"}, []string{"sku_id", "order_id", "count", "created_at", "updated_at"},
		pgx.CopyFromSlice(len(itemsData), func(i int) ([]any, error) {
			return []any{itemsData[i].SKU, itemsData[i].OrderID, itemsData[i].Count, itemsData[i].CreatedAt, itemsData[i].UpdatedAt}, nil
		}))
	if err != nil {
		return 0, fmt.Errorf("query fail: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("transaction commit fail: %w", err)
	}

	return orderID, nil
}

func (or *OrdersRepo) SetStatus(ctx context.Context, orderID int64, status models.OrderStatus) error {
	qry := `UPDATE orders.orders
				SET status=$1, updated_at=$2
				WHERE id=$3`

	tag, err := or.Pool.Exec(ctx, qry, statusToDTO(status), time.Now(), orderID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return models.ErrOrderNotFound
	}

	return nil
}

func (or *OrdersRepo) GetByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	orderQry := `SELECT id, user_id, status
					FROM orders.orders
					WHERE id=$1`
	itemsQry := `SELECT sku_id, order_id, count
					FROM orders.order_items
					WHERE order_id=$1`

	orderRows, err := or.Pool.Query(ctx, orderQry, orderID)
	if err != nil {
		return models.Order{}, err
	}
	defer orderRows.Close()

	oderData, err := pgx.CollectOneRow(orderRows, pgx.RowToStructByName[orderDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, models.ErrOrderNotFound
		}
		return models.Order{}, err
	}

	itemRows, err := or.Pool.Query(ctx, itemsQry, orderID)
	if err != nil {
		return models.Order{}, err
	}
	defer itemRows.Close()

	itemsData, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[itemDTO])
	if err != nil {
		return models.Order{}, err
	}

	return orderToDomain(oderData, itemsData), nil
}

func (or *OrdersRepo) Close() error {
	or.Pool.Close()
	return nil
}
