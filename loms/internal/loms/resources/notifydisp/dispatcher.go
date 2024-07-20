package notifydisp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"route256/loms/config"
	"route256/loms/internal/loms/adapters/notifybox"
	"route256/loms/internal/loms/models"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Dispatcher struct {
	producer  sarama.SyncProducer
	writer    *pgxpool.Pool
	tick      time.Duration
	logger    zerolog.Logger
	notifier  *notifybox.Notifier
	batchSize int
	topic     string
}

func New(cfg config.DispatcherCfg, writer *pgxpool.Pool, logger zerolog.Logger) (*Dispatcher, error) {
	producer, err := sarama.NewSyncProducer(cfg.Brokers, prepareConfig())
	if err != nil {
		return nil, err
	}

	return &Dispatcher{
		producer:  producer,
		writer:    writer,
		tick:      cfg.Tick,
		logger:    logger,
		notifier:  notifybox.New(writer),
		batchSize: cfg.BatchSize,
		topic:     cfg.Topic,
	}, nil
}

func (d *Dispatcher) Close() error {
	return d.producer.Close()
}

func (d *Dispatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(d.tick)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := d.executeDispatch(ctx); err != nil {
				d.logger.Error().Err(err).Send()
			}
		}
	}
}

func prepareConfig() *sarama.Config {
	c := sarama.NewConfig()

	// по ключу
	c.Producer.Partitioner = sarama.NewHashPartitioner
	// acks = -1 (all) - ждем успешной записи на лидере партиции и всех in-sync реплик (настроено в кластере кафки)
	c.Producer.RequiredAcks = sarama.WaitForAll
	// Уменьшаем пропускную способность, тем самым гарантируем строгий порядок отправки сообщений/батчей
	c.Net.MaxOpenRequests = 1
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true

	return c
}

func (d *Dispatcher) executeDispatch(ctx context.Context) error {
	tx, err := d.writer.Begin(ctx)
	if err != nil {
		return fmt.Errorf("connection acquire fail: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				d.logger.Error().Err(err).Caller().Send()
			}
		}
	}()

	notifier := d.notifier.WithTx(tx)

	events, err := notifier.FetchNextBatch(ctx, d.batchSize)
	if err != nil {
		return fmt.Errorf("fetch next batch fail: %w", err)
	}

	if err = d.sendEventsInBatch(events); err != nil {
		return fmt.Errorf("send events fail: %w", err)
	}

	if err = notifier.MarkAsSent(ctx, events); err != nil {
		return fmt.Errorf("mark events as sent fail: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit fail: %w", err)
	}

	return nil
}

func (d *Dispatcher) sendEventsInBatch(events []*models.Event) error {
	messages := make([]*sarama.ProducerMessage, len(events))

	for i, event := range events {
		eventToSend := eventToSend(*event)

		value, err := json.Marshal(eventToSend)
		if err != nil {
			return fmt.Errorf("marshal event fail: %w", err)
		}

		messages[i] = &sarama.ProducerMessage{
			Topic: d.topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("%d", eventToSend.ID)),
			Value: sarama.ByteEncoder(value),
		}
	}

	return d.producer.SendMessages(messages)
}
