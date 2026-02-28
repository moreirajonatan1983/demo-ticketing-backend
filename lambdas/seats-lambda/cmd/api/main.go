package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/demoticketing/seats/internal/adapters/handlers"
	"github.com/demoticketing/seats/internal/adapters/repositories"
	"github.com/demoticketing/seats/internal/core/services"
)

func main() {
	repo := repositories.NewMockSeatRepository()
	service := services.NewSeatService(repo)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
