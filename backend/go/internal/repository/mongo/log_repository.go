package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gino/cars-crud/internal/domain"
	"github.com/gino/cars-crud/internal/repository"
)

type logRepository struct {
	collection *mongo.Collection
}

func NewLogRepository(collection *mongo.Collection) repository.LogRepository {
	return &logRepository{collection: collection}
}

func (r *logRepository) GetAll(ctx context.Context, offset, limit int) ([]domain.RequestLog, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "timestamp", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var logs []domain.RequestLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
