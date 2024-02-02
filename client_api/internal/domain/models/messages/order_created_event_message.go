package messages

type OrderCreatedEvent struct {
	Metadata OrderEventMetadata       `json:"metadata"`
	Content  OrderCreatedEventContent `json:"content"`
}

type OrderEventMetadata struct {
	Version   string `json:"version"`
	EventName string `json:"eventName"`
}

type MetadataEvent struct {
	Metadata OrderEventMetadata `json:"metadata"`
}

type OrderEventContent struct {
	ProductIDs   int     `json:"product_id"`
	OrderedCount int     `json:"ordered_count"`
	UserID       int     `json:"user_id"`
	Price        float32 `json:"price"`
}

type OrderCreatedEventContent struct {
	OrderEventContent
	Text string `json:"text" validate:"required"`
}
