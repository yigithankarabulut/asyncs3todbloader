package productstorage_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/customerror"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/productstorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func Test_productStorage_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Case Success CreateIndex", func(mt *mtest.T) {
		mockCollection := productstorage.New(
			productstorage.WithDB(mt.DB),
			productstorage.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := mockCollection.CreateIndex(context.TODO())
		assert.Nil(t, err)
		mt.ClearEvents()

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err = mockCollection.CreateIndex(context.TODO())
		assert.Nil(t, err)
	})

	mt.Run("Case CreateIndex Error", func(mt *mtest.T) {
		mockCollection := productstorage.New(
			productstorage.WithDB(mt.DB),
			productstorage.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    2,
			Message: "duplicate index value error",
		}))
		err := mockCollection.CreateIndex(context.TODO())
		assert.NotNil(t, err)
		var ce *customerror.Error
		if !assert.ErrorAs(t, err, &ce) {
			t.Fatalf("error should be of type ErrCreateIndexFailed")
		}
		if errors.Is(ce, customerror.ErrCreateIndexFailed) {
			t.Fatalf("error should be of type ErrCreateIndexFailed")
		}
	})

}

func Test_productStorage_CreateIndex(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Case ID is exists error", func(mt *mtest.T) {
		mockCollection := productstorage.New(
			productstorage.WithDB(mt.DB),
			productstorage.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    11000,
			Message: "duplicate key error",
		}))
		err := mockCollection.Create(context.TODO(), model.Product{
			ID:          1234321,
			Brand:       "brand",
			Category:    "category",
			Price:       10,
			Description: "lorem ipsum",
			Title:       "title",
			Url:         "url",
		})
		assert.NotNil(t, err)
		var ce *customerror.Error
		if !assert.ErrorAs(t, err, &ce) {
			t.Fatalf("error should be of type ErrETagExists")
		}
		if errors.Is(ce, customerror.ErrIDExists) {
			t.Fatalf("error should be of type ErrETagExists")
		}
	})

	mt.Run("Case Create Error", func(mt *mtest.T) {
		mockCollection := productstorage.New(
			productstorage.WithDB(mt.DB),
			productstorage.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    2,
			Message: "unknown error",
		}))
		err := mockCollection.Create(context.TODO(), model.Product{
			ID:       1234321,
			Brand:    "brand",
			Category: "category",
		})
		assert.NotNil(t, err)
		var ce *customerror.Error
		if !assert.ErrorAs(t, err, &ce) {
			t.Fatalf("error should be of type ErrCreateObjectInfo")
		}
		if errors.Is(ce, customerror.ErrCreateObjectInfo) {
			t.Fatalf("error should be of type ErrCreateObjectInfo")
		}
	})

	mt.Run("Case Success Create", func(mt *mtest.T) {
		mockCollection := productstorage.New(
			productstorage.WithDB(mt.DB),
			productstorage.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		err := mockCollection.Create(context.TODO(), model.Product{
			ID:          1234321,
			Brand:       "brand",
			Category:    "category",
			Price:       10,
			Description: "lorem ipsum",
			Title:       "title",
			Url:         "url",
		})
		assert.Nil(t, err)
		mt.ClearEvents()
	})
}
