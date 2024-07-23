package consumer

import (
	"context"
	"errors"

	"route256/notifier/config"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
)

type Consumer struct {
	group   sarama.ConsumerGroup
	logger  zerolog.Logger
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func New(cfg config.KafkaCfg, logger zerolog.Logger) (*Consumer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Version = sarama.MaxVersion
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		group:   group,
		logger:  logger,
		handler: NewHandler(logger),
		topics:  cfg.Topic,
	}, nil
}

func (c *Consumer) Close() error {
	return c.group.Close()
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		if err := c.group.Consume(ctx, c.topics, c.handler); err != nil {
			c.logger.Error().Err(err).Msg("Error from consumer")
		}

		if ctx.Err() != nil {
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			return ctx.Err()
		}
	}
}
