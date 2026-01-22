package objects

import "time"

type Payload struct {
	Before any
	After  any
	Source any
}

type User struct {
	Id            string     `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"password_hash"`
	PhoneNo       string     `json:"phone_no"`
	WalletAddress string     `json:"wallet_address"`
	Subscribed    bool       `json:"subscribed"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
