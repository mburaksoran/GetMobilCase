package models

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	StockCount int     `json:"stock_count"`
	Price      float32 `json:"price"`
}
