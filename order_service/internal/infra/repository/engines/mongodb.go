package engines

import (
	"context"
	"fmt"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDbEngine *MongoDbEngine

type MongoDbEngine struct {
	Db     *mongo.Database
	DBname string
}

func GetMongoDbEngine() *MongoDbEngine {
	return mongoDbEngine
}
func SetupMongoDBEngine(cfg *config.AppConfig) (*MongoDbEngine, error) {
	if mongoDbEngine == nil {

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoDbUrl))
		if err != nil {
			fmt.Println(err)
		}
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			fmt.Println(err)
		}
		engine := &MongoDbEngine{}
		engine.Db = client.Database(cfg.MongoDbName)
		mongoDbEngine = engine
		return mongoDbEngine, err
	}

	return mongoDbEngine, nil
}
