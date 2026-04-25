package main

import (
	"context"
	"log"
	"os"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/api"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/repository"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/service"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	ctx := context.TODO()

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == "s3" {
			return aws.Endpoint{
				URL:               "http://localhost:4566",
				SigningRegion:     "us-east-1",
				HostnameImmutable: true,
			}, nil
		}
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
		log.Fatalf("unable to load SDK config, %v", err)
	}

	dbClient := dynamodb.NewFromConfig(cfg)
	sqsClient := sqs.NewFromConfig(cfg)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	userRepo := repository.NewUserRepository(dbClient)

	avatarBucket := os.Getenv("AVATAR_BUCKET_NAME")
	if avatarBucket == "" {
		avatarBucket = "pigbank-user-avatars-04f720f8"
	}

	service.SetRepository(userRepo, sqsClient)
	service.SetS3Client(s3Client, avatarBucket)

	r := gin.Default()
	userHandler := &api.UserHandler{}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)

		authorized := v1.Group("/")
		authorized.Use(api.AuthMiddleware())
		{
			authorized.PUT("/profile/:user_id", userHandler.UpdateProfile)
			authorized.POST("/profile/:user_id/avatar", userHandler.UploadAvatar)
		}
	}

	log.Fatal(r.Run(":8080"))
}
