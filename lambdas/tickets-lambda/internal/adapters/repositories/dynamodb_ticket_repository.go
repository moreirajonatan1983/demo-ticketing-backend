package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/demoticketing/tickets/internal/core/domain"
	"github.com/demoticketing/tickets/internal/core/ports"
)

type dynamoDBTicketRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBTicketRepository(client *dynamodb.Client, tableName string) ports.TicketRepository {
	return &dynamoDBTicketRepository{
		client:    client,
		tableName: tableName,
	}
}

type ticketDTO struct {
	UserID    string `dynamodbav:"user_id"`
	TicketID  string `dynamodbav:"ticket_id"`
	EventName string `dynamodbav:"event_name"`
	Date      string `dynamodbav:"date"`
	Location  string `dynamodbav:"location"`
	Sector    string `dynamodbav:"sector"`
	Status    string `dynamodbav:"status"`
}

func (r *dynamoDBTicketRepository) GetTicketsByUser(userId string) ([]domain.Ticket, error) {
	out, err := r.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("user_id = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userId},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query dynamodb for tickets: %w", err)
	}

	var dtos []ticketDTO
	err = attributevalue.UnmarshalListOfMaps(out.Items, &dtos)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dynamo items: %w", err)
	}

	var tickets []domain.Ticket
	for _, dto := range dtos {
		tickets = append(tickets, domain.Ticket{
			ID:        dto.TicketID,
			EventName: dto.EventName,
			Date:      dto.Date,
			Location:  dto.Location,
			Sector:    dto.Sector,
			Status:    dto.Status,
		})
	}

	return tickets, nil
}
