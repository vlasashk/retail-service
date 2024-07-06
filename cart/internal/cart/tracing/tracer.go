package tracing

import (
	"context"
	"time"

	"route256/cart/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

type Tracer struct {
	*sdktrace.TracerProvider
}

func New(ctx context.Context, cfg config.TelemetryCfg) (Tracer, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Address),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return Tracer{}, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return Tracer{tp}, nil
}

func (t Tracer) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	return t.Shutdown(ctx)
}
