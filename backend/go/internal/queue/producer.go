package queue

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/segmentio/kafka-go"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/pkg/config"
)

type LogProducer struct {
	writer *kafka.Writer
}

func NewLogProducer(cfg *config.Config) *LogProducer {
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    cfg.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}

	return &LogProducer{writer: writer}
}

func (p *LogProducer) Publish(ctx context.Context, log domain.RequestLog) error {
	data, err := json.Marshal(log)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: data,
	})
}

func (p *LogProducer) Close() error {
	return p.writer.Close()
}
