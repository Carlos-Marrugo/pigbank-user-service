package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var s3Client *s3.Client
var avatarBucket string

func SetS3Client(client *s3.Client, bucket string) {
	s3Client = client
	avatarBucket = bucket
}

func UploadAvatar(ctx context.Context, userID string, base64Image string) (string, error) {
	if s3Client == nil {
		return "", fmt.Errorf("S3 client not initialized")
	}

	parts := strings.Split(base64Image, ",")
	var imageData string
	if len(parts) > 1 {
		imageData = parts[1]
	} else {
		imageData = base64Image
	}

	decoded, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	key := fmt.Sprintf("avatars/%s/%s.jpg", userID, uuid.New().String())

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(avatarBucket),
		Key:         aws.String(key),
		Body:        strings.NewReader(string(decoded)),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %v", err)
	}

	url := fmt.Sprintf("http://localhost:4566/%s/%s", avatarBucket, key)

	if userRepo != nil {
		user, err := userRepo.FindByID(ctx, userID)
		if err == nil {
			userRepo.UpdateAvatarURL(ctx, userID, user.Document, url)
		}
	}

	return url, nil
}
