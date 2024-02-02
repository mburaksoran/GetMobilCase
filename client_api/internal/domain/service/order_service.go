package service

import (
	"context"
	"fmt"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/infra/producer"
)

type orderService struct {
	OrderCreateProducer producer.SqsOrderMessageProducer
}

func (os orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	err := os.OrderCreateProducer.OrderCreatedEvent(ctx, order)
	if err != nil {
		fmt.Println("error create order event in order service:", err)
	}
	return err
}
func NewOrderService(sqsProducer producer.SqsOrderMessageProducer) OrderServiceInterface {
	return &orderService{sqsProducer}
}
