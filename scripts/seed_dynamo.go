package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	eventsTable := os.Getenv("EVENTS_TABLE_NAME")
	if eventsTable == "" {
		eventsTable = "events-lambda-EventsTable-123456" // Replace with actual table
	}

	seatsTable := os.Getenv("SEATS_TABLE_NAME")
	if seatsTable == "" {
		seatsTable = "seats-lambda-EventSeatsTable-123456" // Replace with actual table
	}

	fmt.Println("Seeding Events...")
	seedEvents(client, eventsTable)

	fmt.Println("Seeding Seats...")
	seedSeats(client, seatsTable)

	fmt.Println("Seeding Complete!")
}

func seedEvents(client *dynamodb.Client, tableName string) {
	events := []map[string]types.AttributeValue{
		{
			"id":     &types.AttributeValueMemberS{Value: "1"},
			"title":  &types.AttributeValueMemberS{Value: "COLDPLAY - Music of the Spheres"},
			"date":   &types.AttributeValueMemberS{Value: "15 Octubre 2026"},
			"venue":  &types.AttributeValueMemberS{Value: "Estadio Nacional"},
			"image":  &types.AttributeValueMemberS{Value: "https://images.unsplash.com/photo-1540039155733-d7696d819924?auto=format&fit=crop&w=800&q=80"},
			"status": &types.AttributeValueMemberS{Value: "Sold Out"},
		},
		{
			"id":     &types.AttributeValueMemberS{Value: "2"},
			"title":  &types.AttributeValueMemberS{Value: "The Weeknd - After Hours"},
			"date":   &types.AttributeValueMemberS{Value: "02 Noviembre 2026"},
			"venue":  &types.AttributeValueMemberS{Value: "Movistar Arena"},
			"image":  &types.AttributeValueMemberS{Value: "https://images.unsplash.com/photo-1493225457124-a1a2a5f5f4b5?auto=format&fit=crop&w=800&q=80"},
			"status": &types.AttributeValueMemberS{Value: "Ultimos Tickets"},
		},
		{
			"id":     &types.AttributeValueMemberS{Value: "3"},
			"title":  &types.AttributeValueMemberS{Value: "Dua Lipa - Radical Optimism"},
			"date":   &types.AttributeValueMemberS{Value: "10 Diciembre 2026"},
			"venue":  &types.AttributeValueMemberS{Value: "Estadio Bicentenario"},
			"image":  &types.AttributeValueMemberS{Value: "https://images.unsplash.com/photo-1459749411175-04bf5292ceea?auto=format&fit=crop&w=800&q=80"},
			"status": &types.AttributeValueMemberS{Value: "Disponible"},
		},
	}

	for _, item := range events {
		_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})
		if err != nil {
			log.Printf("Failed to insert event %s: %v", item["id"].(*types.AttributeValueMemberS).Value, err)
		}
	}
}

func seedSeats(client *dynamodb.Client, tableName string) {
	eventId := "1"
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	seatsPerRow := 12

	for rIdx, rName := range rows {
		for s := 1; s <= seatsPerRow; s++ {
			status := "available"

			// Simular mock original
			if rIdx < 2 && s > 3 && s < 9 {
				status = "occupied"
			} else if rIdx == 4 && s > 8 {
				status = "occupied"
			} else if rIdx == 6 && (s == 2 || s == 3) {
				status = "processing"
			}

			seatId := fmt.Sprintf("%s%d", rName, s)

			_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
				TableName: aws.String(tableName),
				Item: map[string]types.AttributeValue{
					"event_id": &types.AttributeValueMemberS{Value: eventId},
					"seat_id":  &types.AttributeValueMemberS{Value: seatId},
					"row":      &types.AttributeValueMemberS{Value: rName},
					"number":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", s)},
					"status":   &types.AttributeValueMemberS{Value: status},
				},
			})
			if err != nil {
				log.Printf("Failed to insert seat %s for event %s: %v", seatId, eventId, err)
			}
		}
	}
}
