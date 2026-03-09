package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID     string `json:"uuid" dynamodbav:"uuid"`
	Name     string `json:"name" dynamodbav:"name"`
	LastName string `json:"lastName" dynamodbav:"lastName"`
	Email    string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
	Document string `json:"document" dynamodbav:"document"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func HandleRequest(ctx context.Context, user User) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "Error config", err
	}

	dbClient := dynamodb.NewFromConfig(cfg)

	hashedPass, err := hashPassword(user.Password)
	if err != nil {
		return "Error encriptando", err
	}
	user.Password = hashedPass

	user.UUID = uuid.New().String()

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return "Error mapeo", err
	}

	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("pigbank-users"),
		Item:      item,
	})

	if err != nil {
		log.Printf("Error DB: %v", err)
		return "Error DB", err
	}

	return fmt.Sprintf("Usuario %s registrado y protegido", user.UUID), nil
}

func main() {
	lambda.Start(HandleRequest)
}
