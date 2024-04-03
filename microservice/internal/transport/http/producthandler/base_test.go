package producthandler_test

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	dto "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/dto/product"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/response"
)

var ErrGetProduct = errors.New("error getting product")

type mockProductService struct {
	getErr error
	getRes dto.ProductResponse
}

func (m *mockProductService) GetProduct(ctx context.Context, req dto.GetProductRequest) (dto.ProductResponse, error) {
	return m.getRes, m.getErr
}

type mockValidator struct {
	bindAndValidateErr error
}

func (m *mockValidator) BindAndValidate(c *fiber.Ctx, req interface{}) error {
	return m.bindAndValidateErr
}

type mockResponse struct {
	basicErrorRes response.ErrorResponse
	dataRes       response.DataResponse
}

func (m *mockResponse) BasicError(err interface{}, statusCode int) response.ErrorResponse {
	return m.basicErrorRes
}

func (m *mockResponse) Data(status int, data interface{}) response.DataResponse {
	return m.dataRes
}
