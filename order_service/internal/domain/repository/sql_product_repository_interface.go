package repository

import "github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"

type ProductRepositoryInterface interface {
	CreateProduct(product models.Product) error
	GetProduct(id int) (*models.Product, error)
	DeleteProduct(id int) error
	UpdateStockCount(value int, productId int) error
}
