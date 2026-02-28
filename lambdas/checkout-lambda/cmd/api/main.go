package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/demoticketing/checkout/internal/adapters/handlers"
	"github.com/demoticketing/checkout/internal/core/services"
)

func main() {
	service := services.NewCheckoutService()
	handler := handlers.NewHTTPHandler(service)

	lambda.Start(handler.HandleHTTPRequest)
}
