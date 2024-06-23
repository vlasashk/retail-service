//go:build integration

package suits

import (
	"context"
	"time"

	lomsservicev1 "route256/loms/pkg/api/loms/v1"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultUserID     = 111111
	defaultAlphaSKUID = 1076963 // В stock-data.json 100 всего и 10 зарезервировано
	defaultBetaSKUID  = 1148162 // В stock-data.json 150 всего и 20 зарезервировано
	additionalSKUID   = 1625903 // В stock-data.json 200 всего и 30 зарезервировано
	paySKUID          = 2618151 // В stock-data.json 50 всего и 5 зарезервировано
)

var defaultReq = &lomsservicev1.OrderCreateRequest{
	User: defaultUserID,
	Items: []*lomsservicev1.Item{
		{
			Sku:   defaultAlphaSKUID,
			Count: 10,
		},
		{
			Sku:   defaultBetaSKUID,
			Count: 30,
		},
	},
}

var additionalReq = &lomsservicev1.OrderCreateRequest{
	User: defaultUserID,
	Items: []*lomsservicev1.Item{
		{
			Sku:   defaultAlphaSKUID,
			Count: 30,
		},
		{
			Sku:   additionalSKUID,
			Count: 70,
		},
	},
}

func healthCheck(conn *grpc.ClientConn, attempts int) error {
	var err error
	healthClient := grpc_health_v1.NewHealthClient(conn)

	for attempts > 0 {
		req := &grpc_health_v1.HealthCheckRequest{Service: ""}

		resp, err := healthClient.Check(context.Background(), req)
		if err == nil && resp.Status == grpc_health_v1.HealthCheckResponse_SERVING {
			log.Info().Msg("Service is ready to accept requests")
			return nil
		}

		log.Debug().Err(err).Int("attempts left", attempts).Msg("Service is not available for integration tests")
		time.Sleep(time.Second)
		attempts--
	}

	return err
}
