package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

func RegisterHandler(ctx context.Context, req models.RegisterRequest) (string, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           "http://localhost:4566",
			SigningRegion: "us-east-1",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)

	if err != nil {
		return "Config Error", fmt.Errorf("unable to load SDK config: %v", err)
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)
	repo := repository.NewUserRepository(dbClient)

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return "Encryption Error", err
	}

	user := models.User{
		UUID:     uuid.New().String(),
		Name:     req.Name,
		LastName: req.LastName,
		Email:    req.Email,
		Password: hashedPassword,
		Document: req.Document,
	}

	err = repo.Save(ctx, user)
	if err != nil {
		return "Database Error", fmt.Errorf("failed to save user: %v", err)
	}

	queueURL := os.Getenv("CARD_QUEUE_URL")
	if queueURL == "" {
		queueURL = "http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/create-request-card-sqs"
	}

	cardRequest := map[string]string{
		"userId":  user.UUID,
		"request": "DEBIT",
	}
	body, _ := json.Marshal(cardRequest)

	_, err = sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(string(body)),
	})

	if err != nil {
		return "User saved, but SQS failed", fmt.Errorf("sqs error: %v", err)
	}

	return fmt.Sprintf("User %s registered. Card request sent to SQS.", user.UUID), nil
}

func LoginHandler(ctx context.Context, req models.LoginRequest) (string, error) {
    customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
        return aws.Endpoint{
            URL:           "http://localhost:4566",
            SigningRegion: "us-east-1",
        }, nil
    })

    cfg, _ := config.LoadDefaultConfig(ctx,
        config.WithRegion("us-east-1"),
        config.WithEndpointResolverWithOptions(customResolver),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
    )

    repo := repository.NewUserRepository(dynamodb.NewFromConfig(cfg))

    user, err := repo.FindByEmail(ctx, req.Email)
    if err != nil {
        return "", fmt.Errorf("usuario no encontrado o error en db")
    }

    if !CheckPasswordHash(req.Password, user.Password) {
        return "", fmt.Errorf("credenciales inválidas")
    }

    token, err := GenerateToken(user.Email, user.UUID)
    if err != nil {
        return "", fmt.Errorf("error al generar token")
    }

    return token, nil
}
