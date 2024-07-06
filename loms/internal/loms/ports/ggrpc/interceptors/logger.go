package interceptors

import (
	"context"
	"time"

	"route256/loms/pkg/utils"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func LoggingInterceptor(log zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		traceID, spanID := utils.ExtractTraceInfo(ctx)
		logger := log.With().Str("trace_id", traceID).Str("span_id", spanID).Logger()

		logRequest(logger, req, info)

		defer func(start time.Time) {
			if err != nil {
				logEvent := logger.With().
					Str("method", info.FullMethod).
					Dur("duration", time.Since(start)).Logger()

				st, _ := status.FromError(err)

				switch st.Code() {
				case codes.NotFound, codes.InvalidArgument, codes.FailedPrecondition:
					logEvent.Warn().Err(err).Send()
				case codes.Internal:
					logEvent.Error().Err(err).Send()
				default:
					logEvent.Error().Err(err).Send()
				}

				return
			}

			logResponse(logger, time.Since(start), resp, info)
		}(time.Now())

		return handler(logger.WithContext(ctx), req)
	}
}

func logRequest(logger zerolog.Logger, req any, info *grpc.UnaryServerInfo) {
	reqProto, ok := req.(proto.Message)
	if !ok {
		logger.Warn().Msg("failed to cast request to proto.Message")
	}

	reqBytes, err := protojson.Marshal(reqProto)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to marshal request")
	}

	logger.Debug().
		Str("method", info.FullMethod).
		Bytes("request body", reqBytes).
		Send()
}

func logResponse(logger zerolog.Logger, timeTaken time.Duration, resp any, info *grpc.UnaryServerInfo) {
	respProto, ok := resp.(proto.Message)
	if !ok {
		logger.Warn().Msg("failed to cast response to proto.Message")
	}

	respBytes, err := protojson.Marshal(respProto)
	if err != nil {
		logger.Warn().Err(err).Msg("failed to marshal response")
	}

	logger.Debug().
		Str("method", info.FullMethod).
		Dur("duration", timeTaken).
		Bytes("response", respBytes).
		Msg("received unary request")
}
