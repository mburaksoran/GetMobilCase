package service

import (
	"context"
	"errors"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/repository"
)

type OrderService struct {
	OrderRepository       repository.OrderRepositoryInterface
	ProductRepository     repository.ProductRepositoryInterface
	OrderStatusRepository repository.MongodbOrderRepositoryInterface
}

func (os *OrderService) CreateOrder(order models.Order) error {
	product, err := os.ProductRepository.GetProduct(order.ProductID)
	if err != nil {
		return errors.Join(errors.New("error while gathering products"), err)
	}
	if order.OrderedCount > product.GetStock() {
		return errors.New("ordered Count cannot be bigger than stock count")
	}
	order.Price = float32(order.OrderedCount) * product.GetPrice()

	err = os.ProductRepository.UpdateStockCount(order.OrderedCount, order.ProductID)
	if err != nil {
		return errors.Join(errors.New("error while updating stocks"), err)
	}

	err = os.OrderRepository.CreateOrder(order)
	if err != nil {
		if err != nil {
			return errors.Join(errors.New("error while producing order messages"), err)
		}
	}

	return nil
}

func (os *OrderService) CompleteOrder(order models.Order) error {
	err := os.OrderStatusRepository.Add(order, context.TODO())
	if err != nil {
		return errors.Join(errors.New("error while completing order"), err)
	}
	return nil
}

func NewOrderService(orderRepo repository.OrderRepositoryInterface, productRepo repository.ProductRepositoryInterface, orderStatusRepo repository.MongodbOrderRepositoryInterface) OrderServiceInterface {
	return &OrderService{
		OrderRepository:       orderRepo,
		ProductRepository:     productRepo,
		OrderStatusRepository: orderStatusRepo,
	}
}
