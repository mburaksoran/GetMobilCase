package sqs_client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/service"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/types"
)

func eventHandlerFactory(eventName string, orderService service.OrderServiceInterface) (EventHandler, error) {
	switch eventName {
	case types.OrderCompletedEvent:
		return NewOrderCompletedEventHandler(orderService), nil
	case types.OrderCreatedEvent:
		return NewOrderCreatedEventHandler(orderService), nil
	default:
		return nil, fmt.Errorf("no event handler exists for eventName %s", eventName)
	}
}

type EventHandler interface {
	HandleEvent(ctx context.Context, message *sqs.Message) (err error)
}
