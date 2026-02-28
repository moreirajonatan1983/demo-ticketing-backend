package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/demoticketing/events/internal/adapters/handlers"
	"github.com/demoticketing/events/internal/adapters/repositories"
	"github.com/demoticketing/events/internal/core/services"
)

// @title Events API
// @version 1.0
// @description Ticketera Cloud Events Microservice API
// @host localhost:3000
// @BasePath /

func main() {
	// 1. Dependency Injection setup
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	var client *dynamodb.Client
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint != "" {
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(cfg)
	}
	tableName := os.Getenv("EVENTS_TABLE_NAME")
	if tableName == "" {
		// Fallback for local testing if not provided
		tableName = "EventsTable"
	}

	repo := repositories.NewDynamoDBEventRepository(client, tableName)
	service := services.NewEventService(repo)
	handler := handlers.NewHTTPHandler(service)

	// 2. Start Lambda
	lambda.Start(handler.HandleHTTPRequest)
}
