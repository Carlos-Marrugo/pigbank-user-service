package models

type User struct {
	UUID      string `json:"uuid" dynamodbav:"uuid"`
	Name      string `json:"name" dynamodbav:"name"`
	LastName  string `json:"last_name" dynamodbav:"last_name"`
	Email     string `json:"email" dynamodbav:"email"`
	Password  string `json:"password" dynamodbav:"password"`
	Document  string `json:"document" dynamodbav:"document"`
	Address   string `json:"address" dynamodbav:"address"`
	Phone     string `json:"phone" dynamodbav:"phone"`
	AvatarURL string `json:"avatar_url" dynamodbav:"avatar_url"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Document string `json:"document"`
}

type AvatarRequest struct {
	AvatarBase64 string `json:"avatar_base64"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Address string `json:"address"`
	Phone   string `json:"phone"`
}
