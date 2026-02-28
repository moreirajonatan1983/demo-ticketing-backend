package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/demoticketing/seats/internal/core/domain"
	"github.com/demoticketing/seats/internal/core/ports"
)

type dynamoDBSeatRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBSeatRepository(client *dynamodb.Client, tableName string) ports.SeatRepository {
	return &dynamoDBSeatRepository{
		client:    client,
		tableName: tableName,
	}
}

type seatDTO struct {
	EventID string `dynamodbav:"event_id"`
	SeatID  string `dynamodbav:"seat_id"`
	Row     string `dynamodbav:"row"`
	Number  int    `dynamodbav:"number"`
	Status  string `dynamodbav:"status"`
}

func (r *dynamoDBSeatRepository) GetSeatsByEvent(eventId string) ([]domain.Seat, error) {
	out, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("event_id = :eventId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":eventId": &types.AttributeValueMemberS{Value: eventId},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query seats from dynamodb: %w", err)
	}

	var dtos []seatDTO
	err = attributevalue.UnmarshalListOfMaps(out.Items, &dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamo items: %w", err)
	}

	var seats []domain.Seat
	for _, dto := range dtos {
		seats = append(seats, domain.Seat{
			ID:     dto.SeatID,
			Row:    dto.Row,
			Number: dto.Number,
			Status: dto.Status,
		})
	}

	return seats, nil
}
