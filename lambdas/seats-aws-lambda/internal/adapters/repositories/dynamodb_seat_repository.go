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

func (r *dynamoDBSeatRepository) ReserveSeat(eventId string, seatId string) error {
	_, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: eventId},
			"seat_id":  &types.AttributeValueMemberS{Value: seatId},
		},
		UpdateExpression:    aws.String("SET #s = :processing"),
		ConditionExpression: aws.String("#s = :available OR attribute_not_exists(#s)"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":processing": &types.AttributeValueMemberS{Value: "processing"},
			":available":  &types.AttributeValueMemberS{Value: "available"},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to reserve seat or it was not available: %w", err)
	}

	return nil
}

func (r *dynamoDBSeatRepository) ReleaseSeat(eventId string, seatId string) error {
	_, err := r.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: eventId},
			"seat_id":  &types.AttributeValueMemberS{Value: seatId},
		},
		UpdateExpression:    aws.String("SET #s = :available"),
		ConditionExpression: aws.String("#s = :processing"),
		ExpressionAttributeNames: map[string]string{
			"#s": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":processing": &types.AttributeValueMemberS{Value: "processing"},
			":available":  &types.AttributeValueMemberS{Value: "available"},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to release seat or it wasn't in processing state: %w", err)
	}

	return nil
}
