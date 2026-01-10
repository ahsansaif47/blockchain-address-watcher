package dto

import "time"

type RegisterUserRequest struct {
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	PhoneNo       string     `json:"phone_no"`
	WalletAddress string     `json:"wallet_address"`
	Subscribed    bool       `json:"subscribed"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type RegisterUserResponse struct {
	ID string `json:"id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type UserResponse struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	PhoneNo       string    `json:"phone_no"`
	WalletAddress string    `json:"wallet_address"`
	Subscribed    bool      `json:"subscribed"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// DeleteUserRequest represents the request body for user deletion
type DeleteUserRequest struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
