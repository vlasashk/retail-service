package config

import "github.com/ilyakaznacheev/cleanenv"

type CartConfig struct {
	Address         string `env:"CART_ADDR" env-default:"localhost:8082"`
	LoggerLVL       string `env:"CART_LOG_LVL" env-default:"info"`
	ProductProvider ProductProviderCfg
}

type ProductProviderCfg struct {
	Address     string `env:"PRODUCT_SERVICE_ADDR" env-default:"http://route256.pavl.uk:8080"`
	AccessToken string `env:"PRODUCT_SERVICE_TOKEN" env-default:"testtoken"`
	Retries     int    `env:"PRODUCT_SERVICE_RETRIES" env-default:"3"`
}

func NewCartCfg() (CartConfig, error) {
	var res CartConfig
	if err := cleanenv.ReadEnv(&res); err != nil {
		return CartConfig{}, err
	}
	return res, nil
}
