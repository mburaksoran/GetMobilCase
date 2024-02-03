package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/service"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/repository"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/repository/engines"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/infra/sqs_client"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func main() {

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	lgr := logger.Sugar()

	cfg, err := config.InitFromConfigFile()
	if err != nil {
		lgr.Error(err)
	}

	_, err = prepareSqlDbEngine(cfg, lgr)
	if err != nil {
		lgr.Error(err)
	}

	_, err = prepareMongoDbEngine(cfg, lgr)
	if err != nil {
		lgr.Error(err)
	}

	repoMongo := repository.NewMongodbOrderRepository(cfg, lgr)
	repo := repository.NewOrderRepository(lgr)
	productRepo := repository.NewProductRepository(lgr)

	orderServ := service.NewOrderService(repo, productRepo, repoMongo, lgr)
	client, err := PrepareSQSClient(cfg, lgr)
	if err != nil {
		fmt.Println(err)
	}
	cons := sqs_client.NewConsumer(client, cfg, orderServ, lgr)
	cons.Start()

}

func prepareMongoDbEngine(cfg *config.AppConfig, lgr *zap.SugaredLogger) (*engines.MongoDbEngine, error) {
	return engines.SetupMongoDBEngine(cfg, lgr)
}

func prepareSqlDbEngine(cfg *config.AppConfig, lgr *zap.SugaredLogger) (*engines.SqlDbEngine, error) {
	return engines.SetupSqlDBEngine(cfg, lgr)
}

func PrepareSQSClient(cfg *config.AppConfig, lgr *zap.SugaredLogger) (*sqs.SQS, error) {
	lgr.Info("creating sqs client")
	if cfg.AwsConfig.Region == "" {
		return nil, nil
	}

	awsConfig := aws.NewConfig().
		WithRegion(cfg.AwsConfig.Region)

	if cfg.AwsConfig.Key != "" && cfg.AwsConfig.Secret != "" {
		awsCredentials := credentials.NewStaticCredentials(
			cfg.AwsConfig.Key,
			cfg.AwsConfig.Secret,
			"",
		)
		awsConfig = awsConfig.WithCredentials(awsCredentials)
	}

	awsConfig.Endpoint = aws.String(cfg.SqsHost)

	sess := session.Must(session.NewSession(awsConfig))

	return sqs.New(sess), nil
}
