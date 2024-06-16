package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port       string `env:"LOMS_PORT" env-default:"50000"`
	LoggerLVL  string `env:"LOMS_LOG_LVL" env-default:"debug"`
	Gateway    HTTPGateCfg
	OrdersRepo OrdersRepoCfg
}

type HTTPGateCfg struct {
	SwaggerFilePath string `env:"SWAGGER_FILE_PATH" env-default:"./api/openapiv2/loms.swagger.json"`
	SwaggerDirPath  string `env:"SWAGGER_DIR_PATH" env-default:"./swagger-ui"`
	Address         string `env:"GATEWAY_ADDRESS" env-default:"localhost:8888"`
	LOMSAddress     string `env:"LOMS_SERVICE_ADDRESS"  env-default:"localhost:50000"`
}

type OrdersRepoCfg struct {
	Host     string `env:"ORDERS_DB_HOST" env-default:"localhost"`
	Port     string `env:"ORDERS_DB_PORT" env-default:"5432"`
	Name     string `env:"ORDERS_DB_NAME" env-default:"postgres"`
	User     string `env:"ORDERS_DB_USER" env-default:"postgres"`
	Password string `env:"ORDERS_DB_PASSWORD" env-required:"true"`
}

type StocksRepoCfg struct {
	Schema   string `env:"STOCKS_DB_SCHEMA" env-default:"stocks"`
	Host     string `env:"STOCKS_DB_HOST" env-default:"localhost"`
	Port     string `env:"STOCKS_DB_PORT" env-default:"5432"`
	Name     string `env:"STOCKS_DB_NAME" env-default:"postgres"`
	User     string `env:"STOCKS_DB_USER" env-default:"postgres"`
	Password string `env:"STOCKS_DB_PASSWORD" env-required:"true"`
}

func New() (Config, error) {
	var res Config
	if err := cleanenv.ReadEnv(&res); err != nil {
		return Config{}, err
	}
	return res, nil
}
