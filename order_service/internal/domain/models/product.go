package models

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	StockCount int     `json:"stock_count"`
	Price      float32 `json:"price"`
}

func (p *Product) GetStock() int {
	return p.StockCount
}
func (p *Product) GetPrice() float32 {
	return p.Price
}
