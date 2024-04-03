package productservice_test

import (
	"context"
	"errors"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/model"
)

var ErrGetProductByID = errors.New("error getting product by id")

type mockProductRepository struct {
	getErr error
	getRes model.Product
}

func (m *mockProductRepository) GetProductByID(ctx context.Context, id int) (model.Product, error) {
	return m.getRes, m.getErr
}
