package main

import (
	"github.com/aws/aws-lambda-go/lambda"
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
	repo := repositories.NewMockTicketRepository()
	service := services.NewTicketService(repo)
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
