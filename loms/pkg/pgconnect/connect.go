package pgconnect

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var timeout = 10 * time.Second

func Connect(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pgxCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	pgxCfg.ConnConfig.Tracer = &tracer{}

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping connection pool: %w", err)
	}

	return dbPool, nil
}
