package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

var (
	userRepo  *repository.UserRepository
	sqsClient *sqs.Client
)

func SetRepository(r *repository.UserRepository, s *sqs.Client) {
	userRepo = r
	sqsClient = s
}

func RegisterHandler(ctx context.Context, req models.RegisterRequest) (string, error) {
	if userRepo == nil || sqsClient == nil {
		return "Initialization Error", fmt.Errorf("repository or SQS client not initialized")
	}

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

	err = userRepo.Save(ctx, user)
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
	if userRepo == nil {
		return "", fmt.Errorf("repository not initialized")
	}

	user, err := userRepo.FindByEmail(ctx, req.Email)
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

func UpdateUserProfile(ctx context.Context, userID string, req models.UpdateProfileRequest) error {
	if userRepo == nil {
		return fmt.Errorf("repository not initialized")
	}

	user, err := userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	return userRepo.Update(ctx, userID, user.Document, req.Address, req.Phone)
}
