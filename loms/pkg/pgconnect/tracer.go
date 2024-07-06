package pgconnect

import (
	"context"
	"strings"
	"time"

	"route256/loms/internal/loms/metrics"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type contextKey int

type tracer struct{}

const (
	startTimeKey contextKey = iota
)

func (t *tracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	logger := *zerolog.Ctx(ctx)
	logger.Info().Str("sql", data.SQL).Any("args", data.Args).Msg("Executing command")

	ctx, _ = otel.Tracer("").Start(ctx, data.SQL)

	ctx = context.WithValue(ctx, startTimeKey, time.Now())

	return ctx
}

func (t *tracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	span.End()

	startTime, ok := ctx.Value(startTimeKey).(time.Time)
	if !ok {
		return
	}
	duration := time.Since(startTime).Seconds()

	status := "success"

	if data.Err != nil {
		status = "error"
	}

	metrics.DBQueryDuration.WithLabelValues(status, extractOperation(data.CommandTag)).Observe(duration)
	metrics.DBQueriesTotal.WithLabelValues(status, extractOperation(data.CommandTag)).Inc()
}

func (t *tracer) TraceCopyFromStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromStartData) context.Context {
	logger := *zerolog.Ctx(ctx)
	logger.Info().Str("sql", "CopyFrom").Any("args", data.ColumnNames).Msg("Executing command")

	ctx, _ = otel.Tracer("").Start(ctx, "CopyFrom",
		trace.WithAttributes(
			attribute.StringSlice("target table", data.TableName),
			attribute.StringSlice("columns", data.ColumnNames),
		),
	)

	ctx = context.WithValue(ctx, startTimeKey, time.Now())

	return ctx
}

func (t *tracer) TraceCopyFromEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceCopyFromEndData) {
	span := trace.SpanFromContext(ctx)
	span.End()

	startTime, ok := ctx.Value(startTimeKey).(time.Time)
	if !ok {
		return
	}
	duration := time.Since(startTime).Seconds()

	status := "success"

	if data.Err != nil {
		status = "error"
	}
	metrics.DBQueryDuration.WithLabelValues(status, extractOperation(data.CommandTag)).Observe(duration)
	metrics.DBQueriesTotal.WithLabelValues(status, extractOperation(data.CommandTag)).Inc()
}

func extractOperation(tag pgconn.CommandTag) string {
	parts := strings.SplitN(tag.String(), " ", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}
