package models

type Product struct {
	ID                 int     `json:"product_id"`
	ProductName        string  `json:"product_name"`
	ProductPrice       float64 `json:"product_price"`
	ProductDescription string  `json:"product_description"`
}

type Products struct {
	Products []Product `json:"products"`
}

type ProductsResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    Products `json:"data"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
