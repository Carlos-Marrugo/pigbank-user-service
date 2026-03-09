package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/models"
	"github.com/Carlos-Marrugo/pigbank-user-service/internal/service"
)

func main() {
	req := models.LoginRequest{
		Email:    "jane@doe.com",
		Password: "password123",
	}

	token, err := service.LoginHandler(context.Background(), req)
	if err != nil {
		log.Fatalf("Error en el login: %v", err)
	}

	fmt.Println("--- LOGIN EXITOSO ---")
	fmt.Printf("Tu JWT Token es:\n%s\n", token)
}
