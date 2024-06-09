package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port      string `env:"LOMS_PORT" env-default:"50000"`
	LoggerLVL string `env:"LOMS_LOG_LVL" env-default:"debug"`
}

func New() (Config, error) {
	var res Config
	if err := cleanenv.ReadEnv(&res); err != nil {
		return Config{}, err
	}
	return res, nil
}
