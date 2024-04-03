package producthandler_test

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	productservice "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/service/product"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/transport/http/basehttphandler"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/transport/http/producthandler"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
	"testing"
)

func Test_productHandler_AddRoutes(t *testing.T) {
	app := fiber.New()
	handler := producthandler.New(
		producthandler.WithBaseHttpHandler(
			basehttphandler.New(
				basehttphandler.WithLogger(nil),
				basehttphandler.WithContextTimeout(10),
				basehttphandler.WithPackages(pkg.New(
					pkg.WithValidator(&mockValidator{}),
					pkg.WithResponse(&mockResponse{})),
				),
			),
		),
		producthandler.WithProductService(&mockProductService{
			getErr: fmt.Errorf("%w", mongo.ErrNoDocuments),
		}),
	)
	handler.AddRoutes(app)

	req := httptest.NewRequest("GET", "/api/v1/product/1", nil)

	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test: %v", err)
	}
	if res.StatusCode != 404 {
		t.Fatalf("expected status code 404, got %d", res.StatusCode)
	}
}

func Test_productHandler_GetProduct(t *testing.T) {
	type fields struct {
		BaseHttpHandler *basehttphandler.BaseHttpHandler
		productService  productservice.ProductService
	}
	type args struct {
		app    *fiber.App
		path   string
		body   []byte
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		status int
		body   []byte
	}{
		{
			name: "Test Validation Error should return 400",
			fields: fields{
				BaseHttpHandler: basehttphandler.New(
					basehttphandler.WithLogger(nil),
					basehttphandler.WithContextTimeout(10),
					basehttphandler.WithPackages(pkg.New(
						pkg.WithValidator(&mockValidator{
							bindAndValidateErr: fmt.Errorf("error"),
						}),
						pkg.WithResponse(&mockResponse{})),
					),
				),
				productService: &mockProductService{},
			},
			args: args{
				app:    fiber.New(),
				path:   "/api/v1/product/1",
				method: "GET",
			},
			status: 400,
			body:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := producthandler.New(
				producthandler.WithBaseHttpHandler(tt.fields.BaseHttpHandler),
				producthandler.WithProductService(tt.fields.productService),
			)
			handler.AddRoutes(tt.args.app)
			req := httptest.NewRequest(tt.args.method, tt.args.path, nil)
			res, err := tt.args.app.Test(req)
			if err != nil {
				t.Fatalf("failed to test: %v", err)
			}
			if res.StatusCode != tt.status {
				t.Fatalf("expected status code %d, got %d", tt.status, res.StatusCode)
			}
		})
	}
}
