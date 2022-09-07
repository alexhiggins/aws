package generator

import (
	"context"
	"encoding/json"
	"github.com/alexhiggins/aws/internal/event"
	"github.com/alexhiggins/aws/internal/store"
	"go.uber.org/zap"
)

type InvoiceGenerator struct {
	Logger  *zap.Logger
	Storage store.Storage
}

func (i *InvoiceGenerator) Create(ctx context.Context, message event.Message, fileName string) error {
	msg, err := json.Marshal(message)
	if err != nil {
		i.Logger.Error(
			"Unable to marshal event.Message",
			zap.Error(err),
		)
	}

	return i.Storage.Write(ctx, msg, fileName)
}

func NewInvoiceGenerator(logger *zap.Logger, storage store.Storage) *InvoiceGenerator {
	return &InvoiceGenerator{
		Logger:  logger,
		Storage: storage,
	}
}
