package queue

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/pkg/config"
)

type LogConsumer struct {
	reader     *kafka.Reader
	collection *mongo.Collection
}

func NewLogConsumer(cfg *config.Config, mongoClient *mongo.Client) *LogConsumer {
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    cfg.KafkaTopic,
		GroupID:  "log-consumer-group",
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	collection := mongoClient.Database(cfg.MongoDB).Collection(cfg.MongoCollection)

	return &LogConsumer{
		reader:     reader,
		collection: collection,
	}
}

func (c *LogConsumer) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					log.Printf("kafka consumer read error: %v", err)
					continue
				}

				var reqLog domain.RequestLog
				if err := json.Unmarshal(msg.Value, &reqLog); err != nil {
					log.Printf("kafka consumer unmarshal error: %v", err)
					continue
				}

				if _, err := c.collection.InsertOne(ctx, reqLog); err != nil {
					log.Printf("mongo insert error: %v", err)
				}
			}
		}
	}()
}

func (c *LogConsumer) Close() error {
	return c.reader.Close()
}
