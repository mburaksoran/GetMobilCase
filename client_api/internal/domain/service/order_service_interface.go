package service

import (
	"context"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, order *models.Order) error
}
