package productservice

import (
	"context"
	dto "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/dto/product"
	productrepository "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/repository/product"
)

type ProductService interface {
	GetProduct(ctx context.Context, req dto.GetProductRequest) (dto.ProductResponse, error)
}

type service struct {
	repository productrepository.ProductRepository
}

type Option func(*service)

func WithRepository(repository productrepository.ProductRepository) Option {
	return func(s *service) {
		s.repository = repository
	}
}

func New(opts ...Option) ProductService {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
