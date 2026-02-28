package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/demoticketing/events/internal/adapters/handlers"
	"github.com/demoticketing/events/internal/adapters/repositories"
	"github.com/demoticketing/events/internal/core/services"
)

func main() {
	// 1. Dependency Injection setup
	repo := repositories.NewMockEventRepository()
	service := services.NewEventService(repo)
	handler := handlers.NewHTTPHandler(service)

	// 2. Start Lambda
	lambda.Start(handler.HandleHTTPRequest)
}
