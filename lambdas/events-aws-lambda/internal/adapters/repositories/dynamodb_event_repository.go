package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/demoticketing/events/internal/core/domain"
	"github.com/demoticketing/events/internal/core/ports"
)

type dynamoDBEventRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBEventRepository(client *dynamodb.Client, tableName string) ports.EventRepository {
	return &dynamoDBEventRepository{
		client:    client,
		tableName: tableName,
	}
}

type eventDTO struct {
	ID                string `dynamodbav:"id"`
	Title             string `dynamodbav:"title"`
	Date              string `dynamodbav:"date"`
	Venue             string `dynamodbav:"venue"`
	Image             string `dynamodbav:"image"`
	Status            string `dynamodbav:"status"`
	Description       string `dynamodbav:"description"`
	MaxTicketsPerUser int    `dynamodbav:"max_tickets_per_user"`
}

func (r *dynamoDBEventRepository) GetAll() ([]domain.Event, error) {
	out, err := r.client.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan dynamodb: %w", err)
	}

	var dtos []eventDTO
	err = attributevalue.UnmarshalListOfMaps(out.Items, &dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamo items: %w", err)
	}

	var events []domain.Event
	for _, dto := range dtos {
		events = append(events, domain.Event{
			ID:                dto.ID,
			Title:             dto.Title,
			Date:              dto.Date,
			Venue:             dto.Venue,
			Image:             dto.Image,
			Status:            dto.Status,
			Description:       dto.Description,
			MaxTicketsPerUser: dto.MaxTicketsPerUser,
		})
	}

	return events, nil
}

func (r *dynamoDBEventRepository) GetByID(id string) (*domain.Event, error) {
	out, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item from dynamodb: %w", err)
	}

	if out.Item == nil {
		return nil, nil // Not found
	}

	var dto eventDTO
	err = attributevalue.UnmarshalMap(out.Item, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamo item: %w", err)
	}

	return &domain.Event{
		ID:                dto.ID,
		Title:             dto.Title,
		Date:              dto.Date,
		Venue:             dto.Venue,
		Image:             dto.Image,
		Status:            dto.Status,
		Description:       dto.Description,
		MaxTicketsPerUser: dto.MaxTicketsPerUser,
	}, nil
}
