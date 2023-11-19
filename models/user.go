package models

import "time"

type User struct {
	ID                string    `json:"id"`
	UserName          string    `json:"user_name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	IsVerified        bool      `json:"is_verified"`
	VerificationToken string    `json:"verification_token"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type RegisterInput struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyInput struct {
	Email string `json:"email" binding:"required"`
	Token string `json:"token" binding:"required"`
}

type LoginResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	UserName  string `json:"user_name"`
	TokenType string `json:"token_type"`
	Token     string `json:"token"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	UserName  string    `json:"user_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
