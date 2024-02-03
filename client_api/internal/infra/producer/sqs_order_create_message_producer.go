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
	"go.uber.org/zap"
)

type SqsOrderMessageProducer struct {
	sqs    sqsiface.SQSAPI
	sqsUrl string
	logger *zap.SugaredLogger
}

func NewSqsOrderMessageProducer(sqs sqsiface.SQSAPI, QueueURL string, lgr *zap.SugaredLogger) *SqsOrderMessageProducer {
	lgr.Info("creating sqs client")
	return &SqsOrderMessageProducer{
		sqs:    sqs,
		sqsUrl: QueueURL,
		logger: lgr,
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
		p.logger.Error(err)
		return errors.New("error marshaling message")
	}

	_, err = p.sqs.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageJSON)),
		QueueUrl:    aws.String(p.sqsUrl),
	})

	if err != nil {
		p.logger.Error(err)
		return err
	}

	return nil
}
