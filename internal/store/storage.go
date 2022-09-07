package store

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

type AWS struct {
	Region string
	Bucket string
}

type Storage interface {
	Write(ctx context.Context, contents []byte, fileName string) error
}

type S3Writer struct {
	Logger *zap.Logger
	Client *s3.S3
	Config AWS
}

func (s *S3Writer) Write(ctx context.Context, contents []byte, fileName string) error {
	_, err := s.Client.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(s.Config.Bucket),
		Key:                  aws.String(fileName),
		Body:                 bytes.NewReader(contents),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return err
	}
	return nil
}

type LogWriter struct {
	Logger *zap.Logger
}

func (s *LogWriter) Write(ctx context.Context, contents []byte, fileName string) error {
	s.Logger.Info(
		"Would have written",
		zap.String("fileName", fileName),
		zap.String("contents", string(contents[:])),
	)

	return nil
}

func NewS3Writer(logger *zap.Logger, config AWS) (*S3Writer, error) {
	s3Session, err := session.NewSession(&aws.Config{Region: aws.String(config.Region)})
	if err != nil {
		logger.Fatal(
			"unable to initialize new aws session",
			zap.Error(err),
		)
	}

	return &S3Writer{
		Logger: logger,
		Config: config,
		Client: s3.New(s3Session),
	}, nil
}

func NewLogWriter(logger *zap.Logger) (*LogWriter, error) {
	return &LogWriter{
		Logger: logger,
	}, nil
}
