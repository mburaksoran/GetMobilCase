package sqs_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models/messages"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/service"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/utils"
)

type OrderCreatedEventHandler struct {
	OrderService service.OrderServiceInterface
}

func (rc *OrderCreatedEventHandler) HandleEvent(ctx context.Context, message *sqs.Message) error {

	orderCreatedEvent := messages.OrderCreatedEvent{}
	err := json.Unmarshal([]byte(*message.Body), &orderCreatedEvent)
	if err != nil {
		return fmt.Errorf("OrderCreatedEventHandler - failed to unmarshal %w", err.Error())
	}
	order := utils.OrderCreatedMessageToOrderModel(orderCreatedEvent)
	err = rc.OrderService.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("OrderCreatedEventHandler - failed to CreateOrder %w", err.Error())
	}
	return nil
}

func NewOrderCreatedEventHandler(orderService service.OrderServiceInterface) *OrderCreatedEventHandler {
	return &OrderCreatedEventHandler{OrderService: orderService}
}
