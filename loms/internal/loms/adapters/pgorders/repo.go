package pgorders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"route256/loms/config"
	"route256/loms/internal/loms/adapters/notifybox"
	"route256/loms/internal/loms/adapters/pgorders/pgordersqry"
	"route256/loms/internal/loms/models"
	"route256/loms/pkg/pgcluster"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
)

type OrdersRepo struct {
	Cluster  *pgcluster.Cluster
	notifier *notifybox.Notifier
}

func New(ctx context.Context, cfg config.OrdersRepoCfg) (*OrdersRepo, error) {
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

	masterPool, err := pgconnect.Connect(ctx, masterURL)
	if err != nil {
		return nil, err
	}

	slavePool, err := pgconnect.Connect(ctx, slaveURL)
	if err != nil {
		return nil, err
	}

	cluster := pgcluster.New().SetWriter(masterPool).AddReader(masterPool, slavePool)

	return &OrdersRepo{
		Cluster:  cluster,
		notifier: notifybox.New(masterPool),
	}, nil
}

func (or *OrdersRepo) Create(ctx context.Context, order models.Order) (int64, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Orders: create")
	defer span.End()

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

	if err = or.notifier.WithTx(tx).CreateEvent(ctx, models.Event{
		OrderID: orderID,
		Status:  order.Status.String(),
	}); err != nil {
		return 0, fmt.Errorf("event send fail: %w", err)
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
	ctx, span := otel.Tracer("").Start(ctx, "Orders: set status")
	defer span.End()

	qry := `UPDATE orders.orders
				SET status=$1, updated_at=$2
				WHERE id=$3`
	pool, err := or.Cluster.GetWriter()
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				log.Error().Err(err).Caller().Send()
			}
		}
	}()

	// sqlc не умеет возвращать commandTag (нужен для обработки ошибки)
	tag, err := tx.Exec(ctx, qry, statusToDTO(status), time.Now(), orderID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return models.ErrOrderNotFound
	}

	if err = or.notifier.WithTx(tx).CreateEvent(ctx, models.Event{
		OrderID: orderID,
		Status:  status.String(),
	}); err != nil {
		return fmt.Errorf("event send fail: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit fail: %w", err)
	}

	return nil
}

func (or *OrdersRepo) GetByOrderID(ctx context.Context, orderID int64) (models.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Orders: get by orderID")
	defer span.End()

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
