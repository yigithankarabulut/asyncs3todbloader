package productstorage

import (
	"context"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductStorer interface {
	CreateIndex(ctx context.Context) error
	Create(ctx context.Context, product model.Product) error
}

type productStorage struct {
	productCollectionName string
	db                    *mongo.Database
}

type Option func(*productStorage)

func WithProductCollection(collection string) Option {
	return func(s *productStorage) {
		s.productCollectionName = collection
	}
}

func WithDB(db *mongo.Database) Option {
	return func(s *productStorage) {
		s.db = db
	}
}

func New(opts ...Option) ProductStorer {
	s := &productStorage{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
