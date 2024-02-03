package models

type Order struct {
	ID           int     `json:"id" bson:"id,omitempty"`
	UserID       int     `json:"user_id" bson:"user_id,omitempty" `
	ProductID    int     `json:"product_id" bson:"product_id,omitempty"`
	OrderedCount int     `json:"ordered_count" bson:"ordered_count,omitempty"`
	Price        float32 `json:"price" bson:"price,omitempty"`
}
