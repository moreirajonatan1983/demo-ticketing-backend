package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CheckoutRequest struct {
	EventId string   `json:"event_id"`
	Seats   []string `json:"seats"`
	Method  string   `json:"method"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, OPTIONS",
		"Content-Type":                 "application/json",
	}

	log.Printf("Received checkout payload: %s", request.Body)

	var req CheckoutRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Headers: headers, Body: `{"error": "Invalid JSON"}`}, nil
	}

	// Here we would typically validate payload, interact with DynamoDB to reserve seats lock
	// and trigger AWS Step Functions or EventBridge for asynchronous payment gateway

	responsePayload := map[string]interface{}{
		"message":  "Solicitud Recibida Async",
		"status":   "PROCESSING",
		"order_id": "ORD-XYZ-1234",
	}

	body, _ := json.Marshal(responsePayload)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusAccepted,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
