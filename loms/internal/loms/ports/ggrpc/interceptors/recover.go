package interceptors

import (
	"context"
	"runtime/debug"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoverPanic(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	logger := *zerolog.Ctx(ctx)

	defer func() {
		if rec := recover(); rec != nil {
			logger.Error().Any("panic", rec).Bytes("stack", debug.Stack()).Send()
			err = status.Errorf(codes.Internal, "panic: %v", rec)
		}
	}()

	return handler(ctx, req)
}
