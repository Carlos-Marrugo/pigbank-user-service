package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/service" 
)

func maint() {
	req := models.RegisterRequest{
		Name:     "Jane",
		LastName: "Doe",
		Email:    "jane@doe.com",
		Password: "password123",
		Document: "1234567890",
	}

	res, err := service.RegisterHandler(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Resultado:", res)
}
