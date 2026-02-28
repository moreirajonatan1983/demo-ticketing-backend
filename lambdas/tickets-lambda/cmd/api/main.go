package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/demoticketing/tickets/internal/adapters/handlers"
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

	var client *dynamodb.Client
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint != "" {
		client = dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	} else {
		client = dynamodb.NewFromConfig(cfg)
	}
	tableName := os.Getenv("TICKETS_TABLE_NAME")
	if tableName == "" {
		tableName = "TicketsTable"
	}

	repo := repositories.NewDynamoDBTicketRepository(client, tableName)
	service := services.NewTicketService(repo)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
