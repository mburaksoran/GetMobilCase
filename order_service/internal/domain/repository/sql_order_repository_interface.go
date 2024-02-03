package repository

import "github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"

type OrderRepositoryInterface interface {
	CreateOrder(order models.Order) error
	GetOrder(id int) (*models.Order, error)
	DeleteOrder(id int) error
}
