package productrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *productRepository) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	var product model.Product
	if err := r.db.Collection(r.productCollectionName).FindOne(ctx, bson.M{"id": id}).Decode(&product); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return product, fmt.Errorf("%w", err)
		}
		return product, err
	}
	return product, nil
}
