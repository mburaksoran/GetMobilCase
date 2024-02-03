package sqs_client

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/models/messages"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/service"
	"go.uber.org/zap"
)

type Consumer struct {
	Client  *sqs.SQS
	Cfg     *config.AppConfig
	Service service.OrderServiceInterface
	logger  *zap.SugaredLogger
}

func NewConsumer(client *sqs.SQS, cfg *config.AppConfig, serv service.OrderServiceInterface, lgr *zap.SugaredLogger) Consumer {
	return Consumer{
		Client:  client,
		Cfg:     cfg,
		Service: serv,
		logger:  lgr,
	}
}

func (c *Consumer) Start() {
	for {
		output, err := c.Client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.Cfg.SqsHost),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(10),
		})
		if err != nil {
			c.logger.Error("error while gathering messages from queue", err)
		}

		for _, message := range output.Messages {
			c.handle(message)
		}
	}
}

func (c *Consumer) handle(msg *sqs.Message) {
	metadataEvent := messages.MetadataEvent{}
	err := json.Unmarshal([]byte(*msg.Body), &metadataEvent)
	if err != nil {
		c.logger.Error("error unmarshalling message body", err)
	}
	handler, factoryErr := eventHandlerFactory(metadataEvent.Metadata.EventName, c.Service)
	if factoryErr != nil {
		c.logger.Error("error creating event handler", factoryErr)
	}
	handleErr := handler.HandleEvent(context.Background(), msg)
	if handleErr != nil {
		c.logger.Error("error handle event:", handleErr)
	}
	deleteErr := c.DeleteMessage(msg.ReceiptHandle)
	if deleteErr != nil {
		c.logger.Error("error while deleting consumed messages from queue", err)
	}
}

func (c *Consumer) DeleteMessage(receiptHandle *string) error {
	_, err := c.Client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.Cfg.SqsHost),
		ReceiptHandle: receiptHandle,
	})
	return err
}
