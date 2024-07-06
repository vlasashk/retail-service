package pgstocks

import (
	"errors"
	"fmt"

	"route256/loms/config"
	"route256/loms/internal/loms/adapters/pgstocks/pgstocksqry"
	"route256/loms/internal/loms/models"
	"route256/loms/pkg/pgcluster"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"golang.org/x/net/context"
)

type StocksRepo struct {
	Cluster *pgcluster.Cluster
}

func New(ctx context.Context, cfg config.StocksRepoCfg) (*StocksRepo, error) {
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

	return &StocksRepo{
		Cluster: cluster,
	}, nil
}

func (s *StocksRepo) Reserve(ctx context.Context, order models.Order) error {
	ctx, span := otel.Tracer("").Start(ctx, "Stocks: reserve")
	defer span.End()

	pool, err := s.Cluster.GetWriter()
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

	// Каждый раз создается через New() что бы можно было подменять пулл из кластера
	qtx := pgstocksqry.New(tx)

	for _, item := range order.Items {
		var stocksData pgstocksqry.RetrieveStockForUpdateRow
		stocksData, err = qtx.RetrieveStockForUpdate(ctx, int64(item.SKUid))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.ErrItemNotFound
			}
			return err
		}

		toReserve := int64(item.Count)
		if stocksData.Available < stocksData.Reserved+toReserve {
			err = models.ErrInsufficientStock
			return err
		}

		err = qtx.ReserveStocks(ctx, pgstocksqry.ReserveStocksParams{
			Reserved: toReserve,
			ID:       int64(item.SKUid),
		})
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit fail: %w", err)
	}

	return nil
}

func (s *StocksRepo) ReserveRemove(ctx context.Context, order models.Order) error {
	ctx, span := otel.Tracer("").Start(ctx, "Stocks: reserve remove")
	defer span.End()

	pool, err := s.Cluster.GetWriter()
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

	qtx := pgstocksqry.New(tx)

	for _, item := range order.Items {
		var stocksData pgstocksqry.RetrieveStockForUpdateRow
		stocksData, err = qtx.RetrieveStockForUpdate(ctx, int64(item.SKUid))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.ErrItemNotFound
			}
			return err
		}

		toReserve := int64(item.Count)
		if stocksData.Reserved < toReserve {
			err = models.ErrReservationConflict
			return err
		}

		err = qtx.RemovePayedReservation(ctx, pgstocksqry.RemovePayedReservationParams{
			Available: toReserve,
			ID:        int64(item.SKUid),
		})
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit fail: %w", err)
	}

	return nil
}

func (s *StocksRepo) ReserveCancel(ctx context.Context, order models.Order) error {
	ctx, span := otel.Tracer("").Start(ctx, "Stocks: reserve cancel")
	defer span.End()

	pool, err := s.Cluster.GetWriter()
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

	qtx := pgstocksqry.New(tx)

	for _, item := range order.Items {
		var stocksData pgstocksqry.RetrieveStockForUpdateRow
		stocksData, err = qtx.RetrieveStockForUpdate(ctx, int64(item.SKUid))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.ErrItemNotFound
			}
			return err
		}

		toReserve := int64(item.Count)
		if stocksData.Reserved < toReserve {
			err = models.ErrReservationConflict
			return err
		}

		err = qtx.CancelReservation(ctx, pgstocksqry.CancelReservationParams{
			Reserved: toReserve,
			ID:       int64(item.SKUid),
		})
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit fail: %w", err)
	}

	return nil
}

func (s *StocksRepo) GetBySKU(ctx context.Context, skuID uint32) (int64, error) {
	ctx, span := otel.Tracer("").Start(ctx, "Stocks: get by sku")
	defer span.End()

	pool, err := s.Cluster.GetReader()
	if err != nil {
		return 0, err
	}

	stocksData, err := pgstocksqry.New(pool).GetStockBySKU(ctx, int64(skuID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrItemNotFound
		}
		return 0, err
	}

	return stocksData.Available - stocksData.Reserved, nil
}

func (s *StocksRepo) Close() error {
	s.Cluster.Close()
	return nil
}
