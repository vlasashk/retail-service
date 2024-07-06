package interceptors

import (
	"context"
	"time"

	"route256/loms/internal/loms/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()

	resp, err = handler(ctx, req)

	duration := time.Since(startTime).Seconds()
	statusCode := status.Code(err).String()

	metrics.RequestDuration.WithLabelValues(statusCode, info.FullMethod).Observe(duration)
	metrics.RequestsTotal.WithLabelValues(statusCode, info.FullMethod).Inc()

	return resp, err
}
