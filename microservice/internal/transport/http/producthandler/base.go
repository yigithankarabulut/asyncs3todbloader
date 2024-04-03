package producthandler

import (
	"github.com/gofiber/fiber/v2"
	productservice "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/service/product"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/transport/http/basehttphandler"
)

// ProductHandler is the interface for product handler.
type ProductHandler interface {
	AddRoutes(router fiber.Router)
	GetProduct(c *fiber.Ctx) error
}

// productHandler is the http handler for product operations.
type productHandler struct {
	*basehttphandler.BaseHttpHandler
	productService productservice.ProductService
}

// Option is the option type for user handler.
type Option func(*productHandler)

// WithBaseHttpHandler sets the base http handler option.
func WithBaseHttpHandler(handler *basehttphandler.BaseHttpHandler) Option {
	return func(h *productHandler) {
		h.BaseHttpHandler = handler
	}
}

// WithProductService sets the product service option.
func WithProductService(productService productservice.ProductService) Option {
	return func(h *productHandler) {
		h.productService = productService
	}
}

// New creates a new product handler.
func New(opts ...Option) ProductHandler {
	handler := &productHandler{}
	for _, opt := range opts {
		opt(handler)
	}
	return handler
}
