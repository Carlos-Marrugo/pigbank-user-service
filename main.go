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
)

// User define la estructura para DynamoDB
type User struct {
	UUID     string `json:"uuid" dynamodbav:"uuid"`
	Name     string `json:"name" dynamodbav:"name"`
	LastName string `json:"lastName" dynamodbav:"lastName"`
	Email    string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"` // Pendiente: Encriptar
	Document string `json:"document" dynamodbav:"document"`
}

func HandleRequest(ctx context.Context, user User) (string, error) {
	// 1. Cargar configuración (Conectará a LocalStack por las variables de entorno)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return "Error de configuración", err
	}

	dbClient := dynamodb.NewFromConfig(cfg)

	// Generar UUID si no viene
	user.UUID = uuid.New().String()

	// 2. Mapear objeto a formato DynamoDB
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return "Error mapeando datos", err
	}

	// 3. Guardar en la tabla
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("pigbank-users"),
		Item:      item,
	})

	if err != nil {
		log.Printf("Error DB: %v", err)
		return "Error guardando en base de datos", err
	}

	return fmt.Sprintf("Usuario %s registrado con ID %s", user.Name, user.UUID), nil
}

func main() {
	lambda.Start(HandleRequest)
}
