package productstorage

import (
	"context"
	"fmt"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/customerror"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"github.com/yigithankarabulut/asyncs3todbloader/job/pkg/constant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndex method creates an index for the ID field of the product collection.
func (s *productStorage) CreateIndex(ctx context.Context) error {
	if _, err := s.db.Collection(s.productCollectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"id": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return customerror.New(constant.ErrCreateIndexFailed, true).
			Wrap(fmt.Errorf("productstorage: failed to create index: %w", err)).AddData("err: " + err.Error())
	}
	return nil
}

// Create method creates a product in the database.
func (s *productStorage) Create(ctx context.Context, product model.Product) error {
	if _, err := s.db.Collection(s.productCollectionName).InsertOne(ctx, product); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return customerror.New(constant.ErrIDExists, false).Wrap(fmt.Errorf("productstorage: failed to create product: %w", err)).AddData(product.ID)
		}
		return customerror.New(constant.ErrCreateProduct, true).
			Wrap(fmt.Errorf("productstorage: failed to create product: %w", err)).AddData("err: " + err.Error())
	}
	return nil
}
