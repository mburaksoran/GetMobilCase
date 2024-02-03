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

	sqlEngine, err := prepareSqlDbEngine(cfg)
	if err != nil {
		fmt.Println(err)
	}

	sqlPingErr := sqlEngine.Client.Ping()
	if sqlPingErr != nil {
		fmt.Println("mysql connection failed :", sqlPingErr)
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

func Handler(message *sqs.Message) error {
	fmt.Println("handle", message)

	return nil
}

func pollSqs(chn chan<- *sqs.Message, client *sqs.SQS, cfg *config.AppConfig) {

	for {
		output, err := client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(cfg.SqsHost),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(10),
		})

		if err != nil {
			fmt.Println(err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

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

//func prepareConsumerConfig(cfg *config.AppConfig) Consumer.ConsumerConfig {
//	return Consumer.ConsumerConfig{
//		Type:      Consumer.SyncConsumer,
//		QueueURL:  cfg.SqsHost,
//		MaxWorker: cfg.SqsMaxWorkerCount,
//		MaxMsg:    cfg.SqsMaxMessageCount,
//	}
//}
