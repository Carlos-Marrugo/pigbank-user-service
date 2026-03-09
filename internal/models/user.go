package models

type User struct {
	UUID     string `json:"uuid" dynamodbav:"uuid"`
	Name     string `json:"name" dynamodbav:"name"`
	LastName string `json:"last_name" dynamodbav:"last_name"`
	Email    string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
	Document string `json:"document" dynamodbav:"document"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Document string `json:"document"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
