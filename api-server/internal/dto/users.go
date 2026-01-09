package dto

import "time"

type User struct {
	Email         string
	PasswordHash  string
	PhoneNo       string
	WalletAddress string
	Subscribed    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
