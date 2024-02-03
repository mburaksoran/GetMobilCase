package service

import "github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"

type OrderServiceInterface interface {
	CreateOrder(order models.Order) error
	CompleteOrder(order models.Order) error
}
