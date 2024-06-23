package pgconnect

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type tracer struct {
	log zerolog.Logger
}

var timeout = 10 * time.Second

func Connect(ctx context.Context, connString string, logger zerolog.Logger) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pgxCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	pgxCfg.ConnConfig.Tracer = &tracer{log: logger}

	dbPool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping connection pool: %w", err)
	}

	return dbPool, nil
}

func (t *tracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.log.Info().Str("sql", data.SQL).Any("args", data.Args).Msg("Executing command")
	return ctx
}

func (t *tracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, _ pgx.TraceQueryEndData) {
}
