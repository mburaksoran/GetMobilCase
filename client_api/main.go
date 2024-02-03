package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gofiber/fiber/v2"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/app/client_api"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/app/config"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/service"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/infra/producer"
	"go.uber.org/zap/zapcore"
	"log"
	"time"

	"go.uber.org/zap"
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

	cfg, err := config.InitConfig()
	if err != nil {
		lgr.Error("err in initConfig:", err)
	}

	// Instantiate client.
	client, err := PrepareSQSClient(cfg, lgr)

	if err != nil {
		lgr.Error(err)
	}
	sqsProducer := producer.NewSqsOrderMessageProducer(client, cfg.SqsHost, lgr)
	orderService := service.NewOrderService(*sqsProducer, lgr)
	orderHandler := client_api.NewOrderHandler(orderService, lgr)
	authHandler := client_api.NewGoogleAuthHandler(cfg, lgr)

	app := fiber.New()
	app.Get("/hc", func(c *fiber.Ctx) error {
		return c.SendString("api is running")
	})

	app.Get("/google_login", authHandler.GoogleLogin)

	app.Get("/google_callback", authHandler.GoogleCallback)

	app.Post("/order", orderHandler.CreateOrder)

	app.Listen(":8080")

}

func PrepareSQSClient(cfg *config.AppConfig, lgr *zap.SugaredLogger) (*sqs.SQS, error) {
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
