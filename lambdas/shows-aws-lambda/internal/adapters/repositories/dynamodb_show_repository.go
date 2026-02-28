package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/demoticketing/shows/internal/core/domain"
	"github.com/demoticketing/shows/internal/core/ports"
)

type dynamoDBShowRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBShowRepository(client *dynamodb.Client, tableName string) ports.ShowRepository {
	return &dynamoDBShowRepository{
		client:    client,
		tableName: tableName,
	}
}

type showDTO struct {
	EventID string `dynamodbav:"event_id"`
	ShowID  string `dynamodbav:"id"`
	Date    string `dynamodbav:"date"`
	Time    string `dynamodbav:"time"`
	Status  string `dynamodbav:"status"`
}

func (r *dynamoDBShowRepository) GetShowsByEvent(eventId string) ([]domain.Show, error) {
	out, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("event_id = :eventId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":eventId": &types.AttributeValueMemberS{Value: eventId},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query shows from dynamodb: %w", err)
	}

	var dtos []showDTO
	err = attributevalue.UnmarshalListOfMaps(out.Items, &dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamo items: %w", err)
	}

	var shows []domain.Show
	for _, dto := range dtos {
		shows = append(shows, domain.Show{
			ID:      dto.ShowID,
			EventID: dto.EventID,
			Date:    dto.Date,
			Time:    dto.Time,
			Status:  dto.Status,
		})
	}

	return shows, nil
}
