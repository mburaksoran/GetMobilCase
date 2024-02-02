package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gofiber/fiber/v2"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/app/client_api"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/app/config"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/service"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/infra/producer"
)

func main() {
	app := fiber.New()

	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println("err in initConfig:", err)
	}

	// Instantiate client.
	client, err := PrepareSQSClient(cfg)

	if err != nil {
		fmt.Println(err)
	}
	sqsProducer := producer.NewSqsOrderMessageProducer(client, cfg.SqsHost)
	orderService := service.NewOrderService(*sqsProducer)
	orderHandler := client_api.NewOrderHandler(orderService)
	authHandler := client_api.NewGoogleAuthHandler(cfg)

	app.Get("/hc", func(c *fiber.Ctx) error {
		return c.SendString("api is running")
	})

	app.Get("/google_login", authHandler.GoogleLogin)

	app.Get("/google_callback", authHandler.GoogleCallback)

	app.Post("/order", orderHandler.CreateOrder)

	app.Listen(":8080")

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
