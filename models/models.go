package models

import (
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// can be positive or negative
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Transfer struct {
	ID            int64 `json:"id"`
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	// must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type Users struct {
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	IsEmailVerify  bool      `json:"is_email_verify"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Sessions struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiredAt    time.Time `json:"expired_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type VerifyEmail struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}
