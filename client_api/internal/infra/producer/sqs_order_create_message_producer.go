package producer

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models/messages"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models/types"
)

type SqsOrderMessageProducer struct {
	sqs    sqsiface.SQSAPI
	sqsUrl string
}

func NewSqsOrderMessageProducer(sqs sqsiface.SQSAPI, QueueURL string) *SqsOrderMessageProducer {
	return &SqsOrderMessageProducer{
		sqs:    sqs,
		sqsUrl: QueueURL,
	}
}
func (p SqsOrderMessageProducer) OrderCreatedEvent(ctx context.Context, order *models.Order) error {

	message := messages.OrderCreatedEvent{
		Metadata: messages.OrderEventMetadata{
			Version:   "1",
			EventName: types.OrderCreatedEvent,
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

	_, err = p.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageJSON)),
		QueueUrl:    aws.String(p.sqsUrl),
	})

	if err != nil {
		return err
	}

	return nil
}
