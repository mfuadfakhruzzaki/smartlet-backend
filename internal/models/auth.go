package models

import (
	"time"
)

// AuthToken represents JWT token response - sesuai API spec adalah string
type AuthToken string

// TokenResponse represents full token response with metadata
type TokenResponse struct {
	Token     string    `json:"token"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RegisterRequest represents registration request
type RegisterRequest struct {
	Email      string  `json:"email" validate:"required,email"`
	Name       string  `json:"name" validate:"required"`
	Password   string  `json:"password" validate:"required,min=6"`
	Location   *string `json:"location"`
	NoTelp     *string `json:"no_telp"`
	ImgProfile *string `json:"img_profile"`
	Status     *int    `json:"status"` // 0=pending, 1=approved/active, 2=rejected/suspended
	Role       *int    `json:"role"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	User  User      `json:"user"`
	Token AuthToken `json:"token"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token AuthToken `json:"token"`
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
}

// APIResponse represents generic API response
type APIResponse[T any] struct {
	Data    T      `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Details map[string]string `json:"details,omitempty"`
}