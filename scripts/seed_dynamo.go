package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	eventsTable := os.Getenv("EVENTS_TABLE_NAME")
	if eventsTable == "" {
		eventsTable = "events-aws-lambda-EventsTable-123456" // Replace with actual table
	}

	seatsTable := os.Getenv("SEATS_TABLE_NAME")
	if seatsTable == "" {
		seatsTable = "seats-aws-lambda-EventSeatsTable-123456" // Replace with actual table
	}

	showsTable := os.Getenv("SHOWS_TABLE_NAME")
	if showsTable == "" {
		showsTable = "shows-aws-lambda-ShowsTable-123456"
	}

	fmt.Println("Seeding Events (with S3 Images)...")

	// Create S3 client
	s3Endpoint := os.Getenv("LOCALSTACK_ENDPOINT")
	if s3Endpoint == "" {
		s3Endpoint = "http://localhost:4566"
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3Endpoint)
		o.UsePathStyle = true
	})

	seedEvents(client, s3Client, eventsTable)

	fmt.Println("Seeding Shows...")
	seedShows(client, showsTable)

	fmt.Println("Seeding Seats...")
	seedSeats(client, seatsTable)

	fmt.Println("Seeding Complete!")
}

func seedEvents(client *dynamodb.Client, s3Client *s3.Client, tableName string) {
	bucketName := "ticketera-images-local"

	events := []map[string]interface{}{
		{
			"id":     "1",
			"title":  "COLDPLAY - Music of the Spheres",
			"date":   "15 Octubre 2026",
			"venue":  "Estadio Nacional",
			"image":  "https://images.unsplash.com/photo-1540039155733-d7696d819924?auto=format&fit=crop&w=800&q=80",
			"status": "Sold Out",
			"key":    "events/1/cover.jpg",
		},
		{
			"id":     "2",
			"title":  "The Weeknd - After Hours",
			"date":   "02 Noviembre 2026",
			"venue":  "Movistar Arena",
			"image":  "https://images.unsplash.com/photo-1493225457124-a1a2a5f5f4b5?auto=format&fit=crop&w=800&q=80",
			"status": "Ultimos Tickets",
			"key":    "events/2/cover.jpg",
		},
		{
			"id":     "3",
			"title":  "Dua Lipa - Radical Optimism",
			"date":   "10 Diciembre 2026",
			"venue":  "Estadio Bicentenario",
			"image":  "https://images.unsplash.com/photo-1459749411175-04bf5292ceea?auto=format&fit=crop&w=800&q=80",
			"status": "Disponible",
			"key":    "events/3/cover.jpg",
		},
	}

	for _, item := range events {
		// Download image
		resp, err := http.Get(item["image"].(string))
		if err != nil {
			log.Printf("Failed to download image %s: %v", item["image"], err)
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		// Upload to S3
		_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(item["key"].(string)),
			Body:   strings.NewReader(string(bodyBytes)),
		})
		if err != nil {
			log.Printf("Failed to upload to S3: %v", err)
		}

		// Insert to Dynamo with the key instead of external URL
		dynamoItem := map[string]types.AttributeValue{
			"id":     &types.AttributeValueMemberS{Value: item["id"].(string)},
			"title":  &types.AttributeValueMemberS{Value: item["title"].(string)},
			"date":   &types.AttributeValueMemberS{Value: item["date"].(string)},
			"venue":  &types.AttributeValueMemberS{Value: item["venue"].(string)},
			"image":  &types.AttributeValueMemberS{Value: item["key"].(string)}, // stored in Dynamo!!
			"status": &types.AttributeValueMemberS{Value: item["status"].(string)},
		}

		_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      dynamoItem,
		})
		if err != nil {
			log.Printf("Failed to insert event %s: %v", item["id"], err)
		}
	}
}

func seedShows(client *dynamodb.Client, tableName string) {
	eventId := "1"
	shows := []map[string]types.AttributeValue{
		{
			"event_id": &types.AttributeValueMemberS{Value: eventId},
			"id":       &types.AttributeValueMemberS{Value: "1"},
			"date":     &types.AttributeValueMemberS{Value: "15 Octubre 2026"},
			"time":     &types.AttributeValueMemberS{Value: "21:00 hs"},
			"status":   &types.AttributeValueMemberS{Value: "available"},
		},
		{
			"event_id": &types.AttributeValueMemberS{Value: eventId},
			"id":       &types.AttributeValueMemberS{Value: "2"},
			"date":     &types.AttributeValueMemberS{Value: "16 Octubre 2026"},
			"time":     &types.AttributeValueMemberS{Value: "20:00 hs"},
			"status":   &types.AttributeValueMemberS{Value: "soldout"},
		},
	}

	for _, item := range shows {
		_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item:      item,
		})
		if err != nil {
			log.Printf("Failed to insert show %s: %v", item["id"].(*types.AttributeValueMemberS).Value, err)
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
