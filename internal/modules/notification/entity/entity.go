package entity

import "time"

type NotificationEntity struct {
	ID               uint       `json:"id"`
	ReceiverID       uint       `json:"receiver_id"`
	ReceiverEmail    *string    `json:"receiver_email"`
	Subject          *string    `json:"subject"`
	Message          string     `json:"message"`
	Status           string     `json:"status"`
	NotificationType string     `json:"notification_type"`
	SentAt           *time.Time `json:"sent_at"`
	ReadAt           *time.Time `json:"read_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type NotifyQueryString struct {
	Search    string
	Status    string
	Page      int64
	Limit     int64
	UserID    uint
	OrderBy   string
	OrderType string
	IsRead    bool
}

type JwtUserData struct {
	UserID   uint   `json:"user_id"`
	RoleName string `json:"role_name"`
	Email    string `json:"email"`
}
