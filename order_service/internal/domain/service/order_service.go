package service

import (
	"context"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/repository"
	"go.uber.org/zap"
)

type OrderService struct {
	OrderRepository       repository.OrderRepositoryInterface
	ProductRepository     repository.ProductRepositoryInterface
	OrderStatusRepository repository.MongodbOrderRepositoryInterface
	logger                *zap.SugaredLogger
}

func (os *OrderService) CreateOrder(order models.Order) error {
	product, err := os.ProductRepository.GetProduct(order.ProductID)
	if err != nil {
		os.logger.Error("error while gathering products", err)
		return err
	}
	if order.OrderedCount > product.GetStock() {
		os.logger.Error("ordered Count cannot be bigger than stock count", err)
		return err
	}
	order.Price = float32(order.OrderedCount) * product.GetPrice()

	err = os.ProductRepository.UpdateStockCount(order.OrderedCount, order.ProductID)
	if err != nil {
		os.logger.Error("error while updating stocks", err)
		return err
	}

	err = os.OrderRepository.CreateOrder(order)
	if err != nil {
		if err != nil {
			os.logger.Error("error while producing order messages", err)
			return err
		}
	}

	return nil
}

func (os *OrderService) CompleteOrder(order models.Order) error {
	err := os.OrderStatusRepository.Add(order, context.TODO())
	if err != nil {
		os.logger.Error("error while completing order", err)
		return err
	}
	return nil
}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, productRepo repository.ProductRepositoryInterface, orderStatusRepo repository.MongodbOrderRepositoryInterface, lgr *zap.SugaredLogger) OrderServiceInterface {
	return &OrderService{
		OrderRepository:       orderRepo,
		ProductRepository:     productRepo,
		OrderStatusRepository: orderStatusRepo,
		logger:                lgr,
	}
}
