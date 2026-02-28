package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/demoticketing/shows/internal/adapters/handlers"
	"github.com/demoticketing/shows/internal/adapters/repositories"
	"github.com/demoticketing/shows/internal/core/services"
)

// @title Shows API
// @version 1.0
// @description Ticketera Cloud Shows Microservice API
// @host localhost:3007
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
	tableName := os.Getenv("SHOWS_TABLE_NAME")
	if tableName == "" {
		tableName = "ShowsTable"
	}

	repo := repositories.NewDynamoDBShowRepository(client, tableName)
	service := services.NewShowService(repo)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
