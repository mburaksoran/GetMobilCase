package messages

type OrderCreatedEvent struct {
	Metadata OrderEventMetadata       `json:"metadata"`
	Content  OrderCreatedEventContent `json:"content"`
}

type OrderCreatedEventContent struct {
	OrderEventContent
}
