package sqs_client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/app/config"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/messages"
	"github.com/mburaksoran/GetMobilCase/order_service/internal/domain/service"
)

type Consumer struct {
	Client  *sqs.SQS
	Cfg     *config.AppConfig
	Service service.OrderServiceInterface
}

func (c *Consumer) Start() {
	for {
		output, err := c.Client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(c.Cfg.SqsHost),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(10),
		})
		if err != nil {
			fmt.Println(err)
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
		fmt.Printf("error unmarshalling message body: %s", err.Error())
	}
	handler, factoryErr := eventHandlerFactory(metadataEvent.Metadata.EventName, c.Service)
	if factoryErr != nil {
		fmt.Printf("creating event handler: %s", factoryErr.Error())
	}
	handleErr := handler.HandleEvent(context.Background(), msg)
	if handleErr != nil {
		fmt.Printf("error handle event: %s", handleErr.Error())
	}
	deleteErr := c.DeleteMessage(msg.ReceiptHandle)
	if deleteErr != nil {
		fmt.Printf("error delete event: %s", deleteErr.Error())
	}
}

func (c *Consumer) DeleteMessage(receiptHandle *string) error {
	_, err := c.Client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.Cfg.SqsHost),
		ReceiptHandle: receiptHandle,
	})
	return err
}
