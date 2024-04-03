package productrepository

import (
	"context"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository interface {
	GetProductByID(ctx context.Context, id int) (model.Product, error)
}

type productRepository struct {
	productCollectionName string
	db                    *mongo.Database
}

type Option func(*productRepository)

func WithProductCollection(collection string) Option {
	return func(r *productRepository) {
		r.productCollectionName = collection
	}
}

func WithDB(db *mongo.Database) Option {
	return func(r *productRepository) {
		r.db = db
	}
}

func New(opts ...Option) ProductRepository {
	r := &productRepository{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
