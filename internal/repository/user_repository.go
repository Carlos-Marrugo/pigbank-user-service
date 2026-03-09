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

func (r *UserRepository) FindByID(ctx context.Context, uuid string) (*models.User, error) {
	input := &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("uuid = :u"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":u": &types.AttributeValueMemberS{Value: uuid},
		},
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("user not found with uuid: %s", uuid)
	}

	var user models.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, uuid string, document string, address string, phone string) error {
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"uuid":     &types.AttributeValueMemberS{Value: uuid},
			"document": &types.AttributeValueMemberS{Value: document},
		},
		UpdateExpression: aws.String("SET address = :a, phone = :p"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberS{Value: address},
			":p": &types.AttributeValueMemberS{Value: phone},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	}

	_, err := r.client.UpdateItem(ctx, input)
	return err
}
