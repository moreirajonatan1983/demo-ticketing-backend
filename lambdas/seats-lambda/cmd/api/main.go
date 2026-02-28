package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/demoticketing/seats/internal/adapters/handlers"
	"github.com/demoticketing/seats/internal/adapters/repositories"
	"github.com/demoticketing/seats/internal/core/services"
)

// @title Seats API
// @version 1.0
// @description Ticketera Cloud Seats Microservice API
// @host localhost:3005
// @BasePath /

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	tableName := os.Getenv("SEATS_TABLE_NAME")
	if tableName == "" {
		tableName = "EventSeatsTable"
	}

	repo := repositories.NewDynamoDBSeatRepository(client, tableName)
	service := services.NewSeatService(repo)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
