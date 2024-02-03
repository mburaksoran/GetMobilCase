package repository

import (
	"context"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
)

type MongodbOrderRepositoryInterface interface {
	Add(order models.Order, ctx context.Context) error
	GetById(ctx context.Context, id int) (*models.Order, error)
	Delete(id int, ctx context.Context) error
}
