package pgorders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"route256/loms/config"
	"route256/loms/internal/loms/adapters/pgorders/pgordersqry"
	"route256/loms/internal/loms/models"
	"route256/loms/pkg/pgcluster"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type OrdersRepo struct {
	Cluster *pgcluster.Cluster
}

func New(ctx context.Context, cfg config.OrdersRepoCfg, logger zerolog.Logger) (*OrdersRepo, error) {
	masterURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.HostMaster,
		cfg.PortMaster,
		cfg.Name)

	slaveURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.HostSlave,
		cfg.PortSlave,
		cfg.Name)

	masterPool, err := pgconnect.Connect(ctx, masterURL, logger)
	if err != nil {
		return nil, err
	}

	slavePool, err := pgconnect.Connect(ctx, slaveURL, logger)
	if err != nil {
		return nil, err
	}

	cluster := pgcluster.New().SetWriter(masterPool).AddReader(masterPool, slavePool)

	return &OrdersRepo{
		Cluster: cluster,
	}, nil
}

func (or *OrdersRepo) Create(ctx context.Context, order models.Order) (int64, error) {
	orderData := orderToDTO(order)
	pool, err := or.Cluster.GetWriter()
	if err != nil {
		return 0, err
	}

	tx, err := pool.Begin(ctx)
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

	// Каждый раз создается через New() что бы можно было подменять пулл из кластера
	qtx := pgordersqry.New(tx)

	orderID, err := qtx.CreateOrder(ctx, pgordersqry.CreateOrderParams{
		UserID:    orderData.UserID,
		Status:    orderData.Status,
		CreatedAt: orderData.CreatedAt,
		UpdatedAt: orderData.UpdatedAt,
	})
	if err != nil {
		return 0, fmt.Errorf("query fail: %w", err)
	}

	itemsData := itemsToDTO(orderID, order.Items)

	_, err = qtx.InsertOrderItems(ctx, itemsData)
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
	pool, err := or.Cluster.GetWriter()
	if err != nil {
		return err
	}

	// sqlc не умеет возвращать commandTag (нужен для обработки ошибки)
	tag, err := pool.Exec(ctx, qry, statusToDTO(status), time.Now(), orderID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrOrderNotFound
	}

	return nil
}

func (or *OrdersRepo) GetByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	pool, err := or.Cluster.GetReader()
	if err != nil {
		return models.Order{}, err
	}

	qPool := pgordersqry.New(pool)

	oderData, err := qPool.GetOrderById(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, models.ErrOrderNotFound
		}
		return models.Order{}, err
	}

	itemsData, err := qPool.GetOrderItems(ctx, orderID)
	if err != nil {
		return models.Order{}, err
	}

	return orderToDomain(oderData, itemsData), nil
}

func (or *OrdersRepo) Close() error {
	or.Cluster.Close()
	return nil
}
