package model

type Product struct {
	Id          int     `json:"id"`
	Title       string  `json:"title"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Brand       string  `json:"brand"`
	Url         string  `json:"url"`
	Description string  `json:"description"`
}
