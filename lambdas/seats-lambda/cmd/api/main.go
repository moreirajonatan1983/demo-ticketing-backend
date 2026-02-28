package main

import (
	"github.com/aws/aws-lambda-go/lambda"
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
	repo := repositories.NewMockSeatRepository()
	service := services.NewSeatService(repo)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
