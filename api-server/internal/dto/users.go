package dto

import "time"

type RegisterUserRequest struct {
	Email         string     `json:"email" validate:"required,email,min=5,max=255"`
	Password      string     `json:"password" validate:"required,strong_password,min=8,max=128"`
	PhoneNo       string     `json:"phone_no" validate:"required,phone,min=10,max=20"`
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

type DeleteUserRequest struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Details string            `json:"details,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}
