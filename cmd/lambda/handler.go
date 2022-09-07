package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alexhiggins/aws/internal/event"
	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DocumentGenerator interface {
	Create(ctx context.Context, message event.Message, fileName string) error
}

type LambdaHandler func(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error)

func NewLambdaHandler(logger *zap.Logger, generator DocumentGenerator) LambdaHandler {
	return func(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
		var response events.SQSEventResponse

		for i := range sqsEvent.Records {
			msg := sqsEvent.Records[i]

			var payload events.SNSEntity
			if err := json.Unmarshal([]byte(msg.Body), &payload); err != nil {
				logger.Error(
					"unable to unmarshal SNSEntity",
					zap.Error(err),
					zap.String("messageBody", msg.Body),
				)

				bif := events.SQSBatchItemFailure{ItemIdentifier: msg.MessageId}
				response.BatchItemFailures = append(response.BatchItemFailures, bif)
				continue
			}

			var evt event.Message
			if err := json.Unmarshal([]byte(payload.Message), &evt); err != nil {
				logger.Error(
					"unable to unmarshal invoice message",
					zap.Error(err),
					zap.String("messageBody", payload.Message),
				)
			}

			fileName := fmt.Sprintf("%s-%s.txt", evt.ChipUserId, uuid.New().String())
			if err := generator.Create(ctx, evt, fileName); err != nil {
				logger.Error(
					"unable to generate message",
					zap.Error(err),
				)
				continue
			}
		}

		return response, nil
	}
}
