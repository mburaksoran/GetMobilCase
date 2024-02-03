package main

import (
	"context"
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
)

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.InitFromConfigFile()
	if err != nil {
		fmt.Println(err)
	}

	_, err = prepareSqlDbEngine(cfg)
	if err != nil {
		fmt.Println(err)
	}

	_, err = prepareMongoDbEngine(cfg)
	if err != nil {
		fmt.Println(err)
	}

	repoMongo := repository.NewMongodbOrderRepository(cfg)
	repo := repository.NewOrderRepository()
	productRepo := repository.NewProductRepository()

	orderServ := service.NewOrderService(repo, productRepo, repoMongo)
	client, err := PrepareSQSClient(cfg)
	if err != nil {
		fmt.Println(err)
	}
	cons := sqs_client.Consumer{
		Client:  client,
		Cfg:     cfg,
		Service: orderServ,
	}
	cons.Start()

}

func prepareMongoDbEngine(cfg *config.AppConfig) (*engines.MongoDbEngine, error) {
	return engines.SetupMongoDBEngine(cfg)
}

func prepareSqlDbEngine(cfg *config.AppConfig) (*engines.SqlDbEngine, error) {
	return engines.SetupSqlDBEngine(cfg)
}

func PrepareSQSClient(cfg *config.AppConfig) (*sqs.SQS, error) {
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
