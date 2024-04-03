package dto

// GetProductRequest is the request type for GetProduct.
type GetProductRequest struct {
	ID int `json:"-" query:"-" params:"id" validate:"required"`
}
