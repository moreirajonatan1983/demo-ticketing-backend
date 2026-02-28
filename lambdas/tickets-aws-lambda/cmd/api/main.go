package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/demoticketing/tickets/internal/adapters/handlers"
	"github.com/demoticketing/tickets/internal/adapters/queues"
	"github.com/demoticketing/tickets/internal/adapters/repositories"
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

	repo := repositories.NewDynamoDBTicketRepository(dynamoClient, tableName)
	publisher := queues.NewSQSEventPublisher(sqsClient)
	service := services.NewTicketService(repo, publisher)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
