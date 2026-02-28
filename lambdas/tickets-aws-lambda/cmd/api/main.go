package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/demoticketing/tickets/internal/adapters/handlers"
	"github.com/demoticketing/tickets/internal/adapters/queues"
	"github.com/demoticketing/tickets/internal/adapters/repositories"
	s3adapter "github.com/demoticketing/tickets/internal/adapters/s3"
	"github.com/demoticketing/tickets/internal/core/services"
)

// @title Tickets API
// @version 1.0
// @description Ticketera Cloud Tickets Microservice API
// @host localhost:3006
// @BasePath /

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	endpoint := os.Getenv("AWS_ENDPOINT_URL")

	// DynamoDB client
	var dynamoClient *dynamodb.Client
	if endpoint != "" {
		dynamoClient = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		dynamoClient = dynamodb.NewFromConfig(cfg)
	}

	// SQS client
	sqsEndpoint := os.Getenv("SQS_ENDPOINT_URL")
	if sqsEndpoint == "" {
		sqsEndpoint = endpoint // fallback to same localstack endpoint
	}
	var sqsClient *sqs.Client
	if sqsEndpoint != "" {
		sqsClient = sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.BaseEndpoint = aws.String(sqsEndpoint)
		})
	} else {
		sqsClient = sqs.NewFromConfig(cfg)
	}

	tableName := os.Getenv("TICKETS_TABLE_NAME")
	if tableName == "" {
		tableName = "TicketsTable"
	}

	// S3 client for presigned PDF download URLs
	s3Endpoint := os.Getenv("S3_ENDPOINT_URL")
	if s3Endpoint == "" {
		s3Endpoint = endpoint
	}
	var s3Client *s3.Client
	if s3Endpoint != "" {
		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(s3Endpoint)
			o.UsePathStyle = true
		})
	} else {
		s3Client = s3.NewFromConfig(cfg)
	}

	repo := repositories.NewDynamoDBTicketRepository(dynamoClient, tableName)
	publisher := queues.NewSQSEventPublisher(sqsClient)
	storage := s3adapter.NewS3PresignedURLGenerator(s3Client)
	service := services.NewTicketService(repo, publisher)
	handler := handlers.NewHTTPHandler(service, storage)

	lambda.Start(handler.HandleHTTPRequest)
}
