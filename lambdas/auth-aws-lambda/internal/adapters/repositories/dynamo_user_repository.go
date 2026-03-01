package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/demoticketing/auth/internal/core/domain"
)

// DynamoDBUserRepository implementa UserRepository usando DynamoDB.
type DynamoDBUserRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBUserRepository(client *dynamodb.Client, tableName string) *DynamoDBUserRepository {
	return &DynamoDBUserRepository{client: client, tableName: tableName}
}

func (r *DynamoDBUserRepository) Save(user *domain.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("error marshaling user: %w", err)
	}

	_, err = r.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *DynamoDBUserRepository) FindByEmail(email string) (*domain.User, error) {
	result, err := r.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"email": &types.AttributeValueMemberS{Value: email},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	var user domain.User
	if err := attributevalue.UnmarshalMap(result.Item, &user); err != nil {
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}
	return &user, nil
}
