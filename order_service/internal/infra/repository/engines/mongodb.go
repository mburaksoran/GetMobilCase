package engines

import (
	"context"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var mongoDbEngine *MongoDbEngine

type MongoDbEngine struct {
	Db     *mongo.Database
	DBname string
}

func GetMongoDbEngine() *MongoDbEngine {
	return mongoDbEngine
}
func SetupMongoDBEngine(cfg *config.AppConfig, lgr *zap.SugaredLogger) (*MongoDbEngine, error) {
	if mongoDbEngine == nil {
		lgr.Info("creating mongodb client")
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoDbUrl))
		if err != nil {
			lgr.Error(err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			lgr.Error(err)
		}
		engine := &MongoDbEngine{}
		engine.Db = client.Database(cfg.MongoDbName)
		mongoDbEngine = engine
		return mongoDbEngine, err
	}

	return mongoDbEngine, nil
}
