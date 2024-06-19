package pgstocks

import (
	"errors"
	"fmt"

	"route256/loms/config"
	"route256/loms/internal/loms/models"
	"route256/loms/pkg/pgconnect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type StocksRepo struct {
	Pool *pgxpool.Pool
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

	return &StocksRepo{Pool: pool}, nil
}

func (s *StocksRepo) Reserve(ctx context.Context, order models.Order) error {
	//nolint:goconst //Запросы держать рядом с их вызовом
	retrieveQry := `SELECT available, reserved
						FROM stocks.stocks
						WHERE id = $1 FOR UPDATE`
	updateQry := `UPDATE stocks.stocks
					SET reserved = reserved + $1
					WHERE id = $2`

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Error().Err(err).Caller().Send()
		}
	}()

	for _, item := range order.Items {
		var available, reserved int64
		if err = tx.QueryRow(ctx, retrieveQry, item.SKUid).Scan(&available, &reserved); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.ErrItemNotFound
			}
			return err
		}

		toReserve := int64(item.Count)
		if available < reserved+toReserve {
			err = models.ErrInsufficientStock
			return err
		}

		_, err = tx.Exec(ctx, updateQry, toReserve, item.SKUid)
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
	retrieveQry := `SELECT available, reserved
						FROM stocks.stocks
						WHERE id = $1 FOR UPDATE`
	updateQry := `UPDATE stocks.stocks
					SET available = available - $1,
					    reserved = reserved - $1
					WHERE id = $2`

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Error().Err(err).Caller().Send()
		}
	}()

	for _, item := range order.Items {
		var available, reserved int64
		if err = tx.QueryRow(ctx, retrieveQry, item.SKUid).Scan(&available, &reserved); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.ErrItemNotFound
			}
			return err
		}

		toReserve := int64(item.Count)
		if reserved < toReserve {
			err = models.ErrReservationConflict
			return err
		}

		_, err = tx.Exec(ctx, updateQry, toReserve, item.SKUid)
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
	retrieveQry := `SELECT available, reserved
						FROM stocks.stocks
						WHERE id = $1 FOR UPDATE`
	updateQry := `UPDATE stocks.stocks
					SET reserved = reserved - $1
					WHERE id = $2`

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Error().Err(err).Caller().Send()
		}
	}()

	for _, item := range order.Items {
		var available, reserved int64
		if err = tx.QueryRow(ctx, retrieveQry, item.SKUid).Scan(&available, &reserved); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return models.ErrItemNotFound
			}
			return err
		}

		toReserve := int64(item.Count)
		if reserved < toReserve {
			err = models.ErrReservationConflict
			return err
		}

		_, err = tx.Exec(ctx, updateQry, toReserve, item.SKUid)
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
	retrieveQry := `SELECT available, reserved
						FROM stocks.stocks
						WHERE id = $1`

	var available, reserved int64
	if err := s.Pool.QueryRow(ctx, retrieveQry, skuID).Scan(&available, &reserved); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, models.ErrItemNotFound
		}
		return 0, err
	}

	return available - reserved, nil
}
