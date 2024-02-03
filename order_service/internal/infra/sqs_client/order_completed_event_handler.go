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

type OrderCompletedEventHandler struct {
	OrderService service.OrderServiceInterface
}

func (rc *OrderCompletedEventHandler) HandleEvent(ctx context.Context, message *sqs.Message) error {

	orderCompletedEvent := messages.OrderCompletedEvent{}
	err := json.Unmarshal([]byte(*message.Body), &orderCompletedEvent)
	if err != nil {
		return fmt.Errorf("OrderCompletedEventHandler - failed to unmarshal %w", err.Error())
	}
	order := utils.OrderCompletedMessageToOrderModel(orderCompletedEvent)
	err = rc.OrderService.CompleteOrder(order)
	if err != nil {
		return fmt.Errorf("OrderCompletedEventHandler - failed to CompleteOrder %w", err.Error())
	}
	return nil
}

func NewOrderCompletedEventHandler(orderService service.OrderServiceInterface) *OrderCompletedEventHandler {
	return &OrderCompletedEventHandler{OrderService: orderService}
}
