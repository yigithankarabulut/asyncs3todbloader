package objectinfostorage

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

// CreateIndex method creates an index for the ETag field of the object info collection.
func (s *objectInfoStorage) CreateIndex(ctx context.Context) error {
	if _, err := s.db.Collection(s.collectionName).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"etag": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return customerror.New(constant.ErrCreateIndexFailed, true).Wrap(fmt.Errorf("objectinfostorage: failed to create index: %w", err))
	}
	return nil
}

// Create method creates an object info in the database.
func (s *objectInfoStorage) Create(ctx context.Context, object model.ObjectInfo) error {
	if _, err := s.db.Collection(s.collectionName).InsertOne(ctx, object); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return customerror.New(constant.ErrETagExists, true).Wrap(fmt.Errorf("objectinfostorage: failed to create object info: %w", err)).AddData(object.ETag)
		}
		return customerror.New(constant.ErrCreateObjectInfo, true).
			Wrap(fmt.Errorf("objectinfostorage: failed to create object info: %w", err)).AddData("err: " + err.Error())
	}
	return nil
}
