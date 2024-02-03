package messages

type OrderCompletedEvent struct {
	Metadata OrderEventMetadata         `json:"metadata"`
	Content  OrderCompletedEventContent `json:"content"`
}

type OrderCompletedEventContent struct {
	OrderEventContent
}
