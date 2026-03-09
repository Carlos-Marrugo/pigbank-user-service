package models

type User struct {
	UUID     string `json:"uuid" dynamodbav:"uuid"`
	Name     string `json:"name" dynamodbav:"name"`
	LastName string `json:"lastName" dynamodbav:"lastName"`
	Email    string `json:"email" dynamodbav:"email"`
	Password string `json:"password" dynamodbav:"password"`
	Document string `json:"document" dynamodbav:"document"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Document string `json:"document"`
}
