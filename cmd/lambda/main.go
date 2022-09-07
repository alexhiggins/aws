package main

import (
	"github.com/alexhiggins/aws/internal/generator"
	"github.com/alexhiggins/aws/internal/store"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	storage, err := store.NewS3Writer(
		logger,
		store.AWS{
			Region: os.Getenv("S3_REGION"),
			Bucket: os.Getenv("S3_BUCKET"),
		},
	)

	if err != nil {
		logger.Panic(
			"Unable to initialize storage writer",
			zap.Error(err),
		)
	}

	documentGenerator := generator.NewInvoiceGenerator(logger, storage)

	lambda.Start(NewLambdaHandler(logger, documentGenerator))
}
