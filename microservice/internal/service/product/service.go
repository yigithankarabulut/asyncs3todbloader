package productservice

import (
	"context"
	"fmt"
	dto "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/dto/product"
)

func (s *service) GetProduct(ctx context.Context, req dto.GetProductRequest) (dto.ProductResponse, error) {
	var (
		res dto.ProductResponse
	)
	product, err := s.repository.GetProductByID(ctx, req.ID)
	if err != nil {
		return res, fmt.Errorf("%w. id: %d", err, req.ID)
	}
	res.Convert(product)
	return res, nil
}
