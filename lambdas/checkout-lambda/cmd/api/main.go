package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/demoticketing/checkout/internal/adapters/handlers"
	"github.com/demoticketing/checkout/internal/core/services"
)

// @title Checkout API
// @version 1.0
// @description Ticketera Cloud Checkout Microservice API
// @host localhost:3004
// @BasePath /

func main() {
	service := services.NewCheckoutService()
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
