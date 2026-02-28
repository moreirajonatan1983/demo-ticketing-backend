package main

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
	Venue  string `json:"venue"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

var mockEvents = []Event{
	{
		ID:     1,
		Title:  "COLDPLAY - Music of the Spheres",
		Date:   "15 OCT 2026",
		Venue:  "Estadio Nacional",
		Image:  "https://images.unsplash.com/photo-1540039155733-d7696d819924?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
		Status: "Disponible",
	},
	{
		ID:     2,
		Title:  "THE WEEKND - After Hours Tour",
		Date:   "22 NOV 2026",
		Venue:  "Movistar Arena",
		Image:  "https://images.unsplash.com/photo-1493225457124-a1a2a5f5f4b2?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
		Status: "Pocos Tickets",
	},
	{
		ID:     3,
		Title:  "DUA LIPA - Radical Optimism",
		Date:   "04 DIC 2026",
		Venue:  "Hipódromo",
		Image:  "https://images.unsplash.com/photo-1514525253161-7a46d19cd819?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80",
		Status: "Sold Out",
	},
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Add CORS headers
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, OPTIONS",
		"Content-Type":                 "application/json",
	}

	// Route based on resource/path parameters
	id, hasId := request.PathParameters["id"]

	var body []byte
	var err error

	if hasId {
		// Mock single event selection
		if id == "1" {
			body, err = json.Marshal(mockEvents[0])
		} else {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusNotFound, Headers: headers, Body: `{"error": "Event not found"}`}, nil
		}
	} else {
		// Mock all events
		body, err = json.Marshal(mockEvents)
	}

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Headers: headers, Body: `{"error": "Failed to marshal response"}`}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

func main() {
	lambda.Start(handler)
}
