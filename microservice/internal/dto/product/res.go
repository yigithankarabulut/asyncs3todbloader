package dto

import "github.com/yigithankarabulut/asyncs3todbloader/microservice/model"

// ProductResponse is the request type for GetProductByID.
type ProductResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Brand       string  `json:"brand"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
}

func (c *ProductResponse) Convert(product model.Product) {
	c.ID = product.ID
	c.Title = product.Title
	c.Price = product.Price
	c.Category = product.Category
	c.Brand = product.Brand
	c.Url = product.Url
	c.Description = product.Description
}
