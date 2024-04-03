package producthandler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	dto "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/dto/product"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/releaseinfo"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *productHandler) AddRoutes(router fiber.Router) {
	router.Get(releaseinfo.GetProduct, h.GetProduct)
}

func (h *productHandler) GetProduct(c *fiber.Ctx) error {
	var (
		req dto.GetProductRequest
	)
	if err := h.Validator.BindAndValidate(c, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(h.Response.BasicError(err, fiber.StatusBadRequest))
	}
	res, err := h.productService.GetProduct(c.Context(), req)
	if err != nil {
		var mongoErr mongo.CommandError
		if !errors.As(err, &mongoErr) {
			return c.Status(fiber.StatusNotFound).JSON(h.Response.BasicError(err, fiber.StatusNotFound))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(h.Response.BasicError(err, fiber.StatusInternalServerError))
	}
	return c.Status(fiber.StatusOK).JSON(h.Response.Data(fiber.StatusOK, res))
}
