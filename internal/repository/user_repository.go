package repository

import (
	"context"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
