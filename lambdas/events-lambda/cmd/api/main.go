package main

import (
	"github.com/aws/aws-lambda-go/lambda"
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
	repo := repositories.NewMockEventRepository()
	service := services.NewEventService(repo)
	handler := handlers.NewHTTPHandler(service)

	// 2. Start Lambda
	lambda.Start(handler.HandleHTTPRequest)
}
