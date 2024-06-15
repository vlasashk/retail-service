package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port      string `env:"LOMS_PORT" env-default:"50000"`
	LoggerLVL string `env:"LOMS_LOG_LVL" env-default:"debug"`
	Gateway   HTTPGateCfg
}

type HTTPGateCfg struct {
	SwaggerFilePath string `env:"SWAGGER_FILE_PATH" env-default:"./api/openapiv2/loms.swagger.json"`
	SwaggerDirPath  string `env:"SWAGGER_DIR_PATH" env-default:"./swagger-ui"`
	Address         string `env:"GATEWAY_ADDRESS" env-default:"localhost:8888"`
	LOMSAddress     string `env:"LOMS_SERVICE_ADDRESS"  env-default:"localhost:50000"`
}

func New() (Config, error) {
	var res Config
	if err := cleanenv.ReadEnv(&res); err != nil {
		return Config{}, err
	}
	return res, nil
}
