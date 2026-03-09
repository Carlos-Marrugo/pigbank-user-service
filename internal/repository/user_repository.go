package repository

import (
	"context"
	"fmt"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewUserRepository(db *dynamodb.Client) *UserRepository {
	return &UserRepository{
		client:    db,
		tableName: "pigbank-users",
	}
}

func (r *UserRepository) Save(ctx context.Context, user models.User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	return err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("EmailIndex"),
		KeyConditionExpression: aws.String("email = :e"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":e": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	return &user, err
}
