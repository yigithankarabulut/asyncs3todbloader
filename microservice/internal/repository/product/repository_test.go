package productrepository_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	productrepository "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/repository/product"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"testing"
)

func Test_GetProductByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Case id not found returns error", func(mt *mtest.T) {
		mockCollection := productrepository.New(
			productrepository.WithDB(mt.DB),
			productrepository.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    2,
			Message: "product not found",
		}))
		_, err := mockCollection.GetProductByID(context.TODO(), 1)
		assert.NotNil(t, err)
		if !assert.ErrorAs(t, err, &mongo.ErrNoDocuments) {
			t.Fatalf("Case id not found returns error failed Want: %s, Got: %s", mongo.ErrNoDocuments.Error(), err.Error())
		}
	})

	mt.Run("Case id unknown error returns error", func(mt *mtest.T) {
		mockCollection := productrepository.New(
			productrepository.WithDB(mt.DB),
			productrepository.WithProductCollection("products"),
		)
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Code:    2,
			Message: "unknown error",
		}))
		_, err := mockCollection.GetProductByID(context.TODO(), 1)
		assert.NotNil(t, err)
	})
}
