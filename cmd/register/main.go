package main

import (
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/service"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(service.RegisterHandler)
}
