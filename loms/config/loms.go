package config

import (
	"errors"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port        string `env:"LOMS_PORT" env-default:"50000"`
	LoggerLVL   string `env:"LOMS_LOG_LVL" env-default:"debug"`
	ServiceName string `env:"SERVICE_NAME" env-default:"loms"`
	Gateway     HTTPGateCfg
	OrdersRepo  OrdersRepoCfg
	StocksRepo  StocksRepoCfg
	Telemetry   TelemetryCfg
	Dispatcher  DispatcherCfg
}

type TelemetryCfg struct {
	ServiceName string `env:"SERVICE_NAME" env-default:"loms"`
	Address     string `env:"TELEMETRY_ADDR" env-default:"localhost:4317"`
}

type HTTPGateCfg struct {
	SwaggerFilePath string `env:"SWAGGER_FILE_PATH" env-default:"./api/openapiv2/loms.swagger.json"`
	SwaggerDirPath  string `env:"SWAGGER_DIR_PATH" env-default:"./swagger-ui"`
	Address         string `env:"GATEWAY_ADDRESS" env-default:"localhost:8888"`
	LOMSAddress     string `env:"LOMS_SERVICE_ADDRESS"  env-default:"localhost:50000"`
}

type OrdersRepoCfg struct {
	HostMaster string `env:"ORDERS_DB_MASTER_HOST" env-default:"localhost"`
	PortMaster string `env:"ORDERS_DB_MASTER_PORT" env-default:"5432"`
	HostSlave  string `env:"ORDERS_DB_SLAVE_HOST" env-default:"localhost"`
	PortSlave  string `env:"ORDERS_DB_SLAVE_PORT" env-default:"5433"`
	Name       string `env:"ORDERS_DB_NAME" env-default:"postgres"`
	User       string `env:"ORDERS_DB_USER" env-default:"postgres"`
	Password   string `env:"ORDERS_DB_PASSWORD" env-required:"true"`
}

type StocksRepoCfg struct {
	HostMaster string `env:"STOCKS_DB_MASTER_HOST" env-default:"localhost"`
	PortMaster string `env:"STOCKS_DB_MASTER_PORT" env-default:"5432"`
	HostSlave  string `env:"STOCKS_DB_SLAVE_HOST" env-default:"localhost"`
	PortSlave  string `env:"STOCKS_DB_SLAVE_PORT" env-default:"5433"`
	Name       string `env:"STOCKS_DB_NAME" env-default:"postgres"`
	User       string `env:"STOCKS_DB_USER" env-default:"postgres"`
	Password   string `env:"STOCKS_DB_PASSWORD" env-required:"true"`
}

type DispatcherCfg struct {
	Brokers   []string      `env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	Topic     string        `env:"KAFKA_TOPIC" env-default:"loms.order-events"`
	Tick      time.Duration `env:"DISPATCHER_TICK" env-default:"5s"`
	BatchSize int           `env:"DISPATCHER_BATCH_SIZE" env-default:"3"`
}

// New - parses environment variables
func New(opts ...Option) (Config, error) {
	o := defaultOptions()

	for _, opt := range opts {
		opt.apply(&o)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(o.configPath, &cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
		log.Warn().Err(err).Msg("yaml config file not found")

		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
