package utils

import (
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/messages"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
)

func OrderMessageToOrderModel(event messages.OrderCreatedEvent) models.Order {
	return models.Order{
		UserID:       event.Content.UserID,
		ProductID:    event.Content.ProductIDs,
		OrderedCount: event.Content.OrderedCount,
	}
}
