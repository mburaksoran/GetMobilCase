package utils

import (
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models/messages"
)

func OrderCreatedMessageToOrderModel(event messages.OrderCreatedEvent) models.Order {
	return models.Order{
		UserID:       event.Content.UserID,
		ProductID:    event.Content.ProductIDs,
		OrderedCount: event.Content.OrderedCount,
	}
}

func OrderCompletedMessageToOrderModel(event messages.OrderCompletedEvent) models.Order {
	return models.Order{
		UserID:       event.Content.UserID,
		ProductID:    event.Content.ProductIDs,
		OrderedCount: event.Content.OrderedCount,
	}
}
