package main

import (
	"log"

	"github.com/Carlos-Marrugo/pigbank-user-service/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	r := gin.Default()
	userHandler := &api.UserHandler{}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
	}

	log.Fatal(r.Run(":8080"))
}
