package repository

import (
	"context"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/repository/engines"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MongodbOrderRepository struct {
	mongoDb         *mongo.Database
	mongoDbName     string
	mongoCollection string
	logger          *zap.SugaredLogger
}

func NewMongodbOrderRepository(cfg *config.AppConfig, lgr *zap.SugaredLogger) *MongodbOrderRepository {
	mongoDbEngine := engines.GetMongoDbEngine()
	return &MongodbOrderRepository{
		mongoDb:         mongoDbEngine.Db,
		mongoDbName:     cfg.MongoDbName,
		mongoCollection: cfg.MongoDbCollectionName,
		logger:          lgr,
	}
}

func (app *MongodbOrderRepository) Add(order models.Order, ctx context.Context) error {
	_, err := app.mongoDb.Collection(app.mongoCollection).InsertOne(ctx, fromModel(order))
	if err != nil {
		app.logger.Error("order create err :", err)
	}
	return err
}

func (app *MongodbOrderRepository) GetById(ctx context.Context, id int) (*models.Order, error) {
	var order *models.Order
	err := app.mongoDb.Collection(app.mongoCollection).FindOne(ctx, bson.D{{"id", id}}).Decode(&order)
	if err != nil {
		app.logger.Error(err)
		return nil, err
	}
	return order, nil
}

func (app *MongodbOrderRepository) Delete(id int, ctx context.Context) error {

	_, err := app.mongoDb.Collection(app.mongoCollection).DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		app.logger.Error(err)
		return err
	}
	return nil
}

type order struct {
	UID          primitive.ObjectID `bson:"_id,omitempty"`
	ID           int                `bson:"id,omitempty"`
	UserID       int                `bson:"user_id,omitempty"`
	ProductID    int                `bson:"product_id,omitempty"`
	OrderedCount int                `bson:"ordered_count,omitempty"`
	Price        float32            `bson:"price,omitempty"`
}

func fromModel(data models.Order) order {
	return order{
		ID:           data.ID,
		UserID:       data.UserID,
		ProductID:    data.ProductID,
		OrderedCount: data.OrderedCount,
		Price:        data.Price,
	}
}
