package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/demoticketing/auth/internal/adapters/handlers"
	"github.com/demoticketing/auth/internal/adapters/repositories"
	"github.com/demoticketing/auth/internal/core/services"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
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

	tableName := os.Getenv("USERS_TABLE_NAME")
	if tableName == "" {
		tableName = "UsersTable"
	}

	repo := repositories.NewDynamoDBUserRepository(client, tableName)
	svc := services.NewAuthService(repo)
	handler := handlers.NewHTTPHandler(svc)

	lambda.Start(handler.HandleHTTPRequest)
}
