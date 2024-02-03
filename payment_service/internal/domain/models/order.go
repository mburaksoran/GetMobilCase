package models

type Order struct {
	ID           int     `json:"id"`
	UserID       int     `json:"user_id"`
	ProductID    int     `json:"product_id"`
	OrderedCount int     `json:"ordered_count"`
	Price        float32 `json:"price"`
}
