package objectinfostorage

import (
	"context"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ObjectInfoStorer interface {
	CreateIndex(ctx context.Context) error
	Create(ctx context.Context, objectPartition model.ObjectInfo) error
}

type objectInfoStorage struct {
	collectionName string
	db             *mongo.Database
}

type Option func(*objectInfoStorage)

func WithObjectCollection(collection string) Option {
	return func(s *objectInfoStorage) {
		s.collectionName = collection
	}
}

func WithDB(db *mongo.Database) Option {
	return func(s *objectInfoStorage) {
		s.db = db
	}
}

func New(opts ...Option) ObjectInfoStorer {
	s := &objectInfoStorage{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
