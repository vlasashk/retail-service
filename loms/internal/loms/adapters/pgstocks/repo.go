package pgstocks

import (
	"errors"
	"fmt"

	"route256/loms/config"
	"route256/loms/internal/loms/adapters/pgstocks/pgstocksqry"
	"route256/loms/internal/loms/models"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type StocksRepo struct {
	Pool    *pgxpool.Pool
	queries *pgstocksqry.Queries
}

func New(ctx context.Context, cfg config.StocksRepoCfg, logger zerolog.Logger) (*StocksRepo, error) {
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

	return &StocksRepo{
		Pool:    pool,
		queries: pgstocksqry.New(pool),
	}, nil
}

func (s *StocksRepo) Reserve(ctx context.Context, order models.Order) error {
	tx, err := s.Pool.Begin(ctx)
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

	qtx := s.queries.WithTx(tx)

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
	tx, err := s.Pool.Begin(ctx)
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

	qtx := s.queries.WithTx(tx)

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
	tx, err := s.Pool.Begin(ctx)
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

	qtx := s.queries.WithTx(tx)

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
	stocksData, err := s.queries.GetStockBySKU(ctx, int64(skuID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrItemNotFound
		}
		return 0, err
	}

	return stocksData.Available - stocksData.Reserved, nil
}
