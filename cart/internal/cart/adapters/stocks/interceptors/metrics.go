package interceptors

import (
	"context"
	"time"

	"route256/cart/internal/cart/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const serviceName = "loms"

func Metrics(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	startTime := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	duration := time.Since(startTime).Seconds()
	statusCode := status.Code(err).String()

	metrics.ExternalRequestDuration.WithLabelValues(serviceName, statusCode, method).Observe(duration)
	metrics.ExternalRequestsTotal.WithLabelValues(serviceName, statusCode, method).Inc()

	return err
}
