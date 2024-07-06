package utils

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func ExtractTraceInfo(ctx context.Context) (string, string) {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return "", ""
	}
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	return traceID, spanID
}
