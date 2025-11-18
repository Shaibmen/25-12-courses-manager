package dto

import "time"

type SessionDTO struct {
	ID           string
	Username     string
	RefreshToken string
	IsRevoked    bool
	ExpiresAt    time.Time
}
