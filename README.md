# Pig Bank - User Service (Go)

## Responsibilities

* User registration with bcrypt encryption 
* JWT authentication (24h expiration)
* Profile management (address, phone)
* Avatar upload (Base64 → S3 → URL)
* Event publishing to SQS (card creation)

---

## Tech Stack

* Go 1.25
* Gin Framework
* DynamoDB
* S3
* SQS
* JWT (HS256)

---

## API Endpoints

### Public Routes

| Method | Endpoint           | Description        |
| ------ | ------------------ | ------------------ |
| POST   | `/api/v1/register` | Create new user    |
| POST   | `/api/v1/login`    | Authenticate → JWT |

---

### Protected Routes (Bearer Token)

| Method | Endpoint                           | Description           |
| ------ | ---------------------------------- | --------------------- |
| PUT    | `/api/v1/profile/{user_id}`        | Update address, phone |
| POST   | `/api/v1/profile/{user_id}/avatar` | Upload avatar to S3   |

---

## Event Publishing

On user registration → sends message to SQS:

```json
{
  "userId": "uuid",
  "request": "DEBIT"
}
```

Queue:

```
create-request-card-sqs
```

---

## Local Development

### Environment Variables (`.env`)

```env
AVATAR_BUCKET_NAME=pigbank-user-avatars-xxxx
CARD_QUEUE_URL=http://sqs.localhost:4566/000000000000/create-request-card-sqs
```

---

### Run Service

```bash
go mod tidy
go run main.go
```

Runs on:

```
http://localhost:8081
```

---

### Test Endpoint

```bash
curl -X POST http://localhost:8081/api/v1/register \
-d '{"name":"John","email":"john@test.com","password":"123"}'
```

---

## Database Schema (DynamoDB)

```json
{
  "uuid": "PK",
  "document": "SK",
  "email": "GSI: EmailIndex",
  "password": "bcrypt hash",
  "address": "string",
  "phone": "string",
  "avatar_url": "S3 path"
}
```

---

## Security

* Passwords → bcrypt (cost 14)
* JWT → HS256 (24h expiration)
* Protected routes → Bearer token required

---
