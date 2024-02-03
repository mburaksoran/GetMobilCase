package service

import (
	"context"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/infra/producer"
	"go.uber.org/zap"
)

type orderService struct {
	OrderCreateProducer producer.SqsOrderMessageProducer
	lgr                 *zap.SugaredLogger
}

func (os orderService) CreateOrder(ctx context.Context, order *models.Order) error {
	err := os.OrderCreateProducer.OrderCreatedEvent(ctx, order)
	if err != nil {
		os.lgr.Error("error create order event in order service:", err)
	}
	return err
}
func NewOrderService(sqsProducer producer.SqsOrderMessageProducer, lgr *zap.SugaredLogger) OrderServiceInterface {
	return &orderService{sqsProducer, lgr}
}
