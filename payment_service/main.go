package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/mburaksoran/GetMobilCase/payment_service/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/payment_service/internal/domain/models/messages"
)

const URL string = "http://localstack:4566/000000000000/order_updates"
const OrderCompletedEventType string = "order_completed_event"

func main() {
	client, err := PrepareSQSClient()
	if err != nil {
		fmt.Println(err)
	}
	pollSqs(client)

}
func PrepareSQSClient() (*sqs.SQS, error) {

	awsConfig := aws.NewConfig().
		WithRegion("eu-west-1")

	awsCredentials := credentials.NewStaticCredentials(
		"test",
		"test",
		"",
	)
	awsConfig = awsConfig.WithCredentials(awsCredentials)

	awsConfig.Endpoint = aws.String(URL)

	sess := session.Must(session.NewSession(awsConfig))

	return sqs.New(sess), nil
}

func pollSqs(client *sqs.SQS) {

	for {
		output, err := client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(URL),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(10),
		})

		if err != nil {
			fmt.Println(err)
		}
		for _, message := range output.Messages {
			HandleMessage(message, client)
		}

	}

}

func HandleMessage(message *sqs.Message, client *sqs.SQS) {
	orderCreatedEvent := messages.OrderCreatedEvent{}
	err := json.Unmarshal([]byte(*message.Body), &orderCreatedEvent)
	if err != nil {
		fmt.Errorf("OrderCompletedEvent - failed to unmarshal %w", err.Error())
	}

	order := OrderMessageToOrderModel(orderCreatedEvent)
	OrderCompletedEvent(&order, client)
}

func OrderCompletedEvent(order *models.Order, client *sqs.SQS) error {

	message := messages.OrderCreatedEvent{
		Metadata: messages.OrderEventMetadata{
			Version:   "1",
			EventName: OrderCompletedEventType,
		},
		Content: messages.OrderCreatedEventContent{
			OrderEventContent: messages.OrderEventContent{
				ProductIDs:   order.ProductID,
				OrderedCount: order.OrderedCount,
				UserID:       order.UserID,
				Price:        order.Price,
			},
		},
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return errors.New("error marshaling message")
	}

	_, err = client.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageJSON)),
		QueueUrl:    aws.String(URL),
	})

	if err != nil {
		return err
	}

	return nil
}

func OrderMessageToOrderModel(event messages.OrderCreatedEvent) models.Order {
	return models.Order{
		UserID:       event.Content.UserID,
		ProductID:    event.Content.ProductIDs,
		OrderedCount: event.Content.OrderedCount,
	}
}
