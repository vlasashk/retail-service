package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address     string `env:"NOTIFIER_ADDR" env-default:"localhost:8099"`
	LoggerLVL   string `env:"NOTIFIER_LOG_LVL" env-default:"debug"`
	ServiceName string `env:"SERVICE_NAME" env-default:"notifier"`
	Consumer    KafkaCfg
}

type KafkaCfg struct {
	Brokers []string `env:"KAFKA_BROKERS" env-separator:"," env-default:"localhost:9092"`
	Topic   []string `env:"KAFKA_TOPICS" env-separator:"," env-default:"loms.order-events"`
	GroupID string   `env:"KAFKA_GROUP_ID" env-default:"order_events"`
}

func New() (Config, error) {
	var res Config
	if err := cleanenv.ReadEnv(&res); err != nil {
		return Config{}, err
	}
	return res, nil
}
