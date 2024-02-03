package messages

type OrderCompletedEvent struct {
	Metadata OrderEventMetadata         `json:"metadata"`
	Content  OrderCompletedEventContent `json:"content"`
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

type OrderCompletedEventContent struct {
	OrderEventContent
}
